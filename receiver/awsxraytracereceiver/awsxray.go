package awsxraytracereceiver

import (
	"context"
	"encoding/json"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/xray"
	"github.com/aws/aws-sdk-go-v2/service/xray/types"
	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxraytracereceiver/internal/translator"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/receiver"
	conventions "go.opentelemetry.io/collector/semconv/v1.8.0"
	"go.uber.org/zap"
	"time"
)

const (
	maximumTraceIdsPerQuery = 5
)

type XRayAdapter struct {
	config    Config
	logger    *zap.Logger
	client    *xray.Client
	startTime *time.Time
	endTime   *time.Time
}

func NewXRayAdapter(config Config, set receiver.CreateSettings) (*XRayAdapter, error) {
	var configOpts []func(*awsconfig.LoadOptions) error
	if config.Region != "" {
		configOpts = append(configOpts, func(lo *awsconfig.LoadOptions) error {
			lo.Region = config.Region
			return nil
		})
	}
	awsConf, err := awsconfig.LoadDefaultConfig(context.Background(), configOpts...)
	if err != nil {
		return nil, err
	}
	var xrayOpts []func(*xray.Options)
	if config.Role != "" {
		xrayOpts = append(xrayOpts, func(o *xray.Options) {
			o.Credentials = stscreds.NewAssumeRoleProvider(
				sts.NewFromConfig(awsConf),
				config.Role,
			)
		})
	}
	if config.XRayEndpoint != "" {
		xrayOpts = append(xrayOpts,
			xray.WithEndpointResolver(
				xray.EndpointResolverFromURL(config.XRayEndpoint),
			),
		)
	}
	client := xray.NewFromConfig(awsConf, xrayOpts...)
	return &XRayAdapter{client: client, config: config, logger: set.Logger}, nil
}
func (x *XRayAdapter) ConsumeTraces(ctx context.Context) ([]ptrace.Traces, error) {
	now := time.Now()
	if x.startTime == nil {
		startTime := now.Add(time.Minute * -20)
		x.startTime = &startTime
		x.endTime = &now
	} else {
		x.startTime = x.endTime
		x.endTime = &now
	}
	x.logger.Info("consuming trace from xray ", zap.Time("start-time", *x.startTime), zap.Time("end-time", *x.endTime))
	traceIds, err := x.GetTraceSummaries(ctx)
	if err != nil {
		x.logger.Error("error in getting xray trace id ", zap.Error(err))
		return nil, err
	}
	var traces []types.Trace
	x.logger.Info("number of traces got ", zap.Int("xray-traces", len(traceIds)))
	if len(traceIds) == 0 {
		return nil, nil
	}
	for offset := 0; offset < len(traceIds); offset += maximumTraceIdsPerQuery {
		var nextOffset int
		if offset+maximumTraceIdsPerQuery > len(traceIds) {
			nextOffset = len(traceIds)
		} else {
			nextOffset = offset + maximumTraceIdsPerQuery
		}
		t, err := x.BatchGetTraces(ctx, traceIds[offset:nextOffset])
		if err != nil {
			x.logger.Error("Error in getting traces ", zap.Error(err))
			return nil, err
		}
		traces = append(traces, t...)
	}

	var pTraces []ptrace.Traces
	for _, trace := range traces {
		pTrace, err := x.ConvertTrace(trace)
		if err != nil {
			x.logger.Error("error in converting trace ", zap.Error(err))
			return nil, err
		}
		pTraces = append(pTraces, pTrace)
	}
	return pTraces, nil
}

func (x *XRayAdapter) GetTraceSummariesWithToken(ctx context.Context, nextToken *string) ([]string, *string, error) {
	var traceIds []string
	x.logger.Info("GetTraceSummariesWithToken")
	output, err := x.client.GetTraceSummaries(ctx,
		&xray.GetTraceSummariesInput{StartTime: x.startTime, EndTime: x.endTime, NextToken: nextToken})
	if err != nil {
		x.logger.Error("error in getting trace summaries", zap.Error(err))
		return nil, nil, err
	}
	if *output.TracesProcessedCount > 0 {
		for _, trace := range output.TraceSummaries {
			traceIds = append(traceIds, *trace.Id)
		}
	}
	return traceIds, output.NextToken, nil
}

