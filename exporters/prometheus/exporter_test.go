// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/testutil"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/otlptranslator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus/internal/observ"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

// producerFunc adapts a function to implement metric.Producer.
type producerFunc func(context.Context) ([]metricdata.ScopeMetrics, error)

func (f producerFunc) Produce(ctx context.Context) ([]metricdata.ScopeMetrics, error) { return f(ctx) }

// Helper: scrape with ContinueOnError and return body + status.
func scrapeWithContinueOnError(reg *prometheus.Registry) (int, string) {
	h := promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			ErrorHandling: promhttp.ContinueOnError,
		},
	)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody)
	h.ServeHTTP(rr, req)

	return rr.Code, rr.Body.String()
}

func TestPrometheusExporter(t *testing.T) {
	testCases := []struct {
		name                string
		emptyResource       bool
		customResourceAttrs []attribute.KeyValue
		recordMetrics       func(ctx context.Context, meter otelmetric.Meter)
		options             []Option
		expectedFile        string
		checkMetricFamilies func(t testing.TB, dtos []*dto.MetricFamily)
	}{
		{
			name:         "counter",
			expectedFile: "testdata/counter.txt",
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
			options: []Option{
				WithNamespace("my.dotted.namespace"),
				WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
			},
		},
		{
			name:         "counter that already has the unit suffix",
			expectedFile: "testdata/counter_noutf8_with_unit_suffix.txt",
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
			options: []Option{WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes)},
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
					"foo.dotted",
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
			options: []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			options: []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
		},
		{
			name:         "counter that already has a total suffix",
			expectedFile: "testdata/counter.txt",
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
			options: []Option{
				WithNamespace("my.dotted.namespace"),
				WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
			},
		},
		{
			name:         "counter with suffixes disabled",
			expectedFile: "testdata/counter_disabled_suffix.txt",
			options: []Option{
				WithoutCounterSuffixes(),
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
			},
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
			options:      []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			options:      []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			options:      []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			options: []Option{
				WithoutUnits(),
				WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
			},
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
			options:      []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			options:       []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			options:      []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			name: "without target_info",
			options: []Option{
				WithoutTargetInfo(),
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
			},
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
			name: "without scope_info",
			options: []Option{
				WithoutScopeInfo(),
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
			},
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
			name: "without scope_info and target_info",
			options: []Option{
				WithoutScopeInfo(),
				WithoutTargetInfo(),
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
			},
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
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
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
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
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
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
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
			options: []Option{
				WithNamespace("my.dotted.namespace"),
				WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes),
			},
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
			name:         "counter utf-8 notranslation",
			expectedFile: "testdata/counter_utf8_notranslation.txt",
			options: []Option{
				WithNamespace("my.dotted.namespace"),
				WithTranslationStrategy(otlptranslator.NoTranslation),
			},
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
			options:      []Option{WithTranslationStrategy(otlptranslator.NoUTF8EscapingWithSuffixes)},
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
			ctx := t.Context()
			registry := prometheus.NewRegistry()
			opts := append(tc.options, WithRegisterer(registry))
			exporter, err := New(opts...)
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
					metric.NewView(
						metric.Instrument{Name: "test_dropped_attrs_limit"},
						metric.Stream{
							AttributeFilter: func(kv attribute.KeyValue) bool {
								return kv.Key != "long_attribute"
							},
						},
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
	ctx := t.Context()
	registry := prometheus.NewRegistry()
	exporter, err := New(
		WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
		WithRegisterer(registry),
	)
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

func TestBridgeScopeIgnored(t *testing.T) {
	var handledError error
	eh := otel.ErrorHandlerFunc(func(e error) { handledError = errors.Join(handledError, e) })
	otel.SetErrorHandler(eh)
	ctx := t.Context()
	registry := prometheus.NewRegistry()
	exporter, err := New(
		WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
		WithRegisterer(registry),
	)
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

	fooCounter, err := provider.Meter(bridgeScopeName, otelmetric.WithInstrumentationVersion("v0.1.0")).
		Int64Counter(
			"foo",
			otelmetric.WithUnit("s"),
			otelmetric.WithDescription("meter foo counter"))
	assert.NoError(t, err)
	fooCounter.Add(ctx, 100, otelmetric.WithAttributes(attribute.String("type", "foo")))

	file, err := os.Open("testdata/just_target_info.txt")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, file.Close()) })

	err = testutil.GatherAndCompare(registry, file)
	require.NoError(t, err)

	require.ErrorIs(t, handledError, errBridgeNotSupported)
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
		expectGatherError     bool
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
			expectGatherError: true,
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
			expectGatherError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// initialize registry exporter
			ctx := t.Context()
			registry := prometheus.NewRegistry()
			// This test does not set the Translation Strategy, so it defaults to
			// UnderscoreEscapingWithSuffixes.
			opts := append(
				[]Option{
					WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
				},
				tc.options...,
			)
			exporter, err := New(append(opts, WithRegisterer(registry))...)
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

			if tc.expectGatherError {
				// With improved error handling, conflicting instrument types emit an invalid metric.
				// Gathering should surface an error instead of silently dropping.
				_, err := registry.Gather()
				require.Error(t, err)

				// 2) Also assert what users will see if they opt into ContinueOnError.
				// Compare the HTTP body to an expected file that contains only the valid series
				// (e.g., "target_info" and any non-conflicting families).
				status, body := scrapeWithContinueOnError(registry)
				require.Equal(t, http.StatusOK, status)

				matched := false
				for _, filename := range tc.possibleExpectedFiles {
					want, ferr := os.ReadFile(filename)
					require.NoError(t, ferr)
					if body == string(want) {
						matched = true
						break
					}
				}
				require.Truef(t, matched, "expected export not produced under ContinueOnError; got:\n%s", body)
			} else {
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
			}
		})
	}
}

