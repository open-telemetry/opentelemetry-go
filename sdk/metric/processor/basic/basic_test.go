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

package basic_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/export/metric/metrictest"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	processorTest "go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/resource"
)

func requireNotAfter(t *testing.T, t1, t2 time.Time) {
	require.False(t, t1.After(t2), "expected %v ≤ %v", t1, t2)
}

// TestProcessor tests all the non-error paths in this package.
func TestProcessor(t *testing.T) {
	type exportCase struct {
		kind export.ExportKind
	}
	type instrumentCase struct {
		kind metric.InstrumentKind
	}
	type numberCase struct {
		kind number.Kind
	}
	type aggregatorCase struct {
		kind aggregation.Kind
	}

	for _, tc := range []exportCase{
		{kind: export.CumulativeExportKind},
		{kind: export.DeltaExportKind},
	} {
		t.Run(tc.kind.String(), func(t *testing.T) {
			for _, ic := range []instrumentCase{
				{kind: metric.CounterInstrumentKind},
				{kind: metric.UpDownCounterInstrumentKind},
				{kind: metric.ValueRecorderInstrumentKind},
				{kind: metric.SumObserverInstrumentKind},
				{kind: metric.UpDownSumObserverInstrumentKind},
				{kind: metric.ValueObserverInstrumentKind},
			} {
				t.Run(ic.kind.String(), func(t *testing.T) {
					for _, nc := range []numberCase{
						{kind: number.Int64Kind},
						{kind: number.Float64Kind},
					} {
						t.Run(nc.kind.String(), func(t *testing.T) {
							for _, ac := range []aggregatorCase{
								{kind: aggregation.SumKind},
								{kind: aggregation.MinMaxSumCountKind},
								{kind: aggregation.HistogramKind},
								{kind: aggregation.LastValueKind},
								{kind: aggregation.ExactKind},
							} {
								t.Run(ac.kind.String(), func(t *testing.T) {
									testProcessor(
										t,
										tc.kind,
										ic.kind,
										nc.kind,
										ac.kind,
									)
								})
							}
						})
					}
				})
			}
		})
	}
}

func asNumber(nkind number.Kind, value int64) number.Number {
	if nkind == number.Int64Kind {
		return number.NewInt64Number(value)
	}
	return number.NewFloat64Number(float64(value))
}

func updateFor(t *testing.T, desc *metric.Descriptor, selector export.AggregatorSelector, res *resource.Resource, value int64, labs ...attribute.KeyValue) export.Accumulation {
	ls := attribute.NewSet(labs...)
	var agg export.Aggregator
	selector.AggregatorFor(desc, &agg)
	require.NoError(t, agg.Update(context.Background(), asNumber(desc.NumberKind(), value), desc))

	return export.NewAccumulation(desc, &ls, res, agg)
}

