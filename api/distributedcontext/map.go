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

package distributedcontext

import (
	"go.opentelemetry.io/otel/api/core"
)

type entry struct {
	value core.Value
}

type rawMap map[core.Key]entry

type Map struct {
	m rawMap
}

type MapUpdate struct {
	SingleKV core.KeyValue
	MultiKV  []core.KeyValue
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
	r := make(rawMap, len(m.m)+len(update.MultiKV))
	for k, v := range m.m {
		r[k] = v
	}
	if update.SingleKV.Key.Defined() {
		r[update.SingleKV.Key] = entry{
			value: update.SingleKV.Value,
		}
	}
	for _, kv := range update.MultiKV {
		r[kv.Key] = entry{
			value: kv.Value,
		}
	}
	if len(r) == 0 {
		r = nil
	}
	return newMap(r)
}

func (m Map) Value(k core.Key) (core.Value, bool) {
	entry, ok := m.m[k]
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
