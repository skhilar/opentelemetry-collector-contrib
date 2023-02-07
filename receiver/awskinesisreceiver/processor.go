package awskinesisreceiver

import (
	"context"
	kc "github.com/vmware/vmware-go-kcl-v2/clientlibrary/interfaces"
	"go.opentelemetry.io/collector/consumer"
	"go.uber.org/zap"
)

func newProcessorFactory(nextConsumer consumer.Traces, unmarshaler TracesUnmarshaler, logger *zap.Logger) kc.IRecordProcessorFactory {
	kp := kinesisProcessor{nextConsumer: nextConsumer, unmarshaler: unmarshaler, logger: logger}
	return &kinesisProcessorFactory{kinesisProcessor: &kp}
}

type kinesisProcessorFactory struct {
	kinesisProcessor *kinesisProcessor
}

func (kpf *kinesisProcessorFactory) CreateProcessor() kc.IRecordProcessor {
	return kpf.kinesisProcessor
}

type kinesisProcessor struct {
	nextConsumer consumer.Traces
	unmarshaler  TracesUnmarshaler
	shardId      string
	logger       *zap.Logger
}

func (kp *kinesisProcessor) Initialize(input *kc.InitializationInput) {
	kp.shardId = input.ShardId
	kp.logger.Info("initialized processor with shardId ", zap.String("shardId", kp.shardId))
}

func (kp *kinesisProcessor) ProcessRecords(input *kc.ProcessRecordsInput) {
	records := input.Records
	if len(records) == 0 {
		return
	}
	for _, record := range records {
		traces, err := kp.unmarshaler.Unmarshal(record.Data)
		if err != nil {
			kp.logger.Error("not able to unmarshal traces ", zap.Error(err))
			continue
		}
		err = kp.nextConsumer.ConsumeTraces(context.Background(), traces)
		if err != nil {
			kp.logger.Error("not able to send the trace to next consumer ", zap.Error(err))
		}
	}
	lastSequenceNumber := input.Records[len(input.Records)-1].SequenceNumber
	err := input.Checkpointer.Checkpoint(lastSequenceNumber)
	if err != nil {
		kp.logger.Error("not able to checkpoint sequence number ", zap.Error(err))
	}
}

func (kp *kinesisProcessor) Shutdown(input *kc.ShutdownInput) {
	if input.ShutdownReason == kc.TERMINATE {
		_ = input.Checkpointer.Checkpoint(nil)
	}
}
