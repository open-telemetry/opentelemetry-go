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

package trace

import (
	"context"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/sdk/export"

	apitrace "go.opentelemetry.io/api/trace"
)

const (
	defaultTracerName = "go.opentelemetry.io/sdk/tracer"
)

type currentSpanKeyType struct{}

var (
	currentSpanKey = &currentSpanKeyType{}
)
type ProviderOptions struct {
	syncers  []export.SpanSyncer
	batchers []export.SpanBatcher
	config   Config
}

type ProviderOption func(*ProviderOptions)

type TraceProvider struct {
	mu             sync.Mutex
	namedTracer    map[string]*tracer
	spanProcessors atomic.Value
	config         atomic.Value // access atomically
	currentSpanKey *currentSpanKeyType
}

var _ apitrace.Provider = &TraceProvider{}

func NewProvider(opts ...ProviderOption) (*TraceProvider, error) {

	o := &ProviderOptions{}

	for _, opt := range opts {
		opt(o)
	}

	tp := &TraceProvider{
		namedTracer: make(map[string]*tracer, 0),
	}
	tp.config.Store(&Config{
		DefaultSampler:       ProbabilitySampler(defaultSamplingProbability),
		IDGenerator:          defIDGenerator(),
		MaxAttributesPerSpan: DefaultMaxAttributesPerSpan,
		MaxEventsPerSpan:     DefaultMaxEventsPerSpan,
		MaxLinksPerSpan:      DefaultMaxLinksPerSpan,
	})

	for _, syncer := range o.syncers {
		ssp := NewSimpleSpanProcessor(syncer)
		tp.RegisterSpanProcessor(ssp)
	}

	for _, batcher := range o.batchers {
		bsp, err := NewBatchSpanProcessor(batcher)
		if err != nil {
			return nil, err
		}
		tp.RegisterSpanProcessor(bsp)
	}

	tp.ApplyConfig(o.config)

	return tp, nil
}

func (p *TraceProvider) GetTracer(name string) apitrace.Tracer {
	p.mu.Lock()
	defer p.mu.Unlock()
	if name == "" {
		name = defaultTracerName
	}
	t, ok := p.namedTracer[name]
	if !ok {
		t = &tracer{name: name, provider: p}
		p.namedTracer[name] = t
	}
	return t
}

func (p *TraceProvider) RegisterSpanProcessor(s SpanProcessor) {
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

// UnregisterSpanProcessor removes from the list of SpanProcessors the SpanProcessor that was
// registered with the given name.
func (p *TraceProvider) UnregisterSpanProcessor(s SpanProcessor) {
	mu.Lock()
	defer mu.Unlock()
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

func (p *TraceProvider) ApplyConfig(cfg Config) {
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
	p.config.Store(&c)
}


func (p *TraceProvider) setCurrentSpan(ctx context.Context, span apitrace.Span) context.Context {
	return context.WithValue(ctx, currentSpanKey, span)
}

func (p *TraceProvider) currentSpan(ctx context.Context) apitrace.Span {
	if span, has := ctx.Value(currentSpanKey).(apitrace.Span); has {
		return span
	}
	return apitrace.NoopSpan{}
}

func WithSyncer(syncer export.SpanSyncer) ProviderOption {
	return func (opts *ProviderOptions) {
		opts.syncers = append(opts.syncers, syncer)
	}
}

func WithBatcher(batcher export.SpanBatcher) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.batchers = append(opts.batchers, batcher)
	}
}

func WithConfig(config Config) ProviderOption {
	return func(opts *ProviderOptions) {
		opts.config = config
	}
}
