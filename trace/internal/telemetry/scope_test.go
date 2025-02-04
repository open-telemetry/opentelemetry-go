// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry

import "testing"

func TestScopeEncoding(t *testing.T) {
	scope := &Scope{
		Name:         "go.opentelemetry.io/otel/trace/internal/telemetry/test",
		Version:      "v0.0.1",
		Attrs:        []Attr{String("department", "ops")},
		DroppedAttrs: 1,
	}

	t.Run("CamelCase", runJSONEncodingTests(scope, []byte(`{
		"name": "go.opentelemetry.io/otel/trace/internal/telemetry/test",
		"version": "v0.0.1",
		"attributes": [
			{
				"key": "department",
				"value": {
					"stringValue": "ops"
				}
			}
		],
		"droppedAttributesCount": 1
	}`)))

	t.Run("SnakeCase/Unmarshal", runJSONUnmarshalTest(scope, []byte(`{
		"name": "go.opentelemetry.io/otel/trace/internal/telemetry/test",
		"version": "v0.0.1",
		"attributes": [
			{
				"key": "department",
				"value": {
					"string_value": "ops"
				}
			}
		],
		"dropped_attributes_count": 1
	}`)))
}
