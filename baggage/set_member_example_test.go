// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package baggage

import (
	"fmt"
)

func ExampleBaggage_SetMember() {
	orig, err := Parse("tenant=old")
	if err != nil {
		panic(err)
	}

	m, err := NewMember("tenant", "new")
	if err != nil {
		panic(err)
	}

	updated, err := orig.SetMember(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(orig)
	fmt.Println(updated)
	//Output:
	// tenant=old
	// tenant=new
}
