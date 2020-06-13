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
func NewLastValueRecord(desc *metric.Descriptor, labels *label.Set, value int64) export.Record {
	return export.NewRecord(desc, labels, Resource, LastValueAgg(desc, value))
}

// Convenience method for building a test exported counter record.
func NewCounterRecord(desc *metric.Descriptor, labels *label.Set, value int64) export.Record {
	return export.NewRecord(desc, labels, Resource, CounterAgg(desc, value))
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

	// Set initial lastValue values
	_ = b.Process(NewLastValueRecord(&LastValueADesc, Labels1, 10))
	_ = b.Process(NewLastValueRecord(&LastValueADesc, Labels2, 20))
	_ = b.Process(NewLastValueRecord(&LastValueADesc, Labels3, 30))

	_ = b.Process(NewLastValueRecord(&LastValueBDesc, Labels1, 10))
	_ = b.Process(NewLastValueRecord(&LastValueBDesc, Labels2, 20))
	_ = b.Process(NewLastValueRecord(&LastValueBDesc, Labels3, 30))

	// Another lastValue Set for Labels1
	_ = b.Process(NewLastValueRecord(&LastValueADesc, Labels1, 50))
	_ = b.Process(NewLastValueRecord(&LastValueBDesc, Labels1, 50))

	// Set initial counter values
	_ = b.Process(NewCounterRecord(&CounterADesc, Labels1, 10))
	_ = b.Process(NewCounterRecord(&CounterADesc, Labels2, 20))
	_ = b.Process(NewCounterRecord(&CounterADesc, Labels3, 40))

	_ = b.Process(NewCounterRecord(&CounterBDesc, Labels1, 10))
	_ = b.Process(NewCounterRecord(&CounterBDesc, Labels2, 20))
	_ = b.Process(NewCounterRecord(&CounterBDesc, Labels3, 40))

	// Another counter Add for Labels1
	_ = b.Process(NewCounterRecord(&CounterADesc, Labels1, 50))
	_ = b.Process(NewCounterRecord(&CounterBDesc, Labels1, 50))

	checkpointSet := b.CheckpointSet()

	records := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records.AddTo)

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
	b := simple.New(test.AggregationSelector(), true)

	counterA := NewCounterRecord(&CounterADesc, Labels1, 10)
	_ = b.Process(counterA)

	counterB := NewCounterRecord(&CounterBDesc, Labels1, 10)
	_ = b.Process(counterB)

	checkpointSet := b.CheckpointSet()
	b.FinishedCollection()

	records1 := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records1.AddTo)

	require.EqualValues(t, map[string]float64{
		"a.sum/C=D,G=H/R=V": 10, // labels1
		"b.sum/C=D,G=H/R=V": 10, // labels1
	}, records1.Map)

	alloc := sum.New(4)
	caggA, caggB, ckptA, ckptB := &alloc[0], &alloc[1], &alloc[2], &alloc[3]

	// Test that state was NOT reset
	checkpointSet = b.CheckpointSet()

	records2 := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records2.AddTo)

	require.EqualValues(t, records1.Map, records2.Map)
	b.FinishedCollection()

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
	_ = checkpointSet.ForEach(records3.AddTo)

	require.EqualValues(t, records1.Map, records3.Map)
	b.FinishedCollection()

	// Now process the second update
	_ = b.Process(export.NewRecord(&CounterADesc, Labels1, Resource, ckptA))
	_ = b.Process(export.NewRecord(&CounterBDesc, Labels1, Resource, ckptB))

	checkpointSet = b.CheckpointSet()

	records4 := test.NewOutput(label.DefaultEncoder())
	_ = checkpointSet.ForEach(records4.AddTo)

	require.EqualValues(t, map[string]float64{
		"a.sum/C=D,G=H/R=V": 30,
		"b.sum/C=D,G=H/R=V": 30,
	}, records4.Map)
	b.FinishedCollection()
}
