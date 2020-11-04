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

package stdout_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/export/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testFixture struct {
	t        *testing.T
	ctx      context.Context
	exporter *stdout.Exporter
	output   *bytes.Buffer
}

var testResource = resource.NewWithAttributes(label.String("R", "V"))

func newFixture(t *testing.T, opts ...stdout.Option) testFixture {
	buf := &bytes.Buffer{}
	opts = append(opts, stdout.WithWriter(buf))
	opts = append(opts, stdout.WithoutTimestamps())
	exp, err := stdout.NewExporter(opts...)
	if err != nil {
		t.Fatal("Error building fixture: ", err)
	}
	return testFixture{
		t:        t,
		ctx:      context.Background(),
		exporter: exp,
		output:   buf,
	}
}

func (fix testFixture) Output() string {
	return strings.TrimSpace(fix.output.String())
}

func (fix testFixture) Export(checkpointSet export.CheckpointSet) {
	err := fix.exporter.Export(fix.ctx, checkpointSet)
	if err != nil {
		fix.t.Error("export failed: ", err)
	}
}

func TestStdoutInvalidQuantile(t *testing.T) {
	_, err := stdout.NewExporter(
		stdout.WithQuantiles([]float64{1.1, 0.9}),
	)
	require.Error(t, err, "Invalid quantile error expected")
	require.Equal(t, aggregation.ErrInvalidQuantile, err)
}

