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

	"go.opentelemetry.io/api/core"
)

type tagContent struct {
	value core.Value
	meta  MeasureMetadata
}

type tagMap map[core.Key]tagContent

var _ Map = tagMap{}

func (t tagMap) Apply(update MapUpdate) Map {
	m := make(tagMap, len(t)+len(update.MultiKV)+len(update.MultiMutator))
	for k, v := range t {
		m[k] = v
	}
	if update.SingleKV.Key.Defined() {
		m[update.SingleKV.Key] = tagContent{
			value: update.SingleKV.Value,
		}
	}
	for _, kv := range update.MultiKV {
		m[kv.Key] = tagContent{
			value: kv.Value,
		}
	}
	if update.SingleMutator.Key.Defined() {
		m.apply(update.SingleMutator)
	}
	for _, mutator := range update.MultiMutator {
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

func (m tagMap) apply(mutator Mutator) {
	if m == nil {
		return
	}
	key := mutator.KeyValue.Key
	content := tagContent{
		value: mutator.KeyValue.Value,
		meta:  mutator.MeasureMetadata,
	}
	switch mutator.MutatorOp {
	case INSERT:
		if _, ok := m[key]; !ok {
			m[key] = content
		}
	case UPDATE:
		if _, ok := m[key]; ok {
			m[key] = content
		}
	case UPSERT:
		m[key] = content
	case DELETE:
		delete(m, key)
	}
}

func Insert(kv core.KeyValue) Mutator {
	return Mutator{
		MutatorOp: INSERT,
		KeyValue:  kv,
	}
}

func Update(kv core.KeyValue) Mutator {
	return Mutator{
		MutatorOp: UPDATE,
		KeyValue:  kv,
	}
}

func Upsert(kv core.KeyValue) Mutator {
	return Mutator{
		MutatorOp: UPSERT,
		KeyValue:  kv,
	}
}

func Delete(k core.Key) Mutator {
	return Mutator{
		MutatorOp: DELETE,
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
		keyvals = append(keyvals, k.Variable.Name, v.value.Emit())
	}
	pprof.Do(ctx, pprof.Labels(keyvals...), f)
}
