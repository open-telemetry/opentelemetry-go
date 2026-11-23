// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

func ExampleSet_Filter() {
	set := attribute.NewSet(attribute.String("A", "a"), attribute.String("B", "b"))
	filter := attribute.NewAllowKeysFilter("A")

	filtered, removed := set.Filter(filter)

	fmt.Println(set.Len())
	fmt.Println(filtered.Len())
	fmt.Println(removed[0].Key)

	// Output:
	// 2
	// 1
	// B
}
