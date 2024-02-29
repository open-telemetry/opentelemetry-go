// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import (
	"context"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestNewConfig(t *testing.T) {
	registry := prometheus.NewRegistry()

	aggregationSelector := func(metric.InstrumentKind) metric.Aggregation { return nil }
	producer := &noopProducer{}

	testCases := []struct {
		name       string
		options    []Option
		wantConfig config
	}{
		{
			name:    "Default",
			options: nil,
			wantConfig: config{
				registerer: prometheus.DefaultRegisterer,
			},
		},
		{
			name: "WithRegisterer",
			options: []Option{
				WithRegisterer(registry),
			},
			wantConfig: config{
				registerer: registry,
			},
		},
		{
			name: "WithAggregationSelector",
			options: []Option{
				WithAggregationSelector(aggregationSelector),
			},
			wantConfig: config{
				registerer: prometheus.DefaultRegisterer,
				readerOpts: []metric.ManualReaderOption{metric.WithAggregationSelector(aggregationSelector)},
			},
		},
		{
			name: "WithProducer",
			options: []Option{
				WithProducer(producer),
			},
			wantConfig: config{
				registerer: prometheus.DefaultRegisterer,
				readerOpts: []metric.ManualReaderOption{metric.WithProducer(producer)},
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
				registerer: registry,
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
				registerer: prometheus.DefaultRegisterer,
			},
		},
		{
			name: "without target_info metric",
			options: []Option{
				WithoutTargetInfo(),
			},
			wantConfig: config{
				registerer:        prometheus.DefaultRegisterer,
				disableTargetInfo: true,
			},
		},
		{
			name: "unit suffixes disabled",
			options: []Option{
				WithoutUnits(),
			},
			wantConfig: config{
				registerer:   prometheus.DefaultRegisterer,
				withoutUnits: true,
			},
		},
		{
			name: "with namespace",
			options: []Option{
				WithNamespace("test"),
			},
			wantConfig: config{
				registerer: prometheus.DefaultRegisterer,
				namespace:  "test_",
			},
		},
		{
			name: "with namespace with trailing underscore",
			options: []Option{
				WithNamespace("test_"),
			},
			wantConfig: config{
				registerer: prometheus.DefaultRegisterer,
				namespace:  "test_",
			},
		},
		{
			name: "with unsanitized namespace",
			options: []Option{
				WithNamespace("test/"),
			},
			wantConfig: config{
				registerer: prometheus.DefaultRegisterer,
				namespace:  "test_",
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(tt.options...)
			// only check the length of readerOpts, since they are not comparable
			assert.Equal(t, len(tt.wantConfig.readerOpts), len(cfg.readerOpts))
			cfg.readerOpts = nil
			tt.wantConfig.readerOpts = nil

			assert.Equal(t, tt.wantConfig, cfg)
		})
	}
}

type noopProducer struct{}

func (*noopProducer) Produce(ctx context.Context) ([]metricdata.ScopeMetrics, error) {
	return nil, nil
}
