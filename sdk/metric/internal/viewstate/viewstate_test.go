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
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/internal/test"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

var (
	testLib = instrumentation.Library{
		Name: "test",
	}

	fooToBarView = view.WithClause(
		view.MatchInstrumentName("foo"),
		view.WithName("bar"),
	)

	testHistBoundaries = []float64{1, 2, 3}

	altHistogramConfig = aggregator.Config{
		Histogram: aggregator.HistogramConfig{
			ExplicitBoundaries: testHistBoundaries,
		},
	}

	fooToBarAltHistView = view.WithClause(
		view.MatchInstrumentName("foo"),
		view.WithName("bar"),
		view.WithAggregatorConfig(altHistogramConfig),
	)

	fooToBarFilteredView = view.WithClause(
		view.MatchInstrumentName("foo"),
		view.WithName("bar"),
		view.WithKeys([]attribute.Key{"a", "b"}),
	)

	fooToBarDifferentFiltersViews = []view.Option{
		fooToBarFilteredView,
		view.WithClause(
			view.MatchInstrumentName("bar"),
			view.WithKeys([]attribute.Key{"a"}),
		),
	}

	fooToBarSameFiltersViews = []view.Option{
		fooToBarFilteredView,
		view.WithClause(
			view.MatchInstrumentName("bar"),
			view.WithKeys([]attribute.Key{"a", "b"}),
		),
	}

	dropHistInstView = view.WithClause(
		view.MatchInstrumentKind(sdkinstrument.HistogramKind),
		view.WithAggregation(aggregation.DropKind),
	)

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

	testSequence = data.Sequence{
		Start: startTime,
		Last:  middleTime,
		Now:   endTime,
	}
)

func testCompile(vc *Compiler, name string, ik sdkinstrument.Kind, nk number.Kind, opts ...instrument.Option) (Instrument, error) {
	inst, conflicts := vc.Compile(test.Descriptor(name, ik, nk, opts...))
	return inst, conflicts.AsError()
}

func testCollect(t *testing.T, vc *Compiler) []data.Instrument {
	return test.CollectScope(t, vc.Collectors(), testSequence)
}

func testCollectSequence(t *testing.T, vc *Compiler, seq data.Sequence) []data.Instrument {
	return test.CollectScope(t, vc.Collectors(), seq)
}

func testCollectSequenceReuse(t *testing.T, vc *Compiler, seq data.Sequence, output *data.Scope) []data.Instrument {
	return test.CollectScopeReuse(t, vc.Collectors(), seq, output)
}

// TestDeduplicateNoConflict verifies that two identical instruments
// have the same collector.
func TestDeduplicateNoConflict(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.Equal(t, inst1, inst2)
}

// TestDeduplicateRenameNoConflict verifies that one instrument can be renamed
// such that it becomes identical to another, so no conflict.
func TestDeduplicateRenameNoConflict(t *testing.T) {
	vc := New(testLib, view.New("test", fooToBarView))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "bar", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.Equal(t, inst1, inst2)
}

// TestNoRenameNoConflict verifies that one instrument does not
// conflict with another differently-named instrument.
func TestNoRenameNoConflict(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "bar", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateNumberConflict verifies that two same instruments
// except different number kind conflict.
func TestDuplicateNumberConflict(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Float64Kind)
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflictsError{}))
	require.Equal(t, 1, len(err2.(ViewConflictsError)))
	require.Equal(t, 1, len(err2.(ViewConflictsError)["test"]))
	require.Equal(t, 2, len(err2.(ViewConflictsError)["test"][0].Duplicates))

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateSyncAsyncConflict verifies that two same instruments
// except one synchonous, one asynchronous conflict.
func TestDuplicateSyncAsyncConflict(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "foo", sdkinstrument.CounterObserverKind, number.Float64Kind)
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflictsError{}))

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateUnitConflict verifies that two same instruments
// except different units conflict.
func TestDuplicateUnitConflict(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Float64Kind, instrument.WithUnit("gal_us"))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Float64Kind, instrument.WithUnit("cft_i"))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflictsError{}))
	require.Contains(t, err2.Error(), "test: name \"foo\" conflicts Counter-Float64-MonotonicSum-gal_us")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateMonotonicConflict verifies that two same instruments
