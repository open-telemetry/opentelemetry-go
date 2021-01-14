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
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/codes"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
	"go.opentelemetry.io/otel/internal/matchers"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/trace"
)

func TestSpan(t *testing.T) {
	t.Run("#Tracer", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns the tracer used to start the span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.Tracer()).ToEqual(tracer)
		})
	})

	t.Run("#End", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("ends the span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
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

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			subject.End()

			expectedEndTime, ok := subject.EndTime()
			e.Expect(ok).ToBeTrue()

			subject.End()

			endTime, ok := subject.EndTime()
			e.Expect(ok).ToBeTrue()
			e.Expect(endTime).ToEqual(expectedEndTime)
		})

		t.Run("uses the time from WithTimestamp", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			expectedEndTime := time.Now().AddDate(5, 0, 0)
			subject.End(trace.WithTimestamp(expectedEndTime))

			e.Expect(subject.Ended()).ToBeTrue()

			endTime, ok := subject.EndTime()
			e.Expect(ok).ToBeTrue()

			e.Expect(endTime).ToEqual(expectedEndTime)
		})
	})

	t.Run("#RecordError", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("records an error", func(t *testing.T) {
			t.Parallel()

			scenarios := []struct {
				err error
				typ string
				msg string
			}{
				{
					err: ottest.NewTestError("test error"),
					typ: "go.opentelemetry.io/otel/internal/internaltest.TestError",
					msg: "test error",
				},
				{
					err: errors.New("test error 2"),
					typ: "*errors.errorString",
					msg: "test error 2",
				},
			}

			for _, s := range scenarios {
				e := matchers.NewExpecter(t)

				tracer := tp.Tracer(t.Name())
				_, span := tracer.Start(context.Background(), "test")

				subject, ok := span.(*oteltest.Span)
				e.Expect(ok).ToBeTrue()

				testTime := time.Now()
				subject.RecordError(s.err, trace.WithTimestamp(testTime))

				expectedEvents := []oteltest.Event{{
					Timestamp: testTime,
					Name:      "error",
					Attributes: map[label.Key]label.Value{
						label.Key("error.type"):    label.StringValue(s.typ),
						label.Key("error.message"): label.StringValue(s.msg),
					},
				}}
				e.Expect(subject.Events()).ToEqual(expectedEvents)
				e.Expect(subject.StatusCode()).ToEqual(codes.Error)
				e.Expect(subject.StatusMessage()).ToEqual("")
			}
		})

		t.Run("sets span status if provided", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			errMsg := "test error message"
			testErr := ottest.NewTestError(errMsg)
			testTime := time.Now()
			subject.RecordError(testErr, trace.WithTimestamp(testTime))

			expectedEvents := []oteltest.Event{{
				Timestamp: testTime,
				Name:      "error",
				Attributes: map[label.Key]label.Value{
					label.Key("error.type"):    label.StringValue("go.opentelemetry.io/otel/internal/internaltest.TestError"),
					label.Key("error.message"): label.StringValue(errMsg),
				},
			}}
			e.Expect(subject.Events()).ToEqual(expectedEvents)
			e.Expect(subject.StatusCode()).ToEqual(codes.Error)
		})

		t.Run("cannot be set after the span has ended", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			subject.End()
			subject.RecordError(errors.New("ignored error"))

			e.Expect(len(subject.Events())).ToEqual(0)
		})

		t.Run("has no effect with nil error", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			subject.RecordError(nil)

			e.Expect(len(subject.Events())).ToEqual(0)
		})
	})

	t.Run("#IsRecording", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns true", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.IsRecording()).ToBeTrue()
		})
	})

	t.Run("#SpanContext", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns a valid SpanContext", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.SpanContext().IsValid()).ToBeTrue()
		})

		t.Run("returns a consistent value", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, subject := tracer.Start(context.Background(), "test")

			e.Expect(subject.SpanContext()).ToEqual(subject.SpanContext())
		})
	})

	t.Run("#Name", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns the most recently set name on the span", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			originalName := "test"
			_, span := tracer.Start(context.Background(), originalName)

			subject, ok := span.(*oteltest.Span)
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

			tracer := tp.Tracer(t.Name())
			originalName := "test"
			_, span := tracer.Start(context.Background(), originalName)

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			subject.End()
			subject.SetName("new name")

			e.Expect(subject.Name()).ToEqual(originalName)
		})
	})

	t.Run("#Attributes", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns an empty map by default", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(subject.Attributes()).ToEqual(map[label.Key]label.Value{})
		})

		t.Run("returns the most recently set attributes", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			attr1 := label.String("key1", "value1")
			attr2 := label.String("key2", "value2")
			attr3 := label.String("key3", "value3")
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

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			expectedAttr := label.String("key", "value")
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

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			var wg sync.WaitGroup

			wg.Add(2)

			go func() {
				defer wg.Done()

				subject.SetAttributes(label.String("key", "value"))
			}()

			go func() {
				defer wg.Done()

				subject.Attributes()
			}()

			wg.Wait()
		})
	})

	t.Run("#Links", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns an empty map by default", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(len(subject.Links())).ToEqual(0)
		})
	})

	t.Run("#Events", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns an empty slice by default", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(len(subject.Events())).ToEqual(0)
		})

		t.Run("returns all of the added events", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			event1Name := "event1"
			event1Attributes := []label.KeyValue{
				label.String("event1Attr1", "foo"),
				label.String("event1Attr2", "bar"),
			}

			event1Start := time.Now()
			subject.AddEvent(event1Name, trace.WithAttributes(event1Attributes...))
			event1End := time.Now()

			event2Timestamp := time.Now().AddDate(5, 0, 0)
			event2Name := "event1"
			event2Attributes := []label.KeyValue{
				label.String("event2Attr", "abc"),
			}

			subject.AddEvent(event2Name, trace.WithTimestamp(event2Timestamp), trace.WithAttributes(event2Attributes...))

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

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			subject.AddEvent("test")

			e.Expect(len(subject.Events())).ToEqual(1)

			expectedEvent := subject.Events()[0]

			subject.End()
			subject.AddEvent("should not occur")

			e.Expect(len(subject.Events())).ToEqual(1)
			e.Expect(subject.Events()[0]).ToEqual(expectedEvent)
		})
	})

	t.Run("#Status", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("defaults to OK", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test")

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(subject.StatusCode()).ToEqual(codes.Unset)

			subject.End()

			e.Expect(subject.StatusCode()).ToEqual(codes.Unset)
		})

		statuses := []codes.Code{
			codes.Unset,
			codes.Error,
			codes.Ok,
		}

		for _, status := range statuses {
			t.Run("returns the most recently set status on the span", func(t *testing.T) {
				t.Parallel()

				e := matchers.NewExpecter(t)

				tracer := tp.Tracer(t.Name())
				_, span := tracer.Start(context.Background(), "test")

				subject, ok := span.(*oteltest.Span)
				e.Expect(ok).ToBeTrue()

				subject.SetStatus(codes.Ok, "Ok")
				subject.SetStatus(status, "Yo!")

				e.Expect(subject.StatusCode()).ToEqual(status)
				e.Expect(subject.StatusMessage()).ToEqual("Yo!")
			})

			t.Run("cannot be changed after the span has been ended", func(t *testing.T) {
				t.Parallel()

				e := matchers.NewExpecter(t)

				tracer := tp.Tracer(t.Name())
				_, span := tracer.Start(context.Background(), "test")

				subject, ok := span.(*oteltest.Span)
				e.Expect(ok).ToBeTrue()

				originalStatus := codes.Ok

				subject.SetStatus(originalStatus, "OK")
				subject.End()
				subject.SetStatus(status, fmt.Sprint(status))

				e.Expect(subject.StatusCode()).ToEqual(originalStatus)
				e.Expect(subject.StatusMessage()).ToEqual("OK")
			})
		}
	})

	t.Run("#SpanKind", func(t *testing.T) {
		tp := oteltest.NewTracerProvider()
		t.Run("returns the value given at start", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := tp.Tracer(t.Name())
			_, span := tracer.Start(context.Background(), "test",
				trace.WithSpanKind(trace.SpanKindConsumer))

			subject, ok := span.(*oteltest.Span)
			e.Expect(ok).ToBeTrue()
			subject.End()

			e.Expect(subject.SpanKind()).ToEqual(trace.SpanKindConsumer)
		})
	})
}
