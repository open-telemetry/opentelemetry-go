// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
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

func setup(t *testing.T) (*SLP, func() metricdata.ScopeMetrics) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	orig := otel.GetMeterProvider()
	t.Cleanup(func() {
		otel.SetMeterProvider(orig)
	})

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	otel.SetMeterProvider(mp)

	slp, err := NewSLP(slpComponentID)
	require.NoError(t, err)
	require.NotNil(t, slp)

	return slp, func() metricdata.ScopeMetrics {
		var got metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(t.Context(), &got))
		require.Len(t, got.ScopeMetrics, 1)
		return got.ScopeMetrics[0]
	}
}

func processedMetric(err error) metricdata.Metrics {
	processed := &otelconv.SDKProcessorLogProcessed{}

	attrs := []attribute.KeyValue{
		GetComponentName(slpComponentID),
		processed.AttrComponentType(otelconv.ComponentTypeSimpleLogProcessor),
	}

	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}

	dp := []metricdata.DataPoint[int64]{
		{
			Attributes: attribute.NewSet(attrs...),
			Value:      1,
		},
	}

	return metricdata.Metrics{
		Name:        processed.Name(),
		Description: processed.Description(),
		Unit:        processed.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dp,
		},
	}
}

var Scope = instrumentation.Scope{
	Name:      ScopeName,
	Version:   sdk.Version(),
	SchemaURL: semconv.SchemaURL,
}

func assertMetric(t *testing.T, got metricdata.ScopeMetrics, err error) {
	t.Helper()
	assert.Equal(t, Scope, got.Scope, "unexpected scope")
	m := got.Metrics
	require.Len(t, m, 1, "expected 1 metrics")

	o := metricdatatest.IgnoreTimestamp()
	want := processedMetric(err)

	metricdatatest.AssertEqual(t, want, m[0], o)
}

func TestSLP(t *testing.T) {
	t.Run("NoError", func(t *testing.T) {
		slp, collect := setup(t)
		slp.LogProcessed(t.Context(), nil)
		assertMetric(t, collect(), nil)
	})

	t.Run("Error", func(t *testing.T) {
		processErr := errors.New("error processing log")
		slp, collect := setup(t)
		slp.LogProcessed(t.Context(), processErr)
		assertMetric(t, collect(), processErr)
	})
}
