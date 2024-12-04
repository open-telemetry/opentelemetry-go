// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"slices"
	"time"
)

// eventAttributesInlineCount is the number of attributes that are efficiently
// stored in an array within an Event.
const eventAttributesInlineCount = 5

// Event represents an OpenTelemetry event.
type Event struct {
	// Ensure forward compatibility by explicitly making this not comparable.
	noCmp [0]func() //nolint: unused  // This is indeed used.

	timestamp         time.Time
	observedTimestamp time.Time
	severity          Severity
	severityText      string
	body              Value

	// The fields below are for optimizing the implementation of Attributes and
	// AddAttributes. This design is borrowed from the slog Record type:
	// https://cs.opensource.google/go/go/+/refs/tags/go1.22.0:src/log/slog/record.go;l=20

	// Allocation optimization: an inline array sized to hold
	// the majority of event calls (based on examination of
	// OpenTelemetry Semantic Conventions for Events).
	// It holds the start of the list of attributes.
	front [eventAttributesInlineCount]KeyValue

	// The number of attributes in front.
	nFront int

	// The list of attributes except for those in front.
	// Invariants:
	//   - len(back) > 0 if nFront == len(front)
	//   - Unused array elements are zero-ed. Used to detect mistakes.
	back []KeyValue
}

// Timestamp returns the time when the event occurred.
func (r *Event) Timestamp() time.Time {
	return r.timestamp
}

// SetTimestamp sets the time when the event occurred.
func (r *Event) SetTimestamp(t time.Time) {
	r.timestamp = t
}

// ObservedTimestamp returns the time when the event was observed.
func (r *Event) ObservedTimestamp() time.Time {
	return r.observedTimestamp
}

// SetObservedTimestamp sets the time when the event was observed.
func (r *Event) SetObservedTimestamp(t time.Time) {
	r.observedTimestamp = t
}

// Severity returns the [Severity] of the event.
func (r *Event) Severity() Severity {
	return r.severity
}

// SetSeverity sets the [Severity] level of the event.
func (r *Event) SetSeverity(level Severity) {
	r.severity = level
}

// SeverityText returns severity (also known as log level) text. This is the
// original string representation of the severity as it is known at the source.
func (r *Event) SeverityText() string {
	return r.severityText
}

// SetSeverityText sets severity (also known as log level) text. This is the
// original string representation of the severity as it is known at the source.
func (r *Event) SetSeverityText(text string) {
	r.severityText = text
}

// Body returns the body of the event.
func (r *Event) Body() Value {
	return r.body
}

// SetBody sets the body of the event.
func (r *Event) SetBody(v Value) {
	r.body = v
}

// WalkAttributes walks all attributes the event holds by calling f for
// each on each [KeyValue] in the [Record]. Iteration stops if f returns false.
func (r *Event) WalkAttributes(f func(KeyValue) bool) {
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

// AddAttributes adds attributes to the event.
func (r *Event) AddAttributes(attrs ...KeyValue) {
	var i int
	for i = 0; i < len(attrs) && r.nFront < len(r.front); i++ {
		a := attrs[i]
		r.front[r.nFront] = a
		r.nFront++
	}

	r.back = slices.Grow(r.back, len(attrs[i:]))
	r.back = append(r.back, attrs[i:]...)
}

// AttributesLen returns the number of attributes in the event.
func (r *Event) AttributesLen() int {
	return r.nFront + len(r.back)
}
