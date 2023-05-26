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

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"sync"
	"sync/atomic"
)

const zeroID uint64 = 0

type setRegistry struct {
	sync.RWMutex
	data map[uint64]*setData
}

func newSetRegistry(n int) *setRegistry {
	if n <= 0 {
		return &setRegistry{data: make(map[uint64]*setData)}
	}
	return &setRegistry{data: make(map[uint64]*setData, n)}
}

func (sr *setRegistry) len() int {
	sr.RLock()
	defer sr.RUnlock()
	return len(sr.data)
}

// Load returns the value stored in the registry for a key, or nil if no value
// is present.
//
// It is the callers responsibility to call Decrement when done with value.
func (sr *setRegistry) Load(key uint64) (value *setData) {
	sr.RLock()
	value = sr.data[key]
	value.Increment()
	sr.RUnlock()
	return value
}

// Has returns if the registry holds a value for key.
func (sr *setRegistry) Has(key uint64) (ok bool) {
	sr.RLock()
	_, ok = sr.data[key]
	sr.RUnlock()
	return ok
}

// Store stores data and returns its unique identifying key.
//
// This assumes data is sorted consistently with unique values.
func (sr *setRegistry) Store(data *[]KeyValue) uint64 {
	h := hashKVs(data)
	key := uint64(h)
	rehash := func() {
		h = h.Uint64(key)
		key = uint64(h)
	}

	sr.Lock()
	defer sr.Unlock()
	for {
		if key == zeroID {
			// Reserve zeroID so empty Distinct are empty.
			rehash()
			continue
		}

		stored, collision := sr.data[key]
		switch {
		case !collision:
			sr.data[key] = newSetData(key, sr, data)
			return key
		case equalKVs(stored.data, data):
			slicePool.Put(data)
			stored.Increment()
			return key
		}

		// Re-hash until we find an open key.
		rehash()
	}
}

// Release releases reference to the value stored with key.
func (sr *setRegistry) Release(key uint64) {
	sr.Lock()
	defer sr.Unlock()

	v := sr.data[key]
	if v.decrement() <= 0 {
		delete(sr.data, key)
		v.free()
	}
}

func (sr *setRegistry) delete(key uint64) {
	sr.Lock()
	delete(sr.data, key)
	sr.Unlock()
}

var setDataPool = sync.Pool{New: func() any { return new(setData) }}

type setData struct {
	// nRef is the reference counter.
	nRef atomic.Int64

	key uint64
	reg *setRegistry

	data *[]KeyValue
}

func newSetData(key uint64, reg *setRegistry, data *[]KeyValue) *setData {
	sd := setDataPool.Get().(*setData)
	sd.key = key
	sd.reg = reg
	sd.data = data
	sd.nRef.Store(1)
	return sd
}

func (sd *setData) Len() int {
	if sd == nil {
		return 0
	}
	return len(*sd.data)
}

func (sd *setData) Index(i int) KeyValue {
	if sd == nil {
		return KeyValue{}
	}
	return (*sd.data)[i]
}

func (sd *setData) Increment() {
	if sd == nil {
		return
	}
	sd.nRef.Add(1)
}

func (sd *setData) Decrement() {
	if sd == nil {
		return
	}
	if sd.decrement() <= 0 {
		sd.reg.delete(sd.key)
		sd.free()
	}
}

func (sd *setData) decrement() int64 { return sd.nRef.Add(-1) }

func (sd *setData) free() {
	slicePool.Put(sd.data)
	sd.data = nil
	sd.nRef.Store(0)
	setDataPool.Put(sd)
}
