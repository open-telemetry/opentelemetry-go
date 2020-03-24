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

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/metric"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregator"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
)

type (
	// Encoder is an alternate label encoder to validate grouping logic.
	Encoder struct{}

	// Output collects distinct metric/label set outputs.
	Output map[string]float64

	// testAggregationSelector returns aggregators consistent with
	// the test variables below, needed for testing stateful
	// batchers, which clone Aggregators using AggregatorFor(desc).
	testAggregationSelector struct{}
)

var (
	// LastValueADesc and LastValueBDesc group by "G"
	LastValueADesc = metric.NewDescriptor(
		"lastvalue.a", metric.ObserverKind, core.Int64NumberKind, metric.WithKeys(key.New("G")))
	LastValueBDesc = metric.NewDescriptor(
		"lastvalue.b", metric.ObserverKind, core.Int64NumberKind, metric.WithKeys(key.New("G")))
	// CounterADesc and CounterBDesc group by "C"
	CounterADesc = metric.NewDescriptor(
		"sum.a", metric.CounterKind, core.Int64NumberKind, metric.WithKeys(key.New("C")))
	CounterBDesc = metric.NewDescriptor(
		"sum.b", metric.CounterKind, core.Int64NumberKind, metric.WithKeys(key.New("C")))

	// SdkEncoder uses a non-standard encoder like K1~V1&K2~V2
	SdkEncoder = &Encoder{}
	// GroupEncoder uses the SDK default encoder
	GroupEncoder = sdk.NewDefaultLabelEncoder()

	// LastValue groups are (labels1), (labels2+labels3)
	// Counter groups are (labels1+labels2), (labels3)

	// Labels1 has G=H and C=D
	Labels1 = makeLabels(SdkEncoder, key.String("G", "H"), key.String("C", "D"))
	// Labels2 has C=D and E=F
	Labels2 = makeLabels(SdkEncoder, key.String("C", "D"), key.String("E", "F"))
	// Labels3 is the empty set
	Labels3 = makeLabels(SdkEncoder)
)

// NewAggregationSelector returns a policy that is consistent with the
// test descriptors above.  I.e., it returns sum.New() for counter
// instruments and lastvalue.New for lastValue instruments.
func NewAggregationSelector() export.AggregationSelector {
	return &testAggregationSelector{}
}

func (*testAggregationSelector) AggregatorFor(desc *metric.Descriptor) export.Aggregator {
	switch desc.MetricKind() {
	case metric.CounterKind:
		return sum.New()
	case metric.ObserverKind:
		return lastvalue.New()
	default:
		panic("Invalid descriptor MetricKind for this test")
	}
}

func makeLabels(encoder export.LabelEncoder, labels ...core.KeyValue) export.Labels {
	ls := export.LabelSlice(labels)
	return export.NewLabels(ls, encoder.Encode(ls.Iter()), encoder)
}

func (Encoder) Encode(iter export.LabelIterator) string {
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

// LastValueAgg returns a checkpointed lastValue aggregator w/ the specified descriptor and value.
func LastValueAgg(desc *metric.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	gagg := lastvalue.New()
	_ = gagg.Update(ctx, core.NewInt64Number(v), desc)
	gagg.Checkpoint(ctx, desc)
	return gagg
}

// Convenience method for building a test exported lastValue record.
func NewLastValueRecord(desc *metric.Descriptor, labels export.Labels, value int64) export.Record {
	return export.NewRecord(desc, labels, LastValueAgg(desc, value))
}

// Convenience method for building a test exported counter record.
func NewCounterRecord(desc *metric.Descriptor, labels export.Labels, value int64) export.Record {
	return export.NewRecord(desc, labels, CounterAgg(desc, value))
}

// CounterAgg returns a checkpointed counter aggregator w/ the specified descriptor and value.
func CounterAgg(desc *metric.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	cagg := sum.New()
	_ = cagg.Update(ctx, core.NewInt64Number(v), desc)
	cagg.Checkpoint(ctx, desc)
	return cagg
}

// AddTo adds a name/label-encoding entry with the lastValue or counter
// value to the output map.
func (o Output) AddTo(rec export.Record) error {
	labels := rec.Labels()
	key := fmt.Sprint(rec.Descriptor().Name(), "/", labels.Encoded())
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
	o[key] = value
	return nil
}
