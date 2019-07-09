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

package stats

import (
	"context"

	"github.com/open-telemetry/opentelemetry-go/api/core"
	"github.com/open-telemetry/opentelemetry-go/api/registry"
)

type Measure struct {
	Variable registry.Variable
}

type Measurement struct {
	Measure Measure
	Value   float64
	Scope   core.ScopeID
}

type Interface interface {
	Record(ctx context.Context, m ...Measurement)
	RecordSingle(ctx context.Context, m Measurement)
}

type Recorder struct {
	Scope core.ScopeID
}

var _ Interface = (*Recorder)(nil)

// TODO
// func With(scope scope.Scope) Recorder {
// 	return Recorder{scope.ScopeID()}
// }

func Record(ctx context.Context, m ...Measurement) {
	Recorder{}.Record(ctx, m...)
}

func RecordSingle(ctx context.Context, m Measurement) {
	Recorder{}.RecordSingle(ctx, m)
}

func (r Recorder) Record(ctx context.Context, m ...Measurement) {
	// observer.Record(observer.Event{
	// 	Type:    observer.RECORD_STATS,
	// 	Scope:   r.ScopeID,
	// 	Context: ctx,
	// 	Stats:   m,
	// })
}

func (r Recorder) RecordSingle(ctx context.Context, m Measurement) {
	// observer.Record(observer.Event{
	// 	Type:    observer.RECORD_STATS,
	// 	Scope:   r.ScopeID,
	// 	Context: ctx,
	// 	Stat:    m,
	// })
}

type AnyStatistic struct{}

func (AnyStatistic) String() string {
	return "AnyStatistic"
}

func NewMeasure(name string, opts ...registry.Option) Measure {
	return Measure{
		Variable: registry.Register(name, AnyStatistic{}, opts...),
	}
}

func (m Measure) M(value float64) Measurement {
	return Measurement{
		Measure: m,
		Value:   value,
	}
}
