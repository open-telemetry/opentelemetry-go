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
}

type EndOptions struct {
	EndTime time.Time
}

type EndOption func(*EndOptions)

func WithEndTime(endTime time.Time) EndOption {
	return func(opts *EndOptions) {
		opts.EndTime = endTime
	}
}

type Span interface {
	// Tracer returns tracer used to create this span. Tracer cannot be nil.
	Tracer() Tracer

	// End completes the span. No updates are allowed to span after it
	// ends. The only exception is setting status of the span.
	End(options ...EndOption)

	// AddEvent adds an event to the span.
	AddEvent(ctx context.Context, msg string, attrs ...core.KeyValue)
	// AddEventWithTimestamp adds an event with a custom timestamp
	// to the span.
	AddEventWithTimestamp(ctx context.Context, timestamp time.Time, msg string, attrs ...core.KeyValue)

	// IsRecordingEvents returns true if the span is active and recording events is enabled.
	IsRecordingEvents() bool

	// AddLink adds a link to the span.
	AddLink(link Link)

	// Link creates a link between this span and the other span specified by the SpanContext.
	// It then adds the newly created Link to the span.
	Link(sc core.SpanContext, attrs ...core.KeyValue)

	// SpanContext returns span context of the span. Returned SpanContext is usable
	// even after the span ends.
	SpanContext() core.SpanContext

	// SetStatus sets the status of the span. The status of the span can be updated
	// even after span ends.
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

// Link is used to establish relationship between two spans within the same Trace or
// across different Traces. Few examples of Link usage.
//   1. Batch Processing: A batch of elements may contain elements associated with one
//      or more traces/spans. Since there can only be one parent SpanContext, Link is
//      used to keep reference to SpanContext of all elements in the batch.
//   2. Public Endpoint: A SpanContext in incoming client request on a public endpoint
//      is untrusted from service provider perspective. In such case it is advisable to
//      start a new trace with appropriate sampling decision.
//      However, it is desirable to associate incoming SpanContext to new trace initiated
//      on service provider side so two traces (from Client and from Service Provider) can
//      be correlated.
type Link struct {
	core.SpanContext
	Attributes []core.KeyValue
}

// Start starts a new span using registered global tracer.
func Start(ctx context.Context, name string, opts ...SpanOption) (context.Context, Span) {
	return GlobalTracer().Start(ctx, name, opts...)
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
