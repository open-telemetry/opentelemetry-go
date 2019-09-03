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

type DefaultSpan struct {
	sc core.SpanContext
}

var _ Span = (*DefaultSpan)(nil)

// SpancContext returns an invalid span context.
func (ds *DefaultSpan) SpanContext() core.SpanContext {
	if ds == nil {
		core.EmptySpanContext()
	}
	return ds.sc
}

// IsRecordingEvents always returns false for DefaultSpan.
func (ds *DefaultSpan) IsRecordingEvents() bool {
	return false
}

// SetStatus does nothing.
func (ds *DefaultSpan) SetStatus(status codes.Code) {
}

// SetError does nothing.
func (ds *DefaultSpan) SetError(v bool) {
}

// SetAttribute does nothing.
func (ds *DefaultSpan) SetAttribute(attribute core.KeyValue) {
}

// SetAttributes does nothing.
func (ds *DefaultSpan) SetAttributes(attributes ...core.KeyValue) {
}

// ModifyAttribute does nothing.
func (ds *DefaultSpan) ModifyAttribute(mutator tag.Mutator) {
}

// ModifyAttributes does nothing.
func (ds *DefaultSpan) ModifyAttributes(mutators ...tag.Mutator) {
}

// Finish does nothing.
func (ds *DefaultSpan) Finish() {
}

// Tracer returns noop implementation of Tracer.
func (ds *DefaultSpan) Tracer() Tracer {
	return NoopTracer{}
}

// AddEvent does nothing.
func (ds *DefaultSpan) AddEvent(ctx context.Context, event event.Event) {
}

// Event does nothing.
func (ds *DefaultSpan) Event(ctx context.Context, msg string, attrs ...core.KeyValue) {
}
