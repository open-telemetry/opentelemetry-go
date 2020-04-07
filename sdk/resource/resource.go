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

// Package resource provides functionality for resource, which capture
// identifying information about the entities for which signals are exported.
package resource

import (
	"sort"
	"strings"

	"go.opentelemetry.io/otel/api/core"
)

// Resource describes an entity about which identifying information and metadata is exposed.
type Resource struct {
	sorted []core.KeyValue
	keySet map[core.Key]struct{}
}

// New creates a resource from a set of attributes.
// If there are duplicates keys then the first value of the key is preserved.
func New(kvs ...core.KeyValue) *Resource {
	res := &Resource{keySet: make(map[core.Key]struct{})}
	for _, kv := range kvs {
		// First key wins.
		if _, ok := res.keySet[kv.Key]; !ok {
			res.keySet[kv.Key] = struct{}{}
			res.sorted = append(res.sorted, kv)
		}
	}
	sort.Slice(res.sorted, func(i, j int) bool {
		return res.sorted[i].Key < res.sorted[j].Key
	})
	return res
}

// String implements the Stringer interface and provides a reproducibly
// hashable representation of a Resource.
func (r Resource) String() string {
	// Ensure unique strings if key/value contains '=', ',', or '\'.
	escaper := strings.NewReplacer("=", `\=`, ",", `\,`, `\`, `\\`)

	var b strings.Builder
	// Note: this could be further optimized by precomputing the size of
	// the resulting buffer and adding a call to b.Grow
	b.WriteString("Resource(")
	if len(r.sorted) > 0 {
		b.WriteString(escaper.Replace(string(r.sorted[0].Key)))
		b.WriteRune('=')
		b.WriteString(escaper.Replace(r.sorted[0].Value.Emit()))
		for _, s := range r.sorted[1:] {
			b.WriteRune(',')
			b.WriteString(escaper.Replace(string(s.Key)))
			b.WriteRune('=')
			b.WriteString(escaper.Replace(s.Value.Emit()))
		}

	}
	b.WriteRune(')')

	return b.String()
}

// Attributes returns a copy of attributes from the resource in a sorted order.
func (r Resource) Attributes() []core.KeyValue {
	return append([]core.KeyValue(nil), r.sorted...)
}

// Iter returns an interator of the Resource attributes.
//
// This is ideal to use if you do not want a copy of the attributes.
func (r Resource) Iter() AttributeIterator {
	return NewAttributeIterator(r.sorted)
}

// Equal returns true if other Resource is equal to r.
func (r Resource) Equal(other Resource) bool {
	return r.String() == other.String()
}

// Merge creates a new resource by combining resource a and b.
// If there are common key between resource a and b then value from resource a is preserved.
// If one of the resources is nil then the other resource is returned without creating a new one.
func Merge(a, b *Resource) *Resource {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}

	// Note: the following could be optimized by implementing a dedicated merge sort.

	kvs := make([]core.KeyValue, 0, len(a.sorted)+len(b.sorted))
	kvs = append(kvs, a.sorted...)
	// a overwrites b, so b needs to be at the end.
	kvs = append(kvs, b.sorted...)
	return New(kvs...)
}
