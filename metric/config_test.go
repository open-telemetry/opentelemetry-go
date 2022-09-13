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

package metric // import "go.opentelemetry.io/otel/metric"

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
)

func TestMeterConfig(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		assert.Equal(t, NewMeterConfig(), MeterConfig{})
	})

	t.Run("InstrumentationVersion", func(t *testing.T) {
		v0, v1 := "v0.1.0", "v1.0.0"

		assert.Equal(t, NewMeterConfig(
			WithInstrumentationVersion(v0),
		).InstrumentationVersion(), v0)

		assert.Equal(t, NewMeterConfig(
			WithInstrumentationVersion(v0),
			WithInstrumentationVersion(v1),
		).InstrumentationVersion(), v1, "last option has precedence")
	})

	t.Run("SchemaURL", func(t *testing.T) {
		s120 := "https://opentelemetry.io/schemas/1.2.0"
		s130 := "https://opentelemetry.io/schemas/1.3.0"

		assert.Equal(t, NewMeterConfig(
			WithSchemaURL(s120),
		).SchemaURL(), s120)

		assert.Equal(t, NewMeterConfig(
			WithSchemaURL(s120),
			WithSchemaURL(s130),
		).SchemaURL(), s130, "last option has precedence")
	})

	t.Run("Attributes", func(t *testing.T) {
		one, two := attribute.Int("key", 1), attribute.Int("key", 2)

		assert.Equal(t, NewMeterConfig(
			WithScopeAttributes(one, two),
		).Attributes(), attribute.NewSet(two), "last attribute is used")

		assert.Equal(t, NewMeterConfig(
			WithScopeAttributes(two),
			WithScopeAttributes(one),
		).Attributes(), attribute.NewSet(one), "last option has precedence")
	})
}
