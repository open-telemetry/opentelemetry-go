// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/sdk/resource"
)

const envVarResourceAttributes = "OTEL_RESOURCE_ATTRIBUTES"

type processor struct {
	Name string
	Err  error

	enabledFunc    func(context.Context, EnabledParameters) bool
	onEmitFunc     func(context.Context, *Record) error
	shutdownFunc   func(context.Context) error
	forceFlushFunc func(context.Context) error

	shutdownCalls   int
	forceFlushCalls int

	records []Record
}

func newProcessor(name string) *processor {
	return &processor{Name: name}
}

func (p *processor) Enabled(ctx context.Context, param EnabledParameters) bool {
	if p.enabledFunc != nil {
		return p.enabledFunc(ctx, param)
	}
	return true
}

func (p *processor) OnEmit(ctx context.Context, r *Record) error {
	if p.onEmitFunc != nil {
		return p.onEmitFunc(ctx, r)
	}
	if p.Err != nil {
		return p.Err
	}

	p.records = append(p.records, *r)
	return nil
}

func (p *processor) Shutdown(ctx context.Context) error {
	p.shutdownCalls++
	if p.shutdownFunc != nil {
		return p.shutdownFunc(ctx)
	}
	return p.Err
}

func (p *processor) ForceFlush(ctx context.Context) error {
	p.forceFlushCalls++
	if p.forceFlushFunc != nil {
		return p.forceFlushFunc(ctx)
	}
	return p.Err
}

type fltrProcessor struct {
	*processor

	enabled bool
	params  []EnabledParameters
}

type processorOperation int

const (
	processorEnabled processorOperation = iota
	processorOnEmit
	processorForceFlush
)

type processorBlock struct {
	started     chan struct{}
	release     chan struct{}
	finished    chan struct{}
	releaseOnce sync.Once
	calls       atomic.Int64
	overlap     bool
}

type observedDoneContext struct {
	context.Context
	doneCalled chan struct{}
	once       sync.Once
}

func (c *observedDoneContext) Done() <-chan struct{} {
	c.once.Do(func() { close(c.doneCalled) })
	return c.Context.Done()
}

func newBlockingProcessor(operation processorOperation) (*processor, *processorBlock) {
	proc := newProcessor("first")
	block := &processorBlock{
		started:  make(chan struct{}),
		release:  make(chan struct{}),
		finished: make(chan struct{}),
	}
	wait := func() {
		if block.calls.Add(1) != 1 {
			return
		}
		close(block.started)
		<-block.release
		close(block.finished)
	}

	switch operation {
	case processorEnabled:
		proc.enabledFunc = func(context.Context, EnabledParameters) bool {
			wait()
			return false
		}
	case processorOnEmit:
		proc.onEmitFunc = func(context.Context, *Record) error {
			wait()
			return proc.Err
		}
	case processorForceFlush:
		proc.forceFlushFunc = func(context.Context) error {
			wait()
			return proc.Err
		}
	}
	proc.shutdownFunc = func(context.Context) error {
		select {
		case <-block.finished:
		default:
			block.overlap = true
		}
		return proc.Err
	}
	return proc, block
}

func (b *processorBlock) unblock() {
	b.releaseOnce.Do(func() { close(b.release) })
}

func shutdownWhileBlocked(t *testing.T, provider *LoggerProvider) <-chan error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- provider.Shutdown(t.Context()) }()
	require.Eventually(t, provider.stopped.Load, time.Second, time.Microsecond)
	select {
	case err := <-done:
		require.NoError(t, err)
		t.Fatal("Shutdown returned while a processor operation was blocked")
	default:
	}
	return done
}

func newFltrProcessor(name string, enabled bool) *fltrProcessor {
	return &fltrProcessor{
		processor: newProcessor(name),
		enabled:   enabled,
	}
}

func (p *fltrProcessor) Enabled(_ context.Context, params EnabledParameters) bool {
	p.params = append(p.params, params)
	return p.enabled
}

