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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	defaultTracerName = "go.opentelemetry.io/otel/sdk/tracer"
)

// TODO (MrAlias): unify this API option design:
// https://github.com/open-telemetry/opentelemetry-go/issues/536

// TracerProviderConfig
type TracerProviderConfig struct {
	processors []SpanProcessor
	config     Config
}

type TracerProviderOption func(*TracerProviderConfig)

type TracerProvider struct {
	mu             sync.Mutex
	namedTracer    map[instrumentation.Library]*tracer
	spanProcessors atomic.Value
	config         atomic.Value // access atomically
}

var _ trace.TracerProvider = &TracerProvider{}

// NewTracerProvider returns a new and configured TracerProvider.
//
// By default the returned TracerProvider is configured with:
//  - a ParentBased(AlwaysSample) Sampler
//  - a random number IDGenerator
//  - the resource.Default() Resource
//  - the default SpanLimits.
//
// The passed opts are used to override these default values and configure the
// returned TracerProvider appropriately.
func NewTracerProvider(opts ...TracerProviderOption) *TracerProvider {
	o := &TracerProviderConfig{}

	for _, opt := range opts {
		opt(o)
	}

	tp := &TracerProvider{
		namedTracer: make(map[instrumentation.Library]*tracer),
	}
	tp.config.Store(&Config{
		DefaultSampler: ParentBased(AlwaysSample()),
		IDGenerator:    defaultIDGenerator(),
		SpanLimits: SpanLimits{
			AttributeCountLimit:         DefaultAttributeCountLimit,
			EventCountLimit:             DefaultEventCountLimit,
			LinkCountLimit:              DefaultLinkCountLimit,
			AttributePerEventCountLimit: DefaultAttributePerEventCountLimit,
			AttributePerLinkCountLimit:  DefaultAttributePerLinkCountLimit,
		},
	})

	for _, sp := range o.processors {
		tp.RegisterSpanProcessor(sp)
	}

	tp.ApplyConfig(o.config)

	return tp
}

// Tracer returns a Tracer with the given name and options. If a Tracer for
// the given name and options does not exist it is created, otherwise the
// existing Tracer is returned.
//
// If name is empty, DefaultTracerName is used instead.
//
// This method is safe to be called concurrently.
func (p *TracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	c := trace.NewTracerConfig(opts...)

	p.mu.Lock()
	defer p.mu.Unlock()
	if name == "" {
		name = defaultTracerName
	}
	il := instrumentation.Library{
		Name:    name,
		Version: c.InstrumentationVersion,
	}
	t, ok := p.namedTracer[il]
	if !ok {
		t = &tracer{
			provider:               p,
			instrumentationLibrary: il,
		}
		p.namedTracer[il] = t
	}
	return t
}

// RegisterSpanProcessor adds the given SpanProcessor to the list of SpanProcessors
func (p *TracerProvider) RegisterSpanProcessor(s SpanProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	new := spanProcessorStates{}
	if old, ok := p.spanProcessors.Load().(spanProcessorStates); ok {
		new = append(new, old...)
	}
	newSpanSync := &spanProcessorState{
		sp:    s,
		state: &sync.Once{},
	}
	new = append(new, newSpanSync)
	p.spanProcessors.Store(new)
}

// UnregisterSpanProcessor removes the given SpanProcessor from the list of SpanProcessors
func (p *TracerProvider) UnregisterSpanProcessor(s SpanProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	spss := spanProcessorStates{}
	old, ok := p.spanProcessors.Load().(spanProcessorStates)
	if !ok || len(old) == 0 {
		return
	}
	spss = append(spss, old...)

	// stop the span processor if it is started and remove it from the list
	var stopOnce *spanProcessorState
	var idx int
	for i, sps := range spss {
		if sps.sp == s {
			stopOnce = sps
			idx = i
		}
	}
	if stopOnce != nil {
		stopOnce.state.Do(func() {
			if err := s.Shutdown(context.Background()); err != nil {
				otel.Handle(err)
			}
		})
	}
	if len(spss) > 1 {
		copy(spss[idx:], spss[idx+1:])
	}
	spss[len(spss)-1] = nil
	spss = spss[:len(spss)-1]

	p.spanProcessors.Store(spss)
}