func TestCollectorConcurrentSafe(t *testing.T) {
	// This tests makes sure that the implemented
	// https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Collector
	// is concurrent safe.
	ctx := t.Context()
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

	ctx := t.Context()
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
		strategy              otlptranslator.TranslationStrategyOption
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
			strategy:              otlptranslator.UnderscoreEscapingWithSuffixes,
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
			strategy:              otlptranslator.UnderscoreEscapingWithSuffixes,
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
			strategy:              otlptranslator.NoTranslation,
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
			strategy:              otlptranslator.NoTranslation,
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
			strategy:              otlptranslator.NoTranslation,
		},
		{
			name: "exemplar overflow does not drop exemplar",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				hist, err := meter.Float64Histogram("exponential_histogram")
				require.NoError(t, err)

				// Create attributes that exceed 128-rune limit after accounting for trace_id/span_id
				// trace_id (32) + span_id (16) + "=" (1) = 49 characters
				// Remaining space: 128 - 49 = 79 characters
				// longVal (80 chars) + "=" = 81 chars ✓
				//  81 chars > 79, so long_value should be truncated

				longVal := strings.Repeat("B", 80) // 80 chars

				hist.Record(ctx, 0, otelmetric.WithAttributes(
					attribute.String("long_value", longVal),
				), attrsOpt)
			},
			expectedLabels: expectedEscapedLabels,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// initialize registry exporter
			ctx := t.Context()
			registry := prometheus.NewRegistry()
			exporter, err := New(
				WithRegisterer(registry),
				WithoutTargetInfo(),
				WithoutScopeInfo(),
				WithTranslationStrategy(tc.strategy),
			)
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
	ctx := t.Context()

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
			nil,
			t.Context(),
		)
		// Expect an invalid metric to be sent that carries the scale error.
		var pm prometheus.Metric
		select {
		case pm = <-ch:
		default:
			t.Fatalf("expected an invalid metric to be emitted for invalid scale, but channel was empty")
		}
		var dtoMetric dto.Metric
		werr := pm.Write(&dtoMetric)
		require.ErrorIs(t, werr, errEHScaleBelowMin)
		// The exporter reports via invalid metric, not the global otel error handler.
		assert.NoError(t, capturedError)
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
			nil,
			t.Context(),
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
			nil,
			t.Context(),
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
			nil,
			t.Context(),
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
			nil,
			t.Context(),
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

