// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import sdkmetric "go.opentelemetry.io/otel/sdk/metric"

type bindingOption struct {
	sdkmetric.Option
}

// Experimental prevents the SDK from applying the embedded stable option.
func (bindingOption) Experimental() {}

func (bindingOption) BindingEnabled() bool { return true }

// WithBinding enables experimental bound-instrument support on supported
// synchronous instruments created by a MeterProvider.
//
// This experimental implementation currently supports Int64Counter,
// Float64Counter, Int64UpDownCounter, and Float64UpDownCounter instruments.
// Streams using Sum aggregation resolve bound handles at Bind time; other
// compatible aggregations fall back to recording with the bound attributes.
//
// Combining this option with other experimental MeterProvider meter options
// (such as [WithFinish] or [WithMeterConfigurator]) is not yet supported.
func WithBinding() sdkmetric.Option {
	return bindingOption{}
}
