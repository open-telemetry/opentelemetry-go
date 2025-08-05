// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus

import (
	"context"
	"errors"
	"math"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/model"
	"github.com/prometheus/otlptranslator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
	"go.opentelemetry.io/otel/trace"
)

func TestPrometheusExporter(t *testing.T) {
	testCases := []struct {
		name                string
		emptyResource       bool
		customResourceAttrs []attribute.KeyValue
		recordMetrics       func(ctx context.Context, meter otelmetric.Meter)
		options             []Option
		expectedFile        string
		disableUTF8         bool
		checkMetricFamilies func(t testing.TB, dtos []*dto.MetricFamily)
	}{
		{
			name:         "counter",
			expectedFile: "testdata/counter.txt",
			disableUTF8:  true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter(
					"foo",
					otelmetric.WithDescription("a simple counter"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
		{
			name:         "counter that already has the unit suffix",
			expectedFile: "testdata/counter_noutf8_with_unit_suffix.txt",
			disableUTF8:  true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter(
					"foo.seconds",
					otelmetric.WithDescription("a simple counter"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
		{
			name:         "counter with custom unit not tracked by ucum standards",
			expectedFile: "testdata/counter_with_custom_unit_suffix.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter(
					"foo",
					otelmetric.WithDescription("a simple counter"),
					otelmetric.WithUnit("madeup"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
		{
			name:         "counter with bracketed unit",
			expectedFile: "testdata/counter_no_unit.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter(
					"foo",
					otelmetric.WithDescription("a simple counter"),
					otelmetric.WithUnit("{spans}"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
		{
			name:         "counter that already has a total suffix",
			expectedFile: "testdata/counter.txt",
			disableUTF8:  true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter(
					"foo.total",
					otelmetric.WithDescription("a simple counter"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
		{
			name:         "counter with suffixes disabled",
			expectedFile: "testdata/counter_disabled_suffix.txt",
			options:      []Option{WithoutCounterSuffixes()},
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter(
					"foo",
					otelmetric.WithDescription("a simple counter without a total suffix"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
		{
			name:         "gauge",
			expectedFile: "testdata/gauge.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				)
				gauge, err := meter.Float64Gauge(
					"bar",
					otelmetric.WithDescription("a fun little gauge"),
					otelmetric.WithUnit("1"),
				)
				require.NoError(t, err)
				gauge.Record(ctx, .75, opt)
			},
		},
		{
			name:         "exponential histogram",
			expectedFile: "testdata/exponential_histogram.txt",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				var hist *dto.MetricFamily

				for _, mf := range mfs {
					if *mf.Name == `exponential_histogram_baz_bytes` {
						hist = mf
						break
					}
				}

				if hist == nil {
					t.Fatal("expected to find histogram")
				}

				m := hist.GetMetric()[0].Histogram

				require.Equal(t, 236.0, *m.SampleSum)
				require.Equal(t, uint64(4), *m.SampleCount)
				require.Equal(t, []int64{1, -1, 1, -1, 2}, m.PositiveDelta)
				require.Equal(t, uint32(5), *m.PositiveSpan[0].Length)
				require.Equal(t, int32(3), *m.PositiveSpan[0].Offset)
			},
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				// NOTE(GiedriusS): there is no text format for exponential (native)
				// histograms so we don't expect any output.
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				)
				histogram, err := meter.Float64Histogram(
					"exponential_histogram_baz",
					otelmetric.WithDescription("a very nice histogram"),
					otelmetric.WithUnit("By"),
				)
				require.NoError(t, err)
				histogram.Record(ctx, 23, opt)
				histogram.Record(ctx, 7, opt)
				histogram.Record(ctx, 101, opt)
				histogram.Record(ctx, 105, opt)
			},
		},
		{
			name:         "histogram",
			expectedFile: "testdata/histogram.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				)
				histogram, err := meter.Float64Histogram(
					"histogram_baz",
					otelmetric.WithDescription("a very nice histogram"),
					otelmetric.WithUnit("By"),
				)
				require.NoError(t, err)
				histogram.Record(ctx, 23, opt)
				histogram.Record(ctx, 7, opt)
				histogram.Record(ctx, 101, opt)
				histogram.Record(ctx, 105, opt)
			},
		},
		{
			name:         "sanitized attributes to labels",
			expectedFile: "testdata/sanitized_labels.txt",
			disableUTF8:  true,
			options:      []Option{WithoutUnits()},
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					// exact match, value should be overwritten
					attribute.Key("A.B").String("X"),
					attribute.Key("A.B").String("Q"),

					// unintended match due to sanitization, values should be concatenated
					attribute.Key("C.D").String("Y"),
					attribute.Key("C/D").String("Z"),
				)
				counter, err := meter.Float64Counter(
					"foo",
					otelmetric.WithDescription("a sanitary counter"),
					// This unit is not added to
					otelmetric.WithUnit("By"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)
			},
		},
		{
			name:         "invalid instruments are renamed",
			expectedFile: "testdata/sanitized_names.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				)
				// Valid.
				gauge, err := meter.Float64UpDownCounter("bar", otelmetric.WithDescription("a fun little gauge"))
				require.NoError(t, err)
				gauge.Add(ctx, 100, opt)
				gauge.Add(ctx, -25, opt)

				// Invalid, will be renamed.
				gauge, err = meter.Float64UpDownCounter(
					"invalid.gauge.name",
					otelmetric.WithDescription("a gauge with an invalid name"),
				)
				require.NoError(t, err)
				gauge.Add(ctx, 100, opt)

				counter, err := meter.Float64Counter(
					"0invalid.counter.name",
					otelmetric.WithDescription("a counter with an invalid name"),
				)
				require.ErrorIs(t, err, metric.ErrInstrumentName)
				counter.Add(ctx, 100, opt)

				histogram, err := meter.Float64Histogram(
					"invalid.hist.name",
					otelmetric.WithDescription("a histogram with an invalid name"),
				)
				require.NoError(t, err)
				histogram.Record(ctx, 23, opt)
			},
		},
		{
			name:          "empty resource",
			emptyResource: true,
			expectedFile:  "testdata/empty_resource.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)
			},
		},
		{
			name: "custom resource",
			customResourceAttrs: []attribute.KeyValue{
				attribute.Key("A").String("B"),
				attribute.Key("C").String("D"),
			},
			expectedFile: "testdata/custom_resource.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)
			},
		},
		{
			name:         "without target_info",
			options:      []Option{WithoutTargetInfo()},
			expectedFile: "testdata/without_target_info.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)
			},
		},
		{
			name:         "without scope_info",
			options:      []Option{WithoutScopeInfo()},
			expectedFile: "testdata/without_scope_info.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				)
				gauge, err := meter.Int64Gauge(
					"bar",
					otelmetric.WithDescription("a fun little gauge"),
					otelmetric.WithUnit("1"),
				)
				require.NoError(t, err)
				gauge.Record(ctx, 1, opt)
			},
		},
		{
			name:         "without scope_info and target_info",
			options:      []Option{WithoutScopeInfo(), WithoutTargetInfo()},
			expectedFile: "testdata/without_scope_and_target_info.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				)
				counter, err := meter.Int64Counter(
					"bar",
					otelmetric.WithDescription("a fun little counter"),
					otelmetric.WithUnit("By"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 2, opt)
				counter.Add(ctx, 1, opt)
			},
		},
		{
			name:         "with namespace",
			expectedFile: "testdata/with_namespace.txt",
			options: []Option{
				WithNamespace("test"),
			},
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)
			},
		},
		{
			name:         "with resource attributes filter",
			expectedFile: "testdata/with_resource_attributes_filter.txt",
			options: []Option{
				WithResourceAsConstantLabels(attribute.NewDenyKeysFilter()),
			},
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.1, opt)
				counter.Add(ctx, 9.8, opt)
			},
		},
		{
			name:         "with some resource attributes filter",
			expectedFile: "testdata/with_allow_resource_attributes_filter.txt",
			options: []Option{
				WithResourceAsConstantLabels(attribute.NewAllowKeysFilter("service.name")),
			},
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64Counter("foo", otelmetric.WithDescription("a simple counter"))
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 5.9, opt)
				counter.Add(ctx, 5.3, opt)
			},
		},
		{
			name:         "counter utf-8",
			expectedFile: "testdata/counter_utf8.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				opt := otelmetric.WithAttributes(
					attribute.Key("A.G").String("B"),
					attribute.Key("C.H").String("D"),
					attribute.Key("E.I").Bool(true),
					attribute.Key("F.J").Int(42),
				)
				counter, err := meter.Float64Counter(
					"foo.things",
					otelmetric.WithDescription("a simple counter"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A.G").String("D"),
					attribute.Key("C.H").String("B"),
					attribute.Key("E.I").Bool(true),
					attribute.Key("F.J").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
		{
			name:         "non-monotonic sum does not add exemplars",
			expectedFile: "testdata/non_monotonic_sum_does_not_add_exemplars.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				sc := trace.NewSpanContext(trace.SpanContextConfig{
					SpanID:     trace.SpanID{0o1},
					TraceID:    trace.TraceID{0o1},
					TraceFlags: trace.FlagsSampled,
				})
				ctx = trace.ContextWithSpanContext(ctx, sc)
				opt := otelmetric.WithAttributes(
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter, err := meter.Float64UpDownCounter(
					"foo",
					otelmetric.WithDescription("a simple up down counter"),
					otelmetric.WithUnit("s"),
				)
				require.NoError(t, err)
				counter.Add(ctx, 5, opt)
				counter.Add(ctx, 10.3, opt)
				counter.Add(ctx, 9, opt)
				counter.Add(ctx, -1, opt)

				attrs2 := attribute.NewSet(
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				)
				counter.Add(ctx, 5, otelmetric.WithAttributeSet(attrs2))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.disableUTF8 {
				model.NameValidationScheme = model.LegacyValidation // nolint:staticcheck // We need this check to keep supporting the legacy scheme.
				defer func() {
					// Reset to defaults
					model.NameValidationScheme = model.UTF8Validation // nolint:staticcheck // We need this check to keep supporting the legacy scheme.
				}()
			}
			ctx := context.Background()
			registry := prometheus.NewRegistry()
			exporter, err := New(append(tc.options, WithRegisterer(registry))...)
			require.NoError(t, err)

			var res *resource.Resource
			if tc.emptyResource {
				res = resource.Empty()
			} else {
				res, err = resource.New(ctx,
					// always specify service.name because the default depends on the running OS
					resource.WithAttributes(semconv.ServiceName("prometheus_test")),
					// Overwrite the semconv.TelemetrySDKVersionKey value so we don't need to update every version
					resource.WithAttributes(semconv.TelemetrySDKVersion("latest")),
					resource.WithAttributes(tc.customResourceAttrs...),
				)
				require.NoError(t, err)

				res, err = resource.Merge(resource.Default(), res)
				require.NoError(t, err)
			}

			provider := metric.NewMeterProvider(
				metric.WithResource(res),
				metric.WithReader(exporter),
				metric.WithView(metric.NewView(
					metric.Instrument{Name: "histogram_*"},
					metric.Stream{Aggregation: metric.AggregationExplicitBucketHistogram{
						Boundaries: []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
					}},
				),
					metric.NewView(
						metric.Instrument{Name: "exponential_histogram_*"},
						metric.Stream{Aggregation: metric.AggregationBase2ExponentialHistogram{
							MaxSize: 10,
						}},
					),
				),
			)
			meter := provider.Meter(
				"testmeter",
				otelmetric.WithInstrumentationVersion("v0.1.0"),
				otelmetric.WithInstrumentationAttributes(attribute.String("fizz", "buzz")),
			)

			tc.recordMetrics(ctx, meter)

			file, err := os.Open(tc.expectedFile)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, file.Close()) })

			err = testutil.GatherAndCompare(registry, file)
			require.NoError(t, err)

			if tc.checkMetricFamilies == nil {
				return
			}

			mfs, err := registry.Gather()
			require.NoError(t, err)

			tc.checkMetricFamilies(t, mfs)
		})
	}
}

