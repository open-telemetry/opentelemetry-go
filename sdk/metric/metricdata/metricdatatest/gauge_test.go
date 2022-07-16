package exporttest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestGaugesComparison(t *testing.T) {
	a := export.Gauge{
		DataPoints: []export.DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("a", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      export.Int64(2),
			},
		},
	}

	b := export.Gauge{
		DataPoints: []export.DataPoint{
			{
				Attributes: attribute.NewSet(attribute.Bool("b", true)),
				StartTime:  time.Now(),
				Time:       time.Now(),
				Value:      export.Int64(1),
			},
		},
	}

	AssertGaugesEqual(t, a, a)
	AssertGaugesEqual(t, b, b)

	equal, explination := CompareGauge(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 1, "DataPoints do not match")
}
