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

package metric

import (
	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/scope"
	"github.com/open-telemetry/opentelemetry-go/api/tag"
	"github.com/open-telemetry/opentelemetry-go/exporter/observer"
)

type baseMetric struct {
	measure core.Measure

	mtype   MetricType
	keys    []core.Key
	eventID core.EventID
	status  error // Indicates registry conflict
}

type baseEntry struct {
	base    *baseMetric
	metric  Metric
	eventID core.EventID
}

func initBaseMetric(name string, mtype MetricType, opts []Option, init Metric) Metric {
	var tagOpts []tag.Option
	bm := init.base()

	for _, opt := range opts {
		opt(bm, &tagOpts)
	}

	bm.measure = tag.NewMeasure(name, tagOpts...)
	bm.mtype = mtype

	bm.eventID = observer.Record(observer.Event{
		Type:  observer.NEW_METRIC,
		Scope: bm.measure.DefinitionID().Scope(),
	})

	other, err := GetRegistry().RegisterMetric(init)
	if err != nil {
		bm.status = err
	}
	return other
}

func (bm *baseMetric) base() *baseMetric {
	return bm
}

func (bm *baseMetric) Measure() core.Measure {
	return bm.measure
}

func (bm *baseMetric) Type() MetricType {
	return bm.mtype
}

func (bm *baseMetric) Fields() []core.Key {
	return bm.keys
}

func (bm *baseMetric) Err() error {
	return bm.status
}

func (e *baseEntry) init(m Metric, values []core.KeyValue) {
	e.base = m.base()
	e.metric = m
	e.eventID = scope.New(core.ScopeID{}, values...).ScopeID().EventID
}
