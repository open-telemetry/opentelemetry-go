package key

import (
	"go.opentelemetry.io/api/core"
)

func New(name string) core.Key {
	return core.Key{
		Name: name,
	}
}
