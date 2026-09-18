// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttrEncoding(t *testing.T) {
	attrs := []Attr{
		String("user", "Alice"),
		Bool("admin", true),
		Int64("floor", -2),
		Float64("impact", 0.21362),
		Slice("reports", StringValue("Bob"), StringValue("Dave")),
		Map("favorites", String("food", "hot dog"), Int("number", 13)),
		Bytes("secret", []byte("NUI4RUZGRjc5ODAzODEwM0QyNjlCNjMzODEzRkM2MEM=")),
	}

	t.Run("CamelCase", runJSONEncodingTests(attrs, []byte(`[
		{
			"key": "user",
			"value": {
				"stringValue": "Alice"
			}
		},
		{
			"key": "admin",
			"value": {
				"boolValue": true
			}
		},
		{
			"key": "floor",
			"value": {
				"intValue": "-2"
			}
		},
		{
			"key": "impact",
			"value": {
				"doubleValue": 0.21362
			}
		},
		{
			"key": "reports",
			"value": {
				"arrayValue": {
					"values": [
						{
							"stringValue": "Bob"
						},
						{
							"stringValue": "Dave"
						}
					]
				}
			}
		},
		{
			"key": "favorites",
			"value": {
				"kvlistValue": {
					"values": [
						{
							"key": "food",
							"value": {
								"stringValue": "hot dog"
							}
						},
						{
							"key": "number",
							"value": {
								"intValue": "13"
							}
						}
					]
				}
			}
		},
		{
			"key": "secret",
			"value": {
				"bytesValue": "TlVJNFJVWkdSamM1T0RBek9ERXdNMFF5TmpsQ05qTXpPREV6UmtNMk1FTT0="
			}
		}
	]`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(attrs, []byte(`[
		{
			"key": "user",
			"value": {
				"string_value": "Alice"
			}
		},
		{
			"key": "admin",
			"value": {
				"bool_value": true
			}
		},
		{
			"key": "floor",
			"value": {
				"int_value": "-2"
			}
		},
		{
			"key": "impact",
			"value": {
				"double_value": 0.21362
			}
		},
		{
			"key": "reports",
			"value": {
				"array_value": {
					"values": [
						{
							"string_value": "Bob"
						},
						{
							"string_value": "Dave"
						}
					]
				}
			}
		},
		{
			"key": "favorites",
			"value": {
				"kvlist_value": {
					"values": [
						{
							"key": "food",
							"value": {
								"string_value": "hot dog"
							}
						},
						{
							"key": "number",
							"value": {
								"int_value": "13"
							}
						}
					]
				}
			}
		},
		{
			"key": "secret",
			"value": {
				"bytes_value": "TlVJNFJVWkdSamM1T0RBek9ERXdNMFF5TmpsQ05qTXpPREV6UmtNMk1FTT0="
			}
		}
	]`)))
}

func TestFloat64ValueNonFiniteEncoding(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		encoded []byte
	}{
		{
			name:    "Finite",
			value:   0.21362,
			encoded: []byte(`{"doubleValue":0.21362}`),
		},
		{
			name:    "NaN",
			value:   math.NaN(),
			encoded: []byte(`{"doubleValue":"NaN"}`),
		},
		{
			name:    "PositiveInfinity",
			value:   math.Inf(1),
			encoded: []byte(`{"doubleValue":"Infinity"}`),
		},
		{
			name:    "NegativeInfinity",
			value:   math.Inf(-1),
			encoded: []byte(`{"doubleValue":"-Infinity"}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := Float64Value(test.value)

			got, err := json.Marshal(&value)
			require.NoError(t, err)
			assert.JSONEq(t, string(test.encoded), string(got))

			var decoded Value
			require.NoError(t, json.Unmarshal(test.encoded, &decoded))
			if math.IsNaN(test.value) {
				assert.True(t, math.IsNaN(decoded.AsFloat64()))
			} else {
				assert.Equal(t, value, decoded)
				decodedFloat := protoFloat64(decoded.AsFloat64())
				assert.Equal(t, test.value, (&decodedFloat).Float64())
			}
		})
	}
}
