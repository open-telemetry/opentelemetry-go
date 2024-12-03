// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"
	"fmt"
	"strings"
	"sync"

	logapi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

// Initialize OpenTelemetry Logs SDK and setup logging using a log bridge.
func Example() {
	// Create an exporter that will emit log records.
	// E.g. use go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp
	// to send logs using OTLP over HTTP:
	// exporter, err := otlploghttp.New(ctx)
	var exporter log.Exporter

	// Create a log record processor pipeline.
	processor := log.NewBatchProcessor(exporter)

	// Create a logger provider.
	// You can pass this instance directly when creating a log bridge.
	provider := log.NewLoggerProvider(
		log.WithProcessor(processor),
	)

	// Handle shutdown properly so that nothing leaks.
	defer func() {
		err := provider.Shutdown(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}()

	// Register as global logger provider so that it can be used via global.Meter
	// and accessed using global.GetMeterProvider.
	// Most log bridges use the global logger provider as default.
	// If the global logger provider is not set then a no-op implementation
	// is used, which fails to generate data.
	global.SetLoggerProvider(provider)

	// Use a bridge so that you can emit logs using your Go logging library of preference.
	// E.g. use go.opentelemetry.io/contrib/bridges/otelslog so that you can use log/slog:
	// slog.SetDefault(otelslog.NewLogger("my/pkg/name", otelslog.WithLoggerProvider(provider)))
}

// Use a processor that filters out records based on the provided context.
func ExampleProcessor_filtering() {
	// Existing processor that emits telemetry.
	var processor log.Processor = log.NewBatchProcessor(nil)

	// Wrap the processor so that it ignores processing log records
	// when a context deriving from WithIgnoreLogs is passed
	// to the logging methods.
	processor = &ContextFilterProcessor{Processor: processor}

	// The created processor can then be registered with
	// the OpenTelemetry Logs SDK using the WithProcessor option.
	_ = log.NewLoggerProvider(
		log.WithProcessor(processor),
	)
}

type key struct{}

var ignoreLogsKey key

// WithIgnoreLogs returns a context which is used by [ContextFilterProcessor]
// to filter out log records.
func WithIgnoreLogs(ctx context.Context) context.Context {
	return context.WithValue(ctx, ignoreLogsKey, true)
}

// ContextFilterProcessor filters out logs when a context deriving from
// [WithIgnoreLogs] is passed to its methods.
type ContextFilterProcessor struct {
	log.Processor

	lazyFilter sync.Once
	// Use the experimental FilterProcessor interface
	// (go.opentelemetry.io/otel/sdk/log/internal/x).
	filter filter
}

type filter interface {
	Enabled(ctx context.Context, param logapi.EnabledParameters) bool
}

func (p *ContextFilterProcessor) OnEmit(ctx context.Context, record *log.Record) error {
	if ignoreLogs(ctx) {
		return nil
	}
	return p.Processor.OnEmit(ctx, record)
}

func (p *ContextFilterProcessor) Enabled(ctx context.Context, param logapi.EnabledParameters) bool {
	p.lazyFilter.Do(func() {
		if f, ok := p.Processor.(filter); ok {
			p.filter = f
		}
	})
	return !ignoreLogs(ctx) && (p.filter == nil || p.filter.Enabled(ctx, param))
}

func ignoreLogs(ctx context.Context) bool {
	_, ok := ctx.Value(ignoreLogsKey).(bool)
	return ok
}

// Use a processor which redacts sensitive data from some attributes.
func ExampleProcessor_redact() {
	// Existing processor that emits telemetry.
	var processor log.Processor = log.NewBatchProcessor(nil)

	// Add a processor so that it redacts values from token attributes.
	redactProcessor := &RedactTokensProcessor{}

	// The created processor can then be registered with
	// the OpenTelemetry Logs SDK using the WithProcessor option.
	_ = log.NewLoggerProvider(
		// Order is important here. Redact before handing to the processor.
		log.WithProcessor(redactProcessor),
		log.WithProcessor(processor),
	)
}

// RedactTokensProcessor is a [log.Processor] decorator that redacts values
// from attributes containing "token" in the key.
type RedactTokensProcessor struct{}

// OnEmit redacts values from attributes containing "token" in the key
// by replacing them with a REDACTED value.
func (p *RedactTokensProcessor) OnEmit(ctx context.Context, record *log.Record) error {
	record.WalkAttributes(func(kv logapi.KeyValue) bool {
		if strings.Contains(strings.ToLower(kv.Key), "token") {
			record.AddAttributes(logapi.String(kv.Key, "REDACTED"))
		}
		return true
	})
	return nil
}

// Shutdown returns nil.
func (p *RedactTokensProcessor) Shutdown(ctx context.Context) error {
	return nil
}

// ForceFlush returns nil.
func (p *RedactTokensProcessor) ForceFlush(ctx context.Context) error {
	return nil
}
