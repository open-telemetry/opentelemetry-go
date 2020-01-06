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
			Name:   "record.refcount",
			Offset: unsafe.Offsetof(record{}.refcount),
		},
		{
			Name:   "record.collectedEpoch",
			Offset: unsafe.Offsetof(record{}.collectedEpoch),
		},
		{
			Name:   "record.modifiedEpoch",
			Offset: unsafe.Offsetof(record{}.modifiedEpoch),
		},
		{
			Name:   "record.reclaim",
			Offset: unsafe.Offsetof(record{}.reclaim),
		},
	}
	if !ottest.Aligned8Byte(fields, os.Stderr) {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
