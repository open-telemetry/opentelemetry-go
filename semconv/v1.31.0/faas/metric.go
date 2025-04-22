// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/faas"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	// TriggerPubsub is a function is set to be executed when messages are sent to a
	// messaging system.
	TriggerPubsub TriggerAttr = "pubsub"
	// TriggerTimer is a function is scheduled to be executed regularly.
	TriggerTimer TriggerAttr = "timer"
	// TriggerOther is the if none of the others apply.
	TriggerOther TriggerAttr = "other"
)

// Coldstarts is an instrument used to record metric values conforming to the
// "faas.coldstarts" semantic conventions. It represents the number of invocation
// cold starts.
type Coldstarts struct {
	inst metric.Int64Counter
}

// NewColdstarts returns a new Coldstarts instrument.
func NewColdstarts(m metric.Meter) (Coldstarts, error) {
	i, err := m.Int64Counter(
	    "faas.coldstarts",
	    metric.WithDescription("Number of invocation cold starts"),
	    metric.WithUnit("{coldstart}"),
	)
	if err != nil {
	    return Coldstarts{}, err
	}
	return Coldstarts{i}, nil
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
	attrs ...ColdstartsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m Coldstarts) conv(in []ColdstartsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.coldstartsAttr()
	}
	return out
}

// ColdstartsAttr is an optional attribute for the Coldstarts instrument.
type ColdstartsAttr interface {
    coldstartsAttr() attribute.KeyValue
}

type coldstartsAttr struct {
	kv attribute.KeyValue
}

func (a coldstartsAttr) coldstartsAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Coldstarts) TriggerAttr(val TriggerAttr) ColdstartsAttr {
	return coldstartsAttr{kv: attribute.String("faas.trigger", string(val))}
}

// CPUUsage is an instrument used to record metric values conforming to the
// "faas.cpu_usage" semantic conventions. It represents the distribution of CPU
// usage per invocation.
type CPUUsage struct {
	inst metric.Float64Histogram
}

// NewCPUUsage returns a new CPUUsage instrument.
func NewCPUUsage(m metric.Meter) (CPUUsage, error) {
	i, err := m.Float64Histogram(
	    "faas.cpu_usage",
	    metric.WithDescription("Distribution of CPU usage per invocation"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return CPUUsage{}, err
	}
	return CPUUsage{i}, nil
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
	attrs ...CPUUsageAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m CPUUsage) conv(in []CPUUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.cpuUsageAttr()
	}
	return out
}

// CPUUsageAttr is an optional attribute for the CPUUsage instrument.
type CPUUsageAttr interface {
    cpuUsageAttr() attribute.KeyValue
}

type cpuUsageAttr struct {
	kv attribute.KeyValue
}

func (a cpuUsageAttr) cpuUsageAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (CPUUsage) TriggerAttr(val TriggerAttr) CPUUsageAttr {
	return cpuUsageAttr{kv: attribute.String("faas.trigger", string(val))}
}

// Errors is an instrument used to record metric values conforming to the
// "faas.errors" semantic conventions. It represents the number of invocation
// errors.
type Errors struct {
	inst metric.Int64Counter
}

// NewErrors returns a new Errors instrument.
func NewErrors(m metric.Meter) (Errors, error) {
	i, err := m.Int64Counter(
	    "faas.errors",
	    metric.WithDescription("Number of invocation errors"),
	    metric.WithUnit("{error}"),
	)
	if err != nil {
	    return Errors{}, err
	}
	return Errors{i}, nil
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
	attrs ...ErrorsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m Errors) conv(in []ErrorsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.errorsAttr()
	}
	return out
}

// ErrorsAttr is an optional attribute for the Errors instrument.
type ErrorsAttr interface {
    errorsAttr() attribute.KeyValue
}

type errorsAttr struct {
	kv attribute.KeyValue
}

func (a errorsAttr) errorsAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Errors) TriggerAttr(val TriggerAttr) ErrorsAttr {
	return errorsAttr{kv: attribute.String("faas.trigger", string(val))}
}

// InitDuration is an instrument used to record metric values conforming to the
// "faas.init_duration" semantic conventions. It represents the measures the
// duration of the function's initialization, such as a cold start.
type InitDuration struct {
	inst metric.Float64Histogram
}

// NewInitDuration returns a new InitDuration instrument.
func NewInitDuration(m metric.Meter) (InitDuration, error) {
	i, err := m.Float64Histogram(
	    "faas.init_duration",
	    metric.WithDescription("Measures the duration of the function's initialization, such as a cold start"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return InitDuration{}, err
	}
	return InitDuration{i}, nil
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
	attrs ...InitDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m InitDuration) conv(in []InitDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.initDurationAttr()
	}
	return out
}

// InitDurationAttr is an optional attribute for the InitDuration instrument.
type InitDurationAttr interface {
    initDurationAttr() attribute.KeyValue
}

type initDurationAttr struct {
	kv attribute.KeyValue
}

func (a initDurationAttr) initDurationAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (InitDuration) TriggerAttr(val TriggerAttr) InitDurationAttr {
	return initDurationAttr{kv: attribute.String("faas.trigger", string(val))}
}

// Invocations is an instrument used to record metric values conforming to the
// "faas.invocations" semantic conventions. It represents the number of
// successful invocations.
type Invocations struct {
	inst metric.Int64Counter
}

// NewInvocations returns a new Invocations instrument.
func NewInvocations(m metric.Meter) (Invocations, error) {
	i, err := m.Int64Counter(
	    "faas.invocations",
	    metric.WithDescription("Number of successful invocations"),
	    metric.WithUnit("{invocation}"),
	)
	if err != nil {
	    return Invocations{}, err
	}
	return Invocations{i}, nil
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
	attrs ...InvocationsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m Invocations) conv(in []InvocationsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.invocationsAttr()
	}
	return out
}

