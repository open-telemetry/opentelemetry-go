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

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/exact"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/lastvalue"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	testCounterDesc           = metric.NewDescriptor("counter", metric.CounterInstrumentKind, number.Int64Kind)
	testUpDownCounterDesc     = metric.NewDescriptor("updowncounter", metric.UpDownCounterInstrumentKind, number.Int64Kind)
	testSumObserverDesc       = metric.NewDescriptor("sumobserver", metric.SumObserverInstrumentKind, number.Int64Kind)
	testUpDownSumObserverDesc = metric.NewDescriptor("updownsumobserver", metric.UpDownSumObserverInstrumentKind, number.Int64Kind)
	testValueRecorderDesc     = metric.NewDescriptor("valuerecorder", metric.ValueRecorderInstrumentKind, number.Int64Kind)
	testValueObserverDesc     = metric.NewDescriptor("valueobserver", metric.ValueObserverInstrumentKind, number.Int64Kind)
)

func oneAgg(sel export.AggregatorSelector, desc *metric.Descriptor) export.Aggregator {
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

func TestExactDistribution(t *testing.T) {
	ex := simple.NewWithExactDistribution()
	require.IsType(t, (*exact.Aggregator)(nil), oneAgg(ex, &testValueRecorderDesc))
	testFixedSelectors(t, ex)
}

func TestHistogramDistribution(t *testing.T) {
	hist := simple.NewWithHistogramDistribution()
	require.IsType(t, (*histogram.Aggregator)(nil), oneAgg(hist, &testValueRecorderDesc))
	testFixedSelectors(t, hist)
}
