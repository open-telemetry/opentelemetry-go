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

package metric_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/metric"
)

var (
	syncKinds = []metric.Kind{
		metric.ValueRecorderKind,
		metric.CounterKind,
		metric.UpDownCounterKind,
	}
	asyncKinds = []metric.Kind{
		metric.ValueObserverKind,
		metric.SumObserverKind,
		metric.UpDownSumObserverKind,
	}
	addingKinds = []metric.Kind{
		metric.CounterKind,
		metric.UpDownCounterKind,
		metric.SumObserverKind,
		metric.UpDownSumObserverKind,
	}
	groupingKinds = []metric.Kind{
		metric.ValueRecorderKind,
		metric.ValueObserverKind,
	}

	monotonicKinds = []metric.Kind{
		metric.CounterKind,
		metric.SumObserverKind,
	}

	nonMonotonicKinds = []metric.Kind{
		metric.UpDownCounterKind,
		metric.UpDownSumObserverKind,
		metric.ValueRecorderKind,
		metric.ValueObserverKind,
	}

	precomputedSumKinds = []metric.Kind{
		metric.SumObserverKind,
		metric.UpDownSumObserverKind,
	}

	nonPrecomputedSumKinds = []metric.Kind{
		metric.CounterKind,
		metric.UpDownCounterKind,
		metric.ValueRecorderKind,
		metric.ValueObserverKind,
	}
)

func TestSynchronous(t *testing.T) {
	for _, k := range syncKinds {
		require.True(t, k.Synchronous())
		require.False(t, k.Asynchronous())
	}
	for _, k := range asyncKinds {
		require.True(t, k.Asynchronous())
		require.False(t, k.Synchronous())
	}
}

func TestGrouping(t *testing.T) {
	for _, k := range groupingKinds {
		require.True(t, k.Grouping())
		require.False(t, k.Adding())
	}
	for _, k := range addingKinds {
		require.True(t, k.Adding())
		require.False(t, k.Grouping())
	}
}

func TestMonotonic(t *testing.T) {
	for _, k := range monotonicKinds {
		require.True(t, k.Monotonic())
	}
	for _, k := range nonMonotonicKinds {
		require.False(t, k.Monotonic())
	}
}

func TestPrecomputedSum(t *testing.T) {
	for _, k := range precomputedSumKinds {
		require.True(t, k.PrecomputedSum())
	}
	for _, k := range nonPrecomputedSumKinds {
		require.False(t, k.PrecomputedSum())
	}
}
