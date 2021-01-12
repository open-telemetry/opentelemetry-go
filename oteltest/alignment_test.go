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

package oteltest

import (
	"os"
	"testing"
	"unsafe"

	"go.opentelemetry.io/otel/internal/internaltest"
)

// TestMain ensures struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fields := []internaltest.FieldOffset{
		{
			Name:   "Batch.Measurments",
			Offset: unsafe.Offsetof(Batch{}.Measurements),
		},
		{
			Name:   "Measurement.Number",
			Offset: unsafe.Offsetof(Measurement{}.Number),
		},
	}
	if !internaltest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
