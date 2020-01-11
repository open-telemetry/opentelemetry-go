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

package label

import (
	"go.opentelemetry.io/otel/api/context/internal"
	"go.opentelemetry.io/otel/api/core"
)

// Set represents an immutable set of labels, used to represent the
// "resources" of a Scope.  LabelSets contain a unique mapping from
// Key to Value; duplicates are treated by retaining the last value.
//
// Scopes contain these resources, so that end users will rarely need
// to handle these directly.
//
// Set supports caching the encoded representation of the set of
// labels based on a user-supplied LabelEncoder.
type Set struct {
	set *internal.Set
}

// Empty returns a set with zero keys.
func Empty() Set {
	return Set{internal.EmptySet()}
}

// NewSet constructs a set from a list of KeyValues.  Ordinarily users
// will not construct these directly, as the Scope represents the
// current resources as a label set.
func NewSet(kvs ...core.KeyValue) Set {
	return Set{internal.NewSet(kvs...)}
}

// AddOne adds a single KeyValue to the set.
func (s Set) AddOne(kv core.KeyValue) Set {
	return Set{s.set.AddOne(kv)}

}

// AddMany adds multiple KeyValues to the set.
func (s Set) AddMany(kvs ...core.KeyValue) Set {
	return Set{s.set.AddMany(kvs...)}
}

// Value returns the value associated with the supplied Key and a
// boolean to indicate whether it was found.
func (s Set) Value(k core.Key) (core.Value, bool) {
	return s.set.Value(k)
}

// HasValue returns true if the set contains a value associated with a Key.
func (s Set) HasValue(k core.Key) bool {
	return s.set.HasValue(k)
}

// Len returns the number of labels in the set.
func (s Set) Len() int {
	return s.set.Len()
}

// Ordered returns the label set sorted alphanumerically by Key name.
func (s Set) Ordered() []core.KeyValue {
	return s.set.Ordered()
}

// Foreach calls the provided callback for each label in the set.
func (s Set) Foreach(f func(kv core.KeyValue) bool) {
	for _, kv := range s.set.Ordered() {
		if !f(kv) {
			return
		}
	}
}

// Equals tests whether two sets of labels are identical.
func (s Set) Equals(t Set) bool {
	return s.set.Equals(t.set)
}

// Encoded returns the computed encoding for a label set.  Encoded
// values are cached with the set, to avoid recomputing them.
func (s Set) Encoded(enc core.LabelEncoder) string {
	return s.set.Encoded(enc)
}
