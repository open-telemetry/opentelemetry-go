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

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// snapshot is an record of a logRecords state at a particular checkpointed time.
// It is used as a read-only representation of that state.
type snapshot struct {
	name                  string
	spanContext           trace.SpanContext
	timestamp             time.Time
	attributes            []attribute.KeyValue
	droppedAttributeCount int
	resource              *resource.Resource
	instrumentationScope  instrumentation.Scope
}

var _ ReadOnlyLogRecord = snapshot{}

func (s snapshot) private() {}

// Name returns the name of the span.
func (s snapshot) Name() string {
	return s.name
}

// SpanContext returns the unique SpanContext that identifies the span.
func (s snapshot) SpanContext() trace.SpanContext {
	return s.spanContext
}

// Timestamp returns the time the span started recording.
func (s snapshot) Timestamp() time.Time {
	return s.timestamp
}

// Attributes returns the defining attributes of the span.
func (s snapshot) Attributes() []attribute.KeyValue {
	return s.attributes
}

// InstrumentationScope returns information about the instrumentation
// scope that created the span.
func (s snapshot) InstrumentationScope() instrumentation.Scope {
	return s.instrumentationScope
}

// InstrumentationLibrary returns information about the instrumentation
// library that created the span.
func (s snapshot) InstrumentationLibrary() instrumentation.Library {
	return s.instrumentationScope
}

// Resource returns information about the entity that produced the span.
func (s snapshot) Resource() *resource.Resource {
	return s.resource
}

// DroppedAttributes returns the number of attributes dropped by the span
// due to limits being reached.
func (s snapshot) DroppedAttributes() int {
	return s.droppedAttributeCount
}
