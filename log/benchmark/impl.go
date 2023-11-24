// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package benchmark // import "go.opentelemetry.io/otel/log/benchmark"

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

// writerLogger is a logger that writes to a provided io.Writer without any locking.
// It is intended to represent a high-performance logger that synchronously
// writes text.
type writerLogger struct {
	embedded.Logger
	w io.Writer
}

func (l *writerLogger) Emit(_ context.Context, r log.Record) {
	if !r.Timestamp.IsZero() {
		l.write("timestamp=")
		l.write(strconv.FormatInt(r.Timestamp.Unix(), 10))
		l.write(" ")
	}
	l.write("severity=")
	l.write(strconv.FormatInt(int64(r.Severity), 10))
	l.write(" ")
	l.write("body=")
	l.write(r.Body)
	r.Attributes(func(kv attribute.KeyValue) bool {
		l.write(" ")
		l.write(string(kv.Key))
		l.write("=")
		l.appendValue(kv.Value)
		return true
	})
	l.write("\n")
}

func (l *writerLogger) appendValue(v attribute.Value) {
	switch v.Type() {
	case attribute.STRING:
		l.write(v.AsString())
	case attribute.INT64:
		l.write(strconv.FormatInt(v.AsInt64(), 10)) // strconv.FormatInt allocates memory.
	case attribute.FLOAT64:
		l.write(strconv.FormatFloat(v.AsFloat64(), 'g', -1, 64)) // strconv.FormatFloat allocates memory.
	case attribute.BOOL:
		l.write(strconv.FormatBool(v.AsBool()))
	default:
		panic(fmt.Sprintf("unhandled attribute type: %s", v.Type()))
	}
}

func (l *writerLogger) write(s string) {
	_, _ = io.WriteString(l.w, s)
}
