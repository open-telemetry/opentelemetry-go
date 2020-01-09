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

// TODO
// 0. Fix remaining tests
// 1. Add back sugar
// 2. Make lib/ver obligatory (let namespace be namespace)
// 3. Fix examples
// 4. New examples
//
// 3. Comments
// 4. OTEP draft?

package scope

import (
	"context"

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

type (
	Scope struct {
		*scopeImpl
	}

	scopeImpl struct {
		namespace   core.Namespace
		resources   label.Set
		provider    *Provider
		scopeTracer scopeTracer
		scopeMeter  scopeMeter
	}

	scopeTracer struct {
		*scopeImpl
	}

	scopeMeter struct {
		*scopeImpl
	}

	Provider struct {
		tracer trace.TracerWithNamespace
		meter  metric.MeterWithNamespace
	}
)

var (
	_ trace.Tracer = &scopeTracer{}
	_ metric.Meter = &scopeMeter{}

	nilProvider = &Provider{}
)

func NewProvider(t trace.TracerWithNamespace, m metric.MeterWithNamespace) *Provider {
	return &Provider{
		tracer: t,
		meter:  m,
	}
}

func (p *Provider) Tracer() trace.TracerWithNamespace {
	return p.tracer
}

func (p *Provider) Meter() metric.MeterWithNamespace {
	return p.meter
}

func (p *Provider) New() Scope {
	si := &scopeImpl{
		resources: label.Empty(),
		provider:  p,
	}
	si.scopeMeter.scopeImpl = si
	si.scopeTracer.scopeImpl = si
	return Scope{si}
}

func Empty() Scope {
	return Scope{}
}

func (s Scope) clone() Scope {
	var ri scopeImpl
	if s.scopeImpl != nil {
		ri.provider = s.provider
		ri.resources = s.resources
	} else {
		ri.provider = nilProvider
	}
	ri.scopeMeter.scopeImpl = &ri
	ri.scopeTracer.scopeImpl = &ri
	return Scope{
		scopeImpl: &ri,
	}
}

func (s Scope) WithResources(labels label.Set) Scope {
	r := s.clone()
	r.resources = labels
	return r
}

func (s Scope) AddResources(kvs ...core.KeyValue) Scope {
	if len(kvs) == 0 {
		return s
	}
	r := s.clone()
	r.resources = r.resources.AddMany(kvs...)
	return r
}

func (s Scope) WithNamespace(name core.Namespace) Scope {
	r := s.clone()
	r.namespace = name
	return r
}

func (s Scope) WithMeter(meter metric.MeterWithNamespace) Scope {
	r := s.clone()
	r.provider = NewProvider(r.provider.tracer, meter)
	return r
}

func (s Scope) WithTracer(tracer trace.TracerWithNamespace) Scope {
	r := s.clone()
	r.provider = NewProvider(tracer, r.provider.meter)
	return r
}

func (s Scope) Provider() *Provider {
	if s.scopeImpl == nil {
		return nilProvider
	}
	return s.provider
}

func (s Scope) Resources() label.Set {
	if s.scopeImpl == nil {
		return label.Empty()
	}
	return s.resources
}

func (s Scope) Namespace() core.Namespace {
	if s.scopeImpl == nil {
		return ""
	}
	return s.namespace
}

func (s Scope) Tracer() trace.Tracer {
	if s.scopeImpl == nil {
		return trace.NoopTracer{}
	}
	return &s.scopeTracer
}

func (s Scope) Meter() metric.Meter {
	if s.scopeImpl == nil {
		return metric.NoopMeter{}
	}
	return &s.scopeMeter
}

func (s *scopeImpl) enterScope(ctx context.Context) context.Context {
	o := Current(ctx)
	if o.scopeImpl == s {
		return ctx
	}
	return ContextWithScope(ctx, Scope{s})
}

func (s *scopeImpl) name(n string) core.Name {
	return core.Name{
		Base:      n,
		Namespace: s.namespace,
	}
}

func (t *scopeTracer) Start(
	ctx context.Context,
	name string,
	opts ...trace.StartOption,
) (context.Context, trace.Span) {
	if t.scopeImpl == nil {
		return ctx, trace.NoopSpan{}
	}
	return t.provider.Tracer().Start(t.enterScope(ctx), t.name(name), opts...)
}

func (t *scopeTracer) WithSpan(
	ctx context.Context,
	name string,
	fn func(ctx context.Context) error,
) error {
	if t.scopeImpl == nil {
		return fn(ctx)
	}
	return t.provider.Tracer().WithSpan(t.enterScope(ctx), t.name(name), fn)
}

func (m *scopeMeter) NewInt64Counter(name string, cos ...metric.CounterOptionApplier) metric.Int64Counter {
	return m.provider.Meter().NewInt64Counter(m.name(name), cos...)
}

func (m *scopeMeter) NewFloat64Counter(name string, cos ...metric.CounterOptionApplier) metric.Float64Counter {
	return m.provider.Meter().NewFloat64Counter(m.name(name), cos...)
}

func (m *scopeMeter) NewInt64Gauge(name string, gos ...metric.GaugeOptionApplier) metric.Int64Gauge {
	return m.provider.Meter().NewInt64Gauge(m.name(name), gos...)
}

func (m *scopeMeter) NewFloat64Gauge(name string, gos ...metric.GaugeOptionApplier) metric.Float64Gauge {
	return m.provider.Meter().NewFloat64Gauge(m.name(name), gos...)
}

func (m *scopeMeter) NewInt64Measure(name string, mos ...metric.MeasureOptionApplier) metric.Int64Measure {
	return m.provider.Meter().NewInt64Measure(m.name(name), mos...)
}

func (m *scopeMeter) NewFloat64Measure(name string, mos ...metric.MeasureOptionApplier) metric.Float64Measure {
	return m.provider.Meter().NewFloat64Measure(m.name(name), mos...)
}

func (m *scopeMeter) RecordBatch(ctx context.Context, labels []core.KeyValue, ms ...metric.Measurement) {
	m.provider.Meter().RecordBatch(m.enterScope(ctx), labels, ms...)
}
