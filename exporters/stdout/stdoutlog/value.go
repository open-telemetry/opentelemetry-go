// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"encoding/json"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// valueJSON retains the body's debug type envelope around the shared non-OTLP
// representation.
type valueJSON struct {
	Type  string
	Value any
}

func newValueJSON(v attribute.Value) valueJSON {
	result := valueJSON{Type: v.Type().String()}
	switch v.Type() {
	case attribute.EMPTY:
		return result
	case attribute.STRING, attribute.BYTESLICE:
		result.Value = v.String()
		return result
	case attribute.FLOAT64:
		f := v.AsFloat64()
		if math.IsNaN(f) || math.IsInf(f, 0) {
			result.Value = v.String()
			return result
		}
	}
	// String returns logical text for scalar strings, bytes, and special floats;
	// the remaining types are JSON literals, arrays, or objects.
	result.Value = json.RawMessage(v.String())
	return result
}
