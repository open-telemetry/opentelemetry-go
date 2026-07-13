// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"errors"
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

type shutdownContextKey struct{}

const shutdownContextValue = "shutdown context value"

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
	overlap     bool
}

func newBlockingProcessor(operation processorOperation) (*processor, *processorBlock) {
	proc := newProcessor("first")
	block := &processorBlock{
		started:  make(chan struct{}),
		release:  make(chan struct{}),
		finished: make(chan struct{}),
	}
	wait := func() {
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
	require.Eventually(t, provider.shutdownStarted, time.Second, time.Microsecond)
	select {
	case err := <-done:
		require.NoError(t, err)
		t.Fatal("Shutdown returned while a processor operation was blocked")
	default:
	}
	return done
}

func assertShutdownContext(t *testing.T, ctx context.Context) {
	t.Helper()
	assert.NoError(t, ctx.Err())
	_, hasDeadline := ctx.Deadline()
	assert.False(t, hasDeadline)
	assert.Nil(t, ctx.Done())
	assert.Equal(t, shutdownContextValue, ctx.Value(shutdownContextKey{}))
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
	firstErr := errors.New("first")
	secondErr := errors.New("second")
	var order []string
	first := newProcessor("first")
	first.shutdownFunc = func(context.Context) error {
		order = append(order, first.Name)
		return firstErr
	}
	second := newProcessor("second")
	second.shutdownFunc = func(context.Context) error {
		order = append(order, second.Name)
		return secondErr
	}
	provider := NewLoggerProvider(
		WithProcessor(first),
		WithProcessor(second),
	)

	for range 2 {
		err := provider.Shutdown(t.Context())
		assert.ErrorIs(t, err, firstErr)
		assert.ErrorIs(t, err, secondErr)
	}
	assert.Equal(t, []string{"first", "second"}, order)
	assert.Equal(t, 1, first.shutdownCalls, "first processor Shutdown calls")
	assert.Equal(t, 1, second.shutdownCalls, "second processor Shutdown calls")
	assert.Zero(t, first.forceFlushCalls, "first processor ForceFlush calls")
	assert.Zero(t, second.forceFlushCalls, "second processor ForceFlush calls")
}

func TestLoggerProviderShutdownHonorsContextConcurrentSafe(t *testing.T) {
	t.Run("CanceledWhileWaiting", func(t *testing.T) {
		first, block := newBlockingProcessor(processorOnEmit)
		t.Cleanup(block.unblock)
		second := newProcessor("second")
		var shutdownCtx context.Context
		shutdownCalled := make(chan struct{})
		second.shutdownFunc = func(ctx context.Context) error {
			shutdownCtx = ctx
			close(shutdownCalled)
			return nil
		}
		provider := NewLoggerProvider(WithProcessor(first), WithProcessor(second))
		logger := provider.Logger("test")

		emitDone := make(chan struct{})
		go func() {
			logger.Emit(t.Context(), log.Record{})
			close(emitDone)
		}()
		<-block.started

		valueCtx := context.WithValue(t.Context(), shutdownContextKey{}, shutdownContextValue)
		ctx, cancel := context.WithTimeout(valueCtx, time.Hour)
		shutdownDone := make(chan error, 1)
		go func() {
			shutdownDone <- provider.Shutdown(ctx)
		}()
		require.Eventually(t, provider.shutdownStarted, time.Second, time.Microsecond)
		cancel()

		select {
		case err := <-shutdownDone:
			assert.ErrorIs(t, err, context.Canceled)
		case <-time.After(time.Second):
			t.Fatal("Shutdown did not honor context cancellation")
		}
		select {
		case <-shutdownCalled:
			t.Fatal("processor Shutdown called before admitted operation ended")
		default:
		}

		block.unblock()
		select {
		case <-emitDone:
		case <-time.After(time.Second):
			t.Fatal("Emit did not finish")
		}
		select {
		case <-shutdownCalled:
		case <-time.After(time.Second):
			t.Fatal("processor Shutdown was abandoned after context cancellation")
		}
		assert.Equal(t, 1, second.shutdownCalls)
		assertShutdownContext(t, shutdownCtx)

		require.NoError(t, provider.Shutdown(t.Context()))
		assert.Equal(t, 1, second.shutdownCalls)
	})

	t.Run("CanceledDuringProcessorShutdown", func(t *testing.T) {
		first := newBlockingShutdownProcessor()
		second := newProcessor("second")
		second.Err = assert.AnError
		var shutdownCtx context.Context
		shutdownCalled := make(chan struct{})
		second.shutdownFunc = func(ctx context.Context) error {
			shutdownCtx = ctx
			close(shutdownCalled)
			return second.Err
		}
		provider := NewLoggerProvider(
			WithProcessor(first),
			WithProcessor(second),
		)

		valueCtx := context.WithValue(t.Context(), shutdownContextKey{}, shutdownContextValue)
		ctx, cancel := context.WithCancel(valueCtx)
		shutdownDone := make(chan error, 1)
		go func() {
			shutdownDone <- provider.Shutdown(ctx)
		}()
		<-first.started

		released := false
		defer func() {
			if !released {
				close(first.release)
			}
		}()

		cancel()
		select {
		case err := <-shutdownDone:
			assert.ErrorIs(t, err, context.Canceled)
		case <-time.After(time.Second):
			t.Fatal("Shutdown did not honor context cancellation")
		}
		select {
		case <-shutdownCalled:
			t.Fatal("second processor Shutdown called while first was active")
		default:
		}

		retryCtx, retryCancel := context.WithCancel(t.Context())
		retryCancel()
		assert.ErrorIs(
			t,
			provider.Shutdown(retryCtx),
			context.Canceled,
			"canceled retry did not honor its context",
		)

		retryStarted := make(chan struct{})
		retryDone := make(chan error, 1)
		go func() {
			close(retryStarted)
			retryDone <- provider.Shutdown(t.Context())
		}()
		<-retryStarted
		select {
		case err := <-retryDone:
			t.Fatalf("retry returned before processor shutdown completed: %v", err)
		default:
		}

		close(first.release)
		released = true
		select {
		case <-shutdownCalled:
		case <-time.After(time.Second):
			t.Fatal("processor shutdown sequence was abandoned after context cancellation")
		}
		select {
		case err := <-retryDone:
			assert.ErrorIs(t, err, assert.AnError, "retry did not receive processor error")
		case <-time.After(time.Second):
			t.Fatal("retry did not wait for processor shutdown completion")
		}

		assert.Equal(t, int64(1), first.calls.Load())
		assert.Equal(t, 1, second.shutdownCalls)
		assertShutdownContext(t, shutdownCtx)
	})
}

func TestShutdownStateWait(t *testing.T) {
	done := func(err error) *shutdownState {
		state := &shutdownState{done: make(chan struct{})}
		state.complete(err)
		return state
	}
	pending := &shutdownState{done: make(chan struct{})}
	canceled, cancel := context.WithCancel(t.Context())
	cancel()

	tests := []struct {
		name    string
		ctx     context.Context
		state   *shutdownState
		wantErr error
	}{
		{
			name:  "Done",
			ctx:   t.Context(),
			state: done(nil),
		},
		{
			name:    "Error",
			ctx:     t.Context(),
			state:   done(assert.AnError),
			wantErr: assert.AnError,
		},
		{
			name:    "Canceled",
			ctx:     canceled,
			state:   pending,
			wantErr: context.Canceled,
		},
		{
			name:  "DoneTakesPriority",
			ctx:   canceled,
			state: done(nil),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.state.wait(test.ctx)
			if test.wantErr != nil {
				assert.ErrorIs(t, err, test.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type blockingShutdownProcessor struct {
	processor

	calls       atomic.Int64
	started     chan struct{}
	startedOnce sync.Once
	release     chan struct{}
}

func newBlockingShutdownProcessor() *blockingShutdownProcessor {
	return &blockingShutdownProcessor{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
}

func (p *blockingShutdownProcessor) Shutdown(context.Context) error {
	p.calls.Add(1)
	p.startedOnce.Do(func() { close(p.started) })
	<-p.release
	return p.Err
}

func TestLoggerProviderShutdownOnceConcurrentSafe(t *testing.T) {
	proc := newBlockingShutdownProcessor()
	proc.Err = assert.AnError
	provider := NewLoggerProvider(WithProcessor(proc))
	ctx := t.Context()

	const shutdowns = 100
	shutdownDone := make(chan error, shutdowns)
	for range shutdowns {
		go func() {
			shutdownDone <- provider.Shutdown(ctx)
		}()
	}

	<-proc.started
	close(proc.release)
	for range shutdowns {
		select {
		case err := <-shutdownDone:
			assert.ErrorIs(t, err, assert.AnError)
		case <-time.After(time.Second):
			t.Fatal("concurrent Shutdown did not receive processor result")
		}
	}

	assert.ErrorIs(t, provider.Shutdown(ctx), assert.AnError)
	assert.Equal(t, int64(1), proc.calls.Load())
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
	var calls []string
	second.forceFlushFunc = func(context.Context) error {
		calls = append(calls, "ForceFlush")
		return nil
	}
	second.shutdownFunc = func(context.Context) error {
		calls = append(calls, "Shutdown")
		return nil
	}
	provider := NewLoggerProvider(
		WithProcessor(first),
		WithProcessor(second),
	)
	ctx := t.Context()

	flushDone := make(chan error)
	go func() {
		flushDone <- provider.ForceFlush(ctx)
	}()

	<-block.started
	shutdownDone := shutdownWhileBlocked(t, provider)
	require.NoError(t, provider.ForceFlush(ctx))
	assert.Equal(t, 1, first.forceFlushCalls, "processor ForceFlush called after Shutdown started")
	assert.Zero(t, second.forceFlushCalls, "processor ForceFlush called after Shutdown started")

	block.unblock()
	flushErr := <-flushDone
	shutdownErr := <-shutdownDone

	require.NoError(t, shutdownErr)
	require.NoError(t, flushErr)
	assert.False(t, block.overlap)
	assert.Equal(t, 1, first.shutdownCalls)
	assert.Equal(t, 1, second.shutdownCalls)
	assert.Equal(t, 1, second.forceFlushCalls)
	assert.Equal(t, []string{"ForceFlush", "Shutdown"}, calls)
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
