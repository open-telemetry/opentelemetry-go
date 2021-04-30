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
type SpanSnapshot struct {
	SpanContext trace.SpanContext
	Parent      trace.SpanContext
	SpanKind    trace.SpanKind
	Name        string
	StartTime   time.Time
	// The wall clock time of EndTime will be adjusted to always be offset
	// from StartTime by the duration of the span.
	EndTime    time.Time
	Attributes []attribute.KeyValue
	Events     []Event
	Links      []trace.Link
	Status     Status

	// DroppedAttributeCount contains dropped attributes for the span itself.
	DroppedAttributeCount int
	DroppedEventCount     int
	DroppedLinkCount      int

	// ChildSpanCount holds the number of child span created for this span.
	ChildSpanCount int

	// Resource contains attributes representing an entity that produced this span.
	Resource *resource.Resource

	// InstrumentationLibrary defines the instrumentation library used to
	// provide instrumentation.
	InstrumentationLibrary instrumentation.Library
}
