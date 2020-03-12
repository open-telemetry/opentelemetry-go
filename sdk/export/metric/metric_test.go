// Copyright 2020, OpenTelemetry Authors
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

func TestSliceLabelIterator(t *testing.T) {
	slice := []core.KeyValue{
		key.String("bar", "baz"),
		key.Int("foo", 42),
	}
	iter := NewSliceLabelIterator(slice)

	require.True(t, iter.Next())
	require.Equal(t, key.String("bar", "baz"), iter.Label())
	idx, label := iter.IndexedLabel()
	require.Equal(t, 0, idx)
	require.Equal(t, key.String("bar", "baz"), label)

	require.True(t, iter.Next())
	require.Equal(t, key.Int("foo", 42), iter.Label())
	idx, label = iter.IndexedLabel()
	require.Equal(t, 1, idx)
	require.Equal(t, key.Int("foo", 42), label)

	require.False(t, iter.Next())

	iter = NewSliceLabelIterator(nil)
	require.False(t, iter.Next())

	iter = NewSliceLabelIterator([]core.KeyValue{})
	require.False(t, iter.Next())
}

func TestIteratorToSlice(t *testing.T) {
	slice := []core.KeyValue{
		key.String("bar", "baz"),
		key.Int("foo", 42),
	}
	iter := NewSliceLabelIterator(slice)
	got := IteratorToSlice(iter)
	require.Equal(t, slice, got)

	iter = NewSliceLabelIterator(nil)
	got = IteratorToSlice(iter)
	require.Nil(t, got)
}
