// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package log // import "go.opentelemetry.io/otel/log"

// sliceGrow increases the slice's capacity, if necessary, to guarantee space
// for another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow panics.
//
// This is a copy from https://pkg.go.dev/slices as it is not available in Go 1.20.
func sliceGrow[S ~[]E, E any](s S, n int) S {
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

// sliceEqualFunc reports whether two slices are equal using an equality
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// increasing index order, and the comparison stops at the first index
// for which eq returns false.
//
// This is a copy from https://pkg.go.dev/slices as it is not available in Go 1.20.
func sliceEqualFunc[S1 ~[]E1, S2 ~[]E2, E1, E2 any](s1 S1, s2 S2, eq func(E1, E2) bool) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, v1 := range s1 {
		v2 := s2[i]
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}
