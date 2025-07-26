// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "k8s" namespace.
package k8sconv

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

// NetworkIODirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the network IO operation direction.
type NetworkIODirectionAttr string

var (
	// NetworkIODirectionTransmit is the none.
	NetworkIODirectionTransmit NetworkIODirectionAttr = "transmit"
	// NetworkIODirectionReceive is the none.
	NetworkIODirectionReceive NetworkIODirectionAttr = "receive"
)

// CronJobActiveJobs is an instrument used to record metric values conforming to
// the "k8s.cronjob.active_jobs" semantic conventions. It represents the number
// of actively running jobs for a cronjob.
type CronJobActiveJobs struct {
	metric.Int64UpDownCounter
}

// NewCronJobActiveJobs returns a new CronJobActiveJobs instrument.
func NewCronJobActiveJobs(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (CronJobActiveJobs, error) {
	// Check if the meter is nil.
	if m == nil {
		return CronJobActiveJobs{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.cronjob.active_jobs",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of actively running jobs for a cronjob"),
			metric.WithUnit("{job}"),
		}, opt...)...,
	)
	if err != nil {
	    return CronJobActiveJobs{noop.Int64UpDownCounter{}}, err
	}
	return CronJobActiveJobs{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CronJobActiveJobs) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `active` field of the
// [K8s CronJobStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.cronjob`] resource.
//
// [K8s CronJobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#cronjobstatus-v1-batch
// [`k8s.cronjob`]: ../resource/k8s.md#cronjob
func (m CronJobActiveJobs) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetCurrentScheduledNodes is an instrument used to record metric values
// conforming to the "k8s.daemonset.current_scheduled_nodes" semantic
// conventions. It represents the number of nodes that are running at least 1
// daemon pod and are supposed to run the daemon pod.
type DaemonSetCurrentScheduledNodes struct {
	metric.Int64UpDownCounter
}

