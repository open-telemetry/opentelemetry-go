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

package otlptest

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	exportmetric "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/sum"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Used to avoid implementing locking functions for test
// checkpointsets.
type noopLocker struct{}

// Lock implements sync.Locker, which is needed for
// exportmetric.CheckpointSet.
func (noopLocker) Lock() {}

// Unlock implements sync.Locker, which is needed for
// exportmetric.CheckpointSet.
func (noopLocker) Unlock() {}

// RLock implements exportmetric.CheckpointSet.
func (noopLocker) RLock() {}

// RUnlock implements exportmetric.CheckpointSet.
func (noopLocker) RUnlock() {}

// OneRecordCheckpointSet is a CheckpointSet that returns just one
// filled record. It may be useful for testing driver's metrics
// export.
type OneRecordCheckpointSet struct {
	noopLocker
}

var _ exportmetric.CheckpointSet = OneRecordCheckpointSet{}

// ForEach implements exportmetric.CheckpointSet. It always invokes
// the callback once with always the same record.
func (OneRecordCheckpointSet) ForEach(kindSelector exportmetric.ExportKindSelector, recordFunc func(exportmetric.Record) error) error {
	desc := metric.NewDescriptor(
		"foo",
		metric.CounterInstrumentKind,
		number.Int64Kind,
	)
	res := resource.NewWithAttributes(attribute.String("a", "b"))
	agg := sum.New(1)
	if err := agg[0].Update(context.Background(), number.NewInt64Number(42), &desc); err != nil {
		return err
	}
	start := time.Date(2020, time.December, 8, 19, 15, 0, 0, time.UTC)
	end := time.Date(2020, time.December, 8, 19, 16, 0, 0, time.UTC)
	labels := attribute.NewSet(attribute.String("abc", "def"), attribute.Int64("one", 1))
	rec := exportmetric.NewRecord(&desc, &labels, res, agg[0].Aggregation(), start, end)
	return recordFunc(rec)
}

// SingleSpanSnapshot returns a one-element slice with a snapshot. It
// may be useful for testing driver's trace export.
func SingleSpanSnapshot() []*tracesdk.SpanSnapshot {
	sd := &tracesdk.SpanSnapshot{
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    trace.TraceID{2, 3, 4, 5, 6, 7, 8, 9, 2, 3, 4, 5, 6, 7, 8, 9},
			SpanID:     trace.SpanID{3, 4, 5, 6, 7, 8, 9, 0},
			TraceFlags: trace.FlagsSampled,
		}),
		Parent: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    trace.TraceID{2, 3, 4, 5, 6, 7, 8, 9, 2, 3, 4, 5, 6, 7, 8, 9},
			SpanID:     trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
			TraceFlags: trace.FlagsSampled,
		}),
		SpanKind:                 trace.SpanKindInternal,
		Name:                     "foo",
		StartTime:                time.Date(2020, time.December, 8, 20, 23, 0, 0, time.UTC),
		EndTime:                  time.Date(2020, time.December, 0, 20, 24, 0, 0, time.UTC),
		Attributes:               []attribute.KeyValue{},
		MessageEvents:            []trace.Event{},
		Links:                    []trace.Link{},
		StatusCode:               codes.Ok,
		StatusMessage:            "",
		DroppedAttributeCount:    0,
		DroppedMessageEventCount: 0,
		DroppedLinkCount:         0,
		ChildSpanCount:           0,
		Resource:                 resource.NewWithAttributes(attribute.String("a", "b")),
		InstrumentationLibrary: instrumentation.Library{
			Name:    "bar",
			Version: "0.0.0",
		},
	}
	return []*tracesdk.SpanSnapshot{sd}
}

// EmptyCheckpointSet is a checkpointer that has no records at all.
type EmptyCheckpointSet struct {
	noopLocker
}

var _ exportmetric.CheckpointSet = EmptyCheckpointSet{}

// ForEach implements exportmetric.CheckpointSet. It never invokes the
// callback.
func (EmptyCheckpointSet) ForEach(kindSelector exportmetric.ExportKindSelector, recordFunc func(exportmetric.Record) error) error {
	return nil
}

// FailCheckpointSet is a checkpointer that returns an error during
// ForEach.
type FailCheckpointSet struct {
	noopLocker
}

var _ exportmetric.CheckpointSet = FailCheckpointSet{}

// ForEach implements exportmetric.CheckpointSet. It always fails.
func (FailCheckpointSet) ForEach(kindSelector exportmetric.ExportKindSelector, recordFunc func(exportmetric.Record) error) error {
	return fmt.Errorf("fail")
}
