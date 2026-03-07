// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"cmp"
	"fmt"
	"math"
	"reflect"
	"slices"
	"strings"
	"testing"
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
	func(k string) KeyValue { return Slice(k, []Value{BoolValue(true), IntValue(42)}) },
	func(k string) KeyValue {
		return Slice(k, []Value{StringValue("nested"), SliceValue([]Value{IntValue(1), IntValue(2)})})
	},
	func(k string) KeyValue { return Slice(k, []Value{}) },
	func(k string) KeyValue {
		return Slice(k, []Value{Float64Value(3.14), BoolValue(false), StringValue("test")})
	},
}

func TestHashKVsEquality(t *testing.T) {
	type testcase struct {
		hash uint64
		kvs  []KeyValue
	}

	keys := []string{"k0", "k1"}

	// Test all combinations up to length 3.
	n := len(keyVals)
	result := make([]testcase, 0, 1+len(keys)*(n+(n*n)+(n*n*n)))

	result = append(result, testcase{hashKVs(nil), nil})

	for _, key := range keys {
		for i := range keyVals {
			kvs := []KeyValue{keyVals[i](key)}
			hash := hashKVs(kvs)
			result = append(result, testcase{hash, kvs})

			for j := range keyVals {
				kvs := []KeyValue{
					keyVals[i](key),
					keyVals[j](key),
				}
				hash := hashKVs(kvs)
				result = append(result, testcase{hash, kvs})

				for k := range keyVals {
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
	hI, hJ   uint64
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

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		hashKVs(attrs)
	}
}

func FuzzHashKVs(f *testing.F) {
	// Add seed inputs to ensure coverage of edge cases.
	f.Add("", "", "", "", "", "", 0, int64(0), 0.0, false, uint8(0))
	f.Add("key", "value", "🌍", "test", "bool", "float", -1, int64(-1), -1.0, true, uint8(1))
	f.Add("duplicate", "duplicate", "duplicate", "duplicate", "duplicate", "NaN",
		0, int64(0), math.Inf(1), false, uint8(2))

	f.Fuzz(func(t *testing.T, k1, k2, k3, k4, k5, s string, i int, i64 int64, fVal float64, b bool, sliceType uint8) {
		// Test variable number of attributes (0-10).
		numAttrs := len(k1) % 11 // Use key length to determine number of attributes.
		if numAttrs == 0 && k1 == "" {
			// Test empty set.
			h := hashKVs(nil)
			if h == 0 {
				t.Error("hash of empty slice should not be zero")
			}
			return
		}

		var kvs []KeyValue

		// Add basic types.
		if numAttrs > 0 {
			kvs = append(kvs, String(k1, s))
		}
		if numAttrs > 1 {
			kvs = append(kvs, Int(k2, i))
		}
		if numAttrs > 2 {
			kvs = append(kvs, Int64(k3, i64))
		}
		if numAttrs > 3 {
			kvs = append(kvs, Float64(k4, fVal))
		}
		if numAttrs > 4 {
			kvs = append(kvs, Bool(k5, b))
		}

		// Add slice types based on sliceType parameter
		if numAttrs > 5 {
			switch sliceType % 4 {
			case 0:
				// Test BoolSlice with variable length.
				bools := make([]bool, len(s)%5) // 0-4 elements
				for i := range bools {
					bools[i] = (i+len(k1))%2 == 0
				}
				kvs = append(kvs, BoolSlice("boolslice", bools))
			case 1:
				// Test IntSlice with variable length.
				ints := make([]int, len(s)%6) // 0-5 elements
				for i := range ints {
					ints[i] = i + len(k2)
				}
				kvs = append(kvs, IntSlice("intslice", ints))
			case 2:
				// Test Int64Slice with variable length.
				int64s := make([]int64, len(s)%4) // 0-3 elements
				for i := range int64s {
					int64s[i] = int64(i) + i64
				}
				kvs = append(kvs, Int64Slice("int64slice", int64s))
			case 3:
				// Test Float64Slice with variable length and special values.
				float64s := make([]float64, len(s)%5) // 0-4 elements
				for i := range float64s {
					switch i % 4 {
					case 0:
						float64s[i] = fVal
					case 1:
						float64s[i] = math.Inf(1) // +Inf
					case 2:
						float64s[i] = math.Inf(-1) // -Inf
					case 3:
						float64s[i] = math.NaN() // NaN
					}
				}
				kvs = append(kvs, Float64Slice("float64slice", float64s))
			}
		}

		// Add StringSlice.
		if numAttrs > 6 {
			strings := make([]string, len(k1)%4) // 0-3 elements
			for i := range strings {
				strings[i] = fmt.Sprintf("%s_%d", s, i)
			}
			kvs = append(kvs, StringSlice("stringslice", strings))
		}

		// Add Slice (heterogeneous array).
		if numAttrs > 7 {
			sliceLen := len(k2) % 4 // 0-3 elements
			values := make([]Value, sliceLen)
			for i := range values {
				switch i % 5 {
				case 0:
					values[i] = BoolValue((i+len(k1))%2 == 0)
				case 1:
					values[i] = IntValue(i + len(k2))
				case 2:
					values[i] = StringValue(fmt.Sprintf("item_%d", i))
				case 3:
					values[i] = Float64Value(fVal + float64(i))
				case 4:
					// Nested slice
					values[i] = SliceValue([]Value{IntValue(i), BoolValue(true)})
				}
			}
			kvs = append(kvs, Slice("slice", values))
		}

		// Test duplicate keys (should be handled by Set construction).
		if numAttrs > 8 && k1 != "" {
			kvs = append(kvs, String(k1, "duplicate_key_value"))
		}

		// Add more attributes with Unicode keys.
		if numAttrs > 9 {
			kvs = append(kvs, String("🔑", "unicode_key"))
		}
		if numAttrs > 10 {
			kvs = append(kvs, String("empty", ""))
		}

		// Sort to ensure consistent ordering (as Set would do).
		slices.SortFunc(kvs, func(a, b KeyValue) int {
			return cmp.Compare(string(a.Key), string(b.Key))
		})

		// Remove duplicates (as Set will do).
		if len(kvs) > 1 {
			j := 0
			for i := 1; i < len(kvs); i++ {
				if kvs[j].Key != kvs[i].Key {
					j++
					kvs[j] = kvs[i]
				} else {
					// Keep the later value for duplicate keys.
					kvs[j] = kvs[i]
				}
			}
			kvs = kvs[:j+1]
		}

		// Hash the key-value pairs.
		h1 := hashKVs(kvs)
		h2 := hashKVs(kvs) // Should be deterministic

		if h1 != h2 {
			t.Errorf("hash is not deterministic: %d != %d for kvs=%v", h1, h2, kvs)
		}

		if h1 == 0 && len(kvs) > 0 {
			t.Errorf("hash should not be zero for non-empty input: kvs=%v", kvs)
		}

		// Test that different inputs produce different hashes (most of the time).
		// This is a probabilistic test - collisions are possible but rare.
		if len(kvs) > 0 {
			// Modify one value slightly.
			modifiedKvs := make([]KeyValue, len(kvs))
			copy(modifiedKvs, kvs)
			if len(modifiedKvs) > 0 {
				switch modifiedKvs[0].Value.Type() {
				case STRING:
					modifiedKvs[0] = String(string(modifiedKvs[0].Key), modifiedKvs[0].Value.AsString()+"_modified")
				case INT64:
					modifiedKvs[0] = Int64(string(modifiedKvs[0].Key), modifiedKvs[0].Value.AsInt64()+1)
				case BOOL:
					modifiedKvs[0] = Bool(string(modifiedKvs[0].Key), !modifiedKvs[0].Value.AsBool())
				case FLOAT64:
					val := modifiedKvs[0].Value.AsFloat64()
					if !math.IsNaN(val) && !math.IsInf(val, 0) {
						modifiedKvs[0] = Float64(string(modifiedKvs[0].Key), val+1.0)
					}
				case SLICE:
					origSlice := modifiedKvs[0].Value.AsSlice()
					if len(origSlice) > 0 {
						// Modify the first element in the slice.
						newSlice := make([]Value, len(origSlice))
						copy(newSlice, origSlice)
						switch newSlice[0].Type() {
						case INT64:
							newSlice[0] = Int64Value(newSlice[0].AsInt64() + 1)
						case BOOL:
							newSlice[0] = BoolValue(!newSlice[0].AsBool())
						case STRING:
							newSlice[0] = StringValue(newSlice[0].AsString() + "_mod")
						}
						modifiedKvs[0] = Slice(string(modifiedKvs[0].Key), newSlice)
					}
				}

				h3 := hashKVs(modifiedKvs)
				// Note: We don't assert h1 != h3 because hash collisions are theoretically possible
				// but we can log suspicious cases for manual review.
				if h1 == h3 && !reflect.DeepEqual(kvs, modifiedKvs) {
					t.Logf("Potential hash collision detected: original=%v, modified=%v, hash=%d", kvs, modifiedKvs, h1)
				}
			}
		}
	})
}
