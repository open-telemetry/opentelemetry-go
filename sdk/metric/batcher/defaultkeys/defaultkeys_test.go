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

package defaultkeys_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/batcher/defaultkeys"
	"go.opentelemetry.io/otel/sdk/metric/batcher/test"
)

func TestGroupingStateless(t *testing.T) {
	ctx := context.Background()
	b := defaultkeys.New(test.NewAggregationSelector(), test.GroupEncoder, false)

	_ = b.Process(ctx, export.NewRecord(test.GaugeDesc, test.Labels1, test.GaugeAgg(10)))
	_ = b.Process(ctx, export.NewRecord(test.GaugeDesc, test.Labels2, test.GaugeAgg(20)))
	_ = b.Process(ctx, export.NewRecord(test.GaugeDesc, test.Labels3, test.GaugeAgg(30)))

	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels1, test.CounterAgg(10)))
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels2, test.CounterAgg(20)))
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels3, test.CounterAgg(40)))

	checkpointSet := b.CheckpointSet()
	b.FinishedCollection()

	records := test.Output{}
	checkpointSet.ForEach(records.AddTo)

	// Output gauge should have only the "G=H" and "G=" keys.
	// Output counter should have only the "C=D" and "C=" keys.
	require.EqualValues(t, map[string]int64{
		"counter/C=D": 30, // labels1 + labels2
		"counter/C=":  40, // labels3
		"gauge/G=H":   10, // labels1
		"gauge/G=":    30, // labels3 = last value
	}, records)

	// Verify that state is reset by FinishedCollection()
	checkpointSet = b.CheckpointSet()
	b.FinishedCollection()
	checkpointSet.ForEach(func(rec export.Record) {
		t.Fatal("Unexpected call")
	})
}

func TestGroupingStateful(t *testing.T) {
	ctx := context.Background()
	b := defaultkeys.New(test.NewAggregationSelector(), test.GroupEncoder, true)

	cagg := test.CounterAgg(10)
	_ = b.Process(ctx, export.NewRecord(test.CounterDesc, test.Labels1, cagg))

	checkpointSet := b.CheckpointSet()
	b.FinishedCollection()

	records1 := test.Output{}
	checkpointSet.ForEach(records1.AddTo)

	require.EqualValues(t, map[string]int64{
		"counter/C=D": 10, // labels1
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
		"counter/C=D": 30,
	}, records4)
}
