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

package sdk

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/tag"
	"go.opentelemetry.io/api/trace"
	"go.opentelemetry.io/experimental/streaming/exporter"
)

type span struct {
	tracer  *tracer
	initial exporter.ScopeID
}

// SpancContext returns span context of the span. Return SpanContext is usable
// even after the span is finished.
func (sp *span) SpanContext() core.SpanContext {
	return sp.initial.SpanContext
}

// IsRecordingEvents returns true is the span is active and recording events is enabled.
func (sp *span) IsRecordingEvents() bool {
	return false
}

// SetStatus sets the status of the span.
func (sp *span) SetStatus(status codes.Code) {
	sp.tracer.exporter.Record(exporter.Event{
		Type:   exporter.SET_STATUS,
		Scope:  sp.ScopeID(),
		Status: status,
	})
}

func (sp *span) ScopeID() exporter.ScopeID {
	return sp.initial
}

func (sp *span) SetAttribute(attribute core.KeyValue) {
	sp.tracer.exporter.Record(exporter.Event{
		Type:      exporter.MODIFY_ATTR,
		Scope:     sp.ScopeID(),
		Attribute: attribute,
	})
}

func (sp *span) SetAttributes(attributes ...core.KeyValue) {
	sp.tracer.exporter.Record(exporter.Event{
		Type:       exporter.MODIFY_ATTR,
		Scope:      sp.ScopeID(),
		Attributes: attributes,
	})
}

func (sp *span) ModifyAttribute(mutator tag.Mutator) {
	sp.tracer.exporter.Record(exporter.Event{
		Type:    exporter.MODIFY_ATTR,
		Scope:   sp.ScopeID(),
		Mutator: mutator,
	})
}

func (sp *span) ModifyAttributes(mutators ...tag.Mutator) {
	sp.tracer.exporter.Record(exporter.Event{
		Type:     exporter.MODIFY_ATTR,
		Scope:    sp.ScopeID(),
		Mutators: mutators,
	})
}

func (sp *span) Finish(options ...trace.FinishOption) {
	recovered := recover()
	opts := trace.FinishOptions{}
	for _, opt := range options {
		opt(&opts)
	}
	sp.tracer.exporter.Record(exporter.Event{
		Time:      opts.FinishTime,
		Type:      exporter.FINISH_SPAN,
		Scope:     sp.ScopeID(),
		Recovered: recovered,
	})
	if recovered != nil {
		panic(recovered)
	}
}

func (sp *span) Tracer() trace.Tracer {
	return sp.tracer
}

func (sp *span) AddEvent(ctx context.Context, msg string, attrs ...core.KeyValue) {
	sp.addEventWithTime(ctx, time.Time{}, msg, attrs...)
}

func (sp *span) addEventWithTime(ctx context.Context, timestamp time.Time, msg string, attrs ...core.KeyValue) {
	sp.tracer.exporter.Record(exporter.Event{
		Time:       timestamp,
		Type:       exporter.ADD_EVENT,
		String:     msg,
		Attributes: attrs,
		Context:    ctx,
	})
}

func (sp *span) SetName(name string) {
	sp.tracer.exporter.Record(exporter.Event{
		Type:   exporter.SET_NAME,
		String: name,
	})
}

func (sp *span) AddLink(link trace.Link) {
}

func (sp *span) Link(sc core.SpanContext, attrs ...core.KeyValue) {
}
