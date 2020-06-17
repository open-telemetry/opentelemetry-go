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
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/metric/integrator/test"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Note: This var block and the helpers below will disappear in a
// future PR (see the draft in #799).  The test has been completely
// rewritten there, so this code will simply be dropped.

var (
	// Resource is applied to all test records built in this package.
	Resource = resource.New(kv.String("R", "V"))

	// LastValueADesc and LastValueBDesc group by "G"
	LastValueADesc = metric.NewDescriptor(
		"a.lastvalue", metric.ValueObserverKind, metric.Int64NumberKind)
	LastValueBDesc = metric.NewDescriptor(
		"b.lastvalue", metric.ValueObserverKind, metric.Int64NumberKind)
	// CounterADesc and CounterBDesc group by "C"
	CounterADesc = metric.NewDescriptor(
		"a.sum", metric.CounterKind, metric.Int64NumberKind)
	CounterBDesc = metric.NewDescriptor(
		"b.sum", metric.CounterKind, metric.Int64NumberKind)

	// LastValue groups are (labels1), (labels2+labels3)
	// Counter groups are (labels1+labels2), (labels3)

	// Labels1 has G=H and C=D
	Labels1 = makeLabels(kv.String("G", "H"), kv.String("C", "D"))
	// Labels2 has C=D and E=F
	Labels2 = makeLabels(kv.String("C", "D"), kv.String("E", "F"))
	// Labels3 is the empty set
	Labels3 = makeLabels()
)

func makeLabels(labels ...kv.KeyValue) *label.Set {
	s := label.NewSet(labels...)
	return &s
}

// LastValueAgg returns a checkpointed lastValue aggregator w/ the specified descriptor and value.
func LastValueAgg(desc *metric.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	gagg := &lastvalue.New(1)[0]
	_ = gagg.Update(ctx, metric.NewInt64Number(v), desc)
	return gagg
}

// Convenience method for building a test exported lastValue record.
func NewLastValueAccumulation(desc *metric.Descriptor, labels *label.Set, value int64) export.Accumulation {
	return export.NewAccumulation(desc, labels, Resource, LastValueAgg(desc, value))
}

// Convenience method for building a test exported counter record.
func NewCounterAccumulation(desc *metric.Descriptor, labels *label.Set, value int64) export.Accumulation {
	return export.NewAccumulation(desc, labels, Resource, CounterAgg(desc, value))
}

// CounterAgg returns a checkpointed counter aggregator w/ the specified descriptor and value.
func CounterAgg(desc *metric.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	cagg := &sum.New(1)[0]
	_ = cagg.Update(ctx, metric.NewInt64Number(v), desc)
	return cagg
}

func TestSimpleStateless(t *testing.T) {
	b := simple.New(test.AggregationSelector(), false)

	b.StartCollection()

	// Set initial lastValue values
	_ = b.Process(NewLastValueAccumulation(&LastValueADesc, Labels1, 10))
	_ = b.Process(NewLastValueAccumulation(&LastValueADesc, Labels2, 20))
	_ = b.Process(NewLastValueAccumulation(&LastValueADesc, Labels3, 30))

	_ = b.Process(NewLastValueAccumulation(&LastValueBDesc, Labels1, 10))
	_ = b.Process(NewLastValueAccumulation(&LastValueBDesc, Labels2, 20))
	_ = b.Process(NewLastValueAccumulation(&LastValueBDesc, Labels3, 30))

	// Another lastValue Set for Labels1
	_ = b.Process(NewLastValueAccumulation(&LastValueADesc, Labels1, 50))
	_ = b.Process(NewLastValueAccumulation(&LastValueBDesc, Labels1, 50))

	// Set initial counter values
	_ = b.Process(NewCounterAccumulation(&CounterADesc, Labels1, 10))
	_ = b.Process(NewCounterAccumulation(&CounterADesc, Labels2, 20))
	_ = b.Process(NewCounterAccumulation(&CounterADesc, Labels3, 40))

	_ = b.Process(NewCounterAccumulation(&CounterBDesc, Labels1, 10))
	_ = b.Process(NewCounterAccumulation(&CounterBDesc, Labels2, 20))
	_ = b.Process(NewCounterAccumulation(&CounterBDesc, Labels3, 40))

	// Another counter Add for Labels1
	_ = b.Process(NewCounterAccumulation(&CounterADesc, Labels1, 50))
	_ = b.Process(NewCounterAccumulation(&CounterBDesc, Labels1, 50))

	require.NoError(t, b.FinishCollection())

	checkpointSet := b.CheckpointSet()

	records := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records.AddRecord)

	// Output lastvalue should have only the "G=H" and "G=" keys.
	// Output counter should have only the "C=D" and "C=" keys.
	require.EqualValues(t, map[string]float64{
		"a.sum/C=D,G=H/R=V":       60, // labels1
		"a.sum/C=D,E=F/R=V":       20, // labels2
		"a.sum//R=V":              40, // labels3
		"b.sum/C=D,G=H/R=V":       60, // labels1
		"b.sum/C=D,E=F/R=V":       20, // labels2
		"b.sum//R=V":              40, // labels3
		"a.lastvalue/C=D,G=H/R=V": 50, // labels1
		"a.lastvalue/C=D,E=F/R=V": 20, // labels2
		"a.lastvalue//R=V":        30, // labels3
		"b.lastvalue/C=D,G=H/R=V": 50, // labels1
		"b.lastvalue/C=D,E=F/R=V": 20, // labels2
		"b.lastvalue//R=V":        30, // labels3
	}, records.Map)

	// Verify that state was reset
	b.StartCollection()
	require.NoError(t, b.FinishCollection())
	checkpointSet = b.CheckpointSet()
	_ = checkpointSet.ForEach(func(rec export.Record) error {
		t.Fatal("Unexpected call")
		return nil
	})
}

