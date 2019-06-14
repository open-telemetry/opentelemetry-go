package metric

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/stats"
)

type (
	Float64Gauge struct {
		baseMetric
	}

	Float64Entry struct {
		baseEntry
	}
)

func NewFloat64Gauge(name string, mos ...Option) *Float64Gauge {
	m := initBaseMetric(name, GaugeFloat64, mos, &Float64Gauge{}).(*Float64Gauge)
	return m
}

func (g *Float64Gauge) Gauge(values ...core.KeyValue) Float64Entry {
	var entry Float64Entry
	entry.init(g, values)
	return entry
}

func (g *Float64Gauge) DefinitionID() core.EventID {
	return g.eventID
}

func (g Float64Entry) Set(ctx context.Context, val float64) {
	stats.Record(ctx, g.base.measure.M(val).With(core.ScopeID{
		EventID: g.eventID,
	}))
}