func TestNewLoggerProviderConfiguration(t *testing.T) {
	t.Cleanup(func(orig otel.ErrorHandler) func() {
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			t.Log(err)
		}))
		return func() { otel.SetErrorHandler(orig) }
	}(otel.GetErrorHandler()))

	res := resource.NewSchemaless(attribute.String("key", "value"))
	p0, p1 := newProcessor("0"), newProcessor("1")
	attrCntLim := 12
	attrValLenLim := 21

	testcases := []struct {
		name    string
		envars  map[string]string
		options []LoggerProviderOption
		want    *LoggerProvider
	}{
		{
			name: "Defaults",
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       defaultAttrCntLim,
				attributeValueLengthLimit: defaultAttrValLenLim,
			},
		},
		{
			name: "Options",
			options: []LoggerProviderOption{
				WithResource(res),
				WithProcessor(p0),
				WithProcessor(p1),
				WithAttributeCountLimit(attrCntLim),
				WithAttributeValueLengthLimit(attrValLenLim),
				WithAllowKeyDuplication(),
			},
			want: &LoggerProvider{
				resource:                  res,
				processors:                []Processor{p0, p1},
				attributeCountLimit:       attrCntLim,
				attributeValueLengthLimit: attrValLenLim,
				allowDupKeys:              true,
			},
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarAttrCntLim:    strconv.Itoa(attrCntLim),
				envarAttrValLenLim: strconv.Itoa(attrValLenLim),
			},
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       attrCntLim,
				attributeValueLengthLimit: attrValLenLim,
			},
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarAttrCntLim:    "invalid attributeCountLimit",
				envarAttrValLenLim: "invalid attributeValueLengthLimit",
			},
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       defaultAttrCntLim,
				attributeValueLengthLimit: defaultAttrValLenLim,
			},
		},
		{
			name: "Precedence",
			envars: map[string]string{
				envarAttrCntLim:    strconv.Itoa(100),
				envarAttrValLenLim: strconv.Itoa(101),
			},
			options: []LoggerProviderOption{
				// These override the environment variables.
				WithAttributeCountLimit(attrCntLim),
				WithAttributeValueLengthLimit(attrValLenLim),
			},
			want: &LoggerProvider{
				resource:                  resource.Default(),
				attributeCountLimit:       attrCntLim,
				attributeValueLengthLimit: attrValLenLim,
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, NewLoggerProvider(tc.options...))
		})
	}
}

func mergeResource(t *testing.T, r1, r2 *resource.Resource) *resource.Resource {
	r, err := resource.Merge(r1, r2)
	assert.NoError(t, err)
	return r
}

func TestWithResource(t *testing.T) {
	t.Setenv(envVarResourceAttributes, "key=value,rk5=7")

	cases := []struct {
		name    string
		options []LoggerProviderOption
		want    *resource.Resource
		msg     string
	}{
		{
			name:    "explicitly empty resource",
			options: []LoggerProviderOption{WithResource(resource.Empty())},
			want:    resource.Environment(),
		},
		{
			name:    "uses default if no resource option",
			options: []LoggerProviderOption{},
			want:    resource.Default(),
		},
		{
			name: "explicit resource",
			options: []LoggerProviderOption{
				WithResource(resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk2", 5))),
			},
			want: mergeResource(
				t,
				resource.Environment(),
				resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk2", 5)),
			),
		},
		{
			name: "last resource wins",
			options: []LoggerProviderOption{
				WithResource(resource.NewSchemaless(attribute.String("rk1", "vk1"), attribute.Int64("rk2", 5))),
				WithResource(resource.NewSchemaless(attribute.String("rk3", "rv3"), attribute.Int64("rk4", 10))),
			},
			want: mergeResource(
				t,
				resource.Environment(),
				resource.NewSchemaless(attribute.String("rk3", "rv3"), attribute.Int64("rk4", 10)),
			),
		},
		{
			name: "overlapping attributes with environment resource",
			options: []LoggerProviderOption{
				WithResource(resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk5", 10))),
			},
			want: mergeResource(
				t,
				resource.Environment(),
				resource.NewSchemaless(attribute.String("rk1", "rv1"), attribute.Int64("rk5", 10)),
			),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := newProviderConfig(tc.options).resource
			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("WithResource:\n  -got +want %s", diff)
			}
		})
	}
}

