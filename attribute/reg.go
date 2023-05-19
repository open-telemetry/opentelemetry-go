package attribute

import (
	"fmt"
	"reflect"
	"sync"

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

// sets is a registry of all Set data.
// TODO: optimize initial size.
var sets = newRegistry(10)

type registry struct {
	sync.RWMutex
	data map[uint64][]KeyValue
}

func newRegistry(n int) *registry {
	return &registry{data: make(map[uint64][]KeyValue, n)}
}

// Load returns the value stored in the registry for a key, or nil if no value
// is present. The ok result indicates whether value was found in the registry.
func (r *registry) Load(key uint64) (value []KeyValue, ok bool) {
	r.RLock()
	value, ok = r.data[key]
	r.RUnlock()
	return value, ok
}

// Store stores the value returning a unique identifying key.
//
// This assumes value is sorted consistently.
func (r *registry) Store(value []KeyValue) (key uint64) {
	h := fnv.New()
	for _, kv := range value {
		h = hash(h, kv)
	}
	key = uint64(h)

	r.Lock()
	key = r.unique(key, value)
	r.data[key] = value
	r.Unlock()

	return key
}

// unique ensures key does not collide with any existing key. If it does not,
// the original key will be returned. Otherwise, value is checked to determine
// if it is the same value within the registry. If it is, the original key is
// returned. If the key collides and the values are not equal, the key will be
// re-hashed until an unique key is found, and that unique key will be
// returned.
//
// This function assumes r.Lock is held.
func (r *registry) unique(key uint64, value []KeyValue) uint64 {
	h := fnv.Hash(key)
	for {
		stored, collision := r.data[key]
		if !collision {
			return key
		}

		if equal(stored, value) {
			return key
		}

		// Re-hash until we find an open value.
		h = h.Uint64(key)
		key = uint64(h)
	}
}

// Delete deletes the value for a key.
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
