// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attrnorm

import (
	"fmt"
	"slices"
	"testing"
	"unsafe"

	"go.opentelemetry.io/otel/attribute"
)

func TestSetMayHaveComposite(t *testing.T) {
	dup := attribute.Map("m", attribute.String("k", "first"), attribute.String("k", "second"))
	for _, size := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 14, 20, 21} {
		scalars := make([]attribute.KeyValue, 0, size)
		for i := range size {
			scalars = append(scalars, attribute.Int(fmt.Sprintf("k%02d", i), i))
		}
		set := attribute.NewSet(scalars...)
		if SetMayHaveComposite(&set) && size <= 20 {
			t.Errorf("size %d: scalar-only set reported as possibly composite", size)
		}

		if size == 0 {
			continue
		}
		for _, at := range []int{0, size - 1} {
			kvs := slices.Clone(scalars)
			kvs[at] = attribute.Slice("composite", dup.Value)
			withSlice := attribute.NewSet(kvs...)
			if !SetMayHaveComposite(&withSlice) {
				t.Errorf("size %d, index %d: SLICE value not detected", size, at)
			}

			kvs[at] = dup
			withMap := attribute.NewSet(kvs...)
			if !SetMayHaveComposite(&withMap) {
				t.Errorf("size %d, index %d: MAP value not detected", size, at)
			}
			if _, changed := SetDedup(withMap); !changed {
				t.Errorf("size %d, index %d: duplicate map key not deduplicated", size, at)
			}
		}
	}
}

func TestSetMirrorLayout(t *testing.T) {
	var set attribute.Set
	if got, want := unsafe.Sizeof(set), unsafe.Sizeof(rawSet{}); got != want {
		t.Fatalf("attribute.Set size = %d, rawSet size = %d", got, want)
	}
	set = attribute.NewSet(attribute.String("a", "b"))
	if _, ok := (*rawSet)(unsafe.Pointer(&set)).data.([1]attribute.KeyValue); !ok { //nolint:gosec // Layout check.
		t.Fatal("rawSet does not mirror attribute.Set storage")
	}
}
