// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log_test

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel/log"
	logsdk "go.opentelemetry.io/otel/sdk/log"
)

// Use a processor that filters out records based on the provided context.
func ExampleProcessor_filtering() {
	// Existing processor that emits telemetry.
	var processor logsdk.Processor = logsdk.NewBatchProcessor(nil)

	// Wrap the processor so that it ignores processing log records
	// when a context deriving from WithIgnoreLogs is passed
	// to the logging methods.
	processor = &ContextFilterProcessor{processor}

	// The created processor can then be registered with
	// the OpenTelemetry Logs SDK using the WithProcessor option.
	_ = logsdk.NewLoggerProvider(
		logsdk.WithProcessor(processor),
	)
}

type key struct{}

var igoreLogsKey key

// WithIgnoreLogs returns a context which is used by [ContextFilterProcessor]
// to filter out log records.
func WithIgnoreLogs(ctx context.Context) context.Context {
	return context.WithValue(ctx, igoreLogsKey, true)
}

// ContextFilterProcessor filters out logs when a context deriving from
// [WithIgnoreLogs] is passed to its methods.
type ContextFilterProcessor struct {
	logsdk.Processor
}

func (p *ContextFilterProcessor) OnEmit(ctx context.Context, record logsdk.Record) error {
	if ignoreLogs(ctx) {
		return nil
	}
	return p.Processor.OnEmit(ctx, record)
}

func (p *ContextFilterProcessor) Enabled(ctx context.Context, record logsdk.Record) bool {
	return !ignoreLogs(ctx) && p.Processor.Enabled(ctx, record)
}

func ignoreLogs(ctx context.Context) bool {
	_, ok := ctx.Value(igoreLogsKey).(bool)
	return ok
}

// Use a processor which redacts sensitive data from some attributes.
func ExampleProcessor_redact() {
	// Existing processor that emits telemetry.
	var processor logsdk.Processor = logsdk.NewBatchProcessor(nil)

	// Wrap the processor so that it redacts values from token attributes.
	processor = &RedactTokensProcessor{processor}

	// The created processor can then be registered with
	// the OpenTelemetry Logs SDK using the WithProcessor option.
	_ = logsdk.NewLoggerProvider(
		logsdk.WithProcessor(processor),
	)
}

// RedactTokensProcessor is a [logsdk.Processor] decorator that redacts values
// from attributes containing "token" in the key.
type RedactTokensProcessor struct {
	logsdk.Processor
}

// OnEmit redacts values from attributes containing "token" in the key
// by replacing them with a REDACTED value.
func (p *RedactTokensProcessor) OnEmit(ctx context.Context, record logsdk.Record) error {
	cloned := false
	record.WalkAttributes(func(kv log.KeyValue) bool {
		if strings.Contains(strings.ToLower(kv.Key), "token") {
			if !cloned {
				record = record.Clone()
				cloned = true
			}
			record.AddAttributes(log.String(kv.Key, "REDACTED"))
		}
		return true
	})
	return p.Processor.OnEmit(ctx, record)
}
