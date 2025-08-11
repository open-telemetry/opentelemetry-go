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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/log/internal/x"
	metricSDK "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

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
	_ = s.OnEmit(context.Background(), r)

	require.True(t, e.exportCalled, "exporter Export not called")
	assert.Equal(t, []log.Record{*r}, e.records)
}

func TestSimpleProcessorShutdown(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)
	_ = s.Shutdown(context.Background())
	require.True(t, e.shutdownCalled, "exporter Shutdown not called")
}

func TestSimpleProcessorForceFlush(t *testing.T) {
	e := new(exporter)
	s := log.NewSimpleProcessor(e)
	_ = s.ForceFlush(context.Background())
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
		ctx := context.Background()
		record := new(log.Record)
		assert.NoError(t, s.OnEmit(ctx, record), "OnEmit")
		assert.NoError(t, s.ForceFlush(ctx), "ForceFlush")
		assert.NoError(t, s.Shutdown(ctx), "Shutdown")
	})
}

func TestSimpleProcessorConcurrentSafe(*testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	r := new(log.Record)
	r.SetSeverityText("test")
	ctx := context.Background()
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
		t.Cleanup(func() { otel.SetMeterProvider(originalMP) })
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

		for _, scopeMetrics := range rm.ScopeMetrics {
			for _, m := range scopeMetrics.Metrics {
				if m.Name == "otel.sdk.processor.log.processed" {
					t.Errorf("expected no self-observability metrics when disabled, but found metric: %s", m.Name)
				}
			}
		}
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

		rm := metricdata.ResourceMetrics{}
		err = reader.Collect(context.Background(), &rm)
		require.NoError(t, err)

		var processedMetric *metricdata.ScopeMetrics
		for _, scopeMetrics := range rm.ScopeMetrics {
			for _, m := range scopeMetrics.Metrics {
				if m.Name == "otel.sdk.processor.log.processed" {
					processedMetric = &scopeMetrics
					break
				}
			}
		}

		require.NotNil(t, processedMetric)

		totalCount, hasComponentType, hasComponentName := extractProcessedLogMetricsSuccess(processedMetric)

		assert.Equal(t, int64(3), totalCount)
		assert.True(t, hasComponentType)
		assert.True(t, hasComponentName)
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

		var processedMetric *metricdata.ScopeMetrics
		for _, scopeMetrics := range rm.ScopeMetrics {
			for _, m := range scopeMetrics.Metrics {
				if m.Name == "otel.sdk.processor.log.processed" {
					processedMetric = &scopeMetrics
					break
				}
			}
		}

		require.NotNil(t, processedMetric)

		totalCount, hasErrorType, hasComponentType, hasComponentName := extractProcessedLogMetricsError(
			processedMetric,
		)

		assert.Equal(t, int64(2), totalCount)
		assert.True(t, hasErrorType)
		assert.True(t, hasComponentType)
		assert.True(t, hasComponentName)
	})

	t.Run("self observability enabled", func(t *testing.T) {
		setupCleanMeterProvider(t)

		t.Setenv(x.SelfObservability.Key(), "true")

		otel.SetMeterProvider(noop.NewMeterProvider())

		e := new(exporter)
		s := log.NewSimpleProcessor(e)

		r := new(log.Record)
		r.SetSeverityText("test")
		assert.NotPanics(t, func() {
			_ = s.OnEmit(context.Background(), r)
		})

		require.True(t, e.exportCalled)
		assert.Equal(t, []log.Record{*r}, e.records)
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

func extractProcessedLogMetricsSuccess(
	processedMetric *metricdata.ScopeMetrics,
) (totalCount int64, hasComponentType, hasComponentName bool) {
	for _, m := range processedMetric.Metrics {
		if m.Name != "otel.sdk.processor.log.processed" {
			continue
		}

		data, ok := m.Data.(metricdata.Sum[int64])
		if !ok {
			continue
		}

		for _, dataPoint := range data.DataPoints {
			totalCount += dataPoint.Value
			for _, attr := range dataPoint.Attributes.ToSlice() {
				switch attr.Key {
				case "otel.component.type":
					if attr.Value.AsString() == string(otelconv.ComponentTypeSimpleLogProcessor) {
						hasComponentType = true
					}
				case "otel.component.name":
					if strings.HasPrefix(attr.Value.AsString(), "simple_log_processor/") {
						hasComponentName = true
					}
				}
			}
		}
	}
	return totalCount, hasComponentType, hasComponentName
}

func extractProcessedLogMetricsError(
	processedMetric *metricdata.ScopeMetrics,
) (totalCount int64, hasErrorType, hasComponentType, hasComponentName bool) {
	for _, m := range processedMetric.Metrics {
		if m.Name != "otel.sdk.processor.log.processed" {
			continue
		}

		data, ok := m.Data.(metricdata.Sum[int64])
		if !ok {
			continue
		}

		for _, dataPoint := range data.DataPoints {
			totalCount += dataPoint.Value
			for _, attr := range dataPoint.Attributes.ToSlice() {
				switch attr.Key {
				case "error.type":
					if attr.Value.AsString() == string(otelconv.ErrorTypeOther) {
						hasErrorType = true
					}
				case "otel.component.type":
					if attr.Value.AsString() == string(otelconv.ComponentTypeSimpleLogProcessor) {
						hasComponentType = true
					}
				case "otel.component.name":
					hasComponentName = true
				}
			}
		}
	}
	return totalCount, hasErrorType, hasComponentType, hasComponentName
}

func BenchmarkSimpleProcessorOnEmit(b *testing.B) {
	r := new(log.Record)
	r.SetSeverityText("test")
	ctx := context.Background()
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
