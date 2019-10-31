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

package opentracing

import (
	"context"

	oteltrace "go.opentelemetry.io/otel/api/trace"

	migration "go.opentelemetry.io/otel/bridge/opentracing/migration"
)

type WrapperProvider struct {
	wTracer *WrapperTracer
}

var _ oteltrace.Provider = (*WrapperProvider)(nil)

// GetTracer returns the WrapperTracer associated with the WrapperProvider.
func (p *WrapperProvider) GetTracer(name string) oteltrace.Tracer {
	return p.wTracer
}

// WrapperProvider creates a new trace provider that creates a single
// instance of WrapperTracer that wraps OpenTelemetry tracer.
func NewWrappedProvider(bridge *BridgeTracer, tracer oteltrace.Tracer) *WrapperProvider {
	return &WrapperProvider{
		wTracer: NewWrapperTracer(bridge, tracer),
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
	tracer oteltrace.Tracer
}

var _ oteltrace.Tracer = &WrapperTracer{}
var _ migration.DeferredContextSetupTracerExtension = &WrapperTracer{}

// NewWrapperTracer wraps the passed tracer and also talks to the
// passed bridge tracer when setting up the context with the new
// active OpenTracing span.
func NewWrapperTracer(bridge *BridgeTracer, tracer oteltrace.Tracer) *WrapperTracer {
	return &WrapperTracer{
		bridge: bridge,
		tracer: tracer,
	}
}

func (t *WrapperTracer) otelTracer() oteltrace.Tracer {
	return t.tracer
}

// WithSpan forwards the call to the wrapped tracer with a modified
// body callback, which sets the active OpenTracing span before
// calling the original callback.
func (t *WrapperTracer) WithSpan(ctx context.Context, name string, body func(context.Context) error) error {
	return t.otelTracer().WithSpan(ctx, name, func(ctx context.Context) error {
		span := oteltrace.CurrentSpan(ctx)
		if spanWithExtension, ok := span.(migration.OverrideTracerSpanExtension); ok {
			spanWithExtension.OverrideTracer(t)
		}
		ctx = t.bridge.ContextWithBridgeSpan(ctx, span)
		return body(ctx)
	})
}

// Start forwards the call to the wrapped tracer. It also tries to
// override the tracer of the returned span if the span implements the
// OverrideTracerSpanExtension interface.
func (t *WrapperTracer) Start(ctx context.Context, name string, opts ...oteltrace.SpanOption) (context.Context, oteltrace.Span) {
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
func (t *WrapperTracer) DeferredContextSetupHook(ctx context.Context, span oteltrace.Span) context.Context {
	if tracerWithExtension, ok := t.otelTracer().(migration.DeferredContextSetupTracerExtension); ok {
		ctx = tracerWithExtension.DeferredContextSetupHook(ctx, span)
	}
	ctx = oteltrace.SetCurrentSpan(ctx, span)
	return ctx
}
