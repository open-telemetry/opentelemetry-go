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

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/batcher/test"
	"go.opentelemetry.io/otel/sdk/metric/batcher/ungrouped"
)

// These tests use the original label encoding.

func TestUngroupedStateless(t *testing.T) {
	ctx := context.Background()
	b := ungrouped.New(nil, false)

	_ = b.Process(ctx, test.GaugeDesc, test.Labels1, test.GaugeAgg(10))
	_ = b.Process(ctx, test.GaugeDesc, test.Labels2, test.GaugeAgg(20))
	_ = b.Process(ctx, test.GaugeDesc, test.Labels3, test.GaugeAgg(30))

	_ = b.Process(ctx, test.CounterDesc, test.Labels1, test.CounterAgg(10))
	_ = b.Process(ctx, test.CounterDesc, test.Labels2, test.CounterAgg(20))
	_ = b.Process(ctx, test.CounterDesc, test.Labels3, test.CounterAgg(40))

	processor := b.ReadCheckpoint()

	records := test.Output{}
	processor.Foreach(records.AddTo)

	// Output gauge should have only the "G=H" and "G=" keys.
	// Output counter should have only the "C=D" and "C=" keys.
	require.EqualValues(t, map[string]int64{
		"counter/G~H&C~D": 10, // labels1
		"counter/C~D&E~F": 20, // labels2
		"counter/":        40, // labels3
		"gauge/G~H&C~D":   10, // labels1
		"gauge/C~D&E~F":   20, // labels2
		"gauge/":          30, // labels3
	}, records)

	// Verify that state was reset
	processor = b.ReadCheckpoint()
	processor.Foreach(func(rec export.Record) {
		t.Fatal("Unexpected call")
	})
}

func TestUngroupedStateful(t *testing.T) {
	ctx := context.Background()
	b := ungrouped.New(nil, true)

	_ = b.Process(ctx, test.CounterDesc, test.Labels1, test.CounterAgg(10))

	processor := b.ReadCheckpoint()

	records1 := test.Output{}
	processor.Foreach(records1.AddTo)

	require.EqualValues(t, map[string]int64{
		"counter/G~H&C~D": 10, // labels1
	}, records1)

	// Test that state was NOT reset
	processor = b.ReadCheckpoint()

	records2 := test.Output{}
	processor.Foreach(records2.AddTo)

	require.EqualValues(t, records1, records2)
}
