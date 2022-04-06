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
)

// testInst returns a test instrument descriptor similar to what Meter creates.
func testInst(name string, ik sdkinstrument.Kind, nk number.Kind, opts ...instrument.Option) sdkinstrument.Descriptor {
	cfg := instrument.NewConfig(opts...)
	return sdkinstrument.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())
}

func twoTestReaders() (one, two *metrictest.Exporter, _ []*reader.Reader) {
	exp1 := metrictest.NewExporter()
	exp2 := metrictest.NewExporter()
	rds := []*reader.Reader{
		reader.New(exp1),
		reader.New(exp2),
	}
	return exp1, exp2, rds
}

func oneTestReader() (*metrictest.Exporter, []*reader.Reader) {
	exp := metrictest.NewExporter()
	rds := []*reader.Reader{reader.New(exp)}
	return exp, rds
}

// TestDeduplicateNoConflict verifies that two identical instruments
// have the same collector.
func TestDeduplicateNoConflict(t *testing.T) {
	_, _, rds := twoTestReaders()

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
	_, _, rds := twoTestReaders()

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
	_, _, rds := twoTestReaders()

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
	_, _, rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, DuplicateConflicts{}))
	require.Equal(t, 2, len(err2.(DuplicateConflicts)))

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateSyncAsyncConflict verifies that two same instruments
// except one synchonous, one asynchronous conflict.
func TestDuplicateSyncAsyncConflict(t *testing.T) {
	_, _, rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterObserverKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, DuplicateConflicts{}))

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateUnitConflict verifies that two same instruments
// except different units conflict.
func TestDuplicateUnitConflict(t *testing.T) {
	_, _, rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind, instrument.WithUnit("gal_us")))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind, instrument.WithUnit("cft_i")))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, DuplicateConflicts{}))
	require.Contains(t, err2.Error(), "2 conflict(s) in 2 reader(s)")
	require.Contains(t, err2.Error(), "conflicts Counter-Float64-Sum-gal_us")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateMonotonicConflict verifies that two same instruments
