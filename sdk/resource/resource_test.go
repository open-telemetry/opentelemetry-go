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

package resource_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	kv11 = core.Key("k1").String("v11")
	kv12 = core.Key("k1").String("v12")
	kv21 = core.Key("k2").String("v21")
	kv31 = core.Key("k3").String("v31")
	kv41 = core.Key("k4").String("v41")
)

func TestNew(t *testing.T) {
	cases := []struct {
		name string
		in   []core.KeyValue
		want []core.KeyValue
	}{
		{
			name: "New with common key order1",
			in:   []core.KeyValue{kv11, kv12, kv21},
			want: []core.KeyValue{kv11, kv21},
		},
		{
			name: "New with common key order2",
			in:   []core.KeyValue{kv12, kv11, kv21},
			want: []core.KeyValue{kv12, kv21},
		},
		{
			name: "New with nil",
			in:   nil,
			want: []core.KeyValue{},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res := resource.New(c.in...)
			if diff := cmp.Diff(
				sortedAttributes(res.Attributes()),
				sortedAttributes(c.want),
				cmp.AllowUnexported(core.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	cases := []struct {
		name string
		a, b *resource.Resource
		want []core.KeyValue
	}{
		{
			name: "Merge with no overlap, no nil",
			a:    resource.New(kv11, kv31),
			b:    resource.New(kv21, kv41),
			want: []core.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with common key order1",
			a:    resource.New(kv11),
			b:    resource.New(kv12, kv21),
			want: []core.KeyValue{kv21, kv11},
		},
		{
			name: "Merge with common key order2",
			a:    resource.New(kv12, kv21),
			b:    resource.New(kv11),
			want: []core.KeyValue{kv12, kv21},
		},
		{
			name: "Merge with first resource nil",
			a:    nil,
			b:    resource.New(kv21),
			want: []core.KeyValue{kv21},
		},
		{
			name: "Merge with second resource nil",
			a:    resource.New(kv11),
			b:    nil,
			want: []core.KeyValue{kv11},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res := resource.Merge(c.a, c.b)
			if diff := cmp.Diff(
				sortedAttributes(res.Attributes()),
				sortedAttributes(c.want),
				cmp.AllowUnexported(core.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func sortedAttributes(attrs []core.KeyValue) []core.KeyValue {
	sort.Slice(attrs[:], func(i, j int) bool {
		return attrs[i].Key < attrs[j].Key
	})
	return attrs
}
