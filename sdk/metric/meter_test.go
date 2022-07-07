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

//go:build go1.18
// +build go1.18

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/instrumentation"
)

func TestMeterRegistry(t *testing.T) {
	il0 := instrumentation.Library{Name: "zero"}
	il1 := instrumentation.Library{Name: "one"}

	r := meterRegistry{}
	var m0 *meter
	t.Run("ZeroValueGetDoesNotPanic", func(t *testing.T) {
		assert.NotPanics(t, func() { m0 = r.Get(il0) })
		assert.Equal(t, il0, m0.Library, "uninitialized meter returned")
	})

	m01 := r.Get(il0)
	t.Run("GetSameMeter", func(t *testing.T) {
		assert.Samef(t, m0, m01, "returned different meters: %v", il0)
	})

	m1 := r.Get(il1)
	t.Run("GetDifferentMeter", func(t *testing.T) {
		assert.NotSamef(t, m0, m1, "returned same meters: %v", il1)
	})

	t.Run("RangeComplete", func(t *testing.T) {
		var got []*meter
		r.Range(func(m *meter) bool {
			got = append(got, m)
			return true
		})
		assert.ElementsMatch(t, []*meter{m0, m1}, got)
	})

	t.Run("RangeStopIteration", func(t *testing.T) {
		var i int
		r.Range(func(m *meter) bool {
			i++
			return false
		})
		assert.Equal(t, 1, i, "iteration not stopped after first flase return")
	})
}
