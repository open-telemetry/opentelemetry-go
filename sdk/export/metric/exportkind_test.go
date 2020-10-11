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

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
)

func TestExportKindIdentity(t *testing.T) {
	akind := aggregation.Kind("Noop")

	require.Equal(t, CumulativeExporter, CumulativeExporter.ExportKindFor(nil, akind))
	require.Equal(t, DeltaExporter, DeltaExporter.ExportKindFor(nil, akind))
	require.Equal(t, PassThroughExporter, PassThroughExporter.ExportKindFor(nil, akind))
}

func TestExportKindIncludes(t *testing.T) {
	require.True(t, CumulativeExporter.Includes(CumulativeExporter))
	require.True(t, DeltaExporter.Includes(CumulativeExporter|DeltaExporter))
	require.False(t, DeltaExporter.Includes(PassThroughExporter|CumulativeExporter))
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

func TestExportKindMemoryRequired(t *testing.T) {
	for _, kind := range deltaMemoryKinds {
		require.True(t, DeltaExporter.MemoryRequired(kind))
		require.False(t, CumulativeExporter.MemoryRequired(kind))
		require.False(t, PassThroughExporter.MemoryRequired(kind))
	}

	for _, kind := range cumulativeMemoryKinds {
		require.True(t, CumulativeExporter.MemoryRequired(kind))
		require.False(t, DeltaExporter.MemoryRequired(kind))
		require.False(t, PassThroughExporter.MemoryRequired(kind))
	}
}
