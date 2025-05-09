// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import "testing"

func TestResourceEncoding(t *testing.T) {
	res := &Resource{
		Attrs:        []Attr{String("key", "val")},
		DroppedAttrs: 10,
	}

	t.Run("CamelCase", runJSONEncodingTests(res, []byte(`{
		"attributes": [
			{
				"key": "key",
				"value": {
					"stringValue": "val"
				}
			}
		],
		"droppedAttributesCount": 10
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(res, []byte(`{
		"attributes": [
			{
				"key": "key",
				"value": {
					"string_value": "val"
				}
			}
		],
		"dropped_attributes_count": 10
	}`)))
}
