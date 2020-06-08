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
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/integrator/simple"
	"go.opentelemetry.io/otel/sdk/metric/integrator/test"
	"go.opentelemetry.io/otel/sdk/resource"
)

// TestIntegrator tests all the non-error paths in this package.
func TestIntegrator(t *testing.T) {
	type exportCase struct {
		kind export.ExporterKind
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

func (ts testSelector) AggregatorFor(desc *metric.Descriptor) export.Aggregator {
	switch ts.kind {
	case aggregation.SumKind:
		return sum.New()
	case aggregation.MinMaxSumCountKind:
		return minmaxsumcount.New(desc)
	case aggregation.HistogramKind:
		return histogram.New(desc, nil)
	case aggregation.LastValueKind:
		return lastvalue.New()
	case aggregation.SketchKind:
		return ddsketch.New(desc, nil)
	case aggregation.ExactKind:
		return array.New()
	}
	panic("Unknown aggregation kind")
}

func testSynchronousIntegration(
	t *testing.T,
	ekind export.ExporterKind,
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
		agg := selector.AggregatorFor(desc)
		_ = agg.Update(ctx, asNumber(value), desc)
		agg.Checkpoint(desc)

		//fmt.Printf("AGGREGATOR %T %v\n", agg, agg)
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

					integrator := simple.New(selector, ekind)

					for nc := 0; nc < NCheckpoint; nc++ {

						// The input is 10 per update, scaled by
						// the number of checkpoints for
						// cumulative instruments:
						input := int64(10)
						cumulativeMultiplier := int64(nc + 1)
						if mkind.Cumulative() {
							input *= cumulativeMultiplier
						}

						for na := 0; na < NAccum; na++ {
							_ = integrator.Process(updateFor(&desc1, input, labs1))
							_ = integrator.Process(updateFor(&desc2, input, labs2))
						}

						checkpointSet := integrator.CheckpointSet()

						if nc < NCheckpoint-1 {
							integrator.FinishedCollection()
							continue
						}

						// Test the final checkpoint state.
						records1 := test.NewOutput(label.DefaultEncoder())
						err := checkpointSet.ForEach(ekind, records1.AddRecord)

						// Test for an allowed error:
						if err != nil && err != aggregation.ErrNoSubtraction {
							t.Fatal("unexpected checkpoint error")
						}
						if err == aggregation.ErrNoSubtraction {
							_, canSub := selector.AggregatorFor(&desc1).(export.Subtractor)

							// Allow unsupported subraction case only when it is called for.
							require.True(t, mkind.Cumulative() && ekind == export.DeltaExporter && !canSub)
							return
						}

						var multiplier int64

						if mkind.Asynchronous() {
							// Because async instruments take the last value,
							// the number of accumulators doesn't matter.
							if mkind.Cumulative() {
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

						integrator.FinishedCollection()
					}
				})
			}
		})
	}
}
