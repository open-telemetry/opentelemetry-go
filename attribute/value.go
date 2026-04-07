// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	attribute "go.opentelemetry.io/otel/attribute/internal"
)

//go:generate stringer -type=Type

// Type describes the type of the data Value holds.
type Type int // nolint: revive  // redefines builtin Type.

// Value represents the value part in key-value pairs.
//
// Note that the zero value is a valid empty value.
type Value struct {
	vtype    Type
	numeric  uint64
	stringly string
	slice    any
}

const (
	// EMPTY is used for a Value with no value set.
	EMPTY Type = iota
	// BOOL is a boolean Type Value.
	BOOL
	// INT64 is a 64-bit signed integral Type Value.
	INT64
	// FLOAT64 is a 64-bit floating point Type Value.
	FLOAT64
	// STRING is a string Type Value.
	STRING
	// BOOLSLICE is a slice of booleans Type Value.
	BOOLSLICE
	// INT64SLICE is a slice of 64-bit signed integral numbers Type Value.
	INT64SLICE
	// FLOAT64SLICE is a slice of 64-bit floating point numbers Type Value.
	FLOAT64SLICE
	// STRINGSLICE is a slice of strings Type Value.
	STRINGSLICE
	// BYTESLICE is a slice of bytes Type Value.
	BYTESLICE
	// INVALID is used for a Value with no value set.
	//
	// Deprecated: Use EMPTY instead as an empty value is a valid value.
	INVALID = EMPTY
)

// BoolValue creates a BOOL Value.
func BoolValue(v bool) Value {
	return Value{
		vtype:   BOOL,
		numeric: boolToRaw(v),
	}
}

// BoolSliceValue creates a BOOLSLICE Value.
func BoolSliceValue(v []bool) Value {
	return Value{vtype: BOOLSLICE, slice: attribute.SliceValue(v)}
}

// IntValue creates an INT64 Value.
func IntValue(v int) Value {
	return Int64Value(int64(v))
}

// IntSliceValue creates an INT64SLICE Value.
func IntSliceValue(v []int) Value {
	val := Value{vtype: INT64SLICE}

	// Avoid the common tiny-slice cases from allocating a new slice.
	switch len(v) {
	case 0:
		val.slice = [0]int64{}
	case 1:
		val.slice = [1]int64{int64(v[0])}
	case 2:
		val.slice = [2]int64{int64(v[0]), int64(v[1])}
	case 3:
		val.slice = [3]int64{int64(v[0]), int64(v[1]), int64(v[2])}
	default:
		// Fallback to a new slice for larger slices.
		cp := make([]int64, len(v))
		for i, val := range v {
			cp[i] = int64(val)
		}
		val.slice = attribute.SliceValue(cp)
	}

	return val
}

// Int64Value creates an INT64 Value.
func Int64Value(v int64) Value {
	return Value{
		vtype:   INT64,
		numeric: int64ToRaw(v),
	}
}

// Int64SliceValue creates an INT64SLICE Value.
func Int64SliceValue(v []int64) Value {
	return Value{vtype: INT64SLICE, slice: attribute.SliceValue(v)}
}

// Float64Value creates a FLOAT64 Value.
func Float64Value(v float64) Value {
	return Value{
		vtype:   FLOAT64,
		numeric: float64ToRaw(v),
	}
}

// Float64SliceValue creates a FLOAT64SLICE Value.
func Float64SliceValue(v []float64) Value {
	return Value{vtype: FLOAT64SLICE, slice: attribute.SliceValue(v)}
}

// StringValue creates a STRING Value.
func StringValue(v string) Value {
	return Value{
		vtype:    STRING,
		stringly: v,
	}
}

// StringSliceValue creates a STRINGSLICE Value.
func StringSliceValue(v []string) Value {
	return Value{vtype: STRINGSLICE, slice: attribute.SliceValue(v)}
}

// ByteSliceValue creates a BYTESLICE Value.
func ByteSliceValue(v []byte) Value {
	return Value{
		vtype:    BYTESLICE,
		stringly: string(v),
	}
}

// Type returns a type of the Value.
func (v Value) Type() Type {
	return v.vtype
}

// AsBool returns the bool value. Make sure that the Value's type is
// BOOL.
func (v Value) AsBool() bool {
	return rawToBool(v.numeric)
}

// AsBoolSlice returns the []bool value. Make sure that the Value's type is
// BOOLSLICE.
func (v Value) AsBoolSlice() []bool {
	if v.vtype != BOOLSLICE {
		return nil
	}
	return v.asBoolSlice()
}

func (v Value) asBoolSlice() []bool {
	return attribute.AsSlice[bool](v.slice)
}

