// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/cicd"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	inst metric.Int64UpDownCounter
}

// NewPipelineRunActive returns a new PipelineRunActive instrument.
func NewPipelineRunActive(m metric.Meter) (PipelineRunActive, error) {
	i, err := m.Int64UpDownCounter(
	    "cicd.pipeline.run.active",
	    metric.WithDescription("The number of pipeline runs currently active in the system by state."),
	    metric.WithUnit("{run}"),
	)
	if err != nil {
	    return PipelineRunActive{}, err
	}
	return PipelineRunActive{i}, nil
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
// The cicdPipelineName is the the human readable name of the pipeline within a
// CI/CD system.
//
// The cicdPipelineRunState is the the pipeline run goes through these states
// during its lifecycle.
func (m PipelineRunActive) Add(
    ctx context.Context,
    incr int64,
	pipelineName string,
	pipelineRunState PipelineRunStateAttr,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("cicd.pipeline.name", pipelineName),
			attribute.String("cicd.pipeline.run.state", string(pipelineRunState)),

		),
	)
}

// PipelineRunDuration is an instrument used to record metric values conforming
// to the "cicd.pipeline.run.duration" semantic conventions. It represents the
// duration of a pipeline run grouped by pipeline, state and result.
type PipelineRunDuration struct {
	inst metric.Float64Histogram
}

// NewPipelineRunDuration returns a new PipelineRunDuration instrument.
func NewPipelineRunDuration(m metric.Meter) (PipelineRunDuration, error) {
	i, err := m.Float64Histogram(
	    "cicd.pipeline.run.duration",
	    metric.WithDescription("Duration of a pipeline run grouped by pipeline, state and result."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return PipelineRunDuration{}, err
	}
	return PipelineRunDuration{i}, nil
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
// The cicdPipelineName is the the human readable name of the pipeline within a
// CI/CD system.
//
// The cicdPipelineRunState is the the pipeline run goes through these states
// during its lifecycle.
//
// All additional attrs passed are included in the recorded value.
func (m PipelineRunDuration) Record(
    ctx context.Context,
    val float64,
	pipelineName string,
	pipelineRunState PipelineRunStateAttr,
	attrs ...PipelineRunDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("cicd.pipeline.name", pipelineName),
				attribute.String("cicd.pipeline.run.state", string(pipelineRunState)),
			)...,
		),
	)
}

func (m PipelineRunDuration) conv(in []PipelineRunDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.pipelineRunDurationAttr()
	}
	return out
}

// PipelineRunDurationAttr is an optional attribute for the PipelineRunDuration
// instrument.
type PipelineRunDurationAttr interface {
    pipelineRunDurationAttr() attribute.KeyValue
}

type pipelineRunDurationAttr struct {
	kv attribute.KeyValue
}

func (a pipelineRunDurationAttr) pipelineRunDurationAttr() attribute.KeyValue {
    return a.kv
}

// PipelineResult returns an optional attribute for the "cicd.pipeline.result"
// semantic convention. It represents the result of a pipeline run.
func (PipelineRunDuration) PipelineResultAttr(val PipelineResultAttr) PipelineRunDurationAttr {
	return pipelineRunDurationAttr{kv: attribute.String("cicd.pipeline.result", string(val))}
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (PipelineRunDuration) ErrorTypeAttr(val ErrorTypeAttr) PipelineRunDurationAttr {
	return pipelineRunDurationAttr{kv: attribute.String("error.type", string(val))}
}

// PipelineRunErrors is an instrument used to record metric values conforming to
// the "cicd.pipeline.run.errors" semantic conventions. It represents the number
// of errors encountered in pipeline runs (eg. compile, test failures).
type PipelineRunErrors struct {
	inst metric.Int64Counter
}

// NewPipelineRunErrors returns a new PipelineRunErrors instrument.
func NewPipelineRunErrors(m metric.Meter) (PipelineRunErrors, error) {
	i, err := m.Int64Counter(
	    "cicd.pipeline.run.errors",
	    metric.WithDescription("The number of errors encountered in pipeline runs (eg. compile, test failures)."),
	    metric.WithUnit("{error}"),
	)
	if err != nil {
	    return PipelineRunErrors{}, err
	}
	return PipelineRunErrors{i}, nil
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
// The cicdPipelineName is the the human readable name of the pipeline within a
// CI/CD system.
//
// The errorType is the describes a class of error the operation ended with.
func (m PipelineRunErrors) Add(
    ctx context.Context,
    incr int64,
	pipelineName string,
	errorType ErrorTypeAttr,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("cicd.pipeline.name", pipelineName),
			attribute.String("error.type", string(errorType)),

		),
	)
}

// SystemErrors is an instrument used to record metric values conforming to the
// "cicd.system.errors" semantic conventions. It represents the number of errors
// in a component of the CICD system (eg. controller, scheduler, agent).
type SystemErrors struct {
	inst metric.Int64Counter
}

// NewSystemErrors returns a new SystemErrors instrument.
func NewSystemErrors(m metric.Meter) (SystemErrors, error) {
	i, err := m.Int64Counter(
	    "cicd.system.errors",
	    metric.WithDescription("The number of errors in a component of the CICD system (eg. controller, scheduler, agent)."),
	    metric.WithUnit("{error}"),
	)
	if err != nil {
	    return SystemErrors{}, err
	}
	return SystemErrors{i}, nil
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
// The cicdSystemComponent is the the name of a component of the CICD system.
//
// The errorType is the describes a class of error the operation ended with.
func (m SystemErrors) Add(
    ctx context.Context,
    incr int64,
	systemComponent string,
	errorType ErrorTypeAttr,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("cicd.system.component", systemComponent),
			attribute.String("error.type", string(errorType)),

		),
	)
}

// WorkerCount is an instrument used to record metric values conforming to the
// "cicd.worker.count" semantic conventions. It represents the number of workers
// on the CICD system by state.
type WorkerCount struct {
	inst metric.Int64UpDownCounter
}

// NewWorkerCount returns a new WorkerCount instrument.
func NewWorkerCount(m metric.Meter) (WorkerCount, error) {
	i, err := m.Int64UpDownCounter(
	    "cicd.worker.count",
	    metric.WithDescription("The number of workers on the CICD system by state."),
	    metric.WithUnit("{count}"),
	)
	if err != nil {
	    return WorkerCount{}, err
	}
	return WorkerCount{i}, nil
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
// The cicdWorkerState is the the state of a CICD worker / agent.
func (m WorkerCount) Add(
    ctx context.Context,
    incr int64,
	workerState WorkerStateAttr,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("cicd.worker.state", string(workerState)),

		),
	)
}