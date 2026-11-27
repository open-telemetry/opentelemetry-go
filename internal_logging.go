// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otel

import (
	"github.com/go-logr/logr"

	"go.opentelemetry.io/otel/internal/global"
)

// SetLogger configures the logger used internally to opentelemetry.
//
// Internal logging levels are controlled by the verbosity of the supplied
// logger: Warn messages are logged at V(1), Info messages at V(4), and
// Debug messages at V(8).
//
// For slog users, use logr.FromSlogHandler to adapt a slog.Handler to a
// logr.Logger: https://pkg.go.dev/github.com/go-logr/logr#FromSlogHandler
func SetLogger(logger logr.Logger) {
	global.SetLogger(logger)
}