// except different monotonic values.
func TestDuplicateMonotonicConflict(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "foo", sdkinstrument.UpDownCounterKind, number.Float64Kind)
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflictsError{}))
	require.Contains(t, err2.Error(), "UpDownCounter-Float64-NonMonotonicSum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregatorConfigConflict verifies that two same instruments
// except different aggregator.Config values.
func TestDuplicateAggregatorConfigConflict(t *testing.T) {
	vc := New(testLib, view.New("test", fooToBarAltHistView))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.HistogramKind, number.Float64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "bar", sdkinstrument.HistogramKind, number.Float64Kind)
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflictsError{}))
	require.Contains(t, err2.Error(), "different aggregator configuration")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregatorConfigNoConflict verifies that two same instruments
// with same aggregator.Config values configured in different ways.
func TestDuplicateAggregatorConfigNoConflict(t *testing.T) {
	for _, nk := range numberKinds {
		t.Run(nk.String(), func(t *testing.T) {
			views := view.New(
				"test",
				view.WithDefaultAggregationConfigSelector(
					func(_ sdkinstrument.Kind) (int64Config, float64Config aggregator.Config) {
						if nk == number.Int64Kind {
							return altHistogramConfig, aggregator.Config{}
						}
						return aggregator.Config{}, altHistogramConfig
					},
				),
				fooToBarAltHistView,
			)

			vc := New(testLib, views)

			inst1, err1 := testCompile(vc, "foo", sdkinstrument.HistogramKind, nk)
			require.NoError(t, err1)
			require.NotNil(t, inst1)

			inst2, err2 := testCompile(vc, "bar", sdkinstrument.HistogramKind, nk)
			require.NoError(t, err2)
			require.NotNil(t, inst2)

			require.Equal(t, inst1, inst2)
		})
	}
}

// TestDuplicateAggregationKindConflict verifies that two instruments
// with different aggregation kinds conflict.
func TestDuplicateAggregationKindConflict(t *testing.T) {
	vc := New(testLib, view.New("test", fooToBarView))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.HistogramKind, number.Int64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "bar", sdkinstrument.CounterKind, number.Int64Kind)
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, ViewConflictsError{}))
	require.Contains(t, err2.Error(), "name \"bar\" (original \"foo\") conflicts Histogram-Int64-Histogram, Counter-Int64-MonotonicSum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregationKindNoConflict verifies that two
// instruments with different aggregation kinds do not conflict when
// the view drops one of the instruments.
func TestDuplicateAggregationKindNoConflict(t *testing.T) {
	vc := New(testLib, view.New("test", dropHistInstView))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.HistogramKind, number.Int64Kind)
	require.NoError(t, err1)
	require.Nil(t, inst1) // The viewstate.Instrument is nil, instruments become no-ops.

	inst2, err2 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err2)
	require.NotNil(t, inst2)
}

