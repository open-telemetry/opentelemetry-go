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

package oteltest // import "go.opentelemetry.io/otel/oteltest"

import (
	"sync"

	"go.opentelemetry.io/otel/trace"
)

// TracerProvider is a testing TracerProvider. It is an functioning
// implementation of an OpenTelemetry TracerProvider and can be configured
// with a SpanRecorder that it configure all Tracers it creates to record
// their Spans with.
type TracerProvider struct {
	config config

	tracersMu sync.Mutex
	tracers   map[instrumentation]*Tracer
}

var _ trace.TracerProvider = (*TracerProvider)(nil)

// NewTracerProvider returns a *TracerProvider configured with options.
func NewTracerProvider(options ...Option) *TracerProvider {
	return &TracerProvider{
		config:  newConfig(options...),
		tracers: make(map[instrumentation]*Tracer),
	}
}

type instrumentation struct {
	Name, Version string
}

// Tracer returns an OpenTelemetry Tracer used for testing.
func (p *TracerProvider) Tracer(instName string, opts ...trace.TracerOption) trace.Tracer {
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

// DefaulTracer returns a default tracer for testing purposes.
func DefaultTracer() trace.Tracer {
	return NewTracerProvider().Tracer("")
}
