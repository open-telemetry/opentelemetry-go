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
	"context"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkSpanIDGeneration(b *testing.B) {
	b.ReportAllocs()
	run := func(b *testing.B, workers int) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		x := defaultIDGenerator()
		var wg sync.WaitGroup
		var count int64
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}
					val := atomic.AddInt64(&count, 1)
					if val == int64(b.N) {
						cancel()
						return
					} else if val > int64(b.N) {
						return
					}
					x.NewIDs(ctx)
				}
			}()
		}
		wg.Wait()
	}
	b.Run("1", func(b *testing.B) {
		run(b, 1)
	})
	b.Run("2", func(b *testing.B) {
		run(b, 2)
	})
	b.Run("3", func(b *testing.B) {
		run(b, 3)
	})
	b.Run("5", func(b *testing.B) {
		run(b, 5)
	})
	b.Run("10", func(b *testing.B) {
		run(b, 10)
	})
	b.Run("100", func(b *testing.B) {
		run(b, 100)
	})
	b.Run("1000", func(b *testing.B) {
		run(b, 1000)
	})
}
