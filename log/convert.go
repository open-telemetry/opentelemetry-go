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
		// return v.asInt64Slice()
	case attribute.FLOAT64:
		// return v.AsFloat64()
	case attribute.FLOAT64SLICE:
		// return v.asFloat64Slice()
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
