// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"context"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/errhand"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func TestPrometheusExporter(t *testing.T) {
	eh := errhand.NewGlobal()

	testCases := []struct {
		name               string
		emptyResource      bool
		customResouceAttrs []attribute.KeyValue
		recordMetrics      func(ctx context.Context, meter otelmetric.Meter)
		options            []Option
		expectedFile       string
	}{
		{
			name:         "counter",
			expectedFile: "testdata/counter.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				}
				counter := meter.Float64Counter(
					"foo",
					instrument.WithDescription("a simple counter"),
					instrument.WithUnit(unit.Milliseconds),
				)
				eh.RequireNoErrors(t)
				counter.Add(ctx, 5, attrs...)
				counter.Add(ctx, 10.3, attrs...)
				counter.Add(ctx, 9, attrs...)

				attrs2 := []attribute.KeyValue{
					attribute.Key("A").String("D"),
					attribute.Key("C").String("B"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				}
				counter.Add(ctx, 5, attrs2...)
			},
		},
		{
			name:         "gauge",
			expectedFile: "testdata/gauge.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				}
				gauge := meter.Float64UpDownCounter(
					"bar",
					instrument.WithDescription("a fun little gauge"),
					instrument.WithUnit(unit.Dimensionless),
				)
				eh.RequireNoErrors(t)
				gauge.Add(ctx, 1.0, attrs...)
				gauge.Add(ctx, -.25, attrs...)
			},
		},
		{
			name:         "histogram",
			expectedFile: "testdata/histogram.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				}
				histogram := meter.Float64Histogram(
					"histogram_baz",
					instrument.WithDescription("a very nice histogram"),
					instrument.WithUnit(unit.Bytes),
				)
				eh.RequireNoErrors(t)
				histogram.Record(ctx, 23, attrs...)
				histogram.Record(ctx, 7, attrs...)
				histogram.Record(ctx, 101, attrs...)
				histogram.Record(ctx, 105, attrs...)
			},
		},
		{
			name:         "sanitized attributes to labels",
			expectedFile: "testdata/sanitized_labels.txt",
			options:      []Option{WithoutUnits()},
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					// exact match, value should be overwritten
					attribute.Key("A.B").String("X"),
					attribute.Key("A.B").String("Q"),

					// unintended match due to sanitization, values should be concatenated
					attribute.Key("C.D").String("Y"),
					attribute.Key("C/D").String("Z"),
				}
				counter := meter.Float64Counter(
					"foo",
					instrument.WithDescription("a sanitary counter"),
					// This unit is not added to
					instrument.WithUnit(unit.Bytes),
				)
				eh.RequireNoErrors(t)
				counter.Add(ctx, 5, attrs...)
				counter.Add(ctx, 10.3, attrs...)
				counter.Add(ctx, 9, attrs...)
			},
		},
		{
			name:         "invalid instruments are renamed",
			expectedFile: "testdata/sanitized_names.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				}
				// Valid.
				gauge := meter.Float64UpDownCounter("bar", instrument.WithDescription("a fun little gauge"))
				eh.RequireNoErrors(t)
				gauge.Add(ctx, 100, attrs...)
				gauge.Add(ctx, -25, attrs...)

				// Invalid, will be renamed.
				gauge = meter.Float64UpDownCounter("invalid.gauge.name", instrument.WithDescription("a gauge with an invalid name"))
				eh.RequireNoErrors(t)
				gauge.Add(ctx, 100, attrs...)

				counter := meter.Float64Counter("0invalid.counter.name", instrument.WithDescription("a counter with an invalid name"))
				eh.RequireNoErrors(t)
				counter.Add(ctx, 100, attrs...)

				histogram := meter.Float64Histogram("invalid.hist.name", instrument.WithDescription("a histogram with an invalid name"))
				eh.RequireNoErrors(t)
				histogram.Record(ctx, 23, attrs...)
			},
		},
		{
			name:          "empty resource",
			emptyResource: true,
			expectedFile:  "testdata/empty_resource.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				}
				counter := meter.Float64Counter("foo", instrument.WithDescription("a simple counter"))
				eh.RequireNoErrors(t)
				counter.Add(ctx, 5, attrs...)
				counter.Add(ctx, 10.3, attrs...)
				counter.Add(ctx, 9, attrs...)
			},
		},
		{
			name: "custom resource",
			customResouceAttrs: []attribute.KeyValue{
				attribute.Key("A").String("B"),
				attribute.Key("C").String("D"),
			},
			expectedFile: "testdata/custom_resource.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				}
				counter := meter.Float64Counter("foo", instrument.WithDescription("a simple counter"))
				eh.RequireNoErrors(t)
				counter.Add(ctx, 5, attrs...)
				counter.Add(ctx, 10.3, attrs...)
				counter.Add(ctx, 9, attrs...)
			},
		},
		{
			name:         "without target_info",
			options:      []Option{WithoutTargetInfo()},
			expectedFile: "testdata/without_target_info.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
					attribute.Key("E").Bool(true),
					attribute.Key("F").Int(42),
				}
				counter := meter.Float64Counter("foo", instrument.WithDescription("a simple counter"))
				eh.RequireNoErrors(t)
				counter.Add(ctx, 5, attrs...)
				counter.Add(ctx, 10.3, attrs...)
				counter.Add(ctx, 9, attrs...)
			},
		},
		{
			name:         "without scope_info",
			options:      []Option{WithoutScopeInfo()},
			expectedFile: "testdata/without_scope_info.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				}
				gauge := meter.Int64UpDownCounter(
					"bar",
					instrument.WithDescription("a fun little gauge"),
					instrument.WithUnit(unit.Dimensionless),
				)
				eh.RequireNoErrors(t)
				gauge.Add(ctx, 2, attrs...)
				gauge.Add(ctx, -1, attrs...)
			},
		},
		{
			name:         "without scope_info and target_info",
			options:      []Option{WithoutScopeInfo(), WithoutTargetInfo()},
			expectedFile: "testdata/without_scope_and_target_info.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					attribute.Key("A").String("B"),
					attribute.Key("C").String("D"),
				}
				counter := meter.Int64Counter(
					"bar",
					instrument.WithDescription("a fun little counter"),
					instrument.WithUnit(unit.Bytes),
				)
				eh.RequireNoErrors(t)
				counter.Add(ctx, 2, attrs...)
				counter.Add(ctx, 1, attrs...)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
					resource.WithAttributes(tc.customResouceAttrs...),
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
					metric.Stream{Aggregation: aggregation.ExplicitBucketHistogram{
						Boundaries: []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 1000},
					}},
				)),
			)
			meter := provider.Meter(
				"testmeter",
				otelmetric.WithInstrumentationVersion("v0.1.0"),
			)

			tc.recordMetrics(ctx, meter)

			file, err := os.Open(tc.expectedFile)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, file.Close()) })

			err = testutil.GatherAndCompare(registry, file)
			require.NoError(t, err)
		})
	}
}

