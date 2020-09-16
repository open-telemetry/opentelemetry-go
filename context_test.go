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

package otel_test

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/internal/trace/noop"
	"go.opentelemetry.io/otel/label"
)

func TestSetCurrentSpanOverridesPreviouslySetSpan(t *testing.T) {
	ctx := context.Background()
	originalSpan := noop.Span
	expectedSpan := mockSpan{}

	ctx = otel.ContextWithSpan(ctx, originalSpan)
	ctx = otel.ContextWithSpan(ctx, expectedSpan)

	if span := otel.SpanFromContext(ctx); span != expectedSpan {
		t.Errorf("Want: %v, but have: %v", expectedSpan, span)
	}
}

func TestCurrentSpan(t *testing.T) {
	for _, testcase := range []struct {
		name string
		ctx  context.Context
		want otel.Span
	}{
		{
			name: "CurrentSpan() returns a NoopSpan{} from an empty context",
			ctx:  context.Background(),
			want: noop.Span,
		},
		{
			name: "CurrentSpan() returns current span if set",
			ctx:  otel.ContextWithSpan(context.Background(), mockSpan{}),
			want: mockSpan{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			// proto: CurrentSpan(ctx context.Context) otel.Span
			have := otel.SpanFromContext(testcase.ctx)
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

// a duplicate of otel.NoopSpan for testing
type mockSpan struct{}

var _ otel.Span = mockSpan{}

// SpanContext returns an invalid span context.
func (mockSpan) SpanContext() otel.SpanContext {
	return otel.EmptySpanContext()
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
func (mockSpan) SetAttributes(attributes ...label.KeyValue) {
}

// SetAttribute does nothing.
func (mockSpan) SetAttribute(k string, v interface{}) {
}

// End does nothing.
func (mockSpan) End(options ...otel.SpanOption) {
}

// RecordError does nothing.
func (mockSpan) RecordError(ctx context.Context, err error, opts ...otel.ErrorOption) {
}

// Tracer returns noop implementation of Tracer.
func (mockSpan) Tracer() otel.Tracer {
	return noop.Tracer
}

// Event does nothing.
func (mockSpan) AddEvent(ctx context.Context, name string, attrs ...label.KeyValue) {
}

// AddEventWithTimestamp does nothing.
func (mockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...label.KeyValue) {
}
