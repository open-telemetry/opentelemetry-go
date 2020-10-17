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

package simple_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	testCounterDesc           = otel.NewDescriptor("counter", otel.CounterInstrumentKind, otel.Int64NumberKind)
	testUpDownCounterDesc     = otel.NewDescriptor("updowncounter", otel.UpDownCounterInstrumentKind, otel.Int64NumberKind)
	testSumObserverDesc       = otel.NewDescriptor("sumobserver", otel.SumObserverInstrumentKind, otel.Int64NumberKind)
	testUpDownSumObserverDesc = otel.NewDescriptor("updownsumobserver", otel.UpDownSumObserverInstrumentKind, otel.Int64NumberKind)
	testValueRecorderDesc     = otel.NewDescriptor("valuerecorder", otel.ValueRecorderInstrumentKind, otel.Int64NumberKind)
	testValueObserverDesc     = otel.NewDescriptor("valueobserver", otel.ValueObserverInstrumentKind, otel.Int64NumberKind)
)

func oneAgg(sel export.AggregatorSelector, desc *otel.Descriptor) export.Aggregator {
	var agg export.Aggregator
	sel.AggregatorFor(desc, &agg)
	return agg
}

func testFixedSelectors(t *testing.T, sel export.AggregatorSelector) {
	require.IsType(t, (*lastvalue.Aggregator)(nil), oneAgg(sel, &testValueObserverDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testCounterDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testUpDownCounterDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testSumObserverDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testUpDownSumObserverDesc))
}

func TestInexpensiveDistribution(t *testing.T) {
	inex := simple.NewWithInexpensiveDistribution()
	require.IsType(t, (*minmaxsumcount.Aggregator)(nil), oneAgg(inex, &testValueRecorderDesc))
	testFixedSelectors(t, inex)
}

func TestSketchDistribution(t *testing.T) {
	sk := simple.NewWithSketchDistribution(ddsketch.NewDefaultConfig())
	require.IsType(t, (*ddsketch.Aggregator)(nil), oneAgg(sk, &testValueRecorderDesc))
	testFixedSelectors(t, sk)
}

func TestExactDistribution(t *testing.T) {
	ex := simple.NewWithExactDistribution()
	require.IsType(t, (*array.Aggregator)(nil), oneAgg(ex, &testValueRecorderDesc))
	testFixedSelectors(t, ex)
}

func TestHistogramDistribution(t *testing.T) {
	hist := simple.NewWithHistogramDistribution(nil)
	require.IsType(t, (*histogram.Aggregator)(nil), oneAgg(hist, &testValueRecorderDesc))
	testFixedSelectors(t, hist)
}
