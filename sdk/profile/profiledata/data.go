// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profiledata

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// ResourceProfiles is a collection of ScopeProfiles and the Resource that
// produced them.
type ResourceProfiles struct {
	// Resource describes the entity that collected the profiles.
	Resource *resource.Resource
	// ScopeProfiles holds profiles grouped by instrumentation scope.
	ScopeProfiles []ScopeProfiles
}

// ScopeProfiles is a collection of Profiles produced by a single
// instrumentation scope.
type ScopeProfiles struct {
	// Scope identifies the instrumentation scope that produced Profiles.
	Scope instrumentation.Scope
	// Profiles holds the profiles collected under Scope.
	Profiles []Profile
}

// ValueType describes the kind and unit of a measured value or sampling period.
type ValueType struct {
	// Type is the measurement kind (e.g. "cpu", "inuse_space", "delay").
	Type string
	// Unit is the measurement unit (e.g. "nanoseconds", "bytes", "count").
	Unit string
}

// Attribute is a key-value pair with an optional UCUM unit, mirroring OTLP
// KeyValueAndUnit.
type Attribute struct {
	// Key is the attribute name.
	Key attribute.Key
	// Value is the attribute value.
	Value attribute.Value
	// Unit is an optional unit in UCUM format (empty if implicit or unset).
	Unit string
}

// Link associates a profile sample with a distributed trace span.
type Link struct {
	// TraceID is the trace identifier of the linked span.
	TraceID trace.TraceID
	// SpanID is the span identifier of the linked span.
	SpanID trace.SpanID
}

// IsValid reports whether both TraceID and SpanID are valid.
func (l Link) IsValid() bool {
	return l.TraceID.IsValid() && l.SpanID.IsValid()
}

// Function describes a source-code function referenced by a stack frame line.
type Function struct {
	// Name is the human-readable function name.
	Name string
	// SystemName is the system or mangled function name, if distinct.
	SystemName string
	// Filename is the source file containing the function.
	Filename string
	// StartLine is the 1-based starting line number in Filename (0 if unset).
	StartLine int64
}

// Line represents a single source code location within a frame.
type Line struct {
	// Function is the function containing this line.
	Function Function
	// Line is the 1-based source line number (0 if unset).
	Line int64
	// Column is the 1-based source column number (0 if unset).
	Column int64
}

// Mapping describes an executable or shared library mapped into memory.
type Mapping struct {
	// MemoryStart is the start address at which the binary is loaded.
	MemoryStart uint64
	// MemoryLimit is the exclusive upper bound of the mapped address range.
	MemoryLimit uint64
	// FileOffset is the file offset corresponding to MemoryStart.
	FileOffset uint64
	// Filename is the path to the binary or virtual mapping (e.g. "[vdso]").
	Filename string
	// Attributes holds mapping attributes, such as build IDs.
	Attributes []Attribute
}

// Location represents a single stack frame, optionally containing multiple
// inlined source lines.
type Location struct {
	// Mapping is the binary mapping containing Address, if known.
	Mapping *Mapping
	// Address is the instruction address for this frame, if available.
	Address uint64
	// Lines lists source locations for this frame, innermost inlined function
	// first.
	Lines []Line
	// Attributes holds frame-level attributes (e.g. profile.frame.type).
	Attributes []Attribute
}

// Sample records observations for a single stack trace and attribute context.
type Sample struct {
	// Stack is the call stack for this sample, leaf frame first.
	Stack []*Location
	// Values holds measured values in the unit described by Profile.SampleType.
	Values []int64
	// Timestamps holds optional per-observation UTC timestamps. When both
	// Values and Timestamps are non-empty, they must have equal length.
	Timestamps []time.Time
	// Attributes holds sample-level attributes with unique keys.
	Attributes []Attribute
	// Link optionally links this sample to an active trace span.
	Link Link
}

// Profile represents a single collected profile with one SampleType.
type Profile struct {
	// SampleType describes the type and unit of Sample.Values.
	SampleType ValueType
	// Samples holds the recorded samples.
	Samples []Sample
	// Time is the start time of the profile collection.
	Time time.Time
	// Duration is the collection window duration (0 for instantaneous profiles).
	Duration time.Duration
	// PeriodType describes the type and unit of Period.
	PeriodType ValueType
	// Period is the sampling distance in PeriodType units.
	Period int64
	// ProfileID is an optional 16-byte globally unique profile identifier.
	ProfileID [16]byte
	// DroppedAttributesCount is the number of dropped profile-level attributes.
	DroppedAttributesCount uint32
	// OriginalPayloadFormat identifies the format of OriginalPayload (e.g. "pprof").
	OriginalPayloadFormat string
	// OriginalPayload holds optional raw source payload bytes.
	OriginalPayload []byte
	// Attributes holds profile-level attributes with unique keys.
	Attributes []Attribute
}
