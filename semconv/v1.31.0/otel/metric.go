// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/otel"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	// ComponentTypeOtlpGRPCSpanExporter is the OTLP span exporter over gRPC with
	// protobuf serialization.
	ComponentTypeOtlpGRPCSpanExporter ComponentTypeAttr = "otlp_grpc_span_exporter"
	// ComponentTypeOtlpHTTPSpanExporter is the OTLP span exporter over HTTP with
	// protobuf serialization.
	ComponentTypeOtlpHTTPSpanExporter ComponentTypeAttr = "otlp_http_span_exporter"
	// ComponentTypeOtlpHTTPJSONSpanExporter is the OTLP span exporter over HTTP
	// with JSON serialization.
	ComponentTypeOtlpHTTPJSONSpanExporter ComponentTypeAttr = "otlp_http_json_span_exporter"
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

// SDKExporterSpanExportedCount is an instrument used to record metric values
// conforming to the "otel.sdk.exporter.span.exported.count" semantic
// conventions. It represents the number of spans for which the export has
// finished, either successful or failed.
type SDKExporterSpanExportedCount struct {
	inst metric.Int64Counter
}

// NewSDKExporterSpanExportedCount returns a new SDKExporterSpanExportedCount
// instrument.
func NewSDKExporterSpanExportedCount(m metric.Meter) (SDKExporterSpanExportedCount, error) {
	i, err := m.Int64Counter(
	    "otel.sdk.exporter.span.exported.count",
	    metric.WithDescription("The number of spans for which the export has finished, either successful or failed"),
	    metric.WithUnit("{span}"),
	)
	if err != nil {
	    return SDKExporterSpanExportedCount{}, err
	}
	return SDKExporterSpanExportedCount{i}, nil
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
func (m SDKExporterSpanExportedCount) Add(
    ctx context.Context,
    incr int64,
	attrs ...SDKExporterSpanExportedCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m SDKExporterSpanExportedCount) conv(in []SDKExporterSpanExportedCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.sdkExporterSpanExportedCountAttr()
	}
	return out
}

// SDKExporterSpanExportedCountAttr is an optional attribute for the
// SDKExporterSpanExportedCount instrument.
type SDKExporterSpanExportedCountAttr interface {
    sdkExporterSpanExportedCountAttr() attribute.KeyValue
}

type sdkExporterSpanExportedCountAttr struct {
	kv attribute.KeyValue
}

func (a sdkExporterSpanExportedCountAttr) sdkExporterSpanExportedCountAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (SDKExporterSpanExportedCount) ErrorTypeAttr(val ErrorTypeAttr) SDKExporterSpanExportedCountAttr {
	return sdkExporterSpanExportedCountAttr{kv: attribute.String("error.type", string(val))}
}

// ComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKExporterSpanExportedCount) ComponentNameAttr(val string) SDKExporterSpanExportedCountAttr {
	return sdkExporterSpanExportedCountAttr{kv: attribute.String("otel.component.name", val)}
}

// ComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKExporterSpanExportedCount) ComponentTypeAttr(val ComponentTypeAttr) SDKExporterSpanExportedCountAttr {
	return sdkExporterSpanExportedCountAttr{kv: attribute.String("otel.component.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (SDKExporterSpanExportedCount) ServerAddressAttr(val string) SDKExporterSpanExportedCountAttr {
	return sdkExporterSpanExportedCountAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (SDKExporterSpanExportedCount) ServerPortAttr(val int) SDKExporterSpanExportedCountAttr {
	return sdkExporterSpanExportedCountAttr{kv: attribute.Int("server.port", val)}
}

// SDKExporterSpanInflightCount is an instrument used to record metric values
// conforming to the "otel.sdk.exporter.span.inflight.count" semantic
// conventions. It represents the number of spans which were passed to the
// exporter, but that have not been exported yet (neither successful, nor
// failed).
type SDKExporterSpanInflightCount struct {
	inst metric.Int64UpDownCounter
}

// NewSDKExporterSpanInflightCount returns a new SDKExporterSpanInflightCount
// instrument.
func NewSDKExporterSpanInflightCount(m metric.Meter) (SDKExporterSpanInflightCount, error) {
	i, err := m.Int64UpDownCounter(
	    "otel.sdk.exporter.span.inflight.count",
	    metric.WithDescription("The number of spans which were passed to the exporter, but that have not been exported yet (neither successful, nor failed)"),
	    metric.WithUnit("{span}"),
	)
	if err != nil {
	    return SDKExporterSpanInflightCount{}, err
	}
	return SDKExporterSpanInflightCount{i}, nil
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
func (m SDKExporterSpanInflightCount) Add(
    ctx context.Context,
    incr int64,
	attrs ...SDKExporterSpanInflightCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m SDKExporterSpanInflightCount) conv(in []SDKExporterSpanInflightCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.sdkExporterSpanInflightCountAttr()
	}
	return out
}

// SDKExporterSpanInflightCountAttr is an optional attribute for the
// SDKExporterSpanInflightCount instrument.
type SDKExporterSpanInflightCountAttr interface {
    sdkExporterSpanInflightCountAttr() attribute.KeyValue
}

type sdkExporterSpanInflightCountAttr struct {
	kv attribute.KeyValue
}

func (a sdkExporterSpanInflightCountAttr) sdkExporterSpanInflightCountAttr() attribute.KeyValue {
    return a.kv
}

// ComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKExporterSpanInflightCount) ComponentNameAttr(val string) SDKExporterSpanInflightCountAttr {
	return sdkExporterSpanInflightCountAttr{kv: attribute.String("otel.component.name", val)}
}

// ComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKExporterSpanInflightCount) ComponentTypeAttr(val ComponentTypeAttr) SDKExporterSpanInflightCountAttr {
	return sdkExporterSpanInflightCountAttr{kv: attribute.String("otel.component.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (SDKExporterSpanInflightCount) ServerAddressAttr(val string) SDKExporterSpanInflightCountAttr {
	return sdkExporterSpanInflightCountAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (SDKExporterSpanInflightCount) ServerPortAttr(val int) SDKExporterSpanInflightCountAttr {
	return sdkExporterSpanInflightCountAttr{kv: attribute.Int("server.port", val)}
}

// SDKProcessorSpanProcessedCount is an instrument used to record metric values
// conforming to the "otel.sdk.processor.span.processed.count" semantic
// conventions. It represents the number of spans for which the processing has
// finished, either successful or failed.
type SDKProcessorSpanProcessedCount struct {
	inst metric.Int64Counter
}

// NewSDKProcessorSpanProcessedCount returns a new SDKProcessorSpanProcessedCount
// instrument.
func NewSDKProcessorSpanProcessedCount(m metric.Meter) (SDKProcessorSpanProcessedCount, error) {
	i, err := m.Int64Counter(
	    "otel.sdk.processor.span.processed.count",
	    metric.WithDescription("The number of spans for which the processing has finished, either successful or failed"),
	    metric.WithUnit("{span}"),
	)
	if err != nil {
	    return SDKProcessorSpanProcessedCount{}, err
	}
	return SDKProcessorSpanProcessedCount{i}, nil
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
func (m SDKProcessorSpanProcessedCount) Add(
    ctx context.Context,
    incr int64,
	attrs ...SDKProcessorSpanProcessedCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m SDKProcessorSpanProcessedCount) conv(in []SDKProcessorSpanProcessedCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.sdkProcessorSpanProcessedCountAttr()
	}
	return out
}

// SDKProcessorSpanProcessedCountAttr is an optional attribute for the
// SDKProcessorSpanProcessedCount instrument.
type SDKProcessorSpanProcessedCountAttr interface {
    sdkProcessorSpanProcessedCountAttr() attribute.KeyValue
}

type sdkProcessorSpanProcessedCountAttr struct {
	kv attribute.KeyValue
}

func (a sdkProcessorSpanProcessedCountAttr) sdkProcessorSpanProcessedCountAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents a low-cardinality description of the failure reason.
// SDK Batching Span Processors MUST use `queue_full` for spans dropped due to a
// full queue.
func (SDKProcessorSpanProcessedCount) ErrorTypeAttr(val ErrorTypeAttr) SDKProcessorSpanProcessedCountAttr {
	return sdkProcessorSpanProcessedCountAttr{kv: attribute.String("error.type", string(val))}
}

// ComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorSpanProcessedCount) ComponentNameAttr(val string) SDKProcessorSpanProcessedCountAttr {
	return sdkProcessorSpanProcessedCountAttr{kv: attribute.String("otel.component.name", val)}
}

// ComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorSpanProcessedCount) ComponentTypeAttr(val ComponentTypeAttr) SDKProcessorSpanProcessedCountAttr {
	return sdkProcessorSpanProcessedCountAttr{kv: attribute.String("otel.component.type", string(val))}
}

// SDKProcessorSpanQueueCapacity is an instrument used to record metric values
// conforming to the "otel.sdk.processor.span.queue.capacity" semantic
// conventions. It represents the maximum number of spans the queue of a given
// instance of an SDK span processor can hold.
type SDKProcessorSpanQueueCapacity struct {
	inst metric.Int64UpDownCounter
}

// NewSDKProcessorSpanQueueCapacity returns a new SDKProcessorSpanQueueCapacity
// instrument.
func NewSDKProcessorSpanQueueCapacity(m metric.Meter) (SDKProcessorSpanQueueCapacity, error) {
	i, err := m.Int64UpDownCounter(
	    "otel.sdk.processor.span.queue.capacity",
	    metric.WithDescription("The maximum number of spans the queue of a given instance of an SDK span processor can hold"),
	    metric.WithUnit("{span}"),
	)
	if err != nil {
	    return SDKProcessorSpanQueueCapacity{}, err
	}
	return SDKProcessorSpanQueueCapacity{i}, nil
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
func (m SDKProcessorSpanQueueCapacity) Add(
    ctx context.Context,
    incr int64,
	attrs ...SDKProcessorSpanQueueCapacityAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m SDKProcessorSpanQueueCapacity) conv(in []SDKProcessorSpanQueueCapacityAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.sdkProcessorSpanQueueCapacityAttr()
	}
	return out
}

// SDKProcessorSpanQueueCapacityAttr is an optional attribute for the
// SDKProcessorSpanQueueCapacity instrument.
type SDKProcessorSpanQueueCapacityAttr interface {
    sdkProcessorSpanQueueCapacityAttr() attribute.KeyValue
}

type sdkProcessorSpanQueueCapacityAttr struct {
	kv attribute.KeyValue
}

func (a sdkProcessorSpanQueueCapacityAttr) sdkProcessorSpanQueueCapacityAttr() attribute.KeyValue {
    return a.kv
}

// ComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorSpanQueueCapacity) ComponentNameAttr(val string) SDKProcessorSpanQueueCapacityAttr {
	return sdkProcessorSpanQueueCapacityAttr{kv: attribute.String("otel.component.name", val)}
}

// ComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorSpanQueueCapacity) ComponentTypeAttr(val ComponentTypeAttr) SDKProcessorSpanQueueCapacityAttr {
	return sdkProcessorSpanQueueCapacityAttr{kv: attribute.String("otel.component.type", string(val))}
}

// SDKProcessorSpanQueueSize is an instrument used to record metric values
// conforming to the "otel.sdk.processor.span.queue.size" semantic conventions.
// It represents the number of spans in the queue of a given instance of an SDK
// span processor.
type SDKProcessorSpanQueueSize struct {
	inst metric.Int64UpDownCounter
}

// NewSDKProcessorSpanQueueSize returns a new SDKProcessorSpanQueueSize
// instrument.
func NewSDKProcessorSpanQueueSize(m metric.Meter) (SDKProcessorSpanQueueSize, error) {
	i, err := m.Int64UpDownCounter(
	    "otel.sdk.processor.span.queue.size",
	    metric.WithDescription("The number of spans in the queue of a given instance of an SDK span processor"),
	    metric.WithUnit("{span}"),
	)
	if err != nil {
	    return SDKProcessorSpanQueueSize{}, err
	}
	return SDKProcessorSpanQueueSize{i}, nil
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
func (m SDKProcessorSpanQueueSize) Add(
    ctx context.Context,
    incr int64,
	attrs ...SDKProcessorSpanQueueSizeAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m SDKProcessorSpanQueueSize) conv(in []SDKProcessorSpanQueueSizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.sdkProcessorSpanQueueSizeAttr()
	}
	return out
}

// SDKProcessorSpanQueueSizeAttr is an optional attribute for the
// SDKProcessorSpanQueueSize instrument.
type SDKProcessorSpanQueueSizeAttr interface {
    sdkProcessorSpanQueueSizeAttr() attribute.KeyValue
}

type sdkProcessorSpanQueueSizeAttr struct {
	kv attribute.KeyValue
}

func (a sdkProcessorSpanQueueSizeAttr) sdkProcessorSpanQueueSizeAttr() attribute.KeyValue {
    return a.kv
}

// ComponentName returns an optional attribute for the "otel.component.name"
// semantic convention. It represents a name uniquely identifying the instance of
// the OpenTelemetry component within its containing SDK instance.
func (SDKProcessorSpanQueueSize) ComponentNameAttr(val string) SDKProcessorSpanQueueSizeAttr {
	return sdkProcessorSpanQueueSizeAttr{kv: attribute.String("otel.component.name", val)}
}

