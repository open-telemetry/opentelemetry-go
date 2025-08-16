// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/trace/internal/x"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
	"go.opentelemetry.io/otel/trace"
)

var measureAttrsPool = sync.Pool{
	New: func() any {
		// "component.name" + "component.type" + "error.type"
		const n = 1 + 1 + 1
		s := make([]attribute.KeyValue, 0, n)
		// Return a pointer to a slice instead of a slice itself
		// to avoid allocations on every call.
		return &s
	},
}

// simpleSpanProcessor is a SpanProcessor that synchronously sends all
// completed Spans to a trace.Exporter immediately.
type simpleSpanProcessor struct {
	exporterMu sync.Mutex
	exporter   SpanExporter
	stopOnce   sync.Once

	selfObservabilityEnabled bool
	componentNameAttr        attribute.KeyValue
	spansProcessedCounter    otelconv.SDKProcessorSpanProcessed
}

var _ SpanProcessor = (*simpleSpanProcessor)(nil)

// NewSimpleSpanProcessor returns a new SpanProcessor that will synchronously
// send completed spans to the exporter immediately.
//
// This SpanProcessor is not recommended for production use. The synchronous
// nature of this SpanProcessor makes it good for testing, debugging, or showing
// examples of other features, but it will be slow and have a high computation
// resource usage overhead. The BatchSpanProcessor is recommended for production
// use instead.
func NewSimpleSpanProcessor(exporter SpanExporter) SpanProcessor {
	ssp := &simpleSpanProcessor{
		exporter:                 exporter,
		selfObservabilityEnabled: x.SelfObservability.Enabled(),
	}

	if ssp.selfObservabilityEnabled {
		ssp.componentNameAttr = semconv.OTelComponentName(
			fmt.Sprintf("%s/%d", otelconv.ComponentTypeSimpleSpanProcessor, nextSimpleProcessorID()))

		var err error
		ssp.spansProcessedCounter, err = newInst()
		if err != nil {
			msg := "failed to create self-observability metrics for simple span processor: %w"
			err := fmt.Errorf(msg, err)
			otel.Handle(err)
		}
	}

	global.Warn("SimpleSpanProcessor is not recommended for production use, consider using BatchSpanProcessor instead.")

	return ssp
}

var simpleProcessorIDCounter atomic.Int64

// nextSimpleProcessorID returns an identifier for this simple span processor,
// starting with 0 and incrementing by 1 each time it is called.
func nextSimpleProcessorID() int64 {
	return simpleProcessorIDCounter.Add(1) - 1
}

func newInst() (otelconv.SDKProcessorSpanProcessed, error) {
	meter := otel.GetMeterProvider().Meter(
		selfObsScopeName,
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)
	spansProcessedCounter, err := otelconv.NewSDKProcessorSpanProcessed(meter)
	return spansProcessedCounter, err
}

// OnStart does nothing.
func (*simpleSpanProcessor) OnStart(context.Context, ReadWriteSpan) {}

// OnEnd immediately exports a ReadOnlySpan.
func (ssp *simpleSpanProcessor) OnEnd(s ReadOnlySpan) {
	ssp.exporterMu.Lock()
	defer ssp.exporterMu.Unlock()

	if ssp.exporter != nil && s.SpanContext().TraceFlags().IsSampled() {
		attrs := measureAttrsPool.Get().(*[]attribute.KeyValue)
		defer func() {
			*attrs = (*attrs)[:0] // reset the slice for reuse
			measureAttrsPool.Put(attrs)
		}()
		*attrs = append(*attrs,
			ssp.componentNameAttr,
			ssp.spansProcessedCounter.AttrComponentType(otelconv.ComponentTypeSimpleSpanProcessor))

		err := ssp.exporter.ExportSpans(context.Background(), []ReadOnlySpan{s})
		if err != nil {
			otel.Handle(err)
			*attrs = append(*attrs, semconv.ErrorType(err))
		}
		if ssp.selfObservabilityEnabled {
			// Add the span to the context to ensure the metric is recorded
			// with the correct span context.
			ctx := trace.ContextWithSpanContext(context.Background(), s.SpanContext())
			ssp.spansProcessedCounter.Add(ctx, 1, *attrs...)
		}
	}
}

// Shutdown shuts down the exporter this SimpleSpanProcessor exports to.
func (ssp *simpleSpanProcessor) Shutdown(ctx context.Context) error {
	var err error
	ssp.stopOnce.Do(func() {
		stopFunc := func(exp SpanExporter) (<-chan error, func()) {
			done := make(chan error, 1)
			return done, func() { done <- exp.Shutdown(ctx) }
		}

		// The exporter field of the simpleSpanProcessor needs to be zeroed to
		// signal it is shut down, meaning all subsequent calls to OnEnd will
		// be gracefully ignored. This needs to be done synchronously to avoid
		// any race condition.
		//
		// A closure is used to keep reference to the exporter and then the
		// field is zeroed. This ensures the simpleSpanProcessor is shut down
		// before the exporter. This order is important as it avoids a potential
		// deadlock. If the exporter shut down operation generates a span, that
		// span would need to be exported. Meaning, OnEnd would be called and
		// try acquiring the lock that is held here.
		ssp.exporterMu.Lock()
		done, shutdown := stopFunc(ssp.exporter)
		ssp.exporter = nil
		ssp.exporterMu.Unlock()

		go shutdown()

		// Wait for the exporter to shut down or the deadline to expire.
		select {
		case err = <-done:
		case <-ctx.Done():
			// It is possible for the exporter to have immediately shut down and
			// the context to be done simultaneously. In that case this outer
			// select statement will randomly choose a case. This will result in
			// a different returned error for similar scenarios. Instead, double
			// check if the exporter shut down at the same time and return that
			// error if so. This will ensure consistency as well as ensure
			// the caller knows the exporter shut down successfully (they can
			// already determine if the deadline is expired given they passed
			// the context).
			select {
			case err = <-done:
			default:
				err = ctx.Err()
			}
		}
	})
	return err
}

// ForceFlush does nothing as there is no data to flush.
func (*simpleSpanProcessor) ForceFlush(context.Context) error {
	return nil
}

// MarshalLog is the marshaling function used by the logging system to represent
// this Span Processor.
func (ssp *simpleSpanProcessor) MarshalLog() any {
	return struct {
		Type     string
		Exporter SpanExporter
	}{
		Type:     "SimpleSpanProcessor",
		Exporter: ssp.exporter,
	}
}
