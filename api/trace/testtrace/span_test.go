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
	"testing"
	"time"

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
			_, subject := tracer.Start(context.Background(), "test")

			span, ok := subject.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			e.Expect(span.Ended()).ToBeFalse()

			_, ok = span.EndTime()
			e.Expect(ok).ToBeFalse()

			start := time.Now()

			span.End()

			end := time.Now()

			e.Expect(span.Ended()).ToBeTrue()

			endTime, ok := span.EndTime()
			e.Expect(ok).ToBeTrue()

			e.Expect(endTime).ToBeTemporally(matchers.AfterOrSameTime, start)
			e.Expect(endTime).ToBeTemporally(matchers.BeforeOrSameTime, end)
		})

		t.Run("only takes effect the first time it is called", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			tracer := testtrace.NewTracer()
			_, subject := tracer.Start(context.Background(), "test")

			span, ok := subject.(*testtrace.Span)
			e.Expect(ok).ToBeTrue()

			span.End()

			expectedEndTime, ok := span.EndTime()
			e.Expect(ok).ToBeTrue()

			span.End()

			endTime, ok := span.EndTime()
			e.Expect(ok).ToBeTrue()
			e.Expect(endTime).ToEqual(expectedEndTime)
		})
	})
}