// TestEscapingErrorHandling increases test coverage by exercising some error
// conditions.
func TestEscapingErrorHandling(t *testing.T) {
	// Helper to create a producer that emits a Summary (unsupported) metric.
	makeSummaryProducer := func() metric.Producer {
		return producerFunc(func(_ context.Context) ([]metricdata.ScopeMetrics, error) {
			return []metricdata.ScopeMetrics{
				{
					Metrics: []metricdata.Metrics{
						{
							Name:        "summary_metric",
							Description: "unsupported summary",
							Data:        metricdata.Summary{},
						},
					},
				},
			}, nil
		})
	}
	// Helper to create a producer that emits a metric with an invalid name, to
	// force getName() to fail and exercise reportError at that branch.
	makeBadNameProducer := func() metric.Producer {
		return producerFunc(func(_ context.Context) ([]metricdata.ScopeMetrics, error) {
			return []metricdata.ScopeMetrics{
				{
					Metrics: []metricdata.Metrics{
						{
							Name:        "$%^&", // intentionally invalid; translation should fail normalization
							Description: "bad name for translation",
							// Any supported type is fine; getName runs before add* functions.
							Data: metricdata.Gauge[float64]{
								DataPoints: []metricdata.DataPoint[float64]{
									{Value: 1},
								},
							},
						},
					},
				},
			}, nil
		})
	}
	// Helper to create a producer that emits an ExponentialHistogram with a bad
	// label, to exercise addExponentialHistogramMetric getAttrs error path.
	makeBadEHProducer := func() metric.Producer {
		return producerFunc(func(_ context.Context) ([]metricdata.ScopeMetrics, error) {
			return []metricdata.ScopeMetrics{
				{
					Metrics: []metricdata.Metrics{
						{
							Name:        "exp_hist_metric",
							Description: "bad label",
							Data: metricdata.ExponentialHistogram[float64]{
								DataPoints: []metricdata.ExponentialHistogramDataPoint[float64]{
									{
										Attributes:    attribute.NewSet(attribute.Key("$%^&").String("B")),
										Scale:         0,
										Count:         1,
										ZeroThreshold: 0,
									},
								},
							},
						},
					},
				},
			}, nil
		})
	}
	// Helper to create a producer that emits an ExponentialHistogram with
	// inconsistent bucket counts vs total Count to trigger constructor error in addExponentialHistogramMetric.
	makeBadEHCountProducer := func() metric.Producer {
		return producerFunc(func(_ context.Context) ([]metricdata.ScopeMetrics, error) {
			return []metricdata.ScopeMetrics{
				{
					Metrics: []metricdata.Metrics{
						{
							Name: "exp_hist_metric_bad",
							Data: metricdata.ExponentialHistogram[float64]{
								DataPoints: []metricdata.ExponentialHistogramDataPoint[float64]{
									{
										Scale:         0,
										Count:         0,
										ZeroThreshold: 0,
										PositiveBucket: metricdata.ExponentialBucket{
											Offset: 0,
											Counts: []uint64{1},
										},
									},
								},
							},
						},
					},
				},
			}, nil
		})
	}

	testCases := []struct {
		name                    string
		namespace               string
		counterName             string
		customScopeAttrs        []attribute.KeyValue
		customResourceAttrs     []attribute.KeyValue
		labelName               string
		producer                metric.Producer
		skipInstrument          bool
		record                  func(ctx context.Context, meter otelmetric.Meter) error
		expectNewErr            string
		expectMetricErr         string
		expectGatherErrContains string
		expectGatherErrIs       error
		checkMetricFamilies     func(t testing.TB, dtos []*dto.MetricFamily)
	}{
		{
			name:        "simple happy path",
			counterName: "foo",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.Len(t, mfs, 2)
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						continue
					}
					require.Equal(t, "foo_seconds_total", mf.GetName())
				}
			},
		},
		{
			name:         "bad namespace",
			namespace:    "$%^&",
			counterName:  "foo",
			expectNewErr: `normalization for label name "$%^&" resulted in invalid name "_"`,
		},
		{
			name: "bad translated metric name via producer",
			// Use a producer to emit a metric with an invalid name to trigger getName error.
			producer:       makeBadNameProducer(),
			skipInstrument: true,
			// Error message comes from normalization in the translator; match on a stable substring.
			expectGatherErrContains: "normalization",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				// target_info should still be exported.
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name:        "good namespace, names should be escaped",
			namespace:   "my-strange-namespace",
			counterName: "foo",
			labelName:   "bar",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						continue
					}
					require.Contains(t, mf.GetName(), "my_strange_namespace")
					require.NotContains(t, mf.GetName(), "my-strange-namespace")
				}
			},
		},
		{
			name:        "bad resource attribute",
			counterName: "foo",
			customResourceAttrs: []attribute.KeyValue{
				attribute.Key("$%^&").String("B"),
			},
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.Empty(t, mfs)
			},
		},
		{
			name:        "bad scope metric attribute",
			counterName: "foo",
			customScopeAttrs: []attribute.KeyValue{
				attribute.Key("$%^&").String("B"),
			},
			// With improved error handling, invalid scope label names result in an invalid metric
			// and Gather returns an error containing the normalization failure.
			expectGatherErrContains: `normalization for label name "$%^&" resulted in invalid name "_"`,
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				// target_info should still be exported; metric with bad scope label dropped.
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name:            "bad translated metric name",
			counterName:     "$%^&",
			expectMetricErr: `invalid instrument name: $%^&: must start with a letter`,
		},
		{
			// label names are not translated and therefore not checked until
			// collection time; with improved error handling, we emit an invalid metric and
			// surface the error during Gather.
			name:                    "bad translated label name",
			counterName:             "foo",
			labelName:               "$%^&",
			expectGatherErrContains: `normalization for label name "$%^&" resulted in invalid name "_"`,
		},
		{
			name: "unsupported data type via producer",
			// Use a producer to emit a Summary data point; no SDK instruments.
			producer:          makeSummaryProducer(),
			skipInstrument:    true,
			expectGatherErrIs: errInvalidMetricType,
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name:                    "bad exponential histogram label name via producer",
			producer:                makeBadEHProducer(),
			skipInstrument:          true,
			expectGatherErrContains: `normalization for label name "$%^&" resulted in invalid name "_"`,
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name:                    "exponential histogram constructor error via producer (count mismatch)",
			producer:                makeBadEHCountProducer(),
			skipInstrument:          true,
			expectGatherErrContains: "count",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name: "sum constructor error via duplicate label name",
			record: func(ctx context.Context, meter otelmetric.Meter) error {
				c, err := meter.Int64Counter("sum_metric_dup")
				if err != nil {
					return err
				}
				// Duplicate variable label name with scope label to make Desc invalid.
				c.Add(ctx, 1, otelmetric.WithAttributes(attribute.String(scopeNameLabel, "x")))
				return nil
			},
			expectGatherErrContains: "duplicate label",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name: "gauge constructor error via duplicate label name",
			record: func(ctx context.Context, meter otelmetric.Meter) error {
				g, err := meter.Float64Gauge("gauge_metric_dup")
				if err != nil {
					return err
				}
				g.Record(ctx, 1.0, otelmetric.WithAttributes(attribute.String(scopeNameLabel, "x")))
				return nil
			},
			expectGatherErrContains: "duplicate label",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name: "histogram constructor error via duplicate label name",
			record: func(ctx context.Context, meter otelmetric.Meter) error {
				h, err := meter.Float64Histogram("hist_metric_dup")
				if err != nil {
					return err
				}
				h.Record(ctx, 1.23, otelmetric.WithAttributes(attribute.String(scopeNameLabel, "x")))
				return nil
			},
			expectGatherErrContains: "duplicate label",
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name: "bad gauge label name",
			record: func(ctx context.Context, meter otelmetric.Meter) error {
				g, err := meter.Float64Gauge("gauge_metric")
				if err != nil {
					return err
				}
				g.Record(ctx, 1, otelmetric.WithAttributes(attribute.Key("$%^&").String("B")))
				return nil
			},
			expectGatherErrContains: `normalization for label name "$%^&" resulted in invalid name "_"`,
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
		{
			name: "bad histogram label name",
			record: func(ctx context.Context, meter otelmetric.Meter) error {
				h, err := meter.Float64Histogram("hist_metric")
				if err != nil {
					return err
				}
				h.Record(ctx, 1.23, otelmetric.WithAttributes(attribute.Key("$%^&").String("B")))
				return nil
			},
			expectGatherErrContains: `normalization for label name "$%^&" resulted in invalid name "_"`,
			checkMetricFamilies: func(t testing.TB, mfs []*dto.MetricFamily) {
				require.NotEmpty(t, mfs)
				other := 0
				seenTarget := false
				for _, mf := range mfs {
					if mf.GetName() == "target_info" {
						seenTarget = true
						continue
					}
					other++
				}
				require.True(t, seenTarget)
				require.Equal(t, 0, other)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			registry := prometheus.NewRegistry()

			sc := trace.NewSpanContext(trace.SpanContextConfig{
				SpanID:     trace.SpanID{0o1},
				TraceID:    trace.TraceID{0o1},
				TraceFlags: trace.FlagsSampled,
			})
			ctx = trace.ContextWithSpanContext(ctx, sc)

			opts := []Option{
				WithRegisterer(registry),
				WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
				WithNamespace(tc.namespace),
				WithResourceAsConstantLabels(attribute.NewDenyKeysFilter()),
			}
			if tc.producer != nil {
				opts = append(opts, WithProducer(tc.producer))
			}
			exporter, err := New(opts...)
			if tc.expectNewErr != "" {
				require.ErrorContains(t, err, tc.expectNewErr)
				return
			}
			require.NoError(t, err)
			if !tc.skipInstrument {
				res, err := resource.New(ctx,
					resource.WithAttributes(semconv.ServiceName("prometheus_test")),
					resource.WithAttributes(semconv.TelemetrySDKVersion("latest")),
					resource.WithAttributes(tc.customResourceAttrs...),
				)
				require.NoError(t, err)
				provider := metric.NewMeterProvider(
					metric.WithReader(exporter),
					metric.WithResource(res),
				)
				meter := provider.Meter(
					"meterfoo",
					otelmetric.WithInstrumentationVersion("v0.1.0"),
					otelmetric.WithInstrumentationAttributes(tc.customScopeAttrs...),
				)
				if tc.record != nil {
					err := tc.record(ctx, meter)
					require.NoError(t, err)
				} else {
					fooCounter, err := meter.Int64Counter(
						tc.counterName,
						otelmetric.WithUnit("s"),
						otelmetric.WithDescription(fmt.Sprintf(`meter %q counter`, tc.counterName)))
					if tc.expectMetricErr != "" {
						require.ErrorContains(t, err, tc.expectMetricErr)
						return
					}
					require.NoError(t, err)
					var addOpts []otelmetric.AddOption
					if tc.labelName != "" {
						addOpts = append(addOpts, otelmetric.WithAttributes(attribute.String(tc.labelName, "foo")))
					}
					fooCounter.Add(ctx, 100, addOpts...)
				}
			} else {
				// When skipping instruments, still register the reader so Collect will run.
				_ = metric.NewMeterProvider(metric.WithReader(exporter))
			}
			got, err := registry.Gather()
			if tc.expectGatherErrContains != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectGatherErrContains)
				return
			}
			if tc.expectGatherErrIs != nil {
				require.ErrorIs(t, err, tc.expectGatherErrIs)
				return
			}
			require.NoError(t, err)
			if tc.checkMetricFamilies != nil {
				tc.checkMetricFamilies(t, got)
			}
		})
	}
}

