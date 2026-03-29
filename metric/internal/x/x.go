// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package x facilitates experimental metric API options.
package x // import "go.opentelemetry.io/otel/metric/internal/x"

// ExperimentalOption is an interface used to identify options that are
// experimental and should be ignored by standard configuration builders.
type ExperimentalOption interface {
	experimental()
}
