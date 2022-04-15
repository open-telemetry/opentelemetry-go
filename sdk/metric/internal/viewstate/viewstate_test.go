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

package viewstate

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/number/traits"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/metric/views"
)

var (
	testLib = instrumentation.Library{
		Name: "test",
	}

	fooToBarView = []views.View{
		views.New(
			views.MatchInstrumentName("foo"),
			views.WithName("bar"),
		),
	}

	testHistBoundaries = []float64{1, 2, 3}

	altHistogramConfig = aggregator.Config{
		Histogram: aggregator.HistogramConfig{
			ExplicitBoundaries: testHistBoundaries,
		},
	}

	fooToBarAltHistView = []views.View{
		views.New(
			views.MatchInstrumentName("foo"),
			views.WithName("bar"),
			views.WithAggregatorConfig(altHistogramConfig),
		),
	}

	fooToBarFilteredView = []views.View{
		views.New(
			views.MatchInstrumentName("foo"),
			views.WithName("bar"),
			views.WithKeys("a", "b"),
		),
	}

	fooToBarDifferentFiltersView = []views.View{
		views.New(
			views.MatchInstrumentName("foo"),
			views.WithName("bar"),
			views.WithKeys("a", "b"),
		),
		views.New(
			views.MatchInstrumentName("bar"),
			views.WithKeys("a"),
		),
	}

	fooToBarSameFiltersView = []views.View{
		views.New(
			views.MatchInstrumentName("foo"),
			views.WithName("bar"),
			views.WithKeys("a", "b"),
		),
		views.New(
			views.MatchInstrumentName("bar"),
			views.WithKeys("a", "b"),
		),
	}

	dropHistInstView = []views.View{
		views.New(
			views.MatchInstrumentKind(sdkinstrument.HistogramKind),
			views.WithAggregation(aggregation.DropKind),
		),
	}

	instrumentKinds = []sdkinstrument.Kind{
		sdkinstrument.HistogramKind,
		sdkinstrument.GaugeObserverKind,
		sdkinstrument.CounterKind,
		sdkinstrument.UpDownCounterKind,
		sdkinstrument.CounterObserverKind,
		sdkinstrument.UpDownCounterObserverKind,
	}

	numberKinds = []number.Kind{
		number.Int64Kind,
		number.Float64Kind,
	}

	endTime    = time.Now()
	middleTime = endTime.Add(-time.Millisecond)
	startTime  = endTime.Add(-2 * time.Millisecond)

	testSequence = reader.Sequence{
		Start: startTime,
		Last:  middleTime,
		Now:   endTime,
	}
)

// testInst returns a test instrument descriptor similar to what Meter creates.
func testInst(name string, ik sdkinstrument.Kind, nk number.Kind, opts ...instrument.Option) sdkinstrument.Descriptor {
	cfg := instrument.NewConfig(opts...)
	return sdkinstrument.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())
}

func twoTestReaders(opts ...reader.Option) []*reader.Config {
	return []*reader.Config{
		reader.NewConfig(metrictest.NewReader(), opts...),
		reader.NewConfig(metrictest.NewReader(), opts...),
	}
}

func oneTestReader(opts ...reader.Option) []*reader.Config {
	return []*reader.Config{reader.NewConfig(metrictest.NewReader(), opts...)}
}

// TestDeduplicateNoConflict verifies that two identical instruments
// have the same collector.
func TestDeduplicateNoConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.Equal(t, inst1, inst2)
}

// TestDeduplicateRenameNoConflict verifies that one instrument can be renamed
// such that it becomes identical to another, so no conflict.
func TestDeduplicateRenameNoConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, fooToBarView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.Equal(t, inst1, inst2)
}

