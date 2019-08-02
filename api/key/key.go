package key

import (
	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/registry"
)

type AnyValue struct{}

func (AnyValue) String() string {
	return "AnyValue"
}

func New(name string, opts ...registry.Option) core.Key {
	return core.Key{
		Variable: registry.Register(name, AnyValue{}, opts...),
	}
}

var (
	WithDescription = registry.WithDescription
	WithUnit        = registry.WithUnit
)
