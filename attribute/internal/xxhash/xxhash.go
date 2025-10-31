// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package xxhash provides a wrapper around the xxhash library for attribute hashing.
package xxhash // import "go.opentelemetry.io/otel/attribute/internal/xxhash"

import (
	"encoding/binary"
	"math"

	"github.com/cespare/xxhash/v2"
)

// Hash wraps xxhash.Digest to provide the same interface as the FNV implementation.
type Hash struct {
	d *xxhash.Digest
}

// New returns a new initialized xxHash64 hasher.
func New() Hash {
	return Hash{d: xxhash.New()}
}

func (h Hash) Uint64(val uint64) Hash {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], val)
	// errors from Write are always nil for xxhash
	_, _ = h.d.Write(buf[:])
	return h
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
	// errors from WriteString are always nil for xxhash
	_, _ = h.d.WriteString(val)
	return h
}

// Sum64 returns the current hash value.
func (h Hash) Sum64() uint64 {
	return h.d.Sum64()
}
