// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "cicd" namespace.
package cicdconv

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

// PipelineResultAttr is an attribute conforming to the cicd.pipeline.result
// semantic conventions. It represents the result of a pipeline run.
type PipelineResultAttr string

var (
	// PipelineResultSuccess is the pipeline run finished successfully.
	PipelineResultSuccess PipelineResultAttr = "success"
	// PipelineResultFailure is the pipeline run did not finish successfully, eg.
	// due to a compile error or a failing test. Such failures are usually detected
	// by non-zero exit codes of the tools executed in the pipeline run.
	PipelineResultFailure PipelineResultAttr = "failure"
	// PipelineResultError is the pipeline run failed due to an error in the CICD
	// system, eg. due to the worker being killed.
	PipelineResultError PipelineResultAttr = "error"
	// PipelineResultTimeout is a timeout caused the pipeline run to be interrupted.
	PipelineResultTimeout PipelineResultAttr = "timeout"
	// PipelineResultCancellation is the pipeline run was cancelled, eg. by a user
	// manually cancelling the pipeline run.
	PipelineResultCancellation PipelineResultAttr = "cancellation"
	// PipelineResultSkip is the pipeline run was skipped, eg. due to a precondition
	// not being met.
	PipelineResultSkip PipelineResultAttr = "skip"
)

// PipelineRunStateAttr is an attribute conforming to the cicd.pipeline.run.state
// semantic conventions. It represents the pipeline run goes through these states
// during its lifecycle.
type PipelineRunStateAttr string

var (
	// PipelineRunStatePending is the run pending state spans from the event
	// triggering the pipeline run until the execution of the run starts (eg. time
	// spent in a queue, provisioning agents, creating run resources).
	PipelineRunStatePending PipelineRunStateAttr = "pending"
	// PipelineRunStateExecuting is the executing state spans the execution of any
	// run tasks (eg. build, test).
	PipelineRunStateExecuting PipelineRunStateAttr = "executing"
	// PipelineRunStateFinalizing is the finalizing state spans from when the run
	// has finished executing (eg. cleanup of run resources).
	PipelineRunStateFinalizing PipelineRunStateAttr = "finalizing"
)

// WorkerStateAttr is an attribute conforming to the cicd.worker.state semantic
// conventions. It represents the state of a CICD worker / agent.
type WorkerStateAttr string

var (
	// WorkerStateAvailable is the worker is not performing work for the CICD
	// system. It is available to the CICD system to perform work on (online /
	// idle).
	WorkerStateAvailable WorkerStateAttr = "available"
	// WorkerStateBusy is the worker is performing work for the CICD system.
	WorkerStateBusy WorkerStateAttr = "busy"
	// WorkerStateOffline is the worker is not available to the CICD system
	// (disconnected / down).
	WorkerStateOffline WorkerStateAttr = "offline"
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

// PipelineRunActive is an instrument used to record metric values conforming to
// the "cicd.pipeline.run.active" semantic conventions. It represents the number
// of pipeline runs currently active in the system by state.
type PipelineRunActive struct {
	metric.Int64UpDownCounter
}

// NewPipelineRunActive returns a new PipelineRunActive instrument.
func NewPipelineRunActive(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PipelineRunActive, error) {
	// Check if the meter is nil.
	if m == nil {
		return PipelineRunActive{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"cicd.pipeline.run.active",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of pipeline runs currently active in the system by state."),
			metric.WithUnit("{run}"),
		}, opt...)...,
	)
	if err != nil {
	    return PipelineRunActive{noop.Int64UpDownCounter{}}, err
	}
	return PipelineRunActive{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PipelineRunActive) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PipelineRunActive) Name() string {
	return "cicd.pipeline.run.active"
}

// Unit returns the semantic convention unit of the instrument
func (PipelineRunActive) Unit() string {
	return "{run}"
}

// Description returns the semantic convention description of the instrument
func (PipelineRunActive) Description() string {
	return "The number of pipeline runs currently active in the system by state."
}

