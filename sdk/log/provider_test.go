// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/sdk/resource"
)

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

func (p *processor) OnEmit(ctx context.Context, r *Record) error {
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

type filterProcessor struct {
	*processor

	enabled bool
}

func newFilterProcessor(name string, enabled bool) *filterProcessor {
	return &filterProcessor{
		processor: newProcessor(name),
		enabled:   enabled,
	}
}

func (p *filterProcessor) Enabled(context.Context, Record) bool {
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
			},
			want: &LoggerProvider{
				resource:                  res,
				processors:                []Processor{p0, p1},
				attributeCountLimit:       attrCntLim,
				attributeValueLengthLimit: attrValLenLim,
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

func TestLoggerProviderConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	p := NewLoggerProvider(WithProcessor(newProcessor("0")))
	const name = "testLogger"
	ctx := context.Background()
	for i := 0; i < goRoutineN; i++ {
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
	keysAndValues []interface{}
}

func (l *logSink) Enabled(int) bool { return true }

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
		assert.Equal(t, "", l.keysAndValues[1], "logged name")
	})

	t.Run("Stopped", func(t *testing.T) {
		ctx := context.Background()
		p := NewLoggerProvider()
		_ = p.Shutdown(ctx)
		l := p.Logger("testing")

		assert.NotNil(t, l)
		assert.IsType(t, noop.Logger{}, l)
	})

	t.Run("SameLoggers", func(t *testing.T) {
		p := NewLoggerProvider()

		l0, l1 := p.Logger("l0"), p.Logger("l1")
		l2, l3 := p.Logger("l0"), p.Logger("l1")

		assert.Same(t, l0, l2)
		assert.Same(t, l1, l3)
	})
}

func TestLoggerProviderShutdown(t *testing.T) {
	t.Run("Once", func(t *testing.T) {
		proc := newProcessor("")
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := context.Background()
		require.NoError(t, p.Shutdown(ctx))
		require.Equal(t, 1, proc.shutdownCalls, "processor Shutdown not called")

		require.NoError(t, p.Shutdown(ctx))
		assert.Equal(t, 1, proc.shutdownCalls, "processor Shutdown called multiple times")
	})

	t.Run("Error", func(t *testing.T) {
		proc := newProcessor("")
		proc.Err = assert.AnError
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := context.Background()
		assert.ErrorIs(t, p.Shutdown(ctx), assert.AnError, "processor error not returned")
	})
}

func TestLoggerProviderForceFlush(t *testing.T) {
	t.Run("Stopped", func(t *testing.T) {
		proc := newProcessor("")
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := context.Background()
		require.NoError(t, p.ForceFlush(ctx))
		require.Equal(t, 1, proc.forceFlushCalls, "processor ForceFlush not called")

		require.NoError(t, p.Shutdown(ctx))

		require.NoError(t, p.ForceFlush(ctx))
		assert.Equal(t, 1, proc.forceFlushCalls, "processor ForceFlush called after Shutdown")
	})

	t.Run("Multi", func(t *testing.T) {
		proc := newProcessor("")
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := context.Background()
		require.NoError(t, p.ForceFlush(ctx))
		require.Equal(t, 1, proc.forceFlushCalls, "processor ForceFlush not called")

		require.NoError(t, p.ForceFlush(ctx))
		assert.Equal(t, 2, proc.forceFlushCalls, "processor ForceFlush not called multiple times")
	})

	t.Run("Error", func(t *testing.T) {
		proc := newProcessor("")
		proc.Err = assert.AnError
		p := NewLoggerProvider(WithProcessor(proc))

		ctx := context.Background()
		assert.ErrorIs(t, p.ForceFlush(ctx), assert.AnError, "processor error not returned")
	})
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
	loggers[0].Enabled(context.Background(), log.Record{})
}
