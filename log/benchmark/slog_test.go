// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package benchmark

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slog"
)

func TestSlogHandler(t *testing.T) {
	spy := &spyLogger{}
	l := slog.New(&slogHandler{spy})

	l.Info(testBody, "string", testString)

	want := log.Record{
		Body:     testBody,
		Severity: log.SeverityInfo,
		Attributes: []attribute.KeyValue{
			attribute.String("string", testString),
		},
	}

	assert.NotZero(t, spy.Record.Timestamp, "should set a timestamp")
	spy.Record.Timestamp = time.Time{}
	assert.Equal(t, want, spy.Record)
}

type slogHandler struct {
	Logger log.Logger
}

var slogAttrPool = sync.Pool{
	New: func() interface{} {
		attr := make([]attribute.KeyValue, 0, 5)
		return &attr
	},
}

// Handle handles the Record.
// It should avoid memory allocations whenever possible.
func (h *slogHandler) Handle(_ context.Context, r slog.Record) error {
	record := log.Record{}

	record.Timestamp = r.Time

	record.Body = r.Message

	lvl := convertLevel(r.Level)
	record.Severity = lvl

	ptr := slogAttrPool.Get().(*[]attribute.KeyValue)
	attrs := *ptr
	defer func() {
		*ptr = attrs[:0]
		slogAttrPool.Put(ptr)
	}()
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, convertAttr(a))
		return true
	})
	record.Attributes = attrs

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
