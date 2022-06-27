// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exponential

import (
	"fmt"
	"math/rand"
	"testing"

	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/exponent"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exponential/mapping/logarithm"
)

func benchmarkMapping(b *testing.B, name string, mapper mapping.Mapping) {
	b.Run(fmt.Sprintf("mapping_%s", name), func(b *testing.B) {
		src := rand.New(rand.NewSource(54979))

		for i := 0; i < b.N; i++ {
			_ = mapper.MapToIndex(1 + src.Float64())
		}
	})
}

func benchmarkBoundary(b *testing.B, name string, mapper mapping.Mapping) {
	b.Run(fmt.Sprintf("boundary_%s", name), func(b *testing.B) {
		src := rand.New(rand.NewSource(54979))

		for i := 0; i < b.N; i++ {
			_, _ = mapper.LowerBoundary(int32(src.Int63()))
		}
	})
}

// An earlier draft of this benchmark included a lookup-table based
// implementation:
// https://github.com/open-telemetry/opentelemetry-go-contrib/pull/1353
// That mapping function uses O(2^scale) extra space and falls
// somewhere between the exponent and logarithm methods compared here.
// In the test, lookuptable was 40% faster than logarithm, which did
// not justify the significant extra complexity.

// Benchmarks the MapToIndex function.
func BenchmarkMapping(b *testing.B) {
	em, _ := exponent.NewMapping(-1)
	lm, _ := logarithm.NewMapping(1)
	benchmarkMapping(b, "exponent", em)
	benchmarkMapping(b, "logarithm", lm)
}

// Benchmarks the LowerBoundary function.
func BenchmarkReverseMapping(b *testing.B) {
	em, _ := exponent.NewMapping(-1)
	lm, _ := logarithm.NewMapping(1)
	benchmarkBoundary(b, "exponent", em)
	benchmarkBoundary(b, "logarithm", lm)
}
