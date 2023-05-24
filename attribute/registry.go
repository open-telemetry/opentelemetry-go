package attribute

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

// sets is a registry of all Set data.
// TODO: optimize initial size.
var sets = newDataRegistry(-1)

type value struct {
	// id is the hash of data. This value is stored (in addition to it being a
	// map key) so a pointer that always points to the same place in memory is
	// used. Otherwise, if a pointer to the computed hash value is used it will
	// point to different memory and the finalizer will not be set for a unique
	// identifier of the data.
	id   uint64
	data []KeyValue
}

type registry struct {
	sync.RWMutex
	data map[uint64]*value
}

func newDataRegistry(n int) *registry {
	if n <= 0 {
		return &registry{data: make(map[uint64]*value)}
	}
	return &registry{data: make(map[uint64]*value, n)}
}

// Load returns the value stored in the registry for a key, or nil if no value
// is present. The ok result indicates whether value was found in the registry.
func (r *registry) Load(key uint64) (v []KeyValue, ok bool) {
	var val *value
	r.RLock()
	val, ok = r.data[key]
	r.RUnlock()
	return val.data, ok
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
func (r *registry) Store(data []KeyValue) *uint64 {
	h := fnv.New()
	for _, kv := range data {
		h = hash(h, kv)
	}

	key := uint64(h)

	r.Lock()
	defer r.Unlock()
	for {
		// TODO: reserve 0 so empty Distinct lookups are empty.
		stored, collision := r.data[key]
		if !collision {
			v := &value{data: data, id: key}
			r.data[key] = v
			ptr := &v.id
			runtime.SetFinalizer(ptr, func(k *uint64) { r.Delete(*k) })
			return ptr
		}

		if equal(stored.data, data) {
			return &stored.id
		}

		// Re-hash until we find an open value.
		h = h.Uint64(key)
		key = uint64(h)
	}
}

func (r *registry) Delete(key uint64) {
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
func equal(a, b []KeyValue) bool {
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
