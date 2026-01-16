// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc/internal/observ"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

const (
	ID         = 0
	ServerAddr = "localhost"
	ServerPort = 4317
)

var Target = "dns://" + ServerAddr + ":" + strconv.Itoa(ServerPort)

var Scope = instrumentation.Scope{
	Name:      observ.ScopeName,
	Version:   observ.Version,
	SchemaURL: observ.SchemaURL,
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

func (m *errMeter) Int64UpDownCounter(string, ...mapi.Int64UpDownCounterOption) (mapi.Int64UpDownCounter, error) {
	return nil, m.err
}

func (m *errMeter) Int64Counter(string, ...mapi.Int64CounterOption) (mapi.Int64Counter, error) {
	return nil, m.err
}

func (m *errMeter) Float64Histogram(string, ...mapi.Float64HistogramOption) (mapi.Float64Histogram, error) {
	return nil, m.err
}

func TestNewInstrumentationObservabilityErrors(t *testing.T) {
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })
	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	_, err := observ.NewInstrumentation(ID, Target)
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")

	assert.ErrorContains(t, err, "inflight metric")
	assert.ErrorContains(t, err, "span exported metric")
	assert.ErrorContains(t, err, "operation duration metric")
}

func TestNewInstrumentationObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY.
	got, err := observ.NewInstrumentation(ID, Target)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func setup(t *testing.T) (*observ.Instrumentation, func() metricdata.ScopeMetrics) {
	t.Helper()

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	original := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(original) })

	r := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(r))
	otel.SetMeterProvider(mp)

	inst, err := observ.NewInstrumentation(ID, Target)
	require.NoError(t, err)
	require.NotNil(t, inst)

	return inst, func() metricdata.ScopeMetrics {
		var rm metricdata.ResourceMetrics
		require.NoError(t, r.Collect(t.Context(), &rm))

		require.Len(t, rm.ScopeMetrics, 1)
		return rm.ScopeMetrics[0]
	}
}

func baseAttrs(err error) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(observ.ComponentName(ID)),
		semconv.OTelComponentTypeOtlpGRPCSpanExporter,
		semconv.ServerAddress(ServerAddr),
		semconv.ServerPort(ServerPort),
	}
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}
	return attrs
}

func set(err error) attribute.Set {
	return attribute.NewSet(baseAttrs(err)...)
}

func spanInflight() metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterSpanInflight{}.Name(),
		Description: otelconv.SDKExporterSpanInflight{}.Description(),
		Unit:        otelconv.SDKExporterSpanInflight{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[int64]{
				{Attributes: set(nil), Value: 0},
			},
		},
	}
}

func spanExported(success, total int64, err error) metricdata.Metrics {
	dp := []metricdata.DataPoint[int64]{
		{Attributes: set(nil), Value: success},
	}
	if err != nil {
		dp = append(dp, metricdata.DataPoint[int64]{
			Attributes: set(err),
			Value:      total - success,
		})
	}
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterSpanExported{}.Name(),
		Description: otelconv.SDKExporterSpanExported{}.Description(),
		Unit:        otelconv.SDKExporterSpanExported{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dp,
		},
	}
}

func operationDuration(err error) metricdata.Metrics {
	rpcSet := func(err error) attribute.Set {
		c := int64(status.Code(err))
		return attribute.NewSet(append(
			[]attribute.KeyValue{
				attribute.Int64("rpc.grpc.status_code", c),
			},
			baseAttrs(err)...,
		)...)
	}
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterOperationDuration{}.Name(),
		Description: otelconv.SDKExporterOperationDuration{}.Description(),
		Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{Attributes: rpcSet(err)},
			},
		},
	}
}

func assertMetrics(t *testing.T, got metricdata.ScopeMetrics, spans, success int64, err error) {
	t.Helper()

	assert.Equal(t, Scope, got.Scope, "unexpected scope")

	m := got.Metrics
	require.Len(t, m, 3, "expected 3 metrics")

	o := metricdatatest.IgnoreTimestamp()
	want := spanInflight()
	metricdatatest.AssertEqual(t, want, m[0], o)

	want = spanExported(success, spans, err)
	metricdatatest.AssertEqual(t, want, m[1], o)

	want = operationDuration(err)
	metricdatatest.AssertEqual(t, want, m[2], o, metricdatatest.IgnoreValue())
}

