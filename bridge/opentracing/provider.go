// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opentracing // import "go.opentelemetry.io/otel/bridge/opentracing"

import (
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

// TracerProvider is an OpenTelemetry TracerProvider that wraps an OpenTracing
// Tracer.
type TracerProvider struct {
	embedded.TracerProvider

	bridge   *BridgeTracer
	provider trace.TracerProvider

	tracers map[wrappedTracerKey]*WrapperTracer
	mtx     sync.Mutex
}

var _ trace.TracerProvider = (*TracerProvider)(nil)

// NewTracerProvider returns a new TracerProvider that creates new instances of
// WrapperTracer from the given TracerProvider.
func NewTracerProvider(bridge *BridgeTracer, provider trace.TracerProvider) *TracerProvider {
	return &TracerProvider{
		bridge:   bridge,
		provider: provider,

		tracers: make(map[wrappedTracerKey]*WrapperTracer),
	}
}

type wrappedTracerKey struct {
	name    string
	version string
	schema  string
	attrs   attribute.Set
}

// Tracer creates a WrappedTracer that wraps the OpenTelemetry tracer for each call to
// Tracer(). Repeated calls to Tracer() with the same configuration will look up and
// return an existing instance of WrapperTracer.
func (p *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	p.mtx.Lock()
	defer p.mtx.Unlock()

	c := trace.NewTracerConfig(opts...)
	key := wrappedTracerKey{
		name:    name,
		version: c.InstrumentationVersion(),
		schema:  c.SchemaURL(),
		attrs:   c.InstrumentationAttributes(),
	}

	if t, ok := p.tracers[key]; ok {
		return t
	}

	wrapper := NewWrapperTracer(p.bridge, p.provider.Tracer(name, opts...))
	p.tracers[key] = wrapper
	return wrapper
}
