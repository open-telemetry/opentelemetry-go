package exporttest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export"
)

func TestScopeMetricsComparison(t *testing.T) {
	a := export.ScopeMetrics{
		Scope: instrumentation.Scope{Name: "a"},
	}

	b := export.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "b"},
		Metrics: []export.Metrics{{Name: "b"}},
	}

	AssertScopeMetricsEqual(t, a, a)
	AssertScopeMetricsEqual(t, b, b)

	equal, explination := CompareScopeMetrics(a, b)
	assert.Falsef(t, equal, "%v != %v", a, b)
	assert.Len(t, explination, 2, "Scope and Metrics do not match")
}