func TestInstrumentationExportSpans(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	inst.ExportSpans(t.Context(), n).End(nil, codes.OK)

	assertMetrics(t, collect(), n, n, nil)
}

func TestInstrumentationExportSpansAllErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	c := codes.PermissionDenied
	err := status.Error(c, "go away")
	inst.ExportSpans(t.Context(), n).End(err, c)

	const success = 0
	assertMetrics(t, collect(), n, success, err)
}

func TestInstrumentationExportSpansPartialErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	const success = n - 5

	c := codes.Unavailable
	err := status.Error(c, "temporary failure")
	err = errors.Join(err, &internal.PartialSuccess{RejectedItems: 5})
	inst.ExportSpans(t.Context(), n).End(err, c)

	assertMetrics(t, collect(), n, success, err)
}

func TestInstrumentationExportSpansInvalidPartialErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	pErr := &internal.PartialSuccess{RejectedItems: -5}
	c := codes.Unavailable
	err := errors.Join(status.Error(c, "temporary"), pErr)
	inst.ExportSpans(t.Context(), n).End(err, c)

	// Round -5 to 0.
	success := int64(n) // (n - 0)
	assertMetrics(t, collect(), n, success, err)

	// Note: the metrics are cumulative, so account for the previous
	// ExportSpans call.
	pErr.RejectedItems = n + 5
	inst.ExportSpans(t.Context(), n).End(err, c)

	// Round n+5 to n.
	success += 0 // success + (n - n)
	assertMetrics(t, collect(), n+n, success, err)
}

func TestBaseAttrs(t *testing.T) {
	tests := []struct {
		name   string
		target string
		want   []attribute.KeyValue
	}{
		{
			name:   "HostAndPort",
			target: "dns://localhost:4317",
			want: []attribute.KeyValue{
				semconv.OTelComponentName(observ.ComponentName(ID)),
				semconv.OTelComponentTypeOtlpGRPCSpanExporter,
				semconv.ServerAddress("localhost"),
				semconv.ServerPort(4317),
			},
		},
		{
			name:   "Host",
			target: "dns://localhost",
			want: []attribute.KeyValue{
				semconv.OTelComponentName(observ.ComponentName(ID)),
				semconv.OTelComponentTypeOtlpGRPCSpanExporter,
				semconv.ServerAddress("localhost"),
			},
		},
		{
			name:   "Port",
			target: "dns://:4317",
			want: []attribute.KeyValue{
				semconv.OTelComponentName(observ.ComponentName(ID)),
				semconv.OTelComponentTypeOtlpGRPCSpanExporter,
				semconv.ServerPort(4317),
			},
		},
		{
			name:   "Empty",
			target: "",
			want: []attribute.KeyValue{
				semconv.OTelComponentName(observ.ComponentName(ID)),
				semconv.OTelComponentTypeOtlpGRPCSpanExporter,
			},
		},
		{
			name:   "Invalid",
			target: "dns:///:invalid",
			want: []attribute.KeyValue{
				semconv.OTelComponentName(observ.ComponentName(ID)),
				semconv.OTelComponentTypeOtlpGRPCSpanExporter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := observ.BaseAttrs(ID, tt.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkInstrumentationExportSpans(b *testing.B) {
	setup := func(b *testing.B) *observ.Instrumentation {
		b.Helper()
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		inst, err := observ.NewInstrumentation(ID, Target)
		if err != nil {
			b.Fatalf("failed to create instrumentation: %v", err)
		}
		return inst
	}

	run := func(err error, c codes.Code) func(*testing.B) {
		return func(b *testing.B) {
			inst := setup(b)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.ExportSpans(b.Context(), 10).End(err, c)
			}
		}
	}

	b.Run("NoError", run(nil, codes.OK))
	err := &internal.PartialSuccess{RejectedItems: 6}
	b.Run("PartialError", run(err, codes.Unavailable))
	b.Run("FullError", run(assert.AnError, codes.Aborted))
}
