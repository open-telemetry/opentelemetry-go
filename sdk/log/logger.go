// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/trace"
)

// Compile-time check logger implements metric.log.Logger.
var _ log.Logger = (*logger)(nil)

var keysPool = sync.Pool{
	New: func() any {
		keys := map[string]struct{}{}
		return &keys
	},
}

type logger struct {
	embedded.Logger

	provider *LoggerProvider
	scope    instrumentation.Scope
}

func (l *logger) Emit(ctx context.Context, r log.Record) {
	record := Record{ // This always escapes to the heap.
		resource:                  l.provider.cfg.resource,
		attributeCountLimit:       l.provider.cfg.attributeCountLimit,
		attributeValueLengthLimit: l.provider.cfg.attributeValueLengthLimit,

		scope: &l.scope,

		timestamp:         r.Timestamp(),
		observedTimestamp: r.ObservedTimestamp(),
		severity:          r.Severity(),
		severityText:      r.SeverityText(),
		body:              r.Body(),
	}

	if record.observedTimestamp.Equal(time.Time{}) {
		record.observedTimestamp = time.Now()
	}

	if span := trace.SpanContextFromContext(ctx); span.IsValid() { // This escapes to the heap if there is no span in context.
		record.traceID = span.TraceID()
		record.spanID = span.SpanID()
		record.traceFlags = span.TraceFlags()
	}

	keysPtr := keysPool.Get().(*map[string]struct{})
	defer func() {
		clear(*keysPtr)
		keysPool.Put(keysPtr)
	}()
	keys := *keysPtr
	r.WalkAttributes(func(kv log.KeyValue) bool {
		if _, ok := keys[kv.Key]; ok {
			// Drop duplicated log attribute.
			return true
		}

		keys[kv.Key] = struct{}{}

		if hasDuplicatedKey(kv.Value) {
			// Drop attribute which contains any key-value with duplicated key.
			return true
		}
		record.AddAttributes(kv)
		return true
	})

	for _, processor := range l.provider.cfg.processors {
		if err := processor.OnEmit(ctx, record); err != nil {
			otel.Handle(err)
		}
	}
}

func hasDuplicatedKey(value log.Value) bool {
	if value.Kind() != log.KindMap {
		return false
	}

	keysPtr := keysPool.Get().(*map[string]struct{})
	defer func() {
		clear(*keysPtr)
		keysPool.Put(keysPtr)
	}()
	keys := *keysPtr
	m := value.AsMap()
	for _, kv := range m {
		if _, ok := keys[kv.Key]; ok {
			return true
		}
		keys[kv.Key] = struct{}{}

		if hasDuplicatedKey(kv.Value) {
			return true
		}
	}

	return false
}
