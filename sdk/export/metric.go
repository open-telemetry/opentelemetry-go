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

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/metric"
	"go.opentelemetry.io/api/unit"
)

// Aggregator implements a specific aggregation behavior, e.g., a
// counter, a gauge, a histogram.
type MetricAggregator interface {
	// Update receives a new measured value and incorporates it
	// into the aggregation.
	Update(context.Context, core.Number, MetricRecord)

	// Collect is called during the SDK Collect() to
	// finish one period of aggregation.  Collect() is
	// called in a single-threaded context.  Update()
	// calls may arrive concurrently.
	Collect(context.Context, MetricRecord, MetricBatcher)
}

type MetricKind int8

const (
	CounterMetricKind MetricKind = iota
	GaugeMetricKind
	MeasureMetricKind
)

type Descriptor interface {
	Name() string
	MetricKind() MetricKind
	Keys() []core.Key
	Description() string
	Unit() unit.Unit
	NumberKind() core.NumberKind
	Alternate() bool
	ID() metric.InstrumentID
}

// MetricRecord is the unit of export, pairing a metric
// instrument and set of labels.
type MetricRecord interface {
	// Descriptor() describes the metric instrument.
	Descriptor() Descriptor

	// Labels() describe the labsels corresponding the
	// aggregation being performed.
	Labels() []core.KeyValue
}
