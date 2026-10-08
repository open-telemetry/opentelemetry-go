// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// valueJSON retains the debug type envelope while making non-finite floats
// representable in JSON, without changing shared attribute marshaling.
type valueJSON struct {
	Type  string
	Value any
}

type keyValueJSON struct {
	Key   attribute.Key
	Value valueJSON
}

func newKeyValueJSON(kv attribute.KeyValue) keyValueJSON {
	return keyValueJSON{Key: kv.Key, Value: newValueJSON(kv.Value)}
}

func newValueJSON(v attribute.Value) valueJSON {
	result := valueJSON{Type: v.Type().String()}
	switch v.Type() {
	case attribute.FLOAT64:
		result.Value = float64JSON(v.AsFloat64())
	case attribute.FLOAT64SLICE:
		values := v.AsFloat64Slice()
		result.Value = values
		for _, f := range values {
			if math.IsNaN(f) || math.IsInf(f, 0) {
				converted := make([]any, len(values))
				for i, value := range values {
					converted[i] = float64JSON(value)
				}
				result.Value = converted
				break
			}
		}
	case attribute.SLICE:
		values := v.AsSlice()
		converted := make([]valueJSON, len(values))
		for i, value := range values {
			converted[i] = newValueJSON(value)
		}
		result.Value = converted
	case attribute.MAP:
		values := v.AsMap()
		converted := make([]keyValueJSON, len(values))
		for i, kv := range values {
			converted[i] = newKeyValueJSON(kv)
		}
		result.Value = converted
	default:
		result.Value = v.AsInterface()
	}
	return result
}

func float64JSON(f float64) any {
	switch {
	case math.IsNaN(f):
		return "NaN"
	case math.IsInf(f, 1):
		return "Infinity"
	case math.IsInf(f, -1):
		return "-Infinity"
	default:
		return f
	}
}