func TestSimpleStateful(t *testing.T) {
	ctx := context.Background()
	b := simple.New(test.AggregationSelector(), true)

	b.StartCollection()

	counterA := NewCounterAccumulation(&CounterADesc, Labels1, 10)
	_ = b.Process(counterA)

	counterB := NewCounterAccumulation(&CounterBDesc, Labels1, 10)
	_ = b.Process(counterB)
	require.NoError(t, b.FinishCollection())

	checkpointSet := b.CheckpointSet()

	records1 := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records1.AddRecord)

	require.EqualValues(t, map[string]float64{
		"a.sum/C=D,G=H/R=V": 10, // labels1
		"b.sum/C=D,G=H/R=V": 10, // labels1
	}, records1.Map)

	alloc := sum.New(4)
	caggA, caggB, ckptA, ckptB := &alloc[0], &alloc[1], &alloc[2], &alloc[3]

	// Test that state was NOT reset
	checkpointSet = b.CheckpointSet()

	b.StartCollection()
	require.NoError(t, b.FinishCollection())

	records2 := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records2.AddRecord)

	require.EqualValues(t, records1.Map, records2.Map)

	// Update and re-checkpoint the original record.
	_ = caggA.Update(ctx, metric.NewInt64Number(20), &CounterADesc)
	_ = caggB.Update(ctx, metric.NewInt64Number(20), &CounterBDesc)
	err := caggA.SynchronizedCopy(ckptA, &CounterADesc)
	require.NoError(t, err)
	err = caggB.SynchronizedCopy(ckptB, &CounterBDesc)
	require.NoError(t, err)

	// As yet cagg has not been passed to Integrator.Process.  Should
	// not see an update.
	checkpointSet = b.CheckpointSet()

	records3 := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records3.AddRecord)

	require.EqualValues(t, records1.Map, records3.Map)
	b.StartCollection()

	// Now process the second update
	_ = b.Process(export.NewAccumulation(&CounterADesc, Labels1, Resource, ckptA))
	_ = b.Process(export.NewAccumulation(&CounterBDesc, Labels1, Resource, ckptB))
	require.NoError(t, b.FinishCollection())

	checkpointSet = b.CheckpointSet()

	records4 := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records4.AddRecord)

	require.EqualValues(t, map[string]float64{
		"a.sum/C=D,G=H/R=V": 30,
		"b.sum/C=D,G=H/R=V": 30,
	}, records4.Map)
}

func TestSimpleInconsistent(t *testing.T) {
	// Test double-start
	b := simple.New(test.AggregationSelector(), true)

	b.StartCollection()
	b.StartCollection()
	require.Equal(t, simple.ErrInconsistentState, b.FinishCollection())

	// Test finish without start
	b = simple.New(test.AggregationSelector(), true)

	require.Equal(t, simple.ErrInconsistentState, b.FinishCollection())

	// Test no finish
	b = simple.New(test.AggregationSelector(), true)

	b.StartCollection()
	require.Equal(t, simple.ErrInconsistentState, b.ForEach(func(export.Record) error { return nil }))

	// Test no start
	b = simple.New(test.AggregationSelector(), true)

	require.Equal(t, simple.ErrInconsistentState, b.Process(NewCounterAccumulation(&CounterADesc, Labels1, 10)))
}

func TestSimpleTimestamps(t *testing.T) {
	beforeNew := time.Now()
	b := simple.New(test.AggregationSelector(), true)
	afterNew := time.Now()

	b.StartCollection()
	_ = b.Process(NewCounterAccumulation(&CounterADesc, Labels1, 10))
	require.NoError(t, b.FinishCollection())

	var start1, end1 time.Time

	require.NoError(t, b.ForEach(func(rec export.Record) error {
		start1 = rec.StartTime()
		end1 = rec.EndTime()
		return nil
	}))

	// The first start time is set in the constructor.
	require.True(t, beforeNew.Before(start1))
	require.True(t, afterNew.After(start1))

	for i := 0; i < 2; i++ {
		b.StartCollection()
		require.NoError(t, b.Process(NewCounterAccumulation(&CounterADesc, Labels1, 10)))
		require.NoError(t, b.FinishCollection())

		var start2, end2 time.Time

		require.NoError(t, b.ForEach(func(rec export.Record) error {
			start2 = rec.StartTime()
			end2 = rec.EndTime()
			return nil
		}))

		// Subsequent intervals have their start and end aligned.
		require.Equal(t, start2, end1)
		require.True(t, start1.Before(end1))
		require.True(t, start2.Before(end2))

		start1 = start2
		end1 = end2
	}
}
