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
	str    string
	sorted []core.KeyValue
	keys   map[core.Key]struct{}
}

// New creates a resource from a set of attributes.
// If there are duplicates keys then the first value of the key is preserved.
func New(kvs ...core.KeyValue) *Resource {
	res := &Resource{keys: make(map[core.Key]struct{})}
	for _, kv := range kvs {
		// First key wins.
		if _, ok := res.keys[kv.Key]; !ok {
			res.keys[kv.Key] = struct{}{}
			res.sorted = append(res.sorted, kv)
		}
	}
	sort.Slice(res.sorted, func(i, j int) bool {
		return res.sorted[i].Key < res.sorted[j].Key
	})
	res.str = buildResourceString(res.sorted)
	return res
}

// String implements the Stringer interface and provides a reproducibly
// hashable representation of a Resource.
func (r Resource) String() string {
	return r.str
}

// Attributes returns a copy of attributes from the resource in a sorted order.
func (r Resource) Attributes() []core.KeyValue {
	return append([]core.KeyValue(nil), r.sorted...)
}

// Equal returns true if other Resource is equal to r.
func (r Resource) Equal(other Resource) bool {
	return r.str == other.str
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

	n := len(a.sorted)
	if len(b.sorted) > len(a.sorted) {
		n = len(b.sorted)
	}
	// At a minimum the merge will be as large as the largest resource.
	s := make([]core.KeyValue, 0, n)
	k := make(map[core.Key]struct{}, n)
	ai, bi := 0, 0
	for ; ai < len(a.sorted) && bi < len(b.sorted); ai, bi = ai+1, bi+1 {
		akv := a.sorted[ai]
		s = append(s, akv)
		k[akv.Key] = struct{}{}

		bkv := b.sorted[bi]
		if _, ok := a.keys[bkv.Key]; ok {
			// a overwrites b.
			continue
		}
		k[bkv.Key] = struct{}{}
		s = append(s, bkv)
	}

	for ; ai < len(a.sorted); ai++ {
		akv := a.sorted[ai]
		s = append(s, akv)
		k[akv.Key] = struct{}{}
	}

	for ; bi < len(b.sorted); bi++ {
		bkv := b.sorted[bi]
		if _, ok := k[bkv.Key]; ok {
			continue
		}
		k[bkv.Key] = struct{}{}
		s = append(s, bkv)
	}

	sort.Slice(s, func(i, j int) bool { return s[i].Key < s[j].Key })
	res := &Resource{
		keys:   k,
		sorted: s,
	}
	res.str = buildResourceString(res.sorted)

	return res
}

// buildResourceString returns a string representation of a Resource
// containing kvs.
func buildResourceString(kvs []core.KeyValue) string {
	// Ensure unique strings if key/value contains '=', ',', or '\'.
	escaper := strings.NewReplacer("=", `\=`, ",", `\,`, `\`, `\\`)

	var b strings.Builder
	b.WriteString("Resource(")
	if len(kvs) > 0 {
		b.WriteString(escaper.Replace(string(kvs[0].Key)))
		b.WriteRune('=')
		b.WriteString(escaper.Replace(kvs[0].Value.Emit()))
		for _, s := range kvs[1:] {
			b.WriteRune(',')
			b.WriteString(escaper.Replace(string(s.Key)))
			b.WriteRune('=')
			b.WriteString(escaper.Replace(s.Value.Emit()))
		}

	}
	b.WriteRune(')')
	return b.String()
}
