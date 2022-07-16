package exporttest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestDataPointsComparison(t *testing.T) {
	a := export.DataPoint{
		Attributes: attribute.NewSet(attribute.Bool("a", true)),
		StartTime:  time.Now(),
		Time:       time.Now(),
		Value:      export.Int64(2),
	}

	b := export.DataPoint{
		Attributes: attribute.NewSet(attribute.Bool("b", true)),
		StartTime:  time.Now(),
		Time:       time.Now(),
		Value:      export.Float64(1),
	}

	AssertDataPointsEqual(t, a, a)
	AssertDataPointsEqual(t, b, b)

	equal, explination := CompareDataPoint(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 4, "Attributes, StartTime, Time and Value do not match")
}
