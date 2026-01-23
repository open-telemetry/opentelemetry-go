// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

func TestNextExporterID(t *testing.T) {
	SetSimpleProcessorID(0)

	var expected int64
	for range 10 {
		id := NextSimpleProcessorID()
		assert.Equal(t, expected, id)
		expected++
	}
}

func TestSetExporterID(t *testing.T) {
	SetSimpleProcessorID(0)

	prev := SetSimpleProcessorID(42)
	assert.Equal(t, int64(0), prev)

	id := NextSimpleProcessorID()
	assert.Equal(t, int64(42), id)
}

func TestNextExporterIDConcurrentSafe(t *testing.T) {
	SetSimpleProcessorID(0)

	const goroutines = 100
	const increments = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range increments {
				NextSimpleProcessorID()
			}
		}()
	}

	wg.Wait()

	expected := int64(goroutines * increments)
	id := NextSimpleProcessorID()
	assert.Equal(t, expected, id)
}

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

func (m *errMeter) Float64Histogram(string, ...mapi.Float64HistogramOption) (mapi.Float64Histogram, error) {
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

func TestNewSLPDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY
	bsp, err := NewSLP(slpComponentID)
	assert.NoError(t, err)
	assert.Nil(t, bsp)
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
		GetSLPComponentName(slpComponentID),
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

func BenchmarkSLP(b *testing.B) {
	b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	newSLP := func(b *testing.B) *SLP {
		b.Helper()
		slp, err := NewSLP(slpComponentID)
		require.NoError(b, err)
		require.NotNil(b, slp)
		return slp
	}

	b.Run("LogProcessed", func(b *testing.B) {
		orig := otel.GetMeterProvider()
		b.Cleanup(func() {
			otel.SetMeterProvider(orig)
		})

		otel.SetMeterProvider(noop.NewMeterProvider())

		ssp := newSLP(b)
		ctx := b.Context()

		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ssp.LogProcessed(ctx, nil)
			}
		})
	})

	b.Run("LogProcessedWithError", func(b *testing.B) {
		orig := otel.GetMeterProvider()
		b.Cleanup(func() {
			otel.SetMeterProvider(orig)
		})
		otel.SetMeterProvider(noop.NewMeterProvider())
		slp := newSLP(b)
		ctx := b.Context()

		processErr := errors.New("error processing log")

		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				slp.LogProcessed(ctx, processErr)
			}
		})
	})
}