func testProcessor(
	t *testing.T,
	ekind export.ExportKind,
	mkind metric.InstrumentKind,
	nkind number.Kind,
	akind aggregation.Kind,
) {
	// Note: this selector uses the instrument name to dictate
	// aggregation kind.
	selector := processorTest.AggregatorSelector()
	res := resource.NewWithAttributes(attribute.String("R", "V"))

	labs1 := []attribute.KeyValue{attribute.String("L1", "V")}
	labs2 := []attribute.KeyValue{attribute.String("L2", "V")}

	testBody := func(t *testing.T, hasMemory bool, nAccum, nCheckpoint int) {
		processor := basic.New(selector, export.ConstantExportKindSelector(ekind), basic.WithMemory(hasMemory))

		instSuffix := fmt.Sprint(".", strings.ToLower(akind.String()))

		desc1 := metric.NewDescriptor(fmt.Sprint("inst1", instSuffix), mkind, nkind)
		desc2 := metric.NewDescriptor(fmt.Sprint("inst2", instSuffix), mkind, nkind)

		for nc := 0; nc < nCheckpoint; nc++ {

			// The input is 10 per update, scaled by
			// the number of checkpoints for
			// cumulative instruments:
			input := int64(10)
			cumulativeMultiplier := int64(nc + 1)
			if mkind.PrecomputedSum() {
				input *= cumulativeMultiplier
			}

			processor.StartCollection()

			for na := 0; na < nAccum; na++ {
				_ = processor.Process(updateFor(t, &desc1, selector, res, input, labs1...))
				_ = processor.Process(updateFor(t, &desc2, selector, res, input, labs2...))
			}

			err := processor.FinishCollection()
			if err == aggregation.ErrNoSubtraction {
				var subr export.Aggregator
				selector.AggregatorFor(&desc1, &subr)
				_, canSub := subr.(export.Subtractor)

				// Allow unsupported subraction case only when it is called for.
				require.True(t, mkind.PrecomputedSum() && ekind == export.DeltaExportKind && !canSub)
				return
			} else if err != nil {
				t.Fatal("unexpected FinishCollection error: ", err)
			}

			if nc < nCheckpoint-1 {
				continue
			}

			checkpointSet := processor.CheckpointSet()

			for _, repetitionAfterEmptyInterval := range []bool{false, true} {
				if repetitionAfterEmptyInterval {
					// We're repeating the test after another
					// interval with no updates.
					processor.StartCollection()
					if err := processor.FinishCollection(); err != nil {
						t.Fatal("unexpected collection error: ", err)
					}
				}

				// Test the final checkpoint state.
				records1 := processorTest.NewOutput(attribute.DefaultEncoder())
				err = checkpointSet.ForEach(export.ConstantExportKindSelector(ekind), records1.AddRecord)

				// Test for an allowed error:
				if err != nil && err != aggregation.ErrNoSubtraction {
					t.Fatal("unexpected checkpoint error: ", err)
				}
				var multiplier int64

				if mkind.Asynchronous() {
					// Asynchronous tests accumulate results multiply by the
					// number of Accumulators, unless LastValue aggregation.
					// If a precomputed sum, we expect cumulative inputs.
					if mkind.PrecomputedSum() {
						if ekind == export.DeltaExportKind && akind != aggregation.LastValueKind {
							multiplier = int64(nAccum)
						} else if akind == aggregation.LastValueKind {
							multiplier = cumulativeMultiplier
						} else {
							multiplier = cumulativeMultiplier * int64(nAccum)
						}
					} else {
						if ekind == export.CumulativeExportKind && akind != aggregation.LastValueKind {
							multiplier = cumulativeMultiplier * int64(nAccum)
						} else if akind == aggregation.LastValueKind {
							multiplier = 1
						} else {
							multiplier = int64(nAccum)
						}
					}
				} else {
					// Synchronous accumulate results from multiple accumulators,
					// use that number as the baseline multiplier.
					multiplier = int64(nAccum)
					if ekind == export.CumulativeExportKind {
						// If a cumulative exporter, include prior checkpoints.
						multiplier *= cumulativeMultiplier
					}
					if akind == aggregation.LastValueKind {
						// If a last-value aggregator, set multiplier to 1.0.
						multiplier = 1
					}
				}

				exp := map[string]float64{}
				if hasMemory || !repetitionAfterEmptyInterval {
					exp = map[string]float64{
						fmt.Sprintf("inst1%s/L1=V/R=V", instSuffix): float64(multiplier * 10), // labels1
						fmt.Sprintf("inst2%s/L2=V/R=V", instSuffix): float64(multiplier * 10), // labels2
					}
				}

				require.EqualValues(t, exp, records1.Map(), "with repetition=%v", repetitionAfterEmptyInterval)
			}
		}
	}

	for _, hasMem := range []bool{false, true} {
		t.Run(fmt.Sprintf("HasMemory=%v", hasMem), func(t *testing.T) {
			// For 1 to 3 checkpoints:
			for nAccum := 1; nAccum <= 3; nAccum++ {
				t.Run(fmt.Sprintf("NumAccum=%d", nAccum), func(t *testing.T) {
					// For 1 to 3 accumulators:
					for nCheckpoint := 1; nCheckpoint <= 3; nCheckpoint++ {
						t.Run(fmt.Sprintf("NumCkpt=%d", nCheckpoint), func(t *testing.T) {
							testBody(t, hasMem, nAccum, nCheckpoint)
						})
					}
				})
			}
		})
	}
}

type bogusExporter struct{}

func (bogusExporter) ExportKindFor(*metric.Descriptor, aggregation.Kind) export.ExportKind {
	return 1000000
}

func (bogusExporter) Export(context.Context, export.CheckpointSet) error {
	panic("Not called")
}