// TestNoRenameNoConflict verifies that one instrument does not
// conflict with another differently-named instrument.
func TestNoRenameNoConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateNumberConflict verifies that two same instruments
// except different number kind conflict.
func TestDuplicateNumberConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflicts{}))
	require.Equal(t, 2, len(err2.(ViewConflicts)))

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateSyncAsyncConflict verifies that two same instruments
// except one synchonous, one asynchronous conflict.
func TestDuplicateSyncAsyncConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterObserverKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflicts{}))

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateUnitConflict verifies that two same instruments
// except different units conflict.
func TestDuplicateUnitConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind, instrument.WithUnit("gal_us")))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind, instrument.WithUnit("cft_i")))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflicts{}))
	require.Contains(t, err2.Error(), "2 conflicts in 2 readers")
	require.Contains(t, err2.Error(), "conflicts Counter-Float64-Sum-gal_us")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateMonotonicConflict verifies that two same instruments
// except different monotonic values.
func TestDuplicateMonotonicConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.UpDownCounterKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflicts{}))
	require.Contains(t, err2.Error(), "2 conflicts in 2 readers")
	require.Contains(t, err2.Error(), "UpDownCounter-Float64-Sum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregatorConfigConflict verifies that two same instruments
// except different aggregator.Config values.
func TestDuplicateAggregatorConfigConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, fooToBarAltHistView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.HistogramKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflicts{}))
	require.Contains(t, err2.Error(), "different aggregator configuration")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregatorConfigNoConflict verifies that two same instruments
// with same aggregator.Config values configured in different ways.
func TestDuplicateAggregatorConfigNoConflict(t *testing.T) {
	exp := metrictest.NewReader()

	for _, nk := range numberKinds {
		t.Run(nk.String(), func(t *testing.T) {
			rds := []*reader.Config{
				reader.NewConfig(exp, reader.WithDefaultAggregationConfigFunc(
					func(_ sdkinstrument.Kind) (int64Config, float64Config aggregator.Config) {
						if nk == number.Int64Kind {
							return altHistogramConfig, aggregator.Config{}
						}
						return aggregator.Config{}, altHistogramConfig
					},
				)),
			}

			vc := New(testLib, fooToBarAltHistView, rds)

			inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, nk))
			require.NoError(t, err1)
			require.NotNil(t, inst1)

			inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.HistogramKind, nk))
			require.NoError(t, err2)
			require.NotNil(t, inst2)

			require.Equal(t, inst1, inst2)
		})
	}
}

// TestDuplicateAggregationKindConflict verifies that two instruments
// with different aggregation kinds conflict.
func TestDuplicateAggregationKindConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, fooToBarView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflicts{}))
	require.Contains(t, err2.Error(), "2 conflicts in 2 readers")
	require.Contains(t, err2.Error(), "name \"bar\" (original \"foo\") conflicts Histogram-Int64-Histogram, Counter-Int64-Sum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregationKindOneConflict verifies that two
// instruments with different aggregation kinds do not conflict when
// the reader drops the instrument.
func TestDuplicateAggregationKindOneConflict(t *testing.T) {
	// Let one reader drop histograms
	rds := []*reader.Config{
		reader.NewConfig(metrictest.NewReader(), reader.WithDefaultAggregationKindFunc(func(k sdkinstrument.Kind) aggregation.Kind {
			if k == sdkinstrument.HistogramKind {
				return aggregation.DropKind
			}
			return reader.StandardAggregationKind(k)
		})),
		reader.NewConfig(metrictest.NewReader()),
	}

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflicts{}))
	require.Contains(t, err2.Error(), "name \"foo\" conflicts Histogram-Int64-Histogram, Counter-Int64-Sum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregationKindNoConflict verifies that two
// instruments with different aggregation kinds do not conflict when
// the view drops one of the instruments.
func TestDuplicateAggregationKindNoConflict(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, dropHistInstView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, number.Int64Kind))
	require.NoError(t, err1)
	require.Nil(t, inst1) // The viewstate.Instrument is nil, instruments become no-ops.

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err2)
	require.NotNil(t, inst2)
}

