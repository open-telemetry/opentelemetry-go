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

package otel_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
)

var (
	syncKinds = []otel.Kind{
		otel.ValueRecorderKind,
		otel.CounterKind,
		otel.UpDownCounterKind,
	}
	asyncKinds = []otel.Kind{
		otel.ValueObserverKind,
		otel.SumObserverKind,
		otel.UpDownSumObserverKind,
	}
	addingKinds = []otel.Kind{
		otel.CounterKind,
		otel.UpDownCounterKind,
		otel.SumObserverKind,
		otel.UpDownSumObserverKind,
	}
	groupingKinds = []otel.Kind{
		otel.ValueRecorderKind,
		otel.ValueObserverKind,
	}

	monotonicKinds = []otel.Kind{
		otel.CounterKind,
		otel.SumObserverKind,
	}

	nonMonotonicKinds = []otel.Kind{
		otel.UpDownCounterKind,
		otel.UpDownSumObserverKind,
		otel.ValueRecorderKind,
		otel.ValueObserverKind,
	}

	precomputedSumKinds = []otel.Kind{
		otel.SumObserverKind,
		otel.UpDownSumObserverKind,
	}

	nonPrecomputedSumKinds = []otel.Kind{
		otel.CounterKind,
		otel.UpDownCounterKind,
		otel.ValueRecorderKind,
		otel.ValueObserverKind,
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
