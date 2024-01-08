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
	"fmt"
	"strconv"
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
	userEve    = attribute.String(keyUser, "Eve")
	adminTrue  = attribute.Bool("admin", true)
	adminFalse = attribute.Bool("admin", false)

	alice = attribute.NewSet(userAlice, adminTrue)
	bob   = attribute.NewSet(userBob, adminFalse)
	carol = attribute.NewSet(userCarol, adminFalse)
	dave  = attribute.NewSet(userDave, adminFalse)
	eve   = attribute.NewSet(userEve, adminFalse)

	// Filtered.
	attrFltr = func(kv attribute.KeyValue) bool {
		return kv.Key == attribute.Key(keyUser)
	}
	fltrAlice = attribute.NewSet(userAlice)
	fltrBob   = attribute.NewSet(userBob)
	fltrCarol = attribute.NewSet(userCarol)

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

func TestBuilderFilter(t *testing.T) {
	t.Run("Int64", testBuilderFilter[int64]())
	t.Run("Float64", testBuilderFilter[float64]())
}

func testBuilderFilter[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		run := func(b Builder[N], wantA attribute.Set) func(*testing.T) {
			return func(t *testing.T) {
				t.Helper()

				type customAgg struct {
					metricdata.Aggregation

					Attributes attribute.Set
				}
				var agg metricdata.Aggregation = customAgg{Attributes: alice}

				comp := b.filter(func(dest *metricdata.Aggregation, fltr attribute.Filter) int {
					agg := (*dest).(customAgg)
					agg.Attributes, _ = agg.Attributes.Filter(fltr)
					*dest = agg
					return 0
				})
				comp(&agg)

				assert.Equal(t, wantA, agg.(customAgg).Attributes, "filter not applied")
			}
		}

		t.Run("NoFilter", run(Builder[N]{}, alice))
		t.Run("Filter", run(Builder[N]{Filter: attrFltr}, fltrAlice))
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

func benchmarkFiltered[N int64 | float64](factory func(attribute.Filter) (Measure[N], ComputeAggregation)) func(*testing.B) {
	nAttr := []int{1, 10, 100}       // Number of distinct attribute sets.
	nMeas := []int{1, 10, 100, 1000} // Number of measurements made per attribute set.
	return func(b *testing.B) {
		for _, attributeCap := range nAttr {
			for _, measurements := range nMeas {
				name := fmt.Sprintf("Attributes/%d/Measurements/%d", attributeCap, measurements)
				b.Run(name, func(b *testing.B) {
					attrs := make([]attribute.Set, attributeCap)
					for i := range attrs {
						attrs[i] = attribute.NewSet(
							userAlice,
							attribute.Int("value", i),
						)
					}

					got := &bmarkRes
					ctx := context.Background()
					meas, comp := factory(attrFltr)

					b.ReportAllocs()
					b.ResetTimer()
					for n := 0; n < b.N; n++ {
						for m := 0; m < measurements; m++ {
							for _, attr := range attrs {
								meas(ctx, 1, attr)
							}
						}

						assert.Equal(b, 1, comp(got), "attributes not filtered")
					}
				})
			}
		}
	}
}
