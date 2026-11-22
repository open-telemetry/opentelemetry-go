// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute_test

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

func ExampleNewMergeIterator() {
	first := attribute.NewSet(attribute.String("shared", "first"), attribute.String("one", "1"))
	second := attribute.NewSet(attribute.String("shared", "second"), attribute.String("two", "2"))
	iter := attribute.NewMergeIterator(&first, &second)
	for iter.Next() {
		kv := iter.Attribute()
		fmt.Printf("%s=%s\n", kv.Key, kv.Value.AsString())
	}
	// Output:
	// one=1
	// shared=first
	// two=2
}
