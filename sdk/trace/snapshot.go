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

package trace

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// SpanSnapshot is a snapshot of a span which contains all the information
// collected by the span. Its main purpose is exporting completed spans.
// Although SpanSnapshot fields can be accessed and potentially modified,
// SpanSnapshot should be treated as immutable. Changes to the span from which
// the SpanSnapshot was created are NOT reflected in the SpanSnapshot.
//
// TODO: unexport and rename to snapshot.
// TODO: clean project docs of this type after rename.
type SpanSnapshot struct {
	name                   string
	spanContext            trace.SpanContext
	parent                 trace.SpanContext
	spanKind               trace.SpanKind
	startTime              time.Time
	endTime                time.Time
	attributes             []attribute.KeyValue
	events                 []Event
	links                  []trace.Link
	status                 Status
	childSpanCount         int
	droppedAttributeCount  int
	droppedEventCount      int
	droppedLinkCount       int
	resource               *resource.Resource
	instrumentationLibrary instrumentation.Library
}

var _ ReadOnlySpan = SpanSnapshot{}

func (s SpanSnapshot) private() {}

// Name returns the name of the span.
func (s SpanSnapshot) Name() string {
	return s.name
}

// SpanContext returns the unique SpanContext that identifies the span.
func (s SpanSnapshot) SpanContext() trace.SpanContext {
	return s.spanContext
}

// Parent returns the unique SpanContext that identifies the parent of the
// span if one exists. If the span has no parent the returned SpanContext
// will be invalid.
func (s SpanSnapshot) Parent() trace.SpanContext {
	return s.parent
}

// SpanKind returns the role the span plays in a Trace.
func (s SpanSnapshot) SpanKind() trace.SpanKind {
	return s.spanKind
}

// StartTime returns the time the span started recording.
func (s SpanSnapshot) StartTime() time.Time {
	return s.startTime
}

// EndTime returns the time the span stopped recording. It will be zero if
// the span has not ended.
func (s SpanSnapshot) EndTime() time.Time {
	return s.endTime
}

// Attributes returns the defining attributes of the span.
func (s SpanSnapshot) Attributes() []attribute.KeyValue {
	return s.attributes
}

// Links returns all the links the span has to other spans.
func (s SpanSnapshot) Links() []trace.Link {
	return s.links
}

// Events returns all the events that occurred within in the spans
// lifetime.
func (s SpanSnapshot) Events() []Event {
	return s.events
}

// Status returns the spans status.
func (s SpanSnapshot) Status() Status {
	return s.status
}

// InstrumentationLibrary returns information about the instrumentation
// library that created the span.
func (s SpanSnapshot) InstrumentationLibrary() instrumentation.Library {
	return s.instrumentationLibrary
}

// Resource returns information about the entity that produced the span.
func (s SpanSnapshot) Resource() *resource.Resource {
	return s.resource
}

// DroppedAttributes returns the number of attributes dropped by the span
// due to limits being reached.
func (s SpanSnapshot) DroppedAttributes() int {
	return s.droppedAttributeCount
}

// DroppedLinks returns the number of links dropped by the span due to limits
// being reached.
func (s SpanSnapshot) DroppedLinks() int {
	return s.droppedLinkCount
}

// DroppedEvents returns the number of events dropped by the span due to
// limits being reached.
func (s SpanSnapshot) DroppedEvents() int {
	return s.droppedEventCount
}

// ChildSpanCount returns the count of spans that consider the span a
// direct parent.
func (s SpanSnapshot) ChildSpanCount() int {
	return s.childSpanCount
}

// TODO: remove this.
func (s SpanSnapshot) Snapshot() *SpanSnapshot {
	return &s
}
