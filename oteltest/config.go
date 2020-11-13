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
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel/trace"
)

// defaultSpanContextFunc returns the default SpanContextFunc.
func defaultSpanContextFunc() func(context.Context) trace.SpanContext {
	var traceID, spanID uint64 = 1, 1
	return func(ctx context.Context) trace.SpanContext {
		var sc trace.SpanContext
		if lsc := trace.SpanContextFromContext(ctx); lsc.IsValid() {
			sc = lsc
		} else if rsc := trace.RemoteSpanContextFromContext(ctx); rsc.IsValid() {
			sc = rsc
		} else {
			binary.BigEndian.PutUint64(sc.TraceID[:], atomic.AddUint64(&traceID, 1))
		}
		binary.BigEndian.PutUint64(sc.SpanID[:], atomic.AddUint64(&spanID, 1))
		return sc
	}
}

type config struct {
	// SpanContextFunc returns a SpanContext from an parent Context for a
	// new span.
	SpanContextFunc func(context.Context) trace.SpanContext

	// SpanRecorder keeps track of spans.
	SpanRecorder SpanRecorder
}

func newConfig(opts ...Option) config {
	conf := config{}
	for _, opt := range opts {
		opt.Apply(&conf)
	}
	if conf.SpanContextFunc == nil {
		conf.SpanContextFunc = defaultSpanContextFunc()
	}
	return conf
}

// Option applies an option to a config.
type Option interface {
	Apply(*config)
}

type spanContextFuncOption struct {
	SpanContextFunc func(context.Context) trace.SpanContext
}

func (o spanContextFuncOption) Apply(c *config) {
	c.SpanContextFunc = o.SpanContextFunc
}

// WithSpanContextFunc sets the SpanContextFunc used to generate a new Spans
// context from a parent SpanContext.
func WithSpanContextFunc(f func(context.Context) trace.SpanContext) Option {
	return spanContextFuncOption{f}
}

type spanRecorderOption struct {
	SpanRecorder SpanRecorder
}

func (o spanRecorderOption) Apply(c *config) {
	c.SpanRecorder = o.SpanRecorder
}

// WithSpanRecorder sets the SpanRecorder to use with the TracerProvider for
// testing.
func WithSpanRecorder(sr SpanRecorder) Option {
	return spanRecorderOption{sr}
}

// SpanRecorder performs operations to record a span as it starts and ends.
type SpanRecorder interface {
	// OnStart is called by the Tracer when it starts a Span.
	OnStart(span *Span)
	// OnEnd is called by the Span when it ends.
	OnEnd(span *Span)
}

// StandardSpanRecorder is a SpanRecorder that records all started and ended
// spans in an ordered recording. StandardSpanRecorder is designed to be
// concurrent safe and can by used by multiple goroutines.
type StandardSpanRecorder struct {
	startedMu sync.RWMutex
	started   []*Span

	doneMu sync.RWMutex
	done   []*Span
}

// OnStart records span as started.
func (ssr *StandardSpanRecorder) OnStart(span *Span) {
	ssr.startedMu.Lock()
	defer ssr.startedMu.Unlock()
	ssr.started = append(ssr.started, span)
}

// OnEnd records span as completed.
func (ssr *StandardSpanRecorder) OnEnd(span *Span) {
	ssr.doneMu.Lock()
	defer ssr.doneMu.Unlock()
	ssr.done = append(ssr.done, span)
}

// Started returns a copy of all started Spans in the order they were started.
func (ssr *StandardSpanRecorder) Started() []*Span {
	ssr.startedMu.RLock()
	defer ssr.startedMu.RUnlock()
	started := make([]*Span, len(ssr.started))
	for i := range ssr.started {
		started[i] = ssr.started[i]
	}
	return started
}

// Completed returns a copy of all ended Spans in the order they were ended.
func (ssr *StandardSpanRecorder) Completed() []*Span {
	ssr.doneMu.RLock()
	defer ssr.doneMu.RUnlock()
	done := make([]*Span, len(ssr.done))
	for i := range ssr.done {
		done[i] = ssr.done[i]
	}
	return done
}
