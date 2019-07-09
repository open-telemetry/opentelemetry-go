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
)

type ctxTagsType struct{}

var (
	ctxTagsKey = &ctxTagsType{}
)

type MutatorOp int

const (
	INSERT MutatorOp = iota
	UPDATE
	UPSERT
	DELETE
)

type Mutator struct {
	MutatorOp
	core.KeyValue
	MeasureMetadata
}

type MeasureMetadata struct {
	TTL int // -1 == infinite, 0 == do not propagate
}

func (m Mutator) WithTTL(hops int) Mutator {
	m.TTL = hops
	return m
}

type MapUpdate struct {
	SingleKV      core.KeyValue
	MultiKV       []core.KeyValue
	SingleMutator Mutator
	MultiMutator  []Mutator
}

type Map interface {
	Apply(MapUpdate) Map

	Value(core.Key) (core.Value, bool)
	HasValue(core.Key) bool

	Len() int

	Foreach(func(kv core.KeyValue) bool)
}

func NewEmptyMap() Map {
	return tagMap{}
}

func NewMap(update MapUpdate) Map {
	return NewEmptyMap().Apply(update)
}

func WithMap(ctx context.Context, m Map) context.Context {
	return context.WithValue(ctx, ctxTagsKey, m)
}

func NewContext(ctx context.Context, mutators ...Mutator) context.Context {
	return WithMap(ctx, FromContext(ctx).Apply(MapUpdate{
		MultiMutator: mutators,
	}))
}

func FromContext(ctx context.Context) Map {
	if m, ok := ctx.Value(ctxTagsKey).(Map); ok {
		return m
	}
	return tagMap{}
}
