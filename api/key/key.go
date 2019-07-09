package key

import (
	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/registry"
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
