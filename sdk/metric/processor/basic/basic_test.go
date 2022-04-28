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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregatortest"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	processorTest "go.opentelemetry.io/otel/sdk/metric/processor/processortest"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/resource"
)

func requireNotAfter(t *testing.T, t1, t2 time.Time) {
	require.False(t, t1.After(t2), "expected %v ≤ %v", t1, t2)
}

// TestProcessor tests all the non-error paths in this package.
func TestProcessor(t *testing.T) {
	type exportCase struct {
		kind aggregation.Temporality
	}
	type instrumentCase struct {
		kind sdkapi.InstrumentKind
	}
	type numberCase struct {
		kind number.Kind
	}
	type aggregatorCase struct {
		kind aggregation.Kind
	}

	for _, tc := range []exportCase{
		{kind: aggregation.CumulativeTemporality},
		{kind: aggregation.DeltaTemporality},
	} {
		t.Run(tc.kind.String(), func(t *testing.T) {
			for _, ic := range []instrumentCase{
				{kind: sdkapi.CounterInstrumentKind},
				{kind: sdkapi.UpDownCounterInstrumentKind},
				{kind: sdkapi.HistogramInstrumentKind},
				{kind: sdkapi.CounterObserverInstrumentKind},
				{kind: sdkapi.UpDownCounterObserverInstrumentKind},
				{kind: sdkapi.GaugeObserverInstrumentKind},
			} {
				t.Run(ic.kind.String(), func(t *testing.T) {
					for _, nc := range []numberCase{
						{kind: number.Int64Kind},
						{kind: number.Float64Kind},
					} {
						t.Run(nc.kind.String(), func(t *testing.T) {
							for _, ac := range []aggregatorCase{
								{kind: aggregation.SumKind},
								{kind: aggregation.HistogramKind},
								{kind: aggregation.LastValueKind},
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

func updateFor(t *testing.T, desc *sdkapi.Descriptor, selector export.AggregatorSelector, value int64, labs ...attribute.KeyValue) export.Accumulation {
	ls := attribute.NewSet(labs...)
	var agg aggregator.Aggregator
	selector.AggregatorFor(desc, &agg)
	require.NoError(t, agg.Update(context.Background(), asNumber(desc.NumberKind(), value), desc))

	return export.NewAccumulation(desc, &ls, agg)
}

func testProcessor(
	t *testing.T,
	aggTemp aggregation.Temporality,
	mkind sdkapi.InstrumentKind,
	nkind number.Kind,
	akind aggregation.Kind,
) {
	// This code tests for errors when the export kind is Delta
	// and the instrument kind is PrecomputedSum().
	expectConversion := !(aggTemp == aggregation.DeltaTemporality && mkind.PrecomputedSum())
	requireConversion := func(t *testing.T, err error) {
		if expectConversion {
			require.NoError(t, err)
		} else {
			require.Equal(t, aggregation.ErrNoCumulativeToDelta, err)
		}
	}

	// Note: this selector uses the instrument name to dictate
	// aggregation kind.
	selector := processorTest.AggregatorSelector()

	labs1 := []attribute.KeyValue{attribute.String("L1", "V")}
	labs2 := []attribute.KeyValue{attribute.String("L2", "V")}

	testBody := func(t *testing.T, hasMemory bool, nAccum, nCheckpoint int) {
		processor := basic.New(selector, aggregation.ConstantTemporalitySelector(aggTemp), basic.WithMemory(hasMemory))

		instSuffix := fmt.Sprint(".", strings.ToLower(akind.String()))

		desc1 := metrictest.NewDescriptor(fmt.Sprint("inst1", instSuffix), mkind, nkind)
		desc2 := metrictest.NewDescriptor(fmt.Sprint("inst2", instSuffix), mkind, nkind)

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
				requireConversion(t, processor.Process(updateFor(t, &desc1, selector, input, labs1...)))
				requireConversion(t, processor.Process(updateFor(t, &desc2, selector, input, labs2...)))
			}

			// Note: in case of !expectConversion, we still get no error here
			// because the Process() skipped entering state for those records.
			require.NoError(t, processor.FinishCollection())

			if nc < nCheckpoint-1 {
				continue
			}

			reader := processor.Reader()

			for _, repetitionAfterEmptyInterval := range []bool{false, true} {
				if repetitionAfterEmptyInterval {
					// We're repeating the test after another
					// interval with no updates.
					processor.StartCollection()
					require.NoError(t, processor.FinishCollection())
				}

				// Test the final checkpoint state.
				records1 := processorTest.NewOutput(attribute.DefaultEncoder())
				require.NoError(t, reader.ForEach(aggregation.ConstantTemporalitySelector(aggTemp), records1.AddRecord))

				if !expectConversion {
					require.EqualValues(t, map[string]float64{}, records1.Map())
					continue
				}

				var multiplier int64

				if mkind.Asynchronous() {
					// Asynchronous tests accumulate results multiply by the
					// number of Accumulators, unless LastValue aggregation.
					// If a precomputed sum, we expect cumulative inputs.
					if mkind.PrecomputedSum() {
						require.NotEqual(t, aggTemp, aggregation.DeltaTemporality)
						if akind == aggregation.LastValueKind {
							multiplier = cumulativeMultiplier
						} else {
							multiplier = cumulativeMultiplier * int64(nAccum)
						}
					} else {
						if aggTemp == aggregation.CumulativeTemporality && akind != aggregation.LastValueKind {
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
					if aggTemp == aggregation.CumulativeTemporality {
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
						fmt.Sprintf("inst1%s/L1=V/", instSuffix): float64(multiplier * 10), // attrs1
						fmt.Sprintf("inst2%s/L2=V/", instSuffix): float64(multiplier * 10), // attrs2
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

func (bogusExporter) TemporalityFor(*sdkapi.Descriptor, aggregation.Kind) aggregation.Temporality {
	return 100
}

func (bogusExporter) Export(context.Context, export.Reader) error {
	panic("Not called")
}

func TestBasicInconsistent(t *testing.T) {
	// Test double-start
	b := basic.New(processorTest.AggregatorSelector(), aggregation.StatelessTemporalitySelector())

	b.StartCollection()
	b.StartCollection()
	require.Equal(t, basic.ErrInconsistentState, b.FinishCollection())

	// Test finish without start
	b = basic.New(processorTest.AggregatorSelector(), aggregation.StatelessTemporalitySelector())

	require.Equal(t, basic.ErrInconsistentState, b.FinishCollection())

	// Test no finish
	b = basic.New(processorTest.AggregatorSelector(), aggregation.StatelessTemporalitySelector())

	b.StartCollection()
	require.Equal(
		t,
		basic.ErrInconsistentState,
		b.ForEach(
			aggregation.StatelessTemporalitySelector(),
			func(export.Record) error { return nil },
		),
	)

	// Test no start
	b = basic.New(processorTest.AggregatorSelector(), aggregation.StatelessTemporalitySelector())

	desc := metrictest.NewDescriptor("inst", sdkapi.CounterInstrumentKind, number.Int64Kind)
	accum := export.NewAccumulation(&desc, attribute.EmptySet(), aggregatortest.NoopAggregator{})
	require.Equal(t, basic.ErrInconsistentState, b.Process(accum))

	// Test invalid kind:
	b = basic.New(processorTest.AggregatorSelector(), aggregation.StatelessTemporalitySelector())
	b.StartCollection()
	require.NoError(t, b.Process(accum))
	require.NoError(t, b.FinishCollection())

	err := b.ForEach(
		bogusExporter{},
		func(export.Record) error { return nil },
	)
	require.True(t, errors.Is(err, basic.ErrInvalidTemporality))

}

func TestBasicTimestamps(t *testing.T) {
	beforeNew := time.Now()
	time.Sleep(time.Nanosecond)
	b := basic.New(processorTest.AggregatorSelector(), aggregation.StatelessTemporalitySelector())
	time.Sleep(time.Nanosecond)
	afterNew := time.Now()

	desc := metrictest.NewDescriptor("inst", sdkapi.CounterInstrumentKind, number.Int64Kind)
	accum := export.NewAccumulation(&desc, attribute.EmptySet(), aggregatortest.NoopAggregator{})

	b.StartCollection()
	_ = b.Process(accum)
	require.NoError(t, b.FinishCollection())

	var start1, end1 time.Time

	require.NoError(t, b.ForEach(aggregation.StatelessTemporalitySelector(), func(rec export.Record) error {
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

		require.NoError(t, b.ForEach(aggregation.StatelessTemporalitySelector(), func(rec export.Record) error {
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
	aggTempSel := aggregation.CumulativeTemporalitySelector()

	desc := metrictest.NewDescriptor("inst.sum", sdkapi.CounterInstrumentKind, number.Int64Kind)
	selector := processorTest.AggregatorSelector()

	processor := basic.New(selector, aggTempSel, basic.WithMemory(false))
	reader := processor.Reader()

	for i := 1; i < 3; i++ {
		// Empty interval
		processor.StartCollection()
		require.NoError(t, processor.FinishCollection())

		// Verify zero elements
		records := processorTest.NewOutput(attribute.DefaultEncoder())
		require.NoError(t, reader.ForEach(aggTempSel, records.AddRecord))
		require.EqualValues(t, map[string]float64{}, records.Map())

		// Add 10
		processor.StartCollection()
		_ = processor.Process(updateFor(t, &desc, selector, 10, attribute.String("A", "B")))
		require.NoError(t, processor.FinishCollection())

		// Verify one element
		records = processorTest.NewOutput(attribute.DefaultEncoder())
		require.NoError(t, reader.ForEach(aggTempSel, records.AddRecord))
		require.EqualValues(t, map[string]float64{
			"inst.sum/A=B/": float64(i * 10),
		}, records.Map())
	}
}

func TestMultiObserverSum(t *testing.T) {
	for _, test := range []struct {
		name string
		aggregation.TemporalitySelector
		expectProcessErr error
	}{
		{"cumulative", aggregation.CumulativeTemporalitySelector(), nil},
		{"delta", aggregation.DeltaTemporalitySelector(), aggregation.ErrNoCumulativeToDelta},
	} {
		t.Run(test.name, func(t *testing.T) {
			aggTempSel := test.TemporalitySelector
			desc := metrictest.NewDescriptor("observe.sum", sdkapi.CounterObserverInstrumentKind, number.Int64Kind)
			selector := processorTest.AggregatorSelector()

			processor := basic.New(selector, aggTempSel, basic.WithMemory(false))
			reader := processor.Reader()

			for i := 1; i < 3; i++ {
				// Add i*10*3 times
				processor.StartCollection()
				require.True(t, errors.Is(processor.Process(updateFor(t, &desc, selector, int64(i*10), attribute.String("A", "B"))), test.expectProcessErr))
				require.True(t, errors.Is(processor.Process(updateFor(t, &desc, selector, int64(i*10), attribute.String("A", "B"))), test.expectProcessErr))
				require.True(t, errors.Is(processor.Process(updateFor(t, &desc, selector, int64(i*10), attribute.String("A", "B"))), test.expectProcessErr))
				require.NoError(t, processor.FinishCollection())

				// Verify one element
				records := processorTest.NewOutput(attribute.DefaultEncoder())
				if test.expectProcessErr == nil {
					require.NoError(t, reader.ForEach(aggTempSel, records.AddRecord))
					require.EqualValues(t, map[string]float64{
						"observe.sum/A=B/": float64(3 * 10 * i),
					}, records.Map())
				} else {
					require.NoError(t, reader.ForEach(aggTempSel, records.AddRecord))
					require.EqualValues(t, map[string]float64{}, records.Map())
				}
			}
		})
	}
}

func TestCounterObserverEndToEnd(t *testing.T) {
	ctx := context.Background()
	eselector := aggregation.CumulativeTemporalitySelector()
	proc := basic.New(
		processorTest.AggregatorSelector(),
		eselector,
	)
	accum := sdk.NewAccumulator(proc)
	meter := sdkapi.WrapMeterImpl(accum)

	var calls int64
	ctr, err := meter.AsyncInt64().Counter("observer.sum")
	require.NoError(t, err)
	err = meter.RegisterCallback([]instrument.Asynchronous{ctr}, func(ctx context.Context) {
		calls++
		ctr.Observe(ctx, calls)
	})
	require.NoError(t, err)
	reader := proc.Reader()

	var startTime [3]time.Time
	var endTime [3]time.Time

	for i := range startTime {
		data := proc.Reader()
		data.Lock()
		proc.StartCollection()
		accum.Collect(ctx)
		require.NoError(t, proc.FinishCollection())

		exporter := processortest.New(eselector, attribute.DefaultEncoder())
		require.NoError(t, exporter.Export(ctx, resource.Empty(), processortest.OneInstrumentationLibraryReader(
			instrumentation.Library{
				Name: "test",
			}, reader)))

		require.EqualValues(t, map[string]float64{
			"observer.sum//": float64(i + 1),
		}, exporter.Values())

		var record export.Record
		require.NoError(t, data.ForEach(eselector, func(r export.Record) error {
			record = r
			return nil
		}))

		// Try again, but ask for a Delta
		require.Equal(
			t,
			aggregation.ErrNoCumulativeToDelta,
			data.ForEach(
				aggregation.ConstantTemporalitySelector(aggregation.DeltaTemporality),
				func(r export.Record) error {
					t.Fail()
					return nil
				},
			),
		)

		startTime[i] = record.StartTime()
		endTime[i] = record.EndTime()
		data.Unlock()
	}

	require.Equal(t, startTime[0], startTime[1])
	require.Equal(t, startTime[0], startTime[2])
	requireNotAfter(t, endTime[0], endTime[1])
	requireNotAfter(t, endTime[1], endTime[2])
}
