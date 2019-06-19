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
	"github.com/open-telemetry/opentelemetry-go/api/log"
	"github.com/open-telemetry/opentelemetry-go/api/stats"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

// SpancContext returns span context of the span. Return SpanContext is usable
// even after the span is finished.
func (sp *span) SpanContext() core.SpanContext {
	if sp == nil {
		return core.INVALID_SPAN_CONTEXT
	}
	return sp.spanContext
}

// IsRecordingEvents returns true is the span is active and recording events is enabled.
func (sp *span) IsRecordingEvents() bool {
	return false
}

// SetStatus sets the status of the span.
func (sp *span) SetStatus(status codes.Code) {
	if sp == nil {
		return
	}
	sid := sp.ScopeID()

	observer.Record(observer.Event{
		Type:     observer.SET_STATUS,
		Scope:    sid,
		Sequence: sid.EventID,
		Status:   status,
	})
	sp.status = status
}

func (sp *span) ScopeID() core.ScopeID {
	if sp == nil {
		return core.ScopeID{}
	}
	sp.lock.Lock()
	sid := core.ScopeID{
		EventID:     sp.eventID,
		SpanContext: sp.spanContext,
	}
	sp.lock.Unlock()
	return sid
}

func (sp *span) updateScope() (core.ScopeID, core.EventID) {
	next := observer.NextEventID()

	sp.lock.Lock()
	sid := core.ScopeID{
		EventID:     sp.eventID,
		SpanContext: sp.spanContext,
	}
	sp.eventID = next
	sp.lock.Unlock()

	return sid, next
}

func (sp *span) SetError(v bool) {
	sp.SetAttribute(ErrorKey.Bool(v))
}

func (sp *span) SetAttribute(attribute core.KeyValue) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:      observer.MODIFY_ATTR,
		Scope:     sid,
		Sequence:  next,
		Attribute: attribute,
	})
}

func (sp *span) SetAttributes(attributes ...core.KeyValue) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:       observer.MODIFY_ATTR,
		Scope:      sid,
		Sequence:   next,
		Attributes: attributes,
	})
}

func (sp *span) ModifyAttribute(mutator core.Mutator) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:     observer.MODIFY_ATTR,
		Scope:    sid,
		Sequence: next,
		Mutator:  mutator,
	})
}

func (sp *span) ModifyAttributes(mutators ...core.Mutator) {
	if sp == nil {
		return
	}

	sid, next := sp.updateScope()

	observer.Record(observer.Event{
		Type:     observer.MODIFY_ATTR,
		Scope:    sid,
		Sequence: next,
		Mutators: mutators,
	})
}

func (sp *span) Finish() {
	if sp == nil {
		return
	}
	recovered := recover()
	sp.finishOnce.Do(func() {
		observer.Record(observer.Event{
			Type:      observer.FINISH_SPAN,
			Scope:     sp.ScopeID(),
			Recovered: recovered,
		})
	})
	if recovered != nil {
		panic(recovered)
	}
}

func (sp *span) Tracer() Tracer {
	return sp.tracer
}

func (sp *span) Log(ctx context.Context, msg string, args ...core.KeyValue) {
	log.With(sp).Log(ctx, msg, args...)
}

func (sp *span) Logf(ctx context.Context, fmt string, args ...interface{}) {
	log.With(sp).Logf(ctx, fmt, args...)
}

func (sp *span) Record(ctx context.Context, m ...core.Measurement) {
	stats.With(sp).Record(ctx, m...)
}

func (sp *span) RecordSingle(ctx context.Context, m core.Measurement) {
	stats.With(sp).RecordSingle(ctx, m)
}
