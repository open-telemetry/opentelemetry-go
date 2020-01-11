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

// InContext returns a context with this Scope current.  The current
// Scope's resources will be implicitly associated with metric events
// that happen in the returned context.
//
// When using a Scope's Tracer() or Meter() handle for an API method
// call, the Scope is automatically applied, making it the current
// Scope in the context of the resulting call.
func (s Scope) InContext(ctx context.Context) context.Context {
	return internal.SetScopeImpl(ctx, s)
}

// Current returns the Scope associated with a Context as set by
// Scope.InContext().  If no Scope is located in the context, the
// global scope will be returned.
func Current(ctx context.Context) Scope {
	impl := internal.ScopeImpl(ctx)
	if impl == nil {
		if sc, ok := (*atomic.Value)(atomic.LoadPointer(&internal.GlobalScope)).Load().(Scope); ok {
			return sc
		}
		// If if the global not a Scope, it means the global
		// package was not a dependency, only api/global/internal
		// sets this.
		return Scope{}
	}
	return impl.(Scope)
}

// Labels is a convenience method to return a LabelSet given the
// Context and any additional labels.  `Labels(ctx)` returns the
// current resources.
func Labels(ctx context.Context, labels ...core.KeyValue) label.Set {
	return Current(ctx).AddResources(labels...).Resources()
}

// WithTracerSDK returns a new Scope with just a Tracer attached.
func WithTracerSDK(ti trace.TracerSDK) Scope {
	return Scope{}.WithTracerSDK(ti)
}

// WithMeterSDK returns a new Scope with just a Meter attached.
func WithMeterSDK(ti metric.MeterSDK) Scope {
	return Scope{}.WithMeterSDK(ti)
}

// UnnamedTracer returns a Tracer implementation with an empty namespace,
// as a convenience.
func UnnamedTracer(ti trace.TracerSDK) trace.Tracer {
	return WithTracerSDK(ti).Tracer()
}

// UnnamedMeter returns a Meter implementation with an empty namespace,
// as a convenience.
func UnnamedMeter(ti metric.MeterSDK) metric.Meter {
	return WithMeterSDK(ti).Meter()
}

// NamedTracer returns a Tracer implementation with the specified
// namespace, as a convenience.
func NamedTracer(ti trace.TracerSDK, ns core.Namespace) trace.Tracer {
	return WithTracerSDK(ti).WithNamespace(ns).Tracer()
}

// NamedMeter returns a Tracer implementation with the specified
// namespace, as a convenience.
func NamedMeter(ti metric.MeterSDK, ns core.Namespace) metric.Meter {
	return WithMeterSDK(ti).WithNamespace(ns).Meter()
}