func TestBasicInconsistent(t *testing.T) {
	// Test double-start
	b := basic.New(processorTest.AggregatorSelector(), export.StatelessExportKindSelector())

	b.StartCollection()
	b.StartCollection()
	require.Equal(t, basic.ErrInconsistentState, b.FinishCollection())

	// Test finish without start
	b = basic.New(processorTest.AggregatorSelector(), export.StatelessExportKindSelector())

	require.Equal(t, basic.ErrInconsistentState, b.FinishCollection())

	// Test no finish
	b = basic.New(processorTest.AggregatorSelector(), export.StatelessExportKindSelector())

	b.StartCollection()
	require.Equal(
		t,
		basic.ErrInconsistentState,
		b.ForEach(
			export.StatelessExportKindSelector(),
			func(export.Record) error { return nil },
		),
	)

	// Test no start
	b = basic.New(processorTest.AggregatorSelector(), export.StatelessExportKindSelector())

	desc := metric.NewDescriptor("inst", metric.CounterInstrumentKind, number.Int64Kind)
	accum := export.NewAccumulation(&desc, attribute.EmptySet(), resource.Empty(), metrictest.NoopAggregator{})
	require.Equal(t, basic.ErrInconsistentState, b.Process(accum))

	// Test invalid kind:
	b = basic.New(processorTest.AggregatorSelector(), export.StatelessExportKindSelector())
	b.StartCollection()
	require.NoError(t, b.Process(accum))
	require.NoError(t, b.FinishCollection())

	err := b.ForEach(
		bogusExporter{},
		func(export.Record) error { return nil },
	)
	require.True(t, errors.Is(err, basic.ErrInvalidExportKind))

}

func TestBasicTimestamps(t *testing.T) {
	beforeNew := time.Now()
	time.Sleep(time.Nanosecond)
	b := basic.New(processorTest.AggregatorSelector(), export.StatelessExportKindSelector())
	time.Sleep(time.Nanosecond)
	afterNew := time.Now()

	desc := metric.NewDescriptor("inst", metric.CounterInstrumentKind, number.Int64Kind)
	accum := export.NewAccumulation(&desc, attribute.EmptySet(), resource.Empty(), metrictest.NoopAggregator{})

	b.StartCollection()
	_ = b.Process(accum)
	require.NoError(t, b.FinishCollection())

	var start1, end1 time.Time

	require.NoError(t, b.ForEach(export.StatelessExportKindSelector(), func(rec export.Record) error {
		start1 = rec.StartTime()
		end1 = rec.EndTime()
		return nil
	}))

	// The first start time is set in the constructor.
	requireNotAfter(t, beforeNew, start1)
	requireNotAfter(t, start1, afterNew)

	for i := 0; i < 2; i++ {
		b.StartCollection()
		require.NoError(t, b.Process(accum))
		require.NoError(t, b.FinishCollection())

		var start2, end2 time.Time

		require.NoError(t, b.ForEach(export.StatelessExportKindSelector(), func(rec export.Record) error {
			start2 = rec.StartTime()
			end2 = rec.EndTime()
			return nil
		}))

		// Subsequent intervals have their start and end aligned.
		require.Equal(t, start2, end1)
		requireNotAfter(t, start1, end1)
		requireNotAfter(t, start2, end2)

		start1 = start2
		end1 = end2
	}
}

func TestStatefulNoMemoryCumulative(t *testing.T) {
	res := resource.NewWithAttributes(attribute.String("R", "V"))
	ekindSel := export.CumulativeExportKindSelector()

	desc := metric.NewDescriptor("inst.sum", metric.CounterInstrumentKind, number.Int64Kind)
	selector := processorTest.AggregatorSelector()

	processor := basic.New(selector, ekindSel, basic.WithMemory(false))
	checkpointSet := processor.CheckpointSet()

	for i := 1; i < 3; i++ {
		// Empty interval
		processor.StartCollection()
		require.NoError(t, processor.FinishCollection())

		// Verify zero elements
		records := processorTest.NewOutput(attribute.DefaultEncoder())
		require.NoError(t, checkpointSet.ForEach(ekindSel, records.AddRecord))
		require.EqualValues(t, map[string]float64{}, records.Map())

		// Add 10
		processor.StartCollection()
		_ = processor.Process(updateFor(t, &desc, selector, res, 10, attribute.String("A", "B")))
		require.NoError(t, processor.FinishCollection())

		// Verify one element
		records = processorTest.NewOutput(attribute.DefaultEncoder())
		require.NoError(t, checkpointSet.ForEach(ekindSel, records.AddRecord))
		require.EqualValues(t, map[string]float64{
			"inst.sum/A=B/R=V": float64(i * 10),
		}, records.Map())
	}
}

