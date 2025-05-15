// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "otel" namespace.
package otelconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
)

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the describes a class of error the operation ended
// with.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

// ComponentTypeAttr is an attribute conforming to the otel.component.type
// semantic conventions. It represents a name identifying the type of the
// OpenTelemetry component.
type ComponentTypeAttr string

var (
	// ComponentTypeBatchingSpanProcessor is the builtin SDK Batching Span
	// Processor.
	ComponentTypeBatchingSpanProcessor ComponentTypeAttr = "batching_span_processor"
	// ComponentTypeSimpleSpanProcessor is the builtin SDK Simple Span Processor.
	ComponentTypeSimpleSpanProcessor ComponentTypeAttr = "simple_span_processor"
	// ComponentTypeBatchingLogProcessor is the builtin SDK Batching LogRecord
	// Processor.
	ComponentTypeBatchingLogProcessor ComponentTypeAttr = "batching_log_processor"
	// ComponentTypeSimpleLogProcessor is the builtin SDK Simple LogRecord
	// Processor.
	ComponentTypeSimpleLogProcessor ComponentTypeAttr = "simple_log_processor"
	// ComponentTypeOtlpGRPCSpanExporter is the OTLP span exporter over gRPC with
	// protobuf serialization.
	ComponentTypeOtlpGRPCSpanExporter ComponentTypeAttr = "otlp_grpc_span_exporter"
	// ComponentTypeOtlpHTTPSpanExporter is the OTLP span exporter over HTTP with
	// protobuf serialization.
	ComponentTypeOtlpHTTPSpanExporter ComponentTypeAttr = "otlp_http_span_exporter"
	// ComponentTypeOtlpHTTPJSONSpanExporter is the OTLP span exporter over HTTP
	// with JSON serialization.
	ComponentTypeOtlpHTTPJSONSpanExporter ComponentTypeAttr = "otlp_http_json_span_exporter"
	// ComponentTypeOtlpGRPCLogExporter is the OTLP LogRecord exporter over gRPC
	// with protobuf serialization.
	ComponentTypeOtlpGRPCLogExporter ComponentTypeAttr = "otlp_grpc_log_exporter"
	// ComponentTypeOtlpHTTPLogExporter is the OTLP LogRecord exporter over HTTP
	// with protobuf serialization.
	ComponentTypeOtlpHTTPLogExporter ComponentTypeAttr = "otlp_http_log_exporter"
	// ComponentTypeOtlpHTTPJSONLogExporter is the OTLP LogRecord exporter over HTTP
	// with JSON serialization.
	ComponentTypeOtlpHTTPJSONLogExporter ComponentTypeAttr = "otlp_http_json_log_exporter"
)

// SpanSamplingResultAttr is an attribute conforming to the
// otel.span.sampling_result semantic conventions. It represents the result value
// of the sampler for this span.
type SpanSamplingResultAttr string

var (
	// SpanSamplingResultDrop is the span is not sampled and not recording.
	SpanSamplingResultDrop SpanSamplingResultAttr = "DROP"
	// SpanSamplingResultRecordOnly is the span is not sampled, but recording.
	SpanSamplingResultRecordOnly SpanSamplingResultAttr = "RECORD_ONLY"
	// SpanSamplingResultRecordAndSample is the span is sampled and recording.
	SpanSamplingResultRecordAndSample SpanSamplingResultAttr = "RECORD_AND_SAMPLE"
)

// SDKExporterLogExported is an instrument used to record metric values
// conforming to the "otel.sdk.exporter.log.exported" semantic conventions. It
// represents the number of log records for which the export has finished, either
// successful or failed.
type SDKExporterLogExported struct {
	metric.Int64Counter
}

