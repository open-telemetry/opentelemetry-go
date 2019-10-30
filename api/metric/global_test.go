// Copyright 2019, OpenTelemetry Authors
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

	"go.opentelemetry.io/api/metric"
	mock "go.opentelemetry.io/internal/metric"
)

func TestGlobalMeter(t *testing.T) {
	m := metric.GlobalMeter()
	if _, ok := m.(metric.NoopMeter); !ok {
		t.Errorf("Expected global meter to be a NoopMeter instance, got an instance of %T", m)
	}

	metric.SetGlobalMeter(mock.NewMeter())

	m = metric.GlobalMeter()
	if _, ok := m.(*mock.Meter); !ok {
		t.Errorf("Expected global meter to be a *mock.Meter instance, got an instance of %T", m)
	}
}
