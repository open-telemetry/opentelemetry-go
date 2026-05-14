// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package k8sconv provides types and functionality for OpenTelemetry semantic
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

// ContainerStatusReasonAttr is an attribute conforming to the
// k8s.container.status.reason semantic conventions. It represents the reason for
// the container state. Corresponds to the `reason` field of the:
// [K8s ContainerStateWaiting] or [K8s ContainerStateTerminated].
//
// [K8s ContainerStateWaiting]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstatewaiting-v1-core
// [K8s ContainerStateTerminated]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstateterminated-v1-core
type ContainerStatusReasonAttr string

var (
	// ContainerStatusReasonContainerCreating is the container is being created.
	ContainerStatusReasonContainerCreating ContainerStatusReasonAttr = "ContainerCreating"
	// ContainerStatusReasonCrashLoopBackOff is the container is in a crash loop
	// back off state.
	ContainerStatusReasonCrashLoopBackOff ContainerStatusReasonAttr = "CrashLoopBackOff"
	// ContainerStatusReasonCreateContainerConfigError is the there was an error
	// creating the container configuration.
	ContainerStatusReasonCreateContainerConfigError ContainerStatusReasonAttr = "CreateContainerConfigError"
	// ContainerStatusReasonErrImagePull is the there was an error pulling the
	// container image.
	ContainerStatusReasonErrImagePull ContainerStatusReasonAttr = "ErrImagePull"
	// ContainerStatusReasonImagePullBackOff is the container image pull is in back
	// off state.
	ContainerStatusReasonImagePullBackOff ContainerStatusReasonAttr = "ImagePullBackOff"
	// ContainerStatusReasonOomKilled is the container was killed due to out of
	// memory.
	ContainerStatusReasonOomKilled ContainerStatusReasonAttr = "OOMKilled"
	// ContainerStatusReasonCompleted is the container has completed execution.
	ContainerStatusReasonCompleted ContainerStatusReasonAttr = "Completed"
	// ContainerStatusReasonError is the there was an error with the container.
	ContainerStatusReasonError ContainerStatusReasonAttr = "Error"
	// ContainerStatusReasonContainerCannotRun is the container cannot run.
	ContainerStatusReasonContainerCannotRun ContainerStatusReasonAttr = "ContainerCannotRun"
)

// ContainerStatusStateAttr is an attribute conforming to the
// k8s.container.status.state semantic conventions. It represents the state of
// the container. [K8s ContainerState].
//
// [K8s ContainerState]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstate-v1-core
type ContainerStatusStateAttr string

