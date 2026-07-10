// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"strconv"
	"sync"
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

	shutdownCalls   int
	forceFlushCalls int

	records []Record
}

func newProcessor(name string) *processor {
	return &processor{Name: name}
}

func (*processor) Enabled(context.Context, EnabledParameters) bool {
	return true
}

func (p *processor) OnEmit(_ context.Context, r *Record) error {
	if p.Err != nil {
		return p.Err
	}

	p.records = append(p.records, *r)
	return nil
}

func (p *processor) Shutdown(context.Context) error {
	p.shutdownCalls++
	return p.Err
}

func (p *processor) ForceFlush(context.Context) error {
	p.forceFlushCalls++
	return p.Err
}

type fltrProcessor struct {
	*processor

	enabled bool
	params  []EnabledParameters
}

type shutdownOrderProcessor struct {
	*processor
	order *[]string
}

type shutdownContextProcessor struct {
	*processor

	mu             sync.Mutex
	active         bool
	overlap        bool
	contextErr     error
	emitStarted    chan struct{}
	emitRelease    chan struct{}
	shutdownCalled chan struct{}
}

func newShutdownContextProcessor() *shutdownContextProcessor {
	return &shutdownContextProcessor{
		processor:      newProcessor(""),
		emitStarted:    make(chan struct{}),
		emitRelease:    make(chan struct{}),
		shutdownCalled: make(chan struct{}),
	}
}

func (p *shutdownContextProcessor) OnEmit(ctx context.Context, r *Record) error {
	p.mu.Lock()
	p.active = true
	close(p.emitStarted)
	p.mu.Unlock()

	<-p.emitRelease
	err := p.processor.OnEmit(ctx, r)

	p.mu.Lock()
	p.active = false
	p.mu.Unlock()
	return err
}

func (p *shutdownContextProcessor) Shutdown(ctx context.Context) error {
	p.mu.Lock()
	p.overlap = p.active
	p.contextErr = ctx.Err()
	p.mu.Unlock()
	close(p.shutdownCalled)
	return p.processor.Shutdown(ctx)
}

func (p *shutdownContextProcessor) result() (overlap bool, contextErr error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.overlap, p.contextErr
}

func (p *shutdownOrderProcessor) Shutdown(context.Context) error {
	p.shutdownCalls++
	*p.order = append(*p.order, p.Name)
	return p.Err
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
		assert.Zero(t, proc.forceFlushCalls, "processor ForceFlush called by provider Shutdown")
	})

	t.Run("RegistrationOrder", func(t *testing.T) {
		var order []string
		first := &shutdownOrderProcessor{
			processor: newProcessor("first"),
			order:     &order,
		}
		second := &shutdownOrderProcessor{
			processor: newProcessor("second"),
			order:     &order,
		}
		provider := NewLoggerProvider(
			WithProcessor(first),
			WithProcessor(second),
		)

		require.NoError(t, provider.Shutdown(t.Context()))
		assert.Equal(t, []string{"first", "second"}, order)
	})

	t.Run("Error", func(t *testing.T) {
		proc := newProcessor("")
		proc.Err = assert.AnError
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := t.Context()
		assert.ErrorIs(t, p.Shutdown(ctx), assert.AnError, "processor error not returned")
	})
}

func TestLoggerProviderShutdownWaitsForAdmittedProcessorCallConcurrentSafe(t *testing.T) {
	proc := newShutdownContextProcessor()
	provider := NewLoggerProvider(WithProcessor(proc))

	require.True(t, provider.beginProcessorCall())
	callEnded := false
	defer func() {
		if !callEnded {
			provider.endProcessorCall()
		}
	}()

	ctx := t.Context()

	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- provider.Shutdown(ctx)
	}()

	require.Eventually(t, provider.stopped.Load, time.Second, time.Microsecond)
	assert.False(t, provider.beginProcessorCall(), "processor call admitted after Shutdown")
	select {
	case err := <-shutdownDone:
		require.NoError(t, err)
		t.Fatal("Shutdown returned before the admitted processor call started")
	default:
	}

	emitDone := make(chan error, 1)
	emitReleased := false
	defer func() {
		if !emitReleased {
			close(proc.emitRelease)
		}
	}()
	go func() {
		emitDone <- proc.OnEmit(t.Context(), new(Record))
	}()
	<-proc.emitStarted
	select {
	case <-proc.shutdownCalled:
		t.Fatal("processor Shutdown called while OnEmit was active")
	default:
	}
	close(proc.emitRelease)
	emitReleased = true
	require.NoError(t, <-emitDone)
	select {
	case <-proc.shutdownCalled:
		t.Fatal("processor Shutdown called before the admitted call ended")
	default:
	}
	provider.endProcessorCall()
	callEnded = true

	require.NoError(t, <-shutdownDone)
	assert.Len(t, proc.records, 1)
	assert.Equal(t, 1, proc.shutdownCalls)
	overlap, contextErr := proc.result()
	assert.False(t, overlap)
	assert.NoError(t, contextErr)
}

