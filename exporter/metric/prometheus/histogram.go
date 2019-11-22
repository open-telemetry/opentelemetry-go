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

type histograms struct {
	defaultHistogramBuckets []float64
	registerer              prometheus.Registerer
	histogram               map[metricKey]prometheus.Observer
	histogramVecs           map[*export.Descriptor]*prometheus.HistogramVec
}

func newHistograms(registerer prometheus.Registerer, defaultHistogramBuckets []float64) histograms {
	return histograms{
		registerer:              registerer,
		histogram:               make(map[metricKey]prometheus.Observer),
		histogramVecs:           make(map[*export.Descriptor]*prometheus.HistogramVec),
		defaultHistogramBuckets: defaultHistogramBuckets,
	}
}

func (hi *histograms) export(points aggregator.Points, record export.Record, mKey metricKey) error {
	values, err := points.Points()
	if err != nil {
		return err
	}

	obs, err := hi.getHistogram(record, mKey)
	if err != nil {
		return err
	}

	desc := record.Descriptor()
	for _, v := range values {
		obs.Observe(v.CoerceToFloat64(desc.NumberKind()))
	}
	return nil
}

func (hi *histograms) getHistogram(record export.Record, mKey metricKey) (prometheus.Observer, error) {
	if c, ok := hi.histogram[mKey]; ok {
		return c, nil
	}

	desc := record.Descriptor()
	histogramVec, err := hi.getHistogramVec(desc, record.Labels())
	if err != nil {
		return nil, err
	}

	obs, err := histogramVec.GetMetricWithLabelValues(labelValues(record.Labels())...)
	if err != nil {
		return nil, err
	}

	hi.histogram[mKey] = obs
	return obs, nil
}

func (hi *histograms) getHistogramVec(desc *export.Descriptor, labels export.Labels) (*prometheus.HistogramVec, error) {
	if gv, ok := hi.histogramVecs[desc]; ok {
		return gv, nil
	}

	g := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    sanitize(desc.Name()),
			Help:    desc.Description(),
			Buckets: hi.defaultHistogramBuckets,
		},
		labelsKeys(labels.Ordered()),
	)

	if err := hi.registerer.Register(g); err != nil {
		return nil, err
	}

	hi.histogramVecs[desc] = g
	return g, nil
}
