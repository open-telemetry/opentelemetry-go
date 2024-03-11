// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
)

var (
	ctx            = trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{TraceID: [16]byte{1}, SpanID: [8]byte{42}}))
	testBodyString = "log message"
	testFloat      = 1.2345
	testString     = "7e3b3b2aaeff56a7108fe11e154200dd/7819479873059528190"
	testInt        = 32768
	testBool       = true
	testTimestamp  = time.Date(1988, time.November, 17, 0, 0, 0, 0, time.UTC)
)

var runs = 5

func TestZeroAllocsSimple(t *testing.T) {
	provider := NewLoggerProvider(WithProcessor(NewSimpleProcessor(noopExporter{})))
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.Background())) })
	logger := slog.New(&slogHandler{provider.Logger("log/slog")})

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
			slog.String("string", testString),
			slog.Float64("float", testFloat),
			slog.Int("int", testInt),
			slog.Bool("bool", testBool),
			slog.String("string", testString),
		)
	}))
}

func TestZeroAllocsModifyProcessor(t *testing.T) {
	provider := NewLoggerProvider(WithProcessor(timestampDecorator{NewSimpleProcessor(noopExporter{})}))
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.Background())) })
	logger := slog.New(&slogHandler{provider.Logger("log/slog")})

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
			slog.String("string", testString),
			slog.Float64("float", testFloat),
			slog.Int("int", testInt),
			slog.Bool("bool", testBool),
			slog.String("string", testString),
		)
	}))
}

func TestZeroAllocsBatch(t *testing.T) {
	provider := NewLoggerProvider(WithProcessor(NewBatchingProcessor(noopExporter{})))
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.Background())) })
	logger := slog.New(&slogHandler{provider.Logger("log/slog")})

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
			slog.String("string", testString),
			slog.Float64("float", testFloat),
			slog.Int("int", testInt),
			slog.Bool("bool", testBool),
			slog.String("string", testString),
		)
	}))
}

func TestZeroAllocsNoSpan(t *testing.T) {
	provider := NewLoggerProvider(WithProcessor(NewSimpleProcessor(noopExporter{})))
	t.Cleanup(func() { assert.NoError(t, provider.Shutdown(context.Background())) })
	logger := slog.New(&slogHandler{provider.Logger("log/slog")})

	assert.Equal(t, 0.0, testing.AllocsPerRun(runs, func() {
		logger.LogAttrs(context.Background(), slog.LevelInfo, testBodyString,
			slog.String("string", testString),
			slog.Float64("float", testFloat),
			slog.Int("int", testInt),
			slog.Bool("bool", testBool),
			slog.String("string", testString),
		)
	}))
}

func Benchmark(b *testing.B) {
	for _, call := range []struct {
		name string
		f    func(*slog.Logger)
	}{
		{
			"no attrs",
			func(logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString)
			},
		},
		{
			"3 attrs",
			func(logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
				)
			},
		},
		{
			"5 attrs",
			func(logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
				)
			},
		},
		{
			"10 attrs",
			func(logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
				)
			},
		},
		{
			"40 attrs",
			func(logger *slog.Logger) {
				logger.LogAttrs(ctx, slog.LevelInfo, testBodyString,
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
					slog.String("string", testString),
					slog.Float64("float", testFloat),
					slog.Int("int", testInt),
					slog.Bool("bool", testBool),
					slog.String("string", testString),
				)
			},
		},
	} {
		b.Run(call.name, func(b *testing.B) {
			b.Run("Simple", func(b *testing.B) {
				provider := NewLoggerProvider(WithProcessor(NewSimpleProcessor(noopExporter{})))
				logger := slog.New(&slogHandler{provider.Logger("log/slog")})

				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					call.f(logger)
				}
				_ = provider.Shutdown(context.Background())
			})
			b.Run("Batch", func(b *testing.B) {
				provider := NewLoggerProvider(WithProcessor(NewBatchingProcessor(noopExporter{})))
				logger := slog.New(&slogHandler{provider.Logger("log/slog")})

				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					call.f(logger)
				}
				_ = provider.Shutdown(context.Background())
			})
		})
	}
}

type slogHandler struct {
	Logger log.Logger
}

// Handle handles the Record.
// It should avoid memory allocations whenever possible.
func (h *slogHandler) Handle(ctx context.Context, r slog.Record) error {
	record := log.Record{}

	record.SetTimestamp(r.Time)

	record.SetBody(log.StringValue(r.Message))

	lvl := convertLevel(r.Level)
	record.SetSeverity(lvl)

	r.Attrs(func(a slog.Attr) bool {
		record.AddAttributes(convertAttr(a))
		return true
	})

	h.Logger.Emit(ctx, record)
	return nil
}

// Enabled is implemented as a dummy.
func (h *slogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// WithAttrs is implemented as a dummy.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup is implemented as a dummy.
func (h *slogHandler) WithGroup(name string) slog.Handler {
	return h
}

func convertLevel(l slog.Level) log.Severity {
	return log.Severity(l + 9)
}

func convertAttr(attr slog.Attr) log.KeyValue {
	val := convertValue(attr.Value)
	return log.KeyValue{Key: attr.Key, Value: val}
}

func convertValue(v slog.Value) log.Value {
	switch v.Kind() {
	case slog.KindAny:
		return log.StringValue(fmt.Sprintf("%+v", v.Any()))
	case slog.KindBool:
		return log.BoolValue(v.Bool())
	case slog.KindDuration:
		return log.Int64Value(v.Duration().Nanoseconds())
	case slog.KindFloat64:
		return log.Float64Value(v.Float64())
	case slog.KindInt64:
		return log.Int64Value(v.Int64())
	case slog.KindString:
		return log.StringValue(v.String())
	case slog.KindTime:
		return log.Int64Value(v.Time().UnixNano())
	case slog.KindUint64:
		return log.Int64Value(int64(v.Uint64()))
	default:
		panic(fmt.Sprintf("unhandled attribute kind: %s", v.Kind()))
	}
}

type noopExporter struct{}

func (e noopExporter) Export(_ context.Context, _ []Record) error {
	return nil
}

func (e noopExporter) Shutdown(_ context.Context) error {
	return nil
}

func (e noopExporter) ForceFlush(_ context.Context) error {
	return nil
}

type timestampDecorator struct {
	Processor
}

func (e timestampDecorator) OnEmit(ctx context.Context, r Record) error {
	r = r.Clone()
	r.SetObservedTimestamp(testTimestamp)
	return e.Processor.OnEmit(ctx, r)
}