// except different monotonic values.
func TestDuplicateMonotonicConflict(t *testing.T) {
	_, _, rds := twoTestReaders()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.UpDownCounterKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, DuplicateConflicts{}))
	require.Contains(t, err2.Error(), "2 conflict(s) in 2 reader(s)")
	require.Contains(t, err2.Error(), "UpDownCounter-Float64-Sum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregatorConfigConflict verifies that two same instruments
// except different aggregator.Config values.
func TestDuplicateAggregatorConfigConflict(t *testing.T) {
	_, _, rds := twoTestReaders()

	vc := New(testLib, fooToBarAltHistView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.HistogramKind, number.Float64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, DuplicateConflicts{}))
	require.Contains(t, err2.Error(), "different aggregator configuration")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregatorConfigNoConflict verifies that two same instruments
// with same aggregator.Config values configured in different ways.
func TestDuplicateAggregatorConfigNoConflict(t *testing.T) {
	exp := metrictest.NewExporter()

	for _, nk := range numberKinds {
		t.Run(nk.String(), func(t *testing.T) {
			rds := []*reader.Reader{
				reader.New(exp, reader.WithDefaultAggregationConfigFunc(
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
	_, _, rds := twoTestReaders()

	vc := New(testLib, fooToBarView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, DuplicateConflicts{}))
	require.Contains(t, err2.Error(), "2 conflict(s) in 2 reader(s)")
	require.Contains(t, err2.Error(), "name \"bar\" (original \"foo\") conflicts Histogram-Int64-Histogram, Counter-Int64-Sum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregationKindOneConflict verifies that two
// instruments with different aggregation kinds do not conflict when
// the reader drops the instrument.
func TestDuplicateAggregationKindOneConflict(t *testing.T) {
	exp1, exp2, _ := twoTestReaders()
	// Let one reader drop histograms
	rds := []*reader.Reader{
		reader.New(exp1, reader.WithDefaultAggregationKindFunc(func(k sdkinstrument.Kind) aggregation.Kind {
			if k == sdkinstrument.HistogramKind {
				return aggregation.DropKind
			}
			return reader.StandardAggregationKind(k)
		})),
		reader.New(exp2),
	}

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.HistogramKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.True(t, errors.Is(err2, DuplicateConflicts{}))
	require.Contains(t, err2.Error(), "1 conflict(s), e.g.")
	require.Contains(t, err2.Error(), "name \"foo\" conflicts Histogram-Int64-Histogram, Counter-Int64-Sum")

	require.NotEqual(t, inst1, inst2)
}

// TestDuplicateAggregationKindNoConflict verifies that two
// instruments with different aggregation kinds do not conflict when
// the view drops one of the instruments.
func TestDuplicateAggregationKindNoConflict(t *testing.T) {
	_, _, rds := twoTestReaders()

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
	_, rds := oneTestReader()

	vc := New(testLib, nil, rds)

	inst1, err1 := vc.Compile(testInst("foo", instrumentKinds[0], number.Float64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	for num, ik := range instrumentKinds[1:] {
		inst2, err2 := vc.Compile(testInst("foo", ik, number.Float64Kind))
		require.Error(t, err2)
		require.NotNil(t, inst2)
		require.True(t, errors.Is(err2, DuplicateConflicts{}))
		// The total number of conflicting definitions is 1 in
		// the first place and num+1 for the iterations of this loop.
		require.Equal(t, num+2, len(err2.(DuplicateConflicts)[rds[0]][0]))

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
			_, _, rds := twoTestReaders()

			vc := New(testLib, vws, rds)

			inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
			require.NoError(t, err1)
			require.NotNil(t, inst1)

			inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
			require.Error(t, err2)
			require.NotNil(t, inst2)

			require.True(t, errors.Is(err2, DuplicateConflicts{}))
			require.Contains(t, err2.Error(), "2 conflict(s) in 2 reader(s), e.g.")
			require.Contains(t, err2.Error(), "name \"bar\" (original \"foo\") has conflicts: different attribute filters")
		})
	}
}

// TestDeduplicateSameFilters thests that when one instrument is
// renamed to match another exactly, including filters, they are not
// in conflict.
func TestDeduplicateSameFilters(t *testing.T) {
	_, _, rds := twoTestReaders()

	vc := New(testLib, fooToBarSameFiltersView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

	inst2, err2 := vc.Compile(testInst("bar", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err2)
	require.NotNil(t, inst2)

	require.Equal(t, inst1, inst2)
}

func int64Sum(x int64) aggregation.Sum {
	var s sum.State[int64, traits.Int64]
	var methods sum.Methods[int64, traits.Int64, sum.State[int64, traits.Int64]]
	methods.Init(&s, aggregator.Config{})
	methods.Update(&s, 1)
	return &s
}

// TestDuplicatesMergeDescriptor ensures that the longest description string is used.
func TestDuplicatesMergeDescriptor(t *testing.T) {
	_, rds := oneTestReader()

	vc := New(testLib, fooToBarSameFiltersView, rds)

	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.NoError(t, err1)
	require.NotNil(t, inst1)

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

	end := time.Now()
	middle := end.Add(-time.Millisecond)
	start := end.Add(-2 * time.Millisecond)
	var output []reader.Instrument
	inst1.Collect(rds[0], reader.Sequence{
		Start: start,
		Last:  middle,
		Now:   end,
	}, &output)

	require.Equal(t, 1, len(output))
	require.Equal(t,
		reader.Instrument{
			Descriptor: sdkinstrument.Descriptor{
				Name:        "bar",
				Kind:        sdkinstrument.CounterKind,
				NumberKind:  number.Int64Kind,
				Description: "very long", // Note!
			},
			Temporality: aggregation.CumulativeTemporality,
			Series: []reader.Series{
				{
					Start:       start,
					End:         end,
					Attributes:  attribute.NewSet(),
					Aggregation: int64Sum(1),
				},
			},
		},
		output[0])
}
