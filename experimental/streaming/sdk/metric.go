// Copyright 2019, OpenTelemetry Authors
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

package sdk

import (
	"context"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/experimental/streaming/exporter"
)

type metricHandle struct {
	instrument metric.Instrument
	labels     metricLabels
}

type metricLabels struct {
	sdk   *sdk
	scope exporter.ScopeID
}

func (h *metricHandle) Record(ctx context.Context, value float64) {
	h.labels.sdk.exporter.Record(exporter.Event{
		Type:    exporter.SINGLE_METRIC,
		Context: ctx,
		Scope:   h.labels.scope,
		Measurement: exporter.Measurement{
			Instrument: h.instrument,
			Value:      value,
		},
	})
}

func (m metricLabels) Meter() metric.Meter {
	return m.sdk
}

func (s *sdk) DefineLabels(ctx context.Context, labels ...core.KeyValue) metric.LabelSet {
	return metricLabels{
		sdk:   s,
		scope: s.exporter.NewScope(exporter.ScopeID{}, labels...),
	}
}

func (s *sdk) RecorderFor(ctx context.Context, labels metric.LabelSet, inst metric.Instrument) metric.Recorder {
	mlabels, _ := labels.(metricLabels)

	return &metricHandle{
		instrument: inst,
		labels:     mlabels,
	}
}

func (s *sdk) RecordSingle(ctx context.Context, labels metric.LabelSet, input metric.Measurement) {
	mlabels, _ := labels.(metricLabels)
	s.exporter.Record(exporter.Event{
		Type:    exporter.SINGLE_METRIC,
		Context: ctx,
		Scope:   mlabels.scope,
		Measurement: exporter.Measurement{
			Instrument: input.Instrument,
			Value:      input.Value,
		}})
}

func (s *sdk) RecordBatch(ctx context.Context, labels metric.LabelSet, ms ...metric.Measurement) {
	oms := make([]exporter.Measurement, len(ms))
	mlabels, _ := labels.(metricLabels)

	for i, input := range ms {
		oms[i] = exporter.Measurement{
			Instrument: input.Instrument,
			Value:      input.Value,
		}
	}

	s.exporter.Record(exporter.Event{
		Type:         exporter.BATCH_METRIC,
		Context:      ctx,
		Scope:        mlabels.scope,
		Measurements: oms,
	})
}