func TestStatefulNoMemoryDelta(t *testing.T) {
	res := resource.NewWithAttributes(attribute.String("R", "V"))
	ekindSel := export.DeltaExportKindSelector()

	desc := metric.NewDescriptor("inst.sum", metric.SumObserverInstrumentKind, number.Int64Kind)
	selector := processorTest.AggregatorSelector()

	processor := basic.New(selector, ekindSel, basic.WithMemory(false))
	checkpointSet := processor.CheckpointSet()

	for i := 1; i < 3; i++ {
		// Empty interval
		processor.StartCollection()
		require.NoError(t, processor.FinishCollection())

		// Verify zero elements
		records := processorTest.NewOutput(attribute.DefaultEncoder())
		require.NoError(t, checkpointSet.ForEach(ekindSel, records.AddRecord))
		require.EqualValues(t, map[string]float64{}, records.Map())

		// Add 10
		processor.StartCollection()
		_ = processor.Process(updateFor(t, &desc, selector, res, int64(i*10), attribute.String("A", "B")))
		require.NoError(t, processor.FinishCollection())

		// Verify one element
		records = processorTest.NewOutput(attribute.DefaultEncoder())
		require.NoError(t, checkpointSet.ForEach(ekindSel, records.AddRecord))
		require.EqualValues(t, map[string]float64{
			"inst.sum/A=B/R=V": 10,
		}, records.Map())
	}
}

func TestMultiObserverSum(t *testing.T) {
	for _, ekindSel := range []export.ExportKindSelector{
		export.CumulativeExportKindSelector(),
		export.DeltaExportKindSelector(),
	} {

		res := resource.NewWithAttributes(attribute.String("R", "V"))
		desc := metric.NewDescriptor("observe.sum", metric.SumObserverInstrumentKind, number.Int64Kind)
		selector := processorTest.AggregatorSelector()

		processor := basic.New(selector, ekindSel, basic.WithMemory(false))
		checkpointSet := processor.CheckpointSet()

		for i := 1; i < 3; i++ {
			// Add i*10*3 times
			processor.StartCollection()
			_ = processor.Process(updateFor(t, &desc, selector, res, int64(i*10), attribute.String("A", "B")))
			_ = processor.Process(updateFor(t, &desc, selector, res, int64(i*10), attribute.String("A", "B")))
			_ = processor.Process(updateFor(t, &desc, selector, res, int64(i*10), attribute.String("A", "B")))
			require.NoError(t, processor.FinishCollection())

			// Multiplier is 1 for deltas, otherwise i.
			multiplier := i
			if ekindSel.ExportKindFor(&desc, aggregation.SumKind) == export.DeltaExportKind {
				multiplier = 1
			}

			// Verify one element
			records := processorTest.NewOutput(attribute.DefaultEncoder())
			require.NoError(t, checkpointSet.ForEach(ekindSel, records.AddRecord))
			require.EqualValues(t, map[string]float64{
				"observe.sum/A=B/R=V": float64(3 * 10 * multiplier),
			}, records.Map())
		}
	}
}

func TestSumObserverEndToEnd(t *testing.T) {
	ctx := context.Background()
	eselector := export.CumulativeExportKindSelector()
	proc := basic.New(
		processorTest.AggregatorSelector(),
		eselector,
	)
	accum := sdk.NewAccumulator(proc, resource.Empty())
	meter := metric.WrapMeterImpl(accum, "testing")

	var calls int64
	metric.Must(meter).NewInt64SumObserver("observer.sum",
		func(_ context.Context, result metric.Int64ObserverResult) {
			calls++
			result.Observe(calls)
		},
	)
	data := proc.CheckpointSet()

	var startTime [3]time.Time
	var endTime [3]time.Time

	for i := range startTime {
		data.Lock()
		proc.StartCollection()
		accum.Collect(ctx)
		require.NoError(t, proc.FinishCollection())

		exporter := processortest.NewExporter(eselector, attribute.DefaultEncoder())
		require.NoError(t, exporter.Export(ctx, data))

		require.EqualValues(t, map[string]float64{
			"observer.sum//": float64(i + 1),
		}, exporter.Values())

		var record export.Record
		require.NoError(t, data.ForEach(eselector, func(r export.Record) error {
			record = r
			return nil
		}))

		startTime[i] = record.StartTime()
		endTime[i] = record.EndTime()
		data.Unlock()
	}

	require.Equal(t, startTime[0], startTime[1])
	require.Equal(t, startTime[0], startTime[2])
	requireNotAfter(t, endTime[0], endTime[1])
	requireNotAfter(t, endTime[1], endTime[2])
}
