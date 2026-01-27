// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ_test

import (
	"errors"
	"net/http"
	"strconv"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp/internal/observ"
	"go.opentelemetry.io/otel/internal/global"
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
	ServerPort = 4318
)

var Endpoint = ServerAddr + ":" + strconv.Itoa(ServerPort)

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

	_, err := observ.NewInstrumentation(ID, Endpoint)
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")

	assert.ErrorContains(t, err, "inflight metric")
	assert.ErrorContains(t, err, "span exported metric")
	assert.ErrorContains(t, err, "operation duration metric")
}

func TestNewInstrumentationObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY.
	got, err := observ.NewInstrumentation(ID, Endpoint)
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

	inst, err := observ.NewInstrumentation(ID, Endpoint)
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
		semconv.OTelComponentTypeOtlpHTTPSpanExporter,
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

func operationDuration(err error, statusCode int) metricdata.Metrics {
	httpSet := func(err error, statusCode int) attribute.Set {
		attrs := baseAttrs(err)
		attrs = append(attrs, semconv.HTTPResponseStatusCode(statusCode))
		return attribute.NewSet(attrs...)
	}
	return metricdata.Metrics{
		Name:        otelconv.SDKExporterOperationDuration{}.Name(),
		Description: otelconv.SDKExporterOperationDuration{}.Description(),
		Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{Attributes: httpSet(err, statusCode)},
			},
		},
	}
}

func assertMetrics(t *testing.T, got metricdata.ScopeMetrics, spans, success int64, err error, statusCode int) {
	t.Helper()

	assert.Equal(t, Scope, got.Scope, "unexpected scope")

	m := got.Metrics
	require.Len(t, m, 3, "expected 3 metrics")

	o := metricdatatest.IgnoreTimestamp()
	want := spanInflight()
	metricdatatest.AssertEqual(t, want, m[0], o)

	want = spanExported(success, spans, err)
	metricdatatest.AssertEqual(t, want, m[1], o)

	want = operationDuration(err, statusCode)
	metricdatatest.AssertEqual(t, want, m[2], o, metricdatatest.IgnoreValue())
}

func TestInstrumentationExportSpans(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	inst.ExportSpans(t.Context(), n).End(nil, http.StatusOK)

	assertMetrics(t, collect(), n, n, nil, http.StatusOK)
}

func TestInstrumentationExportSpansAllErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	err := errors.New("http error")
	inst.ExportSpans(t.Context(), n).End(err, http.StatusInternalServerError)

	const success = 0
	assertMetrics(t, collect(), n, success, err, http.StatusInternalServerError)
}

func TestInstrumentationExportSpansPartialErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	const success = n - 5

	err := errors.New("partial failure")
	err = errors.Join(err, &internal.PartialSuccess{RejectedItems: 5})
	inst.ExportSpans(t.Context(), n).End(err, http.StatusOK)

	assertMetrics(t, collect(), n, success, err, http.StatusOK)
}

func TestInstrumentationExportSpansInvalidPartialErrored(t *testing.T) {
	inst, collect := setup(t)

	const n = 10
	pErr := &internal.PartialSuccess{RejectedItems: -5}
	err := errors.Join(errors.New("temporary"), pErr)
	inst.ExportSpans(t.Context(), n).End(err, http.StatusServiceUnavailable)

	// Round -5 to 0.
	success := int64(n) // (n - 0)
	assertMetrics(t, collect(), n, success, err, http.StatusServiceUnavailable)

	// Note: the metrics are cumulative, so account for the previous
	// ExportSpans call.
	pErr.RejectedItems = n + 5
	inst.ExportSpans(t.Context(), n).End(err, http.StatusServiceUnavailable)

	// Round n+5 to n.
	success += 0 // success + (n - n)
	assertMetrics(t, collect(), n+n, success, err, http.StatusServiceUnavailable)
}

