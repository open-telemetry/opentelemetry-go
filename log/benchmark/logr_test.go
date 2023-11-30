// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package benchmark

import (
	"fmt"
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

	assert.Equal(t, testBody, spy.Record.Body)
	assert.Equal(t, log.SeverityInfo, spy.Record.Severity)
	assert.Equal(t, 1, spy.Record.AttributesLen())
	spy.Record.WalkAttributes(func(kv attribute.KeyValue) bool {
		assert.Equal(t, "string", string(kv.Key))
		assert.Equal(t, testString, kv.Value.AsString())
		return true
	})
}

type logrSink struct {
	Logger log.Logger
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
	lvl := log.Severity(9 - level)

	record := log.Record{Severity: lvl, Body: msg}

	if len(keysAndValues)%2 == 1 {
		panic("key without a value")
	}
	kvCount := len(keysAndValues) / 2
	if kvCount > log.AttributesInlineCount {
		attrs := make([]attribute.KeyValue, 0, kvCount)
		for i := 0; i < kvCount; i++ {
			k, ok := keysAndValues[i*2].(string)
			if !ok {
				panic("key is not a string")
			}
			kv := convertKV(k, keysAndValues[i*2+1])
			attrs = append(attrs, kv)
		}
		record.AddAttributes(attrs...)
	} else {
		// special case that avoids heap allocations (hot path)
		for i := 0; i < kvCount; i++ {
			k, ok := keysAndValues[i*2].(string)
			if !ok {
				panic("key is not a string")
			}
			kv := convertKV(k, keysAndValues[i*2+1])
			record.AddAttributes(kv)
		}
	}

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
