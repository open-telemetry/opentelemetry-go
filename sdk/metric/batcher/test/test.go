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
	Encoder struct{}

	Output map[string]int64
)

var (
	// Gauge groups by "G"
	// Counter groups by "C"
	GaugeDesc = export.NewDescriptor(
		"gauge", export.GaugeKind, []core.Key{key.New("G")}, "", "", core.Int64NumberKind, false)
	CounterDesc = export.NewDescriptor(
		"counter", export.CounterKind, []core.Key{key.New("C")}, "", "", core.Int64NumberKind, false)

	// The SDK and the batcher use different encoders in these tests.
	SdkEncoder   = &Encoder{}
	GroupEncoder = sdk.DefaultLabelEncoder()

	// Gauge groups are (labels1), (labels2+labels3)
	// Counter groups are (labels1+labels2), (labels3)
	Labels1 = makeLabels(SdkEncoder, key.String("G", "H"), key.String("C", "D"))
	Labels2 = makeLabels(SdkEncoder, key.String("C", "D"), key.String("E", "F"))
	Labels3 = makeLabels(SdkEncoder)
)

func makeLabels(encoder export.LabelEncoder, labels ...core.KeyValue) export.Labels {
	encoded := encoder.EncodeLabels(labels)
	return export.NewLabels(labels, encoded, encoder)
}

func (Encoder) EncodeLabels(labels []core.KeyValue) string {
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

func GaugeAgg(v int64) export.Aggregator {
	ctx := context.Background()
	gagg := gauge.New()
	_ = gagg.Update(ctx, core.NewInt64Number(v), GaugeDesc)
	gagg.Checkpoint(ctx, CounterDesc)
	return gagg
}

func CounterAgg(v int64) export.Aggregator {
	ctx := context.Background()
	cagg := counter.New()
	_ = cagg.Update(ctx, core.NewInt64Number(v), CounterDesc)
	cagg.Checkpoint(ctx, CounterDesc)
	return cagg
}

func (o Output) AddTo(rec export.Record) {
	labels := rec.Labels()
	key := fmt.Sprint(rec.Descriptor().Name(), "/", labels.Encoded())
	var value int64
	switch t := rec.Aggregator().(type) {
	case *counter.Aggregator:
		value = t.Sum().AsInt64()
	case *gauge.Aggregator:
		value = t.LastValue().AsInt64()
	}
	o[key] = value
}
