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
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

// Exporter is a manually collected exporter for testing the SDK.  It does not
// satisfy the `export.Exporter` interface because it is not intended to be
// used with the periodic collection of the SDK, instead the test should
// manually call `Collect()`
//
// Exporters are not thread safe, and should only be used for testing.
type Exporter struct {
	// Records contains the last metrics collected.
	Records []ExportRecord

	controller          *controller.Controller
	temporalitySelector aggregation.TemporalitySelector
}

// NewTestMeterProvider creates a MeterProvider and Exporter to be used in tests.
func NewTestMeterProvider(opts ...Option) (metric.MeterProvider, *Exporter) {
	cfg := newConfig(opts...)

	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(),
			cfg.temporalitySelector,
		),
		controller.WithCollectPeriod(0),
	)
	exp := &Exporter{
		controller:          c,
		temporalitySelector: cfg.temporalitySelector,
	}

	return c, exp
}

// Library is the same as "sdk/instrumentation".Library but there is
// a package cycle to use it so it is redeclared here.
type Library struct {
	InstrumentationName    string
	InstrumentationVersion string
	SchemaURL              string
}

// ExportRecord represents one collected datapoint from the Exporter.
type ExportRecord struct {
	InstrumentName         string
	InstrumentationLibrary Library
	Attributes             []attribute.KeyValue
	AggregationKind        aggregation.Kind
	NumberKind             number.Kind
	Sum                    number.Number
	Count                  uint64
	Histogram              aggregation.Buckets
	LastValue              number.Number
}

// Collect triggers the SDK's collect methods and then aggregates the data into
// ExportRecords.  This will overwrite any previous collected metrics.
func (e *Exporter) Collect(ctx context.Context) error {
	e.Records = []ExportRecord{}

	err := e.controller.Collect(ctx)
	if err != nil {
		return err
	}

	return e.controller.ForEach(func(l instrumentation.Library, r export.Reader) error {
		lib := Library{
			InstrumentationName:    l.Name,
			InstrumentationVersion: l.Version,
			SchemaURL:              l.SchemaURL,
		}

		return r.ForEach(e.temporalitySelector, func(rec export.Record) error {
			record := ExportRecord{
				InstrumentName:         rec.Descriptor().Name(),
				InstrumentationLibrary: lib,
				Attributes:             rec.Attributes().ToSlice(),
				AggregationKind:        rec.Aggregation().Kind(),
				NumberKind:             rec.Descriptor().NumberKind(),
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

			e.Records = append(e.Records, record)
			return nil
		})
	})
}

// GetRecords returns all Records found by the SDK.
func (e *Exporter) GetRecords() []ExportRecord {
	return e.Records
}

var errNotFound = fmt.Errorf("record not found")

// GetByName returns the first Record with a matching instrument name.
func (e *Exporter) GetByName(name string) (ExportRecord, error) {
	for _, rec := range e.Records {
		if rec.InstrumentName == name {
			return rec, nil
		}
	}
	return ExportRecord{}, errNotFound
}

// GetByNameAndAttributes returns the first Record with a matching name and the sub-set of attributes.
func (e *Exporter) GetByNameAndAttributes(name string, attributes []attribute.KeyValue) (ExportRecord, error) {
	for _, rec := range e.Records {
		if rec.InstrumentName == name && subSet(attributes, rec.Attributes) {
			return rec, nil
		}
	}
	return ExportRecord{}, errNotFound
}

// subSet returns true if attributesA is a subset of attributesB.
func subSet(attributesA, attributesB []attribute.KeyValue) bool {
	b := attribute.NewSet(attributesB...)

	for _, kv := range attributesA {
		if v, found := b.Value(kv.Key); !found || v != kv.Value {
			return false
		}
	}
	return true
}

// NewDescriptor is a test helper for constructing test metric
// descriptors using standard options.
func NewDescriptor(name string, ikind sdkapi.InstrumentKind, nkind number.Kind, opts ...instrument.Option) sdkapi.Descriptor {
	cfg := instrument.NewConfig(opts...)
	return sdkapi.NewDescriptor(name, ikind, nkind, cfg.Description(), cfg.Unit())
}
