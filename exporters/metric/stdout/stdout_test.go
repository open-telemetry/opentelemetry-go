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

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/stdout"
	"go.opentelemetry.io/otel/exporters/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	aggtest "go.opentelemetry.io/otel/sdk/metric/aggregator/test"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testFixture struct {
	t        *testing.T
	ctx      context.Context
	exporter *stdout.Exporter
	output   *bytes.Buffer
}

var testResource = resource.New(kv.String("R", "V"))

func newFixture(t *testing.T, config stdout.Config) testFixture {
	buf := &bytes.Buffer{}
	config.Writer = buf
	config.DoNotPrintTime = true
	exp, err := stdout.NewRawExporter(config)
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
	_, err := stdout.NewRawExporter(stdout.Config{
		Quantiles: []float64{1.1, 0.9},
	})
	require.Error(t, err, "Invalid quantile error expected")
	require.Equal(t, aggregation.ErrInvalidQuantile, err)
}

func TestStdoutTimestamp(t *testing.T) {
	var buf bytes.Buffer
	exporter, err := stdout.NewRawExporter(stdout.Config{
		Writer:         &buf,
		DoNotPrintTime: false,
	})
	if err != nil {
		t.Fatal("Invalid config: ", err)
	}

	before := time.Now()

	checkpointSet := test.NewCheckpointSet(testResource)

	ctx := context.Background()
	desc := metric.NewDescriptor("test.name", metric.ValueObserverKind, metric.Int64NumberKind)

	lvagg, ckpt := test.Unslice2(lastvalue.New(2))

	aggtest.CheckedUpdate(t, lvagg, metric.NewInt64Number(321), &desc)
	require.NoError(t, lvagg.SynchronizedCopy(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt)

	if err := exporter.Export(ctx, checkpointSet); err != nil {
		t.Fatal("Unexpected export error: ", err)
	}

	after := time.Now()

	var printed map[string]interface{}

	if err := json.Unmarshal(buf.Bytes(), &printed); err != nil {
		t.Fatal("JSON parse error: ", err)
	}

	updateTS := printed["time"].(string)
	updateTimestamp, err := time.Parse(time.RFC3339Nano, updateTS)
	if err != nil {
		t.Fatal("JSON parse error: ", updateTS, ": ", err)
	}

	lastValueTS := printed["updates"].([]interface{})[0].(map[string]interface{})["time"].(string)
	lastValueTimestamp, err := time.Parse(time.RFC3339Nano, lastValueTS)
	if err != nil {
		t.Fatal("JSON parse error: ", lastValueTS, ": ", err)
	}

	require.True(t, updateTimestamp.After(before))
	require.True(t, updateTimestamp.Before(after))

	require.True(t, lastValueTimestamp.After(before))
	require.True(t, lastValueTimestamp.Before(after))

	require.True(t, lastValueTimestamp.Before(updateTimestamp))
}

func TestStdoutCounterFormat(t *testing.T) {
	fix := newFixture(t, stdout.Config{})

	checkpointSet := test.NewCheckpointSet(testResource)

	desc := metric.NewDescriptor("test.name", metric.CounterKind, metric.Int64NumberKind)

	cagg, ckpt := test.Unslice2(sum.New(2))

	aggtest.CheckedUpdate(fix.t, cagg, metric.NewInt64Number(123), &desc)
	require.NoError(t, cagg.SynchronizedCopy(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, kv.String("A", "B"), kv.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `{"updates":[{"name":"test.name{R=V,A=B,C=D}","sum":123}]}`, fix.Output())
}

func TestStdoutLastValueFormat(t *testing.T) {
	fix := newFixture(t, stdout.Config{})

	checkpointSet := test.NewCheckpointSet(testResource)

	desc := metric.NewDescriptor("test.name", metric.ValueObserverKind, metric.Float64NumberKind)
	lvagg, ckpt := test.Unslice2(lastvalue.New(2))

	aggtest.CheckedUpdate(fix.t, lvagg, metric.NewFloat64Number(123.456), &desc)
	require.NoError(t, lvagg.SynchronizedCopy(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, kv.String("A", "B"), kv.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `{"updates":[{"name":"test.name{R=V,A=B,C=D}","last":123.456}]}`, fix.Output())
}

func TestStdoutMinMaxSumCount(t *testing.T) {
	fix := newFixture(t, stdout.Config{})

	checkpointSet := test.NewCheckpointSet(testResource)

	desc := metric.NewDescriptor("test.name", metric.ValueRecorderKind, metric.Float64NumberKind)

	magg, ckpt := test.Unslice2(minmaxsumcount.New(2, &desc))

	aggtest.CheckedUpdate(fix.t, magg, metric.NewFloat64Number(123.456), &desc)
	aggtest.CheckedUpdate(fix.t, magg, metric.NewFloat64Number(876.543), &desc)
	require.NoError(t, magg.SynchronizedCopy(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, kv.String("A", "B"), kv.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `{"updates":[{"name":"test.name{R=V,A=B,C=D}","min":123.456,"max":876.543,"sum":999.999,"count":2}]}`, fix.Output())
}

func TestStdoutValueRecorderFormat(t *testing.T) {
	fix := newFixture(t, stdout.Config{
		PrettyPrint: true,
	})

	checkpointSet := test.NewCheckpointSet(testResource)

	desc := metric.NewDescriptor("test.name", metric.ValueRecorderKind, metric.Float64NumberKind)
	aagg, ckpt := test.Unslice2(array.New(2))

	for i := 0; i < 1000; i++ {
		aggtest.CheckedUpdate(fix.t, aagg, metric.NewFloat64Number(float64(i)+0.5), &desc)
	}

	require.NoError(t, aagg.SynchronizedCopy(ckpt, &desc))

	checkpointSet.Add(&desc, ckpt, kv.String("A", "B"), kv.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `{
	"updates": [
		{
			"name": "test.name{R=V,A=B,C=D}",
			"min": 0.5,
			"max": 999.5,
			"sum": 500000,
			"count": 1000,
			"quantiles": [
				{
					"q": 0.5,
					"v": 500.5
				},
				{
					"q": 0.9,
					"v": 900.5
				},
				{
					"q": 0.99,
					"v": 990.5
				}
			]
		}
	]
}`, fix.Output())
}

func TestStdoutNoData(t *testing.T) {
	desc := metric.NewDescriptor("test.name", metric.ValueRecorderKind, metric.Float64NumberKind)

	runTwoAggs := func(agg, ckpt export.Aggregator) {
		t.Run(fmt.Sprintf("%T", agg), func(t *testing.T) {
			t.Parallel()

			fix := newFixture(t, stdout.Config{})

			checkpointSet := test.NewCheckpointSet(testResource)

			require.NoError(t, agg.SynchronizedCopy(ckpt, &desc))

			checkpointSet.Add(&desc, ckpt)

			fix.Export(checkpointSet)

			require.Equal(t, `{"updates":null}`, fix.Output())
		})
	}

	runTwoAggs(test.Unslice2(ddsketch.New(2, &desc, ddsketch.NewDefaultConfig())))
	runTwoAggs(test.Unslice2(minmaxsumcount.New(2, &desc)))
}

func TestStdoutLastValueNotSet(t *testing.T) {
	fix := newFixture(t, stdout.Config{})

	checkpointSet := test.NewCheckpointSet(testResource)

	desc := metric.NewDescriptor("test.name", metric.ValueObserverKind, metric.Float64NumberKind)

	lvagg, ckpt := test.Unslice2(lastvalue.New(2))
	require.NoError(t, lvagg.SynchronizedCopy(ckpt, &desc))

	checkpointSet.Add(&desc, lvagg, kv.String("A", "B"), kv.String("C", "D"))

	fix.Export(checkpointSet)

	require.Equal(t, `{"updates":null}`, fix.Output())
}

func TestStdoutResource(t *testing.T) {
	type testCase struct {
		expect string
		res    *resource.Resource
		attrs  []kv.KeyValue
	}
	newCase := func(expect string, res *resource.Resource, attrs ...kv.KeyValue) testCase {
		return testCase{
			expect: expect,
			res:    res,
			attrs:  attrs,
		}
	}
	testCases := []testCase{
		newCase("R1=V1,R2=V2,A=B,C=D",
			resource.New(kv.String("R1", "V1"), kv.String("R2", "V2")),
			kv.String("A", "B"),
			kv.String("C", "D")),
		newCase("R1=V1,R2=V2",
			resource.New(kv.String("R1", "V1"), kv.String("R2", "V2")),
		),
		newCase("A=B,C=D",
			nil,
			kv.String("A", "B"),
			kv.String("C", "D"),
		),
		// We explicitly do not de-duplicate between resources
		// and metric labels in this exporter.
		newCase("R1=V1,R2=V2,R1=V3,R2=V4",
			resource.New(kv.String("R1", "V1"), kv.String("R2", "V2")),
			kv.String("R1", "V3"),
			kv.String("R2", "V4")),
	}

	for _, tc := range testCases {
		fix := newFixture(t, stdout.Config{})

		checkpointSet := test.NewCheckpointSet(tc.res)

		desc := metric.NewDescriptor("test.name", metric.ValueObserverKind, metric.Float64NumberKind)
		lvagg, ckpt := test.Unslice2(lastvalue.New(2))

		aggtest.CheckedUpdate(fix.t, lvagg, metric.NewFloat64Number(123.456), &desc)
		require.NoError(t, lvagg.SynchronizedCopy(ckpt, &desc))

		checkpointSet.Add(&desc, ckpt, tc.attrs...)

		fix.Export(checkpointSet)

		require.Equal(t, `{"updates":[{"name":"test.name{`+tc.expect+`}","last":123.456}]}`, fix.Output())
	}
}