func TestStdoutTimestamp(t *testing.T) {
	var buf bytes.Buffer
	exporter, err := stdout.NewExporter(
		stdout.WithWriter(&buf),
	)
	if err != nil {
		t.Fatal("Invalid config: ", err)
	}

	before := time.Now()

	checkpointSet := metrictest.NewCheckpointSet(testResource)

	ctx := context.Background()
	desc := otel.NewDescriptor("test.name", otel.ValueObserverInstrumentKind, number.Int64Kind)

	lvagg, ckpt := metrictest.Unslice2(lastvalue.New(2))

	aggregatortest.CheckedUpdate(t, lvagg, number.NewInt64Number(321), &desc)
	require.NoError(t, lvagg.SynchronizedMove(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt)

	if err := exporter.Export(ctx, checkpointSet); err != nil {
		t.Fatal("Unexpected export error: ", err)
	}

	after := time.Now()

	var printed []interface{}
	if err := json.Unmarshal(buf.Bytes(), &printed); err != nil {
		t.Fatal("JSON parse error: ", err)
	}

	require.Len(t, printed, 1)
	lastValue, ok := printed[0].(map[string]interface{})
	require.True(t, ok, "last value format")
	require.Contains(t, lastValue, "Timestamp")
	lastValueTS := lastValue["Timestamp"].(string)
	lastValueTimestamp, err := time.Parse(time.RFC3339Nano, lastValueTS)
	if err != nil {
		t.Fatal("JSON parse error: ", lastValueTS, ": ", err)
	}

	assert.True(t, lastValueTimestamp.After(before))
	assert.True(t, lastValueTimestamp.Before(after))
}

func TestStdoutCounterFormat(t *testing.T) {
	fix := newFixture(t)

	checkpointSet := metrictest.NewCheckpointSet(testResource)

	desc := otel.NewDescriptor("test.name", otel.CounterInstrumentKind, number.Int64Kind)

	cagg, ckpt := metrictest.Unslice2(sum.New(2))

	aggregatortest.CheckedUpdate(fix.t, cagg, number.NewInt64Number(123), &desc)
	require.NoError(t, cagg.SynchronizedMove(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, label.String("A", "B"), label.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `[{"Name":"test.name{R=V,A=B,C=D}","Sum":123}]`, fix.Output())
}

func TestStdoutLastValueFormat(t *testing.T) {
	fix := newFixture(t)

	checkpointSet := metrictest.NewCheckpointSet(testResource)

	desc := otel.NewDescriptor("test.name", otel.ValueObserverInstrumentKind, number.Float64Kind)
	lvagg, ckpt := metrictest.Unslice2(lastvalue.New(2))

	aggregatortest.CheckedUpdate(fix.t, lvagg, number.NewFloat64Number(123.456), &desc)
	require.NoError(t, lvagg.SynchronizedMove(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, label.String("A", "B"), label.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `[{"Name":"test.name{R=V,A=B,C=D}","Last":123.456}]`, fix.Output())
}

func TestStdoutMinMaxSumCount(t *testing.T) {
	fix := newFixture(t)

	checkpointSet := metrictest.NewCheckpointSet(testResource)

	desc := otel.NewDescriptor("test.name", otel.ValueRecorderInstrumentKind, number.Float64Kind)

	magg, ckpt := metrictest.Unslice2(minmaxsumcount.New(2, &desc))

	aggregatortest.CheckedUpdate(fix.t, magg, number.NewFloat64Number(123.456), &desc)
	aggregatortest.CheckedUpdate(fix.t, magg, number.NewFloat64Number(876.543), &desc)
	require.NoError(t, magg.SynchronizedMove(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, label.String("A", "B"), label.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `[{"Name":"test.name{R=V,A=B,C=D}","Min":123.456,"Max":876.543,"Sum":999.999,"Count":2}]`, fix.Output())
}

func TestStdoutValueRecorderFormat(t *testing.T) {
	fix := newFixture(t, stdout.WithPrettyPrint())

	checkpointSet := metrictest.NewCheckpointSet(testResource)

	desc := otel.NewDescriptor("test.name", otel.ValueRecorderInstrumentKind, number.Float64Kind)
	aagg, ckpt := metrictest.Unslice2(array.New(2))

	for i := 0; i < 1000; i++ {
		aggregatortest.CheckedUpdate(fix.t, aagg, number.NewFloat64Number(float64(i)+0.5), &desc)
	}

	require.NoError(t, aagg.SynchronizedMove(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, label.String("A", "B"), label.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `[
	{
		"Name": "test.name{R=V,A=B,C=D}",
		"Min": 0.5,
		"Max": 999.5,
		"Sum": 500000,
		"Count": 1000,
		"Quantiles": [
			{
				"Quantile": 0.5,
				"Value": 500.5
			},
			{
				"Quantile": 0.9,
				"Value": 900.5
			},
			{
				"Quantile": 0.99,
				"Value": 990.5
			}
		]
	}
]`, fix.Output())
}

func TestStdoutNoData(t *testing.T) {
	desc := otel.NewDescriptor("test.name", otel.ValueRecorderInstrumentKind, number.Float64Kind)

	runTwoAggs := func(agg, ckpt export.Aggregator) {
		t.Run(fmt.Sprintf("%T", agg), func(t *testing.T) {
			t.Parallel()

			fix := newFixture(t)

			checkpointSet := metrictest.NewCheckpointSet(testResource)

			require.NoError(t, agg.SynchronizedMove(ckpt, &desc))

			checkpointSet.Add(&desc, ckpt)

			fix.Export(checkpointSet)

			require.Equal(t, "", fix.Output())
		})
	}

	runTwoAggs(metrictest.Unslice2(ddsketch.New(2, &desc, ddsketch.NewDefaultConfig())))
	runTwoAggs(metrictest.Unslice2(minmaxsumcount.New(2, &desc)))
}

func TestStdoutLastValueNotSet(t *testing.T) {
	fix := newFixture(t)

	checkpointSet := metrictest.NewCheckpointSet(testResource)

	desc := otel.NewDescriptor("test.name", otel.ValueObserverInstrumentKind, number.Float64Kind)

	lvagg, ckpt := metrictest.Unslice2(lastvalue.New(2))
	require.NoError(t, lvagg.SynchronizedMove(ckpt, &desc))

	checkpointSet.Add(&desc, lvagg, label.String("A", "B"), label.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, "", fix.Output())
}

func TestStdoutResource(t *testing.T) {
	type testCase struct {
		expect string
		res    *resource.Resource
		attrs  []label.KeyValue
	}
	newCase := func(expect string, res *resource.Resource, attrs ...label.KeyValue) testCase {
		return testCase{
			expect: expect,
			res:    res,
			attrs:  attrs,
		}
	}
	testCases := []testCase{
		newCase("R1=V1,R2=V2,A=B,C=D",
			resource.NewWithAttributes(label.String("R1", "V1"), label.String("R2", "V2")),
			label.String("A", "B"),
			label.String("C", "D")),
		newCase("R1=V1,R2=V2",
			resource.NewWithAttributes(label.String("R1", "V1"), label.String("R2", "V2")),
		),
		newCase("A=B,C=D",
			nil,
			label.String("A", "B"),
			label.String("C", "D"),
		),
		// We explicitly do not de-duplicate between resources
		// and metric labels in this exporter.
		newCase("R1=V1,R2=V2,R1=V3,R2=V4",
			resource.NewWithAttributes(label.String("R1", "V1"), label.String("R2", "V2")),
			label.String("R1", "V3"),
			label.String("R2", "V4")),
	}

	for _, tc := range testCases {
		fix := newFixture(t)

		checkpointSet := metrictest.NewCheckpointSet(tc.res)

		desc := otel.NewDescriptor("test.name", otel.ValueObserverInstrumentKind, number.Float64Kind)
		lvagg, ckpt := metrictest.Unslice2(lastvalue.New(2))

		aggregatortest.CheckedUpdate(fix.t, lvagg, number.NewFloat64Number(123.456), &desc)
		require.NoError(t, lvagg.SynchronizedMove(ckpt, &desc))

		checkpointSet.Add(&desc, ckpt, tc.attrs...)

		fix.Export(checkpointSet)

		require.Equal(t, `[{"Name":"test.name{`+tc.expect+`}","Last":123.456}]`, fix.Output())
	}
}
