// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package benchmark

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

func TestWriterLogger(t *testing.T) {
	sb := &strings.Builder{}
	l := &writerLogger{w: sb}

	r := log.Record{
		Timestamp: testTimestamp,
		Severity:  testSeverity,
		Body:      testBody,
		Attributes: []attribute.KeyValue{
			attribute.String("string", testString),
			attribute.Float64("float", testFloat),
			attribute.Int("int", testInt),
			attribute.Bool("bool", testBool),
		},
	}
	l.Emit(ctx, r)

	want := "timestamp=595728000 severity=9 body=log message string=7e3b3b2aaeff56a7108fe11e154200dd/7819479873059528190 float=1.2345 int=32768 bool=true\n"
	assert.Equal(t, want, sb.String())
}

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
	for _, kv := range r.Attributes {
		l.write(" ")
		l.write(string(kv.Key))
		l.write("=")
		l.appendValue(kv.Value)
	}
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
