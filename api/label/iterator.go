package label

import (
	"go.opentelemetry.io/otel/api/core"
)

// Next moves the iterator to the next label. Returns false if there
// are no more labels.
func (i *Iterator) Next() bool {
	i.idx++
	return i.idx < i.Len()
}

// Label returns current label. Must be called only after Next returns
// true.
func (i *Iterator) Label() core.KeyValue {
	kv, _ := i.storage.Get(i.idx)
	return kv
}

// IndexedLabel returns current index and label. Must be called only
// after Next returns true.
func (i *Iterator) IndexedLabel() (int, core.KeyValue) {
	return i.idx, i.Label()
}

// Len returns a number of labels in the iterator's label storage.
func (i *Iterator) Len() int {
	return i.storage.Len()
}

// ToSlice is a convenience function that creates a slice of labels
// from the passed iterator. The iterator is set up to start from the
// beginning before creating the slice.
func (i *Iterator) ToSlice() []core.KeyValue {
	l := i.Len()
	if l == 0 {
		return nil
	}
	i.idx = -1
	slice := make([]core.KeyValue, 0, l)
	for i.Next() {
		slice = append(slice, i.Label())
	}
	return slice
}
