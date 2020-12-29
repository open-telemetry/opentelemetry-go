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

package metric

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

func TestAggregationTemporalityIncludes(t *testing.T) {
	require.True(t, CumulativeAggregationTemporality.Includes(CumulativeAggregationTemporality))
	require.True(t, DeltaAggregationTemporality.Includes(CumulativeAggregationTemporality|DeltaAggregationTemporality))
}

var deltaMemoryKinds = []metric.InstrumentKind{
	metric.SumObserverInstrumentKind,
	metric.UpDownSumObserverInstrumentKind,
}

var cumulativeMemoryKinds = []metric.InstrumentKind{
	metric.ValueRecorderInstrumentKind,
	metric.ValueObserverInstrumentKind,
	metric.CounterInstrumentKind,
	metric.UpDownCounterInstrumentKind,
}

func TestAggregationTemporalityMemoryRequired(t *testing.T) {
	for _, kind := range deltaMemoryKinds {
		require.True(t, DeltaAggregationTemporality.MemoryRequired(kind))
		require.False(t, CumulativeAggregationTemporality.MemoryRequired(kind))
	}

	for _, kind := range cumulativeMemoryKinds {
		require.True(t, CumulativeAggregationTemporality.MemoryRequired(kind))
		require.False(t, DeltaAggregationTemporality.MemoryRequired(kind))
	}
}

func TestAggregationTemporalitySelectors(t *testing.T) {
	ceks := CumulativeAggregationTemporalitySelector()
	seks := StatelessAggregationTemporalitySelector()

	for _, ikind := range append(deltaMemoryKinds, cumulativeMemoryKinds...) {
		desc := metric.NewDescriptor("instrument", ikind, number.Int64Kind)

		var akind aggregation.Kind
		if ikind.Adding() {
			akind = aggregation.SumKind
		} else {
			akind = aggregation.HistogramKind
		}
		require.Equal(t, CumulativeAggregationTemporality, ceks.AggregationTemporalityFor(&desc, akind))
		require.False(t, seks.AggregationTemporalityFor(&desc, akind).MemoryRequired(ikind))
	}
}
