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

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
)

type (
	// Encoder is an alternate label encoder to validate grouping logic.
	Encoder struct{}

	// Output collects distinct metric/label set outputs.
	Output map[string]int64

	// testAggregationSelector returns aggregators consistent with
	// the test variables below, needed for testing stateful
	// batchers, which clone Aggregators using AggregatorFor(desc).
	testAggregationSelector struct{}
)

var (
	// GaugeDesc groups by "G"
	GaugeDesc = export.NewDescriptor(
		"gauge", export.GaugeKind, []core.Key{key.New("G")}, "", "", core.Int64NumberKind, false)
	// CounterDesc groups by "C"
	CounterDesc = export.NewDescriptor(
		"counter", export.CounterKind, []core.Key{key.New("C")}, "", "", core.Int64NumberKind, false)

	// SdkEncoder uses a non-standard encoder like K1~V1&K2~V2
	SdkEncoder = &Encoder{}
	// GroupEncoder uses the SDK default encoder
	GroupEncoder = sdk.DefaultLabelEncoder()

	// Gauge groups are (labels1), (labels2+labels3)
	// Counter groups are (labels1+labels2), (labels3)

	// Labels1 has G=H and C=D
	Labels1 = makeLabels(SdkEncoder, key.String("G", "H"), key.String("C", "D"))
	// Labels2 has C=D and E=F
	Labels2 = makeLabels(SdkEncoder, key.String("C", "D"), key.String("E", "F"))
	// Labels3 is the empty set
	Labels3 = makeLabels(SdkEncoder)
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

func makeLabels(encoder export.LabelEncoder, labels ...core.KeyValue) export.Labels {
	encoded := encoder.Encode(labels)
	return export.NewLabels(labels, encoded, encoder)
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

// GaugeAgg returns a checkpointed gauge aggregator w/ the specified value.
func GaugeAgg(v int64) export.Aggregator {
	ctx := context.Background()
	gagg := gauge.New()
	_ = gagg.Update(ctx, core.NewInt64Number(v), GaugeDesc)
	gagg.Checkpoint(ctx, CounterDesc)
	return gagg
}

// CounterAgg returns a checkpointed counter aggregator w/ the specified value.
func CounterAgg(v int64) export.Aggregator {
	ctx := context.Background()
	cagg := counter.New()
	_ = cagg.Update(ctx, core.NewInt64Number(v), CounterDesc)
	cagg.Checkpoint(ctx, CounterDesc)
	return cagg
}

// AddTo adds a name/label-encoding entry with the gauge or counter
// value to the output map.
func (o Output) AddTo(rec export.Record) {
	labels := rec.Labels()
	key := fmt.Sprint(rec.Descriptor().Name(), "/", labels.Encoded())
	var value int64
	switch t := rec.Aggregator().(type) {
	case *counter.Aggregator:
		sum, _ := t.Sum()
		value = sum.AsInt64()
	case *gauge.Aggregator:
		lv, _, _ := t.LastValue()
		value = lv.AsInt64()
	}
	o[key] = value
}
