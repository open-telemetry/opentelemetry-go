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
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// ReadOnlyLogRecord allows reading information from the data structure underlying a
// log.LogRecord. It is used in places where reading information from a log record is
// necessary but changing the log record isn't necessary or allowed.
//
// Warning: methods may be added to this interface in minor releases.
type ReadOnlyLogRecord interface {
	// SpanContext returns the unique SpanContext that identifies the span this log record is associated with.
	SpanContext() trace.SpanContext
	// Timestamp returns the time the log record started recording.
	Timestamp() time.Time
	// Attributes returns the defining attributes of the log record.
	// The order of the returned attributes is not guaranteed to be stable across invocations.
	Attributes() []attribute.KeyValue
	// InstrumentationScope returns information about the instrumentation
	// scope that created the log record.
	InstrumentationScope() instrumentation.Scope
	// InstrumentationLibrary returns information about the instrumentation
	// library that created the log record.
	// Deprecated: please use InstrumentationScope instead.
	InstrumentationLibrary() instrumentation.Library
	// Resource returns information about the entity that produced the log record.
	Resource() *resource.Resource
	// DroppedAttributes returns the number of attributes dropped by the log record
	// due to limits being reached.
	DroppedAttributes() int

	// A private method to prevent users implementing the
	// interface and so future additions to it will not
	// violate compatibility.
	private()
}

// ReadWriteLogRecord exposes the same methods as log.LogRecord and in addition allows
// reading information from the underlying data structure.
// This interface exposes the union of the methods of log.LogRecord (which is a
// "write-only" log record) and ReadOnlyLogRecord. New methods for writing or reading log record
// information should be added under log.LogRecord or ReadOnlyLogRecord, respectively.
//
// Warning: methods may be added to this interface in minor releases.
type ReadWriteLogRecord interface {
	log.LogRecord
	ReadOnlyLogRecord
}

// recordingLogRecord is an implementation of the OpenTelemetry Log API
// representing the individual component of a trace that is sampled.
type recordingLogRecord struct {
	// mu protects the contents of this log record.
	mu sync.Mutex

	// timestamp is the time at which this log record was started.
	timestamp time.Time

	// spanContext holds the SpanContext of this log record.
	spanContext trace.SpanContext

	// attributes is a collection of user provided key/values. The collection
	// is constrained by a configurable maximum held by the parent
	// LoggerProvider. When additional attributes are added after this maximum
	// is reached these attributes the user is attempting to add are dropped.
	// This dropped number of attributes is tracked and reported in the
	// ReadOnlyLogRecord exported when the span ends.
	attributes        []attribute.KeyValue
	droppedAttributes int

	// logger is the SDK logger that created this log record.
	logger *logger
}

var _ ReadWriteLogRecord = (*recordingLogRecord)(nil)

// SpanContext returns the SpanContext of this log record.
func (s *recordingLogRecord) SpanContext() trace.SpanContext {
	if s == nil {
		return trace.SpanContext{}
	}
	return s.spanContext
}

// SetAttributes sets attributes of this log record.
//
// If a key from attributes already exists the value associated with that key
// will be overwritten with the value contained in attributes.
//
// If adding attributes to the log record would exceed the maximum amount of
// attributes the span is configured to have, the last added attributes will
// be dropped.
func (s *recordingLogRecord) SetAttributes(attributes ...attribute.KeyValue) {
	s.mu.Lock()
	defer s.mu.Unlock()

	limit := s.logger.provider.logRecordLimits.AttributeCountLimit
	if limit == 0 {
		// No attributes allowed.
		s.droppedAttributes += len(attributes)
		return
	}

	// If adding these attributes could exceed the capacity of s perform a
	// de-duplication and truncation while adding to avoid over allocation.
	if limit > 0 && len(s.attributes)+len(attributes) > limit {
		s.addOverCapAttrs(limit, attributes)
		return
	}

	// Otherwise, add without deduplication. When attributes are read they
	// will be deduplicated, optimizing the operation.
	for _, a := range attributes {
		if !a.Valid() {
			// Drop all invalid attributes.
			s.droppedAttributes++
			continue
		}
		a = truncateAttr(s.logger.provider.logRecordLimits.AttributeValueLengthLimit, a)
		s.attributes = append(s.attributes, a)
	}
}

// addOverCapAttrs adds the attributes attrs to the span s while
// de-duplicating the attributes of s and attrs and dropping attributes that
// exceed the limit.
//
// This method assumes s.mu.Lock is held by the caller.
//
// This method should only be called when there is a possibility that adding
// attrs to s will exceed the limit. Otherwise, attrs should be added to s
// without checking for duplicates and all retrieval methods of the attributes
// for s will de-duplicate as needed.
//
// This method assumes limit is a value > 0. The argument should be validated
// by the caller.
func (s *recordingLogRecord) addOverCapAttrs(limit int, attrs []attribute.KeyValue) {
	// In order to not allocate more capacity to s.attributes than needed,
	// prune and truncate this addition of attributes while adding.

	// Do not set a capacity when creating this map. Benchmark testing has
	// showed this to only add unused memory allocations in general use.
	exists := make(map[attribute.Key]int)
	s.dedupeAttrsFromRecord(&exists)

	// Now that s.attributes is deduplicated, adding unique attributes up to
	// the capacity of s will not over allocate s.attributes.
	for _, a := range attrs {
		if !a.Valid() {
			// Drop all invalid attributes.
			s.droppedAttributes++
			continue
		}

		if idx, ok := exists[a.Key]; ok {
			// Perform all updates before dropping, even when at capacity.
			s.attributes[idx] = a
			continue
		}

		if len(s.attributes) >= limit {
			// Do not just drop all of the remaining attributes, make sure
			// updates are checked and performed.
			s.droppedAttributes++
		} else {
			a = truncateAttr(s.logger.provider.logRecordLimits.AttributeValueLengthLimit, a)
			s.attributes = append(s.attributes, a)
			exists[a.Key] = len(s.attributes) - 1
		}
	}
}

