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
)

const routines = 5

func testAtomic[N int64 | float64](t *testing.T, a Atomic[N]) {
	n := a.Load()
	a.Add(1)
	assert.Equal(t, n+1, a.Load())

	a.Store(n)
	assert.Equal(t, n, a.Load())
}

func TestInt64(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(routines)
	for i := int64(0); i < routines; i++ {
		go func(n int64) {
			defer wg.Done()
			a := NewInt64()
			a.Store(n)
			testAtomic[int64](t, a)
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
			a := NewFloat64()
			a.Store(n)
			testAtomic[float64](t, a)
		}(float64(i))
	}
	wg.Wait()
}
