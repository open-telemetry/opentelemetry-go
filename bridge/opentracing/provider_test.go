// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opentracing

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
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
		tracerFns := []func() trace.Tracer{
			func() trace.Tracer {
				return provider.Tracer(foobar)
			},
			func() trace.Tracer {
				return provider.Tracer(bazbar)
			},
			func() trace.Tracer {
				return provider.Tracer(foobar, trace.WithSchemaURL("https://opentelemetry.io/schemas/1.2.0"))
			},
			func() trace.Tracer {
				return provider.Tracer(foobar, trace.WithInstrumentationAttributes(attribute.String("foo", "bar")))
			},
			func() trace.Tracer {
				return provider.Tracer(foobar, trace.WithSchemaURL("https://opentelemetry.io/schemas/1.2.0"), trace.WithInstrumentationAttributes(attribute.String("foo", "bar")))
			},
		}

		for i, fn1 := range tracerFns {
			for j, fn2 := range tracerFns {
				tracer1, tracer2 := fn1(), fn2()
				if i == j {
					if tracer1 != tracer2 {
						t.Errorf("expected the same tracer, got different tracers; i=%d j=%d", i, j)
					}
				} else {
					if tracer1 == tracer2 {
						t.Errorf("expected different tracers, got the same tracer; i=%d j=%d", i, j)
					}
				}
			}
		}
	})
}
