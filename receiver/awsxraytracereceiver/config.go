package awsxraytracereceiver

import (
	"fmt"
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	Region       string `mapstructure:"region"`
	Role         string `mapstructure:"role"`
	XRayEndpoint string `mapstructure:"xray_endpoint"`
	Interval     string `mapstructure:"interval"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if len(cfg.Region) == 0 {
		return fmt.Errorf("region should be configured")
	}
	return nil
}
