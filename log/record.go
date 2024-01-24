// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log // import "go.opentelemetry.io/otel/log"

import (
	"errors"
	"time"

	"go.opentelemetry.io/otel"
)

// Record represents a log record.
type Record struct {
	timestamp         time.Time
	observedTimestamp time.Time
	severity          Severity
	severityText      string
	body              Value

	// The fields below are for optimizing the implementation of
	// Attributes and AddAttributes.

	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of attributes.
	front [attributesInlineCount]KeyValue

	// The number of attributes in front.
	nFront int

	// The list of attributes except for those in front.
	// Invariants:
	//   - len(back) > 0 if nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []KeyValue
}

const attributesInlineCount = 5

// Timestamp returns the time when the log record occurred.
func (r Record) Timestamp() time.Time {
	return r.timestamp
}

// SetTimestamp sets the time when the log record occurred.
func (r *Record) SetTimestamp(t time.Time) {
	r.timestamp = t
}

// ObservedTimestamp returns the time when the log record was observed.
// If unset the implementation should set it equal to the current time.
func (r Record) ObservedTimestamp() time.Time {
	return r.observedTimestamp
}

// SetObservedTimestamp sets the time when the log record was observed.
// If unset the implementation should set it equal to the current time.
func (r *Record) SetObservedTimestamp(t time.Time) {
	r.observedTimestamp = t
}

// Severity returns the [Severity] of the log record.
func (r Record) Severity() Severity {
	return r.severity
}

// SetSeverity sets the [Severity] of the log record.
// Use the values defined as constants.
func (r *Record) SetSeverity(s Severity) {
	r.severity = s
}

// SeverityText returns severity (also known as log level) text.
// This is the original string representation of the severity
// as it is known at the source.
func (r Record) SeverityText() string {
	return r.severityText
}

// SetSeverityText sets severity (also known as log level) text.
// This is the original string representation of the severity
// as it is known at the source.
func (r *Record) SetSeverityText(s string) {
	r.severityText = s
}

// Body returns the the body of the log record as a strucutured value.
func (r Record) Body() Value {
	return r.body
}

// SetBody sets the the body of the log record as a strucutured value.
func (r *Record) SetBody(v Value) {
	r.body = v
}

// WalkAttributes calls f on each [KeyValue] in the [Record].
// Iteration stops if f returns false.
func (r Record) WalkAttributes(f func(KeyValue) bool) {
	for i := 0; i < r.nFront; i++ {
		if !f(r.front[i]) {
			return
		}
	}
	for _, a := range r.back {
		if !f(a) {
			return
		}
	}
}

var errUnsafeAddAttrs = errors.New("unsafely called AddAttributes on copy of Record made without using Record.Clone")

// AddAttributes appends the given [attribute.KeyValue] to the [Record]'s
// list of [attribute.KeyValue].
// It omits invalid attributes.
func (r *Record) AddAttributes(attrs ...KeyValue) {
	var i int
	for i = 0; i < len(attrs) && r.nFront < len(r.front); i++ {
		a := attrs[i]
		if a.Invalid() {
			continue
		}
		r.front[r.nFront] = a
		r.nFront++
	}
	// Check if a copy was modified by slicing past the end
	// and seeing if the attribute there is non-zero.
	if cap(r.back) > len(r.back) {
		end := r.back[:len(r.back)+1][len(r.back)]
		if !end.Invalid() {
			// Don't panic; copy and muddle through.
			r.back = sliceClip(r.back)
			otel.Handle(errUnsafeAddAttrs)
		}
	}
	ne := countInvalidAttrs(attrs[i:])
	r.back = sliceGrow(r.back, len(attrs[i:])-ne)
	for _, a := range attrs[i:] {
		if a.Invalid() {
			continue
		}
		r.back = append(r.back, a)
	}
}

// Clone returns a copy of the record with no shared state.
// The original record and the clone can both be modified
// without interfering with each other.
func (r Record) Clone() Record {
	r.back = sliceClip(r.back) // prevent append from mutating shared array
	return r
}

// AttributesLen returns the number of attributes in the Record.
func (r Record) AttributesLen() int {
	return r.nFront + len(r.back)
}

// countInvalidAttrs returns the number of invalid attributes.
func countInvalidAttrs(as []KeyValue) int {
	n := 0
	for _, a := range as {
		if a.Invalid() {
			n++
		}
	}
	return n
}
