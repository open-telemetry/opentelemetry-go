// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/k8s"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// NamespacePhaseAttr is an attribute conforming to the k8s.namespace.phase
// semantic conventions. It represents the phase of the K8s namespace.
type NamespacePhaseAttr string

var (
	// NamespacePhaseActive is the active namespace phase as described by [K8s API]
	// .
	//
	// [K8s API]: https://pkg.go.dev/k8s.io/api@v0.31.3/core/v1#NamespacePhase
	NamespacePhaseActive NamespacePhaseAttr = "active"
	// NamespacePhaseTerminating is the terminating namespace phase as described by
	// [K8s API].
	//
	// [K8s API]: https://pkg.go.dev/k8s.io/api@v0.31.3/core/v1#NamespacePhase
	NamespacePhaseTerminating NamespacePhaseAttr = "terminating"
)

// NetworkIoDirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the network IO operation direction.
type NetworkIoDirectionAttr string

var (
	// NetworkIoDirectionTransmit is the none.
	NetworkIoDirectionTransmit NetworkIoDirectionAttr = "transmit"
	// NetworkIoDirectionReceive is the none.
	NetworkIoDirectionReceive NetworkIoDirectionAttr = "receive"
)

// K8SCronJobActiveJobs is an instrument used to record metric values conforming
// to the "k8s.cronjob.active_jobs" semantic conventions. It represents the
// number of actively running jobs for a cronjob.
type CronJobActiveJobs struct {
	inst metric.Int64UpDownCounter
}

// NewCronJobActiveJobs returns a new CronJobActiveJobs instrument.
func NewCronJobActiveJobs(m metric.Meter) (CronJobActiveJobs, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.cronjob.active_jobs",
	    metric.WithDescription("The number of actively running jobs for a cronjob"),
	    metric.WithUnit("{job}"),
	)
	if err != nil {
	    return CronJobActiveJobs{}, err
	}
	return CronJobActiveJobs{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (CronJobActiveJobs) Name() string {
	return "k8s.cronjob.active_jobs"
}

// Unit returns the semantic convention unit of the instrument
func (CronJobActiveJobs) Unit() string {
	return "{job}"
}

// Description returns the semantic convention description of the instrument
func (CronJobActiveJobs) Description() string {
	return "The number of actively running jobs for a cronjob"
}

func (m CronJobActiveJobs) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SDaemonSetCurrentScheduledNodes is an instrument used to record metric
// values conforming to the "k8s.daemonset.current_scheduled_nodes" semantic
// conventions. It represents the number of nodes that are running at least 1
// daemon pod and are supposed to run the daemon pod.
type DaemonSetCurrentScheduledNodes struct {
	inst metric.Int64UpDownCounter
}

// NewDaemonSetCurrentScheduledNodes returns a new DaemonSetCurrentScheduledNodes
// instrument.
func NewDaemonSetCurrentScheduledNodes(m metric.Meter) (DaemonSetCurrentScheduledNodes, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.daemonset.current_scheduled_nodes",
	    metric.WithDescription("Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod"),
	    metric.WithUnit("{node}"),
	)
	if err != nil {
	    return DaemonSetCurrentScheduledNodes{}, err
	}
	return DaemonSetCurrentScheduledNodes{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetCurrentScheduledNodes) Name() string {
	return "k8s.daemonset.current_scheduled_nodes"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetCurrentScheduledNodes) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetCurrentScheduledNodes) Description() string {
	return "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod"
}

func (m DaemonSetCurrentScheduledNodes) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SDaemonSetDesiredScheduledNodes is an instrument used to record metric
// values conforming to the "k8s.daemonset.desired_scheduled_nodes" semantic
// conventions. It represents the number of nodes that should be running the
// daemon pod (including nodes currently running the daemon pod).
type DaemonSetDesiredScheduledNodes struct {
	inst metric.Int64UpDownCounter
}

// NewDaemonSetDesiredScheduledNodes returns a new DaemonSetDesiredScheduledNodes
// instrument.
func NewDaemonSetDesiredScheduledNodes(m metric.Meter) (DaemonSetDesiredScheduledNodes, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.daemonset.desired_scheduled_nodes",
	    metric.WithDescription("Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)"),
	    metric.WithUnit("{node}"),
	)
	if err != nil {
	    return DaemonSetDesiredScheduledNodes{}, err
	}
	return DaemonSetDesiredScheduledNodes{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetDesiredScheduledNodes) Name() string {
	return "k8s.daemonset.desired_scheduled_nodes"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetDesiredScheduledNodes) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetDesiredScheduledNodes) Description() string {
	return "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)"
}

func (m DaemonSetDesiredScheduledNodes) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SDaemonSetMisscheduledNodes is an instrument used to record metric values
// conforming to the "k8s.daemonset.misscheduled_nodes" semantic conventions. It
// represents the number of nodes that are running the daemon pod, but are not
// supposed to run the daemon pod.
type DaemonSetMisscheduledNodes struct {
	inst metric.Int64UpDownCounter
}