// AsInt64 returns the int64 value. Make sure that the Value's type is
// INT64.
func (v Value) AsInt64() int64 {
	return rawToInt64(v.numeric)
}

// AsInt64Slice returns the []int64 value. Make sure that the Value's type is
// INT64SLICE.
func (v Value) AsInt64Slice() []int64 {
	if v.vtype != INT64SLICE {
		return nil
	}
	return v.asInt64Slice()
}

func (v Value) asInt64Slice() []int64 {
	return attribute.AsSlice[int64](v.slice)
}

// AsFloat64 returns the float64 value. Make sure that the Value's
// type is FLOAT64.
func (v Value) AsFloat64() float64 {
	return rawToFloat64(v.numeric)
}

// AsFloat64Slice returns the []float64 value. Make sure that the Value's type is
// FLOAT64SLICE.
func (v Value) AsFloat64Slice() []float64 {
	if v.vtype != FLOAT64SLICE {
		return nil
	}
	return v.asFloat64Slice()
}

func (v Value) asFloat64Slice() []float64 {
	return attribute.AsSlice[float64](v.slice)
}

// AsString returns the string value. Make sure that the Value's type
// is STRING.
func (v Value) AsString() string {
	return v.stringly
}

// AsStringSlice returns the []string value. Make sure that the Value's type is
// STRINGSLICE.
func (v Value) AsStringSlice() []string {
	if v.vtype != STRINGSLICE {
		return nil
	}
	return v.asStringSlice()
}

func (v Value) asStringSlice() []string {
	return attribute.AsSlice[string](v.slice)
}

// AsByteSlice returns the bytes value. Make sure that the Value's type
// is BYTESLICE.
func (v Value) AsByteSlice() []byte {
	if v.vtype != BYTESLICE {
		return nil
	}
	return v.asByteSlice()
}

func (v Value) asByteSlice() []byte {
	return []byte(v.stringly)
}

type unknownValueType struct{}

const lowerhex = "0123456789abcdef"

// AsInterface returns Value's data as any.
func (v Value) AsInterface() any {
	switch v.Type() {
	case BOOL:
		return v.AsBool()
	case BOOLSLICE:
		return v.asBoolSlice()
	case INT64:
		return v.AsInt64()
	case INT64SLICE:
		return v.asInt64Slice()
	case FLOAT64:
		return v.AsFloat64()
	case FLOAT64SLICE:
		return v.asFloat64Slice()
	case STRING:
		return v.stringly
	case STRINGSLICE:
		return v.asStringSlice()
	case BYTESLICE:
		return v.asByteSlice()
	case EMPTY:
		return nil
	}
	return unknownValueType{}
}

// String returns a string representation of Value using the
// [OpenTelemetry AnyValue representation for non-OTLP protocols] rules.
//
// Strings are returned as-is without JSON quoting, booleans and integers use
// JSON literals, floating-point values use JSON numbers except that NaN and
// +/-Infinity are rendered as NaN, Infinity, and -Infinity, byte slices are
// base64-encoded, empty values are the empty string, and slices are encoded as
// JSON arrays. Floating-point special values inside arrays are encoded as JSON
// strings.
//
// [OpenTelemetry AnyValue representation for non-OTLP protocols]: https://opentelemetry.io/docs/specs/otel/common/#anyvalue-representation-for-non-otlp-protocols
func (v Value) String() string {
	switch v.Type() {
	case BOOL:
		return strconv.FormatBool(v.AsBool())
	case BOOLSLICE:
		return formatBoolSliceValue(v.slice)
	case INT64:
		return strconv.FormatInt(v.AsInt64(), 10)
	case INT64SLICE:
		return formatInt64SliceValue(v.slice)
	case FLOAT64:
		return formatFloat64(v.AsFloat64())
	case FLOAT64SLICE:
		return formatFloat64SliceValue(v.slice)
	case STRING:
		return v.stringly
	case STRINGSLICE:
		return formatStringSliceValue(v.slice)
	case BYTESLICE:
		return base64.StdEncoding.EncodeToString(v.asByteSlice())
	case EMPTY:
		return ""
	default:
		return "unknown"
	}
}

