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

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
)

type summaries struct {
	defaultSummaryObjectives map[float64]float64
	registerer               prometheus.Registerer
	summary                  map[metricKey]prometheus.Observer
	summaryVec               map[*export.Descriptor]*prometheus.SummaryVec
}

func newSummaries(registerer prometheus.Registerer, defaultSummaryObjectives map[float64]float64) summaries {
	return summaries{
		registerer:               registerer,
		summary:                  make(map[metricKey]prometheus.Observer),
		summaryVec:               make(map[*export.Descriptor]*prometheus.SummaryVec),
		defaultSummaryObjectives: defaultSummaryObjectives,
	}
}

func (su *summaries) export(points aggregator.Points, record export.Record, mKey metricKey) error {
	values, err := points.Points()
	if err != nil {
		return err
	}

	obs, err := su.getSummary(record, mKey)
	if err != nil {
		return err
	}

	desc := record.Descriptor()
	for _, v := range values {
		obs.Observe(v.CoerceToFloat64(desc.NumberKind()))
	}
	return nil
}

func (su *summaries) getSummary(record export.Record, mKey metricKey) (prometheus.Observer, error) {
	if c, ok := su.summary[mKey]; ok {
		return c, nil
	}

	desc := record.Descriptor()
	histogramVec, err := su.getSummaryVec(desc, record.Labels())
	if err != nil {
		return nil, err
	}

	obs, err := histogramVec.GetMetricWithLabelValues(labelValues(record.Labels())...)
	if err != nil {
		return nil, err
	}

	su.summary[mKey] = obs
	return obs, nil
}

func (su *summaries) getSummaryVec(desc *export.Descriptor, labels export.Labels) (*prometheus.SummaryVec, error) {
	if gv, ok := su.summaryVec[desc]; ok {
		return gv, nil
	}

	g := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       sanitize(desc.Name()),
			Help:       desc.Description(),
			Objectives: su.defaultSummaryObjectives,
		},
		labelsKeys(labels.Ordered()),
	)

	if err := su.registerer.Register(g); err != nil {
		return nil, err
	}

	su.summaryVec[desc] = g
	return g, nil
}
