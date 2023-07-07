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

package oconf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestWithEnvTemporalityPreference(t *testing.T) {
	origReader := DefaultEnvOptionsReader.GetEnv
	tests := []struct {
		name     string
		envValue string
		want     map[metric.InstrumentKind]metricdata.Temporality
	}{
		{
			name:     "default do not set the selector",
			envValue: "",
		},
		{
			name:     "non-normative do not set the selector",
			envValue: "non-normative",
		},
		{
			name:     "cumulative",
			envValue: "cumulative",
			want: map[metric.InstrumentKind]metricdata.Temporality{
				metric.InstrumentKindCounter:                 metricdata.CumulativeTemporality,
				metric.InstrumentKindHistogram:               metricdata.CumulativeTemporality,
				metric.InstrumentKindUpDownCounter:           metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableCounter:       metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableUpDownCounter: metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableGauge:         metricdata.CumulativeTemporality,
			},
		},
		{
			name:     "delta",
			envValue: "delta",
			want: map[metric.InstrumentKind]metricdata.Temporality{
				metric.InstrumentKindCounter:                 metricdata.DeltaTemporality,
				metric.InstrumentKindHistogram:               metricdata.DeltaTemporality,
				metric.InstrumentKindUpDownCounter:           metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableCounter:       metricdata.DeltaTemporality,
				metric.InstrumentKindObservableUpDownCounter: metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableGauge:         metricdata.CumulativeTemporality,
			},
		},
		{
			name:     "lowmemory",
			envValue: "lowmemory",
			want: map[metric.InstrumentKind]metricdata.Temporality{
				metric.InstrumentKindCounter:                 metricdata.DeltaTemporality,
				metric.InstrumentKindHistogram:               metricdata.DeltaTemporality,
				metric.InstrumentKindUpDownCounter:           metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableCounter:       metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableUpDownCounter: metricdata.CumulativeTemporality,
				metric.InstrumentKindObservableGauge:         metricdata.CumulativeTemporality,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DefaultEnvOptionsReader.GetEnv = func(key string) string {
				if key == "OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE" {
					return tt.envValue
				}
				return origReader(key)
			}
			cfg := Config{}
			cfg = ApplyGRPCEnvConfigs(cfg)

			if tt.want == nil {
				// There is no function set, the SDK's default is used.
				assert.Nil(t, cfg.Metrics.TemporalitySelector)
				return
			}

			require.NotNil(t, cfg.Metrics.TemporalitySelector)
			for ik, want := range tt.want {
				assert.Equal(t, want, cfg.Metrics.TemporalitySelector(ik))
			}
		})
	}
	DefaultEnvOptionsReader.GetEnv = origReader
}
