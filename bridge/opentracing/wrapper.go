// Copyright The OpenTelemetry Authors
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

package opentracing // import "go.opentelemetry.io/otel/bridge/opentracing"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/bridge/opentracing/migration"
	"go.opentelemetry.io/otel/trace"
)

// WrapperTracerProvider is an OpenTelemetry TracerProvider that wraps an
// OpenTracing Tracer.
type WrapperTracerProvider struct {
	getWrappedTracer func(name string, opts ...trace.TracerOption) *WrapperTracer
}

var _ trace.TracerProvider = (*WrapperTracerProvider)(nil)

// Tracer creates a WrapperTracer associated with the WrapperTracerProvider.
func (p *WrapperTracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return p.getWrappedTracer(name, opts...)
}

// Deprecated: Use NewTracerProvider(...) instead.
//
// NewWrappedTracerProvider creates a new trace provider that creates a single
// instance of WrapperTracer that wraps OpenTelemetry tracer, and always returns
// it unmodified from Tracer().
func NewWrappedTracerProvider(bridge *BridgeTracer, tracer trace.Tracer) *WrapperTracerProvider {
	wTracer := NewWrapperTracer(bridge, tracer)
	return &WrapperTracerProvider{
		getWrappedTracer: func(_ string, _ ...trace.TracerOption) *WrapperTracer {
			return wTracer
		},
	}
}

type wrappedTracerKey struct {
	name    string
	version string
}

// NewTracerProvider creates a new trace provider that creates new instances of
// WrapperTracer that wraps OpenTelemetry tracer for each call to Tracer().
func NewTracerProvider(bridge *BridgeTracer, provider trace.TracerProvider) *WrapperTracerProvider {
	var (
		mtx     sync.Mutex
		tracers = make(map[wrappedTracerKey]*WrapperTracer)
	)

	return &WrapperTracerProvider{
		getWrappedTracer: func(name string, opts ...trace.TracerOption) *WrapperTracer {
			mtx.Lock()
			defer mtx.Unlock()

			c := trace.NewTracerConfig(opts...)
			key := wrappedTracerKey{
				name:    name,
				version: c.InstrumentationVersion(),
			}

			if t, ok := tracers[key]; ok {
				return t
			}

			wrapper := NewWrapperTracer(bridge, provider.Tracer(name, opts...))
			tracers[key] = wrapper
			return wrapper
		},
	}
}

// WrapperTracer is a wrapper around an OpenTelemetry tracer. It
// mostly forwards the calls to the wrapped tracer, but also does some
// extra steps like setting up a context with the active OpenTracing
// span.
//
// It does not need to be used when the OpenTelemetry tracer is also
// aware how to operate in environment where OpenTracing API is also
// used.
type WrapperTracer struct {
	bridge *BridgeTracer
	tracer trace.Tracer
}

var _ trace.Tracer = &WrapperTracer{}
var _ migration.DeferredContextSetupTracerExtension = &WrapperTracer{}

// NewWrapperTracer wraps the passed tracer and also talks to the
// passed bridge tracer when setting up the context with the new
// active OpenTracing span.
func NewWrapperTracer(bridge *BridgeTracer, tracer trace.Tracer) *WrapperTracer {
	return &WrapperTracer{
		bridge: bridge,
		tracer: tracer,
	}
}

func (t *WrapperTracer) otelTracer() trace.Tracer {
	return t.tracer
}

// Start forwards the call to the wrapped tracer. It also tries to
// override the tracer of the returned span if the span implements the
// OverrideTracerSpanExtension interface.
func (t *WrapperTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := t.otelTracer().Start(ctx, name, opts...)
	if spanWithExtension, ok := span.(migration.OverrideTracerSpanExtension); ok {
		spanWithExtension.OverrideTracer(t)
	}
	if !migration.SkipContextSetup(ctx) {
		ctx = t.bridge.ContextWithBridgeSpan(ctx, span)
	}
	return ctx, span
}

// DeferredContextSetupHook is a part of the implementation of the
// DeferredContextSetupTracerExtension interface. It will try to
// forward the call to the wrapped tracer if it implements the
// interface.
func (t *WrapperTracer) DeferredContextSetupHook(ctx context.Context, span trace.Span) context.Context {
	if tracerWithExtension, ok := t.otelTracer().(migration.DeferredContextSetupTracerExtension); ok {
		ctx = tracerWithExtension.DeferredContextSetupHook(ctx, span)
	}
	ctx = trace.ContextWithSpan(ctx, span)
	return ctx
}