func TestExporterSelfInstrumentation(t *testing.T) {
	testCases := []struct {
		name                  string
		enableObservability   bool
		recordMetrics         func(ctx context.Context, meter otelmetric.Meter)
		expectedObservMetrics []string
		expectedMainMetrics   int
		checkMetrics          func(t *testing.T, mainMetrics []*dto.MetricFamily, observMetrics metricdata.ScopeMetrics)
	}{
		{
			name:                "self instrumentation disabled",
			enableObservability: false,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				counter, err := meter.Int64Counter("test_counter")
				require.NoError(t, err)
				counter.Add(ctx, 1)
			},
			expectedObservMetrics: []string{}, // No observability metrics expected
			expectedMainMetrics:   2,          // test counter + target_info
		},
		{
			name:                "self instrumentation enabled with counter",
			enableObservability: true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				counter, err := meter.Int64Counter("test_counter", otelmetric.WithDescription("test counter"))
				require.NoError(t, err)
				counter.Add(ctx, 1, otelmetric.WithAttributes(attribute.String("key", "value")))
			},
			expectedObservMetrics: []string{
				"otel.sdk.exporter.metric_data_point.inflight",
				"otel.sdk.exporter.metric_data_point.exported",
				"otel.sdk.exporter.operation.duration",
				"otel.sdk.metric_reader.collection.duration",
			},
			expectedMainMetrics: 2, // test counter + target_info
			checkMetrics: func(t *testing.T, _ []*dto.MetricFamily, observMetrics metricdata.ScopeMetrics) {
				// Check that exported metrics include success count
				for _, m := range observMetrics.Metrics {
					if m.Name == "otel.sdk.exporter.metric_data_point.exported" {
						sum, ok := m.Data.(metricdata.Sum[int64])
						require.True(t, ok)
						require.Len(t, sum.DataPoints, 1)
						require.Equal(t, int64(1), sum.DataPoints[0].Value)
					}
				}
			},
		},
		{
			name:                "self instrumentation enabled with gauge",
			enableObservability: true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				gauge, err := meter.Float64Gauge("test_gauge", otelmetric.WithDescription("test gauge"))
				require.NoError(t, err)
				gauge.Record(ctx, 42.5, otelmetric.WithAttributes(attribute.String("key", "value")))
			},
			expectedObservMetrics: []string{
				"otel.sdk.exporter.metric_data_point.inflight",
				"otel.sdk.exporter.metric_data_point.exported",
				"otel.sdk.exporter.operation.duration",
				"otel.sdk.metric_reader.collection.duration",
			},
			expectedMainMetrics: 2,
		},
		{
			name:                "self instrumentation enabled with histogram",
			enableObservability: true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				histogram, err := meter.Float64Histogram("test_histogram", otelmetric.WithDescription("test histogram"))
				require.NoError(t, err)
				histogram.Record(ctx, 1.5, otelmetric.WithAttributes(attribute.String("key", "value")))
				histogram.Record(ctx, 2.5, otelmetric.WithAttributes(attribute.String("key", "value")))
			},
			expectedObservMetrics: []string{
				"otel.sdk.exporter.metric_data_point.inflight",
				"otel.sdk.exporter.metric_data_point.exported",
				"otel.sdk.exporter.operation.duration",
				"otel.sdk.metric_reader.collection.duration",
			},
			expectedMainMetrics: 2,
		},
		{
			name:                "self instrumentation with multiple metrics",
			enableObservability: true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				counter, err := meter.Int64Counter("test_counter")
				require.NoError(t, err)
				counter.Add(ctx, 5, otelmetric.WithAttributes(attribute.String("type", "requests")))
				counter.Add(ctx, 3, otelmetric.WithAttributes(attribute.String("type", "errors")))

				gauge, err := meter.Float64Gauge("test_gauge")
				require.NoError(t, err)
				gauge.Record(ctx, 100.0, otelmetric.WithAttributes(attribute.String("status", "active")))

				histogram, err := meter.Float64Histogram("test_histogram")
				require.NoError(t, err)
				histogram.Record(ctx, 0.1)
				histogram.Record(ctx, 0.2)
				histogram.Record(ctx, 0.3)
			},
			expectedObservMetrics: []string{
				"otel.sdk.exporter.metric_data_point.inflight",
				"otel.sdk.exporter.metric_data_point.exported",
				"otel.sdk.exporter.operation.duration",
				"otel.sdk.metric_reader.collection.duration",
			},
			expectedMainMetrics: 4, // 3 test metrics + target_info
			checkMetrics: func(t *testing.T, _ []*dto.MetricFamily, observMetrics metricdata.ScopeMetrics) {
				// Check that exported metrics track multiple data points
				for _, m := range observMetrics.Metrics {
					if m.Name == "otel.sdk.exporter.metric_data_point.exported" {
						sum, ok := m.Data.(metricdata.Sum[int64])
						require.True(t, ok)
						require.Len(t, sum.DataPoints, 1)
						// Counter: 2 data points, Gauge: 1 data point, Histogram: 1 data point = 4 total
						require.Equal(t, int64(4), sum.DataPoints[0].Value)
					}
				}
			},
		},
		{
			name:                "self instrumentation enabled with up-down counter",
			enableObservability: true,
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				upDownCounter, err := meter.Int64UpDownCounter(
					"test_updown_counter",
					otelmetric.WithDescription("test up-down counter"),
				)
				require.NoError(t, err)
				upDownCounter.Add(ctx, 10, otelmetric.WithAttributes(attribute.String("direction", "up")))
				upDownCounter.Add(ctx, -5, otelmetric.WithAttributes(attribute.String("direction", "down")))
			},
			expectedObservMetrics: []string{
				"otel.sdk.exporter.metric_data_point.inflight",
				"otel.sdk.exporter.metric_data_point.exported",
				"otel.sdk.exporter.operation.duration",
				"otel.sdk.metric_reader.collection.duration",
			},
			expectedMainMetrics: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.enableObservability {
				t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
			} else {
				t.Setenv("OTEL_GO_X_OBSERVABILITY", "")
			}

			// Setup observability metric collection
			var observReader *metric.ManualReader
			var observMetricsFunc func() metricdata.ScopeMetrics

			if tc.enableObservability {
				originalMP := otel.GetMeterProvider()
				defer otel.SetMeterProvider(originalMP)

				observReader = metric.NewManualReader()
				observMP := metric.NewMeterProvider(metric.WithReader(observReader))
				otel.SetMeterProvider(observMP)

				observMetricsFunc = func() metricdata.ScopeMetrics {
					var rm metricdata.ResourceMetrics
					err := observReader.Collect(t.Context(), &rm)
					require.NoError(t, err)

					// Find the Prometheus exporter observability scope specifically.
					for _, sm := range rm.ScopeMetrics {
						if sm.Scope.Name == observ.ScopeName {
							return sm
						}
					}
					// Exporter scope not found (e.g., if disabled or no scrape yet).
					return metricdata.ScopeMetrics{}
				}
			}

			ctx := t.Context()
			registry := prometheus.NewRegistry()

			exporter, err := New(WithRegisterer(registry))
			require.NoError(t, err)

			provider := metric.NewMeterProvider(metric.WithReader(exporter))
			meter := provider.Meter("test", otelmetric.WithInstrumentationVersion("v1.0.0"))

			// Record the test metrics
			tc.recordMetrics(ctx, meter)

			// Collect main metrics
			mainMetrics, err := registry.Gather()
			require.NoError(t, err)

			// Verify the number of main metric families
			assert.Len(t, mainMetrics, tc.expectedMainMetrics)

			// Collect and check observability metrics if enabled
			if tc.enableObservability {
				observMetrics := observMetricsFunc()

				// Check that expected observability metrics are present
				observedMetrics := make(map[string]bool)
				for _, m := range observMetrics.Metrics {
					observedMetrics[m.Name] = true
				}

				for _, expectedMetric := range tc.expectedObservMetrics {
					assert.True(
						t,
						observedMetrics[expectedMetric],
						"Expected observability metric %s not found",
						expectedMetric,
					)
				}

				// Verify observability metrics have expected structure
				expectedScope := instrumentation.Scope{
					Name:      observ.ScopeName,
					Version:   observ.Version,
					SchemaURL: observ.SchemaURL,
				}
				assert.Equal(t, expectedScope, observMetrics.Scope, "Expected observability scope")
				assert.Len(
					t,
					observMetrics.Metrics,
					len(tc.expectedObservMetrics),
					"Expected number of observability metrics",
				)

				// Run custom metric checks if provided
				if tc.checkMetrics != nil {
					tc.checkMetrics(t, mainMetrics, observMetrics)
				}
			}
		})
	}
}

