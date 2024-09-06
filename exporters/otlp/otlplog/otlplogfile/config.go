// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlplogfile // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlplogfile"

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// Option configures a field of the configuration or return an error if needed.
type Option func(*config) (*config, error)

// config contains options for the OTLP Log file exporter.
type config struct {
	// Out is the output where the records should be written.
	out io.WriteCloser
	// Duration represents the interval when the buffer should be flushed.
	flushInterval time.Duration
}

func newConfig(options []Option) (*config, error) {
	c := &config{
		out:           os.Stdout,
		flushInterval: 5 * time.Second,
	}

	var configErr error
	for _, opt := range options {
		if _, err := opt(c); err != nil {
			configErr = errors.Join(configErr, err)
		}
	}

	if configErr != nil {
		return nil, configErr
	}

	return c, nil
}

// WithFile configures a file where the records will be exported.
// An error is returned if the file could not be created or opened.
func WithFile(path string) Option {
	return func(c *config) (*config, error) {
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}

		return WithWriter(file)(c)
	}
}

// WithWriter configures the destination where the exporter should output
// the records. By default, if not specified, stdout is used.
func WithWriter(w io.WriteCloser) Option {
	return func(c *config) (*config, error) {
		c.out = w
		return c, nil
	}
}

// WithFlushInterval configures the duration after which the buffer is periodically flushed to the output.
func WithFlushInterval(flushInterval time.Duration) Option {
	return func(c *config) (*config, error) {
		c.flushInterval = flushInterval
		return c, nil
	}
}
