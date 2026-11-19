// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace_test

import (
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

func ExampleTraceIDFromHex() {
	id, err := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	if err != nil {
		panic(err)
	}
	fmt.Println(id)

	_, err = trace.TraceIDFromHex("00000000000000000000000000000000")
	fmt.Println("all-zero rejected:", err != nil)
	// Output:
	// 4bf92f3577b34da6a3ce929d0e0e4736
	// all-zero rejected: true
}