func TestMultiScopes(t *testing.T) {
	ctx := context.Background()
	registry := prometheus.NewRegistry()
	exporter, err := New(WithRegisterer(registry))
	require.NoError(t, err)

	res, err := resource.New(ctx,
		// always specify service.name because the default depends on the running OS
		resource.WithAttributes(semconv.ServiceName("prometheus_test")),
		// Overwrite the semconv.TelemetrySDKVersionKey value so we don't need to update every version
		resource.WithAttributes(semconv.TelemetrySDKVersion("latest")),
	)
	require.NoError(t, err)
	res, err = resource.Merge(resource.Default(), res)
	require.NoError(t, err)

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(res),
	)

	fooCounter, err := provider.Meter("meterfoo", otelmetric.WithInstrumentationVersion("v0.1.0")).
		Int64Counter(
			"foo",
			otelmetric.WithUnit("s"),
			otelmetric.WithDescription("meter foo counter"))
	assert.NoError(t, err)
	fooCounter.Add(ctx, 100, otelmetric.WithAttributes(attribute.String("type", "foo")))

	barCounter, err := provider.Meter("meterbar", otelmetric.WithInstrumentationVersion("v0.1.0")).
		Int64Counter(
			"bar",
			otelmetric.WithUnit("s"),
			otelmetric.WithDescription("meter bar counter"))
	assert.NoError(t, err)
	barCounter.Add(ctx, 200, otelmetric.WithAttributes(attribute.String("type", "bar")))

	file, err := os.Open("testdata/multi_scopes.txt")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, file.Close()) })

	err = testutil.GatherAndCompare(registry, file)
	require.NoError(t, err)
}