// TestDuplicateMultipleConflicts verifies that multiple duplicate
// instrument conflicts include sufficient explanatory information.
func TestDuplicateMultipleConflicts(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err1 := testCompile(vc, "foo", instrumentKinds[0], number.Float64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	for num, ik := range instrumentKinds[1:] {
		inst2, err2 := testCompile(vc, "foo", ik, number.Float64Kind)
		require.Error(t, err2)
		require.NotNil(t, inst2)
		require.True(t, errors.Is(err2, ViewConflictsError{}))
		// The total number of conflicting definitions is 1 in
		// the first place and num+1 for the iterations of this loop.
		require.Equal(t, num+2, len(err2.(ViewConflictsError)["test"][0].Duplicates))

		if num > 0 {
			require.Contains(t, err2.Error(), fmt.Sprintf("and %d more", num))
		}
	}
}

// TestDuplicateFilterConflicts verifies several cases where
// instruments output the same metric w/ different filters create conflicts.
func TestDuplicateFilterConflicts(t *testing.T) {
	for idx, vws := range [][]view.Option{
		// In the first case, foo has two attribute filters bar has 0.
		[]view.Option{fooToBarFilteredView},
		// In the second case, foo has two attribute filters bar has 1.
		fooToBarDifferentFiltersViews,
	} {
		t.Run(fmt.Sprint(idx), func(t *testing.T) {
			vc := New(testLib, view.New("test", vws...))

			inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
			require.NoError(t, err1)
			require.NotNil(t, inst1)

			inst2, err2 := testCompile(vc, "bar", sdkinstrument.CounterKind, number.Int64Kind)
			require.Error(t, err2)
			require.NotNil(t, inst2)

			require.True(t, errors.Is(err2, ViewConflictsError{}))
			require.Contains(t, err2.Error(), "name \"bar\" (original \"foo\") has conflicts: different attribute filters")
		})
	}
}

// TestDeduplicateSameFilters thests that when one instrument is
// renamed to match another exactly, including filters, they are not
// in conflict.
func TestDeduplicateSameFilters(t *testing.T) {
	vc := New(testLib, view.New("test", fooToBarSameFiltersViews...))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := testCompile(vc, "bar", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.Equal(t, inst1, inst2)
}

// TestDuplicatesMergeDescriptor ensures that the longest description string is used.
func TestDuplicatesMergeDescriptor(t *testing.T) {
	vc := New(testLib, view.New("test", fooToBarSameFiltersViews...))

	inst1, err1 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	// This is the winning description:
	inst2, err2 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind, instrument.WithDescription("very long"))
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	inst3, err3 := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind, instrument.WithDescription("shorter"))
	require.NoError(t, err3)
	require.NotNil(t, inst3)

	require.Equal(t, inst1, inst2)
	require.Equal(t, inst1, inst3)

	accUpp := inst1.NewAccumulator(attribute.NewSet())
	accUpp.(Updater[int64]).Update(1)

	accUpp.SnapshotAndProcess()

	output := testCollect(t, vc)

	require.Equal(t, 1, len(output))
	require.Equal(t, test.Instrument(
		test.Descriptor("bar", sdkinstrument.CounterKind, number.Int64Kind, instrument.WithDescription("very long")),
		test.Point(startTime, endTime, sum.NewMonotonicInt64(1))), output[0],
	)
}

// TestViewDescription ensures that a View can override the description.
func TestViewDescription(t *testing.T) {
	views := view.New(
		"test",
		view.WithClause(
			view.MatchInstrumentName("foo"),
			view.WithDescription("something helpful"),
		),
	)

	vc := New(testLib, views)

	inst1, err1 := testCompile(vc,
		"foo", sdkinstrument.CounterKind, number.Int64Kind,
		instrument.WithDescription("other description"),
	)
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	attrs := []attribute.KeyValue{
		attribute.String("K", "V"),
	}
	accUpp := inst1.NewAccumulator(attribute.NewSet(attrs...))
	accUpp.(Updater[int64]).Update(1)

	accUpp.SnapshotAndProcess()

	output := testCollect(t, vc)

	require.Equal(t, 1, len(output))
	require.Equal(t,
		test.Instrument(
			test.Descriptor(
				"foo", sdkinstrument.CounterKind, number.Int64Kind,
				instrument.WithDescription("something helpful"),
			),
			test.Point(startTime, endTime, sum.NewMonotonicInt64(1), attribute.String("K", "V")),
		),
		output[0],
	)
}

// TestKeyFilters verifies that keys are filtred and metrics are
// correctly aggregated.
func TestKeyFilters(t *testing.T) {
	views := view.New("test",
		view.WithClause(view.WithKeys([]attribute.Key{"a", "b"})),
	)

	vc := New(testLib, views)

	inst, err := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err)
	require.NotNil(t, inst)

	accUpp1 := inst.NewAccumulator(
		attribute.NewSet(attribute.String("a", "1"), attribute.String("b", "2"), attribute.String("c", "3")),
	)
	accUpp2 := inst.NewAccumulator(
		attribute.NewSet(attribute.String("a", "1"), attribute.String("b", "2"), attribute.String("d", "4")),
	)

	accUpp1.(Updater[int64]).Update(1)
	accUpp2.(Updater[int64]).Update(1)
	accUpp1.SnapshotAndProcess()
	accUpp2.SnapshotAndProcess()

	output := testCollect(t, vc)

	require.Equal(t, 1, len(output))
	require.Equal(t, test.Instrument(
		test.Descriptor("foo", sdkinstrument.CounterKind, number.Int64Kind),
		test.Point(
			startTime, endTime, sum.NewMonotonicInt64(2),
			attribute.String("a", "1"), attribute.String("b", "2"),
		)), output[0],
	)
}

