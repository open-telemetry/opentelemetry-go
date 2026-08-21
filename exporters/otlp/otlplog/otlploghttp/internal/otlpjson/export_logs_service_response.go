// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpjson

import (
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

// MarshalExportLogsServiceResponse encodes an ExportLogsServiceResponse as JSON Protobuf encoded bytes.
func MarshalExportLogsServiceResponse(resp *collogpb.ExportLogsServiceResponse) ([]byte, error) {
	return protojson.Marshal(resp)
}

// UnmarshalExportLogsServiceResponse decodes JSON Protobuf encoded payload into an ExportLogsServiceResponse.
func UnmarshalExportLogsServiceResponse(data []byte, resp *collogpb.ExportLogsServiceResponse) error {
	// ignore message fields with unknown names per OTLP specs.
	var unmarshaler protojson.UnmarshalOptions
	unmarshaler.DiscardUnknown = true
	return unmarshaler.Unmarshal(data, resp)
}
