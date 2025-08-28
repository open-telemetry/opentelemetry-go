// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/prometheus/otlptranslator"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestNewConfig(t *testing.T) {
	registry := prometheus.NewRegistry()

	aggregationSelector := func(metric.InstrumentKind) metric.Aggregation { return nil }
	producer := &noopProducer{}

	testCases := []struct {
		name             string
		options          []Option
		wantConfig       config
		legacyValidation bool
	}{
		{
			name:    "Default",
			options: nil,
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
			},
		},
		{
			name: "WithRegisterer",
			options: []Option{
				WithRegisterer(registry),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          registry,
			},
		},
		{
			name: "WithAggregationSelector",
			options: []Option{
				WithAggregationSelector(aggregationSelector),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				readerOpts:          []metric.ManualReaderOption{metric.WithAggregationSelector(aggregationSelector)},
			},
		},
		{
			name: "WithProducer",
			options: []Option{
				WithProducer(producer),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				readerOpts:          []metric.ManualReaderOption{metric.WithProducer(producer)},
			},
		},
		{
			name: "With Multiple Options",
			options: []Option{
				WithRegisterer(registry),
				WithAggregationSelector(aggregationSelector),
				WithProducer(producer),
			},

			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          registry,
				readerOpts: []metric.ManualReaderOption{
					metric.WithAggregationSelector(aggregationSelector),
					metric.WithProducer(producer),
				},
			},
		},
		{
			name: "nil options do nothing",
			options: []Option{
				WithRegisterer(nil),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
			},
		},
		{
			name: "without target_info metric",
			options: []Option{
				WithoutTargetInfo(),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				disableTargetInfo:   true,
			},
		},
		{
			name:             "legacy validation mode default",
			options:          []Option{},
			legacyValidation: true,
			wantConfig: config{
				translationStrategy: otlptranslator.UnderscoreEscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
			},
		},
		{
			name: "legacy validation mode, unit suffixes disabled",
			options: []Option{
				WithoutUnits(),
			},
			legacyValidation: true,
			wantConfig: config{
				translationStrategy: otlptranslator.UnderscoreEscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				withoutUnits:        true,
			},
		},
		{
			name: "legacy validation mode, counter suffixes disabled",
			options: []Option{
				WithoutCounterSuffixes(),
			},
			legacyValidation: true,
			wantConfig: config{
				translationStrategy:    otlptranslator.UnderscoreEscapingWithSuffixes,
				registerer:             prometheus.DefaultRegisterer,
				withoutCounterSuffixes: true,
			},
		},
		{
			name: "unit suffixes disabled",
			options: []Option{
				WithoutUnits(),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				withoutUnits:        true,
			},
		},
		{
			name: "NoTranslation implies no suffixes",
			options: []Option{
				WithTranslationStrategy(otlptranslator.NoTranslation),
			},
			wantConfig: config{
				translationStrategy:    otlptranslator.NoTranslation,
				withoutUnits:           true,
				withoutCounterSuffixes: true,
				registerer:             prometheus.DefaultRegisterer,
			},
		},
		{
			name: "translation strategy does not override unit suffixes disabled",
			options: []Option{
				WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
				WithoutUnits(),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.UnderscoreEscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				withoutUnits:        true,
			},
		},
		{
			name: "translation strategy does not override counter suffixes disabled",
			options: []Option{
				WithTranslationStrategy(otlptranslator.UnderscoreEscapingWithSuffixes),
				WithoutCounterSuffixes(),
			},
			wantConfig: config{
				translationStrategy:    otlptranslator.UnderscoreEscapingWithSuffixes,
				registerer:             prometheus.DefaultRegisterer,
				withoutCounterSuffixes: true,
			},
		},
		{
			name: "with namespace",
			options: []Option{
				WithNamespace("test"),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				namespace:           "test",
			},
		},
		{
			name: "with namespace with trailing underscore",
			options: []Option{
				WithNamespace("test"),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				namespace:           "test",
			},
		},
		{
			name: "with unsanitized namespace",
			options: []Option{
				WithNamespace("test/"),
			},
			wantConfig: config{
				translationStrategy: otlptranslator.NoUTF8EscapingWithSuffixes,
				registerer:          prometheus.DefaultRegisterer,
				namespace:           "test/",
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.legacyValidation {
				//nolint:staticcheck
				model.NameValidationScheme = model.LegacyValidation
			} else {
				//nolint:staticcheck
				model.NameValidationScheme = model.UTF8Validation
			}
			cfg := newConfig(tt.options...)
			// only check the length of readerOpts, since they are not comparable
			assert.Len(t, cfg.readerOpts, len(tt.wantConfig.readerOpts))
			cfg.readerOpts = nil
			tt.wantConfig.readerOpts = nil

			assert.Equal(t, tt.wantConfig, cfg)
		})
	}
}

type noopProducer struct{}

func (*noopProducer) Produce(context.Context) ([]metricdata.ScopeMetrics, error) {
	return nil, nil
}
