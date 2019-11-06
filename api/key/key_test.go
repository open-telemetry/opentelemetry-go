package key_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/key"
)

func TestKeyValueConstructors(t *testing.T) {
	tt := []struct {
		name     string
		actual   otel.KeyValue
		expected otel.KeyValue
	}{
		{
			name:   "Bool",
			actual: key.Bool("k1", true),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Bool(true),
			},
		},
		{
			name:   "Int64",
			actual: key.Int64("k1", 123),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Int64(123),
			},
		},
		{
			name:   "Uint64",
			actual: key.Uint64("k1", 1),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Uint64(1),
			},
		},
		{
			name:   "Float64",
			actual: key.Float64("k1", 123.5),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Float64(123.5),
			},
		},
		{
			name:   "Int32",
			actual: key.Int32("k1", 123),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Int32(123),
			},
		},
		{
			name:   "Uint32",
			actual: key.Uint32("k1", 123),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Uint32(123),
			},
		},
		{
			name:   "Float32",
			actual: key.Float32("k1", 123.5),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Float32(123.5),
			},
		},
		{
			name:   "String",
			actual: key.String("k1", "123.5"),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.String("123.5"),
			},
		},
		{
			name:   "Int",
			actual: key.Int("k1", 123),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Int(123),
			},
		},
		{
			name:   "Uint",
			actual: key.Uint("k1", 123),
			expected: otel.KeyValue{
				Key:   "k1",
				Value: otel.Uint(123),
			},
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			if diff := cmp.Diff(test.actual, test.expected, cmp.AllowUnexported(otel.Value{})); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
