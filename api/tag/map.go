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

package tag

import (
	"go.opentelemetry.io/api/core"
)

type MeasureMetadata struct {
	TTL int // -1 == infinite, 0 == do not propagate
}

type tagContent struct {
	value core.Value
	meta  MeasureMetadata
}

type rawMap map[core.Key]tagContent

type Map struct {
	m rawMap
}

type MapUpdate struct {
	SingleKV      core.KeyValue
	MultiKV       []core.KeyValue
	SingleMutator Mutator
	MultiMutator  []Mutator
}

func newMap(raw rawMap) Map {
	return Map{
		m: raw,
	}
}

func NewEmptyMap() Map {
	return newMap(nil)
}

func NewMap(update MapUpdate) Map {
	return NewEmptyMap().Apply(update)
}

func (m Map) Apply(update MapUpdate) Map {
	r := make(rawMap, len(m.m)+len(update.MultiKV)+len(update.MultiMutator))
	for k, v := range m.m {
		r[k] = v
	}
	if update.SingleKV.Key.Defined() {
		r[update.SingleKV.Key] = tagContent{
			value: update.SingleKV.Value,
		}
	}
	for _, kv := range update.MultiKV {
		r[kv.Key] = tagContent{
			value: kv.Value,
		}
	}
	if update.SingleMutator.Key.Defined() {
		r.apply(update.SingleMutator)
	}
	for _, mutator := range update.MultiMutator {
		r.apply(mutator)
	}
	if len(r) == 0 {
		r = nil
	}
	return newMap(r)
}

func (m Map) Value(k core.Key) (core.Value, bool) {
	entry, ok := m.m[k]
	if !ok {
		entry.value.Type = core.INVALID
	}
	return entry.value, ok
}

func (m Map) HasValue(k core.Key) bool {
	_, has := m.Value(k)
	return has
}

func (m Map) Len() int {
	return len(m.m)
}

func (m Map) Foreach(f func(kv core.KeyValue) bool) {
	for k, v := range m.m {
		if !f(core.KeyValue{
			Key:   k,
			Value: v.value,
		}) {
			return
		}
	}
}

func (r rawMap) apply(mutator Mutator) {
	key := mutator.KeyValue.Key
	content := tagContent{
		value: mutator.KeyValue.Value,
		meta:  mutator.MeasureMetadata,
	}
	switch mutator.MutatorOp {
	case INSERT:
		if _, ok := r[key]; !ok {
			r[key] = content
		}
	case UPDATE:
		if _, ok := r[key]; ok {
			r[key] = content
		}
	case UPSERT:
		r[key] = content
	case DELETE:
		delete(r, key)
	}
}
