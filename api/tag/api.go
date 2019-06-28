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

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/unit"
)

type Map interface {
	// TODO combine these four into a struct
	Apply(a1 core.KeyValue, attributes []core.KeyValue, m1 core.Mutator, mutators []core.Mutator) Map

	Value(core.Key) (core.Value, bool)
	HasValue(core.Key) bool

	Len() int

	Foreach(func(kv core.KeyValue) bool)
}

type Option func(*registeredKey)

var (
	EmptyMap = NewMap(core.KeyValue{}, nil, core.Mutator{}, nil)
)

func New(name string, opts ...Option) core.Key { // TODO rename NewKey?
	return register(name, opts)
}

func NewMeasure(name string, opts ...Option) core.Measure {
	return measure{
		rk: register(name, opts),
	}
}

func NewMap(a1 core.KeyValue, attributes []core.KeyValue, m1 core.Mutator, mutators []core.Mutator) Map {
	var t tagMap
	return t.Apply(a1, attributes, m1, mutators)
}

func WithMap(ctx context.Context, m Map) context.Context {
	return context.WithValue(ctx, ctxTagsKey, m)
}

func NewContext(ctx context.Context, mutators ...core.Mutator) context.Context {
	return WithMap(ctx, FromContext(ctx).Apply(
		core.KeyValue{}, nil,
		core.Mutator{}, mutators,
	))
}

func FromContext(ctx context.Context) Map {
	if m, ok := ctx.Value(ctxTagsKey).(Map); ok {
		return m
	}
	return tagMap{}
}

// WithDescription applies provided description.
func WithDescription(desc string) Option {
	return func(rk *registeredKey) {
		rk.desc = desc
	}
}

// WithUnit applies provided unit.
func WithUnit(unit unit.Unit) Option {
	return func(rk *registeredKey) {
		rk.unit = unit
	}
}
