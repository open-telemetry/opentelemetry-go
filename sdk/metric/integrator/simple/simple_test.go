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

package simple_test

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/api/metric"

	"github.com/stretchr/testify/require"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/metric/integrator/test"
)

// These tests use the ../test label encoding.

func TestSimpleStateless(t *testing.T) {
	ctx := context.Background()
	b := simple.New(test.NewAggregationSelector(), false)

	// Set initial lastValue values
	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueADesc, test.Labels1, 10))
	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueADesc, test.Labels2, 20))
	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueADesc, test.Labels3, 30))

	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueBDesc, test.Labels1, 10))
	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueBDesc, test.Labels2, 20))
	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueBDesc, test.Labels3, 30))

	// Another lastValue Set for Labels1
	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueADesc, test.Labels1, 50))
	_ = b.Process(ctx, test.NewLastValueRecord(&test.LastValueBDesc, test.Labels1, 50))

	// Set initial counter values
	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterADesc, test.Labels1, 10))
	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterADesc, test.Labels2, 20))
	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterADesc, test.Labels3, 40))

	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterBDesc, test.Labels1, 10))
	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterBDesc, test.Labels2, 20))
	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterBDesc, test.Labels3, 40))

	// Another counter Add for Labels1
	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterADesc, test.Labels1, 50))
	_ = b.Process(ctx, test.NewCounterRecord(&test.CounterBDesc, test.Labels1, 50))

	checkpointSet := b.CheckpointSet()

	records := test.NewOutput(test.SdkEncoder)
	_ = checkpointSet.ForEach(records.AddTo)

	// Output lastvalue should have only the "G=H" and "G=" keys.
	// Output counter should have only the "C=D" and "C=" keys.
	require.EqualValues(t, map[string]float64{
		"sum.a/C~D&G~H/R~V":       60, // labels1
		"sum.a/C~D&E~F/R~V":       20, // labels2
		"sum.a//R~V":              40, // labels3
		"sum.b/C~D&G~H/R~V":       60, // labels1
		"sum.b/C~D&E~F/R~V":       20, // labels2
		"sum.b//R~V":              40, // labels3
		"lastvalue.a/C~D&G~H/R~V": 50, // labels1
		"lastvalue.a/C~D&E~F/R~V": 20, // labels2
		"lastvalue.a//R~V":        30, // labels3
		"lastvalue.b/C~D&G~H/R~V": 50, // labels1
		"lastvalue.b/C~D&E~F/R~V": 20, // labels2
		"lastvalue.b//R~V":        30, // labels3
	}, records.Map)
	b.FinishedCollection()

	// Verify that state was reset
	checkpointSet = b.CheckpointSet()
	_ = checkpointSet.ForEach(func(rec export.Record) error {
		t.Fatal("Unexpected call")
		return nil
	})
	b.FinishedCollection()
}

func TestSimpleStateful(t *testing.T) {
	ctx := context.Background()
	b := simple.New(test.NewAggregationSelector(), true)

	counterA := test.NewCounterRecord(&test.CounterADesc, test.Labels1, 10)
	caggA := counterA.Aggregator()
	_ = b.Process(ctx, counterA)

	counterB := test.NewCounterRecord(&test.CounterBDesc, test.Labels1, 10)
	caggB := counterB.Aggregator()
	_ = b.Process(ctx, counterB)

	checkpointSet := b.CheckpointSet()
	b.FinishedCollection()

	records1 := test.NewOutput(test.SdkEncoder)
	_ = checkpointSet.ForEach(records1.AddTo)

	require.EqualValues(t, map[string]float64{
		"sum.a/C~D&G~H/R~V": 10, // labels1
		"sum.b/C~D&G~H/R~V": 10, // labels1
	}, records1.Map)

	// Test that state was NOT reset
	checkpointSet = b.CheckpointSet()

	records2 := test.NewOutput(test.SdkEncoder)
	_ = checkpointSet.ForEach(records2.AddTo)

	require.EqualValues(t, records1.Map, records2.Map)
	b.FinishedCollection()

	// Update and re-checkpoint the original record.
	_ = caggA.Update(ctx, metric.NewInt64Number(20), &test.CounterADesc)
	_ = caggB.Update(ctx, metric.NewInt64Number(20), &test.CounterBDesc)
	caggA.Checkpoint(ctx, &test.CounterADesc)
	caggB.Checkpoint(ctx, &test.CounterBDesc)

	// As yet cagg has not been passed to Integrator.Process.  Should
	// not see an update.
	checkpointSet = b.CheckpointSet()

	records3 := test.NewOutput(test.SdkEncoder)
	_ = checkpointSet.ForEach(records3.AddTo)

	require.EqualValues(t, records1.Map, records3.Map)
	b.FinishedCollection()

	// Now process the second update
	_ = b.Process(ctx, export.NewRecord(&test.CounterADesc, test.Labels1, test.Resource, caggA))
	_ = b.Process(ctx, export.NewRecord(&test.CounterBDesc, test.Labels1, test.Resource, caggB))

	checkpointSet = b.CheckpointSet()

	records4 := test.NewOutput(test.SdkEncoder)
	_ = checkpointSet.ForEach(records4.AddTo)

	require.EqualValues(t, map[string]float64{
		"sum.a/C~D&G~H/R~V": 30,
		"sum.b/C~D&G~H/R~V": 30,
	}, records4.Map)
	b.FinishedCollection()
}