// TestTwoViewsOneInt64Instrument verifies that multiple int64
// instrument behaviors work; in this case, viewing a Sum in each
// of three independent dimensions.
func TestTwoViewsOneInt64Instrument(t *testing.T) {
	views := view.New(
		"test",
		view.WithClause(
			view.MatchInstrumentName("foo"),
			view.WithName("foo_a"),
			view.WithKeys([]attribute.Key{"a"}),
		),
		view.WithClause(
			view.MatchInstrumentName("foo"),
			view.WithName("foo_b"),
			view.WithKeys([]attribute.Key{"b"}),
		),
		view.WithClause(
			view.MatchInstrumentName("foo"),
			view.WithName("foo_c"),
			view.WithKeys([]attribute.Key{"c"}),
		),
	)

	vc := New(testLib, views)

	inst, err := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err)

	for _, acc := range []Accumulator{
		inst.NewAccumulator(attribute.NewSet(attribute.String("a", "1"), attribute.String("b", "1"))),
		inst.NewAccumulator(attribute.NewSet(attribute.String("a", "1"), attribute.String("b", "2"))),
		inst.NewAccumulator(attribute.NewSet(attribute.String("a", "2"), attribute.String("b", "1"))),
		inst.NewAccumulator(attribute.NewSet(attribute.String("a", "2"), attribute.String("b", "2"))),
	} {
		acc.(Updater[int64]).Update(1)
		acc.SnapshotAndProcess()
	}

	output := testCollect(t, vc)

	test.RequireEqualMetrics(t,
		output,
		test.Instrument(
			test.Descriptor("foo_a", sdkinstrument.CounterKind, number.Int64Kind),
			test.Point(
				startTime, endTime, sum.NewMonotonicInt64(2), attribute.String("a", "1"),
			),
			test.Point(
				startTime, endTime, sum.NewMonotonicInt64(2), attribute.String("a", "2"),
			),
		),
		test.Instrument(
			test.Descriptor("foo_b", sdkinstrument.CounterKind, number.Int64Kind),
			test.Point(
				startTime, endTime, sum.NewMonotonicInt64(2), attribute.String("b", "1"),
			),
			test.Point(
				startTime, endTime, sum.NewMonotonicInt64(2), attribute.String("b", "2"),
			),
		),
		test.Instrument(
			test.Descriptor("foo_c", sdkinstrument.CounterKind, number.Int64Kind),
			test.Point(
				startTime, endTime, sum.NewMonotonicInt64(4),
			),
		),
	)
}

// TestHistogramTwoAggregations verifies that two float64 instrument
// behaviors are correctly combined, in this case one sum and one histogram.
func TestHistogramTwoAggregations(t *testing.T) {
	views := view.New(
		"test",
		view.WithClause(
			view.MatchInstrumentName("foo"),
			view.WithName("foo_sum"),
			view.WithAggregation(aggregation.MonotonicSumKind),
			view.WithKeys([]attribute.Key{}),
		),
		view.WithClause(
			view.MatchInstrumentName("foo"),
			view.WithName("foo_hist"),
			view.WithAggregation(aggregation.HistogramKind),
		),
	)

	vc := New(testLib, views)

	inst, err := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err)

	acc := inst.NewAccumulator(attribute.NewSet())
	acc.(Updater[float64]).Update(1)
	acc.(Updater[float64]).Update(2)
	acc.(Updater[float64]).Update(3)
	acc.(Updater[float64]).Update(4)
	acc.SnapshotAndProcess()

	output := testCollect(t, vc)

	test.RequireEqualMetrics(t, output,
		test.Instrument(
			test.Descriptor("foo_sum", sdkinstrument.CounterKind, number.Float64Kind),
			test.Point(
				startTime, endTime, sum.NewMonotonicFloat64(10),
			),
		),
		test.Instrument(
			test.Descriptor("foo_hist", sdkinstrument.CounterKind, number.Float64Kind),
			test.Point(
				startTime, endTime, histogram.NewFloat64(nil, 1, 2, 3, 4),
			),
		),
	)
}

