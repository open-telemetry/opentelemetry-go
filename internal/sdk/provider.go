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

package sdk

import (
	"context"
	"encoding/binary"
	"sync/atomic"

	"go.opentelemetry.io/otel/api/trace"
)

func defaultSpanContextFunc() func(context.Context) trace.SpanContext {
	var traceID, spanID uint64 = 1, 1
	return func(ctx context.Context) trace.SpanContext {
		var sc trace.SpanContext
		if lsc := trace.SpanFromContext(ctx).SpanContext(); lsc.IsValid() {
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

type TracingConfig struct {
	// SpanContextFunc returns a SpanContext from an parent Context for a
	// new span.
	SpanContextFunc func(context.Context) trace.SpanContext

	// SpanRecorder keeps track of spans.
	SpanRecorder SpanRecorder
}

type TracingOption interface {
	Apply(*TracingConfig)
}

type spanContextFuncOption struct {
	SpanContextFunc func(context.Context) trace.SpanContext
}

func (o spanContextFuncOption) Apply(c *TracingConfig) {
	c.SpanContextFunc = o.SpanContextFunc
}

func WithSpanContextFunc(f func(context.Context) trace.SpanContext) TracingOption {
	return spanContextFuncOption{f}
}

type spanRecorderOption struct {
	SpanRecorder SpanRecorder
}

func (o spanRecorderOption) Apply(c *TracingConfig) {
	c.SpanRecorder = o.SpanRecorder
}

func WithSpanRecorder(sr SpanRecorder) TracingOption {
	return spanRecorderOption{sr}
}

type SpanRecorder interface {
	OnStart(*Span)
	OnEnd(*Span)
}

type TraceProvider struct {
	Config TracingConfig
}

var _ trace.Provider = TraceProvider{}

func NewTraceProvider(opts ...TracingOption) TraceProvider {
	conf := TracingConfig{}
	for _, opt := range opts {
		opt.Apply(&conf)
	}
	if conf.SpanContextFunc == nil {
		conf.SpanContextFunc = defaultSpanContextFunc()
	}
	return TraceProvider{Config: conf}
}

func (p TraceProvider) Tracer(instName string, opts ...trace.TracerOption) trace.Tracer {
	conf := new(trace.TracerConfig)
	for _, o := range opts {
		o(conf)
	}
	return &Tracer{
		Name:    instName,
		Version: conf.InstrumentationVersion,
		Config:  &p.Config,
	}
}