// Emit returns a string representation of Value's data.
func (v Value) Emit() string {
	switch v.Type() {
	case BOOLSLICE:
		return fmt.Sprint(v.asBoolSlice())
	case BOOL:
		return strconv.FormatBool(v.AsBool())
	case INT64SLICE:
		j, err := json.Marshal(v.asInt64Slice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asInt64Slice())
		}
		return string(j)
	case INT64:
		return strconv.FormatInt(v.AsInt64(), 10)
	case FLOAT64SLICE:
		j, err := json.Marshal(v.asFloat64Slice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asFloat64Slice())
		}
		return string(j)
	case FLOAT64:
		return fmt.Sprint(v.AsFloat64())
	case STRINGSLICE:
		j, err := json.Marshal(v.asStringSlice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asStringSlice())
		}
		return string(j)
	case STRING:
		return v.stringly
	case BYTESLICE:
		return base64.StdEncoding.EncodeToString(v.asByteSlice())
	case EMPTY:
		return ""
	default:
		return "unknown"
	}
}

func formatBoolSlice(vals []bool) string {
	if len(vals) == 0 {
		return "[]"
	}

	var b strings.Builder
	b.Grow(2 + len(vals)*5)
	_ = b.WriteByte('[')
	for i, val := range vals {
		if i > 0 {
			_ = b.WriteByte(',')
		}
		if val {
			_, _ = b.WriteString("true")
		} else {
			_, _ = b.WriteString("false")
		}
	}
	_ = b.WriteByte(']')
	return b.String()
}

func formatBoolSliceValue(v any) string {
	switch vals := v.(type) {
	case [0]bool:
		return "[]"
	case [1]bool:
		return formatBoolSlice(vals[:])
	case [2]bool:
		return formatBoolSlice(vals[:])
	case [3]bool:
		return formatBoolSlice(vals[:])
	default:
		return formatBoolSliceReflect(v)
	}
}

func formatInt64Slice(vals []int64) string {
	if len(vals) == 0 {
		return "[]"
	}

	var b strings.Builder
	b.Grow(2 + len(vals)*20)
	_ = b.WriteByte('[')

	var scratch [20]byte
	for i, val := range vals {
		if i > 0 {
			_ = b.WriteByte(',')
		}
		out := strconv.AppendInt(scratch[:0], val, 10)
		_, _ = b.Write(out)
	}

	_ = b.WriteByte(']')
	return b.String()
}

func formatInt64SliceValue(v any) string {
	switch vals := v.(type) {
	case [0]int64:
		return "[]"
	case [1]int64:
		return formatInt64Slice(vals[:])
	case [2]int64:
		return formatInt64Slice(vals[:])
	case [3]int64:
		return formatInt64Slice(vals[:])
	default:
		return formatInt64SliceReflect(v)
	}
}

func formatFloat64(v float64) string {
	switch {
	case math.IsNaN(v):
		return "NaN"
	case math.IsInf(v, 1):
		return "Infinity"
	case math.IsInf(v, -1):
		return "-Infinity"
	default:
		return strconv.FormatFloat(v, 'g', -1, 64)
	}
}

func formatFloat64Slice(vals []float64) string {
	if len(vals) == 0 {
		return "[]"
	}

	var b strings.Builder
	b.Grow(2 + len(vals)*24)
	_ = b.WriteByte('[')

	var scratch [24]byte
	for i, val := range vals {
		if i > 0 {
			_ = b.WriteByte(',')
		}

		switch {
		case math.IsNaN(val):
			_, _ = b.WriteString(`"NaN"`)
		case math.IsInf(val, 1):
			_, _ = b.WriteString(`"Infinity"`)
		case math.IsInf(val, -1):
			_, _ = b.WriteString(`"-Infinity"`)
		default:
			out := strconv.AppendFloat(scratch[:0], val, 'g', -1, 64)
			_, _ = b.Write(out)
		}
	}

	_ = b.WriteByte(']')
	return b.String()
}

func formatFloat64SliceValue(v any) string {
	switch vals := v.(type) {
	case [0]float64:
		return "[]"
	case [1]float64:
		return formatFloat64Slice(vals[:])
	case [2]float64:
		return formatFloat64Slice(vals[:])
	case [3]float64:
		return formatFloat64Slice(vals[:])
	default:
		return formatFloat64SliceReflect(v)
	}
}

func formatStringSlice(vals []string) string {
	if len(vals) == 0 {
		return "[]"
	}

	size := 2
	for _, val := range vals {
		size += len(val) + 2
	}
	size += len(vals) - 1

	var b strings.Builder
	b.Grow(size)
	_ = b.WriteByte('[')
	for i, val := range vals {
		if i > 0 {
			_ = b.WriteByte(',')
		}
		appendJSONString(&b, val)
	}
	_ = b.WriteByte(']')
	return b.String()
}

func formatStringSliceValue(v any) string {
	switch vals := v.(type) {
	case [0]string:
		return "[]"
	case [1]string:
		return formatStringSlice(vals[:])
	case [2]string:
		return formatStringSlice(vals[:])
	case [3]string:
		return formatStringSlice(vals[:])
	default:
		return formatStringSliceReflect(v)
	}
}

