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

package otel

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
)

type NoopSpan struct {
}

var _ Span = (*NoopSpan)(nil)

// SpancContext returns an invalid span context.
func (NoopSpan) SpanContext() SpanContext {
	return EmptySpanContext()
}

// IsRecording always returns false for NoopSpan.
func (NoopSpan) IsRecording() bool {
	return false
}

// SetStatus does nothing.
func (NoopSpan) SetStatus(status codes.Code) {
}

// SetError does nothing.
func (NoopSpan) SetError(v bool) {
}

// SetAttribute does nothing.
func (NoopSpan) SetAttribute(attribute KeyValue) {
}

// SetAttributes does nothing.
func (NoopSpan) SetAttributes(attributes ...KeyValue) {
}

// End does nothing.
func (NoopSpan) End(options ...EndOption) {
}

// Tracer returns noop implementation of Tracer.
func (NoopSpan) Tracer() Tracer {
	return NoopTracer{}
}

// AddEvent does nothing.
func (NoopSpan) AddEvent(ctx context.Context, msg string, attrs ...KeyValue) {
}

// AddEventWithTimestamp does nothing.
func (NoopSpan) AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...KeyValue) {
}

// SetName does nothing.
func (NoopSpan) SetName(name string) {
}

// AddLink does nothing.
func (NoopSpan) AddLink(link Link) {
}

// Link does nothing.
func (NoopSpan) Link(sc SpanContext, attrs ...KeyValue) {
}
