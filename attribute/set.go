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
	"encoding/json"
	"reflect"
	"sort"
	"sync"
)

var slicePool = sync.Pool{New: func() any { return new([]KeyValue) }}

func getSlice(length, capacity int) *[]KeyValue {
	v := slicePool.Get().(*[]KeyValue)
	if cap(*v) < capacity {
		*v = make([]KeyValue, length, capacity)
	}
	*v = (*v)[:length]
	return v
}

type (
	// Set is the representation for a distinct attribute set. It manages an
	// immutable set of attributes, with an internal cache for storing
	// attribute encodings.
	//
	// This type supports the Equivalent method of comparison using values of
	// type Distinct.
	Set struct {
		// id is the key to the sets registry where this set's data is.
		//
		// A pointer is used so the finalizer can handle reference-counting for
		// the sets registry while still being optimized as a map key.
		id *uint64
	}

	// Distinct wraps a variable-size array of KeyValue, constructed with keys
	// in sorted order. This can be used as a map key or for equality checking
	// between Sets.
	Distinct struct {
		// id is the key to the sets registry where this set's data is.
		id uint64
	}

	// Filter supports removing certain attributes from attribute sets. When
	// the filter returns true, the attribute will be kept in the filtered
	// attribute set. When the filter returns false, the attribute is excluded
	// from the filtered attribute set, and the attribute instead appears in
	// the removed list of excluded attributes.
	Filter func(KeyValue) bool

	// Sortable implements sort.Interface, used for sorting KeyValue. This is
	// an exported type to support a memory optimization. A pointer to one of
	// these is needed for the call to sort.Stable(), which the caller may
	// provide in order to avoid an allocation. See NewSetWithSortable().
	Sortable []KeyValue
)

var (
	// keyValueType is used in computeDistinctReflect.
	keyValueType = reflect.TypeOf(KeyValue{})

	// emptySet is returned for empty attribute sets.
	emptySet = &Set{}

	// sortables is a pool of Sortables used to create Sets with a user does
	// not provide one.
	sortables = sync.Pool{
		New: func() interface{} { return new(Sortable) },
	}
)

// EmptySet returns a reference to a Set with no elements.
//
// This is a convenience provided for optimized calling utility.
func EmptySet() *Set {
	return emptySet
}

// Valid returns true if this value refers to a valid Set.
func (d Distinct) Valid() bool { return sets.Has(d.id) }

// Len returns the number of attributes in this set.
func (l *Set) Len() int {
	if l == nil || l.id == nil {
		return 0
	}
	v, ok := sets.Load(*l.id)
	if !ok {
		return 0
	}
	return len(*v)
}

// Get returns the KeyValue at ordered position idx in this set.
func (l *Set) Get(idx int) (KeyValue, bool) {
	if l == nil || l.id == nil {
		return KeyValue{}, false
	}
	v, ok := sets.Load(*l.id)
	if ok && idx >= 0 && idx < len(*v) {
		// Note: The Go compiler successfully avoids an allocation for
		// the interface{} conversion here:
		return (*v)[idx], true
	}

	return KeyValue{}, false
}

// Value returns the value of a specified key in this set.
func (l *Set) Value(k Key) (Value, bool) {
	if l == nil || l.id == nil {
		return Value{}, false
	}
	v, ok := sets.Load(*l.id)
	if !ok {
		return Value{}, false
	}
	idx := sort.Search(len(*v), func(idx int) bool {
		return (*v)[idx].Key >= k
	})
	if idx >= len(*v) {
		return Value{}, false
	}
	if k == (*v)[idx].Key {
		return (*v)[idx].Value, true
	}
	return Value{}, false
}

// HasValue tests whether a key is defined in this set.
func (l *Set) HasValue(k Key) bool {
	_, ok := l.Value(k)
	return ok
}

// Iter returns an iterator for visiting the attributes in this set.
func (l *Set) Iter() Iterator {
	if l == nil || l.id == nil {
		return Iterator{idx: -1}
	}
	v, _ := sets.Load(*l.id)
	return Iterator{
		storage: *v,
		idx:     -1,
	}
}

// ToSlice returns the set of attributes belonging to this set, sorted, where
// keys appear no more than once.
func (l *Set) ToSlice() []KeyValue {
	if l == nil || l.id == nil {
		return nil
	}
	v, ok := sets.Load(*l.id)
	if !ok {
		return nil
	}
	// Ensure our copy is immutable.
	dest := make([]KeyValue, len(*v))
	copy(dest, *v)
	return dest
}

func (l *Set) toSlice() *[]KeyValue {
	if l == nil || l.id == nil {
		return nil
	}
	v, ok := sets.Load(*l.id)
	if !ok {
		return nil
	}
	dest := getSlice(len(*v), len(*v))
	copy(*dest, *v)
	return dest
}

// Equivalent returns a value that may be used as a map key. The Distinct type
// guarantees that the result will equal the equivalent. Distinct value of any
// attribute set with the same elements as this, where sets are made unique by
// choosing the last value in the input for any given key.
func (l *Set) Equivalent() Distinct {
	if l == nil || l.id == nil {
		return Distinct{}
	}
	return Distinct{*l.id}
}

// Equals returns true if the argument set is equivalent to this set.
func (l *Set) Equals(o *Set) bool {
	if l == nil || o == nil {
		return l == o
	}
	if l.id == nil || o.id == nil {
		return l.id == o.id
	}
	return *l.id == *o.id
}

// Encoded returns the encoded form of this set, according to encoder.
func (l *Set) Encoded(encoder Encoder) string {
	if l == nil || encoder == nil {
		return ""
	}

	return encoder.Encode(l.Iter())
}

func empty() Set {
	return Set{}
}

