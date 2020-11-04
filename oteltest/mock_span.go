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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
)

// MockSpan is a mock span used in association with MockTracer for testing purpose only.
type MockSpan struct {
	StatusMsg string
	Name      string
	Status    codes.Code
	sc        otel.SpanContext
	tracer    otel.Tracer
}

var _ otel.Span = (*MockSpan)(nil)

// SpanContext returns associated label.SpanContext. If the receiver is nil it returns
// an empty label.SpanContext
func (ms *MockSpan) SpanContext() otel.SpanContext {
	if ms == nil {
		return otel.SpanContext{}
	}
	return ms.sc
}

// IsRecording always returns false for MockSpan.
func (ms *MockSpan) IsRecording() bool { return false }

// SetStatus does nothing.
func (ms *MockSpan) SetStatus(status codes.Code, msg string) {
	ms.Status = status
	ms.StatusMsg = msg
}

// SetError does nothing.
func (ms *MockSpan) SetError(v bool) {}

// SetAttributes does nothing.
func (ms *MockSpan) SetAttributes(attributes ...label.KeyValue) {}

// End does nothing.
func (ms *MockSpan) End(options ...otel.SpanOption) {}

// RecordError does nothing.
func (ms *MockSpan) RecordError(err error, opts ...otel.EventOption) {}

// SetName sets the span name.
func (ms *MockSpan) SetName(name string) { ms.Name = name }

// Tracer returns MockTracer implementation of Tracer.
func (ms *MockSpan) Tracer() otel.Tracer { return ms.tracer }

// AddEvent does nothing.
func (ms *MockSpan) AddEvent(string, ...otel.EventOption) {}
