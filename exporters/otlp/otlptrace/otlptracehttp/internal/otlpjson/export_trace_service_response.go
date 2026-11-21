// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpjson

import (
	"encoding/json"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal"
)

// exportResponseJSON corresponds to the JSON Protobuf encoding of the OTLP
// ExportTraceServiceResponse message.
type exportResponseJSON struct {
	PartialSuccess *exportPartialSuccessJSON `json:"partialSuccess,omitempty"`
}

// exportPartialSuccessJSON corresponds to the JSON Protobuf encoding of the
// OTLP ExportTracePartialSuccess message. The rejected span count is an
// int64, which ProtoJSON encodes as a quoted decimal string.
type exportPartialSuccessJSON struct {
	ErrorMessage  string `json:"errorMessage,omitempty"`
	RejectedSpans int64  `json:"rejectedSpans,omitempty,string"`
}

// MarshalExportTraceServiceResponse encodes an OTLP ExportTraceServiceResponse
// as JSON Protobuf encoded bytes.
func MarshalExportTraceServiceResponse(resp *internal.ExportResponse) ([]byte, error) {
	out := exportResponseJSON{}
	if ps := resp.GetPartialSuccess(); ps != nil {
		out.PartialSuccess = &exportPartialSuccessJSON{
			ErrorMessage:  ps.ErrorMessage,
			RejectedSpans: ps.RejectedItems,
		}
	}
	return json.Marshal(out)
}

// UnmarshalExportTraceServiceResponse decodes a JSON Protobuf encoded
// ExportTraceServiceResponse payload. Unknown message fields are ignored per
// the OTLP specs.
func UnmarshalExportTraceServiceResponse(data []byte) (*internal.ExportResponse, error) {
	var jr exportResponseJSON
	if err := json.Unmarshal(data, &jr); err != nil {
		return nil, err
	}
	resp := &internal.ExportResponse{}
	if jr.PartialSuccess != nil {
		resp.PartialSuccess = &internal.ExportPartialSuccess{
			ErrorMessage:  jr.PartialSuccess.ErrorMessage,
			RejectedItems: jr.PartialSuccess.RejectedSpans,
		}
	}
	return resp, nil
}
