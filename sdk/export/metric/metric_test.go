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

package metric

import (
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"

	"github.com/stretchr/testify/require"
)

var testSlice = []core.KeyValue{
	key.String("bar", "baz"),
	key.Int("foo", 42),
}

func TestLabelIterator(t *testing.T) {
	iter := LabelSlice(testSlice).Iter()
	require.Equal(t, 2, iter.Len())

	require.True(t, iter.Next())
	require.Equal(t, key.String("bar", "baz"), iter.Label())
	idx, label := iter.IndexedLabel()
	require.Equal(t, 0, idx)
	require.Equal(t, key.String("bar", "baz"), label)
	require.Equal(t, 2, iter.Len())

	require.True(t, iter.Next())
	require.Equal(t, key.Int("foo", 42), iter.Label())
	idx, label = iter.IndexedLabel()
	require.Equal(t, 1, idx)
	require.Equal(t, key.Int("foo", 42), label)
	require.Equal(t, 2, iter.Len())

	require.False(t, iter.Next())
	require.Equal(t, 2, iter.Len())
}

func TestEmptyLabelIterator(t *testing.T) {
	iter := LabelSlice(nil).Iter()
	require.Equal(t, 0, iter.Len())
	require.False(t, iter.Next())
}

func TestIteratorToSlice(t *testing.T) {
	iter := LabelSlice(testSlice).Iter()
	got := IteratorToSlice(iter)
	require.Equal(t, testSlice, got)

	iter = LabelSlice(nil).Iter()
	got = IteratorToSlice(iter)
	require.Nil(t, got)
}