// NewDaemonSetMisscheduledNodes returns a new DaemonSetMisscheduledNodes
// instrument.
func NewDaemonSetMisscheduledNodes(m metric.Meter) (DaemonSetMisscheduledNodes, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.daemonset.misscheduled_nodes",
	    metric.WithDescription("Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod"),
	    metric.WithUnit("{node}"),
	)
	if err != nil {
	    return DaemonSetMisscheduledNodes{}, err
	}
	return DaemonSetMisscheduledNodes{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetMisscheduledNodes) Name() string {
	return "k8s.daemonset.misscheduled_nodes"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetMisscheduledNodes) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetMisscheduledNodes) Description() string {
	return "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod"
}

func (m DaemonSetMisscheduledNodes) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SDaemonSetReadyNodes is an instrument used to record metric values
// conforming to the "k8s.daemonset.ready_nodes" semantic conventions. It
// represents the number of nodes that should be running the daemon pod and have
// one or more of the daemon pod running and ready.
type DaemonSetReadyNodes struct {
	inst metric.Int64UpDownCounter
}

// NewDaemonSetReadyNodes returns a new DaemonSetReadyNodes instrument.
func NewDaemonSetReadyNodes(m metric.Meter) (DaemonSetReadyNodes, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.daemonset.ready_nodes",
	    metric.WithDescription("Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready"),
	    metric.WithUnit("{node}"),
	)
	if err != nil {
	    return DaemonSetReadyNodes{}, err
	}
	return DaemonSetReadyNodes{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetReadyNodes) Name() string {
	return "k8s.daemonset.ready_nodes"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetReadyNodes) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetReadyNodes) Description() string {
	return "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready"
}

func (m DaemonSetReadyNodes) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SDeploymentAvailablePods is an instrument used to record metric values
// conforming to the "k8s.deployment.available_pods" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this deployment.
type DeploymentAvailablePods struct {
	inst metric.Int64UpDownCounter
}

// NewDeploymentAvailablePods returns a new DeploymentAvailablePods instrument.
func NewDeploymentAvailablePods(m metric.Meter) (DeploymentAvailablePods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.deployment.available_pods",
	    metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return DeploymentAvailablePods{}, err
	}
	return DeploymentAvailablePods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DeploymentAvailablePods) Name() string {
	return "k8s.deployment.available_pods"
}

// Unit returns the semantic convention unit of the instrument
func (DeploymentAvailablePods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (DeploymentAvailablePods) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment"
}

