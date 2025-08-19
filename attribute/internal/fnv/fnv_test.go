// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fnv

import (
	"encoding/binary"
	"hash/fnv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringHashCorrectness(t *testing.T) {
	input := []string{"", "a", "ab", "abc", "世界"}

	refH := fnv.New64a()
	for _, in := range input {
		h := New()
		got := h.String(in)

		refH.Reset()
		n, err := refH.Write([]byte(in))
		require.NoError(t, err)
		require.Equalf(t, len(in), n, "wrote only %d out of %d bytes", n, len(in))
		want := refH.Sum64()

		assert.Equal(t, want, uint64(got), in)
	}
}

func TestUint64HashCorrectness(t *testing.T) {
	input := []uint64{0, 10, 312984238623, 1024}

	buf := make([]byte, 8)
	refH := fnv.New64a()
	for _, in := range input {
		h := New()
		got := h.Uint64(in)

		refH.Reset()
		binary.BigEndian.PutUint64(buf, in)
		n, err := refH.Write(buf)
		require.NoError(t, err)
		require.Equalf(t, 8, n, "wrote only %d out of 8 bytes", n)
		want := refH.Sum64()

		assert.Equal(t, want, uint64(got), in)
	}
}

func TestIntegrity(t *testing.T) {
	data := []byte{'1', '2', 3, 4, 5, 6, 7, 8, 9, 10}
	h0 := New()
	want := h0.String(string(data))

	h1 := New()
	got := h1.String(string(data[:2]))
	num := binary.BigEndian.Uint64(data[2:])
	got = got.Uint64(num)

	assert.Equal(t, want, got)
}

var result Hash

func BenchmarkStringKB(b *testing.B) {
	b.SetBytes(1024)
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i)
	}
	s := string(data)
	h := New()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result = h.String(s)
	}
}

func BenchmarkUint64KB(b *testing.B) {
	b.SetBytes(8)
	i := uint64(192386739218721)
	h := New()

	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		result = h.Uint64(i)
	}
}
