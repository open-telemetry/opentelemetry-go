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
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/core"
)

func TestAttributeIterator(t *testing.T) {
	one := core.Key("one").String("1")
	two := core.Key("two").Int(2)
	iter := NewAttributeIterator([]core.KeyValue{one, two})
	require.Equal(t, 2, iter.Len())

	require.True(t, iter.Next())
	require.Equal(t, one, iter.Attribute())
	idx, attr := iter.IndexedAttribute()
	require.Equal(t, 0, idx)
	require.Equal(t, one, attr)
	require.Equal(t, 2, iter.Len())

	require.True(t, iter.Next())
	require.Equal(t, two, iter.Attribute())
	idx, attr = iter.IndexedAttribute()
	require.Equal(t, 1, idx)
	require.Equal(t, two, attr)
	require.Equal(t, 2, iter.Len())

	require.False(t, iter.Next())
	require.Equal(t, 2, iter.Len())
}

func TestEmptyAttributeIterator(t *testing.T) {
	iter := NewAttributeIterator(nil)
	require.Equal(t, 0, iter.Len())
	require.False(t, iter.Next())
}
