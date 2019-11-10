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

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/batcher/defaultkeys"
	"go.opentelemetry.io/otel/sdk/metric/batcher/test"
)

func TestGroupingStateless(t *testing.T) {
	ctx := context.Background()
	b := defaultkeys.New(nil, test.GroupEncoder, false)

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
		"counter/C=D": 30, // labels1 + labels2
		"counter/C=":  40, // labels3
		"gauge/G=H":   10, // labels1
		"gauge/G=":    30, // labels3 = last value
	}, records)

	// Verify that state was reset
	processor = b.ReadCheckpoint()
	processor.Foreach(func(rec export.Record) {
		t.Fatal("Unexpected call")
	})
}

func TestGroupingStateful(t *testing.T) {
	ctx := context.Background()
	b := defaultkeys.New(nil, test.GroupEncoder, true)

	_ = b.Process(ctx, test.CounterDesc, test.Labels1, test.CounterAgg(10))

	processor := b.ReadCheckpoint()

	records1 := test.Output{}
	processor.Foreach(records1.AddTo)

	require.EqualValues(t, map[string]int64{
		"counter/C=D": 10, // labels1
	}, records1)

	// Test that state was NOT reset
	processor = b.ReadCheckpoint()

	records2 := test.Output{}
	processor.Foreach(records2.AddTo)

	require.EqualValues(t, records1, records2)
}
