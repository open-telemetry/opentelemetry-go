// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/x"

import (
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/resource"
)

const defaultCardinalityLimit = 2000

type config struct {
	resource         *resource.Resource
	readers          []Reader
	views            []View
	exemplarFilter   exemplar.Filter
	cardinalityLimit int
}

func newConfig(options []Option) config {
	cfg := config{
		resource:         resource.Default(),
		exemplarFilter:   exemplar.TraceBasedFilter,
		cardinalityLimit: defaultCardinalityLimit,
	}
	for _, option := range options {
		cfg = option.apply(cfg)
	}
	return cfg
}

// Option configures a MeterProvider.
type Option interface{ apply(config) config }

type optionFunc func(config) config

func (f optionFunc) apply(cfg config) config { return f(cfg) }

// WithResource associates res with a MeterProvider.
func WithResource(res *resource.Resource) Option {
	return optionFunc(func(cfg config) config {
		if res != nil {
			cfg.resource = res
		}
		return cfg
	})
}

// WithReader associates reader with a MeterProvider.
func WithReader(reader Reader) Option {
	return optionFunc(func(cfg config) config {
		if reader != nil {
			cfg.readers = append(cfg.readers, reader)
		}
		return cfg
	})
}

// WithView associates views with a MeterProvider.
func WithView(views ...View) Option {
	return optionFunc(func(cfg config) config {
		cfg.views = append(cfg.views, views...)
		return cfg
	})
}

// WithExemplarFilter configures the exemplar filter.
func WithExemplarFilter(filter exemplar.Filter) Option {
	return optionFunc(func(cfg config) config {
		if filter != nil {
			cfg.exemplarFilter = filter
		}
		return cfg
	})
}

// WithCardinalityLimit sets the global per-stream cardinality limit. Values
// less than or equal to zero disable the limit.
func WithCardinalityLimit(limit int) Option {
	return optionFunc(func(cfg config) config {
		cfg.cardinalityLimit = limit
		return cfg
	})
}