func TestExporterSelfInstrumentationErrors(t *testing.T) {
	testCases := []struct {
		name                 string
		setupError           func() (metric.Reader, func())
		expectedMinMetrics   int // Minimum expected metrics in error scenarios
		checkErrorAttributes func(t *testing.T, observMetrics metricdata.ScopeMetrics)
	}{
		{
			name: "reader shutdown error",
			setupError: func() (metric.Reader, func()) {
				reader := metric.NewManualReader()
				//nolint:usetesting // required to avoid getting a canceled context at cleanup.
				return reader, func() { _ = reader.Shutdown(context.Background()) }
			},
			expectedMinMetrics: 1, // At least some metrics should be present
		},
		{
			name: "reader not registered error",
			setupError: func() (metric.Reader, func()) {
				reader := metric.NewManualReader()
				// Don't register the reader with a provider
				return reader, func() {}
			},
			expectedMinMetrics: 1, // At least some metrics should be present
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Enable observability
			t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

			// Setup observability metric collection
			originalMP := otel.GetMeterProvider()
			defer otel.SetMeterProvider(originalMP)

			observReader := metric.NewManualReader()
			observMP := metric.NewMeterProvider(metric.WithReader(observReader))
			otel.SetMeterProvider(observMP)

			registry := prometheus.NewRegistry()

			reader, cleanup := tc.setupError()
			defer cleanup()

			// Create exporter with the error-prone reader
			cfg := newConfig()
			cfg.registerer = registry

			collector := &collector{
				reader:                   reader,
				disableTargetInfo:        cfg.disableTargetInfo,
				withoutUnits:             cfg.withoutUnits,
				withoutCounterSuffixes:   cfg.withoutCounterSuffixes,
				disableScopeInfo:         cfg.disableScopeInfo,
				metricFamilies:           make(map[string]*dto.MetricFamily),
				namespace:                cfg.namespace,
				resourceAttributesFilter: cfg.resourceAttributesFilter,
				metricNamer:              otlptranslator.NewMetricNamer(cfg.namespace, cfg.translationStrategy),
			}

			var err error
			collector.inst, err = observ.NewInstrumentation(0)
			require.NoError(t, err)

			err = registry.Register(collector)
			require.NoError(t, err)

			// Trigger collection which should encounter the error
			_, err = registry.Gather()
			require.NoError(t, err)

			// Collect observability metrics
			var observMetrics metricdata.ResourceMetrics
			err = observReader.Collect(t.Context(), &observMetrics)
			require.NoError(t, err)

			if len(observMetrics.ScopeMetrics) > 0 {
				// Verify observability metrics are still present (at least some)
				scopeMetrics := observMetrics.ScopeMetrics[0]
				foundObservMetrics := 0
				for _, m := range scopeMetrics.Metrics {
					switch m.Name {
					case "otel.sdk.exporter.metric_data_point.inflight",
						"otel.sdk.exporter.metric_data_point.exported",
						"otel.sdk.exporter.operation.duration",
						"otel.sdk.metric_reader.collection.duration":
						foundObservMetrics++
					}
				}
				assert.GreaterOrEqual(
					t,
					foundObservMetrics,
					tc.expectedMinMetrics,
					"Should have at least some observability metrics even with errors",
				)

				if tc.checkErrorAttributes != nil {
					tc.checkErrorAttributes(t, scopeMetrics)
				}
			}
		})
	}
}

