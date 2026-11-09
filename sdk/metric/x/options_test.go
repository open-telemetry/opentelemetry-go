// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	apiMetric "go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
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

func TestMeterProviderComposableViewEnabled(t *testing.T) {
	t.Run("OptInInstrumentEnabledAndComposed", func(t *testing.T) {
		r := metric.NewManualReader()
		enabled := true
		v1 := NewView(
			metric.Instrument{Name: "optin.counter"},
			Stream{Enabled: &enabled},
		)
		v2 := metric.NewView(
			metric.Instrument{Name: "optin.counter"},
			metric.Stream{Description: "composed opt-in counter"},
		)
		mp := metric.NewMeterProvider(
			metric.WithReader(r),
			WithView(v1),
			metric.WithView(v2),
			WithViewMatchingMode(ViewMatchingModeComposable),
		)

		ctx := t.Context()
		m := mp.Meter("test")
		ctr, err := m.Int64Counter("optin.counter", metricx.WithOptIn())
		require.NoError(t, err)
		assert.True(t, ctr.Enabled(ctx))

		ctr.Add(ctx, 7)

		var rm metricdata.ResourceMetrics
		err = r.Collect(ctx, &rm)
		require.NoError(t, err)
		require.Len(t, rm.ScopeMetrics, 1)
		require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
		assert.Equal(t, "optin.counter", rm.ScopeMetrics[0].Metrics[0].Name)
		assert.Equal(t, "composed opt-in counter", rm.ScopeMetrics[0].Metrics[0].Description)
	})

	t.Run("StandardInstrumentDisabledByComposableView", func(t *testing.T) {
		r := metric.NewManualReader()
		disabled := false
		v1 := metric.NewView(
			metric.Instrument{Name: "std.counter"},
			metric.Stream{Description: "custom desc"},
		)
		v2 := NewView(
			metric.Instrument{Name: "std.counter"},
			Stream{Enabled: &disabled},
		)
		mp := metric.NewMeterProvider(
			metric.WithReader(r),
			metric.WithView(v1),
			WithView(v2),
			WithViewMatchingMode(ViewMatchingModeComposable),
		)

		ctx := t.Context()
		m := mp.Meter("test")
		ctr, err := m.Int64Counter("std.counter")
		require.NoError(t, err)
		assert.False(t, ctr.Enabled(ctx))

		ctr.Add(ctx, 7)

		var rm metricdata.ResourceMetrics
		err = r.Collect(ctx, &rm)
		require.NoError(t, err)
		assert.Empty(t, rm.ScopeMetrics)
	})
}

func TestMeterProviderComposableUnspecifiedStreamNameAndProperties(t *testing.T) {
	r := metric.NewManualReader()
	// Renaming views create two target stream groups.
	vRenameMS := metric.NewView(
		metric.Instrument{Name: "req.duration"},
		metric.Stream{Name: "req.duration.ms", Unit: "ms"},
	)
	vRenameS := metric.NewView(
		metric.Instrument{Name: "req.duration"},
		metric.Stream{Name: "req.duration.s", Unit: "s"},
	)
	// View returning instrument's own Name (via NewView) leaves Name and Unit unspecified.
	vInstrumentName := metric.NewView(
		metric.Instrument{Name: "req.duration"},
		metric.Stream{Description: "shared description"},
	)
	// Direct View returning empty Name and instrument's own Description/Unit leaves Name, Description,
	// and Unit unspecified while contributing an AttributeFilter to both stream groups.
	vEmptyName := metric.View(func(inst metric.Instrument) (metric.Stream, bool) {
		if inst.Name == "req.duration" {
			return metric.Stream{
				Name:            "",
				Description:     inst.Description,
				Unit:            inst.Unit,
				AttributeFilter: attribute.NewAllowKeysFilter("keep"),
			}, true
		}
		return metric.Stream{}, false
	})

	mp := metric.NewMeterProvider(
		metric.WithReader(r),
		metric.WithView(vRenameMS, vRenameS, vInstrumentName, vEmptyName),
		WithViewMatchingMode(ViewMatchingModeComposable),
	)

	ctx := t.Context()
	m := mp.Meter("test")
	ctr, err := m.Int64Counter(
		"req.duration",
		apiMetric.WithDescription("orig description"),
		apiMetric.WithUnit("1"),
	)
	require.NoError(t, err)

	ctr.Add(ctx, 10, apiMetric.WithAttributes(
		attribute.String("keep", "yes"),
		attribute.String("drop", "no"),
	))

	var rm metricdata.ResourceMetrics
	err = r.Collect(ctx, &rm)
	require.NoError(t, err)
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 2)

	byName := make(map[string]metricdata.Metrics, 2)
	for _, mMetric := range rm.ScopeMetrics[0].Metrics {
		byName[mMetric.Name] = mMetric
	}
	require.Contains(t, byName, "req.duration.ms")
	require.Contains(t, byName, "req.duration.s")

	for name, wantUnit := range map[string]string{"req.duration.ms": "ms", "req.duration.s": "s"} {
		mMetric := byName[name]
		assert.Equal(t, "shared description", mMetric.Description)
		assert.Equal(t, wantUnit, mMetric.Unit)
		sum, ok := mMetric.Data.(metricdata.Sum[int64])
		require.True(t, ok)
		require.Len(t, sum.DataPoints, 1)
		assert.Equal(t, attribute.NewSet(attribute.String("keep", "yes")), sum.DataPoints[0].Attributes)
	}
}