// InvocationsAttr is an optional attribute for the Invocations instrument.
type InvocationsAttr interface {
    invocationsAttr() attribute.KeyValue
}

type invocationsAttr struct {
	kv attribute.KeyValue
}

func (a invocationsAttr) invocationsAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Invocations) TriggerAttr(val TriggerAttr) InvocationsAttr {
	return invocationsAttr{kv: attribute.String("faas.trigger", string(val))}
}

// InvokeDuration is an instrument used to record metric values conforming to the
// "faas.invoke_duration" semantic conventions. It represents the measures the
// duration of the function's logic execution.
type InvokeDuration struct {
	inst metric.Float64Histogram
}

// NewInvokeDuration returns a new InvokeDuration instrument.
func NewInvokeDuration(m metric.Meter) (InvokeDuration, error) {
	i, err := m.Float64Histogram(
	    "faas.invoke_duration",
	    metric.WithDescription("Measures the duration of the function's logic execution"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return InvokeDuration{}, err
	}
	return InvokeDuration{i}, nil
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
	attrs ...InvokeDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m InvokeDuration) conv(in []InvokeDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.invokeDurationAttr()
	}
	return out
}

// InvokeDurationAttr is an optional attribute for the InvokeDuration instrument.
type InvokeDurationAttr interface {
    invokeDurationAttr() attribute.KeyValue
}

type invokeDurationAttr struct {
	kv attribute.KeyValue
}

func (a invokeDurationAttr) invokeDurationAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (InvokeDuration) TriggerAttr(val TriggerAttr) InvokeDurationAttr {
	return invokeDurationAttr{kv: attribute.String("faas.trigger", string(val))}
}

// MemUsage is an instrument used to record metric values conforming to the
// "faas.mem_usage" semantic conventions. It represents the distribution of max
// memory usage per invocation.
type MemUsage struct {
	inst metric.Int64Histogram
}

// NewMemUsage returns a new MemUsage instrument.
func NewMemUsage(m metric.Meter) (MemUsage, error) {
	i, err := m.Int64Histogram(
	    "faas.mem_usage",
	    metric.WithDescription("Distribution of max memory usage per invocation"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return MemUsage{}, err
	}
	return MemUsage{i}, nil
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
	attrs ...MemUsageAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m MemUsage) conv(in []MemUsageAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.memUsageAttr()
	}
	return out
}

// MemUsageAttr is an optional attribute for the MemUsage instrument.
type MemUsageAttr interface {
    memUsageAttr() attribute.KeyValue
}

type memUsageAttr struct {
	kv attribute.KeyValue
}

func (a memUsageAttr) memUsageAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (MemUsage) TriggerAttr(val TriggerAttr) MemUsageAttr {
	return memUsageAttr{kv: attribute.String("faas.trigger", string(val))}
}

// NetIo is an instrument used to record metric values conforming to the
// "faas.net_io" semantic conventions. It represents the distribution of net I/O
// usage per invocation.
type NetIo struct {
	inst metric.Int64Histogram
}

// NewNetIo returns a new NetIo instrument.
func NewNetIo(m metric.Meter) (NetIo, error) {
	i, err := m.Int64Histogram(
	    "faas.net_io",
	    metric.WithDescription("Distribution of net I/O usage per invocation"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return NetIo{}, err
	}
	return NetIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NetIo) Name() string {
	return "faas.net_io"
}

// Unit returns the semantic convention unit of the instrument
func (NetIo) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NetIo) Description() string {
	return "Distribution of net I/O usage per invocation"
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m NetIo) Record(
    ctx context.Context,
    val int64,
	attrs ...NetIoAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NetIo) conv(in []NetIoAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.netIoAttr()
	}
	return out
}

// NetIoAttr is an optional attribute for the NetIo instrument.
type NetIoAttr interface {
    netIoAttr() attribute.KeyValue
}

type netIoAttr struct {
	kv attribute.KeyValue
}

func (a netIoAttr) netIoAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (NetIo) TriggerAttr(val TriggerAttr) NetIoAttr {
	return netIoAttr{kv: attribute.String("faas.trigger", string(val))}
}

// Timeouts is an instrument used to record metric values conforming to the
// "faas.timeouts" semantic conventions. It represents the number of invocation
// timeouts.
type Timeouts struct {
	inst metric.Int64Counter
}

// NewTimeouts returns a new Timeouts instrument.
func NewTimeouts(m metric.Meter) (Timeouts, error) {
	i, err := m.Int64Counter(
	    "faas.timeouts",
	    metric.WithDescription("Number of invocation timeouts"),
	    metric.WithUnit("{timeout}"),
	)
	if err != nil {
	    return Timeouts{}, err
	}
	return Timeouts{i}, nil
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
	attrs ...TimeoutsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m Timeouts) conv(in []TimeoutsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.timeoutsAttr()
	}
	return out
}

// TimeoutsAttr is an optional attribute for the Timeouts instrument.
type TimeoutsAttr interface {
    timeoutsAttr() attribute.KeyValue
}

type timeoutsAttr struct {
	kv attribute.KeyValue
}

func (a timeoutsAttr) timeoutsAttr() attribute.KeyValue {
    return a.kv
}

// Trigger returns an optional attribute for the "faas.trigger" semantic
// convention. It represents the type of the trigger which caused this function
// invocation.
func (Timeouts) TriggerAttr(val TriggerAttr) TimeoutsAttr {
	return timeoutsAttr{kv: attribute.String("faas.trigger", string(val))}
}