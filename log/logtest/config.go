// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/log/logtest"

type config struct {
	enabledFn enabledFn
}

func newConfig(options []Option) config {
	var c config
	for _, opt := range options {
		c = opt.apply(c)
	}

	return c
}

// Option configures a [Recorder].
type Option interface {
	apply(config) config
}

type optFunc func(config) config

func (f optFunc) apply(c config) config { return f(c) }

// WithEnabledFn allows configuring whether the recorder enables specific log entries or not.
//
// By default, every log entry will be enabled.
func WithEnabledFn(fn enabledFn) Option {
	return optFunc(func(c config) config {
		c.enabledFn = fn
		return c
	})
}
