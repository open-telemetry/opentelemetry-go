// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutmetric

import (
	"encoding/json"
	"io"
	"os"

	"go.opentelemetry.io/otel/sdk/metric"
)

// config contains options for the exporter.
type config struct {
	prettyPrint         bool
	encoder             *encoderHolder
	temporalitySelector metric.TemporalitySelector
	aggregationSelector metric.AggregationSelector
	redactTimestamps    bool
}

// newConfig creates a validated config configured with options.
func newConfig(options ...Option) config {
	cfg := config{}
	for _, opt := range options {
		cfg = opt.apply(cfg)
	}

	if cfg.encoder == nil {
		enc := newJSONEncoder(os.Stdout)
		cfg.encoder = &encoderHolder{encoder: enc}
	}

	if cfg.prettyPrint {
		switch e := cfg.encoder.encoder.(type) {
		case *json.Encoder:
			e.SetIndent("", "\t")
		case *jsonEncoder:
			e.SetIndent("", "\t")
		}
	}

	if cfg.temporalitySelector == nil {
		cfg.temporalitySelector = metric.DefaultTemporalitySelector
	}

	if cfg.aggregationSelector == nil {
		cfg.aggregationSelector = metric.DefaultAggregationSelector
	}

	return cfg
}

// Option sets exporter option values.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (o optionFunc) apply(c config) config {
	return o(c)
}

// WithEncoder sets the exporter to use encoder to encode all the metric
// data-types to an output.
func WithEncoder(encoder Encoder) Option {
	return optionFunc(func(c config) config {
		if encoder != nil {
			c.encoder = &encoderHolder{encoder: encoder}
		}
		return c
	})
}

// WithWriter sets the export stream destination.
// Using this option overrides any previously set encoder.
// If w is nil, this option has no effect.
func WithWriter(w io.Writer) Option {
	if w == nil {
		return WithEncoder(nil)
	}
	return optionFunc(func(c config) config {
		c.encoder = &encoderHolder{encoder: newJSONEncoder(w)}
		return c
	})
}

// WithPrettyPrint prettifies the emitted output.
// This option works with the default encoder, WithWriter, or a *json.Encoder
// passed to WithEncoder.
func WithPrettyPrint() Option {
	return optionFunc(func(c config) config {
		c.prettyPrint = true
		return c
	})
}

// WithTemporalitySelector sets the TemporalitySelector the exporter will use
// to determine the Temporality of an instrument based on its kind. If this
// option is not used, the exporter will use the DefaultTemporalitySelector
// from the go.opentelemetry.io/otel/sdk/metric package.
func WithTemporalitySelector(selector metric.TemporalitySelector) Option {
	return temporalitySelectorOption{selector: selector}
}

type temporalitySelectorOption struct {
	selector metric.TemporalitySelector
}

func (t temporalitySelectorOption) apply(c config) config {
	c.temporalitySelector = t.selector
	return c
}

// WithAggregationSelector sets the AggregationSelector the exporter will use
// to determine the aggregation to use for an instrument based on its kind. If
// this option is not used, the exporter will use the
// DefaultAggregationSelector from the go.opentelemetry.io/otel/sdk/metric
// package or the aggregation explicitly passed for a view matching an
// instrument.
func WithAggregationSelector(selector metric.AggregationSelector) Option {
	return aggregationSelectorOption{selector: selector}
}

type aggregationSelectorOption struct {
	selector metric.AggregationSelector
}

func (t aggregationSelectorOption) apply(c config) config {
	c.aggregationSelector = t.selector
	return c
}

// WithoutTimestamps sets all timestamps to zero in the output stream.
func WithoutTimestamps() Option {
	return optionFunc(func(c config) config {
		c.redactTimestamps = true
		return c
	})
}
