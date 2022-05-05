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
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/data"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/view"
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
	reader  *sdkmetric.ManualReader
}

// NewTestMeterProvider creates a MeterProvider and Exporter to be used in tests.
func NewTestMeterProvider(opts ...Option) (metric.MeterProvider, *Exporter) {
	exp := &Exporter{
		reader: sdkmetric.NewManualReader("inmemory"),
	}
	cfg := newConfig(opts...)
	sdk := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			exp.reader,
			view.WithDefaultAggregationTemporalitySelector(cfg.temporalitySelector),
		),
	)
	return sdk, exp
}

// ExportRecord represents one collected datapoint from the Exporter.
type ExportRecord struct {
	InstrumentName         string
	InstrumentationLibrary instrumentation.Library
	Attributes             []attribute.KeyValue
	AggregationKind        aggregation.Kind
	NumberKind             number.Kind
	StartTime              time.Time
	EndTime                time.Time
	Aggregation            aggregation.Aggregation
}

func toExportRecord(m data.Metrics, s data.Scope, inst data.Instrument, p data.Point) ExportRecord {
	return ExportRecord{
		InstrumentName:         inst.Descriptor.Name,
		InstrumentationLibrary: s.Library,
		Attributes:             p.Attributes.ToSlice(),
		AggregationKind:        p.Aggregation.Kind(),
		NumberKind:             inst.Descriptor.NumberKind,
		StartTime:              p.Start,
		EndTime:                p.End,
		Aggregation:            p.Aggregation,
	}
}

// Collect triggers the SDK's collect methods and then aggregates the data into
// ExportRecords.  This will overwrite any previous collected metrics.
func (e *Exporter) Collect(ctx context.Context) error {
	e.Records = []ExportRecord{}

	m := e.reader.Produce(nil)

	for _, scope := range m.Scopes {
		for _, instr := range scope.Instruments {
			for _, point := range instr.Points {
				e.Records = append(e.Records, toExportRecord(m, scope, instr, point))
			}
		}
	}

	return nil
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

// GetByNameAndAttributes returns the first Record with a matching name and set of Attributes.
func (e *Exporter) GetByNameAndAttributes(name string, attributes []attribute.KeyValue) (ExportRecord, error) {
	for _, rec := range e.Records {
		if rec.InstrumentName == name && subSet(attributes, rec.Attributes) {
			return rec, nil
		}
	}
	return ExportRecord{}, errNotFound
}

// subSet returns true if A is a subset of B
func subSet(attributesA, attributesB []attribute.KeyValue) bool {
	b := attribute.NewSet(attributesB...)

	for _, kv := range attributesA {
		if v, found := b.Value(kv.Key); !found || v != kv.Value {
			return false
		}
	}
	return true
}
