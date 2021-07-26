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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	"testing"

	"go.opentelemetry.io/otel/internal/internaltest"
)

// Harness is a testing harness used to test implementations of the
// OpenTelemetry API.
//
// Deprecated: this will be removed in the next major release.
type Harness = internaltest.Harness

// NewHarness returns an instantiated *Harness using t.
//
// Deprecated: this will be removed in the next major release.
func NewHarness(t *testing.T) *Harness {
	return internaltest.NewHarness(t)
}
