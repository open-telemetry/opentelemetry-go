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

package reducer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	processorTest "go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/processor/reducer"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	kvs1 = []attribute.KeyValue{
		attribute.Int("A", 1),
		attribute.Int("B", 2),
		attribute.Int("C", 3),
	}
	kvs2 = []attribute.KeyValue{
		attribute.Int("A", 1),
		attribute.Int("B", 0),
		attribute.Int("C", 3),
	}
)

type testFilter struct{}

func (testFilter) LabelFilterFor(_ *sdkapi.Descriptor) attribute.Filter {
	return func(label attribute.KeyValue) bool {
		return label.Key == "A" || label.Key == "C"
	}
}

func generateData(impl sdkapi.MeterImpl) {
	ctx := context.Background()
	meter := metric.WrapMeterImpl(impl)

	counter := metric.Must(meter).NewFloat64Counter("counter.sum")

	_ = metric.Must(meter).NewInt64CounterObserver("observer.sum",
		func(_ context.Context, result metric.Int64ObserverResult) {
			result.Observe(10, kvs1...)
			result.Observe(10, kvs2...)
		},
	)

	counter.Add(ctx, 100, kvs1...)
	counter.Add(ctx, 100, kvs2...)
}

func TestFilterProcessor(t *testing.T) {
	testProc := processorTest.NewProcessor(
		processorTest.AggregatorSelector(),
		attribute.DefaultEncoder(),
	)
	accum := metricsdk.NewAccumulator(
		reducer.New(testFilter{}, processorTest.NewCheckpointer(testProc)),
	)
	generateData(accum)

	accum.Collect(context.Background())

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=1,C=3/":  200,
		"observer.sum/A=1,C=3/": 20,
	}, testProc.Values())
}

// Test a filter with the ../basic Processor.
func TestFilterBasicProcessor(t *testing.T) {
	basicProc := basic.New(processorTest.AggregatorSelector(), aggregation.CumulativeTemporalitySelector())
	accum := metricsdk.NewAccumulator(
		reducer.New(testFilter{}, basicProc),
	)
	exporter := processorTest.New(basicProc, attribute.DefaultEncoder())

	generateData(accum)

	basicProc.StartCollection()
	accum.Collect(context.Background())
	if err := basicProc.FinishCollection(); err != nil {
		t.Error(err)
	}

	res := resource.NewSchemaless(attribute.String("R", "V"))
	require.NoError(t, exporter.Export(context.Background(), res, processortest.OneInstrumentationLibraryReader(instrumentation.Library{
		Name: "test",
	}, basicProc.Reader())))

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=1,C=3/R=V":  200,
		"observer.sum/A=1,C=3/R=V": 20,
	}, exporter.Values())
}
