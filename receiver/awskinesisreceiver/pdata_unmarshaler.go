package awskinesisreceiver

import (
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type pdataTracesUnmarshaler struct {
	ptrace.Unmarshaler
	encoding string
}

func (p pdataTracesUnmarshaler) Unmarshal(buf []byte) (ptrace.Traces, error) {
	return p.Unmarshaler.UnmarshalTraces(buf)
}

func (p pdataTracesUnmarshaler) Encoding() string {
	return p.encoding
}

func newPdataTracesUnmarshaler(unmarshaler ptrace.Unmarshaler, encoding string) TracesUnmarshaler {
	return pdataTracesUnmarshaler{
		Unmarshaler: unmarshaler,
		encoding:    encoding,
	}
}
