// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package benchmark

import (
	"fmt"
	"sync"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

func TestLogrSink(t *testing.T) {
	spy := &spyLogger{}

	l := logr.New(&logrSink{spy})

	l.Info(testBody, "string", testString)

	want := log.Record{
		Body:     testBody,
		Severity: log.SeverityInfo,
		Attributes: []attribute.KeyValue{
			attribute.String("string", testString),
		},
	}

	assert.Equal(t, want, spy.Record)
}

type logrSink struct {
	Logger log.Logger
}

var logrAttrPool = sync.Pool{
	New: func() interface{} {
		attr := make([]attribute.KeyValue, 0, 5)
		return &attr
	},
}

// Init is implementated as a dummy.
func (s *logrSink) Init(info logr.RuntimeInfo) {
}

// Enabled is implementated as a dummy.
func (s *logrSink) Enabled(level int) bool {
	return true
}

// Info logs a non-error message with the given key/value pairs as context.
// It should avoid memory allocations whenever possible.
func (s *logrSink) Info(level int, msg string, keysAndValues ...any) {
	record := log.Record{}

	record.Body = msg

	lvl := log.Severity(9 - level)
	record.Severity = lvl

	if len(keysAndValues)%2 == 1 {
		panic("key without a value")
	}
	kvCount := len(keysAndValues) / 2
	ptr := logrAttrPool.Get().(*[]attribute.KeyValue)
	attrs := *ptr
	defer func() {
		*ptr = attrs[:0]
		logrAttrPool.Put(ptr)
	}()
	for i := 0; i < kvCount; i++ {
		k, ok := keysAndValues[i*2].(string)
		if !ok {
			panic("key is not a string")
		}
		kv := convertKV(k, keysAndValues[i*2+1])
		attrs = append(attrs, kv)
	}
	record.Attributes = attrs

	s.Logger.Emit(ctx, record)
}

// Error is implementated as a dummy.
func (s *logrSink) Error(err error, msg string, keysAndValues ...any) {
}

// WithValues is implementated as a dummy.
func (s *logrSink) WithValues(keysAndValues ...any) logr.LogSink {
	return s
}

// WithName is implementated as a dummy.
func (s *logrSink) WithName(name string) logr.LogSink {
	return s
}

func convertKV(k string, v interface{}) attribute.KeyValue {
	switch val := v.(type) {
	case bool:
		return attribute.Bool(k, val)
	case float64:
		return attribute.Float64(k, val)
	case int:
		return attribute.Int(k, val)
	case string:
		return attribute.String(k, val)
	default:
		panic(fmt.Sprintf("unhandled value type: %T", val))
	}
}
