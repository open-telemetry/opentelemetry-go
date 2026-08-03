// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package tracetransform

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestNilResource(t *testing.T) {
	arena := NewArena(16)
	assert.Empty(t, Resource(nil, arena))
}

func TestEmptyResource(t *testing.T) {
	arena := NewArena(16)
	assert.Empty(t, Resource(&resource.Resource{}, arena))
}

/*
* This does not include any testing on the ordering of Resource Attributes.
* They are stored as a map internally to the Resource and their order is not
* guaranteed.
 */

func TestResourceAttributes(t *testing.T) {
	attrs := []attribute.KeyValue{attribute.Int("one", 1), attribute.Int("two", 2)}

	arena := NewArena(16)
	got := Resource(resource.NewSchemaless(attrs...), arena).GetAttributes()
	if !assert.Len(t, attrs, 2) {
		return
	}
	assert.ElementsMatch(t, KeyValues(attrs, arena), got)
}