// ApplyConfig changes the configuration of the provider.
// If a field in the configuration is empty or nil then its original value is preserved.
func (p *TracerProvider) ApplyConfig(cfg Config) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c := *p.config.Load().(*Config)
	if cfg.DefaultSampler != nil {
		c.DefaultSampler = cfg.DefaultSampler
	}
	if cfg.IDGenerator != nil {
		c.IDGenerator = cfg.IDGenerator
	}
	if cfg.SpanLimits.EventCountLimit > 0 {
		c.SpanLimits.EventCountLimit = cfg.SpanLimits.EventCountLimit
	}
	if cfg.SpanLimits.AttributeCountLimit > 0 {
		c.SpanLimits.AttributeCountLimit = cfg.SpanLimits.AttributeCountLimit
	}
	if cfg.SpanLimits.LinkCountLimit > 0 {
		c.SpanLimits.LinkCountLimit = cfg.SpanLimits.LinkCountLimit
	}
	if cfg.SpanLimits.AttributePerEventCountLimit > 0 {
		c.SpanLimits.AttributePerEventCountLimit = cfg.SpanLimits.AttributePerEventCountLimit
	}
	if cfg.SpanLimits.AttributePerLinkCountLimit > 0 {
		c.SpanLimits.AttributePerLinkCountLimit = cfg.SpanLimits.AttributePerLinkCountLimit
	}
	c.Resource = cfg.Resource
	if c.Resource == nil {
		c.Resource = resource.Default()
	}
	p.config.Store(&c)
}

// ForceFlush immediately exports all spans that have not yet been exported for
// all the registered span processors.
func (p *TracerProvider) ForceFlush(ctx context.Context) error {
	spss, ok := p.spanProcessors.Load().(spanProcessorStates)
	if !ok {
		return fmt.Errorf("failed to load span processors")
	}
	if len(spss) == 0 {
		return nil
	}

	for _, sps := range spss {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := sps.sp.ForceFlush(ctx); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown shuts down the span processors in the order they were registered.
func (p *TracerProvider) Shutdown(ctx context.Context) error {
	spss, ok := p.spanProcessors.Load().(spanProcessorStates)
	if !ok {
		return fmt.Errorf("failed to load span processors")
	}
	if len(spss) == 0 {
		return nil
	}

	for _, sps := range spss {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var err error
		sps.state.Do(func() {
			err = sps.sp.Shutdown(ctx)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// WithSyncer registers the exporter with the TracerProvider using a
// SimpleSpanProcessor.
func WithSyncer(e export.SpanExporter) TracerProviderOption {
	return WithSpanProcessor(NewSimpleSpanProcessor(e))
}

// WithBatcher registers the exporter with the TracerProvider using a
// BatchSpanProcessor configured with the passed opts.
func WithBatcher(e export.SpanExporter, opts ...BatchSpanProcessorOption) TracerProviderOption {
	return WithSpanProcessor(NewBatchSpanProcessor(e, opts...))
}

// WithSpanProcessor registers the SpanProcessor with a TracerProvider.
func WithSpanProcessor(sp SpanProcessor) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.processors = append(opts.processors, sp)
	}
}

// WithResource returns a TracerProviderOption that will configure the
// Resource r as a TracerProvider's Resource. The configured Resource is
// referenced by all the Tracers the TracerProvider creates. It represents the
// entity producing telemetry.
//
// If this option is not used, the TracerProvider will use the
// resource.Default() Resource by default.
func WithResource(r *resource.Resource) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.config.Resource = r
	}
}

// WithIDGenerator returns a TracerProviderOption that will configure the
// IDGenerator g as a TracerProvider's IDGenerator. The configured IDGenerator
// is used by the Tracers the TracerProvider creates to generate new Span and
// Trace IDs.
//
// If this option is not used, the TracerProvider will use a random number
// IDGenerator by default.
func WithIDGenerator(g IDGenerator) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.config.IDGenerator = g
	}
}

// WithSampler returns a TracerProviderOption that will configure the Sampler
// s as a TracerProvider's Sampler. The configured Sampler is used by the
// Tracers the TracerProvider creates to make their sampling decisions for the
// Spans they create.
//
// If this option is not used, the TracerProvider will use a
// ParentBased(AlwaysSample) Sampler by default.
func WithSampler(s Sampler) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.config.DefaultSampler = s
	}
}

// WithSpanLimits returns a TracerProviderOption that will configure the
// SpanLimits sl as a TracerProvider's SpanLimits. The configured SpanLimits
// are used used by the Tracers the TracerProvider and the Spans they create
// to limit tracing resources used.
//
// If this option is not used, the TracerProvider will use the default
// SpanLimits.
func WithSpanLimits(sl SpanLimits) TracerProviderOption {
	return func(opts *TracerProviderConfig) {
		opts.config.SpanLimits = sl
	}
}
