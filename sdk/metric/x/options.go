// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"strconv"

	"go.opentelemetry.io/otel/sdk/metric"
)

// ViewMatchingMode controls how multiple matching Views are applied to an Instrument.
type ViewMatchingMode int

const (
	// ViewMatchingModeIndependent specifies that each matching View creates a
	// separate metric stream independently. This is the default behavior.
	ViewMatchingModeIndependent ViewMatchingMode = iota

	// ViewMatchingModeComposable specifies that matching Views are combined
	// (merged) to modify metric streams rather than creating independent streams.
	ViewMatchingModeComposable
)

// String returns the string representation of the ViewMatchingMode.
func (m ViewMatchingMode) String() string {
	switch m {
	case ViewMatchingModeIndependent:
		return "independent"
	case ViewMatchingModeComposable:
		return "composable"
	default:
		return "unknown(" + strconv.Itoa(int(m)) + ")"
	}
}

type viewMatchingModeOption struct {
	metric.Option
	mode ViewMatchingMode
}

// Experimental prevents the API from panicking when the option is used.
func (viewMatchingModeOption) Experimental() {}

// ViewMatchingMode returns the configured ViewMatchingMode as an integer.
func (o viewMatchingModeOption) ViewMatchingMode() int {
	return int(o.mode)
}

// WithViewMatchingMode returns a metric.Option that configures the ViewMatchingMode
// for a MeterProvider.
func WithViewMatchingMode(mode ViewMatchingMode) metric.Option {
	return viewMatchingModeOption{mode: mode}
}
