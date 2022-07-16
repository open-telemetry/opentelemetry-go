package exporttest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestMetricsComparison(t *testing.T) {
	a := export.Metrics{
		Name:        "a",
		Description: "a desc",
		Unit:        unit.Dimensionless,
	}

	b := export.Metrics{
		Name:        "b",
		Description: "b desc",
		Unit:        unit.Bytes,
		Data: export.Gauge{
			DataPoints: []export.DataPoint{
				{
					Attributes: attribute.NewSet(attribute.Bool("b", true)),
					StartTime:  time.Now(),
					Time:       time.Now(),
					Value:      export.Int64(1),
				},
			},
		},
	}

	AssertMetricsEqual(t, a, a)
	AssertMetricsEqual(t, b, b)

	equal, explination := CompareMetrics(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 4, "Name, Description, Unit, and Data do not match")
}
