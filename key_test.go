package otel_test

import (
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel"
)

func TestValue(t *testing.T) {
	k := otel.Key("test")
	bli := getBitlessInfo(42)
	for _, testcase := range []struct {
		name      string
		value     otel.Value
		wantType  otel.ValueType
		wantValue interface{}
	}{
		{
			name:      "Key.Bool() correctly returns keys's internal bool value",
			value:     k.Bool(true).Value,
			wantType:  otel.BOOL,
			wantValue: true,
		},
		{
			name:      "Key.Int64() correctly returns keys's internal int64 value",
			value:     k.Int64(42).Value,
			wantType:  otel.INT64,
			wantValue: int64(42),
		},
		{
			name:      "Key.Uint64() correctly returns keys's internal uint64 value",
			value:     k.Uint64(42).Value,
			wantType:  otel.UINT64,
			wantValue: uint64(42),
		},
		{
			name:      "Key.Float64() correctly returns keys's internal float64 value",
			value:     k.Float64(42.1).Value,
			wantType:  otel.FLOAT64,
			wantValue: float64(42.1),
		},
		{
			name:      "Key.Int32() correctly returns keys's internal int32 value",
			value:     k.Int32(42).Value,
			wantType:  otel.INT32,
			wantValue: int32(42),
		},
		{
			name:      "Key.Uint32() correctly returns keys's internal uint32 value",
			value:     k.Uint32(42).Value,
			wantType:  otel.UINT32,
			wantValue: uint32(42),
		},
		{
			name:      "Key.Float32() correctly returns keys's internal float32 value",
			value:     k.Float32(42.1).Value,
			wantType:  otel.FLOAT32,
			wantValue: float32(42.1),
		},
		{
			name:      "Key.String() correctly returns keys's internal string value",
			value:     k.String("foo").Value,
			wantType:  otel.STRING,
			wantValue: "foo",
		},
		{
			name:      "Key.Int() correctly returns keys's internal signed integral value",
			value:     k.Int(bli.intValue).Value,
			wantType:  bli.signedType,
			wantValue: bli.signedValue,
		},
		{
			name:      "Key.Uint() correctly returns keys's internal unsigned integral value",
			value:     k.Uint(bli.uintValue).Value,
			wantType:  bli.unsignedType,
			wantValue: bli.unsignedValue,
		},
	} {
		t.Logf("Running test case %s", testcase.name)
		if testcase.value.Type() != testcase.wantType {
			t.Errorf("wrong value type, got %#v, expected %#v", testcase.value.Type(), testcase.wantType)
		}
		got := testcase.value.AsInterface()
		if diff := cmp.Diff(testcase.wantValue, got); diff != "" {
			t.Errorf("+got, -want: %s", diff)
		}
	}
}

type bitlessInfo struct {
	intValue      int
	uintValue     uint
	signedType    otel.ValueType
	unsignedType  otel.ValueType
	signedValue   interface{}
	unsignedValue interface{}
}

func getBitlessInfo(i int) bitlessInfo {
	if unsafe.Sizeof(i) == 4 {
		return bitlessInfo{
			intValue:      i,
			uintValue:     uint(i),
			signedType:    otel.INT32,
			unsignedType:  otel.UINT32,
			signedValue:   int32(i),
			unsignedValue: uint32(i),
		}
	}
	return bitlessInfo{
		intValue:      i,
		uintValue:     uint(i),
		signedType:    otel.INT64,
		unsignedType:  otel.UINT64,
		signedValue:   int64(i),
		unsignedValue: uint64(i),
	}
}

func TestDefined(t *testing.T) {
	for _, testcase := range []struct {
		name string
		k    otel.Key
		want bool
	}{
		{
			name: "Key.Defined() returns true when len(v.Name) != 0",
			k:    otel.Key("foo"),
			want: true,
		},
		{
			name: "Key.Defined() returns false when len(v.Name) == 0",
			k:    otel.Key(""),
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//func (k otel.Key) Defined() bool {
			have := testcase.k.Defined()
			if have != testcase.want {
				t.Errorf("Want: %v, but have: %v", testcase.want, have)
			}
		})
	}
}

func TestEmit(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    otel.Value
		want string
	}{
		{
			name: `test Key.Emit() can emit a string representing self.BOOL`,
			v:    otel.Bool(true),
			want: "true",
		},
		{
			name: `test Key.Emit() can emit a string representing self.INT32`,
			v:    otel.Int32(42),
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.INT64`,
			v:    otel.Int64(42),
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.UINT32`,
			v:    otel.Uint32(42),
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.UINT64`,
			v:    otel.Uint64(42),
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.FLOAT32`,
			v:    otel.Float32(42.1),
			want: "42.1",
		},
		{
			name: `test Key.Emit() can emit a string representing self.FLOAT64`,
			v:    otel.Float64(42.1),
			want: "42.1",
		},
		{
			name: `test Key.Emit() can emit a string representing self.STRING`,
			v:    otel.String("foo"),
			want: "foo",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (v otel.Value) Emit() string {
			have := testcase.v.Emit()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}

func BenchmarkEmitBool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := otel.Bool(i%2 == 0)
		_ = n.Emit()
	}
}

func BenchmarkEmitInt64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := otel.Int64(int64(i))
		_ = n.Emit()
	}
}

func BenchmarkEmitUInt64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := otel.Uint64(uint64(i))
		_ = n.Emit()
	}
}

func BenchmarkEmitFloat64(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := otel.Float64(float64(i))
		_ = n.Emit()
	}
}

func BenchmarkEmitFloat32(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n := otel.Float32(float32(i))
		_ = n.Emit()
	}
}