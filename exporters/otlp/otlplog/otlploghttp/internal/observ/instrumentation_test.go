// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal"
	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
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

	err := internal.PartialSuccess{RejectedItems: success}
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

func TestInstrumentation(t *testing.T) {
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
