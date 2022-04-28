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

package processortest // import "go.opentelemetry.io/otel/sdk/metric/processor/processortest"

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	// mapKey is the unique key for a metric, consisting of its unique
	// descriptor, distinct attributes, and distinct resource attributes.
	mapKey struct {
		desc     *sdkapi.Descriptor
		attrs    attribute.Distinct
		resource attribute.Distinct
	}

	// mapValue is value stored in a processor used to produce a
	// Reader.
	mapValue struct {
		attrs      *attribute.Set
		resource   *resource.Resource
		aggregator aggregator.Aggregator
	}

	// Output implements export.Reader.
	Output struct {
		m           map[mapKey]mapValue
		attrEncoder attribute.Encoder
		sync.RWMutex
	}

	// testAggregatorSelector returns aggregators consistent with
	// the test variables below, needed for testing stateful
	// processors, which clone Aggregators using AggregatorFor(desc).
	testAggregatorSelector struct{}

	// testCheckpointer is a export.Checkpointer.
	testCheckpointer struct {
		started  int
		finished int
		*Processor
	}

	// Processor is a testing implementation of export.Processor that
	// assembles its results as a map[string]float64.
	Processor struct {
		export.AggregatorSelector
		output *Output
	}

	// Exporter is a testing implementation of export.Exporter that
	// assembles its results as a map[string]float64.
	Exporter struct {
		aggregation.TemporalitySelector
		output      *Output
		exportCount int

		// InjectErr supports returning conditional errors from
		// the Export() routine.  This must be set before the
		// Exporter is first used.
		InjectErr func(export.Record) error
	}
)

type testFactory struct {
	selector export.AggregatorSelector
	encoder  attribute.Encoder
}

func NewCheckpointerFactory(selector export.AggregatorSelector, encoder attribute.Encoder) export.CheckpointerFactory {
	return testFactory{
		selector: selector,
		encoder:  encoder,
	}
}

func NewCheckpointer(p *Processor) export.Checkpointer {
	return &testCheckpointer{
		Processor: p,
	}
}

func (f testFactory) NewCheckpointer() export.Checkpointer {
	return NewCheckpointer(NewProcessor(f.selector, f.encoder))
}

// NewProcessor returns a new testing Processor implementation.
// Verify expected outputs using Values(), e.g.:
//
//     require.EqualValues(t, map[string]float64{
//         "counter.sum/A=1,B=2/R=V": 100,
//     }, processor.Values())
//
// Where in the example A=1,B=2 is the encoded attributes and R=V is the
// encoded resource value.
func NewProcessor(selector export.AggregatorSelector, encoder attribute.Encoder) *Processor {
	return &Processor{
		AggregatorSelector: selector,
		output:             NewOutput(encoder),
	}
}

// Process implements export.Processor.
func (p *Processor) Process(accum export.Accumulation) error {
	return p.output.AddAccumulation(accum)
}

// Values returns the mapping from attribute set to point values for the
// accumulations that were processed. Point values are chosen as either the
// Sum or the LastValue, whichever is implemented. (All the built-in
// Aggregators implement one of these interfaces.)
func (p *Processor) Values() map[string]float64 {
	return p.output.Map()
}

// Reset clears the state of this test processor.
func (p *Processor) Reset() {
	p.output.Reset()
}

// StartCollection implements export.Checkpointer.
func (c *testCheckpointer) StartCollection() {
	if c.started != c.finished {
		panic(fmt.Sprintf("collection was already started: %d != %d", c.started, c.finished))
	}

	c.started++
}

// FinishCollection implements export.Checkpointer.
func (c *testCheckpointer) FinishCollection() error {
	if c.started-1 != c.finished {
		return fmt.Errorf("collection was not started: %d != %d", c.started, c.finished)
	}

	c.finished++
	return nil
}

// Reader implements export.Checkpointer.
func (c *testCheckpointer) Reader() export.Reader {
	return c.Processor.output
}

