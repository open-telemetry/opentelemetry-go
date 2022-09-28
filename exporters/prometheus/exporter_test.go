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
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"
)

func TestPrometheusExporter(t *testing.T) {
	testCases := []struct {
		name          string
		recordMetrics func(ctx context.Context, meter otelmetric.Meter)
		expectedFile  string
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
				counter, err := meter.SyncFloat64().Counter("foo", instrument.WithDescription("a simple counter"))
				require.NoError(t, err)
				counter.Add(ctx, 5, attrs...)
				counter.Add(ctx, 10.3, attrs...)
				counter.Add(ctx, 9, attrs...)
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
				gauge, err := meter.SyncFloat64().UpDownCounter("bar", instrument.WithDescription("a fun little gauge"))
				require.NoError(t, err)
				gauge.Add(ctx, 100, attrs...)
				gauge.Add(ctx, -25, attrs...)
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
				histogram, err := meter.SyncFloat64().Histogram("baz", instrument.WithDescription("a very nice histogram"))
				require.NoError(t, err)
				histogram.Record(ctx, 23, attrs...)
				histogram.Record(ctx, 7, attrs...)
				histogram.Record(ctx, 101, attrs...)
				histogram.Record(ctx, 105, attrs...)
			},
		},
		{
			name:         "sanitized attributes to labels",
			expectedFile: "testdata/sanitized_labels.txt",
			recordMetrics: func(ctx context.Context, meter otelmetric.Meter) {
				attrs := []attribute.KeyValue{
					// exact match, value should be overwritten
					attribute.Key("A.B").String("X"),
					attribute.Key("A.B").String("Q"),

					// unintended match due to sanitization, values should be concatenated
					attribute.Key("C.D").String("Y"),
					attribute.Key("C/D").String("Z"),
				}
				counter, err := meter.SyncFloat64().Counter("foo", instrument.WithDescription("a sanitary counter"))
				require.NoError(t, err)
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
				gauge, err := meter.SyncFloat64().UpDownCounter("bar", instrument.WithDescription("a fun little gauge"))
				require.NoError(t, err)
				gauge.Add(ctx, 100, attrs...)
				gauge.Add(ctx, -25, attrs...)

				// Invalid, will be renamed.
				gauge, err = meter.SyncFloat64().UpDownCounter("invalid.gauge.name", instrument.WithDescription("a gauge with an invalid name"))
				require.NoError(t, err)
				gauge.Add(ctx, 100, attrs...)

				counter, err := meter.SyncFloat64().Counter("0invalid.counter.name", instrument.WithDescription("a counter with an invalid name"))
				require.NoError(t, err)
				counter.Add(ctx, 100, attrs...)

				histogram, err := meter.SyncFloat64().Histogram("invalid.hist.name", instrument.WithDescription("a histogram with an invalid name"))
				require.NoError(t, err)
				histogram.Record(ctx, 23, attrs...)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			registry := prometheus.NewRegistry() // This is the default behavior, this is used to manually gather.

			exporter, err := New(WithRegistry(registry))
			require.NoError(t, err)
			provider := metric.NewMeterProvider(metric.WithReader(exporter))
			meter := provider.Meter("testmeter")

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
