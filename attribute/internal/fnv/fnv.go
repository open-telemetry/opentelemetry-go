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

// Package fnv1 provides an efficient and allocation free implementation of the
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

// Taken from "hash/fnv".
const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

// Hash is an FNV-1a hash with appropriate hashing functions for methods.
type Hash uint64

// New64 returns a new initialized 64-bit FNV-1a Hash. Its value is laid out in
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

func (h Hash) Bool(val bool) Hash {
	if val {
		return h.Uint64(1)
	}
	return h.Uint64(0)
}

func (h Hash) Float64(val float64) Hash {
	return h.Uint64(math.Float64bits(val))
}

func (h Hash) Int64(val int64) Hash {
	return h.Uint64(uint64(val))
}

func (h Hash) String(val string) Hash {
	v := uint64(h)
	for _, c := range val {
		v ^= uint64(c)
		v *= prime64
	}
	return Hash(v)
}
