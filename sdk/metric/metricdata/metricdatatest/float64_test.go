package exporttest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestFloat64sComparison(t *testing.T) {
	a := export.Float64(-1)
	b := export.Float64(2)

	AssertFloat64sEqual(t, a, a)
	AssertFloat64sEqual(t, b, b)

	equal, explination := CompareFloat64(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 1, "Value does not match")
}