// TestAllKeysFilter tests that view.WithKeys([]attribute.Key{})
// correctly erases all keys.
func TestAllKeysFilter(t *testing.T) {
	views := view.New(
		"test",
		view.WithClause(view.WithKeys([]attribute.Key{})),
	)

	vc := New(testLib, views)

	inst, err := testCompile(vc, "foo", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err)

	acc1 := inst.NewAccumulator(attribute.NewSet(attribute.String("a", "1")))
	acc1.(Updater[float64]).Update(1)
	acc1.SnapshotAndProcess()

	acc2 := inst.NewAccumulator(attribute.NewSet(attribute.String("b", "2")))
	acc2.(Updater[float64]).Update(1)
	acc2.SnapshotAndProcess()

	output := testCollect(t, vc)

	test.RequireEqualMetrics(t, output,
		test.Instrument(
			test.Descriptor("foo", sdkinstrument.CounterKind, number.Float64Kind),
			test.Point(
				startTime, endTime, sum.NewMonotonicFloat64(2),
			),
		),
	)
}

// TestAnySumAggregation checks that the proper aggregation inference
// is performed for each of the inbstrument types when
// aggregation.AnySum kind is configured.
func TestAnySumAggregation(t *testing.T) {
	views := view.New(
		"test",
		view.WithClause(view.WithAggregation(aggregation.AnySumKind)),
	)

	vc := New(testLib, views)

	for _, ik := range []sdkinstrument.Kind{
		sdkinstrument.CounterKind,
		sdkinstrument.CounterObserverKind,
		sdkinstrument.UpDownCounterKind,
		sdkinstrument.UpDownCounterObserverKind,
		sdkinstrument.HistogramKind,
		sdkinstrument.GaugeObserverKind,
	} {
		inst, err := testCompile(vc, ik.String(), ik, number.Float64Kind)
		if ik == sdkinstrument.GaugeObserverKind {
			// semantic conflict, Gauge can't handle AnySum aggregation!
			require.Error(t, err)
			require.Contains(t,
				err.Error(),
				"GaugeObserver instrument incompatible with Undefined aggregation",
			)
		} else {
			require.NoError(t, err)
		}

		acc := inst.NewAccumulator(attribute.NewSet())
		acc.(Updater[float64]).Update(1)
		acc.SnapshotAndProcess()
	}

	output := testCollect(t, vc)

	test.RequireEqualMetrics(t, output,
		test.Instrument(
			test.Descriptor("CounterKind", sdkinstrument.CounterKind, number.Float64Kind),
			test.Point(startTime, endTime, sum.NewMonotonicFloat64(1)), // AnySum -> Monotonic
		),
		test.Instrument(
			test.Descriptor("CounterObserverKind", sdkinstrument.CounterObserverKind, number.Float64Kind),
			test.Point(startTime, endTime, sum.NewMonotonicFloat64(1)), // AnySum -> Monotonic
		),
		test.Instrument(
			test.Descriptor("UpDownCounterKind", sdkinstrument.UpDownCounterKind, number.Float64Kind),
			test.Point(startTime, endTime, sum.NewNonMonotonicFloat64(1)), // AnySum -> Non-Monotonic
		),
		test.Instrument(
			test.Descriptor("UpDownCounterObserverKind", sdkinstrument.UpDownCounterObserverKind, number.Float64Kind),
			test.Point(startTime, endTime, sum.NewNonMonotonicFloat64(1)), // AnySum -> Non-Monotonic
		),
		test.Instrument(
			test.Descriptor("HistogramKind", sdkinstrument.HistogramKind, number.Float64Kind),
			test.Point(startTime, endTime, sum.NewMonotonicFloat64(1)), // Histogram to Monotonic Sum
		),
		test.Instrument(
			test.Descriptor("GaugeObserverKind", sdkinstrument.GaugeObserverKind, number.Float64Kind),
			test.Point(startTime, endTime, gauge.NewFloat64(1)), // This stays a Gauge!
		),
	)
}

