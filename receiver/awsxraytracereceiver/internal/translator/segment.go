package translator

import (
	"encoding/hex"
	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"math"
	"strings"
	"time"
)

const (
	xrayTraceIDDelim = "-"
)

type SegmentBuilder struct {
	OTLPSpanResourceBuilder
	parentSpanID pcommon.SpanID
	traceID      pcommon.TraceID
}

func (s *SegmentBuilder) WithParentSegmentID(ID *string) *SegmentBuilder {
	s.parentSpanID = pcommon.NewSpanIDEmpty()
	if ID != nil {
		hex.Decode(s.parentSpanID[:], []byte(*ID))
	}
	return s
}
func (s *SegmentBuilder) WithTraceID(ID string) *SegmentBuilder {
	s.traceID = makeTraceIDFromAWSXrayTraceID(ID)
	return s
}
func (s *SegmentBuilder) WithSpanKind(spanKind ptrace.SpanKind) *SegmentBuilder {
	s.pSpan.SetKind(spanKind)
	return s
}
func (s *SegmentBuilder) WithSegment(segment *awsxray.Segment) *SegmentBuilder {
	s.pSpan.SetStartTimestamp(makeTimestamp(*segment.StartTime))
	s.pSpan.SetEndTimestamp(makeTimestamp(*segment.EndTime))
	spanID := pcommon.NewSpanIDEmpty()
	hex.Decode(spanID[:], []byte(*segment.ID))
	s.pSpan.SetSpanID(spanID)
	s.pSpan.SetTraceID(s.traceID)
	if !s.parentSpanID.IsEmpty() {
		s.pSpan.SetParentSpanID(s.parentSpanID)
	}
	s.pSpan.SetName(*segment.Name)
	return s
}

//Need to validate if the implementation is correct
func makeTraceIDFromAWSXrayTraceID(traceID string) pcommon.TraceID {
	parts := strings.Split(traceID, xrayTraceIDDelim)
	convertedTraceID := pcommon.NewTraceIDEmpty()
	hex.Decode(convertedTraceID[0:4], []byte(parts[1]))
	hex.Decode(convertedTraceID[4:], []byte(parts[2]))
	return convertedTraceID
}

func makeTimestamp(ts float64) pcommon.Timestamp {
	sec, nSec := math.Modf(ts)
	t := time.Unix(int64(sec), int64(nSec*(1e9)))
	return pcommon.NewTimestampFromTime(t)
}
