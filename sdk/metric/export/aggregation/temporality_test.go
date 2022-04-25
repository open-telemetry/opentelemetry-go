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

package aggregation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
)

func TestTemporalityIncludes(t *testing.T) {
	require.True(t, CumulativeTemporality.Includes(CumulativeTemporality))
	require.True(t, DeltaTemporality.Includes(CumulativeTemporality|DeltaTemporality))
}

var deltaMemoryTemporalties = []sdkapi.InstrumentKind{
	sdkapi.CounterObserverInstrumentKind,
	sdkapi.UpDownCounterObserverInstrumentKind,
}

var cumulativeMemoryTemporalties = []sdkapi.InstrumentKind{
	sdkapi.HistogramInstrumentKind,
	sdkapi.GaugeObserverInstrumentKind,
	sdkapi.CounterInstrumentKind,
	sdkapi.UpDownCounterInstrumentKind,
}

func TestTemporalityMemoryRequired(t *testing.T) {
	for _, kind := range deltaMemoryTemporalties {
		require.True(t, DeltaTemporality.MemoryRequired(kind))
		require.False(t, CumulativeTemporality.MemoryRequired(kind))
	}

	for _, kind := range cumulativeMemoryTemporalties {
		require.True(t, CumulativeTemporality.MemoryRequired(kind))
		require.False(t, DeltaTemporality.MemoryRequired(kind))
	}
}

func TestTemporalitySelectors(t *testing.T) {
	cAggTemp := CumulativeTemporalitySelector()
	dAggTemp := DeltaTemporalitySelector()
	sAggTemp := StatelessTemporalitySelector()

	for _, ikind := range append(deltaMemoryTemporalties, cumulativeMemoryTemporalties...) {
		desc := sdkapi.NewDescriptor("instrument", ikind, number.Int64Kind, "", "")

		var akind Kind
		if ikind.Adding() {
			akind = SumKind
		} else {
			akind = HistogramKind
		}
		require.Equal(t, CumulativeTemporality, cAggTemp.TemporalityFor(&desc, akind))
		require.Equal(t, DeltaTemporality, dAggTemp.TemporalityFor(&desc, akind))
		require.False(t, sAggTemp.TemporalityFor(&desc, akind).MemoryRequired(ikind))
	}
}
