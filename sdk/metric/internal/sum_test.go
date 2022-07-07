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

//go:build go1.18
// +build go1.18

package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

const goroutines = 5

func testSumAggregation[N int64 | float64](t *testing.T, agg Aggregator[N]) {
	const increments = 30
	additions := map[attribute.Set]N{
		attribute.NewSet(
			attribute.String("user", "alice"), attribute.Bool("admin", true),
		): 1,
		attribute.NewSet(
			attribute.String("user", "bob"), attribute.Bool("admin", false),
		): -1,
		attribute.NewSet(
			attribute.String("user", "carol"), attribute.Bool("admin", false),
		): 2,
	}

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < increments; j++ {
				for attrs, n := range additions {
					agg.Aggregate(n, &attrs)
				}
			}
		}()
	}
	wg.Wait()

	extra := make(map[attribute.Set]struct{})
	got := make(map[attribute.Set]N)
	flush := agg.flush()
	for _, a := range flush {
		got[*a.Attributes] = a.Value.(SingleValue[N]).Value
		extra[*a.Attributes] = struct{}{}
	}

	for attr, v := range additions {
		name := attr.Encoded(attribute.DefaultEncoder())
		t.Run(name, func(t *testing.T) {
			require.Contains(t, got, attr)
			delete(extra, attr)
			assert.Equal(t, v*increments*goroutines, got[attr])
		})
	}

	assert.Lenf(t, extra, 0, "unknown values added: %v", extra)
}

func TestInt64Sum(t *testing.T)   { testSumAggregation(t, NewSum[int64]()) }
func TestFloat64Sum(t *testing.T) { testSumAggregation(t, NewSum[float64]()) }

func benchmarkSumAggregation[N int64 | float64](b *testing.B, agg Aggregator[N], count int) {
	attrs := make([]attribute.Set, count)
	for i := range attrs {
		attrs[i] = attribute.NewSet(attribute.Int("value", i))
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, attr := range attrs {
			agg.Aggregate(1, &attr)
		}
		agg.flush()
	}
}

func BenchmarkInt64Sum(b *testing.B) {
	for _, n := range []int{10, 50, 100} {
		b.Run(fmt.Sprintf("count-%d", n), func(b *testing.B) {
			benchmarkSumAggregation(b, NewSum(NewInt64), n)
		})
	}
}
func BenchmarkFloat64Sum(b *testing.B) {
	for _, n := range []int{10, 50, 100} {
		b.Run(fmt.Sprintf("count-%d", n), func(b *testing.B) {
			benchmarkSumAggregation(b, NewSum(NewFloat64), n)
		})
	}
}

var aggsStore []Aggregation

// This isn't a perfect benchmark, because we don't get consistant writes. I would probably remove it for production.
func benchmarkSumAggregationParallel[N int64 | float64](b *testing.B, agg Aggregator[N]) {
	attrs := make([]attribute.Set, 100)
	for i := range attrs {
		attrs[i] = attribute.NewSet(attribute.Int("value", i))
	}

	ctx, cancel := context.WithCancel(context.Background())
	b.Cleanup(cancel)

	for i := 0; i < 4; i++ {
		go func(i int) {
			for {
				if ctx.Err() != nil {
					return
				}
				for j := 0; j < 25; j++ {
					agg.Aggregate(1, &attrs[i*25+j])
				}
			}
		}(i)
	}

	agg.flush()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		aggsStore = agg.flush()
		time.Sleep(time.Microsecond)

	}
}

func BenchmarkInt64SumParallel(b *testing.B) {
	benchmarkSumAggregationParallel(b, NewSum(NewInt64))
}
func BenchmarkFloat64SumParallel(b *testing.B) {
	benchmarkSumAggregationParallel(b, NewSum(NewFloat64))
}