func TestSantitizeName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"nam€_with_3_width_rune", "nam__with_3_width_rune"},
		{"`", "_"},
		{
			`! "#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWKYZ[]\^_abcdefghijklmnopqrstuvwkyz{|}~`,
			`________________0123456789:______ABCDEFGHIJKLMNOPQRSTUVWKYZ_____abcdefghijklmnopqrstuvwkyz____`,
		},

		// Test cases taken from
		// https://github.com/prometheus/common/blob/dfbc25bd00225c70aca0d94c3c4bb7744f28ace0/model/metric_test.go#L85-L136
		{"Avalid_23name", "Avalid_23name"},
		{"_Avalid_23name", "_Avalid_23name"},
		{"1valid_23name", "_1valid_23name"},
		{"avalid_23name", "avalid_23name"},
		{"Ava:lid_23name", "Ava:lid_23name"},
		{"a lid_23name", "a_lid_23name"},
		{":leading_colon", ":leading_colon"},
		{"colon:in:the:middle", "colon:in:the:middle"},
		{"", ""},
	}

	for _, test := range tests {
		require.Equalf(t, test.want, sanitizeName(test.input), "input: %q", test.input)
	}
}

func TestMultiScopes(t *testing.T) {
	eh := errhand.NewGlobal()

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

	fooCounter := provider.Meter("meterfoo", otelmetric.WithInstrumentationVersion("v0.1.0")).
		Int64Counter(
			"foo",
			instrument.WithUnit(unit.Milliseconds),
			instrument.WithDescription("meter foo counter"))
	eh.AssertNoErrors(t)
	fooCounter.Add(ctx, 100, attribute.String("type", "foo"))

	barCounter := provider.Meter("meterbar", otelmetric.WithInstrumentationVersion("v0.1.0")).
		Int64Counter(
			"bar",
			instrument.WithUnit(unit.Milliseconds),
			instrument.WithDescription("meter bar counter"))
	eh.AssertNoErrors(t)
	barCounter.Add(ctx, 200, attribute.String("type", "bar"))

	file, err := os.Open("testdata/multi_scopes.txt")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, file.Close()) })

	err = testutil.GatherAndCompare(registry, file)
	require.NoError(t, err)
}