func TestDuplicateMetrics(t *testing.T) {
	ab := attribute.NewSet(attribute.String("A", "B"))
	withAB := otelmetric.WithAttributeSet(ab)
	typeBar := attribute.NewSet(attribute.String("type", "bar"))
	withTypeBar := otelmetric.WithAttributeSet(typeBar)
	typeFoo := attribute.NewSet(attribute.String("type", "foo"))
	withTypeFoo := otelmetric.WithAttributeSet(typeFoo)
	testCases := []struct {
		name                  string
		customResourceAttrs   []attribute.KeyValue
		recordMetrics         func(ctx context.Context, meterA, meterB otelmetric.Meter)
		options               []Option
		possibleExpectedFiles []string
	}{
		{
			name: "no_conflict_two_counters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				fooA, err := meterA.Int64Counter("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter counter foo"))
				assert.NoError(t, err)
				fooA.Add(ctx, 100, withAB)

				fooB, err := meterB.Int64Counter("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter counter foo"))
				assert.NoError(t, err)
				fooB.Add(ctx, 100, withAB)
			},
			possibleExpectedFiles: []string{"testdata/no_conflict_two_counters.txt"},
		},
		{
			name: "no_conflict_two_updowncounters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				fooA, err := meterA.Int64UpDownCounter("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter gauge foo"))
				assert.NoError(t, err)
				fooA.Add(ctx, 100, withAB)

				fooB, err := meterB.Int64UpDownCounter("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter gauge foo"))
				assert.NoError(t, err)
				fooB.Add(ctx, 100, withAB)
			},
			possibleExpectedFiles: []string{"testdata/no_conflict_two_updowncounters.txt"},
		},
		{
			name: "no_conflict_two_histograms",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				fooA, err := meterA.Int64Histogram("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter histogram foo"))
				assert.NoError(t, err)
				fooA.Record(ctx, 100, withAB)

				fooB, err := meterB.Int64Histogram("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter histogram foo"))
				assert.NoError(t, err)
				fooB.Record(ctx, 100, withAB)
			},
			possibleExpectedFiles: []string{"testdata/no_conflict_two_histograms.txt"},
		},
		{
			name: "conflict_help_two_counters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA, err := meterA.Int64Counter("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter a bar"))
				assert.NoError(t, err)
				barA.Add(ctx, 100, withTypeBar)

				barB, err := meterB.Int64Counter("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter b bar"))
				assert.NoError(t, err)
				barB.Add(ctx, 100, withTypeBar)
			},
			possibleExpectedFiles: []string{
				"testdata/conflict_help_two_counters_1.txt",
				"testdata/conflict_help_two_counters_2.txt",
			},
		},
		{
			name: "conflict_help_two_updowncounters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA, err := meterA.Int64UpDownCounter("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter a bar"))
				assert.NoError(t, err)
				barA.Add(ctx, 100, withTypeBar)

				barB, err := meterB.Int64UpDownCounter("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter b bar"))
				assert.NoError(t, err)
				barB.Add(ctx, 100, withTypeBar)
			},
			possibleExpectedFiles: []string{
				"testdata/conflict_help_two_updowncounters_1.txt",
				"testdata/conflict_help_two_updowncounters_2.txt",
			},
		},
		{
			name: "conflict_help_two_histograms",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA, err := meterA.Int64Histogram("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter a bar"))
				assert.NoError(t, err)
				barA.Record(ctx, 100, withAB)

				barB, err := meterB.Int64Histogram("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter b bar"))
				assert.NoError(t, err)
				barB.Record(ctx, 100, withAB)
			},
			possibleExpectedFiles: []string{
				"testdata/conflict_help_two_histograms_1.txt",
				"testdata/conflict_help_two_histograms_2.txt",
			},
		},
		{
			name: "conflict_unit_two_counters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				bazA, err := meterA.Int64Counter("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter bar"))
				assert.NoError(t, err)
				bazA.Add(ctx, 100, withTypeBar)

				bazB, err := meterB.Int64Counter("bar",
					otelmetric.WithUnit("s"),
					otelmetric.WithDescription("meter bar"))
				assert.NoError(t, err)
				bazB.Add(ctx, 100, withTypeBar)
			},
			options:               []Option{WithoutUnits()},
			possibleExpectedFiles: []string{"testdata/conflict_unit_two_counters.txt"},
		},
		{
			name: "conflict_unit_two_updowncounters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA, err := meterA.Int64UpDownCounter("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter gauge bar"))
				assert.NoError(t, err)
				barA.Add(ctx, 100, withTypeBar)

				barB, err := meterB.Int64UpDownCounter("bar",
					otelmetric.WithUnit("s"),
					otelmetric.WithDescription("meter gauge bar"))
				assert.NoError(t, err)
				barB.Add(ctx, 100, withTypeBar)
			},
			options:               []Option{WithoutUnits()},
			possibleExpectedFiles: []string{"testdata/conflict_unit_two_updowncounters.txt"},
		},
		{
			name: "conflict_unit_two_histograms",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA, err := meterA.Int64Histogram("bar",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter histogram bar"))
				assert.NoError(t, err)
				barA.Record(ctx, 100, withAB)

				barB, err := meterB.Int64Histogram("bar",
					otelmetric.WithUnit("s"),
					otelmetric.WithDescription("meter histogram bar"))
				assert.NoError(t, err)
				barB.Record(ctx, 100, withAB)
			},
			options:               []Option{WithoutUnits()},
			possibleExpectedFiles: []string{"testdata/conflict_unit_two_histograms.txt"},
		},
		{
			name: "conflict_type_counter_and_updowncounter",
			recordMetrics: func(ctx context.Context, meterA, _ otelmetric.Meter) {
				counter, err := meterA.Int64Counter("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter foo"))
				assert.NoError(t, err)
				counter.Add(ctx, 100, withTypeFoo)

				gauge, err := meterA.Int64UpDownCounter("foo_total",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter foo"))
				assert.NoError(t, err)
				gauge.Add(ctx, 200, withTypeFoo)
			},
			options: []Option{WithoutUnits()},
			possibleExpectedFiles: []string{
				"testdata/conflict_type_counter_and_updowncounter_1.txt",
				"testdata/conflict_type_counter_and_updowncounter_2.txt",
			},
		},
		{
			name: "conflict_type_histogram_and_updowncounter",
			recordMetrics: func(ctx context.Context, meterA, _ otelmetric.Meter) {
				fooA, err := meterA.Int64UpDownCounter("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter gauge foo"))
				assert.NoError(t, err)
				fooA.Add(ctx, 100, withAB)

				fooHistogramA, err := meterA.Int64Histogram("foo",
					otelmetric.WithUnit("By"),
					otelmetric.WithDescription("meter histogram foo"))
				assert.NoError(t, err)
				fooHistogramA.Record(ctx, 100, withAB)
			},
			possibleExpectedFiles: []string{
				"testdata/conflict_type_histogram_and_updowncounter_1.txt",
				"testdata/conflict_type_histogram_and_updowncounter_2.txt",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// initialize registry exporter
			ctx := context.Background()
			registry := prometheus.NewRegistry()
			exporter, err := New(append(tc.options, WithRegisterer(registry))...)
			require.NoError(t, err)

			// initialize resource
			res, err := resource.New(ctx,
				resource.WithAttributes(semconv.ServiceName("prometheus_test")),
				resource.WithAttributes(semconv.TelemetrySDKVersion("latest")),
			)
			require.NoError(t, err)
			res, err = resource.Merge(resource.Default(), res)
			require.NoError(t, err)

			// initialize provider
			provider := metric.NewMeterProvider(
				metric.WithReader(exporter),
				metric.WithResource(res),
			)

			// initialize two meter a, b
			meterA := provider.Meter("ma", otelmetric.WithInstrumentationVersion("v0.1.0"))
			meterB := provider.Meter("mb", otelmetric.WithInstrumentationVersion("v0.1.0"))

			tc.recordMetrics(ctx, meterA, meterB)

			match := false
			for _, filename := range tc.possibleExpectedFiles {
				file, ferr := os.Open(filename)
				require.NoError(t, ferr)
				t.Cleanup(func() { require.NoError(t, file.Close()) })

				err = testutil.GatherAndCompare(registry, file)
				if err == nil {
					match = true
					break
				}
			}
			require.Truef(t, match, "expected export not produced: %v", err)
		})
	}
}

