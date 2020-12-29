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

package oteltest_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel/internal/matchers"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/trace"
)

func TestTracer(t *testing.T) {
	tp := oteltest.NewTracerProvider()

	oteltest.NewHarness(t).TestTracer(func() func() trace.Tracer {
		tp := oteltest.NewTracerProvider()
		var i uint64
		return func() trace.Tracer {
			return tp.Tracer(fmt.Sprintf("tracer %d", atomic.AddUint64(&i, 1)))
		}
	}())

	t.Run("#Start", func(t *testing.T) {
		testTracedSpan(t, func(tracer trace.Tracer, name string) (trace.Span, error) {
			_, span := tracer.Start(context.Background(), name)

			return span, nil
		})

		t.Run("uses the start time from WithTimestamp", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			expectedStartTime := time.Now().AddDate(5, 0, 0)

			subject := tp.Tracer(t.Name())
			_, span := subject.Start(context.Background(), "test", trace.WithTimestamp(expectedStartTime))

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(testSpan.StartTime()).ToEqual(expectedStartTime)
		})

		t.Run("uses the attributes from WithAttributes", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			attr1 := label.String("a", "1")
			attr2 := label.String("b", "2")

			subject := tp.Tracer(t.Name())
			_, span := subject.Start(context.Background(), "test", trace.WithAttributes(attr1, attr2))

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			attributes := testSpan.Attributes()
			e.Expect(attributes[attr1.Key]).ToEqual(attr1.Value)
			e.Expect(attributes[attr2.Key]).ToEqual(attr2.Value)
		})

		t.Run("uses the current span from context as parent", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := tp.Tracer(t.Name())

			parent, parentSpan := subject.Start(context.Background(), "parent")
			parentSpanContext := parentSpan.SpanContext()

			_, span := subject.Start(parent, "child")

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).ToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID()).ToEqual(parentSpanContext.SpanID)
		})

		t.Run("uses the current span from context as parent, even if it has remote span context", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := tp.Tracer(t.Name())

			parent, parentSpan := subject.Start(context.Background(), "parent")
			_, remoteParentSpan := subject.Start(context.Background(), "remote not-a-parent")
			parent = trace.ContextWithRemoteSpanContext(parent, remoteParentSpan.SpanContext())
			parentSpanContext := parentSpan.SpanContext()

			_, span := subject.Start(parent, "child")

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).ToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID()).ToEqual(parentSpanContext.SpanID)
		})

		t.Run("uses the remote span context from context as parent, if current span is missing", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := tp.Tracer(t.Name())

			_, remoteParentSpan := subject.Start(context.Background(), "remote parent")
			parent := trace.ContextWithRemoteSpanContext(context.Background(), remoteParentSpan.SpanContext())
			remoteParentSpanContext := remoteParentSpan.SpanContext()

			_, span := subject.Start(parent, "child")

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).ToEqual(remoteParentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(remoteParentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID()).ToEqual(remoteParentSpanContext.SpanID)
		})

		t.Run("creates new root when both current span and remote span context are missing", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := tp.Tracer(t.Name())

			_, parentSpan := subject.Start(context.Background(), "not-a-parent")
			_, remoteParentSpan := subject.Start(context.Background(), "remote not-a-parent")
			parentSpanContext := parentSpan.SpanContext()
			remoteParentSpanContext := remoteParentSpan.SpanContext()

			_, span := subject.Start(context.Background(), "child")

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).NotToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.TraceID).NotToEqual(remoteParentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(childSpanContext.SpanID).NotToEqual(remoteParentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID().IsValid()).ToBeFalse()
		})

		t.Run("creates new root when requested, even if both current span and remote span context are in context", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := tp.Tracer(t.Name())

			parentCtx, parentSpan := subject.Start(context.Background(), "not-a-parent")
			_, remoteParentSpan := subject.Start(context.Background(), "remote not-a-parent")
			parentSpanContext := parentSpan.SpanContext()
			remoteParentSpanContext := remoteParentSpan.SpanContext()
			parentCtx = trace.ContextWithRemoteSpanContext(parentCtx, remoteParentSpanContext)

			_, span := subject.Start(parentCtx, "child", trace.WithNewRoot())

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).NotToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.TraceID).NotToEqual(remoteParentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(childSpanContext.SpanID).NotToEqual(remoteParentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID().IsValid()).ToBeFalse()

			expectedLinks := []trace.Link{
				{
					SpanContext: parentSpanContext,
					Attributes: []label.KeyValue{
						label.String("ignored-on-demand", "current"),
					},
				},
				{
					SpanContext: remoteParentSpanContext,
					Attributes: []label.KeyValue{
						label.String("ignored-on-demand", "remote"),
					},
				},
			}
			gotLinks := testSpan.Links()
			e.Expect(gotLinks).ToMatchInAnyOrder(expectedLinks)
		})

		t.Run("uses the links provided through WithLinks", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := tp.Tracer(t.Name())

			_, span := subject.Start(context.Background(), "link1")
			link1 := trace.Link{
				SpanContext: span.SpanContext(),
				Attributes: []label.KeyValue{
					label.String("a", "1"),
				},
			}

			_, span = subject.Start(context.Background(), "link2")
			link2 := trace.Link{
				SpanContext: span.SpanContext(),
				Attributes: []label.KeyValue{
					label.String("b", "2"),
				},
			}

			_, span = subject.Start(context.Background(), "test", trace.WithLinks(link1, link2))

			testSpan, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			links := testSpan.Links()
			e.Expect(links[0].Attributes).ToEqual(link1.Attributes)
			e.Expect(links[1].Attributes).ToEqual(link2.Attributes)
		})
	})
}

