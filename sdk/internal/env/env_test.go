// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ottest "go.opentelemetry.io/otel/sdk/internal/internaltest"
)

func TestEnvParse(t *testing.T) {
	testCases := []struct {
		name string
		keys []string
		f    func(int) int
	}{
		{
			name: "BatchSpanProcessorScheduleDelay",
			keys: []string{BatchSpanProcessorScheduleDelayKey},
			f:    BatchSpanProcessorScheduleDelay,
		},

		{
			name: "BatchSpanProcessorExportTimeout",
			keys: []string{BatchSpanProcessorExportTimeoutKey},
			f:    BatchSpanProcessorExportTimeout,
		},

		{
			name: "BatchSpanProcessorMaxQueueSize",
			keys: []string{BatchSpanProcessorMaxQueueSizeKey},
			f:    BatchSpanProcessorMaxQueueSize,
		},

		{
			name: "BatchSpanProcessorMaxExportBatchSize",
			keys: []string{BatchSpanProcessorMaxExportBatchSizeKey},
			f:    BatchSpanProcessorMaxExportBatchSize,
		},

		{
			name: "SpanAttributeValueLength",
			keys: []string{SpanAttributeValueLengthKey, AttributeValueLengthKey},
			f:    SpanAttributeValueLength,
		},

		{
			name: "SpanAttributeCount",
			keys: []string{SpanAttributeCountKey, AttributeCountKey},
			f:    SpanAttributeCount,
		},

		{
			name: "SpanEventCount",
			keys: []string{SpanEventCountKey},
			f:    SpanEventCount,
		},

		{
			name: "SpanEventAttributeCount",
			keys: []string{SpanEventAttributeCountKey},
			f:    SpanEventAttributeCount,
		},

		{
			name: "SpanLinkCount",
			keys: []string{SpanLinkCountKey},
			f:    SpanLinkCount,
		},

		{
			name: "SpanLinkAttributeCount",
			keys: []string{SpanLinkAttributeCountKey},
			f:    SpanLinkAttributeCount,
		},
	}

	const (
		defVal    = 500
		envVal    = 2500
		envValStr = "2500"
		invalid   = "localhost"
		empty     = ""
	)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, key := range tc.keys {
				t.Run(key, func(t *testing.T) {
					envStore := ottest.NewEnvStore()
					t.Cleanup(func() { require.NoError(t, envStore.Restore()) })
					envStore.Record(key)

					assert.Equal(t, defVal, tc.f(defVal), "environment variable unset")

					t.Setenv(key, envValStr)
					assert.Equal(t, envVal, tc.f(defVal), "environment variable set/valid")

					t.Setenv(key, invalid)
					assert.Equal(t, defVal, tc.f(defVal), "invalid value")

					t.Setenv(key, empty)
					assert.Equal(t, defVal, tc.f(defVal), "empty value")
				})
			}
		})
	}
}
