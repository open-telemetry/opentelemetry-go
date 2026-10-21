// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import sdkmetric "go.opentelemetry.io/otel/sdk/metric"

type finishingOption struct {
	sdkmetric.Option
}

// Experimental prevents the SDK from applying the embedded stable option.
func (finishingOption) Experimental() {}

func (finishingOption) FinishEnabled() bool { return true }

// WithFinish enables experimental exact-attribute Finish support on supported
// synchronous instruments created by a MeterProvider.
//
// This experimental implementation currently supports Int64Counter and
// Float64Counter instruments using Sum aggregation, and Int64Gauge and
// Float64Gauge instruments using LastValue aggregation. Finish is a no-op for
// other aggregations.
func WithFinish() sdkmetric.Option {
	return finishingOption{}
}
