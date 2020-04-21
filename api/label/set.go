// Copyright 2019, OpenTelemetry Authors
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

package label // import "go.opentelemetry.io/otel/api/label"

import (
	"reflect"
	"sort"
	"sync"

	"go.opentelemetry.io/otel/api/core"
)

const maxConcurrentEncoders = 3

type (
	// Set is the internal representation for LabelSet.  It manages an
	// immutable set of labels with an internal cache for storing encoded
	// labels.
	Set struct {
		equivalent Equivalent

		lock     sync.Mutex
		encoders [maxConcurrentEncoders]EncoderID
		encoded  [maxConcurrentEncoders]string
	}

	// Iterator allows iterating over an ordered set of labels.
	Iterator struct {
		storage Storage
		idx     int
	}

	// Storage provides an access to the ordered labels.
	Storage interface {
		// Len returns a number of labels in the storage.
		Len() int
		// Get gets a label from a passed index.
		Get(int) (core.KeyValue, bool)
	}

	Sortable []core.KeyValue

	Equivalent struct {
		iface interface{}
	}
)

var (
	_ Storage = (*Set)(nil)

	keyValueType = reflect.TypeOf(core.KeyValue{})

	emptySet = Set{
		equivalent: Equivalent{
			iface: [0]core.KeyValue{},
		},
	}
)

func EmptySet() Set {
	return emptySet
}

func (e Equivalent) reflect() reflect.Value {
	return reflect.ValueOf(e.iface)
}

// Len returns the number of labels in this set.
func (l *Set) Len() int {
	if l == nil {
		return 0
	}
	return l.equivalent.reflect().Len()
}

// Get returns the KeyValue at ordered position `idx` in this set.
func (l *Set) Get(idx int) (core.KeyValue, bool) {
	if l == nil {
		return core.KeyValue{}, false
	}
	value := l.equivalent.reflect()

	if idx >= 0 && idx < value.Len() {
		// Note: The Go compiler successfully avoids an allocation for
		// the interface{} conversion here:
		return value.Index(idx).Interface().(core.KeyValue), true
	}

	return core.KeyValue{}, false
}

// Value returns the value of a specified key in this set.
func (l *Set) Value(k core.Key) (core.Value, bool) {
	if l == nil {
		return core.Value{}, false
	}
	value := l.equivalent.reflect()
	vlen := value.Len()

	idx := sort.Search(vlen, func(idx int) bool {
		return value.Index(idx).Interface().(core.KeyValue).Key >= k
	})
	if idx >= vlen {
		return core.Value{}, false
	}
	kv := value.Index(idx).Interface().(core.KeyValue)
	if k == kv.Key {
		return kv.Value, true
	}
	return core.Value{}, false
}

// HasValue tests whether a key is defined in this set.
func (l *Set) HasValue(k core.Key) bool {
	if l == nil {
		return false
	}
	_, ok := l.Value(k)
	return ok
}

func (l *Set) Iter() Iterator {
	return Iterator{
		storage: l,
		idx:     -1,
	}
}

// Equivalent returns a value that may be used as a map key.
// Equivalent guarantees that the result will equal the Equivalent
// value of any label set with the same elements as this set.
func (l *Set) Equivalent() Equivalent {
	if l == nil {
		return emptySet.equivalent
	}
	return l.equivalent
}

// Equals returns true if the argument set is equivalent to this set.
func (l *Set) Equals(o *Set) bool {
	return l.Equivalent() == o.Equivalent()
}

// Encoded is a pre-encoded form of the ordered labels.
func (l *Set) Encoded(encoder Encoder) string {
	if l == nil || encoder == nil {
		return ""
	}

	id := encoder.ID()
	if !id.Valid() {
		// Invalid IDs are not cached.
		return encoder.Encode(l.Iter())
	}

	var lookup *string
	l.lock.Lock()
	for idx := 0; idx < maxConcurrentEncoders; idx++ {
		if l.encoders[idx] == id {
			lookup = &l.encoded[idx]
			break
		}
	}
	l.lock.Unlock()

	if lookup != nil {
		return *lookup
	}

	r := encoder.Encode(l.Iter())

	l.lock.Lock()
	defer l.lock.Unlock()

	for idx := 0; idx < maxConcurrentEncoders; idx++ {
		if l.encoders[idx] == id {
			return l.encoded[idx]
		}
		if !l.encoders[idx].Valid() {
			l.encoders[idx] = id
			l.encoded[idx] = r
			return r
		}
	}

	// TODO: This is a performance cliff.  Find a way for this to
	// generate a warning.
	return r
}

func NewSet(kvs ...core.KeyValue) Set {
	return NewSetWithSortable(kvs, new(Sortable))
}

func NewSetWithSortable(kvs []core.KeyValue, tmp *Sortable) Set {
	// Check for empty set.
	if len(kvs) == 0 {
		return emptySet
	}

	*tmp = kvs

	// Sort and de-duplicate.  Note: this use of `tmp`
	// avoids an allocation because it is a pointer.
	sort.Stable(tmp)

	*tmp = nil

	// Modify the input slice, ensuring the following:
	// - Last-value wins semantics
	// - Caller sees the reordering, but doesn't lose values
	// - Repeated call preserve last-value wins.
	// This requires that the stable result be placed in the end of the
	// input slice, while duplicate values are swapped to the beginning.

	position := len(kvs) - 1
	offset := position - 1

	for ; offset >= 0; offset-- {
		if kvs[offset].Key == kvs[position].Key {
			continue
		}
		kvs[offset], kvs[position-1] = kvs[position-1], kvs[offset]
		position = position - 1
	}

	return Set{
		equivalent: computeEquivalent(kvs[position:]),
	}
}

func computeEquivalent(kvs []core.KeyValue) Equivalent {
	iface := computeEquivalentFixed(kvs)
	if iface == nil {
		iface = computeEquivalentReflect(kvs)
	}
	return Equivalent{
		iface: iface,
	}
}

func computeEquivalentFixed(kvs []core.KeyValue) interface{} {
	switch len(kvs) {
	case 1:
		ptr := new([1]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 2:
		ptr := new([2]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 3:
		ptr := new([3]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 4:
		ptr := new([4]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 5:
		ptr := new([5]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 6:
		ptr := new([6]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 7:
		ptr := new([7]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 8:
		ptr := new([8]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 9:
		ptr := new([9]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	case 10:
		ptr := new([10]core.KeyValue)
		copy((*ptr)[:], kvs)
		return *ptr
	default:
		return nil
	}
}

func computeEquivalentReflect(kvs []core.KeyValue) interface{} {
	at := reflect.New(reflect.ArrayOf(len(kvs), keyValueType)).Elem()
	for i, kv := range kvs {
		*(at.Index(i).Addr().Interface().(*core.KeyValue)) = kv
	}
	return at.Interface()
}

func (l *Sortable) Len() int {
	return len(*l)
}

func (l *Sortable) Swap(i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}

func (l *Sortable) Less(i, j int) bool {
	return (*l)[i].Key < (*l)[j].Key
}
