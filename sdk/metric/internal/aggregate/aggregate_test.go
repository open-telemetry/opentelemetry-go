// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

var (
	keyUser    = "user"
	userAlice  = attribute.String(keyUser, "Alice")
	userBob    = attribute.String(keyUser, "Bob")
	userCarol  = attribute.String(keyUser, "Carol")
	userDave   = attribute.String(keyUser, "Dave")
	adminTrue  = attribute.Bool("admin", true)
	adminFalse = attribute.Bool("admin", false)

	alice = attribute.NewSet(userAlice, adminTrue)
	bob   = attribute.NewSet(userBob, adminFalse)
	carol = attribute.NewSet(userCarol, adminFalse)
	dave  = attribute.NewSet(userDave, adminFalse)

	// Filtered.
	attrFltr = func(kv attribute.KeyValue) bool {
		return kv.Key == attribute.Key(keyUser)
	}
	fltrAlice = attribute.NewSet(userAlice)
	fltrBob   = attribute.NewSet(userBob)

	// Sat Jan 01 2000 00:00:00 GMT+0000.
	y2k = time.Unix(946684800, 0)
)

// y2kPlus returns the timestamp at n seconds past Sat Jan 01 2000 00:00:00 GMT+0000.
func y2kPlus(n int64) time.Time {
	d := time.Duration(n) * time.Second
	return y2k.Add(d)
}

// clock is a test clock. It provides a predictable value for now() that can be
// reset.
type clock struct {
	ticks atomic.Int64
}

// Now returns the mocked time starting at y2kPlus(0). Each call to Now will
// increment the returned value by one second.
func (c *clock) Now() time.Time {
	old := c.ticks.Add(1) - 1
	return y2kPlus(old)
}

// Reset resets the clock c to tick from y2kPlus(0).
func (c *clock) Reset() { c.ticks.Store(0) }

// Register registers clock c's Now method as the now var. It returns an
// unregister func that should be called to restore the original now value.
func (c *clock) Register() (unregister func()) {
	orig := now
	now = c.Now
	return func() { now = orig }
}

func dropExemplars[N int64 | float64](attr attribute.Set) FilteredExemplarReservoir[N] {
	return dropReservoir[N](attr)
}

func TestBuilderFilter(t *testing.T) {
	t.Run("Int64", testBuilderFilter[int64]())
	t.Run("Float64", testBuilderFilter[float64]())
}

func testBuilderFilter[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		value, attr := N(1), alice
		run := func(b Builder[N], wantF attribute.Set, wantD []attribute.KeyValue) func(*testing.T) {
			return func(t *testing.T) {
				t.Helper()

				meas := b.filter(
					func(_ context.Context, v N, _ attribute.Distinct, f attribute.Set, getKVs func() []attribute.KeyValue, d []attribute.KeyValue) {
						assert.Equal(t, value, v, "measured incorrect value")
						if f.Len() == 0 && getKVs != nil {
							kvs := getKVs()
							if len(kvs) > 0 {
								f = attribute.NewSet(kvs...)
							}
						}
						assert.Equal(t, wantF, f, "measured incorrect filtered attributes")
						assert.ElementsMatch(t, wantD, d, "measured incorrect dropped attributes")
					},
				)

				t.Run("Set", func(t *testing.T) {
					meas(t.Context(), value, attr.Equivalent(), attr, nil)
				})

				t.Run("KVs", func(t *testing.T) {
					kvs := attr.ToSlice()
					meas(t.Context(), value, attribute.NewDistinctFromSorted(kvs), *attribute.EmptySet(), kvs)
				})
			}
		}

		t.Run("NoFilter", run(Builder[N]{}, attr, nil))
		t.Run("Filter", run(Builder[N]{Filter: attrFltr}, fltrAlice, []attribute.KeyValue{adminTrue}))
	}
}

type arg[N int64 | float64] struct {
	ctx context.Context

	value N
	attr  attribute.Set
}

type output struct {
	n   int
	agg metricdata.Aggregation
}

type teststep[N int64 | float64] struct {
	input  []arg[N]
	expect output
}

func test[N int64 | float64](meas Measure[N], comp ComputeAggregation, steps []teststep[N]) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		got := new(metricdata.Aggregation)
		for i, step := range steps {
			for _, args := range step.input {
				meas(args.ctx, args.value, args.attr.Equivalent(), args.attr, nil)
			}

			t.Logf("step: %d", i)
			assert.Equal(t, step.expect.n, comp(got), "incorrect data size")
			metricdatatest.AssertAggregationsEqual(t, step.expect.agg, *got)
		}
	}
}

