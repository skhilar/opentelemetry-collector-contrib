package awskinesisreceiver

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	guuid "github.com/google/uuid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awskinesisreceiver/internal/decompressor"
	chk "github.com/vmware/vmware-go-kcl-v2/clientlibrary/checkpoint"
	cfg "github.com/vmware/vmware-go-kcl-v2/clientlibrary/config"
	wk "github.com/vmware/vmware-go-kcl-v2/clientlibrary/worker"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

type kinesisConsumer struct {
	logger       *zap.Logger
	config       Config
	settings     receiver.CreateSettings
	unmarshalers map[string]TracesUnmarshaler
	worker       *wk.Worker
}

var _ receiver.Traces = (*kinesisConsumer)(nil)

type kinesisClientOptions struct {
	NewKinesisClient func(conf aws.Config, opts ...func(*kinesis.Options)) *kinesis.Client
}

type dynamoDBClientOptions struct {
	NewDynamoDBClient func(conf aws.Config, opts ...func(*dynamodb.Options)) *dynamodb.Client
}

func newTracesReceiver(config Config, set receiver.CreateSettings, unmarshalers map[string]TracesUnmarshaler,
	nextConsumer consumer.Traces) (*kinesisConsumer, error) {
	unmarshaler := unmarshalers[config.Encoding]
	if unmarshaler == nil {
		return nil, fmt.Errorf("unrecognized encoding")
	}
	id, err := guuid.NewUUID()
	if err != nil {
		return nil, err
	}
	logger := log{set.Logger}
	kinesisOptions := &kinesisClientOptions{
		NewKinesisClient: kinesis.NewFromConfig,
	}

	dynamoOptions := &dynamoDBClientOptions{NewDynamoDBClient: dynamodb.NewFromConfig}
	var configOpts []func(*awsconfig.LoadOptions) error
	if config.AWS.Region != "" {
		configOpts = append(configOpts, func(lo *awsconfig.LoadOptions) error {
			lo.Region = config.AWS.Region
			return nil
		})
	}

	awsConf, err := awsconfig.LoadDefaultConfig(context.Background(), configOpts...)
	if err != nil {
		return nil, err
	}

	var kinesisOpts []func(*kinesis.Options)
	var dynamoOpts []func(*dynamodb.Options)
	if config.AWS.Role != "" {
		kinesisOpts = append(kinesisOpts, func(o *kinesis.Options) {
			o.Credentials = stscreds.NewAssumeRoleProvider(
				sts.NewFromConfig(awsConf),
				config.AWS.Role,
			)
		})
		dynamoOpts = append(dynamoOpts, func(o *dynamodb.Options) {
			o.Credentials = stscreds.NewAssumeRoleProvider(
				sts.NewFromConfig(awsConf),
				config.AWS.Role,
			)
		})
	}

	if config.AWS.KinesisEndpoint != "" {
		kinesisOpts = append(kinesisOpts,
			kinesis.WithEndpointResolver(
				kinesis.EndpointResolverFromURL(config.AWS.KinesisEndpoint),
			),
		)
	}

	if config.AWS.DynamoDBEndpoint != "" {
		dynamoOpts = append(dynamoOpts,
			dynamodb.WithEndpointResolver(
				dynamodb.EndpointResolverFromURL(config.AWS.KinesisEndpoint),
			),
		)
	}
	kclConfig := cfg.NewKinesisClientLibConfig(config.AWS.ConsumerGroupName, config.AWS.StreamName, config.AWS.Region,
		id.String()).
		WithInitialPositionInStream(positionMap[config.AWS.PositionInStream]).
		WithMaxRecords(config.AWS.MaxRecordSize).
		WithMaxLeasesForWorker(cfg.DefaultMaxLeasesForWorker).
		WithShardSyncIntervalMillis(config.AWS.Interval).
		WithFailoverTimeMillis(cfg.DefaultFailoverTimeMillis).
		WithKinesisEndpoint(config.AWS.KinesisEndpoint).
		WithLogger(logger)

	decompressor := decompressor.NewDecompressor(config.Compression)
	worker := wk.NewWorker(newProcessorFactory(nextConsumer, unmarshaler, decompressor, set.Logger), kclConfig)
	worker.WithKinesis(kinesisOptions.NewKinesisClient(awsConf, kinesisOpts...))
	chkPointer := chk.NewDynamoCheckpoint(kclConfig).WithDynamoDB(dynamoOptions.NewDynamoDBClient(awsConf, dynamoOpts...))
	worker.WithCheckpointer(chkPointer)
	return &kinesisConsumer{logger: set.Logger,
		config:       config,
		unmarshalers: unmarshalers,
		worker:       worker,
		settings:     set,
	}, nil
}

func (c *kinesisConsumer) Start(_ context.Context, _ component.Host) error {
	c.logger.Debug("starting to process traces")
	return c.worker.Start()
}

func (c *kinesisConsumer) Shutdown(context.Context) error {
	c.logger.Debug("shutting down trace receiver")
	c.worker.Shutdown()
	return nil
}
