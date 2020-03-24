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

package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestInflightOperationMustEndBeforeSwap(t *testing.T) {
	var swapped bool
	ch := make(chan struct{})

	l := StateLocker{}
	op1 := l.Start()

	go func() {
		l.SwapActiveState(func() {})
		swapped = true
		ch <- struct{}{}
	}()

	require.False(t, swapped, "Swap should wait the end of the in-flight operation.")

	l.End(op1)

	select {
	case <-ch:
		require.True(t, swapped, "Swap should've been completed. ")
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Swap was not concluded after 50 milliseconds.")
	}
}

func TestEnsureIndexIsConsistent(t *testing.T) {
	l := StateLocker{}
	op1 := l.Start()
	l.End(op1)

	l.SwapActiveState(func() {})

	op2 := l.Start()
	l.End(op2)

	op3 := l.Start()
	l.End(op3)

	l.SwapActiveState(func() {})

	op4 := l.Start()
	l.End(op4)

	require.Equal(t, op1, op4, "two operations separated by two swaps should have the same index.")
	require.Equal(t, op2, op3, "two operations with no swap in between should have the same index.")

	require.Equal(t, 0, op1, "first index should be 0")
	require.Equal(t, 1, op2, "second index should be 1")
}

func TestTwoSwapsCanHappenWithoutOperationsInBetween(t *testing.T) {
	l := StateLocker{}

	require.Equal(t, 1, l.ColdIdx(), "first cold index should be 1")
	l.SwapActiveState(func() {})
	require.Equal(t, 0, l.ColdIdx(), "second cold index should be 0")
	l.SwapActiveState(func() {})
	require.Equal(t, 1, l.ColdIdx(), "third cold index should be 1")
}

func BenchmarkStateLocker_StartEnd(b *testing.B) {
	l := StateLocker{}

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		l.End(l.Start())
	}
}

func BenchmarkStateLocker_SwapActiveState(b *testing.B) {

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		l := StateLocker{}
		l.SwapActiveState(func() {})
	}
}
