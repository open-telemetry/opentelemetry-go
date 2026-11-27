// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import "fmt"

func ExampleTraceState_Insert() {
	original, err := ParseTraceState("a=1,b=2")
	if err != nil {
		panic(err)
	}
	updated, err := original.Insert("b", "3")
	if err != nil {
		panic(err)
	}
	fmt.Println(original)
	fmt.Println(updated)
	// Output:
	// a=1,b=2
	// b=3,a=1
}
