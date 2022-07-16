package exporttest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestHistogramsComparison(t *testing.T) {
	a := export.Histogram{
		Temporality: export.CumulativeTemporality,
	}

	b := export.Histogram{
		Temporality: export.DeltaTemporality,
		DataPoints: []export.HistogramDataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("b", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
			},
		},
	}

	AssertHistogramsEqual(t, a, a)
	AssertHistogramsEqual(t, b, b)

	equal, explination := CompareHistogram(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 2, "Temporality and DataPoints do not match")
}
