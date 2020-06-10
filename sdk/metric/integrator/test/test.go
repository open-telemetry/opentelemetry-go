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

	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type (
	// Output collects distinct metric/label set outputs.
	Output struct {
		Map          map[string]float64
		labelEncoder label.Encoder
	}

	// testAggregationSelector returns aggregators consistent with
	// the test variables below, needed for testing stateful
	// integrators, which clone Aggregators using AggregatorFor(desc).
	testAggregationSelector struct{}
)

func NewOutput(labelEncoder label.Encoder) Output {
	return Output{
		Map:          make(map[string]float64),
		labelEncoder: labelEncoder,
	}
}

// NewAggregationSelector returns a policy that is consistent with the
// test descriptors above.  I.e., it returns sum.New() for counter
// instruments and lastvalue.New() for lastValue instruments.
func NewAggregationSelector() export.AggregationSelector {
	return &testAggregationSelector{}
}

func (*testAggregationSelector) AggregatorFor(desc *metric.Descriptor, aggPtrs ...*export.Aggregator) {
	for _, aggp := range aggPtrs {
		switch desc.MetricKind() {
		case metric.CounterKind:
			*aggp = &sum.New(1)[0]
		case metric.ValueObserverKind:
			*aggp = &lastvalue.New(1)[0]
		default:
			panic("Invalid descriptor MetricKind for this test")
		}
	}
}

// AddTo adds a name/label-encoding entry with the lastValue or counter
// value to the output map.
func (o Output) AddTo(rec export.Record) error {
	encoded := rec.Labels().Encoded(o.labelEncoder)
	rencoded := rec.Resource().Encoded(o.labelEncoder)
	key := fmt.Sprint(rec.Descriptor().Name(), "/", encoded, "/", rencoded)
	var value float64

	if s, ok := rec.Aggregator().(aggregator.Sum); ok {
		sum, _ := s.Sum()
		value = sum.CoerceToFloat64(rec.Descriptor().NumberKind())
	} else if l, ok := rec.Aggregator().(aggregator.LastValue); ok {
		last, _, _ := l.LastValue()
		value = last.CoerceToFloat64(rec.Descriptor().NumberKind())
	} else {
		panic(fmt.Sprintf("Unhandled aggregator type: %T", rec.Aggregator()))
	}
	o.Map[key] = value
	return nil
}
