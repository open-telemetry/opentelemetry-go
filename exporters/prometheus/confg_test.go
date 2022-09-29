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
)

func TestNewConfig(t *testing.T) {
	registry := prometheus.NewRegistry()

	testCases := []struct {
		name           string
		options        []Option
		wantRegisterer prometheus.Registerer
		wantGatherer   prometheus.Gatherer
	}{
		{
			name:           "Default",
			options:        nil,
			wantRegisterer: prometheus.DefaultRegisterer,
			wantGatherer:   prometheus.DefaultGatherer,
		},

		{
			name: "WithGatherer",
			options: []Option{
				WithGatherer(registry),
			},
			wantRegisterer: prometheus.DefaultRegisterer,
			wantGatherer:   registry,
		},
		{
			name: "WithRegisterer",
			options: []Option{
				WithRegisterer(registry),
			},
			wantRegisterer: registry,
			wantGatherer:   prometheus.DefaultGatherer,
		},
		{
			name: "Multiple Options",
			options: []Option{
				WithGatherer(registry),
				WithRegisterer(registry),
			},
			wantRegisterer: registry,
			wantGatherer:   registry,
		},
		{
			name: "nil options do nothing",
			options: []Option{
				WithGatherer(nil),
				WithRegisterer(nil),
			},
			wantRegisterer: prometheus.DefaultRegisterer,
			wantGatherer:   prometheus.DefaultGatherer,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			cfg := newConfig(tt.options...)

			// If no Registry is provided you should get the DefaultRegisterer and DefaultGatherer.
			assert.Equal(t, tt.wantRegisterer, cfg.registerer)
			assert.Equal(t, tt.wantGatherer, cfg.gatherer)
		})
	}
}
