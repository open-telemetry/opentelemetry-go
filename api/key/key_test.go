package key_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/key"
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
				Key:   "k1",
				Value: core.Bool(true),
			},
		},
		{
			name:   "Int64",
			actual: key.Int64("k1", 123),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Int64(123),
			},
		},
		{
			name:   "Uint64",
			actual: key.Uint64("k1", 1),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Uint64(1),
			},
		},
		{
			name:   "Float64",
			actual: key.Float64("k1", 123.5),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Float64(123.5),
			},
		},
		{
			name:   "Int32",
			actual: key.Int32("k1", 123),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Int32(123),
			},
		},
		{
			name:   "Uint32",
			actual: key.Uint32("k1", 123),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Uint32(123),
			},
		},
		{
			name:   "Float32",
			actual: key.Float32("k1", 123.5),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Float32(123.5),
			},
		},
		{
			name:   "String",
			actual: key.String("k1", "123.5"),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.String("123.5"),
			},
		},
		{
			name:   "Int",
			actual: key.Int("k1", 123),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Int(123),
			},
		},
		{
			name:   "Uint",
			actual: key.Uint("k1", 123),
			expected: core.KeyValue{
				Key:   "k1",
				Value: core.Uint(123),
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.actual, test.expected, cmp.AllowUnexported(core.Value{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
