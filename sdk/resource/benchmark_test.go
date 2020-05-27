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

package resource_test

import (
	"fmt"
	"math/rand"
	"testing"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

const conflict = 0.5

func makeLabels(n int) (_, _ *resource.Resource) {
	used := map[string]bool{}
	l1 := make([]kv.KeyValue, n)
	l2 := make([]kv.KeyValue, n)
	for i := 0; i < n; i++ {
		var k string
		for {
			k = fmt.Sprint("k", rand.Intn(1000000000))
			if !used[k] {
				used[k] = true
				break
			}
		}
		l1[i] = kv.String(k, fmt.Sprint("v", rand.Intn(1000000000)))

		if rand.Float64() < conflict {
			l2[i] = l1[i]
		} else {
			l2[i] = kv.String(k, fmt.Sprint("v", rand.Intn(1000000000)))
		}

	}
	return resource.New(l1...), resource.New(l2...)
}

func benchmarkMergeResource(b *testing.B, size int) {
	r1, r2 := makeLabels(size)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = resource.Merge(r1, r2)
	}
}

func BenchmarkMergeResource_1(b *testing.B) {
	benchmarkMergeResource(b, 1)
}
func BenchmarkMergeResource_2(b *testing.B) {
	benchmarkMergeResource(b, 2)
}
func BenchmarkMergeResource_3(b *testing.B) {
	benchmarkMergeResource(b, 3)
}
func BenchmarkMergeResource_4(b *testing.B) {
	benchmarkMergeResource(b, 4)
}
func BenchmarkMergeResource_6(b *testing.B) {
	benchmarkMergeResource(b, 6)
}
func BenchmarkMergeResource_8(b *testing.B) {
	benchmarkMergeResource(b, 8)
}
func BenchmarkMergeResource_16(b *testing.B) {
	benchmarkMergeResource(b, 16)
}
