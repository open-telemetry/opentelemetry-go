// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricExportBatchSize(t *testing.T) {
	const key = "OTEL_GO_X_METRIC_EXPORT_BATCH_SIZE"
	require.Contains(t, MetricExportBatchSize.Keys(), key)

	tests := []struct {
		name    string
		value   string
		enabled bool
		want    int
	}{
		{name: "empty", value: "", enabled: false, want: 0},
		{name: "invalid", value: "invalid", enabled: false, want: 0},
		{name: "zero", value: "0", enabled: false, want: 0},
		{name: "negative", value: "-10", enabled: false, want: 0},
		{name: "valid small", value: "10", enabled: true, want: 10},
		{name: "valid large", value: "200", enabled: true, want: 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(key, tt.value)
			assert.Equal(t, tt.enabled, MetricExportBatchSize.Enabled())
			got, ok := MetricExportBatchSize.Lookup()
			assert.Equal(t, tt.enabled, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}
