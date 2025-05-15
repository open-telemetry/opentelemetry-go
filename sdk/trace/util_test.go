// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
)

func basicTracerProvider(t *testing.T) *TracerProvider {
	tp := NewTracerProvider(WithSampler(AlwaysSample()))
	t.Cleanup(func() {
		assert.NoError(t, tp.Shutdown(context.Background()))
	})
	return tp
}

type testError string

var _ error = testError("")

func newTestError(s string) error {
	return testError(s)
}

func (e testError) Error() string {
	return string(e)
}

// harness is a testing harness used to test implementations of the
// OpenTelemetry API.
type harness struct {
	t *testing.T
}

// newHarness returns an instantiated *harness using t.
func newHarness(t *testing.T) *harness {
	return &harness{
		t: t,
	}
}

func (h *harness) testSpan(tracerFactory func() trace.Tracer) {
	methods := map[string]func(span trace.Span){
		"#End": func(span trace.Span) {
			span.End()
		},
		"#AddEvent": func(span trace.Span) {
			span.AddEvent("test event")
		},
		"#AddEventWithTimestamp": func(span trace.Span) {
			span.AddEvent("test event", trace.WithTimestamp(time.Now().Add(1*time.Second)))
		},
		"#SetStatus": func(span trace.Span) {
			span.SetStatus(codes.Error, "internal")
		},
		"#SetName": func(span trace.Span) {
			span.SetName("new name")
		},
		"#SetAttributes": func(span trace.Span) {
			span.SetAttributes(attribute.String("key1", "value"), attribute.Int("key2", 123))
		},
	}
	mechanisms := map[string]func() trace.Span{
		"Span created via Tracer#Start": func() trace.Span {
			tracer := tracerFactory()
			_, subject := tracer.Start(context.Background(), "test")

			return subject
		},
		"Span created via span.TracerProvider()": func() trace.Span {
			ctx, spanA := tracerFactory().Start(context.Background(), "span1")

			_, spanB := spanA.TracerProvider().Tracer("second").Start(ctx, "span2")
			return spanB
		},
	}

	for mechanismName, mechanism := range mechanisms {
		h.t.Run(mechanismName, func(t *testing.T) {
			for methodName, method := range methods {
				t.Run(methodName, func(t *testing.T) {
					t.Run("is thread-safe", func(t *testing.T) {
						t.Parallel()

						span := mechanism()

						wg := &sync.WaitGroup{}
						wg.Add(2)

						go func() {
							defer wg.Done()

							method(span)
						}()

						go func() {
							defer wg.Done()

							method(span)
						}()

						wg.Wait()
					})
				})
			}

			t.Run("#End", func(t *testing.T) {
				t.Run("can be called multiple times", func(t *testing.T) {
					t.Parallel()

					span := mechanism()

					span.End()
					span.End()
				})
			})
		})
	}
}

type testCtxKey struct{}
