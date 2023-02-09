package awskinesisreceiver

import (
	"fmt"
	"go.opentelemetry.io/collector/component"
)

// AWSConfig contains AWS specific configuration such as awskinesis stream, region, etc.
type AWSConfig struct {
	ConsumerGroupName string `mapstructure:"consumer_group_name"`
	StreamName        string `mapstructure:"stream_name"`
	KinesisEndpoint   string `mapstructure:"kinesis_endpoint"`
	DynamoDBEndpoint  string `mapstructure:"dynamodb_endpoint"`
	Region            string `mapstructure:"region"`
	Role              string `mapstructure:"role"`
	MaxRecordSize     int    `mapstructure:"max_record_size"`
	Interval          int    `mapstructure:"interval"`
	PositionInStream  string `mapstructure:"position_in_stream"`
}

type Encoding struct {
	Name string `mapstructure:"name"`
}

// Config contains the main configuration options for the awskinesis receiver
type Config struct {
	Encoding    string    `mapstructure:"encoding"`
	Compression string    `mapstructure:"compression"`
	AWS         AWSConfig `mapstructure:"aws"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if len(cfg.AWS.Region) == 0 {
		return fmt.Errorf("region should be configured")
	}
	if len(cfg.AWS.StreamName) == 0 {
		return fmt.Errorf("stream name should be configured")
	}
	return nil
}
