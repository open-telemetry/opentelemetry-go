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

// This package provide test routines for the LabelIterator
// implementations.
package test // import "go.opentelemetry.io/otel/sdk/export/metric/test"

import (
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
	export "go.opentelemetry.io/otel/sdk/export/metric"

	"github.com/stretchr/testify/require"
)

// IteratorProvider is responsible for providing LabelIterators
// for testing.
type IteratorProvider interface {
	// Iterators should return a list of iterators that will yield
	// the same labels in the same order as specified in passed
	// slice. It should always return the same list of
	// iterators. Empty iterator list is also valid - no tests
	// will be run in this case.
	Iterators([]core.KeyValue) []export.LabelIterator
	// EmptyIterators should return a list of empty iterators. It
	// should always return the same list of iterators. Empty
	// iterator list is also valid - no tests will be run in this
	// case.
	EmptyIterators() []export.LabelIterator
}

type iteratorTestFunc func(*testing.T, export.LabelIterator)

// RunLabelIteratorTests will test iterators given by the passed
// provider. It will call provider functions multiple times.
func RunLabelIteratorTests(t *testing.T, provider IteratorProvider) {
	slice := []core.KeyValue{
		key.String("bar", "baz"),
		key.Int("foo", 42),
	}
	iteratorFuncs := []iteratorTestFunc{
		testLabelIterator,
		testLabelIteratorReset,
	}
	for _, f := range iteratorFuncs {
		for _, iter := range provider.Iterators(slice) {
			f(t, iter)
		}
	}
	emptyIteratorFuncs := []iteratorTestFunc{
		testEmptyLabelIterator,
		testEmptyLabelIteratorReset,
	}
	for _, f := range emptyIteratorFuncs {
		for _, iter := range provider.EmptyIterators() {
			f(t, iter)
		}
	}
}

func testLabelIterator(t *testing.T, iter export.LabelIterator) {
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

func testEmptyLabelIterator(t *testing.T, iter export.LabelIterator) {
	require.Equal(t, 0, iter.Len())
	require.False(t, iter.Next())
}
