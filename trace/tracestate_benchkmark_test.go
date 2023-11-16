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
	"testing"
)

func BenchmarkTraceStateParse(b *testing.B) {
	for _, test := range testcases {
		b.Run(test.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				ParseTraceState(test.in)
			}
		})
	}
}

// Insert 因为语义发生变化了，我们没有做正则校验，所以一定是比 otel 的实现快的，并且因为不平等，没必要再做 benchmark 了

func BenchmarkTraceStateString(b *testing.B) {
	for _, test := range testcases {
		if len(test.tracestate.list) == 0 {
			continue
		}
		b.Run(test.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = test.tracestate.String()
			}
		})
	}
}
