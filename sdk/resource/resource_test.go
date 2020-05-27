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

package resource_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/api/kv/value"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	kv11 = kv.String("k1", "v11")
	kv12 = kv.String("k1", "v12")
	kv21 = kv.String("k2", "v21")
	kv31 = kv.String("k3", "v31")
	kv41 = kv.String("k4", "v41")
)

func TestNew(t *testing.T) {
	cases := []struct {
		name string
		in   []kv.KeyValue
		want []kv.KeyValue
	}{
		{
			name: "Key with common key order1",
			in:   []kv.KeyValue{kv12, kv11, kv21},
			want: []kv.KeyValue{kv11, kv21},
		},
		{
			name: "Key with common key order2",
			in:   []kv.KeyValue{kv11, kv12, kv21},
			want: []kv.KeyValue{kv12, kv21},
		},
		{
			name: "Key with nil",
			in:   nil,
			want: nil,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res := resource.New(c.in...)
			if diff := cmp.Diff(
				res.Attributes(),
				c.want,
				cmp.AllowUnexported(value.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	cases := []struct {
		name string
		a, b *resource.Resource
		want []kv.KeyValue
	}{
		{
			name: "Merge with no overlap, no nil",
			a:    resource.New(kv11, kv31),
			b:    resource.New(kv21, kv41),
			want: []kv.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with no overlap, no nil, not interleaved",
			a:    resource.New(kv11, kv21),
			b:    resource.New(kv31, kv41),
			want: []kv.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with common key order1",
			a:    resource.New(kv11),
			b:    resource.New(kv12, kv21),
			want: []kv.KeyValue{kv11, kv21},
		},
		{
			name: "Merge with common key order2",
			a:    resource.New(kv12, kv21),
			b:    resource.New(kv11),
			want: []kv.KeyValue{kv12, kv21},
		},
		{
			name: "Merge with common key order4",
			a:    resource.New(kv11, kv21, kv41),
			b:    resource.New(kv31, kv41),
			want: []kv.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with no keys",
			a:    resource.New(),
			b:    resource.New(),
			want: nil,
		},
		{
			name: "Merge with first resource no keys",
			a:    resource.New(),
			b:    resource.New(kv21),
			want: []kv.KeyValue{kv21},
		},
		{
			name: "Merge with second resource no keys",
			a:    resource.New(kv11),
			b:    resource.New(),
			want: []kv.KeyValue{kv11},
		},
		{
			name: "Merge with first resource nil",
			a:    nil,
			b:    resource.New(kv21),
			want: []kv.KeyValue{kv21},
		},
		{
			name: "Merge with second resource nil",
			a:    resource.New(kv11),
			b:    nil,
			want: []kv.KeyValue{kv11},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res := resource.Merge(c.a, c.b)
			if diff := cmp.Diff(
				res.Attributes(),
				c.want,
				cmp.AllowUnexported(value.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func TestString(t *testing.T) {
	for _, test := range []struct {
		kvs  []kv.KeyValue
		want string
	}{
		{
			kvs:  nil,
			want: "",
		},
		{
			kvs:  []kv.KeyValue{},
			want: "",
		},
		{
			kvs:  []kv.KeyValue{kv11},
			want: "k1=v11",
		},
		{
			kvs:  []kv.KeyValue{kv11, kv12},
			want: "k1=v12",
		},
		{
			kvs:  []kv.KeyValue{kv11, kv21},
			want: "k1=v11,k2=v21",
		},
		{
			kvs:  []kv.KeyValue{kv21, kv11},
			want: "k1=v11,k2=v21",
		},
		{
			kvs:  []kv.KeyValue{kv11, kv21, kv31},
			want: "k1=v11,k2=v21,k3=v31",
		},
		{
			kvs:  []kv.KeyValue{kv31, kv11, kv21},
			want: "k1=v11,k2=v21,k3=v31",
		},
		{
			kvs:  []kv.KeyValue{kv.String("A", "a"), kv.String("B", "b")},
			want: "A=a,B=b",
		},
		{
			kvs:  []kv.KeyValue{kv.String("A", "a,B=b")},
			want: `A=a\,B\=b`,
		},
		{
			kvs:  []kv.KeyValue{kv.String("A", `a,B\=b`)},
			want: `A=a\,B\\\=b`,
		},
		{
			kvs:  []kv.KeyValue{kv.String("A=a,B", `b`)},
			want: `A\=a\,B=b`,
		},
		{
			kvs:  []kv.KeyValue{kv.String(`A=a\,B`, `b`)},
			want: `A\=a\\\,B=b`,
		},
	} {
		if got := resource.New(test.kvs...).String(); got != test.want {
			t.Errorf("Resource(%v).String() = %q, want %q", test.kvs, got, test.want)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	r := resource.New(kv.Int64("A", 1), kv.String("C", "D"))
	data, err := json.Marshal(r)
	require.NoError(t, err)
	require.Equal(t,
		`[{"Key":"A","Value":{"Type":"INT64","Value":1}},{"Key":"C","Value":{"Type":"STRING","Value":"D"}}]`,
		string(data))
}
