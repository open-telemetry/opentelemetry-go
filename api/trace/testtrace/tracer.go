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

package testtrace

import (
	"context"

	"go.opentelemetry.io/api/trace"
)

var _ trace.Tracer = (*Tracer)(nil)

type Tracer struct {
	generator Generator
}

func NewTracer(opts ...TracerOption) *Tracer {
	c := newTracerConfig(opts...)

	return &Tracer{
		generator: c.generator,
	}
}

func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	return nil, nil
}

func (t *Tracer) WithSpan(ctx context.Context, name string, body func(ctx context.Context) error) error {
	return nil
}

type TracerOption func(*tracerConfig)

func TracerWithGenerator(generator Generator) TracerOption {
	return func(c *tracerConfig) {
		c.generator = generator
	}
}

type tracerConfig struct {
	generator Generator
}

func newTracerConfig(opts ...TracerOption) tracerConfig {
	var c tracerConfig
	defaultOpts := []TracerOption{
		TracerWithGenerator(NewCountGenerator()),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(&c)
	}

	return c
}