// TestDuplicateAsyncMeasurementsIngored tests that asynchronous
// instrument accumulators keep only the last observed value, while
// synchronous instruments correctly snapshotAndProcess them all.
func TestDuplicateAsyncMeasurementsIngored(t *testing.T) {
	vc := New(testLib, view.New("test"))

	inst1, err := testCompile(vc, "async", sdkinstrument.CounterObserverKind, number.Float64Kind)
	require.NoError(t, err)

	inst2, err := testCompile(vc, "sync", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err)

	for _, inst := range []Instrument{inst1, inst2} {
		acc := inst.NewAccumulator(attribute.NewSet())
		acc.(Updater[float64]).Update(1)
		acc.(Updater[float64]).Update(10)
		acc.(Updater[float64]).Update(100)
		acc.(Updater[float64]).Update(1000)
		acc.(Updater[float64]).Update(10000)
		acc.(Updater[float64]).Update(100000)
		acc.SnapshotAndProcess()
	}

	output := testCollect(t, vc)

	test.RequireEqualMetrics(t, output,
		test.Instrument(
			test.Descriptor("async", sdkinstrument.CounterObserverKind, number.Float64Kind),
			test.Point(
				startTime, endTime, sum.NewMonotonicFloat64(100000),
			),
		),
		test.Instrument(
			test.Descriptor("sync", sdkinstrument.CounterKind, number.Float64Kind),
			test.Point(
				startTime, endTime, sum.NewMonotonicFloat64(111111),
			),
		),
	)
}

// TestCumulativeTemporality ensures that synchronous instruments
// snapshotAndProcess data over time, whereas asynchronous instruments do not.
func TestCumulativeTemporality(t *testing.T) {
	views := view.New(
		"test",
		view.WithClause(
			// Dropping all keys
			view.WithKeys([]attribute.Key{}),
		),
		view.WithDefaultAggregationTemporalitySelector(view.StandardTemporality),
	)

	vc := New(testLib, views)

	inst1, err := testCompile(vc, "sync", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err)

	inst2, err := testCompile(vc, "async", sdkinstrument.CounterObserverKind, number.Float64Kind)
	require.NoError(t, err)

	setA := attribute.NewSet(attribute.String("A", "1"))
	setB := attribute.NewSet(attribute.String("B", "1"))

	for rounds := 1; rounds <= 2; rounds++ {
		for _, acc := range []Accumulator{
			inst1.NewAccumulator(setA),
			inst1.NewAccumulator(setB),
			inst2.NewAccumulator(setA),
			inst2.NewAccumulator(setB),
		} {
			acc.(Updater[float64]).Update(1)
			acc.SnapshotAndProcess()
		}

		test.RequireEqualMetrics(t, testCollect(t, vc),
			test.Instrument(
				test.Descriptor("sync", sdkinstrument.CounterKind, number.Float64Kind),
				test.Point(
					// Because synchronous instruments snapshotAndProcess, the
					// rounds multiplier is used here but not in the case below.
					startTime, endTime, sum.NewMonotonicFloat64(float64(rounds)*2),
				),
			),
			test.Instrument(
				test.Descriptor("async", sdkinstrument.CounterObserverKind, number.Float64Kind),
				test.Point(
					startTime, endTime, sum.NewMonotonicFloat64(2),
				),
			),
		)
	}
}

