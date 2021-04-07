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
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv"
)

var (
	kv11 = attribute.String("k1", "v11")
	kv12 = attribute.String("k1", "v12")
	kv21 = attribute.String("k2", "v21")
	kv31 = attribute.String("k3", "v31")
	kv41 = attribute.String("k4", "v41")
	kv42 = attribute.String("k4", "")
)

func TestNew(t *testing.T) {
	cases := []struct {
		name string
		in   []attribute.KeyValue
		want []attribute.KeyValue
	}{
		{
			name: "Key with common key order1",
			in:   []attribute.KeyValue{kv12, kv11, kv21},
			want: []attribute.KeyValue{kv11, kv21},
		},
		{
			name: "Key with common key order2",
			in:   []attribute.KeyValue{kv11, kv12, kv21},
			want: []attribute.KeyValue{kv12, kv21},
		},
		{
			name: "Key with nil",
			in:   nil,
			want: nil,
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res := resource.NewWithAttributes(c.in...)
			if diff := cmp.Diff(
				res.Attributes(),
				c.want,
				cmp.AllowUnexported(attribute.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	cases := []struct {
		name string
		a, b *resource.Resource
		want []attribute.KeyValue
	}{
		{
			name: "Merge with no overlap, no nil",
			a:    resource.NewWithAttributes(kv11, kv31),
			b:    resource.NewWithAttributes(kv21, kv41),
			want: []attribute.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with no overlap, no nil, not interleaved",
			a:    resource.NewWithAttributes(kv11, kv21),
			b:    resource.NewWithAttributes(kv31, kv41),
			want: []attribute.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with common key order1",
			a:    resource.NewWithAttributes(kv11),
			b:    resource.NewWithAttributes(kv12, kv21),
			want: []attribute.KeyValue{kv12, kv21},
		},
		{
			name: "Merge with common key order2",
			a:    resource.NewWithAttributes(kv12, kv21),
			b:    resource.NewWithAttributes(kv11),
			want: []attribute.KeyValue{kv11, kv21},
		},
		{
			name: "Merge with common key order4",
			a:    resource.NewWithAttributes(kv11, kv21, kv41),
			b:    resource.NewWithAttributes(kv31, kv41),
			want: []attribute.KeyValue{kv11, kv21, kv31, kv41},
		},
		{
			name: "Merge with no keys",
			a:    resource.NewWithAttributes(),
			b:    resource.NewWithAttributes(),
			want: nil,
		},
		{
			name: "Merge with first resource no keys",
			a:    resource.NewWithAttributes(),
			b:    resource.NewWithAttributes(kv21),
			want: []attribute.KeyValue{kv21},
		},
		{
			name: "Merge with second resource no keys",
			a:    resource.NewWithAttributes(kv11),
			b:    resource.NewWithAttributes(),
			want: []attribute.KeyValue{kv11},
		},
		{
			name: "Merge with first resource nil",
			a:    nil,
			b:    resource.NewWithAttributes(kv21),
			want: []attribute.KeyValue{kv21},
		},
		{
			name: "Merge with second resource nil",
			a:    resource.NewWithAttributes(kv11),
			b:    nil,
			want: []attribute.KeyValue{kv11},
		},
		{
			name: "Merge with first resource value empty string",
			a:    resource.NewWithAttributes(kv42),
			b:    resource.NewWithAttributes(kv41),
			want: []attribute.KeyValue{kv41},
		},
		{
			name: "Merge with second resource value empty string",
			a:    resource.NewWithAttributes(kv41),
			b:    resource.NewWithAttributes(kv42),
			want: []attribute.KeyValue{kv42},
		},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("case-%s", c.name), func(t *testing.T) {
			res := resource.Merge(c.a, c.b)
			if diff := cmp.Diff(
				res.Attributes(),
				c.want,
				cmp.AllowUnexported(attribute.Value{})); diff != "" {
				t.Fatalf("unwanted result: diff %+v,", diff)
			}
		})
	}
}

func TestDefault(t *testing.T) {
	res := resource.Default()
	require.False(t, res.Equal(resource.Empty()))
	require.True(t, res.Set().HasValue(semconv.ServiceNameKey))

	serviceName, _ := res.Set().Value(semconv.ServiceNameKey)
	require.True(t, strings.HasPrefix(serviceName.AsString(), "unknown_service:"))
	require.Greaterf(t, len(serviceName.AsString()), len("unknown_service:"),
		"default service.name should include executable name")

	require.Contains(t, res.Attributes(), semconv.TelemetrySDKLanguageGo)
	require.Contains(t, res.Attributes(), semconv.TelemetrySDKVersionKey.String(otel.Version()))
	require.Contains(t, res.Attributes(), semconv.TelemetrySDKNameKey.String("opentelemetry"))
}

func TestString(t *testing.T) {
	for _, test := range []struct {
		kvs  []attribute.KeyValue
		want string
	}{
		{
			kvs:  nil,
			want: "",
		},
		{
			kvs:  []attribute.KeyValue{},
			want: "",
		},
		{
			kvs:  []attribute.KeyValue{kv11},
			want: "k1=v11",
		},
		{
			kvs:  []attribute.KeyValue{kv11, kv12},
			want: "k1=v12",
		},
		{
			kvs:  []attribute.KeyValue{kv11, kv21},
			want: "k1=v11,k2=v21",
		},
		{
			kvs:  []attribute.KeyValue{kv21, kv11},
			want: "k1=v11,k2=v21",
		},
		{
			kvs:  []attribute.KeyValue{kv11, kv21, kv31},
			want: "k1=v11,k2=v21,k3=v31",
		},
		{
			kvs:  []attribute.KeyValue{kv31, kv11, kv21},
			want: "k1=v11,k2=v21,k3=v31",
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A", "a"), attribute.String("B", "b")},
			want: "A=a,B=b",
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A", "a,B=b")},
			want: `A=a\,B\=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A", `a,B\=b`)},
			want: `A=a\,B\\\=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("A=a,B", `b`)},
			want: `A\=a\,B=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String(`A=a\,B`, `b`)},
			want: `A\=a\\\,B=b`,
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("", "invalid")},
			want: "",
		},
		{
			kvs:  []attribute.KeyValue{attribute.String("", "invalid"), attribute.String("B", "b")},
			want: "B=b",
		},
	} {
		if got := resource.NewWithAttributes(test.kvs...).String(); got != test.want {
			t.Errorf("Resource(%v).String() = %q, want %q", test.kvs, got, test.want)
		}
	}
}

func TestMarshalJSON(t *testing.T) {
	r := resource.NewWithAttributes(attribute.Int64("A", 1), attribute.String("C", "D"))
	data, err := json.Marshal(r)
	require.NoError(t, err)
	require.Equal(t,
		`[{"Key":"A","Value":{"Type":"INT64","Value":1}},{"Key":"C","Value":{"Type":"STRING","Value":"D"}}]`,
		string(data))
}
