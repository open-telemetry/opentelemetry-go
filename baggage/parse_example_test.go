// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage_test

import (
	"fmt"

	"go.opentelemetry.io/otel/baggage"
)

func ExampleParse() {
	b, err := baggage.Parse("key=value,debug=inva lid")

	fmt.Println(err != nil)
	fmt.Println(b.Member("key").Value())

	// Output:
	// true
	// value
}
