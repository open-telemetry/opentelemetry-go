// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x_test

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/internal/x"
)

// MyExperimentalOption implements both metric.AddOption and
// x.ExperimentalOption by embedding them. This statically implements both
// types, allowing it to be passed to standard API functions like NewAddConfig
// without needing to define methods or run into unexported method constraints.
type MyExperimentalOption struct {
	x.ExperimentalOption
	metric.AddOption
}

func ExampleExperimentalOption() {
	// Users can pass standard options and experimental options together
	// because both implement metric.AddOption.
	opts := []metric.AddOption{
		metric.WithAttributes(attribute.String("key", "value")),
		MyExperimentalOption{},
	}

	// NewAddConfig is a standard configuration builder in the metrics API.
	// It ignores options that implement x.ExperimentalOption so this doesn't
	// panic. The configuration will only contain the standard option.
	config := metric.NewAddConfig(opts)
	attrs := config.Attributes()
	fmt.Printf("Number of attributes: %d\n", attrs.Len())

	// A consumer of our ExperimentalOption can detect and act on it.
	for _, opt := range opts {
		if _, ok := opt.(x.ExperimentalOption); ok {
			fmt.Println("Experimental option found")
		}
	}

	// Output:
	// Number of attributes: 1
	// Experimental option found
}
