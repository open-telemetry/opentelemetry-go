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
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/label"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/resource"
)

type (
	// Encoder is an alternate label encoder to validate grouping logic.
	Encoder struct{}

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

var (
	// Resource is applied to all test records built in this package.
	Resource = resource.New(kv.String("R", "V"))

	// LastValueADesc and LastValueBDesc group by "G"
	LastValueADesc = metric.NewDescriptor(
		"lastvalue.a", metric.ValueObserverKind, metric.Int64NumberKind)
	LastValueBDesc = metric.NewDescriptor(
		"lastvalue.b", metric.ValueObserverKind, metric.Int64NumberKind)
	// CounterADesc and CounterBDesc group by "C"
	CounterADesc = metric.NewDescriptor(
		"sum.a", metric.CounterKind, metric.Int64NumberKind)
	CounterBDesc = metric.NewDescriptor(
		"sum.b", metric.CounterKind, metric.Int64NumberKind)

	// SdkEncoder uses a non-standard encoder like K1~V1&K2~V2
	SdkEncoder = &Encoder{}
	// GroupEncoder uses the SDK default encoder
	GroupEncoder = label.DefaultEncoder()

	// LastValue groups are (labels1), (labels2+labels3)
	// Counter groups are (labels1+labels2), (labels3)

	// Labels1 has G=H and C=D
	Labels1 = makeLabels(kv.String("G", "H"), kv.String("C", "D"))
	// Labels2 has C=D and E=F
	Labels2 = makeLabels(kv.String("C", "D"), kv.String("E", "F"))
	// Labels3 is the empty set
	Labels3 = makeLabels()

	testLabelEncoderID = label.NewEncoderID()
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

func (*testAggregationSelector) AggregatorFor(desc *metric.Descriptor) export.Aggregator {
	switch desc.MetricKind() {
	case metric.CounterKind:
		return sum.New()
	case metric.ValueObserverKind:
		return lastvalue.New()
	default:
		panic("Invalid descriptor MetricKind for this test")
	}
}

func makeLabels(labels ...kv.KeyValue) *label.Set {
	s := label.NewSet(labels...)
	return &s
}

func (Encoder) Encode(iter label.Iterator) string {
	var sb strings.Builder
	for iter.Next() {
		i, l := iter.IndexedLabel()
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(string(l.Key))
		sb.WriteString("~")
		sb.WriteString(l.Value.Emit())
	}
	return sb.String()
}

func (Encoder) ID() label.EncoderID {
	return testLabelEncoderID
}

// LastValueAgg returns a checkpointed lastValue aggregator w/ the specified descriptor and value.
func LastValueAgg(desc *metric.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	gagg := lastvalue.New()
	_ = gagg.Update(ctx, metric.NewInt64Number(v), desc)
	gagg.Checkpoint(ctx, desc)
	return gagg
}

// Convenience method for building a test exported lastValue record.
func NewLastValueRecord(desc *metric.Descriptor, labels *label.Set, value int64) export.Record {
	return export.NewRecord(desc, labels, Resource, LastValueAgg(desc, value))
}

// Convenience method for building a test exported counter record.
func NewCounterRecord(desc *metric.Descriptor, labels *label.Set, value int64) export.Record {
	return export.NewRecord(desc, labels, Resource, CounterAgg(desc, value))
}

// CounterAgg returns a checkpointed counter aggregator w/ the specified descriptor and value.
func CounterAgg(desc *metric.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	cagg := sum.New()
	_ = cagg.Update(ctx, metric.NewInt64Number(v), desc)
	cagg.Checkpoint(ctx, desc)
	return cagg
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
