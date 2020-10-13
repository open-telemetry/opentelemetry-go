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

package processortest_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	processorTest "go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/resource"
)

func generateTestData(proc export.Processor) {
	ctx := context.Background()
	accum := metricsdk.NewAccumulator(
		proc,
		metricsdk.WithResource(
			resource.New(label.String("R", "V")),
		),
	)
	meter := otel.WrapMeterImpl(accum, "testing")

	counter := otel.Must(meter).NewFloat64Counter("counter.sum")

	_ = otel.Must(meter).NewInt64SumObserver("observer.sum",
		func(_ context.Context, result otel.Int64ObserverResult) {
			result.Observe(10, label.String("K1", "V1"))
			result.Observe(11, label.String("K1", "V2"))
		},
	)

	counter.Add(ctx, 100, label.String("K1", "V1"))
	counter.Add(ctx, 101, label.String("K1", "V2"))

	accum.Collect(ctx)
}

func TestProcessorTesting(t *testing.T) {
	// Test the Processor test helper using a real Accumulator to
	// generate Accumulations.
	testProc := processorTest.NewProcessor(
		processorTest.AggregatorSelector(),
		label.DefaultEncoder(),
	)
	checkpointer := processorTest.Checkpointer(testProc)

	generateTestData(checkpointer)

	expect := map[string]float64{
		"counter.sum/K1=V1/R=V":  100,
		"counter.sum/K1=V2/R=V":  101,
		"observer.sum/K1=V1/R=V": 10,
		"observer.sum/K1=V2/R=V": 11,
	}

	// Validate the processor's checkpoint directly.
	require.EqualValues(t, expect, testProc.Values())

	// Export the data and validate it again.
	exporter := processorTest.NewExporter(
		export.PassThroughExporter,
		label.DefaultEncoder(),
	)

	err := exporter.Export(context.Background(), checkpointer.CheckpointSet())
	require.NoError(t, err)
	require.EqualValues(t, expect, exporter.Values())
}
