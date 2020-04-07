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
