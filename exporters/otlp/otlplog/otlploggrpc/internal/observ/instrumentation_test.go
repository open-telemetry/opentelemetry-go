// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc/internal/observ"

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	mapi "go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
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
		em, err := NewInstrumentation(
			"newExportMetricsTest",
			"newExportMetricsTest/1",
			"newExportMetricsTest",
			"localhost:8080",
		)
		require.NoError(t, err)
		assert.ElementsMatch(t, []attribute.KeyValue{
			semconv.OTelComponentName("newExportMetricsTest/1"),
			semconv.OTelComponentTypeKey.String("newExportMetricsTest"),
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

		_, err := NewInstrumentation(
			"newExportMetrics",
			"newExportMetrics/1",
			"newExportMetrics",
			"localhost:8080",
		)
		require.ErrorIs(t, err, assert.AnError, "new instrument errors")

		assert.ErrorContains(t, err, "inflight metric")
		assert.ErrorContains(t, err, "span exported metric")
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