// NewDaemonSetCurrentScheduledNodes returns a new DaemonSetCurrentScheduledNodes
// instrument.
func NewDaemonSetCurrentScheduledNodes(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetCurrentScheduledNodes, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetCurrentScheduledNodes{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.current_scheduled_nodes",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod"),
			metric.WithUnit("{node}"),
		}, opt...)...,
	)
	if err != nil {
	    return DaemonSetCurrentScheduledNodes{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetCurrentScheduledNodes{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetCurrentScheduledNodes) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `currentNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.daemonset`] resource.
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
// [`k8s.daemonset`]: ../resource/k8s.md#daemonset
func (m DaemonSetCurrentScheduledNodes) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetDesiredScheduledNodes is an instrument used to record metric values
// conforming to the "k8s.daemonset.desired_scheduled_nodes" semantic
// conventions. It represents the number of nodes that should be running the
// daemon pod (including nodes currently running the daemon pod).
type DaemonSetDesiredScheduledNodes struct {
	metric.Int64UpDownCounter
}

// NewDaemonSetDesiredScheduledNodes returns a new DaemonSetDesiredScheduledNodes
// instrument.
func NewDaemonSetDesiredScheduledNodes(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetDesiredScheduledNodes, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetDesiredScheduledNodes{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.desired_scheduled_nodes",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)"),
			metric.WithUnit("{node}"),
		}, opt...)...,
	)
	if err != nil {
	    return DaemonSetDesiredScheduledNodes{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetDesiredScheduledNodes{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetDesiredScheduledNodes) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `desiredNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.daemonset`] resource.
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
// [`k8s.daemonset`]: ../resource/k8s.md#daemonset
func (m DaemonSetDesiredScheduledNodes) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetMisscheduledNodes is an instrument used to record metric values
// conforming to the "k8s.daemonset.misscheduled_nodes" semantic conventions. It
// represents the number of nodes that are running the daemon pod, but are not
// supposed to run the daemon pod.
type DaemonSetMisscheduledNodes struct {
	metric.Int64UpDownCounter
}

// NewDaemonSetMisscheduledNodes returns a new DaemonSetMisscheduledNodes
// instrument.
func NewDaemonSetMisscheduledNodes(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetMisscheduledNodes, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetMisscheduledNodes{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.misscheduled_nodes",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod"),
			metric.WithUnit("{node}"),
		}, opt...)...,
	)
	if err != nil {
	    return DaemonSetMisscheduledNodes{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetMisscheduledNodes{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetMisscheduledNodes) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `numberMisscheduled` field of the
// [K8s DaemonSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.daemonset`] resource.
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
// [`k8s.daemonset`]: ../resource/k8s.md#daemonset
func (m DaemonSetMisscheduledNodes) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetReadyNodes is an instrument used to record metric values conforming
// to the "k8s.daemonset.ready_nodes" semantic conventions. It represents the
// number of nodes that should be running the daemon pod and have one or more of
// the daemon pod running and ready.
type DaemonSetReadyNodes struct {
	metric.Int64UpDownCounter
}

// NewDaemonSetReadyNodes returns a new DaemonSetReadyNodes instrument.
func NewDaemonSetReadyNodes(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetReadyNodes, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetReadyNodes{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.ready_nodes",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready"),
			metric.WithUnit("{node}"),
		}, opt...)...,
	)
	if err != nil {
	    return DaemonSetReadyNodes{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetReadyNodes{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetReadyNodes) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `numberReady` field of the
// [K8s DaemonSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.daemonset`] resource.
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
// [`k8s.daemonset`]: ../resource/k8s.md#daemonset
func (m DaemonSetReadyNodes) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DeploymentAvailablePods is an instrument used to record metric values
// conforming to the "k8s.deployment.available_pods" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this deployment.
type DeploymentAvailablePods struct {
	metric.Int64UpDownCounter
}

// NewDeploymentAvailablePods returns a new DeploymentAvailablePods instrument.
func NewDeploymentAvailablePods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DeploymentAvailablePods, error) {
	// Check if the meter is nil.
	if m == nil {
		return DeploymentAvailablePods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.deployment.available_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return DeploymentAvailablePods{noop.Int64UpDownCounter{}}, err
	}
	return DeploymentAvailablePods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DeploymentAvailablePods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s DeploymentStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.deployment`] resource.
//
// [K8s DeploymentStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentstatus-v1-apps
// [`k8s.deployment`]: ../resource/k8s.md#deployment
func (m DeploymentAvailablePods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DeploymentDesiredPods is an instrument used to record metric values conforming
// to the "k8s.deployment.desired_pods" semantic conventions. It represents the
// number of desired replica pods in this deployment.
type DeploymentDesiredPods struct {
	metric.Int64UpDownCounter
}

// NewDeploymentDesiredPods returns a new DeploymentDesiredPods instrument.
func NewDeploymentDesiredPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DeploymentDesiredPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return DeploymentDesiredPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.deployment.desired_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of desired replica pods in this deployment"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return DeploymentDesiredPods{noop.Int64UpDownCounter{}}, err
	}
	return DeploymentDesiredPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DeploymentDesiredPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `replicas` field of the
// [K8s DeploymentSpec].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.deployment`] resource.
//
// [K8s DeploymentSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentspec-v1-apps
// [`k8s.deployment`]: ../resource/k8s.md#deployment
func (m DeploymentDesiredPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPACurrentPods is an instrument used to record metric values conforming to the
// "k8s.hpa.current_pods" semantic conventions. It represents the current number
// of replica pods managed by this horizontal pod autoscaler, as last seen by the
// autoscaler.
type HPACurrentPods struct {
	metric.Int64UpDownCounter
}

// NewHPACurrentPods returns a new HPACurrentPods instrument.
func NewHPACurrentPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPACurrentPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPACurrentPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.current_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return HPACurrentPods{noop.Int64UpDownCounter{}}, err
	}
	return HPACurrentPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPACurrentPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPACurrentPods) Name() string {
	return "k8s.hpa.current_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HPACurrentPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPACurrentPods) Description() string {
	return "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler"
}

