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

package correlation

import (
	"go.opentelemetry.io/otel/api/core"
)

// TODO Comments needed! This was formerly known as distributedcontext.Map

type rawMap map[core.Key]core.Value
type keySet map[core.Key]struct{}

type Map struct {
	m rawMap
}

type MapUpdate struct {
	DropSingleK core.Key
	DropMultiK  []core.Key

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
	delSet, addSet := getModificationSets(update)

	mapSizeDiff := 0
	for k := range addSet {
		if _, ok := m.m[k]; !ok {
			mapSizeDiff++
		}
	}
	for k := range delSet {
		if _, ok := m.m[k]; ok {
			if _, inAddSet := addSet[k]; !inAddSet {
				mapSizeDiff--
			}
		}
	}

	r := make(rawMap, len(m.m)+mapSizeDiff)
	for k, v := range m.m {
		// do not copy items we want to drop
		if _, ok := delSet[k]; ok {
			continue
		}
		// do not copy items we would overwrite
		if _, ok := addSet[k]; ok {
			continue
		}
		r[k] = v
	}
	if update.SingleKV.Key.Defined() {
		r[update.SingleKV.Key] = update.SingleKV.Value
	}
	for _, kv := range update.MultiKV {
		r[kv.Key] = kv.Value
	}
	if len(r) == 0 {
		r = nil
	}
	return newMap(r)
}

func getModificationSets(update MapUpdate) (keySet, keySet) {
	deletionsCount := len(update.DropMultiK)
	if update.DropSingleK.Defined() {
		deletionsCount++
	}
	var delSet keySet
	if deletionsCount > 0 {
		delSet = make(map[core.Key]struct{}, deletionsCount)
		for _, k := range update.DropMultiK {
			delSet[k] = struct{}{}
		}
		if update.DropSingleK.Defined() {
			delSet[update.DropSingleK] = struct{}{}
		}
	}

	additionsCount := len(update.MultiKV)
	if update.SingleKV.Key.Defined() {
		additionsCount++
	}
	var addSet keySet
	if additionsCount > 0 {
		addSet = make(map[core.Key]struct{}, additionsCount)
		for _, k := range update.MultiKV {
			addSet[k.Key] = struct{}{}
		}
		if update.SingleKV.Key.Defined() {
			addSet[update.SingleKV.Key] = struct{}{}
		}
	}

	return delSet, addSet
}

func (m Map) Value(k core.Key) (core.Value, bool) {
	value, ok := m.m[k]
	return value, ok
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
			Value: v,
		}) {
			return
		}
	}
}
