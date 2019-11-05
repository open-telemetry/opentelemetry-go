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

package metric // import "go.opentelemetry.io/otel/sdk/export/metric"

import (
	"context"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/unit"
)

// Batcher is responsible for deciding which kind of aggregation
// to use and gathering exported results from the SDK.  The standard SDK
// supports binding only one of these interfaces, i.e., a single exporter.
//
// Multiple-exporters could be implemented by implementing this interface
// for a group of Batcher.
type Batcher interface {
	// AggregatorFor should return the kind of aggregator
	// suited to the requested export.  Returning `nil`
	// indicates to ignore the metric update.
	//
	// Note: This is context-free because the handle should not be
	// bound to the incoming context.  This call should not block.
	AggregatorFor(Record) Aggregator

	// Export receives pairs of records and aggregators
	// during the SDK Collect().  Exporter implementations
	// must access the specific aggregator to receive the
	// exporter data, since the format of the data varies
	// by aggregation.
	Export(context.Context, Record, Aggregator)
}

// Aggregator implements a specific aggregation behavior, e.g.,
// a counter, a gauge, a histogram.
type Aggregator interface {
	// Update receives a new measured value and incorporates it
	// into the aggregation.
	Update(context.Context, core.Number, Record)

	// Collect is called during the SDK Collect() to
	// finish one period of aggregation.  Collect() is
	// called in a single-threaded context.  Update()
	// calls may arrive concurrently.
	Collect(context.Context, Record, Batcher)

	// Merge combines state from two aggregators into one.
	Merge(Aggregator, *Descriptor)
}

// Record is the unit of export, pairing a metric
// instrument and set of labels.
type Record interface {
	// Descriptor() describes the metric instrument.
	Descriptor() *Descriptor

	// Labels() describe the labsels corresponding the
	// aggregation being performed.
	Labels() []core.KeyValue
}

// Kind describes the kind of instrument.
type Kind int8

const (
	CounterKind Kind = iota
	GaugeKind
	MeasureKind
)

// Descriptor describes a metric instrument to the exporter.
type Descriptor struct {
	name        string
	metricKind  Kind
	keys        []core.Key
	description string
	unit        unit.Unit
	numberKind  core.NumberKind
	alternate   bool
}

// NewDescriptor builds a new descriptor, for use by `Meter`
// implementations to interface with a metric export pipeline.
func NewDescriptor(
	name string,
	metricKind Kind,
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

func (d *Descriptor) MetricKind() Kind {
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
