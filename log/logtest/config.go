// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

import (
	"go.opentelemetry.io/otel/log"
)

type config struct {
	minSeverity log.Severity
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.apply(c)
	}

	return c
}

// Option configures a [Hook].
type Option interface {
	apply(config) config
}

type optFunc func(config) config

func (f optFunc) apply(c config) config { return f(c) }

// WithMinSeverity returns an [Option] that configures the minimum severity the
// recorder will return true for when Enabled is called.
//
// By default, the recorder will be enabled for all levels.
func WithMinSeverity(l log.Severity) Option {
	return optFunc(func(c config) config {
		c.minSeverity = l
		return c
	})
}