func TestDuplicateMetrics(t *testing.T) {
	eh := errhand.NewGlobal()

	testCases := []struct {
		name                  string
		customResouceAttrs    []attribute.KeyValue
		recordMetrics         func(ctx context.Context, meterA, meterB otelmetric.Meter)
		options               []Option
		possibleExpectedFiles []string
	}{
		{
			name: "no_conflict_two_counters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				fooA := meterA.Int64Counter("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter counter foo"))
				eh.AssertNoErrors(t)
				fooA.Add(ctx, 100, attribute.String("A", "B"))

				fooB := meterB.Int64Counter("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter counter foo"))
				eh.AssertNoErrors(t)
				fooB.Add(ctx, 100, attribute.String("A", "B"))
			},
			possibleExpectedFiles: []string{"testdata/no_conflict_two_counters.txt"},
		},
		{
			name: "no_conflict_two_updowncounters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				fooA := meterA.Int64UpDownCounter("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter gauge foo"))
				eh.AssertNoErrors(t)
				fooA.Add(ctx, 100, attribute.String("A", "B"))

				fooB := meterB.Int64UpDownCounter("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter gauge foo"))
				eh.AssertNoErrors(t)
				fooB.Add(ctx, 100, attribute.String("A", "B"))
			},
			possibleExpectedFiles: []string{"testdata/no_conflict_two_updowncounters.txt"},
		},
		{
			name: "no_conflict_two_histograms",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				fooA := meterA.Int64Histogram("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter histogram foo"))
				eh.AssertNoErrors(t)
				fooA.Record(ctx, 100, attribute.String("A", "B"))

				fooB := meterB.Int64Histogram("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter histogram foo"))
				eh.AssertNoErrors(t)
				fooB.Record(ctx, 100, attribute.String("A", "B"))
			},
			possibleExpectedFiles: []string{"testdata/no_conflict_two_histograms.txt"},
		},
		{
			name: "conflict_help_two_counters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA := meterA.Int64Counter("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter a bar"))
				eh.AssertNoErrors(t)
				barA.Add(ctx, 100, attribute.String("type", "bar"))

				barB := meterB.Int64Counter("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter b bar"))
				eh.AssertNoErrors(t)
				barB.Add(ctx, 100, attribute.String("type", "bar"))
			},
			possibleExpectedFiles: []string{
				"testdata/conflict_help_two_counters_1.txt",
				"testdata/conflict_help_two_counters_2.txt",
			},
		},
		{
			name: "conflict_help_two_updowncounters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA := meterA.Int64UpDownCounter("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter a bar"))
				eh.AssertNoErrors(t)
				barA.Add(ctx, 100, attribute.String("type", "bar"))

				barB := meterB.Int64UpDownCounter("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter b bar"))
				eh.AssertNoErrors(t)
				barB.Add(ctx, 100, attribute.String("type", "bar"))
			},
			possibleExpectedFiles: []string{
				"testdata/conflict_help_two_updowncounters_1.txt",
				"testdata/conflict_help_two_updowncounters_2.txt",
			},
		},
		{
			name: "conflict_help_two_histograms",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA := meterA.Int64Histogram("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter a bar"))
				eh.AssertNoErrors(t)
				barA.Record(ctx, 100, attribute.String("A", "B"))

				barB := meterB.Int64Histogram("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter b bar"))
				eh.AssertNoErrors(t)
				barB.Record(ctx, 100, attribute.String("A", "B"))
			},
			possibleExpectedFiles: []string{
				"testdata/conflict_help_two_histograms_1.txt",
				"testdata/conflict_help_two_histograms_2.txt",
			},
		},
		{
			name: "conflict_unit_two_counters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				bazA := meterA.Int64Counter("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter bar"))
				eh.AssertNoErrors(t)
				bazA.Add(ctx, 100, attribute.String("type", "bar"))

				bazB := meterB.Int64Counter("bar",
					instrument.WithUnit(unit.Milliseconds),
					instrument.WithDescription("meter bar"))
				eh.AssertNoErrors(t)
				bazB.Add(ctx, 100, attribute.String("type", "bar"))
			},
			options:               []Option{WithoutUnits()},
			possibleExpectedFiles: []string{"testdata/conflict_unit_two_counters.txt"},
		},
		{
			name: "conflict_unit_two_updowncounters",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA := meterA.Int64UpDownCounter("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter gauge bar"))
				eh.AssertNoErrors(t)
				barA.Add(ctx, 100, attribute.String("type", "bar"))

				barB := meterB.Int64UpDownCounter("bar",
					instrument.WithUnit(unit.Milliseconds),
					instrument.WithDescription("meter gauge bar"))
				eh.AssertNoErrors(t)
				barB.Add(ctx, 100, attribute.String("type", "bar"))
			},
			options:               []Option{WithoutUnits()},
			possibleExpectedFiles: []string{"testdata/conflict_unit_two_updowncounters.txt"},
		},
		{
			name: "conflict_unit_two_histograms",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				barA := meterA.Int64Histogram("bar",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter histogram bar"))
				eh.AssertNoErrors(t)
				barA.Record(ctx, 100, attribute.String("A", "B"))

				barB := meterB.Int64Histogram("bar",
					instrument.WithUnit(unit.Milliseconds),
					instrument.WithDescription("meter histogram bar"))
				eh.AssertNoErrors(t)
				barB.Record(ctx, 100, attribute.String("A", "B"))
			},
			options:               []Option{WithoutUnits()},
			possibleExpectedFiles: []string{"testdata/conflict_unit_two_histograms.txt"},
		},
		{
			name: "conflict_type_counter_and_updowncounter",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				counter := meterA.Int64Counter("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter foo"))
				eh.AssertNoErrors(t)
				counter.Add(ctx, 100, attribute.String("type", "foo"))

				gauge := meterA.Int64UpDownCounter("foo_total",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter foo"))
				eh.AssertNoErrors(t)
				gauge.Add(ctx, 200, attribute.String("type", "foo"))
			},
			options: []Option{WithoutUnits()},
			possibleExpectedFiles: []string{
				"testdata/conflict_type_counter_and_updowncounter_1.txt",
				"testdata/conflict_type_counter_and_updowncounter_2.txt",
			},
		},
		{
			name: "conflict_type_histogram_and_updowncounter",
			recordMetrics: func(ctx context.Context, meterA, meterB otelmetric.Meter) {
				fooA := meterA.Int64UpDownCounter("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter gauge foo"))
				eh.AssertNoErrors(t)
				fooA.Add(ctx, 100, attribute.String("A", "B"))

				fooHistogramA := meterA.Int64Histogram("foo",
					instrument.WithUnit(unit.Bytes),
					instrument.WithDescription("meter histogram foo"))
				eh.AssertNoErrors(t)
				fooHistogramA.Record(ctx, 100, attribute.String("A", "B"))
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

			var match = false
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
