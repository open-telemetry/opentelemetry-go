// Copyright 2019, OpenTelemetry Authors
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

package ungrouped_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/batcher/test"
	"go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"
)

// These tests use the ../test label encoding.

func TestUngroupedStateless(t *testing.T) {
	ctx := context.Background()
	b := ungrouped.New(test.NewAggregationSelector(), false)

	// Set initial gauge values
	_ = b.Process(ctx, export.NewRecord(test.GaugeDesc, test.Labels1, test.GaugeAgg(10)))
	_ = b.Process(ctx, export.NewRecord(test.GaugeDesc, test.Labels2, test.GaugeAgg(20)))
	_ = b.Process(ctx, export.NewRecord(test.GaugeDesc, test.Labels3, test.GaugeAgg(30)))

	// Another gauge Set for Labels1
	_ = b.Process(ctx, export.NewRecord(test.GaugeDesc, test.Labels1, test.GaugeAgg(50)))

	// Set initial counter values
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels1, test.CounterAgg(10)))
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels2, test.CounterAgg(20)))
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels3, test.CounterAgg(40)))

	// Another counter Add for Labels1
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels1, test.CounterAgg(50)))

	checkpointSet := b.CheckpointSet()
	b.FinishedCollection()

	records := test.Output{}
	checkpointSet.ForEach(records.AddTo)

	// Output gauge should have only the "G=H" and "G=" keys.
	// Output counter should have only the "C=D" and "C=" keys.
	require.EqualValues(t, map[string]int64{
		"counter/G~H&C~D": 60, // labels1
		"counter/C~D&E~F": 20, // labels2
		"counter/":        40, // labels3
		"gauge/G~H&C~D":   50, // labels1
		"gauge/C~D&E~F":   20, // labels2
		"gauge/":          30, // labels3
	}, records)

	// Verify that state was reset
	checkpointSet = b.CheckpointSet()
	b.FinishedCollection()
	checkpointSet.ForEach(func(rec export.Record) {
		t.Fatal("Unexpected call")
	})
}

func TestUngroupedStateful(t *testing.T) {
	ctx := context.Background()
	b := ungrouped.New(test.NewAggregationSelector(), true)

	cagg := test.CounterAgg(10)
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels1, cagg))

	checkpointSet := b.CheckpointSet()
	b.FinishedCollection()

	records1 := test.Output{}
	checkpointSet.ForEach(records1.AddTo)

	require.EqualValues(t, map[string]int64{
		"counter/G~H&C~D": 10, // labels1
	}, records1)

	// Test that state was NOT reset
	checkpointSet = b.CheckpointSet()
	b.FinishedCollection()

	records2 := test.Output{}
	checkpointSet.ForEach(records2.AddTo)

	require.EqualValues(t, records1, records2)

	// Update and re-checkpoint the original record.
	_ = cagg.Update(ctx, core.NewInt64Number(20), test.CounterDesc)
	cagg.Checkpoint(ctx, test.CounterDesc)

	// As yet cagg has not been passed to Batcher.Process.  Should
	// not see an update.
	checkpointSet = b.CheckpointSet()
	b.FinishedCollection()

	records3 := test.Output{}
	checkpointSet.ForEach(records3.AddTo)

	require.EqualValues(t, records1, records3)

	// Now process the second update
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels1, cagg))

	checkpointSet = b.CheckpointSet()
	b.FinishedCollection()

	records4 := test.Output{}
	checkpointSet.ForEach(records4.AddTo)

	require.EqualValues(t, map[string]int64{
		"counter/G~H&C~D": 30,
	}, records4)
}
