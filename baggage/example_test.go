// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage_test

import (
	"fmt"

	"go.opentelemetry.io/otel/baggage"
)

func ExampleBaggage_SetMember() {
	orig, err := baggage.Parse("tenant=old")
	if err != nil {
		panic(err)
	}

	m, err := baggage.NewMember("tenant", "new")
	if err != nil {
		panic(err)
	}

	updated, err := orig.SetMember(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(orig)
	fmt.Println(updated)
	// Output:
	// tenant=old
	// tenant=new
}

func ExampleNewMemberRaw() {
	m, err := baggage.NewMemberRaw("greeting", "hello, world")
	if err != nil {
		panic(err)
	}

	fmt.Println(m.Value())
	fmt.Println(m.String())
	// Output:
	// hello, world
	// greeting=hello%2C%20world
}
