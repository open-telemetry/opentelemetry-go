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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	k0, k1 := "one", "two"
	v0, v1 := 1, 2

	c := cache[string, int]{}

	var got int
	require.NotPanics(t, func() {
		got = c.Lookup(k0, func() int { return v0 })
	}, "zero-value cache panics on Lookup")
	assert.Equal(t, v0, got, "zero-value cache did not return fallback")

	assert.Equal(t, v0, c.Lookup(k0, func() int { return v1 }), "existing key")

	assert.Equal(t, v1, c.Lookup(k1, func() int { return v1 }), "non-existing key")
}

func TestCacheConcurrentSafe(t *testing.T) {
	const (
		key        = "k"
		goroutines = 10
	)

	c := cache[string, int]{}
	var wg sync.WaitGroup
	for n := 0; n < goroutines; n++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			assert.NotPanics(t, func() {
				c.Lookup(key, func() int { return i })
			})
		}(n)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		assert.Fail(t, "timeout")
	}
}
