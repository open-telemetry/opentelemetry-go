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

package test

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

type (
	// Encoder is an alternate label encoder to validate grouping logic.
	Encoder struct{}

	// Output collects distinct metric/label set outputs.
	Output struct {
		Encoder label.Encoder
		Values  map[string]int64
	}

	// testAggregationSelector returns aggregators consistent with
	// the test variables below, needed for testing stateful
	// batchers, which clone Aggregators using AggregatorFor(desc).
	testAggregationSelector struct{}
)

var (
	// GaugeADesc and GaugeBDesc group by "G"
	GaugeADesc = export.NewDescriptor(
		"gauge.a", export.GaugeKind, []core.Key{key.New("G")}, "", "", core.Int64NumberKind, false)
	GaugeBDesc = export.NewDescriptor(
		"gauge.b", export.GaugeKind, []core.Key{key.New("G")}, "", "", core.Int64NumberKind, false)
	// CounterADesc and CounterBDesc group by "C"
	CounterADesc = export.NewDescriptor(
		"counter.a", export.CounterKind, []core.Key{key.New("C")}, "", "", core.Int64NumberKind, false)
	CounterBDesc = export.NewDescriptor(
		"counter.b", export.CounterKind, []core.Key{key.New("C")}, "", "", core.Int64NumberKind, false)

	// SdkEncoder uses a non-standard encoder like K1~V1&K2~V2
	SdkEncoder = &Encoder{}
	// GroupEncoder uses the SDK default encoder
	GroupEncoder = label.NewDefaultEncoder()

	// Gauge groups are (labels1), (labels2+labels3)
	// Counter groups are (labels1+labels2), (labels3)

	// Labels1 has G=H and C=D
	Labels1 = label.NewSet(key.String("G", "H"), key.String("C", "D"))
	// Labels2 has C=D and E=F
	Labels2 = label.NewSet(key.String("C", "D"), key.String("E", "F"))
	// Labels3 is the empty set
	Labels3 = label.NewSet()
)

// NewAggregationSelector returns a policy that is consistent with the
// test descriptors above.  I.e., it returns counter.New() for counter
// instruments and gauge.New for gauge instruments.
func NewAggregationSelector() export.AggregationSelector {
	return &testAggregationSelector{}
}

func (*testAggregationSelector) AggregatorFor(desc *export.Descriptor) export.Aggregator {
	switch desc.MetricKind() {
	case export.CounterKind:
		return counter.New()
	case export.GaugeKind:
		return gauge.New()
	default:
		panic("Invalid descriptor MetricKind for this test")
	}
}

func (Encoder) Encode(labels []core.KeyValue) string {
	var sb strings.Builder
	for i, l := range labels {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(string(l.Key))
		sb.WriteString("~")
		sb.WriteString(l.Value.Emit())
	}
	return sb.String()
}

// GaugeAgg returns a checkpointed gauge aggregator w/ the specified descriptor and value.
func GaugeAgg(desc *export.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	gagg := gauge.New()
	_ = gagg.Update(ctx, core.NewInt64Number(v), desc)
	gagg.Checkpoint(ctx, desc)
	return gagg
}

// Convenience method for building a test exported gauge record.
func NewGaugeRecord(desc *export.Descriptor, labels label.Set, value int64) export.Record {
	return export.NewRecord(desc, labels, GaugeAgg(desc, value))
}

// Convenience method for building a test exported counter record.
func NewCounterRecord(desc *export.Descriptor, labels label.Set, value int64) export.Record {
	return export.NewRecord(desc, labels, CounterAgg(desc, value))
}

// CounterAgg returns a checkpointed counter aggregator w/ the specified descriptor and value.
func CounterAgg(desc *export.Descriptor, v int64) export.Aggregator {
	ctx := context.Background()
	cagg := counter.New()
	_ = cagg.Update(ctx, core.NewInt64Number(v), desc)
	cagg.Checkpoint(ctx, desc)
	return cagg
}

func NewOutput(encoder label.Encoder) *Output {
	return &Output{
		Values:  map[string]int64{},
		Encoder: encoder,
	}
}

// AddTo adds a name/label-encoding entry with the gauge or counter
// value to the output map.
func (o Output) AddTo(rec export.Record) {
	labels := rec.Labels()
	key := fmt.Sprint(rec.Descriptor().Name(), "/", labels.Encoded(o.Encoder))
	var value int64
	switch t := rec.Aggregator().(type) {
	case *counter.Aggregator:
		sum, _ := t.Sum()
		value = sum.AsInt64()
	case *gauge.Aggregator:
		lv, _, _ := t.LastValue()
		value = lv.AsInt64()
	}
	o.Values[key] = value
}
