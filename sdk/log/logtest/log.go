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

package logtest // import "go.opentelemetry.io/otel/sdk/trace/logtest"

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	logsdk "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// LogRecordStubs is a slice of LogRecordStub use for testing an SDK.
type LogRecordStubs []LogRecordStub

// LogRecordStubsFromReadOnlyLogRecords returns LogRecordStubs populated from ro.
func LogRecordStubsFromReadOnlyLogRecords(ro []logsdk.ReadOnlyLogRecord) LogRecordStubs {
	if len(ro) == 0 {
		return nil
	}

	s := make(LogRecordStubs, 0, len(ro))
	for _, r := range ro {
		s = append(s, SpanStubFromReadOnlySpan(r))
	}

	return s
}

// Snapshots returns s as a slice of ReadOnlySpans.
func (s LogRecordStubs) Snapshots() []logsdk.ReadOnlyLogRecord {
	if len(s) == 0 {
		return nil
	}

	ro := make([]logsdk.ReadOnlyLogRecord, len(s))
	for i := 0; i < len(s); i++ {
		ro[i] = s[i].Snapshot()
	}
	return ro
}

// LogRecordStub is a stand-in for a Span.
type LogRecordStub struct {
	SpanContext            trace.SpanContext
	Timestamp              time.Time
	EndTime                time.Time
	Attributes             []attribute.KeyValue
	DroppedAttributes      int
	Resource               *resource.Resource
	InstrumentationLibrary instrumentation.Library
}

// SpanStubFromReadOnlySpan returns a LogRecordStub populated from ro.
func SpanStubFromReadOnlySpan(ro logsdk.ReadOnlyLogRecord) LogRecordStub {
	if ro == nil {
		return LogRecordStub{}
	}

	return LogRecordStub{
		SpanContext:            ro.SpanContext(),
		Timestamp:              ro.Timestamp(),
		Attributes:             ro.Attributes(),
		DroppedAttributes:      ro.DroppedAttributes(),
		Resource:               ro.Resource(),
		InstrumentationLibrary: ro.InstrumentationScope(),
	}
}

// Snapshot returns a read-only copy of the LogRecordStub.
func (s LogRecordStub) Snapshot() logsdk.ReadOnlyLogRecord {
	return logRecordSnapshot{
		spanContext:          s.SpanContext,
		timestamp:            s.Timestamp,
		attributes:           s.Attributes,
		droppedAttributes:    s.DroppedAttributes,
		resource:             s.Resource,
		instrumentationScope: s.InstrumentationLibrary,
	}
}

type logRecordSnapshot struct {
	// Embed the interface to implement the private method.
	logsdk.ReadOnlyLogRecord

	spanContext          trace.SpanContext
	timestamp            time.Time
	attributes           []attribute.KeyValue
	droppedAttributes    int
	resource             *resource.Resource
	instrumentationScope instrumentation.Scope
}

func (s logRecordSnapshot) SpanContext() trace.SpanContext   { return s.spanContext }
func (s logRecordSnapshot) Timestamp() time.Time             { return s.timestamp }
func (s logRecordSnapshot) Attributes() []attribute.KeyValue { return s.attributes }
func (s logRecordSnapshot) DroppedAttributes() int           { return s.droppedAttributes }
func (s logRecordSnapshot) Resource() *resource.Resource     { return s.resource }
func (s logRecordSnapshot) InstrumentationScope() instrumentation.Scope {
	return s.instrumentationScope
}
func (s logRecordSnapshot) InstrumentationLibrary() instrumentation.Library {
	return s.instrumentationScope
}
