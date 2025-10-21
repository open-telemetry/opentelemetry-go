// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"testing"
)

type errMeterProvider struct {
	mapi.MeterProvider
	err error
}

func (m *errMeterProvider) Meter(string, ...mapi.MeterOption) mapi.Meter {
	return &errMeter{err: m.err}
}

type errMeter struct {
	mapi.Meter
	err error
}

func (m *errMeter) Int64Counter(string, ...mapi.Int64CounterOption) (mapi.Int64Counter, error) {
	return nil, m.err
}

const slpComponentID = 0

func TestNewSLPError(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	errMp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(errMp)

	_, err := NewSLP(slpComponentID)
	require.ErrorIs(t, err, assert.AnError)
	assert.ErrorContains(t, err, "failed to create a processed log metric")
}

func setup(t *testing.T) func() metricdata.ScopeMetrics {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() {
		otel.SetMeterProvider(orig)
	})

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(mp)

	return func() metricdata.ScopeMetrics {
		var got metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(t.Context(), &got))
		if len(got.ScopeMetrics) != 1 {
			return metricdata.ScopeMetrics{}
		}
		return got.ScopeMetrics[0]
	}
}
