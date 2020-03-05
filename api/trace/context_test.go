package trace_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/core"
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
func (mockSpan) SpanContext() core.SpanContext {
	return core.EmptySpanContext()
}

// IsRecording always returns false for mockSpan.
func (mockSpan) IsRecording() bool {
	return false
}

// SetStatus does nothing.
func (mockSpan) SetStatus(status codes.Code) {
}

// SetName does nothing.
func (mockSpan) SetName(name string) {
}

// SetError does nothing.
func (mockSpan) SetError(v bool) {
}

// SetAttributes does nothing.
func (mockSpan) SetAttributes(attributes ...core.KeyValue) {
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
func (mockSpan) AddEvent(ctx context.Context, name string, attrs ...core.KeyValue) {
}

// AddEventWithTimestamp does nothing.
func (mockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, name string, attrs ...core.KeyValue) {
}
