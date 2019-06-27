// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metric

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/stats"
)

type Float64Gauge struct {
	baseMetric
}

type Float64Entry struct {
	baseEntry
}

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