// Add adds incr to the existing count.
//
// The pipelineName is the the human readable name of the pipeline within a CI/CD
// system.
//
// The pipelineRunState is the the pipeline run goes through these states during
// its lifecycle.
func (m PipelineRunActive) Add(
	ctx context.Context,
	incr int64,
	pipelineName string,
	pipelineRunState PipelineRunStateAttr,
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
			append(
				attrs,
				attribute.String("cicd.pipeline.name", pipelineName),
				attribute.String("cicd.pipeline.run.state", string(pipelineRunState)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PipelineRunDuration is an instrument used to record metric values conforming
// to the "cicd.pipeline.run.duration" semantic conventions. It represents the
// duration of a pipeline run grouped by pipeline, state and result.
type PipelineRunDuration struct {
	metric.Float64Histogram
}

// NewPipelineRunDuration returns a new PipelineRunDuration instrument.
func NewPipelineRunDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (PipelineRunDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return PipelineRunDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"cicd.pipeline.run.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("Duration of a pipeline run grouped by pipeline, state and result."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return PipelineRunDuration{noop.Float64Histogram{}}, err
	}
	return PipelineRunDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PipelineRunDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (PipelineRunDuration) Name() string {
	return "cicd.pipeline.run.duration"
}

// Unit returns the semantic convention unit of the instrument
func (PipelineRunDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (PipelineRunDuration) Description() string {
	return "Duration of a pipeline run grouped by pipeline, state and result."
}

// Record records val to the current distribution.
//
// The pipelineName is the the human readable name of the pipeline within a CI/CD
// system.
//
// The pipelineRunState is the the pipeline run goes through these states during
// its lifecycle.
//
// All additional attrs passed are included in the recorded value.
func (m PipelineRunDuration) Record(
	ctx context.Context,
	val float64,
	pipelineName string,
	pipelineRunState PipelineRunStateAttr,
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
			append(
				attrs,
				attribute.String("cicd.pipeline.name", pipelineName),
				attribute.String("cicd.pipeline.run.state", string(pipelineRunState)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrPipelineResult returns an optional attribute for the
// "cicd.pipeline.result" semantic convention. It represents the result of a
// pipeline run.
func (PipelineRunDuration) AttrPipelineResult(val PipelineResultAttr) attribute.KeyValue {
	return attribute.String("cicd.pipeline.result", string(val))
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (PipelineRunDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// PipelineRunErrors is an instrument used to record metric values conforming to
// the "cicd.pipeline.run.errors" semantic conventions. It represents the number
// of errors encountered in pipeline runs (eg. compile, test failures).
type PipelineRunErrors struct {
	metric.Int64Counter
}

// NewPipelineRunErrors returns a new PipelineRunErrors instrument.
func NewPipelineRunErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (PipelineRunErrors, error) {
	// Check if the meter is nil.
	if m == nil {
		return PipelineRunErrors{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"cicd.pipeline.run.errors",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of errors encountered in pipeline runs (eg. compile, test failures)."),
			metric.WithUnit("{error}"),
		}, opt...)...,
	)
	if err != nil {
	    return PipelineRunErrors{noop.Int64Counter{}}, err
	}
	return PipelineRunErrors{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PipelineRunErrors) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (PipelineRunErrors) Name() string {
	return "cicd.pipeline.run.errors"
}

// Unit returns the semantic convention unit of the instrument
func (PipelineRunErrors) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (PipelineRunErrors) Description() string {
	return "The number of errors encountered in pipeline runs (eg. compile, test failures)."
}

// Add adds incr to the existing count.
//
// The pipelineName is the the human readable name of the pipeline within a CI/CD
// system.
//
// The errorType is the describes a class of error the operation ended with.
//
// There might be errors in a pipeline run that are non fatal (eg. they are
// suppressed) or in a parallel stage multiple stages could have a fatal error.
// This means that this error count might not be the same as the count of metric
// `cicd.pipeline.run.duration` with run result `failure`.
func (m PipelineRunErrors) Add(
	ctx context.Context,
	incr int64,
	pipelineName string,
	errorType ErrorTypeAttr,
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
			append(
				attrs,
				attribute.String("cicd.pipeline.name", pipelineName),
				attribute.String("error.type", string(errorType)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// SystemErrors is an instrument used to record metric values conforming to the
// "cicd.system.errors" semantic conventions. It represents the number of errors
// in a component of the CICD system (eg. controller, scheduler, agent).
type SystemErrors struct {
	metric.Int64Counter
}

// NewSystemErrors returns a new SystemErrors instrument.
func NewSystemErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (SystemErrors, error) {
	// Check if the meter is nil.
	if m == nil {
		return SystemErrors{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"cicd.system.errors",
		append([]metric.Int64CounterOption{
			metric.WithDescription("The number of errors in a component of the CICD system (eg. controller, scheduler, agent)."),
			metric.WithUnit("{error}"),
		}, opt...)...,
	)
	if err != nil {
	    return SystemErrors{noop.Int64Counter{}}, err
	}
	return SystemErrors{i}, nil
}

// Inst returns the underlying metric instrument.
func (m SystemErrors) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (SystemErrors) Name() string {
	return "cicd.system.errors"
}

// Unit returns the semantic convention unit of the instrument
func (SystemErrors) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (SystemErrors) Description() string {
	return "The number of errors in a component of the CICD system (eg. controller, scheduler, agent)."
}

// Add adds incr to the existing count.
//
// The systemComponent is the the name of a component of the CICD system.
//
// The errorType is the describes a class of error the operation ended with.
//
// Errors in pipeline run execution are explicitly excluded. Ie a test failure is
// not counted in this metric.
func (m SystemErrors) Add(
	ctx context.Context,
	incr int64,
	systemComponent string,
	errorType ErrorTypeAttr,
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
			append(
				attrs,
				attribute.String("cicd.system.component", systemComponent),
				attribute.String("error.type", string(errorType)),
			)...,
		),
	)

	m.Int64Counter.Add(ctx, incr, *o...)
}

// WorkerCount is an instrument used to record metric values conforming to the
// "cicd.worker.count" semantic conventions. It represents the number of workers
// on the CICD system by state.
type WorkerCount struct {
	metric.Int64UpDownCounter
}

// NewWorkerCount returns a new WorkerCount instrument.
func NewWorkerCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (WorkerCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return WorkerCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"cicd.worker.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of workers on the CICD system by state."),
			metric.WithUnit("{count}"),
		}, opt...)...,
	)
	if err != nil {
	    return WorkerCount{noop.Int64UpDownCounter{}}, err
	}
	return WorkerCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m WorkerCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (WorkerCount) Name() string {
	return "cicd.worker.count"
}

// Unit returns the semantic convention unit of the instrument
func (WorkerCount) Unit() string {
	return "{count}"
}

// Description returns the semantic convention description of the instrument
func (WorkerCount) Description() string {
	return "The number of workers on the CICD system by state."
}

// Add adds incr to the existing count.
//
// The workerState is the the state of a CICD worker / agent.
func (m WorkerCount) Add(
	ctx context.Context,
	incr int64,
	workerState WorkerStateAttr,
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
			append(
				attrs,
				attribute.String("cicd.worker.state", string(workerState)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}