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

package resource

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
	"go.opentelemetry.io/otel/schema/v1.1/types"
	"go.opentelemetry.io/otel/sdk/resource/internal/schema"
)

var (
	testURLV3 = "testing-URL/1.3.0"
	testURLV2 = "testing-URL/1.2.0"
	testURLV1 = "testing-URL/1.1.0"
	testURLV0 = "testing-URL/1.0.0"

	testSchemas = map[string]*ast.Schema{
		testURLV3: {
			FileFormat: "1.1.0",
			SchemaURL:  testURLV3,
			Versions: map[types.TelemetryVersion]ast.VersionDef{
				types.TelemetryVersion("1.3.0"): {
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
				types.TelemetryVersion("1.2.0"): {
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
				types.TelemetryVersion("1.1.0"): {
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
				types.TelemetryVersion("1.0.0"): {},
			},
		},
	}
)

func TestConverterResource(t *testing.T) {
	v0 := NewWithAttributes(
		testURLV0,
		attribute.Bool("foo", true),
		attribute.Bool("untouched", true),
	)

	v3 := NewWithAttributes(
		testURLV3,
		attribute.Bool("qux", true),
		attribute.Bool("untouched", true),
	)

	c := newConverter(schema.NewStaticClient(testSchemas))
	ctx := context.Background()

	t.Run("Upgrade", func(t *testing.T) {
		got, err := c.Resource(ctx, testURLV3, v0)
		require.NoError(t, err)
		assert.Equal(t, v3, got)
	})

	t.Run("Downgrade", func(t *testing.T) {
		got, err := c.Resource(ctx, testURLV0, v3)
		require.NoError(t, err)
		assert.Equal(t, v0, got)
	})
}

func TestConverterMergeResources(t *testing.T) {
	v0 := NewWithAttributes(testURLV0, attribute.Bool("foo", true))
	v1 := NewWithAttributes(testURLV1, attribute.Bool("untranslated", false))
	v2 := NewWithAttributes(testURLV2, attribute.Bool("untranslated", true))

	v3 := NewWithAttributes(
		testURLV3,
		attribute.Bool("qux", true),
		attribute.Bool("untranslated", true),
	)

	c := newConverter(schema.NewStaticClient(testSchemas))
	ctx := context.Background()

	got, err := c.MergeResources(ctx, testURLV3, v0, v1, v2)
	require.NoError(t, err)
	assert.Equal(t, v3, got)
}
