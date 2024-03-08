// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/trace"
)

var recordsPool = sync.Pool{
	New: func() any {
		b := make([]*Record, 1)
		return &b
	},
}

// Compile-time check logger implements metric.log.Logger.
var _ log.Logger = (*logger)(nil)

type logger struct {
	embedded.Logger

	provider *LoggerProvider
	scope    instrumentation.Scope
}

func (l *logger) Emit(ctx context.Context, r log.Record) {
	record := &Record{
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

	if span := trace.SpanContextFromContext(ctx); span.IsValid() {
		record.traceID = span.TraceID()
		record.spanID = span.SpanID()
		record.traceFlags = span.TraceFlags()
	}

	r.WalkAttributes(func(kv log.KeyValue) bool {
		record.AddAttributes(kv)
		return true
	})

	records := recordsPool.Get().(*[]*Record)
	(*records)[0] = record
	for _, expoter := range l.provider.cfg.exporters {
		expoter.Export(ctx, *records)
	}
	recordsPool.Put(records)
}
