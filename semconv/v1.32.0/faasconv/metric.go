// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "faas" namespace.
package faasconv

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

// TriggerAttr is an attribute conforming to the faas.trigger semantic
// conventions. It represents the type of the trigger which caused this function
// invocation.
type TriggerAttr string

var (
	// TriggerDatasource is a response to some data source operation such as a
	// database or filesystem read/write.
	TriggerDatasource TriggerAttr = "datasource"
	// TriggerHTTP is the to provide an answer to an inbound HTTP request.
	TriggerHTTP TriggerAttr = "http"
	// TriggerPubSub is a function is set to be executed when messages are sent to a
	// messaging system.
	TriggerPubSub TriggerAttr = "pubsub"
	// TriggerTimer is a function is scheduled to be executed regularly.
	TriggerTimer TriggerAttr = "timer"
	// TriggerOther is the if none of the others apply.
	TriggerOther TriggerAttr = "other"
)

// Coldstarts is an instrument used to record metric values conforming to the
// "faas.coldstarts" semantic conventions. It represents the number of invocation
// cold starts.
type Coldstarts struct {
	metric.Int64Counter
}

// NewColdstarts returns a new Coldstarts instrument.
func NewColdstarts(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Coldstarts, error) {
	// Check if the meter is nil.
	if m == nil {
		return Coldstarts{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"faas.coldstarts",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Number of invocation cold starts"),
			metric.WithUnit("{coldstart}"),
		}, opt...)...,
	)
	if err != nil {
	    return Coldstarts{noop.Int64Counter{}}, err
	}
	return Coldstarts{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Coldstarts) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (Coldstarts) Name() string {
	return "faas.coldstarts"
}

// Unit returns the semantic convention unit of the instrument
func (Coldstarts) Unit() string {
	return "{coldstart}"
}

// Description returns the semantic convention description of the instrument
func (Coldstarts) Description() string {
	return "Number of invocation cold starts"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m Coldstarts) Add(
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

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Coldstarts) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// CPUUsage is an instrument used to record metric values conforming to the
// "faas.cpu_usage" semantic conventions. It represents the distribution of CPU
// usage per invocation.
type CPUUsage struct {
	metric.Float64Histogram
}

