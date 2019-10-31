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

package trace

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/otel/api/core"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

// MockSpan is a mock span used in association with MockTracer for testing purpose only.
type MockSpan struct {
	sc     core.SpanContext
	tracer apitrace.Tracer
}

var _ apitrace.Span = (*MockSpan)(nil)

// SpanContext returns associated core.SpanContext. If the receiver is nil it returns
// an empty core.SpanContext
func (ms *MockSpan) SpanContext() core.SpanContext {
	if ms == nil {
		core.EmptySpanContext()
	}
	return ms.sc
}

// IsRecording always returns false for MockSpan.
func (ms *MockSpan) IsRecording() bool {
	return false
}

// SetStatus does nothing.
func (ms *MockSpan) SetStatus(status codes.Code) {
}

// SetError does nothing.
func (ms *MockSpan) SetError(v bool) {
}

// SetAttribute does nothing.
func (ms *MockSpan) SetAttribute(attribute core.KeyValue) {
}

// SetAttributes does nothing.
func (ms *MockSpan) SetAttributes(attributes ...core.KeyValue) {
}

// End does nothing.
func (ms *MockSpan) End(options ...apitrace.EndOption) {
}

// SetName does nothing.
func (ms *MockSpan) SetName(name string) {
}

// Tracer returns MockTracer implementation of Tracer.
func (ms *MockSpan) Tracer() apitrace.Tracer {
	return ms.tracer
}

// AddEvent does nothing.
func (ms *MockSpan) AddEvent(ctx context.Context, msg string, attrs ...core.KeyValue) {
}

// AddEvent does nothing.
func (ms *MockSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...core.KeyValue) {
}

// AddLink does nothing.
func (ms *MockSpan) AddLink(link apitrace.Link) {
}

// Link does nothing.
func (ms *MockSpan) Link(sc core.SpanContext, attrs ...core.KeyValue) {
}
