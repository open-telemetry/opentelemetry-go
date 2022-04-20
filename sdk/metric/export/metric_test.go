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

package export

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
)

var testSlice = []attribute.KeyValue{
	attribute.String("bar", "baz"),
	attribute.Int("foo", 42),
}

func newIter(slice []attribute.KeyValue) attribute.Iterator {
	attrs := attribute.NewSet(slice...)
	return attrs.Iter()
}

func TestAttributeIterator(t *testing.T) {
	iter := newIter(testSlice)
	require.Equal(t, 2, iter.Len())

	require.True(t, iter.Next())
	require.Equal(t, attribute.String("bar", "baz"), iter.Attribute())
	idx, kv := iter.IndexedAttribute()
	require.Equal(t, 0, idx)
	require.Equal(t, attribute.String("bar", "baz"), kv)
	require.Equal(t, 2, iter.Len())

	require.True(t, iter.Next())
	require.Equal(t, attribute.Int("foo", 42), iter.Attribute())
	idx, kv = iter.IndexedAttribute()
	require.Equal(t, 1, idx)
	require.Equal(t, attribute.Int("foo", 42), kv)
	require.Equal(t, 2, iter.Len())

	require.False(t, iter.Next())
	require.Equal(t, 2, iter.Len())
}

func TestEmptyAttributeIterator(t *testing.T) {
	iter := newIter(nil)
	require.Equal(t, 0, iter.Len())
	require.False(t, iter.Next())
}

func TestIteratorToSlice(t *testing.T) {
	iter := newIter(testSlice)
	got := iter.ToSlice()
	require.Equal(t, testSlice, got)

	iter = newIter(nil)
	got = iter.ToSlice()
	require.Nil(t, got)
}