// NewSet returns a new Set. See the documentation for
// NewSetWithSortableFiltered for more details.
//
// Except for empty sets, this method adds an additional allocation compared
// with calls that include a Sortable.
func NewSet(kvs ...KeyValue) Set {
	// Check for empty set.
	if len(kvs) == 0 {
		return empty()
	}
	srt := sortables.Get().(*Sortable)
	s, _ := NewSetWithSortableFiltered(kvs, srt, nil)
	sortables.Put(srt)
	return s
}

// NewSetWithSortable returns a new Set. See the documentation for
// NewSetWithSortableFiltered for more details.
//
// This call includes a Sortable option as a memory optimization.
func NewSetWithSortable(kvs []KeyValue, tmp *Sortable) Set {
	// Check for empty set.
	if len(kvs) == 0 {
		return empty()
	}
	s, _ := NewSetWithSortableFiltered(kvs, tmp, nil)
	return s
}

// NewSetWithFiltered returns a new Set. See the documentation for
// NewSetWithSortableFiltered for more details.
//
// This call includes a Filter to include/exclude attribute keys from the
// return value. Excluded keys are returned as a slice of attribute values.
func NewSetWithFiltered(kvs []KeyValue, filter Filter) (Set, []KeyValue) {
	// Check for empty set.
	if len(kvs) == 0 {
		return empty(), nil
	}
	srt := sortables.Get().(*Sortable)
	s, filtered := NewSetWithSortableFiltered(kvs, srt, filter)
	sortables.Put(srt)
	return s, filtered
}

// NewSetWithSortableFiltered returns a new Set.
//
// Duplicate keys are eliminated by taking the last value.  This
// re-orders the input slice so that unique last-values are contiguous
// at the end of the slice.
//
// This ensures the following:
//
// - Last-value-wins semantics
// - Caller sees the reordering, but doesn't lose values
// - Repeated call preserve last-value wins.
//
// Note that methods are defined on Set, although this returns Set. Callers
// can avoid memory allocations by:
//
// - allocating a Sortable for use as a temporary in this method
// - allocating a Set for storing the return value of this constructor.
//
// The result maintains a cache of encoded attributes, by attribute.EncoderID.
// This value should not be copied after its first use.
//
// The second []KeyValue return value is a list of attributes that were
// excluded by the Filter (if non-nil).
func NewSetWithSortableFiltered(kvs []KeyValue, tmp *Sortable, filter Filter) (Set, []KeyValue) {
	// Check for empty set.
	if len(kvs) == 0 {
		return empty(), nil
	}

	data := getSlice(len(kvs), len(kvs))
	copy(*data, kvs)

	*tmp = *data
	// Stable sort so the following de-duplication can implement
	// last-value-wins semantics.
	sort.Stable(tmp)
	*tmp = nil

	position := len(*data) - 1
	offset := position - 1

	// The requirements stated above require that the stable
	// result be placed in the end of the input slice, while
	// overwritten values are swapped to the beginning.
	//
	// De-duplicate with last-value-wins semantics.  Preserve
	// duplicate values at the beginning of the input slice.
	for ; offset >= 0; offset-- {
		if (*data)[offset].Key == (*data)[position].Key {
			continue
		}
		position--
		(*data)[offset], (*data)[position] = (*data)[position], (*data)[offset]
	}

	var (
		set     Set
		dropped []KeyValue
	)
	*data = (*data)[position:]
	if filter != nil {
		set, dropped = filterSet(data, filter)
	} else {
		set = newSet(data)
	}

	return set, dropped
}

// filterSet reorders kvs so that included keys are contiguous at the end of
// the slice, while excluded keys precede the included keys.
func filterSet(kvs *[]KeyValue, filter Filter) (Set, []KeyValue) {
	var excluded []KeyValue

	// Move attributes that do not match the filter so they're adjacent before
	// calling computeDistinct().
	distinctPosition := len(*kvs)

	// Swap indistinct keys forward and distinct keys toward the
	// end of the slice.
	offset := len(*kvs) - 1
	for ; offset >= 0; offset-- {
		if filter((*kvs)[offset]) {
			distinctPosition--
			(*kvs)[offset], (*kvs)[distinctPosition] = (*kvs)[distinctPosition], (*kvs)[offset]
			continue
		}
	}
	excluded = (*kvs)[:distinctPosition]

	*kvs = (*kvs)[distinctPosition:]
	set := newSet(kvs)
	return set, excluded
}

// Filter returns a filtered copy of this Set. See the documentation for
// NewSetWithSortableFiltered for more details.
func (l *Set) Filter(re Filter) (Set, []KeyValue) {
	if re == nil {
		return *l, nil
	}

	// Note: This could be refactored to avoid the temporary slice
	// allocation, if it proves to be expensive.
	return filterSet(l.toSlice(), re)
}

// MarshalJSON returns the JSON encoding of the Set.
func (l *Set) MarshalJSON() ([]byte, error) {
	if l == nil || l.id == nil {
		return nil, nil
	}
	v, _ := sets.Load(*l.id)
	return json.Marshal(v)
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (l Set) MarshalLog() interface{} {
	kvs := make(map[string]string)
	s := l.toSlice()
	for _, kv := range *s {
		kvs[string(kv.Key)] = kv.Value.Emit()
	}
	slicePool.Put(s)
	return kvs
}

// Len implements sort.Interface.
func (l *Sortable) Len() int {
	return len(*l)
}

// Swap implements sort.Interface.
func (l *Sortable) Swap(i, j int) {
	(*l)[i], (*l)[j] = (*l)[j], (*l)[i]
}

// Less implements sort.Interface.
func (l *Sortable) Less(i, j int) bool {
	return (*l)[i].Key < (*l)[j].Key
}
