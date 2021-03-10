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

package global_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/trace"
)

func TestTraceWithSDK(t *testing.T) {
	global.ResetForTest()

	ctx := context.Background()
	gtp := otel.GetTracerProvider()
	tracer1 := gtp.Tracer("pre")
	// This is started before an SDK was registered and should be dropped.
	_, span1 := tracer1.Start(ctx, "span1")

	sr := new(oteltest.SpanRecorder)
	tp := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	otel.SetTracerProvider(tp)

	// This span was started before initialization, it is expected to be dropped.
	span1.End()

	// The existing Tracer should have been configured to now use the configured SDK.
	_, span2 := tracer1.Start(ctx, "span2")
	span2.End()

	// The global TracerProvider should now create Tracers that also use the newly configured SDK.
	tracer2 := gtp.Tracer("post")
	_, span3 := tracer2.Start(ctx, "span3")
	span3.End()

	filterNames := func(spans []*oteltest.Span) []string {
		names := make([]string, len(spans))
		for i := range spans {
			names[i] = spans[i].Name()
		}
		return names
	}
	expected := []string{"span2", "span3"}
	assert.ElementsMatch(t, expected, filterNames(sr.Started()))
	assert.ElementsMatch(t, expected, filterNames(sr.Completed()))
}

type fnTracerProvider struct {
	tracer func(string, ...trace.TracerOption) trace.Tracer
}

func (fn fnTracerProvider) Tracer(instrumentationName string, opts ...trace.TracerOption) trace.Tracer {
	return fn.tracer(instrumentationName, opts...)
}

type fnTracer struct {
	start func(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span)
}

func (fn fnTracer) Start(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span) {
	return fn.start(ctx, spanName, opts...)
}

func TestTraceProviderDelegates(t *testing.T) {
	global.ResetForTest()

	// Retrieve the placeholder TracerProvider.
	gtp := otel.GetTracerProvider()

	// Configure it with a spy.
	called := false
	otel.SetTracerProvider(fnTracerProvider{
		tracer: func(name string, opts ...trace.TracerOption) trace.Tracer {
			called = true
			assert.Equal(t, "abc", name)
			assert.Equal(t, []trace.TracerOption{trace.WithInstrumentationVersion("xyz")}, opts)
			return trace.NewNoopTracerProvider().Tracer("")
		},
	})

	gtp.Tracer("abc", trace.WithInstrumentationVersion("xyz"))
	assert.True(t, called, "expected configured TraceProvider to be called")
}

func TestTraceProviderDelegatesConcurrentSafe(t *testing.T) {
	global.ResetForTest()

	// Retrieve the placeholder TracerProvider.
	gtp := otel.GetTracerProvider()

	done := make(chan struct{})
	quit := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-time.After(1 * time.Millisecond):
				gtp.Tracer("abc", trace.WithInstrumentationVersion("xyz"))
			case <-quit:
				return
			}
		}
	}()

	// Wait for the goroutine to make some calls before installing the provider.
	<-time.After(100 * time.Millisecond)

	// Configure it with a spy.
	called := int32(0)
	otel.SetTracerProvider(fnTracerProvider{
		tracer: func(name string, opts ...trace.TracerOption) trace.Tracer {
			newVal := atomic.AddInt32(&called, 1)
			assert.Equal(t, "abc", name)
			assert.Equal(t, []trace.TracerOption{trace.WithInstrumentationVersion("xyz")}, opts)
			if newVal == 10 {
				// Signal the goroutine to finish.
				close(quit)
			}
			return trace.NewNoopTracerProvider().Tracer("")
		},
	})

	// Wait for the go routine to finish
	<-done

	assert.LessOrEqual(t, int32(10), atomic.LoadInt32(&called), "expected configured TraceProvider to be called")
}

func TestTracerDelegatesConcurrentSafe(t *testing.T) {
	global.ResetForTest()

	// Retrieve the placeholder TracerProvider.
	gtp := otel.GetTracerProvider()
	tracer := gtp.Tracer("abc", trace.WithInstrumentationVersion("xyz"))

	done := make(chan struct{})
	quit := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case <-time.After(1 * time.Millisecond):
				tracer.Start(context.Background(), "name")
			case <-quit:
				return
			}
		}
	}()

	// Wait for the goroutine to make some calls before installing the provider.
	<-time.After(100 * time.Millisecond)

	// Configure it with a spy.
	called := int32(0)
	otel.SetTracerProvider(fnTracerProvider{
		tracer: func(name string, opts ...trace.TracerOption) trace.Tracer {
			assert.Equal(t, "abc", name)
			assert.Equal(t, []trace.TracerOption{trace.WithInstrumentationVersion("xyz")}, opts)
			return fnTracer{
				start: func(ctx context.Context, spanName string, opts ...trace.SpanOption) (context.Context, trace.Span) {
					newVal := atomic.AddInt32(&called, 1)
					assert.Equal(t, "name", spanName)
					if newVal == 10 {
						// Signal the goroutine to finish.
						close(quit)
					}
					return trace.NewNoopTracerProvider().Tracer("").Start(ctx, spanName)
				},
			}
		},
	})

	// Wait for the go routine to finish
	<-done

	assert.LessOrEqual(t, int32(10), atomic.LoadInt32(&called), "expected configured TraceProvider to be called")
}

func TestTraceProviderDelegatesSameInstance(t *testing.T) {
	global.ResetForTest()

	// Retrieve the placeholder TracerProvider.
	gtp := otel.GetTracerProvider()
	tracer := gtp.Tracer("abc", trace.WithInstrumentationVersion("xyz"))
	assert.Same(t, tracer, gtp.Tracer("abc", trace.WithInstrumentationVersion("xyz")))
	assert.Same(t, tracer, gtp.Tracer("abc", trace.WithInstrumentationVersion("xyz")))

	otel.SetTracerProvider(fnTracerProvider{
		tracer: func(name string, opts ...trace.TracerOption) trace.Tracer {
			return trace.NewNoopTracerProvider().Tracer("")
		},
	})

	assert.NotSame(t, tracer, gtp.Tracer("abc", trace.WithInstrumentationVersion("xyz")))
}
