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

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"strconv"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

var (
	keyUser    = "user"
	userAlice  = attribute.String(keyUser, "Alice")
	userBob    = attribute.String(keyUser, "Bob")
	adminTrue  = attribute.Bool("admin", true)
	adminFalse = attribute.Bool("admin", false)

	alice = attribute.NewSet(userAlice, adminTrue)
	bob   = attribute.NewSet(userBob, adminFalse)

	// Filtered.
	attrFltr = func(kv attribute.KeyValue) bool {
		return kv.Key == attribute.Key(keyUser)
	}
	fltrAlice = attribute.NewSet(userAlice)
	fltrBob   = attribute.NewSet(userBob)

	// Sat Jan 01 2000 00:00:00 GMT+0000.
	staticTime    = time.Unix(946684800, 0)
	staticNowFunc = func() time.Time { return staticTime }
	// Pass to t.Cleanup to override the now function with staticNowFunc and
	// revert once the test completes. E.g. t.Cleanup(mockTime(now)).
	mockTime = func(orig func() time.Time) (cleanup func()) {
		now = staticNowFunc
		return func() { now = orig }
	}
)

func newRes[N int64 | float64]() exemplar.Reservoir[N] {
	fltr := func(_ context.Context, v N, attr attribute.Set) bool {
		return attr == alice && v == 2
	}

	return exemplar.Filtered(exemplar.FixedSize[N](10), fltr)
}

func TestAggregate(t *testing.T) {
	t.Cleanup(mockTime(now))

	t.Run("Int64/LastValue", testLastValue[int64]())
	t.Run("Float64/LastValue", testLastValue[float64]())

	t.Run("Int64/DeltaSum", testDeltaSum[int64]())
	t.Run("Float64/DeltaSum", testDeltaSum[float64]())

	t.Run("Int64/CumulativeSum", testCumulativeSum[int64]())
	t.Run("Float64/CumulativeSum", testCumulativeSum[float64]())

	t.Run("Int64/DeltaPrecomputedSum", testDeltaPrecomputedSum[int64]())
	t.Run("Float64/DeltaPrecomputedSum", testDeltaPrecomputedSum[float64]())

	t.Run("Int64/CumulativePrecomputedSum", testCumulativePrecomputedSum[int64]())
	t.Run("Float64/CumulativePrecomputedSum", testCumulativePrecomputedSum[float64]())
}

func testDeltaSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:   metricdata.DeltaTemporality,
		Filter:        attrFltr,
		ReservoirFunc: newRes[N],
	}.Sum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, alice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.DeltaTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      4,
						Exemplars: []metricdata.Exemplar[N]{{
							FilteredAttributes: []attribute.KeyValue{adminTrue},
							Time:               staticTime,
							Value:              2,
						}},
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      -11,
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.DeltaTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      10,
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      3,
					},
				},
			},
		},
	})
}

func testCumulativeSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:   metricdata.CumulativeTemporality,
		Filter:        attrFltr,
		ReservoirFunc: newRes[N],
	}.Sum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, alice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      4,
						Exemplars: []metricdata.Exemplar[N]{{
							FilteredAttributes: []attribute.KeyValue{adminTrue},
							Time:               staticTime,
							Value:              2,
						}},
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      -11,
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      14,
						Exemplars: []metricdata.Exemplar[N]{{
							FilteredAttributes: []attribute.KeyValue{adminTrue},
							Time:               staticTime,
							Value:              2,
						}},
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      -8,
					},
				},
			},
		},
	})
}

func testDeltaPrecomputedSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:   metricdata.DeltaTemporality,
		Filter:        attrFltr,
		ReservoirFunc: newRes[N],
	}.PrecomputedSum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.DeltaTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      3,
						Exemplars: []metricdata.Exemplar[N]{{
							FilteredAttributes: []attribute.KeyValue{adminTrue},
							Time:               staticTime,
							Value:              2,
						}},
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      -10,
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice},
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.DeltaTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      8,
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      13,
					},
				},
			},
		},
	})
}

