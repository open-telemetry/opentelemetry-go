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

package test

import (
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type (
	nameWithNumKind struct {
		name       string
		numberKind metric.NumberKind
	}

	// Output collects distinct metric/label set outputs.
	//
	// TODO(#872) make this internal.
	Output struct {
		m            map[nameWithNumKind]export.Aggregator
		labelEncoder label.Encoder
	}

	// testAggregatorSelector returns aggregators consistent with
	// the test variables below, needed for testing stateful
	// processors, which clone Aggregators using AggregatorFor(desc).
	testAggregatorSelector struct{}

	// testExportKindSelector is a ExportKindSelector
	testExportKindSelector export.ExportKind

	// Processor is a testing implementation of export.Processor that
	// assembles its results as a map[string]float64.
	Processor struct {
		export.AggregatorSelector
		output Output
	}

	// Exporter is a testing implementation of export.Exporter that
	// assembles its results as a map[string]float64.
	Exporter struct {
		export.ExportKindSelector
		proc export.Processor
	}
)

// NewProcessor returns a new testing Processor implementation.
// Verify expected outputs using Values(), e.g.:
//
//     require.EqualValues(t, map[string]float64{
//         "counter.sum/A=1,B=2/R=V": 100,
//     }, processor.Values())
//
// Where in the example A=1,B=2 is the encoded labels and R=V is the
// encoded resource value.
func NewProcessor(selector export.AggregatorSelector, encoder label.Encoder) *Processor {
	return &Processor{
		AggregatorSelector: selector,
		output:             NewOutput(encoder),
	}
}

// Process implements export.Processor.
func (p *Processor) Process(accum export.Accumulation) error {
	return p.output.AddAccumulation(accum)
}

// Values returns the mapping from label set to point values for the
// accumulations that were processed.  Point values are chosen as
// either the Sum or the LastValue, whichever is implemented.  (All
// the built-in Aggregators implement one of these interfaces.)
func (p *Processor) Values() map[string]float64 {
	return p.output.Map()
}

// NewExporter returns a new testing Exporter implementation.
// Verify exporter outputs using Values(), e.g.,:
//
//     require.EqualValues(t, map[string]float64{
//         "counter.sum/A=1,B=2/R=V": 100,
//     }, exporter.Values())
//
// Where in the example A=1,B=2 is the encoded labels and R=V is the
// encoded resource value.
func NewExporter(proc export.Processor, selector export.ExportKindSelector, encoder label.Encoder) *Exporter {
	return &Exporter{
		ExportKindSelector: selector,
		proc:               proc,
	}
}

// Values returns the mapping from label set to point values for the
// accumulations that were processed.  Point values are chosen as
// either the Sum or the LastValue, whichever is implemented.  (All
// the built-in Aggregators implement one of these interfaces.)
func (e *Exporter) Values(ckpt export.CheckpointSet) map[string]float64 {
	output := NewOutput(e.ExportKindSelector)
	err := ckpt.ForEach(e.ExportKindSelector, func(r Record) error {
		return e.output.AddRecord(r)
	})
	return output.Map()
}

// NewOutput is a helper for testing an expected set of Accumulations
// (from an Accumulator) or an expected set of Records (from a
// Processor).  If testing with an Accumulator, it may be simpler to
// use the test Processor in this package.
//
// TODO(#872): This class should be made internal, and callers should either
// use a test Processor or a test Exporter to use these facilities.
func NewOutput(labelEncoder label.Encoder) Output {
	return Output{
		m:            make(map[nameWithNumKind]export.Aggregator),
		labelEncoder: labelEncoder,
	}
}

// AggregatorSelector returns a policy that is consistent with the
// test descriptors above.  I.e., it returns sum.New() for counter
// instruments and lastvalue.New() for lastValue instruments.
func AggregatorSelector() export.AggregatorSelector {
	return testAggregatorSelector{}
}

func (testAggregatorSelector) AggregatorFor(desc *metric.Descriptor, aggPtrs ...*export.Aggregator) {

	switch {
	case strings.HasSuffix(desc.Name(), ".disabled"):
		for i := range aggPtrs {
			*aggPtrs[i] = nil
		}
	case strings.HasSuffix(desc.Name(), ".sum"):
		aggs := sum.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	case strings.HasSuffix(desc.Name(), ".minmaxsumcount"):
		aggs := minmaxsumcount.New(len(aggPtrs), desc)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	case strings.HasSuffix(desc.Name(), ".lastvalue"):
		aggs := lastvalue.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	case strings.HasSuffix(desc.Name(), ".sketch"):
		aggs := ddsketch.New(len(aggPtrs), desc, ddsketch.NewDefaultConfig())
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	case strings.HasSuffix(desc.Name(), ".histogram"):
		aggs := histogram.New(len(aggPtrs), desc, nil)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	case strings.HasSuffix(desc.Name(), ".exact"):
		aggs := array.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		panic(fmt.Sprint("Invalid instrument name for test AggregatorSelector: ", desc.Name()))
	}
}

// ExportKindSelector returns a policy with fixed export kind.
func ExportKindSelector(kind export.ExportKind) export.ExportKindSelector {
	return testExportKindSelector(kind)
}

// ExportKindFor implements export.ExportKindSelector.
func (kind testExportKindSelector) ExportKindFor(*metric.Descriptor, aggregation.Kind) export.ExportKind {
	return export.ExportKind(kind)
}

// AddRecord adds a string representation of the exported metric data
// to a map for use in testing.  The value taken from the record is
// either the Sum() or the LastValue() of its Aggregation(), whichever
// is defined.  Record timestamps are ignored.
func (o Output) AddRecord(rec export.Record) error {
	encoded := rec.Labels().Encoded(o.labelEncoder)
	rencoded := rec.Resource().Encoded(o.labelEncoder)
	key := nameWithNumKind{
		name:       fmt.Sprint(rec.Descriptor().Name(), "/", encoded, "/", rencoded),
		numberKind: rec.Descriptor().NumberKind(),
	}

	if _, ok := o.m[key]; !ok {
		var agg export.Aggregator
		testAggregatorSelector{}.AggregatorFor(rec.Descriptor(), &agg)
		o.m[key] = agg
	}
	return o.m[key].Merge(rec.Aggregation().(export.Aggregator), rec.Descriptor())
}

// Map returns the calculated values for test validation from a set of
// Accumulations or a set of Records.  When mapping records or
// accumulations into floating point values, the Sum() or LastValue()
// is chosen, whichever is implemented by the underlying Aggregator.
func (o Output) Map() map[string]float64 {
	r := make(map[string]float64)
	for nnk, agg := range o.m {
		value := 0.0
		if s, ok := agg.(aggregation.Sum); ok {
			sum, _ := s.Sum()
			value = sum.CoerceToFloat64(nnk.numberKind)
		} else if l, ok := agg.(aggregation.LastValue); ok {
			last, _, _ := l.LastValue()
			value = last.CoerceToFloat64(nnk.numberKind)
		} else {
			panic(fmt.Sprintf("Unhandled aggregator type: %T", agg))
		}
		r[nnk.name] = value
	}
	return r
}

// AddAccumulation adds a string representation of the exported metric
// data to a map for use in testing.  The value taken from the
// accumulation is either the Sum() or the LastValue() of its
// Aggregator().Aggregation(), whichever is defined.
func (o Output) AddAccumulation(acc export.Accumulation) error {
	return o.AddRecord(
		export.NewRecord(
			acc.Descriptor(),
			acc.Labels(),
			acc.Resource(),
			acc.Aggregator().Aggregation(),
			time.Time{},
			time.Time{},
		),
	)
}
