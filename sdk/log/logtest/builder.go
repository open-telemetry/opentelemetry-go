// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/sdk/log/logtest"

import "go.opentelemetry.io/otel/sdk/log"

// RecordBuilder creates records for testing purposes.
// [RecordBuilder] can be used for testing
// [go.opentelemetry.io/otel/sdk/log.Exporter]
// and [go.opentelemetry.io/otel/sdk/log.Processor]
// implementations.
type RecordBuilder struct {
	record *record
}

// Record builds a new record.
func (b *RecordBuilder) Record() log.Record {
	if b.record == nil {
		b.record = new(record)
	}
	r := log.Record(b.record.Clone())
	return &r
}
