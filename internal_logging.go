// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel

import (
	"github.com/go-logr/logr"

	"go.opentelemetry.io/otel/internal/global"
)

// SetLogger configures the logger used internally to opentelemetry.
//
// The internal logger uses the following verbosity levels:
//   - V(1): Warn messages.
//   - V(4): Info messages.
//   - V(8): Debug messages.
//
// For log/slog users, [logr.FromSlogHandler] adapts a slog.Handler to a
// [logr.Logger].
func SetLogger(logger logr.Logger) {
	global.SetLogger(logger)
}
