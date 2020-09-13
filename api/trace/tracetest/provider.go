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

package tracetest

import (
	"sync"

	"go.opentelemetry.io/otel/api/trace"
)

type Provider struct {
	config config

	tracersMu sync.Mutex
	tracers   map[instrumentation]*Tracer
}

var _ trace.Provider = (*Provider)(nil)

func NewProvider(opts ...Option) *Provider {
	return &Provider{
		config:  newConfig(opts...),
		tracers: make(map[instrumentation]*Tracer),
	}
}

type instrumentation struct {
	Name, Version string
}

func (p *Provider) Tracer(instName string, opts ...trace.TracerOption) trace.Tracer {
	conf := trace.NewTracerConfig(opts...)

	inst := instrumentation{
		Name:    instName,
		Version: conf.InstrumentationVersion,
	}
	p.tracersMu.Lock()
	defer p.tracersMu.Unlock()
	t, ok := p.tracers[inst]
	if !ok {
		t = &Tracer{
			Name:    instName,
			Version: conf.InstrumentationVersion,
			config:  &p.config,
		}
		p.tracers[inst] = t
	}
	return t
}