// Add adds incr to the existing count.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.hpa`] resource.
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
// [`k8s.hpa`]: ../resource/k8s.md#horizontalpodautoscaler
func (m HPACurrentPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPADesiredPods is an instrument used to record metric values conforming to the
// "k8s.hpa.desired_pods" semantic conventions. It represents the desired number
// of replica pods managed by this horizontal pod autoscaler, as last calculated
// by the autoscaler.
type HPADesiredPods struct {
	metric.Int64UpDownCounter
}

// NewHPADesiredPods returns a new HPADesiredPods instrument.
func NewHPADesiredPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPADesiredPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPADesiredPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.desired_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return HPADesiredPods{noop.Int64UpDownCounter{}}, err
	}
	return HPADesiredPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPADesiredPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPADesiredPods) Name() string {
	return "k8s.hpa.desired_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HPADesiredPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPADesiredPods) Description() string {
	return "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler"
}

// Add adds incr to the existing count.
//
// This metric aligns with the `desiredReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.hpa`] resource.
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
// [`k8s.hpa`]: ../resource/k8s.md#horizontalpodautoscaler
func (m HPADesiredPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAMaxPods is an instrument used to record metric values conforming to the
// "k8s.hpa.max_pods" semantic conventions. It represents the upper limit for the
// number of replica pods to which the autoscaler can scale up.
type HPAMaxPods struct {
	metric.Int64UpDownCounter
}

// NewHPAMaxPods returns a new HPAMaxPods instrument.
func NewHPAMaxPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPAMaxPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMaxPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.max_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The upper limit for the number of replica pods to which the autoscaler can scale up"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return HPAMaxPods{noop.Int64UpDownCounter{}}, err
	}
	return HPAMaxPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMaxPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAMaxPods) Name() string {
	return "k8s.hpa.max_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMaxPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAMaxPods) Description() string {
	return "The upper limit for the number of replica pods to which the autoscaler can scale up"
}

// Add adds incr to the existing count.
//
// This metric aligns with the `maxReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.hpa`] resource.
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
// [`k8s.hpa`]: ../resource/k8s.md#horizontalpodautoscaler
func (m HPAMaxPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAMinPods is an instrument used to record metric values conforming to the
// "k8s.hpa.min_pods" semantic conventions. It represents the lower limit for the
// number of replica pods to which the autoscaler can scale down.
type HPAMinPods struct {
	metric.Int64UpDownCounter
}

// NewHPAMinPods returns a new HPAMinPods instrument.
func NewHPAMinPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPAMinPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMinPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.min_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The lower limit for the number of replica pods to which the autoscaler can scale down"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return HPAMinPods{noop.Int64UpDownCounter{}}, err
	}
	return HPAMinPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMinPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAMinPods) Name() string {
	return "k8s.hpa.min_pods"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMinPods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAMinPods) Description() string {
	return "The lower limit for the number of replica pods to which the autoscaler can scale down"
}

// Add adds incr to the existing count.
//
// This metric aligns with the `minReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.hpa`] resource.
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
// [`k8s.hpa`]: ../resource/k8s.md#horizontalpodautoscaler
func (m HPAMinPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobActivePods is an instrument used to record metric values conforming to the
// "k8s.job.active_pods" semantic conventions. It represents the number of
// pending and actively running pods for a job.
type JobActivePods struct {
	metric.Int64UpDownCounter
}

// NewJobActivePods returns a new JobActivePods instrument.
func NewJobActivePods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobActivePods, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobActivePods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.active_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of pending and actively running pods for a job"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return JobActivePods{noop.Int64UpDownCounter{}}, err
	}
	return JobActivePods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobActivePods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `active` field of the
