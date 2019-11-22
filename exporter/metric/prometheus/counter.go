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

type counters struct {
	registerer  prometheus.Registerer
	counters    map[metricKey]prometheus.Counter
	counterVecs map[*export.Descriptor]*prometheus.CounterVec
}

func newCounters(registerer prometheus.Registerer) counters {
	return counters{
		registerer:  registerer,
		counters:    make(map[metricKey]prometheus.Counter),
		counterVecs: make(map[*export.Descriptor]*prometheus.CounterVec),
	}
}

func (co *counters) export(sum aggregator.Sum, record export.Record, mKey metricKey) error {
	value, err := sum.Sum()
	if err != nil {
		return err
	}

	c, err := co.getCounter(record, mKey)
	if err != nil {
		return err
	}

	desc := record.Descriptor()
	c.Add(value.CoerceToFloat64(desc.NumberKind()))

	return nil
}

func (co *counters) getCounter(record export.Record, mKey metricKey) (prometheus.Counter, error) {
	if c, ok := co.counters[mKey]; ok {
		return c, nil
	}

	desc := record.Descriptor()
	counterVec, err := co.getCounterVec(desc, record.Labels())
	if err != nil {
		return nil, err
	}

	counter, err := counterVec.GetMetricWithLabelValues(labelValues(record.Labels())...)
	if err != nil {
		return nil, err
	}

	co.counters[mKey] = counter
	return counter, nil
}

func (co *counters) getCounterVec(desc *export.Descriptor, labels export.Labels) (*prometheus.CounterVec, error) {
	if c, ok := co.counterVecs[desc]; ok {
		return c, nil
	}

	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: sanitize(desc.Name()),
			Help: desc.Description(),
		},
		labelsKeys(labels.Ordered()),
	)

	if err := co.registerer.Register(c); err != nil {
		return nil, err
	}

	co.counterVecs[desc] = c
	return c, nil
}