// TestDuplicateMultipleConflicts verifies that multiple duplicate
// instrument conflicts include sufficient explanatory information.
func TestDuplicateMultipleConflicts(t *testing.T) {
	rds := oneTestReader()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", instrumentKinds[0], number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	for num, ik := range instrumentKinds[1:] {
		inst2, err2 := vc.Compile(testInst("foo", ik, number.Float64Kind))
		require.Error(t, err2)
		require.NotNil(t, inst2)
		require.True(t, errors.Is(err2, ViewConflicts{}))
		// The total number of conflicting definitions is 1 in
		// the first place and num+1 for the iterations of this loop.
		require.Equal(t, num+2, len(err2.(ViewConflicts)[rds[0]][0].Duplicates))

		if num > 0 {
			require.Contains(t, err2.Error(), fmt.Sprintf("and %d more", num))
		}
	}
}

// TestDuplicateFilterConflicts verifies several cases where
// instruments output the same metric w/ different filters create conflicts.
func TestDuplicateFilterConflicts(t *testing.T) {
	for idx, vws := range [][]views.View{
		fooToBarFilteredView,
		fooToBarDifferentFiltersView,
	} {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			rds := twoTestReaders()

			vc := New(testLib, vws, rds)

			inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
			require.NoError(t, err1)
			require.NotNil(t, inst1)

			inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
			require.Error(t, err2)
			require.NotNil(t, inst2)

			require.True(t, errors.Is(err2, ViewConflicts{}))
			require.Contains(t, err2.Error(), "2 conflicts in 2 readers, e.g.")
			require.Contains(t, err2.Error(), "name \"bar\" (original \"foo\") has conflicts: different attribute filters")
		})
	}
}

// TestDeduplicateSameFilters thests that when one instrument is
// renamed to match another exactly, including filters, they are not
// in conflict.
func TestDeduplicateSameFilters(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, fooToBarSameFiltersView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.Equal(t, inst1, inst2)
}

func int64MonoSum(x int64) aggregation.Sum {
	var s sum.State[int64, traits.Int64, sum.Monotonic]
	var methods sum.Methods[int64, traits.Int64, sum.Monotonic, sum.State[int64, traits.Int64, sum.Monotonic]]
	methods.Init(&s, aggregator.Config{})
	methods.Update(&s, x)
	return &s
}

func float64MonoSum(x float64) aggregation.Sum {
	var s sum.State[float64, traits.Float64, sum.Monotonic]
	var methods sum.Methods[float64, traits.Float64, sum.Monotonic, sum.State[float64, traits.Float64, sum.Monotonic]]
	methods.Init(&s, aggregator.Config{})
	methods.Update(&s, x)
	return &s
}

// TestDuplicatesMergeDescriptor ensures that the longest description string is used.
func TestDuplicatesMergeDescriptor(t *testing.T) {
	rds := oneTestReader()

	vc := New(testLib, fooToBarSameFiltersView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	// This is the winning description:
	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind, instrument.WithDescription("very long")))
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	inst3, err3 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind, instrument.WithDescription("shorter")))
	require.NoError(t, err3)
	require.NotNil(t, inst3)

	require.Equal(t, inst1, inst2)
	require.Equal(t, inst1, inst3)

	accUpp := inst1.NewAccumulator(attribute.NewSet(), rds[0])
	accUpp.(Updater[int64]).Update(1)

	accUpp.Accumulate()

	output := testCollect(t, vc, rds[0])

	require.Equal(t, 1, len(output))
	require.Equal(t, testCumulative(
		testInst("bar", sdkinstrument.CounterKind, number.Int64Kind, instrument.WithDescription("very long")),
		testPoint(startTime, endTime, int64MonoSum(1))), output[0],
	)
}

func testCollect(t *testing.T, vc *Compiler, r *reader.Config) []reader.Instrument {
	var output []reader.Instrument
	for _, coll := range vc.Collectors(r) {
		coll.Collect(r, testSequence, &output)
	}
	return output
}

func testInstrument(desc sdkinstrument.Descriptor, temporality aggregation.Temporality, points ...reader.Point) reader.Instrument {
	return reader.Instrument{
		Descriptor:  desc,
		Temporality: temporality,
		Points:      points,
	}
}

