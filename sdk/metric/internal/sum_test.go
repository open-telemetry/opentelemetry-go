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
	"sync"
	"testing"

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
	for _, a := range agg.flush() {
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

func TestInt64Sum(t *testing.T)   { testSumAggregation(t, NewSum[int64](NewInt64)) }
func TestFloat64Sum(t *testing.T) { testSumAggregation(t, NewSum[float64](NewFloat64)) }
