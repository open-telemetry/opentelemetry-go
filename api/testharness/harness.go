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

package testharness

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/kv"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/internal/matchers"
)

type Harness struct {
	t *testing.T
}

func NewHarness(t *testing.T) *Harness {
	return &Harness{
		t: t,
	}
}

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
			e.Expect(span.SpanContext()).NotToEqual(trace.EmptySpanContext())
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

	h.t.Run("#WithSpan", func(t *testing.T) {
		t.Run("returns nil if the body does not return an error", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			err := subject.WithSpan(context.Background(), "test", func(ctx context.Context) error {
				return nil
			})

			e.Expect(err).ToBeNil()
		})

		t.Run("propagates the error from the body if the body errors", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			expectedErr := errors.New("test error")

			err := subject.WithSpan(context.Background(), "test", func(ctx context.Context) error {
				return expectedErr
			})

			e.Expect(err).ToMatchError(expectedErr)
		})

		t.Run("propagates the original context to the body", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			ctxKey := testCtxKey{}
			ctxValue := "ctx value"
			ctx := context.WithValue(context.Background(), ctxKey, ctxValue)

			var actualCtx context.Context

			err := subject.WithSpan(ctx, "test", func(ctx context.Context) error {
				actualCtx = ctx

				return nil
			})

			e.Expect(err).ToBeNil()

			e.Expect(actualCtx.Value(ctxKey)).ToEqual(ctxValue)
		})

		t.Run("stores a span containing the expected properties on the context provided to the body function", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			var span trace.Span

			err := subject.WithSpan(context.Background(), "test", func(ctx context.Context) error {
				span = trace.SpanFromContext(ctx)

				return nil
			})

			e.Expect(err).ToBeNil()

			e.Expect(span).NotToBeNil()

			e.Expect(span.Tracer()).ToEqual(subject)
			e.Expect(span.SpanContext().IsValid()).ToBeTrue()
		})

		t.Run("starts spans with unique trace and span IDs", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			var span1 trace.Span
			var span2 trace.Span

			err := subject.WithSpan(context.Background(), "span1", func(ctx context.Context) error {
				span1 = trace.SpanFromContext(ctx)

				return nil
			})

			e.Expect(err).ToBeNil()

			err = subject.WithSpan(context.Background(), "span2", func(ctx context.Context) error {
				span2 = trace.SpanFromContext(ctx)

				return nil
			})

			e.Expect(err).ToBeNil()

			sc1 := span1.SpanContext()
			sc2 := span2.SpanContext()

			e.Expect(sc1.TraceID).NotToEqual(sc2.TraceID)
			e.Expect(sc1.SpanID).NotToEqual(sc2.SpanID)
		})

		t.Run("propagates a parent's trace ID through the context", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)
			subject := subjectFactory()

			ctx, parent := subject.Start(context.Background(), "parent")

			var child trace.Span

			err := subject.WithSpan(ctx, "child", func(ctx context.Context) error {
				child = trace.SpanFromContext(ctx)

				return nil
			})

			e.Expect(err).ToBeNil()

			psc := parent.SpanContext()
			csc := child.SpanContext()

			e.Expect(csc.TraceID).ToEqual(psc.TraceID)
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
			span.AddEvent(context.Background(), "test event")
		},
		"#AddEventWithTimestamp": func(span trace.Span) {
			span.AddEventWithTimestamp(context.Background(), time.Now(), "test event")
		},
		"#SetStatus": func(span trace.Span) {
			span.SetStatus(codes.Internal, "internal")
		},
		"#SetName": func(span trace.Span) {
			span.SetName("new name")
		},
		"#SetAttributes": func(span trace.Span) {
			span.SetAttributes(kv.String("key1", "value"), kv.Int("key2", 123))
		},
	}
	var mechanisms = map[string]func() trace.Span{
		"Span created via Tracer#Start": func() trace.Span {
			tracer := tracerFactory()
			_, subject := tracer.Start(context.Background(), "test")

			return subject
		},
		"Span created via Tracer#WithSpan": func() trace.Span {
			tracer := tracerFactory()

			var actualCtx context.Context

			_ = tracer.WithSpan(context.Background(), "test", func(ctx context.Context) error {
				actualCtx = ctx

				return nil
			})

			return trace.SpanFromContext(actualCtx)
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
