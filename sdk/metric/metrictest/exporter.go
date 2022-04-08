// Copyright The OpenTelemetry Authors
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

package metrictest // import "go.opentelemetry.io/otel/sdk/metric/metrictest"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// Exporter is a manually collected exporter for testing the SDK.  It does not
// satisfy the `export.Exporter` interface because it is not intended to be
// used with the periodic collection of the SDK, instead the test should
// manually call `Collect()`
//
// Exporters are not thread safe, and should only be used for testing.
type Exporter struct {
	exports []ExportRecord
	// resource *resource.Resource

	controller *controller.Controller
}

func NewTestMeterProvider() (metric.MeterProvider, *Exporter) {
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithCollectPeriod(0),
	)
	exp := &Exporter{

		controller: c,
	}

	return c, exp
}

// ExportRecord represents one collected datapoint from the Exporter.
type ExportRecord struct {
	InstrumentName         string
	InstrumentationLibrary Library
	Labels                 []attribute.KeyValue
	AggregationKind        aggregation.Kind
	Sum                    number.Number
	Count                  uint64
	Histogram              aggregation.Buckets
	LastValue              number.Number
}

// Collect triggers the SDK's collect methods and then aggregates the data into
// `ExportRecord`s.
func (e *Exporter) Collect(ctx context.Context) error {
	e.exports = []ExportRecord{}

	e.controller.Collect(ctx)

	e.controller.ForEach(func(l instrumentation.Library, r export.Reader) error {
		lib := Library{
			InstrumentationName:    l.Name,
			InstrumentationVersion: l.Version,
			SchemaURL:              l.SchemaURL,
		}

		r.ForEach(aggregation.CumulativeTemporalitySelector(), func(rec export.Record) error {
			record := ExportRecord{
				InstrumentName:         rec.Descriptor().Name(),
				InstrumentationLibrary: lib,
				Labels:                 rec.Labels().ToSlice(),
				AggregationKind:        rec.Aggregation().Kind(),
			}

			var err error
			switch agg := rec.Aggregation().(type) {
			case aggregation.Histogram:
				record.AggregationKind = aggregation.HistogramKind
				record.Histogram, err = agg.Histogram()
				if err != nil {
					return err
				}
				record.Sum, err = agg.Sum()
				if err != nil {
					return err
				}
				record.Count, err = agg.Count()
				if err != nil {
					return err
				}
			case aggregation.Count:
				record.Count, err = agg.Count()
				if err != nil {
					return err
				}
			case aggregation.LastValue:
				record.LastValue, _, err = agg.LastValue()
				if err != nil {
					return err
				}
			case aggregation.Sum:
				record.Sum, err = agg.Sum()
				if err != nil {
					return err
				}
			}

			e.exports = append(e.exports, record)
			return nil
		})
		return nil
	})
	return nil
}

// GetRecords returns all Records found by the SDK
func (e *Exporter) GetRecords() []ExportRecord {
	return e.exports
}

var ErrNotFound = fmt.Errorf("record not found")

// GetByName returns the first Record with a matching name.
func (e *Exporter) GetByName(name string) (ExportRecord, error) {
	for _, rec := range e.exports {
		if rec.InstrumentName == name {
			return rec, nil
		}
	}
	return ExportRecord{}, ErrNotFound
}

// GetByNameAndLabels returns the first Record with a matching name and set of labels
func (e *Exporter) GetByNameAndLabels(name string, labels []attribute.KeyValue) (ExportRecord, error) {
	for _, rec := range e.exports {
		if rec.InstrumentName == name && labelsMatch(labels, rec.Labels) {
			return rec, nil
		}
	}
	return ExportRecord{}, ErrNotFound
}

func labelsMatch(labelsA, labelsB []attribute.KeyValue) bool {
	if len(labelsA) == len(labelsB) {
		return false
	}
	for i := range labelsA {
		if labelsA[i] != labelsB[i] {
			return false
		}
	}

	return true
}