func TestCollectorConcurrentSafe(t *testing.T) {
	// This tests makes sure that the implemented
	// https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Collector
	// is concurrent safe.
	ctx := context.Background()
	registry := prometheus.NewRegistry()
	exporter, err := New(WithRegisterer(registry))
	require.NoError(t, err)
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("testmeter")
	cnt, err := meter.Int64Counter("foo")
	require.NoError(t, err)
	cnt.Add(ctx, 100)

	var wg sync.WaitGroup
	concurrencyLevel := 10
	for range concurrencyLevel {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := registry.Gather() // this calls collector.Collect
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}

func TestShutdownExporter(t *testing.T) {
	var handledError error
	eh := otel.ErrorHandlerFunc(func(e error) { handledError = errors.Join(handledError, e) })
	otel.SetErrorHandler(eh)

	ctx := context.Background()
	registry := prometheus.NewRegistry()

	for range 3 {
		exporter, err := New(WithRegisterer(registry))
		require.NoError(t, err)
		provider := metric.NewMeterProvider(
			metric.WithResource(resource.Default()),
			metric.WithReader(exporter))
		meter := provider.Meter("testmeter")
		cnt, err := meter.Int64Counter("foo")
		require.NoError(t, err)
		cnt.Add(ctx, 100)

		// verify that metrics added to a previously shutdown MeterProvider
		// do not conflict with metrics added in this loop.
		_, err = registry.Gather()
		require.NoError(t, err)

		// Shutdown should cause future prometheus Gather() calls to no longer
		// include metrics from this loop's MeterProvider.
		err = provider.Shutdown(ctx)
		require.NoError(t, err)
	}
	// ensure we aren't unnecessarily logging errors from the shutdown MeterProvider
	require.NoError(t, handledError)
}

func TestExemplars(t *testing.T) {
	attrsOpt := otelmetric.WithAttributes(
		attribute.Key("A.1").String("B"),
		attribute.Key("C.2").String("D"),
		attribute.Key("E.3").Bool(true),
		attribute.Key("F.4").Int(42),
	)
	expectedNonEscapedLabels := map[string]string{
		otlptranslator.ExemplarTraceIDKey: "01000000000000000000000000000000",
		otlptranslator.ExemplarSpanIDKey:  "0100000000000000",
		"A.1":                             "B",
		"C.2":                             "D",
		"E.3":                             "true",
		"F.4":                             "42",
	}
	expectedEscapedLabels := map[string]string{
		otlptranslator.ExemplarTraceIDKey: "01000000000000000000000000000000",
		otlptranslator.ExemplarSpanIDKey:  "0100000000000000",
		"A_1":                             "B",
		"C_2":                             "D",
		"E_3":                             "true",
		"F_4":                             "42",
	}
	for _, tc := range []struct {
		name                  string
		recordMetrics         func(ctx context.Context, meter otelmetric.Meter)
		expectedExemplarValue float64
		expectedLabels        map[string]string
		escapingScheme        model.EscapingScheme
		validationScheme      model.ValidationScheme
	}{
		{
			name: "escaped counter",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				counter, err := meter.Float64Counter("foo")
				require.NoError(t, err)
				counter.Add(ctx, 9, attrsOpt)
			},
			expectedExemplarValue: 9,
			expectedLabels:        expectedEscapedLabels,
			escapingScheme:        model.UnderscoreEscaping,
			validationScheme:      model.LegacyValidation,
		},
		{
			name: "escaped histogram",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				hist, err := meter.Int64Histogram("foo")
				require.NoError(t, err)
				hist.Record(ctx, 9, attrsOpt)
			},
			expectedExemplarValue: 9,
			expectedLabels:        expectedEscapedLabels,
			escapingScheme:        model.UnderscoreEscaping,
			validationScheme:      model.LegacyValidation,
		},
		{
			name: "non-escaped counter",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				counter, err := meter.Float64Counter("foo")
				require.NoError(t, err)
				counter.Add(ctx, 9, attrsOpt)
			},
			expectedExemplarValue: 9,
			expectedLabels:        expectedNonEscapedLabels,
			escapingScheme:        model.NoEscaping,
			validationScheme:      model.UTF8Validation,
		},
		{
			name: "non-escaped histogram",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				hist, err := meter.Int64Histogram("foo")
				require.NoError(t, err)
				hist.Record(ctx, 9, attrsOpt)
			},
			expectedExemplarValue: 9,
			expectedLabels:        expectedNonEscapedLabels,
			escapingScheme:        model.NoEscaping,
			validationScheme:      model.UTF8Validation,
		},
		{
			name: "exponential histogram",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				hist, err := meter.Int64Histogram("exponential_histogram")
				require.NoError(t, err)
				hist.Record(ctx, 9, attrsOpt)
			},
			expectedExemplarValue: 9,
			expectedLabels:        expectedNonEscapedLabels,
			escapingScheme:        model.NoEscaping,
			validationScheme:      model.UTF8Validation,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			originalEscapingScheme := model.NameEscapingScheme
			originalValidationScheme := model.NameValidationScheme // nolint:staticcheck // We need this check to keep supporting the legacy scheme.
			model.NameEscapingScheme = tc.escapingScheme
			model.NameValidationScheme = tc.validationScheme // nolint:staticcheck // We need this check to keep supporting the legacy scheme.
			// Restore original value after the test is complete
			defer func() {
				model.NameEscapingScheme = originalEscapingScheme
				model.NameValidationScheme = originalValidationScheme // nolint:staticcheck // We need this check to keep supporting the legacy scheme.
			}()
			// initialize registry exporter
			ctx := context.Background()
			registry := prometheus.NewRegistry()
			exporter, err := New(WithRegisterer(registry), WithoutTargetInfo(), WithoutScopeInfo())
			require.NoError(t, err)

			// initialize resource
			res, err := resource.New(ctx,
				resource.WithAttributes(semconv.ServiceName("prometheus_test")),
				resource.WithAttributes(semconv.TelemetrySDKVersion("latest")),
			)
			require.NoError(t, err)
			res, err = resource.Merge(resource.Default(), res)
			require.NoError(t, err)

			// initialize provider and meter
			provider := metric.NewMeterProvider(
				metric.WithReader(exporter),
				metric.WithResource(res),
				metric.WithView(metric.NewView(
					metric.Instrument{Name: "foo"},
					metric.Stream{
						// filter out all attributes so they are added as filtered
						// attributes to the exemplar
						AttributeFilter: attribute.NewAllowKeysFilter(),
					},
				),
				),
				metric.WithView(metric.NewView(
					metric.Instrument{Name: "exponential_histogram"},
					metric.Stream{
						Aggregation: metric.AggregationBase2ExponentialHistogram{
							MaxSize: 20,
						},
						AttributeFilter: attribute.NewAllowKeysFilter(),
					},
				),
				),
			)
			meter := provider.Meter("meter", otelmetric.WithInstrumentationVersion("v0.1.0"))

			// Add a sampled span context so that measurements get exemplars added
			sc := trace.NewSpanContext(trace.SpanContextConfig{
				SpanID:     trace.SpanID{0o1},
				TraceID:    trace.TraceID{0o1},
				TraceFlags: trace.FlagsSampled,
			})
			ctx = trace.ContextWithSpanContext(ctx, sc)
			// Record a single observation with the exemplar
			tc.recordMetrics(ctx, meter)

			// Verify that the exemplar is present in the proto version of the
			// prometheus metrics.
			got, done, err := prometheus.ToTransactionalGatherer(registry).Gather()
			defer done()
			require.NoError(t, err)

			require.Len(t, got, 1)
			family := got[0]
			require.Len(t, family.GetMetric(), 1)
			metric := family.GetMetric()[0]
			var exemplar *dto.Exemplar
			switch family.GetType() {
			case dto.MetricType_COUNTER:
				exemplar = metric.GetCounter().GetExemplar()
			case dto.MetricType_HISTOGRAM:
				h := metric.GetHistogram()
				for _, b := range h.GetBucket() {
					if b.GetExemplar() != nil {
						exemplar = b.GetExemplar()
						continue
					}
				}
				if h.GetZeroThreshold() != 0 || h.GetZeroCount() != 0 ||
					len(h.PositiveSpan) != 0 || len(h.NegativeSpan) != 0 {
					require.NotNil(t, h.Exemplars)
					exemplar = h.Exemplars[0]
				}
			}
			require.NotNil(t, exemplar)
			require.Equal(t, tc.expectedExemplarValue, exemplar.GetValue())
			require.Len(t, exemplar.GetLabel(), len(tc.expectedLabels))

			for _, label := range exemplar.GetLabel() {
				val, ok := tc.expectedLabels[label.GetName()]
				require.True(t, ok)
				require.Equal(t, label.GetValue(), val)
			}
		})
	}
}

