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

	"go.opentelemetry.io/otel/metric/metrictest"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

func TestTemporalityIncludes(t *testing.T) {
	require.True(t, CumulativeTemporality.Includes(CumulativeTemporality))
	require.True(t, DeltaTemporality.Includes(CumulativeTemporality|DeltaTemporality))
}

var deltaMemoryKinds = []sdkapi.InstrumentKind{
	sdkapi.CounterObserverInstrumentKind,
	sdkapi.UpDownCounterObserverInstrumentKind,
}

var cumulativeMemoryKinds = []sdkapi.InstrumentKind{
	sdkapi.HistogramInstrumentKind,
	sdkapi.GaugeObserverInstrumentKind,
	sdkapi.CounterInstrumentKind,
	sdkapi.UpDownCounterInstrumentKind,
}

func TestTemporalityMemoryRequired(t *testing.T) {
	for _, kind := range deltaMemoryKinds {
		require.True(t, DeltaTemporality.MemoryRequired(kind))
		require.False(t, CumulativeTemporality.MemoryRequired(kind))
	}

	for _, kind := range cumulativeMemoryKinds {
		require.True(t, CumulativeTemporality.MemoryRequired(kind))
		require.False(t, DeltaTemporality.MemoryRequired(kind))
	}
}

func TestTemporalitySelectors(t *testing.T) {
	ceks := CumulativeTemporalitySelector()
	deks := DeltaTemporalitySelector()
	seks := StatelessTemporalitySelector()

	for _, ikind := range append(deltaMemoryKinds, cumulativeMemoryKinds...) {
		desc := metrictest.NewDescriptor("instrument", ikind, number.Int64Kind)

		var akind aggregation.Kind
		if ikind.Adding() {
			akind = aggregation.SumKind
		} else {
			akind = aggregation.HistogramKind
		}
		require.Equal(t, CumulativeTemporality, ceks.TemporalityFor(&desc, akind))
		require.Equal(t, DeltaTemporality, deks.TemporalityFor(&desc, akind))
		require.False(t, seks.TemporalityFor(&desc, akind).MemoryRequired(ikind))
	}
}
