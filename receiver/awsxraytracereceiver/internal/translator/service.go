package translator

import (
	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

type ServiceBuilder struct {
	OTLPSpanResourceBuilder
}

func (s *ServiceBuilder) WithService(service *awsxray.ServiceData) *ServiceBuilder {
	if service != nil {
		s.pResource.Attributes().PutStr(conventions.AttributeServiceVersion, *service.Version)
	}
	return s
}