// AggregatorSelector returns a policy that is consistent with the
// test descriptors above.  I.e., it returns sum.New() for counter
// instruments and lastvalue.New() for lastValue instruments.
func AggregatorSelector() export.AggregatorSelector {
	return testAggregatorSelector{}
}

// AggregatorFor implements export.AggregatorSelector.
func (testAggregatorSelector) AggregatorFor(desc *sdkapi.Descriptor, aggPtrs ...*aggregator.Aggregator) {

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
	case strings.HasSuffix(desc.Name(), ".lastvalue"):
		aggs := lastvalue.New(len(aggPtrs))
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	case strings.HasSuffix(desc.Name(), ".histogram"):
		aggs := histogram.New(len(aggPtrs), desc)
		for i := range aggPtrs {
			*aggPtrs[i] = &aggs[i]
		}
	default:
		panic(fmt.Sprint("Invalid instrument name for test AggregatorSelector: ", desc.Name()))
	}
}

// NewOutput is a helper for testing an expected set of Accumulations
// (from an Accumulator) or an expected set of Records (from a
// Processor).  If testing with an Accumulator, it may be simpler to
// use the test Processor in this package.
func NewOutput(attrEncoder attribute.Encoder) *Output {
	return &Output{
		m:           make(map[mapKey]mapValue),
		attrEncoder: attrEncoder,
	}
}

// ForEach implements export.Reader.
func (o *Output) ForEach(_ aggregation.TemporalitySelector, ff func(export.Record) error) error {
	for key, value := range o.m {
		if err := ff(export.NewRecord(
			key.desc,
			value.attrs,
			value.aggregator.Aggregation(),
			time.Time{},
			time.Time{},
		)); err != nil {
			return err
		}
	}
	return nil
}

// AddRecord adds a string representation of the exported metric data
// to a map for use in testing.  The value taken from the record is
// either the Sum() or the LastValue() of its Aggregation(), whichever
// is defined.  Record timestamps are ignored.
func (o *Output) AddRecord(rec export.Record) error {
	return o.AddRecordWithResource(rec, resource.Empty())
}

func (o *Output) AddInstrumentationLibraryRecord(_ instrumentation.Library, rec export.Record) error {
	return o.AddRecordWithResource(rec, resource.Empty())
}

func (o *Output) AddRecordWithResource(rec export.Record, res *resource.Resource) error {
	key := mapKey{
		desc:     rec.Descriptor(),
		attrs:    rec.Attributes().Equivalent(),
		resource: res.Equivalent(),
	}
	if _, ok := o.m[key]; !ok {
		var agg aggregator.Aggregator
		testAggregatorSelector{}.AggregatorFor(rec.Descriptor(), &agg)
		o.m[key] = mapValue{
			aggregator: agg,
			attrs:      rec.Attributes(),
			resource:   res,
		}
	}
	return o.m[key].aggregator.Merge(rec.Aggregation().(aggregator.Aggregator), rec.Descriptor())
}

