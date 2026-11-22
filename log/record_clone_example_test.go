// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func ExampleRecord_Clone() {
	attr0 := attribute.Bool("0", true)

	r0 := log.Record{}
	r0.AddAttributes(attr0)

	// Clone and add additional attribute to the clone
	attr1 := attribute.Int64("1", 2)

	r1 := r0.Clone()
	r1.AddAttributes(attr1)

	fmt.Println(r0.AttributesLen())
	fmt.Println(r1.AttributesLen())
	// Output:
	//1
	//2
}
