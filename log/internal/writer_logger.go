// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/log/internal"

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/trace"
)

// writerLogger is a logger that writes to a provided io.Writer without any locking.
// It is intended to represent a high-performance logger that synchronously
// writes text.
type writerLogger struct {
	embedded.Logger
	w io.Writer
}

func (l *writerLogger) Emit(ctx context.Context, r log.Record) {
	if !r.Timestamp().IsZero() {
		l.write("timestamp=")
		l.write(strconv.FormatInt(r.Timestamp().Unix(), 10))
		l.write(" ")
	}
	l.write("severity=")
	l.write(strconv.FormatInt(int64(r.Severity()), 10))
	l.write(" ")

	if !r.Body().Empty() {
		l.write("body=")
		l.appendValue(r.Body())
	}

	r.WalkAttributes(func(kv log.KeyValue) bool {
		l.write(" ")
		l.write(string(kv.Key))
		l.write("=")
		l.appendValue(kv.Value)
		return true
	})

	span := trace.SpanContextFromContext(ctx)
	if span.IsValid() {
		l.write(" traced=true")
	}

	l.write("\n")
}

func (l *writerLogger) appendValue(v log.Value) {
	switch v.Kind() {
	case log.KindString:
		l.write(v.String())
	case log.KindInt64:
		l.write(strconv.FormatInt(v.Int64(), 10)) // strconv.FormatInt allocates memory.
	case log.KindFloat64:
		l.write(strconv.FormatFloat(v.Float64(), 'g', -1, 64)) // strconv.FormatFloat allocates memory.
	case log.KindBool:
		l.write(strconv.FormatBool(v.Bool()))
	case log.KindBytes:
		l.write(fmt.Sprint(v.Bytes()))
	case log.KindMap:
		l.write(fmt.Sprint(v.Map()))
	case log.KindEmpty:
		l.write("<nil>")
	default:
		panic(fmt.Sprintf("unhandled value kind: %s", v.Kind()))
	}
}

func (l *writerLogger) write(s string) {
	_, _ = io.WriteString(l.w, s)
}
