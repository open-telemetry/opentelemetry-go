// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/otlptranslator"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
)

// config contains options for the exporter.
type config struct {
	registerer               prometheus.Registerer
	disableTargetInfo        bool
	translationStrategy      otlptranslator.TranslationStrategyOption
	readerOpts               []metric.ManualReaderOption
	disableScopeInfo         bool
	namespace                string
	resourceAttributesFilter attribute.Filter
}

// newConfig creates a validated config configured with options.
func newConfig(opts ...Option) config {
	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.translationStrategy == "" {
		cfg.translationStrategy = otlptranslator.UnderscoreEscapingWithSuffixes
	}

	if cfg.registerer == nil {
		cfg.registerer = prometheus.DefaultRegisterer
	}

	return cfg
}

// Option sets exporter option values.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (fn optionFunc) apply(cfg config) config {
	return fn(cfg)
}

// WithRegisterer configures which prometheus Registerer the Exporter will
// register with.  If no registerer is used the prometheus DefaultRegisterer is
// used.
func WithRegisterer(reg prometheus.Registerer) Option {
	return optionFunc(func(cfg config) config {
		cfg.registerer = reg
		return cfg
	})
}

// WithAggregationSelector configure the Aggregation Selector the exporter will
// use. If no AggregationSelector is provided the DefaultAggregationSelector is
// used.
// If agg is nil, this option does nothing.
func WithAggregationSelector(agg metric.AggregationSelector) Option {
	return optionFunc(func(cfg config) config {
		cfg.readerOpts = append(cfg.readerOpts, metric.WithAggregationSelector(agg))
		return cfg
	})
}

// WithProducer configure the metric Producer the exporter will use as a source
// of external metric data.
// If producer is nil, this option does nothing.
func WithProducer(producer metric.Producer) Option {
	return optionFunc(func(cfg config) config {
		cfg.readerOpts = append(cfg.readerOpts, metric.WithProducer(producer))
		return cfg
	})
}

// WithoutTargetInfo configures the Exporter to not export the resource target_info metric.
// If not specified, the Exporter will create a target_info metric containing
// the metrics' resource.Resource attributes.
func WithoutTargetInfo() Option {
	return optionFunc(func(cfg config) config {
		cfg.disableTargetInfo = true
		return cfg
	})
}

// WithTranslationStrategy provides a standardized way to define how metric and
// label names should be handled during translation to Prometheus format. See:
// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.48.0/specification/metrics/sdk_exporters/prometheus.md#configuration.
// The recommended approach is to use either
// [otlptranslator.UnderscoreEscapingWithSuffixes] for full Prometheus-style
// compatibility or [otlptranslator.NoTranslation] for OpenTelemetry-style names.
//
// If this option is not used, the exporter defaults to
// [otlptranslator.UnderscoreEscapingWithSuffixes]. Users wanting a different
// translation strategy should specify it explicitly.
func WithTranslationStrategy(strategy otlptranslator.TranslationStrategyOption) Option {
	return optionFunc(func(cfg config) config {
		cfg.translationStrategy = strategy
		return cfg
	})
}

// WithoutScopeInfo configures the Exporter to not export
// labels about Instrumentation Scope to all metric points.
func WithoutScopeInfo() Option {
	return optionFunc(func(cfg config) config {
		cfg.disableScopeInfo = true
		return cfg
	})
}

// WithNamespace configures the Exporter to prefix metric with the given
// namespace. Metadata metrics such as target_info are not prefixed since these
// have special behavior based on their name. Namespaces will be prepended even
// if [otlptranslator.NoTranslation] is set as a translation strategy. If the provided namespace
// is empty, nothing will be prepended to metric names.
func WithNamespace(ns string) Option {
	return optionFunc(func(cfg config) config {
		cfg.namespace = ns
		return cfg
	})
}

// WithResourceAsConstantLabels configures the Exporter to add the resource attributes the
// resourceFilter returns true for as attributes on all exported metrics.
//
// The does not affect the target info generated from resource attributes.
func WithResourceAsConstantLabels(resourceFilter attribute.Filter) Option {
	return optionFunc(func(cfg config) config {
		cfg.resourceAttributesFilter = resourceFilter
		return cfg
	})
}
