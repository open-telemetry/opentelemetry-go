// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"go.opentelemetry.io/otel/attribute"
)

// ConvertAttributeValue converts [attribute.Value) to [Value].
func ConvertAttributeValue(value attribute.Value) Value {
	switch value.Type() {
	case attribute.INVALID:
		return Value{}
	case attribute.BOOL:
		return BoolValue(value.AsBool())
	case attribute.BOOLSLICE:
		val := value.AsBoolSlice()
		res := make([]Value, 0, len(val))
		for _, v := range val {
			res = append(res, BoolValue(v))
		}
		return SliceValue(res...)
	case attribute.INT64:
		return Int64Value(value.AsInt64())
	case attribute.INT64SLICE:
		val := value.AsInt64Slice()
		res := make([]Value, 0, len(val))
		for _, v := range val {
			res = append(res, Int64Value(v))
		}
		return SliceValue(res...)
	case attribute.FLOAT64:
		return Float64Value(value.AsFloat64())
	case attribute.FLOAT64SLICE:
		val := value.AsFloat64Slice()
		res := make([]Value, 0, len(val))
		for _, v := range val {
			res = append(res, Float64Value(v))
		}
		return SliceValue(res...)
	case attribute.STRING:
		// return v.stringly
	case attribute.STRINGSLICE:
		// return v.asStringSlice()
	}
	panic("unknown attribute type")
}

// ConvertAttributeKeyValue converts [attribute.KeyValue) to [KeyValue].
func ConvertAttributeKeyValue(kv attribute.KeyValue) KeyValue {
	return KeyValue{}
}
