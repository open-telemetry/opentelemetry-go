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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	processorTest "go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/processor/reducer"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	kvs1 = []label.KeyValue{
		label.Int("A", 1),
		label.Int("B", 2),
		label.Int("C", 3),
	}
	kvs2 = []label.KeyValue{
		label.Int("A", 1),
		label.Int("B", 0),
		label.Int("C", 3),
	}
)

type testFilter struct{}

func (testFilter) LabelFilterFor(_ *otel.Descriptor) label.Filter {
	return func(label label.KeyValue) bool {
		return label.Key == "A" || label.Key == "C"
	}
}

func generateData(impl otel.MeterImpl) {
	ctx := context.Background()
	meter := otel.WrapMeterImpl(impl, "testing")

	counter := otel.Must(meter).NewFloat64Counter("counter.sum")

	_ = otel.Must(meter).NewInt64SumObserver("observer.sum",
		func(_ context.Context, result otel.Int64ObserverResult) {
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
		label.DefaultEncoder(),
	)
	accum := metricsdk.NewAccumulator(
		reducer.New(testFilter{}, processorTest.Checkpointer(testProc)),
		metricsdk.WithResource(
			resource.New(label.String("R", "V")),
		),
	)
	generateData(accum)

	accum.Collect(context.Background())

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=1,C=3/R=V":  200,
		"observer.sum/A=1,C=3/R=V": 20,
	}, testProc.Values())
}

// Test a filter with the ../basic Processor.
func TestFilterBasicProcessor(t *testing.T) {
	basicProc := basic.New(processorTest.AggregatorSelector(), export.CumulativeExporter)
	accum := metricsdk.NewAccumulator(
		reducer.New(testFilter{}, basicProc),
		metricsdk.WithResource(
			resource.New(label.String("R", "V")),
		),
	)
	exporter := processorTest.NewExporter(basicProc, label.DefaultEncoder())

	generateData(accum)

	basicProc.StartCollection()
	accum.Collect(context.Background())
	if err := basicProc.FinishCollection(); err != nil {
		t.Error(err)
	}

	require.NoError(t, exporter.Export(context.Background(), basicProc.CheckpointSet()))

	require.EqualValues(t, map[string]float64{
		"counter.sum/A=1,C=3/R=V":  200,
		"observer.sum/A=1,C=3/R=V": 20,
	}, exporter.Values())
}