// NewCPUUsage returns a new CPUUsage instrument.
func NewCPUUsage(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (CPUUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return CPUUsage{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"faas.cpu_usage",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Distribution of CPU usage per invocation"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return CPUUsage{noop.Float64Histogram{}}, err
	}
	return CPUUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CPUUsage) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (CPUUsage) Name() string {
	return "faas.cpu_usage"
}

// Unit returns the semantic convention unit of the instrument
func (CPUUsage) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (CPUUsage) Description() string {
	return "Distribution of CPU usage per invocation"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m CPUUsage) Record(
	ctx context.Context,
	val float64,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (CPUUsage) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// Errors is an instrument used to record metric values conforming to the
// "faas.errors" semantic conventions. It represents the number of invocation
// errors.
type Errors struct {
	metric.Int64Counter
}

// NewErrors returns a new Errors instrument.
func NewErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Errors, error) {
	// Check if the meter is nil.
	if m == nil {
		return Errors{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"faas.errors",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Number of invocation errors"),
			metric.WithUnit("{error}"),
		}, opt...)...,
	)
	if err != nil {
	    return Errors{noop.Int64Counter{}}, err
	}
	return Errors{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Errors) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (Errors) Name() string {
	return "faas.errors"
}

// Unit returns the semantic convention unit of the instrument
func (Errors) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (Errors) Description() string {
	return "Number of invocation errors"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m Errors) Add(
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

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Errors) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// InitDuration is an instrument used to record metric values conforming to the
// "faas.init_duration" semantic conventions. It represents the measures the
// duration of the function's initialization, such as a cold start.
type InitDuration struct {
	metric.Float64Histogram
}

// NewInitDuration returns a new InitDuration instrument.
func NewInitDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (InitDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return InitDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"faas.init_duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Measures the duration of the function's initialization, such as a cold start"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return InitDuration{noop.Float64Histogram{}}, err
	}
	return InitDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m InitDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (InitDuration) Name() string {
	return "faas.init_duration"
}

// Unit returns the semantic convention unit of the instrument
func (InitDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (InitDuration) Description() string {
	return "Measures the duration of the function's initialization, such as a cold start"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m InitDuration) Record(
	ctx context.Context,
	val float64,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (InitDuration) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// Invocations is an instrument used to record metric values conforming to the
// "faas.invocations" semantic conventions. It represents the number of
// successful invocations.
type Invocations struct {
	metric.Int64Counter
}

// NewInvocations returns a new Invocations instrument.
func NewInvocations(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Invocations, error) {
	// Check if the meter is nil.
	if m == nil {
		return Invocations{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"faas.invocations",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Number of successful invocations"),
			metric.WithUnit("{invocation}"),
		}, opt...)...,
	)
	if err != nil {
	    return Invocations{noop.Int64Counter{}}, err
	}
	return Invocations{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Invocations) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (Invocations) Name() string {
	return "faas.invocations"
}

// Unit returns the semantic convention unit of the instrument
func (Invocations) Unit() string {
	return "{invocation}"
}

// Description returns the semantic convention description of the instrument
func (Invocations) Description() string {
	return "Number of successful invocations"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m Invocations) Add(
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

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Invocations) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// InvokeDuration is an instrument used to record metric values conforming to the
// "faas.invoke_duration" semantic conventions. It represents the measures the
// duration of the function's logic execution.
type InvokeDuration struct {
	metric.Float64Histogram
}

// NewInvokeDuration returns a new InvokeDuration instrument.
func NewInvokeDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (InvokeDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return InvokeDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"faas.invoke_duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Measures the duration of the function's logic execution"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return InvokeDuration{noop.Float64Histogram{}}, err
	}
	return InvokeDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m InvokeDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (InvokeDuration) Name() string {
	return "faas.invoke_duration"
}

// Unit returns the semantic convention unit of the instrument
func (InvokeDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (InvokeDuration) Description() string {
	return "Measures the duration of the function's logic execution"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m InvokeDuration) Record(
	ctx context.Context,
	val float64,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (InvokeDuration) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// MemUsage is an instrument used to record metric values conforming to the
// "faas.mem_usage" semantic conventions. It represents the distribution of max
// memory usage per invocation.
type MemUsage struct {
	metric.Int64Histogram
}

// NewMemUsage returns a new MemUsage instrument.
func NewMemUsage(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (MemUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return MemUsage{noop.Int64Histogram{}}, nil
	}

	i, err := m.Int64Histogram(
		"faas.mem_usage",
		append([]metric.Int64HistogramOption{
			metric.WithDescription("Distribution of max memory usage per invocation"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return MemUsage{noop.Int64Histogram{}}, err
	}
	return MemUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m MemUsage) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (MemUsage) Name() string {
	return "faas.mem_usage"
}

// Unit returns the semantic convention unit of the instrument
func (MemUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (MemUsage) Description() string {
	return "Distribution of max memory usage per invocation"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m MemUsage) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (MemUsage) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// NetIO is an instrument used to record metric values conforming to the
// "faas.net_io" semantic conventions. It represents the distribution of net I/O
// usage per invocation.
type NetIO struct {
	metric.Int64Histogram
}

// NewNetIO returns a new NetIO instrument.
func NewNetIO(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (NetIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return NetIO{noop.Int64Histogram{}}, nil
	}

	i, err := m.Int64Histogram(
		"faas.net_io",
		append([]metric.Int64HistogramOption{
			metric.WithDescription("Distribution of net I/O usage per invocation"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return NetIO{noop.Int64Histogram{}}, err
	}
	return NetIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NetIO) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (NetIO) Name() string {
	return "faas.net_io"
}

// Unit returns the semantic convention unit of the instrument
func (NetIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetIO) Description() string {
	return "Distribution of net I/O usage per invocation"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m NetIO) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(
		*o,
		metric.WithAttributes(
			attrs...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (NetIO) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}

// Timeouts is an instrument used to record metric values conforming to the
// "faas.timeouts" semantic conventions. It represents the number of invocation
// timeouts.
type Timeouts struct {
	metric.Int64Counter
}

// NewTimeouts returns a new Timeouts instrument.
func NewTimeouts(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Timeouts, error) {
	// Check if the meter is nil.
	if m == nil {
		return Timeouts{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"faas.timeouts",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Number of invocation timeouts"),
			metric.WithUnit("{timeout}"),
		}, opt...)...,
	)
	if err != nil {
	    return Timeouts{noop.Int64Counter{}}, err
	}
	return Timeouts{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Timeouts) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (Timeouts) Name() string {
	return "faas.timeouts"
}

// Unit returns the semantic convention unit of the instrument
func (Timeouts) Unit() string {
	return "{timeout}"
}

// Description returns the semantic convention description of the instrument
func (Timeouts) Description() string {
	return "Number of invocation timeouts"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m Timeouts) Add(
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

// AttrTrigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Timeouts) AttrTrigger(val TriggerAttr) attribute.KeyValue {
	return attribute.String("faas.trigger", string(val))
}