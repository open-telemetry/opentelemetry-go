// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package attribute // import "go.opentelemetry.io/otel/attribute"

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute/internal/fnv"
)

type generator[T any] struct {
	New   func(string, T) KeyValue
	Value T
}

func (g generator[T]) Build(key string) KeyValue {
	return g.New(key, g.Value)
}

type builder interface {
	Build(key string) KeyValue
}

// keyVals is all the KeyValue generators that are used for testing. This is
// not []KeyValue so test detect if internal fields of the KeyValue that differ
// for different instances are used in the hashing.
var keyVals = []builder{
	generator[bool]{Bool, true},
	generator[bool]{Bool, false},
	generator[[]bool]{BoolSlice, []bool{false, true}},
	generator[[]bool]{BoolSlice, []bool{true, true, false}},
	generator[int]{Int, -1278},
	generator[int]{Int, 0}, // Should be different than false above.
	generator[[]int]{IntSlice, []int{3, 23, 21, -8, 0}},
	generator[[]int]{IntSlice, []int{1}},
	generator[int64]{Int64, 1}, // Should be different from true and []int{1}.
	generator[int64]{Int64, 29369},
	generator[[]int64]{Int64Slice, []int64{3826, -38, -29, -1}},
	generator[[]int64]{Int64Slice, []int64{8, -328, 29, 0}},
	generator[float64]{Float64, -0.3812381},
	generator[float64]{Float64, 1e32},
	generator[[]float64]{Float64Slice, []float64{0.1, -3.8, -29., 0.3321}},
	generator[[]float64]{Float64Slice, []float64{-13e8, -32.8, 4., 1e28}},
	generator[string]{String, "foo"},
	generator[string]{String, "bar"},
	generator[[]string]{StringSlice, []string{"foo", "bar", "baz"}},
	generator[[]string]{StringSlice, []string{"[]i1"}},
}

func TestEquivalence(t *testing.T) {
	const key = "key"

	// Test all combinations up to length 3.
	n := len(keyVals)
	kvs0 := make([][]KeyValue, 0, 1+n+(n*n)+(n*n*n))
	kvs1 := make([][]KeyValue, 0, 1+n+(n*n)+(n*n*n))

	kvs0 = append(kvs0, []KeyValue{})
	kvs1 = append(kvs1, []KeyValue{})

	for i := 0; i < len(keyVals); i++ {
		kvs0 = append(kvs0, []KeyValue{keyVals[i].Build(key)})
		kvs1 = append(kvs1, []KeyValue{keyVals[i].Build(key)})

		for j := 0; j < len(keyVals); j++ {
			kvs0 = append(kvs0, []KeyValue{
				keyVals[i].Build(key),
				keyVals[j].Build(key),
			})
			kvs1 = append(kvs1, []KeyValue{
				keyVals[i].Build(key),
				keyVals[j].Build(key),
			})

			for k := 0; k < len(keyVals); k++ {
				kvs0 = append(kvs0, []KeyValue{
					keyVals[i].Build(key),
					keyVals[j].Build(key),
					keyVals[k].Build(key),
				})
				kvs1 = append(kvs1, []KeyValue{
					keyVals[i].Build(key),
					keyVals[j].Build(key),
					keyVals[k].Build(key),
				})
			}
		}
	}

	if testing.Short() {
		// If running with -short, evaluate a random subset.
		const reducedLen = 100

		cp0, cp1 := kvs0[:0], kvs1[:0]
		seen := make(map[int]struct{})
		for i := 0; i < reducedLen; i++ {
			n := rand.Intn(len(kvs0))
			if _, ok := seen[n]; ok {
				// Choose another.
				i--
				continue
			}
			seen[n] = struct{}{}
			cp0 = append(cp0, kvs0[i])
			cp1 = append(cp1, kvs1[i])
		}
		kvs0, kvs1 = cp0, cp1
	}

	for i, kv0 := range kvs0 {
		for j, kv1 := range kvs1 {
			h0, h1 := hashKVs(kv0), hashKVs(kv1)
			if i == j {
				assert.Equal(t, h0, h1, msg{"!=", i, j, h0, h1, kv0, kv1})
			} else {
				assert.NotEqual(t, h0, h1, msg{"==", i, j, h0, h1, kv0, kv1})
			}
		}
	}
}

type msg struct {
	cmp      string
	i, j     int
	h0, h1   fnv.Hash
	kv0, kv1 []KeyValue
}

func (m msg) String() string {
	return fmt.Sprintf(
		"(%d: %d)%s %s (%d: %d)%s",
		m.i, m.h0, m.slice(m.kv0), m.cmp, m.j, m.h1, m.slice(m.kv1),
	)
}

func (m msg) slice(kvs []KeyValue) string {
	if len(kvs) == 0 {
		return "[]"
	}

	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(string(kvs[0].Key))
	_, _ = b.WriteRune(':')
	_, _ = b.WriteString(kvs[0].Value.Emit())
	for _, kv := range kvs[1:] {
		_, _ = b.WriteRune(',')
		_, _ = b.WriteString(string(kv.Key))
		_, _ = b.WriteRune(':')
		_, _ = b.WriteString(kv.Value.Emit())
	}
	_, _ = b.WriteRune(']')
	return b.String()
}
