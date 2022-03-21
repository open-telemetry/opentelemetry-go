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

package sdkapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric/number"
)

func TestMeasurementGetters(t *testing.T) {
	num := number.NewFloat64Number(1.5)
	si := NewNoopSyncInstrument()
	meas := NewMeasurement(si, num)

	require.Equal(t, si, meas.SyncImpl())
	require.Equal(t, num, meas.Number())
}

func TestObservationGetters(t *testing.T) {
	num := number.NewFloat64Number(1.5)
	ai := NewNoopAsyncInstrument()
	obs := NewObservation(ai, num)

	require.Equal(t, ai, obs.AsyncImpl())
	require.Equal(t, num, obs.Number())
}
