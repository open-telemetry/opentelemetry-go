// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/log/internal"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/log"

	"golang.org/x/exp/slog"
)

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