func TestLoggerProviderShutdownHonorsContext(t *testing.T) {
	t.Run("CanceledWhileWaiting", func(t *testing.T) {
		first := &blockingForceFlushProcessor{
			processor: newProcessor("first"),
			started:   make(chan struct{}),
			release:   make(chan struct{}),
			finished:  make(chan struct{}),
		}
		second := &orderedForceFlushProcessor{processor: newProcessor("second")}
		provider := NewLoggerProvider(
			WithProcessor(first),
			WithProcessor(second),
		)

		flushDone := make(chan error, 1)
		go func() {
			flushDone <- provider.ForceFlush(t.Context())
		}()
		<-first.started

		released := false
		defer func() {
			if !released {
				close(first.release)
			}
		}()

		ctx, cancel := context.WithCancel(t.Context())
		shutdownDone := make(chan error, 1)
		go func() {
			shutdownDone <- provider.Shutdown(ctx)
		}()
		require.Eventually(t, provider.stopped.Load, time.Second, time.Microsecond)
		cancel()

		select {
		case err := <-shutdownDone:
			assert.ErrorIs(t, err, context.Canceled)
		case <-time.After(time.Second):
			t.Fatal("Shutdown did not honor context cancellation")
		}
		assert.Zero(t, first.shutdownCalls)
		assert.Zero(t, second.shutdownCalls)

		close(first.release)
		released = true
		require.NoError(t, <-flushDone)
		assert.Equal(t, 1, second.forceFlushCalls)
		assert.Equal(t, []string{"ForceFlush"}, second.calls)

		require.NoError(t, provider.Shutdown(t.Context()))
		assert.Zero(t, first.shutdownCalls)
		assert.Zero(t, second.shutdownCalls)
	})
}

type shutdownDetectingProcessor struct {
	processor

	mu      sync.Mutex
	running bool
	calls   int
	overlap bool
	started chan struct{}
	release chan struct{}
}

func newShutdownDetectingProcessor() *shutdownDetectingProcessor {
	return &shutdownDetectingProcessor{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
}

func (p *shutdownDetectingProcessor) Shutdown(context.Context) error {
	p.mu.Lock()
	p.calls++
	if p.running {
		p.overlap = true
	}
	p.running = true
	if p.calls == 1 {
		close(p.started)
	}
	p.mu.Unlock()

	<-p.release

	p.mu.Lock()
	p.running = false
	p.mu.Unlock()
	return nil
}

func (p *shutdownDetectingProcessor) result() (calls int, overlap bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.calls, p.overlap
}

func TestLoggerProviderShutdownOnceConcurrentSafe(t *testing.T) {
	proc := newShutdownDetectingProcessor()
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
	for range shutdowns - 1 {
		require.NoError(t, <-shutdownDone)
	}
	close(proc.release)
	require.NoError(t, <-shutdownDone)

	for range shutdowns {
		require.NoError(t, provider.Shutdown(ctx))
	}

	calls, overlap := proc.result()
	assert.Equal(t, 1, calls)
	assert.False(t, overlap)
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

type blockingForceFlushProcessor struct {
	*processor
	started  chan struct{}
	release  chan struct{}
	finished chan struct{}
	overlap  bool
}

func (p *blockingForceFlushProcessor) ForceFlush(context.Context) error {
	p.forceFlushCalls++
	close(p.started)
	<-p.release
	close(p.finished)
	return p.Err
}

func (p *blockingForceFlushProcessor) Shutdown(ctx context.Context) error {
	select {
	case <-p.finished:
	default:
		p.overlap = true
	}
	return p.processor.Shutdown(ctx)
}

type orderedForceFlushProcessor struct {
	*processor
	calls []string
}

func (p *orderedForceFlushProcessor) ForceFlush(context.Context) error {
	p.forceFlushCalls++
	p.calls = append(p.calls, "ForceFlush")
	return p.Err
}

func (p *orderedForceFlushProcessor) Shutdown(context.Context) error {
	p.shutdownCalls++
	p.calls = append(p.calls, "Shutdown")
	return p.Err
}

func TestLoggerProviderForceFlushShutdownConcurrentSafe(t *testing.T) {
	first := &blockingForceFlushProcessor{
		processor: newProcessor("first"),
		started:   make(chan struct{}),
		release:   make(chan struct{}),
		finished:  make(chan struct{}),
	}
	second := &orderedForceFlushProcessor{processor: newProcessor("second")}
	provider := NewLoggerProvider(
		WithProcessor(first),
		WithProcessor(second),
	)
	ctx := t.Context()

	flushDone := make(chan error)
	go func() {
		flushDone <- provider.ForceFlush(ctx)
	}()

	<-first.started
	shutdownDone := make(chan error, 1)
	go func() {
		shutdownDone <- provider.Shutdown(ctx)
	}()
	require.Eventually(t, provider.stopped.Load, time.Second, time.Microsecond)
	select {
	case err := <-shutdownDone:
		close(first.release)
		<-flushDone
		require.NoError(t, err)
		t.Fatal("Shutdown returned before ForceFlush completed")
	default:
	}

	close(first.release)
	flushErr := <-flushDone
	shutdownErr := <-shutdownDone

	require.NoError(t, shutdownErr)
	require.NoError(t, flushErr)
	assert.False(t, first.overlap)
	assert.Equal(t, 1, first.shutdownCalls)
	assert.Equal(t, 1, second.shutdownCalls)
	assert.Equal(t, 1, second.forceFlushCalls)
	assert.Equal(t, []string{"ForceFlush", "Shutdown"}, second.calls)
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
