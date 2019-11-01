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

package export

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
)

// MetricAggregator implements a specific aggregation behavior, e.g.,
// a counter, a gauge, a histogram.
type MetricAggregator interface {
	// Update receives a new measured value and incorporates it
	// into the aggregation.
	Update(context.Context, core.Number, MetricRecord)

	// Collect is called during the SDK Collect() to
	// finish one period of aggregation.  Collect() is
	// called in a single-threaded context.  Update()
	// calls may arrive concurrently.
	Collect(context.Context, MetricRecord, MetricBatcher)

	// Merge combines state from two aggregators into one.
	Merge(MetricAggregator, *Descriptor)
}

// MetricRecord is the unit of export, pairing a metric
// instrument and set of labels.
type MetricRecord interface {
	// Descriptor() describes the metric instrument.
	Descriptor() *Descriptor

	// Labels() describe the labsels corresponding the
	// aggregation being performed.
	Labels() []core.KeyValue
}

// MetricKind describes the kind of instrument.
type MetricKind int8

const (
	CounterMetricKind MetricKind = iota
	GaugeMetricKind
	MeasureMetricKind
)

// Descriptor describes a metric instrument to the exporter.
type Descriptor struct {
	name        string
	metricKind  MetricKind
	keys        []core.Key
	description string
	unit        unit.Unit
	numberKind  core.NumberKind
	alternate   bool
}

// NewDescriptor builds a new descriptor, for use by `Meter`
// implementations.
func NewDescriptor(
	name string,
	metricKind MetricKind,
	keys []core.Key,
	description string,
	unit unit.Unit,
	numberKind core.NumberKind,
	alternate bool,
) *Descriptor {
	return &Descriptor{
		name:        name,
		metricKind:  metricKind,
		keys:        keys,
		description: description,
		unit:        unit,
		numberKind:  numberKind,
		alternate:   alternate,
	}
}

func (d *Descriptor) Name() string {
	return d.name
}

func (d *Descriptor) MetricKind() MetricKind {
	return d.metricKind
}

func (d *Descriptor) Keys() []core.Key {
	return d.keys
}

func (d *Descriptor) Description() string {
	return d.description
}

func (d *Descriptor) Unit() unit.Unit {
	return d.unit
}

func (d *Descriptor) NumberKind() core.NumberKind {
	return d.numberKind
}

func (d *Descriptor) Alternate() bool {
	return d.alternate
}
