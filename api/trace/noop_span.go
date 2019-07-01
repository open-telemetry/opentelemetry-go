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

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/event"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/api/stats"
)

type noopSpan struct {
}

var _ Span = (*noopSpan)(nil)
var _ stats.Interface = (*noopSpan)(nil)
var _ scope.Mutable = (*noopSpan)(nil)

// SpancContext returns an invalid span context.
func (sp *noopSpan) SpanContext() core.SpanContext {
	return core.INVALID_SPAN_CONTEXT
}

// IsRecordingEvents always returns false for noopSpan.
func (sp *noopSpan) IsRecordingEvents() bool {
	return false
}

// SetStatus does nothing.
func (sp *noopSpan) SetStatus(status codes.Code) {
	return
}

// ScopeID returns and empty ScopeID.
func (sp *noopSpan) ScopeID() core.ScopeID {
	return core.ScopeID{}
}

// SetError does nothing.
func (sp *noopSpan) SetError(v bool) {
	return
}

// SetAttribute does nothing.
func (sp *noopSpan) SetAttribute(attribute core.KeyValue) {
	return
}

// SetAttributes does nothing.
func (sp *noopSpan) SetAttributes(attributes ...core.KeyValue) {
	return
}

// ModifyAttribute does nothing.
func (sp *noopSpan) ModifyAttribute(mutator core.Mutator) {
	return
}

// ModifyAttributes does nothing.
func (sp *noopSpan) ModifyAttributes(mutators ...core.Mutator) {
	return
}

// Finish does nothing.
func (sp *noopSpan) Finish() {
	return
}

// Tracer returns noop implementation of Tracer.
func (sp *noopSpan) Tracer() Tracer {
	return t
}

// AddEvent does nothing.
func (sp *noopSpan) AddEvent(ctx context.Context, event event.Event) {
}

// Record does nothing.
func (sp *noopSpan) Record(ctx context.Context, m ...core.Measurement) {
}

// RecordSingle does nothing.
func (sp *noopSpan) RecordSingle(ctx context.Context, m core.Measurement) {
}