func testCumulative(desc sdkinstrument.Descriptor, points ...reader.Point) reader.Instrument {
	return testInstrument(desc, aggregation.CumulativeTemporality, points...)
}

func testPoint(start, end time.Time, agg aggregation.Aggregation, kvs ...attribute.KeyValue) reader.Point {
	attrs := attribute.NewSet(kvs...)
	return reader.Point{
		Start:       start,
		End:         end,
		Attributes:  attrs,
		Aggregation: agg,
	}
}

// TestViewDescription ensures that a View can override the description.
func TestViewDescription(t *testing.T) {
	rds := oneTestReader()

	vc := New(testLib, []views.View{
		views.New(
			views.MatchInstrumentName("foo"),
			views.WithDescription("something helpful"),
		),
	}, rds)

	inst1, err1 := vc.Compile(testInst(
		"foo", sdkinstrument.CounterKind, number.Int64Kind,
		instrument.WithDescription("other description"),
	))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	attrs := []attribute.KeyValue{
		attribute.String("K", "V"),
	}
	accUpp := inst1.NewAccumulator(attribute.NewSet(attrs...), rds[0])
	accUpp.(Updater[int64]).Update(1)

	accUpp.Accumulate()

	output := testCollect(t, vc, rds[0])

	require.Equal(t, 1, len(output))
	require.Equal(t,
		testCumulative(
			testInst(
				"foo", sdkinstrument.CounterKind, number.Int64Kind,
				instrument.WithDescription("something helpful"),
			),
			testPoint(startTime, endTime, int64MonoSum(1), attribute.String("K", "V")),
		),
		output[0],
	)
}

// TestKeyFilters verifies that keys are filtred and metrics are
// correctly aggregated.
func TestKeyFilters(t *testing.T) {
	rds := oneTestReader()

	vc := New(testLib, []views.View{
		views.New(views.WithKeys("a", "b")),
	}, rds)

	inst, err := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err)
	require.NotNil(t, inst)

	accUpp1 := inst.NewAccumulator(
		attribute.NewSet(attribute.String("a", "1"), attribute.String("b", "2"), attribute.String("c", "3")),
		rds[0],
	)
	accUpp2 := inst.NewAccumulator(
		attribute.NewSet(attribute.String("a", "1"), attribute.String("b", "2"), attribute.String("d", "4")),
		rds[0],
	)

	accUpp1.(Updater[int64]).Update(1)
	accUpp2.(Updater[int64]).Update(1)
	accUpp1.Accumulate()
	accUpp2.Accumulate()

	output := testCollect(t, vc, rds[0])

	require.Equal(t, 1, len(output))
	require.Equal(t, testCumulative(
		testInst("foo", sdkinstrument.CounterKind, number.Int64Kind),
		testPoint(
			startTime, endTime, int64MonoSum(2),
			attribute.String("a", "1"), attribute.String("b", "2"),
		)), output[0],
	)
}

// TestTwoCounterReaders performs alternating reads from two readers,
// they see 10, 20, 30, 40.
func TestTwoCounterReaders(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	scntr, _ := vc.Compile(testInst("sync_counter", sdkinstrument.CounterKind, number.Int64Kind))

	sup01 := scntr.NewAccumulator(attribute.NewSet(), nil)

	for twice := int64(0); twice < 2; twice++ {

		sup01.(Updater[int64]).Update(10)

		sup01.Accumulate()

		// Collect reader 0 reads 10 or 30
		output := testCollect(t, vc, rds[0])

		require.Equal(t, 1, len(output))
		require.Equal(t,
			testCumulative(
				testInst("sync_counter", sdkinstrument.CounterKind, number.Int64Kind),
				testPoint(startTime, endTime, int64MonoSum(10+twice*20)),
			),
			output[0],
		)

		sup01.(Updater[int64]).Update(10)

		sup01.Accumulate()

		// Collect reader 1 reads 20 or 40
		output = testCollect(t, vc, rds[1])
		require.Equal(t, 1, len(output))
		require.Equal(t,
			testCumulative(
				testInst("sync_counter", sdkinstrument.CounterKind, number.Int64Kind),
				testPoint(startTime, endTime, int64MonoSum(20+twice*20)),
			),
			output[0],
		)
	}
}

