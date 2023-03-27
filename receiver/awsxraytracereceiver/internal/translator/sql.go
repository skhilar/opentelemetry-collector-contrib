package translator

import (
	awsxray "github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/xray"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

type SQLBuilder struct {
	OTLPSpanResourceBuilder
}

func (b *SQLBuilder) WithSQL(data *awsxray.SQLData) *SQLBuilder {
	if data != nil {
		if data.User != nil {
			b.pSpan.Attributes().PutStr(conventions.AttributeDBUser, *data.User)
		}
		if data.ConnectionString != nil {
			b.pSpan.Attributes().PutStr(conventions.AttributeDBConnectionString, *data.ConnectionString)
		}
		if data.DatabaseType != nil {
			b.pSpan.Attributes().PutStr(conventions.AttributeDBSystem, *data.DatabaseType)
		}
		if data.SanitizedQuery != nil {
			b.pSpan.Attributes().PutStr(conventions.AttributeDBStatement, *data.SanitizedQuery)
		}
		if data.User != nil {
			b.pSpan.SetName(*data.URL)
		}
	}
	return b
}