func (m DeploymentAvailablePods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SDeploymentDesiredPods is an instrument used to record metric values
// conforming to the "k8s.deployment.desired_pods" semantic conventions. It
// represents the number of desired replica pods in this deployment.
type DeploymentDesiredPods struct {
	inst metric.Int64UpDownCounter
}

// NewDeploymentDesiredPods returns a new DeploymentDesiredPods instrument.
func NewDeploymentDesiredPods(m metric.Meter) (DeploymentDesiredPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.deployment.desired_pods",
	    metric.WithDescription("Number of desired replica pods in this deployment"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return DeploymentDesiredPods{}, err
	}
	return DeploymentDesiredPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (DeploymentDesiredPods) Name() string {
	return "k8s.deployment.desired_pods"
}

// Unit returns the semantic convention unit of the instrument
func (DeploymentDesiredPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (DeploymentDesiredPods) Description() string {
	return "Number of desired replica pods in this deployment"
}

func (m DeploymentDesiredPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SHpaCurrentPods is an instrument used to record metric values conforming to
// the "k8s.hpa.current_pods" semantic conventions. It represents the current
// number of replica pods managed by this horizontal pod autoscaler, as last seen
// by the autoscaler.
type HpaCurrentPods struct {
	inst metric.Int64UpDownCounter
}

// NewHpaCurrentPods returns a new HpaCurrentPods instrument.
func NewHpaCurrentPods(m metric.Meter) (HpaCurrentPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.hpa.current_pods",
	    metric.WithDescription("Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return HpaCurrentPods{}, err
	}
	return HpaCurrentPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HpaCurrentPods) Name() string {
	return "k8s.hpa.current_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HpaCurrentPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HpaCurrentPods) Description() string {
	return "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler"
}

func (m HpaCurrentPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SHpaDesiredPods is an instrument used to record metric values conforming to
// the "k8s.hpa.desired_pods" semantic conventions. It represents the desired
// number of replica pods managed by this horizontal pod autoscaler, as last
// calculated by the autoscaler.
type HpaDesiredPods struct {
	inst metric.Int64UpDownCounter
}

// NewHpaDesiredPods returns a new HpaDesiredPods instrument.
func NewHpaDesiredPods(m metric.Meter) (HpaDesiredPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.hpa.desired_pods",
	    metric.WithDescription("Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return HpaDesiredPods{}, err
	}
	return HpaDesiredPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HpaDesiredPods) Name() string {
	return "k8s.hpa.desired_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HpaDesiredPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HpaDesiredPods) Description() string {
	return "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler"
}

func (m HpaDesiredPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SHpaMaxPods is an instrument used to record metric values conforming to the
// "k8s.hpa.max_pods" semantic conventions. It represents the upper limit for the
// number of replica pods to which the autoscaler can scale up.
type HpaMaxPods struct {
	inst metric.Int64UpDownCounter
}

// NewHpaMaxPods returns a new HpaMaxPods instrument.
func NewHpaMaxPods(m metric.Meter) (HpaMaxPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.hpa.max_pods",
	    metric.WithDescription("The upper limit for the number of replica pods to which the autoscaler can scale up"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return HpaMaxPods{}, err
	}
	return HpaMaxPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HpaMaxPods) Name() string {
	return "k8s.hpa.max_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HpaMaxPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HpaMaxPods) Description() string {
	return "The upper limit for the number of replica pods to which the autoscaler can scale up"
}

func (m HpaMaxPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SHpaMinPods is an instrument used to record metric values conforming to the
// "k8s.hpa.min_pods" semantic conventions. It represents the lower limit for the
// number of replica pods to which the autoscaler can scale down.
type HpaMinPods struct {
	inst metric.Int64UpDownCounter
}

// NewHpaMinPods returns a new HpaMinPods instrument.
func NewHpaMinPods(m metric.Meter) (HpaMinPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.hpa.min_pods",
	    metric.WithDescription("The lower limit for the number of replica pods to which the autoscaler can scale down"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return HpaMinPods{}, err
	}
	return HpaMinPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HpaMinPods) Name() string {
	return "k8s.hpa.min_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HpaMinPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HpaMinPods) Description() string {
	return "The lower limit for the number of replica pods to which the autoscaler can scale down"
}

func (m HpaMinPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SJobActivePods is an instrument used to record metric values conforming to
// the "k8s.job.active_pods" semantic conventions. It represents the number of
// pending and actively running pods for a job.
type JobActivePods struct {
	inst metric.Int64UpDownCounter
}

// NewJobActivePods returns a new JobActivePods instrument.
func NewJobActivePods(m metric.Meter) (JobActivePods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.job.active_pods",
	    metric.WithDescription("The number of pending and actively running pods for a job"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return JobActivePods{}, err
	}
	return JobActivePods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (JobActivePods) Name() string {
	return "k8s.job.active_pods"
}

// Unit returns the semantic convention unit of the instrument
func (JobActivePods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobActivePods) Description() string {
	return "The number of pending and actively running pods for a job"
}

func (m JobActivePods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SJobDesiredSuccessfulPods is an instrument used to record metric values
// conforming to the "k8s.job.desired_successful_pods" semantic conventions. It
// represents the desired number of successfully finished pods the job should be
// run with.
type JobDesiredSuccessfulPods struct {
	inst metric.Int64UpDownCounter
}

// NewJobDesiredSuccessfulPods returns a new JobDesiredSuccessfulPods instrument.
func NewJobDesiredSuccessfulPods(m metric.Meter) (JobDesiredSuccessfulPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.job.desired_successful_pods",
	    metric.WithDescription("The desired number of successfully finished pods the job should be run with"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return JobDesiredSuccessfulPods{}, err
	}
	return JobDesiredSuccessfulPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (JobDesiredSuccessfulPods) Name() string {
	return "k8s.job.desired_successful_pods"
}

// Unit returns the semantic convention unit of the instrument
func (JobDesiredSuccessfulPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobDesiredSuccessfulPods) Description() string {
	return "The desired number of successfully finished pods the job should be run with"
}

func (m JobDesiredSuccessfulPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SJobFailedPods is an instrument used to record metric values conforming to
// the "k8s.job.failed_pods" semantic conventions. It represents the number of
// pods which reached phase Failed for a job.
type JobFailedPods struct {
	inst metric.Int64UpDownCounter
}

// NewJobFailedPods returns a new JobFailedPods instrument.
func NewJobFailedPods(m metric.Meter) (JobFailedPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.job.failed_pods",
	    metric.WithDescription("The number of pods which reached phase Failed for a job"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return JobFailedPods{}, err
	}
	return JobFailedPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (JobFailedPods) Name() string {
	return "k8s.job.failed_pods"
}

// Unit returns the semantic convention unit of the instrument
func (JobFailedPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobFailedPods) Description() string {
	return "The number of pods which reached phase Failed for a job"
}

func (m JobFailedPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SJobMaxParallelPods is an instrument used to record metric values conforming
// to the "k8s.job.max_parallel_pods" semantic conventions. It represents the max
// desired number of pods the job should run at any given time.
type JobMaxParallelPods struct {
	inst metric.Int64UpDownCounter
}

// NewJobMaxParallelPods returns a new JobMaxParallelPods instrument.
func NewJobMaxParallelPods(m metric.Meter) (JobMaxParallelPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.job.max_parallel_pods",
	    metric.WithDescription("The max desired number of pods the job should run at any given time"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return JobMaxParallelPods{}, err
	}
	return JobMaxParallelPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (JobMaxParallelPods) Name() string {
	return "k8s.job.max_parallel_pods"
}

// Unit returns the semantic convention unit of the instrument
func (JobMaxParallelPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobMaxParallelPods) Description() string {
	return "The max desired number of pods the job should run at any given time"
}

func (m JobMaxParallelPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SJobSuccessfulPods is an instrument used to record metric values conforming
// to the "k8s.job.successful_pods" semantic conventions. It represents the
// number of pods which reached phase Succeeded for a job.
type JobSuccessfulPods struct {
	inst metric.Int64UpDownCounter
}

// NewJobSuccessfulPods returns a new JobSuccessfulPods instrument.
func NewJobSuccessfulPods(m metric.Meter) (JobSuccessfulPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.job.successful_pods",
	    metric.WithDescription("The number of pods which reached phase Succeeded for a job"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return JobSuccessfulPods{}, err
	}
	return JobSuccessfulPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (JobSuccessfulPods) Name() string {
	return "k8s.job.successful_pods"
}

// Unit returns the semantic convention unit of the instrument
func (JobSuccessfulPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobSuccessfulPods) Description() string {
	return "The number of pods which reached phase Succeeded for a job"
}

func (m JobSuccessfulPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SNamespacePhase is an instrument used to record metric values conforming to
// the "k8s.namespace.phase" semantic conventions. It represents the describes
// number of K8s namespaces that are currently in a given phase.
type NamespacePhase struct {
	inst metric.Int64UpDownCounter
}

// NewNamespacePhase returns a new NamespacePhase instrument.
func NewNamespacePhase(m metric.Meter) (NamespacePhase, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.namespace.phase",
	    metric.WithDescription("Describes number of K8s namespaces that are currently in a given phase."),
	    metric.WithUnit("{namespace}"),
	)
	if err != nil {
	    return NamespacePhase{}, err
	}
	return NamespacePhase{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NamespacePhase) Name() string {
	return "k8s.namespace.phase"
}

// Unit returns the semantic convention unit of the instrument
func (NamespacePhase) Unit() string {
	return "{namespace}"
}

// Description returns the semantic convention description of the instrument
func (NamespacePhase) Description() string {
	return "Describes number of K8s namespaces that are currently in a given phase."
}

// Add adds incr to the existing count.
//
// The k8sNamespacePhase is the the phase of the K8s namespace.
func (m NamespacePhase) Add(
    ctx context.Context,
    incr int64,
	namespacePhase NamespacePhaseAttr,

) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(			attribute.String("k8s.namespace.phase", string(namespacePhase)),

		),
	)
}

// K8SNodeCPUTime is an instrument used to record metric values conforming to the
// "k8s.node.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type NodeCPUTime struct {
	inst metric.Float64Counter
}

// NewNodeCPUTime returns a new NodeCPUTime instrument.
func NewNodeCPUTime(m metric.Meter) (NodeCPUTime, error) {
	i, err := m.Float64Counter(
	    "k8s.node.cpu.time",
	    metric.WithDescription("Total CPU time consumed"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return NodeCPUTime{}, err
	}
	return NodeCPUTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NodeCPUTime) Name() string {
	return "k8s.node.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (NodeCPUTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (NodeCPUTime) Description() string {
	return "Total CPU time consumed"
}

func (m NodeCPUTime) Add(ctx context.Context, incr float64) {
    m.inst.Add(ctx, incr)
}

// K8SNodeCPUUsage is an instrument used to record metric values conforming to
// the "k8s.node.cpu.usage" semantic conventions. It represents the node's CPU
// usage, measured in cpus. Range from 0 to the number of allocatable CPUs.
type NodeCPUUsage struct {
	inst metric.Int64Gauge
}

// NewNodeCPUUsage returns a new NodeCPUUsage instrument.
func NewNodeCPUUsage(m metric.Meter) (NodeCPUUsage, error) {
	i, err := m.Int64Gauge(
	    "k8s.node.cpu.usage",
	    metric.WithDescription("Node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"),
	    metric.WithUnit("{cpu}"),
	)
	if err != nil {
	    return NodeCPUUsage{}, err
	}
	return NodeCPUUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NodeCPUUsage) Name() string {
	return "k8s.node.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeCPUUsage) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (NodeCPUUsage) Description() string {
	return "Node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"
}

func (m NodeCPUUsage) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// K8SNodeMemoryUsage is an instrument used to record metric values conforming to
// the "k8s.node.memory.usage" semantic conventions. It represents the memory
// usage of the Node.
type NodeMemoryUsage struct {
	inst metric.Int64Gauge
}

// NewNodeMemoryUsage returns a new NodeMemoryUsage instrument.
func NewNodeMemoryUsage(m metric.Meter) (NodeMemoryUsage, error) {
	i, err := m.Int64Gauge(
	    "k8s.node.memory.usage",
	    metric.WithDescription("Memory usage of the Node"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return NodeMemoryUsage{}, err
	}
	return NodeMemoryUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryUsage) Name() string {
	return "k8s.node.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryUsage) Description() string {
	return "Memory usage of the Node"
}

func (m NodeMemoryUsage) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// K8SNodeNetworkErrors is an instrument used to record metric values conforming
// to the "k8s.node.network.errors" semantic conventions. It represents the node
// network errors.
type NodeNetworkErrors struct {
	inst metric.Int64Counter
}

// NewNodeNetworkErrors returns a new NodeNetworkErrors instrument.
func NewNodeNetworkErrors(m metric.Meter) (NodeNetworkErrors, error) {
	i, err := m.Int64Counter(
	    "k8s.node.network.errors",
	    metric.WithDescription("Node network errors"),
	    metric.WithUnit("{error}"),
	)
	if err != nil {
	    return NodeNetworkErrors{}, err
	}
	return NodeNetworkErrors{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NodeNetworkErrors) Name() string {
	return "k8s.node.network.errors"
}

// Unit returns the semantic convention unit of the instrument
func (NodeNetworkErrors) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (NodeNetworkErrors) Description() string {
	return "Node network errors"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NodeNetworkErrors) Add(
    ctx context.Context,
    incr int64,
	attrs ...NodeNetworkErrorsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NodeNetworkErrors) conv(in []NodeNetworkErrorsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.nodeNetworkErrorsAttr()
	}
	return out
}

// NodeNetworkErrorsAttr is an optional attribute for the NodeNetworkErrors
// instrument.
type NodeNetworkErrorsAttr interface {
    nodeNetworkErrorsAttr() attribute.KeyValue
}

type nodeNetworkErrorsAttr struct {
	kv attribute.KeyValue
}

func (a nodeNetworkErrorsAttr) nodeNetworkErrorsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NodeNetworkErrors) NetworkInterfaceName(val string) NodeNetworkErrorsAttr {
	return nodeNetworkErrorsAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NodeNetworkErrors) NetworkIoDirection(val NetworkIoDirectionAttr) NodeNetworkErrorsAttr {
	return nodeNetworkErrorsAttr{kv: attribute.String("network.io.direction", string(val))}
}

// K8SNodeNetworkIo is an instrument used to record metric values conforming to
// the "k8s.node.network.io" semantic conventions. It represents the network
// bytes for the Node.
type NodeNetworkIo struct {
	inst metric.Int64Counter
}

// NewNodeNetworkIo returns a new NodeNetworkIo instrument.
func NewNodeNetworkIo(m metric.Meter) (NodeNetworkIo, error) {
	i, err := m.Int64Counter(
	    "k8s.node.network.io",
	    metric.WithDescription("Network bytes for the Node"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return NodeNetworkIo{}, err
	}
	return NodeNetworkIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NodeNetworkIo) Name() string {
	return "k8s.node.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NodeNetworkIo) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeNetworkIo) Description() string {
	return "Network bytes for the Node"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NodeNetworkIo) Add(
    ctx context.Context,
    incr int64,
	attrs ...NodeNetworkIoAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m NodeNetworkIo) conv(in []NodeNetworkIoAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.nodeNetworkIoAttr()
	}
	return out
}

// NodeNetworkIoAttr is an optional attribute for the NodeNetworkIo instrument.
type NodeNetworkIoAttr interface {
    nodeNetworkIoAttr() attribute.KeyValue
}

type nodeNetworkIoAttr struct {
	kv attribute.KeyValue
}

func (a nodeNetworkIoAttr) nodeNetworkIoAttr() attribute.KeyValue {
    return a.kv
}

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NodeNetworkIo) NetworkInterfaceName(val string) NodeNetworkIoAttr {
	return nodeNetworkIoAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NodeNetworkIo) NetworkIoDirection(val NetworkIoDirectionAttr) NodeNetworkIoAttr {
	return nodeNetworkIoAttr{kv: attribute.String("network.io.direction", string(val))}
}

// K8SNodeUptime is an instrument used to record metric values conforming to the
// "k8s.node.uptime" semantic conventions. It represents the time the Node has
// been running.
type NodeUptime struct {
	inst metric.Float64Gauge
}

// NewNodeUptime returns a new NodeUptime instrument.
func NewNodeUptime(m metric.Meter) (NodeUptime, error) {
	i, err := m.Float64Gauge(
	    "k8s.node.uptime",
	    metric.WithDescription("The time the Node has been running"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return NodeUptime{}, err
	}
	return NodeUptime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (NodeUptime) Name() string {
	return "k8s.node.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (NodeUptime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (NodeUptime) Description() string {
	return "The time the Node has been running"
}

func (m NodeUptime) Record(ctx context.Context, val float64) {
    m.inst.Record(ctx, val)
}

// K8SPodCPUTime is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type PodCPUTime struct {
	inst metric.Float64Counter
}

// NewPodCPUTime returns a new PodCPUTime instrument.
func NewPodCPUTime(m metric.Meter) (PodCPUTime, error) {
	i, err := m.Float64Counter(
	    "k8s.pod.cpu.time",
	    metric.WithDescription("Total CPU time consumed"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return PodCPUTime{}, err
	}
	return PodCPUTime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PodCPUTime) Name() string {
	return "k8s.pod.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPUTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (PodCPUTime) Description() string {
	return "Total CPU time consumed"
}

func (m PodCPUTime) Add(ctx context.Context, incr float64) {
    m.inst.Add(ctx, incr)
}

// K8SPodCPUUsage is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.usage" semantic conventions. It represents the pod's CPU usage,
// measured in cpus. Range from 0 to the number of allocatable CPUs.
type PodCPUUsage struct {
	inst metric.Int64Gauge
}

// NewPodCPUUsage returns a new PodCPUUsage instrument.
func NewPodCPUUsage(m metric.Meter) (PodCPUUsage, error) {
	i, err := m.Int64Gauge(
	    "k8s.pod.cpu.usage",
	    metric.WithDescription("Pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"),
	    metric.WithUnit("{cpu}"),
	)
	if err != nil {
	    return PodCPUUsage{}, err
	}
	return PodCPUUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PodCPUUsage) Name() string {
	return "k8s.pod.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPUUsage) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPUUsage) Description() string {
	return "Pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"
}

func (m PodCPUUsage) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// K8SPodMemoryUsage is an instrument used to record metric values conforming to
// the "k8s.pod.memory.usage" semantic conventions. It represents the memory
// usage of the Pod.
type PodMemoryUsage struct {
	inst metric.Int64Gauge
}

// NewPodMemoryUsage returns a new PodMemoryUsage instrument.
func NewPodMemoryUsage(m metric.Meter) (PodMemoryUsage, error) {
	i, err := m.Int64Gauge(
	    "k8s.pod.memory.usage",
	    metric.WithDescription("Memory usage of the Pod"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return PodMemoryUsage{}, err
	}
	return PodMemoryUsage{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryUsage) Name() string {
	return "k8s.pod.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryUsage) Description() string {
	return "Memory usage of the Pod"
}

func (m PodMemoryUsage) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// K8SPodNetworkErrors is an instrument used to record metric values conforming
// to the "k8s.pod.network.errors" semantic conventions. It represents the pod
// network errors.
type PodNetworkErrors struct {
	inst metric.Int64Counter
}

// NewPodNetworkErrors returns a new PodNetworkErrors instrument.
func NewPodNetworkErrors(m metric.Meter) (PodNetworkErrors, error) {
	i, err := m.Int64Counter(
	    "k8s.pod.network.errors",
	    metric.WithDescription("Pod network errors"),
	    metric.WithUnit("{error}"),
	)
	if err != nil {
	    return PodNetworkErrors{}, err
	}
	return PodNetworkErrors{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PodNetworkErrors) Name() string {
	return "k8s.pod.network.errors"
}

// Unit returns the semantic convention unit of the instrument
func (PodNetworkErrors) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (PodNetworkErrors) Description() string {
	return "Pod network errors"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m PodNetworkErrors) Add(
    ctx context.Context,
    incr int64,
	attrs ...PodNetworkErrorsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m PodNetworkErrors) conv(in []PodNetworkErrorsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.podNetworkErrorsAttr()
	}
	return out
}

// PodNetworkErrorsAttr is an optional attribute for the PodNetworkErrors
// instrument.
type PodNetworkErrorsAttr interface {
    podNetworkErrorsAttr() attribute.KeyValue
}

type podNetworkErrorsAttr struct {
	kv attribute.KeyValue
}

func (a podNetworkErrorsAttr) podNetworkErrorsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (PodNetworkErrors) NetworkInterfaceName(val string) PodNetworkErrorsAttr {
	return podNetworkErrorsAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (PodNetworkErrors) NetworkIoDirection(val NetworkIoDirectionAttr) PodNetworkErrorsAttr {
	return podNetworkErrorsAttr{kv: attribute.String("network.io.direction", string(val))}
}

// K8SPodNetworkIo is an instrument used to record metric values conforming to
// the "k8s.pod.network.io" semantic conventions. It represents the network bytes
// for the Pod.
type PodNetworkIo struct {
	inst metric.Int64Counter
}

// NewPodNetworkIo returns a new PodNetworkIo instrument.
func NewPodNetworkIo(m metric.Meter) (PodNetworkIo, error) {
	i, err := m.Int64Counter(
	    "k8s.pod.network.io",
	    metric.WithDescription("Network bytes for the Pod"),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return PodNetworkIo{}, err
	}
	return PodNetworkIo{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PodNetworkIo) Name() string {
	return "k8s.pod.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (PodNetworkIo) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodNetworkIo) Description() string {
	return "Network bytes for the Pod"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m PodNetworkIo) Add(
    ctx context.Context,
    incr int64,
	attrs ...PodNetworkIoAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m PodNetworkIo) conv(in []PodNetworkIoAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.podNetworkIoAttr()
	}
	return out
}

// PodNetworkIoAttr is an optional attribute for the PodNetworkIo instrument.
type PodNetworkIoAttr interface {
    podNetworkIoAttr() attribute.KeyValue
}

type podNetworkIoAttr struct {
	kv attribute.KeyValue
}

func (a podNetworkIoAttr) podNetworkIoAttr() attribute.KeyValue {
    return a.kv
}

// NetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (PodNetworkIo) NetworkInterfaceName(val string) PodNetworkIoAttr {
	return podNetworkIoAttr{kv: attribute.String("network.interface.name", val)}
}

// NetworkIoDirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (PodNetworkIo) NetworkIoDirection(val NetworkIoDirectionAttr) PodNetworkIoAttr {
	return podNetworkIoAttr{kv: attribute.String("network.io.direction", string(val))}
}

// K8SPodUptime is an instrument used to record metric values conforming to the
// "k8s.pod.uptime" semantic conventions. It represents the time the Pod has been
// running.
type PodUptime struct {
	inst metric.Float64Gauge
}

// NewPodUptime returns a new PodUptime instrument.
func NewPodUptime(m metric.Meter) (PodUptime, error) {
	i, err := m.Float64Gauge(
	    "k8s.pod.uptime",
	    metric.WithDescription("The time the Pod has been running"),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return PodUptime{}, err
	}
	return PodUptime{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (PodUptime) Name() string {
	return "k8s.pod.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (PodUptime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (PodUptime) Description() string {
	return "The time the Pod has been running"
}

func (m PodUptime) Record(ctx context.Context, val float64) {
    m.inst.Record(ctx, val)
}

// K8SReplicaSetAvailablePods is an instrument used to record metric values
// conforming to the "k8s.replicaset.available_pods" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this replicaset.
type ReplicaSetAvailablePods struct {
	inst metric.Int64UpDownCounter
}

// NewReplicaSetAvailablePods returns a new ReplicaSetAvailablePods instrument.
func NewReplicaSetAvailablePods(m metric.Meter) (ReplicaSetAvailablePods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.replicaset.available_pods",
	    metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return ReplicaSetAvailablePods{}, err
	}
	return ReplicaSetAvailablePods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ReplicaSetAvailablePods) Name() string {
	return "k8s.replicaset.available_pods"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicaSetAvailablePods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicaSetAvailablePods) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset"
}

func (m ReplicaSetAvailablePods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SReplicaSetDesiredPods is an instrument used to record metric values
// conforming to the "k8s.replicaset.desired_pods" semantic conventions. It
// represents the number of desired replica pods in this replicaset.
type ReplicaSetDesiredPods struct {
	inst metric.Int64UpDownCounter
}

// NewReplicaSetDesiredPods returns a new ReplicaSetDesiredPods instrument.
func NewReplicaSetDesiredPods(m metric.Meter) (ReplicaSetDesiredPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.replicaset.desired_pods",
	    metric.WithDescription("Number of desired replica pods in this replicaset"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return ReplicaSetDesiredPods{}, err
	}
	return ReplicaSetDesiredPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ReplicaSetDesiredPods) Name() string {
	return "k8s.replicaset.desired_pods"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicaSetDesiredPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicaSetDesiredPods) Description() string {
	return "Number of desired replica pods in this replicaset"
}

func (m ReplicaSetDesiredPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SReplicationControllerAvailablePods is an instrument used to record metric
// values conforming to the "k8s.replication_controller.available_pods" semantic
// conventions. It represents the deprecated, use
// `k8s.replicationcontroller.available_pods` instead.
type ReplicationControllerAvailablePods struct {
	inst metric.Int64UpDownCounter
}

// NewReplicationControllerAvailablePods returns a new
// ReplicationControllerAvailablePods instrument.
func NewReplicationControllerAvailablePods(m metric.Meter) (ReplicationControllerAvailablePods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.replication_controller.available_pods",
	    metric.WithDescription("Deprecated, use `k8s.replicationcontroller.available_pods` instead."),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return ReplicationControllerAvailablePods{}, err
	}
	return ReplicationControllerAvailablePods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerAvailablePods) Name() string {
	return "k8s.replication_controller.available_pods"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerAvailablePods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerAvailablePods) Description() string {
	return "Deprecated, use `k8s.replicationcontroller.available_pods` instead."
}

func (m ReplicationControllerAvailablePods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SReplicationControllerDesiredPods is an instrument used to record metric
// values conforming to the "k8s.replication_controller.desired_pods" semantic
// conventions. It represents the deprecated, use
// `k8s.replicationcontroller.desired_pods` instead.
type ReplicationControllerDesiredPods struct {
	inst metric.Int64UpDownCounter
}

// NewReplicationControllerDesiredPods returns a new
// ReplicationControllerDesiredPods instrument.
func NewReplicationControllerDesiredPods(m metric.Meter) (ReplicationControllerDesiredPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.replication_controller.desired_pods",
	    metric.WithDescription("Deprecated, use `k8s.replicationcontroller.desired_pods` instead."),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return ReplicationControllerDesiredPods{}, err
	}
	return ReplicationControllerDesiredPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerDesiredPods) Name() string {
	return "k8s.replication_controller.desired_pods"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerDesiredPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerDesiredPods) Description() string {
	return "Deprecated, use `k8s.replicationcontroller.desired_pods` instead."
}

func (m ReplicationControllerDesiredPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SReplicationControllerAvailablePods is an instrument used to record metric
// values conforming to the "k8s.replicationcontroller.available_pods" semantic
// conventions. It represents the total number of available replica pods (ready
// for at least minReadySeconds) targeted by this replication controller.
type ReplicationControllerAvailablePods struct {
	inst metric.Int64UpDownCounter
}

// NewReplicationControllerAvailablePods returns a new
// ReplicationControllerAvailablePods instrument.
func NewReplicationControllerAvailablePods(m metric.Meter) (ReplicationControllerAvailablePods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.replicationcontroller.available_pods",
	    metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return ReplicationControllerAvailablePods{}, err
	}
	return ReplicationControllerAvailablePods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerAvailablePods) Name() string {
	return "k8s.replicationcontroller.available_pods"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerAvailablePods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerAvailablePods) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller"
}

func (m ReplicationControllerAvailablePods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SReplicationControllerDesiredPods is an instrument used to record metric
// values conforming to the "k8s.replicationcontroller.desired_pods" semantic
// conventions. It represents the number of desired replica pods in this
// replication controller.
type ReplicationControllerDesiredPods struct {
	inst metric.Int64UpDownCounter
}

// NewReplicationControllerDesiredPods returns a new
// ReplicationControllerDesiredPods instrument.
func NewReplicationControllerDesiredPods(m metric.Meter) (ReplicationControllerDesiredPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.replicationcontroller.desired_pods",
	    metric.WithDescription("Number of desired replica pods in this replication controller"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return ReplicationControllerDesiredPods{}, err
	}
	return ReplicationControllerDesiredPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerDesiredPods) Name() string {
	return "k8s.replicationcontroller.desired_pods"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerDesiredPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerDesiredPods) Description() string {
	return "Number of desired replica pods in this replication controller"
}

func (m ReplicationControllerDesiredPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SStatefulSetCurrentPods is an instrument used to record metric values
// conforming to the "k8s.statefulset.current_pods" semantic conventions. It
// represents the number of replica pods created by the statefulset controller
// from the statefulset version indicated by currentRevision.
type StatefulSetCurrentPods struct {
	inst metric.Int64UpDownCounter
}

// NewStatefulSetCurrentPods returns a new StatefulSetCurrentPods instrument.
func NewStatefulSetCurrentPods(m metric.Meter) (StatefulSetCurrentPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.statefulset.current_pods",
	    metric.WithDescription("The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return StatefulSetCurrentPods{}, err
	}
	return StatefulSetCurrentPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetCurrentPods) Name() string {
	return "k8s.statefulset.current_pods"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetCurrentPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetCurrentPods) Description() string {
	return "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision"
}

func (m StatefulSetCurrentPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SStatefulSetDesiredPods is an instrument used to record metric values
// conforming to the "k8s.statefulset.desired_pods" semantic conventions. It
// represents the number of desired replica pods in this statefulset.
type StatefulSetDesiredPods struct {
	inst metric.Int64UpDownCounter
}

// NewStatefulSetDesiredPods returns a new StatefulSetDesiredPods instrument.
func NewStatefulSetDesiredPods(m metric.Meter) (StatefulSetDesiredPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.statefulset.desired_pods",
	    metric.WithDescription("Number of desired replica pods in this statefulset"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return StatefulSetDesiredPods{}, err
	}
	return StatefulSetDesiredPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetDesiredPods) Name() string {
	return "k8s.statefulset.desired_pods"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetDesiredPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetDesiredPods) Description() string {
	return "Number of desired replica pods in this statefulset"
}

func (m StatefulSetDesiredPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SStatefulSetReadyPods is an instrument used to record metric values
// conforming to the "k8s.statefulset.ready_pods" semantic conventions. It
// represents the number of replica pods created for this statefulset with a
// Ready Condition.
type StatefulSetReadyPods struct {
	inst metric.Int64UpDownCounter
}

// NewStatefulSetReadyPods returns a new StatefulSetReadyPods instrument.
func NewStatefulSetReadyPods(m metric.Meter) (StatefulSetReadyPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.statefulset.ready_pods",
	    metric.WithDescription("The number of replica pods created for this statefulset with a Ready Condition"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return StatefulSetReadyPods{}, err
	}
	return StatefulSetReadyPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetReadyPods) Name() string {
	return "k8s.statefulset.ready_pods"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetReadyPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetReadyPods) Description() string {
	return "The number of replica pods created for this statefulset with a Ready Condition"
}

func (m StatefulSetReadyPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}

// K8SStatefulSetUpdatedPods is an instrument used to record metric values
// conforming to the "k8s.statefulset.updated_pods" semantic conventions. It
// represents the number of replica pods created by the statefulset controller
// from the statefulset version indicated by updateRevision.
type StatefulSetUpdatedPods struct {
	inst metric.Int64UpDownCounter
}

// NewStatefulSetUpdatedPods returns a new StatefulSetUpdatedPods instrument.
func NewStatefulSetUpdatedPods(m metric.Meter) (StatefulSetUpdatedPods, error) {
	i, err := m.Int64UpDownCounter(
	    "k8s.statefulset.updated_pods",
	    metric.WithDescription("Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision"),
	    metric.WithUnit("{pod}"),
	)
	if err != nil {
	    return StatefulSetUpdatedPods{}, err
	}
	return StatefulSetUpdatedPods{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetUpdatedPods) Name() string {
	return "k8s.statefulset.updated_pods"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetUpdatedPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetUpdatedPods) Description() string {
	return "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision"
}

func (m StatefulSetUpdatedPods) Add(ctx context.Context, incr int64) {
    m.inst.Add(ctx, incr)
}