// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlplogfile // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile"

import "time"

type fnOpt func(config) config

func (f fnOpt) applyOption(c config) config { return f(c) }

// Option sets the configuration value for an Exporter.
type Option interface {
	applyOption(config) config
}

// config contains options for the OTLP Log file exporter.
type config struct {
	// Path to a file on disk where records must be appended.
	// This file is preferably a json line file as stated in the specification.
	// See: https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/protocol/file-exporter.md#json-lines-file
	// See: https://jsonlines.org
	path string
	// Duration represents the interval when the buffer should be flushed.
	flushInterval time.Duration
}

func newConfig(options []Option) config {
	c := config{
		path:          "/var/log/opentelemetry/logs.jsonl",
		flushInterval: 5 * time.Second,
	}
	for _, opt := range options {
		c = opt.applyOption(c)
	}
	return c
}

// WithFlushInterval configures the duration after which the buffer is periodically flushed to the disk.
func WithFlushInterval(flushInterval time.Duration) Option {
	return fnOpt(func(c config) config {
		c.flushInterval = flushInterval
		return c
	})
}

// WithPath defines a path to a file where the log records will be written.
// If not set, will default to /var/log/opentelemetry/logs.jsonl.
func WithPath(path string) Option {
	return fnOpt(func(c config) config {
		c.path = path
		return c
	})
}
