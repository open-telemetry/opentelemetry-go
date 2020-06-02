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

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	testCounterDesc       = metric.NewDescriptor("counter", metric.CounterKind, metric.Int64NumberKind)
	testValueRecorderDesc = metric.NewDescriptor("valuerecorder", metric.ValueRecorderKind, metric.Int64NumberKind)
	testValueObserverDesc = metric.NewDescriptor("valueobserver", metric.ValueObserverKind, metric.Int64NumberKind)
)

func TestInexpensiveDistribution(t *testing.T) {
	inex := simple.NewWithInexpensiveDistribution()
	require.NotPanics(t, func() { _ = inex.AggregatorFor(&testCounterDesc).(*sum.Aggregator) })
	require.NotPanics(t, func() { _ = inex.AggregatorFor(&testValueRecorderDesc).(*minmaxsumcount.Aggregator) })
	require.NotPanics(t, func() { _ = inex.AggregatorFor(&testValueObserverDesc).(*minmaxsumcount.Aggregator) })
}

func TestSketchDistribution(t *testing.T) {
	sk := simple.NewWithSketchDistribution(ddsketch.NewDefaultConfig())
	require.NotPanics(t, func() { _ = sk.AggregatorFor(&testCounterDesc).(*sum.Aggregator) })
	require.NotPanics(t, func() { _ = sk.AggregatorFor(&testValueRecorderDesc).(*ddsketch.Aggregator) })
	require.NotPanics(t, func() { _ = sk.AggregatorFor(&testValueObserverDesc).(*ddsketch.Aggregator) })
}

func TestExactDistribution(t *testing.T) {
	ex := simple.NewWithExactDistribution()
	require.NotPanics(t, func() { _ = ex.AggregatorFor(&testCounterDesc).(*sum.Aggregator) })
	require.NotPanics(t, func() { _ = ex.AggregatorFor(&testValueRecorderDesc).(*array.Aggregator) })
	require.NotPanics(t, func() { _ = ex.AggregatorFor(&testValueObserverDesc).(*array.Aggregator) })
}

func TestHistogramDistribution(t *testing.T) {
	ex := simple.NewWithHistogramDistribution(nil)
	require.NotPanics(t, func() { _ = ex.AggregatorFor(&testCounterDesc).(*sum.Aggregator) })
	require.NotPanics(t, func() { _ = ex.AggregatorFor(&testValueRecorderDesc).(*histogram.Aggregator) })
	require.NotPanics(t, func() { _ = ex.AggregatorFor(&testValueObserverDesc).(*histogram.Aggregator) })
}