// [K8s JobStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.job`] resource.
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
// [`k8s.job`]: ../resource/k8s.md#job
func (m JobActivePods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobDesiredSuccessfulPods is an instrument used to record metric values
// conforming to the "k8s.job.desired_successful_pods" semantic conventions. It
// represents the desired number of successfully finished pods the job should be
// run with.
type JobDesiredSuccessfulPods struct {
	metric.Int64UpDownCounter
}

// NewJobDesiredSuccessfulPods returns a new JobDesiredSuccessfulPods instrument.
func NewJobDesiredSuccessfulPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobDesiredSuccessfulPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobDesiredSuccessfulPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.desired_successful_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The desired number of successfully finished pods the job should be run with"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return JobDesiredSuccessfulPods{noop.Int64UpDownCounter{}}, err
	}
	return JobDesiredSuccessfulPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobDesiredSuccessfulPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `completions` field of the
// [K8s JobSpec].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.job`] resource.
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
// [`k8s.job`]: ../resource/k8s.md#job
func (m JobDesiredSuccessfulPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobFailedPods is an instrument used to record metric values conforming to the
// "k8s.job.failed_pods" semantic conventions. It represents the number of pods
// which reached phase Failed for a job.
type JobFailedPods struct {
	metric.Int64UpDownCounter
}

// NewJobFailedPods returns a new JobFailedPods instrument.
func NewJobFailedPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobFailedPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobFailedPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.failed_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of pods which reached phase Failed for a job"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return JobFailedPods{noop.Int64UpDownCounter{}}, err
	}
	return JobFailedPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobFailedPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `failed` field of the
// [K8s JobStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.job`] resource.
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
// [`k8s.job`]: ../resource/k8s.md#job
func (m JobFailedPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobMaxParallelPods is an instrument used to record metric values conforming to
// the "k8s.job.max_parallel_pods" semantic conventions. It represents the max
// desired number of pods the job should run at any given time.
type JobMaxParallelPods struct {
	metric.Int64UpDownCounter
}

// NewJobMaxParallelPods returns a new JobMaxParallelPods instrument.
func NewJobMaxParallelPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobMaxParallelPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobMaxParallelPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.max_parallel_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The max desired number of pods the job should run at any given time"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return JobMaxParallelPods{noop.Int64UpDownCounter{}}, err
	}
	return JobMaxParallelPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobMaxParallelPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `parallelism` field of the
// [K8s JobSpec].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.job`] resource.
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
// [`k8s.job`]: ../resource/k8s.md#job
func (m JobMaxParallelPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobSuccessfulPods is an instrument used to record metric values conforming to
// the "k8s.job.successful_pods" semantic conventions. It represents the number
// of pods which reached phase Succeeded for a job.
type JobSuccessfulPods struct {
	metric.Int64UpDownCounter
}

// NewJobSuccessfulPods returns a new JobSuccessfulPods instrument.
func NewJobSuccessfulPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobSuccessfulPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobSuccessfulPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.successful_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of pods which reached phase Succeeded for a job"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return JobSuccessfulPods{noop.Int64UpDownCounter{}}, err
	}
	return JobSuccessfulPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobSuccessfulPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `succeeded` field of the
// [K8s JobStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.job`] resource.
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
// [`k8s.job`]: ../resource/k8s.md#job
func (m JobSuccessfulPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NamespacePhase is an instrument used to record metric values conforming to the
// "k8s.namespace.phase" semantic conventions. It represents the describes number
// of K8s namespaces that are currently in a given phase.
type NamespacePhase struct {
	metric.Int64UpDownCounter
}

// NewNamespacePhase returns a new NamespacePhase instrument.
func NewNamespacePhase(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NamespacePhase, error) {
	// Check if the meter is nil.
	if m == nil {
		return NamespacePhase{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.namespace.phase",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Describes number of K8s namespaces that are currently in a given phase."),
			metric.WithUnit("{namespace}"),
		}, opt...)...,
	)
	if err != nil {
	    return NamespacePhase{noop.Int64UpDownCounter{}}, err
	}
	return NamespacePhase{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NamespacePhase) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
