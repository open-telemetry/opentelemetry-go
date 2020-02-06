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
			Name:   "record.refMapped.value",
			Offset: unsafe.Offsetof(record{}.refMapped.value),
		},
		{
			Name:   "record.modified",
			Offset: unsafe.Offsetof(record{}.modified),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
