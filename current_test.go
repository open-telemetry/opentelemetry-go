package otel_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel"
)

func TestSetCurrentSpanOverridesPreviouslySetSpan(t *testing.T) {
	originalSpan := otel.NoopSpan{}
	expectedSpan := mockSpan{}

	ctx := context.Background()

	ctx = otel.SetCurrentSpan(ctx, originalSpan)
	ctx = otel.SetCurrentSpan(ctx, expectedSpan)

	if span := otel.CurrentSpan(ctx); span != expectedSpan {
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
			want: otel.NoopSpan{},
		},
		{
			name: "CurrentSpan() returns current span if set",
			ctx:  otel.SetCurrentSpan(context.Background(), mockSpan{}),
			want: mockSpan{},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			// proto: CurrentSpan(ctx context.Context) otel.Span
			have := otel.CurrentSpan(testcase.ctx)
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
func (mockSpan) SetStatus(status codes.Code) {
}

// SetName does nothing.
func (mockSpan) SetName(name string) {
}

// SetError does nothing.
func (mockSpan) SetError(v bool) {
}

// SetAttribute does nothing.
func (mockSpan) SetAttribute(attribute otel.KeyValue) {
}

// SetAttributes does nothing.
func (mockSpan) SetAttributes(attributes ...otel.KeyValue) {
}

// End does nothing.
func (mockSpan) End(options ...otel.EndOption) {
}

// Tracer returns noop implementation of Tracer.
func (mockSpan) Tracer() otel.Tracer {
	return otel.NoopTracer{}
}

// Event does nothing.
func (mockSpan) AddEvent(ctx context.Context, msg string, attrs ...otel.KeyValue) {
}

// AddEventWithTimestamp does nothing.
func (mockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...otel.KeyValue) {
}

// AddLink does nothing.
func (mockSpan) AddLink(link otel.Link) {
}

// Link does nothing.
func (mockSpan) Link(sc otel.SpanContext, attrs ...otel.KeyValue) {
}
