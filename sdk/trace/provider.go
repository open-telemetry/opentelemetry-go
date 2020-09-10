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

package trace

import (
	"sync"
	"sync/atomic"

	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel/api/trace"
	apitrace "go.opentelemetry.io/otel/api/trace"
)

const (
	defaultTracerName = "go.opentelemetry.io/otel/sdk/tracer"
)

// TODO (MrAlias): unify this API option design:
// https://github.com/open-telemetry/opentelemetry-go/issues/536

// ProviderOptions
type ProviderOptions struct {
	processors []SpanProcessor
	config     Config
}

type ProviderOption func(*ProviderOptions)

type Provider struct {
	mu             sync.Mutex
	namedTracer    map[instrumentation.Library]*tracer
	spanProcessors atomic.Value
	config         atomic.Value // access atomically
}

var _ apitrace.Provider = &Provider{}

// NewProvider creates an instance of trace provider. Optional
// parameter configures the provider with common options applicable
// to all tracer instances that will be created by this provider.
func NewProvider(opts ...ProviderOption) *Provider {
	o := &ProviderOptions{}

	for _, opt := range opts {
		opt(o)
	}

	tp := &Provider{
		namedTracer: make(map[instrumentation.Library]*tracer),
	}
	tp.config.Store(&Config{
		DefaultSampler:       ParentSample(AlwaysSample()),
		IDGenerator:          defIDGenerator(),
		MaxAttributesPerSpan: DefaultMaxAttributesPerSpan,
		MaxEventsPerSpan:     DefaultMaxEventsPerSpan,
		MaxLinksPerSpan:      DefaultMaxLinksPerSpan,
	})

	for _, sp := range o.processors {
		tp.RegisterSpanProcessor(sp)
	}

	tp.ApplyConfig(o.config)

	return tp
}

// Tracer with the given name. If a tracer for the given name does not exist,
// it is created first. If the name is empty, DefaultTracerName is used.
func (p *Provider) Tracer(name string, opts ...apitrace.TracerOption) apitrace.Tracer {
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
func (p *Provider) RegisterSpanProcessor(s SpanProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	new := make(spanProcessorMap)
	if old, ok := p.spanProcessors.Load().(spanProcessorMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	new[s] = &sync.Once{}
	p.spanProcessors.Store(new)
}

// UnregisterSpanProcessor removes the given SpanProcessor from the list of SpanProcessors
func (p *Provider) UnregisterSpanProcessor(s SpanProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	new := make(spanProcessorMap)
	if old, ok := p.spanProcessors.Load().(spanProcessorMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	if stopOnce, ok := new[s]; ok && stopOnce != nil {
		stopOnce.Do(func() {
			s.Shutdown()
		})
	}
	delete(new, s)
	p.spanProcessors.Store(new)
}

// ApplyConfig changes the configuration of the provider.
// If a field in the configuration is empty or nil then its original value is preserved.
func (p *Provider) ApplyConfig(cfg Config) {
	p.mu.Lock()
	defer p.mu.Unlock()
	c := *p.config.Load().(*Config)
	if cfg.DefaultSampler != nil {
		c.DefaultSampler = cfg.DefaultSampler
	}
	if cfg.IDGenerator != nil {
		c.IDGenerator = cfg.IDGenerator
	}
	if cfg.MaxEventsPerSpan > 0 {
		c.MaxEventsPerSpan = cfg.MaxEventsPerSpan
	}
	if cfg.MaxAttributesPerSpan > 0 {
		c.MaxAttributesPerSpan = cfg.MaxAttributesPerSpan
	}
	if cfg.MaxLinksPerSpan > 0 {
		c.MaxLinksPerSpan = cfg.MaxLinksPerSpan
	}
	if cfg.Resource != nil {
		c.Resource = cfg.Resource
	}
	p.config.Store(&c)
}

// WithSyncer registers the exporter with the Provider using a
// SimpleSpanProcessor.
func WithSyncer(e export.SpanExporter) ProviderOption {
	return WithSpanProcessor(NewSimpleSpanProcessor(e))
}

// WithBatcher registers the exporter with the Provider using a
// BatchSpanProcessor configured with the passed opts.
func WithBatcher(e export.SpanExporter, opts ...BatchSpanProcessorOption) ProviderOption {
	return WithSpanProcessor(NewBatchSpanProcessor(e, opts...))
}

// WithSpanProcessor registers the SpanProcessor with a Provider.
func WithSpanProcessor(sp SpanProcessor) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.processors = append(opts.processors, sp)
	}
}

// WithConfig option sets the configuration to provider.
func WithConfig(config Config) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.config = config
	}
}

// WithResource option attaches a resource to the provider.
// The resource is added to the span when it is started.
func WithResource(r *resource.Resource) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.config.Resource = r
	}
}