func TestMapDeduplication(t *testing.T) {
	dup := attribute.Map(
		"map",
		attribute.String("key", "first"),
		attribute.String("key", "second"),
	)
	dedup := attribute.Map("map", attribute.String("key", "second"))

	res := resource.NewSchemaless(dup)

	t.Run("Resource", func(t *testing.T) {
		got := newProviderConfig([]LoggerProviderOption{WithResource(res)}).resource
		assert.Equal(t, []attribute.KeyValue{dedup}, got.Attributes())
	})

	t.Run("ResourceAlwaysDeduplicates", func(t *testing.T) {
		got := newProviderConfig([]LoggerProviderOption{
			WithResource(res),
			WithAllowKeyDuplication(),
		}).resource
		assert.Equal(t, []attribute.KeyValue{dedup}, got.Attributes())
	})

	t.Run("Scope", func(t *testing.T) {
		p := newProcessor("processor")
		lp := NewLoggerProvider(WithProcessor(p))
		l := lp.Logger("scope", log.WithInstrumentationAttributes(dup))

		l.Emit(t.Context(), log.Record{})

		require.Len(t, p.records, 1)
		assert.Equal(t, attribute.NewSet(dedup), p.records[0].InstrumentationScope().Attributes)
	})

	t.Run("ScopeWithAllowKeyDuplication", func(t *testing.T) {
		p := newProcessor("processor")
		lp := NewLoggerProvider(WithProcessor(p), WithAllowKeyDuplication())
		l := lp.Logger("scope", log.WithInstrumentationAttributes(dup))

		l.Emit(t.Context(), log.Record{})

		require.Len(t, p.records, 1)
		assert.Equal(t, attribute.NewSet(dup), p.records[0].InstrumentationScope().Attributes)
	})
}

func TestLoggerProviderConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	p := NewLoggerProvider(WithProcessor(newProcessor("0")))
	const name = "testLogger"
	ctx := t.Context()
	for range goRoutineN {
		go func() {
			defer wg.Done()

			_ = p.Logger(name)
			_ = p.Shutdown(ctx)
			_ = p.ForceFlush(ctx)
		}()
	}

	wg.Wait()
}

type logSink struct {
	logr.LogSink

	level         int
	msg           string
	keysAndValues []any
}

func (*logSink) Enabled(int) bool { return true }

func (l *logSink) Info(level int, msg string, keysAndValues ...any) {
	l.level, l.msg, l.keysAndValues = level, msg, keysAndValues
	l.LogSink.Info(level, msg, keysAndValues...)
}

func TestLoggerProviderLogger(t *testing.T) {
	t.Run("InvalidName", func(t *testing.T) {
		l := &logSink{LogSink: testr.New(t).GetSink()}
		t.Cleanup(func(orig logr.Logger) func() {
			global.SetLogger(logr.New(l))
			return func() { global.SetLogger(orig) }
		}(global.GetLogger()))

		_ = NewLoggerProvider().Logger("")
		assert.Equal(t, 1, l.level, "logged level")
		assert.Equal(t, "Invalid Logger name.", l.msg, "logged message")
		require.Len(t, l.keysAndValues, 2, "logged key values")
		assert.Empty(t, l.keysAndValues[1], "logged name")
	})

	t.Run("Stopped", func(t *testing.T) {
		ctx := t.Context()
		p := NewLoggerProvider()
		_ = p.Shutdown(ctx)
		l := p.Logger("testing")

		assert.NotNil(t, l)
		assert.IsType(t, noop.Logger{}, l)
	})

	t.Run("SameLoggers", func(t *testing.T) {
		p := NewLoggerProvider()

		l0, l1, l2 := p.Logger(
			"l0",
		), p.Logger(
			"l1",
		), p.Logger(
			"l0",
			log.WithInstrumentationAttributes(attribute.String("foo", "bar")),
		)
		assert.NotSame(t, l0, l1)
		assert.NotSame(t, l0, l2)
		assert.NotSame(t, l1, l2)

		l3, l4, l5 := p.Logger(
			"l0",
		), p.Logger(
			"l1",
		), p.Logger(
			"l0",
			log.WithInstrumentationAttributes(attribute.String("foo", "bar")),
		)
		assert.Same(t, l0, l3)
		assert.Same(t, l1, l4)
		assert.Same(t, l2, l5)
	})
}

