// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opentracing

import (
	"testing"

	"go.opentelemetry.io/otel/bridge/opentracing/internal"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
)

type namedMockTracer struct {
	name string
	*internal.MockTracer
}

type namedMockTracerProvider struct{ embedded.TracerProvider }

var _ trace.TracerProvider = (*namedMockTracerProvider)(nil)

// Tracer returns the WrapperTracer associated with the WrapperTracerProvider.
func (p *namedMockTracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return &namedMockTracer{
		name:       name,
		MockTracer: internal.NewMockTracer(),
	}
}

func TestTracerProvider(t *testing.T) {
	// assertMockTracerName casts tracer into a named mock tracer provided by
	// namedMockTracerProvider, and asserts against its name
	assertMockTracerName := func(t *testing.T, tracer trace.Tracer, name string) {
		// Unwrap the tracer
		wrapped := tracer.(*WrapperTracer)
		tracer = wrapped.tracer

		// Cast into the underlying type and assert
		if mock, ok := tracer.(*namedMockTracer); ok {
			if name != mock.name {
				t.Errorf("expected name %q, got %q", name, mock.name)
			}
		} else if !ok {
			t.Errorf("expected *namedMockTracer, got %T", mock)
		}
	}

	var (
		foobar   = "foobar"
		bazbar   = "bazbar"
		provider = NewTracerProvider(nil, &namedMockTracerProvider{})
	)

	t.Run("Tracers should be created with foobar from provider", func(t *testing.T) {
		tracer := provider.Tracer(foobar)
		assertMockTracerName(t, tracer, foobar)
	})

	t.Run("Repeated requests to create a tracer should provide the existing tracer", func(t *testing.T) {
		tracer1 := provider.Tracer(foobar)
		assertMockTracerName(t, tracer1, foobar)
		tracer2 := provider.Tracer(foobar)
		assertMockTracerName(t, tracer2, foobar)
		tracer3 := provider.Tracer(bazbar)
		assertMockTracerName(t, tracer3, bazbar)

		if tracer1 != tracer2 {
			t.Errorf("expected the same tracer, got different tracers")
		}
		if tracer1 == tracer3 || tracer2 == tracer3 {
			t.Errorf("expected different tracers, got the same tracer")
		}
	})
}
