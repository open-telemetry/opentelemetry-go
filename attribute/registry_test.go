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

package attribute

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func wait(d time.Duration, done func() bool) {
	timer := time.NewTimer(d)
	defer timer.Stop()

	// Ensure on slow systems this is attempted at least once.
	runtime.GC()
	runtime.Gosched()
	if done() {
		return
	}

	func() {
		for {
			select {
			case <-timer.C:
				return
			default:
				runtime.GC()
				runtime.Gosched()
				if done() {
					return
				}
			}
		}
	}()
}

func TestRegistry(t *testing.T) {
	data0 := []KeyValue{Int("one", 1), Int("two", 2)}
	data1 := []KeyValue{Int("one", 1), Int("two", 2)}
	data2 := []KeyValue{String("A", "a"), String("B", "b")}
	reg := newRegistry(-1)

	t.Run("Store", func(t *testing.T) {
		// First entry.
		s0 := reg.newSet(data0)
		k0 := s0.id
		require.NotNil(t, k0, "invalid first key")

		assert.Equal(t, 1, reg.len(), "registry should hold only one entry")
		if assert.True(t, reg.Has(*k0), "data not stored in registry") {
			v, ok := reg.Load(*k0)
			assert.True(t, ok, "Load returned different state from Has")
			assert.Equal(t, data0, v, "incorrect data stored")
		}

		// Second entry (same value as the first).
		s1 := reg.newSet(data1)
		k1 := s1.id
		require.NotNil(t, k1, "invalid second key")

		assert.Truef(t, k0 == k1, "different keys for the same data: %v, %v", k0, k1)
		assert.Equal(t, 1, reg.len(), "registry should hold only one entry")
		if assert.True(t, reg.Has(*k0), "original data removed from registry") {
			v, ok := reg.Load(*k1)
			assert.True(t, ok, "Load returned different state from Has")
			assert.Equal(t, data0, v, "data corrupted")
		}

		// Third entry (different than the previous two).
		s2 := reg.newSet(data2)
		k2 := s2.id
		require.NotNil(t, k1, "invalid third key")

		assert.False(t, k0 == k2, "same keys for the different data")
		assert.Equal(t, 2, reg.len(), "registry should hold only two entry")
		if assert.True(t, reg.Has(*k0), "original data overwrote in registry") {
			v, ok := reg.Load(*k0)
			assert.True(t, ok, "Load returned different state from Has")
			assert.Equal(t, data0, v, "data corrupted")
		}
		if assert.True(t, reg.Has(*k2), "second data set not stored in registry") {
			v, ok := reg.Load(*k2)
			assert.True(t, ok, "Load returned different state from Has")
			assert.Equal(t, data2, v, "incorrect data stored")
		}
	})

	// Leaving the scope holding the keys should mean the GC will try to
	// reclaim them, and their finalizers should run (deleting the entries from
	// the registry).
	wait(time.Second*2, func() bool { return reg.len() == 0 })
	if !assert.Equalf(t, 0, reg.len(), "registry should be empty: %#v", reg.data) {
		// Reset manually for the next tests.
		reg = newRegistry(-1)
	}

	t.Run("Scope", func(t *testing.T) {
		var k *uint64
		{
			localS := reg.newSet(data0)
			localK := localS.id
			require.NotNil(t, localK, "invalid local key")
			assert.True(t, reg.Has(*localK), "data not stored in registry")

			// Should have no effect.
			runtime.GC()
			runtime.Gosched()
			assert.Truef(t, reg.Has(*localK), "premature clear: %#v", reg.data)

			// Copy pointer.
			k = localK

			// Should have no effect.
			runtime.GC()
			runtime.Gosched()
			assert.Truef(t, reg.Has(*k), "premature clear, copy held: %#v", reg.data)
		}

		// Should have no effect, k is still in scope.
		runtime.GC()
		runtime.Gosched()

		assert.Equal(t, 1, reg.len())
		assert.True(t, reg.Has(*k), "data cleared when reference still exist")
		runtime.KeepAlive(k)
	})

	wait(time.Second*2, func() bool { return reg.len() == 0 })
	assert.Equalf(t, 0, reg.len(), "registry should be empty: %#v", reg.data)
}
