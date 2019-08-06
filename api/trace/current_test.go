package trace

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/event"
	"go.opentelemetry.io/api/tag"
)

func TestSetCurrentSpan(t *testing.T) {
	for _, testcase := range []struct {
		name string
		span Span
	}{
		{
			name: "set noop span",
			span: noopSpan{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//func SetCurrentSpan(ctx context.Context, span Span) context.Context {
			have := SetCurrentSpan(context.Background(), testcase.span)
			if have.Value(currentSpanKey) != testcase.span {
				t.Errorf("Want: %v, but have: %v", testcase.span, have)
			}
		})
	}
}

func TestCurrentSpan(t *testing.T) {
	for _, testcase := range []struct {
		name string
		ctx  context.Context
		want Span
	}{
		{
			name: "empty context",
			ctx:  context.Background(),
			want: noopSpan{},
		},
		{
			name: "context with mockSpan",
			ctx:  context.WithValue(context.Background(), currentSpanKey, mockSpan{}),
			want: mockSpan{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			// proto: CurrentSpan(ctx context.Context) Span
			have := CurrentSpan(testcase.ctx)
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

// a duplicate of noopSpan for testing
type mockSpan struct{}

// SpanContext returns an invalid span context.
func (mockSpan) SpanContext() core.SpanContext {
	return core.INVALID_SPAN_CONTEXT
}

// IsRecordingEvents always returns false for mockSpan.
func (mockSpan) IsRecordingEvents() bool {
	return false
}

// SetStatus does nothing.
func (mockSpan) SetStatus(status codes.Code) {
}

// SetError does nothing.
func (mockSpan) SetError(v bool) {
}

// SetAttribute does nothing.
func (mockSpan) SetAttribute(attribute core.KeyValue) {
}

// SetAttributes does nothing.
func (mockSpan) SetAttributes(attributes ...core.KeyValue) {
}

// ModifyAttribute does nothing.
func (mockSpan) ModifyAttribute(mutator tag.Mutator) {
}

// ModifyAttributes does nothing.
func (mockSpan) ModifyAttributes(mutators ...tag.Mutator) {
}

// Finish does nothing.
func (mockSpan) Finish() {
}

// Tracer returns noop implementation of Tracer.
func (mockSpan) Tracer() Tracer {
	return noopTracer{}
}

// AddEvent does nothing.
func (mockSpan) AddEvent(ctx context.Context, event event.Event) {
}

// Event does nothing.
func (mockSpan) Event(ctx context.Context, msg string, attrs ...core.KeyValue) {
}