func getConcurrentVals[N int64 | float64]() []N {
	// Keep length of v in sync with concurrentNumRecords
	// and expectedConcurrentSum.
	switch any(*new(N)).(type) {
	case float64:
		v := []float64{2.5, 6.1, 4.4, 10.0, 22.0, -3.5, -6.5, 3.0, -6.0}
		return any(v).([]N)
	default:
		v := []int64{2, 6, 4, 10, 22, -3, -6, 3, -6}
		return any(v).([]N)
	}
}

const (
	concurrentValsSum       = 32
	concurrentNumGoroutines = 10
	concurrentNumRecords    = 90 // Multiple of 9 (length of values sequences)
	expectedConcurrentCount = uint64(concurrentNumGoroutines * concurrentNumRecords)
)

func expectedConcurrentSum[N int64 | float64]() N {
	return N(int64(concurrentNumGoroutines) * int64(concurrentNumRecords/9) * concurrentValsSum)
}

// testAggregationConcurrentSafe provides a unified stress test for all generic aggregators
// by generating high contention, cardinality limit overflow, and validating exact results.
func testAggregationConcurrentSafe[N int64 | float64](
	meas Measure[N],
	comp ComputeAggregation,
	validate func(t *testing.T, aggs []metricdata.Aggregation),
) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ctx := t.Context()
		var wg sync.WaitGroup

		// Use 10 different attribute sets to force overflow on the AggregationLimit
		// which is typically set to 3.
		attrs := make([]attribute.Set, concurrentNumGoroutines)
		for i := range attrs {
			attrs[i] = attribute.NewSet(attribute.String(keyUser, strconv.Itoa(i)))
		}

		vals := getConcurrentVals[N]()

		wg.Add(concurrentNumGoroutines)
		for i := range concurrentNumGoroutines {
			go func(id int) {
				defer wg.Done()
				// Each goroutine records to a distinct attribute set
				attr := attrs[id]

				for j := range concurrentNumRecords {
					meas(ctx, vals[j%len(vals)], attr.Equivalent(), attr, nil)
				}
			}(i)
		}

		var results []metricdata.Aggregation

		// Run computation concurrently with measurements to stress hot/cold swaps
		wg.Go(func() {
			for range concurrentNumRecords {
				got := new(metricdata.Aggregation)
				comp(got)
				results = append(results, *got)
			}
		})

		wg.Wait()

		// Final flush to get final values
		got := new(metricdata.Aggregation)
		comp(got)
		results = append(results, *got)

		validate(t, results)
	}
}

func assertSumEqual[N int64 | float64](t *testing.T, expected, actual N) {
	if _, ok := any(*new(N)).(float64); ok {
		assert.InDelta(t, float64(expected), float64(actual), 0.0001)
	} else {
		assert.Equal(t, expected, actual)
	}
}

func benchmarkAggregate[N int64 | float64](factory func() (Measure[N], ComputeAggregation)) func(*testing.B) {
	counts := []int{1, 10, 100}
	return func(b *testing.B) {
		for _, n := range counts {
			b.Run(strconv.Itoa(n), func(b *testing.B) {
				benchmarkAggregateN(b, factory, n)
			})
		}
	}
}

var bmarkRes metricdata.Aggregation

func benchmarkAggregateN[N int64 | float64](b *testing.B, factory func() (Measure[N], ComputeAggregation), count int) {
	ctx := b.Context()
	attrs := make([]attribute.Set, count)
	for i := range attrs {
		attrs[i] = attribute.NewSet(attribute.Int("value", i))
	}

	b.Run("Measure", func(b *testing.B) {
		got := &bmarkRes
		meas, comp := factory()
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			for _, attr := range attrs {
				meas(ctx, 1, attr.Equivalent(), attr, nil)
			}
		}

		comp(got)
	})

	b.Run("ComputeAggregation", func(b *testing.B) {
		comps := make([]ComputeAggregation, b.N)
		for n := range comps {
			meas, comp := factory()
			for _, attr := range attrs {
				meas(ctx, 1, attr.Equivalent(), attr, nil)
			}
			comps[n] = comp
		}

		got := &bmarkRes
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			comps[n](got)
		}
	})
}
