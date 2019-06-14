package scope

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

type (
	Scope interface {
		ScopeID() core.ScopeID
	}

	Mutable interface {
		Scope

		SetAttribute(core.KeyValue)
		SetAttributes(...core.KeyValue)

		ModifyAttribute(core.Mutator)
		ModifyAttributes(...core.Mutator)
	}

	scopeIdent struct {
		id core.ScopeID
	}

	scopeKeyType struct{}
)

var (
	scopeKey   = &scopeKeyType{}
	emptyScope = &scopeIdent{}
)

func SetActive(ctx context.Context, scope Scope) context.Context {
	return context.WithValue(ctx, scopeKey, scope)
}

func Active(ctx context.Context) Scope {
	if scope, has := ctx.Value(scopeKey).(Scope); has {
		return scope
	}
	return emptyScope
}

func (s *scopeIdent) ScopeID() core.ScopeID {
	if s == nil {
		return core.ScopeID{}
	}
	return s.id
}

func New(parent core.ScopeID, attributes ...core.KeyValue) Scope {
	eventID := observer.Record(observer.Event{
		Type:       observer.NEW_SCOPE,
		Scope:      parent,
		Attributes: attributes,
	})
	return &scopeIdent{
		id: core.ScopeID{
			EventID:     eventID,
			SpanContext: parent.SpanContext,
		},
	}
}
