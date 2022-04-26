package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkinstrument"
)

func Descriptor(name string, ik sdkinstrument.Kind, nk number.Kind, opts ...instrument.Option) sdkinstrument.Descriptor {
	cfg := instrument.NewConfig(opts...)
	return sdkinstrument.NewDescriptor(name, ik, nk, cfg.Description(), cfg.Unit())
}

func Point(start, end time.Time, agg aggregation.Aggregation, kvs ...attribute.KeyValue) data.Point {
	attrs := attribute.NewSet(kvs...)
	return data.Point{
		Start:       start,
		End:         end,
		Attributes:  attrs,
		Aggregation: agg,
	}
}

func Instrument(desc sdkinstrument.Descriptor, points ...data.Point) data.Instrument {
	return data.Instrument{
		Descriptor: desc,
		Points:     points,
	}
}

func CollectScope(t *testing.T, collectors []data.Collector, seq data.Sequence) []data.Instrument {
	var output data.Scope
	return CollectScopeReuse(t, collectors, seq, &output)
}

func CollectScopeReuse(t *testing.T, collectors []data.Collector, seq data.Sequence, output *data.Scope) []data.Instrument {
	output.Reset()
	for _, coll := range collectors {
		coll.Collect(seq, &output.Instruments)
	}
	return output.Instruments
}

// RequireEqualMetrics checks that an output equals the expected
// output, where the points are taken to be unordered.  Instrument
// order is expected to match because the compiler preserves the order
// of views and instruments as they are compiled.
func RequireEqualMetrics(
	t *testing.T,
	output []data.Instrument,
	expected ...data.Instrument) {
	t.Helper()

	require.Equal(t, len(expected), len(output))

	for idx := range output {
		require.Equal(t, expected[idx].Descriptor, output[idx].Descriptor)
		require.ElementsMatch(t, expected[idx].Points, output[idx].Points)
	}
}

func OTelErrors() *[]error {
	errors := new([]error)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		*errors = append(*errors, err)
	}))
	return errors
}
