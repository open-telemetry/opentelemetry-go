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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/internal/matchers"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/trace"
)

// Harness is a testing harness used to test implementations of the
// OpenTelemetry API.
type Harness struct {
	t *testing.T
}

// NewHarness returns an instantiated *Harness using t.
func NewHarness(t *testing.T) *Harness {
	return &Harness{
		t: t,
	}
}

// TestTracer runs validation tests for an implementation of the OpenTelemetry
// Tracer API.
func (h *Harness) TestTracer(subjectFactory func() trace.Tracer) {
	h.t.Run("#Start", func(t *testing.T) {
		t.Run("propagates the original context", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			ctxKey := testCtxKey{}
			ctxValue := "ctx value"
			ctx := context.WithValue(context.Background(), ctxKey, ctxValue)

			ctx, _ = subject.Start(ctx, "test")

			e.Expect(ctx.Value(ctxKey)).ToEqual(ctxValue)
		})

		t.Run("returns a span containing the expected properties", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			_, span := subject.Start(context.Background(), "test")

			e.Expect(span).NotToBeNil()

			e.Expect(span.Tracer()).ToEqual(subject)
			e.Expect(span.SpanContext().IsValid()).ToBeTrue()
		})

		t.Run("stores the span on the provided context", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			ctx, span := subject.Start(context.Background(), "test")

			e.Expect(span).NotToBeNil()
			e.Expect(span.SpanContext()).NotToEqual(trace.SpanContext{})
			e.Expect(trace.SpanFromContext(ctx)).ToEqual(span)
		})

		t.Run("starts spans with unique trace and span IDs", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			_, span1 := subject.Start(context.Background(), "span1")
			_, span2 := subject.Start(context.Background(), "span2")

			sc1 := span1.SpanContext()
			sc2 := span2.SpanContext()

			e.Expect(sc1.TraceID).NotToEqual(sc2.TraceID)
			e.Expect(sc1.SpanID).NotToEqual(sc2.SpanID)
		})

		t.Run("records the span if specified", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			_, span := subject.Start(context.Background(), "span", trace.WithRecord())

			e.Expect(span.IsRecording()).ToBeTrue()
		})

		t.Run("propagates a parent's trace ID through the context", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			ctx, parent := subject.Start(context.Background(), "parent")
			_, child := subject.Start(ctx, "child")

			psc := parent.SpanContext()
			csc := child.SpanContext()

			e.Expect(csc.TraceID).ToEqual(psc.TraceID)
			e.Expect(csc.SpanID).NotToEqual(psc.SpanID)
		})

		t.Run("ignores parent's trace ID when new root is requested", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			ctx, parent := subject.Start(context.Background(), "parent")
			_, child := subject.Start(ctx, "child", trace.WithNewRoot())

			psc := parent.SpanContext()
			csc := child.SpanContext()

			e.Expect(csc.TraceID).NotToEqual(psc.TraceID)
			e.Expect(csc.SpanID).NotToEqual(psc.SpanID)
		})

		t.Run("propagates remote parent's trace ID through the context", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			_, remoteParent := subject.Start(context.Background(), "remote parent")
			parentCtx := trace.ContextWithRemoteSpanContext(context.Background(), remoteParent.SpanContext())
			_, child := subject.Start(parentCtx, "child")

			psc := remoteParent.SpanContext()
			csc := child.SpanContext()

			e.Expect(csc.TraceID).ToEqual(psc.TraceID)
			e.Expect(csc.SpanID).NotToEqual(psc.SpanID)
		})

		t.Run("ignores remote parent's trace ID when new root is requested", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			_, remoteParent := subject.Start(context.Background(), "remote parent")
			parentCtx := trace.ContextWithRemoteSpanContext(context.Background(), remoteParent.SpanContext())
			_, child := subject.Start(parentCtx, "child", trace.WithNewRoot())

			psc := remoteParent.SpanContext()
			csc := child.SpanContext()

			e.Expect(csc.TraceID).NotToEqual(psc.TraceID)
			e.Expect(csc.SpanID).NotToEqual(psc.SpanID)
		})
	})

	h.testSpan(subjectFactory)
}

func (h *Harness) testSpan(tracerFactory func() trace.Tracer) {
	var methods = map[string]func(span trace.Span){
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
			span.SetAttributes(label.String("key1", "value"), label.Int("key2", 123))
		},
	}
	var mechanisms = map[string]func() trace.Span{
		"Span created via Tracer#Start": func() trace.Span {
			tracer := tracerFactory()
			_, subject := tracer.Start(context.Background(), "test")

			return subject
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