func TestExporterSelfInstrumentationConcurrency(t *testing.T) {
	// Enable observability
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Setup observability metric collection
	originalMP := otel.GetMeterProvider()
	defer otel.SetMeterProvider(originalMP)

	observReader := metric.NewManualReader()
	observMP := metric.NewMeterProvider(metric.WithReader(observReader))
	otel.SetMeterProvider(observMP)

	ctx := t.Context()
	registry := prometheus.NewRegistry()

	exporter, err := New(WithRegisterer(registry))
	require.NoError(t, err)

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("concurrent_test", otelmetric.WithInstrumentationVersion("v1.0.0"))

	counter, err := meter.Int64Counter("concurrent_counter")
	require.NoError(t, err)

	// Run concurrent operations
	const numGoroutines = 10
	const numOperations = 100
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				counter.Add(ctx, 1, otelmetric.WithAttributes(attribute.Int("goroutine", id)))

				// Occasionally trigger collection
				if j%10 == 0 {
					_, _ = registry.Gather()
				}
			}
		}(i)
	}

	wg.Wait()

	// Final collection
	_, err = registry.Gather()
	require.NoError(t, err)

	// Collect observability metrics
	var observMetrics metricdata.ResourceMetrics
	err = observReader.Collect(t.Context(), &observMetrics)
	require.NoError(t, err)

	if len(observMetrics.ScopeMetrics) > 0 {
		scopeMetrics := observMetrics.ScopeMetrics[0]
		// Verify observability metrics are present and have reasonable values
		for _, m := range scopeMetrics.Metrics {
			switch m.Name {
			case "otel.sdk.exporter.metric_data_point.exported":
				sum, ok := m.Data.(metricdata.Sum[int64])
				require.True(t, ok)
				require.NotEmpty(t, sum.DataPoints)
				// Should have exported many data points
				assert.Positive(t, sum.DataPoints[0].Value)
			case "otel.sdk.exporter.operation.duration":
				hist, ok := m.Data.(metricdata.Histogram[float64])
				require.True(t, ok)
				require.NotEmpty(t, hist.DataPoints)
				// Should have recorded operation durations
				assert.Positive(t, hist.DataPoints[0].Count)
			case "otel.sdk.metric_reader.collection.duration":
				hist, ok := m.Data.(metricdata.Histogram[float64])
				require.True(t, ok)
				require.NotEmpty(t, hist.DataPoints)
				// Should have recorded collection durations
				assert.Positive(t, hist.DataPoints[0].Count)
			}
		}
	}
}

