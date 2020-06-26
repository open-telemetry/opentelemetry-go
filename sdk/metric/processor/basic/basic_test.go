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
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	exportTest "go.opentelemetry.io/otel/exporters/metric/test"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/processor/test"
	"go.opentelemetry.io/otel/sdk/resource"
)

// TestProcessor tests all the non-error paths in this package.
func TestProcessor(t *testing.T) {
	type exportCase struct {
		kind export.ExportKind
	}
	type instrumentCase struct {
		kind metric.Kind
	}
	type numberCase struct {
		kind metric.NumberKind
	}
	type aggregatorCase struct {
		kind aggregation.Kind
	}

	for _, tc := range []exportCase{
		{kind: export.PassThroughExporter},
		{kind: export.CumulativeExporter},
		{kind: export.DeltaExporter},
	} {
		t.Run(tc.kind.String(), func(t *testing.T) {
			for _, ic := range []instrumentCase{
				{kind: metric.CounterKind},
				{kind: metric.UpDownCounterKind},
				{kind: metric.ValueRecorderKind},
				{kind: metric.SumObserverKind},
				{kind: metric.UpDownSumObserverKind},
				{kind: metric.ValueObserverKind},
			} {
				t.Run(ic.kind.String(), func(t *testing.T) {
					for _, nc := range []numberCase{
						{kind: metric.Int64NumberKind},
						{kind: metric.Float64NumberKind},
					} {
						t.Run(nc.kind.String(), func(t *testing.T) {
							for _, ac := range []aggregatorCase{
								{kind: aggregation.SumKind},
								{kind: aggregation.MinMaxSumCountKind},
								{kind: aggregation.HistogramKind},
								{kind: aggregation.LastValueKind},
								{kind: aggregation.ExactKind},
								{kind: aggregation.SketchKind},
							} {
								t.Run(ac.kind.String(), func(t *testing.T) {
									testSynchronousIntegration(
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

type testSelector struct {
	kind aggregation.Kind
}

func (ts testSelector) AggregatorFor(desc *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	for i := range aggPtrs {
		switch ts.kind {
		case aggregation.SumKind:
			*aggPtrs[i] = &sum.New(1)[0]
		case aggregation.MinMaxSumCountKind:
			*aggPtrs[i] = &minmaxsumcount.New(1, desc)[0]
		case aggregation.HistogramKind:
			*aggPtrs[i] = &histogram.New(1, desc, nil)[0]
		case aggregation.LastValueKind:
			*aggPtrs[i] = &lastvalue.New(1)[0]
		case aggregation.SketchKind:
			*aggPtrs[i] = &ddsketch.New(1, desc, nil)[0]
		case aggregation.ExactKind:
			*aggPtrs[i] = &array.New(1)[0]
		}
	}
}

func testSynchronousIntegration(
	t *testing.T,
	ekind export.ExportKind,
	mkind metric.Kind,
	nkind metric.NumberKind,
	akind aggregation.Kind,
) {
	ctx := context.Background()
	selector := testSelector{akind}
	res := resource.New(kv.String("R", "V"))

	asNumber := func(value int64) metric.Number {
		if nkind == metric.Int64NumberKind {
			return metric.NewInt64Number(value)
		}
		return metric.NewFloat64Number(float64(value))
	}

	updateFor := func(desc *metric.Descriptor, value int64, labs []kv.KeyValue) export.Accumulation {
		ls := label.NewSet(labs...)
		var agg export.Aggregator
		selector.AggregatorFor(desc, &agg)
		_ = agg.Update(ctx, asNumber(value), desc)

		return export.NewAccumulation(desc, &ls, res, agg)
	}

	labs1 := []kv.KeyValue{kv.String("L1", "V")}
	labs2 := []kv.KeyValue{kv.String("L2", "V")}

	desc1 := metric.NewDescriptor("inst1", mkind, nkind)
	desc2 := metric.NewDescriptor("inst2", mkind, nkind)

	// For 1 to 3 checkpoints:
	for NAccum := 1; NAccum <= 3; NAccum++ {
		t.Run(fmt.Sprintf("NumAccum=%d", NAccum), func(t *testing.T) {
			// For 1 to 3 accumulators:
			for NCheckpoint := 1; NCheckpoint <= 3; NCheckpoint++ {
				t.Run(fmt.Sprintf("NumCkpt=%d", NCheckpoint), func(t *testing.T) {

					processor := basic.New(selector, ekind)

					for nc := 0; nc < NCheckpoint; nc++ {

						// The input is 10 per update, scaled by
						// the number of checkpoints for
						// cumulative instruments:
						input := int64(10)
						cumulativeMultiplier := int64(nc + 1)
						if mkind.PrecomputedSum() {
							input *= cumulativeMultiplier
						}

						processor.StartCollection()

						for na := 0; na < NAccum; na++ {
							_ = processor.Process(updateFor(&desc1, input, labs1))
							_ = processor.Process(updateFor(&desc2, input, labs2))
						}

						err := processor.FinishCollection()
						if err == aggregation.ErrNoSubtraction {
							var subr export.Aggregator
							selector.AggregatorFor(&desc1, &subr)
							_, canSub := subr.(export.Subtractor)

							// Allow unsupported subraction case only when it is called for.
							require.True(t, mkind.PrecomputedSum() && ekind == export.DeltaExporter && !canSub)
							return
						} else if err != nil {
							t.Fatal(fmt.Sprint("unexpected FinishCollection error: ", err))
						}

						if nc < NCheckpoint-1 {
							continue
						}

						checkpointSet := processor.CheckpointSet()

						// Test the final checkpoint state.
						records1 := test.NewOutput(label.DefaultEncoder())
						err = checkpointSet.ForEach(ekind, records1.AddRecord)

						// Test for an allowed error:
						if err != nil && err != aggregation.ErrNoSubtraction {
							t.Fatal(fmt.Sprint("unexpected checkpoint error: ", err))
						}
						var multiplier int64

						if mkind.Asynchronous() {
							// Because async instruments take the last value,
							// the number of accumulators doesn't matter.
							if mkind.PrecomputedSum() {
								if ekind == export.DeltaExporter {
									multiplier = 1
								} else {
									multiplier = cumulativeMultiplier
								}
							} else {
								if ekind == export.CumulativeExporter && akind != aggregation.LastValueKind {
									multiplier = cumulativeMultiplier
								} else {
									multiplier = 1
								}
							}
						} else {
							// Synchronous accumulate results from multiple accumulators,
							// use that number as the baseline multiplier.
							multiplier = int64(NAccum)
							if ekind == export.CumulativeExporter {
								// If a cumulative exporter, include prior checkpoints.
								multiplier *= cumulativeMultiplier
							}
							if akind == aggregation.LastValueKind {
								// If a last-value aggregator, set multiplier to 1.0.
								multiplier = 1
							}
						}

						require.EqualValues(t, map[string]float64{
							"inst1/L1=V/R=V": float64(multiplier * 10), // labels1
							"inst2/L2=V/R=V": float64(multiplier * 10), // labels2
						}, records1.Map)
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
	b := basic.New(test.AggregatorSelector(), export.PassThroughExporter)

	b.StartCollection()
	b.StartCollection()
	require.Equal(t, basic.ErrInconsistentState, b.FinishCollection())

	// Test finish without start
	b = basic.New(test.AggregatorSelector(), export.PassThroughExporter)

	require.Equal(t, basic.ErrInconsistentState, b.FinishCollection())

	// Test no finish
	b = basic.New(test.AggregatorSelector(), export.PassThroughExporter)

	b.StartCollection()
	require.Equal(
		t,
		basic.ErrInconsistentState,
		b.ForEach(
			export.PassThroughExporter,
			func(export.Record) error { return nil },
		),
	)

	// Test no start
	b = basic.New(test.AggregatorSelector(), export.PassThroughExporter)

	desc := metric.NewDescriptor("inst", metric.CounterKind, metric.Int64NumberKind)
	accum := export.NewAccumulation(&desc, label.EmptySet(), resource.Empty(), exportTest.NoopAggregator{})
	require.Equal(t, basic.ErrInconsistentState, b.Process(accum))

	// Test invalid kind:
	b = basic.New(test.AggregatorSelector(), export.PassThroughExporter)
	b.StartCollection()
	require.NoError(t, b.Process(accum))
	require.NoError(t, b.FinishCollection())

	err := b.ForEach(
		bogusExporter{},
		func(export.Record) error { return nil },
	)
	require.True(t, errors.Is(err, basic.ErrInvalidExporterKind))

}

func TestBasicTimestamps(t *testing.T) {
	beforeNew := time.Now()
	b := basic.New(test.AggregatorSelector(), export.PassThroughExporter)
	afterNew := time.Now()

	desc := metric.NewDescriptor("inst", metric.CounterKind, metric.Int64NumberKind)
	accum := export.NewAccumulation(&desc, label.EmptySet(), resource.Empty(), exportTest.NoopAggregator{})

	b.StartCollection()
	_ = b.Process(accum)
	require.NoError(t, b.FinishCollection())

	var start1, end1 time.Time

	require.NoError(t, b.ForEach(export.PassThroughExporter, func(rec export.Record) error {
		start1 = rec.StartTime()
		end1 = rec.EndTime()
		return nil
	}))

	// The first start time is set in the constructor.
	require.True(t, beforeNew.Before(start1))
	require.True(t, afterNew.After(start1))

	for i := 0; i < 2; i++ {
		b.StartCollection()
		require.NoError(t, b.Process(accum))
		require.NoError(t, b.FinishCollection())

		var start2, end2 time.Time

		require.NoError(t, b.ForEach(export.PassThroughExporter, func(rec export.Record) error {
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
