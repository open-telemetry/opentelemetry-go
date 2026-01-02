// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

const (
	ID = 0
)

type errMeterProvider struct {
	metric.MeterProvider
	err error
}

func (m *errMeterProvider) Meter(string, ...metric.MeterOption) metric.Meter {
	return &errMeter{err: m.err}
}

type errMeter struct {
	metric.Meter
	err error
}

func (e *errMeter) Int64UpDownCounter(string, ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return nil, e.err
}

func (e *errMeter) Int64Counter(string, ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return nil, e.err
}

func (e *errMeter) Float64Histogram(string, ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return nil, e.err
}

func TestNewInstrumentation(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
	t.Run("NoError", func(t *testing.T) {
		inst, err := NewInstrumentation(ID)
		require.NoError(t, err)
		assert.NotNil(t, inst.inflight, "logInflightMetric should be created")
		assert.NotNil(t, inst.exported, "logExportedMetric should be created")
		assert.NotNil(t, inst.duration, "logExportedDurationMetric should be created")
	})

	t.Run("error", func(t *testing.T) {
		orig := otel.GetMeterProvider()
		t.Cleanup(func() {
			otel.SetMeterProvider(orig)
		})
		otel.SetMeterProvider(&errMeterProvider{
			err: assert.AnError,
		})
		_, err := NewInstrumentation(ID)
		require.ErrorIs(t, err, assert.AnError, "new instrument errors")
		assert.ErrorContains(t, err, "inflight metric")
		assert.ErrorContains(t, err, "exported metric")
		assert.ErrorContains(t, err, "duration metric")
	})
}

func set(err error) attribute.Set {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(GetComponentName(ID)),
		semconv.OTelComponentNameKey.String(ComponentType),
	}
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}
	return attribute.NewSet(attrs...)
}

func logInflight() metricdata.Metrics {
	inflight := otelconv.SDKExporterLogInflight{}
	return metricdata.Metrics{
		Name:        inflight.Name(),
		Description: inflight.Description(),
		Unit:        inflight.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: set(nil), Value: 0},
			},
		},
	}
}

func logExported(success, total int64, err error) metricdata.Metrics {
	dp := []metricdata.DataPoint[int64]{
		{Attributes: set(nil), Value: success},
	}

	if err != nil {
		dp = append(dp, metricdata.DataPoint[int64]{
			Attributes: set(err),
			Value:      total - success,
		})
	}

	exported := otelconv.SDKExporterLogExported{}
	return metricdata.Metrics{
		Name:        exported.Name(),
		Description: exported.Description(),
		Unit:        exported.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dp,
		},
	}
}

func logExportedDuration(err error) metricdata.Metrics {
	attrs := set(err)

	duration := otelconv.SDKExporterOperationDuration{}
	return metricdata.Metrics{
		Name:        duration.Name(),
		Description: duration.Description(),
		Unit:        duration.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{Attributes: attrs},
			},
		},
	}
}

func setup(t *testing.T) (*Instrumentation, func() metricdata.ScopeMetrics) {
	t.Helper()
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	original := otel.GetMeterProvider()
	t.Cleanup(func() {
		otel.SetMeterProvider(original)
	})

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	otel.SetMeterProvider(mp)

	inst, err := NewInstrumentation(ID)
	require.NoError(t, err)
	require.NotNil(t, inst)

	return inst, func() metricdata.ScopeMetrics {
		var rm metricdata.ResourceMetrics
		require.NoError(t, reader.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		return rm.ScopeMetrics[0]
	}
}

var Scope = instrumentation.Scope{
	Name:      ScopeName,
	Version:   internal.Version,
	SchemaURL: semconv.SchemaURL,
}

func assertMetrics(
	t *testing.T,
	got metricdata.ScopeMetrics,
	logs int64,
	success int64,
	err error,
) {
	t.Helper()
	assert.Equal(t, Scope, got.Scope)

	m := got.Metrics
	require.Len(t, m, 3)

	o := metricdatatest.IgnoreTimestamp()
	want := logInflight()
	metricdatatest.AssertEqual(t, want, m[0], o)

	want = logExported(success, logs, err)
	metricdatatest.AssertEqual(t, want, m[1], o)

	want = logExportedDuration(err)
	metricdatatest.AssertEqual(t, want, m[2], metricdatatest.IgnoreValue(), o)
}

func TestInstrumentationExportLogs(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	inst.ExportLogs(t.Context(), n).End(nil)
	assertMetrics(t, collect(), n, n, nil)
}

func TestInstrumentationExportLogPartialErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	const success = 5

	err := internal.PartialSuccess{RejectedItems: n - success}
	inst.ExportLogs(t.Context(), n).End(err)

	assertMetrics(t, collect(), n, success, err)
}

func TestInstrumentationExportLogAllErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	const success = 0
	inst.ExportLogs(t.Context(), n).End(assert.AnError)

	assertMetrics(t, collect(), n, success, assert.AnError)
}

func TestInstrumentationExportLogsInvalidPartialErrored(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	err := internal.PartialSuccess{RejectedItems: -5}
	inst.ExportLogs(t.Context(), n).End(err)

	success := int64(n)
	assertMetrics(t, collect(), n, success, err)

	err.RejectedItems = n + 5
	inst.ExportLogs(t.Context(), n).End(err)

	success += 0
	assertMetrics(t, collect(), n+n, success, err)
}

func BenchmarkInstrumentationExportLogs(b *testing.B) {
	setup := func(tb *testing.B) *Instrumentation {
		tb.Helper()
		tb.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		inst, err := NewInstrumentation(ID)
		if err != nil {
			tb.Fatalf("failed to create instrumentation: %v", err)
		}
		return inst
	}
	run := func(err error) func(*testing.B) {
		return func(b *testing.B) {
			inst := setup(b)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.ExportLogs(b.Context(), 10).End(err)
			}
		}
	}
	b.Run("NoError", run(nil))
	b.Run("PartialError", run(&internal.PartialSuccess{RejectedItems: 6}))
	b.Run("FullError", run(assert.AnError))
}
