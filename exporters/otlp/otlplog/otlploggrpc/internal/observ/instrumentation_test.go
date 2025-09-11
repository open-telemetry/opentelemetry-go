// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/observ"

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	mapi "go.opentelemetry.io/otel/metric"
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
	t.Run("No Error", func(t *testing.T) {
		em, err := NewInstrumentation(ID, "localhost:8080")
		require.NoError(t, err)
		assert.ElementsMatch(t, []attribute.KeyValue{
			semconv.OTelComponentName(GetComponentName(ID)),
			semconv.OTelComponentTypeKey.String(string(otelconv.ComponentTypeOtlpGRPCLogExporter)),
			semconv.ServerAddress("localhost"),
			semconv.ServerPort(8080),
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

		_, err := NewInstrumentation(ID, "localhost:8080")
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
		{
			name:   "Simple host port",
			target: "localhost:10001",
			want:   []attribute.KeyValue{semconv.ServerAddress("localhost"), semconv.ServerPort(10001)},
		},
		{
			name:   "Host without port",
			target: "example.com",
			want:   []attribute.KeyValue{semconv.ServerAddress("example.com")},
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
		semconv.RPCGRPCStatusCodeKey.Int64(int64(code)),
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
		require.NoError(t, r.Collect(context.Background(), &rm))
		require.Len(t, rm.ScopeMetrics, 1)
		return rm.ScopeMetrics[0]
	}
}

var Scope = instrumentation.Scope{
	Name:      ScopeName,
	Version:   sdk.Version(),
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
	end := inst.ExportLogs(context.Background(), n)
	end(nil, n, codes.OK)
	assertMetrics(t, collect(), n, n, nil, codes.OK)
}

func TestInstrumentationExportLogPartialErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	end := inst.ExportLogs(context.Background(), n)
	const success = 5
	end(assert.AnError, success, codes.Canceled)

	assertMetrics(t, collect(), n, success, assert.AnError, codes.Canceled)
}

func TestInstrumentationExportLogAllErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	end := inst.ExportLogs(context.Background(), n)
	const success = 0
	end(assert.AnError, success, codes.Canceled)

	assertMetrics(t, collect(), n, success, assert.AnError, codes.Canceled)
}

func BenchmarkInstrumentationExportLogs(b *testing.B) {
	b.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
	inst, err := NewInstrumentation(ID, TARGET)
	if err != nil {
		b.Fatalf("failed to create instrumentation: %v", err)
	}

	var end ExportLogDone
	err = errors.New("benchmark error")

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		end = inst.ExportLogs(context.Background(), 10)
		end(err, 4, codes.Canceled)
	}
	_ = end
}
