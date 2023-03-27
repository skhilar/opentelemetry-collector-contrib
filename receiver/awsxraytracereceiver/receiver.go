package awsxraytracereceiver

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
	"time"
)

type xrayReceiver struct {
	host         component.Host
	cancel       context.CancelFunc
	logger       *zap.Logger
	nextConsumer consumer.Traces
	config       Config
	xray         *XRayAdapter
}

func newXRayReceiver(config Config, set receiver.CreateSettings, nextConsumer consumer.Traces) (*xrayReceiver, error) {
	xrayAdapter, err := NewXRayAdapter(config, set)
	if err != nil {
		return nil, err
	}
	return &xrayReceiver{logger: set.Logger, config: config, xray: xrayAdapter, nextConsumer: nextConsumer}, nil
}

func (r *xrayReceiver) Start(ctx context.Context, host component.Host) error {
	r.host = host
	ctx, r.cancel = context.WithCancel(ctx)

	interval, _ := time.ParseDuration(r.config.Interval)
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.logger.Info("processing the trace!")
				tds, err := r.xray.ConsumeTraces(ctx)
				if err != nil {
					r.logger.Error("error in getting traces from xray", zap.Error(err))
				}
				for _, td := range tds {
					r.nextConsumer.ConsumeTraces(ctx, td)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (r *xrayReceiver) Shutdown(ctx context.Context) error {
	r.cancel()
	return nil
}