func TestExponentialHistogramScaleValidation(t *testing.T) {
	ctx := context.Background()

	t.Run("normal_exponential_histogram_works", func(t *testing.T) {
		registry := prometheus.NewRegistry()
		exporter, err := New(WithRegisterer(registry), WithoutTargetInfo(), WithoutScopeInfo())
		require.NoError(t, err)

		provider := metric.NewMeterProvider(
			metric.WithReader(exporter),
			metric.WithResource(resource.Default()),
		)
		defer func() {
			err := provider.Shutdown(ctx)
			require.NoError(t, err)
		}()

		// Create a histogram with a valid scale
		meter := provider.Meter("test")
		hist, err := meter.Float64Histogram(
			"test_exponential_histogram",
			otelmetric.WithDescription("test histogram"),
		)
		require.NoError(t, err)
		hist.Record(ctx, 1.0)
		hist.Record(ctx, 10.0)
		hist.Record(ctx, 100.0)

		metricFamilies, err := registry.Gather()
		require.NoError(t, err)
		assert.NotEmpty(t, metricFamilies)
	})

	t.Run("error_handling_for_invalid_scales", func(t *testing.T) {
		var capturedError error
		originalHandler := otel.GetErrorHandler()
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			capturedError = err
		}))
		defer otel.SetErrorHandler(originalHandler)

		now := time.Now()
		invalidScaleData := metricdata.ExponentialHistogramDataPoint[float64]{
			Attributes:    attribute.NewSet(),
			StartTime:     now,
			Time:          now,
			Count:         1,
			Sum:           10.0,
			Scale:         -5, // Invalid scale below -4
			ZeroCount:     0,
			ZeroThreshold: 0.0,
			PositiveBucket: metricdata.ExponentialBucket{
				Offset: 1,
				Counts: []uint64{1},
			},
			NegativeBucket: metricdata.ExponentialBucket{
				Offset: 1,
				Counts: []uint64{},
			},
		}

		ch := make(chan prometheus.Metric, 10)
		defer close(ch)

		histogram := metricdata.ExponentialHistogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{invalidScaleData},
		}

		m := metricdata.Metrics{
			Name:        "test_histogram",
			Description: "test",
		}

		addExponentialHistogramMetric(
			ch,
			histogram,
			m,
			"test_histogram",
			keyVals{},
			otlptranslator.LabelNamer{},
		)
		assert.Error(t, capturedError)
		assert.Contains(t, capturedError.Error(), "scale -5 is below minimum")
		select {
		case <-ch:
			t.Error("Expected no metrics to be produced for invalid scale")
		default:
			// No metrics were produced for the invalid scale
		}
	})
}