// truncateAttr returns a truncated version of attr. Only string and string
// slice attribute values are truncated. String values are truncated to at
// most a length of limit. Each string slice value is truncated in this fashion
// (the slice length itself is unaffected).
//
// No truncation is perfromed for a negative limit.
func truncateAttr(limit int, attr attribute.KeyValue) attribute.KeyValue {
	if limit < 0 {
		return attr
	}
	switch attr.Value.Type() {
	case attribute.STRING:
		if v := attr.Value.AsString(); len(v) > limit {
			return attr.Key.String(safeTruncate(v, limit))
		}
	case attribute.STRINGSLICE:
		v := attr.Value.AsStringSlice()
		for i := range v {
			if len(v[i]) > limit {
				v[i] = safeTruncate(v[i], limit)
			}
		}
		return attr.Key.StringSlice(v)
	}
	return attr
}

// safeTruncate truncates the string and guarantees valid UTF-8 is returned.
func safeTruncate(input string, limit int) string {
	if trunc, ok := safeTruncateValidUTF8(input, limit); ok {
		return trunc
	}
	trunc, _ := safeTruncateValidUTF8(strings.ToValidUTF8(input, ""), limit)
	return trunc
}

// safeTruncateValidUTF8 returns a copy of the input string safely truncated to
// limit. The truncation is ensured to occur at the bounds of complete UTF-8
// characters. If invalid encoding of UTF-8 is encountered, input is returned
// with false, otherwise, the truncated input will be returned with true.
func safeTruncateValidUTF8(input string, limit int) (string, bool) {
	for cnt := 0; cnt <= limit; {
		r, size := utf8.DecodeRuneInString(input[cnt:])
		if r == utf8.RuneError {
			return input, false
		}

		if cnt+size > limit {
			return input[:cnt], true
		}
		cnt += size
	}
	return input, true
}

// Timestamp returns the time this log record started.
func (s *recordingLogRecord) Timestamp() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.timestamp
}

// Attributes returns the attributes of this span.
//
// The order of the returned attributes is not guaranteed to be stable.
func (s *recordingLogRecord) Attributes() []attribute.KeyValue {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dedupeAttrs()
	return s.attributes
}

// dedupeAttrs deduplicates the attributes of s to fit capacity.
//
// This method assumes s.mu.Lock is held by the caller.
func (s *recordingLogRecord) dedupeAttrs() {
	// Do not set a capacity when creating this map. Benchmark testing has
	// showed this to only add unused memory allocations in general use.
	exists := make(map[attribute.Key]int)
	s.dedupeAttrsFromRecord(&exists)
}

// dedupeAttrsFromRecord deduplicates the attributes of s to fit capacity
// using record as the record of unique attribute keys to their index.
//
// This method assumes s.mu.Lock is held by the caller.
func (s *recordingLogRecord) dedupeAttrsFromRecord(record *map[attribute.Key]int) {
	// Use the fact that slices share the same backing array.
	unique := s.attributes[:0]
	for _, a := range s.attributes {
		if idx, ok := (*record)[a.Key]; ok {
			unique[idx] = a
		} else {
			unique = append(unique, a)
			(*record)[a.Key] = len(unique) - 1
		}
	}
	// s.attributes have element types of attribute.KeyValue. These types are
	// not pointers and they themselves do not contain pointer fields,
	// therefore the duplicate values do not need to be zeroed for them to be
	// garbage collected.
	s.attributes = unique
}

// InstrumentationScope returns the instrumentation.Scope associated with
// the Tracer that created this span.
func (s *recordingLogRecord) InstrumentationScope() instrumentation.Scope {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.logger.instrumentationScope
}

// InstrumentationLibrary returns the instrumentation.Library associated with
// the Tracer that created this span.
func (s *recordingLogRecord) InstrumentationLibrary() instrumentation.Library {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.logger.instrumentationScope
}

// Resource returns the Resource associated with the Tracer that created this
// span.
func (s *recordingLogRecord) Resource() *resource.Resource {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.logger.provider.resource
}

// DroppedAttributes returns the number of attributes dropped by the span
// due to limits being reached.
func (s *recordingLogRecord) DroppedAttributes() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.droppedAttributes
}

// LoggerProvider returns a log.LoggerProvider that can be used to generate
// additional Spans on the same telemetry pipeline as the current Span.
func (s *recordingLogRecord) LoggerProvider() log.LoggerProvider {
	return s.logger.provider
}

// snapshot creates a read-only copy of the current state of the span.
func (s *recordingLogRecord) snapshot() ReadOnlyLogRecord {
	var sd snapshot
	s.mu.Lock()
	defer s.mu.Unlock()

	sd.instrumentationScope = s.logger.instrumentationScope
	sd.resource = s.logger.provider.resource
	sd.spanContext = s.spanContext
	sd.timestamp = s.timestamp

	if len(s.attributes) > 0 {
		s.dedupeAttrs()
		sd.attributes = s.attributes
	}
	sd.droppedAttributeCount = s.droppedAttributes
	return &sd
}

func (*recordingLogRecord) private() {}