func TestExporterSelfInstrumentationExemplarHandling(t *testing.T) {
	// Enable observability
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Setup observability metric collection
	originalMP := otel.GetMeterProvider()
	defer otel.SetMeterProvider(originalMP)

	observReader := metric.NewManualReader()
	observMP := metric.NewMeterProvider(metric.WithReader(observReader))
	otel.SetMeterProvider(observMP)

	ctx := t.Context()
	registry := prometheus.NewRegistry()

	exporter, err := New(WithRegisterer(registry))
	require.NoError(t, err)

	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("exemplar_test", otelmetric.WithInstrumentationVersion("v1.0.0"))

	// Create trace context for exemplars
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:     [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
		TraceFlags: trace.FlagsSampled,
	})
	ctx = trace.ContextWithSpanContext(ctx, sc)

	histogram, err := meter.Float64Histogram("test_histogram_with_exemplars")
	require.NoError(t, err)

	// Record values that should generate exemplars
	histogram.Record(ctx, 1.0, otelmetric.WithAttributes(attribute.String("key", "value1")))
	histogram.Record(ctx, 2.0, otelmetric.WithAttributes(attribute.String("key", "value2")))

	// Collect metrics
	got, err := registry.Gather()
	require.NoError(t, err)

	// Verify that metrics are collected without errors even when exemplars are present
	foundTestHistogram := false

	for _, mf := range got {
		if *mf.Name == "test_histogram_with_exemplars" {
			foundTestHistogram = true
			assert.Equal(t, dto.MetricType_HISTOGRAM, *mf.Type)
		}
	}

	assert.True(t, foundTestHistogram, "Test histogram should be present")

	// Collect observability metrics
	var observMetrics metricdata.ResourceMetrics
	err = observReader.Collect(t.Context(), &observMetrics)
	require.NoError(t, err)

	if len(observMetrics.ScopeMetrics) > 0 {
		expectedMetrics := map[string]bool{
			"otel.sdk.exporter.metric_data_point.inflight": false,
			"otel.sdk.exporter.metric_data_point.exported": false,
			"otel.sdk.exporter.operation.duration":         false,
			"otel.sdk.metric_reader.collection.duration":   false,
		}

		// Check all scope metrics, not just the first one
		for _, scopeMetrics := range observMetrics.ScopeMetrics {
			for _, m := range scopeMetrics.Metrics {
				if _, exists := expectedMetrics[m.Name]; exists {
					expectedMetrics[m.Name] = true
				}
			}
		}

		foundObservabilityMetrics := 0
		for _, found := range expectedMetrics {
			if found {
				foundObservabilityMetrics++
			}
		}
		assert.Equal(t, 4, foundObservabilityMetrics, "All observability metrics should be present")
	}
}

