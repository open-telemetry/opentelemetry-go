package exporttest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestResourceMetricsComparison(t *testing.T) {
	a := export.ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("resource", "a")),
	}

	b := export.ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("resource", "b")),
		ScopeMetrics: []export.ScopeMetrics{
			{Scope: instrumentation.Scope{Name: "b"}},
		},
	}

	AssertResourceMetricsEqual(t, a, a)
	AssertResourceMetricsEqual(t, b, b)

	equal, explination := CompareResourceMetrics(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 2, "Resource and ScopeMetrics do not match")
}
