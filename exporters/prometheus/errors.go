// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import "errors"

// Sentinel errors for consistent error checks in tests.
var (
	ErrInvalidMetricType = errors.New("invalid metric type")
	ErrInvalidMetric     = errors.New("invalid metric")
	ErrEHScaleBelowMin   = errors.New("exponential histogram scale below minimum supported")
)
