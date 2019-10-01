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
	descriptor metric.Descriptor
	labels     metricLabels
}

var _ metric.Handle = &metricHandle{}

type metricLabels struct {
	sdk   *sdk
	scope exporter.ScopeID
}

var _ metric.LabelSet = &metricLabels{}

func (h *metricHandle) RecordFloat(ctx context.Context, value float64) {
	h.labels.sdk.exporter.Record(exporter.Event{
		Type:    exporter.SINGLE_METRIC,
		Context: ctx,
		Scope:   h.labels.scope,
		Measurement: metric.Measurement{
			Descriptor: h.descriptor,
			Value:      metric.NewFloat64MeasurementValue(value),
		},
	})
}

func (h *metricHandle) RecordInt(ctx context.Context, value int64) {
	h.labels.sdk.exporter.Record(exporter.Event{
		Type:    exporter.SINGLE_METRIC,
		Context: ctx,
		Scope:   h.labels.scope,
		Measurement: metric.Measurement{
			Descriptor: h.descriptor,
			Value:      metric.NewInt64MeasurementValue(value),
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

func (s *sdk) NewHandle(ctx context.Context, descriptor metric.Descriptor, labels metric.LabelSet) metric.Handle {
	mlabels, _ := labels.(metricLabels)

	return &metricHandle{
		descriptor: descriptor,
		labels:     mlabels,
	}
}

func (s *sdk) DeleteHandle(handle metric.Handle) {
}

func (s *sdk) RecordBatch(ctx context.Context, labels metric.LabelSet, ms ...metric.Measurement) {
	eventType := exporter.BATCH_METRIC
	if len(ms) == 1 {
		eventType = exporter.SINGLE_METRIC
	}
	oms := make([]metric.Measurement, len(ms))
	mlabels, _ := labels.(metricLabels)

	copy(oms, ms)

	s.exporter.Record(exporter.Event{
		Type:         eventType,
		Context:      ctx,
		Scope:        mlabels.scope,
		Measurements: oms,
	})
}
