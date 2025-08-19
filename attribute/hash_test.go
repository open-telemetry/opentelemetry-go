// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"fmt"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

// keyVals is all the KeyValue generators that are used for testing. This is
// not []KeyValue so different keys can be used with the test Values.
var keyVals = []func(string) KeyValue{
	func(k string) KeyValue { return Bool(k, true) },
	func(k string) KeyValue { return Bool(k, false) },
	func(k string) KeyValue { return BoolSlice(k, []bool{false, true}) },
	func(k string) KeyValue { return BoolSlice(k, []bool{true, true, false}) },
	func(k string) KeyValue { return Int(k, -1278) },
	func(k string) KeyValue { return Int(k, 0) }, // Should be different than false above.
	func(k string) KeyValue { return IntSlice(k, []int{3, 23, 21, -8, 0}) },
	func(k string) KeyValue { return IntSlice(k, []int{1}) },
	func(k string) KeyValue { return Int64(k, 1) }, // Should be different from true and []int{1}.
	func(k string) KeyValue { return Int64(k, 29369) },
	func(k string) KeyValue { return Int64Slice(k, []int64{3826, -38, -29, -1}) },
	func(k string) KeyValue { return Int64Slice(k, []int64{8, -328, 29, 0}) },
	func(k string) KeyValue { return Float64(k, -0.3812381) },
	func(k string) KeyValue { return Float64(k, 1e32) },
	func(k string) KeyValue { return Float64Slice(k, []float64{0.1, -3.8, -29., 0.3321}) },
	func(k string) KeyValue { return Float64Slice(k, []float64{-13e8, -32.8, 4., 1e28}) },
	func(k string) KeyValue { return String(k, "foo") },
	func(k string) KeyValue { return String(k, "bar") },
	func(k string) KeyValue { return StringSlice(k, []string{"foo", "bar", "baz"}) },
	func(k string) KeyValue { return StringSlice(k, []string{"[]i1"}) },
}

func TestHashKVsEquality(t *testing.T) {
	type testcase struct {
		hash fnv.Hash
		kvs  []KeyValue
	}

	keys := []string{"k0", "k1"}

	// Test all combinations up to length 3.
	n := len(keyVals)
	result := make([]testcase, 0, 1+len(keys)*(n+(n*n)+(n*n*n)))

	result = append(result, testcase{hashKVs(nil), nil})

	for _, key := range keys {
		for i := 0; i < len(keyVals); i++ {
			kvs := []KeyValue{keyVals[i](key)}
			hash := hashKVs(kvs)
			result = append(result, testcase{hash, kvs})

			for j := 0; j < len(keyVals); j++ {
				kvs := []KeyValue{
					keyVals[i](key),
					keyVals[j](key),
				}
				hash := hashKVs(kvs)
				result = append(result, testcase{hash, kvs})

				for k := 0; k < len(keyVals); k++ {
					kvs := []KeyValue{
						keyVals[i](key),
						keyVals[j](key),
						keyVals[k](key),
					}
					hash := hashKVs(kvs)
					result = append(result, testcase{hash, kvs})
				}
			}
		}
	}

	for i := 0; i < len(result); i++ {
		hI, kvI := result[i].hash, result[i].kvs
		for j := 0; j < len(result); j++ {
			hJ, kvJ := result[j].hash, result[j].kvs
			m := msg{i: i, j: j, hI: hI, hJ: hJ, kvI: kvI, kvJ: kvJ}
			if i == j {
				m.cmp = "=="
				if hI != hJ {
					t.Errorf("hashes not equal: %s", m)
				}
			} else {
				m.cmp = "!="
				if hI == hJ {
					// Do not use testify/assert here. It is slow.
					t.Errorf("hashes equal: %s", m)
				}
			}
		}
	}
}

type msg struct {
	cmp      string
	i, j     int
	hI, hJ   fnv.Hash
	kvI, kvJ []KeyValue
}

func (m msg) String() string {
	return fmt.Sprintf(
		"(%d: %d)%s %s (%d: %d)%s",
		m.i, m.hI, slice(m.kvI), m.cmp, m.j, m.hJ, slice(m.kvJ),
	)
}

func slice(kvs []KeyValue) string {
	if len(kvs) == 0 {
		return "[]"
	}

	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(string(kvs[0].Key))
	_, _ = b.WriteRune(':')
	_, _ = b.WriteString(kvs[0].Value.Emit())
	for _, kv := range kvs[1:] {
		_, _ = b.WriteRune(',')
		_, _ = b.WriteString(string(kv.Key))
		_, _ = b.WriteRune(':')
		_, _ = b.WriteString(kv.Value.Emit())
	}
	_, _ = b.WriteRune(']')
	return b.String()
}

func BenchmarkHashKVs(b *testing.B) {
	attrs := make([]KeyValue, len(keyVals))
	for i := range keyVals {
		attrs[i] = keyVals[i]("k")
	}

	var h fnv.Hash

	b.ResetTimer()
	b.ReportAllocs()
	for range b.N {
		h = hashKVs(attrs)
	}

	_ = h
}
