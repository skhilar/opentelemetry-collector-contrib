package awsxraytracereceiver

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"time"
)

const (
	typeStr         = "awsxraytrace"
	defaultInterval = 2 * time.Minute
	stability       = component.StabilityLevelAlpha
	defaultRegion   = "us-east-1"
)

func createDefaultConfig() component.Config {
	return &Config{
		Interval: string(defaultInterval),
		Region:   defaultRegion,
	}
}

func createTracesReceiver(_ context.Context, params receiver.CreateSettings, baseCfg component.Config, consumer consumer.Traces) (receiver.Traces, error) {
	if consumer == nil {
		return nil, component.ErrNilNextConsumer
	}
	cfg := baseCfg.(*Config)
	traceRcvr, err := newXRayReceiver(*cfg, params, consumer)
	if err != nil {
		return nil, err
	}
	return traceRcvr, nil
}

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, stability),
	)
}
