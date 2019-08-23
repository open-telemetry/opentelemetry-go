package core_test

import (
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/registry"
)

func TestBool(t *testing.T) {
	for _, testcase := range []struct {
		name string
		v    bool
		want core.Value
	}{
		{
			name: "Key.Bool() correctly returns keys's internal bool value",
			v:    true,
			want: core.Value{
				Type: core.BOOL,
				Bool: true,
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) core.Bool(v bool) KeyValue {}
			have := core.Key{}.Bool(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.Int64() correctly returns keys's internal int64 value",
			v:    int64(42),
			want: core.Value{
				Type:  core.INT64,
				Int64: int64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Int64(v int64) KeyValue {
			have := core.Key{}.Int64(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.Uint64() correctly returns keys's internal uint64 value",
			v:    uint64(42),
			want: core.Value{
				Type:   core.UINT64,
				Uint64: uint64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Uint64(v uint64) KeyValue {
			have := core.Key{}.Uint64(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.float64() correctly returns keys's internal floa64 value",
			v:    float64(42.1),
			want: core.Value{
				Type:    core.FLOAT64,
				Float64: float64(42.1),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Float64(v float64) KeyValue {
			have := core.Key{}.Float64(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.int32() correctly returns keys's internal int32 value",
			v:    int32(42),
			want: core.Value{
				Type:  core.INT32,
				Int64: int64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Int32(v int32) KeyValue {
			have := core.Key{}.Int32(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.uint32() correctly returns keys's internal uint32 value",
			v:    uint32(42),
			want: core.Value{
				Type:   core.UINT32,
				Uint64: uint64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Uint32(v uint32) KeyValue {
			have := core.Key{}.Uint32(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.float32() correctly returns keys's internal float32 value",
			v:    float32(42.0),
			want: core.Value{
				Type:    core.FLOAT32,
				Float64: float64(42.0),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Float32(v float32) KeyValue {
			have := core.Key{}.Float32(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.String() correctly returns keys's internal string value",
			v:    "foo",
			want: core.Value{
				Type:   core.STRING,
				String: "foo",
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) String(v string) KeyValue {
			have := core.Key{}.String(testcase.v)
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
		want core.Value
	}{
		{
			name: "Key.Bytes() correctly returns keys's internal []byte value",
			v:    []byte{'f', 'o', 'o'},
			want: core.Value{
				Type:  core.BYTES,
				Bytes: []byte{'f', 'o', 'o'},
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Bytes(v []byte) KeyValue {
			have := core.Key{}.Bytes(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestInt(t *testing.T) {
	WTYPE := core.INT64
	if unsafe.Sizeof(int(42)) == 4 {
		// switch the desired value-type depending on system int byte-size
		WTYPE = core.INT32
	}

	for _, testcase := range []struct {
		name string
		v    int
		want core.Value
	}{
		{
			name: "Key.Int() correctly returns keys's internal int64 value",
			v:    int(42),
			want: core.Value{
				Type:  WTYPE,
				Int64: int64(42),
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Int(v int) KeyValue {
			have := core.Key{}.Int(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestUint(t *testing.T) {
	WTYPE := core.UINT64
	if unsafe.Sizeof(uint(42)) == 4 {
		// switch the desired value-type depending on system int byte-size
		WTYPE = core.UINT32
	}

	for _, testcase := range []struct {
		name string
		v    uint
		want core.Value
	}{
		{
			name: "Key.Uint() correctly returns keys's internal uint64 value",
			v:    uint(42),
			want: core.Value{
				Type:   WTYPE,
				Uint64: 42,
			},
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (k core.Key) Uint(v uint) KeyValue {
			have := core.Key{}.Uint(testcase.v)
			if diff := cmp.Diff(testcase.want, have.Value); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestDefined(t *testing.T) {
	for _, testcase := range []struct {
		name string
		k    core.Key
		want bool
	}{
		{
			name: "Key.Defined() returns true when len(v.Name) != 0",
			k: core.Key{
				registry.Variable{
					Name: "foo",
				},
			},
			want: true,
		},
		{
			name: "Key.Defined() returns false when len(v.Name) == 0",
			k:    core.Key{registry.Variable{}},
			want: false,
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//func (k core.Key) Defined() bool {
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
		v    core.Value
		want string
	}{
		{
			name: `test Key.Emit() can emit a string representing self.BOOL`,
			v: core.Value{
				Type: core.BOOL,
				Bool: true,
			},
			want: "true",
		},
		{
			name: `test Key.Emit() can emit a string representing self.INT32`,
			v: core.Value{
				Type:  core.INT32,
				Int64: 42,
			},
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.INT64`,
			v: core.Value{
				Type:  core.INT64,
				Int64: 42,
			},
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.UINT32`,
			v: core.Value{
				Type:   core.UINT32,
				Uint64: 42,
			},
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.UINT64`,
			v: core.Value{
				Type:   core.UINT64,
				Uint64: 42,
			},
			want: "42",
		},
		{
			name: `test Key.Emit() can emit a string representing self.FLOAT32`,
			v: core.Value{
				Type:    core.FLOAT32,
				Float64: 42.1,
			},
			want: "42.1",
		},
		{
			name: `test Key.Emit() can emit a string representing self.FLOAT64`,
			v: core.Value{
				Type:    core.FLOAT64,
				Float64: 42.1,
			},
			want: "42.1",
		},
		{
			name: `test Key.Emit() can emit a string representing self.STRING`,
			v: core.Value{
				Type:   core.STRING,
				String: "foo",
			},
			want: "foo",
		},
		{
			name: `test Key.Emit() can emit a string representing self.BYTES`,
			v: core.Value{
				Type:  core.BYTES,
				Bytes: []byte{'f', 'o', 'o'},
			},
			want: "foo",
		},
	} {
		t.Run(testcase.name, func(t *testing.T) {
			//proto: func (v core.Value) Emit() string {
			have := testcase.v.Emit()
			if have != testcase.want {
				t.Errorf("Want: %s, but have: %s", testcase.want, have)
			}
		})
	}
}
