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

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	processorTest "go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/resource"
)

func generateTestData(t *testing.T, proc export.Processor) {
	ctx := context.Background()
	accum := metricsdk.NewAccumulator(proc)
	meter := sdkapi.WrapMeterImpl(accum)

	counter, err := meter.SyncFloat64().Counter("counter.sum")
	require.NoError(t, err)

	counter.Add(ctx, 100, attribute.String("K1", "V1"))
	counter.Add(ctx, 101, attribute.String("K1", "V2"))

	counterObserver, err := meter.AsyncInt64().Counter("observer.sum")
	require.NoError(t, err)

	err = meter.RegisterCallback([]instrument.Asynchronous{counterObserver}, func(ctx context.Context) {
		counterObserver.Observe(ctx, 10, attribute.String("K1", "V1"))
		counterObserver.Observe(ctx, 11, attribute.String("K1", "V2"))
	})
	require.NoError(t, err)

	accum.Collect(ctx)
}

func TestProcessorTesting(t *testing.T) {
	// Test the Processor test helper using a real Accumulator to
	// generate Accumulations.
	checkpointer := processorTest.NewCheckpointer(
		processorTest.NewProcessor(
			processorTest.AggregatorSelector(),
			attribute.DefaultEncoder(),
		),
	)
	generateTestData(t, checkpointer)

	res := resource.NewSchemaless(attribute.String("R", "V"))
	expect := map[string]float64{
		"counter.sum/K1=V1/R=V":  100,
		"counter.sum/K1=V2/R=V":  101,
		"observer.sum/K1=V1/R=V": 10,
		"observer.sum/K1=V2/R=V": 11,
	}

	// Export the data and validate it again.
	exporter := processorTest.New(
		aggregation.StatelessTemporalitySelector(),
		attribute.DefaultEncoder(),
	)

	err := exporter.Export(context.Background(), res, processortest.OneInstrumentationLibraryReader(
		instrumentation.Library{
			Name: "test",
		},
		checkpointer.Reader(),
	))
	require.NoError(t, err)
	require.EqualValues(t, expect, exporter.Values())
}
