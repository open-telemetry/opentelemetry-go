// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlpjson

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	collogpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

func TestUnmarshalExportLogsServiceResponse_IgnoresUnknownFields(t *testing.T) {
	input := `{"partialSuccess":{"rejectedLogRecords":"3","errorMessage":"err","futureField":true},"unknownTop":42}`

	resp := &collogpb.ExportLogsServiceResponse{}
	err := UnmarshalExportLogsServiceResponse([]byte(input), resp)
	require.NoError(t, err)
	require.NotNil(t, resp.PartialSuccess)
	assert.Equal(t, int64(3), resp.PartialSuccess.RejectedLogRecords)
	assert.Equal(t, "err", resp.PartialSuccess.ErrorMessage)
}

func TestExportLogsServiceResponseRoundTrip(t *testing.T) {
	original := &collogpb.ExportLogsServiceResponse{
		PartialSuccess: &collogpb.ExportLogsPartialSuccess{
			RejectedLogRecords: 42,
			ErrorMessage:       "quota exceeded",
		},
	}

	data, err := MarshalExportLogsServiceResponse(original)
	require.NoError(t, err)

	decoded := &collogpb.ExportLogsServiceResponse{}
	err = UnmarshalExportLogsServiceResponse(data, decoded)
	require.NoError(t, err)

	require.NotNil(t, decoded.PartialSuccess)
	assert.Equal(t, original.PartialSuccess.RejectedLogRecords, decoded.PartialSuccess.RejectedLogRecords)
	assert.Equal(t, original.PartialSuccess.ErrorMessage, decoded.PartialSuccess.ErrorMessage)
}