func TestLoggerProviderShutdown(t *testing.T) {
	t.Run("Once", func(t *testing.T) {
		proc := newProcessor("")
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := t.Context()
		require.NoError(t, p.Shutdown(ctx))
		require.Equal(t, 1, proc.shutdownCalls, "processor Shutdown not called")

		require.NoError(t, p.Shutdown(ctx))
		assert.Equal(t, 1, proc.shutdownCalls, "processor Shutdown called multiple times")
	})

	t.Run("Error", func(t *testing.T) {
		proc := newProcessor("")
		proc.Err = assert.AnError
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := t.Context()
		assert.ErrorIs(t, p.Shutdown(ctx), assert.AnError, "processor error not returned")
	})

	t.Run("CanceledWhileProcessorOperationActive", func(t *testing.T) {
		proc, block := newBlockingProcessor(processorForceFlush)
		t.Cleanup(block.unblock)
		provider := NewLoggerProvider(WithProcessor(proc))

		flushDone := make(chan error, 1)
		go func() { flushDone <- provider.ForceFlush(t.Context()) }()
		<-block.started

		baseCtx, cancel := context.WithCancel(t.Context())
		ctx := &observedDoneContext{
			Context:    baseCtx,
			doneCalled: make(chan struct{}),
		}
		shutdownDone := make(chan error, 1)
		go func() { shutdownDone <- provider.Shutdown(ctx) }()
		<-ctx.doneCalled
		cancel()

		select {
		case err := <-shutdownDone:
			assert.ErrorIs(t, err, context.Canceled)
		case <-time.After(time.Second):
			t.Fatal("Shutdown did not honor context cancellation while waiting")
		}
		assert.Zero(t, proc.shutdownCalls, "processor Shutdown called before active operation completed")

		block.unblock()
		require.NoError(t, <-flushDone)
	})
}

func TestLoggerProviderForceFlush(t *testing.T) {
	t.Run("Stopped", func(t *testing.T) {
		proc := newProcessor("")
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := t.Context()
		require.NoError(t, p.ForceFlush(ctx))
		require.Equal(t, 1, proc.forceFlushCalls, "processor ForceFlush not called")

		require.NoError(t, p.Shutdown(ctx))

		require.NoError(t, p.ForceFlush(ctx))
		assert.Equal(t, 1, proc.forceFlushCalls, "processor ForceFlush called after Shutdown")
	})

	t.Run("Multi", func(t *testing.T) {
		proc := newProcessor("")
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := t.Context()
		require.NoError(t, p.ForceFlush(ctx))
		require.Equal(t, 1, proc.forceFlushCalls, "processor ForceFlush not called")

		require.NoError(t, p.ForceFlush(ctx))
		assert.Equal(t, 2, proc.forceFlushCalls, "processor ForceFlush not called multiple times")
	})

	t.Run("Error", func(t *testing.T) {
		proc := newProcessor("")
		proc.Err = assert.AnError
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := t.Context()
		assert.ErrorIs(t, p.ForceFlush(ctx), assert.AnError, "processor error not returned")
	})
}

func TestLoggerProviderForceFlushShutdownConcurrentSafe(t *testing.T) {
	first, block := newBlockingProcessor(processorForceFlush)
	t.Cleanup(block.unblock)
	second := newProcessor("second")
	provider := NewLoggerProvider(
		WithProcessor(first),
		WithProcessor(second),
	)
	ctx := t.Context()

	flushDone := make(chan error, 1)
	go func() {
		flushDone <- provider.ForceFlush(ctx)
	}()

	<-block.started
	shutdownDone := shutdownWhileBlocked(t, provider)
	require.NoError(t, provider.ForceFlush(ctx))
	assert.Equal(t, int64(1), block.calls.Load(), "ForceFlush admitted after shutdown started")

	block.unblock()
	flushErr := <-flushDone
	shutdownErr := <-shutdownDone

	require.NoError(t, flushErr)
	require.NoError(t, shutdownErr)
	assert.False(t, block.overlap)
	assert.Equal(t, 1, first.forceFlushCalls)
	assert.Equal(t, 1, second.forceFlushCalls, "admitted ForceFlush did not complete")
	assert.Equal(t, 1, first.shutdownCalls)
	assert.Equal(t, 1, second.shutdownCalls)
}

func BenchmarkLoggerProviderLogger(b *testing.B) {
	p := NewLoggerProvider()
	names := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		names[i] = fmt.Sprintf("%d logger", i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	loggers := make([]log.Logger, b.N)
	for i := 0; i < b.N; i++ {
		loggers[i] = p.Logger(names[i])
	}

	b.StopTimer()
	loggers[0].Enabled(b.Context(), log.EnabledParameters{})
}

type testExperimentalOption struct {
	LoggerProviderOption
}

func (testExperimentalOption) Experimental() {}

func TestExperimentalOptionSafe(t *testing.T) {
	var opt testExperimentalOption

	assert.NotPanics(t, func() { _ = NewLoggerProvider(opt) })
}
