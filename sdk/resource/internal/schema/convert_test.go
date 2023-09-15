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

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/schema/v1.0/ast"
)

func TestUpgrade(t *testing.T) {
	tforms := []transform{
		// Handle empty
		{Version: semver.New(0, 8, 0, "", "")},
		{Version: semver.New(0, 9, 0, "", "")},
		{Version: semver.New(1, 0, 0, "", "")},
		{
			Version: semver.New(1, 1, 0, "", ""),
			All: ast.Attributes{
				Changes: []ast.AttributeChange{
					{
						RenameAttributes: &ast.RenameAttributes{
							AttributeMap: ast.AttributeMap{"foo": "bar"},
						},
					},
				},
			},
			Resources: ast.Attributes{
				Changes: []ast.AttributeChange{
					// These are expected to be applied in order.
					{
						RenameAttributes: &ast.RenameAttributes{
							AttributeMap: ast.AttributeMap{"bar": "baz"},
						},
					},
					{
						RenameAttributes: &ast.RenameAttributes{
							AttributeMap: ast.AttributeMap{"baz": "qux"},
						},
					},
				},
			},
		},
		{Version: semver.New(1, 2, 0, "", "")},
		{
			Version: semver.New(1, 3, 0, "", ""),
			// v1.1.0 should be applied before this, making it a no-op.
			All: ast.Attributes{
				Changes: []ast.AttributeChange{
					{
						RenameAttributes: &ast.RenameAttributes{
							AttributeMap: ast.AttributeMap{"foo": "v1.3.0"},
						},
					},
				},
			},
		},
	}

	attr := []attribute.KeyValue{
		attribute.Bool("foo", true),
		attribute.Bool("untouched", true),
	}
	err := upgrade(tforms, attr)
	require.NoError(t, err)

	want := []attribute.KeyValue{
		attribute.Bool("qux", true),
		attribute.Bool("untouched", true),
	}
	assert.Equal(t, want, attr)
}

func TestSlice(t *testing.T) {
	tforms := []transform{
		{Version: semver.New(1, 4, 0, "", "")},
		{Version: semver.New(1, 4, 1, "", "")},
		{Version: semver.New(1, 5, 0, "", "")},
		{Version: semver.New(1, 6, 1, "", "")},
		{Version: semver.New(1, 7, 0, "", "")},
		{Version: semver.New(1, 8, 0, "", "")},
	}

	testcases := []struct {
		min, max *semver.Version
		want     []transform
	}{
		{
			min:  semver.New(1, 4, 0, "", ""),
			max:  semver.New(1, 8, 0, "", ""),
			want: tforms,
		},
		{
			min:  semver.New(1, 4, 1, "", ""),
			max:  semver.New(1, 8, 0, "", ""),
			want: tforms[1:],
		},
		{
			min:  semver.New(1, 4, 2, "", ""),
			max:  semver.New(1, 8, 0, "", ""),
			want: tforms[2:],
		},
		{
			min:  semver.New(1, 4, 2, "", ""),
			max:  semver.New(1, 7, 0, "", ""),
			want: tforms[2 : len(tforms)-1],
		},
		{
			min:  semver.New(1, 4, 2, "", ""),
			max:  semver.New(1, 6, 3, "", ""),
			want: tforms[2 : len(tforms)-2],
		},
		{
			min:  semver.New(1, 8, 0, "", ""),
			max:  semver.New(1, 4, 0, "", ""),
			want: nil,
		},
	}

	for _, tc := range testcases {
		got := slice(tforms, tc.min, tc.max)
		assert.Equal(t, tc.want, got)
	}
}

func TestAttributes(t *testing.T) {
	attr := []attribute.KeyValue{
		attribute.Bool("foo", true),
		attribute.Bool("bar", false),
	}

	a := newAttributes(attr)
	a.Rename("bar", "baz")
	assert.Equal(t, []attribute.KeyValue{
		attribute.Bool("foo", true),
		attribute.Bool("baz", false),
	}, attr)

	a.Rename("baz", "foo")
	assert.Equal(t, []attribute.KeyValue{
		attribute.Bool("foo", true),
		attribute.Bool("foo", false),
	}, attr)

	a.Rename("foo", "bar")
	assert.Equal(t, []attribute.KeyValue{
		attribute.Bool("bar", true),
		attribute.Bool("bar", false),
	}, attr)
}
