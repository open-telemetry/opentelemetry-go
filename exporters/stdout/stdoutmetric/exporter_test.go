// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric_test // import "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric/internal/counter"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

func testEncoderOption() stdoutmetric.Option {
	// Discard export output for testing.
	enc := json.NewEncoder(io.Discard)
	return stdoutmetric.WithEncoder(enc)
}

// failingEncoder always returns an error when Encode is called.
type failingEncoder struct{}

func (failingEncoder) Encode(any) error {
	return errors.New("encoding failed")
}

func testCtxErrHonored(factory func(*testing.T) func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx := t.Context()

		t.Run("DeadlineExceeded", func(t *testing.T) {
			innerCtx, innerCancel := context.WithTimeout(ctx, time.Nanosecond)
			t.Cleanup(innerCancel)
			<-innerCtx.Done()

			f := factory(t)
			assert.ErrorIs(t, f(innerCtx), context.DeadlineExceeded)
		})

		t.Run("Canceled", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory(t)
			assert.ErrorIs(t, f(innerCtx), context.Canceled)
		})

		t.Run("NoError", func(t *testing.T) {
			f := factory(t)
			assert.NoError(t, f(ctx))
		})
	}
}

func testCtxErrIgnored(factory func(*testing.T) func(context.Context) error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		ctx := t.Context()

		t.Run("Canceled Ignored", func(t *testing.T) {
			innerCtx, innerCancel := context.WithCancel(ctx)
			innerCancel()

			f := factory(t)
			assert.NoError(t, f(innerCtx))
		})

		t.Run("NoError", func(t *testing.T) {
			f := factory(t)
			assert.NoError(t, f(ctx))
		})
	}
}

func TestExporterExportHonorsContextErrors(t *testing.T) {
	t.Run("Export", testCtxErrHonored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return func(ctx context.Context) error {
			data := new(metricdata.ResourceMetrics)
			return exp.Export(ctx, data)
		}
	}))
}

func TestExporterForceFlushIgnoresContextErrors(t *testing.T) {
	t.Run("ForceFlush", testCtxErrIgnored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return exp.ForceFlush
	}))
}

func TestExporterShutdownIgnoresContextErrors(t *testing.T) {
	t.Run("Shutdown", testCtxErrIgnored(func(t *testing.T) func(context.Context) error {
		exp, err := stdoutmetric.New(testEncoderOption())
		require.NoError(t, err)
		return exp.Shutdown
	}))
}

func TestShutdownExporterReturnsShutdownErrorOnExport(t *testing.T) {
	var (
		data     = new(metricdata.ResourceMetrics)
		ctx      = t.Context()
		exp, err = stdoutmetric.New(testEncoderOption())
	)
	require.NoError(t, err)
	require.NoError(t, exp.Shutdown(ctx))
	assert.EqualError(t, exp.Export(ctx, data), "exporter shutdown")
}

func deltaSelector(metric.InstrumentKind) metricdata.Temporality {
	return metricdata.DeltaTemporality
}

func TestExportWithOptions(t *testing.T) {
	var (
		data = new(metricdata.ResourceMetrics)
		ctx  = t.Context()
	)

	for _, tt := range []struct {
		name string
		opts []stdoutmetric.Option

		expectedData string
	}{
		{
			name:         "with no options",
			expectedData: "{\"Resource\":null,\"ScopeMetrics\":null}\n",
		},
		{
			name: "with pretty print",
			opts: []stdoutmetric.Option{
				stdoutmetric.WithPrettyPrint(),
			},
			expectedData: "{\n\t\"Resource\": null,\n\t\"ScopeMetrics\": null\n}\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			opts := append(tt.opts, stdoutmetric.WithWriter(&b))

			exp, err := stdoutmetric.New(opts...)
			require.NoError(t, err)
			require.NoError(t, exp.Export(ctx, data))

			assert.Equal(t, tt.expectedData, b.String())
		})
	}
}

func TestTemporalitySelector(t *testing.T) {
	exp, err := stdoutmetric.New(
		testEncoderOption(),
		stdoutmetric.WithTemporalitySelector(deltaSelector),
	)
	require.NoError(t, err)

	var unknownKind metric.InstrumentKind
	assert.Equal(t, metricdata.DeltaTemporality, exp.Temporality(unknownKind))
}

func dropSelector(metric.InstrumentKind) metric.Aggregation {
	return metric.AggregationDrop{}
}

func TestAggregationSelector(t *testing.T) {
	exp, err := stdoutmetric.New(
		testEncoderOption(),
		stdoutmetric.WithAggregationSelector(dropSelector),
	)
	require.NoError(t, err)

	var unknownKind metric.InstrumentKind
	assert.Equal(t, metric.AggregationDrop{}, exp.Aggregation(unknownKind))
}

