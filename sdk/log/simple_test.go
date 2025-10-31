// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"

	"go.opentelemetry.io/otel/attribute"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/internal/x"
	metricSDK "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

// extractAttributeValue extracts the value of a specific attribute from the first DataPoint found.
func extractAttributeValue(data any, attrKey attribute.Key) (attribute.Value, bool) {
	switch d := data.(type) {
	case metricdata.ResourceMetrics:
		for _, scope := range d.ScopeMetrics {
			for _, m := range scope.Metrics {
				if val, ok := extractAttributeValue(m.Data, attrKey); ok {
					return val, true
				}
			}
		}
	case metricdata.Sum[int64]:
		for _, dp := range d.DataPoints {
			if val, ok := dp.Attributes.Value(attrKey); ok {
				return val, true
			}
		}
	case metricdata.Sum[float64]:
		for _, dp := range d.DataPoints {
			if val, ok := dp.Attributes.Value(attrKey); ok {
				return val, true
			}
		}
	case metricdata.Gauge[int64]:
		for _, dp := range d.DataPoints {
			if val, ok := dp.Attributes.Value(attrKey); ok {
				return val, true
			}
		}
	case metricdata.Gauge[float64]:
		for _, dp := range d.DataPoints {
			if val, ok := dp.Attributes.Value(attrKey); ok {
				return val, true
			}
		}
	}
	return attribute.Value{}, false
}

type exporter struct {
	records []log.Record

	exportCalled     bool
	shutdownCalled   bool
	forceFlushCalled bool
}

func (e *exporter) Export(_ context.Context, r []log.Record) error {
	e.records = r
	e.exportCalled = true
	return nil
}

func (e *exporter) Shutdown(context.Context) error {
	e.shutdownCalled = true
	return nil
}

func (e *exporter) ForceFlush(context.Context) error {
	e.forceFlushCalled = true
	return nil
}

func TestSimpleProcessorOnEmit(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)

	r := new(log.Record)
	r.SetSeverityText("test")
	_ = s.OnEmit(t.Context(), r)

	require.True(t, e.exportCalled, "exporter Export not called")
	assert.Equal(t, []log.Record{*r}, e.records)
}

func TestSimpleProcessorShutdown(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)
	_ = s.Shutdown(t.Context())
	require.True(t, e.shutdownCalled, "exporter Shutdown not called")
}

func TestSimpleProcessorForceFlush(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)
	_ = s.ForceFlush(t.Context())
	require.True(t, e.forceFlushCalled, "exporter ForceFlush not called")
}

type writerExporter struct {
	io.Writer
}

func (e *writerExporter) Export(_ context.Context, records []log.Record) error {
	for _, r := range records {
		_, _ = io.WriteString(e.Writer, r.Body().String())
	}
	return nil
}

func (*writerExporter) Shutdown(context.Context) error {
	return nil
}

func (*writerExporter) ForceFlush(context.Context) error {
	return nil
}

func TestSimpleProcessorEmpty(t *testing.T) {
	assert.NotPanics(t, func() {
		var s log.SimpleProcessor
		ctx := t.Context()
		record := new(log.Record)
		assert.NoError(t, s.OnEmit(ctx, record), "OnEmit")
		assert.NoError(t, s.ForceFlush(ctx), "ForceFlush")
		assert.NoError(t, s.Shutdown(ctx), "Shutdown")
	})
}

func TestSimpleProcessorConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	r := new(log.Record)
	r.SetSeverityText("test")
	ctx := t.Context()
	e := &writerExporter{new(strings.Builder)}
	s := log.NewSimpleProcessor(e)
	for range goRoutineN {
		go func() {
			defer wg.Done()

			_ = s.OnEmit(ctx, r)
			_ = s.Shutdown(ctx)
			_ = s.ForceFlush(ctx)
		}()
	}

	wg.Wait()
}

type errorExporter struct {
	err error
}

func (e *errorExporter) Export(_ context.Context, _ []log.Record) error {
	return e.err
}

func (*errorExporter) Shutdown(context.Context) error {
	return nil
}

func (*errorExporter) ForceFlush(context.Context) error {
	return nil
}

type failingMeterProvider struct {
	noop.MeterProvider
}

func (*failingMeterProvider) Meter(name string, opts ...metric.MeterOption) metric.Meter {
	return &failingMeter{Meter: noop.NewMeterProvider().Meter(name, opts...)}
}

type failingMeter struct {
	metric.Meter
}

func (*failingMeter) Int64Counter(_ string, _ ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return nil, errors.New("failed to create counter")
}

