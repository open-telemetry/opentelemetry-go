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

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/internal"
)

func ContextWithScope(ctx context.Context, sc Scope) context.Context {
	return internal.SetScopeImpl(ctx, sc.scopeImpl)
}

func Current(ctx context.Context) Scope {
	impl := internal.ScopeImpl(ctx)
	if impl == nil {
		return Empty()
	}
	return Scope{internal.ScopeImpl(ctx).(*scopeImpl)}
}

func Labels(ctx context.Context, labels ...core.KeyValue) label.Set {
	return Current(ctx).AddResources(labels...).Resources()
}
