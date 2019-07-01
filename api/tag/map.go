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
	"context"
	"runtime/pprof"

	"github.com/open-telemetry/opentelemetry-go/api/core"
)

type tagContent struct {
	value core.Value
	meta  core.MeasureMetadata
}

type tagMap map[core.Key]tagContent

var _ Map = (*tagMap)(nil)

func (t tagMap) Apply(a1 core.KeyValue, attributes []core.KeyValue, m1 core.Mutator, mutators []core.Mutator) Map {
	m := make(tagMap, len(t)+len(attributes)+len(mutators))
	for k, v := range t {
		m[k] = v
	}
	if a1.Key != nil {
		m[a1.Key] = tagContent{
			value: a1.Value,
		}
	}
	for _, kv := range attributes {
		m[kv.Key] = tagContent{
			value: kv.Value,
		}
	}
	if m1.KeyValue.Key != nil {
		m.apply(m1)
	}
	for _, mutator := range mutators {
		m.apply(mutator)
	}
	return m
}

func (m tagMap) Value(k core.Key) (core.Value, bool) {
	entry, ok := m[k]
	if !ok {
		entry.value.Type = core.INVALID
	}
	return entry.value, ok
}

func (m tagMap) HasValue(k core.Key) bool {
	_, has := m.Value(k)
	return has
}

func (m tagMap) Len() int {
	return len(m)
}

func (m tagMap) Foreach(f func(kv core.KeyValue) bool) {
	for k, v := range m {
		if !f(core.KeyValue{
			Key:   k,
			Value: v.value,
		}) {
			return
		}
	}
}

func (m tagMap) apply(mutator core.Mutator) {
	if m == nil {
		return
	}
	key := mutator.KeyValue.Key
	content := tagContent{
		value: mutator.KeyValue.Value,
		meta:  mutator.MeasureMetadata,
	}
	switch mutator.MutatorOp {
	case core.INSERT:
		if _, ok := m[key]; !ok {
			m[key] = content
		}
	case core.UPDATE:
		if _, ok := m[key]; ok {
			m[key] = content
		}
	case core.UPSERT:
		m[key] = content
	case core.DELETE:
		delete(m, key)
	}
}

func Insert(kv core.KeyValue) core.Mutator {
	return core.Mutator{
		MutatorOp: core.INSERT,
		KeyValue:  kv,
	}
}

func Update(kv core.KeyValue) core.Mutator {
	return core.Mutator{
		MutatorOp: core.UPDATE,
		KeyValue:  kv,
	}
}

func Upsert(kv core.KeyValue) core.Mutator {
	return core.Mutator{
		MutatorOp: core.UPSERT,
		KeyValue:  kv,
	}
}

func Delete(k core.Key) core.Mutator {
	return core.Mutator{
		MutatorOp: core.DELETE,
		KeyValue: core.KeyValue{
			Key: k,
		},
	}
}

// Note: the golang pprof.Do API forces this memory allocation, we
// should file an issue about that.  (There's a TODO in the source.)
func Do(ctx context.Context, f func(ctx context.Context)) {
	m := FromContext(ctx).(tagMap)
	keyvals := make([]string, 0, 2*len(m))
	for k, v := range m {
		keyvals = append(keyvals, k.Name(), v.value.Emit())
	}
	pprof.Do(ctx, pprof.Labels(keyvals...), f)
}
