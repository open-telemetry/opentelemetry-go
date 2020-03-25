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

package correlation

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
)

type testCase struct {
	name    string
	value   MapUpdate
	init    []int
	wantKVs []core.KeyValue
}

func TestMap(t *testing.T) {
	for _, testcase := range getTestCases() {
		t.Logf("Running test case %s", testcase.name)
		var got Map
		if len(testcase.init) > 0 {
			got = makeTestMap(testcase.init).Apply(testcase.value)
		} else {
			got = NewMap(testcase.value)
		}
		for _, s := range testcase.wantKVs {
			if ok := got.HasValue(s.Key); !ok {
				t.Errorf("Expected Key %s to have Value", s.Key)
			}
			if g, ok := got.Value(s.Key); !ok || g != s.Value {
				t.Errorf("+got: %v, -want: %v", g, s.Value)
			}
		}
		// test Foreach()
		got.Foreach(func(kv core.KeyValue) bool {
			for _, want := range testcase.wantKVs {
				if kv == want {
					return false
				}
			}
			t.Errorf("Expected kv %v, but not found", kv)
			return true
		})
		if l, exp := got.Len(), len(testcase.wantKVs); l != exp {
			t.Errorf("+got: %d, -want: %d", l, exp)
		}
	}
}

func TestSizeComputation(t *testing.T) {
	for _, testcase := range getTestCases() {
		t.Logf("Running test case %s", testcase.name)
		var initMap Map
		if len(testcase.init) > 0 {
			initMap = makeTestMap(testcase.init)
		} else {
			initMap = NewEmptyMap()
		}
		gotMap := initMap.Apply(testcase.value)

		delSet, addSet := getModificationSets(testcase.value)
		mapSize := getNewMapSize(initMap.m, delSet, addSet)

		if gotMap.Len() != mapSize {
			t.Errorf("Expected computed size to be %d, got %d", gotMap.Len(), mapSize)
		}
	}
}

func getTestCases() []testCase {
	return []testCase{
		{
			name: "New map with MultiKV",
			value: MapUpdate{MultiKV: []core.KeyValue{
				key.Int64("key1", 1),
				key.String("key2", "val2")},
			},
			init: []int{},
			wantKVs: []core.KeyValue{
				key.Int64("key1", 1),
				key.String("key2", "val2"),
			},
		},
		{
			name:  "New map with SingleKV",
			value: MapUpdate{SingleKV: key.String("key1", "val1")},
			init:  []int{},
			wantKVs: []core.KeyValue{
				key.String("key1", "val1"),
			},
		},
		{
			name: "New map with both add fields",
			value: MapUpdate{SingleKV: key.Int64("key1", 3),
				MultiKV: []core.KeyValue{
					key.String("key1", ""),
					key.String("key2", "val2")},
			},
			init: []int{},
			wantKVs: []core.KeyValue{
				key.String("key1", ""),
				key.String("key2", "val2"),
			},
		},
		{
			name:    "New map with empty MapUpdate",
			value:   MapUpdate{},
			init:    []int{},
			wantKVs: []core.KeyValue{},
		},
		{
			name:    "New map with DropSingleK",
			value:   MapUpdate{DropSingleK: core.Key("key1")},
			init:    []int{},
			wantKVs: []core.KeyValue{},
		},
		{
			name: "New map with DropMultiK",
			value: MapUpdate{DropMultiK: []core.Key{
				core.Key("key1"), core.Key("key2"),
			}},
			init:    []int{},
			wantKVs: []core.KeyValue{},
		},
		{
			name: "New map with both drop fields",
			value: MapUpdate{
				DropSingleK: core.Key("key1"),
				DropMultiK: []core.Key{
					core.Key("key1"),
					core.Key("key2"),
				},
			},
			init:    []int{},
			wantKVs: []core.KeyValue{},
		},
		{
			name: "New map with all fields",
			value: MapUpdate{
				DropSingleK: core.Key("key1"),
				DropMultiK: []core.Key{
					core.Key("key1"),
					core.Key("key2"),
				},
				SingleKV: key.String("key4", "val4"),
				MultiKV: []core.KeyValue{
					key.String("key1", ""),
					key.String("key2", "val2"),
					key.String("key3", "val3"),
				},
			},
			init: []int{},
			wantKVs: []core.KeyValue{
				key.String("key1", ""),
				key.String("key2", "val2"),
				key.String("key3", "val3"),
				key.String("key4", "val4"),
			},
		},
		{
			name: "Existing map with MultiKV",
			value: MapUpdate{MultiKV: []core.KeyValue{
				key.Int64("key1", 1),
				key.String("key2", "val2")},
			},
			init: []int{5},
			wantKVs: []core.KeyValue{
				key.Int64("key1", 1),
				key.String("key2", "val2"),
				key.Int("key5", 5),
			},
		},
		{
			name:  "Existing map with SingleKV",
			value: MapUpdate{SingleKV: key.String("key1", "val1")},
			init:  []int{5},
			wantKVs: []core.KeyValue{
				key.String("key1", "val1"),
				key.Int("key5", 5),
			},
		},
		{
			name: "Existing map with both add fields",
			value: MapUpdate{SingleKV: key.Int64("key1", 3),
				MultiKV: []core.KeyValue{
					key.String("key1", ""),
					key.String("key2", "val2")},
			},
			init: []int{5},
			wantKVs: []core.KeyValue{
				key.String("key1", ""),
				key.String("key2", "val2"),
				key.Int("key5", 5),
			},
		},
		{
			name:  "Existing map with empty MapUpdate",
			value: MapUpdate{},
			init:  []int{5},
			wantKVs: []core.KeyValue{
				key.Int("key5", 5),
			},
		},
		{
			name:  "Existing map with DropSingleK",
			value: MapUpdate{DropSingleK: core.Key("key1")},
			init:  []int{1, 5},
			wantKVs: []core.KeyValue{
				key.Int("key5", 5),
			},
		},
		{
			name: "Existing map with DropMultiK",
			value: MapUpdate{DropMultiK: []core.Key{
				core.Key("key1"), core.Key("key2"),
			}},
			init: []int{1, 5},
			wantKVs: []core.KeyValue{
				key.Int("key5", 5),
			},
		},
		{
			name: "Existing map with both drop fields",
			value: MapUpdate{
				DropSingleK: core.Key("key1"),
				DropMultiK: []core.Key{
					core.Key("key1"),
					core.Key("key2"),
				},
			},
			init: []int{1, 2, 5},
			wantKVs: []core.KeyValue{
				key.Int("key5", 5),
			},
		},
		{
			name: "Existing map with all the fields",
			value: MapUpdate{
				DropSingleK: core.Key("key1"),
				DropMultiK: []core.Key{
					core.Key("key1"),
					core.Key("key2"),
					core.Key("key5"),
					core.Key("key6"),
				},
				SingleKV: key.String("key4", "val4"),
				MultiKV: []core.KeyValue{
					key.String("key1", ""),
					key.String("key2", "val2"),
					key.String("key3", "val3"),
				},
			},
			init: []int{5, 6, 7},
			wantKVs: []core.KeyValue{
				key.String("key1", ""),
				key.String("key2", "val2"),
				key.String("key3", "val3"),
				key.String("key4", "val4"),
				key.Int("key7", 7),
			},
		},
	}
}

func makeTestMap(ints []int) Map {
	r := make(rawMap, len(ints))
	for _, v := range ints {
		r[core.Key(fmt.Sprintf("key%d", v))] = core.Int(v)
	}
	return newMap(r)
}