func TestSimpleProcessorSelfObservability(t *testing.T) {
	originalMP := otel.GetMeterProvider()
	setupCleanMeterProvider := func(t *testing.T) {
		t.Cleanup(func() {
			otel.SetMeterProvider(originalMP)
		})
	}

	t.Run("self observability disabled", func(t *testing.T) {
		setupCleanMeterProvider(t)

		reader := metricSDK.NewManualReader()
		mp := metricSDK.NewMeterProvider(metricSDK.WithReader(reader))
		otel.SetMeterProvider(mp)

		e := new(exporter)
		s := log.NewSimpleProcessor(e)

		r := new(log.Record)
		r.SetSeverityText("test")
		_ = s.OnEmit(context.Background(), r)

		require.True(t, e.exportCalled)
		assert.Equal(t, []log.Record{*r}, e.records)

		rm := metricdata.ResourceMetrics{}
		err := reader.Collect(context.Background(), &rm)
		require.NoError(t, err)

		expected := metricdata.ResourceMetrics{
			Resource:     rm.Resource,
			ScopeMetrics: []metricdata.ScopeMetrics{},
		}

		metricdatatest.AssertEqual(t, expected, rm, metricdatatest.IgnoreTimestamp())
	})

	t.Run("self observability enabled without error", func(t *testing.T) {
		setupCleanMeterProvider(t)

		t.Setenv(x.SelfObservability.Key(), "true")

		reader := metricSDK.NewManualReader()
		mp := metricSDK.NewMeterProvider(metricSDK.WithReader(reader))
		otel.SetMeterProvider(mp)

		e := new(exporter)
		s := log.NewSimpleProcessor(e)

		r := new(log.Record)
		r.SetSeverityText("test")

		var err error
		err = s.OnEmit(context.Background(), r)
		require.NoError(t, err)

		err = s.OnEmit(context.Background(), r)
		require.NoError(t, err)

		err = s.OnEmit(context.Background(), r)
		require.NoError(t, err)

		// First collect to get the actual component name
		rm := metricdata.ResourceMetrics{}
		err = reader.Collect(context.Background(), &rm)
		require.NoError(t, err)
		require.Len(t, rm.ScopeMetrics, 1)

		// Extract the actual component name and use it in expected structure
		componentVal, ok := extractAttributeValue(rm, "otel.component.name")
		require.True(t, ok, "component name attribute should be present")
		actualComponentName := componentVal.AsString()

		expected := metricdata.ResourceMetrics{
			ScopeMetrics: []metricdata.ScopeMetrics{
				{
					Scope: instrumentation.Scope{
						Name:      "go.opentelemetry.io/otel/sdk/log",
						Version:   sdk.Version(),
						SchemaURL: semconv.SchemaURL,
					},
					Metrics: []metricdata.Metrics{
						{
							Name:        "otel.sdk.processor.log.processed",
							Description: "The number of log records for which the processing has finished, either successful or failed",
							Unit:        "{log_record}",
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value: 3,
										Attributes: attribute.NewSet(
											attribute.String(
												"otel.component.type",
												string(otelconv.ComponentTypeSimpleLogProcessor),
											),
											attribute.String("otel.component.name", actualComponentName),
										),
										Exemplars: []metricdata.Exemplar[int64]{},
									},
								},
								Temporality: metricdata.CumulativeTemporality,
								IsMonotonic: true,
							},
						},
					},
				},
			},
		}

		require.Len(t, rm.ScopeMetrics, 1)
		metricdatatest.AssertEqual(t, expected.ScopeMetrics[0], rm.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())
	})

	t.Run("self observability enabled with error", func(t *testing.T) {
		setupCleanMeterProvider(t)

		t.Setenv(x.SelfObservability.Key(), "true")

		reader := metricSDK.NewManualReader()
		mp := metricSDK.NewMeterProvider(metricSDK.WithReader(reader))
		otel.SetMeterProvider(mp)

		e := &errorExporter{err: errors.New("export failed")}
		s := log.NewSimpleProcessor(e)

		r := new(log.Record)
		r.SetSeverityText("test")

		var err error
		err = s.OnEmit(context.Background(), r)
		require.Error(t, err)
		assert.Equal(t, "export failed", err.Error())

		err = s.OnEmit(context.Background(), r)
		require.Error(t, err)
		assert.Equal(t, "export failed", err.Error())

		rm := metricdata.ResourceMetrics{}
		collectErr := reader.Collect(context.Background(), &rm)
		require.NoError(t, collectErr)
		require.Len(t, rm.ScopeMetrics, 1)

		componentVal, ok := extractAttributeValue(rm, "otel.component.name")
		require.True(t, ok)
		actualComponentName := componentVal.AsString()

		expected := metricdata.ResourceMetrics{
			ScopeMetrics: []metricdata.ScopeMetrics{
				{
					Scope: instrumentation.Scope{
						Name:      "go.opentelemetry.io/otel/sdk/log",
						Version:   sdk.Version(),
						SchemaURL: semconv.SchemaURL,
					},
					Metrics: []metricdata.Metrics{
						{
							Name:        "otel.sdk.processor.log.processed",
							Description: "The number of log records for which the processing has finished, either successful or failed",
							Unit:        "{log_record}",
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value: 2,
										Attributes: attribute.NewSet(
											attribute.String(
												"otel.component.type",
												string(otelconv.ComponentTypeSimpleLogProcessor),
											),
											attribute.String("otel.component.name", actualComponentName),
											attribute.String("error.type", string(otelconv.ErrorTypeOther)),
										),
										Exemplars: []metricdata.Exemplar[int64]{},
									},
								},
								Temporality: metricdata.CumulativeTemporality,
								IsMonotonic: true,
							},
						},
					},
				},
			},
		}

		require.Len(t, rm.ScopeMetrics, 1)
		metricdatatest.AssertEqual(t, expected.ScopeMetrics[0], rm.ScopeMetrics[0], metricdatatest.IgnoreTimestamp())
	})

	t.Run("self observability metric creation error handled", func(t *testing.T) {
		setupCleanMeterProvider(t)

		t.Setenv(x.SelfObservability.Key(), "true")

		failingMP := &failingMeterProvider{}
		otel.SetMeterProvider(failingMP)

		assert.NotPanics(t, func() {
			e := new(exporter)
			s := log.NewSimpleProcessor(e)

			r := new(log.Record)
			r.SetSeverityText("test")
			_ = s.OnEmit(context.Background(), r)

			require.True(t, e.exportCalled)
			assert.Equal(t, []log.Record{*r}, e.records)
		})
	})
}

func BenchmarkSimpleProcessorOnEmit(b *testing.B) {
	r := new(log.Record)
	r.SetSeverityText("test")
	ctx := b.Context()
	s := log.NewSimpleProcessor(nil)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var out error

		for pb.Next() {
			out = s.OnEmit(ctx, r)
		}

		_ = out
	})
}
