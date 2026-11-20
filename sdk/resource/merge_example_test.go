// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource_test

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

func ExampleMerge() {
	a := resource.NewSchemaless(attribute.String("service.name", "checkout"))
	b := resource.NewSchemaless(attribute.String("service.name", "payments"))

	merged, err := resource.Merge(a, b)
	if err != nil {
		panic(err)
	}

	v, _ := merged.Set().Value("service.name")
	fmt.Println(v.AsString())
	// Output:
	// payments
}
