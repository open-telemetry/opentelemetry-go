package viewstate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
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
