package exporttest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestHistogramDataPointsComparison(t *testing.T) {
	a := export.HistogramDataPoint{
		Attributes:   attribute.NewSet(attribute.Bool("a", true)),
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Sum:          2,
	}

	max, min := 99.0, 3.
	b := export.HistogramDataPoint{
		Attributes:   attribute.NewSet(attribute.Bool("b", true)),
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          &max,
		Min:          &min,
		Sum:          3,
	}

	AssertHistogramDataPointsEqual(t, a, a)
	AssertHistogramDataPointsEqual(t, b, b)

	equal, explination := CompareHistogramDataPoint(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 9, "Attributes, StartTime, Time, Count, Bounds, BucketCounts, Max, Min, and Sum do not match")
}
