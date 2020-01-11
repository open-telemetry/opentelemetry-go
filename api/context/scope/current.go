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
	"sync/atomic"

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/internal"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

// ContextWithScope returns a context with a new current Scope.  The
// active Scope's resources will be implicitly associated with metric
// events that happen in the returned context.
//
// When using a Scope's Tracer() or Meter() handle for an API method
// call, the Scope is automatically applied, making it the current
// Scope for the resulting call.
func ContextWithScope(ctx context.Context, sc Scope) context.Context {
	return internal.SetScopeImpl(ctx, sc.scopeImpl)
}

// Current returns the Scope associated with a Context as set by
// ContextWithScope.  If no Scope is located in the context, the
// global scope will be used.
func Current(ctx context.Context) Scope {
	impl := internal.ScopeImpl(ctx)
	if impl == nil {
		// If if the global not a Scope, it means the global package was not loaded
		if sc, ok := (*atomic.Value)(atomic.LoadPointer(&internal.GlobalScope)).Load().(Scope); ok {
			return sc
		}
		return Scope{}
	}
	return Scope{internal.ScopeImpl(ctx).(*scopeImpl)}
}

// Labels is a convenience method to return a LabelSet given the
// Context and any additional labels.  `Labels(ctx)` returns the
// current resources.
func Labels(ctx context.Context, labels ...core.KeyValue) label.Set {
	return Current(ctx).AddResources(labels...).Resources()
}

// WithTracer returns a new Scope with just a Tracer attached.
func WithTracer(ti trace.TracerSDK) Scope {
	return Scope{}.WithTracer(ti)
}

// UnnamedTracer returns a Tracer implementation with an empty namespace,
// as a convenience.
func UnnamedTracer(ti trace.TracerSDK) trace.Tracer {
	return WithTracer(ti).Tracer()
}

// NamedTracer returns a Tracer implementation with the specified
// namespace, as a convenience.
func NamedTracer(ti trace.TracerSDK, ns core.Namespace) trace.Tracer {
	return WithTracer(ti).WithNamespace(ns).Tracer()
}

// WithMeter returns a new Scope with just a Meter attached.
func WithMeter(ti metric.MeterSDK) Scope {
	return Scope{}.WithMeter(ti)
}

// UnnamedMeter returns a Meter implementation with an empty namespace,
// as a convenience.
func UnnamedMeter(ti metric.MeterSDK) metric.Meter {
	return WithMeter(ti).Meter()
}

// NamedMeter returns a Tracer implementation with the specified
// namespace, as a convenience.
func NamedMeter(ti metric.MeterSDK, ns core.Namespace) metric.Meter {
	return WithMeter(ti).WithNamespace(ns).Meter()
}

// InContext returns a context for this scope.  Uses of the global
// Tracer methods (e.g., trace.Start) and Meter methods (e.g.,
// metric.NewInt64Counter) will be constructed using the namespace
// from this scope.
//
// Uses of the global meter.RecordBatch will use the resources of this
// scope from the context.
func (s Scope) InContext(ctx context.Context) context.Context {
	return ContextWithScope(ctx, s)
}
