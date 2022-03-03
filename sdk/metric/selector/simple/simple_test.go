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

	"go.opentelemetry.io/otel/sdk/metric/aggregator"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/metrictest"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	testCounterDesc               = metrictest.NewDescriptor("counter", sdkapi.CounterInstrumentKind, number.Int64Kind)
	testUpDownCounterDesc         = metrictest.NewDescriptor("updowncounter", sdkapi.UpDownCounterInstrumentKind, number.Int64Kind)
	testCounterObserverDesc       = metrictest.NewDescriptor("counterobserver", sdkapi.CounterObserverInstrumentKind, number.Int64Kind)
	testUpDownCounterObserverDesc = metrictest.NewDescriptor("updowncounterobserver", sdkapi.UpDownCounterObserverInstrumentKind, number.Int64Kind)
	testHistogramDesc             = metrictest.NewDescriptor("histogram", sdkapi.HistogramInstrumentKind, number.Int64Kind)
	testGaugeObserverDesc         = metrictest.NewDescriptor("gauge", sdkapi.GaugeObserverInstrumentKind, number.Int64Kind)
)

func oneAgg(sel export.AggregatorSelector, desc *sdkapi.Descriptor) aggregator.Aggregator {
	var agg aggregator.Aggregator
	sel.AggregatorFor(desc, &agg)
	return agg
}

func testFixedSelectors(t *testing.T, sel export.AggregatorSelector) {
	require.IsType(t, (*lastvalue.Aggregator)(nil), oneAgg(sel, &testGaugeObserverDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testCounterDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testUpDownCounterDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testCounterObserverDesc))
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(sel, &testUpDownCounterObserverDesc))
}

func TestInexpensiveDistribution(t *testing.T) {
	inex := simple.NewWithInexpensiveDistribution()
	require.IsType(t, (*sum.Aggregator)(nil), oneAgg(inex, &testHistogramDesc))
	testFixedSelectors(t, inex)
}

func TestHistogramDistribution(t *testing.T) {
	hist := simple.NewWithHistogramDistribution()
	require.IsType(t, (*histogram.Aggregator)(nil), oneAgg(hist, &testHistogramDesc))
	testFixedSelectors(t, hist)
}