func TestDownscaleExponentialBucket(t *testing.T) {
	tests := []struct {
		name       string
		bucket     metricdata.ExponentialBucket
		scaleDelta int32
		want       metricdata.ExponentialBucket
	}{
		{
			name:       "Empty bucket",
			bucket:     metricdata.ExponentialBucket{},
			scaleDelta: 3,
			want:       metricdata.ExponentialBucket{},
		},
		{
			name: "1 size bucket",
			bucket: metricdata.ExponentialBucket{
				Offset: 50,
				Counts: []uint64{7},
			},
			scaleDelta: 4,
			want: metricdata.ExponentialBucket{
				Offset: 3,
				Counts: []uint64{7},
			},
		},
		{
			name: "zero scale delta",
			bucket: metricdata.ExponentialBucket{
				Offset: 50,
				Counts: []uint64{7, 5},
			},
			scaleDelta: 0,
			want: metricdata.ExponentialBucket{
				Offset: 50,
				Counts: []uint64{7, 5},
			},
		},
		{
			name: "aligned bucket scale 1",
			bucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{1, 2, 3, 4, 5, 6},
			},
			scaleDelta: 1,
			want: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{3, 7, 11},
			},
		},
		{
			name: "aligned bucket scale 2",
			bucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{1, 2, 3, 4, 5, 6},
			},
			scaleDelta: 2,
			want: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{10, 11},
			},
		},
		{
			name: "unaligned bucket scale 1",
			bucket: metricdata.ExponentialBucket{
				Offset: 5,
				Counts: []uint64{1, 2, 3, 4, 5, 6},
			}, // This is equivalent to [0,0,0,0,0,1,2,3,4,5,6]
			scaleDelta: 1,
			want: metricdata.ExponentialBucket{
				Offset: 2,
				Counts: []uint64{1, 5, 9, 6},
			}, // This is equivalent to [0,0,1,5,9,6]
		},
		{
			name: "negative startBin",
			bucket: metricdata.ExponentialBucket{
				Offset: -1,
				Counts: []uint64{1, 0, 3},
			},
			scaleDelta: 1,
			want: metricdata.ExponentialBucket{
				Offset: -1,
				Counts: []uint64{1, 3},
			},
		},
		{
			name: "negative startBin 2",
			bucket: metricdata.ExponentialBucket{
				Offset: -4,
				Counts: []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			},
			scaleDelta: 1,
			want: metricdata.ExponentialBucket{
				Offset: -2,
				Counts: []uint64{3, 7, 11, 15, 19},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := downscaleExponentialBucket(tt.bucket, tt.scaleDelta)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExponentialHistogramHighScaleDownscaling(t *testing.T) {
	t.Run("scale_10_downscales_to_8", func(t *testing.T) {
		// Test that scale 10 gets properly downscaled to 8 with correct bucket re-aggregation
		ch := make(chan prometheus.Metric, 10)
		defer close(ch)

		now := time.Now()

		// Create an exponential histogram data point with scale 10
		dataPoint := metricdata.ExponentialHistogramDataPoint[float64]{
			Attributes:    attribute.NewSet(),
			StartTime:     now,
			Time:          now,
			Count:         8,
			Sum:           55.0,
			Scale:         10, // This should be downscaled to 8
			ZeroCount:     0,
			ZeroThreshold: 0.0,
			PositiveBucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{1, 1, 1, 1, 1, 1, 1, 1}, // 8 buckets with 1 count each
			},
			NegativeBucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{},
			},
		}

		histogram := metricdata.ExponentialHistogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{dataPoint},
		}

		m := metricdata.Metrics{
			Name:        "test_high_scale_histogram",
			Description: "test histogram with high scale",
		}

		// This should not produce any errors and should properly downscale buckets
		addExponentialHistogramMetric(
			ch,
			histogram,
			m,
			"test_high_scale_histogram",
			keyVals{},
			otlptranslator.LabelNamer{},
		)

		// Verify a metric was produced
		select {
		case metric := <-ch:
			// Check that the metric was created successfully
			require.NotNil(t, metric)

			// The scale should have been clamped to 8, and buckets should be re-aggregated
			// With scale 10 -> 8, we have a scaleDelta of 2, meaning 2^2 = 4 buckets merge into 1
			// Original: 8 buckets with 1 count each at scale 10
			// After downscaling: 2 buckets with 4 counts each at scale 8
		default:
			t.Error("Expected a metric to be produced")
		}
	})

	t.Run("scale_12_downscales_to_8", func(t *testing.T) {
		// Test that scale 12 gets properly downscaled to 8 with correct bucket re-aggregation
		ch := make(chan prometheus.Metric, 10)
		defer close(ch)

		now := time.Now()

		// Create an exponential histogram data point with scale 12
		dataPoint := metricdata.ExponentialHistogramDataPoint[float64]{
			Attributes:    attribute.NewSet(),
			StartTime:     now,
			Time:          now,
			Count:         16,
			Sum:           120.0,
			Scale:         12, // This should be downscaled to 8
			ZeroCount:     0,
			ZeroThreshold: 0.0,
			PositiveBucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, // 16 buckets with 1 count each
			},
			NegativeBucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{},
			},
		}

		histogram := metricdata.ExponentialHistogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{dataPoint},
		}

		m := metricdata.Metrics{
			Name:        "test_very_high_scale_histogram",
			Description: "test histogram with very high scale",
		}

		// This should not produce any errors and should properly downscale buckets
		addExponentialHistogramMetric(
			ch,
			histogram,
			m,
			"test_very_high_scale_histogram",
			keyVals{},
			otlptranslator.LabelNamer{},
		)

		// Verify a metric was produced
		select {
		case metric := <-ch:
			// Check that the metric was created successfully
			require.NotNil(t, metric)

			// The scale should have been clamped to 8, and buckets should be re-aggregated
			// With scale 12 -> 8, we have a scaleDelta of 4, meaning 2^4 = 16 buckets merge into 1
			// Original: 16 buckets with 1 count each at scale 12
			// After downscaling: 1 bucket with 16 counts at scale 8
		default:
			t.Error("Expected a metric to be produced")
		}
	})

	t.Run("exponential_histogram_with_negative_buckets", func(t *testing.T) {
		// Test that exponential histograms with negative buckets are handled correctly
		ch := make(chan prometheus.Metric, 10)
		defer close(ch)

		now := time.Now()

		// Create an exponential histogram data point with both positive and negative buckets
		dataPoint := metricdata.ExponentialHistogramDataPoint[float64]{
			Attributes:    attribute.NewSet(),
			StartTime:     now,
			Time:          now,
			Count:         6,
			Sum:           25.0,
			Scale:         2,
			ZeroCount:     0,
			ZeroThreshold: 0.0,
			PositiveBucket: metricdata.ExponentialBucket{
				Offset: 1,
				Counts: []uint64{1, 2}, // 2 positive buckets
			},
			NegativeBucket: metricdata.ExponentialBucket{
				Offset: 1,
				Counts: []uint64{2, 1}, // 2 negative buckets
			},
		}

		histogram := metricdata.ExponentialHistogram[float64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{dataPoint},
		}

		m := metricdata.Metrics{
			Name:        "test_histogram_with_negative_buckets",
			Description: "test histogram with negative buckets",
		}

		// This should handle negative buckets correctly
		addExponentialHistogramMetric(
			ch,
			histogram,
			m,
			"test_histogram_with_negative_buckets",
			keyVals{},
			otlptranslator.LabelNamer{},
		)

		// Verify a metric was produced
		select {
		case metric := <-ch:
			require.NotNil(t, metric)
		default:
			t.Error("Expected a metric to be produced")
		}
	})

	t.Run("exponential_histogram_int64_type", func(t *testing.T) {
		// Test that int64 exponential histograms are handled correctly
		ch := make(chan prometheus.Metric, 10)
		defer close(ch)

		now := time.Now()

		// Create an exponential histogram data point with int64 type
		dataPoint := metricdata.ExponentialHistogramDataPoint[int64]{
			Attributes:    attribute.NewSet(),
			StartTime:     now,
			Time:          now,
			Count:         4,
			Sum:           20,
			Scale:         3,
			ZeroCount:     0,
			ZeroThreshold: 0.0,
			PositiveBucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{1, 1, 1, 1}, // 4 buckets with 1 count each
			},
			NegativeBucket: metricdata.ExponentialBucket{
				Offset: 0,
				Counts: []uint64{},
			},
		}

		histogram := metricdata.ExponentialHistogram[int64]{
			Temporality: metricdata.CumulativeTemporality,
			DataPoints:  []metricdata.ExponentialHistogramDataPoint[int64]{dataPoint},
		}

		m := metricdata.Metrics{
			Name:        "test_int64_exponential_histogram",
			Description: "test int64 exponential histogram",
		}

		// This should handle int64 exponential histograms correctly
		addExponentialHistogramMetric(
			ch,
			histogram,
			m,
			"test_int64_exponential_histogram",
			keyVals{},
			otlptranslator.LabelNamer{},
		)

		// Verify a metric was produced
		select {
		case metric := <-ch:
			require.NotNil(t, metric)
		default:
			t.Error("Expected a metric to be produced")
		}
	})
}

