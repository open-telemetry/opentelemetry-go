// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package benchmark

import (
	"context"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestSlogHandler(t *testing.T) {
	spy := &spyLogger{}
	l := slog.New(&slogHandler{spy})

	l.Info(testBody, "string", testString)

	assert.Equal(t, testBody, spy.Record.Body)
	assert.Equal(t, log.SeverityInfo, spy.Record.Severity)
	assert.Equal(t, 1, spy.Record.AttributesLen())
	spy.Record.Attributes(func(kv attribute.KeyValue) bool {
		assert.Equal(t, "string", string(kv.Key))
		assert.Equal(t, testString, kv.Value.AsString())
		return true
	})
}

type slogHandler struct {
	Logger log.Logger
}

// Handle handles the Record.
// It should avoid memory allocations whenever possible.
func (h *slogHandler) Handle(_ context.Context, r slog.Record) error {
	lvl := convertLevel(r.Level)

	record := log.Record{Timestamp: r.Time, Severity: lvl, Body: r.Message}

	if r.NumAttrs() > log.AttributesInlineCount {
		attrs := make([]attribute.KeyValue, 0, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, convertAttr(a))
			return true
		})
		record.AddAttributes(attrs...)
	} else {
		// special case that avoids heap allocations (hot path)
		r.Attrs(func(a slog.Attr) bool {
			record.AddAttributes(convertAttr(a))
			return true
		})
	}

	h.Logger.Emit(context.Background(), record)
	return nil
}

// Enabled is implementated as a dummy.
func (h *slogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

// WithAttrs is implementated as a dummy.
func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup is implementated as a dummy.
func (h *slogHandler) WithGroup(name string) slog.Handler {
	return h
}

func convertLevel(l slog.Level) log.Severity {
	return log.Severity(l + 9)
}

func convertAttr(attr slog.Attr) attribute.KeyValue {
	val := convertValue(attr.Value)
	return attribute.KeyValue{Key: attribute.Key(attr.Key), Value: val}
}

func convertValue(v slog.Value) attribute.Value {
	switch v.Kind() {
	case slog.KindAny:
		return attribute.StringValue(fmt.Sprintf("%+v", v.Any()))
	case slog.KindBool:
		return attribute.BoolValue(v.Bool())
	case slog.KindDuration:
		return attribute.Int64Value(v.Duration().Nanoseconds())
	case slog.KindFloat64:
		return attribute.Float64Value(v.Float64())
	case slog.KindInt64:
		return attribute.Int64Value(v.Int64())
	case slog.KindString:
		return attribute.StringValue(v.String())
	case slog.KindTime:
		return attribute.Int64Value(v.Time().UnixNano())
	case slog.KindUint64:
		return attribute.Int64Value(int64(v.Uint64()))
	default:
		panic(fmt.Sprintf("unhandled attribute kind: %s", v.Kind()))
	}
}
