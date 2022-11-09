// Copyright The OpenTelemetry Authors
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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"fmt"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/view"
)

// instProvider provides all OpenTelemetry instruments.
type instProvider[N int64 | float64] struct {
	scope   instrumentation.Scope
	pipes   pipelines
	resolve resolver[N]
}

func newInstProvider[N int64 | float64](s instrumentation.Scope, p pipelines, c instrumentCache[N]) *instProvider[N] {
	r := newResolver(p, c)
	return &instProvider[N]{
		scope:   s,
		pipes:   p,
		resolve: r,
	}
}

func (p *instProvider[N]) Lookup(kind view.InstrumentKind, name string, opts []metric.InstrumentOption) (*instrumentImpl[N], error) {
	cfg := metric.NewInstrumentConfig(opts...)
	return p.lookup(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        kind,
	}, cfg.Unit())
}

func (p *instProvider[N]) LookupObservable(kind view.InstrumentKind, name string, opts []metric.ObservableOption) (*instrumentImpl[N], error) {
	cfg := metric.NewObservableConfig(opts...)
	inst, err := p.lookup(view.Instrument{
		Scope:       p.scope,
		Name:        name,
		Description: cfg.Description(),
		Kind:        kind,
	}, cfg.Unit())
	if err != nil {
		return nil, err
	}

	if inst == nil {
		// Drop aggregator.
		return nil, nil
	}

	for _, cBack := range cfg.Callbacks() {
		p.pipes.registerCallback(cBack)
	}

	return inst, nil
}

func (p *instProvider[N]) lookup(inst view.Instrument, u unit.Unit) (*instrumentImpl[N], error) {
	aggs, err := p.resolve.Aggregators(inst, u)
	if len(aggs) == 0 && err != nil {
		err = fmt.Errorf("instrument does not match any view: %w", err)
	}
	return &instrumentImpl[N]{aggregators: aggs}, err
}
