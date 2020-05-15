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

package trace_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
)

func TestSetCurrentSpanOverridesPreviouslySetSpan(t *testing.T) {
	originalSpan := trace.NoopSpan{}
	expectedSpan := mockSpan{}

	ctx := context.Background()

	ctx = trace.ContextWithSpan(ctx, originalSpan)
	ctx = trace.ContextWithSpan(ctx, expectedSpan)

	if span := trace.SpanFromContext(ctx); span != expectedSpan {
		t.Errorf("Want: %v, but have: %v", expectedSpan, span)
	}
}

func TestCurrentSpan(t *testing.T) {
	for _, testcase := range []struct {
		name string
		ctx  context.Context
		want trace.Span
	}{
		{
			name: "CurrentSpan() returns a NoopSpan{} from an empty context",
			ctx:  context.Background(),
			want: trace.NoopSpan{},
		},
		{
			name: "CurrentSpan() returns current span if set",
			ctx:  trace.ContextWithSpan(context.Background(), mockSpan{}),
			want: mockSpan{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			// proto: CurrentSpan(ctx context.Context) trace.Span
			have := trace.SpanFromContext(testcase.ctx)
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

// a duplicate of trace.NoopSpan for testing
type mockSpan struct{}

var _ trace.Span = mockSpan{}

// SpanContext returns an invalid span context.
func (mockSpan) SpanContext() trace.SpanContext {
	return trace.EmptySpanContext()
}

// IsRecording always returns false for mockSpan.
func (mockSpan) IsRecording() bool {
	return false
}

// SetStatus does nothing.
func (mockSpan) SetStatus(status codes.Code, msg string) {
}

// SetName does nothing.
func (mockSpan) SetName(name string) {
}

// SetError does nothing.
func (mockSpan) SetError(v bool) {
}

// SetAttributes does nothing.
func (mockSpan) SetAttributes(attributes ...kv.KeyValue) {
}

// SetAttribute does nothing.
func (mockSpan) SetAttribute(k string, v interface{}) {
}

// End does nothing.
func (mockSpan) End(options ...trace.EndOption) {
}

// RecordError does nothing.
func (mockSpan) RecordError(ctx context.Context, err error, opts ...trace.ErrorOption) {
}

// Tracer returns noop implementation of Tracer.
func (mockSpan) Tracer() trace.Tracer {
	return trace.NoopTracer{}
}

// Event does nothing.
func (mockSpan) AddEvent(ctx context.Context, name string, attrs ...kv.KeyValue) {
}

// AddEventWithTimestamp does nothing.
func (mockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...kv.KeyValue) {
}
