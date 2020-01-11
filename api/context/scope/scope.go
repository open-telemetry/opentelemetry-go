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

package scope

import (
	"context"

	"go.opentelemetry.io/otel/api/context/label"
	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
)

type (
	// Scope is a container for static context related to
	// OpenTelemetry, including access to the effective Tracer()
	// and Meter() instances, resources, and namespace.
	//
	// Scopes are used to configure the OpenTelemetry
	// implementation locally in code.  Use Scopes to inject an
	// OpenTelemetry dependency into third-party libraries,
	// pre-bound to a set of resource labels.
	//
	// Scope provides a Meter and Tracer instance that, when used,
	// switch into the corresponding Scope.  The Tracer and Meter
	// provided by the Scope auomatically swiches into the Scope.
	// As a result, Spans created through a Scope.Tracer() take
	// place in the Scope's namespace and inherit the Scope's
	// resources.
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

	// Provider is an immutable description of the SDK that
	// provides Tracer and Meter interfaces.
	Provider struct {
		tracer trace.TracerSDK
		meter  metric.MeterSDK
	}
)

var (
	_ trace.Tracer = &scopeTracer{}
	_ metric.Meter = &scopeMeter{}

	nilProvider = &Provider{}
)

// NewProvider constructs an SDK provider.
func NewProvider(t trace.TracerSDK, m metric.MeterSDK) *Provider {
	return &Provider{
		tracer: t,
		meter:  m,
	}
}

// Tracer returns a Tracer with the namespace and resources of this Scope.
func (p *Provider) Tracer() trace.TracerSDK {
	return p.tracer
}

// Meter returns a Meter with the namespace and resources of this Scope.
func (p *Provider) Meter() metric.MeterSDK {
	return p.meter
}

// New returns a Scope for this Provider, with empty resources and
// namespace.
func (p *Provider) New() Scope {
	si := &scopeImpl{
		resources: label.Empty(),
		provider:  p,
	}
	si.scopeMeter.scopeImpl = si
	si.scopeTracer.scopeImpl = si
	return Scope{si}
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

// AddResources returns a Scope with the addition of new resource labels.
func (s Scope) AddResources(kvs ...core.KeyValue) Scope {
	if len(kvs) == 0 {
		return s
	}
	r := s.clone()
	r.resources = r.resources.AddMany(kvs...)
	return r
}

// WithNamespace returns a Scope with the namespace set to `name`.
func (s Scope) WithNamespace(name core.Namespace) Scope {
	r := s.clone()
	r.namespace = name
	return r
}

// WithMeter returns a Scope with the effective Meter SDK set.
func (s Scope) WithMeter(meter metric.MeterSDK) Scope {
	r := s.clone()
	r.provider = NewProvider(r.provider.tracer, meter)
	return r
}

// WithTracer returns a Tracer with the effective Tracer SDK set.
func (s Scope) WithTracer(tracer trace.TracerSDK) Scope {
	r := s.clone()
	r.provider = NewProvider(tracer, r.provider.meter)
	return r
}

// Provider returns the underlying interfaces that provide the SDK.
func (s Scope) Provider() *Provider {
	if s.scopeImpl == nil {
		return nilProvider
	}
	return s.provider
}

// Resources returns the label set of this Scope.
func (s Scope) Resources() label.Set {
	if s.scopeImpl == nil {
		return label.Empty()
	}
	return s.resources
}

// Namespace returns the namespace of this Scope.
func (s Scope) Namespace() core.Namespace {
	if s.scopeImpl == nil {
		return ""
	}
	return s.namespace
}

// Tracer returns the effective Tracer of this Scope.
func (s Scope) Tracer() trace.Tracer {
	if s.scopeImpl == nil {
		return trace.NoopTracer{}
	}
	return &s.scopeTracer
}

// Meter returns the effective Meter of this Scope.
func (s Scope) Meter() metric.Meter {
	if s.scopeImpl == nil {
		return metric.NoopMeter{}
	}
	return &s.scopeMeter
}

func (s *scopeImpl) tracer() trace.TracerSDK {
	if s == nil {
		return trace.NoopTracerSDK{}
	}
	return s.provider.Tracer()
}

func (s *scopeImpl) meter() metric.MeterSDK {
	if s == nil {
		return metric.NoopMeterSDK{}
	}
	return s.provider.Meter()
}

func (s *scopeImpl) enterScope(ctx context.Context) context.Context {
	o := Current(ctx)
	if o.scopeImpl == s {
		return ctx
	}
	return Scope{s}.InContext(ctx)
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
	return t.tracer().Start(t.enterScope(ctx), t.name(name), opts...)
}

func (t *scopeTracer) WithSpan(
	ctx context.Context,
	name string,
	fn func(ctx context.Context) error,
) error {
	return t.tracer().WithSpan(t.enterScope(ctx), t.name(name), fn)
}

func (m *scopeMeter) NewInt64Counter(name string, cos ...metric.CounterOptionApplier) metric.Int64Counter {
	return m.meter().NewInt64Counter(m.name(name), cos...)
}

func (m *scopeMeter) NewFloat64Counter(name string, cos ...metric.CounterOptionApplier) metric.Float64Counter {
	return m.meter().NewFloat64Counter(m.name(name), cos...)
}

func (m *scopeMeter) NewInt64Gauge(name string, gos ...metric.GaugeOptionApplier) metric.Int64Gauge {
	return m.meter().NewInt64Gauge(m.name(name), gos...)
}

func (m *scopeMeter) NewFloat64Gauge(name string, gos ...metric.GaugeOptionApplier) metric.Float64Gauge {
	return m.meter().NewFloat64Gauge(m.name(name), gos...)
}

func (m *scopeMeter) NewInt64Measure(name string, mos ...metric.MeasureOptionApplier) metric.Int64Measure {
	return m.meter().NewInt64Measure(m.name(name), mos...)
}

func (m *scopeMeter) NewFloat64Measure(name string, mos ...metric.MeasureOptionApplier) metric.Float64Measure {
	return m.meter().NewFloat64Measure(m.name(name), mos...)
}

func (m *scopeMeter) RecordBatch(ctx context.Context, labels []core.KeyValue, ms ...metric.Measurement) {
	m.meter().RecordBatch(m.enterScope(ctx), labels, ms...)
}
