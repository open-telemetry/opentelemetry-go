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
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/event"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/api/stats"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
)

type Tracer interface {
	Start(context.Context, string, ...SpanOption) (context.Context, Span)

	// WithSpan wraps the execution of the function body with a span.
	// It starts a new span and sets it as an active span in the context.
	// It then executes the body. It closes the span before returning the execution result.
	// TODO: Should it restore the previous span?
	WithSpan(
		ctx context.Context,
		operation string,
		body func(ctx context.Context) error,
	) error

	// TODO: Do we need WithService and WithComponent?
	WithService(name string) Tracer
	WithComponent(name string) Tracer

	// WithResources attaches resource attributes to the Tracer.
	WithResources(res ...core.KeyValue) Tracer

	// Note: see https://github.com/opentracing/opentracing-go/issues/127
	Inject(context.Context, Span, Injector)

	// ScopeID returns the resource scope of this tracer.
	scope.Scope
}

type Span interface {
	scope.Mutable

	stats.Interface

	// Tracer returns tracer used to create this span. Tracer cannot be nil.
	Tracer() Tracer

	// Finish completes the span. No updates are allowed to span after it
	// finishes. The only exception is setting status of the span.
	Finish()

	// AddEvent adds an event to the span.
	AddEvent(ctx context.Context, event event.Event)

	// IsRecordingEvents returns true if the span is active and recording events is enabled.
	IsRecordingEvents() bool

	// SpancContext returns span context of the span. Return SpanContext is usable
	// even after the span is finished.
	SpanContext() core.SpanContext

	// SetStatus sets the status of the span. The status of the span can be updated
	// even after span is finished.
	SetStatus(codes.Code)
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

var (
	// The process global tracer could have process-wide resource
	// tags applied directly, or we can have a SetGlobal tracer to
	// install a default tracer w/ resources.
	global atomic.Value

	// TODO: create NOOP Tracer and register it instead of creating empty tracer here.
	nt = &noopTracer{}
)

const (
	ChildOfRelationship RelationshipType = iota
	FollowsFromRelationship
)

// GlobalTracer return tracer registered with global registry.
// If no tracer is registered then an instance of noop Tracer is returned.
func GlobalTracer() Tracer {
	if t := global.Load(); t != nil {
		return t.(Tracer)
	}
	return nt
}

// SetGlobalTracer sets provided tracer as a global tracer.
func SetGlobalTracer(t Tracer) {
	global.Store(t)
}

// Start starts a new span using registered global tracer.
func Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	return GlobalTracer().Start(ctx, name, opts...)
}

// Active returns current span from the context.
func Active(ctx context.Context) Span {
	span, _ := scope.Active(ctx).(Span)
	return span
}

// Inject is convenient function to inject current span context using injector.
// Injector is expected to serialize span context and inject it in to a carrier.
// An example of a carrier is http request.
func Inject(ctx context.Context, injector Injector) {
	span := Active(ctx)
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
