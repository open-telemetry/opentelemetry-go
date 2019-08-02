package core

import (
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/api/registry"
)

func TestBool(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    bool
		want Value
	}{
		{
			name: "value: true",
			v:    true,
			want: Value{
				Type: BOOL,
				Bool: true,
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Bool(v bool) KeyValue {}
			have := Key{}.Bool(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    int64
		want Value
	}{
		{
			name: "value: int64(42)",
			v:    int64(42),
			want: Value{
				Type:  INT64,
				Int64: int64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Int64(v int64) KeyValue {
			have := Key{}.Int64(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUint64(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    uint64
		want Value
	}{
		{
			name: "value: uint64(42)",
			v:    uint64(42),
			want: Value{
				Type:   UINT64,
				Uint64: uint64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Uint64(v uint64) KeyValue {
			have := Key{}.Uint64(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    float64
		want Value
	}{
		{
			name: "value: float64(42.1)",
			v:    float64(42.1),
			want: Value{
				Type:    FLOAT64,
				Float64: float64(42.1),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Float64(v float64) KeyValue {
			have := Key{}.Float64(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt32(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    int32
		want Value
	}{
		{
			name: "value: int32(42)",
			v:    int32(42),
			want: Value{
				Type:  INT32,
				Int64: int64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Int32(v int32) KeyValue {
			have := Key{}.Int32(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUint32(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    uint32
		want Value
	}{
		{
			name: "value: uint32(42)",
			v:    uint32(42),
			want: Value{
				Type:   UINT32,
				Uint64: uint64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Uint32(v uint32) KeyValue {
			have := Key{}.Uint32(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestFloat32(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    float32
		want Value
	}{
		{
			name: "value: float32(42.0)",
			v:    float32(42.0),
			want: Value{
				Type:    FLOAT32,
				Float64: float64(42.0),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Float32(v float32) KeyValue {
			have := Key{}.Float32(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestString(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    string
		want Value
	}{
		{
			name: `value: string("foo")`,
			v:    "foo",
			want: Value{
				Type:   STRING,
				String: "foo",
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) String(v string) KeyValue {
			have := Key{}.String(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestBytes(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    []byte
		want Value
	}{
		{
			name: `value: []byte{'f','o','o'}`,
			v:    []byte{'f', 'o', 'o'},
			want: Value{
				Type:  BYTES,
				Bytes: []byte{'f', 'o', 'o'},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Bytes(v []byte) KeyValue {
			have := Key{}.Bytes(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt(t *testing.T) {
	WTYPE := INT64
	if unsafe.Sizeof(int(42)) == 4 {
		// switch the desired value-type depending on system int byte-size
		WTYPE = INT32
	}

	for _, testcase := range []struct {
		name string
		v    int
		want Value
	}{
		{
			name: `value: int(42)`,
			v:    int(42),
			want: Value{
				Type:  WTYPE,
				Int64: int64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Int(v int) KeyValue {
			have := Key{}.Int(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUint(t *testing.T) {
	WTYPE := UINT64
	if unsafe.Sizeof(uint(42)) == 4 {
		// switch the desired value-type depending on system int byte-size
		WTYPE = UINT32
	}

	for _, testcase := range []struct {
		name string
		v    uint
		want Value
	}{
		{
			name: `value: uint(42)`,
			v:    uint(42),
			want: Value{
				Type:   WTYPE,
				Uint64: 42,
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k Key) Uint(v uint) KeyValue {
			have := Key{}.Uint(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestDefined(t *testing.T) {
	for _, testcase := range []struct {
		name string
		k    Key
		want bool
	}{
		{
			name: `Key Defined`,
			k: Key{
				registry.Variable{
					Name: "foo",
				},
			},
			want: true,
		},
		{
			name: `Key not Defined`,
			k:    Key{registry.Variable{}},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//func (k Key) Defined() bool {
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
		v    Value
		want string
	}{
		{
			name: `bool`,
			v: Value{
				Type: BOOL,
				Bool: true,
			},
			want: "true",
		},
		{
			name: `int32`,
			v: Value{
				Type:  INT32,
				Int64: 42,
			},
			want: "42",
		},
		{
			name: `int64`,
			v: Value{
				Type:  INT64,
				Int64: 42,
			},
			want: "42",
		},
		{
			name: `uint32`,
			v: Value{
				Type:   UINT32,
				Uint64: 42,
			},
			want: "42",
		},
		{
			name: `uint64`,
			v: Value{
				Type:   UINT64,
				Uint64: 42,
			},
			want: "42",
		},
		{
			name: `float32`,
			v: Value{
				Type:    FLOAT32,
				Float64: 42.1,
			},
			want: "42.1",
		},
		{
			name: `float64`,
			v: Value{
				Type:    FLOAT64,
				Float64: 42.1,
			},
			want: "42.1",
		},
		{
			name: `string`,
			v: Value{
				Type:   STRING,
				String: "foo",
			},
			want: "foo",
		},
		{
			name: `bytes`,
			v: Value{
				Type:  BYTES,
				Bytes: []byte{'f', 'o', 'o'},
			},
			want: "foo",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (v Value) Emit() string {
			have := testcase.v.Emit()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}