func testCumulativePrecomputedSum[N int64 | float64]() func(t *testing.T) {
	mono := false
	in, out := Builder[N]{
		Temporality:   metricdata.CumulativeTemporality,
		Filter:        attrFltr,
		ReservoirFunc: newRes[N],
	}.PrecomputedSum(mono)
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      3,
						Exemplars: []metricdata.Exemplar[N]{{
							FilteredAttributes: []attribute.KeyValue{adminTrue},
							Time:               staticTime,
							Value:              2,
						}},
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      -10,
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 1, fltrAlice},
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			want: metricdata.Sum[N]{
				IsMonotonic: mono,
				Temporality: metricdata.CumulativeTemporality,
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      11,
						Exemplars: []metricdata.Exemplar[N]{{
							FilteredAttributes: []attribute.KeyValue{adminTrue},
							Time:               staticTime,
							Value:              2,
						}},
					},
					{
						Attributes: fltrBob,
						StartTime:  staticTime,
						Time:       staticTime,
						Value:      3,
					},
				},
			},
		},
	})
}

func testLastValue[N int64 | float64]() func(t *testing.T) {
	in, out := Builder[N]{
		Filter:        attrFltr,
		ReservoirFunc: newRes[N],
	}.LastValue()
	ctx := context.Background()
	return test[N](in, out, []teststep[N]{
		{
			input: []arg[N]{
				{ctx, 1, alice},
				{ctx, -1, bob},
				{ctx, 1, fltrAlice},
				{ctx, 2, alice},
				{ctx, -10, bob},
			},
			want: metricdata.Gauge[N]{
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						Time:       staticTime,
						Value:      2,
						Exemplars: []metricdata.Exemplar[N]{{
							FilteredAttributes: []attribute.KeyValue{adminTrue},
							Time:               staticTime,
							Value:              2,
						}},
					},
					{
						Attributes: fltrBob,
						Time:       staticTime,
						Value:      -10,
					},
				},
			},
		},
		{
			input: []arg[N]{
				{ctx, 10, alice},
				{ctx, 3, bob},
			},
			want: metricdata.Gauge[N]{
				DataPoints: []metricdata.DataPoint[N]{
					{
						Attributes: fltrAlice,
						Time:       staticTime,
						Value:      10,
					},
					{
						Attributes: fltrBob,
						Time:       staticTime,
						Value:      3,
					},
				},
			},
		},
	})
}

type arg[N int64 | float64] struct {
	ctx   context.Context
	value N
	attr  attribute.Set
}

type teststep[N int64 | float64] struct {
	input []arg[N]
	want  metricdata.Aggregation
}

func test[N int64 | float64](in Input[N], out Output, steps []teststep[N]) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		got := new(metricdata.Aggregation)
		for _, step := range steps {
			for _, args := range step.input {
				in(args.ctx, args.value, args.attr)
			}

			out(got)
			metricdatatest.AssertAggregationsEqual(t, step.want, *got)
		}
	}
}

var bmarkResults metricdata.Aggregation

func benchmarkAggregatorN[N int64 | float64](b *testing.B, factory func() (Input[N], Output), count int) {
	ctx := context.Background()
	attrs := make([]attribute.Set, count)
	for i := range attrs {
		attrs[i] = attribute.NewSet(attribute.Int("value", i))
	}

	b.Run("Aggregate", func(b *testing.B) {
		got := &bmarkResults
		in, out := factory()
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			for _, attr := range attrs {
				in(ctx, 1, attr)
			}
		}

		out(got)
	})

	b.Run("Aggregations", func(b *testing.B) {
		outs := make([]Output, b.N)
		for n := range outs {
			in, out := factory()
			for _, attr := range attrs {
				in(ctx, 1, attr)
			}
			outs[n] = out
		}

		got := &bmarkResults
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			outs[n](got)
		}
	})
}

func benchmarkAggregator[N int64 | float64](factory func() (Input[N], Output)) func(*testing.B) {
	counts := []int{1, 10, 100}
	return func(b *testing.B) {
		for _, n := range counts {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				benchmarkAggregatorN(b, factory, n)
			})
		}
	}
}

func BenchmarkSum(b *testing.B) {
	b.Run("Int64", benchmarkSum[int64])
}

func benchmarkSum[N int64 | float64](b *testing.B) {
	// The monotonic argument is only used to annotate the Sum returned from
	// the Aggregation method. It should not have an effect on operational
	// performance, therefore, only monotonic=false is benchmarked here.
	factory := func() (Input[N], Output) {
		b := Builder[N]{Temporality: metricdata.DeltaTemporality}
		return b.Sum(false)
	}
	b.Run("Delta", benchmarkAggregator(factory))

	factory = func() (Input[N], Output) {
		b := Builder[N]{
			Temporality: metricdata.DeltaTemporality,
			Filter: func(kv attribute.KeyValue) bool {
				return kv.Key == attribute.Key(keyUser)
			},
		}
		return b.Sum(false)
	}
	b.Run("Filtered", benchmarkAggregator(factory))
}
