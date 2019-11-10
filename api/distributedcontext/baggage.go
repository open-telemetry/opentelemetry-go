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

type Baggage struct {
	m map[core.Key]core.Value
}

type BaggageUpdate struct {
	DropSingleK core.Key
	DropMultiK  []core.Key

	SingleKV core.KeyValue
	MultiKV  []core.KeyValue
	Map      Baggage
}

func NewBaggage() Baggage {
	return Baggage{
		m: nil,
	}
}

func (m Baggage) Apply(update BaggageUpdate) Baggage {
	r := make(map[core.Key]core.Value, len(m.m)+len(update.MultiKV)+update.Map.Len())
	for k, v := range m.m {
		r[k] = v
	}
	if update.DropSingleK.Defined() {
		delete(r, update.DropSingleK)
	}
	for _, k := range update.DropMultiK {
		delete(r, k)
	}
	if update.SingleKV.Key.Defined() {
		r[update.SingleKV.Key] = update.SingleKV.Value
	}
	for _, kv := range update.MultiKV {
		r[kv.Key] = kv.Value
	}
	for k, v := range update.Map.m {
		r[k] = v
	}
	if len(r) == 0 {
		r = nil
	}
	return Baggage{
		m: r,
	}
}

func (m Baggage) Value(k core.Key) (core.Value, bool) {
	v, ok := m.m[k]
	return v, ok
}

func (m Baggage) HasValue(k core.Key) bool {
	_, has := m.Value(k)
	return has
}

func (m Baggage) Len() int {
	return len(m.m)
}

func (m Baggage) Foreach(f func(kv core.KeyValue) bool) {
	for k, v := range m.m {
		if !f(core.KeyValue{
			Key:   k,
			Value: v,
		}) {
			return
		}
	}
}

func (m Baggage) KeyValues() []core.KeyValue {
	a := make([]core.KeyValue, 0, len(m.m))
	for k, v := range m.m {
		a = append(a, core.KeyValue{
			Key:   k,
			Value: v,
		})
	}
	return a
}
