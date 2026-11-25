// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func ExampleMerge() {
	a := resource.NewSchemaless(attribute.String("service.name", "first"))
	b := resource.NewSchemaless(attribute.String("service.name", "second"))

	merged, _ := resource.Merge(a, b)

	for _, kv := range merged.Attributes() {
		fmt.Println(kv.Key, kv.Value.AsString())
	}

	// Output: service.name second
}