// TestTwoCounterObserverReaders performs alternating reads from two readers,
// they see 101, 102, 103, 104.
func TestTwoCounterObserverReaders(t *testing.T) {
	rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	scntr, _ := vc.Compile(testInst("async_counter", sdkinstrument.CounterObserverKind, number.Float64Kind))

	for twice := 0.0; twice < 2; twice++ {

		aup0 := scntr.NewAccumulator(attribute.NewSet(), rds[0])
		aup0.(Updater[float64]).Update(101 + twice*2)
		aup0.Accumulate()

		// Collect reader 0
		output := testCollect(t, vc, rds[0])

		require.Equal(t, 1, len(output))
		require.Equal(t,
			testCumulative(
				testInst("async_counter", sdkinstrument.CounterObserverKind, number.Float64Kind),
				testPoint(startTime, endTime, float64MonoSum(101+twice*2)),
			),
			output[0],
		)

		aup1 := scntr.NewAccumulator(attribute.NewSet(), rds[1])
		aup1.(Updater[float64]).Update(102 + twice*2)
		aup1.Accumulate()

		// Collect reader 1
		output = testCollect(t, vc, rds[1])
		require.Equal(t, 1, len(output))
		require.Equal(t,
			testCumulative(
				testInst("async_counter", sdkinstrument.CounterObserverKind, number.Float64Kind),
				testPoint(startTime, endTime, float64MonoSum(102+twice*2)),
			),
			output[0],
		)
	}
}

func TestSemanticIncompat(t *testing.T) {
	// rds := oneTestReader()

	// vc := New(testLib, []views.View{
	// 	views.New(
	// 		views.MatchInstrumentName("gauge"),
	// 		views.WithAggregation("gauge"),
	// 	),
	// 	views.New(
	// 		views.MatchInstrumentName("sum"),
	// 		views.WithAggregation("sum"),
	// 	),
	// 	views.New(
	// 		views.MatchInstrumentName("hist"),
	// 		views.WithAggregation("histogram"),
	// 	),
	// }, rds)

	// type pair struct {
	// 	inst sdkinstrument.Kind
	// 	agg  aggregation.Kind
	// }

	// cant := []pair{
	// 	// Gauge observers can't become sums or histograms
	// 	{sdkinstrument.GaugeObserverKind, aggregation.SumKind},
	// 	{sdkinstrument.GaugeObserverKind, aggregation.HistogramKind},

	// 	// UpDownCounters can't become histograms or gauges
	// 	{sdkinstrument.UpDownCounterKind, aggregation.HistogramKind},
	// 	{sdkinstrument.UpDownCounterKind, aggregation.GaugeKind},

	// 	// (UpDown)CounterObservers can't become histograms
	// 	{sdkinstrument.UpDownCounterObserverKind, aggregation.HistogramKind},
	// 	{sdkinstrument.CounterObserverKind, aggregation.HistogramKind},
	// }

	// _, err := vc.Compile()
	// require.Equal(t, "GaugeKind instrument incompatible with sum aggregation", err.Error())

	// _, err = vc.Compile()
	// require.Equal(t, "GaugeKind instrument incompatible with histogram aggregation", err.Error())

	// // Counters and histograms can't become gauges
	// _, err = vc.Compile(testInst("gauge", sdkinstrument.CounterKind, number.Int64Kind))
	// require.Equal(t, "CounterKind instrument incompatible with gauge aggregation", err.Error())

	// _, err = vc.Compile(testInst("gauge", sdkinstrument.HistogramKind, number.Int64Kind))
	// require.Equal(t, "HistogramKind instrument incompatible with gauge aggregation", err.Error())

	// sdkinstrument.CounterObserverKind
}
