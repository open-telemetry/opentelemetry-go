// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/internal/baggage"
)

func TestItem(t *testing.T) {
	properties := []baggage.Property{{Key: "state", Value: "on", HasValue: true}}
	tests := []struct {
		name       string
		item       baggage.Item
		value      string
		metadata   string
		properties []baggage.Property
	}{
		{name: "ZeroValue", item: baggage.Item{}},
		{
			name:     "Metadata",
			item:     baggage.NewItemWithMetadata("value", "state=on"),
			value:    "value",
			metadata: "state=on",
		},
		{
			name:       "Properties",
			item:       baggage.NewItemWithProperties("value", properties),
			value:      "value",
			properties: properties,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.value, test.item.Value())
			assert.Equal(t, test.metadata, test.item.Metadata())
			assert.Equal(t, test.properties, test.item.Properties())
		})
	}
}