func TestBaseAttrs(t *testing.T) {
	tests := []struct {
		endpoint string
		host     string
		port     int
	}{
		// Empty.
		{endpoint: "", host: "", port: -1},

		// Only a port.
		{endpoint: ":4318", host: "", port: 4318},

		// Hostname.
		{endpoint: "localhost:4318", host: "localhost", port: 4318},
		{endpoint: "localhost", host: "localhost", port: -1},

		// IPv4 address.
		{endpoint: "127.0.0.1:4318", host: "127.0.0.1", port: 4318},
		{endpoint: "127.0.0.1", host: "127.0.0.1", port: -1},

		// IPv6 address.
		{endpoint: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", host: "2001:db8:85a3::8a2e:370:7334", port: -1},
		{endpoint: "2001:db8:85a3:0:0:8a2e:370:7334", host: "2001:db8:85a3::8a2e:370:7334", port: -1},
		{endpoint: "2001:db8:85a3::8a2e:370:7334", host: "2001:db8:85a3::8a2e:370:7334", port: -1},
		{endpoint: "[2001:db8:85a3::8a2e:370:7334]", host: "2001:db8:85a3::8a2e:370:7334", port: -1},
		{endpoint: "[::1]:9090", host: "::1", port: 9090},

		// Port edge cases.
		{endpoint: "example.com:0", host: "example.com", port: 0},
		{endpoint: "example.com:65535", host: "example.com", port: 65535},

		// Case insensitive.
		{endpoint: "ExAmPlE.COM:8080", host: "ExAmPlE.COM", port: 8080},
	}
	for _, tt := range tests {
		got := observ.BaseAttrs(ID, tt.endpoint)
		want := []attribute.KeyValue{
			semconv.OTelComponentName(observ.ComponentName(ID)),
			semconv.OTelComponentTypeOtlpHTTPSpanExporter,
		}

		if tt.host != "" {
			want = append(want, semconv.ServerAddress(tt.host))
		}
		if tt.port != -1 {
			want = append(want, semconv.ServerPort(tt.port))
		}
		assert.Equal(t, want, got)
	}
}

type logSink struct {
	logr.LogSink

	level         int
	msg           string
	keysAndValues []any
}

func (*logSink) Enabled(int) bool { return true }

func (l *logSink) Info(level int, msg string, keysAndValues ...any) {
	l.level, l.msg, l.keysAndValues = level, msg, keysAndValues
	l.LogSink.Info(level, msg, keysAndValues...)
}

func TestBaseAttrsError(t *testing.T) {
	endpoints := []string{
		"example.com:invalid",   // Non-numeric port.
		"example.com:8080:9090", // Multiple colons in port.
		"example.com:99999",     // Port out of range.
		"example.com:-1",        // Port out of range.
	}
	for _, endpoint := range endpoints {
		l := &logSink{LogSink: testr.New(t).GetSink()}
		t.Cleanup(func(orig logr.Logger) func() {
			global.SetLogger(logr.New(l))
			return func() { global.SetLogger(orig) }
		}(global.GetLogger()))

		// Set the logger as global so BaseAttrs can log the error.
		got := observ.BaseAttrs(ID, endpoint)
		want := []attribute.KeyValue{
			semconv.OTelComponentName(observ.ComponentName(ID)),
			semconv.OTelComponentTypeOtlpHTTPSpanExporter,
		}
		assert.Equal(t, want, got)

		assert.Equal(t, 8, l.level, "expected Debug log level")
		assert.Equal(t, "failed to parse endpoint", l.msg)
	}
}

func BenchmarkInstrumentationExportSpans(b *testing.B) {
	setup := func(b *testing.B) *observ.Instrumentation {
		b.Helper()
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		inst, err := observ.NewInstrumentation(ID, Endpoint)
		if err != nil {
			b.Fatalf("failed to create instrumentation: %v", err)
		}
		return inst
	}

	run := func(err error, statusCode int) func(*testing.B) {
		return func(b *testing.B) {
			inst := setup(b)
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				inst.ExportSpans(b.Context(), 10).End(err, statusCode)
			}
		}
	}

	b.Run("NoError", run(nil, http.StatusOK))
	err := &internal.PartialSuccess{RejectedItems: 6}
	b.Run("PartialError", run(err, http.StatusOK))
	b.Run("FullError", run(assert.AnError, http.StatusInternalServerError))
}
