// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import "testing"

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
