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

package oteltest

import (
	"context"
	"encoding/binary"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
)

// defaultSpanContextFunc returns the default SpanContextFunc.
func defaultSpanContextFunc() func(context.Context) otel.SpanContext {
	var traceID, spanID uint64 = 1, 1
	return func(ctx context.Context) otel.SpanContext {
		var sc otel.SpanContext
		if lsc := otel.SpanFromContext(ctx).SpanContext(); lsc.IsValid() {
			sc = lsc
		} else if rsc := otel.RemoteSpanContextFromContext(ctx); rsc.IsValid() {
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
	SpanContextFunc func(context.Context) otel.SpanContext

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

type Option interface {
	Apply(*config)
}

type spanContextFuncOption struct {
	SpanContextFunc func(context.Context) otel.SpanContext
}

func (o spanContextFuncOption) Apply(c *config) {
	c.SpanContextFunc = o.SpanContextFunc
}

func WithSpanContextFunc(f func(context.Context) otel.SpanContext) Option {
	return spanContextFuncOption{f}
}

type spanRecorderOption struct {
	SpanRecorder SpanRecorder
}

func (o spanRecorderOption) Apply(c *config) {
	c.SpanRecorder = o.SpanRecorder
}

func WithSpanRecorder(sr SpanRecorder) Option {
	return spanRecorderOption{sr}
}

type SpanRecorder interface {
	OnStart(span *Span)
	OnEnd(span *Span)
}

type StandardSpanRecorder struct {
	startedMu sync.RWMutex
	started   []*Span

	doneMu sync.RWMutex
	done   []*Span
}

func (ssr *StandardSpanRecorder) OnStart(span *Span) {
	ssr.startedMu.Lock()
	defer ssr.startedMu.Unlock()
	ssr.started = append(ssr.started, span)
}

func (ssr *StandardSpanRecorder) OnEnd(span *Span) {
	ssr.doneMu.Lock()
	defer ssr.doneMu.Unlock()
	ssr.done = append(ssr.done, span)
}

func (ssr *StandardSpanRecorder) Started() []*Span {
	ssr.startedMu.RLock()
	defer ssr.startedMu.RUnlock()
	started := make([]*Span, len(ssr.started))
	for i := range ssr.started {
		started[i] = ssr.started[i]
	}
	return started
}

func (ssr *StandardSpanRecorder) Completed() []*Span {
	ssr.doneMu.RLock()
	defer ssr.doneMu.RUnlock()
	done := make([]*Span, len(ssr.done))
	for i := range ssr.done {
		done[i] = ssr.done[i]
	}
	return done
}