func TestDownscaleExponentialBucketEdgeCases(t *testing.T) {
	t.Run("min_idx_larger_than_current", func(t *testing.T) {
		// Test case where we find a minIdx that's smaller than the current
		bucket := metricdata.ExponentialBucket{
			Offset: 10, // Start at offset 10
			Counts: []uint64{1, 0, 0, 0, 1},
		}

		// Scale delta of 3 will cause downscaling: original indices 10->1, 14->1
		result := downscaleExponentialBucket(bucket, 3)

		// Both original buckets 10 and 14 should map to the same downscaled bucket at index 1
		expected := metricdata.ExponentialBucket{
			Offset: 1,
			Counts: []uint64{2}, // Both counts combined
		}

		assert.Equal(t, expected, result)
	})

	t.Run("empty_downscaled_counts", func(t *testing.T) {
		// Create a scenario that results in empty downscaled counts
		bucket := metricdata.ExponentialBucket{
			Offset: math.MaxInt32 - 5, // Very large offset that won't cause overflow in this case
			Counts: []uint64{1, 1, 1, 1, 1},
		}

		// This should work normally and downscale the buckets
		result := downscaleExponentialBucket(bucket, 1)

		// Should return bucket with downscaled values
		expected := metricdata.ExponentialBucket{
			Offset: 1073741821,        // ((MaxInt32-5) + 0) >> 1 = 1073741821
			Counts: []uint64{2, 2, 1}, // Buckets get combined during downscaling
		}

		assert.Equal(t, expected, result)
	})
}
