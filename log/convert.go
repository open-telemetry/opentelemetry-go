// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import (
	"go.opentelemetry.io/otel/attribute"
)

// ConvertAttributeValue converts [attribute.Value) to [Value].
func ConvertAttributeValue(v attribute.Value) Value {
	switch v.Type() {
	case attribute.INVALID:
		return Value{}
	case attribute.BOOL:
		return BoolValue(v.AsBool())
	case attribute.BOOLSLICE:
		// return v.asBoolSlice()
	case attribute.INT64:
		// return v.AsInt64()
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
