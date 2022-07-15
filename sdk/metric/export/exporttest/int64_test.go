package exporttest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestInt64sComparison(t *testing.T) {
	a := export.Int64(-1)
	b := export.Int64(2)

	AssertInt64sEqual(t, a, a)
	AssertInt64sEqual(t, b, b)

	equal, explination := CompareInt64(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 1, "Value does not match")
}
