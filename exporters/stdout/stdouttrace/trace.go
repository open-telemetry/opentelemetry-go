// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdouttrace // import "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace/internal/x"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.36.0"
	"go.opentelemetry.io/otel/semconv/v1.36.0/otelconv"
)

// otelComponentType is a name identifying the type of the OpenTelemetry component.
const otelComponentType = "stdout_trace_exporter"

var zeroTime time.Time

var _ trace.SpanExporter = &Exporter{}

// New creates an Exporter with the passed options.
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options...)

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	exporter := &Exporter{
		encoder:    enc,
		timestamps: cfg.Timestamps,
	}

	if !x.SelfObservability.Enabled() {
		return exporter, nil
	}

	exporter.selfObservabilityEnabled = true
	exporter.selfObservabilityAttrs = []attribute.KeyValue{
		semconv.OTelComponentName(fmt.Sprintf("%s/%d", otelComponentType, nextExporterID())),
		semconv.OTelComponentTypeKey.String(otelComponentType),
	}

	mp := otel.GetMeterProvider()
	m := mp.Meter(
		"go.opentelemetry.io/otel/exporters/stdout/stdouttrace",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err, e error
	if exporter.spanInflightMetric, e = otelconv.NewSDKExporterSpanInflight(m); e != nil {
		e = fmt.Errorf("failed to create span inflight metric: %w", e)
		err = errors.Join(err, e)
	}
	if exporter.spanExportedMetric, e = otelconv.NewSDKExporterSpanExported(m); e != nil {
		e = fmt.Errorf("failed to create span exported metric: %w", e)
		err = errors.Join(err, e)
	}
	if exporter.operationDurationMetric, e = otelconv.NewSDKExporterOperationDuration(m); e != nil {
		e = fmt.Errorf("failed to create operation duration metric: %w", e)
		err = errors.Join(err, e)
	}

	return exporter, err
}

// Exporter is an implementation of trace.SpanSyncer that writes spans to stdout.
type Exporter struct {
	encoder    *json.Encoder
	encoderMu  sync.Mutex
	timestamps bool

	stoppedMu sync.RWMutex
	stopped   bool

	selfObservabilityEnabled bool
	selfObservabilityAttrs   []attribute.KeyValue // selfObservability common attributes
	spanInflightMetric       otelconv.SDKExporterSpanInflight
	spanExportedMetric       otelconv.SDKExporterSpanExported
	operationDurationMetric  otelconv.SDKExporterOperationDuration
}

// ExportSpans writes spans in json format to stdout.
func (e *Exporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) (err error) {
	if e.selfObservabilityEnabled {
		count := int64(len(spans))

		e.spanInflightMetric.Add(ctx, count, e.selfObservabilityAttrs...)
		defer func(starting time.Time) {
			// additional attributes for self-observability,
			// only spanExportedMetric and operationDurationMetric are supported
			addAttrs := make([]attribute.KeyValue, len(e.selfObservabilityAttrs), len(e.selfObservabilityAttrs)+1)
			copy(addAttrs, e.selfObservabilityAttrs)
			if err != nil {
				addAttrs = append(addAttrs, semconv.ErrorType(err))
			}

			e.spanInflightMetric.Add(ctx, -count, e.selfObservabilityAttrs...)
			e.spanExportedMetric.Add(ctx, count, addAttrs...)
			e.operationDurationMetric.Record(ctx, time.Since(starting).Seconds(), addAttrs...)
		}(time.Now())
	}

	if err := ctx.Err(); err != nil {
		return err
	}
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		return nil
	}

	if len(spans) == 0 {
		return nil
	}

	stubs := tracetest.SpanStubsFromReadOnlySpans(spans)

	e.encoderMu.Lock()
	defer e.encoderMu.Unlock()
	for i := range stubs {
		stub := &stubs[i]
		// Remove timestamps
		if !e.timestamps {
			stub.StartTime = zeroTime
			stub.EndTime = zeroTime
			for j := range stub.Events {
				ev := &stub.Events[j]
				ev.Time = zeroTime
			}
		}

		// Encode span stubs, one by one
		if err := e.encoder.Encode(stub); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown is called to stop the exporter, it performs no action.
func (e *Exporter) Shutdown(context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	return nil
}

// MarshalLog is the marshaling function used by the logging system to represent this Exporter.
func (e *Exporter) MarshalLog() any {
	return struct {
		Type           string
		WithTimestamps bool
	}{
		Type:           "stdout",
		WithTimestamps: e.timestamps,
	}
}

var exporterIDCounter atomic.Int64

// nextExporterID returns a new unique ID for an exporter.
// the starting value is 0, and it increments by 1 for each call.
func nextExporterID() int64 {
	return exporterIDCounter.Add(1) - 1
}
