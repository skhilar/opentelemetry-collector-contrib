package translator

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type OTLPSpan struct {
	*ptrace.Span
}

type OTLPResource struct {
	*pcommon.Resource
}

type OTLPSpanResource struct {
	span     *ptrace.Span
	resource *pcommon.Resource
}

type OTLPSpanResourceBuilder struct {
	pSpan     *OTLPSpan
	pResource *OTLPResource
}

func NewOTLPSpanResourceBuilder(span *ptrace.Span, resource *pcommon.Resource) *OTLPSpanResourceBuilder {
	return &OTLPSpanResourceBuilder{pSpan: &OTLPSpan{span}, pResource: &OTLPResource{resource}}
}

func (b *OTLPSpanResourceBuilder) Segment() *SegmentBuilder {
	s := &SegmentBuilder{parentSpanID: pcommon.NewSpanIDEmpty()}
	s.OTLPSpanResourceBuilder = *b
	return s
}

func (b *OTLPSpanResourceBuilder) Http() *HttpBuilder {
	return &HttpBuilder{*b}
}

func (b *OTLPSpanResourceBuilder) Sql() *SQLBuilder {
	return &SQLBuilder{*b}
}

func (b *OTLPSpanResourceBuilder) Service() *ServiceBuilder {
	return &ServiceBuilder{*b}
}

func (b *OTLPSpanResourceBuilder) Cause() *CauseBuilder {
	return &CauseBuilder{*b}
}

func (b *OTLPSpanResourceBuilder) Aws() *AWSBuilder {
	return &AWSBuilder{*b}
}

func (b *OTLPSpanResourceBuilder) Build() *OTLPSpanResource {
	return &OTLPSpanResource{span: b.pSpan.Span, resource: b.pResource.Resource}
}
