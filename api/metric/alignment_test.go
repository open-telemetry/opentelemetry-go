package metric

import (
	"os"
	"testing"
	"unsafe"

	ottest "go.opentelemetry.io/otel/internal/testing"
)

// Ensure struct alignment prior to running tests.
func TestMain(m *testing.M) {
	fields := []ottest.FieldOffset{
		{
			Name:   "Measurement.number",
			Offset: unsafe.Offsetof(Measurement{}.number),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
