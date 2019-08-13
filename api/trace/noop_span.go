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

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/event"
	"go.opentelemetry.io/api/tag"
)

type noopSpan struct {
}

var _ Span = (*noopSpan)(nil)

// SpancContext returns an invalid span context.
func (noopSpan) SpanContext() core.SpanContext {
	return core.EmptySpanContext()
}

// IsRecordingEvents always returns false for noopSpan.
func (noopSpan) IsRecordingEvents() bool {
	return false
}

// SetStatus does nothing.
func (noopSpan) SetStatus(status codes.Code) {
}

// SetError does nothing.
func (noopSpan) SetError(v bool) {
}

// SetAttribute does nothing.
func (noopSpan) SetAttribute(attribute core.KeyValue) {
}

// SetAttributes does nothing.
func (noopSpan) SetAttributes(attributes ...core.KeyValue) {
}

// ModifyAttribute does nothing.
func (noopSpan) ModifyAttribute(mutator tag.Mutator) {
}

// ModifyAttributes does nothing.
func (noopSpan) ModifyAttributes(mutators ...tag.Mutator) {
}

// Finish does nothing.
func (noopSpan) Finish() {
}

// Tracer returns noop implementation of Tracer.
func (noopSpan) Tracer() Tracer {
	return noopTracer{}
}

// AddEvent does nothing.
func (noopSpan) AddEvent(ctx context.Context, event event.Event) {
}

// Event does nothing.
func (noopSpan) Event(ctx context.Context, msg string, attrs ...core.KeyValue) {
}
