package resource

import "go.opentelemetry.io/otel/api/core"

// AttributeIterator allows iterating over an ordered set of Resource attributes.
//
// The typical use of the iterator assuming a Resource named `res`, is
// something like the following:
//
//     for iter := res.Iter(); iter.Next(); {
//       attr := iter.Attribute()
//       // or, if an index is needed:
//       // idx, attr := iter.IndexedAttribute()
//
//       // ...
//     }
type AttributeIterator struct {
	attrs []core.KeyValue
	idx   int
}

// NewAttributeIterator creates an iterator going over a passed attrs.
func NewAttributeIterator(attrs []core.KeyValue) AttributeIterator {
	return AttributeIterator{attrs: attrs, idx: -1}
}

// Next moves the iterator to the next attribute.
// Returns false if there are no more attributes.
func (i *AttributeIterator) Next() bool {
	i.idx++
	return i.idx < i.Len()
}

// Attribute returns current attribute.
//
// Must be called only after Next returns true.
func (i *AttributeIterator) Attribute() core.KeyValue {
	return i.attrs[i.idx]
}

// IndexedAttribute returns current index and attribute.
//
// Must be called only after Next returns true.
func (i *AttributeIterator) IndexedAttribute() (int, core.KeyValue) {
	return i.idx, i.Attribute()
}

// Len returns a number of attributes.
func (i *AttributeIterator) Len() int {
	return len(i.attrs)
}
