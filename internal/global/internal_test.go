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

package global_test

import (
	"os"
	"testing"

	"go.opentelemetry.io/otel/internal/global"
	ottest "go.opentelemetry.io/otel/internal/internaltest"
)

// Ensure struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fieldsMap := global.AtomicFieldOffsets()
	fields := make([]ottest.FieldOffset, 0, len(fieldsMap))
	for name, offset := range fieldsMap {
		fields = append(fields, ottest.FieldOffset{
			Name:   name,
			Offset: offset,
		})
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