func (x *XRayAdapter) GetTraceSummaries(ctx context.Context) ([]string, error) {
	var traceIds []string
	var isNext = true
	var nextToken *string
	x.logger.Info("GetTraceSummaries")
	for isNext {
		traces, token, err := x.GetTraceSummariesWithToken(ctx, nextToken)
		if err != nil {
			x.logger.Error("Error in getting trace summaries", zap.Error(err))
			return traceIds, err
		}
		nextToken = token
		if nextToken == nil {
			isNext = false
		}
		traceIds = append(traceIds, traces...)
	}
	return traceIds, nil
}

func (x *XRayAdapter) BatchGetTracesWithToken(ctx context.Context, traceIds []string, nextToken *string) ([]types.Trace, *string, error) {
	x.logger.Info("BatchGetTracesWithToken")
	out, err := x.client.BatchGetTraces(ctx, &xray.BatchGetTracesInput{TraceIds: traceIds, NextToken: nextToken})
	if err != nil {
		x.logger.Error("error in getting traces", zap.Error(err))
		return nil, nil, err
	}
	return out.Traces, out.NextToken, nil
}

func (x *XRayAdapter) BatchGetTraces(ctx context.Context, traceIds []string) ([]types.Trace, error) {
	var xrayTraces []types.Trace
	var isNext = true
	var nextToken *string
	x.logger.Info("BatchGetTraces")
	for isNext {
		traces, token, err := x.BatchGetTracesWithToken(ctx, traceIds, nextToken)
		if err != nil {
			x.logger.Error("error in getting traces", zap.Error(err))
			return xrayTraces, err
		}
		nextToken = token
		if nextToken == nil {
			isNext = false
		}
		xrayTraces = append(xrayTraces, traces...)
	}
	return xrayTraces, nil
}

func (x *XRayAdapter) ConvertTrace(trace types.Trace) (ptrace.Traces, error) {
	pTrace := ptrace.NewTraces()
	var segments []awsxray.Segment
	x.logger.Info("ConvertTrace")
	for _, segment := range trace.Segments {
		var xraySegment = &awsxray.Segment{}
		err := json.Unmarshal([]byte(*segment.Document), xraySegment)
		if err != nil {
			return pTrace, err
		}
		segments = append(segments, *xraySegment)
	}
	for _, segment := range segments {
		rSpan := pTrace.ResourceSpans().AppendEmpty()
		resource := rSpan.Resource()
		resource.Attributes().PutStr(conventions.AttributeServiceName, *segment.Name)
		pSpan := rSpan.ScopeSpans().AppendEmpty().Spans().AppendEmpty()
		spnBuilder := translator.NewOTLPSpanResourceBuilder(&pSpan, &resource)
		spnBuilder.Segment().WithTraceID(*segment.TraceID).
			WithParentSegmentID(segment.ParentID).WithSegment(&segment).WithSpanKind(ptrace.SpanKindServer).Aws().WithAws(segment.AWS).Cause().
			WithCause(segment.Cause).Service().WithService(segment.Service).Sql().
			WithSQL(segment.SQL).Http().WithHttp(segment.HTTP).Build()
		if len(segment.Subsegments) > 0 {
			for _, subSegment := range segment.Subsegments {
				rSpan := pTrace.ResourceSpans().AppendEmpty()
				resource := rSpan.Resource()
				resource.Attributes().PutStr(conventions.AttributeServiceName, *segment.Name)
				pSpan := rSpan.ScopeSpans().AppendEmpty().Spans().AppendEmpty()
				spnBuilder := translator.NewOTLPSpanResourceBuilder(&pSpan, &resource)
				spnBuilder.Segment().WithSpanKind(ptrace.SpanKindClient).WithTraceID(*segment.TraceID).WithParentSegmentID(segment.ID).WithSegment(&subSegment).Aws().WithAws(subSegment.AWS).Cause().
					WithCause(subSegment.Cause).Service().WithService(subSegment.Service).Sql().
					WithSQL(subSegment.SQL).Http().WithHttp(subSegment.HTTP).Build()
			}
		}

	}
	return pTrace, nil
}
