package translator

import (
	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

type HttpBuilder struct {
	OTLPSpanResourceBuilder
}

func (h *HttpBuilder) WithHttp(data *awsxray.HTTPData) *HttpBuilder {
	if data != nil {
		//h.pSpan.SetKind(ptrace.SpanKindServer) // To be addressed. We need to decide kind is server or client
		if data.Request != nil {
			if data.Request.Method != nil {
				h.pSpan.Attributes().PutStr(conventions.AttributeHTTPMethod, *data.Request.Method)
			}
			if data.Request.ClientIP != nil {
				h.pSpan.Attributes().PutStr(conventions.AttributeHTTPClientIP, *data.Request.ClientIP)
			}
			if data.Request.UserAgent != nil {
				h.pSpan.Attributes().PutStr(conventions.AttributeHTTPUserAgent, *data.Request.UserAgent)
			}
			if data.Request.URL != nil {
				h.pSpan.Attributes().PutStr(conventions.AttributeHTTPURL, *data.Request.URL)
			}
		}
		if data.Response != nil {
			if data.Response.Status != nil {
				h.pSpan.Attributes().PutInt(conventions.AttributeHTTPStatusCode, *data.Response.Status)
			}
			if data.Response.ContentLength != nil {
				value := data.Response.ContentLength.(float64)
				h.pSpan.Attributes().PutDouble(conventions.AttributeMessagingMessagePayloadSizeBytes, value)
			}
		}

	}
	return h
}
