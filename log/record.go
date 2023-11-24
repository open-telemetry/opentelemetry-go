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
	"go.opentelemetry.io/otel/attribute"
)

var errUnsafeAddAttrs = errors.New("unsafely called AddAttrs on copy of Record made without using Record.Clone")

// Record TODO: comment.
// TODO: Add unit tests.
type Record struct {
	// TODO: comment.
	Timestamp time.Time

	// TODO: comment.
	ObservedTimestamp time.Time

	// TODO: comment.
	Severity Severity

	// TODO: comment.
	SeverityText string

	// TODO: comment.
	Body string

	// The fields below are for optimizing the implementation of
	// Attributes and AddAttributes.

	// Allocation optimization: an inline array sized to hold
	// the majority of log calls (based on examination of open-source
	// code). It holds the start of the list of attributes.
	front [nAttrsInline]attribute.KeyValue

	// The number of attributes in front.
	nFront int

	// The list of attributes except for those in front.
	// Invariants:
	//   - len(back) > 0 iff nFront == len(front)
	//   - Unused array elements are zero. Used to detect mistakes.
	back []attribute.KeyValue
}

const nAttrsInline = 5

// Severity TODO: comment.
type Severity int

// TODO: comment.
const (
	SeverityUndefined Severity = iota
	SeverityTrace
	SeverityTrace2
	SeverityTrace3
	SeverityTrace4
	SeverityDebug
	SeverityDebug2
	SeverityDebug3
	SeverityDebug4
	SeverityInfo
	SeverityInfo2
	SeverityInfo3
	SeverityInfo4
	SeverityWarn
	SeverityWarn2
	SeverityWarn3
	SeverityWarn4
	SeverityError
	SeverityError2
	SeverityError3
	SeverityError4
	SeverityFatal
	SeverityFatal2
	SeverityFatal3
	SeverityFatal4
)

// Attributes calls f on each [attribute.KeyValue] in the [Record].
// Iteration stops if f returns false.
func (r Record) Attributes(f func(attribute.KeyValue) bool) {
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

// AddAttributes appends the given [attribute.KeyValue] to the [Record]'s list of [attribute.KeyValue].
// It omits invalid attributes.
func (r *Record) AddAttributes(attrs ...attribute.KeyValue) {
	var i int
	for i = 0; i < len(attrs) && r.nFront < len(r.front); i++ {
		a := attrs[i]
		if !a.Valid() {
			continue
		}
		r.front[r.nFront] = a
		r.nFront++
	}
	// Check if a copy was modified by slicing past the end
	// and seeing if the Attr there is non-zero.
	if cap(r.back) > len(r.back) {
		end := r.back[:len(r.back)+1][len(r.back)]
		if end.Valid() {
			// Don't panic; copy and muddle through.
			r.back = sliceClip(r.back)
			otel.Handle(errUnsafeAddAttrs)
		}
	}
	ne := countInvalidAttrs(attrs[i:])
	r.back = sliceGrow(r.back, len(attrs[i:])-ne)
	for _, a := range attrs[i:] {
		if a.Valid() {
			r.back = append(r.back, a)
		}
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
func countInvalidAttrs(as []attribute.KeyValue) int {
	n := 0
	for _, a := range as {
		if !a.Valid() {
			n++
		}
	}
	return n
}

// sliceGrow increases the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
//
// This is a copy from https://pkg.go.dev/slices as it is not available in Go 1.20.
func sliceGrow[S ~[]E, E any](s S, n int) S {
	if n < 0 {
		panic("cannot be negative")
	}
	if n -= cap(s) - len(s); n > 0 {
		s = append(s[:cap(s)], make([]E, n)...)[:len(s)]
	}
	return s
}

// sliceClip removes unused capacity from the slice, returning s[:len(s):len(s)].
//
// This is a copy from https://pkg.go.dev/slices as it is not available in Go 1.20.
func sliceClip[S ~[]E, E any](s S) S {
	return s[:len(s):len(s)]
}