// Map returns the calculated values for test validation from a set of
// Accumulations or a set of Records.  When mapping records or
// accumulations into floating point values, the Sum() or LastValue()
// is chosen, whichever is implemented by the underlying Aggregator.
func (o *Output) Map() map[string]float64 {
	r := make(map[string]float64)
	err := o.ForEach(aggregation.StatelessTemporalitySelector(), func(record export.Record) error {
		for key, entry := range o.m {
			encoded := entry.attrs.Encoded(o.attrEncoder)
			rencoded := entry.resource.Encoded(o.attrEncoder)
			value := 0.0
			if s, ok := entry.aggregator.(aggregation.Sum); ok {
				sum, _ := s.Sum()
				value = sum.CoerceToFloat64(key.desc.NumberKind())
			} else if l, ok := entry.aggregator.(aggregation.LastValue); ok {
				last, _, _ := l.LastValue()
				value = last.CoerceToFloat64(key.desc.NumberKind())
			} else {
				panic(fmt.Sprintf("Unhandled aggregator type: %T", entry.aggregator))
			}
			name := fmt.Sprint(key.desc.Name(), "/", encoded, "/", rencoded)
			r[name] = value
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprint("Unexpected processor error: ", err))
	}
	return r
}

// Reset restores the Output to its initial state, with no accumulated
// metric data.
func (o *Output) Reset() {
	o.m = map[mapKey]mapValue{}
}

// AddAccumulation adds a string representation of the exported metric
// data to a map for use in testing.  The value taken from the
// accumulation is either the Sum() or the LastValue() of its
// Aggregator().Aggregation(), whichever is defined.
func (o *Output) AddAccumulation(acc export.Accumulation) error {
	return o.AddRecord(
		export.NewRecord(
			acc.Descriptor(),
			acc.Attributes(),
			acc.Aggregator().Aggregation(),
			time.Time{},
			time.Time{},
		),
	)
}

// New returns a new testing Exporter implementation.
// Verify exporter outputs using Values(), e.g.,:
//
//     require.EqualValues(t, map[string]float64{
//         "counter.sum/A=1,B=2/R=V": 100,
//     }, exporter.Values())
//
// Where in the example A=1,B=2 is the encoded attributes and R=V is the
// encoded resource value.
func New(selector aggregation.TemporalitySelector, encoder attribute.Encoder) *Exporter {
	return &Exporter{
		TemporalitySelector: selector,
		output:              NewOutput(encoder),
	}
}

func (e *Exporter) Export(_ context.Context, res *resource.Resource, ckpt export.InstrumentationLibraryReader) error {
	e.output.Lock()
	defer e.output.Unlock()
	e.exportCount++
	return ckpt.ForEach(func(library instrumentation.Library, mr export.Reader) error {
		return mr.ForEach(e.TemporalitySelector, func(r export.Record) error {
			if e.InjectErr != nil {
				if err := e.InjectErr(r); err != nil {
					return err
				}
			}
			return e.output.AddRecordWithResource(r, res)
		})
	})
}

// Values returns the mapping from attribute set to point values for the
// accumulations that were processed. Point values are chosen as either the
// Sum or the LastValue, whichever is implemented. (All the built-in
// Aggregators implement one of these interfaces.)
func (e *Exporter) Values() map[string]float64 {
	e.output.Lock()
	defer e.output.Unlock()
	return e.output.Map()
}

// ExportCount returns the number of times Export() has been called
// since the last Reset().
func (e *Exporter) ExportCount() int {
	e.output.Lock()
	defer e.output.Unlock()
	return e.exportCount
}

// Reset sets the exporter's output to the initial, empty state and
// resets the export count to zero.
func (e *Exporter) Reset() {
	e.output.Lock()
	defer e.output.Unlock()
	e.output.Reset()
	e.exportCount = 0
}

func OneInstrumentationLibraryReader(l instrumentation.Library, r export.Reader) export.InstrumentationLibraryReader {
	return oneLibraryReader{l, r}
}

type oneLibraryReader struct {
	library instrumentation.Library
	reader  export.Reader
}

func (o oneLibraryReader) ForEach(readerFunc func(instrumentation.Library, export.Reader) error) error {
	return readerFunc(o.library, o.reader)
}

func MultiInstrumentationLibraryReader(records map[instrumentation.Library][]export.Record) export.InstrumentationLibraryReader {
	return instrumentationLibraryReader{records: records}
}

type instrumentationLibraryReader struct {
	records map[instrumentation.Library][]export.Record
}

var _ export.InstrumentationLibraryReader = instrumentationLibraryReader{}

func (m instrumentationLibraryReader) ForEach(fn func(instrumentation.Library, export.Reader) error) error {
	for library, records := range m.records {
		if err := fn(library, &metricReader{records: records}); err != nil {
			return err
		}
	}
	return nil
}

type metricReader struct {
	sync.RWMutex
	records []export.Record
}

var _ export.Reader = &metricReader{}

func (m *metricReader) ForEach(_ aggregation.TemporalitySelector, fn func(export.Record) error) error {
	for _, record := range m.records {
		if err := fn(record); err != nil && err != aggregation.ErrNoData {
			return err
		}
	}
	return nil
}
