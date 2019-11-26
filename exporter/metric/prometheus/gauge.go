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

type gauges struct {
	registerer prometheus.Registerer
	gauges     map[metricKey]prometheus.Gauge
	gaugeVecs  map[*export.Descriptor]*prometheus.GaugeVec
}

func newGauges(registerer prometheus.Registerer) gauges {
	return gauges{
		registerer: registerer,
		gauges:     make(map[metricKey]prometheus.Gauge),
		gaugeVecs:  make(map[*export.Descriptor]*prometheus.GaugeVec),
	}
}

func (ga *gauges) export(gauge aggregator.LastValue, record export.Record, mKey metricKey) error {
	lv, _, err := gauge.LastValue()
	if err != nil {
		return err
	}

	g, err := ga.getGauge(record, mKey)
	if err != nil {
		return err
	}

	desc := record.Descriptor()
	g.Set(lv.CoerceToFloat64(desc.NumberKind()))

	return nil
}

func (ga *gauges) getGauge(record export.Record, mKey metricKey) (prometheus.Gauge, error) {
	if c, ok := ga.gauges[mKey]; ok {
		return c, nil
	}

	desc := record.Descriptor()
	gaugeVec, err := ga.getGaugeVec(desc, record.Labels())
	if err != nil {
		return nil, err
	}

	gauge, err := gaugeVec.GetMetricWithLabelValues(labelValues(record.Labels())...)
	if err != nil {
		return nil, err
	}

	ga.gauges[mKey] = gauge
	return gauge, nil
}

func (ga *gauges) getGaugeVec(desc *export.Descriptor, labels export.Labels) (*prometheus.GaugeVec, error) {
	if gv, ok := ga.gaugeVecs[desc]; ok {
		return gv, nil
	}

	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: sanitize(desc.Name()),
			Help: desc.Description(),
		},
		labelsKeys(labels.Ordered()),
	)

	if err := ga.registerer.Register(g); err != nil {
		return nil, err
	}

	ga.gaugeVecs[desc] = g
	return g, nil
}
