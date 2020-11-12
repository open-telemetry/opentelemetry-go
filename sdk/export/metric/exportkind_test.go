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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric/number"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

func TestExportKindIncludes(t *testing.T) {
	require.True(t, CumulativeExportKind.Includes(CumulativeExportKind))
	require.True(t, DeltaExportKind.Includes(CumulativeExportKind|DeltaExportKind))
}

var deltaMemoryKinds = []otel.InstrumentKind{
	otel.SumObserverInstrumentKind,
	otel.UpDownSumObserverInstrumentKind,
}

var cumulativeMemoryKinds = []otel.InstrumentKind{
	otel.ValueRecorderInstrumentKind,
	otel.ValueObserverInstrumentKind,
	otel.CounterInstrumentKind,
	otel.UpDownCounterInstrumentKind,
}

func TestExportKindMemoryRequired(t *testing.T) {
	for _, kind := range deltaMemoryKinds {
		require.True(t, DeltaExportKind.MemoryRequired(kind))
		require.False(t, CumulativeExportKind.MemoryRequired(kind))
	}

	for _, kind := range cumulativeMemoryKinds {
		require.True(t, CumulativeExportKind.MemoryRequired(kind))
		require.False(t, DeltaExportKind.MemoryRequired(kind))
	}
}

func TestExportKindSelectors(t *testing.T) {
	ceks := CumulativeExportKindSelector()
	deks := DeltaExportKindSelector()
	seks := StatelessExportKindSelector()

	for _, ikind := range append(deltaMemoryKinds, cumulativeMemoryKinds...) {
		desc := otel.NewDescriptor("instrument", ikind, number.Int64Kind)

		var akind aggregation.Kind
		if ikind.Adding() {
			akind = aggregation.SumKind
		} else {
			akind = aggregation.HistogramKind
		}
		require.Equal(t, CumulativeExportKind, ceks.ExportKindFor(&desc, akind))
		require.Equal(t, DeltaExportKind, deks.ExportKindFor(&desc, akind))
		require.False(t, seks.ExportKindFor(&desc, akind).MemoryRequired(ikind))
	}
}