// TestDeltaTemporality ensures that synchronous instruments
// snapshotAndProcess data over time, whereas asynchronous instruments do not.
func TestDeltaTemporalityCounter(t *testing.T) {
	views := view.New(
		"test",
		view.WithClause(
			// Dropping all keys
			view.WithKeys([]attribute.Key{}),
		),
		view.WithDefaultAggregationTemporalitySelector(view.DeltaPreferredTemporality),
	)

	vc := New(testLib, views)

	inst1, err := testCompile(vc, "sync", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err)

	inst2, err := testCompile(vc, "async", sdkinstrument.CounterObserverKind, number.Float64Kind)
	require.NoError(t, err)

	setA := attribute.NewSet(attribute.String("A", "1"))
	setB := attribute.NewSet(attribute.String("B", "1"))

	seq := testSequence

	for rounds := 1; rounds <= 3; rounds++ {
		for _, acc := range []Accumulator{
			inst1.NewAccumulator(setA),
			inst1.NewAccumulator(setB),
			inst2.NewAccumulator(setA),
			inst2.NewAccumulator(setB),
		} {
			acc.(Updater[float64]).Update(float64(rounds))
			acc.SnapshotAndProcess()
		}

		test.RequireEqualMetrics(t, testCollectSequence(t, vc, seq),
			test.Instrument(
				test.Descriptor("sync", sdkinstrument.CounterKind, number.Float64Kind),
				test.Point(
					// By construction, the change is rounds per attribute set == 2*rounds
					seq.Last, seq.Now, sum.NewMonotonicFloat64(2*float64(rounds)),
				),
			),
			test.Instrument(
				test.Descriptor("async", sdkinstrument.CounterObserverKind, number.Float64Kind),
				test.Point(
					// By construction, the change is 1 per attribute set == 2
					seq.Last, seq.Now, sum.NewMonotonicFloat64(2),
				),
			),
		)

		// Update the test sequence
		seq.Last = seq.Now
		seq.Now = time.Now()
	}
}

// TestDeltaTemporalityGauge ensures that the asynchronous gauge
// when used with delta temporalty only reports changed values.
func TestDeltaTemporalityGauge(t *testing.T) {
	views := view.New(
		"test",
		view.WithDefaultAggregationTemporalitySelector(view.DeltaPreferredTemporality),
	)

	vc := New(testLib, views)

	instF, err := testCompile(vc, "gaugeF", sdkinstrument.GaugeObserverKind, number.Float64Kind)
	require.NoError(t, err)

	instI, err := testCompile(vc, "gaugeI", sdkinstrument.GaugeObserverKind, number.Int64Kind)
	require.NoError(t, err)

	set := attribute.NewSet()

	observe := func(x int) {
		accI := instI.NewAccumulator(set)
		accI.(Updater[int64]).Update(int64(x))
		accI.SnapshotAndProcess()

		accF := instF.NewAccumulator(set)
		accF.(Updater[float64]).Update(float64(x))
		accF.SnapshotAndProcess()
	}

	expectValues := func(x int, seq data.Sequence) {
		test.RequireEqualMetrics(t,
			testCollectSequence(t, vc, seq),
			test.Instrument(
				test.Descriptor("gaugeF", sdkinstrument.GaugeObserverKind, number.Float64Kind),
				test.Point(seq.Last, seq.Now, gauge.NewFloat64(float64(x))),
			),
			test.Instrument(
				test.Descriptor("gaugeI", sdkinstrument.GaugeObserverKind, number.Int64Kind),
				test.Point(seq.Last, seq.Now, gauge.NewInt64(int64(x))),
			),
		)
	}
	expectNone := func(seq data.Sequence) {
		test.RequireEqualMetrics(t,
			testCollectSequence(t, vc, seq),
			test.Instrument(
				test.Descriptor("gaugeF", sdkinstrument.GaugeObserverKind, number.Float64Kind),
			),
			test.Instrument(
				test.Descriptor("gaugeI", sdkinstrument.GaugeObserverKind, number.Int64Kind),
			),
		)
	}
	seq := testSequence
	tick := func() {
		// Update the test sequence
		seq.Last = seq.Now
		seq.Now = time.Now()
	}

	observe(10)
	expectValues(10, seq)
	tick()

	observe(10)
	expectNone(seq)
	tick()

	observe(10)
	expectNone(seq)
	tick()

	observe(11)
	expectValues(11, seq)
	tick()

	observe(11)
	expectNone(seq)
	tick()

	observe(10)
	expectValues(10, seq)
	tick()
}

