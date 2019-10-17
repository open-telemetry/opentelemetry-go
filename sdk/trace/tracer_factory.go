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
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/sdk/export"

	apitrace "go.opentelemetry.io/api/trace"
)

const (
	defaultTracerName = "go.opentelemetry.io/sdk/tracer"
)

type ProviderOptions struct {
	syncers  []export.SpanSyncer
	batchers []export.SpanBatcher
}

type ProviderOption func(*ProviderOptions)

type traceProvider struct {
	o              *ProviderOptions
	mu             sync.Mutex
	namedTracer    map[string]*tracer
	spanProcessors atomic.Value
}

var _ apitrace.Provider = &traceProvider{}

func NewProvider(opts ...ProviderOption) (apitrace.Provider, error) {

	o := &ProviderOptions{}
	for _, opt := range opts {
		opt(o)
	}
	tp := &traceProvider{o: o}

	new := make(spanProcessorMap)
	for _, syncer := range o.syncers {
		ssp := NewSimpleSpanProcessor(syncer)
		// TODO(rghetia): if unregister is not required then there is no need for sync.Once
		new[ssp] = &sync.Once{}
	}

	for _, batcher := range o.batchers {
		bsp, err := NewBatchSpanProcessor(batcher)
		if err != nil {
			return nil, err
		}
		new[bsp] = &sync.Once{}
	}
	// TODO (rghetia): if span processors are only register during construction then simple
	// map is sufficient.
	p.spanProcessors.Store(new)
	return tp, nil
}

func (p *traceProvider) GetTracer(name string) apitrace.Tracer {
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

func (p *traceProvider) RegisterSpanProcessor(s SpanProcessor) {
	p.mu.Lock()
	defer p.mu.Unlock()
	new := make(spanProcessorMap)
	if old, ok := spanProcessors.Load().(spanProcessorMap); ok {
		for k, v := range old {
			new[k] = v
		}
	}
	new[s] = &sync.Once{}
	p.spanProcessors.Store(new)
}