func TestExporter_Export_SelfObservability(t *testing.T) {
	componentNameAttr := semconv.OTelComponentName("go.opentelemetry.io/otel/exporters/stdout/stdoutmetric.exporter/0")
	componentTypeAttr := semconv.OTelComponentTypeKey.String(
		"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric.exporter",
	)
	wantErr := errors.New("encoding failed")

	tests := []struct {
		name                     string
		ctx                      context.Context
		exporterOpts             []stdoutmetric.Option
		selfObservabilityEnabled bool
		expectedExportedCount    int64
		inflightAttrs            attribute.Set
		attributes               attribute.Set
		wantErr                  error
	}{
		{
			name:                     "Enabled",
			ctx:                      context.Background(),
			exporterOpts:             []stdoutmetric.Option{testEncoderOption()},
			selfObservabilityEnabled: true,
			expectedExportedCount:    19,
			inflightAttrs:            attribute.NewSet(componentNameAttr, componentTypeAttr),
			attributes:               attribute.NewSet(componentNameAttr, componentTypeAttr),
		},
		{
			name:                     "Disabled",
			ctx:                      context.Background(),
			exporterOpts:             []stdoutmetric.Option{testEncoderOption()},
			selfObservabilityEnabled: false,
			expectedExportedCount:    0,
		},
		{
			name:                     "EncodingError",
			ctx:                      context.Background(),
			exporterOpts:             []stdoutmetric.Option{stdoutmetric.WithEncoder(failingEncoder{})},
			selfObservabilityEnabled: true,
			expectedExportedCount:    19,
			inflightAttrs:            attribute.NewSet(componentNameAttr, componentTypeAttr),
			attributes: attribute.NewSet(
				componentNameAttr,
				componentTypeAttr,
				semconv.ErrorType(wantErr),
			),
			wantErr: wantErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", strconv.FormatBool(tt.selfObservabilityEnabled))
			reader := metric.NewManualReader()
			mp := metric.NewMeterProvider(metric.WithReader(reader))
			origMp := otel.GetMeterProvider()
			otel.SetMeterProvider(mp)
			t.Cleanup(func() { otel.SetMeterProvider(origMp) })

			exp, err := stdoutmetric.New(tt.exporterOpts...)
			require.NoError(t, err)
			rm := &metricdata.ResourceMetrics{ScopeMetrics: scopeMetrics()}

			err = exp.Export(tt.ctx, rm)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
			}

			var metrics metricdata.ResourceMetrics
			err = reader.Collect(tt.ctx, &metrics)
			require.NoError(t, err)

			if !tt.selfObservabilityEnabled {
				assert.Empty(t, metrics.ScopeMetrics)
			} else {
				assert.Len(t, metrics.ScopeMetrics, 1)
				expectedMetrics := metricdata.ScopeMetrics{
					Scope: instrumentation.Scope{
						Name:      "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric",
						Version:   sdk.Version(),
						SchemaURL: semconv.SchemaURL,
					},
					Metrics: []metricdata.Metrics{
						{
							Name:        otelconv.SDKExporterMetricDataPointInflight{}.Name(),
							Description: otelconv.SDKExporterMetricDataPointInflight{}.Description(),
							Unit:        otelconv.SDKExporterMetricDataPointInflight{}.Unit(),
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value:      0,
										Attributes: tt.inflightAttrs,
									},
								},
								Temporality: metricdata.CumulativeTemporality,
							},
						},
						{
							Name:        otelconv.SDKExporterMetricDataPointExported{}.Name(),
							Description: otelconv.SDKExporterMetricDataPointExported{}.Description(),
							Unit:        otelconv.SDKExporterMetricDataPointExported{}.Unit(),
							Data: metricdata.Sum[int64]{
								DataPoints: []metricdata.DataPoint[int64]{
									{
										Value:      tt.expectedExportedCount,
										Attributes: tt.attributes,
									},
								},
								Temporality: metricdata.CumulativeTemporality,
								IsMonotonic: true,
							},
						},
						{
							Name:        otelconv.SDKExporterOperationDuration{}.Name(),
							Description: otelconv.SDKExporterOperationDuration{}.Description(),
							Unit:        otelconv.SDKExporterOperationDuration{}.Unit(),
							Data: metricdata.Histogram[float64]{
								DataPoints: []metricdata.HistogramDataPoint[float64]{
									{
										Attributes: tt.attributes,
									},
								},
								Temporality: metricdata.CumulativeTemporality,
							},
						},
					},
				}
				assert.Equal(t, expectedMetrics.Scope, metrics.ScopeMetrics[0].Scope)
				metricdatatest.AssertEqual(t, expectedMetrics.Metrics[0], metrics.ScopeMetrics[0].Metrics[0], metricdatatest.IgnoreTimestamp())
				metricdatatest.AssertEqual(t, expectedMetrics.Metrics[1], metrics.ScopeMetrics[0].Metrics[1], metricdatatest.IgnoreTimestamp())
				metricdatatest.AssertEqual(t, expectedMetrics.Metrics[2], metrics.ScopeMetrics[0].Metrics[2], metricdatatest.IgnoreTimestamp(), metricdatatest.IgnoreValue())
			}
			_ = counter.SetExporterID(0)
		})
	}
}