func testTracedSpan(t *testing.T, fn func(tracer trace.Tracer, name string) (trace.Span, error)) {
	tp := oteltest.NewTracerProvider()
	t.Run("starts a span with the expected name", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		subject := tp.Tracer(t.Name())

		expectedName := "test name"
		span, err := fn(subject, expectedName)

		e.Expect(err).ToBeNil()

		testSpan, ok := span.(*oteltest.Span)
		e.Expect(ok).ToBeTrue()

		e.Expect(testSpan.Name()).ToEqual(expectedName)
	})

	t.Run("uses the current time as the start time", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		subject := tp.Tracer(t.Name())

		start := time.Now()
		span, err := fn(subject, "test")
		end := time.Now()

		e.Expect(err).ToBeNil()

		testSpan, ok := span.(*oteltest.Span)
		e.Expect(ok).ToBeTrue()

		e.Expect(testSpan.StartTime()).ToBeTemporally(matchers.AfterOrSameTime, start)
		e.Expect(testSpan.StartTime()).ToBeTemporally(matchers.BeforeOrSameTime, end)
	})

	t.Run("calls SpanRecorder.OnStart", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		sr := new(oteltest.StandardSpanRecorder)
		subject := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr)).Tracer(t.Name())
		subject.Start(context.Background(), "span1")

		e.Expect(len(sr.Started())).ToEqual(1)

		span, err := fn(subject, "span2")
		e.Expect(err).ToBeNil()

		spans := sr.Started()

		e.Expect(len(spans)).ToEqual(2)
		e.Expect(spans[1]).ToEqual(span)
	})

	t.Run("can be run concurrently with another call", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		sr := new(oteltest.StandardSpanRecorder)
		subject := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr)).Tracer(t.Name())

		numSpans := 2

		var wg sync.WaitGroup

		wg.Add(numSpans)

		for i := 0; i < numSpans; i++ {
			go func() {
				_, err := fn(subject, "test")
				e.Expect(err).ToBeNil()

				wg.Done()
			}()
		}

		wg.Wait()

		e.Expect(len(sr.Started())).ToEqual(numSpans)
	})
}
