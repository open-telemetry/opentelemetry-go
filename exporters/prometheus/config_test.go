// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus // import "go.opentelemetry.io/otel/exporters/prometheus"

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

func TestNewConfig(t *testing.T) {
	registry := prometheus.NewRegistry()

	aggregationSelector := func(metric.InstrumentKind) aggregation.Aggregation { return nil }

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
			},
		},
		{
			name: "With Multiple Options",
			options: []Option{
				WithRegisterer(registry),
				WithAggregationSelector(aggregationSelector),
			},

			wantConfig: config{
				registerer: registry,
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
			// tested by TestConfigManualReaderOptions
			cfg.aggregation = nil

			assert.Equal(t, tt.wantConfig, cfg)
		})
	}
}

func TestConfigManualReaderOptions(t *testing.T) {
	aggregationSelector := func(metric.InstrumentKind) aggregation.Aggregation { return nil }

	testCases := []struct {
		name            string
		config          config
		wantOptionCount int
	}{
		{
			name:            "Default",
			config:          config{},
			wantOptionCount: 0,
		},

		{
			name:            "WithAggregationSelector",
			config:          config{aggregation: aggregationSelector},
			wantOptionCount: 1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			opts := tt.config.manualReaderOptions()
			assert.Len(t, opts, tt.wantOptionCount)
		})
	}
}
