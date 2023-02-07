package awskinesisreceiver

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
)

const (
	typeStr         = "kinesisrcv"
	stability       = component.StabilityLevelAlpha
	defaultEncoding = "otlp"
)

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
		Encoding: defaultEncoding,
		AWS:      AWSConfig{MaxRecordSize: 500, Interval: 5000},
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
