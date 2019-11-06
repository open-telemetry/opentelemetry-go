// Copyright 2019, OpenTelemetry Authors
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

package distributedcontext

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestMap(t *testing.T) {
	for _, testcase := range []struct {
		name    string
		value   MapUpdate
		init    []int
		wantKVs []otel.KeyValue
	}{
		{
			name: "NewMap with MultiKV",
			value: MapUpdate{MultiKV: []otel.KeyValue{
				otel.Key("key1").Int64(1),
				otel.Key("key2").String("val2"),
			}},
			init: []int{},
			wantKVs: []otel.KeyValue{
				otel.Key("key1").Int64(1),
				otel.Key("key2").String("val2"),
			},
		},
		{
			name:  "NewMap with SingleKV",
			value: MapUpdate{SingleKV: otel.Key("key1").String("val1")},
			init:  []int{},
			wantKVs: []otel.KeyValue{
				otel.Key("key1").String("val1"),
			},
		},
		{
			name: "NewMap with MapUpdate",
			value: MapUpdate{SingleKV: otel.Key("key1").Int64(3),
				MultiKV: []otel.KeyValue{
					otel.Key("key1").String(""),
					otel.Key("key2").String("val2"),
				}},
			init: []int{},
			wantKVs: []otel.KeyValue{
				otel.Key("key1").String(""),
				otel.Key("key2").String("val2"),
			},
		},
		{
			name:    "NewMap with empty MapUpdate",
			value:   MapUpdate{MultiKV: []otel.KeyValue{}},
			init:    []int{},
			wantKVs: []otel.KeyValue{},
		},
		{
			name: "Map with MultiKV",
			value: MapUpdate{MultiKV: []otel.KeyValue{
				otel.Key("key1").Int64(1),
				otel.Key("key2").String("val2"),
			}},
			init: []int{5},
			wantKVs: []otel.KeyValue{
				otel.Key("key1").Int64(1),
				otel.Key("key2").String("val2"),
				otel.Key("key5").Int(5),
			},
		},
		{
			name:  "Map with SingleKV",
			value: MapUpdate{SingleKV: otel.Key("key1").String("val1")},
			init:  []int{5},
			wantKVs: []otel.KeyValue{
				otel.Key("key1").String("val1"),
				otel.Key("key5").Int(5),
			},
		},
		{
			name: "Map with MapUpdate",
			value: MapUpdate{SingleKV: otel.Key("key1").Int64(3),
				MultiKV: []otel.KeyValue{
					otel.Key("key1").String(""),
					otel.Key("key2").String("val2"),
				}},
			init: []int{5},
			wantKVs: []otel.KeyValue{
				otel.Key("key1").String(""),
				otel.Key("key2").String("val2"),
				otel.Key("key5").Int(5),
			},
		},
		{
			name:  "Map with empty MapUpdate",
			value: MapUpdate{MultiKV: []otel.KeyValue{}},
			init:  []int{5},
			wantKVs: []otel.KeyValue{
				otel.Key("key5").Int(5),
			},
		},
	} {
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
		got.Foreach(func(kv otel.KeyValue) bool {
			for _, want := range testcase.wantKVs {
				if kv == want {
					return false
				}
			}
			t.Errorf("Expected kv %v, but not found", kv)
			return true
		})
		if len, exp := got.Len(), len(testcase.wantKVs); len != exp {
			t.Errorf("+got: %d, -want: %d", len, exp)
		}
	}
}

func makeTestMap(ints []int) Map {
	r := make(rawMap, len(ints))
	for _, v := range ints {
		r[otel.Key(fmt.Sprintf("key%d", v))] = entry{
			value: otel.Int(v),
		}
	}
	return newMap(r)
}
