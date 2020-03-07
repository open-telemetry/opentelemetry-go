package internal_test

import (
	"os"
	"testing"

	"go.opentelemetry.io/otel/api/global/internal"
	ottest "go.opentelemetry.io/otel/internal/testing"
)

// Ensure struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fieldsMap := internal.AtomicFieldOffsets()
	fields := make([]ottest.FieldOffset, 0, len(fieldsMap))
	for name, offset := range fieldsMap {
		fields = append(fields, ottest.FieldOffset{
			Name:   name,
			Offset: offset,
		})
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
