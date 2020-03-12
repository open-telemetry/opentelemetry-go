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

package simple_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/array"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/ddsketch"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

var (
	testCounterDesc  = export.NewDescriptor("counter", export.CounterKind, nil, "", "", core.Int64NumberKind)
	testMeasureDesc  = export.NewDescriptor("measure", export.MeasureKind, nil, "", "", core.Int64NumberKind)
	testObserverDesc = export.NewDescriptor("observer", export.ObserverKind, nil, "", "", core.Int64NumberKind)
)

func TestInexpensiveMeasure(t *testing.T) {
	inex := simple.NewWithInexpensiveMeasure()
	require.NotPanics(t, func() { _ = inex.AggregatorFor(testCounterDesc).(*sum.Aggregator) })
	require.NotPanics(t, func() { _ = inex.AggregatorFor(testMeasureDesc).(*minmaxsumcount.Aggregator) })
	require.NotPanics(t, func() { _ = inex.AggregatorFor(testObserverDesc).(*minmaxsumcount.Aggregator) })
}

func TestSketchMeasure(t *testing.T) {
	sk := simple.NewWithSketchMeasure(ddsketch.NewDefaultConfig())
	require.NotPanics(t, func() { _ = sk.AggregatorFor(testCounterDesc).(*sum.Aggregator) })
	require.NotPanics(t, func() { _ = sk.AggregatorFor(testMeasureDesc).(*ddsketch.Aggregator) })
	require.NotPanics(t, func() { _ = sk.AggregatorFor(testObserverDesc).(*ddsketch.Aggregator) })
}

func TestExactMeasure(t *testing.T) {
	ex := simple.NewWithExactMeasure()
	require.NotPanics(t, func() { _ = ex.AggregatorFor(testCounterDesc).(*sum.Aggregator) })
	require.NotPanics(t, func() { _ = ex.AggregatorFor(testMeasureDesc).(*array.Aggregator) })
	require.NotPanics(t, func() { _ = ex.AggregatorFor(testObserverDesc).(*array.Aggregator) })
}
