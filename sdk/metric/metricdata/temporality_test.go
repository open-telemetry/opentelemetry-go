// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metricdata_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestTemporalityString(t *testing.T) {
	testCases := []struct {
		name        string
		temporality metricdata.Temporality
		want        string
	}{
		{
			name:        "undefined",
			temporality: metricdata.Temporality(0),
			want:        "undefinedTemporality",
		},
		{
			name:        "cumulative",
			temporality: metricdata.CumulativeTemporality,
			want:        "CumulativeTemporality",
		},
		{
			name:        "delta",
			temporality: metricdata.DeltaTemporality,
			want:        "DeltaTemporality",
		},
		{
			name:        "out-of-range positive",
			temporality: metricdata.Temporality(99),
			want:        "Temporality(99)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.temporality.String())
		})
	}
}

func TestTemporalityMarshalText(t *testing.T) {
	testCases := []struct {
		name        string
		temporality metricdata.Temporality
		want        string
	}{
		{
			name:        "cumulative",
			temporality: metricdata.CumulativeTemporality,
			want:        "CumulativeTemporality",
		},
		{
			name:        "delta",
			temporality: metricdata.DeltaTemporality,
			want:        "DeltaTemporality",
		},
		{
			name:        "undefined",
			temporality: metricdata.Temporality(0),
			want:        "undefinedTemporality",
		},
		{
			name:        "unknown",
			temporality: metricdata.Temporality(255),
			want:        "Temporality(255)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.temporality.MarshalText()
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(data))
		})
	}
}
