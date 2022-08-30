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

package tracetransform // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/tracetransform"

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
)

func TestInstrumentationScope(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		assert.Nil(t, InstrumentationScope(instrumentation.Scope{}))
	})

	t.Run("Mapping", func(t *testing.T) {
		var (
			name    = "instrumentation name"
			version = "v0.1.0"
			attr    = attribute.NewSet(attribute.String("domain", "trace"))
			attrPb  = Iterator(attr.Iter())
		)
		expected := &commonpb.InstrumentationScope{
			Name:       name,
			Version:    version,
			Attributes: attrPb,
		}
		actual := InstrumentationScope(instrumentation.Scope{
			Name:       name,
			Version:    version,
			SchemaURL:  "http://this.is.mapped.elsewhere.com",
			Attributes: attr,
		})
		assert.Equal(t, expected, actual)
	})
}
