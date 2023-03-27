package translator

import (
	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	"time"
)

type CauseBuilder struct {
	OTLPSpanResourceBuilder
}

func (c *CauseBuilder) WithCause(cause *awsxray.CauseData) *CauseBuilder {
	if cause != nil {
		for _, e := range cause.Exceptions {
			event := c.pSpan.Events().AppendEmpty()
			event.SetTimestamp(pcommon.NewTimestampFromTime(time.Now()))
			event.SetName(*e.ID)
			event.Attributes().PutStr(conventions.AttributeExceptionType, *e.Type)
			event.Attributes().PutStr(conventions.AttributeExceptionMessage, *e.Message)
		}
	}
	return c
}
