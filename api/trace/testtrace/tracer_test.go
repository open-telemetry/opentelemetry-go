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

package testtrace_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/testharness"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/api/trace/testtrace"
	"go.opentelemetry.io/otel/internal/matchers"
)

func TestTracer(t *testing.T) {
	testharness.NewHarness(t).TestTracer(func() trace.Tracer {
		return testtrace.NewTracer()
	})

	t.Run("#Start", func(t *testing.T) {
		testTracedSpan(t, func(tracer trace.Tracer, name string) (trace.Span, error) {
			_, span := tracer.Start(context.Background(), name)

			return span, nil
		})

		t.Run("uses the start time from WithStartTime", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			expectedStartTime := time.Now().AddDate(5, 0, 0)

			subject := testtrace.NewTracer()
			_, span := subject.Start(context.Background(), "test", trace.WithStartTime(expectedStartTime))

			testSpan, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(testSpan.StartTime()).ToEqual(expectedStartTime)
		})

		t.Run("uses the attributes from WithAttributes", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			attr1 := core.Key("a").String("1")
			attr2 := core.Key("b").String("2")

			subject := testtrace.NewTracer()
			_, span := subject.Start(context.Background(), "test", trace.WithAttributes(attr1, attr2))

			testSpan, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			attributes := testSpan.Attributes()
			e.Expect(attributes[attr1.Key]).ToEqual(attr1.Value)
			e.Expect(attributes[attr2.Key]).ToEqual(attr2.Value)
		})

		t.Run("uses the parent's span context from ChildOf", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := testtrace.NewTracer()

			_, parent := subject.Start(context.Background(), "parent")
			parentSpanContext := parent.SpanContext()

			_, span := subject.Start(context.Background(), "child", trace.ChildOf(parentSpanContext))

			testSpan, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).ToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID()).ToEqual(parentSpanContext.SpanID)
		})

		t.Run("defers to ChildOf if the provided context also contains a parent span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := testtrace.NewTracer()

			_, parent := subject.Start(context.Background(), "parent")
			parentSpanContext := parent.SpanContext()

			ctx, _ := subject.Start(context.Background(), "should be ignored")
			_, span := subject.Start(ctx, "child", trace.ChildOf(parentSpanContext))

			testSpan, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).ToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID()).ToEqual(parentSpanContext.SpanID)
		})

		t.Run("uses the parent's span context from FollowsFrom", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := testtrace.NewTracer()

			_, parent := subject.Start(context.Background(), "parent")
			parentSpanContext := parent.SpanContext()

			_, span := subject.Start(context.Background(), "child", trace.FollowsFrom(parentSpanContext))

			testSpan, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).ToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID()).ToEqual(parentSpanContext.SpanID)
		})

		t.Run("defers to FollowsFrom if the provided context also contains a parent span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := testtrace.NewTracer()

			_, parent := subject.Start(context.Background(), "parent")
			parentSpanContext := parent.SpanContext()

			ctx, _ := subject.Start(context.Background(), "should be ignored")
			_, span := subject.Start(ctx, "child", trace.FollowsFrom(parentSpanContext))

			testSpan, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			childSpanContext := testSpan.SpanContext()
			e.Expect(childSpanContext.TraceID).ToEqual(parentSpanContext.TraceID)
			e.Expect(childSpanContext.SpanID).NotToEqual(parentSpanContext.SpanID)
			e.Expect(testSpan.ParentSpanID()).ToEqual(parentSpanContext.SpanID)
		})

		t.Run("uses the links provided through LinkedTo", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := testtrace.NewTracer()

			_, span := subject.Start(context.Background(), "link1")
			link1 := trace.Link{
				SpanContext: span.SpanContext(),
				Attributes: []core.KeyValue{
					core.Key("a").String("1"),
				},
			}

			_, span = subject.Start(context.Background(), "link2")
			link2 := trace.Link{
				SpanContext: span.SpanContext(),
				Attributes: []core.KeyValue{
					core.Key("b").String("2"),
				},
			}

			_, span = subject.Start(context.Background(), "test", trace.LinkedTo(link1.SpanContext, link1.Attributes...), trace.LinkedTo(link2.SpanContext, link2.Attributes...))

			testSpan, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			links := testSpan.Links()
			e.Expect(links[link1.SpanContext]).ToEqual(link1.Attributes)
			e.Expect(links[link2.SpanContext]).ToEqual(link2.Attributes)
		})
	})

	t.Run("#WithSpan", func(t *testing.T) {
		testTracedSpan(t, func(tracer trace.Tracer, name string) (trace.Span, error) {
			var span trace.Span

			err := tracer.WithSpan(context.Background(), name, func(ctx context.Context) error {
				span = trace.CurrentSpan(ctx)

				return nil
			})

			return span, err
		})
	})
}

func testTracedSpan(t *testing.T, fn func(tracer trace.Tracer, name string) (trace.Span, error)) {
	t.Run("starts a span with the expected name", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		subject := testtrace.NewTracer()

		expectedName := "test name"
		span, err := fn(subject, expectedName)

		e.Expect(err).ToBeNil()

		testSpan, ok := span.(*testtrace.Span)
		e.Expect(ok).ToBeTrue()

		e.Expect(testSpan.Name()).ToEqual(expectedName)
	})

	t.Run("uses the current time as the start time", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		subject := testtrace.NewTracer()

		start := time.Now()
		span, err := fn(subject, "test")
		end := time.Now()

		e.Expect(err).ToBeNil()

		testSpan, ok := span.(*testtrace.Span)
		e.Expect(ok).ToBeTrue()

		e.Expect(testSpan.StartTime()).ToBeTemporally(matchers.AfterOrSameTime, start)
		e.Expect(testSpan.StartTime()).ToBeTemporally(matchers.BeforeOrSameTime, end)
	})

	t.Run("appends the span to the list of Spans", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		subject := testtrace.NewTracer()
		subject.Start(context.Background(), "span1")

		e.Expect(len(subject.Spans())).ToEqual(1)

		span, err := fn(subject, "span2")
		e.Expect(err).ToBeNil()

		spans := subject.Spans()

		e.Expect(len(spans)).ToEqual(2)
		e.Expect(spans[1]).ToEqual(span)
	})

	t.Run("can be run concurrently with another call", func(t *testing.T) {
		t.Parallel()

		e := matchers.NewExpecter(t)

		subject := testtrace.NewTracer()

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

		e.Expect(len(subject.Spans())).ToEqual(numSpans)
	})
}
