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
		t.Run("starts a span with the expected name", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := testtrace.NewTracer()

			expectedName := "test name"

			_, span := subject.Start(context.Background(), expectedName)

			testSpan, ok := span.(*testtrace.Span)

			e.Expect(ok).ToBeTrue()
			e.Expect(testSpan.Name()).ToEqual(expectedName)
		})

		t.Run("uses the current time as the start time", func(t *testing.T) {
			t.Parallel()

			e := matchers.NewExpecter(t)

			subject := testtrace.NewTracer()

			start := time.Now()

			_, span := subject.Start(context.Background(), "test")

			end := time.Now()

			testSpan, ok := span.(*testtrace.Span)

			e.Expect(ok).ToBeTrue()
			e.Expect(testSpan.StartTime()).ToBeTemporally(matchers.AfterOrSameTime, start)
			e.Expect(testSpan.StartTime()).ToBeTemporally(matchers.BeforeOrSameTime, end)
		})
	})
}