func appendJSONString(dst *strings.Builder, s string) {
	_ = dst.WriteByte('"')
	start := 0

	for i := 0; i < len(s); {
		if c := s[i]; c < utf8.RuneSelf {
			if c >= 0x20 && c != '\\' && c != '"' {
				i++
				continue
			}

			if start < i {
				_, _ = dst.WriteString(s[start:i])
			}

			switch c {
			case '\\', '"':
				_ = dst.WriteByte('\\')
				_ = dst.WriteByte(c)
			case '\b':
				_, _ = dst.WriteString(`\b`)
			case '\f':
				_, _ = dst.WriteString(`\f`)
			case '\n':
				_, _ = dst.WriteString(`\n`)
			case '\r':
				_, _ = dst.WriteString(`\r`)
			case '\t':
				_, _ = dst.WriteString(`\t`)
			default:
				_, _ = dst.WriteString(`\u00`)
				_ = dst.WriteByte(lowerhex[c>>4])
				_ = dst.WriteByte(lowerhex[c&0x0f])
			}

			i++
			start = i
			continue
		}

		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			if start < i {
				_, _ = dst.WriteString(s[start:i])
			}
			_, _ = dst.WriteString(`\ufffd`)
			i++
			start = i
			continue
		}

		if r == '\u2028' || r == '\u2029' {
			if start < i {
				_, _ = dst.WriteString(s[start:i])
			}
			_, _ = dst.WriteString(`\u202`)
			_ = dst.WriteByte(lowerhex[r&0x0f])
			i += size
			start = i
			continue
		}

		i += size
	}

	if start < len(s) {
		_, _ = dst.WriteString(s[start:])
	}
	_ = dst.WriteByte('"')
}

func formatBoolSliceReflect(v any) string {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Array {
		return "unknown"
	}
	if rv.Len() == 0 {
		return "[]"
	}

	var b strings.Builder
	b.Grow(2 + rv.Len()*5)
	_ = b.WriteByte('[')
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			_ = b.WriteByte(',')
		}
		if rv.Index(i).Bool() {
			_, _ = b.WriteString("true")
		} else {
			_, _ = b.WriteString("false")
		}
	}
	_ = b.WriteByte(']')
	return b.String()
}

func formatInt64SliceReflect(v any) string {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Array {
		return "unknown"
	}
	if rv.Len() == 0 {
		return "[]"
	}

	var b strings.Builder
	b.Grow(2 + rv.Len()*20)
	_ = b.WriteByte('[')

	var scratch [20]byte
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			_ = b.WriteByte(',')
		}
		out := strconv.AppendInt(scratch[:0], rv.Index(i).Int(), 10)
		_, _ = b.Write(out)
	}

	_ = b.WriteByte(']')
	return b.String()
}

func formatFloat64SliceReflect(v any) string {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Array {
		return "unknown"
	}
	if rv.Len() == 0 {
		return "[]"
	}

	var b strings.Builder
	b.Grow(2 + rv.Len()*24)
	_ = b.WriteByte('[')

	var scratch [24]byte
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			_ = b.WriteByte(',')
		}
		val := rv.Index(i).Float()
		switch {
		case math.IsNaN(val):
			_, _ = b.WriteString(`"NaN"`)
		case math.IsInf(val, 1):
			_, _ = b.WriteString(`"Infinity"`)
		case math.IsInf(val, -1):
			_, _ = b.WriteString(`"-Infinity"`)
		default:
			out := strconv.AppendFloat(scratch[:0], val, 'g', -1, 64)
			_, _ = b.Write(out)
		}
	}

	_ = b.WriteByte(']')
	return b.String()
}

func formatStringSliceReflect(v any) string {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() || rv.Kind() != reflect.Array {
		return "unknown"
	}
	if rv.Len() == 0 {
		return "[]"
	}

	size := rv.Len() + 1
	for i := 0; i < rv.Len(); i++ {
		size += len(rv.Index(i).String()) + 2
	}

	var b strings.Builder
	b.Grow(size)
	_ = b.WriteByte('[')
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			_ = b.WriteByte(',')
		}
		appendJSONString(&b, rv.Index(i).String())
	}
	_ = b.WriteByte(']')
	return b.String()
}

// MarshalJSON returns the JSON encoding of the Value.
func (v Value) MarshalJSON() ([]byte, error) {
	var jsonVal struct {
		Type  string
		Value any
	}
	jsonVal.Type = v.Type().String()
	jsonVal.Value = v.AsInterface()
	return json.Marshal(jsonVal)
}