// ComponentType returns an optional attribute for the "otel.component.type"
// semantic convention. It represents a name identifying the type of the
// OpenTelemetry component.
func (SDKProcessorSpanQueueSize) ComponentTypeAttr(val ComponentTypeAttr) SDKProcessorSpanQueueSizeAttr {
	return sdkProcessorSpanQueueSizeAttr{kv: attribute.String("otel.component.type", string(val))}
}

// SDKSpanEndedCount is an instrument used to record metric values conforming to
// the "otel.sdk.span.ended.count" semantic conventions. It represents the number
// of created spans for which the end operation was called.
type SDKSpanEndedCount struct {
	inst metric.Int64Counter
}

// NewSDKSpanEndedCount returns a new SDKSpanEndedCount instrument.
func NewSDKSpanEndedCount(m metric.Meter) (SDKSpanEndedCount, error) {
	i, err := m.Int64Counter(
	    "otel.sdk.span.ended.count",
	    metric.WithDescription("The number of created spans for which the end operation was called"),
	    metric.WithUnit("{span}"),
	)
	if err != nil {
	    return SDKSpanEndedCount{}, err
	}
	return SDKSpanEndedCount{i}, nil
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
func (m SDKSpanEndedCount) Add(
    ctx context.Context,
    incr int64,
	attrs ...SDKSpanEndedCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m SDKSpanEndedCount) conv(in []SDKSpanEndedCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.sdkSpanEndedCountAttr()
	}
	return out
}

// SDKSpanEndedCountAttr is an optional attribute for the SDKSpanEndedCount
// instrument.
type SDKSpanEndedCountAttr interface {
    sdkSpanEndedCountAttr() attribute.KeyValue
}

type sdkSpanEndedCountAttr struct {
	kv attribute.KeyValue
}

func (a sdkSpanEndedCountAttr) sdkSpanEndedCountAttr() attribute.KeyValue {
    return a.kv
}

// SpanSamplingResult returns an optional attribute for the
// "otel.span.sampling_result" semantic convention. It represents the result
// value of the sampler for this span.
func (SDKSpanEndedCount) SpanSamplingResultAttr(val SpanSamplingResultAttr) SDKSpanEndedCountAttr {
	return sdkSpanEndedCountAttr{kv: attribute.String("otel.span.sampling_result", string(val))}
}

// SDKSpanLiveCount is an instrument used to record metric values conforming to
// the "otel.sdk.span.live.count" semantic conventions. It represents the number
// of created spans for which the end operation has not been called yet.
type SDKSpanLiveCount struct {
	inst metric.Int64UpDownCounter
}

// NewSDKSpanLiveCount returns a new SDKSpanLiveCount instrument.
func NewSDKSpanLiveCount(m metric.Meter) (SDKSpanLiveCount, error) {
	i, err := m.Int64UpDownCounter(
	    "otel.sdk.span.live.count",
	    metric.WithDescription("The number of created spans for which the end operation has not been called yet"),
	    metric.WithUnit("{span}"),
	)
	if err != nil {
	    return SDKSpanLiveCount{}, err
	}
	return SDKSpanLiveCount{i}, nil
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
func (m SDKSpanLiveCount) Add(
    ctx context.Context,
    incr int64,
	attrs ...SDKSpanLiveCountAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m SDKSpanLiveCount) conv(in []SDKSpanLiveCountAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.sdkSpanLiveCountAttr()
	}
	return out
}

// SDKSpanLiveCountAttr is an optional attribute for the SDKSpanLiveCount
// instrument.
type SDKSpanLiveCountAttr interface {
    sdkSpanLiveCountAttr() attribute.KeyValue
}

type sdkSpanLiveCountAttr struct {
	kv attribute.KeyValue
}

func (a sdkSpanLiveCountAttr) sdkSpanLiveCountAttr() attribute.KeyValue {
    return a.kv
}

// SpanSamplingResult returns an optional attribute for the
// "otel.span.sampling_result" semantic convention. It represents the result
// value of the sampler for this span.
func (SDKSpanLiveCount) SpanSamplingResultAttr(val SpanSamplingResultAttr) SDKSpanLiveCountAttr {
	return sdkSpanLiveCountAttr{kv: attribute.String("otel.span.sampling_result", string(val))}
}