// NewSDKExporterLogExported returns a new SDKExporterLogExported instrument.
func NewSDKExporterLogExported(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (SDKExporterLogExported, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKExporterLogExported{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"otel.sdk.exporter.log.exported",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of log records for which the export has finished, either successful or failed"),
			metric.WithUnit("{log_record}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKExporterLogExported{noop.Int64Counter{}}, err
	}
	return SDKExporterLogExported{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKExporterLogExported) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (SDKExporterLogExported) Name() string {
	return "otel.sdk.exporter.log.exported"
}

// Unit returns the semantic convention unit of the instrument
func (SDKExporterLogExported) Unit() string {
	return "{log_record}"
}

// Description returns the semantic convention description of the instrument
func (SDKExporterLogExported) Description() string {
	return "The number of log records for which the export has finished, either successful or failed"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For successful exports, `error.type` MUST NOT be set. For failed exports,
// `error.type` must contain the failure cause.
// For exporters with partial success semantics (e.g. OTLP with
// `rejected_log_records`), rejected log records must count as failed and only
// non-rejected log records count as success.
// If no rejection reason is available, `rejected` SHOULD be used as value for
// `error.type`.
func (m SDKExporterLogExported) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (SDKExporterLogExported) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKExporterLogExported) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKExporterLogExported) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (SDKExporterLogExported) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (SDKExporterLogExported) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// SDKExporterLogInflight is an instrument used to record metric values
// conforming to the "otel.sdk.exporter.log.inflight" semantic conventions. It
// represents the number of log records which were passed to the exporter, but
// that have not been exported yet (neither successful, nor failed).
type SDKExporterLogInflight struct {
	metric.Int64UpDownCounter
}

// NewSDKExporterLogInflight returns a new SDKExporterLogInflight instrument.
func NewSDKExporterLogInflight(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (SDKExporterLogInflight, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKExporterLogInflight{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"otel.sdk.exporter.log.inflight",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of log records which were passed to the exporter, but that have not been exported yet (neither successful, nor failed)"),
			metric.WithUnit("{log_record}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKExporterLogInflight{noop.Int64UpDownCounter{}}, err
	}
	return SDKExporterLogInflight{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKExporterLogInflight) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (SDKExporterLogInflight) Name() string {
	return "otel.sdk.exporter.log.inflight"
}

// Unit returns the semantic convention unit of the instrument
func (SDKExporterLogInflight) Unit() string {
	return "{log_record}"
}

// Description returns the semantic convention description of the instrument
func (SDKExporterLogInflight) Description() string {
	return "The number of log records which were passed to the exporter, but that have not been exported yet (neither successful, nor failed)"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For successful exports, `error.type` MUST NOT be set. For failed exports,
// `error.type` must contain the failure cause.
func (m SDKExporterLogInflight) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKExporterLogInflight) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKExporterLogInflight) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (SDKExporterLogInflight) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (SDKExporterLogInflight) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// SDKExporterSpanExportedCount is an instrument used to record metric values
// conforming to the "otel.sdk.exporter.span.exported.count" semantic
// conventions. It represents the number of spans for which the export has
// finished, either successful or failed.
type SDKExporterSpanExportedCount struct {
	metric.Int64Counter
}

// NewSDKExporterSpanExportedCount returns a new SDKExporterSpanExportedCount
// instrument.
func NewSDKExporterSpanExportedCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (SDKExporterSpanExportedCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKExporterSpanExportedCount{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"otel.sdk.exporter.span.exported.count",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of spans for which the export has finished, either successful or failed"),
			metric.WithUnit("{span}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKExporterSpanExportedCount{noop.Int64Counter{}}, err
	}
	return SDKExporterSpanExportedCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKExporterSpanExportedCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (SDKExporterSpanExportedCount) Name() string {
	return "otel.sdk.exporter.span.exported.count"
}

// Unit returns the semantic convention unit of the instrument
func (SDKExporterSpanExportedCount) Unit() string {
	return "{span}"
}

// Description returns the semantic convention description of the instrument
func (SDKExporterSpanExportedCount) Description() string {
	return "The number of spans for which the export has finished, either successful or failed"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For successful exports, `error.type` MUST NOT be set. For failed exports,
// `error.type` must contain the failure cause.
// For exporters with partial success semantics (e.g. OTLP with `rejected_spans`
// ), rejected spans must count as failed and only non-rejected spans count as
// success.
// If no rejection reason is available, `rejected` SHOULD be used as value for
// `error.type`.
func (m SDKExporterSpanExportedCount) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (SDKExporterSpanExportedCount) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKExporterSpanExportedCount) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKExporterSpanExportedCount) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (SDKExporterSpanExportedCount) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (SDKExporterSpanExportedCount) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// SDKExporterSpanInflightCount is an instrument used to record metric values
// conforming to the "otel.sdk.exporter.span.inflight.count" semantic
// conventions. It represents the number of spans which were passed to the
// exporter, but that have not been exported yet (neither successful, nor
// failed).
type SDKExporterSpanInflightCount struct {
	metric.Int64UpDownCounter
}

// NewSDKExporterSpanInflightCount returns a new SDKExporterSpanInflightCount
// instrument.
func NewSDKExporterSpanInflightCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (SDKExporterSpanInflightCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKExporterSpanInflightCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"otel.sdk.exporter.span.inflight.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed)"),
			metric.WithUnit("{span}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKExporterSpanInflightCount{noop.Int64UpDownCounter{}}, err
	}
	return SDKExporterSpanInflightCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKExporterSpanInflightCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (SDKExporterSpanInflightCount) Name() string {
	return "otel.sdk.exporter.span.inflight.count"
}

// Unit returns the semantic convention unit of the instrument
func (SDKExporterSpanInflightCount) Unit() string {
	return "{span}"
}

// Description returns the semantic convention description of the instrument
func (SDKExporterSpanInflightCount) Description() string {
	return "The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed)"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For successful exports, `error.type` MUST NOT be set. For failed exports,
// `error.type` must contain the failure cause.
func (m SDKExporterSpanInflightCount) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKExporterSpanInflightCount) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKExporterSpanInflightCount) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (SDKExporterSpanInflightCount) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (SDKExporterSpanInflightCount) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// SDKLogCreated is an instrument used to record metric values conforming to the
// "otel.sdk.log.created" semantic conventions. It represents the number of logs
// submitted to enabled SDK Loggers.
type SDKLogCreated struct {
	metric.Int64Counter
}

// NewSDKLogCreated returns a new SDKLogCreated instrument.
func NewSDKLogCreated(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (SDKLogCreated, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKLogCreated{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"otel.sdk.log.created",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of logs submitted to enabled SDK Loggers"),
			metric.WithUnit("{log_record}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKLogCreated{noop.Int64Counter{}}, err
	}
	return SDKLogCreated{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKLogCreated) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (SDKLogCreated) Name() string {
	return "otel.sdk.log.created"
}

// Unit returns the semantic convention unit of the instrument
func (SDKLogCreated) Unit() string {
	return "{log_record}"
}

// Description returns the semantic convention description of the instrument
func (SDKLogCreated) Description() string {
	return "The number of logs submitted to enabled SDK Loggers"
}

// Add adds incr to the existing count.
func (m SDKLogCreated) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// SDKProcessorLogProcessed is an instrument used to record metric values
// conforming to the "otel.sdk.processor.log.processed" semantic conventions. It
// represents the number of log records for which the processing has finished,
// either successful or failed.
type SDKProcessorLogProcessed struct {
	metric.Int64Counter
}

// NewSDKProcessorLogProcessed returns a new SDKProcessorLogProcessed instrument.
func NewSDKProcessorLogProcessed(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (SDKProcessorLogProcessed, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKProcessorLogProcessed{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"otel.sdk.processor.log.processed",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of log records for which the processing has finished, either successful or failed"),
			metric.WithUnit("{log_record}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKProcessorLogProcessed{noop.Int64Counter{}}, err
	}
	return SDKProcessorLogProcessed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKProcessorLogProcessed) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (SDKProcessorLogProcessed) Name() string {
	return "otel.sdk.processor.log.processed"
}

// Unit returns the semantic convention unit of the instrument
func (SDKProcessorLogProcessed) Unit() string {
	return "{log_record}"
}

// Description returns the semantic convention description of the instrument
func (SDKProcessorLogProcessed) Description() string {
	return "The number of log records for which the processing has finished, either successful or failed"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For successful processing, `error.type` MUST NOT be set. For failed
// processing, `error.type` must contain the failure cause.
// For the SDK Simple and Batching Log Record Processor a log record is
// considered to be processed already when it has been submitted to the exporter,
// not when the corresponding export call has finished.
func (m SDKProcessorLogProcessed) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents a low-cardinality description of the failure reason.
// SDK Batching Log Record Processors MUST use `queue_full` for log records
// dropped due to a full queue.
func (SDKProcessorLogProcessed) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorLogProcessed) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorLogProcessed) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// SDKProcessorLogQueueCapacity is an instrument used to record metric values
// conforming to the "otel.sdk.processor.log.queue.capacity" semantic
// conventions. It represents the maximum number of log records the queue of a
// given instance of an SDK Log Record processor can hold.
type SDKProcessorLogQueueCapacity struct {
	metric.Int64UpDownCounter
}

// NewSDKProcessorLogQueueCapacity returns a new SDKProcessorLogQueueCapacity
// instrument.
func NewSDKProcessorLogQueueCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (SDKProcessorLogQueueCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKProcessorLogQueueCapacity{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"otel.sdk.processor.log.queue.capacity",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The maximum number of log records the queue of a given instance of an SDK Log Record processor can hold"),
			metric.WithUnit("{log_record}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKProcessorLogQueueCapacity{noop.Int64UpDownCounter{}}, err
	}
	return SDKProcessorLogQueueCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKProcessorLogQueueCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (SDKProcessorLogQueueCapacity) Name() string {
	return "otel.sdk.processor.log.queue.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (SDKProcessorLogQueueCapacity) Unit() string {
	return "{log_record}"
}

// Description returns the semantic convention description of the instrument
func (SDKProcessorLogQueueCapacity) Description() string {
	return "The maximum number of log records the queue of a given instance of an SDK Log Record processor can hold"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// Only applies to Log Record processors which use a queue, e.g. the SDK Batching
// Log Record Processor.
func (m SDKProcessorLogQueueCapacity) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorLogQueueCapacity) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorLogQueueCapacity) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// SDKProcessorLogQueueSize is an instrument used to record metric values
// conforming to the "otel.sdk.processor.log.queue.size" semantic conventions. It
// represents the number of log records in the queue of a given instance of an
// SDK log processor.
type SDKProcessorLogQueueSize struct {
	metric.Int64UpDownCounter
}

// NewSDKProcessorLogQueueSize returns a new SDKProcessorLogQueueSize instrument.
func NewSDKProcessorLogQueueSize(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (SDKProcessorLogQueueSize, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKProcessorLogQueueSize{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"otel.sdk.processor.log.queue.size",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of log records in the queue of a given instance of an SDK log processor"),
			metric.WithUnit("{log_record}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKProcessorLogQueueSize{noop.Int64UpDownCounter{}}, err
	}
	return SDKProcessorLogQueueSize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKProcessorLogQueueSize) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (SDKProcessorLogQueueSize) Name() string {
	return "otel.sdk.processor.log.queue.size"
}

// Unit returns the semantic convention unit of the instrument
func (SDKProcessorLogQueueSize) Unit() string {
	return "{log_record}"
}

// Description returns the semantic convention description of the instrument
func (SDKProcessorLogQueueSize) Description() string {
	return "The number of log records in the queue of a given instance of an SDK log processor"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// Only applies to log record processors which use a queue, e.g. the SDK Batching
// Log Record Processor.
func (m SDKProcessorLogQueueSize) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorLogQueueSize) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorLogQueueSize) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// SDKProcessorSpanProcessedCount is an instrument used to record metric values
// conforming to the "otel.sdk.processor.span.processed.count" semantic
// conventions. It represents the number of spans for which the processing has
// finished, either successful or failed.
type SDKProcessorSpanProcessedCount struct {
	metric.Int64Counter
}

// NewSDKProcessorSpanProcessedCount returns a new SDKProcessorSpanProcessedCount
// instrument.
func NewSDKProcessorSpanProcessedCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (SDKProcessorSpanProcessedCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKProcessorSpanProcessedCount{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"otel.sdk.processor.span.processed.count",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of spans for which the processing has finished, either successful or failed"),
			metric.WithUnit("{span}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKProcessorSpanProcessedCount{noop.Int64Counter{}}, err
	}
	return SDKProcessorSpanProcessedCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKProcessorSpanProcessedCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (SDKProcessorSpanProcessedCount) Name() string {
	return "otel.sdk.processor.span.processed.count"
}

// Unit returns the semantic convention unit of the instrument
func (SDKProcessorSpanProcessedCount) Unit() string {
	return "{span}"
}

// Description returns the semantic convention description of the instrument
func (SDKProcessorSpanProcessedCount) Description() string {
	return "The number of spans for which the processing has finished, either successful or failed"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For successful processing, `error.type` MUST NOT be set. For failed
// processing, `error.type` must contain the failure cause.
// For the SDK Simple and Batching Span Processor a span is considered to be
// processed already when it has been submitted to the exporter, not when the
// corresponding export call has finished.
func (m SDKProcessorSpanProcessedCount) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents a low-cardinality description of the failure reason.
// SDK Batching Span Processors MUST use `queue_full` for spans dropped due to a
// full queue.
func (SDKProcessorSpanProcessedCount) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorSpanProcessedCount) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorSpanProcessedCount) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// SDKProcessorSpanQueueCapacity is an instrument used to record metric values
// conforming to the "otel.sdk.processor.span.queue.capacity" semantic
// conventions. It represents the maximum number of spans the queue of a given
// instance of an SDK span processor can hold.
type SDKProcessorSpanQueueCapacity struct {
	metric.Int64UpDownCounter
}

// NewSDKProcessorSpanQueueCapacity returns a new SDKProcessorSpanQueueCapacity
// instrument.
func NewSDKProcessorSpanQueueCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (SDKProcessorSpanQueueCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKProcessorSpanQueueCapacity{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"otel.sdk.processor.span.queue.capacity",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The maximum number of spans the queue of a given instance of an SDK span processor can hold"),
			metric.WithUnit("{span}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKProcessorSpanQueueCapacity{noop.Int64UpDownCounter{}}, err
	}
	return SDKProcessorSpanQueueCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKProcessorSpanQueueCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (SDKProcessorSpanQueueCapacity) Name() string {
	return "otel.sdk.processor.span.queue.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (SDKProcessorSpanQueueCapacity) Unit() string {
	return "{span}"
}

// Description returns the semantic convention description of the instrument
func (SDKProcessorSpanQueueCapacity) Description() string {
	return "The maximum number of spans the queue of a given instance of an SDK span processor can hold"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// Only applies to span processors which use a queue, e.g. the SDK Batching Span
// Processor.
func (m SDKProcessorSpanQueueCapacity) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorSpanQueueCapacity) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorSpanQueueCapacity) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// SDKProcessorSpanQueueSize is an instrument used to record metric values
// conforming to the "otel.sdk.processor.span.queue.size" semantic conventions.
// It represents the number of spans in the queue of a given instance of an SDK
// span processor.
type SDKProcessorSpanQueueSize struct {
	metric.Int64UpDownCounter
}

// NewSDKProcessorSpanQueueSize returns a new SDKProcessorSpanQueueSize
// instrument.
func NewSDKProcessorSpanQueueSize(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (SDKProcessorSpanQueueSize, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKProcessorSpanQueueSize{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"otel.sdk.processor.span.queue.size",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of spans in the queue of a given instance of an SDK span processor"),
			metric.WithUnit("{span}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKProcessorSpanQueueSize{noop.Int64UpDownCounter{}}, err
	}
	return SDKProcessorSpanQueueSize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKProcessorSpanQueueSize) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (SDKProcessorSpanQueueSize) Name() string {
	return "otel.sdk.processor.span.queue.size"
}

// Unit returns the semantic convention unit of the instrument
func (SDKProcessorSpanQueueSize) Unit() string {
	return "{span}"
}

// Description returns the semantic convention description of the instrument
func (SDKProcessorSpanQueueSize) Description() string {
	return "The number of spans in the queue of a given instance of an SDK span processor"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// Only applies to span processors which use a queue, e.g. the SDK Batching Span
// Processor.
func (m SDKProcessorSpanQueueSize) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorSpanQueueSize) AttrComponentName(val string) attribute.KeyValue {
	return attribute.String("otel.component.name", val)
}

// AttrComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorSpanQueueSize) AttrComponentType(val ComponentTypeAttr) attribute.KeyValue {
	return attribute.String("otel.component.type", string(val))
}

// SDKSpanEndedCount is an instrument used to record metric values conforming to
// the "otel.sdk.span.ended.count" semantic conventions. It represents the number
// of created spans for which the end operation was called.
type SDKSpanEndedCount struct {
	metric.Int64Counter
}

// NewSDKSpanEndedCount returns a new SDKSpanEndedCount instrument.
func NewSDKSpanEndedCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (SDKSpanEndedCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKSpanEndedCount{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"otel.sdk.span.ended.count",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of created spans for which the end operation was called"),
			metric.WithUnit("{span}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKSpanEndedCount{noop.Int64Counter{}}, err
	}
	return SDKSpanEndedCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKSpanEndedCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (SDKSpanEndedCount) Name() string {
	return "otel.sdk.span.ended.count"
}

// Unit returns the semantic convention unit of the instrument
func (SDKSpanEndedCount) Unit() string {
	return "{span}"
}

// Description returns the semantic convention description of the instrument
func (SDKSpanEndedCount) Description() string {
	return "The number of created spans for which the end operation was called"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For spans with `recording=true`: Implementations MUST record both
// `otel.sdk.span.live.count` and `otel.sdk.span.ended.count`.
// For spans with `recording=false`: If implementations decide to record this
// metric, they MUST also record `otel.sdk.span.live.count`.
func (m SDKSpanEndedCount) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrSpanSamplingResult returns an optional attribute for the
// "otel.span.sampling_result" semantic convention. It represents the result
// value of the sampler for this span.
func (SDKSpanEndedCount) AttrSpanSamplingResult(val SpanSamplingResultAttr) attribute.KeyValue {
	return attribute.String("otel.span.sampling_result", string(val))
}

// SDKSpanLiveCount is an instrument used to record metric values conforming to
// the "otel.sdk.span.live.count" semantic conventions. It represents the number
// of created spans for which the end operation has not been called yet.
type SDKSpanLiveCount struct {
	metric.Int64UpDownCounter
}

// NewSDKSpanLiveCount returns a new SDKSpanLiveCount instrument.
func NewSDKSpanLiveCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (SDKSpanLiveCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return SDKSpanLiveCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"otel.sdk.span.live.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of created spans for which the end operation has not been called yet"),
			metric.WithUnit("{span}"),
		}, opt...)...,
	)
	if err != nil {
	    return SDKSpanLiveCount{noop.Int64UpDownCounter{}}, err
	}
	return SDKSpanLiveCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SDKSpanLiveCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (SDKSpanLiveCount) Name() string {
	return "otel.sdk.span.live.count"
}

// Unit returns the semantic convention unit of the instrument
func (SDKSpanLiveCount) Unit() string {
	return "{span}"
}

// Description returns the semantic convention description of the instrument
func (SDKSpanLiveCount) Description() string {
	return "The number of created spans for which the end operation has not been called yet"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
//
// For spans with `recording=true`: Implementations MUST record both
// `otel.sdk.span.live.count` and `otel.sdk.span.ended.count`.
// For spans with `recording=false`: If implementations decide to record this
// metric, they MUST also record `otel.sdk.span.ended.count`.
func (m SDKSpanLiveCount) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrSpanSamplingResult returns an optional attribute for the
// "otel.span.sampling_result" semantic convention. It represents the result
// value of the sampler for this span.
func (SDKSpanLiveCount) AttrSpanSamplingResult(val SpanSamplingResultAttr) attribute.KeyValue {
	return attribute.String("otel.span.sampling_result", string(val))
}