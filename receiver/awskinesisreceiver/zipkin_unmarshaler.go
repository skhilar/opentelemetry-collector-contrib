package awskinesisreceiver

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin/zipkinv1"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin/zipkinv2"
)

const (
	zipkinProtobufEncoding = "zipkin_proto"
	zipkinJSONEncoding     = "zipkin_json"
	zipkinThriftEncoding   = "zipkin_thrift"
)

func newZipkinProtobufUnmarshaler() TracesUnmarshaler {
	return newPdataTracesUnmarshaler(zipkinv2.NewProtobufTracesUnmarshaler(false, false), zipkinProtobufEncoding)
}

func newZipkinJSONUnmarshaler() TracesUnmarshaler {
	return newPdataTracesUnmarshaler(zipkinv2.NewJSONTracesUnmarshaler(false), zipkinJSONEncoding)
}

func newZipkinThriftUnmarshaler() TracesUnmarshaler {
	return newPdataTracesUnmarshaler(zipkinv1.NewThriftTracesUnmarshaler(), zipkinThriftEncoding)
}
