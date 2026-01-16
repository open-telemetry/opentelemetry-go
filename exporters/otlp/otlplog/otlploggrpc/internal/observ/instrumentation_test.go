// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/observ"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

const (
	ID     = 0
	TARGET = "localhost:8080"
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

func (m *errMeter) Int64UpDownCounter(string, ...mapi.Int64UpDownCounterOption) (mapi.Int64UpDownCounter, error) {
	return nil, m.err
}

func (m *errMeter) Int64Counter(string, ...mapi.Int64CounterOption) (mapi.Int64Counter, error) {
	return nil, m.err
}

func (m *errMeter) Float64Histogram(string, ...mapi.Float64HistogramOption) (mapi.Float64Histogram, error) {
	return nil, m.err
}

func TestNewExporterMetrics(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	t.Run("No Error", func(t *testing.T) {
		em, err := NewInstrumentation(ID, "dns:///example.com:42")
		require.NoError(t, err)
		assert.ElementsMatch(t, []attribute.KeyValue{
			semconv.OTelComponentName(GetComponentName(ID)),
			semconv.OTelComponentTypeKey.String(string(otelconv.ComponentTypeOtlpGRPCLogExporter)),
			semconv.ServerAddress("example.com"),
			semconv.ServerPort(42),
		}, em.presetAttrs)

		assert.NotNil(t, em.logInflightMetric, "logInflightMetric should be created")
		assert.NotNil(t, em.logExportedMetric, "logExportedMetric should be created")
		assert.NotNil(t, em.logExportedDurationMetric, "logExportedDurationMetric should be created")
	})

	t.Run("Error", func(t *testing.T) {
		orig := otel.GetMeterProvider()
		t.Cleanup(func() { otel.SetMeterProvider(orig) })
		mp := &errMeterProvider{err: assert.AnError}
		otel.SetMeterProvider(mp)

		_, err := NewInstrumentation(ID, "dns:///:8080")
		require.ErrorIs(t, err, assert.AnError, "new instrument errors")

		assert.ErrorContains(t, err, "inflight metric")
		assert.ErrorContains(t, err, "log exported metric")
		assert.ErrorContains(t, err, "operation duration metric")
	})
}

func TestServerAddrAttrs(t *testing.T) {
	testcases := []struct {
		name   string
		target string
		want   []attribute.KeyValue
	}{
		{
			name:   "Unix socket",
			target: "unix:///tmp/grpc.sock",
			want:   []attribute.KeyValue{semconv.ServerAddress("/tmp/grpc.sock")},
		},
		{
			name:   "DNS with port",
			target: "dns:///localhost:8080",
			want:   []attribute.KeyValue{semconv.ServerAddress("localhost"), semconv.ServerPort(8080)},
		},
		{
			name:   "Dns with endpoint host:port",
			target: "dns://8.8.8.8/example.com:4",
			want:   []attribute.KeyValue{semconv.ServerAddress("example.com"), semconv.ServerPort(4)},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			attrs := ServerAddrAttrs(tc.target)
			assert.Equal(t, tc.want, attrs)
		})
	}
}

func set(err error) attribute.Set {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(GetComponentName(ID)),
		semconv.OTelComponentTypeKey.String(string(otelconv.ComponentTypeOtlpGRPCLogExporter)),
	}
	attrs = append(attrs, ServerAddrAttrs(TARGET)...)
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}
	return attribute.NewSet(attrs...)
}

func logInflightMetrics() metricdata.Metrics {
	m := otelconv.SDKExporterLogInflight{}
	return metricdata.Metrics{
		Name:        m.Name(),
		Description: m.Description(),
		Unit:        m.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: set(nil), Value: 0},
			},
		},
	}
}

func logExportedMetrics(success, total int64, err error) metricdata.Metrics {
	dp := []metricdata.DataPoint[int64]{
		{Attributes: set(nil), Value: success},
	}

	if err != nil {
		dp = append(dp, metricdata.DataPoint[int64]{
			Attributes: set(err),
			Value:      total - success,
		})
	}

	m := otelconv.SDKExporterLogExported{}
	return metricdata.Metrics{
		Name:        m.Name(),
		Description: m.Description(),
		Unit:        m.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dp,
		},
	}
}

func logOperationDurationMetrics(err error, code codes.Code) metricdata.Metrics {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(GetComponentName(ID)),
		semconv.OTelComponentTypeKey.String(string(otelconv.ComponentTypeOtlpGRPCLogExporter)),
		attribute.Int64("rpc.grpc.status_code", int64(code)),
	}
	attrs = append(attrs, ServerAddrAttrs(TARGET)...)
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}

	m := otelconv.SDKExporterOperationDuration{}
	return metricdata.Metrics{
		Name:        m.Name(),
		Description: m.Description(),
		Unit:        m.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{Attributes: attribute.NewSet(attrs...)},
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

	r := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(r))
	otel.SetMeterProvider(mp)

	inst, err := NewInstrumentation(ID, TARGET)
	require.NoError(t, err)
	require.NotNil(t, inst)

	return inst, func() metricdata.ScopeMetrics {
		var rm metricdata.ResourceMetrics
		require.NoError(t, r.Collect(t.Context(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		return rm.ScopeMetrics[0]
	}
}

var Scope = instrumentation.Scope{
	Name:      ScopeName,
	Version:   Version,
	SchemaURL: semconv.SchemaURL,
}

func assertMetrics(
	t *testing.T,
	got metricdata.ScopeMetrics,
	spans int64,
	success int64,
	err error,
	code codes.Code,
) {
	t.Helper()

	assert.Equal(t, Scope, got.Scope, "unexpected scope")

	m := got.Metrics
	require.Len(t, m, 3, "expected 3 metrics")

	o := metricdatatest.IgnoreTimestamp()
	want := logInflightMetrics()
	metricdatatest.AssertEqual(t, want, m[0], o)

	want = logExportedMetrics(success, spans, err)
	metricdatatest.AssertEqual(t, want, m[1], o)

	want = logOperationDurationMetrics(err, code)
	metricdatatest.AssertEqual(t, want, m[2], metricdatatest.IgnoreValue(), o)
}

func TestInstrumentationExportLogs(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	inst.ExportLogs(t.Context(), n).End(nil)
	assertMetrics(t, collect(), n, n, nil, codes.OK)
}

func TestInstrumentationExportLogPartialErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	const success = 5

	err := internal.PartialSuccess{RejectedItems: n - success}
	inst.ExportLogs(t.Context(), n).End(err)

	assertMetrics(t, collect(), n, success, err, status.Code(err))
}

func TestInstrumentationExportLogAllErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	const success = 0
	inst.ExportLogs(t.Context(), n).End(assert.AnError)

	assertMetrics(t, collect(), n, success, assert.AnError, status.Code(assert.AnError))
}

func TestInstrumentationExportLogsInvalidPartialErrored(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	err := internal.PartialSuccess{RejectedItems: -5}
	inst.ExportLogs(t.Context(), n).End(err)

	success := int64(n)
	assertMetrics(t, collect(), n, success, err, status.Code(err))

	err.RejectedItems = n + 5
	inst.ExportLogs(t.Context(), n).End(err)

	success += 0
	assertMetrics(t, collect(), n+n, success, err, status.Code(err))
}

func BenchmarkInstrumentationExportLogs(b *testing.B) {
	setup := func(tb *testing.B) *Instrumentation {
		tb.Helper()
		tb.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		inst, err := NewInstrumentation(ID, TARGET)
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

func BenchmarkSetPresetAttrs(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := range b.N {
		getPresetAttrs(int64(i), "dns:///192.168.1.1:8080")
	}
}
