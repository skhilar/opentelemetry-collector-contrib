package awskinesisreceiver

import (
	"context"
	cfg "github.com/vmware/vmware-go-kcl-v2/clientlibrary/config"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr                 = "awskinesis"
	stability               = component.StabilityLevelAlpha
	defaultEncoding         = "otlp"
	defaultPositionInStream = "TRIM_HORIZON"
	defaultConsumerGroup    = "collector"
	defaultCompression      = "none"
)

var positionMap = map[string]cfg.InitialPositionInStream{
	"LATEST":       cfg.LATEST,
	"TRIM_HORIZON": cfg.TRIM_HORIZON,
	"AT_TIMESTAMP": cfg.AT_TIMESTAMP,
}

type FactoryOption func(factory *kinesisReceiverFactory)

func WithTracesUnmarshalers(tracesUnmarshalers ...TracesUnmarshaler) FactoryOption {
	return func(factory *kinesisReceiverFactory) {
		for _, unmarshaler := range tracesUnmarshalers {
			factory.tracesUnmarshalers[unmarshaler.Encoding()] = unmarshaler
		}
	}
}

func NewFactory(options ...FactoryOption) receiver.Factory {
	f := &kinesisReceiverFactory{
		tracesUnmarshalers: defaultTracesUnmarshalers(),
	}
	for _, o := range options {
		o(f)
	}
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithTraces(f.createTracesReceiver, stability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		Encoding:    defaultEncoding,
		Compression: defaultCompression,
		AWS: AWSConfig{
			MaxRecordSize:     500,
			Interval:          5000,
			ConsumerGroupName: defaultConsumerGroup,
			PositionInStream:  defaultPositionInStream,
		},
	}
}

type kinesisReceiverFactory struct {
	tracesUnmarshalers map[string]TracesUnmarshaler
}

func (f *kinesisReceiverFactory) createTracesReceiver(_ context.Context, set receiver.CreateSettings,
	cfg component.Config, nextConsumer consumer.Traces) (receiver.Traces, error) {
	c := cfg.(*Config)
	r, err := newTracesReceiver(*c, set, f.tracesUnmarshalers, nextConsumer)
	if err != nil {
		return nil, err
	}
	return r, nil
}