// The namespacePhase is the the phase of the K8s namespace.
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.namespace`] resource.
//
// [`k8s.namespace`]: ../resource/k8s.md#namespace
func (m NamespacePhase) Add(
	ctx context.Context,
	incr int64,
	namespacePhase NamespacePhaseAttr,
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
				attribute.String("k8s.namespace.phase", string(namespacePhase)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeCPUTime is an instrument used to record metric values conforming to the
// "k8s.node.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type NodeCPUTime struct {
	metric.Float64Counter
}

// NewNodeCPUTime returns a new NodeCPUTime instrument.
func NewNodeCPUTime(
	m metric.Meter,
	opt ...metric.Float64CounterOption,
) (NodeCPUTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUTime{noop.Float64Counter{}}, nil
	}

	i, err := m.Float64Counter(
		"k8s.node.cpu.time",
		append([]metric.Float64CounterOption{
			metric.WithDescription("Total CPU time consumed"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeCPUTime{noop.Float64Counter{}}, err
	}
	return NodeCPUTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUTime) Inst() metric.Float64Counter {
	return m.Float64Counter
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

// Add adds incr to the existing count.
//
// Total CPU time consumed by the specific Node on all available CPU cores
func (m NodeCPUTime) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// NodeCPUUsage is an instrument used to record metric values conforming to the
// "k8s.node.cpu.usage" semantic conventions. It represents the node's CPU usage,
// measured in cpus. Range from 0 to the number of allocatable CPUs.
type NodeCPUUsage struct {
	metric.Int64Gauge
}

// NewNodeCPUUsage returns a new NodeCPUUsage instrument.
func NewNodeCPUUsage(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (NodeCPUUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUUsage{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"k8s.node.cpu.usage",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeCPUUsage{noop.Int64Gauge{}}, err
	}
	return NodeCPUUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUUsage) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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

// Record records val to the current distribution.
//
// CPU usage of the specific Node on all available CPU cores, averaged over the
// sample window
func (m NodeCPUUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// NodeMemoryUsage is an instrument used to record metric values conforming to
// the "k8s.node.memory.usage" semantic conventions. It represents the memory
// usage of the Node.
type NodeMemoryUsage struct {
	metric.Int64Gauge
}

// NewNodeMemoryUsage returns a new NodeMemoryUsage instrument.
func NewNodeMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (NodeMemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryUsage{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"k8s.node.memory.usage",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Memory usage of the Node"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeMemoryUsage{noop.Int64Gauge{}}, err
	}
	return NodeMemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryUsage) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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

// Record records val to the current distribution.
//
// Total memory usage of the Node
func (m NodeMemoryUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// NodeNetworkErrors is an instrument used to record metric values conforming to
// the "k8s.node.network.errors" semantic conventions. It represents the node
// network errors.
type NodeNetworkErrors struct {
	metric.Int64Counter
}

// NewNodeNetworkErrors returns a new NodeNetworkErrors instrument.
func NewNodeNetworkErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NodeNetworkErrors, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeNetworkErrors{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"k8s.node.network.errors",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Node network errors"),
			metric.WithUnit("{error}"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeNetworkErrors{noop.Int64Counter{}}, err
	}
	return NodeNetworkErrors{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeNetworkErrors) Inst() metric.Int64Counter {
	return m.Int64Counter
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

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NodeNetworkErrors) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NodeNetworkErrors) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NodeNetworkIO is an instrument used to record metric values conforming to the
// "k8s.node.network.io" semantic conventions. It represents the network bytes
// for the Node.
type NodeNetworkIO struct {
	metric.Int64Counter
}

// NewNodeNetworkIO returns a new NodeNetworkIO instrument.
func NewNodeNetworkIO(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NodeNetworkIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeNetworkIO{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"k8s.node.network.io",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Network bytes for the Node"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeNetworkIO{noop.Int64Counter{}}, err
	}
	return NodeNetworkIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeNetworkIO) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (NodeNetworkIO) Name() string {
	return "k8s.node.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NodeNetworkIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeNetworkIO) Description() string {
	return "Network bytes for the Node"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m NodeNetworkIO) Add(
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

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NodeNetworkIO) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (NodeNetworkIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NodeUptime is an instrument used to record metric values conforming to the
// "k8s.node.uptime" semantic conventions. It represents the time the Node has
// been running.
type NodeUptime struct {
	metric.Float64Gauge
}

// NewNodeUptime returns a new NodeUptime instrument.
func NewNodeUptime(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (NodeUptime, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeUptime{noop.Float64Gauge{}}, nil
	}

	i, err := m.Float64Gauge(
		"k8s.node.uptime",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("The time the Node has been running"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeUptime{noop.Float64Gauge{}}, err
	}
	return NodeUptime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeUptime) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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

// Record records val to the current distribution.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m NodeUptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodCPUTime is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type PodCPUTime struct {
	metric.Float64Counter
}

// NewPodCPUTime returns a new PodCPUTime instrument.
func NewPodCPUTime(
	m metric.Meter,
	opt ...metric.Float64CounterOption,
) (PodCPUTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPUTime{noop.Float64Counter{}}, nil
	}

	i, err := m.Float64Counter(
		"k8s.pod.cpu.time",
		append([]metric.Float64CounterOption{
			metric.WithDescription("Total CPU time consumed"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return PodCPUTime{noop.Float64Counter{}}, err
	}
	return PodCPUTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPUTime) Inst() metric.Float64Counter {
	return m.Float64Counter
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

// Add adds incr to the existing count.
//
// Total CPU time consumed by the specific Pod on all available CPU cores
func (m PodCPUTime) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// PodCPUUsage is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.usage" semantic conventions. It represents the pod's CPU usage,
// measured in cpus. Range from 0 to the number of allocatable CPUs.
type PodCPUUsage struct {
	metric.Int64Gauge
}

// NewPodCPUUsage returns a new PodCPUUsage instrument.
func NewPodCPUUsage(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (PodCPUUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPUUsage{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"k8s.pod.cpu.usage",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs"),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
	)
	if err != nil {
	    return PodCPUUsage{noop.Int64Gauge{}}, err
	}
	return PodCPUUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPUUsage) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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

// Record records val to the current distribution.
//
// CPU usage of the specific Pod on all available CPU cores, averaged over the
// sample window
func (m PodCPUUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// PodMemoryUsage is an instrument used to record metric values conforming to the
// "k8s.pod.memory.usage" semantic conventions. It represents the memory usage of
// the Pod.
type PodMemoryUsage struct {
	metric.Int64Gauge
}

// NewPodMemoryUsage returns a new PodMemoryUsage instrument.
func NewPodMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (PodMemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryUsage{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"k8s.pod.memory.usage",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Memory usage of the Pod"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return PodMemoryUsage{noop.Int64Gauge{}}, err
	}
	return PodMemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryUsage) Inst() metric.Int64Gauge {
	return m.Int64Gauge
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

// Record records val to the current distribution.
//
// Total memory usage of the Pod
func (m PodMemoryUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// PodNetworkErrors is an instrument used to record metric values conforming to
// the "k8s.pod.network.errors" semantic conventions. It represents the pod
// network errors.
type PodNetworkErrors struct {
	metric.Int64Counter
}

// NewPodNetworkErrors returns a new PodNetworkErrors instrument.
func NewPodNetworkErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (PodNetworkErrors, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodNetworkErrors{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"k8s.pod.network.errors",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Pod network errors"),
			metric.WithUnit("{error}"),
		}, opt...)...,
	)
	if err != nil {
	    return PodNetworkErrors{noop.Int64Counter{}}, err
	}
	return PodNetworkErrors{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodNetworkErrors) Inst() metric.Int64Counter {
	return m.Int64Counter
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

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (PodNetworkErrors) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (PodNetworkErrors) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// PodNetworkIO is an instrument used to record metric values conforming to the
// "k8s.pod.network.io" semantic conventions. It represents the network bytes for
// the Pod.
type PodNetworkIO struct {
	metric.Int64Counter
}

// NewPodNetworkIO returns a new PodNetworkIO instrument.
func NewPodNetworkIO(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (PodNetworkIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodNetworkIO{noop.Int64Counter{}}, nil
	}

	i, err := m.Int64Counter(
		"k8s.pod.network.io",
		append([]metric.Int64CounterOption{
			metric.WithDescription("Network bytes for the Pod"),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return PodNetworkIO{noop.Int64Counter{}}, err
	}
	return PodNetworkIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodNetworkIO) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (PodNetworkIO) Name() string {
	return "k8s.pod.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (PodNetworkIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodNetworkIO) Description() string {
	return "Network bytes for the Pod"
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m PodNetworkIO) Add(
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

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (PodNetworkIO) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (PodNetworkIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// PodUptime is an instrument used to record metric values conforming to the
// "k8s.pod.uptime" semantic conventions. It represents the time the Pod has been
// running.
type PodUptime struct {
	metric.Float64Gauge
}

// NewPodUptime returns a new PodUptime instrument.
func NewPodUptime(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (PodUptime, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodUptime{noop.Float64Gauge{}}, nil
	}

	i, err := m.Float64Gauge(
		"k8s.pod.uptime",
		append([]metric.Float64GaugeOption{
			metric.WithDescription("The time the Pod has been running"),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return PodUptime{noop.Float64Gauge{}}, err
	}
	return PodUptime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodUptime) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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

// Record records val to the current distribution.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m PodUptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// ReplicaSetAvailablePods is an instrument used to record metric values
// conforming to the "k8s.replicaset.available_pods" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this replicaset.
type ReplicaSetAvailablePods struct {
	metric.Int64UpDownCounter
}

// NewReplicaSetAvailablePods returns a new ReplicaSetAvailablePods instrument.
func NewReplicaSetAvailablePods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicaSetAvailablePods, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicaSetAvailablePods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicaset.available_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return ReplicaSetAvailablePods{noop.Int64UpDownCounter{}}, err
	}
	return ReplicaSetAvailablePods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicaSetAvailablePods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicaSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.replicaset`] resource.
