package viewstate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

var oneConflict = Conflict{
	Semantic: SemanticError{
		InstrumentKind:  sdkinstrument.CounterKind,
		AggregationKind: aggregation.GaugeKind,
	},
}

// TestViewConflictsError exercises the code paths that construct example
// error messages from duplicate instrument conditions.
func TestViewConflictsError(t *testing.T) {
	// Note: These all use "no conflicts" strings, which happens
	// under artificial conditions such as conflicts w/ < 2 examples
	// and allows testing the code that avoids lengthy messages
	// when there is only one conflict or only one reader.
	var err error
	err = ViewConflicts{}
	require.Equal(t, noConflictsString, err.Error())
	require.True(t, errors.Is(err, ViewConflicts{}))

	rd1 := reader.New(metrictest.NewExporter())
	rd2 := reader.New(metrictest.NewExporter())

	require.True(t, errors.Is(oneConflict.Semantic, SemanticError{}))

	// This is a synthetic case, for the sake of coverage.
	err = ViewConflicts{
		rd1: []Conflict{},
	}
	require.Equal(t, noConflictsString, err.Error())

	// Note: This test ignores duplicates, one semantic error is
	// enough to test the ViewConflicts logic.
	oneError := oneConflict.Semantic.Error()

	err = ViewConflicts{
		rd1: []Conflict{
			oneConflict,
		},
	}
	require.Equal(t, oneError, err.Error())

	err = ViewConflicts{
		rd1: []Conflict{
			oneConflict,
			oneConflict,
		},
	}
	require.Equal(t, "2 conflicts, e.g. "+oneError, err.Error())

	err = ViewConflicts{
		rd1: []Conflict{
			oneConflict,
		},
		rd2: []Conflict{
			oneConflict,
		},
	}
	require.Equal(t, "2 conflicts in 2 readers, e.g. "+oneError, err.Error())
}

// TestConflictError tests that both semantic errors and duplicate
// conflicts are printed.  Note this uses the real library to generate
// the conflict, to avoid creating a relatively large test-only type.
func TestConflictError(t *testing.T) {
	rds := []*reader.Reader{
		reader.New(metrictest.NewExporter(), reader.WithDefaultAggregationKindFunc(func(k sdkinstrument.Kind) aggregation.Kind {
			return aggregation.GaugeKind
		})),
	}

	vc := New(testLib, nil, rds)

	// Create a synchronous then an asynchronous counter
	inst1, err1 := vc.Compile(testInst("foo", sdkinstrument.CounterObserverKind, number.Int64Kind))
	require.Error(t, err1)
	require.NotNil(t, inst1)
	require.Equal(t, "CounterObserverKind instrument incompatible with Gauge aggregation", err1.Error())

	inst2, err2 := vc.Compile(testInst("foo", sdkinstrument.CounterKind, number.Int64Kind))
	require.Error(t, err2)
	require.NotNil(t, inst2)
	require.Equal(t, "CounterKind instrument incompatible with Gauge aggregation; "+
		"name \"foo\" conflicts CounterObserver-Int64-Gauge, Counter-Int64-Sum", err2.Error())

	require.NotEqual(t, inst1, inst2)
}