var (
	// ContainerStatusStateTerminated is the container has terminated.
	ContainerStatusStateTerminated ContainerStatusStateAttr = "terminated"
	// ContainerStatusStateRunning is the container is running.
	ContainerStatusStateRunning ContainerStatusStateAttr = "running"
	// ContainerStatusStateWaiting is the container is waiting.
	ContainerStatusStateWaiting ContainerStatusStateAttr = "waiting"
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

// NodeConditionStatusAttr is an attribute conforming to the
// k8s.node.condition.status semantic conventions. It represents the status of
// the condition, one of True, False, Unknown.
type NodeConditionStatusAttr string

var (
	// NodeConditionStatusConditionTrue is the standardized value "true" of
	// NodeConditionStatusAttr.
	NodeConditionStatusConditionTrue NodeConditionStatusAttr = "true"
	// NodeConditionStatusConditionFalse is the standardized value "false" of
	// NodeConditionStatusAttr.
	NodeConditionStatusConditionFalse NodeConditionStatusAttr = "false"
	// NodeConditionStatusConditionUnknown is the standardized value "unknown" of
	// NodeConditionStatusAttr.
	NodeConditionStatusConditionUnknown NodeConditionStatusAttr = "unknown"
)

// NodeConditionTypeAttr is an attribute conforming to the
// k8s.node.condition.type semantic conventions. It represents the condition type
// of a K8s Node.
type NodeConditionTypeAttr string

var (
	// NodeConditionTypeReady is the node is healthy and ready to accept pods.
	NodeConditionTypeReady NodeConditionTypeAttr = "Ready"
	// NodeConditionTypeDiskPressure is the pressure exists on the disk size—that
	// is, if the disk capacity is low.
	NodeConditionTypeDiskPressure NodeConditionTypeAttr = "DiskPressure"
	// NodeConditionTypeMemoryPressure is the pressure exists on the node
	// memory—that is, if the node memory is low.
	NodeConditionTypeMemoryPressure NodeConditionTypeAttr = "MemoryPressure"
	// NodeConditionTypePIDPressure is the pressure exists on the processes—that
	// is, if there are too many processes on the node.
	NodeConditionTypePIDPressure NodeConditionTypeAttr = "PIDPressure"
	// NodeConditionTypeNetworkUnavailable is the network for the node is not
	// correctly configured.
	NodeConditionTypeNetworkUnavailable NodeConditionTypeAttr = "NetworkUnavailable"
)

// PodStatusPhaseAttr is an attribute conforming to the k8s.pod.status.phase
// semantic conventions. It represents the phase for the pod. Corresponds to the
// `phase` field of the: [K8s PodStatus].
//
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#podstatus-v1-core
type PodStatusPhaseAttr string

var (
	// PodStatusPhasePending is the pod has been accepted by the system, but one or
	// more of the containers has not been started. This includes time before being
	// bound to a node, as well as time spent pulling images onto the host.
	PodStatusPhasePending PodStatusPhaseAttr = "Pending"
	// PodStatusPhaseRunning is the pod has been bound to a node and all of the
	// containers have been started. At least one container is still running or is
	// in the process of being restarted.
	PodStatusPhaseRunning PodStatusPhaseAttr = "Running"
	// PodStatusPhaseSucceeded is the all containers in the pod have voluntarily
	// terminated with a container exit code of 0, and the system is not going to
	// restart any of these containers.
	PodStatusPhaseSucceeded PodStatusPhaseAttr = "Succeeded"
	// PodStatusPhaseFailed is the all containers in the pod have terminated, and at
	// least one container has terminated in a failure (exited with a non-zero exit
	// code or was stopped by the system).
	PodStatusPhaseFailed PodStatusPhaseAttr = "Failed"
	// PodStatusPhaseUnknown is the for some reason the state of the pod could not
	// be obtained, typically due to an error in communicating with the host of the
	// pod.
	PodStatusPhaseUnknown PodStatusPhaseAttr = "Unknown"
)

// PodStatusReasonAttr is an attribute conforming to the k8s.pod.status.reason
// semantic conventions. It represents the reason for the pod state. Corresponds
// to the `reason` field of the: [K8s PodStatus].
//
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#podstatus-v1-core
type PodStatusReasonAttr string

var (
	// PodStatusReasonEvicted is the pod is evicted.
	PodStatusReasonEvicted PodStatusReasonAttr = "Evicted"
	// PodStatusReasonNodeAffinity is the pod is in a status because of its node
	// affinity.
	PodStatusReasonNodeAffinity PodStatusReasonAttr = "NodeAffinity"
	// PodStatusReasonNodeLost is the reason on a pod when its state cannot be
	// confirmed as kubelet is unresponsive on the node it is (was) running.
	PodStatusReasonNodeLost PodStatusReasonAttr = "NodeLost"
	// PodStatusReasonShutdown is the node is shutdown.
	PodStatusReasonShutdown PodStatusReasonAttr = "Shutdown"
	// PodStatusReasonUnexpectedAdmissionError is the pod was rejected admission to
	// the node because of an error during admission that could not be categorized.
	PodStatusReasonUnexpectedAdmissionError PodStatusReasonAttr = "UnexpectedAdmissionError"
)

// VolumeTypeAttr is an attribute conforming to the k8s.volume.type semantic
// conventions. It represents the type of the K8s volume.
type VolumeTypeAttr string

var (
	// VolumeTypePersistentVolumeClaim is a [persistentVolumeClaim] volume.
	//
	// [persistentVolumeClaim]: https://v1-30.docs.kubernetes.io/docs/concepts/storage/volumes/#persistentvolumeclaim
	VolumeTypePersistentVolumeClaim VolumeTypeAttr = "persistentVolumeClaim"
	// VolumeTypeConfigMap is a [configMap] volume.
	//
	// [configMap]: https://v1-30.docs.kubernetes.io/docs/concepts/storage/volumes/#configmap
	VolumeTypeConfigMap VolumeTypeAttr = "configMap"
	// VolumeTypeDownwardAPI is a [downwardAPI] volume.
	//
	// [downwardAPI]: https://v1-30.docs.kubernetes.io/docs/concepts/storage/volumes/#downwardapi
	VolumeTypeDownwardAPI VolumeTypeAttr = "downwardAPI"
	// VolumeTypeEmptyDir is an [emptyDir] volume.
	//
	// [emptyDir]: https://v1-30.docs.kubernetes.io/docs/concepts/storage/volumes/#emptydir
	VolumeTypeEmptyDir VolumeTypeAttr = "emptyDir"
	// VolumeTypeSecret is a [secret] volume.
	//
	// [secret]: https://v1-30.docs.kubernetes.io/docs/concepts/storage/volumes/#secret
	VolumeTypeSecret VolumeTypeAttr = "secret"
	// VolumeTypeLocal is a [local] volume.
	//
	// [local]: https://v1-30.docs.kubernetes.io/docs/concepts/storage/volumes/#local
	VolumeTypeLocal VolumeTypeAttr = "local"
)

// NetworkIODirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the network IO operation direction.
type NetworkIODirectionAttr string

var (
	// NetworkIODirectionTransmit is the standardized value "transmit" of
	// NetworkIODirectionAttr.
	NetworkIODirectionTransmit NetworkIODirectionAttr = "transmit"
	// NetworkIODirectionReceive is the standardized value "receive" of
	// NetworkIODirectionAttr.
	NetworkIODirectionReceive NetworkIODirectionAttr = "receive"
)

// SystemPagingFaultTypeAttr is an attribute conforming to the
// system.paging.fault.type semantic conventions. It represents the paging fault
// type.
type SystemPagingFaultTypeAttr string

var (
	// SystemPagingFaultTypeMajor is the standardized value "major" of
	// SystemPagingFaultTypeAttr.
	SystemPagingFaultTypeMajor SystemPagingFaultTypeAttr = "major"
	// SystemPagingFaultTypeMinor is the standardized value "minor" of
	// SystemPagingFaultTypeAttr.
	SystemPagingFaultTypeMinor SystemPagingFaultTypeAttr = "minor"
)

// ContainerCPULimit is an instrument used to record metric values conforming to
// the "k8s.container.cpu.limit" semantic conventions. It represents the maximum
// CPU resource limit set for the container.
type ContainerCPULimit struct {
	metric.Int64UpDownCounter
}

var newContainerCPULimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit set for the container."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPULimit returns a new ContainerCPULimit instrument.
func NewContainerCPULimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerCPULimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimit{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitOpts
	} else {
		opt = append(opt, newContainerCPULimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.cpu.limit",
		opt...,
	)
	if err != nil {
		return ContainerCPULimit{noop.Int64UpDownCounter{}}, err
	}
	return ContainerCPULimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimit) Name() string {
	return "k8s.container.cpu.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimit) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimit) Description() string {
	return "Maximum CPU resource limit set for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerCPULimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerCPULimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerCPULimitUtilization is an instrument used to record metric values
// conforming to the "k8s.container.cpu.limit_utilization" semantic conventions.
// It represents the ratio of container CPU usage to its CPU limit.
type ContainerCPULimitUtilization struct {
	metric.Int64Gauge
}

var newContainerCPULimitUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("The ratio of container CPU usage to its CPU limit."),
	metric.WithUnit("1"),
}

// NewContainerCPULimitUtilization returns a new ContainerCPULimitUtilization
// instrument.
func NewContainerCPULimitUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (ContainerCPULimitUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimitUtilization{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitUtilizationOpts
	} else {
		opt = append(opt, newContainerCPULimitUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.container.cpu.limit_utilization",
		opt...,
	)
	if err != nil {
		return ContainerCPULimitUtilization{noop.Int64Gauge{}}, err
	}
	return ContainerCPULimitUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimitUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimitUtilization) Name() string {
	return "k8s.container.cpu.limit_utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimitUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimitUtilization) Description() string {
	return "The ratio of container CPU usage to its CPU limit."
}

// Record records val to the current distribution for attrs.
//
// The value range is [0.0,1.0]. A value of 1.0 means the container is using 100%
// of its CPU limit. If the CPU limit is not set, this metric SHOULD NOT be
// emitted for that container.
func (m ContainerCPULimitUtilization) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The value range is [0.0,1.0]. A value of 1.0 means the container is using 100%
// of its CPU limit. If the CPU limit is not set, this metric SHOULD NOT be
// emitted for that container.
func (m ContainerCPULimitUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// ContainerCPURequest is an instrument used to record metric values conforming
// to the "k8s.container.cpu.request" semantic conventions. It represents the CPU
// resource requested for the container.
type ContainerCPURequest struct {
	metric.Int64UpDownCounter
}

var newContainerCPURequestOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("CPU resource requested for the container."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPURequest returns a new ContainerCPURequest instrument.
func NewContainerCPURequest(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerCPURequest, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequest{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestOpts
	} else {
		opt = append(opt, newContainerCPURequestOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.cpu.request",
		opt...,
	)
	if err != nil {
		return ContainerCPURequest{noop.Int64UpDownCounter{}}, err
	}
	return ContainerCPURequest{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequest) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequest) Name() string {
	return "k8s.container.cpu.request"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequest) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequest) Description() string {
	return "CPU resource requested for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerCPURequest) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerCPURequest) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerCPURequestUtilization is an instrument used to record metric values
// conforming to the "k8s.container.cpu.request_utilization" semantic
// conventions. It represents the ratio of container CPU usage to its CPU
// request.
type ContainerCPURequestUtilization struct {
	metric.Int64Gauge
}

var newContainerCPURequestUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("The ratio of container CPU usage to its CPU request."),
	metric.WithUnit("1"),
}

// NewContainerCPURequestUtilization returns a new ContainerCPURequestUtilization
// instrument.
func NewContainerCPURequestUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (ContainerCPURequestUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequestUtilization{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestUtilizationOpts
	} else {
		opt = append(opt, newContainerCPURequestUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.container.cpu.request_utilization",
		opt...,
	)
	if err != nil {
		return ContainerCPURequestUtilization{noop.Int64Gauge{}}, err
	}
	return ContainerCPURequestUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequestUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequestUtilization) Name() string {
	return "k8s.container.cpu.request_utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequestUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequestUtilization) Description() string {
	return "The ratio of container CPU usage to its CPU request."
}

// Record records val to the current distribution for attrs.
func (m ContainerCPURequestUtilization) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ContainerCPURequestUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// ContainerEphemeralStorageLimit is an instrument used to record metric values
// conforming to the "k8s.container.ephemeral_storage.limit" semantic
// conventions. It represents the maximum ephemeral storage resource limit set
// for the container.
type ContainerEphemeralStorageLimit struct {
	metric.Int64UpDownCounter
}

var newContainerEphemeralStorageLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum ephemeral storage resource limit set for the container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageLimit returns a new ContainerEphemeralStorageLimit
// instrument.
func NewContainerEphemeralStorageLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerEphemeralStorageLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerEphemeralStorageLimit{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerEphemeralStorageLimitOpts
	} else {
		opt = append(opt, newContainerEphemeralStorageLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.ephemeral_storage.limit",
		opt...,
	)
	if err != nil {
		return ContainerEphemeralStorageLimit{noop.Int64UpDownCounter{}}, err
	}
	return ContainerEphemeralStorageLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerEphemeralStorageLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerEphemeralStorageLimit) Name() string {
	return "k8s.container.ephemeral_storage.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerEphemeralStorageLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerEphemeralStorageLimit) Description() string {
	return "Maximum ephemeral storage resource limit set for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerEphemeralStorageLimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerEphemeralStorageLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerEphemeralStorageRequest is an instrument used to record metric values
// conforming to the "k8s.container.ephemeral_storage.request" semantic
// conventions. It represents the ephemeral storage resource requested for the
// container.
type ContainerEphemeralStorageRequest struct {
	metric.Int64UpDownCounter
}

var newContainerEphemeralStorageRequestOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Ephemeral storage resource requested for the container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageRequest returns a new
// ContainerEphemeralStorageRequest instrument.
func NewContainerEphemeralStorageRequest(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerEphemeralStorageRequest, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerEphemeralStorageRequest{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerEphemeralStorageRequestOpts
	} else {
		opt = append(opt, newContainerEphemeralStorageRequestOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.ephemeral_storage.request",
		opt...,
	)
	if err != nil {
		return ContainerEphemeralStorageRequest{noop.Int64UpDownCounter{}}, err
	}
	return ContainerEphemeralStorageRequest{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerEphemeralStorageRequest) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerEphemeralStorageRequest) Name() string {
	return "k8s.container.ephemeral_storage.request"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerEphemeralStorageRequest) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerEphemeralStorageRequest) Description() string {
	return "Ephemeral storage resource requested for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerEphemeralStorageRequest) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerEphemeralStorageRequest) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerMemoryLimit is an instrument used to record metric values conforming
// to the "k8s.container.memory.limit" semantic conventions. It represents the
// maximum memory resource limit set for the container.
type ContainerMemoryLimit struct {
	metric.Int64UpDownCounter
}

var newContainerMemoryLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit set for the container."),
	metric.WithUnit("By"),
}

// NewContainerMemoryLimit returns a new ContainerMemoryLimit instrument.
func NewContainerMemoryLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimit{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryLimitOpts
	} else {
		opt = append(opt, newContainerMemoryLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.limit",
		opt...,
	)
	if err != nil {
		return ContainerMemoryLimit{noop.Int64UpDownCounter{}}, err
	}
	return ContainerMemoryLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryLimit) Name() string {
	return "k8s.container.memory.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryLimit) Description() string {
	return "Maximum memory resource limit set for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerMemoryLimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerMemoryLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerMemoryRequest is an instrument used to record metric values
// conforming to the "k8s.container.memory.request" semantic conventions. It
// represents the memory resource requested for the container.
type ContainerMemoryRequest struct {
	metric.Int64UpDownCounter
}

var newContainerMemoryRequestOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Memory resource requested for the container."),
	metric.WithUnit("By"),
}

// NewContainerMemoryRequest returns a new ContainerMemoryRequest instrument.
func NewContainerMemoryRequest(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryRequest, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequest{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryRequestOpts
	} else {
		opt = append(opt, newContainerMemoryRequestOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.request",
		opt...,
	)
	if err != nil {
		return ContainerMemoryRequest{noop.Int64UpDownCounter{}}, err
	}
	return ContainerMemoryRequest{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryRequest) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryRequest) Name() string {
	return "k8s.container.memory.request"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryRequest) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryRequest) Description() string {
	return "Memory resource requested for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerMemoryRequest) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerMemoryRequest) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerReady is an instrument used to record metric values conforming to the
// "k8s.container.ready" semantic conventions. It represents the indicates
// whether the container is currently marked as ready to accept traffic, based on
// its readiness probe (1 = ready, 0 = not ready).
type ContainerReady struct {
	metric.Int64UpDownCounter
}

var newContainerReadyOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready)."),
	metric.WithUnit("{container}"),
}

// NewContainerReady returns a new ContainerReady instrument.
func NewContainerReady(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerReady, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerReady{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerReadyOpts
	} else {
		opt = append(opt, newContainerReadyOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.ready",
		opt...,
	)
	if err != nil {
		return ContainerReady{noop.Int64UpDownCounter{}}, err
	}
	return ContainerReady{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerReady) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerReady) Name() string {
	return "k8s.container.ready"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerReady) Unit() string {
	return "{container}"
}

// Description returns the semantic convention description of the instrument
func (ContainerReady) Description() string {
	return "Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready)."
}

// Add adds incr to the existing count for attrs.
//
// This metric SHOULD reflect the value of the `ready` field in the
// [K8s ContainerStatus].
//
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstatus-v1-core
func (m ContainerReady) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric SHOULD reflect the value of the `ready` field in the
// [K8s ContainerStatus].
//
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstatus-v1-core
func (m ContainerReady) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerRestartCount is an instrument used to record metric values conforming
// to the "k8s.container.restart.count" semantic conventions. It represents the
// describes how many times the container has restarted (since the last counter
// reset).
type ContainerRestartCount struct {
	metric.Int64UpDownCounter
}

var newContainerRestartCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes how many times the container has restarted (since the last counter reset)."),
	metric.WithUnit("{restart}"),
}

// NewContainerRestartCount returns a new ContainerRestartCount instrument.
func NewContainerRestartCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerRestartCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerRestartCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerRestartCountOpts
	} else {
		opt = append(opt, newContainerRestartCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.restart.count",
		opt...,
	)
	if err != nil {
		return ContainerRestartCount{noop.Int64UpDownCounter{}}, err
	}
	return ContainerRestartCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerRestartCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerRestartCount) Name() string {
	return "k8s.container.restart.count"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerRestartCount) Unit() string {
	return "{restart}"
}

// Description returns the semantic convention description of the instrument
func (ContainerRestartCount) Description() string {
	return "Describes how many times the container has restarted (since the last counter reset)."
}

// Add adds incr to the existing count for attrs.
//
// This value is pulled directly from the K8s API and the value can go
// indefinitely high and be reset to 0
// at any time depending on how your kubelet is configured to prune dead
// containers.
// It is best to not depend too much on the exact value but rather look at it as
// either == 0, in which case you can conclude there were no restarts in the
// recent past, or > 0, in which case
// you can conclude there were restarts in the recent past, and not try and
// analyze the value beyond that.
func (m ContainerRestartCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This value is pulled directly from the K8s API and the value can go
// indefinitely high and be reset to 0
// at any time depending on how your kubelet is configured to prune dead
// containers.
// It is best to not depend too much on the exact value but rather look at it as
// either == 0, in which case you can conclude there were no restarts in the
// recent past, or > 0, in which case
// you can conclude there were restarts in the recent past, and not try and
// analyze the value beyond that.
func (m ContainerRestartCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStatusReason is an instrument used to record metric values conforming
// to the "k8s.container.status.reason" semantic conventions. It represents the
// describes the number of K8s containers that are currently in a state for a
// given reason.
type ContainerStatusReason struct {
	metric.Int64UpDownCounter
}

var newContainerStatusReasonOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the number of K8s containers that are currently in a state for a given reason."),
	metric.WithUnit("{container}"),
}

// NewContainerStatusReason returns a new ContainerStatusReason instrument.
func NewContainerStatusReason(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStatusReason, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStatusReason{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStatusReasonOpts
	} else {
		opt = append(opt, newContainerStatusReasonOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.status.reason",
		opt...,
	)
	if err != nil {
		return ContainerStatusReason{noop.Int64UpDownCounter{}}, err
	}
	return ContainerStatusReason{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStatusReason) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStatusReason) Name() string {
	return "k8s.container.status.reason"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStatusReason) Unit() string {
	return "{container}"
}

// Description returns the semantic convention description of the instrument
func (ContainerStatusReason) Description() string {
	return "Describes the number of K8s containers that are currently in a state for a given reason."
}

// Add adds incr to the existing count for attrs.
//
// The containerStatusReason is the the reason for the container state.
// Corresponds to the `reason` field of the: [K8s ContainerStateWaiting] or
// [K8s ContainerStateTerminated]
//
// All possible container state reasons will be reported at each time interval to
// avoid missing metrics.
// Only the value corresponding to the current state reason will be non-zero.
//
// [K8s ContainerStateWaiting]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstatewaiting-v1-core
// [K8s ContainerStateTerminated]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstateterminated-v1-core
func (m ContainerStatusReason) Add(
	ctx context.Context,
	incr int64,
	containerStatusReason ContainerStatusReasonAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.container.status.reason", string(containerStatusReason)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// All possible container state reasons will be reported at each time interval to
// avoid missing metrics.
// Only the value corresponding to the current state reason will be non-zero.
func (m ContainerStatusReason) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStatusState is an instrument used to record metric values conforming
// to the "k8s.container.status.state" semantic conventions. It represents the
// describes the number of K8s containers that are currently in a given state.
type ContainerStatusState struct {
	metric.Int64UpDownCounter
}

var newContainerStatusStateOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the number of K8s containers that are currently in a given state."),
	metric.WithUnit("{container}"),
}

// NewContainerStatusState returns a new ContainerStatusState instrument.
func NewContainerStatusState(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStatusState, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStatusState{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStatusStateOpts
	} else {
		opt = append(opt, newContainerStatusStateOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.status.state",
		opt...,
	)
	if err != nil {
		return ContainerStatusState{noop.Int64UpDownCounter{}}, err
	}
	return ContainerStatusState{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStatusState) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStatusState) Name() string {
	return "k8s.container.status.state"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStatusState) Unit() string {
	return "{container}"
}

// Description returns the semantic convention description of the instrument
func (ContainerStatusState) Description() string {
	return "Describes the number of K8s containers that are currently in a given state."
}

// Add adds incr to the existing count for attrs.
//
// The containerStatusState is the the state of the container.
// [K8s ContainerState]
//
// All possible container states will be reported at each time interval to avoid
// missing metrics.
// Only the value corresponding to the current state will be non-zero.
//
// [K8s ContainerState]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstate-v1-core
func (m ContainerStatusState) Add(
	ctx context.Context,
	incr int64,
	containerStatusState ContainerStatusStateAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.container.status.state", string(containerStatusState)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// All possible container states will be reported at each time interval to avoid
// missing metrics.
// Only the value corresponding to the current state will be non-zero.
func (m ContainerStatusState) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStorageLimit is an instrument used to record metric values conforming
// to the "k8s.container.storage.limit" semantic conventions. It represents the
// maximum storage resource limit set for the container.
type ContainerStorageLimit struct {
	metric.Int64UpDownCounter
}

var newContainerStorageLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum storage resource limit set for the container."),
	metric.WithUnit("By"),
}

// NewContainerStorageLimit returns a new ContainerStorageLimit instrument.
func NewContainerStorageLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStorageLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStorageLimit{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStorageLimitOpts
	} else {
		opt = append(opt, newContainerStorageLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.storage.limit",
		opt...,
	)
	if err != nil {
		return ContainerStorageLimit{noop.Int64UpDownCounter{}}, err
	}
	return ContainerStorageLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStorageLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStorageLimit) Name() string {
	return "k8s.container.storage.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStorageLimit) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerStorageLimit) Description() string {
	return "Maximum storage resource limit set for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerStorageLimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerStorageLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStorageRequest is an instrument used to record metric values
// conforming to the "k8s.container.storage.request" semantic conventions. It
// represents the storage resource requested for the container.
type ContainerStorageRequest struct {
	metric.Int64UpDownCounter
}

var newContainerStorageRequestOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Storage resource requested for the container."),
	metric.WithUnit("By"),
}

// NewContainerStorageRequest returns a new ContainerStorageRequest instrument.
func NewContainerStorageRequest(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStorageRequest, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStorageRequest{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStorageRequestOpts
	} else {
		opt = append(opt, newContainerStorageRequestOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.storage.request",
		opt...,
	)
	if err != nil {
		return ContainerStorageRequest{noop.Int64UpDownCounter{}}, err
	}
	return ContainerStorageRequest{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStorageRequest) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStorageRequest) Name() string {
	return "k8s.container.storage.request"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStorageRequest) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerStorageRequest) Description() string {
	return "Storage resource requested for the container."
}

// Add adds incr to the existing count for attrs.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerStorageRequest) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// See
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#resourcerequirements-v1-core
// for details.
func (m ContainerStorageRequest) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// CronJobJobActive is an instrument used to record metric values conforming to
// the "k8s.cronjob.job.active" semantic conventions. It represents the number of
// actively running jobs for a cronjob.
type CronJobJobActive struct {
	metric.Int64UpDownCounter
}

var newCronJobJobActiveOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of actively running jobs for a cronjob."),
	metric.WithUnit("{job}"),
}

// NewCronJobJobActive returns a new CronJobJobActive instrument.
func NewCronJobJobActive(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (CronJobJobActive, error) {
	// Check if the meter is nil.
	if m == nil {
		return CronJobJobActive{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newCronJobJobActiveOpts
	} else {
		opt = append(opt, newCronJobJobActiveOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.cronjob.job.active",
		opt...,
	)
	if err != nil {
		return CronJobJobActive{noop.Int64UpDownCounter{}}, err
	}
	return CronJobJobActive{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CronJobJobActive) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (CronJobJobActive) Name() string {
	return "k8s.cronjob.job.active"
}

// Unit returns the semantic convention unit of the instrument
func (CronJobJobActive) Unit() string {
	return "{job}"
}

// Description returns the semantic convention description of the instrument
func (CronJobJobActive) Description() string {
	return "The number of actively running jobs for a cronjob."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `active` field of the
// [K8s CronJobStatus].
//
// [K8s CronJobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#cronjobstatus-v1-batch
func (m CronJobJobActive) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `active` field of the
// [K8s CronJobStatus].
//
// [K8s CronJobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#cronjobstatus-v1-batch
func (m CronJobJobActive) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeCurrentScheduled is an instrument used to record metric values
// conforming to the "k8s.daemonset.node.current_scheduled" semantic conventions.
// It represents the number of nodes that are running at least 1 daemon pod and
// are supposed to run the daemon pod.
type DaemonSetNodeCurrentScheduled struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeCurrentScheduledOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeCurrentScheduled returns a new DaemonSetNodeCurrentScheduled
// instrument.
func NewDaemonSetNodeCurrentScheduled(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetNodeCurrentScheduled, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeCurrentScheduled{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeCurrentScheduledOpts
	} else {
		opt = append(opt, newDaemonSetNodeCurrentScheduledOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.node.current_scheduled",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeCurrentScheduled{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetNodeCurrentScheduled{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeCurrentScheduled) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeCurrentScheduled) Name() string {
	return "k8s.daemonset.node.current_scheduled"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeCurrentScheduled) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeCurrentScheduled) Description() string {
	return "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `currentNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeCurrentScheduled) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeCurrentScheduled) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeDesiredScheduled is an instrument used to record metric values
// conforming to the "k8s.daemonset.node.desired_scheduled" semantic conventions.
// It represents the number of nodes that should be running the daemon pod
// (including nodes currently running the daemon pod).
type DaemonSetNodeDesiredScheduled struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeDesiredScheduledOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeDesiredScheduled returns a new DaemonSetNodeDesiredScheduled
// instrument.
func NewDaemonSetNodeDesiredScheduled(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetNodeDesiredScheduled, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeDesiredScheduled{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeDesiredScheduledOpts
	} else {
		opt = append(opt, newDaemonSetNodeDesiredScheduledOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.node.desired_scheduled",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeDesiredScheduled{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetNodeDesiredScheduled{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeDesiredScheduled) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeDesiredScheduled) Name() string {
	return "k8s.daemonset.node.desired_scheduled"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeDesiredScheduled) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeDesiredScheduled) Description() string {
	return "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `desiredNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeDesiredScheduled) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `desiredNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeDesiredScheduled) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeMisscheduled is an instrument used to record metric values
// conforming to the "k8s.daemonset.node.misscheduled" semantic conventions. It
// represents the number of nodes that are running the daemon pod, but are not
// supposed to run the daemon pod.
type DaemonSetNodeMisscheduled struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeMisscheduledOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeMisscheduled returns a new DaemonSetNodeMisscheduled
// instrument.
func NewDaemonSetNodeMisscheduled(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetNodeMisscheduled, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeMisscheduled{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeMisscheduledOpts
	} else {
		opt = append(opt, newDaemonSetNodeMisscheduledOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.node.misscheduled",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeMisscheduled{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetNodeMisscheduled{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeMisscheduled) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeMisscheduled) Name() string {
	return "k8s.daemonset.node.misscheduled"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeMisscheduled) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeMisscheduled) Description() string {
	return "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `numberMisscheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeMisscheduled) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `numberMisscheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeMisscheduled) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeReady is an instrument used to record metric values conforming to
// the "k8s.daemonset.node.ready" semantic conventions. It represents the number
// of nodes that should be running the daemon pod and have one or more of the
// daemon pod running and ready.
type DaemonSetNodeReady struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeReadyOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeReady returns a new DaemonSetNodeReady instrument.
func NewDaemonSetNodeReady(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DaemonSetNodeReady, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeReady{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeReadyOpts
	} else {
		opt = append(opt, newDaemonSetNodeReadyOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.daemonset.node.ready",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeReady{noop.Int64UpDownCounter{}}, err
	}
	return DaemonSetNodeReady{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeReady) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeReady) Name() string {
	return "k8s.daemonset.node.ready"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeReady) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeReady) Description() string {
	return "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `numberReady` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeReady) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `numberReady` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetNodeReady) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DeploymentPodAvailable is an instrument used to record metric values
// conforming to the "k8s.deployment.pod.available" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this deployment.
type DeploymentPodAvailable struct {
	metric.Int64UpDownCounter
}

var newDeploymentPodAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment."),
	metric.WithUnit("{pod}"),
}

// NewDeploymentPodAvailable returns a new DeploymentPodAvailable instrument.
func NewDeploymentPodAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DeploymentPodAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DeploymentPodAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDeploymentPodAvailableOpts
	} else {
		opt = append(opt, newDeploymentPodAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.deployment.pod.available",
		opt...,
	)
	if err != nil {
		return DeploymentPodAvailable{noop.Int64UpDownCounter{}}, err
	}
	return DeploymentPodAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DeploymentPodAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DeploymentPodAvailable) Name() string {
	return "k8s.deployment.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (DeploymentPodAvailable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (DeploymentPodAvailable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s DeploymentStatus].
//
// [K8s DeploymentStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentstatus-v1-apps
func (m DeploymentPodAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s DeploymentStatus].
//
// [K8s DeploymentStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentstatus-v1-apps
func (m DeploymentPodAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DeploymentPodDesired is an instrument used to record metric values conforming
// to the "k8s.deployment.pod.desired" semantic conventions. It represents the
// number of desired replica pods in this deployment.
type DeploymentPodDesired struct {
	metric.Int64UpDownCounter
}

var newDeploymentPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this deployment."),
	metric.WithUnit("{pod}"),
}

// NewDeploymentPodDesired returns a new DeploymentPodDesired instrument.
func NewDeploymentPodDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (DeploymentPodDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return DeploymentPodDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDeploymentPodDesiredOpts
	} else {
		opt = append(opt, newDeploymentPodDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.deployment.pod.desired",
		opt...,
	)
	if err != nil {
		return DeploymentPodDesired{noop.Int64UpDownCounter{}}, err
	}
	return DeploymentPodDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DeploymentPodDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DeploymentPodDesired) Name() string {
	return "k8s.deployment.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (DeploymentPodDesired) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (DeploymentPodDesired) Description() string {
	return "Number of desired replica pods in this deployment."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s DeploymentSpec].
//
// [K8s DeploymentSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentspec-v1-apps
func (m DeploymentPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s DeploymentSpec].
//
// [K8s DeploymentSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentspec-v1-apps
func (m DeploymentPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAMetricTargetCPUAverageUtilization is an instrument used to record metric
// values conforming to the "k8s.hpa.metric.target.cpu.average_utilization"
// semantic conventions. It represents the target average utilization, in
// percentage, for CPU resource in HPA config.
type HPAMetricTargetCPUAverageUtilization struct {
	metric.Int64Gauge
}

var newHPAMetricTargetCPUAverageUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Target average utilization, in percentage, for CPU resource in HPA config."),
	metric.WithUnit("1"),
}

// NewHPAMetricTargetCPUAverageUtilization returns a new
// HPAMetricTargetCPUAverageUtilization instrument.
func NewHPAMetricTargetCPUAverageUtilization(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (HPAMetricTargetCPUAverageUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMetricTargetCPUAverageUtilization{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAMetricTargetCPUAverageUtilizationOpts
	} else {
		opt = append(opt, newHPAMetricTargetCPUAverageUtilizationOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.hpa.metric.target.cpu.average_utilization",
		opt...,
	)
	if err != nil {
		return HPAMetricTargetCPUAverageUtilization{noop.Int64Gauge{}}, err
	}
	return HPAMetricTargetCPUAverageUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMetricTargetCPUAverageUtilization) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (HPAMetricTargetCPUAverageUtilization) Name() string {
	return "k8s.hpa.metric.target.cpu.average_utilization"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMetricTargetCPUAverageUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (HPAMetricTargetCPUAverageUtilization) Description() string {
	return "Target average utilization, in percentage, for CPU resource in HPA config."
}

// Record records val to the current distribution for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric aligns with the `averageUtilization` field of the
// [K8s HPA MetricTarget].
// If the type of the metric is [`ContainerResource`],
// the `k8s.container.name` attribute MUST be set to identify the specific
// container within the pod to which the metric applies.
//
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageUtilization) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

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

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// This metric aligns with the `averageUtilization` field of the
// [K8s HPA MetricTarget].
// If the type of the metric is [`ContainerResource`],
// the `k8s.container.name` attribute MUST be set to identify the specific
// container within the pod to which the metric applies.
//
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrContainerName returns an optional attribute for the "k8s.container.name"
// semantic convention. It represents the name of the Container from Pod
// specification, must be unique within a Pod. Container runtime usually uses
// different globally unique name (`container.name`).
func (HPAMetricTargetCPUAverageUtilization) AttrContainerName(val string) attribute.KeyValue {
	return attribute.String("k8s.container.name", val)
}

// AttrHPAMetricType returns an optional attribute for the "k8s.hpa.metric.type"
// semantic convention. It represents the type of metric source for the
// horizontal pod autoscaler.
func (HPAMetricTargetCPUAverageUtilization) AttrHPAMetricType(val string) attribute.KeyValue {
	return attribute.String("k8s.hpa.metric.type", val)
}

// HPAMetricTargetCPUAverageValue is an instrument used to record metric values
// conforming to the "k8s.hpa.metric.target.cpu.average_value" semantic
// conventions. It represents the target average value for CPU resource in HPA
// config.
type HPAMetricTargetCPUAverageValue struct {
	metric.Int64Gauge
}

var newHPAMetricTargetCPUAverageValueOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Target average value for CPU resource in HPA config."),
	metric.WithUnit("{cpu}"),
}

// NewHPAMetricTargetCPUAverageValue returns a new HPAMetricTargetCPUAverageValue
// instrument.
func NewHPAMetricTargetCPUAverageValue(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (HPAMetricTargetCPUAverageValue, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMetricTargetCPUAverageValue{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAMetricTargetCPUAverageValueOpts
	} else {
		opt = append(opt, newHPAMetricTargetCPUAverageValueOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.hpa.metric.target.cpu.average_value",
		opt...,
	)
	if err != nil {
		return HPAMetricTargetCPUAverageValue{noop.Int64Gauge{}}, err
	}
	return HPAMetricTargetCPUAverageValue{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMetricTargetCPUAverageValue) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (HPAMetricTargetCPUAverageValue) Name() string {
	return "k8s.hpa.metric.target.cpu.average_value"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMetricTargetCPUAverageValue) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (HPAMetricTargetCPUAverageValue) Description() string {
	return "Target average value for CPU resource in HPA config."
}

// Record records val to the current distribution for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric aligns with the `averageValue` field of the
// [K8s HPA MetricTarget].
// If the type of the metric is [`ContainerResource`],
// the `k8s.container.name` attribute MUST be set to identify the specific
// container within the pod to which the metric applies.
//
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageValue) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

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

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// This metric aligns with the `averageValue` field of the
// [K8s HPA MetricTarget].
// If the type of the metric is [`ContainerResource`],
// the `k8s.container.name` attribute MUST be set to identify the specific
// container within the pod to which the metric applies.
//
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageValue) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrContainerName returns an optional attribute for the "k8s.container.name"
// semantic convention. It represents the name of the Container from Pod
// specification, must be unique within a Pod. Container runtime usually uses
// different globally unique name (`container.name`).
func (HPAMetricTargetCPUAverageValue) AttrContainerName(val string) attribute.KeyValue {
	return attribute.String("k8s.container.name", val)
}

// AttrHPAMetricType returns an optional attribute for the "k8s.hpa.metric.type"
// semantic convention. It represents the type of metric source for the
// horizontal pod autoscaler.
func (HPAMetricTargetCPUAverageValue) AttrHPAMetricType(val string) attribute.KeyValue {
	return attribute.String("k8s.hpa.metric.type", val)
}

// HPAMetricTargetCPUValue is an instrument used to record metric values
// conforming to the "k8s.hpa.metric.target.cpu.value" semantic conventions. It
// represents the target value for CPU resource in HPA config.
type HPAMetricTargetCPUValue struct {
	metric.Int64Gauge
}

var newHPAMetricTargetCPUValueOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Target value for CPU resource in HPA config."),
	metric.WithUnit("{cpu}"),
}

// NewHPAMetricTargetCPUValue returns a new HPAMetricTargetCPUValue instrument.
func NewHPAMetricTargetCPUValue(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (HPAMetricTargetCPUValue, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMetricTargetCPUValue{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAMetricTargetCPUValueOpts
	} else {
		opt = append(opt, newHPAMetricTargetCPUValueOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.hpa.metric.target.cpu.value",
		opt...,
	)
	if err != nil {
		return HPAMetricTargetCPUValue{noop.Int64Gauge{}}, err
	}
	return HPAMetricTargetCPUValue{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMetricTargetCPUValue) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (HPAMetricTargetCPUValue) Name() string {
	return "k8s.hpa.metric.target.cpu.value"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMetricTargetCPUValue) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (HPAMetricTargetCPUValue) Description() string {
	return "Target value for CPU resource in HPA config."
}

// Record records val to the current distribution for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric aligns with the `value` field of the
// [K8s HPA MetricTarget].
// If the type of the metric is [`ContainerResource`],
// the `k8s.container.name` attribute MUST be set to identify the specific
// container within the pod to which the metric applies.
//
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUValue) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

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

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// This metric aligns with the `value` field of the
// [K8s HPA MetricTarget].
// If the type of the metric is [`ContainerResource`],
// the `k8s.container.name` attribute MUST be set to identify the specific
// container within the pod to which the metric applies.
//
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUValue) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrContainerName returns an optional attribute for the "k8s.container.name"
// semantic convention. It represents the name of the Container from Pod
// specification, must be unique within a Pod. Container runtime usually uses
// different globally unique name (`container.name`).
func (HPAMetricTargetCPUValue) AttrContainerName(val string) attribute.KeyValue {
	return attribute.String("k8s.container.name", val)
}

// AttrHPAMetricType returns an optional attribute for the "k8s.hpa.metric.type"
// semantic convention. It represents the type of metric source for the
// horizontal pod autoscaler.
func (HPAMetricTargetCPUValue) AttrHPAMetricType(val string) attribute.KeyValue {
	return attribute.String("k8s.hpa.metric.type", val)
}

// HPAPodCurrent is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.current" semantic conventions. It represents the current number
// of replica pods managed by this horizontal pod autoscaler, as last seen by the
// autoscaler.
type HPAPodCurrent struct {
	metric.Int64UpDownCounter
}

var newHPAPodCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodCurrent returns a new HPAPodCurrent instrument.
func NewHPAPodCurrent(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPAPodCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodCurrent{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodCurrentOpts
	} else {
		opt = append(opt, newHPAPodCurrentOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.pod.current",
		opt...,
	)
	if err != nil {
		return HPAPodCurrent{noop.Int64UpDownCounter{}}, err
	}
	return HPAPodCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodCurrent) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodCurrent) Name() string {
	return "k8s.hpa.pod.current"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodCurrent) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodCurrent) Description() string {
	return "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAPodDesired is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.desired" semantic conventions. It represents the desired number
// of replica pods managed by this horizontal pod autoscaler, as last calculated
// by the autoscaler.
type HPAPodDesired struct {
	metric.Int64UpDownCounter
}

var newHPAPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodDesired returns a new HPAPodDesired instrument.
func NewHPAPodDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPAPodDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodDesiredOpts
	} else {
		opt = append(opt, newHPAPodDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.pod.desired",
		opt...,
	)
	if err != nil {
		return HPAPodDesired{noop.Int64UpDownCounter{}}, err
	}
	return HPAPodDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodDesired) Name() string {
	return "k8s.hpa.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodDesired) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodDesired) Description() string {
	return "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `desiredReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `desiredReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAPodMax is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.max" semantic conventions. It represents the upper limit for the
// number of replica pods to which the autoscaler can scale up.
type HPAPodMax struct {
	metric.Int64UpDownCounter
}

var newHPAPodMaxOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The upper limit for the number of replica pods to which the autoscaler can scale up."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodMax returns a new HPAPodMax instrument.
func NewHPAPodMax(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPAPodMax, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodMax{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodMaxOpts
	} else {
		opt = append(opt, newHPAPodMaxOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.pod.max",
		opt...,
	)
	if err != nil {
		return HPAPodMax{noop.Int64UpDownCounter{}}, err
	}
	return HPAPodMax{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodMax) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodMax) Name() string {
	return "k8s.hpa.pod.max"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodMax) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodMax) Description() string {
	return "The upper limit for the number of replica pods to which the autoscaler can scale up."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `maxReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMax) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `maxReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMax) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAPodMin is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.min" semantic conventions. It represents the lower limit for the
// number of replica pods to which the autoscaler can scale down.
type HPAPodMin struct {
	metric.Int64UpDownCounter
}

var newHPAPodMinOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The lower limit for the number of replica pods to which the autoscaler can scale down."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodMin returns a new HPAPodMin instrument.
func NewHPAPodMin(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (HPAPodMin, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodMin{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodMinOpts
	} else {
		opt = append(opt, newHPAPodMinOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.hpa.pod.min",
		opt...,
	)
	if err != nil {
		return HPAPodMin{noop.Int64UpDownCounter{}}, err
	}
	return HPAPodMin{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodMin) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodMin) Name() string {
	return "k8s.hpa.pod.min"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodMin) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodMin) Description() string {
	return "The lower limit for the number of replica pods to which the autoscaler can scale down."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `minReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMin) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `minReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMin) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodActive is an instrument used to record metric values conforming to the
// "k8s.job.pod.active" semantic conventions. It represents the number of pending
// and actively running pods for a job.
type JobPodActive struct {
	metric.Int64UpDownCounter
}

var newJobPodActiveOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of pending and actively running pods for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodActive returns a new JobPodActive instrument.
func NewJobPodActive(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobPodActive, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodActive{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodActiveOpts
	} else {
		opt = append(opt, newJobPodActiveOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.pod.active",
		opt...,
	)
	if err != nil {
		return JobPodActive{noop.Int64UpDownCounter{}}, err
	}
	return JobPodActive{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodActive) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodActive) Name() string {
	return "k8s.job.pod.active"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodActive) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodActive) Description() string {
	return "The number of pending and actively running pods for a job."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `active` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobPodActive) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `active` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobPodActive) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodDesiredSuccessful is an instrument used to record metric values
// conforming to the "k8s.job.pod.desired_successful" semantic conventions. It
// represents the desired number of successfully finished pods the job should be
// run with.
type JobPodDesiredSuccessful struct {
	metric.Int64UpDownCounter
}

var newJobPodDesiredSuccessfulOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The desired number of successfully finished pods the job should be run with."),
	metric.WithUnit("{pod}"),
}

// NewJobPodDesiredSuccessful returns a new JobPodDesiredSuccessful instrument.
func NewJobPodDesiredSuccessful(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobPodDesiredSuccessful, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodDesiredSuccessful{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodDesiredSuccessfulOpts
	} else {
		opt = append(opt, newJobPodDesiredSuccessfulOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.pod.desired_successful",
		opt...,
	)
	if err != nil {
		return JobPodDesiredSuccessful{noop.Int64UpDownCounter{}}, err
	}
	return JobPodDesiredSuccessful{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodDesiredSuccessful) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodDesiredSuccessful) Name() string {
	return "k8s.job.pod.desired_successful"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodDesiredSuccessful) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodDesiredSuccessful) Description() string {
	return "The desired number of successfully finished pods the job should be run with."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `completions` field of the
// [K8s JobSpec]..
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
func (m JobPodDesiredSuccessful) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `completions` field of the
// [K8s JobSpec]..
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
func (m JobPodDesiredSuccessful) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodFailed is an instrument used to record metric values conforming to the
// "k8s.job.pod.failed" semantic conventions. It represents the number of pods
// which reached phase Failed for a job.
type JobPodFailed struct {
	metric.Int64UpDownCounter
}

var newJobPodFailedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of pods which reached phase Failed for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodFailed returns a new JobPodFailed instrument.
func NewJobPodFailed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobPodFailed, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodFailed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodFailedOpts
	} else {
		opt = append(opt, newJobPodFailedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.pod.failed",
		opt...,
	)
	if err != nil {
		return JobPodFailed{noop.Int64UpDownCounter{}}, err
	}
	return JobPodFailed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodFailed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodFailed) Name() string {
	return "k8s.job.pod.failed"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodFailed) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodFailed) Description() string {
	return "The number of pods which reached phase Failed for a job."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `failed` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobPodFailed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `failed` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobPodFailed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodMaxParallel is an instrument used to record metric values conforming to
// the "k8s.job.pod.max_parallel" semantic conventions. It represents the max
// desired number of pods the job should run at any given time.
type JobPodMaxParallel struct {
	metric.Int64UpDownCounter
}

var newJobPodMaxParallelOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The max desired number of pods the job should run at any given time."),
	metric.WithUnit("{pod}"),
}

// NewJobPodMaxParallel returns a new JobPodMaxParallel instrument.
func NewJobPodMaxParallel(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobPodMaxParallel, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodMaxParallel{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodMaxParallelOpts
	} else {
		opt = append(opt, newJobPodMaxParallelOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.pod.max_parallel",
		opt...,
	)
	if err != nil {
		return JobPodMaxParallel{noop.Int64UpDownCounter{}}, err
	}
	return JobPodMaxParallel{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodMaxParallel) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodMaxParallel) Name() string {
	return "k8s.job.pod.max_parallel"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodMaxParallel) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodMaxParallel) Description() string {
	return "The max desired number of pods the job should run at any given time."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `parallelism` field of the
// [K8s JobSpec].
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
func (m JobPodMaxParallel) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `parallelism` field of the
// [K8s JobSpec].
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
func (m JobPodMaxParallel) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodSuccessful is an instrument used to record metric values conforming to
// the "k8s.job.pod.successful" semantic conventions. It represents the number of
// pods which reached phase Succeeded for a job.
type JobPodSuccessful struct {
	metric.Int64UpDownCounter
}

var newJobPodSuccessfulOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of pods which reached phase Succeeded for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodSuccessful returns a new JobPodSuccessful instrument.
func NewJobPodSuccessful(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (JobPodSuccessful, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodSuccessful{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodSuccessfulOpts
	} else {
		opt = append(opt, newJobPodSuccessfulOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.job.pod.successful",
		opt...,
	)
	if err != nil {
		return JobPodSuccessful{noop.Int64UpDownCounter{}}, err
	}
	return JobPodSuccessful{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodSuccessful) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodSuccessful) Name() string {
	return "k8s.job.pod.successful"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodSuccessful) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodSuccessful) Description() string {
	return "The number of pods which reached phase Succeeded for a job."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `succeeded` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobPodSuccessful) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `succeeded` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobPodSuccessful) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NamespacePhase is an instrument used to record metric values conforming to the
// "k8s.namespace.phase" semantic conventions. It represents the describes number
// of K8s namespaces that are currently in a given phase.
type NamespacePhase struct {
	metric.Int64UpDownCounter
}

var newNamespacePhaseOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes number of K8s namespaces that are currently in a given phase."),
	metric.WithUnit("{namespace}"),
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

	if len(opt) == 0 {
		opt = newNamespacePhaseOpts
	} else {
		opt = append(opt, newNamespacePhaseOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.namespace.phase",
		opt...,
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

// Add adds incr to the existing count for attrs.
//
// The namespacePhase is the the phase of the K8s namespace.
func (m NamespacePhase) Add(
	ctx context.Context,
	incr int64,
	namespacePhase NamespacePhaseAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
func (m NamespacePhase) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeConditionStatus is an instrument used to record metric values conforming
// to the "k8s.node.condition.status" semantic conventions. It represents the
// describes the condition of a particular Node.
type NodeConditionStatus struct {
	metric.Int64UpDownCounter
}

var newNodeConditionStatusOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the condition of a particular Node."),
	metric.WithUnit("{node}"),
}

// NewNodeConditionStatus returns a new NodeConditionStatus instrument.
func NewNodeConditionStatus(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeConditionStatus, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeConditionStatus{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeConditionStatusOpts
	} else {
		opt = append(opt, newNodeConditionStatusOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.condition.status",
		opt...,
	)
	if err != nil {
		return NodeConditionStatus{noop.Int64UpDownCounter{}}, err
	}
	return NodeConditionStatus{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeConditionStatus) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeConditionStatus) Name() string {
	return "k8s.node.condition.status"
}

// Unit returns the semantic convention unit of the instrument
func (NodeConditionStatus) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (NodeConditionStatus) Description() string {
	return "Describes the condition of a particular Node."
}

// Add adds incr to the existing count for attrs.
//
// The nodeConditionStatus is the the status of the condition, one of True,
// False, Unknown.
//
// The nodeConditionType is the the condition type of a K8s Node.
//
// All possible node condition pairs (type and status) will be reported at each
// time interval to avoid missing metrics. Condition pairs corresponding to the
// current conditions' statuses will be non-zero.
func (m NodeConditionStatus) Add(
	ctx context.Context,
	incr int64,
	nodeConditionStatus NodeConditionStatusAttr,
	nodeConditionType NodeConditionTypeAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.node.condition.status", string(nodeConditionStatus)),
				attribute.String("k8s.node.condition.type", string(nodeConditionType)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// All possible node condition pairs (type and status) will be reported at each
// time interval to avoid missing metrics. Condition pairs corresponding to the
// current conditions' statuses will be non-zero.
func (m NodeConditionStatus) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeCPUAllocatable is an instrument used to record metric values conforming to
// the "k8s.node.cpu.allocatable" semantic conventions. It represents the amount
// of cpu allocatable on the node.
type NodeCPUAllocatable struct {
	metric.Int64UpDownCounter
}

var newNodeCPUAllocatableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Amount of cpu allocatable on the node."),
	metric.WithUnit("{cpu}"),
}

// NewNodeCPUAllocatable returns a new NodeCPUAllocatable instrument.
func NewNodeCPUAllocatable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeCPUAllocatable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUAllocatable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeCPUAllocatableOpts
	} else {
		opt = append(opt, newNodeCPUAllocatableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.cpu.allocatable",
		opt...,
	)
	if err != nil {
		return NodeCPUAllocatable{noop.Int64UpDownCounter{}}, err
	}
	return NodeCPUAllocatable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUAllocatable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeCPUAllocatable) Name() string {
	return "k8s.node.cpu.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodeCPUAllocatable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (NodeCPUAllocatable) Description() string {
	return "Amount of cpu allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeCPUAllocatable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
func (m NodeCPUAllocatable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeCPUTime is an instrument used to record metric values conforming to the
// "k8s.node.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type NodeCPUTime struct {
	metric.Float64Counter
}

var newNodeCPUTimeOpts = []metric.Float64CounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
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

	if len(opt) == 0 {
		opt = newNodeCPUTimeOpts
	} else {
		opt = append(opt, newNodeCPUTimeOpts...)
	}

	i, err := m.Float64Counter(
		"k8s.node.cpu.time",
		opt...,
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
	return "Total CPU time consumed."
}

// Add adds incr to the existing count for attrs.
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

// AddSet adds incr to the existing count for set.
//
// Total CPU time consumed by the specific Node on all available CPU cores
func (m NodeCPUTime) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// NodeCPUUsage is an instrument used to record metric values conforming to the
// "k8s.node.cpu.usage" semantic conventions. It represents the node's CPU usage,
// measured in cpus. Range from 0 to the number of allocatable CPUs.
type NodeCPUUsage struct {
	metric.Int64Gauge
}

var newNodeCPUUsageOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
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

	if len(opt) == 0 {
		opt = newNodeCPUUsageOpts
	} else {
		opt = append(opt, newNodeCPUUsageOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.node.cpu.usage",
		opt...,
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
	return "Node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."
}

// Record records val to the current distribution for attrs.
//
// CPU usage of the specific Node on all available CPU cores, averaged over the
// sample window
func (m NodeCPUUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// CPU usage of the specific Node on all available CPU cores, averaged over the
// sample window
func (m NodeCPUUsage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// NodeEphemeralStorageAllocatable is an instrument used to record metric values
// conforming to the "k8s.node.ephemeral_storage.allocatable" semantic
// conventions. It represents the amount of ephemeral-storage allocatable on the
// node.
type NodeEphemeralStorageAllocatable struct {
	metric.Int64UpDownCounter
}

var newNodeEphemeralStorageAllocatableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Amount of ephemeral-storage allocatable on the node."),
	metric.WithUnit("By"),
}

// NewNodeEphemeralStorageAllocatable returns a new
// NodeEphemeralStorageAllocatable instrument.
func NewNodeEphemeralStorageAllocatable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeEphemeralStorageAllocatable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeEphemeralStorageAllocatable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeEphemeralStorageAllocatableOpts
	} else {
		opt = append(opt, newNodeEphemeralStorageAllocatableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.ephemeral_storage.allocatable",
		opt...,
	)
	if err != nil {
		return NodeEphemeralStorageAllocatable{noop.Int64UpDownCounter{}}, err
	}
	return NodeEphemeralStorageAllocatable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeEphemeralStorageAllocatable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeEphemeralStorageAllocatable) Name() string {
	return "k8s.node.ephemeral_storage.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodeEphemeralStorageAllocatable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeEphemeralStorageAllocatable) Description() string {
	return "Amount of ephemeral-storage allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeEphemeralStorageAllocatable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
func (m NodeEphemeralStorageAllocatable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemAvailable is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.available" semantic conventions. It
// represents the node filesystem available bytes.
type NodeFilesystemAvailable struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemAvailable returns a new NodeFilesystemAvailable instrument.
func NewNodeFilesystemAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeFilesystemAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemAvailableOpts
	} else {
		opt = append(opt, newNodeFilesystemAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.available",
		opt...,
	)
	if err != nil {
		return NodeFilesystemAvailable{noop.Int64UpDownCounter{}}, err
	}
	return NodeFilesystemAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemAvailable) Name() string {
	return "k8s.node.filesystem.available"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemAvailable) Description() string {
	return "Node filesystem available bytes."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the
// [FsStats.AvailableBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.AvailableBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemCapacity is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.capacity" semantic conventions. It
// represents the node filesystem capacity.
type NodeFilesystemCapacity struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node filesystem capacity."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemCapacity returns a new NodeFilesystemCapacity instrument.
func NewNodeFilesystemCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeFilesystemCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemCapacity{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemCapacityOpts
	} else {
		opt = append(opt, newNodeFilesystemCapacityOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.capacity",
		opt...,
	)
	if err != nil {
		return NodeFilesystemCapacity{noop.Int64UpDownCounter{}}, err
	}
	return NodeFilesystemCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemCapacity) Name() string {
	return "k8s.node.filesystem.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemCapacity) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemCapacity) Description() string {
	return "Node filesystem capacity."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the
// [FsStats.CapacityBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemCapacity) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.CapacityBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemUsage is an instrument used to record metric values conforming
// to the "k8s.node.filesystem.usage" semantic conventions. It represents the
// node filesystem usage.
type NodeFilesystemUsage struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node filesystem usage."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemUsage returns a new NodeFilesystemUsage instrument.
func NewNodeFilesystemUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeFilesystemUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemUsageOpts
	} else {
		opt = append(opt, newNodeFilesystemUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.usage",
		opt...,
	)
	if err != nil {
		return NodeFilesystemUsage{noop.Int64UpDownCounter{}}, err
	}
	return NodeFilesystemUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemUsage) Name() string {
	return "k8s.node.filesystem.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemUsage) Description() string {
	return "Node filesystem usage."
}

// Add adds incr to the existing count for attrs.
//
// This may not equal capacity - available.
//
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This may not equal capacity - available.
//
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryAllocatable is an instrument used to record metric values conforming
// to the "k8s.node.memory.allocatable" semantic conventions. It represents the
// amount of memory allocatable on the node.
type NodeMemoryAllocatable struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryAllocatableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Amount of memory allocatable on the node."),
	metric.WithUnit("By"),
}

// NewNodeMemoryAllocatable returns a new NodeMemoryAllocatable instrument.
func NewNodeMemoryAllocatable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeMemoryAllocatable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryAllocatable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryAllocatableOpts
	} else {
		opt = append(opt, newNodeMemoryAllocatableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.memory.allocatable",
		opt...,
	)
	if err != nil {
		return NodeMemoryAllocatable{noop.Int64UpDownCounter{}}, err
	}
	return NodeMemoryAllocatable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryAllocatable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryAllocatable) Name() string {
	return "k8s.node.memory.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryAllocatable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryAllocatable) Description() string {
	return "Amount of memory allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeMemoryAllocatable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
func (m NodeMemoryAllocatable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryAvailable is an instrument used to record metric values conforming
// to the "k8s.node.memory.available" semantic conventions. It represents the
// node memory available.
type NodeMemoryAvailable struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node memory available."),
	metric.WithUnit("By"),
}

// NewNodeMemoryAvailable returns a new NodeMemoryAvailable instrument.
func NewNodeMemoryAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeMemoryAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryAvailableOpts
	} else {
		opt = append(opt, newNodeMemoryAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.memory.available",
		opt...,
	)
	if err != nil {
		return NodeMemoryAvailable{noop.Int64UpDownCounter{}}, err
	}
	return NodeMemoryAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryAvailable) Name() string {
	return "k8s.node.memory.available"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryAvailable) Description() string {
	return "Node memory available."
}

// Add adds incr to the existing count for attrs.
//
// Available memory for use. This is defined as the memory limit -
// workingSetBytes. If memory limit is undefined, the available bytes is omitted.
// This metric is derived from the [MemoryStats.AvailableBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// Available memory for use. This is defined as the memory limit -
// workingSetBytes. If memory limit is undefined, the available bytes is omitted.
// This metric is derived from the [MemoryStats.AvailableBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryPagingFaults is an instrument used to record metric values
// conforming to the "k8s.node.memory.paging.faults" semantic conventions. It
// represents the node memory paging faults.
type NodeMemoryPagingFaults struct {
	metric.Int64Counter
}

var newNodeMemoryPagingFaultsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Node memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewNodeMemoryPagingFaults returns a new NodeMemoryPagingFaults instrument.
func NewNodeMemoryPagingFaults(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NodeMemoryPagingFaults, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryPagingFaults{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryPagingFaultsOpts
	} else {
		opt = append(opt, newNodeMemoryPagingFaultsOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.node.memory.paging.faults",
		opt...,
	)
	if err != nil {
		return NodeMemoryPagingFaults{noop.Int64Counter{}}, err
	}
	return NodeMemoryPagingFaults{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryPagingFaults) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryPagingFaults) Name() string {
	return "k8s.node.memory.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryPagingFaults) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryPagingFaults) Description() string {
	return "Node memory paging faults."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Cumulative number of major/minor page faults.
// This metric is derived from the [MemoryStats.PageFaults] and
// [MemoryStats.MajorPageFaults] fields of the [NodeStats.Memory] of the
// Kubelet's stats API.
//
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryPagingFaults) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
//
// Cumulative number of major/minor page faults.
// This metric is derived from the [MemoryStats.PageFaults] and
// [MemoryStats.MajorPageFaults] fields of the [NodeStats.Memory] of the
// Kubelet's stats API.
//
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryPagingFaults) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (NodeMemoryPagingFaults) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// NodeMemoryRss is an instrument used to record metric values conforming to the
// "k8s.node.memory.rss" semantic conventions. It represents the node memory RSS.
type NodeMemoryRss struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryRssOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node memory RSS."),
	metric.WithUnit("By"),
}

// NewNodeMemoryRss returns a new NodeMemoryRss instrument.
func NewNodeMemoryRss(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeMemoryRss, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryRss{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryRssOpts
	} else {
		opt = append(opt, newNodeMemoryRssOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.memory.rss",
		opt...,
	)
	if err != nil {
		return NodeMemoryRss{noop.Int64UpDownCounter{}}, err
	}
	return NodeMemoryRss{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryRss) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryRss) Name() string {
	return "k8s.node.memory.rss"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryRss) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryRss) Description() string {
	return "Node memory RSS."
}

// Add adds incr to the existing count for attrs.
//
// The amount of anonymous and swap cache memory (includes transparent
// hugepages).
// This metric is derived from the [MemoryStats.RSSBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.RSSBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryRss) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// The amount of anonymous and swap cache memory (includes transparent
// hugepages).
// This metric is derived from the [MemoryStats.RSSBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.RSSBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryRss) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryUsage is an instrument used to record metric values conforming to
// the "k8s.node.memory.usage" semantic conventions. It represents the memory
// usage of the Node.
type NodeMemoryUsage struct {
	metric.Int64Gauge
}

var newNodeMemoryUsageOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Memory usage of the Node."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newNodeMemoryUsageOpts
	} else {
		opt = append(opt, newNodeMemoryUsageOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.node.memory.usage",
		opt...,
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
	return "Memory usage of the Node."
}

// Record records val to the current distribution for attrs.
//
// Total memory usage of the Node
func (m NodeMemoryUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Total memory usage of the Node
func (m NodeMemoryUsage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// NodeMemoryWorkingSet is an instrument used to record metric values conforming
// to the "k8s.node.memory.working_set" semantic conventions. It represents the
// node memory working set.
type NodeMemoryWorkingSet struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryWorkingSetOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node memory working set."),
	metric.WithUnit("By"),
}

// NewNodeMemoryWorkingSet returns a new NodeMemoryWorkingSet instrument.
func NewNodeMemoryWorkingSet(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeMemoryWorkingSet, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryWorkingSet{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryWorkingSetOpts
	} else {
		opt = append(opt, newNodeMemoryWorkingSetOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.memory.working_set",
		opt...,
	)
	if err != nil {
		return NodeMemoryWorkingSet{noop.Int64UpDownCounter{}}, err
	}
	return NodeMemoryWorkingSet{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryWorkingSet) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryWorkingSet) Name() string {
	return "k8s.node.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryWorkingSet) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryWorkingSet) Description() string {
	return "Node memory working set."
}

// Add adds incr to the existing count for attrs.
//
// The amount of working set memory. This includes recently accessed memory,
// dirty memory, and kernel memory. WorkingSetBytes is <= UsageBytes.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryWorkingSet) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// The amount of working set memory. This includes recently accessed memory,
// dirty memory, and kernel memory. WorkingSetBytes is <= UsageBytes.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryWorkingSet) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeNetworkErrors is an instrument used to record metric values conforming to
// the "k8s.node.network.errors" semantic conventions. It represents the node
// network errors.
type NodeNetworkErrors struct {
	metric.Int64Counter
}

var newNodeNetworkErrorsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Node network errors."),
	metric.WithUnit("{error}"),
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

	if len(opt) == 0 {
		opt = newNodeNetworkErrorsOpts
	} else {
		opt = append(opt, newNodeNetworkErrorsOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.node.network.errors",
		opt...,
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
	return "Node network errors."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m NodeNetworkErrors) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
func (m NodeNetworkErrors) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
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

var newNodeNetworkIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Network bytes for the Node."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newNodeNetworkIOOpts
	} else {
		opt = append(opt, newNodeNetworkIOOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.node.network.io",
		opt...,
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
	return "Network bytes for the Node."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m NodeNetworkIO) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
func (m NodeNetworkIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
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

// NodePodAllocatable is an instrument used to record metric values conforming to
// the "k8s.node.pod.allocatable" semantic conventions. It represents the amount
// of pods allocatable on the node.
type NodePodAllocatable struct {
	metric.Int64UpDownCounter
}

var newNodePodAllocatableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Amount of pods allocatable on the node."),
	metric.WithUnit("{pod}"),
}

// NewNodePodAllocatable returns a new NodePodAllocatable instrument.
func NewNodePodAllocatable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodePodAllocatable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePodAllocatable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePodAllocatableOpts
	} else {
		opt = append(opt, newNodePodAllocatableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.pod.allocatable",
		opt...,
	)
	if err != nil {
		return NodePodAllocatable{noop.Int64UpDownCounter{}}, err
	}
	return NodePodAllocatable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePodAllocatable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodePodAllocatable) Name() string {
	return "k8s.node.pod.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodePodAllocatable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (NodePodAllocatable) Description() string {
	return "Amount of pods allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodePodAllocatable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
func (m NodePodAllocatable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeUptime is an instrument used to record metric values conforming to the
// "k8s.node.uptime" semantic conventions. It represents the time the Node has
// been running.
type NodeUptime struct {
	metric.Float64Gauge
}

var newNodeUptimeOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The time the Node has been running."),
	metric.WithUnit("s"),
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

	if len(opt) == 0 {
		opt = newNodeUptimeOpts
	} else {
		opt = append(opt, newNodeUptimeOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.node.uptime",
		opt...,
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
	return "The time the Node has been running."
}

// Record records val to the current distribution for attrs.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m NodeUptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m NodeUptime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodCPUTime is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
type PodCPUTime struct {
	metric.Float64Counter
}

var newPodCPUTimeOpts = []metric.Float64CounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
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

	if len(opt) == 0 {
		opt = newPodCPUTimeOpts
	} else {
		opt = append(opt, newPodCPUTimeOpts...)
	}

	i, err := m.Float64Counter(
		"k8s.pod.cpu.time",
		opt...,
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
	return "Total CPU time consumed."
}

// Add adds incr to the existing count for attrs.
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

// AddSet adds incr to the existing count for set.
//
// Total CPU time consumed by the specific Pod on all available CPU cores
func (m PodCPUTime) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// PodCPUUsage is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.usage" semantic conventions. It represents the pod's CPU usage,
// measured in cpus. Range from 0 to the number of allocatable CPUs.
type PodCPUUsage struct {
	metric.Int64Gauge
}

var newPodCPUUsageOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
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

	if len(opt) == 0 {
		opt = newPodCPUUsageOpts
	} else {
		opt = append(opt, newPodCPUUsageOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.pod.cpu.usage",
		opt...,
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
	return "Pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."
}

// Record records val to the current distribution for attrs.
//
// CPU usage of the specific Pod on all available CPU cores, averaged over the
// sample window
func (m PodCPUUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// CPU usage of the specific Pod on all available CPU cores, averaged over the
// sample window
func (m PodCPUUsage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// PodFilesystemAvailable is an instrument used to record metric values
// conforming to the "k8s.pod.filesystem.available" semantic conventions. It
// represents the pod filesystem available bytes.
type PodFilesystemAvailable struct {
	metric.Int64UpDownCounter
}

var newPodFilesystemAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewPodFilesystemAvailable returns a new PodFilesystemAvailable instrument.
func NewPodFilesystemAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodFilesystemAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodFilesystemAvailableOpts
	} else {
		opt = append(opt, newPodFilesystemAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.filesystem.available",
		opt...,
	)
	if err != nil {
		return PodFilesystemAvailable{noop.Int64UpDownCounter{}}, err
	}
	return PodFilesystemAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodFilesystemAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodFilesystemAvailable) Name() string {
	return "k8s.pod.filesystem.available"
}

// Unit returns the semantic convention unit of the instrument
func (PodFilesystemAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodFilesystemAvailable) Description() string {
	return "Pod filesystem available bytes."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the
// [FsStats.AvailableBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.AvailableBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodFilesystemCapacity is an instrument used to record metric values conforming
// to the "k8s.pod.filesystem.capacity" semantic conventions. It represents the
// pod filesystem capacity.
type PodFilesystemCapacity struct {
	metric.Int64UpDownCounter
}

var newPodFilesystemCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod filesystem capacity."),
	metric.WithUnit("By"),
}

// NewPodFilesystemCapacity returns a new PodFilesystemCapacity instrument.
func NewPodFilesystemCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodFilesystemCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemCapacity{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodFilesystemCapacityOpts
	} else {
		opt = append(opt, newPodFilesystemCapacityOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.filesystem.capacity",
		opt...,
	)
	if err != nil {
		return PodFilesystemCapacity{noop.Int64UpDownCounter{}}, err
	}
	return PodFilesystemCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodFilesystemCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodFilesystemCapacity) Name() string {
	return "k8s.pod.filesystem.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PodFilesystemCapacity) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodFilesystemCapacity) Description() string {
	return "Pod filesystem capacity."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the
// [FsStats.CapacityBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemCapacity) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.CapacityBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodFilesystemUsage is an instrument used to record metric values conforming to
// the "k8s.pod.filesystem.usage" semantic conventions. It represents the pod
// filesystem usage.
type PodFilesystemUsage struct {
	metric.Int64UpDownCounter
}

var newPodFilesystemUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod filesystem usage."),
	metric.WithUnit("By"),
}

// NewPodFilesystemUsage returns a new PodFilesystemUsage instrument.
func NewPodFilesystemUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodFilesystemUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodFilesystemUsageOpts
	} else {
		opt = append(opt, newPodFilesystemUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.filesystem.usage",
		opt...,
	)
	if err != nil {
		return PodFilesystemUsage{noop.Int64UpDownCounter{}}, err
	}
	return PodFilesystemUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodFilesystemUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodFilesystemUsage) Name() string {
	return "k8s.pod.filesystem.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodFilesystemUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodFilesystemUsage) Description() string {
	return "Pod filesystem usage."
}

// Add adds incr to the existing count for attrs.
//
// This may not equal capacity - available.
//
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This may not equal capacity - available.
//
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryAvailable is an instrument used to record metric values conforming to
// the "k8s.pod.memory.available" semantic conventions. It represents the pod
// memory available.
type PodMemoryAvailable struct {
	metric.Int64UpDownCounter
}

var newPodMemoryAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod memory available."),
	metric.WithUnit("By"),
}

// NewPodMemoryAvailable returns a new PodMemoryAvailable instrument.
func NewPodMemoryAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryAvailableOpts
	} else {
		opt = append(opt, newPodMemoryAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.available",
		opt...,
	)
	if err != nil {
		return PodMemoryAvailable{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryAvailable) Name() string {
	return "k8s.pod.memory.available"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryAvailable) Description() string {
	return "Pod memory available."
}

// Add adds incr to the existing count for attrs.
//
// Available memory for use. This is defined as the memory limit -
// workingSetBytes. If memory limit is undefined, the available bytes is omitted.
// This metric is derived from the [MemoryStats.AvailableBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// Available memory for use. This is defined as the memory limit -
// workingSetBytes. If memory limit is undefined, the available bytes is omitted.
// This metric is derived from the [MemoryStats.AvailableBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryPagingFaults is an instrument used to record metric values conforming
// to the "k8s.pod.memory.paging.faults" semantic conventions. It represents the
// pod memory paging faults.
type PodMemoryPagingFaults struct {
	metric.Int64Counter
}

var newPodMemoryPagingFaultsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Pod memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewPodMemoryPagingFaults returns a new PodMemoryPagingFaults instrument.
func NewPodMemoryPagingFaults(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (PodMemoryPagingFaults, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryPagingFaults{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryPagingFaultsOpts
	} else {
		opt = append(opt, newPodMemoryPagingFaultsOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.pod.memory.paging.faults",
		opt...,
	)
	if err != nil {
		return PodMemoryPagingFaults{noop.Int64Counter{}}, err
	}
	return PodMemoryPagingFaults{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryPagingFaults) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryPagingFaults) Name() string {
	return "k8s.pod.memory.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryPagingFaults) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryPagingFaults) Description() string {
	return "Pod memory paging faults."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Cumulative number of major/minor page faults.
// This metric is derived from the [MemoryStats.PageFaults] and
// [MemoryStats.MajorPageFaults] field of the [PodStats.Memory] of the Kubelet's
// stats API.
//
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryPagingFaults) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
//
// Cumulative number of major/minor page faults.
// This metric is derived from the [MemoryStats.PageFaults] and
// [MemoryStats.MajorPageFaults] field of the [PodStats.Memory] of the Kubelet's
// stats API.
//
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryPagingFaults) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (PodMemoryPagingFaults) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// PodMemoryRss is an instrument used to record metric values conforming to the
// "k8s.pod.memory.rss" semantic conventions. It represents the pod memory RSS.
type PodMemoryRss struct {
	metric.Int64UpDownCounter
}

var newPodMemoryRssOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod memory RSS."),
	metric.WithUnit("By"),
}

// NewPodMemoryRss returns a new PodMemoryRss instrument.
func NewPodMemoryRss(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryRss, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRss{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRssOpts
	} else {
		opt = append(opt, newPodMemoryRssOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.rss",
		opt...,
	)
	if err != nil {
		return PodMemoryRss{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryRss{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRss) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRss) Name() string {
	return "k8s.pod.memory.rss"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRss) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRss) Description() string {
	return "Pod memory RSS."
}

// Add adds incr to the existing count for attrs.
//
// The amount of anonymous and swap cache memory (includes transparent
// hugepages).
// This metric is derived from the [MemoryStats.RSSBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.RSSBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryRss) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// The amount of anonymous and swap cache memory (includes transparent
// hugepages).
// This metric is derived from the [MemoryStats.RSSBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.RSSBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryRss) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryUsage is an instrument used to record metric values conforming to the
// "k8s.pod.memory.usage" semantic conventions. It represents the memory usage of
// the Pod.
type PodMemoryUsage struct {
	metric.Int64Gauge
}

var newPodMemoryUsageOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Memory usage of the Pod."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newPodMemoryUsageOpts
	} else {
		opt = append(opt, newPodMemoryUsageOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.pod.memory.usage",
		opt...,
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
	return "Memory usage of the Pod."
}

// Record records val to the current distribution for attrs.
//
// Total memory usage of the Pod
func (m PodMemoryUsage) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Total memory usage of the Pod
func (m PodMemoryUsage) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// PodMemoryWorkingSet is an instrument used to record metric values conforming
// to the "k8s.pod.memory.working_set" semantic conventions. It represents the
// pod memory working set.
type PodMemoryWorkingSet struct {
	metric.Int64UpDownCounter
}

var newPodMemoryWorkingSetOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod memory working set."),
	metric.WithUnit("By"),
}

// NewPodMemoryWorkingSet returns a new PodMemoryWorkingSet instrument.
func NewPodMemoryWorkingSet(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryWorkingSet, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryWorkingSet{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryWorkingSetOpts
	} else {
		opt = append(opt, newPodMemoryWorkingSetOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.working_set",
		opt...,
	)
	if err != nil {
		return PodMemoryWorkingSet{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryWorkingSet{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryWorkingSet) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryWorkingSet) Name() string {
	return "k8s.pod.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryWorkingSet) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryWorkingSet) Description() string {
	return "Pod memory working set."
}

// Add adds incr to the existing count for attrs.
//
// The amount of working set memory. This includes recently accessed memory,
// dirty memory, and kernel memory. WorkingSetBytes is <= UsageBytes.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryWorkingSet) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// The amount of working set memory. This includes recently accessed memory,
// dirty memory, and kernel memory. WorkingSetBytes is <= UsageBytes.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryWorkingSet) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodNetworkErrors is an instrument used to record metric values conforming to
// the "k8s.pod.network.errors" semantic conventions. It represents the pod
// network errors.
type PodNetworkErrors struct {
	metric.Int64Counter
}

var newPodNetworkErrorsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Pod network errors."),
	metric.WithUnit("{error}"),
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

	if len(opt) == 0 {
		opt = newPodNetworkErrorsOpts
	} else {
		opt = append(opt, newPodNetworkErrorsOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.pod.network.errors",
		opt...,
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
	return "Pod network errors."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m PodNetworkErrors) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
func (m PodNetworkErrors) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
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

var newPodNetworkIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Network bytes for the Pod."),
	metric.WithUnit("By"),
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

	if len(opt) == 0 {
		opt = newPodNetworkIOOpts
	} else {
		opt = append(opt, newPodNetworkIOOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.pod.network.io",
		opt...,
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
	return "Network bytes for the Pod."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m PodNetworkIO) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
func (m PodNetworkIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
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

// PodStatusPhase is an instrument used to record metric values conforming to the
// "k8s.pod.status.phase" semantic conventions. It represents the describes
// number of K8s Pods that are currently in a given phase.
type PodStatusPhase struct {
	metric.Int64UpDownCounter
}

var newPodStatusPhaseOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes number of K8s Pods that are currently in a given phase."),
	metric.WithUnit("{pod}"),
}

// NewPodStatusPhase returns a new PodStatusPhase instrument.
func NewPodStatusPhase(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodStatusPhase, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodStatusPhase{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodStatusPhaseOpts
	} else {
		opt = append(opt, newPodStatusPhaseOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.status.phase",
		opt...,
	)
	if err != nil {
		return PodStatusPhase{noop.Int64UpDownCounter{}}, err
	}
	return PodStatusPhase{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodStatusPhase) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodStatusPhase) Name() string {
	return "k8s.pod.status.phase"
}

// Unit returns the semantic convention unit of the instrument
func (PodStatusPhase) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (PodStatusPhase) Description() string {
	return "Describes number of K8s Pods that are currently in a given phase."
}

// Add adds incr to the existing count for attrs.
//
// The podStatusPhase is the the phase for the pod. Corresponds to the `phase`
// field of the: [K8s PodStatus]
//
// All possible pod phases will be reported at each time interval to avoid
// missing metrics.
// Only the value corresponding to the current phase will be non-zero.
//
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#podstatus-v1-core
func (m PodStatusPhase) Add(
	ctx context.Context,
	incr int64,
	podStatusPhase PodStatusPhaseAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.pod.status.phase", string(podStatusPhase)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// All possible pod phases will be reported at each time interval to avoid
// missing metrics.
// Only the value corresponding to the current phase will be non-zero.
func (m PodStatusPhase) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodStatusReason is an instrument used to record metric values conforming to
// the "k8s.pod.status.reason" semantic conventions. It represents the describes
// the number of K8s Pods that are currently in a state for a given reason.
type PodStatusReason struct {
	metric.Int64UpDownCounter
}

var newPodStatusReasonOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the number of K8s Pods that are currently in a state for a given reason."),
	metric.WithUnit("{pod}"),
}

// NewPodStatusReason returns a new PodStatusReason instrument.
func NewPodStatusReason(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodStatusReason, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodStatusReason{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodStatusReasonOpts
	} else {
		opt = append(opt, newPodStatusReasonOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.status.reason",
		opt...,
	)
	if err != nil {
		return PodStatusReason{noop.Int64UpDownCounter{}}, err
	}
	return PodStatusReason{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodStatusReason) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodStatusReason) Name() string {
	return "k8s.pod.status.reason"
}

// Unit returns the semantic convention unit of the instrument
func (PodStatusReason) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (PodStatusReason) Description() string {
	return "Describes the number of K8s Pods that are currently in a state for a given reason."
}

// Add adds incr to the existing count for attrs.
//
// The podStatusReason is the the reason for the pod state. Corresponds to the
// `reason` field of the: [K8s PodStatus]
//
// All possible pod status reasons will be reported at each time interval to
// avoid missing metrics.
// Only the value corresponding to the current reason will be non-zero.
//
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#podstatus-v1-core
func (m PodStatusReason) Add(
	ctx context.Context,
	incr int64,
	podStatusReason PodStatusReasonAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.pod.status.reason", string(podStatusReason)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// All possible pod status reasons will be reported at each time interval to
// avoid missing metrics.
// Only the value corresponding to the current reason will be non-zero.
func (m PodStatusReason) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodUptime is an instrument used to record metric values conforming to the
// "k8s.pod.uptime" semantic conventions. It represents the time the Pod has been
// running.
type PodUptime struct {
	metric.Float64Gauge
}

var newPodUptimeOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The time the Pod has been running."),
	metric.WithUnit("s"),
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

	if len(opt) == 0 {
		opt = newPodUptimeOpts
	} else {
		opt = append(opt, newPodUptimeOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.pod.uptime",
		opt...,
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
	return "The time the Pod has been running."
}

// Record records val to the current distribution for attrs.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m PodUptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available.
// The actual accuracy would depend on the instrumentation and operating system.
func (m PodUptime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodVolumeAvailable is an instrument used to record metric values conforming to
// the "k8s.pod.volume.available" semantic conventions. It represents the pod
// volume storage space available.
type PodVolumeAvailable struct {
	metric.Int64UpDownCounter
}

var newPodVolumeAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod volume storage space available."),
	metric.WithUnit("By"),
}

// NewPodVolumeAvailable returns a new PodVolumeAvailable instrument.
func NewPodVolumeAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeAvailableOpts
	} else {
		opt = append(opt, newPodVolumeAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.available",
		opt...,
	)
	if err != nil {
		return PodVolumeAvailable{noop.Int64UpDownCounter{}}, err
	}
	return PodVolumeAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeAvailable) Name() string {
	return "k8s.pod.volume.available"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeAvailable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeAvailable) Description() string {
	return "Pod volume storage space available."
}

// Add adds incr to the existing count for attrs.
//
// The volumeName is the the name of the K8s volume.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is derived from the
// [VolumeStats.AvailableBytes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeAvailable) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.volume.name", volumeName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [VolumeStats.AvailableBytes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeAvailable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeCapacity is an instrument used to record metric values conforming to
// the "k8s.pod.volume.capacity" semantic conventions. It represents the pod
// volume total capacity.
type PodVolumeCapacity struct {
	metric.Int64UpDownCounter
}

var newPodVolumeCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod volume total capacity."),
	metric.WithUnit("By"),
}

// NewPodVolumeCapacity returns a new PodVolumeCapacity instrument.
func NewPodVolumeCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeCapacity{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeCapacityOpts
	} else {
		opt = append(opt, newPodVolumeCapacityOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.capacity",
		opt...,
	)
	if err != nil {
		return PodVolumeCapacity{noop.Int64UpDownCounter{}}, err
	}
	return PodVolumeCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeCapacity) Name() string {
	return "k8s.pod.volume.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeCapacity) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeCapacity) Description() string {
	return "Pod volume total capacity."
}

// Add adds incr to the existing count for attrs.
//
// The volumeName is the the name of the K8s volume.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is derived from the
// [VolumeStats.CapacityBytes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeCapacity) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.volume.name", volumeName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [VolumeStats.CapacityBytes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeCapacity) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeCount is an instrument used to record metric values conforming
// to the "k8s.pod.volume.inode.count" semantic conventions. It represents the
// total inodes in the filesystem of the Pod's volume.
type PodVolumeInodeCount struct {
	metric.Int64UpDownCounter
}

var newPodVolumeInodeCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The total inodes in the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeCount returns a new PodVolumeInodeCount instrument.
func NewPodVolumeInodeCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeInodeCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeInodeCountOpts
	} else {
		opt = append(opt, newPodVolumeInodeCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.inode.count",
		opt...,
	)
	if err != nil {
		return PodVolumeInodeCount{noop.Int64UpDownCounter{}}, err
	}
	return PodVolumeInodeCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeInodeCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeInodeCount) Name() string {
	return "k8s.pod.volume.inode.count"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeInodeCount) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeInodeCount) Description() string {
	return "The total inodes in the filesystem of the Pod's volume."
}

// Add adds incr to the existing count for attrs.
//
// The volumeName is the the name of the K8s volume.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is derived from the
// [VolumeStats.Inodes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.Inodes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeCount) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.volume.name", volumeName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [VolumeStats.Inodes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.Inodes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeCount) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeFree is an instrument used to record metric values conforming to
// the "k8s.pod.volume.inode.free" semantic conventions. It represents the free
// inodes in the filesystem of the Pod's volume.
type PodVolumeInodeFree struct {
	metric.Int64UpDownCounter
}

var newPodVolumeInodeFreeOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The free inodes in the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeFree returns a new PodVolumeInodeFree instrument.
func NewPodVolumeInodeFree(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeInodeFree, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeFree{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeInodeFreeOpts
	} else {
		opt = append(opt, newPodVolumeInodeFreeOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.inode.free",
		opt...,
	)
	if err != nil {
		return PodVolumeInodeFree{noop.Int64UpDownCounter{}}, err
	}
	return PodVolumeInodeFree{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeInodeFree) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeInodeFree) Name() string {
	return "k8s.pod.volume.inode.free"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeInodeFree) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeInodeFree) Description() string {
	return "The free inodes in the filesystem of the Pod's volume."
}

// Add adds incr to the existing count for attrs.
//
// The volumeName is the the name of the K8s volume.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is derived from the
// [VolumeStats.InodesFree] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.InodesFree]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeFree) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.volume.name", volumeName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [VolumeStats.InodesFree] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.InodesFree]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeFree) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeFree) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeUsed is an instrument used to record metric values conforming to
// the "k8s.pod.volume.inode.used" semantic conventions. It represents the inodes
// used by the filesystem of the Pod's volume.
type PodVolumeInodeUsed struct {
	metric.Int64UpDownCounter
}

var newPodVolumeInodeUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The inodes used by the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeUsed returns a new PodVolumeInodeUsed instrument.
func NewPodVolumeInodeUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeInodeUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeInodeUsedOpts
	} else {
		opt = append(opt, newPodVolumeInodeUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.inode.used",
		opt...,
	)
	if err != nil {
		return PodVolumeInodeUsed{noop.Int64UpDownCounter{}}, err
	}
	return PodVolumeInodeUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeInodeUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeInodeUsed) Name() string {
	return "k8s.pod.volume.inode.used"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeInodeUsed) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeInodeUsed) Description() string {
	return "The inodes used by the filesystem of the Pod's volume."
}

// Add adds incr to the existing count for attrs.
//
// The volumeName is the the name of the K8s volume.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is derived from the
// [VolumeStats.InodesUsed] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// This may not be equal to `inodes - free` because filesystem may share inodes
// with other filesystems.
//
// [VolumeStats.InodesUsed]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeUsed) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.volume.name", volumeName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [VolumeStats.InodesUsed] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// This may not be equal to `inodes - free` because filesystem may share inodes
// with other filesystems.
//
// [VolumeStats.InodesUsed]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeUsed) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeUsage is an instrument used to record metric values conforming to the
// "k8s.pod.volume.usage" semantic conventions. It represents the pod volume
// usage.
type PodVolumeUsage struct {
	metric.Int64UpDownCounter
}

var newPodVolumeUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod volume usage."),
	metric.WithUnit("By"),
}

// NewPodVolumeUsage returns a new PodVolumeUsage instrument.
func NewPodVolumeUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeUsageOpts
	} else {
		opt = append(opt, newPodVolumeUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.usage",
		opt...,
	)
	if err != nil {
		return PodVolumeUsage{noop.Int64UpDownCounter{}}, err
	}
	return PodVolumeUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeUsage) Name() string {
	return "k8s.pod.volume.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeUsage) Description() string {
	return "Pod volume usage."
}

// Add adds incr to the existing count for attrs.
//
// The volumeName is the the name of the K8s volume.
//
// All additional attrs passed are included in the recorded value.
//
// This may not equal capacity - available.
//
// This metric is derived from the
// [VolumeStats.UsedBytes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeUsage) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.volume.name", volumeName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This may not equal capacity - available.
//
// This metric is derived from the
// [VolumeStats.UsedBytes] field
// of the [PodStats] of the
// Kubelet's stats API.
//
// [VolumeStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeUsage) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// ReplicaSetPodAvailable is an instrument used to record metric values
// conforming to the "k8s.replicaset.pod.available" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this replicaset.
type ReplicaSetPodAvailable struct {
	metric.Int64UpDownCounter
}

var newReplicaSetPodAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset."),
	metric.WithUnit("{pod}"),
}

// NewReplicaSetPodAvailable returns a new ReplicaSetPodAvailable instrument.
func NewReplicaSetPodAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicaSetPodAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicaSetPodAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicaSetPodAvailableOpts
	} else {
		opt = append(opt, newReplicaSetPodAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicaset.pod.available",
		opt...,
	)
	if err != nil {
		return ReplicaSetPodAvailable{noop.Int64UpDownCounter{}}, err
	}
	return ReplicaSetPodAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicaSetPodAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicaSetPodAvailable) Name() string {
	return "k8s.replicaset.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicaSetPodAvailable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicaSetPodAvailable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicaSetStatus].
//
// [K8s ReplicaSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetstatus-v1-apps
func (m ReplicaSetPodAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicaSetStatus].
//
// [K8s ReplicaSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetstatus-v1-apps
func (m ReplicaSetPodAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicaSetPodDesired is an instrument used to record metric values conforming
// to the "k8s.replicaset.pod.desired" semantic conventions. It represents the
// number of desired replica pods in this replicaset.
type ReplicaSetPodDesired struct {
	metric.Int64UpDownCounter
}

var newReplicaSetPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this replicaset."),
	metric.WithUnit("{pod}"),
}

// NewReplicaSetPodDesired returns a new ReplicaSetPodDesired instrument.
func NewReplicaSetPodDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicaSetPodDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicaSetPodDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicaSetPodDesiredOpts
	} else {
		opt = append(opt, newReplicaSetPodDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicaset.pod.desired",
		opt...,
	)
	if err != nil {
		return ReplicaSetPodDesired{noop.Int64UpDownCounter{}}, err
	}
	return ReplicaSetPodDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicaSetPodDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicaSetPodDesired) Name() string {
	return "k8s.replicaset.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicaSetPodDesired) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicaSetPodDesired) Description() string {
	return "Number of desired replica pods in this replicaset."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicaSetSpec].
//
// [K8s ReplicaSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetspec-v1-apps
func (m ReplicaSetPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicaSetSpec].
//
// [K8s ReplicaSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetspec-v1-apps
func (m ReplicaSetPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicationControllerPodAvailable is an instrument used to record metric
// values conforming to the "k8s.replicationcontroller.pod.available" semantic
// conventions. It represents the total number of available replica pods (ready
// for at least minReadySeconds) targeted by this replication controller.
type ReplicationControllerPodAvailable struct {
	metric.Int64UpDownCounter
}

var newReplicationControllerPodAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller."),
	metric.WithUnit("{pod}"),
}

// NewReplicationControllerPodAvailable returns a new
// ReplicationControllerPodAvailable instrument.
func NewReplicationControllerPodAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicationControllerPodAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicationControllerPodAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicationControllerPodAvailableOpts
	} else {
		opt = append(opt, newReplicationControllerPodAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicationcontroller.pod.available",
		opt...,
	)
	if err != nil {
		return ReplicationControllerPodAvailable{noop.Int64UpDownCounter{}}, err
	}
	return ReplicationControllerPodAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicationControllerPodAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerPodAvailable) Name() string {
	return "k8s.replicationcontroller.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerPodAvailable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerPodAvailable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicationControllerStatus]
//
// [K8s ReplicationControllerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerstatus-v1-core
func (m ReplicationControllerPodAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicationControllerStatus]
//
// [K8s ReplicationControllerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerstatus-v1-core
func (m ReplicationControllerPodAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicationControllerPodDesired is an instrument used to record metric values
// conforming to the "k8s.replicationcontroller.pod.desired" semantic
// conventions. It represents the number of desired replica pods in this
// replication controller.
type ReplicationControllerPodDesired struct {
	metric.Int64UpDownCounter
}

var newReplicationControllerPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this replication controller."),
	metric.WithUnit("{pod}"),
}

// NewReplicationControllerPodDesired returns a new
// ReplicationControllerPodDesired instrument.
func NewReplicationControllerPodDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ReplicationControllerPodDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicationControllerPodDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicationControllerPodDesiredOpts
	} else {
		opt = append(opt, newReplicationControllerPodDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.replicationcontroller.pod.desired",
		opt...,
	)
	if err != nil {
		return ReplicationControllerPodDesired{noop.Int64UpDownCounter{}}, err
	}
	return ReplicationControllerPodDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicationControllerPodDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerPodDesired) Name() string {
	return "k8s.replicationcontroller.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerPodDesired) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerPodDesired) Description() string {
	return "Number of desired replica pods in this replication controller."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicationControllerSpec]
//
// [K8s ReplicationControllerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerspec-v1-core
func (m ReplicationControllerPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicationControllerSpec]
//
// [K8s ReplicationControllerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerspec-v1-core
func (m ReplicationControllerPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPULimitHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.limit.hard" semantic conventions. It
// represents the CPU limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaCPULimitHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaCPULimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPULimitHard returns a new ResourceQuotaCPULimitHard
// instrument.
func NewResourceQuotaCPULimitHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaCPULimitHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPULimitHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPULimitHardOpts
	} else {
		opt = append(opt, newResourceQuotaCPULimitHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.limit.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPULimitHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaCPULimitHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPULimitHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPULimitHard) Name() string {
	return "k8s.resourcequota.cpu.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPULimitHard) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPULimitHard) Description() string {
	return "The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPULimitUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.limit.used" semantic conventions. It
// represents the CPU limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaCPULimitUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaCPULimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPULimitUsed returns a new ResourceQuotaCPULimitUsed
// instrument.
func NewResourceQuotaCPULimitUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaCPULimitUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPULimitUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPULimitUsedOpts
	} else {
		opt = append(opt, newResourceQuotaCPULimitUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.limit.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPULimitUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaCPULimitUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPULimitUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPULimitUsed) Name() string {
	return "k8s.resourcequota.cpu.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPULimitUsed) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPULimitUsed) Description() string {
	return "The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPURequestHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.request.hard" semantic conventions.
// It represents the CPU requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaCPURequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaCPURequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPURequestHard returns a new ResourceQuotaCPURequestHard
// instrument.
func NewResourceQuotaCPURequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaCPURequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPURequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPURequestHardOpts
	} else {
		opt = append(opt, newResourceQuotaCPURequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPURequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaCPURequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPURequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPURequestHard) Name() string {
	return "k8s.resourcequota.cpu.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPURequestHard) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPURequestHard) Description() string {
	return "The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPURequestUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.request.used" semantic conventions.
// It represents the CPU requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaCPURequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaCPURequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPURequestUsed returns a new ResourceQuotaCPURequestUsed
// instrument.
func NewResourceQuotaCPURequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaCPURequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPURequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPURequestUsedOpts
	} else {
		opt = append(opt, newResourceQuotaCPURequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPURequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaCPURequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPURequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPURequestUsed) Name() string {
	return "k8s.resourcequota.cpu.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPURequestUsed) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPURequestUsed) Description() string {
	return "The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageLimitHard is an instrument used to record metric
// values conforming to the "k8s.resourcequota.ephemeral_storage.limit.hard"
// semantic conventions. It represents the sum of local ephemeral storage limits
// in the namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaEphemeralStorageLimitHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageLimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageLimitHard returns a new
// ResourceQuotaEphemeralStorageLimitHard instrument.
func NewResourceQuotaEphemeralStorageLimitHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaEphemeralStorageLimitHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageLimitHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageLimitHardOpts
	} else {
		opt = append(opt, newResourceQuotaEphemeralStorageLimitHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.limit.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageLimitHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageLimitHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageLimitHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageLimitHard) Name() string {
	return "k8s.resourcequota.ephemeral_storage.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageLimitHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageLimitHard) Description() string {
	return "The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageLimitUsed is an instrument used to record metric
// values conforming to the "k8s.resourcequota.ephemeral_storage.limit.used"
// semantic conventions. It represents the sum of local ephemeral storage limits
// in the namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaEphemeralStorageLimitUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageLimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageLimitUsed returns a new
// ResourceQuotaEphemeralStorageLimitUsed instrument.
func NewResourceQuotaEphemeralStorageLimitUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaEphemeralStorageLimitUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageLimitUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageLimitUsedOpts
	} else {
		opt = append(opt, newResourceQuotaEphemeralStorageLimitUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.limit.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageLimitUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageLimitUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageLimitUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageLimitUsed) Name() string {
	return "k8s.resourcequota.ephemeral_storage.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageLimitUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageLimitUsed) Description() string {
	return "The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageRequestHard is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.request.hard" semantic conventions. It
// represents the sum of local ephemeral storage requests in the namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaEphemeralStorageRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageRequestHard returns a new
// ResourceQuotaEphemeralStorageRequestHard instrument.
func NewResourceQuotaEphemeralStorageRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaEphemeralStorageRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageRequestHardOpts
	} else {
		opt = append(opt, newResourceQuotaEphemeralStorageRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageRequestHard) Name() string {
	return "k8s.resourcequota.ephemeral_storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageRequestHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageRequestHard) Description() string {
	return "The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageRequestUsed is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.request.used" semantic conventions. It
// represents the sum of local ephemeral storage requests in the namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaEphemeralStorageRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageRequestUsed returns a new
// ResourceQuotaEphemeralStorageRequestUsed instrument.
func NewResourceQuotaEphemeralStorageRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaEphemeralStorageRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageRequestUsedOpts
	} else {
		opt = append(opt, newResourceQuotaEphemeralStorageRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageRequestUsed) Name() string {
	return "k8s.resourcequota.ephemeral_storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageRequestUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageRequestUsed) Description() string {
	return "The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaHugepageCountRequestHard is an instrument used to record metric
// values conforming to the "k8s.resourcequota.hugepage_count.request.hard"
// semantic conventions. It represents the huge page requests in a specific
// namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaHugepageCountRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaHugepageCountRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{hugepage}"),
}

// NewResourceQuotaHugepageCountRequestHard returns a new
// ResourceQuotaHugepageCountRequestHard instrument.
func NewResourceQuotaHugepageCountRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaHugepageCountRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaHugepageCountRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaHugepageCountRequestHardOpts
	} else {
		opt = append(opt, newResourceQuotaHugepageCountRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.hugepage_count.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaHugepageCountRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaHugepageCountRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaHugepageCountRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaHugepageCountRequestHard) Name() string {
	return "k8s.resourcequota.hugepage_count.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaHugepageCountRequestHard) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaHugepageCountRequestHard) Description() string {
	return "The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// The hugepageSize is the the size (identifier) of the K8s huge page.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestHard) Add(
	ctx context.Context,
	incr int64,
	hugepageSize string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.hugepage.size", hugepageSize),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaHugepageCountRequestUsed is an instrument used to record metric
// values conforming to the "k8s.resourcequota.hugepage_count.request.used"
// semantic conventions. It represents the huge page requests in a specific
// namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaHugepageCountRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaHugepageCountRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{hugepage}"),
}

// NewResourceQuotaHugepageCountRequestUsed returns a new
// ResourceQuotaHugepageCountRequestUsed instrument.
func NewResourceQuotaHugepageCountRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaHugepageCountRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaHugepageCountRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaHugepageCountRequestUsedOpts
	} else {
		opt = append(opt, newResourceQuotaHugepageCountRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.hugepage_count.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaHugepageCountRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaHugepageCountRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaHugepageCountRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaHugepageCountRequestUsed) Name() string {
	return "k8s.resourcequota.hugepage_count.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaHugepageCountRequestUsed) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaHugepageCountRequestUsed) Description() string {
	return "The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// The hugepageSize is the the size (identifier) of the K8s huge page.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestUsed) Add(
	ctx context.Context,
	incr int64,
	hugepageSize string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.hugepage.size", hugepageSize),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryLimitHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.limit.hard" semantic conventions.
// It represents the memory limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaMemoryLimitHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryLimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryLimitHard returns a new ResourceQuotaMemoryLimitHard
// instrument.
func NewResourceQuotaMemoryLimitHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaMemoryLimitHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryLimitHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryLimitHardOpts
	} else {
		opt = append(opt, newResourceQuotaMemoryLimitHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.limit.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryLimitHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaMemoryLimitHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryLimitHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryLimitHard) Name() string {
	return "k8s.resourcequota.memory.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryLimitHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryLimitHard) Description() string {
	return "The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryLimitUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.limit.used" semantic conventions.
// It represents the memory limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaMemoryLimitUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryLimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryLimitUsed returns a new ResourceQuotaMemoryLimitUsed
// instrument.
func NewResourceQuotaMemoryLimitUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaMemoryLimitUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryLimitUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryLimitUsedOpts
	} else {
		opt = append(opt, newResourceQuotaMemoryLimitUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.limit.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryLimitUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaMemoryLimitUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryLimitUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryLimitUsed) Name() string {
	return "k8s.resourcequota.memory.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryLimitUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryLimitUsed) Description() string {
	return "The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryRequestHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.request.hard" semantic
// conventions. It represents the memory requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaMemoryRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryRequestHard returns a new ResourceQuotaMemoryRequestHard
// instrument.
func NewResourceQuotaMemoryRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaMemoryRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryRequestHardOpts
	} else {
		opt = append(opt, newResourceQuotaMemoryRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaMemoryRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryRequestHard) Name() string {
	return "k8s.resourcequota.memory.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryRequestHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryRequestHard) Description() string {
	return "The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryRequestUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.request.used" semantic
// conventions. It represents the memory requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaMemoryRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryRequestUsed returns a new ResourceQuotaMemoryRequestUsed
// instrument.
func NewResourceQuotaMemoryRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaMemoryRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryRequestUsedOpts
	} else {
		opt = append(opt, newResourceQuotaMemoryRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaMemoryRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryRequestUsed) Name() string {
	return "k8s.resourcequota.memory.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryRequestUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryRequestUsed) Description() string {
	return "The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaObjectCountHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.object_count.hard" semantic conventions.
// It represents the object count limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaObjectCountHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaObjectCountHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{object}"),
}

// NewResourceQuotaObjectCountHard returns a new ResourceQuotaObjectCountHard
// instrument.
func NewResourceQuotaObjectCountHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaObjectCountHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaObjectCountHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaObjectCountHardOpts
	} else {
		opt = append(opt, newResourceQuotaObjectCountHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.object_count.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaObjectCountHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaObjectCountHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaObjectCountHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaObjectCountHard) Name() string {
	return "k8s.resourcequota.object_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaObjectCountHard) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaObjectCountHard) Description() string {
	return "The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// The resourcequotaResourceName is the the name of the K8s resource a resource
// quota defines.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountHard) Add(
	ctx context.Context,
	incr int64,
	resourcequotaResourceName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.resourcequota.resource_name", resourcequotaResourceName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaObjectCountUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.object_count.used" semantic conventions.
// It represents the object count limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaObjectCountUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaObjectCountUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{object}"),
}

// NewResourceQuotaObjectCountUsed returns a new ResourceQuotaObjectCountUsed
// instrument.
func NewResourceQuotaObjectCountUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaObjectCountUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaObjectCountUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaObjectCountUsedOpts
	} else {
		opt = append(opt, newResourceQuotaObjectCountUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.object_count.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaObjectCountUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaObjectCountUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaObjectCountUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaObjectCountUsed) Name() string {
	return "k8s.resourcequota.object_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaObjectCountUsed) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaObjectCountUsed) Description() string {
	return "The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// The resourcequotaResourceName is the the name of the K8s resource a resource
// quota defines.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountUsed) Add(
	ctx context.Context,
	incr int64,
	resourcequotaResourceName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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
				attribute.String("k8s.resourcequota.resource_name", resourcequotaResourceName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaPersistentvolumeclaimCountHard is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.persistentvolumeclaim_count.hard" semantic conventions. It
// represents the total number of PersistentVolumeClaims that can exist in the
// namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaPersistentvolumeclaimCountHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaPersistentvolumeclaimCountHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewResourceQuotaPersistentvolumeclaimCountHard returns a new
// ResourceQuotaPersistentvolumeclaimCountHard instrument.
func NewResourceQuotaPersistentvolumeclaimCountHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaPersistentvolumeclaimCountHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaPersistentvolumeclaimCountHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaPersistentvolumeclaimCountHardOpts
	} else {
		opt = append(opt, newResourceQuotaPersistentvolumeclaimCountHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.persistentvolumeclaim_count.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaPersistentvolumeclaimCountHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaPersistentvolumeclaimCountHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaPersistentvolumeclaimCountHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaPersistentvolumeclaimCountHard) Name() string {
	return "k8s.resourcequota.persistentvolumeclaim_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaPersistentvolumeclaimCountHard) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaPersistentvolumeclaimCountHard) Description() string {
	return "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountHard) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ResourceQuotaPersistentvolumeclaimCountHard) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaPersistentvolumeclaimCountUsed is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.persistentvolumeclaim_count.used" semantic conventions. It
// represents the total number of PersistentVolumeClaims that can exist in the
// namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaPersistentvolumeclaimCountUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaPersistentvolumeclaimCountUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewResourceQuotaPersistentvolumeclaimCountUsed returns a new
// ResourceQuotaPersistentvolumeclaimCountUsed instrument.
func NewResourceQuotaPersistentvolumeclaimCountUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaPersistentvolumeclaimCountUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaPersistentvolumeclaimCountUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaPersistentvolumeclaimCountUsedOpts
	} else {
		opt = append(opt, newResourceQuotaPersistentvolumeclaimCountUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.persistentvolumeclaim_count.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaPersistentvolumeclaimCountUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaPersistentvolumeclaimCountUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaPersistentvolumeclaimCountUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaPersistentvolumeclaimCountUsed) Name() string {
	return "k8s.resourcequota.persistentvolumeclaim_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaPersistentvolumeclaimCountUsed) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaPersistentvolumeclaimCountUsed) Description() string {
	return "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountUsed) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ResourceQuotaPersistentvolumeclaimCountUsed) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaStorageRequestHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.storage.request.hard" semantic
// conventions. It represents the storage requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
type ResourceQuotaStorageRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaStorageRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaStorageRequestHard returns a new
// ResourceQuotaStorageRequestHard instrument.
func NewResourceQuotaStorageRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaStorageRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaStorageRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaStorageRequestHardOpts
	} else {
		opt = append(opt, newResourceQuotaStorageRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.storage.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaStorageRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaStorageRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaStorageRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaStorageRequestHard) Name() string {
	return "k8s.resourcequota.storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaStorageRequestHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaStorageRequestHard) Description() string {
	return "The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestHard) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ResourceQuotaStorageRequestHard) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaStorageRequestUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.storage.request.used" semantic
// conventions. It represents the storage requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
type ResourceQuotaStorageRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaStorageRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaStorageRequestUsed returns a new
// ResourceQuotaStorageRequestUsed instrument.
func NewResourceQuotaStorageRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ResourceQuotaStorageRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaStorageRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaStorageRequestUsedOpts
	} else {
		opt = append(opt, newResourceQuotaStorageRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.storage.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaStorageRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ResourceQuotaStorageRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaStorageRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaStorageRequestUsed) Name() string {
	return "k8s.resourcequota.storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaStorageRequestUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaStorageRequestUsed) Description() string {
	return "The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestUsed) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

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

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ResourceQuotaStorageRequestUsed) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// StatefulSetPodCurrent is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.current" semantic conventions. It represents the
// number of replica pods created by the statefulset controller from the
// statefulset version indicated by currentRevision.
type StatefulSetPodCurrent struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodCurrent returns a new StatefulSetPodCurrent instrument.
func NewStatefulSetPodCurrent(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetPodCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodCurrent{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodCurrentOpts
	} else {
		opt = append(opt, newStatefulSetPodCurrentOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.pod.current",
		opt...,
	)
	if err != nil {
		return StatefulSetPodCurrent{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetPodCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodCurrent) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodCurrent) Name() string {
	return "k8s.statefulset.pod.current"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodCurrent) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodCurrent) Description() string {
	return "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetPodCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetPodCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodDesired is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.desired" semantic conventions. It represents the
// number of desired replica pods in this statefulset.
type StatefulSetPodDesired struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this statefulset."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodDesired returns a new StatefulSetPodDesired instrument.
func NewStatefulSetPodDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetPodDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodDesiredOpts
	} else {
		opt = append(opt, newStatefulSetPodDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.pod.desired",
		opt...,
	)
	if err != nil {
		return StatefulSetPodDesired{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetPodDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodDesired) Name() string {
	return "k8s.statefulset.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodDesired) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodDesired) Description() string {
	return "Number of desired replica pods in this statefulset."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s StatefulSetSpec].
//
// [K8s StatefulSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetspec-v1-apps
func (m StatefulSetPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s StatefulSetSpec].
//
// [K8s StatefulSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetspec-v1-apps
func (m StatefulSetPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodReady is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.ready" semantic conventions. It represents the
// number of replica pods created for this statefulset with a Ready Condition.
type StatefulSetPodReady struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodReadyOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of replica pods created for this statefulset with a Ready Condition."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodReady returns a new StatefulSetPodReady instrument.
func NewStatefulSetPodReady(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetPodReady, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodReady{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodReadyOpts
	} else {
		opt = append(opt, newStatefulSetPodReadyOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.pod.ready",
		opt...,
	)
	if err != nil {
		return StatefulSetPodReady{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetPodReady{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodReady) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodReady) Name() string {
	return "k8s.statefulset.pod.ready"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodReady) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodReady) Description() string {
	return "The number of replica pods created for this statefulset with a Ready Condition."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `readyReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetPodReady) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `readyReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetPodReady) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodUpdated is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.updated" semantic conventions. It represents the
// number of replica pods created by the statefulset controller from the
// statefulset version indicated by updateRevision.
type StatefulSetPodUpdated struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodUpdatedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodUpdated returns a new StatefulSetPodUpdated instrument.
func NewStatefulSetPodUpdated(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetPodUpdated, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodUpdated{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodUpdatedOpts
	} else {
		opt = append(opt, newStatefulSetPodUpdatedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.pod.updated",
		opt...,
	)
	if err != nil {
		return StatefulSetPodUpdated{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetPodUpdated{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodUpdated) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodUpdated) Name() string {
	return "k8s.statefulset.pod.updated"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodUpdated) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodUpdated) Description() string {
	return "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `updatedReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetPodUpdated) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `updatedReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetPodUpdated) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}
