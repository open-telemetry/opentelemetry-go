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

package trace

import (
	"fmt"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func attributes(n int) []attribute.KeyValue {
	a := make([]attribute.KeyValue, n)
	for i := 0; i < n; i++ {
		a[i] = attribute.Int(fmt.Sprint(i), i)
	}
	return a
}
func benchmarkSetAttributes(b *testing.B, i int) {
	attrs := attributes(i)

	s := &span{
		startTime:  time.Now(),
		attributes: newAttributesMap(i),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {

		s.SetAttributes(attrs...)
	}
}

func BenchmarkSpan_SetAttributes_1(b *testing.B)    { benchmarkSetAttributes(b, 1) }
func BenchmarkSpan_SetAttributes_10(b *testing.B)   { benchmarkSetAttributes(b, 10) }
func BenchmarkSpan_SetAttributes_100(b *testing.B)  { benchmarkSetAttributes(b, 100) }
func BenchmarkSpan_SetAttributes_1000(b *testing.B) { benchmarkSetAttributes(b, 1000) }

func BenchmarkSpan_SetAttribute(b *testing.B) {
	attr := attribute.Int("1", 1)

	s := &span{
		startTime:  time.Now(),
		attributes: newAttributesMap(1),
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		s.SetAttribute(attr)
	}
}
