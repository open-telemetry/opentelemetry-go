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

package metric

import (
	"testing"
)

func TestGlobalMeter(t *testing.T) {
	m := GlobalMeter()
	if _, ok := m.(noopMeter); !ok {
		t.Errorf("Expected global meter to be a noopMeter instance, got an instance of %T", m)
	}

	SetGlobalMeter(newMockMeter())

	m = GlobalMeter()
	if _, ok := m.(*mockMeter); !ok {
		t.Errorf("Expected global meter to be a *mockMetric.MockMeter instance, got an instance of %T", m)
	}
}
