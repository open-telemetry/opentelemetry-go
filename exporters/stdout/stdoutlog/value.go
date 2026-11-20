// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog

import (
	"encoding/json"
	"math"

	"go.opentelemetry.io/otel/attribute"
)

// newValueJSON encodes v using the non-OTLP representation, ready to embed
// as a JSON value in a log record.
func newValueJSON(v attribute.Value) json.RawMessage {
	var quote bool
	switch v.Type() {
	case attribute.EMPTY:
		// A nil RawMessage is encoded as JSON null.
		return nil
	case attribute.STRING, attribute.BYTESLICE:
		// String returns unquoted text for strings and Base64 text for bytes.
		// Both need JSON string quoting and escaping.
		quote = true
	case attribute.FLOAT64:
		// Non-finite floats must be quoted as "NaN", "Infinity", or "-Infinity"
		// because JSON has no numeric literals for these values.
		f := v.AsFloat64()
		quote = math.IsNaN(f) || math.IsInf(f, 0)
	}
	if quote {
		// Marshal applies JSON string escaping and replaces invalid UTF-8.
		// Marshaling a string cannot fail.
		value, _ := json.Marshal(v.String())
		return value
	}
	// For all other types, String returns a complete JSON value, including
	// quoting nested strings and non-finite floats in arrays and maps.
	// Embed it directly to preserve JSON types and avoid double encoding.
	return json.RawMessage(v.String())
}
