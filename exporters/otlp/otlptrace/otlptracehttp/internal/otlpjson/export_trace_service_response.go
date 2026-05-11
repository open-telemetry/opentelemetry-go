// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package otlpjson implements OTLP JSON Protobuf encoding for trace data.
//
// The encoding conforms to the OTLP specs
// (https://opentelemetry.io/docs/specs/otlp/#json-protobuf-encoding):
//   - trace ID and span ID byte arrays are encoded as case-insensitive hex-encoded strings
//   - enum values encoded as integers
//   - field names in lowerCamelCase
//   - 64-bit integers encoded as quoted decimal strings (ProtoJSON specs)
package otlpjson // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/otlpjson"

import (
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// MarshalExportTraceServiceResponse encodes an ExportTraceServiceResponse as JSON Protobuf encoded bytes.
func MarshalExportTraceServiceResponse(resp *coltracepb.ExportTraceServiceResponse) ([]byte, error) {
	return protojson.Marshal(resp)
}

// UnmarshalExportTraceServiceResponse decodes JSON Protobuf encoded payload into an ExportTraceServiceResponse.
func UnmarshalExportTraceServiceResponse(data []byte, resp *coltracepb.ExportTraceServiceResponse) error {
	// ignore message fields with unknown names per OTLP specs.
	var unmarshaler protojson.UnmarshalOptions
	unmarshaler.DiscardUnknown = true
	return unmarshaler.Unmarshal(data, resp)
}
