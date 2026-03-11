// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package xxhash

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegrity(t *testing.T) {
	data := []byte{'1', '2', 3, 4, 5, 6, 7, 8, 9, 10}
	h0 := New()
	want := h0.String(string(data))

	h1 := New()
	got := h1.String(string(data[:2]))
	num := binary.LittleEndian.Uint64(data[2:])
	got = got.Uint64(num)

	assert.Equal(t, want.Sum64(), got.Sum64())
}

func TestNew(t *testing.T) {
	h1 := New()
	h2 := New()

	// Test that the underlying digest is properly initialized.
	if h1.Sum64() != h2.Sum64() {
		t.Errorf("New() should return consistent initial value: %d != %d", h1.Sum64(), h2.Sum64())
	}
}

func TestUint64(t *testing.T) {
	h1 := New().Uint64(42)
	h2 := New().Uint64(42)
	if h1.Sum64() != h2.Sum64() {
		t.Errorf("Uint64() should be deterministic: %d != %d", h1.Sum64(), h2.Sum64())
	}

	h3 := New().Uint64(43)
	if h1.Sum64() == h3.Sum64() {
		t.Errorf("Different inputs should produce different hashes: %d == %d", h1.Sum64(), h3.Sum64())
	}
}

func TestBool(t *testing.T) {
	h1 := New().Bool(true)
	h2 := New().Bool(true)
	if h1.Sum64() != h2.Sum64() {
		t.Errorf("Bool() should be deterministic: %d != %d", h1.Sum64(), h2.Sum64())
	}

	h3 := New().Bool(false)
	if h1.Sum64() == h3.Sum64() {
		t.Errorf("Different bool values should produce different hashes: %d == %d", h1.Sum64(), h3.Sum64())
	}
}

func TestFloat64(t *testing.T) {
	h1 := New().Float64(3.14)
	h2 := New().Float64(3.14)
	if h1.Sum64() != h2.Sum64() {
		t.Errorf("Float64() should be deterministic: %d != %d", h1.Sum64(), h2.Sum64())
	}

	h3 := New().Float64(2.71)
	if h1.Sum64() == h3.Sum64() {
		t.Errorf("Different float values should produce different hashes: %d == %d", h1.Sum64(), h3.Sum64())
	}
}

func TestInt64(t *testing.T) {
	h1 := New().Int64(42)
	h2 := New().Int64(42)
	if h1.Sum64() != h2.Sum64() {
		t.Errorf("Int64() should be deterministic: %d != %d", h1.Sum64(), h2.Sum64())
	}

	h3 := New().Int64(43)
	if h1.Sum64() == h3.Sum64() {
		t.Errorf("Different int64 values should produce different hashes: %d == %d", h1.Sum64(), h3.Sum64())
	}
}

func TestString(t *testing.T) {
	h1 := New().String("hello")
	h2 := New().String("hello")
	if h1.Sum64() != h2.Sum64() {
		t.Errorf("String() should be deterministic: %d != %d", h1.Sum64(), h2.Sum64())
	}

	h3 := New().String("world")
	if h1.Sum64() == h3.Sum64() {
		t.Errorf("Different strings should produce different hashes: %d == %d", h1.Sum64(), h3.Sum64())
	}
}

func TestChaining(t *testing.T) {
	// Test that methods can be chained and produce different results
	h1 := New().String("key").Uint64(42).Bool(true)
	h2 := New().String("key").Uint64(42).Bool(true)
	h3 := New().String("key").Uint64(43).Bool(true)

	if h1.Sum64() != h2.Sum64() {
		t.Errorf("Chained operations should be deterministic: %d != %d", h1.Sum64(), h2.Sum64())
	}

	if h1.Sum64() == h3.Sum64() {
		t.Errorf("Different chained operations should produce different hashes: %d == %d", h1.Sum64(), h3.Sum64())
	}
}

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
	for b.Loop() {
		h.String(s)
	}
}

func BenchmarkUint64KB(b *testing.B) {
	b.SetBytes(8)
	i := uint64(192386739218721)
	h := New()

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		h.Uint64(i)
	}
}

func BenchmarkUint64(b *testing.B) {
	h := New()

	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		h = h.Uint64(uint64(i))
	}
}

func BenchmarkString(b *testing.B) {
	h := New()
	str := "benchmark_string_value"

	b.ReportAllocs()
	for b.Loop() {
		h = h.String(str)
	}
}

func BenchmarkBool(b *testing.B) {
	h := New()

	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		h = h.Bool(i%2 == 0)
	}
}

func BenchmarkFloat64(b *testing.B) {
	h := New()

	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		h = h.Float64(float64(i) * 3.14159)
	}
}

func BenchmarkInt64(b *testing.B) {
	h := New()

	b.ReportAllocs()
	for i := 0; b.Loop(); i++ {
		h = h.Int64(int64(i))
	}
}

func BenchmarkSum64(b *testing.B) {
	h := New().String("key").Uint64(42).Bool(true)

	b.ReportAllocs()
	for b.Loop() {
		_ = h.Sum64()
	}
}
