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

package internal

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const routines = 5

func testAtomic[N int64 | float64](t *testing.T, a Atomic[N]) {
	n, clone := a.Load(), a.Clone()
	require.Equal(t, n, clone.Load(), "Clone() did not copy value")

	clone.Add(1)
	assert.Equal(t, n, a.Load(), "Clone() returned original")
	assert.Equal(t, n+1, clone.Load())

	clone.Store(n)
	assert.Equal(t, n, clone.Load())
}

func TestInt64(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(routines)
	for i := int64(0); i < routines; i++ {
		go func(n int64) {
			defer wg.Done()
			testAtomic[int64](t, NewInt64(n))
		}(i)
	}
	wg.Wait()
}

func TestFloat64(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(routines)
	for i := 0; i < routines; i++ {
		go func(n float64) {
			defer wg.Done()
			testAtomic[float64](t, NewFloat64(n))
		}(float64(i))
	}
	wg.Wait()
}
