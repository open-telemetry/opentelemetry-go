// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package baggage provides the base types used internally by the public
// `go.opentelemetry.io/otel/baggage` package to store and retrieve baggage.
package baggage

// List is the collection of baggage members. The W3C allows for duplicates,
// but OpenTelemetry does not, therefore, this is represented as a map.
type List map[string]Item

// Item is the value and metadata properties part of a list-member.
type Item struct {
	Value      string
	Properties []Property
}

// Property is a metadata entry for a list-member.
type Property struct {
	Key, Value string

	// HasValue indicates if a zero-value value means the property does not
	// have a value or if it was the zero-value.
	HasValue bool
}
