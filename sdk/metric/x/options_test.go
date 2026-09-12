// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	apiMetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestWithViewMatchingMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     ViewMatchingMode
		wantMode int
		wantStr  string
	}{
		{
			name:     "Independent",
			mode:     ViewMatchingModeIndependent,
			wantMode: 0,
			wantStr:  "independent",
		},
		{
			name:     "Composable",
			mode:     ViewMatchingModeComposable,
			wantMode: 1,
			wantStr:  "composable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := WithViewMatchingMode(tt.mode)
			require.NotNil(t, opt)

			modeOpt, ok := opt.(viewMatchingModeOption)
			require.True(t, ok)
			assert.Equal(t, tt.wantMode, modeOpt.ViewMatchingMode())
			assert.Equal(t, tt.wantStr, tt.mode.String())

			// Verify it implements experimentalOption
			type experimentalOption interface {
				Experimental()
			}
			_, ok = opt.(experimentalOption)
			assert.True(t, ok)
		})
	}
}

func TestViewMatchingModeStringUnknown(t *testing.T) {
	assert.Equal(t, "unknown(99)", ViewMatchingMode(99).String())
}

func TestMeterProviderIntegration(t *testing.T) {
	r := metric.NewManualReader()
	v1 := metric.NewView(
		metric.Instrument{Name: "http.server.duration"},
		metric.Stream{Name: "http.latency", Unit: "ms"},
	)
	v2 := metric.NewView(
		metric.Instrument{Name: "http.server.duration"},
		metric.Stream{Description: "request latency", AttributeFilter: attribute.NewDenyKeysFilter("http.flavor")},
	)
	mp := metric.NewMeterProvider(
		metric.WithReader(r),
		metric.WithView(v1, v2),
		WithViewMatchingMode(ViewMatchingModeComposable),
	)

	ctx := t.Context()
	m := mp.Meter("test")
	hist, err := m.Float64Histogram("http.server.duration")
	require.NoError(t, err)

	hist.Record(ctx, 42.0, apiMetric.WithAttributes(
		attribute.String("http.method", "GET"),
		attribute.String("http.flavor", "1.1"),
	))

	var rm metricdata.ResourceMetrics
	err = r.Collect(ctx, &rm)
	require.NoError(t, err)
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)

	mMetric := rm.ScopeMetrics[0].Metrics[0]
	assert.Equal(t, "http.latency", mMetric.Name)
	assert.Equal(t, "request latency", mMetric.Description)
	assert.Equal(t, "ms", mMetric.Unit)

	histData, ok := mMetric.Data.(metricdata.Histogram[float64])
	require.True(t, ok)
	require.Len(t, histData.DataPoints, 1)
	dp := histData.DataPoints[0]
	assert.True(t, dp.Attributes.HasValue("http.method"))
	assert.False(t, dp.Attributes.HasValue("http.flavor"))
}
