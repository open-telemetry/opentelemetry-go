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

package stdoutmetric_test

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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/resource"
)

type testFixture struct {
	t        *testing.T
	ctx      context.Context
	cont     *controller.Controller
	meter    metric.Meter
	exporter *stdoutmetric.Exporter
	output   *bytes.Buffer
}

var testResource = resource.NewSchemaless(attribute.String("R", "V"))

func newFixture(t *testing.T, opts ...stdoutmetric.Option) testFixture {
	return newFixtureWithResource(t, testResource, opts...)
}

func newFixtureWithResource(t *testing.T, res *resource.Resource, opts ...stdoutmetric.Option) testFixture {
	buf := &bytes.Buffer{}
	opts = append(opts, stdoutmetric.WithWriter(buf))
	opts = append(opts, stdoutmetric.WithoutTimestamps())
	exp, err := stdoutmetric.New(opts...)
	if err != nil {
		t.Fatal("Error building fixture: ", err)
	}
	aggSel := processortest.AggregatorSelector()
	proc := processor.NewFactory(aggSel, aggregation.StatelessTemporalitySelector())
	cont := controller.New(proc,
		controller.WithExporter(exp),
		controller.WithResource(res),
	)
	ctx := context.Background()
	require.NoError(t, cont.Start(ctx))
	meter := cont.Meter("test")

	return testFixture{
		t:        t,
		ctx:      ctx,
		exporter: exp,
		cont:     cont,
		meter:    meter,
		output:   buf,
	}
}

func (fix testFixture) Output() string {
	return strings.TrimSpace(fix.output.String())
}

func TestStdoutTimestamp(t *testing.T) {
	var buf bytes.Buffer
	aggSel := processortest.AggregatorSelector()
	proc := processor.NewFactory(aggSel, aggregation.CumulativeTemporalitySelector())
	exporter, err := stdoutmetric.New(
		stdoutmetric.WithWriter(&buf),
	)
	if err != nil {
		t.Fatal("Invalid config: ", err)
	}
	cont := controller.New(proc,
		controller.WithExporter(exporter),
		controller.WithResource(testResource),
	)
	ctx := context.Background()

	require.NoError(t, cont.Start(ctx))
	meter := cont.Meter("test")
	counter, err := meter.SyncInt64().Counter("name.lastvalue")
	require.NoError(t, err)

	before := time.Now()
	// Ensure the timestamp is after before.
	time.Sleep(time.Nanosecond)

	counter.Add(ctx, 1)

	require.NoError(t, cont.Stop(ctx))

	// Ensure the timestamp is before after.
	time.Sleep(time.Nanosecond)
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

	counter, err := fix.meter.SyncInt64().Counter("name.sum")
	require.NoError(t, err)
	counter.Add(fix.ctx, 123, attribute.String("A", "B"), attribute.String("C", "D"))

	require.NoError(t, fix.cont.Stop(fix.ctx))

	require.Equal(t, `[{"Name":"name.sum{R=V,instrumentation.name=test,A=B,C=D}","Sum":123}]`, fix.Output())
}

func TestStdoutLastValueFormat(t *testing.T) {
	fix := newFixture(t)

	counter, err := fix.meter.SyncFloat64().Counter("name.lastvalue")
	require.NoError(t, err)
	counter.Add(fix.ctx, 123.456, attribute.String("A", "B"), attribute.String("C", "D"))

	require.NoError(t, fix.cont.Stop(fix.ctx))

	require.Equal(t, `[{"Name":"name.lastvalue{R=V,instrumentation.name=test,A=B,C=D}","Last":123.456}]`, fix.Output())
}

func TestStdoutHistogramFormat(t *testing.T) {
	fix := newFixture(t, stdoutmetric.WithPrettyPrint())

	inst, err := fix.meter.SyncFloat64().Histogram("name.histogram")
	require.NoError(t, err)

	for i := 0; i < 1000; i++ {
		inst.Record(fix.ctx, float64(i)+0.5, attribute.String("A", "B"), attribute.String("C", "D"))
	}
	require.NoError(t, fix.cont.Stop(fix.ctx))

	// TODO: Stdout does not export `Count` for histogram, nor the buckets.
	require.Equal(t, `[
	{
		"Name": "name.histogram{R=V,instrumentation.name=test,A=B,C=D}",
		"Sum": 500000
	}
]`, fix.Output())
}

func TestStdoutNoData(t *testing.T) {
	runTwoAggs := func(aggName string) {
		t.Run(aggName, func(t *testing.T) {
			t.Parallel()

			fix := newFixture(t)
			_, err := fix.meter.SyncFloat64().Counter(fmt.Sprint("name.", aggName))
			require.NoError(t, err)
			require.NoError(t, fix.cont.Stop(fix.ctx))

			require.Equal(t, "", fix.Output())
		})
	}

	runTwoAggs("lastvalue")
}

func TestStdoutResource(t *testing.T) {
	type testCase struct {
		name   string
		expect string
		res    *resource.Resource
		attrs  []attribute.KeyValue
	}
	newCase := func(name, expect string, res *resource.Resource, attrs ...attribute.KeyValue) testCase {
		return testCase{
			name:   name,
			expect: expect,
			res:    res,
			attrs:  attrs,
		}
	}
	testCases := []testCase{
		newCase("resource and attribute",
			"R1=V1,R2=V2,instrumentation.name=test,A=B,C=D",
			resource.NewSchemaless(attribute.String("R1", "V1"), attribute.String("R2", "V2")),
			attribute.String("A", "B"),
			attribute.String("C", "D")),
		newCase("only resource",
			"R1=V1,R2=V2,instrumentation.name=test",
			resource.NewSchemaless(attribute.String("R1", "V1"), attribute.String("R2", "V2")),
		),
		newCase("empty resource",
			"instrumentation.name=test,A=B,C=D",
			resource.Empty(),
			attribute.String("A", "B"),
			attribute.String("C", "D"),
		),
		newCase("default resource",
			fmt.Sprint(resource.Default().Encoded(attribute.DefaultEncoder()),
				",instrumentation.name=test,A=B,C=D"),
			resource.Default(),
			attribute.String("A", "B"),
			attribute.String("C", "D"),
		),
		// We explicitly do not de-duplicate between resources
		// and metric attributes in this exporter.
		newCase("resource deduplication",
			"R1=V1,R2=V2,instrumentation.name=test,R1=V3,R2=V4",
			resource.NewSchemaless(attribute.String("R1", "V1"), attribute.String("R2", "V2")),
			attribute.String("R1", "V3"),
			attribute.String("R2", "V4")),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			fix := newFixtureWithResource(t, tc.res)

			counter, err := fix.meter.SyncFloat64().Counter("name.lastvalue")
			require.NoError(t, err)
			counter.Add(ctx, 123.456, tc.attrs...)

			require.NoError(t, fix.cont.Stop(fix.ctx))

			require.Equal(t, `[{"Name":"name.lastvalue{`+tc.expect+`}","Last":123.456}]`, fix.Output())
		})
	}
}
