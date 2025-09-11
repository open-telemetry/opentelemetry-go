// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

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
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/counter"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal/x"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/log"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/semconv/v1.37.0/otelconv"
)

// otelComponentType is a name identifying the type of the OpenTelemetry component.
const otelComponentType = "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

var _ log.Exporter = &Exporter{}

// Exporter writes JSON-encoded log records to an [io.Writer] ([os.Stdout] by default).
// Exporter must be created with [New].
type Exporter struct {
	encoder    atomic.Pointer[json.Encoder]
	timestamps bool
	inst       *instrumentationImpl
}

type instrumentationImpl struct {
	attrs                   []attribute.KeyValue
	inflightMetric          otelconv.SDKExporterLogInflight
	exportedMetric          otelconv.SDKExporterLogExported
	operationDurationMetric otelconv.SDKExporterOperationDuration
}

// New creates an [Exporter].
func New(options ...Option) (*Exporter, error) {
	cfg := newConfig(options)

	enc := json.NewEncoder(cfg.Writer)
	if cfg.PrettyPrint {
		enc.SetIndent("", "\t")
	}

	e := Exporter{
		timestamps: cfg.Timestamps,
	}
	e.encoder.Store(enc)
	selfObs, err := newInstrumentation()
	if err != nil {
		return nil, err
	}
	e.inst = selfObs

	return &e, nil
}

func newInstrumentation() (*instrumentationImpl, error) {
	if !x.SelfObservability.Enabled() {
		return nil, nil
	}

	inst := &instrumentationImpl{
		attrs: []attribute.KeyValue{
			semconv.OTelComponentName(fmt.Sprintf("%s/%d", otelComponentType, counter.NextExporterID())),
			semconv.OTelComponentTypeKey.String(otelComponentType),
		},
	}

	mp := otel.GetMeterProvider()
	m := mp.Meter("go.opentelemetry.io/otel/exporters/stdout/stdoutlog",
		metric.WithInstrumentationVersion(sdk.Version()),
		metric.WithSchemaURL(semconv.SchemaURL),
	)

	var err, e error
	inst.inflightMetric, e = otelconv.NewSDKExporterLogInflight(m)
	err = errors.Join(err, e)

	inst.exportedMetric, e = otelconv.NewSDKExporterLogExported(m)
	err = errors.Join(err, e)

	inst.operationDurationMetric, e = otelconv.NewSDKExporterOperationDuration(m)
	err = errors.Join(err, e)

	return inst, err
}

// Export exports log records to writer.
func (e *Exporter) Export(ctx context.Context, records []log.Record) error {
	var err error
	if e.inst != nil && x.SelfObservability.Enabled() {
		err = e.exportWithSelfObservability(ctx, records)
	} else {
		err = e.exportWithoutSelfObservability(ctx, records)
	}
	return err
}

const bufferSize = 4

var attrPool = sync.Pool{
	New: func() any {
		buf := make([]attribute.KeyValue, 0, bufferSize)
		return &buf
	},
}

// exportWithSelfObservability exports logs with self-observability metrics.
func (e *Exporter) exportWithSelfObservability(ctx context.Context, records []log.Record) (err error) {
	count := int64(len(records))
	start := time.Now()

	e.inst.inflightMetric.Add(ctx, count, e.inst.attrs...)

	bufPtrAny := attrPool.Get()
	bufPtr, ok := bufPtrAny.(*[]attribute.KeyValue)
	if !ok || bufPtr == nil {
		bufPtr = &[]attribute.KeyValue{}
	}
	defer func() {
		*bufPtr = (*bufPtr)[:0]
		attrPool.Put(bufPtr)
	}()

	defer func() {
		addAttrs := (*bufPtr)[:0]
		addAttrs = append(addAttrs, e.inst.attrs...)

		if err != nil {
			addAttrs = append(addAttrs, semconv.ErrorType(err))
		}
		e.inst.exportedMetric.Add(ctx, count, addAttrs...)
		e.inst.inflightMetric.Add(ctx, -count, e.inst.attrs...)
		e.inst.operationDurationMetric.Record(ctx, time.Since(start).Seconds(), addAttrs...)

		*bufPtr = addAttrs
	}()

	err = e.exportWithoutSelfObservability(ctx, records)
	return err
}

// exportWithoutSelfObservability exports logs without self-observability metrics.
func (e *Exporter) exportWithoutSelfObservability(ctx context.Context, records []log.Record) error {
	enc := e.encoder.Load()
	if enc == nil {
		return nil
	}
	for _, record := range records {
		// Honor context cancellation.
		if err := ctx.Err(); err != nil {
			return err
		}

		recordJSON := e.newRecordJSON(record)
		if err := enc.Encode(recordJSON); err != nil {
			return err
		}
	}
	return nil
}

// Shutdown shuts down the Exporter.
// Calls to Export will perform no operation after this is called.
func (e *Exporter) Shutdown(context.Context) error {
	e.encoder.Store(nil)
	return nil
}

// ForceFlush performs no action.
func (*Exporter) ForceFlush(context.Context) error {
	return nil
}
