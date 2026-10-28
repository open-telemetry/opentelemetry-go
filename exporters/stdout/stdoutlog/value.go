// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"encoding/json"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

func newValueJSON(v attribute.Value) json.RawMessage {
	var quote bool
	switch v.Type() {
	case attribute.EMPTY:
		return nil
	case attribute.STRING, attribute.BYTESLICE:
		quote = true
	case attribute.FLOAT64:
		f := v.AsFloat64()
		quote = math.IsNaN(f) || math.IsInf(f, 0)
	}
	if quote {
		// Marshaling a string cannot fail.
		value, _ := json.Marshal(v.String())
		return value
	}
	// String returns logical text for scalar strings, bytes, and special floats;
	// the remaining types are JSON literals, arrays, or objects.
	return json.RawMessage(v.String())
}
