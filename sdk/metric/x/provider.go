// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import sdkmetric "go.opentelemetry.io/otel/sdk/metric"

// NewMeterProvider returns a stable-backed provider whose Int64Counter values
// implement the experimental binding and finishing contracts.
func NewMeterProvider(options ...sdkmetric.Option) *sdkmetric.MeterProvider {
	options = append(options, boundCounterOption{Option: sdkmetric.WithView()})
	return sdkmetric.NewMeterProvider(options...)
}
