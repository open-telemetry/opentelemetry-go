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
