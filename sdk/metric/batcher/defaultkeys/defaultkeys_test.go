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

package defaultkeys_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/counter"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/gauge"
	"go.opentelemetry.io/otel/sdk/metric/batcher/defaultkeys"
)

type (
	testEncoder struct{}

	testOutput map[string]int64
)

var (
	// Gauge groups by "G"
	// Counter groups by "C"
	testGaugeDesc = export.NewDescriptor(
		"gauge", export.GaugeKind, []core.Key{key.New("G")}, "", "", core.Int64NumberKind, false)
	testCounterDesc = export.NewDescriptor(
		"counter", export.CounterKind, []core.Key{key.New("C")}, "", "", core.Int64NumberKind, false)

	// The SDK and the batcher use different encoders in these tests.
	sdkEncoder   = &testEncoder{}
	groupEncoder = sdk.DefaultLabelEncoder()

	// Gauge groups are (labels1), (labels2+labels3)
	// Counter groups are (labels1+labels2), (labels3)
	labels1 = makeLabels(sdkEncoder, key.String("G", "H"), key.String("C", "D"))
	labels2 = makeLabels(sdkEncoder, key.String("C", "D"), key.String("E", "F"))
	labels3 = makeLabels(sdkEncoder)
)

func makeLabels(encoder export.LabelEncoder, labels ...core.KeyValue) export.Labels {
	encoded := encoder.EncodeLabels(labels)
	return export.NewLabels(labels, encoded, encoder)
}

func (testEncoder) EncodeLabels(labels []core.KeyValue) string {
	return fmt.Sprint(labels)
}

func gaugeAgg(v int64) export.Aggregator {
	ctx := context.Background()
	gagg := gauge.New()
	_ = gagg.Update(ctx, core.NewInt64Number(v), testGaugeDesc)
	gagg.Checkpoint(ctx, testCounterDesc)
	return gagg
}

func counterAgg(v int64) export.Aggregator {
	ctx := context.Background()
	cagg := counter.New()
	_ = cagg.Update(ctx, core.NewInt64Number(v), testCounterDesc)
	cagg.Checkpoint(ctx, testCounterDesc)
	return cagg
}

func (o testOutput) addTo(rec export.Record) {
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

func TestGroupingStateless(t *testing.T) {
	ctx := context.Background()
	b := defaultkeys.New(nil, groupEncoder, false)

	b.Process(ctx, testGaugeDesc, labels1, gaugeAgg(10))
	b.Process(ctx, testGaugeDesc, labels2, gaugeAgg(20))
	b.Process(ctx, testGaugeDesc, labels3, gaugeAgg(30))

	b.Process(ctx, testCounterDesc, labels1, counterAgg(10))
	b.Process(ctx, testCounterDesc, labels2, counterAgg(20))
	b.Process(ctx, testCounterDesc, labels3, counterAgg(40))

	processor := b.ReadCheckpoint()

	records := testOutput{}
	processor.Foreach(records.addTo)

	// Output gauge should have only the "G=H" and "G=" keys.
	// Output counter should have only the "C=D" and "C=" keys.
	require.EqualValues(t, map[string]int64{
		"counter/C=D": 30, // labels1 + labels2
		"counter/C=":  40, // labels3
		"gauge/G=H":   10, // labels1
		"gauge/G=":    30, // labels3 = last value
	}, records)

	// Verify that state was reset
	processor = b.ReadCheckpoint()
	processor.Foreach(func(rec export.Record) {
		t.Fatal("Unexpected call")
	})
}

func TestGroupingStateful(t *testing.T) {
	ctx := context.Background()
	b := defaultkeys.New(nil, groupEncoder, true)

	b.Process(ctx, testCounterDesc, labels1, counterAgg(10))

	processor := b.ReadCheckpoint()

	records1 := testOutput{}
	processor.Foreach(records1.addTo)

	require.EqualValues(t, map[string]int64{
		"counter/C=D": 10, // labels1
	}, records1)

	// Test that state was NOT reset
	processor = b.ReadCheckpoint()

	records2 := testOutput{}
	processor.Foreach(records2.addTo)

	require.EqualValues(t, records1, records2)
}
