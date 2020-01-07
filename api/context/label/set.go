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

type Set struct {
	set *internal.Set
}

func Empty() Set {
	return Set{internal.EmptySet()}
}

func NewSet(kvs ...core.KeyValue) Set {
	return Set{internal.NewSet(kvs...)}
}

func (s Set) AddOne(kv core.KeyValue) Set {
	return Set{s.set.AddOne(kv)}

}

func (s Set) AddMany(kvs ...core.KeyValue) Set {
	return Set{s.set.AddMany(kvs...)}
}

func (s Set) Value(k core.Key) (core.Value, bool) {
	return s.set.Value(k)
}

func (s Set) HasValue(k core.Key) bool {
	return s.set.HasValue(k)
}

func (s Set) Len() int {
	return s.set.Len()
}

func (s Set) Ordered() []core.KeyValue {
	return s.set.Ordered()
}

func (s Set) Foreach(f func(kv core.KeyValue) bool) {
	for _, kv := range s.set.Ordered() {
		if !f(kv) {
			return
		}
	}
}

func (s Set) Equals(t Set) bool {
	return s.set.Equals(t.set)
}

func (s Set) Encoded(enc core.LabelEncoder) string {
	return s.set.Encoded(enc)
}