//
// [K8s ReplicaSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetstatus-v1-apps
// [`k8s.replicaset`]: ../resource/k8s.md#replicaset
func (m ReplicaSetAvailablePods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicaSetDesiredPods is an instrument used to record metric values conforming
// to the "k8s.replicaset.desired_pods" semantic conventions. It represents the
// number of desired replica pods in this replicaset.
type ReplicaSetDesiredPods struct {
	metric.Int64UpDownCounter
}

// NewReplicaSetDesiredPods returns a new ReplicaSetDesiredPods instrument.
func NewReplicaSetDesiredPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicaSetDesiredPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicaSetDesiredPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicaset.desired_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of desired replica pods in this replicaset"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return ReplicaSetDesiredPods{noop.Int64UpDownCounter{}}, err
	}
	return ReplicaSetDesiredPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicaSetDesiredPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicaSetSpec].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.replicaset`] resource.
//
// [K8s ReplicaSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetspec-v1-apps
// [`k8s.replicaset`]: ../resource/k8s.md#replicaset
func (m ReplicaSetDesiredPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicationControllerAvailablePods is an instrument used to record metric
// values conforming to the "k8s.replicationcontroller.available_pods" semantic
// conventions. It represents the total number of available replica pods (ready
// for at least minReadySeconds) targeted by this replication controller.
type ReplicationControllerAvailablePods struct {
	metric.Int64UpDownCounter
}

// NewReplicationControllerAvailablePods returns a new
// ReplicationControllerAvailablePods instrument.
func NewReplicationControllerAvailablePods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicationControllerAvailablePods, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicationControllerAvailablePods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicationcontroller.available_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return ReplicationControllerAvailablePods{noop.Int64UpDownCounter{}}, err
	}
	return ReplicationControllerAvailablePods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicationControllerAvailablePods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicationControllerStatus]
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.replicationcontroller`] resource.
//
// [K8s ReplicationControllerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerstatus-v1-core
// [`k8s.replicationcontroller`]: ../resource/k8s.md#replicationcontroller
func (m ReplicationControllerAvailablePods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicationControllerDesiredPods is an instrument used to record metric values
// conforming to the "k8s.replicationcontroller.desired_pods" semantic
// conventions. It represents the number of desired replica pods in this
// replication controller.
type ReplicationControllerDesiredPods struct {
	metric.Int64UpDownCounter
}

// NewReplicationControllerDesiredPods returns a new
// ReplicationControllerDesiredPods instrument.
func NewReplicationControllerDesiredPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicationControllerDesiredPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicationControllerDesiredPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicationcontroller.desired_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of desired replica pods in this replication controller"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return ReplicationControllerDesiredPods{noop.Int64UpDownCounter{}}, err
	}
	return ReplicationControllerDesiredPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicationControllerDesiredPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicationControllerSpec]
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.replicationcontroller`] resource.
//
// [K8s ReplicationControllerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerspec-v1-core
// [`k8s.replicationcontroller`]: ../resource/k8s.md#replicationcontroller
func (m ReplicationControllerDesiredPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetCurrentPods is an instrument used to record metric values
// conforming to the "k8s.statefulset.current_pods" semantic conventions. It
// represents the number of replica pods created by the statefulset controller
// from the statefulset version indicated by currentRevision.
type StatefulSetCurrentPods struct {
	metric.Int64UpDownCounter
}

// NewStatefulSetCurrentPods returns a new StatefulSetCurrentPods instrument.
func NewStatefulSetCurrentPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetCurrentPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetCurrentPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.current_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return StatefulSetCurrentPods{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetCurrentPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetCurrentPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s StatefulSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.statefulset`] resource.
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
// [`k8s.statefulset`]: ../resource/k8s.md#statefulset
func (m StatefulSetCurrentPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetDesiredPods is an instrument used to record metric values
// conforming to the "k8s.statefulset.desired_pods" semantic conventions. It
// represents the number of desired replica pods in this statefulset.
type StatefulSetDesiredPods struct {
	metric.Int64UpDownCounter
}

// NewStatefulSetDesiredPods returns a new StatefulSetDesiredPods instrument.
func NewStatefulSetDesiredPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetDesiredPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetDesiredPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.desired_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of desired replica pods in this statefulset"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return StatefulSetDesiredPods{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetDesiredPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetDesiredPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `replicas` field of the
// [K8s StatefulSetSpec].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.statefulset`] resource.
//
// [K8s StatefulSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetspec-v1-apps
// [`k8s.statefulset`]: ../resource/k8s.md#statefulset
func (m StatefulSetDesiredPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetReadyPods is an instrument used to record metric values conforming
// to the "k8s.statefulset.ready_pods" semantic conventions. It represents the
// number of replica pods created for this statefulset with a Ready Condition.
type StatefulSetReadyPods struct {
	metric.Int64UpDownCounter
}

// NewStatefulSetReadyPods returns a new StatefulSetReadyPods instrument.
func NewStatefulSetReadyPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetReadyPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetReadyPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.ready_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The number of replica pods created for this statefulset with a Ready Condition"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return StatefulSetReadyPods{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetReadyPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetReadyPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `readyReplicas` field of the
// [K8s StatefulSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.statefulset`] resource.
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
// [`k8s.statefulset`]: ../resource/k8s.md#statefulset
func (m StatefulSetReadyPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetUpdatedPods is an instrument used to record metric values
// conforming to the "k8s.statefulset.updated_pods" semantic conventions. It
// represents the number of replica pods created by the statefulset controller
// from the statefulset version indicated by updateRevision.
type StatefulSetUpdatedPods struct {
	metric.Int64UpDownCounter
}

// NewStatefulSetUpdatedPods returns a new StatefulSetUpdatedPods instrument.
func NewStatefulSetUpdatedPods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetUpdatedPods, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetUpdatedPods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.updated_pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision"),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return StatefulSetUpdatedPods{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetUpdatedPods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetUpdatedPods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count.
//
// This metric aligns with the `updatedReplicas` field of the
// [K8s StatefulSetStatus].
//
// This metric SHOULD, at a minimum, be reported against a
// [`k8s.statefulset`] resource.
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
// [`k8s.statefulset`]: ../resource/k8s.md#statefulset
func (m StatefulSetUpdatedPods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}