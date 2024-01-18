// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/log/internal"

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
)

type logrSink struct {
	Logger log.Logger
}

// Init is implemented as a dummy.
func (s *logrSink) Init(info logr.RuntimeInfo) {
}

// Enabled is implemented as a dummy.
func (s *logrSink) Enabled(level int) bool {
	return true
}

// Info logs a non-error message with the given key/value pairs as context.
// It should avoid memory allocations whenever possible.
func (s *logrSink) Info(level int, msg string, keysAndValues ...any) {
	record := log.Record{}

	record.SetBody(msg)

	lvl := log.Severity(9 - level)
	record.SetSeverity(lvl)

	if len(keysAndValues)%2 == 1 {
		panic("key without a value")
	}
	kvCount := len(keysAndValues) / 2
	ctx := context.Background()
	for i := 0; i < kvCount; i++ {
		k, ok := keysAndValues[i*2].(string)
		if !ok {
			panic("key is not a string")
		}
		v := keysAndValues[i*2+1]
		if vCtx, ok := v.(context.Context); ok {
			// Special case when a field is of context.Context type.
			ctx = vCtx
			continue
		}
		kv := convertKV(k, v)
		record.AddAttributes(kv)
	}

	s.Logger.Emit(ctx, record)
}

// Error is implemented as a dummy.
func (s *logrSink) Error(err error, msg string, keysAndValues ...any) {
}

// WithValues is implemented as a dummy.
func (s *logrSink) WithValues(keysAndValues ...any) logr.LogSink {
	return s
}

// WithName is implemented as a dummy.
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
