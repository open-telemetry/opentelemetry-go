// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	metricx "go.opentelemetry.io/otel/metric/x"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestExperimentalGaugeFinishConstruction(t *testing.T) {
	for _, number := range []string{"Int64", "Float64"} {
		t.Run(number, func(t *testing.T) {
			create := func(m api.Meter, name string, opts ...api.InstrumentOption) (metricx.Finisher, error) {
				var instrument any
				var err error
				if number == "Int64" {
					options := make([]api.Int64GaugeOption, len(opts))
					for i, opt := range opts {
						options[i] = opt
					}
					instrument, err = m.Int64Gauge(name, options...)
				} else {
					options := make([]api.Float64GaugeOption, len(opts))
					for i, opt := range opts {
						options[i] = opt
					}
					instrument, err = m.Float64Gauge(name, options...)
				}
				finisher, ok := instrument.(metricx.Finisher)
				require.True(t, ok)
				return finisher, err
			}
			t.Run("OptionsAndShutdown", func(t *testing.T) {
				reader := NewManualReader()
				provider := NewMeterProvider(WithReader(reader), testFinishOption(true))
				meter := provider.Meter("test")
				opts := []api.InstrumentOption{
					api.WithDescription("Active entities"),
					api.WithUnit("{entity}"),
					metricx.WithDefaultAttributes("entity"),
				}
				gauge, err := create(meter, "active", opts...)
				require.NoError(t, err)
				again, err := create(meter, "active", opts...)
				require.NoError(t, err)
				assert.Same(t, gauge, again)
				attrs := []attribute.KeyValue{attribute.String("zone", "east"), attribute.String("entity", "one")}
				record := func() {
					if number == "Int64" {
						gauge.(api.Int64Gauge).Record(t.Context(), -3, api.WithAttributes(attrs...))
					} else {
						gauge.(api.Float64Gauge).Record(t.Context(), -0.75, api.WithAttributes(attrs...))
					}
				}
				record()
				gauge.Finish(t.Context(), attrs...)
				var rm metricdata.ResourceMetrics
				require.NoError(t, reader.Collect(t.Context(), &rm))
				require.Len(t, rm.ScopeMetrics, 1)
				require.Len(t, rm.ScopeMetrics[0].Metrics, 1)
				m := rm.ScopeMetrics[0].Metrics[0]
				assert.Equal(t, "Active entities", m.Description)
				assert.Equal(t, "{entity}", m.Unit)
				want := attribute.NewSet(attribute.String("entity", "one"))
				if number == "Int64" {
					data := m.Data.(metricdata.Gauge[int64])
					require.Len(t, data.DataPoints, 1)
					assert.Equal(t, want, data.DataPoints[0].Attributes)
				} else {
					data := m.Data.(metricdata.Gauge[float64])
					require.Len(t, data.DataPoints, 1)
					assert.Equal(t, want, data.DataPoints[0].Attributes)
				}
				require.NoError(t, reader.Collect(t.Context(), &rm))
				assert.Empty(t, rm.ScopeMetrics)
				record() // Leave live state for shutdown to retire.
				require.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context())))
				record()
				gauge.Finish(t.Context(), attrs...)
				require.NoError(t, provider.pipes[0].produce(t.Context(), &rm))
				assert.Empty(t, rm.ScopeMetrics, "shutdown must stop the Finish-aware aggregation")
			})
			t.Run("Errors", func(t *testing.T) {
				provider := NewMeterProvider(testFinishOption(true))
				t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.WithoutCancel(t.Context()))) })
				_, err := create(provider.Meter("test"), "invalid name")
				assert.ErrorIs(t, err, ErrInstrumentName)
				reader := NewManualReader()
				failedProvider := NewMeterProvider(
					WithReader(reader),
					testFinishOption(true),
					WithView(NewView(Instrument{Name: "active"}, Stream{Aggregation: testUnknownAggregation{}})),
				)
				t.Cleanup(func() { assert.NoError(t, failedProvider.Shutdown(context.WithoutCancel(t.Context()))) })
				gauge, err := create(failedProvider.Meter("test"), "active")
				assert.ErrorIs(t, err, errUnknownAggregation)
				assert.NotPanics(t, func() { gauge.Finish(t.Context()) })
			})
		})
	}
}
