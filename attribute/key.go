// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

// Key represents the key part in key-value pairs. It's a string. The
// allowed character set in the key depends on the use of the key.
type Key string

// Bool creates a [KeyValue] with k and a [BOOL] value.
//
// If creating both a key and value at the same time, use the package-level
// Bool function.
func (k Key) Bool(v bool) KeyValue {
	return KeyValue{
		Key:   k,
		Value: BoolValue(v),
	}
}

// BoolSlice creates a [KeyValue] with k and a [BOOLSLICE] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// If creating both a key and value at the same time, use the package-level
// [BoolSlice] function.
func (k Key) BoolSlice(v []bool) KeyValue {
	return KeyValue{
		Key:   k,
		Value: BoolSliceValue(v),
	}
}

// Int creates a [KeyValue] with k and an [INT64] value.
//
// If creating both a key and value at the same time, use the package-level [Int]
// function.
func (k Key) Int(v int) KeyValue {
	return KeyValue{
		Key:   k,
		Value: IntValue(v),
	}
}

// IntSlice creates a [KeyValue] with k and an [INT64SLICE] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// If creating both a key and value at the same time, use the package-level
// [IntSlice] function.
func (k Key) IntSlice(v []int) KeyValue {
	return KeyValue{
		Key:   k,
		Value: IntSliceValue(v),
	}
}

// Int64 creates a [KeyValue] with k and an [INT64] value.
//
// If creating both a key and value at the same time, use the package-level
// [Int64] function.
func (k Key) Int64(v int64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Int64Value(v),
	}
}

// Int64Slice creates a [KeyValue] with k and an [INT64SLICE] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// If creating both a key and value at the same time, use the package-level
// [Int64Slice] function.
func (k Key) Int64Slice(v []int64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Int64SliceValue(v),
	}
}

// Float64 creates a [KeyValue] with k and a [FLOAT64] value.
//
// If creating both a key and value at the same time, use the package-level
// [Float64] function.
func (k Key) Float64(v float64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Float64Value(v),
	}
}

// Float64Slice creates a [KeyValue] with k and a [FLOAT64SLICE] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// If creating both a key and value at the same time, use the package-level
// [Float64Slice] function.
func (k Key) Float64Slice(v []float64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Float64SliceValue(v),
	}
}

// String creates a [KeyValue] with k and a [STRING] value.
//
// If creating both a key and value at the same time, use the package-level
// [String] function.
func (k Key) String(v string) KeyValue {
	return KeyValue{
		Key:   k,
		Value: StringValue(v),
	}
}

// StringSlice creates a [KeyValue] with k and a [STRINGSLICE] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// If creating both a key and value at the same time, use the package-level
// [StringSlice] function.
func (k Key) StringSlice(v []string) KeyValue {
	return KeyValue{
		Key:   k,
		Value: StringSliceValue(v),
	}
}

// ByteSlice creates a [KeyValue] with k and a [BYTESLICE] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// If creating both a key and value at the same time, use the package-level
// [ByteSlice] function.
func (k Key) ByteSlice(v []byte) KeyValue {
	return KeyValue{
		Key:   k,
		Value: ByteSliceValue(v),
	}
}

// Slice creates a [KeyValue] with k and a [SLICE] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// If creating both a key and value at the same time, use the package-level
// [Slice] function.
func (k Key) Slice(v ...Value) KeyValue {
	return KeyValue{
		Key:   k,
		Value: SliceValue(v...),
	}
}

// Map creates a [KeyValue] with k and a [MAP] value.
//
// Note that many observability backends are not optimized to query, index, or
// aggregate complex attribute values. Complex values may also carry
// additional performance overhead. Prefer primitive values when
// possible.
//
// Users should avoid providing duplicate keys; many receivers handle maps
// containing duplicate keys unpredictably.
//
// If creating both a key and value at the same time, use the package-level [Map]
// function.
func (k Key) Map(v ...KeyValue) KeyValue {
	return KeyValue{
		Key:   k,
		Value: MapValue(v...),
	}
}

// Defined reports whether the key is not empty.
func (k Key) Defined() bool {
	return len(k) != 0
}
