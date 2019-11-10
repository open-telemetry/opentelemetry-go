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

type HopLimit int

const (
	NoPropagation HopLimit = iota
	UnlimitedPropagation
)

type CorrelationValue struct {
	core.Value
	HopLimit HopLimit
}

type Correlation struct {
	core.KeyValue
	HopLimit HopLimit
}

func (c *Correlation) CorrelationValue() CorrelationValue {
	return CorrelationValue{
		Value:    c.Value,
		HopLimit: c.HopLimit,
	}
}

type Correlations struct {
	m map[core.Key]CorrelationValue
}

type CorrelationsUpdate struct {
	SingleKV Correlation
	MultiKV  []Correlation
	Map      Correlations
}

func NewCorrelations() Correlations {
	return Correlations{
		m: nil,
	}
}

func (m Correlations) Apply(update CorrelationsUpdate) Correlations {
	r := make(map[core.Key]CorrelationValue, len(m.m)+len(update.MultiKV)+update.Map.Len())
	for k, v := range m.m {
		r[k] = v
	}
	if update.SingleKV.Key.Defined() {
		r[update.SingleKV.Key] = update.SingleKV.CorrelationValue()
	}
	for _, kv := range update.MultiKV {
		r[kv.Key] = kv.CorrelationValue()
	}
	for k, v := range update.Map.m {
		r[k] = v
	}
	if len(r) == 0 {
		r = nil
	}
	return Correlations{
		m: r,
	}
}

func (m Correlations) Value(k core.Key) (CorrelationValue, bool) {
	v, ok := m.m[k]
	return v, ok
}

func (m Correlations) HasValue(k core.Key) bool {
	_, has := m.Value(k)
	return has
}

func (m Correlations) Len() int {
	return len(m.m)
}

func (m Correlations) Foreach(f func(kv Correlation) bool) {
	for k, v := range m.m {
		if !f(Correlation{
			KeyValue: core.KeyValue{
				Key:   k,
				Value: v.Value,
			},
			HopLimit: v.HopLimit,
		}) {
			return
		}
	}
}

func (m Correlations) Correlations() []Correlation {
	a := make([]Correlation, 0, len(m.m))
	for k, v := range m.m {
		a = append(a, Correlation{
			KeyValue: core.KeyValue{
				Key:   k,
				Value: v.Value,
			},
			HopLimit: v.HopLimit,
		})
	}
	return a
}
