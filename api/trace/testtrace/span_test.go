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

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/api/trace/testtrace"
	"go.opentelemetry.io/otel/internal/matchers"
)

func TestSpan(t *testing.T) {
	t.Run("#Tracer", func(t *testing.T) {
		t.Run("returns the tracer used to start the span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.Tracer()).ToEqual(tracer)
		})
	})

	t.Run("#End", func(t *testing.T) {
		t.Run("ends the span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(subject.Ended()).ToBeFalse()

			_, ok = subject.EndTime()
			e.Expect(ok).ToBeFalse()

			start := time.Now()

			subject.End()

			end := time.Now()

			e.Expect(subject.Ended()).ToBeTrue()

			endTime, ok := subject.EndTime()
			e.Expect(ok).ToBeTrue()

			e.Expect(endTime).ToBeTemporally(matchers.AfterOrSameTime, start)
			e.Expect(endTime).ToBeTemporally(matchers.BeforeOrSameTime, end)
		})

		t.Run("only takes effect the first time it is called", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			subject.End()

			expectedEndTime, ok := subject.EndTime()
			e.Expect(ok).ToBeTrue()

			subject.End()

			endTime, ok := subject.EndTime()
			e.Expect(ok).ToBeTrue()
			e.Expect(endTime).ToEqual(expectedEndTime)
		})

		t.Run("uses the time from WithEndTime", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			expectedEndTime := time.Now().AddDate(5, 0, 0)
			subject.End(trace.WithEndTime(expectedEndTime))

			e.Expect(subject.Ended()).ToBeTrue()

			endTime, ok := subject.EndTime()
			e.Expect(ok).ToBeTrue()

			e.Expect(endTime).ToEqual(expectedEndTime)
		})
	})

	t.Run("#IsRecording", func(t *testing.T) {
		t.Run("returns true", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.IsRecording()).ToBeTrue()
		})
	})

	t.Run("#SpanContext", func(t *testing.T) {
		t.Run("returns a valid SpanContext", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.SpanContext().IsValid()).ToBeTrue()
		})

		t.Run("returns a consistent value", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.SpanContext()).ToEqual(subject.SpanContext())
		})
	})

	t.Run("#Name", func(t *testing.T) {
		t.Run("returns the most recently set name on the span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			originalName := "test"
			_, span := tracer.Start(context.Background(), originalName)

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(subject.Name()).ToEqual(originalName)

			subject.SetName("in-between")

			newName := "new name"

			subject.SetName(newName)

			e.Expect(subject.Name()).ToEqual(newName)
		})

		t.Run("cannot be changed after the span has been ended", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			originalName := "test"
			_, span := tracer.Start(context.Background(), originalName)

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			subject.End()
			subject.SetName("new name")

			e.Expect(subject.Name()).ToEqual(originalName)
		})
	})

	t.Run("#Attributes", func(t *testing.T) {
		t.Run("returns an empty map by default", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(subject.Attributes()).ToEqual(map[core.Key]core.Value{})
		})

		t.Run("returns the most recently set attributes", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			attr1 := core.Key("key1").String("value1")
			attr2 := core.Key("key2").String("value2")
			attr3 := core.Key("key3").String("value3")
			unexpectedAttr := attr2.Key.String("unexpected")

			subject.SetAttributes(attr1, unexpectedAttr, attr3)
			subject.SetAttributes(attr2)

			attributes := subject.Attributes()

			e.Expect(attributes[attr1.Key]).ToEqual(attr1.Value)
			e.Expect(attributes[attr2.Key]).ToEqual(attr2.Value)
			e.Expect(attributes[attr3.Key]).ToEqual(attr3.Value)
		})

		t.Run("cannot be changed after the span has been ended", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			expectedAttr := core.Key("key").String("value")
			subject.SetAttributes(expectedAttr)
			subject.End()

			unexpectedAttr := expectedAttr.Key.String("unexpected")
			subject.SetAttributes(unexpectedAttr)
			subject.End()

			attributes := subject.Attributes()
			e.Expect(attributes[expectedAttr.Key]).ToEqual(expectedAttr.Value)
		})

		t.Run("can be used concurrently with setter", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			var wg sync.WaitGroup

			wg.Add(2)

			go func() {
				defer wg.Done()

				subject.SetAttributes(core.Key("key").String("value"))
			}()

			go func() {
				defer wg.Done()

				subject.Attributes()
			}()

			wg.Wait()
		})
	})

	t.Run("#Links", func(t *testing.T) {
		t.Run("returns an empty map by default", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(len(subject.Links())).ToEqual(0)
		})
	})

	t.Run("#Events", func(t *testing.T) {
		t.Run("returns an empty slice by default", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(len(subject.Events())).ToEqual(0)
		})

		t.Run("returns all of the added events", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			event1Name := "event1"
			event1Attributes := []core.KeyValue{
				core.Key("event1Attr1").String("foo"),
				core.Key("event1Attr2").String("bar"),
			}

			event1Start := time.Now()
			subject.AddEvent(context.Background(), event1Name, event1Attributes...)
			event1End := time.Now()

			event2Timestamp := time.Now().AddDate(5, 0, 0)
			event2Name := "event1"
			event2Attributes := []core.KeyValue{
				core.Key("event2Attr").String("abc"),
			}

			subject.AddEventWithTimestamp(context.Background(), event2Timestamp, event2Name, event2Attributes...)

			events := subject.Events()

			e.Expect(len(events)).ToEqual(2)

			event1 := events[0]

			e.Expect(event1.Timestamp).ToBeTemporally(matchers.AfterOrSameTime, event1Start)
			e.Expect(event1.Timestamp).ToBeTemporally(matchers.BeforeOrSameTime, event1End)
			e.Expect(event1.Name).ToEqual(event1Name)

			for _, attr := range event1Attributes {
				e.Expect(event1.Attributes[attr.Key]).ToEqual(attr.Value)
			}

			event2 := events[1]

			e.Expect(event2.Timestamp).ToEqual(event2Timestamp)
			e.Expect(event2.Name).ToEqual(event2Name)

			for _, attr := range event2Attributes {
				e.Expect(event2.Attributes[attr.Key]).ToEqual(attr.Value)
			}
		})

		t.Run("cannot be changed after the span has been ended", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			subject.AddEvent(context.Background(), "test")

			e.Expect(len(subject.Events())).ToEqual(1)

			expectedEvent := subject.Events()[0]

			subject.End()
			subject.AddEvent(context.Background(), "should not occur")

			e.Expect(len(subject.Events())).ToEqual(1)
			e.Expect(subject.Events()[0]).ToEqual(expectedEvent)
		})
	})

	t.Run("#Status", func(t *testing.T) {
		t.Run("defaults to OK", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(subject.Status()).ToEqual(codes.OK)

			subject.End()

			e.Expect(subject.Status()).ToEqual(codes.OK)
		})

		statuses := []codes.Code{
			codes.OK,
			codes.Canceled,
			codes.Unknown,
			codes.InvalidArgument,
			codes.DeadlineExceeded,
			codes.NotFound,
			codes.AlreadyExists,
			codes.PermissionDenied,
			codes.ResourceExhausted,
			codes.FailedPrecondition,
			codes.Aborted,
			codes.OutOfRange,
			codes.Unimplemented,
			codes.Internal,
			codes.Unavailable,
			codes.DataLoss,
			codes.Unauthenticated,
		}

		for _, status := range statuses {
			t.Run("returns the most recently set status on the span", func(t *testing.T) {
				t.Parallel()

				e := matchers.NewExpecter(t)

				tracer := testtrace.NewTracer()
				_, span := tracer.Start(context.Background(), "test")

				subject, ok := span.(*testtrace.Span)
				e.Expect(ok).ToBeTrue()

				subject.SetStatus(codes.OK)
				subject.SetStatus(status)

				e.Expect(subject.Status()).ToEqual(status)
			})

			t.Run("cannot be changed after the span has been ended", func(t *testing.T) {
				t.Parallel()

				e := matchers.NewExpecter(t)

				tracer := testtrace.NewTracer()
				_, span := tracer.Start(context.Background(), "test")

				subject, ok := span.(*testtrace.Span)
				e.Expect(ok).ToBeTrue()

				originalStatus := codes.OK

				subject.SetStatus(originalStatus)
				subject.End()
				subject.SetStatus(status)

				e.Expect(subject.Status()).ToEqual(originalStatus)
			})
		}
	})
}
