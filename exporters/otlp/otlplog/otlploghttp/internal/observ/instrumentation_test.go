// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal"
	"go.opentelemetry.io/otel/internal/global"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
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

func TestNewInstrumentationObservabilityErrors(t *testing.T) {
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })
	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	_, err := NewInstrumentation(ID, TARGET)
	require.ErrorIs(t, err, assert.AnError, "new instrument errors")

	assert.ErrorContains(t, err, "inflight metric")
	assert.ErrorContains(t, err, "exported metric")
	assert.ErrorContains(t, err, "operation duration metric")
}

func TestNewInstrumentationObservabilityDisabled(t *testing.T) {
	// Do not set OTEL_GO_X_OBSERVABILITY.
	got, err := NewInstrumentation(ID, TARGET)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func set(err error) attribute.Set {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(GetComponentName(ID)),
		semconv.OTelComponentTypeKey.String(string(otelconv.ComponentTypeOtlpHTTPLogExporter)),
	}
	attrs = append(attrs, ServerAddrAttrs(TARGET)...)
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}
	return attribute.NewSet(attrs...)
}

func inflightMetric() metricdata.Metrics {
	inflight := otelconv.SDKExporterLogInflight{}

	return metricdata.Metrics{
		Name:        inflight.Name(),
		Description: inflight.Description(),
		Unit:        inflight.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.DataPoint[int64]{
				{
					Attributes: set(nil),
					Value:      0,
				},
			},
		},
	}
}