func scopeMetrics() []metricdata.ScopeMetrics {
	return []metricdata.ScopeMetrics{
		{
			Metrics: []metricdata.Metrics{
				{
					Name: "gauge_int64",
					Data: metricdata.Gauge[int64]{
						DataPoints: []metricdata.DataPoint[int64]{{Value: 1}, {Value: 2}},
					},
				},
				{
					Name: "gauge_float64",
					Data: metricdata.Gauge[float64]{
						DataPoints: []metricdata.DataPoint[float64]{
							{Value: 1.0},
							{Value: 2.0},
							{Value: 3.0},
						},
					},
				},
				{
					Name: "sum_int64",
					Data: metricdata.Sum[int64]{
						DataPoints: []metricdata.DataPoint[int64]{{Value: 10}},
					},
				},
				{
					Name: "sum_float64",
					Data: metricdata.Sum[float64]{
						DataPoints: []metricdata.DataPoint[float64]{{Value: 10.5}, {Value: 20.5}},
					},
				},
				{
					Name: "histogram_int64",
					Data: metricdata.Histogram[int64]{
						DataPoints: []metricdata.HistogramDataPoint[int64]{
							{Count: 1},
							{Count: 2},
							{Count: 3},
						},
					},
				},
				{
					Name: "histogram_float64",
					Data: metricdata.Histogram[float64]{
						DataPoints: []metricdata.HistogramDataPoint[float64]{{Count: 1}},
					},
				},
				{
					Name: "exponential_histogram_int64",
					Data: metricdata.ExponentialHistogram[int64]{
						DataPoints: []metricdata.ExponentialHistogramDataPoint[int64]{
							{Count: 1},
							{Count: 2},
						},
					},
				},
				{
					Name: "exponential_histogram_float64",
					Data: metricdata.ExponentialHistogram[float64]{
						DataPoints: []metricdata.ExponentialHistogramDataPoint[float64]{
							{Count: 1},
							{Count: 2},
							{Count: 3},
							{Count: 4},
						},
					},
				},
				{
					Name: "summary",
					Data: metricdata.Summary{
						DataPoints: []metricdata.SummaryDataPoint{{Count: 1}},
					},
				},
			},
		},
	}
}

func TestExporter_Export_EncodingErrorTracking(t *testing.T) {
	t.Setenv("OTEL_GO_X_SELF_OBSERVABILITY", "true")
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	origMp := otel.GetMeterProvider()
	otel.SetMeterProvider(mp)
	t.Cleanup(func() { otel.SetMeterProvider(origMp) })

	exp, err := stdoutmetric.New(stdoutmetric.WithEncoder(failingEncoder{}))
	assert.NoError(t, err)

	rm := &metricdata.ResourceMetrics{
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Metrics: []metricdata.Metrics{
					{
						Name: "test_gauge",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{{Value: 1}, {Value: 2}},
						},
					},
				},
			},
		},
	}

	ctx := context.Background()
	err = exp.Export(ctx, rm)
	assert.EqualError(t, err, "encoding failed")

	var metrics metricdata.ResourceMetrics
	err = reader.Collect(ctx, &metrics)
	require.NoError(t, err)

	var foundErrorType bool
	for _, sm := range metrics.ScopeMetrics {
		for _, m := range sm.Metrics {
			x := otelconv.SDKExporterMetricDataPointExported{}.Name()
			if m.Name == x {
				if sum, ok := m.Data.(metricdata.Sum[int64]); ok {
					for _, dp := range sum.DataPoints {
						var attr attribute.Value
						attr, foundErrorType = dp.Attributes.Value(semconv.ErrorTypeKey)
						assert.Equal(t, "*errors.errorString", attr.AsString())
					}
				}
			}
		}
	}
	assert.True(t, foundErrorType)
}
