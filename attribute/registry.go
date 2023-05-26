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
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

// sets is a registry of all Set data.
// TODO: optimize initial size.
var sets = newRegistry(-1)

func newSet(data *[]KeyValue) Set {
	s := sets.Store(data)
	// Set the finalizer here so putID can be defined at the package scope.
	// Otherwise, and anonymous function, or method, is needed and that will
	// make an allocation.
	runtime.SetFinalizer(s.id, putID)
	return s
}

var idPool = sync.Pool{New: func() any { return new(uint64) }}

func getID() *uint64 {
	return idPool.Get().(*uint64)
}

func putID(id *uint64) {
	sets.Release(*id)
	idPool.Put(id)
}

var referencePool = sync.Pool{
	New: func() any { return new(reference) },
}

type reference struct {
	// nRef is the reference counter.
	nRef atomic.Int64

	key uint64
	reg *registry

	data *[]KeyValue
}

func newReference(key uint64, reg *registry, data *[]KeyValue) *reference {
	r := referencePool.Get().(*reference)
	r.key = key
	r.reg = reg
	r.data = data
	r.nRef.Store(1)
	return r
}

func (r *reference) Len() int {
	if r == nil {
		return 0
	}
	return len(*r.data)
}

func (r *reference) Index(i int) KeyValue {
	if r == nil {
		return KeyValue{}
	}
	return (*r.data)[i]
}

func (r *reference) Increment() {
	if r == nil {
		return
	}
	r.nRef.Add(1)
}

func (r *reference) Decrement() {
	if r == nil {
		return
	}
	if r.decrement() <= 0 {
		r.reg.delete(r.key)
		r.free()
	}
}

func (r *reference) decrement() int64 { return r.nRef.Add(-1) }

func (r *reference) free() {
	slicePool.Put(r.data)
	r.data = nil
	r.nRef.Store(0)
	referencePool.Put(r)
}

type registry struct {
	sync.RWMutex
	data map[uint64]*reference
}

func newRegistry(n int) *registry {
	if n <= 0 {
		return &registry{data: make(map[uint64]*reference)}
	}
	return &registry{data: make(map[uint64]*reference, n)}
}

func (r *registry) len() int {
	r.RLock()
	defer r.RUnlock()
	return len(r.data)
}

// Load returns the value stored in the registry for a key, or nil if no value
// is present.
//
// It is the callers responsibility to call Decrement when done with value.
func (r *registry) Load(key uint64) (value *reference) {
	r.RLock()
	value = r.data[key]
	value.Increment()
	r.RUnlock()
	return value
}

func (r *registry) Has(key uint64) (ok bool) {
	r.RLock()
	_, ok = r.data[key]
	r.RUnlock()
	return ok
}

// Store stores data and returns a pointer to its unique identifying key.
//
// This assumes data is sorted consistently with unique values.
func (r *registry) Store(data *[]KeyValue) Set {
	h := fnv.New()
	for _, kv := range *data {
		h = hash(h, kv)
	}

	key := uint64(h)

	r.Lock()
	defer r.Unlock()
	for {
		// TODO: reserve 0 so empty Distinct lookups are empty.
		stored, collision := r.data[key]
		switch {
		case !collision:
			r.data[key] = newReference(key, r, data)
			id := getID()
			*id = key
			return Set{id: id}
		case equal(stored.data, data):
			slicePool.Put(data)
			stored.Increment()
			id := getID()
			*id = key
			return Set{id: id}
		}

		// Re-hash until we find an open value.
		h = h.Uint64(key)
		key = uint64(h)
	}
}

// Release reference to the value stored with key.
func (r *registry) Release(key uint64) {
	r.Lock()
	v := r.data[key]
	if v.decrement() <= 0 {
		delete(r.data, key)
		v.free()
	}
	r.Unlock()
}

func (r *registry) delete(key uint64) {
	r.Lock()
	delete(r.data, key)
	r.Unlock()
}

// hash returns the hash of kv with h as the base.
func hash(h fnv.Hash, kv KeyValue) fnv.Hash {
	h = h.String(string(kv.Key))

	switch kv.Value.Type() {
	case BOOL:
		h = h.Bool(kv.Value.AsBool())
	case INT64:
		h = h.Int64(kv.Value.AsInt64())
	case FLOAT64:
		h = h.Float64(kv.Value.AsFloat64())
	case STRING:
		h = h.String(kv.Value.AsString())
	case BOOLSLICE:
		// Avoid allocating a new []bool with AsBoolSlice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Bool(rv.Index(i).Bool())
		}
	case INT64SLICE:
		// Avoid allocating a new []int64 with AsInt64Slice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Int64(rv.Index(i).Int())
		}
	case FLOAT64SLICE:
		// Avoid allocating a new []float64 with AsFloat64Slice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.Float64(rv.Index(i).Float())
		}
	case STRINGSLICE:
		// Avoid allocating a new []string with AsStringSlice.
		rv := reflect.ValueOf(kv.Value.slice)
		for i := 0; i < rv.Len(); i++ {
			h = h.String(rv.Index(i).String())
		}
	default:
		// Logging is an alternative, but using the internal logger here
		// causes an import cycle so it is not done.
		v := kv.Value.AsInterface()
		msg := fmt.Sprintf("unknown value type: %[1]v (%[1]T)", v)
		panic(msg)
	}
	return h
}

// equal returns if the sorted slices of []KeyValue a and b are equal, or not.
func equal(aPtr, bPtr *[]KeyValue) bool {
	a, b := *aPtr, *bPtr
	if len(a) != len(b) {
		return false
	}

	for i, aKV := range a {
		bKV := b[i]

		if aKV.Key != bKV.Key {
			return false
		}

		aVal, bVal := aKV.Value, bKV.Value
		if aVal.Type() != bVal.Type() {
			return false
		}

		switch aVal.Type() {
		case BOOL, INT64, FLOAT64:
			if aVal.numeric != bVal.numeric {
				return false
			}
		case STRING:
			if aVal.stringly != bVal.stringly {
				return false
			}
		case BOOLSLICE, INT64SLICE, FLOAT64SLICE, STRINGSLICE:
			if aVal.slice != bVal.slice {
				return false
			}
		default:
			// Logging is an alternative, but using the internal logger here
			// causes an import cycle so it is not done.
			v := aVal.AsInterface()
			msg := fmt.Sprintf("unknown value type: %[1]v (%[1]T)", v)
			panic(msg)
		}
	}

	return true
}
