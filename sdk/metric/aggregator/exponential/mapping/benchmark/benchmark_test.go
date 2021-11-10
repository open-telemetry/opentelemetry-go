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

package benchmark

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
			_, _ = mapper.LowerBoundary(src.Int63())
		}
	})
}

func BenchmarkMapping(b *testing.B) {
	// None of these have time complexity dependent on scale.
	benchmarkMapping(b, "exponent", exponent.NewMapping(-1))
	benchmarkMapping(b, "logarithm", logarithm.NewMapping(3))
}

func BenchmarkBoundary(b *testing.B) {
	// None of these have time complexity dependent on scale.
	benchmarkBoundary(b, "exponent", exponent.NewMapping(-1))
	benchmarkBoundary(b, "logarithm", logarithm.NewMapping(3))
}
