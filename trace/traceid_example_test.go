// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import "fmt"

func ExampleTraceIDFromHex() {
	id, err := TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	fmt.Println(id, err == nil)

	_, err = TraceIDFromHex("00000000000000000000000000000000")
	fmt.Println(err != nil)

	// Output:
	// 4bf92f3577b34da6a3ce929d0e0e4736 true
	// true
}
