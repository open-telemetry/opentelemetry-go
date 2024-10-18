// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"strconv"
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

				meas := b.filter(func(_ context.Context, v N, f attribute.Set, d []attribute.KeyValue) {
					assert.Equal(t, value, v, "measured incorrect value")
					assert.Equal(t, wantF, f, "measured incorrect filtered attributes")
					assert.ElementsMatch(t, wantD, d, "measured incorrect dropped attributes")
				})
				meas(context.Background(), value, attr)
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
				meas(args.ctx, args.value, args.attr)
			}

			t.Logf("step: %d", i)
			assert.Equal(t, step.expect.n, comp(got), "incorrect data size")
			metricdatatest.AssertAggregationsEqual(t, step.expect.agg, *got)
		}
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
	ctx := context.Background()
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
				meas(ctx, 1, attr)
			}
		}

		comp(got)
	})

	b.Run("ComputeAggregation", func(b *testing.B) {
		comps := make([]ComputeAggregation, b.N)
		for n := range comps {
			meas, comp := factory()
			for _, attr := range attrs {
				meas(ctx, 1, attr)
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
