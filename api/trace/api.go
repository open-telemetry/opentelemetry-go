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

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/event"
	"go.opentelemetry.io/api/tag"
)

type Tracer interface {
	// Start a span.
	Start(context.Context, string, ...SpanOption) (context.Context, Span)

	// WithSpan wraps the execution of the function body with a span.
	// It starts a new span and sets it as an active span in the context.
	// It then executes the body. It closes the span before returning the execution result.
	WithSpan(
		ctx context.Context,
		operation string,
		body func(ctx context.Context) error,
	) error

	// TODO: Do we need WithService and WithComponent?
	// TODO: Can we make these helpers (based on WithResources)?
	WithService(name string) Tracer
	WithComponent(name string) Tracer

	// WithResources attaches resource attributes to the Tracer.
	WithResources(res ...core.KeyValue) Tracer

	// Note: see https://github.com/opentracing/opentracing-go/issues/127
	Inject(context.Context, Span, Injector)
}

type Span interface {
	// Tracer returns tracer used to create this span. Tracer cannot be nil.
	Tracer() Tracer

	// Finish completes the span. No updates are allowed to span after it
	// finishes. The only exception is setting status of the span.
	Finish()

	// AddEvent adds an event to the span.
	AddEvent(ctx context.Context, event event.Event)
	// AddEvent records an event to the span.
	Event(ctx context.Context, msg string, attrs ...core.KeyValue)

	// IsRecordingEvents returns true if the span is active and recording events is enabled.
	IsRecordingEvents() bool

	// SpancContext returns span context of the span. Return SpanContext is usable
	// even after the span is finished.
	SpanContext() core.SpanContext

	// SetStatus sets the status of the span. The status of the span can be updated
	// even after span is finished.
	SetStatus(codes.Code)

	// SetName sets the name of the span.
	SetName(name string)

	// Set span attributes
	SetAttribute(core.KeyValue)
	SetAttributes(...core.KeyValue)

	// Modify and delete span attributes
	ModifyAttribute(tag.Mutator)
	ModifyAttributes(...tag.Mutator)
}

type Injector interface {
	// Inject serializes span context and tag.Map and inserts them in to
	// carrier associated with the injector. For example in case of http request,
	// span context could added to the request (carrier) as W3C Trace context header.
	Inject(core.SpanContext, tag.Map)
}

// SpanOption apply changes to SpanOptions.
type SpanOption func(*SpanOptions)

// SpanOptions provides options to set properties of span at the time of starting
// a new span.
type SpanOptions struct {
	Attributes  []core.KeyValue
	StartTime   time.Time
	Reference   Reference
	RecordEvent bool
}

// Reference is used to establish relationship between newly created span and the
// other span. The other span could be related as a parent or linked or any other
// future relationship type.
type Reference struct {
	core.SpanContext
	RelationshipType
}

type RelationshipType int

const (
	ChildOfRelationship RelationshipType = iota
	FollowsFromRelationship
)

// Start starts a new span using registered global tracer.
func Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	return GlobalTracer().Start(ctx, name, opts...)
}

// Inject is convenient function to inject current span context using injector.
// Injector is expected to serialize span context and inject it in to a carrier.
// An example of a carrier is http request.
func Inject(ctx context.Context, injector Injector) {
	span := CurrentSpan(ctx)
	if span == nil {
		return
	}

	span.Tracer().Inject(ctx, span, injector)
}

// WithStartTime sets the start time of the span to provided time t, when it is started.
// In absensce of this option, wall clock time is used as start time.
// This option is typically used when starting of the span is delayed.
func WithStartTime(t time.Time) SpanOption {
	return func(o *SpanOptions) {
		o.StartTime = t
	}
}

// WithAttributes sets attributes to span. These attributes provides additional
// data about the span.
func WithAttributes(attrs ...core.KeyValue) SpanOption {
	return func(o *SpanOptions) {
		o.Attributes = attrs
	}
}

// WithRecordEvents enables recording of the events while the span is active.
// In the absence of this option, RecordEvent is set to false, disabling any recording of
// the events.
func WithRecordEvents() SpanOption {
	return func(o *SpanOptions) {
		o.RecordEvent = true
	}
}

// ChildOf. TODO: do we need this?.
func ChildOf(sc core.SpanContext) SpanOption {
	return func(o *SpanOptions) {
		o.Reference = Reference{
			SpanContext:      sc,
			RelationshipType: ChildOfRelationship,
		}
	}
}

// FollowsFrom. TODO: do we need this?.
func FollowsFrom(sc core.SpanContext) SpanOption {
	return func(o *SpanOptions) {
		o.Reference = Reference{
			SpanContext:      sc,
			RelationshipType: FollowsFromRelationship,
		}
	}
}
