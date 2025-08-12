// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fnv provides an efficient and allocation free implementation of the
// FNV-1a, non-cryptographic hash functions created by Glenn Fowler, Landon
// Curt Noll, and Phong Vo. See
// https://en.wikipedia.org/wiki/Fowler-Noll-Vo_hash_function.
//
// This implementation is provided as an alternative to "hash/fnv". The
// built-in implementation requires two allocations per Write for a string (one
// for the hash pointer and the other to convert a string to a []byte). This
// implementation is more efficientient and does not require any allocations.
package fnv // import "go.opentelemetry.io/otel/attribute/internal/fnv"

import (
	"math"
)

// Taken from "hash/fnv". Verified at:
//
//   - https://datatracker.ietf.org/doc/html/draft-eastlake-fnv-17.html
//   - http://www.isthe.com/chongo/tech/comp/fnv/index.html#FNV-param
const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// Hash is an FNV-1a hash with appropriate hashing functions for methods.
type Hash uint64

// New returns a new initialized 64-bit FNV-1a Hash. Its value is laid out in
// big-endian byte order.
func New() Hash {
	return offset64
}

func (h Hash) Uint64(val uint64) Hash {
	v := uint64(h)
	v = (v ^ ((val >> 56) & 0xFF)) * prime64
	v = (v ^ ((val >> 48) & 0xFF)) * prime64
	v = (v ^ ((val >> 40) & 0xFF)) * prime64
	v = (v ^ ((val >> 32) & 0xFF)) * prime64
	v = (v ^ ((val >> 24) & 0xFF)) * prime64
	v = (v ^ ((val >> 16) & 0xFF)) * prime64
	v = (v ^ ((val >> 8) & 0xFF)) * prime64
	v = (v ^ ((val >> 0) & 0xFF)) * prime64
	return Hash(v)
}

func (h Hash) Bool(val bool) Hash { // nolint:revive  // val is not a flag.
	if val {
		return h.Uint64(1)
	}
	return h.Uint64(0)
}

func (h Hash) Float64(val float64) Hash {
	return h.Uint64(math.Float64bits(val))
}

func (h Hash) Int64(val int64) Hash {
	return h.Uint64(uint64(val)) // nolint:gosec // overflow doesn't matter since we are hashing.
}

func (h Hash) String(val string) Hash {
	v := uint64(h)
	for _, c := range val {
		v ^= uint64(c)
		v *= prime64
	}
	return Hash(v)
}
