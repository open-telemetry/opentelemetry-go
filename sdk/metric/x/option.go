// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import sdkmetric "go.opentelemetry.io/otel/sdk/metric"

// WithBinding enables experimental binding on supported synchronous
// instruments created by a stable SDK MeterProvider.
//
// This proof of concept currently supports Int64Counter only.
func WithBinding() sdkmetric.Option {
	return bindingOption{Option: sdkmetric.WithView()}
}

// WithFinish enables experimental exact-attribute Finish on supported
// synchronous instruments created by a stable SDK MeterProvider.
//
// This proof of concept currently supports Int64Counter only.
func WithFinish() sdkmetric.Option {
	return finishingOption{Option: sdkmetric.WithView()}
}