func exportedMetric(err error, total, success int64) metricdata.Metrics {
	dp := []metricdata.DataPoint[int64]{
		{
			Attributes: set(nil),
			Value:      success,
		},
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

func operationDurationMetric(err error, code int) metricdata.Metrics {
	attrs := []attribute.KeyValue{
		semconv.OTelComponentName(GetComponentName(ID)),
		semconv.OTelComponentTypeOtlpHTTPLogExporter,
		semconv.HTTPResponseStatusCode(code),
	}
	attrs = append(attrs, ServerAddrAttrs(TARGET)...)
	if err != nil {
		attrs = append(attrs, semconv.ErrorType(err))
	}

	operation := otelconv.SDKExporterOperationDuration{}

	return metricdata.Metrics{
		Name:        operation.Name(),
		Description: operation.Description(),
		Unit:        operation.Unit(),
		Data: metricdata.Histogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{
					Attributes: attribute.NewSet(attrs...),
				},
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
	count int64,
	success int64,
	err error,
	code int,
) {
	t.Helper()
	assert.Equal(t, Scope, got.Scope, "unexpected scope")
	m := got.Metrics
	require.Len(t, m, 3, "expected 3 metrics")

	o := metricdatatest.IgnoreTimestamp()
	want := inflightMetric()
	metricdatatest.AssertEqual(t, want, m[0], o)

	want = exportedMetric(err, count, success)
	metricdatatest.AssertEqual(t, want, m[1], o)

	want = operationDurationMetric(err, code)
	metricdatatest.AssertEqual(t, want, m[2], metricdatatest.IgnoreValue(), o)
}

func TestInstrumentationExportedLogs(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	inst.ExportLogs(t.Context(), n).End(nil, http.StatusOK)
	assertMetrics(t, collect(), n, n, nil, http.StatusOK)
}

func TestInstrumentationExportLogsPartialErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	const success = 5

	err := internal.PartialSuccess{RejectedItems: n - success}
	inst.ExportLogs(t.Context(), n).End(err, http.StatusPartialContent)

	assertMetrics(t, collect(), n, success, err, http.StatusPartialContent)
}

func TestInstrumentationExportLogAllErrors(t *testing.T) {
	inst, collect := setup(t)
	const n = 10
	const success = 0

	inst.ExportLogs(t.Context(), n).End(assert.AnError, http.StatusUnauthorized)

	assertMetrics(t, collect(), n, success, assert.AnError, http.StatusUnauthorized)
}

func TestInstrumentationExportLogsInvalidPartialErrored(t *testing.T) {
	inst, collect := setup(t)
	const n = 10

	err := internal.PartialSuccess{RejectedItems: -5}
	inst.ExportLogs(t.Context(), n).End(err, http.StatusPartialContent)

	success := n
	assertMetrics(t, collect(), n, int64(success), err, http.StatusPartialContent)

	err.RejectedItems = n + 5
	inst.ExportLogs(t.Context(), n).End(err, http.StatusPartialContent)

	success += 0
	assertMetrics(t, collect(), n+n, int64(success), err, http.StatusPartialContent)
}

type spy struct {
	enabled  bool
	called   *bool
	panicMsg string
}

func (s spy) Enabled(context.Context) bool { return s.enabled }

func (s spy) markCalled() {
	if !s.enabled {
		panic(s.panicMsg)
	}
	if s.called != nil {
		*s.called = true
	}
}

type upDownCounterSpy struct {
	embedded.Int64UpDownCounter
	spy
}

func (c upDownCounterSpy) Add(context.Context, int64, ...mapi.AddOption) { c.markCalled() }

type counterSpy struct {
	embedded.Int64Counter
	spy
}

func (c counterSpy) Add(context.Context, int64, ...mapi.AddOption) { c.markCalled() }

type histogramSpy struct {
	embedded.Float64Histogram
	spy
}

func (h histogramSpy) Record(context.Context, float64, ...mapi.RecordOption) { h.markCalled() }

func TestEndSkipsDisabledInstruments(t *testing.T) {
	const n = 10

	tests := []struct {
		name         string
		disable      string // "inflight" | "exported" | "duration"
		wantInflight bool
		wantExported bool
		wantDuration bool
	}{
		{name: "inflight disabled", disable: "inflight", wantExported: true, wantDuration: true},
		{name: "exported disabled", disable: "exported", wantInflight: true, wantDuration: true},
		{name: "duration disabled", disable: "duration", wantInflight: true, wantExported: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var inflightCalled, exportedCalled, durationCalled bool

			inst := &Instrumentation{
				inflightMetric: upDownCounterSpy{
					spy: spy{
						enabled:  tt.disable != "inflight",
						called:   &inflightCalled,
						panicMsg: "inflight Add called while disabled",
					},
				},
				exportedMetric: counterSpy{
					spy: spy{
						enabled:  tt.disable != "exported",
						called:   &exportedCalled,
						panicMsg: "exported Add called while disabled",
					},
				},
				operationDuration: histogramSpy{
					spy: spy{
						enabled:  tt.disable != "duration",
						called:   &durationCalled,
						panicMsg: "duration Record called while disabled",
					},
				},
			}

			// If your Enabled() guards are missing, this will panic in the disabled case.
			inst.ExportLogs(t.Context(), n).End(nil, http.StatusOK)

			require.Equal(t, tt.wantInflight, inflightCalled)
			require.Equal(t, tt.wantExported, exportedCalled)
			require.Equal(t, tt.wantDuration, durationCalled)
		})
	}
}

func TestSetPresetAttrs(t *testing.T) {
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
		got := setPresetAttrs(GetComponentName(ID), tt.endpoint)
		want := []attribute.KeyValue{
			semconv.OTelComponentName(GetComponentName(ID)),
			semconv.OTelComponentTypeOtlpHTTPLogExporter,
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

func TestSetPresetAttrsError(t *testing.T) {
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
		got := setPresetAttrs(GetComponentName(ID), endpoint)
		want := []attribute.KeyValue{
			semconv.OTelComponentName(GetComponentName(ID)),
			semconv.OTelComponentTypeOtlpHTTPLogExporter,
		}
		assert.Equal(t, want, got)

		assert.Equal(t, 8, l.level, "expected Debug log level")
		assert.Equal(t, "failed to parse target", l.msg)
	}
}

func BenchmarkInstrumentationExportLogs(b *testing.B) {
	setup := func(b *testing.B) *Instrumentation {
		b.Helper()
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		inst, err := NewInstrumentation(ID, TARGET)
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
				inst.ExportLogs(b.Context(), 10).End(err, statusCode)
			}
		}
	}

	b.Run("NoError", run(nil, http.StatusOK))
	err := &internal.PartialSuccess{RejectedItems: 6}
	b.Run("PartialError", run(err, http.StatusOK))
	b.Run("FullError", run(assert.AnError, http.StatusInternalServerError))
}
