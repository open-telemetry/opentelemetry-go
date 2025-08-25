// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package pool provides common pools for semantic convention implementations.
package pool // import "go.opentelemetry.io/otel/semconv/internal/pool"

import (
	"slices"
	"sync"

	"go.opentelemetry.io/otel/attribute"
)

// MaxCapAttrsSmall is the maximum capacity of a slice that will be taken from
// the small pool. Slices with a capacity greater than this will be taken from
// the medium pool.
const MaxCapAttrsSmall = 8

var (
	// attrSlice is a pool for slices of attribute.KeyValue that are less than
	// or equal to the cutoffSmall.
	attrSliceSmall = &sync.Pool{
		New: func() any {
			var s []attribute.KeyValue
			// Return pointer to a slice instead of a slice to avoid an
			// allocation when putting it back in the pool.
			return &s
		},
	}

	attrSliceMedium = &sync.Pool{
		New: func() any {
			// Instead of allocating with make, we return a nil slice here.
			// GetAtrrSlice will grow it to the needed size in a single step.
			var s []attribute.KeyValue
			// Return pointer to a slice instead of a slice to avoid an
			// allocation when putting it back in the pool.
			return &s
		},
	}
)

// GetAttrSlice returns a pointer to a slice of [attribute.KeyValue]s from a
// common pool. The slice will have length zero and a capacity of at least n.
//
// The returned slice must be returned to the pool with [PutAttrSlice] when no
// longer needed.
func GetAttrSlice(n int) *[]attribute.KeyValue {
	var a *[]attribute.KeyValue
	if n <= MaxCapAttrsSmall {
		a = attrSliceSmall.Get().(*[]attribute.KeyValue)
	} else {
		a = attrSliceMedium.Get().(*[]attribute.KeyValue)
	}
	*a = (*a)[:0] // reset slice length to zero.

	if n > cap(*a) {
		*a = slices.Grow(*a, n)
	}

	return a
}

// PutAttrSlice returns a to the common pool.
func PutAttrSlice(a *[]attribute.KeyValue) {
	if cap(*a) <= MaxCapAttrsSmall {
		attrSliceSmall.Put(a)
	} else {
		attrSliceMedium.Put(a)
	}
}