func TestExporterSelfInstrumentationInitErrors(t *testing.T) {
	// Test when NewInstrumentation returns an error
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Set up a meter provider that will cause NewInstrumentation to fail
	original := otel.GetMeterProvider()
	defer otel.SetMeterProvider(original)

	errMP := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(errMP)

	registry := prometheus.NewRegistry()

	// This should fail during exporter creation due to instrumentation init error
	_, err := New(WithRegisterer(registry))
	require.ErrorIs(t, err, assert.AnError, "Expected NewInstrumentation error to be propagated")
}

// Helper types for testing NewInstrumentation errors.
type errMeterProvider struct {
	otelmetric.MeterProvider
	err error
}

func (m *errMeterProvider) Meter(string, ...otelmetric.MeterOption) otelmetric.Meter {
	return &errMeter{err: m.err}
}

type errMeter struct {
	otelmetric.Meter
	err error
}

func (m *errMeter) Int64UpDownCounter(
	string,
	...otelmetric.Int64UpDownCounterOption,
) (otelmetric.Int64UpDownCounter, error) {
	return nil, m.err
}

func (m *errMeter) Int64Counter(string, ...otelmetric.Int64CounterOption) (otelmetric.Int64Counter, error) {
	return nil, m.err
}

func (m *errMeter) Float64Histogram(string, ...otelmetric.Float64HistogramOption) (otelmetric.Float64Histogram, error) {
	return nil, m.err
}
