package tag

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/unit"
)

type (
	Map interface {
		// TODO combine these four into a struct
		Apply(a1 core.KeyValue, attributes []core.KeyValue, m1 core.Mutator, mutators []core.Mutator) Map

		Value(core.Key) (core.Value, bool)
		HasValue(core.Key) bool

		Len() int

		Foreach(func(kv core.KeyValue) bool)
	}

	Option func(*registeredKey)
)

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

func (input tagMap) Apply(a1 core.KeyValue, attributes []core.KeyValue, m1 core.Mutator, mutators []core.Mutator) Map {
	m := make(tagMap, len(input)+len(attributes)+len(mutators))
	for k, v := range input {
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
