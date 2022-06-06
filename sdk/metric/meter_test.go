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

//go:build go1.17
// +build go1.17

package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

func TestMeterRegistryGetDoesNotPanicForZeroValue(t *testing.T) {
	r := meterRegistry{}
	assert.NotPanics(t, func() { _, _ = r.Get(instrumentation.Library{}) })
}

func TestMeterRegistry(t *testing.T) {
	il0 := instrumentation.Library{Name: "zero"}
	il1 := instrumentation.Library{Name: "one"}

	r := meterRegistry{}
	m0, ok := r.Get(il0)
	assert.Falsef(t, ok, "meter was in registry: %v", il0)

	m01, ok := r.Get(il0)
	assert.Truef(t, ok, "meter was not in registry: %v", il0)
	assert.Samef(t, m0, m01, "returned different meters: %v", il0)

	m1, ok := r.Get(il1)
	assert.Falsef(t, ok, "meter was in registry: %v", il1)
	assert.NotSamef(t, m0, m1, "returned same meters: %v", il1)

	var got []*meter
	r.Range(func(m *meter) bool {
		got = append(got, m)
		return true
	})
	assert.Equal(t, []*meter{m0, m1}, got)
}
