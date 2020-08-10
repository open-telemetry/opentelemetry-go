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

package processortest

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
	Output struct {
		m            map[nameWithNumKind]export.Aggregator
		labelEncoder label.Encoder
	}

	// testAggregatorSelector returns aggregators consistent with
	// the test variables below, needed for testing stateful
	// processors, which clone Aggregators using AggregatorFor(desc).
	testAggregatorSelector struct{}
)

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
