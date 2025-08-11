// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/log/internal/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

// Compile-time check SimpleProcessor implements Processor.
var _ Processor = (*SimpleProcessor)(nil)

// simpleProcessorIDCounter is used to generate unique component names.
var simpleProcessorIDCounter atomic.Uint64

// SimpleProcessor is an processor that synchronously exports log records.
//
// Use [NewSimpleProcessor] to create a SimpleProcessor.
type SimpleProcessor struct {
	mu       sync.Mutex
	exporter Exporter

	selfObservabilityEnabled bool
	processedMetric          otelconv.SDKProcessorLogProcessed
	componentName            string

	noCmp [0]func() //nolint: unused  // This is indeed used.
}

// NewSimpleProcessor is a simple Processor adapter.
//
// This Processor is not recommended for production use due to its synchronous
// nature, which makes it suitable for testing, debugging, or demonstrating
// other features, but can lead to slow performance and high computational
// overhead. For production environments, it is recommended to use
// [NewBatchProcessor] instead. However, there may be exceptions where certain
// [Exporter] implementations perform better with this Processor.
func NewSimpleProcessor(exporter Exporter, _ ...SimpleProcessorOption) *SimpleProcessor {
	s := &SimpleProcessor{
		exporter: exporter,
		componentName: fmt.Sprintf(
			"%s/%d",
			string(otelconv.ComponentTypeSimpleLogProcessor),
			simpleProcessorIDCounter.Add(1)-1,
		),
	}
	s.initSelfObservability()
	return s
}

func (s *SimpleProcessor) initSelfObservability() {
	if !x.SelfObservability.Enabled() {
		return
	}

	s.selfObservabilityEnabled = true
	mp := otel.GetMeterProvider()
	m := mp.Meter("go.opentelemetry.io/otel/sdk/log",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL))

	var err error
	if s.processedMetric, err = otelconv.NewSDKProcessorLogProcessed(m); err != nil {
		otel.Handle(err)
	}
}

var simpleProcRecordsPool = sync.Pool{
	New: func() any {
		records := make([]Record, 1)
		return &records
	},
}

// OnEmit batches provided log record.
func (s *SimpleProcessor) OnEmit(ctx context.Context, r *Record) (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	records := simpleProcRecordsPool.Get().(*[]Record)
	(*records)[0] = *r
	defer func() {
		simpleProcRecordsPool.Put(records)
	}()

	defer func() {
		if s.selfObservabilityEnabled {
			attrs := make([]attribute.KeyValue, 2, 3)
			attrs[0] = s.processedMetric.AttrComponentType(otelconv.ComponentTypeSimpleLogProcessor)
			attrs[1] = s.processedMetric.AttrComponentName(s.componentName)
			if err != nil {
				attrs = append(attrs, s.processedMetric.AttrErrorType(otelconv.ErrorTypeOther))
			}
			s.processedMetric.Add(context.Background(), int64(len(*records)), attrs...)
		}
	}()

	if s.exporter == nil {
		return nil
	}

	return s.exporter.Export(ctx, *records)
}

// Shutdown shuts down the exporter.
func (s *SimpleProcessor) Shutdown(ctx context.Context) error {
	if s.exporter == nil {
		return nil
	}

	return s.exporter.Shutdown(ctx)
}

// ForceFlush flushes the exporter.
func (s *SimpleProcessor) ForceFlush(ctx context.Context) error {
	if s.exporter == nil {
		return nil
	}

	return s.exporter.ForceFlush(ctx)
}

// SimpleProcessorOption applies a configuration to a [SimpleProcessor].
type SimpleProcessorOption interface {
	apply()
}

// ResetSimpleProcessorIDCounterForTesting resets the global ID counter for testing purposes.
func ResetSimpleProcessorIDCounterForTesting() {
	simpleProcessorIDCounter.Store(0)
}
