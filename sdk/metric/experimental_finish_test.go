// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type testFinishOption bool

func (testFinishOption) apply(config) config {
	panic("experimental options must not be applied as stable options")
}

func (testFinishOption) Experimental() {}

func (o testFinishOption) FinishEnabled() bool { return bool(o) }

type testUnknownAggregation struct{}

func (a testUnknownAggregation) copy() Aggregation { return a }

func (testUnknownAggregation) err() error { return nil }

func TestFinishOptionDisabled(t *testing.T) {
	factory, shutdown := newExperimentalMeterFactory([]experimentalOption{testFinishOption(false)})

	assert.NotNil(t, factory)
	assert.Nil(t, shutdown)
}

func TestExperimentalInt64CounterFinish(t *testing.T) {
	reader := NewManualReader()
	provider := NewMeterProvider(
		WithReader(reader),
		WithView(
			NewView(
				Instrument{Name: "requests"},
				Stream{
					Name:            "requests.by_route",
					AttributeFilter: attribute.NewAllowKeysFilter("route"),
				},
			),
			NewView(
				Instrument{Name: "requests"},
				Stream{
					Name:        "requests.dropped",
					Aggregation: AggregationDrop{},
				},
			),
		),
		testFinishOption(true),
	)

	counter, err := provider.Meter("test").Int64Counter("requests")
	require.NoError(t, err)
	finisher, ok := counter.(metricx.Finisher)
	require.True(t, ok)

	attrs := []attribute.KeyValue{
		attribute.String("zone", "east"),
		attribute.String("route", "/old"),
		attribute.String("route", "/orders"),
	}
	counter.Add(t.Context(), 3, api.WithAttributes(attrs...))
	finisher.Finish(t.Context(), attrs...)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
	sum, ok := rm.ScopeMetrics[0].Metrics[0].Data.(metricdata.Sum[int64])
	require.True(t, ok)
	require.Len(t, sum.DataPoints, 1)
	assert.Equal(t, int64(3), sum.DataPoints[0].Value)
	assert.Equal(t, attribute.NewSet(attribute.String("route", "/orders")), sum.DataPoints[0].Attributes)

	assert.NoError(t, provider.Shutdown(t.Context()))
}

func TestExperimentalInt64CounterFinishAggregationError(t *testing.T) {
	reader := NewManualReader()
	provider := NewMeterProvider(
		WithReader(reader),
		WithView(NewView(
			Instrument{Name: "requests"},
			Stream{Aggregation: testUnknownAggregation{}},
		)),
		testFinishOption(true),
	)
	t.Cleanup(func() { _ = provider.Shutdown(t.Context()) })

	_, err := provider.Meter("test").Int64Counter("requests")
	assert.ErrorIs(t, err, errUnknownAggregation)
}
