package key_test

import (
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/key"
)

func TestKeyValueConstructors(t *testing.T) {

	tt := []struct {
		name     string
		actual   core.KeyValue
		expected core.KeyValue
	}{
		{
			name:   "Bool",
			actual: key.Bool("k1", true),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type: core.BOOL,
					Bool: true,
				},
			},
		},
		{
			name:   "Int64",
			actual: key.Int64("k1", 123),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type:  core.INT64,
					Int64: 123,
				},
			},
		},
		{
			name:   "Uint64",
			actual: key.Uint64("k1", 1),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type:   core.UINT64,
					Uint64: 1,
				},
			},
		},
		{
			name:   "Float64",
			actual: key.Float64("k1", 123.5),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type:    core.FLOAT64,
					Float64: 123.5,
				},
			},
		},
		{
			name:   "Int32",
			actual: key.Int32("k1", 123),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type:  core.INT32,
					Int64: 123,
				},
			},
		},
		{
			name:   "Uint32",
			actual: key.Uint32("k1", 123),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type:   core.UINT32,
					Uint64: 123,
				},
			},
		},
		{
			name:   "Float32",
			actual: key.Float32("k1", 123.5),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type:    core.FLOAT32,
					Float64: 123.5,
				},
			},
		},
		{
			name:   "Bytes",
			actual: key.Bytes("k1", []byte("v1")),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Type:  core.BYTES,
					Bytes: []byte("v1"),
				},
			},
		},
		{
			name:   "Int",
			actual: key.Int("k1", 123),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Int64: 123,
					Type:  IntType(123),
				},
			},
		},
		{
			name:   "Uint",
			actual: key.Uint("k1", 123),
			expected: core.KeyValue{
				Key: "k1",
				Value: core.Value{
					Uint64: 123,
					Type:   UintType(123),
				},
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.actual, test.expected); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

// IntType returns the core.ValueType depending on system int byte-size
func IntType(v int) core.ValueType {
	if unsafe.Sizeof(v) == 4 {
		return core.INT32
	}
	return core.INT64
}

// UintType returns the core.ValueType depending on system uint byte-size
func UintType(v uint) core.ValueType {
	if unsafe.Sizeof(v) == 4 {
		return core.UINT32
	}
	return core.UINT64
}
