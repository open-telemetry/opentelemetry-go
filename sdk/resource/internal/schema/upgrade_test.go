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

package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
	"go.opentelemetry.io/otel/schema/v1.1/types"
)

var schema = &ast.Schema{
	FileFormat: "1.1.0",
	SchemaURL:  "https://opentelemetry.io/schemas/1.3.0",
	Versions: map[types.TelemetryVersion]ast.VersionDef{
		"1.3.0": {
			// v1.1.0 should be applied before this, making it a no-op.
			All: ast10.Attributes{
				Changes: []ast10.AttributeChange{
					{
						RenameAttributes: &ast10.RenameAttributes{
							AttributeMap: ast10.AttributeMap{"foo": "v1.3.0"},
						},
					},
				},
			},
		},
		"1.2.0": {
			// This should not apply to a resource.
			Spans: ast10.Spans{
				Changes: []ast10.SpansChange{
					{
						RenameAttributes: &ast10.AttributeMapForSpans{
							AttributeMap: ast10.AttributeMap{"qux": "v1.2.0"},
						},
					},
				},
			},
		},
		"1.1.0": {
			All: ast10.Attributes{
				Changes: []ast10.AttributeChange{
					{
						RenameAttributes: &ast10.RenameAttributes{
							AttributeMap: ast10.AttributeMap{"foo": "bar"},
						},
					},
				},
			},
			Resources: ast10.Attributes{
				Changes: []ast10.AttributeChange{
					// These are expected to be applied in order.
					{
						RenameAttributes: &ast10.RenameAttributes{
							AttributeMap: ast10.AttributeMap{"bar": "baz"},
						},
					},
					{
						RenameAttributes: &ast10.RenameAttributes{
							AttributeMap: ast10.AttributeMap{"baz": "qux"},
						},
					},
				},
			},
		},
		// Handle empty
		"1.0.0": {},
	},
}

func v0Attr() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Bool("foo", true),
		attribute.Bool("untouched", true),
	}
}

func v3Attr() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.Bool("qux", true),
		attribute.Bool("untouched", true),
	}
}

func TestUpgrade(t *testing.T) {
	attr := v0Attr()
	err := Upgrade(schema, attr)
	require.NoError(t, err)
	assert.Equal(t, v3Attr(), attr)
}

func TestDowngrade(t *testing.T) {
	attr := v3Attr()
	err := Downgrade(schema, "https://opentelemetry.io/schemas/0.9.0", attr)
	require.NoError(t, err)

	assert.Equal(t, v0Attr(), attr)
}