// TestSyncDeltaTemporalityCounter ensures that counter and updowncounter
// are skip points with delta temporality and no change.
func TestSyncDeltaTemporalityCounter(t *testing.T) {
	views := view.New(
		"test",
		view.WithDefaultAggregationTemporalitySelector(
			func(ik sdkinstrument.Kind) aggregation.Temporality {
				return aggregation.DeltaTemporality // Always delta
			}),
	)

	vc := New(testLib, views)

	instCF, err := testCompile(vc, "counterF", sdkinstrument.CounterKind, number.Float64Kind)
	require.NoError(t, err)

	instCI, err := testCompile(vc, "counterI", sdkinstrument.CounterKind, number.Int64Kind)
	require.NoError(t, err)

	instUF, err := testCompile(vc, "updowncounterF", sdkinstrument.UpDownCounterKind, number.Float64Kind)
	require.NoError(t, err)

	instUI, err := testCompile(vc, "updowncounterI", sdkinstrument.UpDownCounterKind, number.Int64Kind)
	require.NoError(t, err)

	set := attribute.NewSet()

	var output data.Scope

	observe := func(mono, nonMono int) {
		accCI := instCI.NewAccumulator(set)
		accCI.(Updater[int64]).Update(int64(mono))
		accCI.SnapshotAndProcess()

		accCF := instCF.NewAccumulator(set)
		accCF.(Updater[float64]).Update(float64(mono))
		accCF.SnapshotAndProcess()

		accUI := instUI.NewAccumulator(set)
		accUI.(Updater[int64]).Update(int64(nonMono))
		accUI.SnapshotAndProcess()

		accUF := instUF.NewAccumulator(set)
		accUF.(Updater[float64]).Update(float64(nonMono))
		accUF.SnapshotAndProcess()
	}

	expectValues := func(mono, nonMono int, seq data.Sequence) {
		test.RequireEqualMetrics(t,
			testCollectSequenceReuse(t, vc, seq, &output),
			test.Instrument(
				test.Descriptor("counterF", sdkinstrument.CounterKind, number.Float64Kind),
				test.Point(seq.Last, seq.Now, sum.NewMonotonicFloat64(float64(mono))),
			),
			test.Instrument(
				test.Descriptor("counterI", sdkinstrument.CounterKind, number.Int64Kind),
				test.Point(seq.Last, seq.Now, sum.NewMonotonicInt64(int64(mono))),
			),
			test.Instrument(
				test.Descriptor("updowncounterF", sdkinstrument.UpDownCounterKind, number.Float64Kind),
				test.Point(seq.Last, seq.Now, sum.NewNonMonotonicFloat64(float64(nonMono))),
			),
			test.Instrument(
				test.Descriptor("updowncounterI", sdkinstrument.UpDownCounterKind, number.Int64Kind),
				test.Point(seq.Last, seq.Now, sum.NewNonMonotonicInt64(int64(nonMono))),
			),
		)
	}
	expectNone := func(seq data.Sequence) {
		test.RequireEqualMetrics(t,
			testCollectSequenceReuse(t, vc, seq, &output),
			test.Instrument(
				test.Descriptor("counterF", sdkinstrument.CounterKind, number.Float64Kind),
			),
			test.Instrument(
				test.Descriptor("counterI", sdkinstrument.CounterKind, number.Int64Kind),
			),
			test.Instrument(
				test.Descriptor("updowncounterF", sdkinstrument.UpDownCounterKind, number.Float64Kind),
			),
			test.Instrument(
				test.Descriptor("updowncounterI", sdkinstrument.UpDownCounterKind, number.Int64Kind),
			),
		)
	}
	seq := testSequence
	tick := func() {
		// Update the test sequence
		seq.Last = seq.Now
		seq.Now = time.Now()
	}

	observe(10, 10)
	expectValues(10, 10, seq)
	tick()

	observe(0, 100)
	observe(0, -100)
	expectNone(seq)
	tick()

	observe(100, 100)
	expectValues(100, 100, seq)
	tick()
}
