// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpjson

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

func TestMarshalExportTraceServiceResponse_Nil(t *testing.T) {
	data, err := MarshalExportTraceServiceResponse(nil)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(data))
}

func TestMarshalExportTraceServiceResponse_Empty(t *testing.T) {
	data, err := MarshalExportTraceServiceResponse(&coltracepb.ExportTraceServiceResponse{})
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(data))
}

func TestMarshalExportTraceServiceResponse_PartialSuccess(t *testing.T) {
	resp := &coltracepb.ExportTraceServiceResponse{
		PartialSuccess: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: 5,
			ErrorMessage:  "resource exhausted",
		},
	}

	data, err := MarshalExportTraceServiceResponse(resp)
	require.NoError(t, err)

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var m map[string]any
	require.NoError(t, dec.Decode(&m))

	ps, ok := m["partialSuccess"].(map[string]any)
	require.True(t, ok, "expected partialSuccess field")

	rejected, ok := ps["rejectedSpans"].(string)
	require.True(t, ok, "rejectedSpans must be a quoted string, got %T", ps["rejectedSpans"])
	assert.Equal(t, "5", rejected)

	errMsg, ok := ps["errorMessage"].(string)
	require.True(t, ok)
	assert.Equal(t, "resource exhausted", errMsg)

	// Field names must be camelCase.
	_, hasSnake := ps["rejected_spans"]
	assert.False(t, hasSnake, "must not use snake_case")
	_, hasSnake = ps["error_message"]
	assert.False(t, hasSnake, "must not use snake_case")
}

func TestMarshalExportTraceServiceResponse_ZeroRejected(t *testing.T) {
	resp := &coltracepb.ExportTraceServiceResponse{
		PartialSuccess: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: 0,
			ErrorMessage:  "partial",
		},
	}

	data, err := MarshalExportTraceServiceResponse(resp)
	require.NoError(t, err)

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var m map[string]any
	require.NoError(t, dec.Decode(&m))

	ps := m["partialSuccess"].(map[string]any)
	_, hasRejected := ps["rejectedSpans"]
	assert.False(t, hasRejected, "zero rejectedSpans should be omitted")
}

func TestUnmarshalExportTraceServiceResponse_PartialSuccess(t *testing.T) {
	input := `{"partialSuccess":{"rejectedSpans":"5","errorMessage":"resource exhausted"}}`

	resp := &coltracepb.ExportTraceServiceResponse{}
	err := UnmarshalExportTraceServiceResponse([]byte(input), resp)
	require.NoError(t, err)
	require.NotNil(t, resp.PartialSuccess)
	assert.Equal(t, int64(5), resp.PartialSuccess.RejectedSpans)
	assert.Equal(t, "resource exhausted", resp.PartialSuccess.ErrorMessage)
}

func TestUnmarshalExportTraceServiceResponse_Empty(t *testing.T) {
	resp := &coltracepb.ExportTraceServiceResponse{}
	err := UnmarshalExportTraceServiceResponse([]byte(`{}`), resp)
	require.NoError(t, err)
	assert.Nil(t, resp.PartialSuccess)
}

func TestUnmarshalExportTraceServiceResponse_IgnoresUnknownFields(t *testing.T) {
	input := `{"partialSuccess":{"rejectedSpans":"3","errorMessage":"err","futureField":true},"unknownTop":42}`

	resp := &coltracepb.ExportTraceServiceResponse{}
	err := UnmarshalExportTraceServiceResponse([]byte(input), resp)
	require.NoError(t, err)
	require.NotNil(t, resp.PartialSuccess)
	assert.Equal(t, int64(3), resp.PartialSuccess.RejectedSpans)
	assert.Equal(t, "err", resp.PartialSuccess.ErrorMessage)
}

func TestExportTraceServiceResponseRoundTrip(t *testing.T) {
	original := &coltracepb.ExportTraceServiceResponse{
		PartialSuccess: &coltracepb.ExportTracePartialSuccess{
			RejectedSpans: 42,
			ErrorMessage:  "quota exceeded",
		},
	}

	data, err := MarshalExportTraceServiceResponse(original)
	require.NoError(t, err)

	decoded := &coltracepb.ExportTraceServiceResponse{}
	err = UnmarshalExportTraceServiceResponse(data, decoded)
	require.NoError(t, err)

	require.NotNil(t, decoded.PartialSuccess)
	assert.Equal(t, original.PartialSuccess.RejectedSpans, decoded.PartialSuccess.RejectedSpans)
	assert.Equal(t, original.PartialSuccess.ErrorMessage, decoded.PartialSuccess.ErrorMessage)
}
