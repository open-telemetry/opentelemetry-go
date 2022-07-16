package exporttest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestSumsComparison(t *testing.T) {
	a := export.Sum{
		Temporality: export.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints: []export.DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("a", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      export.Int64(2),
			},
		},
	}

	b := export.Sum{
		Temporality: export.DeltaTemporality,
		IsMonotonic: false,
		DataPoints: []export.DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("b", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      export.Int64(1),
			},
		},
	}

	AssertSumsEqual(t, a, a)
	AssertSumsEqual(t, b, b)

	equal, explination := CompareSum(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 3, "Temporality, IsMonotonic, and DataPoints do not match")
}
