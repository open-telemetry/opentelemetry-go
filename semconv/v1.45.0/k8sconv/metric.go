// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package k8sconv provides types and functionality for OpenTelemetry semantic
// conventions in the "k8s" namespace.
package k8sconv

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/internal/metricpool"
)

// ContainerEphemeralStorageFsTypeAttr is an attribute conforming to the
// k8s.container.ephemeral_storage.fs_type semantic conventions. It represents
// the type of file system component for ephemeral storage.
type ContainerEphemeralStorageFsTypeAttr string

var (
	// ContainerEphemeralStorageFsTypeRootfs is the for the container's writable
	// layer usage.
	ContainerEphemeralStorageFsTypeRootfs ContainerEphemeralStorageFsTypeAttr = "rootfs"
	// ContainerEphemeralStorageFsTypeLogs is the for the container's log files
	// usage (stdout/stderr).
	ContainerEphemeralStorageFsTypeLogs ContainerEphemeralStorageFsTypeAttr = "logs"
)

// ContainerStatusReasonAttr is an attribute conforming to the
// k8s.container.status.reason semantic conventions. It represents the reason for
// the container state. Corresponds to the `reason` field of the:
// [K8s ContainerStateWaiting] or [K8s ContainerStateTerminated].
//
// [K8s ContainerStateWaiting]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatewaiting-v1-core
// [K8s ContainerStateTerminated]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstateterminated-v1-core
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
// [K8s ContainerState]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstate-v1-core
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

// PersistentvolumeStatusPhaseAttr is an attribute conforming to the
// k8s.persistentvolume.status.phase semantic conventions. It represents the
// phase of the PersistentVolume.
type PersistentvolumeStatusPhaseAttr string

var (
	// PersistentvolumeStatusPhaseAvailable is the volume is available and not yet
	// bound to a claim.
	PersistentvolumeStatusPhaseAvailable PersistentvolumeStatusPhaseAttr = "Available"
	// PersistentvolumeStatusPhaseBound is the volume is bound to a claim.
	PersistentvolumeStatusPhaseBound PersistentvolumeStatusPhaseAttr = "Bound"
	// PersistentvolumeStatusPhaseFailed is the volume has failed its automatic
	// reclamation.
	PersistentvolumeStatusPhaseFailed PersistentvolumeStatusPhaseAttr = "Failed"
	// PersistentvolumeStatusPhasePending is the volume is being provisioned.
	PersistentvolumeStatusPhasePending PersistentvolumeStatusPhaseAttr = "Pending"
	// PersistentvolumeStatusPhaseReleased is the claim has been deleted but the
	// volume is not yet available.
	PersistentvolumeStatusPhaseReleased PersistentvolumeStatusPhaseAttr = "Released"
)

// PersistentvolumeclaimStatusPhaseAttr is an attribute conforming to the
// k8s.persistentvolumeclaim.status.phase semantic conventions. It represents the
// phase of the PersistentVolumeClaim.
type PersistentvolumeclaimStatusPhaseAttr string

var (
	// PersistentvolumeclaimStatusPhaseBound is the claim is bound to a volume.
	PersistentvolumeclaimStatusPhaseBound PersistentvolumeclaimStatusPhaseAttr = "Bound"
	// PersistentvolumeclaimStatusPhaseLost is the claim has lost its underlying
	// volume (the volume does not exist anymore).
	PersistentvolumeclaimStatusPhaseLost PersistentvolumeclaimStatusPhaseAttr = "Lost"
	// PersistentvolumeclaimStatusPhasePending is the claim has not yet been bound
	// to a volume.
	PersistentvolumeclaimStatusPhasePending PersistentvolumeclaimStatusPhaseAttr = "Pending"
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

// ServiceEndpointAddressTypeAttr is an attribute conforming to the
// k8s.service.endpoint.address_type semantic conventions. It represents the
// address type of the service endpoint.
type ServiceEndpointAddressTypeAttr string

var (
	// ServiceEndpointAddressTypeIPv4 is the IPv4 address type.
	ServiceEndpointAddressTypeIPv4 ServiceEndpointAddressTypeAttr = "IPv4"
	// ServiceEndpointAddressTypeIPv6 is the IPv6 address type.
	ServiceEndpointAddressTypeIPv6 ServiceEndpointAddressTypeAttr = "IPv6"
	// ServiceEndpointAddressTypeFqdn is the FQDN address type.
	ServiceEndpointAddressTypeFqdn ServiceEndpointAddressTypeAttr = "FQDN"
)

// ServiceEndpointConditionAttr is an attribute conforming to the
// k8s.service.endpoint.condition semantic conventions. It represents the
// condition of the service endpoint.
type ServiceEndpointConditionAttr string

var (
	// ServiceEndpointConditionReady is the endpoint is ready to receive new
	// connections.
	ServiceEndpointConditionReady ServiceEndpointConditionAttr = "ready"
	// ServiceEndpointConditionServing is the endpoint is currently handling
	// traffic.
	ServiceEndpointConditionServing ServiceEndpointConditionAttr = "serving"
	// ServiceEndpointConditionTerminating is the endpoint is in the process of
	// shutting down.
	ServiceEndpointConditionTerminating ServiceEndpointConditionAttr = "terminating"
)

// VolumeTypeAttr is an attribute conforming to the k8s.volume.type semantic
// conventions. It represents the type of the K8s volume.
type VolumeTypeAttr string

var (
	// VolumeTypePersistentVolumeClaim is a [persistentVolumeClaim] volume.
	//
	// [persistentVolumeClaim]: https://kubernetes.io/docs/concepts/storage/volumes/#persistentvolumeclaim
	VolumeTypePersistentVolumeClaim VolumeTypeAttr = "persistentVolumeClaim"
	// VolumeTypeConfigMap is a [configMap] volume.
	//
	// [configMap]: https://kubernetes.io/docs/concepts/storage/volumes/#configmap
	VolumeTypeConfigMap VolumeTypeAttr = "configMap"
	// VolumeTypeDownwardAPI is a [downwardAPI] volume.
	//
	// [downwardAPI]: https://kubernetes.io/docs/concepts/storage/volumes/#downwardapi
	VolumeTypeDownwardAPI VolumeTypeAttr = "downwardAPI"
	// VolumeTypeEmptyDir is an [emptyDir] volume.
	//
	// [emptyDir]: https://kubernetes.io/docs/concepts/storage/volumes/#emptydir
	VolumeTypeEmptyDir VolumeTypeAttr = "emptyDir"
	// VolumeTypeSecret is a [secret] volume.
	//
	// [secret]: https://kubernetes.io/docs/concepts/storage/volumes/#secret
	VolumeTypeSecret VolumeTypeAttr = "secret"
	// VolumeTypeLocal is a [local] volume.
	//
	// [local]: https://kubernetes.io/docs/concepts/storage/volumes/#local
	VolumeTypeLocal VolumeTypeAttr = "local"
)

// NetworkIODirectionAttr is an attribute conforming to the network.io.direction
// semantic conventions. It represents the direction of traffic from the
// perspective of the observing host's physical or virtual network interface. It
// should not be used to represent the logical direction of a stateful connection
// or network flow.
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

// ContainerCPULimitCurrent is an instrument used to record metric values
// conforming to the "k8s.container.cpu.limit.current" semantic conventions. It
// represents the maximum CPU resource limit currently configured for a running
// container.
//
// Stability: Development
type ContainerCPULimitCurrent struct {
	metric.Float64UpDownCounter
}

var newContainerCPULimitCurrentOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit currently configured for a running container."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPULimitCurrent returns a new ContainerCPULimitCurrent instrument.
//
// Stability: Development
func NewContainerCPULimitCurrent(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ContainerCPULimitCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimitCurrent{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPULimitCurrentOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.container.cpu.limit.current",
		opt...,
	)
	if err != nil {
		return ContainerCPULimitCurrent{noop.Float64UpDownCounter{}}, err
	}
	return ContainerCPULimitCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimitCurrent) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimitCurrent) Name() string {
	return "k8s.container.cpu.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimitCurrent) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimitCurrent) Description() string {
	return "Maximum CPU resource limit currently configured for a running container."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPULimitCurrent) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPULimitCurrent) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerCPULimitCurrentObservable is an instrument used to record metric
// values conforming to the "k8s.container.cpu.limit.current" semantic
// conventions. It represents the maximum CPU resource limit currently configured
// for a running container.
//
// Stability: Development
type ContainerCPULimitCurrentObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newContainerCPULimitCurrentObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit currently configured for a running container."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPULimitCurrentObservable returns a new
// ContainerCPULimitCurrentObservable instrument.
//
// Stability: Development
func NewContainerCPULimitCurrentObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ContainerCPULimitCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimitCurrentObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPULimitCurrentObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.container.cpu.limit.current",
		opt...,
	)
	if err != nil {
		return ContainerCPULimitCurrentObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ContainerCPULimitCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimitCurrentObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimitCurrentObservable) Name() string {
	return "k8s.container.cpu.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimitCurrentObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimitCurrentObservable) Description() string {
	return "Maximum CPU resource limit currently configured for a running container."
}

// ContainerCPULimitDesired is an instrument used to record metric values
// conforming to the "k8s.container.cpu.limit.desired" semantic conventions. It
// represents the maximum CPU resource limit as defined by the container spec.
//
// Stability: Development
type ContainerCPULimitDesired struct {
	metric.Float64UpDownCounter
}

var newContainerCPULimitDesiredOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit as defined by the container spec."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPULimitDesired returns a new ContainerCPULimitDesired instrument.
//
// Stability: Development
func NewContainerCPULimitDesired(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ContainerCPULimitDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimitDesired{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPULimitDesiredOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.container.cpu.limit.desired",
		opt...,
	)
	if err != nil {
		return ContainerCPULimitDesired{noop.Float64UpDownCounter{}}, err
	}
	return ContainerCPULimitDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimitDesired) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimitDesired) Name() string {
	return "k8s.container.cpu.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimitDesired) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimitDesired) Description() string {
	return "Maximum CPU resource limit as defined by the container spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPULimitDesired) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPULimitDesired) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerCPULimitDesiredObservable is an instrument used to record metric
// values conforming to the "k8s.container.cpu.limit.desired" semantic
// conventions. It represents the maximum CPU resource limit as defined by the
// container spec.
//
// Stability: Development
type ContainerCPULimitDesiredObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newContainerCPULimitDesiredObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit as defined by the container spec."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPULimitDesiredObservable returns a new
// ContainerCPULimitDesiredObservable instrument.
//
// Stability: Development
func NewContainerCPULimitDesiredObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ContainerCPULimitDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimitDesiredObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPULimitDesiredObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.container.cpu.limit.desired",
		opt...,
	)
	if err != nil {
		return ContainerCPULimitDesiredObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ContainerCPULimitDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimitDesiredObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimitDesiredObservable) Name() string {
	return "k8s.container.cpu.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimitDesiredObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimitDesiredObservable) Description() string {
	return "Maximum CPU resource limit as defined by the container spec."
}

// ContainerCPULimitUtilization is an instrument used to record metric values
// conforming to the "k8s.container.cpu.limit.utilization" semantic conventions.
// It represents the ratio of container CPU usage to its current CPU limit.
//
// Stability: Development
type ContainerCPULimitUtilization struct {
	metric.Float64Gauge
}

var newContainerCPULimitUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of container CPU usage to its current CPU limit."),
	metric.WithUnit("1"),
}

// NewContainerCPULimitUtilization returns a new ContainerCPULimitUtilization
// instrument.
//
// Stability: Development
func NewContainerCPULimitUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (ContainerCPULimitUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimitUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPULimitUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.container.cpu.limit.utilization",
		opt...,
	)
	if err != nil {
		return ContainerCPULimitUtilization{noop.Float64Gauge{}}, err
	}
	return ContainerCPULimitUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimitUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimitUtilization) Name() string {
	return "k8s.container.cpu.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimitUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimitUtilization) Description() string {
	return "The ratio of container CPU usage to its current CPU limit."
}

// Record records val to the current distribution for attrs.
//
// The current CPU limit reflects the actual resources applied to the container,
// as reported by
// [ContainerStatus].
// The value range is [0.0,1.0]. A value of 1.0 means the container is using 100%
// of its actual CPU limit.
// If the CPU limit is not set, this metric SHOULD NOT be emitted for that
// container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerCPULimitUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The current CPU limit reflects the actual resources applied to the container,
// as reported by
// [ContainerStatus].
// The value range is [0.0,1.0]. A value of 1.0 means the container is using 100%
// of its actual CPU limit.
// If the CPU limit is not set, this metric SHOULD NOT be emitted for that
// container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerCPULimitUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// ContainerCPULimitUtilizationObservable is an instrument used to record metric
// values conforming to the "k8s.container.cpu.limit.utilization" semantic
// conventions. It represents the ratio of container CPU usage to its current CPU
// limit.
//
// Stability: Development
type ContainerCPULimitUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newContainerCPULimitUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of container CPU usage to its current CPU limit."),
	metric.WithUnit("1"),
}

// NewContainerCPULimitUtilizationObservable returns a new
// ContainerCPULimitUtilizationObservable instrument.
//
// Stability: Development
func NewContainerCPULimitUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (ContainerCPULimitUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPULimitUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPULimitUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPULimitUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.container.cpu.limit.utilization",
		opt...,
	)
	if err != nil {
		return ContainerCPULimitUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return ContainerCPULimitUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPULimitUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPULimitUtilizationObservable) Name() string {
	return "k8s.container.cpu.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPULimitUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPULimitUtilizationObservable) Description() string {
	return "The ratio of container CPU usage to its current CPU limit."
}

// ContainerCPURequestCurrent is an instrument used to record metric values
// conforming to the "k8s.container.cpu.request.current" semantic conventions. It
// represents the CPU resource requested currently configured for a running
// container.
//
// Stability: Development
type ContainerCPURequestCurrent struct {
	metric.Float64UpDownCounter
}

var newContainerCPURequestCurrentOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("CPU resource requested currently configured for a running container."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPURequestCurrent returns a new ContainerCPURequestCurrent
// instrument.
//
// Stability: Development
func NewContainerCPURequestCurrent(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ContainerCPURequestCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequestCurrent{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPURequestCurrentOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.container.cpu.request.current",
		opt...,
	)
	if err != nil {
		return ContainerCPURequestCurrent{noop.Float64UpDownCounter{}}, err
	}
	return ContainerCPURequestCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequestCurrent) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequestCurrent) Name() string {
	return "k8s.container.cpu.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequestCurrent) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequestCurrent) Description() string {
	return "CPU resource requested currently configured for a running container."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPURequestCurrent) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPURequestCurrent) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerCPURequestCurrentObservable is an instrument used to record metric
// values conforming to the "k8s.container.cpu.request.current" semantic
// conventions. It represents the CPU resource requested currently configured for
// a running container.
//
// Stability: Development
type ContainerCPURequestCurrentObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newContainerCPURequestCurrentObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("CPU resource requested currently configured for a running container."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPURequestCurrentObservable returns a new
// ContainerCPURequestCurrentObservable instrument.
//
// Stability: Development
func NewContainerCPURequestCurrentObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ContainerCPURequestCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequestCurrentObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPURequestCurrentObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.container.cpu.request.current",
		opt...,
	)
	if err != nil {
		return ContainerCPURequestCurrentObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ContainerCPURequestCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequestCurrentObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequestCurrentObservable) Name() string {
	return "k8s.container.cpu.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequestCurrentObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequestCurrentObservable) Description() string {
	return "CPU resource requested currently configured for a running container."
}

// ContainerCPURequestDesired is an instrument used to record metric values
// conforming to the "k8s.container.cpu.request.desired" semantic conventions. It
// represents the CPU resource requested as defined by the container spec.
//
// Stability: Development
type ContainerCPURequestDesired struct {
	metric.Float64UpDownCounter
}

var newContainerCPURequestDesiredOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("CPU resource requested as defined by the container spec."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPURequestDesired returns a new ContainerCPURequestDesired
// instrument.
//
// Stability: Development
func NewContainerCPURequestDesired(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ContainerCPURequestDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequestDesired{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPURequestDesiredOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.container.cpu.request.desired",
		opt...,
	)
	if err != nil {
		return ContainerCPURequestDesired{noop.Float64UpDownCounter{}}, err
	}
	return ContainerCPURequestDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequestDesired) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequestDesired) Name() string {
	return "k8s.container.cpu.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequestDesired) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequestDesired) Description() string {
	return "CPU resource requested as defined by the container spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPURequestDesired) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerCPURequestDesired) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerCPURequestDesiredObservable is an instrument used to record metric
// values conforming to the "k8s.container.cpu.request.desired" semantic
// conventions. It represents the CPU resource requested as defined by the
// container spec.
//
// Stability: Development
type ContainerCPURequestDesiredObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newContainerCPURequestDesiredObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("CPU resource requested as defined by the container spec."),
	metric.WithUnit("{cpu}"),
}

// NewContainerCPURequestDesiredObservable returns a new
// ContainerCPURequestDesiredObservable instrument.
//
// Stability: Development
func NewContainerCPURequestDesiredObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ContainerCPURequestDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequestDesiredObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPURequestDesiredObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.container.cpu.request.desired",
		opt...,
	)
	if err != nil {
		return ContainerCPURequestDesiredObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ContainerCPURequestDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequestDesiredObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequestDesiredObservable) Name() string {
	return "k8s.container.cpu.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequestDesiredObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequestDesiredObservable) Description() string {
	return "CPU resource requested as defined by the container spec."
}

// ContainerCPURequestUtilization is an instrument used to record metric values
// conforming to the "k8s.container.cpu.request.utilization" semantic
// conventions. It represents the ratio of container CPU usage to its current CPU
// request.
//
// Stability: Development
type ContainerCPURequestUtilization struct {
	metric.Float64Gauge
}

var newContainerCPURequestUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of container CPU usage to its current CPU request."),
	metric.WithUnit("1"),
}

// NewContainerCPURequestUtilization returns a new ContainerCPURequestUtilization
// instrument.
//
// Stability: Development
func NewContainerCPURequestUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (ContainerCPURequestUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequestUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPURequestUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.container.cpu.request.utilization",
		opt...,
	)
	if err != nil {
		return ContainerCPURequestUtilization{noop.Float64Gauge{}}, err
	}
	return ContainerCPURequestUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequestUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequestUtilization) Name() string {
	return "k8s.container.cpu.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequestUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequestUtilization) Description() string {
	return "The ratio of container CPU usage to its current CPU request."
}

// Record records val to the current distribution for attrs.
//
// The current CPU request reflects the request applied to the running container,
// as reported by
// [ContainerStatus].
// A value of 1.0 means the container is using exactly its CPU request. Since a
// CPU request is a scheduling
// reservation and not a hard cap, values above 1.0 are possible and indicate the
// container is using more CPU
// than requested. If the CPU request is not set, this metric SHOULD NOT be
// emitted for that container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerCPURequestUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The current CPU request reflects the request applied to the running container,
// as reported by
// [ContainerStatus].
// A value of 1.0 means the container is using exactly its CPU request. Since a
// CPU request is a scheduling
// reservation and not a hard cap, values above 1.0 are possible and indicate the
// container is using more CPU
// than requested. If the CPU request is not set, this metric SHOULD NOT be
// emitted for that container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerCPURequestUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// ContainerCPURequestUtilizationObservable is an instrument used to record
// metric values conforming to the "k8s.container.cpu.request.utilization"
// semantic conventions. It represents the ratio of container CPU usage to its
// current CPU request.
//
// Stability: Development
type ContainerCPURequestUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newContainerCPURequestUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of container CPU usage to its current CPU request."),
	metric.WithUnit("1"),
}

// NewContainerCPURequestUtilizationObservable returns a new
// ContainerCPURequestUtilizationObservable instrument.
//
// Stability: Development
func NewContainerCPURequestUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (ContainerCPURequestUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerCPURequestUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerCPURequestUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerCPURequestUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.container.cpu.request.utilization",
		opt...,
	)
	if err != nil {
		return ContainerCPURequestUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return ContainerCPURequestUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerCPURequestUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerCPURequestUtilizationObservable) Name() string {
	return "k8s.container.cpu.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerCPURequestUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerCPURequestUtilizationObservable) Description() string {
	return "The ratio of container CPU usage to its current CPU request."
}

// ContainerEphemeralStorageLimit is an instrument used to record metric values
// conforming to the "k8s.container.ephemeral_storage.limit" semantic
// conventions. It represents the maximum ephemeral storage resource limit set
// for the container.
//
// Stability: Development
type ContainerEphemeralStorageLimit struct {
	metric.Int64UpDownCounter
}

var newContainerEphemeralStorageLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum ephemeral storage resource limit set for the container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageLimit returns a new ContainerEphemeralStorageLimit
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerEphemeralStorageLimitOpts...)
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
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerEphemeralStorageLimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerEphemeralStorageLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerEphemeralStorageLimitObservable is an instrument used to record
// metric values conforming to the "k8s.container.ephemeral_storage.limit"
// semantic conventions. It represents the maximum ephemeral storage resource
// limit set for the container.
//
// Stability: Development
type ContainerEphemeralStorageLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerEphemeralStorageLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum ephemeral storage resource limit set for the container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageLimitObservable returns a new
// ContainerEphemeralStorageLimitObservable instrument.
//
// Stability: Development
func NewContainerEphemeralStorageLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerEphemeralStorageLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerEphemeralStorageLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerEphemeralStorageLimitObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerEphemeralStorageLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.ephemeral_storage.limit",
		opt...,
	)
	if err != nil {
		return ContainerEphemeralStorageLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerEphemeralStorageLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerEphemeralStorageLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerEphemeralStorageLimitObservable) Name() string {
	return "k8s.container.ephemeral_storage.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerEphemeralStorageLimitObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerEphemeralStorageLimitObservable) Description() string {
	return "Maximum ephemeral storage resource limit set for the container."
}

// ContainerEphemeralStorageRequest is an instrument used to record metric values
// conforming to the "k8s.container.ephemeral_storage.request" semantic
// conventions. It represents the ephemeral storage resource requested for the
// container.
//
// Stability: Development
type ContainerEphemeralStorageRequest struct {
	metric.Int64UpDownCounter
}

var newContainerEphemeralStorageRequestOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Ephemeral storage resource requested for the container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageRequest returns a new
// ContainerEphemeralStorageRequest instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerEphemeralStorageRequestOpts...)
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
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerEphemeralStorageRequest) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerEphemeralStorageRequest) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerEphemeralStorageRequestObservable is an instrument used to record
// metric values conforming to the "k8s.container.ephemeral_storage.request"
// semantic conventions. It represents the ephemeral storage resource requested
// for the container.
//
// Stability: Development
type ContainerEphemeralStorageRequestObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerEphemeralStorageRequestObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Ephemeral storage resource requested for the container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageRequestObservable returns a new
// ContainerEphemeralStorageRequestObservable instrument.
//
// Stability: Development
func NewContainerEphemeralStorageRequestObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerEphemeralStorageRequestObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerEphemeralStorageRequestObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerEphemeralStorageRequestObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerEphemeralStorageRequestObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.ephemeral_storage.request",
		opt...,
	)
	if err != nil {
		return ContainerEphemeralStorageRequestObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerEphemeralStorageRequestObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerEphemeralStorageRequestObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerEphemeralStorageRequestObservable) Name() string {
	return "k8s.container.ephemeral_storage.request"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerEphemeralStorageRequestObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerEphemeralStorageRequestObservable) Description() string {
	return "Ephemeral storage resource requested for the container."
}

// ContainerEphemeralStorageUsage is an instrument used to record metric values
// conforming to the "k8s.container.ephemeral_storage.usage" semantic
// conventions. It represents the ephemeral storage used by a container.
//
// Stability: Development
type ContainerEphemeralStorageUsage struct {
	metric.Int64UpDownCounter
}

var newContainerEphemeralStorageUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The ephemeral storage used by a container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageUsage returns a new ContainerEphemeralStorageUsage
// instrument.
//
// Stability: Development
func NewContainerEphemeralStorageUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerEphemeralStorageUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerEphemeralStorageUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerEphemeralStorageUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerEphemeralStorageUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.ephemeral_storage.usage",
		opt...,
	)
	if err != nil {
		return ContainerEphemeralStorageUsage{noop.Int64UpDownCounter{}}, err
	}
	return ContainerEphemeralStorageUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerEphemeralStorageUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerEphemeralStorageUsage) Name() string {
	return "k8s.container.ephemeral_storage.usage"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerEphemeralStorageUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerEphemeralStorageUsage) Description() string {
	return "The ephemeral storage used by a container."
}

// Add adds incr to the existing count for attrs.
//
// The containerEphemeralStorageFsType is the the type of file system component
// for ephemeral storage.
//
// The value for this metric can be compared against
// `metric.k8s.container.ephemeral_storage.request` and
// `metric.k8s.container.ephemeral_storage.limit`.
func (m ContainerEphemeralStorageUsage) Add(
	ctx context.Context,
	incr int64,
	containerEphemeralStorageFsType ContainerEphemeralStorageFsTypeAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.container.ephemeral_storage.fs_type", string(containerEphemeralStorageFsType)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("k8s.container.ephemeral_storage.fs_type", string(containerEphemeralStorageFsType)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// The value for this metric can be compared against
// `metric.k8s.container.ephemeral_storage.request` and
// `metric.k8s.container.ephemeral_storage.limit`.
func (m ContainerEphemeralStorageUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerEphemeralStorageUsageObservable is an instrument used to record
// metric values conforming to the "k8s.container.ephemeral_storage.usage"
// semantic conventions. It represents the ephemeral storage used by a container.
//
// Stability: Development
type ContainerEphemeralStorageUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerEphemeralStorageUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The ephemeral storage used by a container."),
	metric.WithUnit("By"),
}

// NewContainerEphemeralStorageUsageObservable returns a new
// ContainerEphemeralStorageUsageObservable instrument.
//
// Stability: Development
func NewContainerEphemeralStorageUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerEphemeralStorageUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerEphemeralStorageUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerEphemeralStorageUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerEphemeralStorageUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.ephemeral_storage.usage",
		opt...,
	)
	if err != nil {
		return ContainerEphemeralStorageUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerEphemeralStorageUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerEphemeralStorageUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerEphemeralStorageUsageObservable) Name() string {
	return "k8s.container.ephemeral_storage.usage"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerEphemeralStorageUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerEphemeralStorageUsageObservable) Description() string {
	return "The ephemeral storage used by a container."
}

// AttrContainerEphemeralStorageFsType returns a required attribute for the
// "k8s.container.ephemeral_storage.fs_type" semantic convention. It represents
// the type of file system component for ephemeral storage.
func (ContainerEphemeralStorageUsageObservable) AttrContainerEphemeralStorageFsType(val ContainerEphemeralStorageFsTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.container.ephemeral_storage.fs_type", string(val))
}

// ContainerMemoryLimitCurrent is an instrument used to record metric values
// conforming to the "k8s.container.memory.limit.current" semantic conventions.
// It represents the maximum memory resource limit currently configured for a
// running container.
//
// Stability: Development
type ContainerMemoryLimitCurrent struct {
	metric.Int64UpDownCounter
}

var newContainerMemoryLimitCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit currently configured for a running container."),
	metric.WithUnit("By"),
}

// NewContainerMemoryLimitCurrent returns a new ContainerMemoryLimitCurrent
// instrument.
//
// Stability: Development
func NewContainerMemoryLimitCurrent(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryLimitCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimitCurrent{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryLimitCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryLimitCurrentOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.limit.current",
		opt...,
	)
	if err != nil {
		return ContainerMemoryLimitCurrent{noop.Int64UpDownCounter{}}, err
	}
	return ContainerMemoryLimitCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryLimitCurrent) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryLimitCurrent) Name() string {
	return "k8s.container.memory.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryLimitCurrent) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryLimitCurrent) Description() string {
	return "Maximum memory resource limit currently configured for a running container."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryLimitCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryLimitCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerMemoryLimitCurrentObservable is an instrument used to record metric
// values conforming to the "k8s.container.memory.limit.current" semantic
// conventions. It represents the maximum memory resource limit currently
// configured for a running container.
//
// Stability: Development
type ContainerMemoryLimitCurrentObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerMemoryLimitCurrentObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit currently configured for a running container."),
	metric.WithUnit("By"),
}

// NewContainerMemoryLimitCurrentObservable returns a new
// ContainerMemoryLimitCurrentObservable instrument.
//
// Stability: Development
func NewContainerMemoryLimitCurrentObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerMemoryLimitCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimitCurrentObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryLimitCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryLimitCurrentObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.memory.limit.current",
		opt...,
	)
	if err != nil {
		return ContainerMemoryLimitCurrentObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerMemoryLimitCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryLimitCurrentObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryLimitCurrentObservable) Name() string {
	return "k8s.container.memory.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryLimitCurrentObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryLimitCurrentObservable) Description() string {
	return "Maximum memory resource limit currently configured for a running container."
}

// ContainerMemoryLimitDesired is an instrument used to record metric values
// conforming to the "k8s.container.memory.limit.desired" semantic conventions.
// It represents the maximum memory resource limit as defined by the container
// spec.
//
// Stability: Development
type ContainerMemoryLimitDesired struct {
	metric.Int64UpDownCounter
}

var newContainerMemoryLimitDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit as defined by the container spec."),
	metric.WithUnit("By"),
}

// NewContainerMemoryLimitDesired returns a new ContainerMemoryLimitDesired
// instrument.
//
// Stability: Development
func NewContainerMemoryLimitDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryLimitDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimitDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryLimitDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryLimitDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.limit.desired",
		opt...,
	)
	if err != nil {
		return ContainerMemoryLimitDesired{noop.Int64UpDownCounter{}}, err
	}
	return ContainerMemoryLimitDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryLimitDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryLimitDesired) Name() string {
	return "k8s.container.memory.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryLimitDesired) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryLimitDesired) Description() string {
	return "Maximum memory resource limit as defined by the container spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryLimitDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryLimitDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerMemoryLimitDesiredObservable is an instrument used to record metric
// values conforming to the "k8s.container.memory.limit.desired" semantic
// conventions. It represents the maximum memory resource limit as defined by the
// container spec.
//
// Stability: Development
type ContainerMemoryLimitDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerMemoryLimitDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit as defined by the container spec."),
	metric.WithUnit("By"),
}

// NewContainerMemoryLimitDesiredObservable returns a new
// ContainerMemoryLimitDesiredObservable instrument.
//
// Stability: Development
func NewContainerMemoryLimitDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerMemoryLimitDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimitDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryLimitDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryLimitDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.memory.limit.desired",
		opt...,
	)
	if err != nil {
		return ContainerMemoryLimitDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerMemoryLimitDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryLimitDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryLimitDesiredObservable) Name() string {
	return "k8s.container.memory.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryLimitDesiredObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryLimitDesiredObservable) Description() string {
	return "Maximum memory resource limit as defined by the container spec."
}

// ContainerMemoryLimitUtilization is an instrument used to record metric values
// conforming to the "k8s.container.memory.limit.utilization" semantic
// conventions. It represents the ratio of container memory working set to its
// current memory limit.
//
// Stability: Development
type ContainerMemoryLimitUtilization struct {
	metric.Float64Gauge
}

var newContainerMemoryLimitUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of container memory working set to its current memory limit."),
	metric.WithUnit("1"),
}

// NewContainerMemoryLimitUtilization returns a new
// ContainerMemoryLimitUtilization instrument.
//
// Stability: Development
func NewContainerMemoryLimitUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (ContainerMemoryLimitUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimitUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryLimitUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryLimitUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.container.memory.limit.utilization",
		opt...,
	)
	if err != nil {
		return ContainerMemoryLimitUtilization{noop.Float64Gauge{}}, err
	}
	return ContainerMemoryLimitUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryLimitUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryLimitUtilization) Name() string {
	return "k8s.container.memory.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryLimitUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryLimitUtilization) Description() string {
	return "The ratio of container memory working set to its current memory limit."
}

// Record records val to the current distribution for attrs.
//
// The numerator is the container memory working set, as reported by
// `container.memory.working_set`. The current memory limit
// reflects the actual resources applied to the container, as reported by
// [ContainerStatus].
// A value of 1.0 means the container is using 100% of its actual memory limit.
// Values above 1.0 can occur transiently, since memory usage can briefly exceed
// the limit.
// If the memory limit is not set, this metric SHOULD NOT be emitted for that
// container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerMemoryLimitUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The numerator is the container memory working set, as reported by
// `container.memory.working_set`. The current memory limit
// reflects the actual resources applied to the container, as reported by
// [ContainerStatus].
// A value of 1.0 means the container is using 100% of its actual memory limit.
// Values above 1.0 can occur transiently, since memory usage can briefly exceed
// the limit.
// If the memory limit is not set, this metric SHOULD NOT be emitted for that
// container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerMemoryLimitUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// ContainerMemoryLimitUtilizationObservable is an instrument used to record
// metric values conforming to the "k8s.container.memory.limit.utilization"
// semantic conventions. It represents the ratio of container memory working set
// to its current memory limit.
//
// Stability: Development
type ContainerMemoryLimitUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newContainerMemoryLimitUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of container memory working set to its current memory limit."),
	metric.WithUnit("1"),
}

// NewContainerMemoryLimitUtilizationObservable returns a new
// ContainerMemoryLimitUtilizationObservable instrument.
//
// Stability: Development
func NewContainerMemoryLimitUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (ContainerMemoryLimitUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimitUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryLimitUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryLimitUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.container.memory.limit.utilization",
		opt...,
	)
	if err != nil {
		return ContainerMemoryLimitUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return ContainerMemoryLimitUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryLimitUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryLimitUtilizationObservable) Name() string {
	return "k8s.container.memory.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryLimitUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryLimitUtilizationObservable) Description() string {
	return "The ratio of container memory working set to its current memory limit."
}

// ContainerMemoryRequestCurrent is an instrument used to record metric values
// conforming to the "k8s.container.memory.request.current" semantic conventions.
// It represents the memory resource request currently configured for a running
// container.
//
// Stability: Development
type ContainerMemoryRequestCurrent struct {
	metric.Int64UpDownCounter
}

var newContainerMemoryRequestCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Memory resource request currently configured for a running container."),
	metric.WithUnit("By"),
}

// NewContainerMemoryRequestCurrent returns a new ContainerMemoryRequestCurrent
// instrument.
//
// Stability: Development
func NewContainerMemoryRequestCurrent(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryRequestCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequestCurrent{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryRequestCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryRequestCurrentOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.request.current",
		opt...,
	)
	if err != nil {
		return ContainerMemoryRequestCurrent{noop.Int64UpDownCounter{}}, err
	}
	return ContainerMemoryRequestCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryRequestCurrent) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryRequestCurrent) Name() string {
	return "k8s.container.memory.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryRequestCurrent) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryRequestCurrent) Description() string {
	return "Memory resource request currently configured for a running container."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryRequestCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s ContainerStatus]
// (status.containerStatuses[*].resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryRequestCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerMemoryRequestCurrentObservable is an instrument used to record metric
// values conforming to the "k8s.container.memory.request.current" semantic
// conventions. It represents the memory resource request currently configured
// for a running container.
//
// Stability: Development
type ContainerMemoryRequestCurrentObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerMemoryRequestCurrentObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Memory resource request currently configured for a running container."),
	metric.WithUnit("By"),
}

// NewContainerMemoryRequestCurrentObservable returns a new
// ContainerMemoryRequestCurrentObservable instrument.
//
// Stability: Development
func NewContainerMemoryRequestCurrentObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerMemoryRequestCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequestCurrentObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryRequestCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryRequestCurrentObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.memory.request.current",
		opt...,
	)
	if err != nil {
		return ContainerMemoryRequestCurrentObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerMemoryRequestCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryRequestCurrentObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryRequestCurrentObservable) Name() string {
	return "k8s.container.memory.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryRequestCurrentObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryRequestCurrentObservable) Description() string {
	return "Memory resource request currently configured for a running container."
}

// ContainerMemoryRequestDesired is an instrument used to record metric values
// conforming to the "k8s.container.memory.request.desired" semantic conventions.
// It represents the memory resource requested as defined by the container spec.
//
// Stability: Development
type ContainerMemoryRequestDesired struct {
	metric.Int64UpDownCounter
}

var newContainerMemoryRequestDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Memory resource requested as defined by the container spec."),
	metric.WithUnit("By"),
}

// NewContainerMemoryRequestDesired returns a new ContainerMemoryRequestDesired
// instrument.
//
// Stability: Development
func NewContainerMemoryRequestDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryRequestDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequestDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryRequestDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryRequestDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.request.desired",
		opt...,
	)
	if err != nil {
		return ContainerMemoryRequestDesired{noop.Int64UpDownCounter{}}, err
	}
	return ContainerMemoryRequestDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryRequestDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryRequestDesired) Name() string {
	return "k8s.container.memory.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryRequestDesired) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryRequestDesired) Description() string {
	return "Memory resource requested as defined by the container spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryRequestDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s Container]
// (spec.containers[*].resources). Also see `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
// [K8s Container]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#container-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m ContainerMemoryRequestDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerMemoryRequestDesiredObservable is an instrument used to record metric
// values conforming to the "k8s.container.memory.request.desired" semantic
// conventions. It represents the memory resource requested as defined by the
// container spec.
//
// Stability: Development
type ContainerMemoryRequestDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerMemoryRequestDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Memory resource requested as defined by the container spec."),
	metric.WithUnit("By"),
}

// NewContainerMemoryRequestDesiredObservable returns a new
// ContainerMemoryRequestDesiredObservable instrument.
//
// Stability: Development
func NewContainerMemoryRequestDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerMemoryRequestDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequestDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryRequestDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryRequestDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.memory.request.desired",
		opt...,
	)
	if err != nil {
		return ContainerMemoryRequestDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerMemoryRequestDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryRequestDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryRequestDesiredObservable) Name() string {
	return "k8s.container.memory.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryRequestDesiredObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryRequestDesiredObservable) Description() string {
	return "Memory resource requested as defined by the container spec."
}

// ContainerMemoryRequestUtilization is an instrument used to record metric
// values conforming to the "k8s.container.memory.request.utilization" semantic
// conventions. It represents the ratio of container memory working set to its
// current memory request.
//
// Stability: Development
type ContainerMemoryRequestUtilization struct {
	metric.Float64Gauge
}

var newContainerMemoryRequestUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of container memory working set to its current memory request."),
	metric.WithUnit("1"),
}

// NewContainerMemoryRequestUtilization returns a new
// ContainerMemoryRequestUtilization instrument.
//
// Stability: Development
func NewContainerMemoryRequestUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (ContainerMemoryRequestUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequestUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryRequestUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryRequestUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.container.memory.request.utilization",
		opt...,
	)
	if err != nil {
		return ContainerMemoryRequestUtilization{noop.Float64Gauge{}}, err
	}
	return ContainerMemoryRequestUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryRequestUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryRequestUtilization) Name() string {
	return "k8s.container.memory.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryRequestUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryRequestUtilization) Description() string {
	return "The ratio of container memory working set to its current memory request."
}

// Record records val to the current distribution for attrs.
//
// The numerator is the container memory working set, as reported by
// `container.memory.working_set`. The current memory request
// reflects the request applied to the running container, as reported by
// [ContainerStatus].
// A value of 1.0 means the container is using exactly its memory request. Since
// a memory request is a
// scheduling reservation and not a hard cap, values above 1.0 are possible and
// indicate the container is
// using more memory than requested. If the memory request is not set, this
// metric SHOULD NOT be emitted for
// that container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerMemoryRequestUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The numerator is the container memory working set, as reported by
// `container.memory.working_set`. The current memory request
// reflects the request applied to the running container, as reported by
// [ContainerStatus].
// A value of 1.0 means the container is using exactly its memory request. Since
// a memory request is a
// scheduling reservation and not a hard cap, values above 1.0 are possible and
// indicate the container is
// using more memory than requested. If the memory request is not set, this
// metric SHOULD NOT be emitted for
// that container.
//
// [ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerMemoryRequestUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// ContainerMemoryRequestUtilizationObservable is an instrument used to record
// metric values conforming to the "k8s.container.memory.request.utilization"
// semantic conventions. It represents the ratio of container memory working set
// to its current memory request.
//
// Stability: Development
type ContainerMemoryRequestUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newContainerMemoryRequestUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of container memory working set to its current memory request."),
	metric.WithUnit("1"),
}

// NewContainerMemoryRequestUtilizationObservable returns a new
// ContainerMemoryRequestUtilizationObservable instrument.
//
// Stability: Development
func NewContainerMemoryRequestUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (ContainerMemoryRequestUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequestUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerMemoryRequestUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerMemoryRequestUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.container.memory.request.utilization",
		opt...,
	)
	if err != nil {
		return ContainerMemoryRequestUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return ContainerMemoryRequestUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerMemoryRequestUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (ContainerMemoryRequestUtilizationObservable) Name() string {
	return "k8s.container.memory.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerMemoryRequestUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (ContainerMemoryRequestUtilizationObservable) Description() string {
	return "The ratio of container memory working set to its current memory request."
}

// ContainerReady is an instrument used to record metric values conforming to the
// "k8s.container.ready" semantic conventions. It represents the indicates
// whether the container is currently marked as ready to accept traffic, based on
// its readiness probe (1 = ready, 0 = not ready).
//
// Stability: Development
type ContainerReady struct {
	metric.Int64UpDownCounter
}

var newContainerReadyOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready)."),
	metric.WithUnit("{container}"),
}

// NewContainerReady returns a new ContainerReady instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerReadyOpts...)
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
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerReady) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric SHOULD reflect the value of the `ready` field in the
// [K8s ContainerStatus].
//
// [K8s ContainerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatus-v1-core
func (m ContainerReady) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerReadyObservable is an instrument used to record metric values
// conforming to the "k8s.container.ready" semantic conventions. It represents
// the indicates whether the container is currently marked as ready to accept
// traffic, based on its readiness probe (1 = ready, 0 = not ready).
//
// Stability: Development
type ContainerReadyObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerReadyObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready)."),
	metric.WithUnit("{container}"),
}

// NewContainerReadyObservable returns a new ContainerReadyObservable instrument.
//
// Stability: Development
func NewContainerReadyObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerReadyObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerReadyObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerReadyObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerReadyObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.ready",
		opt...,
	)
	if err != nil {
		return ContainerReadyObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerReadyObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerReadyObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerReadyObservable) Name() string {
	return "k8s.container.ready"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerReadyObservable) Unit() string {
	return "{container}"
}

// Description returns the semantic convention description of the instrument
func (ContainerReadyObservable) Description() string {
	return "Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready)."
}

// ContainerRestartCount is an instrument used to record metric values conforming
// to the "k8s.container.restart.count" semantic conventions. It represents the
// describes how many times the container has restarted (since the last counter
// reset).
//
// Stability: Development
type ContainerRestartCount struct {
	metric.Int64UpDownCounter
}

var newContainerRestartCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes how many times the container has restarted (since the last counter reset)."),
	metric.WithUnit("{restart}"),
}

// NewContainerRestartCount returns a new ContainerRestartCount instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerRestartCountOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerRestartCountObservable is an instrument used to record metric values
// conforming to the "k8s.container.restart.count" semantic conventions. It
// represents the describes how many times the container has restarted (since the
// last counter reset).
//
// Stability: Development
type ContainerRestartCountObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerRestartCountObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Describes how many times the container has restarted (since the last counter reset)."),
	metric.WithUnit("{restart}"),
}

// NewContainerRestartCountObservable returns a new
// ContainerRestartCountObservable instrument.
//
// Stability: Development
func NewContainerRestartCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerRestartCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerRestartCountObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerRestartCountObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerRestartCountObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.restart.count",
		opt...,
	)
	if err != nil {
		return ContainerRestartCountObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerRestartCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerRestartCountObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerRestartCountObservable) Name() string {
	return "k8s.container.restart.count"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerRestartCountObservable) Unit() string {
	return "{restart}"
}

// Description returns the semantic convention description of the instrument
func (ContainerRestartCountObservable) Description() string {
	return "Describes how many times the container has restarted (since the last counter reset)."
}

// ContainerStatusReason is an instrument used to record metric values conforming
// to the "k8s.container.status.reason" semantic conventions. It represents the
// describes the number of K8s containers that are currently in a state for a
// given reason.
//
// Stability: Development
type ContainerStatusReason struct {
	metric.Int64UpDownCounter
}

var newContainerStatusReasonOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the number of K8s containers that are currently in a state for a given reason."),
	metric.WithUnit("{container}"),
}

// NewContainerStatusReason returns a new ContainerStatusReason instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerStatusReasonOpts...)
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
// [K8s ContainerStateWaiting]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatewaiting-v1-core
// [K8s ContainerStateTerminated]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstateterminated-v1-core
func (m ContainerStatusReason) Add(
	ctx context.Context,
	incr int64,
	containerStatusReason ContainerStatusReasonAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.container.status.reason", string(containerStatusReason)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStatusReasonObservable is an instrument used to record metric values
// conforming to the "k8s.container.status.reason" semantic conventions. It
// represents the describes the number of K8s containers that are currently in a
// state for a given reason.
//
// Stability: Development
type ContainerStatusReasonObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerStatusReasonObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Describes the number of K8s containers that are currently in a state for a given reason."),
	metric.WithUnit("{container}"),
}

// NewContainerStatusReasonObservable returns a new
// ContainerStatusReasonObservable instrument.
//
// Stability: Development
func NewContainerStatusReasonObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerStatusReasonObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStatusReasonObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStatusReasonObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerStatusReasonObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.status.reason",
		opt...,
	)
	if err != nil {
		return ContainerStatusReasonObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerStatusReasonObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStatusReasonObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStatusReasonObservable) Name() string {
	return "k8s.container.status.reason"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStatusReasonObservable) Unit() string {
	return "{container}"
}

// Description returns the semantic convention description of the instrument
func (ContainerStatusReasonObservable) Description() string {
	return "Describes the number of K8s containers that are currently in a state for a given reason."
}

// AttrContainerStatusReason returns a required attribute for the
// "k8s.container.status.reason" semantic convention. It represents the reason
// for the container state. Corresponds to the `reason` field of the:
// [K8s ContainerStateWaiting] or [K8s ContainerStateTerminated].
//
// [K8s ContainerStateWaiting]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstatewaiting-v1-core
// [K8s ContainerStateTerminated]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstateterminated-v1-core
func (ContainerStatusReasonObservable) AttrContainerStatusReason(val ContainerStatusReasonAttr) attribute.KeyValue {
	return attribute.String("k8s.container.status.reason", string(val))
}

// ContainerStatusState is an instrument used to record metric values conforming
// to the "k8s.container.status.state" semantic conventions. It represents the
// describes the number of K8s containers that are currently in a given state.
//
// Stability: Development
type ContainerStatusState struct {
	metric.Int64UpDownCounter
}

var newContainerStatusStateOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the number of K8s containers that are currently in a given state."),
	metric.WithUnit("{container}"),
}

// NewContainerStatusState returns a new ContainerStatusState instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerStatusStateOpts...)
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
// [K8s ContainerState]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstate-v1-core
func (m ContainerStatusState) Add(
	ctx context.Context,
	incr int64,
	containerStatusState ContainerStatusStateAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.container.status.state", string(containerStatusState)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStatusStateObservable is an instrument used to record metric values
// conforming to the "k8s.container.status.state" semantic conventions. It
// represents the describes the number of K8s containers that are currently in a
// given state.
//
// Stability: Development
type ContainerStatusStateObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerStatusStateObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Describes the number of K8s containers that are currently in a given state."),
	metric.WithUnit("{container}"),
}

// NewContainerStatusStateObservable returns a new ContainerStatusStateObservable
// instrument.
//
// Stability: Development
func NewContainerStatusStateObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerStatusStateObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStatusStateObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStatusStateObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerStatusStateObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.status.state",
		opt...,
	)
	if err != nil {
		return ContainerStatusStateObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerStatusStateObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStatusStateObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStatusStateObservable) Name() string {
	return "k8s.container.status.state"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStatusStateObservable) Unit() string {
	return "{container}"
}

// Description returns the semantic convention description of the instrument
func (ContainerStatusStateObservable) Description() string {
	return "Describes the number of K8s containers that are currently in a given state."
}

// AttrContainerStatusState returns a required attribute for the
// "k8s.container.status.state" semantic convention. It represents the state of
// the container. [K8s ContainerState].
//
// [K8s ContainerState]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#containerstate-v1-core
func (ContainerStatusStateObservable) AttrContainerStatusState(val ContainerStatusStateAttr) attribute.KeyValue {
	return attribute.String("k8s.container.status.state", string(val))
}

// ContainerStorageLimit is an instrument used to record metric values conforming
// to the "k8s.container.storage.limit" semantic conventions. It represents the
// maximum storage resource limit set for the container.
//
// Stability: Development
type ContainerStorageLimit struct {
	metric.Int64UpDownCounter
}

var newContainerStorageLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum storage resource limit set for the container."),
	metric.WithUnit("By"),
}

// NewContainerStorageLimit returns a new ContainerStorageLimit instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerStorageLimitOpts...)
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
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerStorageLimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerStorageLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStorageLimitObservable is an instrument used to record metric values
// conforming to the "k8s.container.storage.limit" semantic conventions. It
// represents the maximum storage resource limit set for the container.
//
// Stability: Development
type ContainerStorageLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerStorageLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum storage resource limit set for the container."),
	metric.WithUnit("By"),
}

// NewContainerStorageLimitObservable returns a new
// ContainerStorageLimitObservable instrument.
//
// Stability: Development
func NewContainerStorageLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerStorageLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStorageLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStorageLimitObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerStorageLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.storage.limit",
		opt...,
	)
	if err != nil {
		return ContainerStorageLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerStorageLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStorageLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStorageLimitObservable) Name() string {
	return "k8s.container.storage.limit"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStorageLimitObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerStorageLimitObservable) Description() string {
	return "Maximum storage resource limit set for the container."
}

// ContainerStorageRequest is an instrument used to record metric values
// conforming to the "k8s.container.storage.request" semantic conventions. It
// represents the storage resource requested for the container.
//
// Stability: Development
type ContainerStorageRequest struct {
	metric.Int64UpDownCounter
}

var newContainerStorageRequestOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Storage resource requested for the container."),
	metric.WithUnit("By"),
}

// NewContainerStorageRequest returns a new ContainerStorageRequest instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newContainerStorageRequestOpts...)
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
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerStorageRequest) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// See the [Kubernetes ResourceRequirements API] for details.
//
// [Kubernetes ResourceRequirements API]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcerequirements-v1-core
func (m ContainerStorageRequest) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ContainerStorageRequestObservable is an instrument used to record metric
// values conforming to the "k8s.container.storage.request" semantic conventions.
// It represents the storage resource requested for the container.
//
// Stability: Development
type ContainerStorageRequestObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newContainerStorageRequestObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Storage resource requested for the container."),
	metric.WithUnit("By"),
}

// NewContainerStorageRequestObservable returns a new
// ContainerStorageRequestObservable instrument.
//
// Stability: Development
func NewContainerStorageRequestObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ContainerStorageRequestObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStorageRequestObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newContainerStorageRequestObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newContainerStorageRequestObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.container.storage.request",
		opt...,
	)
	if err != nil {
		return ContainerStorageRequestObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ContainerStorageRequestObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ContainerStorageRequestObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ContainerStorageRequestObservable) Name() string {
	return "k8s.container.storage.request"
}

// Unit returns the semantic convention unit of the instrument
func (ContainerStorageRequestObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ContainerStorageRequestObservable) Description() string {
	return "Storage resource requested for the container."
}

// CronJobJobActive is an instrument used to record metric values conforming to
// the "k8s.cronjob.job.active" semantic conventions. It represents the number of
// actively running jobs for a cronjob.
//
// Stability: Development
type CronJobJobActive struct {
	metric.Int64UpDownCounter
}

var newCronJobJobActiveOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of actively running jobs for a cronjob."),
	metric.WithUnit("{job}"),
}

// NewCronJobJobActive returns a new CronJobJobActive instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newCronJobJobActiveOpts...)
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
// [K8s CronJobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#cronjobstatus-v1-batch
func (m CronJobJobActive) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `active` field of the
// [K8s CronJobStatus].
//
// [K8s CronJobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#cronjobstatus-v1-batch
func (m CronJobJobActive) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// CronJobJobActiveObservable is an instrument used to record metric values
// conforming to the "k8s.cronjob.job.active" semantic conventions. It represents
// the number of actively running jobs for a cronjob.
//
// Stability: Development
type CronJobJobActiveObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newCronJobJobActiveObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of actively running jobs for a cronjob."),
	metric.WithUnit("{job}"),
}

// NewCronJobJobActiveObservable returns a new CronJobJobActiveObservable
// instrument.
//
// Stability: Development
func NewCronJobJobActiveObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (CronJobJobActiveObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return CronJobJobActiveObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newCronJobJobActiveObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newCronJobJobActiveObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.cronjob.job.active",
		opt...,
	)
	if err != nil {
		return CronJobJobActiveObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return CronJobJobActiveObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m CronJobJobActiveObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (CronJobJobActiveObservable) Name() string {
	return "k8s.cronjob.job.active"
}

// Unit returns the semantic convention unit of the instrument
func (CronJobJobActiveObservable) Unit() string {
	return "{job}"
}

// Description returns the semantic convention description of the instrument
func (CronJobJobActiveObservable) Description() string {
	return "The number of actively running jobs for a cronjob."
}

// DaemonSetNodeCurrentScheduled is an instrument used to record metric values
// conforming to the "k8s.daemonset.node.current_scheduled" semantic conventions.
// It represents the number of nodes that are running at least 1 daemon pod and
// are supposed to run the daemon pod.
//
// Stability: Development
type DaemonSetNodeCurrentScheduled struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeCurrentScheduledOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeCurrentScheduled returns a new DaemonSetNodeCurrentScheduled
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeCurrentScheduledOpts...)
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
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeCurrentScheduled) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeCurrentScheduled) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeCurrentScheduledObservable is an instrument used to record metric
// values conforming to the "k8s.daemonset.node.current_scheduled" semantic
// conventions. It represents the number of nodes that are running at least 1
// daemon pod and are supposed to run the daemon pod.
//
// Stability: Development
type DaemonSetNodeCurrentScheduledObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newDaemonSetNodeCurrentScheduledObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeCurrentScheduledObservable returns a new
// DaemonSetNodeCurrentScheduledObservable instrument.
//
// Stability: Development
func NewDaemonSetNodeCurrentScheduledObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (DaemonSetNodeCurrentScheduledObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeCurrentScheduledObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeCurrentScheduledObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeCurrentScheduledObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.daemonset.node.current_scheduled",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeCurrentScheduledObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return DaemonSetNodeCurrentScheduledObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeCurrentScheduledObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeCurrentScheduledObservable) Name() string {
	return "k8s.daemonset.node.current_scheduled"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeCurrentScheduledObservable) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeCurrentScheduledObservable) Description() string {
	return "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod."
}

// DaemonSetNodeDesiredScheduled is an instrument used to record metric values
// conforming to the "k8s.daemonset.node.desired_scheduled" semantic conventions.
// It represents the number of nodes that should be running the daemon pod
// (including nodes currently running the daemon pod).
//
// Stability: Development
type DaemonSetNodeDesiredScheduled struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeDesiredScheduledOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeDesiredScheduled returns a new DaemonSetNodeDesiredScheduled
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeDesiredScheduledOpts...)
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
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeDesiredScheduled) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `desiredNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeDesiredScheduled) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeDesiredScheduledObservable is an instrument used to record metric
// values conforming to the "k8s.daemonset.node.desired_scheduled" semantic
// conventions. It represents the number of nodes that should be running the
// daemon pod (including nodes currently running the daemon pod).
//
// Stability: Development
type DaemonSetNodeDesiredScheduledObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newDaemonSetNodeDesiredScheduledObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeDesiredScheduledObservable returns a new
// DaemonSetNodeDesiredScheduledObservable instrument.
//
// Stability: Development
func NewDaemonSetNodeDesiredScheduledObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (DaemonSetNodeDesiredScheduledObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeDesiredScheduledObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeDesiredScheduledObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeDesiredScheduledObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.daemonset.node.desired_scheduled",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeDesiredScheduledObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return DaemonSetNodeDesiredScheduledObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeDesiredScheduledObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeDesiredScheduledObservable) Name() string {
	return "k8s.daemonset.node.desired_scheduled"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeDesiredScheduledObservable) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeDesiredScheduledObservable) Description() string {
	return "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)."
}

// DaemonSetNodeMisscheduled is an instrument used to record metric values
// conforming to the "k8s.daemonset.node.misscheduled" semantic conventions. It
// represents the number of nodes that are running the daemon pod, but are not
// supposed to run the daemon pod.
//
// Stability: Development
type DaemonSetNodeMisscheduled struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeMisscheduledOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeMisscheduled returns a new DaemonSetNodeMisscheduled
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeMisscheduledOpts...)
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
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeMisscheduled) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `numberMisscheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeMisscheduled) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeMisscheduledObservable is an instrument used to record metric
// values conforming to the "k8s.daemonset.node.misscheduled" semantic
// conventions. It represents the number of nodes that are running the daemon
// pod, but are not supposed to run the daemon pod.
//
// Stability: Development
type DaemonSetNodeMisscheduledObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newDaemonSetNodeMisscheduledObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeMisscheduledObservable returns a new
// DaemonSetNodeMisscheduledObservable instrument.
//
// Stability: Development
func NewDaemonSetNodeMisscheduledObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (DaemonSetNodeMisscheduledObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeMisscheduledObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeMisscheduledObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeMisscheduledObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.daemonset.node.misscheduled",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeMisscheduledObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return DaemonSetNodeMisscheduledObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeMisscheduledObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeMisscheduledObservable) Name() string {
	return "k8s.daemonset.node.misscheduled"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeMisscheduledObservable) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeMisscheduledObservable) Description() string {
	return "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod."
}

// DaemonSetNodeReady is an instrument used to record metric values conforming to
// the "k8s.daemonset.node.ready" semantic conventions. It represents the number
// of nodes that should be running the daemon pod and have one or more of the
// daemon pod running and ready.
//
// Stability: Development
type DaemonSetNodeReady struct {
	metric.Int64UpDownCounter
}

var newDaemonSetNodeReadyOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeReady returns a new DaemonSetNodeReady instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeReadyOpts...)
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
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeReady) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `numberReady` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#daemonsetstatus-v1-apps
func (m DaemonSetNodeReady) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DaemonSetNodeReadyObservable is an instrument used to record metric values
// conforming to the "k8s.daemonset.node.ready" semantic conventions. It
// represents the number of nodes that should be running the daemon pod and have
// one or more of the daemon pod running and ready.
//
// Stability: Development
type DaemonSetNodeReadyObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newDaemonSetNodeReadyObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready."),
	metric.WithUnit("{node}"),
}

// NewDaemonSetNodeReadyObservable returns a new DaemonSetNodeReadyObservable
// instrument.
//
// Stability: Development
func NewDaemonSetNodeReadyObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (DaemonSetNodeReadyObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DaemonSetNodeReadyObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDaemonSetNodeReadyObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newDaemonSetNodeReadyObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.daemonset.node.ready",
		opt...,
	)
	if err != nil {
		return DaemonSetNodeReadyObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return DaemonSetNodeReadyObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DaemonSetNodeReadyObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DaemonSetNodeReadyObservable) Name() string {
	return "k8s.daemonset.node.ready"
}

// Unit returns the semantic convention unit of the instrument
func (DaemonSetNodeReadyObservable) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (DaemonSetNodeReadyObservable) Description() string {
	return "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready."
}

// DeploymentPodAvailable is an instrument used to record metric values
// conforming to the "k8s.deployment.pod.available" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this deployment.
//
// Stability: Development
type DeploymentPodAvailable struct {
	metric.Int64UpDownCounter
}

var newDeploymentPodAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment."),
	metric.WithUnit("{pod}"),
}

// NewDeploymentPodAvailable returns a new DeploymentPodAvailable instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newDeploymentPodAvailableOpts...)
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
// [K8s DeploymentStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#deploymentstatus-v1-apps
func (m DeploymentPodAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s DeploymentStatus].
//
// [K8s DeploymentStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#deploymentstatus-v1-apps
func (m DeploymentPodAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DeploymentPodAvailableObservable is an instrument used to record metric values
// conforming to the "k8s.deployment.pod.available" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this deployment.
//
// Stability: Development
type DeploymentPodAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newDeploymentPodAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment."),
	metric.WithUnit("{pod}"),
}

// NewDeploymentPodAvailableObservable returns a new
// DeploymentPodAvailableObservable instrument.
//
// Stability: Development
func NewDeploymentPodAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (DeploymentPodAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DeploymentPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDeploymentPodAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newDeploymentPodAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.deployment.pod.available",
		opt...,
	)
	if err != nil {
		return DeploymentPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return DeploymentPodAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DeploymentPodAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DeploymentPodAvailableObservable) Name() string {
	return "k8s.deployment.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (DeploymentPodAvailableObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (DeploymentPodAvailableObservable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment."
}

// DeploymentPodDesired is an instrument used to record metric values conforming
// to the "k8s.deployment.pod.desired" semantic conventions. It represents the
// number of desired replica pods in this deployment.
//
// Stability: Development
type DeploymentPodDesired struct {
	metric.Int64UpDownCounter
}

var newDeploymentPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this deployment."),
	metric.WithUnit("{pod}"),
}

// NewDeploymentPodDesired returns a new DeploymentPodDesired instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newDeploymentPodDesiredOpts...)
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
// [K8s DeploymentSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#deploymentspec-v1-apps
func (m DeploymentPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s DeploymentSpec].
//
// [K8s DeploymentSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#deploymentspec-v1-apps
func (m DeploymentPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// DeploymentPodDesiredObservable is an instrument used to record metric values
// conforming to the "k8s.deployment.pod.desired" semantic conventions. It
// represents the number of desired replica pods in this deployment.
//
// Stability: Development
type DeploymentPodDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newDeploymentPodDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this deployment."),
	metric.WithUnit("{pod}"),
}

// NewDeploymentPodDesiredObservable returns a new DeploymentPodDesiredObservable
// instrument.
//
// Stability: Development
func NewDeploymentPodDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (DeploymentPodDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return DeploymentPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newDeploymentPodDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newDeploymentPodDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.deployment.pod.desired",
		opt...,
	)
	if err != nil {
		return DeploymentPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return DeploymentPodDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m DeploymentPodDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (DeploymentPodDesiredObservable) Name() string {
	return "k8s.deployment.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (DeploymentPodDesiredObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (DeploymentPodDesiredObservable) Description() string {
	return "Number of desired replica pods in this deployment."
}

// HPAMetricTargetCPUAverageUtilization is an instrument used to record metric
// values conforming to the "k8s.hpa.metric.target.cpu.average_utilization"
// semantic conventions. It represents the target average utilization, in
// percentage, for CPU resource in HPA config.
//
// Stability: Development
type HPAMetricTargetCPUAverageUtilization struct {
	metric.Int64Gauge
}

var newHPAMetricTargetCPUAverageUtilizationOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Target average utilization, in percentage, for CPU resource in HPA config."),
	metric.WithUnit("1"),
}

// NewHPAMetricTargetCPUAverageUtilization returns a new
// HPAMetricTargetCPUAverageUtilization instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newHPAMetricTargetCPUAverageUtilizationOpts...)
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
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageUtilization) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

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
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageUtilization) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

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

// HPAMetricTargetCPUAverageUtilizationObservable is an instrument used to record
// metric values conforming to the
// "k8s.hpa.metric.target.cpu.average_utilization" semantic conventions. It
// represents the target average utilization, in percentage, for CPU resource in
// HPA config.
//
// Stability: Development
type HPAMetricTargetCPUAverageUtilizationObservable struct {
	metric.Int64ObservableGauge
}

var newHPAMetricTargetCPUAverageUtilizationObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Target average utilization, in percentage, for CPU resource in HPA config."),
	metric.WithUnit("1"),
}

// NewHPAMetricTargetCPUAverageUtilizationObservable returns a new
// HPAMetricTargetCPUAverageUtilizationObservable instrument.
//
// Stability: Development
func NewHPAMetricTargetCPUAverageUtilizationObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (HPAMetricTargetCPUAverageUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMetricTargetCPUAverageUtilizationObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAMetricTargetCPUAverageUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newHPAMetricTargetCPUAverageUtilizationObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"k8s.hpa.metric.target.cpu.average_utilization",
		opt...,
	)
	if err != nil {
		return HPAMetricTargetCPUAverageUtilizationObservable{noop.Int64ObservableGauge{}}, err
	}
	return HPAMetricTargetCPUAverageUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMetricTargetCPUAverageUtilizationObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (HPAMetricTargetCPUAverageUtilizationObservable) Name() string {
	return "k8s.hpa.metric.target.cpu.average_utilization"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMetricTargetCPUAverageUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (HPAMetricTargetCPUAverageUtilizationObservable) Description() string {
	return "Target average utilization, in percentage, for CPU resource in HPA config."
}

// AttrContainerName returns an optional attribute for the "k8s.container.name"
// semantic convention. It represents the name of the Container from Pod
// specification, must be unique within a Pod. Container runtime usually uses
// different globally unique name (`container.name`).
func (HPAMetricTargetCPUAverageUtilizationObservable) AttrContainerName(val string) attribute.KeyValue {
	return attribute.String("k8s.container.name", val)
}

// AttrHPAMetricType returns an optional attribute for the "k8s.hpa.metric.type"
// semantic convention. It represents the type of metric source for the
// horizontal pod autoscaler.
func (HPAMetricTargetCPUAverageUtilizationObservable) AttrHPAMetricType(val string) attribute.KeyValue {
	return attribute.String("k8s.hpa.metric.type", val)
}

// HPAMetricTargetCPUAverageValue is an instrument used to record metric values
// conforming to the "k8s.hpa.metric.target.cpu.average_value" semantic
// conventions. It represents the target average value for CPU resource in HPA
// config.
//
// Stability: Development
type HPAMetricTargetCPUAverageValue struct {
	metric.Int64Gauge
}

var newHPAMetricTargetCPUAverageValueOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Target average value for CPU resource in HPA config."),
	metric.WithUnit("{cpu}"),
}

// NewHPAMetricTargetCPUAverageValue returns a new HPAMetricTargetCPUAverageValue
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newHPAMetricTargetCPUAverageValueOpts...)
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
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageValue) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

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
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUAverageValue) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

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

// HPAMetricTargetCPUAverageValueObservable is an instrument used to record
// metric values conforming to the "k8s.hpa.metric.target.cpu.average_value"
// semantic conventions. It represents the target average value for CPU resource
// in HPA config.
//
// Stability: Development
type HPAMetricTargetCPUAverageValueObservable struct {
	metric.Int64ObservableGauge
}

var newHPAMetricTargetCPUAverageValueObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Target average value for CPU resource in HPA config."),
	metric.WithUnit("{cpu}"),
}

// NewHPAMetricTargetCPUAverageValueObservable returns a new
// HPAMetricTargetCPUAverageValueObservable instrument.
//
// Stability: Development
func NewHPAMetricTargetCPUAverageValueObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (HPAMetricTargetCPUAverageValueObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMetricTargetCPUAverageValueObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAMetricTargetCPUAverageValueObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newHPAMetricTargetCPUAverageValueObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"k8s.hpa.metric.target.cpu.average_value",
		opt...,
	)
	if err != nil {
		return HPAMetricTargetCPUAverageValueObservable{noop.Int64ObservableGauge{}}, err
	}
	return HPAMetricTargetCPUAverageValueObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMetricTargetCPUAverageValueObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (HPAMetricTargetCPUAverageValueObservable) Name() string {
	return "k8s.hpa.metric.target.cpu.average_value"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMetricTargetCPUAverageValueObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (HPAMetricTargetCPUAverageValueObservable) Description() string {
	return "Target average value for CPU resource in HPA config."
}

// AttrContainerName returns an optional attribute for the "k8s.container.name"
// semantic convention. It represents the name of the Container from Pod
// specification, must be unique within a Pod. Container runtime usually uses
// different globally unique name (`container.name`).
func (HPAMetricTargetCPUAverageValueObservable) AttrContainerName(val string) attribute.KeyValue {
	return attribute.String("k8s.container.name", val)
}

// AttrHPAMetricType returns an optional attribute for the "k8s.hpa.metric.type"
// semantic convention. It represents the type of metric source for the
// horizontal pod autoscaler.
func (HPAMetricTargetCPUAverageValueObservable) AttrHPAMetricType(val string) attribute.KeyValue {
	return attribute.String("k8s.hpa.metric.type", val)
}

// HPAMetricTargetCPUValue is an instrument used to record metric values
// conforming to the "k8s.hpa.metric.target.cpu.value" semantic conventions. It
// represents the target value for CPU resource in HPA config.
//
// Stability: Development
type HPAMetricTargetCPUValue struct {
	metric.Int64Gauge
}

var newHPAMetricTargetCPUValueOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Target value for CPU resource in HPA config."),
	metric.WithUnit("{cpu}"),
}

// NewHPAMetricTargetCPUValue returns a new HPAMetricTargetCPUValue instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newHPAMetricTargetCPUValueOpts...)
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
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUValue) Record(
	ctx context.Context,
	val int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

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
// [K8s HPA MetricTarget]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#metrictarget-v2-autoscaling
// [`ContainerResource`]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/#support-for-metrics-apis
func (m HPAMetricTargetCPUValue) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

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

// HPAMetricTargetCPUValueObservable is an instrument used to record metric
// values conforming to the "k8s.hpa.metric.target.cpu.value" semantic
// conventions. It represents the target value for CPU resource in HPA config.
//
// Stability: Development
type HPAMetricTargetCPUValueObservable struct {
	metric.Int64ObservableGauge
}

var newHPAMetricTargetCPUValueObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Target value for CPU resource in HPA config."),
	metric.WithUnit("{cpu}"),
}

// NewHPAMetricTargetCPUValueObservable returns a new
// HPAMetricTargetCPUValueObservable instrument.
//
// Stability: Development
func NewHPAMetricTargetCPUValueObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (HPAMetricTargetCPUValueObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMetricTargetCPUValueObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAMetricTargetCPUValueObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newHPAMetricTargetCPUValueObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"k8s.hpa.metric.target.cpu.value",
		opt...,
	)
	if err != nil {
		return HPAMetricTargetCPUValueObservable{noop.Int64ObservableGauge{}}, err
	}
	return HPAMetricTargetCPUValueObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAMetricTargetCPUValueObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (HPAMetricTargetCPUValueObservable) Name() string {
	return "k8s.hpa.metric.target.cpu.value"
}

// Unit returns the semantic convention unit of the instrument
func (HPAMetricTargetCPUValueObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (HPAMetricTargetCPUValueObservable) Description() string {
	return "Target value for CPU resource in HPA config."
}

// AttrContainerName returns an optional attribute for the "k8s.container.name"
// semantic convention. It represents the name of the Container from Pod
// specification, must be unique within a Pod. Container runtime usually uses
// different globally unique name (`container.name`).
func (HPAMetricTargetCPUValueObservable) AttrContainerName(val string) attribute.KeyValue {
	return attribute.String("k8s.container.name", val)
}

// AttrHPAMetricType returns an optional attribute for the "k8s.hpa.metric.type"
// semantic convention. It represents the type of metric source for the
// horizontal pod autoscaler.
func (HPAMetricTargetCPUValueObservable) AttrHPAMetricType(val string) attribute.KeyValue {
	return attribute.String("k8s.hpa.metric.type", val)
}

// HPAPodCurrent is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.current" semantic conventions. It represents the current number
// of replica pods managed by this horizontal pod autoscaler, as last seen by the
// autoscaler.
//
// Stability: Development
type HPAPodCurrent struct {
	metric.Int64UpDownCounter
}

var newHPAPodCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodCurrent returns a new HPAPodCurrent instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newHPAPodCurrentOpts...)
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
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAPodCurrentObservable is an instrument used to record metric values
// conforming to the "k8s.hpa.pod.current" semantic conventions. It represents
// the current number of replica pods managed by this horizontal pod autoscaler,
// as last seen by the autoscaler.
//
// Stability: Development
type HPAPodCurrentObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newHPAPodCurrentObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodCurrentObservable returns a new HPAPodCurrentObservable instrument.
//
// Stability: Development
func NewHPAPodCurrentObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (HPAPodCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodCurrentObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newHPAPodCurrentObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.hpa.pod.current",
		opt...,
	)
	if err != nil {
		return HPAPodCurrentObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return HPAPodCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodCurrentObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodCurrentObservable) Name() string {
	return "k8s.hpa.pod.current"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodCurrentObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodCurrentObservable) Description() string {
	return "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler."
}

// HPAPodDesired is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.desired" semantic conventions. It represents the desired number
// of replica pods managed by this horizontal pod autoscaler, as last calculated
// by the autoscaler.
//
// Stability: Development
type HPAPodDesired struct {
	metric.Int64UpDownCounter
}

var newHPAPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodDesired returns a new HPAPodDesired instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newHPAPodDesiredOpts...)
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
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `desiredReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPAPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAPodDesiredObservable is an instrument used to record metric values
// conforming to the "k8s.hpa.pod.desired" semantic conventions. It represents
// the desired number of replica pods managed by this horizontal pod autoscaler,
// as last calculated by the autoscaler.
//
// Stability: Development
type HPAPodDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newHPAPodDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodDesiredObservable returns a new HPAPodDesiredObservable instrument.
//
// Stability: Development
func NewHPAPodDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (HPAPodDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newHPAPodDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.hpa.pod.desired",
		opt...,
	)
	if err != nil {
		return HPAPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return HPAPodDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodDesiredObservable) Name() string {
	return "k8s.hpa.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodDesiredObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodDesiredObservable) Description() string {
	return "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler."
}

// HPAPodMax is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.max" semantic conventions. It represents the upper limit for the
// number of replica pods to which the autoscaler can scale up.
//
// Stability: Development
type HPAPodMax struct {
	metric.Int64UpDownCounter
}

var newHPAPodMaxOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The upper limit for the number of replica pods to which the autoscaler can scale up."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodMax returns a new HPAPodMax instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newHPAPodMaxOpts...)
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
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMax) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `maxReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMax) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAPodMaxObservable is an instrument used to record metric values conforming
// to the "k8s.hpa.pod.max" semantic conventions. It represents the upper limit
// for the number of replica pods to which the autoscaler can scale up.
//
// Stability: Development
type HPAPodMaxObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newHPAPodMaxObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The upper limit for the number of replica pods to which the autoscaler can scale up."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodMaxObservable returns a new HPAPodMaxObservable instrument.
//
// Stability: Development
func NewHPAPodMaxObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (HPAPodMaxObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodMaxObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodMaxObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newHPAPodMaxObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.hpa.pod.max",
		opt...,
	)
	if err != nil {
		return HPAPodMaxObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return HPAPodMaxObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodMaxObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodMaxObservable) Name() string {
	return "k8s.hpa.pod.max"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodMaxObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodMaxObservable) Description() string {
	return "The upper limit for the number of replica pods to which the autoscaler can scale up."
}

// HPAPodMin is an instrument used to record metric values conforming to the
// "k8s.hpa.pod.min" semantic conventions. It represents the lower limit for the
// number of replica pods to which the autoscaler can scale down.
//
// Stability: Development
type HPAPodMin struct {
	metric.Int64UpDownCounter
}

var newHPAPodMinOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The lower limit for the number of replica pods to which the autoscaler can scale down."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodMin returns a new HPAPodMin instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newHPAPodMinOpts...)
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
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMin) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `minReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAPodMin) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// HPAPodMinObservable is an instrument used to record metric values conforming
// to the "k8s.hpa.pod.min" semantic conventions. It represents the lower limit
// for the number of replica pods to which the autoscaler can scale down.
//
// Stability: Development
type HPAPodMinObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newHPAPodMinObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The lower limit for the number of replica pods to which the autoscaler can scale down."),
	metric.WithUnit("{pod}"),
}

// NewHPAPodMinObservable returns a new HPAPodMinObservable instrument.
//
// Stability: Development
func NewHPAPodMinObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (HPAPodMinObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAPodMinObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newHPAPodMinObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newHPAPodMinObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.hpa.pod.min",
		opt...,
	)
	if err != nil {
		return HPAPodMinObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return HPAPodMinObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HPAPodMinObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (HPAPodMinObservable) Name() string {
	return "k8s.hpa.pod.min"
}

// Unit returns the semantic convention unit of the instrument
func (HPAPodMinObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (HPAPodMinObservable) Description() string {
	return "The lower limit for the number of replica pods to which the autoscaler can scale down."
}

// JobPodActive is an instrument used to record metric values conforming to the
// "k8s.job.pod.active" semantic conventions. It represents the number of pending
// and actively running pods for a job.
//
// Stability: Development
type JobPodActive struct {
	metric.Int64UpDownCounter
}

var newJobPodActiveOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of pending and actively running pods for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodActive returns a new JobPodActive instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newJobPodActiveOpts...)
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
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobstatus-v1-batch
func (m JobPodActive) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `active` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobstatus-v1-batch
func (m JobPodActive) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodActiveObservable is an instrument used to record metric values
// conforming to the "k8s.job.pod.active" semantic conventions. It represents the
// number of pending and actively running pods for a job.
//
// Stability: Development
type JobPodActiveObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newJobPodActiveObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of pending and actively running pods for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodActiveObservable returns a new JobPodActiveObservable instrument.
//
// Stability: Development
func NewJobPodActiveObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (JobPodActiveObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodActiveObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodActiveObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newJobPodActiveObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.job.pod.active",
		opt...,
	)
	if err != nil {
		return JobPodActiveObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return JobPodActiveObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodActiveObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodActiveObservable) Name() string {
	return "k8s.job.pod.active"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodActiveObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodActiveObservable) Description() string {
	return "The number of pending and actively running pods for a job."
}

// JobPodDesiredSuccessful is an instrument used to record metric values
// conforming to the "k8s.job.pod.desired_successful" semantic conventions. It
// represents the desired number of successfully finished pods the job should be
// run with.
//
// Stability: Development
type JobPodDesiredSuccessful struct {
	metric.Int64UpDownCounter
}

var newJobPodDesiredSuccessfulOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The desired number of successfully finished pods the job should be run with."),
	metric.WithUnit("{pod}"),
}

// NewJobPodDesiredSuccessful returns a new JobPodDesiredSuccessful instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newJobPodDesiredSuccessfulOpts...)
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
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobspec-v1-batch
func (m JobPodDesiredSuccessful) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `completions` field of the
// [K8s JobSpec]..
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobspec-v1-batch
func (m JobPodDesiredSuccessful) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodDesiredSuccessfulObservable is an instrument used to record metric
// values conforming to the "k8s.job.pod.desired_successful" semantic
// conventions. It represents the desired number of successfully finished pods
// the job should be run with.
//
// Stability: Development
type JobPodDesiredSuccessfulObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newJobPodDesiredSuccessfulObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The desired number of successfully finished pods the job should be run with."),
	metric.WithUnit("{pod}"),
}

// NewJobPodDesiredSuccessfulObservable returns a new
// JobPodDesiredSuccessfulObservable instrument.
//
// Stability: Development
func NewJobPodDesiredSuccessfulObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (JobPodDesiredSuccessfulObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodDesiredSuccessfulObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodDesiredSuccessfulObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newJobPodDesiredSuccessfulObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.job.pod.desired_successful",
		opt...,
	)
	if err != nil {
		return JobPodDesiredSuccessfulObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return JobPodDesiredSuccessfulObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodDesiredSuccessfulObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodDesiredSuccessfulObservable) Name() string {
	return "k8s.job.pod.desired_successful"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodDesiredSuccessfulObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodDesiredSuccessfulObservable) Description() string {
	return "The desired number of successfully finished pods the job should be run with."
}

// JobPodFailed is an instrument used to record metric values conforming to the
// "k8s.job.pod.failed" semantic conventions. It represents the number of pods
// which reached phase Failed for a job.
//
// Stability: Development
type JobPodFailed struct {
	metric.Int64UpDownCounter
}

var newJobPodFailedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of pods which reached phase Failed for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodFailed returns a new JobPodFailed instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newJobPodFailedOpts...)
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
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobstatus-v1-batch
func (m JobPodFailed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `failed` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobstatus-v1-batch
func (m JobPodFailed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodFailedObservable is an instrument used to record metric values
// conforming to the "k8s.job.pod.failed" semantic conventions. It represents the
// number of pods which reached phase Failed for a job.
//
// Stability: Development
type JobPodFailedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newJobPodFailedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of pods which reached phase Failed for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodFailedObservable returns a new JobPodFailedObservable instrument.
//
// Stability: Development
func NewJobPodFailedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (JobPodFailedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodFailedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodFailedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newJobPodFailedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.job.pod.failed",
		opt...,
	)
	if err != nil {
		return JobPodFailedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return JobPodFailedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodFailedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodFailedObservable) Name() string {
	return "k8s.job.pod.failed"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodFailedObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodFailedObservable) Description() string {
	return "The number of pods which reached phase Failed for a job."
}

// JobPodMaxParallel is an instrument used to record metric values conforming to
// the "k8s.job.pod.max_parallel" semantic conventions. It represents the max
// desired number of pods the job should run at any given time.
//
// Stability: Development
type JobPodMaxParallel struct {
	metric.Int64UpDownCounter
}

var newJobPodMaxParallelOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The max desired number of pods the job should run at any given time."),
	metric.WithUnit("{pod}"),
}

// NewJobPodMaxParallel returns a new JobPodMaxParallel instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newJobPodMaxParallelOpts...)
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
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobspec-v1-batch
func (m JobPodMaxParallel) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `parallelism` field of the
// [K8s JobSpec].
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobspec-v1-batch
func (m JobPodMaxParallel) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodMaxParallelObservable is an instrument used to record metric values
// conforming to the "k8s.job.pod.max_parallel" semantic conventions. It
// represents the max desired number of pods the job should run at any given
// time.
//
// Stability: Development
type JobPodMaxParallelObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newJobPodMaxParallelObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The max desired number of pods the job should run at any given time."),
	metric.WithUnit("{pod}"),
}

// NewJobPodMaxParallelObservable returns a new JobPodMaxParallelObservable
// instrument.
//
// Stability: Development
func NewJobPodMaxParallelObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (JobPodMaxParallelObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodMaxParallelObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodMaxParallelObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newJobPodMaxParallelObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.job.pod.max_parallel",
		opt...,
	)
	if err != nil {
		return JobPodMaxParallelObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return JobPodMaxParallelObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodMaxParallelObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodMaxParallelObservable) Name() string {
	return "k8s.job.pod.max_parallel"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodMaxParallelObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodMaxParallelObservable) Description() string {
	return "The max desired number of pods the job should run at any given time."
}

// JobPodSuccessful is an instrument used to record metric values conforming to
// the "k8s.job.pod.successful" semantic conventions. It represents the number of
// pods which reached phase Succeeded for a job.
//
// Stability: Development
type JobPodSuccessful struct {
	metric.Int64UpDownCounter
}

var newJobPodSuccessfulOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of pods which reached phase Succeeded for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodSuccessful returns a new JobPodSuccessful instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newJobPodSuccessfulOpts...)
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
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobstatus-v1-batch
func (m JobPodSuccessful) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `succeeded` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#jobstatus-v1-batch
func (m JobPodSuccessful) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// JobPodSuccessfulObservable is an instrument used to record metric values
// conforming to the "k8s.job.pod.successful" semantic conventions. It represents
// the number of pods which reached phase Succeeded for a job.
//
// Stability: Development
type JobPodSuccessfulObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newJobPodSuccessfulObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of pods which reached phase Succeeded for a job."),
	metric.WithUnit("{pod}"),
}

// NewJobPodSuccessfulObservable returns a new JobPodSuccessfulObservable
// instrument.
//
// Stability: Development
func NewJobPodSuccessfulObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (JobPodSuccessfulObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return JobPodSuccessfulObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJobPodSuccessfulObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newJobPodSuccessfulObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.job.pod.successful",
		opt...,
	)
	if err != nil {
		return JobPodSuccessfulObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return JobPodSuccessfulObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m JobPodSuccessfulObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (JobPodSuccessfulObservable) Name() string {
	return "k8s.job.pod.successful"
}

// Unit returns the semantic convention unit of the instrument
func (JobPodSuccessfulObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (JobPodSuccessfulObservable) Description() string {
	return "The number of pods which reached phase Succeeded for a job."
}

// NamespacePhase is an instrument used to record metric values conforming to the
// "k8s.namespace.phase" semantic conventions. It represents the describes number
// of K8s namespaces that are currently in a given phase.
//
// Stability: Development
type NamespacePhase struct {
	metric.Int64UpDownCounter
}

var newNamespacePhaseOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes number of K8s namespaces that are currently in a given phase."),
	metric.WithUnit("{namespace}"),
}

// NewNamespacePhase returns a new NamespacePhase instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newNamespacePhaseOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.namespace.phase", string(namespacePhase)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("k8s.namespace.phase", string(namespacePhase)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NamespacePhase) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NamespacePhaseObservable is an instrument used to record metric values
// conforming to the "k8s.namespace.phase" semantic conventions. It represents
// the describes number of K8s namespaces that are currently in a given phase.
//
// Stability: Development
type NamespacePhaseObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNamespacePhaseObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Describes number of K8s namespaces that are currently in a given phase."),
	metric.WithUnit("{namespace}"),
}

// NewNamespacePhaseObservable returns a new NamespacePhaseObservable instrument.
//
// Stability: Development
func NewNamespacePhaseObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NamespacePhaseObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NamespacePhaseObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNamespacePhaseObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNamespacePhaseObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.namespace.phase",
		opt...,
	)
	if err != nil {
		return NamespacePhaseObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NamespacePhaseObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NamespacePhaseObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NamespacePhaseObservable) Name() string {
	return "k8s.namespace.phase"
}

// Unit returns the semantic convention unit of the instrument
func (NamespacePhaseObservable) Unit() string {
	return "{namespace}"
}

// Description returns the semantic convention description of the instrument
func (NamespacePhaseObservable) Description() string {
	return "Describes number of K8s namespaces that are currently in a given phase."
}

// AttrNamespacePhase returns a required attribute for the "k8s.namespace.phase"
// semantic convention. It represents the phase of the K8s namespace.
func (NamespacePhaseObservable) AttrNamespacePhase(val NamespacePhaseAttr) attribute.KeyValue {
	return attribute.String("k8s.namespace.phase", string(val))
}

// NodeConditionStatus is an instrument used to record metric values conforming
// to the "k8s.node.condition.status" semantic conventions. It represents the
// describes the condition of a particular Node.
//
// Stability: Development
type NodeConditionStatus struct {
	metric.Int64UpDownCounter
}

var newNodeConditionStatusOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the condition of a particular Node."),
	metric.WithUnit("{node}"),
}

// NewNodeConditionStatus returns a new NodeConditionStatus instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newNodeConditionStatusOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.node.condition.status", string(nodeConditionStatus)),
			attribute.String("k8s.node.condition.type", string(nodeConditionType)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeConditionStatusObservable is an instrument used to record metric values
// conforming to the "k8s.node.condition.status" semantic conventions. It
// represents the describes the condition of a particular Node.
//
// Stability: Development
type NodeConditionStatusObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeConditionStatusObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Describes the condition of a particular Node."),
	metric.WithUnit("{node}"),
}

// NewNodeConditionStatusObservable returns a new NodeConditionStatusObservable
// instrument.
//
// Stability: Development
func NewNodeConditionStatusObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeConditionStatusObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeConditionStatusObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeConditionStatusObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeConditionStatusObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.condition.status",
		opt...,
	)
	if err != nil {
		return NodeConditionStatusObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeConditionStatusObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeConditionStatusObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeConditionStatusObservable) Name() string {
	return "k8s.node.condition.status"
}

// Unit returns the semantic convention unit of the instrument
func (NodeConditionStatusObservable) Unit() string {
	return "{node}"
}

// Description returns the semantic convention description of the instrument
func (NodeConditionStatusObservable) Description() string {
	return "Describes the condition of a particular Node."
}

// AttrNodeConditionStatus returns a required attribute for the
// "k8s.node.condition.status" semantic convention. It represents the status of
// the condition, one of True, False, Unknown.
func (NodeConditionStatusObservable) AttrNodeConditionStatus(val NodeConditionStatusAttr) attribute.KeyValue {
	return attribute.String("k8s.node.condition.status", string(val))
}

// AttrNodeConditionType returns a required attribute for the
// "k8s.node.condition.type" semantic convention. It represents the condition
// type of a K8s Node.
func (NodeConditionStatusObservable) AttrNodeConditionType(val NodeConditionTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.node.condition.type", string(val))
}

// NodeCPUAllocatable is an instrument used to record metric values conforming to
// the "k8s.node.cpu.allocatable" semantic conventions. It represents the amount
// of CPU allocatable on the node.
//
// Stability: Development
type NodeCPUAllocatable struct {
	metric.Float64UpDownCounter
}

var newNodeCPUAllocatableOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("Amount of CPU allocatable on the node."),
	metric.WithUnit("{cpu}"),
}

// NewNodeCPUAllocatable returns a new NodeCPUAllocatable instrument.
//
// Stability: Development
func NewNodeCPUAllocatable(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (NodeCPUAllocatable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUAllocatable{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeCPUAllocatableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeCPUAllocatableOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.node.cpu.allocatable",
		opt...,
	)
	if err != nil {
		return NodeCPUAllocatable{noop.Float64UpDownCounter{}}, err
	}
	return NodeCPUAllocatable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUAllocatable) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
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
	return "Amount of CPU allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeCPUAllocatable) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NodeCPUAllocatable) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeCPUAllocatableObservable is an instrument used to record metric values
// conforming to the "k8s.node.cpu.allocatable" semantic conventions. It
// represents the amount of CPU allocatable on the node.
//
// Stability: Development
type NodeCPUAllocatableObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newNodeCPUAllocatableObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("Amount of CPU allocatable on the node."),
	metric.WithUnit("{cpu}"),
}

// NewNodeCPUAllocatableObservable returns a new NodeCPUAllocatableObservable
// instrument.
//
// Stability: Development
func NewNodeCPUAllocatableObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (NodeCPUAllocatableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUAllocatableObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeCPUAllocatableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeCPUAllocatableObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.node.cpu.allocatable",
		opt...,
	)
	if err != nil {
		return NodeCPUAllocatableObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return NodeCPUAllocatableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUAllocatableObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeCPUAllocatableObservable) Name() string {
	return "k8s.node.cpu.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodeCPUAllocatableObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (NodeCPUAllocatableObservable) Description() string {
	return "Amount of CPU allocatable on the node."
}

// NodeCPUTime is an instrument used to record metric values conforming to the
// "k8s.node.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
//
// Stability: Release_Candidate
type NodeCPUTime struct {
	metric.Float64Counter
}

var newNodeCPUTimeOpts = []metric.Float64CounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
}

// NewNodeCPUTime returns a new NodeCPUTime instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeCPUTimeOpts...)
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
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Total CPU time consumed by the specific Node on all available CPU cores
func (m NodeCPUTime) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// NodeCPUTimeObservable is an instrument used to record metric values conforming
// to the "k8s.node.cpu.time" semantic conventions. It represents the total CPU
// time consumed.
//
// Stability: Release_Candidate
type NodeCPUTimeObservable struct {
	metric.Float64ObservableCounter
}

var newNodeCPUTimeObservableOpts = []metric.Float64ObservableCounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
}

// NewNodeCPUTimeObservable returns a new NodeCPUTimeObservable instrument.
//
// Stability: Release_Candidate
func NewNodeCPUTimeObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableCounterOption,
) (NodeCPUTimeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUTimeObservable{noop.Float64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeCPUTimeObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeCPUTimeObservableOpts...)
	}

	i, err := m.Float64ObservableCounter(
		"k8s.node.cpu.time",
		opt...,
	)
	if err != nil {
		return NodeCPUTimeObservable{noop.Float64ObservableCounter{}}, err
	}
	return NodeCPUTimeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUTimeObservable) Inst() metric.Float64ObservableCounter {
	return m.Float64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeCPUTimeObservable) Name() string {
	return "k8s.node.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (NodeCPUTimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (NodeCPUTimeObservable) Description() string {
	return "Total CPU time consumed."
}

// NodeCPUUsage is an instrument used to record metric values conforming to the
// "k8s.node.cpu.usage" semantic conventions. It represents the node's CPU usage,
// measured in CPUs. Range from 0 to the number of allocatable CPUs.
//
// Stability: Release_Candidate
type NodeCPUUsage struct {
	metric.Float64Gauge
}

var newNodeCPUUsageOpts = []metric.Float64GaugeOption{
	metric.WithDescription("Node's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
}

// NewNodeCPUUsage returns a new NodeCPUUsage instrument.
//
// Stability: Release_Candidate
func NewNodeCPUUsage(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (NodeCPUUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUUsage{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeCPUUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeCPUUsageOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.node.cpu.usage",
		opt...,
	)
	if err != nil {
		return NodeCPUUsage{noop.Float64Gauge{}}, err
	}
	return NodeCPUUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUUsage) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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
	return "Node's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."
}

// Record records val to the current distribution for attrs.
//
// CPU usage of the specific Node on all available CPU cores. It is calculated as
// the change in cumulative CPU time (k8s.node.cpu.time) over a measurement
// interval, divided by the elapsed time: usageCores = (cpuTimeEnd -
// cpuTimeStart) / elapsedSeconds
func (m NodeCPUUsage) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// CPU usage of the specific Node on all available CPU cores. It is calculated as
// the change in cumulative CPU time (k8s.node.cpu.time) over a measurement
// interval, divided by the elapsed time: usageCores = (cpuTimeEnd -
// cpuTimeStart) / elapsedSeconds
func (m NodeCPUUsage) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// NodeCPUUsageObservable is an instrument used to record metric values
// conforming to the "k8s.node.cpu.usage" semantic conventions. It represents the
// node's CPU usage, measured in CPUs. Range from 0 to the number of allocatable
// CPUs.
//
// Stability: Release_Candidate
type NodeCPUUsageObservable struct {
	metric.Float64ObservableGauge
}

var newNodeCPUUsageObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("Node's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
}

// NewNodeCPUUsageObservable returns a new NodeCPUUsageObservable instrument.
//
// Stability: Release_Candidate
func NewNodeCPUUsageObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (NodeCPUUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeCPUUsageObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeCPUUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeCPUUsageObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.node.cpu.usage",
		opt...,
	)
	if err != nil {
		return NodeCPUUsageObservable{noop.Float64ObservableGauge{}}, err
	}
	return NodeCPUUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeCPUUsageObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (NodeCPUUsageObservable) Name() string {
	return "k8s.node.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeCPUUsageObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (NodeCPUUsageObservable) Description() string {
	return "Node's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."
}

// NodeEphemeralStorageAllocatable is an instrument used to record metric values
// conforming to the "k8s.node.ephemeral_storage.allocatable" semantic
// conventions. It represents the amount of ephemeral-storage allocatable on the
// node.
//
// Stability: Development
type NodeEphemeralStorageAllocatable struct {
	metric.Int64UpDownCounter
}

var newNodeEphemeralStorageAllocatableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Amount of ephemeral-storage allocatable on the node."),
	metric.WithUnit("By"),
}

// NewNodeEphemeralStorageAllocatable returns a new
// NodeEphemeralStorageAllocatable instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newNodeEphemeralStorageAllocatableOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NodeEphemeralStorageAllocatable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeEphemeralStorageAllocatableObservable is an instrument used to record
// metric values conforming to the "k8s.node.ephemeral_storage.allocatable"
// semantic conventions. It represents the amount of ephemeral-storage
// allocatable on the node.
//
// Stability: Development
type NodeEphemeralStorageAllocatableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeEphemeralStorageAllocatableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Amount of ephemeral-storage allocatable on the node."),
	metric.WithUnit("By"),
}

// NewNodeEphemeralStorageAllocatableObservable returns a new
// NodeEphemeralStorageAllocatableObservable instrument.
//
// Stability: Development
func NewNodeEphemeralStorageAllocatableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeEphemeralStorageAllocatableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeEphemeralStorageAllocatableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeEphemeralStorageAllocatableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeEphemeralStorageAllocatableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.ephemeral_storage.allocatable",
		opt...,
	)
	if err != nil {
		return NodeEphemeralStorageAllocatableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeEphemeralStorageAllocatableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeEphemeralStorageAllocatableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeEphemeralStorageAllocatableObservable) Name() string {
	return "k8s.node.ephemeral_storage.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodeEphemeralStorageAllocatableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeEphemeralStorageAllocatableObservable) Description() string {
	return "Amount of ephemeral-storage allocatable on the node."
}

// NodeFilesystemAvailable is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.available" semantic conventions. It
// represents the node filesystem available bytes.
//
// Stability: Release_Candidate
type NodeFilesystemAvailable struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemAvailable returns a new NodeFilesystemAvailable instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemAvailableOpts...)
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
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemAvailableObservable is an instrument used to record metric
// values conforming to the "k8s.node.filesystem.available" semantic conventions.
// It represents the node filesystem available bytes.
//
// Stability: Release_Candidate
type NodeFilesystemAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeFilesystemAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Node filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemAvailableObservable returns a new
// NodeFilesystemAvailableObservable instrument.
//
// Stability: Release_Candidate
func NewNodeFilesystemAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeFilesystemAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.filesystem.available",
		opt...,
	)
	if err != nil {
		return NodeFilesystemAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeFilesystemAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemAvailableObservable) Name() string {
	return "k8s.node.filesystem.available"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemAvailableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemAvailableObservable) Description() string {
	return "Node filesystem available bytes."
}

// NodeFilesystemCapacity is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.capacity" semantic conventions. It
// represents the node filesystem capacity.
//
// Stability: Release_Candidate
type NodeFilesystemCapacity struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node filesystem capacity."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemCapacity returns a new NodeFilesystemCapacity instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemCapacityOpts...)
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
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemCapacity) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemCapacityObservable is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.capacity" semantic conventions. It
// represents the node filesystem capacity.
//
// Stability: Release_Candidate
type NodeFilesystemCapacityObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeFilesystemCapacityObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Node filesystem capacity."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemCapacityObservable returns a new
// NodeFilesystemCapacityObservable instrument.
//
// Stability: Release_Candidate
func NewNodeFilesystemCapacityObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeFilesystemCapacityObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemCapacityObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemCapacityObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemCapacityObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.filesystem.capacity",
		opt...,
	)
	if err != nil {
		return NodeFilesystemCapacityObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeFilesystemCapacityObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemCapacityObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemCapacityObservable) Name() string {
	return "k8s.node.filesystem.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemCapacityObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemCapacityObservable) Description() string {
	return "Node filesystem capacity."
}

// NodeFilesystemInodeCount is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.inode.count" semantic conventions. It
// represents the total inodes in the node's root filesystem.
//
// Stability: Release_Candidate
type NodeFilesystemInodeCount struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemInodeCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The total inodes in the node's root filesystem."),
	metric.WithUnit("{inode}"),
}

// NewNodeFilesystemInodeCount returns a new NodeFilesystemInodeCount instrument.
//
// Stability: Release_Candidate
func NewNodeFilesystemInodeCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeFilesystemInodeCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemInodeCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemInodeCountOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemInodeCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.inode.count",
		opt...,
	)
	if err != nil {
		return NodeFilesystemInodeCount{noop.Int64UpDownCounter{}}, err
	}
	return NodeFilesystemInodeCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemInodeCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemInodeCount) Name() string {
	return "k8s.node.filesystem.inode.count"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemInodeCount) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemInodeCount) Description() string {
	return "The total inodes in the node's root filesystem."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the
// [FsStats.Inodes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.Inodes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemInodeCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.Inodes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.Inodes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemInodeCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemInodeCountObservable is an instrument used to record metric
// values conforming to the "k8s.node.filesystem.inode.count" semantic
// conventions. It represents the total inodes in the node's root filesystem.
//
// Stability: Release_Candidate
type NodeFilesystemInodeCountObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeFilesystemInodeCountObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The total inodes in the node's root filesystem."),
	metric.WithUnit("{inode}"),
}

// NewNodeFilesystemInodeCountObservable returns a new
// NodeFilesystemInodeCountObservable instrument.
//
// Stability: Release_Candidate
func NewNodeFilesystemInodeCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeFilesystemInodeCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemInodeCountObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemInodeCountObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemInodeCountObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.filesystem.inode.count",
		opt...,
	)
	if err != nil {
		return NodeFilesystemInodeCountObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeFilesystemInodeCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemInodeCountObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemInodeCountObservable) Name() string {
	return "k8s.node.filesystem.inode.count"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemInodeCountObservable) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemInodeCountObservable) Description() string {
	return "The total inodes in the node's root filesystem."
}

// NodeFilesystemInodeFree is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.inode.free" semantic conventions. It
// represents the free inodes in the node's root filesystem.
//
// Stability: Release_Candidate
type NodeFilesystemInodeFree struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemInodeFreeOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The free inodes in the node's root filesystem."),
	metric.WithUnit("{inode}"),
}

// NewNodeFilesystemInodeFree returns a new NodeFilesystemInodeFree instrument.
//
// Stability: Release_Candidate
func NewNodeFilesystemInodeFree(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeFilesystemInodeFree, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemInodeFree{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemInodeFreeOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemInodeFreeOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.inode.free",
		opt...,
	)
	if err != nil {
		return NodeFilesystemInodeFree{noop.Int64UpDownCounter{}}, err
	}
	return NodeFilesystemInodeFree{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemInodeFree) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemInodeFree) Name() string {
	return "k8s.node.filesystem.inode.free"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemInodeFree) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemInodeFree) Description() string {
	return "The free inodes in the node's root filesystem."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the
// [FsStats.InodesFree] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.InodesFree]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemInodeFree) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.InodesFree] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// [FsStats.InodesFree]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemInodeFree) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemInodeFreeObservable is an instrument used to record metric
// values conforming to the "k8s.node.filesystem.inode.free" semantic
// conventions. It represents the free inodes in the node's root filesystem.
//
// Stability: Release_Candidate
type NodeFilesystemInodeFreeObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeFilesystemInodeFreeObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The free inodes in the node's root filesystem."),
	metric.WithUnit("{inode}"),
}

// NewNodeFilesystemInodeFreeObservable returns a new
// NodeFilesystemInodeFreeObservable instrument.
//
// Stability: Release_Candidate
func NewNodeFilesystemInodeFreeObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeFilesystemInodeFreeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemInodeFreeObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemInodeFreeObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemInodeFreeObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.filesystem.inode.free",
		opt...,
	)
	if err != nil {
		return NodeFilesystemInodeFreeObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeFilesystemInodeFreeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemInodeFreeObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemInodeFreeObservable) Name() string {
	return "k8s.node.filesystem.inode.free"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemInodeFreeObservable) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemInodeFreeObservable) Description() string {
	return "The free inodes in the node's root filesystem."
}

// NodeFilesystemUsage is an instrument used to record metric values conforming
// to the "k8s.node.filesystem.usage" semantic conventions. It represents the
// node filesystem usage.
//
// Stability: Release_Candidate
type NodeFilesystemUsage struct {
	metric.Int64UpDownCounter
}

var newNodeFilesystemUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node filesystem usage."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemUsage returns a new NodeFilesystemUsage instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemUsageOpts...)
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
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// Capacity and available are reported for the backing filesystem,
// which may be shared with other workloads; therefore,
// usage may not equal `capacity - available`.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [NodeStats.Fs]
// of the Kubelet's stats API.
//
// Capacity and available are reported for the backing filesystem,
// which may be shared with other workloads; therefore,
// usage may not equal `capacity - available`.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [NodeStats.Fs]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeFilesystemUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeFilesystemUsageObservable is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.usage" semantic conventions. It
// represents the node filesystem usage.
//
// Stability: Release_Candidate
type NodeFilesystemUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeFilesystemUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Node filesystem usage."),
	metric.WithUnit("By"),
}

// NewNodeFilesystemUsageObservable returns a new NodeFilesystemUsageObservable
// instrument.
//
// Stability: Release_Candidate
func NewNodeFilesystemUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeFilesystemUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeFilesystemUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeFilesystemUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.filesystem.usage",
		opt...,
	)
	if err != nil {
		return NodeFilesystemUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeFilesystemUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeFilesystemUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeFilesystemUsageObservable) Name() string {
	return "k8s.node.filesystem.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeFilesystemUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeFilesystemUsageObservable) Description() string {
	return "Node filesystem usage."
}

// NodeMemoryAllocatable is an instrument used to record metric values conforming
// to the "k8s.node.memory.allocatable" semantic conventions. It represents the
// amount of memory allocatable on the node.
//
// Stability: Development
type NodeMemoryAllocatable struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryAllocatableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Amount of memory allocatable on the node."),
	metric.WithUnit("By"),
}

// NewNodeMemoryAllocatable returns a new NodeMemoryAllocatable instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryAllocatableOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NodeMemoryAllocatable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryAllocatableObservable is an instrument used to record metric values
// conforming to the "k8s.node.memory.allocatable" semantic conventions. It
// represents the amount of memory allocatable on the node.
//
// Stability: Development
type NodeMemoryAllocatableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeMemoryAllocatableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Amount of memory allocatable on the node."),
	metric.WithUnit("By"),
}

// NewNodeMemoryAllocatableObservable returns a new
// NodeMemoryAllocatableObservable instrument.
//
// Stability: Development
func NewNodeMemoryAllocatableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeMemoryAllocatableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryAllocatableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryAllocatableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryAllocatableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.memory.allocatable",
		opt...,
	)
	if err != nil {
		return NodeMemoryAllocatableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeMemoryAllocatableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryAllocatableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryAllocatableObservable) Name() string {
	return "k8s.node.memory.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryAllocatableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryAllocatableObservable) Description() string {
	return "Amount of memory allocatable on the node."
}

// NodeMemoryAvailable is an instrument used to record metric values conforming
// to the "k8s.node.memory.available" semantic conventions. It represents the
// node memory available.
//
// Stability: Release_Candidate
type NodeMemoryAvailable struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node memory available."),
	metric.WithUnit("By"),
}

// NewNodeMemoryAvailable returns a new NodeMemoryAvailable instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryAvailableOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryAvailableObservable is an instrument used to record metric values
// conforming to the "k8s.node.memory.available" semantic conventions. It
// represents the node memory available.
//
// Stability: Release_Candidate
type NodeMemoryAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeMemoryAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Node memory available."),
	metric.WithUnit("By"),
}

// NewNodeMemoryAvailableObservable returns a new NodeMemoryAvailableObservable
// instrument.
//
// Stability: Release_Candidate
func NewNodeMemoryAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeMemoryAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.memory.available",
		opt...,
	)
	if err != nil {
		return NodeMemoryAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeMemoryAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryAvailableObservable) Name() string {
	return "k8s.node.memory.available"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryAvailableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryAvailableObservable) Description() string {
	return "Node memory available."
}

// NodeMemoryRss is an instrument used to record metric values conforming to the
// "k8s.node.memory.rss" semantic conventions. It represents the node memory RSS.
//
// Stability: Release_Candidate
type NodeMemoryRss struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryRssOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node memory RSS."),
	metric.WithUnit("By"),
}

// NewNodeMemoryRss returns a new NodeMemoryRss instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryRssOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryRssObservable is an instrument used to record metric values
// conforming to the "k8s.node.memory.rss" semantic conventions. It represents
// the node memory RSS.
//
// Stability: Release_Candidate
type NodeMemoryRssObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeMemoryRssObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Node memory RSS."),
	metric.WithUnit("By"),
}

// NewNodeMemoryRssObservable returns a new NodeMemoryRssObservable instrument.
//
// Stability: Release_Candidate
func NewNodeMemoryRssObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeMemoryRssObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryRssObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryRssObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryRssObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.memory.rss",
		opt...,
	)
	if err != nil {
		return NodeMemoryRssObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeMemoryRssObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryRssObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryRssObservable) Name() string {
	return "k8s.node.memory.rss"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryRssObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryRssObservable) Description() string {
	return "Node memory RSS."
}

// NodeMemoryUsage is an instrument used to record metric values conforming to
// the "k8s.node.memory.usage" semantic conventions. It represents the memory
// usage of the Node.
//
// Stability: Release_Candidate
type NodeMemoryUsage struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Memory usage of the Node."),
	metric.WithUnit("By"),
}

// NewNodeMemoryUsage returns a new NodeMemoryUsage instrument.
//
// Stability: Release_Candidate
func NewNodeMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeMemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.memory.usage",
		opt...,
	)
	if err != nil {
		return NodeMemoryUsage{noop.Int64UpDownCounter{}}, err
	}
	return NodeMemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count for attrs.
//
// Total memory in use. This includes all memory regardless of when it was
// accessed.
// This metric is derived from the [MemoryStats.UsageBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.UsageBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Total memory in use. This includes all memory regardless of when it was
// accessed.
// This metric is derived from the [MemoryStats.UsageBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.UsageBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryUsageObservable is an instrument used to record metric values
// conforming to the "k8s.node.memory.usage" semantic conventions. It represents
// the memory usage of the Node.
//
// Stability: Release_Candidate
type NodeMemoryUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeMemoryUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Memory usage of the Node."),
	metric.WithUnit("By"),
}

// NewNodeMemoryUsageObservable returns a new NodeMemoryUsageObservable
// instrument.
//
// Stability: Release_Candidate
func NewNodeMemoryUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeMemoryUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.memory.usage",
		opt...,
	)
	if err != nil {
		return NodeMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeMemoryUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryUsageObservable) Name() string {
	return "k8s.node.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryUsageObservable) Description() string {
	return "Memory usage of the Node."
}

// NodeMemoryWorkingSet is an instrument used to record metric values conforming
// to the "k8s.node.memory.working_set" semantic conventions. It represents the
// node memory working set.
//
// Stability: Release_Candidate
type NodeMemoryWorkingSet struct {
	metric.Int64UpDownCounter
}

var newNodeMemoryWorkingSetOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node memory working set."),
	metric.WithUnit("By"),
}

// NewNodeMemoryWorkingSet returns a new NodeMemoryWorkingSet instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryWorkingSetOpts...)
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
// dirty memory, and kernel memory. On Linux, the working set is calculated by
// excluding inactive file-backed memory from memory usage (`inactive_file` on
// cgroup v2, `total_inactive_file` on cgroup v1), so active file-backed memory
// remains included. On Linux, the kubelet uses the node working set when
// determining memory availability for the `memory.available` eviction signal.
// `WorkingSetBytes` is less than or equal to `UsageBytes`.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryWorkingSet) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// The amount of working set memory. This includes recently accessed memory,
// dirty memory, and kernel memory. On Linux, the working set is calculated by
// excluding inactive file-backed memory from memory usage (`inactive_file` on
// cgroup v2, `total_inactive_file` on cgroup v1), so active file-backed memory
// remains included. On Linux, the kubelet uses the node working set when
// determining memory availability for the `memory.available` eviction signal.
// `WorkingSetBytes` is less than or equal to `UsageBytes`.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [NodeStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodeMemoryWorkingSet) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeMemoryWorkingSetObservable is an instrument used to record metric values
// conforming to the "k8s.node.memory.working_set" semantic conventions. It
// represents the node memory working set.
//
// Stability: Release_Candidate
type NodeMemoryWorkingSetObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeMemoryWorkingSetObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Node memory working set."),
	metric.WithUnit("By"),
}

// NewNodeMemoryWorkingSetObservable returns a new NodeMemoryWorkingSetObservable
// instrument.
//
// Stability: Release_Candidate
func NewNodeMemoryWorkingSetObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeMemoryWorkingSetObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeMemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeMemoryWorkingSetObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeMemoryWorkingSetObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.memory.working_set",
		opt...,
	)
	if err != nil {
		return NodeMemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeMemoryWorkingSetObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeMemoryWorkingSetObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeMemoryWorkingSetObservable) Name() string {
	return "k8s.node.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (NodeMemoryWorkingSetObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeMemoryWorkingSetObservable) Description() string {
	return "Node memory working set."
}

// NodeNetworkErrors is an instrument used to record metric values conforming to
// the "k8s.node.network.errors" semantic conventions. It represents the node
// network errors.
//
// Stability: Release_Candidate
type NodeNetworkErrors struct {
	metric.Int64Counter
}

var newNodeNetworkErrorsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Node network errors."),
	metric.WithUnit("{error}"),
}

// NewNodeNetworkErrors returns a new NodeNetworkErrors instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeNetworkErrorsOpts...)
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (NodeNetworkErrors) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NodeNetworkErrorsObservable is an instrument used to record metric values
// conforming to the "k8s.node.network.errors" semantic conventions. It
// represents the node network errors.
//
// Stability: Release_Candidate
type NodeNetworkErrorsObservable struct {
	metric.Int64ObservableCounter
}

var newNodeNetworkErrorsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Node network errors."),
	metric.WithUnit("{error}"),
}

// NewNodeNetworkErrorsObservable returns a new NodeNetworkErrorsObservable
// instrument.
//
// Stability: Release_Candidate
func NewNodeNetworkErrorsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (NodeNetworkErrorsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeNetworkErrorsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeNetworkErrorsObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeNetworkErrorsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"k8s.node.network.errors",
		opt...,
	)
	if err != nil {
		return NodeNetworkErrorsObservable{noop.Int64ObservableCounter{}}, err
	}
	return NodeNetworkErrorsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeNetworkErrorsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeNetworkErrorsObservable) Name() string {
	return "k8s.node.network.errors"
}

// Unit returns the semantic convention unit of the instrument
func (NodeNetworkErrorsObservable) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (NodeNetworkErrorsObservable) Description() string {
	return "Node network errors."
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NodeNetworkErrorsObservable) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (NodeNetworkErrorsObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NodeNetworkIO is an instrument used to record metric values conforming to the
// "k8s.node.network.io" semantic conventions. It represents the network bytes
// for the Node.
//
// Stability: Release_Candidate
type NodeNetworkIO struct {
	metric.Int64Counter
}

var newNodeNetworkIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Network bytes for the Node."),
	metric.WithUnit("By"),
}

// NewNodeNetworkIO returns a new NodeNetworkIO instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeNetworkIOOpts...)
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (NodeNetworkIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NodeNetworkIOObservable is an instrument used to record metric values
// conforming to the "k8s.node.network.io" semantic conventions. It represents
// the network bytes for the Node.
//
// Stability: Release_Candidate
type NodeNetworkIOObservable struct {
	metric.Int64ObservableCounter
}

var newNodeNetworkIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Network bytes for the Node."),
	metric.WithUnit("By"),
}

// NewNodeNetworkIOObservable returns a new NodeNetworkIOObservable instrument.
//
// Stability: Release_Candidate
func NewNodeNetworkIOObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (NodeNetworkIOObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeNetworkIOObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeNetworkIOObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeNetworkIOObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"k8s.node.network.io",
		opt...,
	)
	if err != nil {
		return NodeNetworkIOObservable{noop.Int64ObservableCounter{}}, err
	}
	return NodeNetworkIOObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeNetworkIOObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeNetworkIOObservable) Name() string {
	return "k8s.node.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (NodeNetworkIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeNetworkIOObservable) Description() string {
	return "Network bytes for the Node."
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (NodeNetworkIOObservable) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (NodeNetworkIOObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// NodePagingFaults is an instrument used to record metric values conforming to
// the "k8s.node.paging.faults" semantic conventions. It represents the node
// memory paging faults.
//
// Stability: Release_Candidate
type NodePagingFaults struct {
	metric.Int64Counter
}

var newNodePagingFaultsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Node memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewNodePagingFaults returns a new NodePagingFaults instrument.
//
// Stability: Release_Candidate
func NewNodePagingFaults(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (NodePagingFaults, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePagingFaults{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePagingFaultsOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodePagingFaultsOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.node.paging.faults",
		opt...,
	)
	if err != nil {
		return NodePagingFaults{noop.Int64Counter{}}, err
	}
	return NodePagingFaults{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePagingFaults) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (NodePagingFaults) Name() string {
	return "k8s.node.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (NodePagingFaults) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (NodePagingFaults) Description() string {
	return "Node memory paging faults."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Cumulative number of major and minor page faults.
// The [MemoryStats.MajorPageFaults] field of the [NodeStats.Memory] of the
// Kubelet's stats API maps directly to `system.paging.fault.type=major`.
// The [MemoryStats.PageFaults] field reports the total number of page faults
// (equivalent to `pgfault` in the cgroup `memory.stat`), with minor and major
// combined. The minor page faults MUST therefore be derived as
// `MemoryStats.PageFaults - MemoryStats.MajorPageFaults`, so that summing the
// metric over `system.paging.fault.type` yields the total number of page faults.
//
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
func (m NodePagingFaults) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// Cumulative number of major and minor page faults.
// The [MemoryStats.MajorPageFaults] field of the [NodeStats.Memory] of the
// Kubelet's stats API maps directly to `system.paging.fault.type=major`.
// The [MemoryStats.PageFaults] field reports the total number of page faults
// (equivalent to `pgfault` in the cgroup `memory.stat`), with minor and major
// combined. The minor page faults MUST therefore be derived as
// `MemoryStats.PageFaults - MemoryStats.MajorPageFaults`, so that summing the
// metric over `system.paging.fault.type` yields the total number of page faults.
//
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [NodeStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#NodeStats
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
func (m NodePagingFaults) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (NodePagingFaults) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// NodePagingFaultsObservable is an instrument used to record metric values
// conforming to the "k8s.node.paging.faults" semantic conventions. It represents
// the node memory paging faults.
//
// Stability: Release_Candidate
type NodePagingFaultsObservable struct {
	metric.Int64ObservableCounter
}

var newNodePagingFaultsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Node memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewNodePagingFaultsObservable returns a new NodePagingFaultsObservable
// instrument.
//
// Stability: Release_Candidate
func NewNodePagingFaultsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (NodePagingFaultsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePagingFaultsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePagingFaultsObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodePagingFaultsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"k8s.node.paging.faults",
		opt...,
	)
	if err != nil {
		return NodePagingFaultsObservable{noop.Int64ObservableCounter{}}, err
	}
	return NodePagingFaultsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePagingFaultsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (NodePagingFaultsObservable) Name() string {
	return "k8s.node.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (NodePagingFaultsObservable) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (NodePagingFaultsObservable) Description() string {
	return "Node memory paging faults."
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (NodePagingFaultsObservable) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// NodePIDLimit is an instrument used to record metric values conforming to the
// "k8s.node.pid.limit" semantic conventions. It represents the total number of
// processes/threads allowed by the operating system on the node.
//
// Stability: Development
type NodePIDLimit struct {
	metric.Int64UpDownCounter
}

var newNodePIDLimitOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of processes/threads allowed by the operating system on the node."),
	metric.WithUnit("{thread}"),
}

// NewNodePIDLimit returns a new NodePIDLimit instrument.
//
// Stability: Development
func NewNodePIDLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodePIDLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePIDLimit{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePIDLimitOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodePIDLimitOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.pid.limit",
		opt...,
	)
	if err != nil {
		return NodePIDLimit{noop.Int64UpDownCounter{}}, err
	}
	return NodePIDLimit{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePIDLimit) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodePIDLimit) Name() string {
	return "k8s.node.pid.limit"
}

// Unit returns the semantic convention unit of the instrument
func (NodePIDLimit) Unit() string {
	return "{thread}"
}

// Description returns the semantic convention description of the instrument
func (NodePIDLimit) Description() string {
	return "Total number of processes/threads allowed by the operating system on the node."
}

// Add adds incr to the existing count for attrs.
//
// Derived from the Kubelet Summary API
// ([RlimitStats.MaxPID])
// field of the [NodeStats.Rlimit],
// representing the system-wide task limit (the lower of `kernel.pid_max` and
// `kernel.threads-max`).
//
// [RlimitStats.MaxPID]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#RlimitStats
// [NodeStats.Rlimit]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodePIDLimit) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Derived from the Kubelet Summary API
// ([RlimitStats.MaxPID])
// field of the [NodeStats.Rlimit],
// representing the system-wide task limit (the lower of `kernel.pid_max` and
// `kernel.threads-max`).
//
// [RlimitStats.MaxPID]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#RlimitStats
// [NodeStats.Rlimit]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodePIDLimit) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodePIDLimitObservable is an instrument used to record metric values
// conforming to the "k8s.node.pid.limit" semantic conventions. It represents the
// total number of processes/threads allowed by the operating system on the node.
//
// Stability: Development
type NodePIDLimitObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodePIDLimitObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Total number of processes/threads allowed by the operating system on the node."),
	metric.WithUnit("{thread}"),
}

// NewNodePIDLimitObservable returns a new NodePIDLimitObservable instrument.
//
// Stability: Development
func NewNodePIDLimitObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodePIDLimitObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePIDLimitObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePIDLimitObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodePIDLimitObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.pid.limit",
		opt...,
	)
	if err != nil {
		return NodePIDLimitObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodePIDLimitObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePIDLimitObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodePIDLimitObservable) Name() string {
	return "k8s.node.pid.limit"
}

// Unit returns the semantic convention unit of the instrument
func (NodePIDLimitObservable) Unit() string {
	return "{thread}"
}

// Description returns the semantic convention description of the instrument
func (NodePIDLimitObservable) Description() string {
	return "Total number of processes/threads allowed by the operating system on the node."
}

// NodePIDUsage is an instrument used to record metric values conforming to the
// "k8s.node.pid.usage" semantic conventions. It represents the total number of
// existing processes and threads on the node.
//
// Stability: Development
type NodePIDUsage struct {
	metric.Int64UpDownCounter
}

var newNodePIDUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of existing processes and threads on the node."),
	metric.WithUnit("{thread}"),
}

// NewNodePIDUsage returns a new NodePIDUsage instrument.
//
// Stability: Development
func NewNodePIDUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodePIDUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePIDUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePIDUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodePIDUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.pid.usage",
		opt...,
	)
	if err != nil {
		return NodePIDUsage{noop.Int64UpDownCounter{}}, err
	}
	return NodePIDUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePIDUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodePIDUsage) Name() string {
	return "k8s.node.pid.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodePIDUsage) Unit() string {
	return "{thread}"
}

// Description returns the semantic convention description of the instrument
func (NodePIDUsage) Description() string {
	return "Total number of existing processes and threads on the node."
}

// Add adds incr to the existing count for attrs.
//
// Derived from the Kubelet Summary API
// ([RlimitStats.NumOfRunningProcesses])
// field of the [NodeStats.Rlimit],
// which reads the total count of scheduling entities (threads/processes across
// all states) from `/proc/loadavg`.
//
// [RlimitStats.NumOfRunningProcesses]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#RlimitStats
// [NodeStats.Rlimit]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodePIDUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Derived from the Kubelet Summary API
// ([RlimitStats.NumOfRunningProcesses])
// field of the [NodeStats.Rlimit],
// which reads the total count of scheduling entities (threads/processes across
// all states) from `/proc/loadavg`.
//
// [RlimitStats.NumOfRunningProcesses]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#RlimitStats
// [NodeStats.Rlimit]: https://pkg.go.dev/k8s.io/kubelet@v0.33.0/pkg/apis/stats/v1alpha1#NodeStats
func (m NodePIDUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodePIDUsageObservable is an instrument used to record metric values
// conforming to the "k8s.node.pid.usage" semantic conventions. It represents the
// total number of existing processes and threads on the node.
//
// Stability: Development
type NodePIDUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodePIDUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Total number of existing processes and threads on the node."),
	metric.WithUnit("{thread}"),
}

// NewNodePIDUsageObservable returns a new NodePIDUsageObservable instrument.
//
// Stability: Development
func NewNodePIDUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodePIDUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePIDUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePIDUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodePIDUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.pid.usage",
		opt...,
	)
	if err != nil {
		return NodePIDUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodePIDUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePIDUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodePIDUsageObservable) Name() string {
	return "k8s.node.pid.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodePIDUsageObservable) Unit() string {
	return "{thread}"
}

// Description returns the semantic convention description of the instrument
func (NodePIDUsageObservable) Description() string {
	return "Total number of existing processes and threads on the node."
}

// NodePodAllocatable is an instrument used to record metric values conforming to
// the "k8s.node.pod.allocatable" semantic conventions. It represents the amount
// of pods allocatable on the node.
//
// Stability: Development
type NodePodAllocatable struct {
	metric.Int64UpDownCounter
}

var newNodePodAllocatableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Amount of pods allocatable on the node."),
	metric.WithUnit("{pod}"),
}

// NewNodePodAllocatable returns a new NodePodAllocatable instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newNodePodAllocatableOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m NodePodAllocatable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodePodAllocatableObservable is an instrument used to record metric values
// conforming to the "k8s.node.pod.allocatable" semantic conventions. It
// represents the amount of pods allocatable on the node.
//
// Stability: Development
type NodePodAllocatableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodePodAllocatableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Amount of pods allocatable on the node."),
	metric.WithUnit("{pod}"),
}

// NewNodePodAllocatableObservable returns a new NodePodAllocatableObservable
// instrument.
//
// Stability: Development
func NewNodePodAllocatableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodePodAllocatableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodePodAllocatableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodePodAllocatableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodePodAllocatableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.pod.allocatable",
		opt...,
	)
	if err != nil {
		return NodePodAllocatableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodePodAllocatableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodePodAllocatableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodePodAllocatableObservable) Name() string {
	return "k8s.node.pod.allocatable"
}

// Unit returns the semantic convention unit of the instrument
func (NodePodAllocatableObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (NodePodAllocatableObservable) Description() string {
	return "Amount of pods allocatable on the node."
}

// NodeSystemContainerCPUTime is an instrument used to record metric values
// conforming to the "k8s.node.system_container.cpu.time" semantic conventions.
// It represents the node's system container CPU time.
//
// Stability: Development
type NodeSystemContainerCPUTime struct {
	metric.Float64Counter
}

var newNodeSystemContainerCPUTimeOpts = []metric.Float64CounterOption{
	metric.WithDescription("Node's system container CPU time."),
	metric.WithUnit("s"),
}

// NewNodeSystemContainerCPUTime returns a new NodeSystemContainerCPUTime
// instrument.
//
// Stability: Development
func NewNodeSystemContainerCPUTime(
	m metric.Meter,
	opt ...metric.Float64CounterOption,
) (NodeSystemContainerCPUTime, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerCPUTime{noop.Float64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerCPUTimeOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerCPUTimeOpts...)
	}

	i, err := m.Float64Counter(
		"k8s.node.system_container.cpu.time",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerCPUTime{noop.Float64Counter{}}, err
	}
	return NodeSystemContainerCPUTime{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerCPUTime) Inst() metric.Float64Counter {
	return m.Float64Counter
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerCPUTime) Name() string {
	return "k8s.node.system_container.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerCPUTime) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerCPUTime) Description() string {
	return "Node's system container CPU time."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the [CPUStats.UsageCoreNanoSeconds] field of the
// [ContainerStats] of [Node.SystemContainers] of the Kubelet's stats API.
//
// [CPUStats.UsageCoreNanoSeconds]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L236
// [ContainerStats]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L157C6-L157C20
// [Node.SystemContainers]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L40
func (m NodeSystemContainerCPUTime) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the [CPUStats.UsageCoreNanoSeconds] field of the
// [ContainerStats] of [Node.SystemContainers] of the Kubelet's stats API.
//
// [CPUStats.UsageCoreNanoSeconds]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L236
// [ContainerStats]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L157C6-L157C20
// [Node.SystemContainers]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L40
func (m NodeSystemContainerCPUTime) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// NodeSystemContainerCPUTimeObservable is an instrument used to record metric
// values conforming to the "k8s.node.system_container.cpu.time" semantic
// conventions. It represents the node's system container CPU time.
//
// Stability: Development
type NodeSystemContainerCPUTimeObservable struct {
	metric.Float64ObservableCounter
}

var newNodeSystemContainerCPUTimeObservableOpts = []metric.Float64ObservableCounterOption{
	metric.WithDescription("Node's system container CPU time."),
	metric.WithUnit("s"),
}

// NewNodeSystemContainerCPUTimeObservable returns a new
// NodeSystemContainerCPUTimeObservable instrument.
//
// Stability: Development
func NewNodeSystemContainerCPUTimeObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableCounterOption,
) (NodeSystemContainerCPUTimeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerCPUTimeObservable{noop.Float64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerCPUTimeObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerCPUTimeObservableOpts...)
	}

	i, err := m.Float64ObservableCounter(
		"k8s.node.system_container.cpu.time",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerCPUTimeObservable{noop.Float64ObservableCounter{}}, err
	}
	return NodeSystemContainerCPUTimeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerCPUTimeObservable) Inst() metric.Float64ObservableCounter {
	return m.Float64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerCPUTimeObservable) Name() string {
	return "k8s.node.system_container.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerCPUTimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerCPUTimeObservable) Description() string {
	return "Node's system container CPU time."
}

// NodeSystemContainerCPUUsage is an instrument used to record metric values
// conforming to the "k8s.node.system_container.cpu.usage" semantic conventions.
// It represents the node's system container CPU usage, measured in CPUs.
//
// Stability: Development
type NodeSystemContainerCPUUsage struct {
	metric.Float64Gauge
}

var newNodeSystemContainerCPUUsageOpts = []metric.Float64GaugeOption{
	metric.WithDescription("Node's system container CPU usage, measured in CPUs."),
	metric.WithUnit("{cpu}"),
}

// NewNodeSystemContainerCPUUsage returns a new NodeSystemContainerCPUUsage
// instrument.
//
// Stability: Development
func NewNodeSystemContainerCPUUsage(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (NodeSystemContainerCPUUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerCPUUsage{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerCPUUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerCPUUsageOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.node.system_container.cpu.usage",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerCPUUsage{noop.Float64Gauge{}}, err
	}
	return NodeSystemContainerCPUUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerCPUUsage) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerCPUUsage) Name() string {
	return "k8s.node.system_container.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerCPUUsage) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerCPUUsage) Description() string {
	return "Node's system container CPU usage, measured in CPUs."
}

// Record records val to the current distribution for attrs.
//
// CPU usage of the specific System Container on all available CPU cores. It is
// calculated as the change in cumulative CPU time
// (k8s.node.system_container.cpu.time) over a measurement interval, divided by
// the elapsed time: usageCores = (cpuTimeEnd - cpuTimeStart) / elapsedSeconds
func (m NodeSystemContainerCPUUsage) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// CPU usage of the specific System Container on all available CPU cores. It is
// calculated as the change in cumulative CPU time
// (k8s.node.system_container.cpu.time) over a measurement interval, divided by
// the elapsed time: usageCores = (cpuTimeEnd - cpuTimeStart) / elapsedSeconds
func (m NodeSystemContainerCPUUsage) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// NodeSystemContainerCPUUsageObservable is an instrument used to record metric
// values conforming to the "k8s.node.system_container.cpu.usage" semantic
// conventions. It represents the node's system container CPU usage, measured in
// CPUs.
//
// Stability: Development
type NodeSystemContainerCPUUsageObservable struct {
	metric.Float64ObservableGauge
}

var newNodeSystemContainerCPUUsageObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("Node's system container CPU usage, measured in CPUs."),
	metric.WithUnit("{cpu}"),
}

// NewNodeSystemContainerCPUUsageObservable returns a new
// NodeSystemContainerCPUUsageObservable instrument.
//
// Stability: Development
func NewNodeSystemContainerCPUUsageObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (NodeSystemContainerCPUUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerCPUUsageObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerCPUUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerCPUUsageObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.node.system_container.cpu.usage",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerCPUUsageObservable{noop.Float64ObservableGauge{}}, err
	}
	return NodeSystemContainerCPUUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerCPUUsageObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerCPUUsageObservable) Name() string {
	return "k8s.node.system_container.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerCPUUsageObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerCPUUsageObservable) Description() string {
	return "Node's system container CPU usage, measured in CPUs."
}

// NodeSystemContainerMemoryUsage is an instrument used to record metric values
// conforming to the "k8s.node.system_container.memory.usage" semantic
// conventions. It represents the node's system container memory usage.
//
// Stability: Development
type NodeSystemContainerMemoryUsage struct {
	metric.Int64UpDownCounter
}

var newNodeSystemContainerMemoryUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Node's system container memory usage."),
	metric.WithUnit("By"),
}

// NewNodeSystemContainerMemoryUsage returns a new NodeSystemContainerMemoryUsage
// instrument.
//
// Stability: Development
func NewNodeSystemContainerMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeSystemContainerMemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerMemoryUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerMemoryUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerMemoryUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.system_container.memory.usage",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerMemoryUsage{noop.Int64UpDownCounter{}}, err
	}
	return NodeSystemContainerMemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerMemoryUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerMemoryUsage) Name() string {
	return "k8s.node.system_container.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerMemoryUsage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerMemoryUsage) Description() string {
	return "Node's system container memory usage."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the [MemoryStats.UsageBytes] field of the
// [ContainerStats] of [Node.SystemContainers] of the Kubelet's stats API.
//
// [MemoryStats.UsageBytes]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L252
// [ContainerStats]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L157C6-L157C20
// [Node.SystemContainers]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L40
func (m NodeSystemContainerMemoryUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the [MemoryStats.UsageBytes] field of the
// [ContainerStats] of [Node.SystemContainers] of the Kubelet's stats API.
//
// [MemoryStats.UsageBytes]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L252
// [ContainerStats]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L157C6-L157C20
// [Node.SystemContainers]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L40
func (m NodeSystemContainerMemoryUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeSystemContainerMemoryUsageObservable is an instrument used to record
// metric values conforming to the "k8s.node.system_container.memory.usage"
// semantic conventions. It represents the node's system container memory usage.
//
// Stability: Development
type NodeSystemContainerMemoryUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeSystemContainerMemoryUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Node's system container memory usage."),
	metric.WithUnit("By"),
}

// NewNodeSystemContainerMemoryUsageObservable returns a new
// NodeSystemContainerMemoryUsageObservable instrument.
//
// Stability: Development
func NewNodeSystemContainerMemoryUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeSystemContainerMemoryUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerMemoryUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerMemoryUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.system_container.memory.usage",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeSystemContainerMemoryUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerMemoryUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerMemoryUsageObservable) Name() string {
	return "k8s.node.system_container.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerMemoryUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerMemoryUsageObservable) Description() string {
	return "Node's system container memory usage."
}

// NodeSystemContainerMemoryWorkingSet is an instrument used to record metric
// values conforming to the "k8s.node.system_container.memory.working_set"
// semantic conventions. It represents the amount of working set memory.
//
// Stability: Development
type NodeSystemContainerMemoryWorkingSet struct {
	metric.Int64UpDownCounter
}

var newNodeSystemContainerMemoryWorkingSetOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The amount of working set memory."),
	metric.WithUnit("By"),
}

// NewNodeSystemContainerMemoryWorkingSet returns a new
// NodeSystemContainerMemoryWorkingSet instrument.
//
// Stability: Development
func NewNodeSystemContainerMemoryWorkingSet(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeSystemContainerMemoryWorkingSet, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerMemoryWorkingSet{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerMemoryWorkingSetOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerMemoryWorkingSetOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.system_container.memory.working_set",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerMemoryWorkingSet{noop.Int64UpDownCounter{}}, err
	}
	return NodeSystemContainerMemoryWorkingSet{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerMemoryWorkingSet) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerMemoryWorkingSet) Name() string {
	return "k8s.node.system_container.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerMemoryWorkingSet) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerMemoryWorkingSet) Description() string {
	return "The amount of working set memory."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [ContainerStats] of [Node.SystemContainers] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L256
// [ContainerStats]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L157C6-L157C20
// [Node.SystemContainers]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L40
func (m NodeSystemContainerMemoryWorkingSet) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [ContainerStats] of [Node.SystemContainers] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L256
// [ContainerStats]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L157C6-L157C20
// [Node.SystemContainers]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L40
func (m NodeSystemContainerMemoryWorkingSet) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// NodeSystemContainerMemoryWorkingSetObservable is an instrument used to record
// metric values conforming to the "k8s.node.system_container.memory.working_set"
// semantic conventions. It represents the amount of working set memory.
//
// Stability: Development
type NodeSystemContainerMemoryWorkingSetObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newNodeSystemContainerMemoryWorkingSetObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The amount of working set memory."),
	metric.WithUnit("By"),
}

// NewNodeSystemContainerMemoryWorkingSetObservable returns a new
// NodeSystemContainerMemoryWorkingSetObservable instrument.
//
// Stability: Development
func NewNodeSystemContainerMemoryWorkingSetObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (NodeSystemContainerMemoryWorkingSetObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeSystemContainerMemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeSystemContainerMemoryWorkingSetObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeSystemContainerMemoryWorkingSetObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.node.system_container.memory.working_set",
		opt...,
	)
	if err != nil {
		return NodeSystemContainerMemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return NodeSystemContainerMemoryWorkingSetObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeSystemContainerMemoryWorkingSetObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeSystemContainerMemoryWorkingSetObservable) Name() string {
	return "k8s.node.system_container.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (NodeSystemContainerMemoryWorkingSetObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeSystemContainerMemoryWorkingSetObservable) Description() string {
	return "The amount of working set memory."
}

// NodeUptime is an instrument used to record metric values conforming to the
// "k8s.node.uptime" semantic conventions. It represents the time the Node has
// been running.
//
// Stability: Release_Candidate
type NodeUptime struct {
	metric.Float64Gauge
}

var newNodeUptimeOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The time the Node has been running."),
	metric.WithUnit("s"),
}

// NewNodeUptime returns a new NodeUptime instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newNodeUptimeOpts...)
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
// seconds as a floating point number with the highest precision available. The
// actual accuracy would depend on the instrumentation and operating system.
// This metric is derived from the [NodeStats.StartTime] field of Kubelet's stats
// API, which is set to the [system boot time]. This value resets when the node
// reboots.
//
// [NodeStats.StartTime]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L42
// [system boot time]: https://github.com/kubernetes/kubernetes/blob/v1.36.3/pkg/kubelet/server/stats/summary.go#L113
func (m NodeUptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available. The
// actual accuracy would depend on the instrumentation and operating system.
// This metric is derived from the [NodeStats.StartTime] field of Kubelet's stats
// API, which is set to the [system boot time]. This value resets when the node
// reboots.
//
// [NodeStats.StartTime]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L42
// [system boot time]: https://github.com/kubernetes/kubernetes/blob/v1.36.3/pkg/kubelet/server/stats/summary.go#L113
func (m NodeUptime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// NodeUptimeObservable is an instrument used to record metric values conforming
// to the "k8s.node.uptime" semantic conventions. It represents the time the Node
// has been running.
//
// Stability: Release_Candidate
type NodeUptimeObservable struct {
	metric.Float64ObservableGauge
}

var newNodeUptimeObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The time the Node has been running."),
	metric.WithUnit("s"),
}

// NewNodeUptimeObservable returns a new NodeUptimeObservable instrument.
//
// Stability: Release_Candidate
func NewNodeUptimeObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (NodeUptimeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeUptimeObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newNodeUptimeObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newNodeUptimeObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.node.uptime",
		opt...,
	)
	if err != nil {
		return NodeUptimeObservable{noop.Float64ObservableGauge{}}, err
	}
	return NodeUptimeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeUptimeObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (NodeUptimeObservable) Name() string {
	return "k8s.node.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (NodeUptimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (NodeUptimeObservable) Description() string {
	return "The time the Node has been running."
}

// PersistentvolumeStatusPhase is an instrument used to record metric values
// conforming to the "k8s.persistentvolume.status.phase" semantic conventions. It
// represents the number of PersistentVolumes in a given phase.
//
// Stability: Development
type PersistentvolumeStatusPhase struct {
	metric.Int64UpDownCounter
}

var newPersistentvolumeStatusPhaseOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of PersistentVolumes in a given phase."),
	metric.WithUnit("{persistentvolume}"),
}

// NewPersistentvolumeStatusPhase returns a new PersistentvolumeStatusPhase
// instrument.
//
// Stability: Development
func NewPersistentvolumeStatusPhase(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PersistentvolumeStatusPhase, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeStatusPhase{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeStatusPhaseOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeStatusPhaseOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.persistentvolume.status.phase",
		opt...,
	)
	if err != nil {
		return PersistentvolumeStatusPhase{noop.Int64UpDownCounter{}}, err
	}
	return PersistentvolumeStatusPhase{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeStatusPhase) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeStatusPhase) Name() string {
	return "k8s.persistentvolume.status.phase"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeStatusPhase) Unit() string {
	return "{persistentvolume}"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeStatusPhase) Description() string {
	return "Number of PersistentVolumes in a given phase."
}

// Add adds incr to the existing count for attrs.
//
// The persistentvolumeStatusPhase is the the phase of the PersistentVolume.
//
// All possible phases should be reported at each interval to avoid gaps in the
// time series.
// This metric is derived from the `.status.phase` field of the
// [K8s PersistentVolumeStatus].
//
// [K8s PersistentVolumeStatus]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-v1/#PersistentVolumeStatus
func (m PersistentvolumeStatusPhase) Add(
	ctx context.Context,
	incr int64,
	persistentvolumeStatusPhase PersistentvolumeStatusPhaseAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.persistentvolume.status.phase", string(persistentvolumeStatusPhase)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("k8s.persistentvolume.status.phase", string(persistentvolumeStatusPhase)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// All possible phases should be reported at each interval to avoid gaps in the
// time series.
// This metric is derived from the `.status.phase` field of the
// [K8s PersistentVolumeStatus].
//
// [K8s PersistentVolumeStatus]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-v1/#PersistentVolumeStatus
func (m PersistentvolumeStatusPhase) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PersistentvolumeStatusPhaseObservable is an instrument used to record metric
// values conforming to the "k8s.persistentvolume.status.phase" semantic
// conventions. It represents the number of PersistentVolumes in a given phase.
//
// Stability: Development
type PersistentvolumeStatusPhaseObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPersistentvolumeStatusPhaseObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of PersistentVolumes in a given phase."),
	metric.WithUnit("{persistentvolume}"),
}

// NewPersistentvolumeStatusPhaseObservable returns a new
// PersistentvolumeStatusPhaseObservable instrument.
//
// Stability: Development
func NewPersistentvolumeStatusPhaseObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PersistentvolumeStatusPhaseObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeStatusPhaseObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeStatusPhaseObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeStatusPhaseObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.persistentvolume.status.phase",
		opt...,
	)
	if err != nil {
		return PersistentvolumeStatusPhaseObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PersistentvolumeStatusPhaseObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeStatusPhaseObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeStatusPhaseObservable) Name() string {
	return "k8s.persistentvolume.status.phase"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeStatusPhaseObservable) Unit() string {
	return "{persistentvolume}"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeStatusPhaseObservable) Description() string {
	return "Number of PersistentVolumes in a given phase."
}

// AttrPersistentvolumeStatusPhase returns a required attribute for the
// "k8s.persistentvolume.status.phase" semantic convention. It represents the
// phase of the PersistentVolume.
func (PersistentvolumeStatusPhaseObservable) AttrPersistentvolumeStatusPhase(val PersistentvolumeStatusPhaseAttr) attribute.KeyValue {
	return attribute.String("k8s.persistentvolume.status.phase", string(val))
}

// PersistentvolumeStorageCapacity is an instrument used to record metric values
// conforming to the "k8s.persistentvolume.storage.capacity" semantic
// conventions. It represents the storage capacity of the PersistentVolume.
//
// Stability: Development
type PersistentvolumeStorageCapacity struct {
	metric.Int64UpDownCounter
}

var newPersistentvolumeStorageCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The storage capacity of the PersistentVolume."),
	metric.WithUnit("By"),
}

// NewPersistentvolumeStorageCapacity returns a new
// PersistentvolumeStorageCapacity instrument.
//
// Stability: Development
func NewPersistentvolumeStorageCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PersistentvolumeStorageCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeStorageCapacity{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeStorageCapacityOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeStorageCapacityOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.persistentvolume.storage.capacity",
		opt...,
	)
	if err != nil {
		return PersistentvolumeStorageCapacity{noop.Int64UpDownCounter{}}, err
	}
	return PersistentvolumeStorageCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeStorageCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeStorageCapacity) Name() string {
	return "k8s.persistentvolume.storage.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeStorageCapacity) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeStorageCapacity) Description() string {
	return "The storage capacity of the PersistentVolume."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the `.spec.capacity.storage` field of the
// [K8s PersistentVolumeSpec].
//
// [K8s PersistentVolumeSpec]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-v1/#PersistentVolumeSpec
func (m PersistentvolumeStorageCapacity) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the `.spec.capacity.storage` field of the
// [K8s PersistentVolumeSpec].
//
// [K8s PersistentVolumeSpec]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-v1/#PersistentVolumeSpec
func (m PersistentvolumeStorageCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PersistentvolumeStorageCapacityObservable is an instrument used to record
// metric values conforming to the "k8s.persistentvolume.storage.capacity"
// semantic conventions. It represents the storage capacity of the
// PersistentVolume.
//
// Stability: Development
type PersistentvolumeStorageCapacityObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPersistentvolumeStorageCapacityObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The storage capacity of the PersistentVolume."),
	metric.WithUnit("By"),
}

// NewPersistentvolumeStorageCapacityObservable returns a new
// PersistentvolumeStorageCapacityObservable instrument.
//
// Stability: Development
func NewPersistentvolumeStorageCapacityObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PersistentvolumeStorageCapacityObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeStorageCapacityObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeStorageCapacityObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeStorageCapacityObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.persistentvolume.storage.capacity",
		opt...,
	)
	if err != nil {
		return PersistentvolumeStorageCapacityObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PersistentvolumeStorageCapacityObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeStorageCapacityObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeStorageCapacityObservable) Name() string {
	return "k8s.persistentvolume.storage.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeStorageCapacityObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeStorageCapacityObservable) Description() string {
	return "The storage capacity of the PersistentVolume."
}

// PersistentvolumeclaimStatusPhase is an instrument used to record metric values
// conforming to the "k8s.persistentvolumeclaim.status.phase" semantic
// conventions. It represents the number of PersistentVolumeClaims in a given
// phase.
//
// Stability: Development
type PersistentvolumeclaimStatusPhase struct {
	metric.Int64UpDownCounter
}

var newPersistentvolumeclaimStatusPhaseOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of PersistentVolumeClaims in a given phase."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewPersistentvolumeclaimStatusPhase returns a new
// PersistentvolumeclaimStatusPhase instrument.
//
// Stability: Development
func NewPersistentvolumeclaimStatusPhase(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PersistentvolumeclaimStatusPhase, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeclaimStatusPhase{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeclaimStatusPhaseOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeclaimStatusPhaseOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.persistentvolumeclaim.status.phase",
		opt...,
	)
	if err != nil {
		return PersistentvolumeclaimStatusPhase{noop.Int64UpDownCounter{}}, err
	}
	return PersistentvolumeclaimStatusPhase{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeclaimStatusPhase) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeclaimStatusPhase) Name() string {
	return "k8s.persistentvolumeclaim.status.phase"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeclaimStatusPhase) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeclaimStatusPhase) Description() string {
	return "Number of PersistentVolumeClaims in a given phase."
}

// Add adds incr to the existing count for attrs.
//
// The persistentvolumeclaimStatusPhase is the the phase of the
// PersistentVolumeClaim.
//
// All possible phases should be reported at each interval to avoid gaps in the
// time series.
// This metric is derived from the `.status.phase` field of the
// [K8s PersistentVolumeClaimStatus].
//
// [K8s PersistentVolumeClaimStatus]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#PersistentVolumeClaimStatus
func (m PersistentvolumeclaimStatusPhase) Add(
	ctx context.Context,
	incr int64,
	persistentvolumeclaimStatusPhase PersistentvolumeclaimStatusPhaseAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.persistentvolumeclaim.status.phase", string(persistentvolumeclaimStatusPhase)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("k8s.persistentvolumeclaim.status.phase", string(persistentvolumeclaimStatusPhase)),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// All possible phases should be reported at each interval to avoid gaps in the
// time series.
// This metric is derived from the `.status.phase` field of the
// [K8s PersistentVolumeClaimStatus].
//
// [K8s PersistentVolumeClaimStatus]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#PersistentVolumeClaimStatus
func (m PersistentvolumeclaimStatusPhase) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PersistentvolumeclaimStatusPhaseObservable is an instrument used to record
// metric values conforming to the "k8s.persistentvolumeclaim.status.phase"
// semantic conventions. It represents the number of PersistentVolumeClaims in a
// given phase.
//
// Stability: Development
type PersistentvolumeclaimStatusPhaseObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPersistentvolumeclaimStatusPhaseObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of PersistentVolumeClaims in a given phase."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewPersistentvolumeclaimStatusPhaseObservable returns a new
// PersistentvolumeclaimStatusPhaseObservable instrument.
//
// Stability: Development
func NewPersistentvolumeclaimStatusPhaseObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PersistentvolumeclaimStatusPhaseObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeclaimStatusPhaseObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeclaimStatusPhaseObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeclaimStatusPhaseObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.persistentvolumeclaim.status.phase",
		opt...,
	)
	if err != nil {
		return PersistentvolumeclaimStatusPhaseObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PersistentvolumeclaimStatusPhaseObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeclaimStatusPhaseObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeclaimStatusPhaseObservable) Name() string {
	return "k8s.persistentvolumeclaim.status.phase"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeclaimStatusPhaseObservable) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeclaimStatusPhaseObservable) Description() string {
	return "Number of PersistentVolumeClaims in a given phase."
}

// AttrPersistentvolumeclaimStatusPhase returns a required attribute for the
// "k8s.persistentvolumeclaim.status.phase" semantic convention. It represents
// the phase of the PersistentVolumeClaim.
func (PersistentvolumeclaimStatusPhaseObservable) AttrPersistentvolumeclaimStatusPhase(val PersistentvolumeclaimStatusPhaseAttr) attribute.KeyValue {
	return attribute.String("k8s.persistentvolumeclaim.status.phase", string(val))
}

// PersistentvolumeclaimStorageCapacity is an instrument used to record metric
// values conforming to the "k8s.persistentvolumeclaim.storage.capacity" semantic
// conventions. It represents the actual storage capacity provisioned for the
// PersistentVolumeClaim.
//
// Stability: Development
type PersistentvolumeclaimStorageCapacity struct {
	metric.Int64UpDownCounter
}

var newPersistentvolumeclaimStorageCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The actual storage capacity provisioned for the PersistentVolumeClaim."),
	metric.WithUnit("By"),
}

// NewPersistentvolumeclaimStorageCapacity returns a new
// PersistentvolumeclaimStorageCapacity instrument.
//
// Stability: Development
func NewPersistentvolumeclaimStorageCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PersistentvolumeclaimStorageCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeclaimStorageCapacity{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeclaimStorageCapacityOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeclaimStorageCapacityOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.persistentvolumeclaim.storage.capacity",
		opt...,
	)
	if err != nil {
		return PersistentvolumeclaimStorageCapacity{noop.Int64UpDownCounter{}}, err
	}
	return PersistentvolumeclaimStorageCapacity{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeclaimStorageCapacity) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeclaimStorageCapacity) Name() string {
	return "k8s.persistentvolumeclaim.storage.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeclaimStorageCapacity) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeclaimStorageCapacity) Description() string {
	return "The actual storage capacity provisioned for the PersistentVolumeClaim."
}

// Add adds incr to the existing count for attrs.
//
// Only available when the PVC is bound. May differ from the requested capacity
// due to provisioner rounding.
// This metric is derived from the `.status.capacity.storage` field of the
// [K8s PersistentVolumeClaimStatus].
//
// [K8s PersistentVolumeClaimStatus]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#PersistentVolumeClaimStatus
func (m PersistentvolumeclaimStorageCapacity) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Only available when the PVC is bound. May differ from the requested capacity
// due to provisioner rounding.
// This metric is derived from the `.status.capacity.storage` field of the
// [K8s PersistentVolumeClaimStatus].
//
// [K8s PersistentVolumeClaimStatus]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#PersistentVolumeClaimStatus
func (m PersistentvolumeclaimStorageCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PersistentvolumeclaimStorageCapacityObservable is an instrument used to record
// metric values conforming to the "k8s.persistentvolumeclaim.storage.capacity"
// semantic conventions. It represents the actual storage capacity provisioned
// for the PersistentVolumeClaim.
//
// Stability: Development
type PersistentvolumeclaimStorageCapacityObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPersistentvolumeclaimStorageCapacityObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The actual storage capacity provisioned for the PersistentVolumeClaim."),
	metric.WithUnit("By"),
}

// NewPersistentvolumeclaimStorageCapacityObservable returns a new
// PersistentvolumeclaimStorageCapacityObservable instrument.
//
// Stability: Development
func NewPersistentvolumeclaimStorageCapacityObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PersistentvolumeclaimStorageCapacityObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeclaimStorageCapacityObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeclaimStorageCapacityObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeclaimStorageCapacityObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.persistentvolumeclaim.storage.capacity",
		opt...,
	)
	if err != nil {
		return PersistentvolumeclaimStorageCapacityObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PersistentvolumeclaimStorageCapacityObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeclaimStorageCapacityObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeclaimStorageCapacityObservable) Name() string {
	return "k8s.persistentvolumeclaim.storage.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeclaimStorageCapacityObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeclaimStorageCapacityObservable) Description() string {
	return "The actual storage capacity provisioned for the PersistentVolumeClaim."
}

// PersistentvolumeclaimStorageRequest is an instrument used to record metric
// values conforming to the "k8s.persistentvolumeclaim.storage.request" semantic
// conventions. It represents the storage requested by the PersistentVolumeClaim.
//
// Stability: Development
type PersistentvolumeclaimStorageRequest struct {
	metric.Int64UpDownCounter
}

var newPersistentvolumeclaimStorageRequestOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The storage requested by the PersistentVolumeClaim."),
	metric.WithUnit("By"),
}

// NewPersistentvolumeclaimStorageRequest returns a new
// PersistentvolumeclaimStorageRequest instrument.
//
// Stability: Development
func NewPersistentvolumeclaimStorageRequest(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PersistentvolumeclaimStorageRequest, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeclaimStorageRequest{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeclaimStorageRequestOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeclaimStorageRequestOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.persistentvolumeclaim.storage.request",
		opt...,
	)
	if err != nil {
		return PersistentvolumeclaimStorageRequest{noop.Int64UpDownCounter{}}, err
	}
	return PersistentvolumeclaimStorageRequest{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeclaimStorageRequest) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeclaimStorageRequest) Name() string {
	return "k8s.persistentvolumeclaim.storage.request"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeclaimStorageRequest) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeclaimStorageRequest) Description() string {
	return "The storage requested by the PersistentVolumeClaim."
}

// Add adds incr to the existing count for attrs.
//
// This metric is derived from the `.spec.resources.requests.storage` field of
// the [K8s PersistentVolumeClaimSpec].
//
// [K8s PersistentVolumeClaimSpec]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#PersistentVolumeClaimSpec
func (m PersistentvolumeclaimStorageRequest) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the `.spec.resources.requests.storage` field of
// the [K8s PersistentVolumeClaimSpec].
//
// [K8s PersistentVolumeClaimSpec]: https://kubernetes.io/docs/reference/kubernetes-api/config-and-storage-resources/persistent-volume-claim-v1/#PersistentVolumeClaimSpec
func (m PersistentvolumeclaimStorageRequest) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PersistentvolumeclaimStorageRequestObservable is an instrument used to record
// metric values conforming to the "k8s.persistentvolumeclaim.storage.request"
// semantic conventions. It represents the storage requested by the
// PersistentVolumeClaim.
//
// Stability: Development
type PersistentvolumeclaimStorageRequestObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPersistentvolumeclaimStorageRequestObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The storage requested by the PersistentVolumeClaim."),
	metric.WithUnit("By"),
}

// NewPersistentvolumeclaimStorageRequestObservable returns a new
// PersistentvolumeclaimStorageRequestObservable instrument.
//
// Stability: Development
func NewPersistentvolumeclaimStorageRequestObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PersistentvolumeclaimStorageRequestObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PersistentvolumeclaimStorageRequestObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPersistentvolumeclaimStorageRequestObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPersistentvolumeclaimStorageRequestObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.persistentvolumeclaim.storage.request",
		opt...,
	)
	if err != nil {
		return PersistentvolumeclaimStorageRequestObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PersistentvolumeclaimStorageRequestObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PersistentvolumeclaimStorageRequestObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PersistentvolumeclaimStorageRequestObservable) Name() string {
	return "k8s.persistentvolumeclaim.storage.request"
}

// Unit returns the semantic convention unit of the instrument
func (PersistentvolumeclaimStorageRequestObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PersistentvolumeclaimStorageRequestObservable) Description() string {
	return "The storage requested by the PersistentVolumeClaim."
}

// PodCPULimitCurrent is an instrument used to record metric values conforming to
// the "k8s.pod.cpu.limit.current" semantic conventions. It represents the
// maximum CPU resource limit currently configured for a running Pod.
//
// Stability: Development
type PodCPULimitCurrent struct {
	metric.Float64UpDownCounter
}

var newPodCPULimitCurrentOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit currently configured for a running Pod."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPULimitCurrent returns a new PodCPULimitCurrent instrument.
//
// Stability: Development
func NewPodCPULimitCurrent(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (PodCPULimitCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPULimitCurrent{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPULimitCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPULimitCurrentOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.pod.cpu.limit.current",
		opt...,
	)
	if err != nil {
		return PodCPULimitCurrent{noop.Float64UpDownCounter{}}, err
	}
	return PodCPULimitCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPULimitCurrent) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPULimitCurrent) Name() string {
	return "k8s.pod.cpu.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPULimitCurrent) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPULimitCurrent) Description() string {
	return "Maximum CPU resource limit currently configured for a running Pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPULimitCurrent) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPULimitCurrent) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// PodCPULimitCurrentObservable is an instrument used to record metric values
// conforming to the "k8s.pod.cpu.limit.current" semantic conventions. It
// represents the maximum CPU resource limit currently configured for a running
// Pod.
//
// Stability: Development
type PodCPULimitCurrentObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newPodCPULimitCurrentObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit currently configured for a running Pod."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPULimitCurrentObservable returns a new PodCPULimitCurrentObservable
// instrument.
//
// Stability: Development
func NewPodCPULimitCurrentObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (PodCPULimitCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPULimitCurrentObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPULimitCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPULimitCurrentObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.pod.cpu.limit.current",
		opt...,
	)
	if err != nil {
		return PodCPULimitCurrentObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return PodCPULimitCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPULimitCurrentObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPULimitCurrentObservable) Name() string {
	return "k8s.pod.cpu.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPULimitCurrentObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPULimitCurrentObservable) Description() string {
	return "Maximum CPU resource limit currently configured for a running Pod."
}

// PodCPULimitDesired is an instrument used to record metric values conforming to
// the "k8s.pod.cpu.limit.desired" semantic conventions. It represents the
// maximum CPU resource limit as defined by the Pod spec.
//
// Stability: Development
type PodCPULimitDesired struct {
	metric.Float64UpDownCounter
}

var newPodCPULimitDesiredOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit as defined by the Pod spec."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPULimitDesired returns a new PodCPULimitDesired instrument.
//
// Stability: Development
func NewPodCPULimitDesired(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (PodCPULimitDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPULimitDesired{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPULimitDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPULimitDesiredOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.pod.cpu.limit.desired",
		opt...,
	)
	if err != nil {
		return PodCPULimitDesired{noop.Float64UpDownCounter{}}, err
	}
	return PodCPULimitDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPULimitDesired) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPULimitDesired) Name() string {
	return "k8s.pod.cpu.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPULimitDesired) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPULimitDesired) Description() string {
	return "Maximum CPU resource limit as defined by the Pod spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level limits
// based on the aggregated limits of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPULimitDesired) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level limits
// based on the aggregated limits of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPULimitDesired) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// PodCPULimitDesiredObservable is an instrument used to record metric values
// conforming to the "k8s.pod.cpu.limit.desired" semantic conventions. It
// represents the maximum CPU resource limit as defined by the Pod spec.
//
// Stability: Development
type PodCPULimitDesiredObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newPodCPULimitDesiredObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum CPU resource limit as defined by the Pod spec."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPULimitDesiredObservable returns a new PodCPULimitDesiredObservable
// instrument.
//
// Stability: Development
func NewPodCPULimitDesiredObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (PodCPULimitDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPULimitDesiredObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPULimitDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPULimitDesiredObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.pod.cpu.limit.desired",
		opt...,
	)
	if err != nil {
		return PodCPULimitDesiredObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return PodCPULimitDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPULimitDesiredObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPULimitDesiredObservable) Name() string {
	return "k8s.pod.cpu.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPULimitDesiredObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPULimitDesiredObservable) Description() string {
	return "Maximum CPU resource limit as defined by the Pod spec."
}

// PodCPULimitUtilization is an instrument used to record metric values
// conforming to the "k8s.pod.cpu.limit.utilization" semantic conventions. It
// represents the ratio of Pod CPU usage to its current CPU limit.
//
// Stability: Development
type PodCPULimitUtilization struct {
	metric.Float64Gauge
}

var newPodCPULimitUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of Pod CPU usage to its current CPU limit."),
	metric.WithUnit("1"),
}

// NewPodCPULimitUtilization returns a new PodCPULimitUtilization instrument.
//
// Stability: Development
func NewPodCPULimitUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (PodCPULimitUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPULimitUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPULimitUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPULimitUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.pod.cpu.limit.utilization",
		opt...,
	)
	if err != nil {
		return PodCPULimitUtilization{noop.Float64Gauge{}}, err
	}
	return PodCPULimitUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPULimitUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PodCPULimitUtilization) Name() string {
	return "k8s.pod.cpu.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPULimitUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodCPULimitUtilization) Description() string {
	return "The ratio of Pod CPU usage to its current CPU limit."
}

// Record records val to the current distribution for attrs.
//
// The current CPU limit reflects the actual resources applied at the Pod level,
// as reported by
// [PodStatus].
// The value range is [0.0,1.0]. A value of 1.0 means the Pod is using 100% of
// its actual CPU limit.
// If the CPU limit is not set, this metric SHOULD NOT be emitted for that Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodCPULimitUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The current CPU limit reflects the actual resources applied at the Pod level,
// as reported by
// [PodStatus].
// The value range is [0.0,1.0]. A value of 1.0 means the Pod is using 100% of
// its actual CPU limit.
// If the CPU limit is not set, this metric SHOULD NOT be emitted for that Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodCPULimitUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodCPULimitUtilizationObservable is an instrument used to record metric values
// conforming to the "k8s.pod.cpu.limit.utilization" semantic conventions. It
// represents the ratio of Pod CPU usage to its current CPU limit.
//
// Stability: Development
type PodCPULimitUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newPodCPULimitUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of Pod CPU usage to its current CPU limit."),
	metric.WithUnit("1"),
}

// NewPodCPULimitUtilizationObservable returns a new
// PodCPULimitUtilizationObservable instrument.
//
// Stability: Development
func NewPodCPULimitUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (PodCPULimitUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPULimitUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPULimitUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPULimitUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.pod.cpu.limit.utilization",
		opt...,
	)
	if err != nil {
		return PodCPULimitUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return PodCPULimitUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPULimitUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PodCPULimitUtilizationObservable) Name() string {
	return "k8s.pod.cpu.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPULimitUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodCPULimitUtilizationObservable) Description() string {
	return "The ratio of Pod CPU usage to its current CPU limit."
}

// PodCPURequestCurrent is an instrument used to record metric values conforming
// to the "k8s.pod.cpu.request.current" semantic conventions. It represents the
// CPU resource requested currently configured for a running Pod.
//
// Stability: Development
type PodCPURequestCurrent struct {
	metric.Float64UpDownCounter
}

var newPodCPURequestCurrentOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("CPU resource requested currently configured for a running Pod."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPURequestCurrent returns a new PodCPURequestCurrent instrument.
//
// Stability: Development
func NewPodCPURequestCurrent(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (PodCPURequestCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPURequestCurrent{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPURequestCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPURequestCurrentOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.pod.cpu.request.current",
		opt...,
	)
	if err != nil {
		return PodCPURequestCurrent{noop.Float64UpDownCounter{}}, err
	}
	return PodCPURequestCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPURequestCurrent) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPURequestCurrent) Name() string {
	return "k8s.pod.cpu.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPURequestCurrent) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPURequestCurrent) Description() string {
	return "CPU resource requested currently configured for a running Pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPURequestCurrent) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPURequestCurrent) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// PodCPURequestCurrentObservable is an instrument used to record metric values
// conforming to the "k8s.pod.cpu.request.current" semantic conventions. It
// represents the CPU resource requested currently configured for a running Pod.
//
// Stability: Development
type PodCPURequestCurrentObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newPodCPURequestCurrentObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("CPU resource requested currently configured for a running Pod."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPURequestCurrentObservable returns a new PodCPURequestCurrentObservable
// instrument.
//
// Stability: Development
func NewPodCPURequestCurrentObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (PodCPURequestCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPURequestCurrentObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPURequestCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPURequestCurrentObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.pod.cpu.request.current",
		opt...,
	)
	if err != nil {
		return PodCPURequestCurrentObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return PodCPURequestCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPURequestCurrentObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPURequestCurrentObservable) Name() string {
	return "k8s.pod.cpu.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPURequestCurrentObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPURequestCurrentObservable) Description() string {
	return "CPU resource requested currently configured for a running Pod."
}

// PodCPURequestDesired is an instrument used to record metric values conforming
// to the "k8s.pod.cpu.request.desired" semantic conventions. It represents the
// CPU resource requested as defined by the Pod spec.
//
// Stability: Development
type PodCPURequestDesired struct {
	metric.Float64UpDownCounter
}

var newPodCPURequestDesiredOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("CPU resource requested as defined by the Pod spec."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPURequestDesired returns a new PodCPURequestDesired instrument.
//
// Stability: Development
func NewPodCPURequestDesired(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (PodCPURequestDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPURequestDesired{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPURequestDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPURequestDesiredOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.pod.cpu.request.desired",
		opt...,
	)
	if err != nil {
		return PodCPURequestDesired{noop.Float64UpDownCounter{}}, err
	}
	return PodCPURequestDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPURequestDesired) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPURequestDesired) Name() string {
	return "k8s.pod.cpu.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPURequestDesired) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPURequestDesired) Description() string {
	return "CPU resource requested as defined by the Pod spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level requests
// based on the aggregated requests of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPURequestDesired) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level requests
// based on the aggregated requests of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodCPURequestDesired) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// PodCPURequestDesiredObservable is an instrument used to record metric values
// conforming to the "k8s.pod.cpu.request.desired" semantic conventions. It
// represents the CPU resource requested as defined by the Pod spec.
//
// Stability: Development
type PodCPURequestDesiredObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newPodCPURequestDesiredObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("CPU resource requested as defined by the Pod spec."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPURequestDesiredObservable returns a new PodCPURequestDesiredObservable
// instrument.
//
// Stability: Development
func NewPodCPURequestDesiredObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (PodCPURequestDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPURequestDesiredObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPURequestDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPURequestDesiredObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.pod.cpu.request.desired",
		opt...,
	)
	if err != nil {
		return PodCPURequestDesiredObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return PodCPURequestDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPURequestDesiredObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPURequestDesiredObservable) Name() string {
	return "k8s.pod.cpu.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPURequestDesiredObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPURequestDesiredObservable) Description() string {
	return "CPU resource requested as defined by the Pod spec."
}

// PodCPURequestUtilization is an instrument used to record metric values
// conforming to the "k8s.pod.cpu.request.utilization" semantic conventions. It
// represents the ratio of Pod CPU usage to its current CPU request.
//
// Stability: Development
type PodCPURequestUtilization struct {
	metric.Float64Gauge
}

var newPodCPURequestUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of Pod CPU usage to its current CPU request."),
	metric.WithUnit("1"),
}

// NewPodCPURequestUtilization returns a new PodCPURequestUtilization instrument.
//
// Stability: Development
func NewPodCPURequestUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (PodCPURequestUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPURequestUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPURequestUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPURequestUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.pod.cpu.request.utilization",
		opt...,
	)
	if err != nil {
		return PodCPURequestUtilization{noop.Float64Gauge{}}, err
	}
	return PodCPURequestUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPURequestUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PodCPURequestUtilization) Name() string {
	return "k8s.pod.cpu.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPURequestUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodCPURequestUtilization) Description() string {
	return "The ratio of Pod CPU usage to its current CPU request."
}

// Record records val to the current distribution for attrs.
//
// The current CPU request reflects the request applied at the Pod level, as
// reported by
// [PodStatus].
// A value of 1.0 means the Pod is using exactly its CPU request. Since a CPU
// request is a scheduling
// reservation and not a hard cap, values above 1.0 are possible and indicate the
// Pod is using more CPU
// than requested. If the CPU request is not set, this metric SHOULD NOT be
// emitted for that Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodCPURequestUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The current CPU request reflects the request applied at the Pod level, as
// reported by
// [PodStatus].
// A value of 1.0 means the Pod is using exactly its CPU request. Since a CPU
// request is a scheduling
// reservation and not a hard cap, values above 1.0 are possible and indicate the
// Pod is using more CPU
// than requested. If the CPU request is not set, this metric SHOULD NOT be
// emitted for that Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodCPURequestUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodCPURequestUtilizationObservable is an instrument used to record metric
// values conforming to the "k8s.pod.cpu.request.utilization" semantic
// conventions. It represents the ratio of Pod CPU usage to its current CPU
// request.
//
// Stability: Development
type PodCPURequestUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newPodCPURequestUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of Pod CPU usage to its current CPU request."),
	metric.WithUnit("1"),
}

// NewPodCPURequestUtilizationObservable returns a new
// PodCPURequestUtilizationObservable instrument.
//
// Stability: Development
func NewPodCPURequestUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (PodCPURequestUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPURequestUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPURequestUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPURequestUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.pod.cpu.request.utilization",
		opt...,
	)
	if err != nil {
		return PodCPURequestUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return PodCPURequestUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPURequestUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PodCPURequestUtilizationObservable) Name() string {
	return "k8s.pod.cpu.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPURequestUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodCPURequestUtilizationObservable) Description() string {
	return "The ratio of Pod CPU usage to its current CPU request."
}

// PodCPUTime is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.time" semantic conventions. It represents the total CPU time
// consumed.
//
// Stability: Release_Candidate
type PodCPUTime struct {
	metric.Float64Counter
}

var newPodCPUTimeOpts = []metric.Float64CounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
}

// NewPodCPUTime returns a new PodCPUTime instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodCPUTimeOpts...)
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
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Total CPU time consumed by the specific Pod on all available CPU cores
func (m PodCPUTime) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Counter.Add(ctx, incr, *o...)
}

// PodCPUTimeObservable is an instrument used to record metric values conforming
// to the "k8s.pod.cpu.time" semantic conventions. It represents the total CPU
// time consumed.
//
// Stability: Release_Candidate
type PodCPUTimeObservable struct {
	metric.Float64ObservableCounter
}

var newPodCPUTimeObservableOpts = []metric.Float64ObservableCounterOption{
	metric.WithDescription("Total CPU time consumed."),
	metric.WithUnit("s"),
}

// NewPodCPUTimeObservable returns a new PodCPUTimeObservable instrument.
//
// Stability: Release_Candidate
func NewPodCPUTimeObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableCounterOption,
) (PodCPUTimeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPUTimeObservable{noop.Float64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPUTimeObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPUTimeObservableOpts...)
	}

	i, err := m.Float64ObservableCounter(
		"k8s.pod.cpu.time",
		opt...,
	)
	if err != nil {
		return PodCPUTimeObservable{noop.Float64ObservableCounter{}}, err
	}
	return PodCPUTimeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPUTimeObservable) Inst() metric.Float64ObservableCounter {
	return m.Float64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (PodCPUTimeObservable) Name() string {
	return "k8s.pod.cpu.time"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPUTimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (PodCPUTimeObservable) Description() string {
	return "Total CPU time consumed."
}

// PodCPUUsage is an instrument used to record metric values conforming to the
// "k8s.pod.cpu.usage" semantic conventions. It represents the pod's CPU usage,
// measured in CPUs. Range from 0 to the number of allocatable CPUs.
//
// Stability: Release_Candidate
type PodCPUUsage struct {
	metric.Float64Gauge
}

var newPodCPUUsageOpts = []metric.Float64GaugeOption{
	metric.WithDescription("Pod's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPUUsage returns a new PodCPUUsage instrument.
//
// Stability: Release_Candidate
func NewPodCPUUsage(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (PodCPUUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPUUsage{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPUUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPUUsageOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.pod.cpu.usage",
		opt...,
	)
	if err != nil {
		return PodCPUUsage{noop.Float64Gauge{}}, err
	}
	return PodCPUUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPUUsage) Inst() metric.Float64Gauge {
	return m.Float64Gauge
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
	return "Pod's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."
}

// Record records val to the current distribution for attrs.
//
// CPU usage of the specific Pod on all available CPU cores. It is calculated as
// the change in cumulative CPU time (k8s.pod.cpu.time) over a measurement
// interval, divided by the elapsed time: usageCores = (cpuTimeEnd -
// cpuTimeStart) / elapsedSeconds
func (m PodCPUUsage) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// CPU usage of the specific Pod on all available CPU cores. It is calculated as
// the change in cumulative CPU time (k8s.pod.cpu.time) over a measurement
// interval, divided by the elapsed time: usageCores = (cpuTimeEnd -
// cpuTimeStart) / elapsedSeconds
func (m PodCPUUsage) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodCPUUsageObservable is an instrument used to record metric values conforming
// to the "k8s.pod.cpu.usage" semantic conventions. It represents the pod's CPU
// usage, measured in CPUs. Range from 0 to the number of allocatable CPUs.
//
// Stability: Release_Candidate
type PodCPUUsageObservable struct {
	metric.Float64ObservableGauge
}

var newPodCPUUsageObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("Pod's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."),
	metric.WithUnit("{cpu}"),
}

// NewPodCPUUsageObservable returns a new PodCPUUsageObservable instrument.
//
// Stability: Release_Candidate
func NewPodCPUUsageObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (PodCPUUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodCPUUsageObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodCPUUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodCPUUsageObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.pod.cpu.usage",
		opt...,
	)
	if err != nil {
		return PodCPUUsageObservable{noop.Float64ObservableGauge{}}, err
	}
	return PodCPUUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodCPUUsageObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PodCPUUsageObservable) Name() string {
	return "k8s.pod.cpu.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodCPUUsageObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (PodCPUUsageObservable) Description() string {
	return "Pod's CPU usage, measured in CPUs. Range from 0 to the number of allocatable CPUs."
}

// PodFilesystemAvailable is an instrument used to record metric values
// conforming to the "k8s.pod.filesystem.available" semantic conventions. It
// represents the pod filesystem available bytes.
//
// Stability: Release_Candidate
type PodFilesystemAvailable struct {
	metric.Int64UpDownCounter
}

var newPodFilesystemAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewPodFilesystemAvailable returns a new PodFilesystemAvailable instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodFilesystemAvailableOpts...)
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
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [FsStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodFilesystemAvailableObservable is an instrument used to record metric values
// conforming to the "k8s.pod.filesystem.available" semantic conventions. It
// represents the pod filesystem available bytes.
//
// Stability: Release_Candidate
type PodFilesystemAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodFilesystemAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod filesystem available bytes."),
	metric.WithUnit("By"),
}

// NewPodFilesystemAvailableObservable returns a new
// PodFilesystemAvailableObservable instrument.
//
// Stability: Release_Candidate
func NewPodFilesystemAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodFilesystemAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodFilesystemAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodFilesystemAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.filesystem.available",
		opt...,
	)
	if err != nil {
		return PodFilesystemAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodFilesystemAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodFilesystemAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodFilesystemAvailableObservable) Name() string {
	return "k8s.pod.filesystem.available"
}

// Unit returns the semantic convention unit of the instrument
func (PodFilesystemAvailableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodFilesystemAvailableObservable) Description() string {
	return "Pod filesystem available bytes."
}

// PodFilesystemCapacity is an instrument used to record metric values conforming
// to the "k8s.pod.filesystem.capacity" semantic conventions. It represents the
// pod filesystem capacity.
//
// Stability: Release_Candidate
type PodFilesystemCapacity struct {
	metric.Int64UpDownCounter
}

var newPodFilesystemCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod filesystem capacity."),
	metric.WithUnit("By"),
}

// NewPodFilesystemCapacity returns a new PodFilesystemCapacity instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodFilesystemCapacityOpts...)
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
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemCapacity) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [FsStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodFilesystemCapacityObservable is an instrument used to record metric values
// conforming to the "k8s.pod.filesystem.capacity" semantic conventions. It
// represents the pod filesystem capacity.
//
// Stability: Release_Candidate
type PodFilesystemCapacityObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodFilesystemCapacityObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod filesystem capacity."),
	metric.WithUnit("By"),
}

// NewPodFilesystemCapacityObservable returns a new
// PodFilesystemCapacityObservable instrument.
//
// Stability: Release_Candidate
func NewPodFilesystemCapacityObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodFilesystemCapacityObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemCapacityObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodFilesystemCapacityObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodFilesystemCapacityObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.filesystem.capacity",
		opt...,
	)
	if err != nil {
		return PodFilesystemCapacityObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodFilesystemCapacityObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodFilesystemCapacityObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodFilesystemCapacityObservable) Name() string {
	return "k8s.pod.filesystem.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PodFilesystemCapacityObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodFilesystemCapacityObservable) Description() string {
	return "Pod filesystem capacity."
}

// PodFilesystemUsage is an instrument used to record metric values conforming to
// the "k8s.pod.filesystem.usage" semantic conventions. It represents the pod
// filesystem usage.
//
// Stability: Release_Candidate
type PodFilesystemUsage struct {
	metric.Int64UpDownCounter
}

var newPodFilesystemUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod filesystem usage."),
	metric.WithUnit("By"),
}

// NewPodFilesystemUsage returns a new PodFilesystemUsage instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodFilesystemUsageOpts...)
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
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// Capacity and available are reported for the backing filesystem, which may be
// shared
// with other workloads; therefore, usage may not equal `capacity - available`.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is derived from the
// [FsStats.UsedBytes] field
// of the [PodStats.EphemeralStorage]
// of the Kubelet's stats API.
//
// Capacity and available are reported for the backing filesystem, which may be
// shared
// with other workloads; therefore, usage may not equal `capacity - available`.
//
// [FsStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#FsStats
// [PodStats.EphemeralStorage]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodFilesystemUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodFilesystemUsageObservable is an instrument used to record metric values
// conforming to the "k8s.pod.filesystem.usage" semantic conventions. It
// represents the pod filesystem usage.
//
// Stability: Release_Candidate
type PodFilesystemUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodFilesystemUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod filesystem usage."),
	metric.WithUnit("By"),
}

// NewPodFilesystemUsageObservable returns a new PodFilesystemUsageObservable
// instrument.
//
// Stability: Release_Candidate
func NewPodFilesystemUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodFilesystemUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodFilesystemUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodFilesystemUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.filesystem.usage",
		opt...,
	)
	if err != nil {
		return PodFilesystemUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodFilesystemUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodFilesystemUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodFilesystemUsageObservable) Name() string {
	return "k8s.pod.filesystem.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodFilesystemUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodFilesystemUsageObservable) Description() string {
	return "Pod filesystem usage."
}

// PodMemoryAvailable is an instrument used to record metric values conforming to
// the "k8s.pod.memory.available" semantic conventions. It represents the pod
// memory available.
//
// Stability: Release_Candidate
type PodMemoryAvailable struct {
	metric.Int64UpDownCounter
}

var newPodMemoryAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod memory available."),
	metric.WithUnit("By"),
}

// NewPodMemoryAvailable returns a new PodMemoryAvailable instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodMemoryAvailableOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryAvailableObservable is an instrument used to record metric values
// conforming to the "k8s.pod.memory.available" semantic conventions. It
// represents the pod memory available.
//
// Stability: Release_Candidate
type PodMemoryAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod memory available."),
	metric.WithUnit("By"),
}

// NewPodMemoryAvailableObservable returns a new PodMemoryAvailableObservable
// instrument.
//
// Stability: Release_Candidate
func NewPodMemoryAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.available",
		opt...,
	)
	if err != nil {
		return PodMemoryAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryAvailableObservable) Name() string {
	return "k8s.pod.memory.available"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryAvailableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryAvailableObservable) Description() string {
	return "Pod memory available."
}

// PodMemoryLimitCurrent is an instrument used to record metric values conforming
// to the "k8s.pod.memory.limit.current" semantic conventions. It represents the
// maximum memory resource limit currently configured for a running Pod.
//
// Stability: Development
type PodMemoryLimitCurrent struct {
	metric.Int64UpDownCounter
}

var newPodMemoryLimitCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit currently configured for a running Pod."),
	metric.WithUnit("By"),
}

// NewPodMemoryLimitCurrent returns a new PodMemoryLimitCurrent instrument.
//
// Stability: Development
func NewPodMemoryLimitCurrent(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryLimitCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryLimitCurrent{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryLimitCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryLimitCurrentOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.limit.current",
		opt...,
	)
	if err != nil {
		return PodMemoryLimitCurrent{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryLimitCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryLimitCurrent) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryLimitCurrent) Name() string {
	return "k8s.pod.memory.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryLimitCurrent) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryLimitCurrent) Description() string {
	return "Maximum memory resource limit currently configured for a running Pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryLimitCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryLimitCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryLimitCurrentObservable is an instrument used to record metric values
// conforming to the "k8s.pod.memory.limit.current" semantic conventions. It
// represents the maximum memory resource limit currently configured for a
// running Pod.
//
// Stability: Development
type PodMemoryLimitCurrentObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryLimitCurrentObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit currently configured for a running Pod."),
	metric.WithUnit("By"),
}

// NewPodMemoryLimitCurrentObservable returns a new
// PodMemoryLimitCurrentObservable instrument.
//
// Stability: Development
func NewPodMemoryLimitCurrentObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryLimitCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryLimitCurrentObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryLimitCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryLimitCurrentObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.limit.current",
		opt...,
	)
	if err != nil {
		return PodMemoryLimitCurrentObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryLimitCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryLimitCurrentObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryLimitCurrentObservable) Name() string {
	return "k8s.pod.memory.limit.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryLimitCurrentObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryLimitCurrentObservable) Description() string {
	return "Maximum memory resource limit currently configured for a running Pod."
}

// PodMemoryLimitDesired is an instrument used to record metric values conforming
// to the "k8s.pod.memory.limit.desired" semantic conventions. It represents the
// maximum memory resource limit as defined by the Pod spec.
//
// Stability: Development
type PodMemoryLimitDesired struct {
	metric.Int64UpDownCounter
}

var newPodMemoryLimitDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit as defined by the Pod spec."),
	metric.WithUnit("By"),
}

// NewPodMemoryLimitDesired returns a new PodMemoryLimitDesired instrument.
//
// Stability: Development
func NewPodMemoryLimitDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryLimitDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryLimitDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryLimitDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryLimitDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.limit.desired",
		opt...,
	)
	if err != nil {
		return PodMemoryLimitDesired{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryLimitDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryLimitDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryLimitDesired) Name() string {
	return "k8s.pod.memory.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryLimitDesired) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryLimitDesired) Description() string {
	return "Maximum memory resource limit as defined by the Pod spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level limits
// based on the aggregated limits of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryLimitDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the limit in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level limits
// based on the aggregated limits of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryLimitDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryLimitDesiredObservable is an instrument used to record metric values
// conforming to the "k8s.pod.memory.limit.desired" semantic conventions. It
// represents the maximum memory resource limit as defined by the Pod spec.
//
// Stability: Development
type PodMemoryLimitDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryLimitDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Maximum memory resource limit as defined by the Pod spec."),
	metric.WithUnit("By"),
}

// NewPodMemoryLimitDesiredObservable returns a new
// PodMemoryLimitDesiredObservable instrument.
//
// Stability: Development
func NewPodMemoryLimitDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryLimitDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryLimitDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryLimitDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryLimitDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.limit.desired",
		opt...,
	)
	if err != nil {
		return PodMemoryLimitDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryLimitDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryLimitDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryLimitDesiredObservable) Name() string {
	return "k8s.pod.memory.limit.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryLimitDesiredObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryLimitDesiredObservable) Description() string {
	return "Maximum memory resource limit as defined by the Pod spec."
}

// PodMemoryLimitUtilization is an instrument used to record metric values
// conforming to the "k8s.pod.memory.limit.utilization" semantic conventions. It
// represents the ratio of Pod memory working set to its current memory limit.
//
// Stability: Development
type PodMemoryLimitUtilization struct {
	metric.Float64Gauge
}

var newPodMemoryLimitUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of Pod memory working set to its current memory limit."),
	metric.WithUnit("1"),
}

// NewPodMemoryLimitUtilization returns a new PodMemoryLimitUtilization
// instrument.
//
// Stability: Development
func NewPodMemoryLimitUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (PodMemoryLimitUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryLimitUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryLimitUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryLimitUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.pod.memory.limit.utilization",
		opt...,
	)
	if err != nil {
		return PodMemoryLimitUtilization{noop.Float64Gauge{}}, err
	}
	return PodMemoryLimitUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryLimitUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryLimitUtilization) Name() string {
	return "k8s.pod.memory.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryLimitUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryLimitUtilization) Description() string {
	return "The ratio of Pod memory working set to its current memory limit."
}

// Record records val to the current distribution for attrs.
//
// The numerator is the Pod memory working set, as reported by
// `k8s.pod.memory.working_set`. The current memory limit
// reflects the actual resources applied at the Pod level, as reported by
// [PodStatus].
// A value of 1.0 means the Pod is using 100% of its actual memory limit. Values
// above 1.0 can occur transiently, since memory usage can briefly exceed the
// limit.
// If the memory limit is not set, this metric SHOULD NOT be emitted for that
// Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodMemoryLimitUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The numerator is the Pod memory working set, as reported by
// `k8s.pod.memory.working_set`. The current memory limit
// reflects the actual resources applied at the Pod level, as reported by
// [PodStatus].
// A value of 1.0 means the Pod is using 100% of its actual memory limit. Values
// above 1.0 can occur transiently, since memory usage can briefly exceed the
// limit.
// If the memory limit is not set, this metric SHOULD NOT be emitted for that
// Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodMemoryLimitUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodMemoryLimitUtilizationObservable is an instrument used to record metric
// values conforming to the "k8s.pod.memory.limit.utilization" semantic
// conventions. It represents the ratio of Pod memory working set to its current
// memory limit.
//
// Stability: Development
type PodMemoryLimitUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newPodMemoryLimitUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of Pod memory working set to its current memory limit."),
	metric.WithUnit("1"),
}

// NewPodMemoryLimitUtilizationObservable returns a new
// PodMemoryLimitUtilizationObservable instrument.
//
// Stability: Development
func NewPodMemoryLimitUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (PodMemoryLimitUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryLimitUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryLimitUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryLimitUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.pod.memory.limit.utilization",
		opt...,
	)
	if err != nil {
		return PodMemoryLimitUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return PodMemoryLimitUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryLimitUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryLimitUtilizationObservable) Name() string {
	return "k8s.pod.memory.limit.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryLimitUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryLimitUtilizationObservable) Description() string {
	return "The ratio of Pod memory working set to its current memory limit."
}

// PodMemoryRequestCurrent is an instrument used to record metric values
// conforming to the "k8s.pod.memory.request.current" semantic conventions. It
// represents the memory resource request currently configured for a running Pod.
//
// Stability: Development
type PodMemoryRequestCurrent struct {
	metric.Int64UpDownCounter
}

var newPodMemoryRequestCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Memory resource request currently configured for a running Pod."),
	metric.WithUnit("By"),
}

// NewPodMemoryRequestCurrent returns a new PodMemoryRequestCurrent instrument.
//
// Stability: Development
func NewPodMemoryRequestCurrent(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryRequestCurrent, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRequestCurrent{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRequestCurrentOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRequestCurrentOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.request.current",
		opt...,
	)
	if err != nil {
		return PodMemoryRequestCurrent{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryRequestCurrent{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRequestCurrent) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRequestCurrent) Name() string {
	return "k8s.pod.memory.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRequestCurrent) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRequestCurrent) Description() string {
	return "Memory resource request currently configured for a running Pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryRequestCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodStatus]
// (status.resources). Also see `Actual Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
//
//	for more details.
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryRequestCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryRequestCurrentObservable is an instrument used to record metric
// values conforming to the "k8s.pod.memory.request.current" semantic
// conventions. It represents the memory resource request currently configured
// for a running Pod.
//
// Stability: Development
type PodMemoryRequestCurrentObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryRequestCurrentObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Memory resource request currently configured for a running Pod."),
	metric.WithUnit("By"),
}

// NewPodMemoryRequestCurrentObservable returns a new
// PodMemoryRequestCurrentObservable instrument.
//
// Stability: Development
func NewPodMemoryRequestCurrentObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryRequestCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRequestCurrentObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRequestCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRequestCurrentObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.request.current",
		opt...,
	)
	if err != nil {
		return PodMemoryRequestCurrentObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryRequestCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRequestCurrentObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRequestCurrentObservable) Name() string {
	return "k8s.pod.memory.request.current"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRequestCurrentObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRequestCurrentObservable) Description() string {
	return "Memory resource request currently configured for a running Pod."
}

// PodMemoryRequestDesired is an instrument used to record metric values
// conforming to the "k8s.pod.memory.request.desired" semantic conventions. It
// represents the memory resource requested as defined by the Pod spec.
//
// Stability: Development
type PodMemoryRequestDesired struct {
	metric.Int64UpDownCounter
}

var newPodMemoryRequestDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Memory resource requested as defined by the Pod spec."),
	metric.WithUnit("By"),
}

// NewPodMemoryRequestDesired returns a new PodMemoryRequestDesired instrument.
//
// Stability: Development
func NewPodMemoryRequestDesired(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryRequestDesired, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRequestDesired{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRequestDesiredOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRequestDesiredOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.request.desired",
		opt...,
	)
	if err != nil {
		return PodMemoryRequestDesired{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryRequestDesired{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRequestDesired) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRequestDesired) Name() string {
	return "k8s.pod.memory.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRequestDesired) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRequestDesired) Description() string {
	return "Memory resource requested as defined by the Pod spec."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level requests
// based on the aggregated requests of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryRequestDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the request in the
// [`resources`] field of
// [K8s PodSpec]
// (spec.resources). See
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]
// for more details on Pod level resources. Note that Kubernetes can default
// unset Pod level requests
// based on the aggregated requests of the Pod's containers. Also see
// `Desired Resources` in
//
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]
// .
//
// [`resources`]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#resourcerequirements-v1-core
// [K8s PodSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podspec-v1-core
// [https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/assign-pod-level-resources/
// [https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/]: https://kubernetes.io/docs/tasks/configure-pod-container/resize-container-resources/
func (m PodMemoryRequestDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryRequestDesiredObservable is an instrument used to record metric
// values conforming to the "k8s.pod.memory.request.desired" semantic
// conventions. It represents the memory resource requested as defined by the Pod
// spec.
//
// Stability: Development
type PodMemoryRequestDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryRequestDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Memory resource requested as defined by the Pod spec."),
	metric.WithUnit("By"),
}

// NewPodMemoryRequestDesiredObservable returns a new
// PodMemoryRequestDesiredObservable instrument.
//
// Stability: Development
func NewPodMemoryRequestDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryRequestDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRequestDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRequestDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRequestDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.request.desired",
		opt...,
	)
	if err != nil {
		return PodMemoryRequestDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryRequestDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRequestDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRequestDesiredObservable) Name() string {
	return "k8s.pod.memory.request.desired"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRequestDesiredObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRequestDesiredObservable) Description() string {
	return "Memory resource requested as defined by the Pod spec."
}

// PodMemoryRequestUtilization is an instrument used to record metric values
// conforming to the "k8s.pod.memory.request.utilization" semantic conventions.
// It represents the ratio of Pod memory working set to its current memory
// request.
//
// Stability: Development
type PodMemoryRequestUtilization struct {
	metric.Float64Gauge
}

var newPodMemoryRequestUtilizationOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The ratio of Pod memory working set to its current memory request."),
	metric.WithUnit("1"),
}

// NewPodMemoryRequestUtilization returns a new PodMemoryRequestUtilization
// instrument.
//
// Stability: Development
func NewPodMemoryRequestUtilization(
	m metric.Meter,
	opt ...metric.Float64GaugeOption,
) (PodMemoryRequestUtilization, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRequestUtilization{noop.Float64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRequestUtilizationOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRequestUtilizationOpts...)
	}

	i, err := m.Float64Gauge(
		"k8s.pod.memory.request.utilization",
		opt...,
	)
	if err != nil {
		return PodMemoryRequestUtilization{noop.Float64Gauge{}}, err
	}
	return PodMemoryRequestUtilization{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRequestUtilization) Inst() metric.Float64Gauge {
	return m.Float64Gauge
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRequestUtilization) Name() string {
	return "k8s.pod.memory.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRequestUtilization) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRequestUtilization) Description() string {
	return "The ratio of Pod memory working set to its current memory request."
}

// Record records val to the current distribution for attrs.
//
// The numerator is the Pod memory working set, as reported by
// `k8s.pod.memory.working_set`. The current memory request
// reflects the request applied at the Pod level, as reported by
// [PodStatus].
// A value of 1.0 means the Pod is using exactly its memory request. Since a
// memory request is a
// scheduling reservation and not a hard cap, values above 1.0 are possible and
// indicate the Pod is using
// more memory than requested. If the memory request is not set, this metric
// SHOULD NOT be emitted for that Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodMemoryRequestUtilization) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// The numerator is the Pod memory working set, as reported by
// `k8s.pod.memory.working_set`. The current memory request
// reflects the request applied at the Pod level, as reported by
// [PodStatus].
// A value of 1.0 means the Pod is using exactly its memory request. Since a
// memory request is a
// scheduling reservation and not a hard cap, values above 1.0 are possible and
// indicate the Pod is using
// more memory than requested. If the memory request is not set, this metric
// SHOULD NOT be emitted for that Pod.
//
// [PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.35/#podstatus-v1-core
func (m PodMemoryRequestUtilization) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodMemoryRequestUtilizationObservable is an instrument used to record metric
// values conforming to the "k8s.pod.memory.request.utilization" semantic
// conventions. It represents the ratio of Pod memory working set to its current
// memory request.
//
// Stability: Development
type PodMemoryRequestUtilizationObservable struct {
	metric.Float64ObservableGauge
}

var newPodMemoryRequestUtilizationObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The ratio of Pod memory working set to its current memory request."),
	metric.WithUnit("1"),
}

// NewPodMemoryRequestUtilizationObservable returns a new
// PodMemoryRequestUtilizationObservable instrument.
//
// Stability: Development
func NewPodMemoryRequestUtilizationObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (PodMemoryRequestUtilizationObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRequestUtilizationObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRequestUtilizationObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRequestUtilizationObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.pod.memory.request.utilization",
		opt...,
	)
	if err != nil {
		return PodMemoryRequestUtilizationObservable{noop.Float64ObservableGauge{}}, err
	}
	return PodMemoryRequestUtilizationObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRequestUtilizationObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRequestUtilizationObservable) Name() string {
	return "k8s.pod.memory.request.utilization"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRequestUtilizationObservable) Unit() string {
	return "1"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRequestUtilizationObservable) Description() string {
	return "The ratio of Pod memory working set to its current memory request."
}

// PodMemoryRss is an instrument used to record metric values conforming to the
// "k8s.pod.memory.rss" semantic conventions. It represents the pod memory RSS.
//
// Stability: Release_Candidate
type PodMemoryRss struct {
	metric.Int64UpDownCounter
}

var newPodMemoryRssOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod memory RSS."),
	metric.WithUnit("By"),
}

// NewPodMemoryRss returns a new PodMemoryRss instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRssOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryRssObservable is an instrument used to record metric values
// conforming to the "k8s.pod.memory.rss" semantic conventions. It represents the
// pod memory RSS.
//
// Stability: Release_Candidate
type PodMemoryRssObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryRssObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod memory RSS."),
	metric.WithUnit("By"),
}

// NewPodMemoryRssObservable returns a new PodMemoryRssObservable instrument.
//
// Stability: Release_Candidate
func NewPodMemoryRssObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryRssObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryRssObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryRssObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryRssObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.rss",
		opt...,
	)
	if err != nil {
		return PodMemoryRssObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryRssObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryRssObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryRssObservable) Name() string {
	return "k8s.pod.memory.rss"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryRssObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryRssObservable) Description() string {
	return "Pod memory RSS."
}

// PodMemoryUsage is an instrument used to record metric values conforming to the
// "k8s.pod.memory.usage" semantic conventions. It represents the memory usage of
// the Pod.
//
// Stability: Release_Candidate
type PodMemoryUsage struct {
	metric.Int64UpDownCounter
}

var newPodMemoryUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Memory usage of the Pod."),
	metric.WithUnit("By"),
}

// NewPodMemoryUsage returns a new PodMemoryUsage instrument.
//
// Stability: Release_Candidate
func NewPodMemoryUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodMemoryUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryUsage{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryUsageOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryUsageOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.memory.usage",
		opt...,
	)
	if err != nil {
		return PodMemoryUsage{noop.Int64UpDownCounter{}}, err
	}
	return PodMemoryUsage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryUsage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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

// Add adds incr to the existing count for attrs.
//
// Total memory in use. This includes all memory regardless of when it was
// accessed.
// This metric is derived from the [MemoryStats.UsageBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.UsageBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryUsage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Total memory in use. This includes all memory regardless of when it was
// accessed.
// This metric is derived from the [MemoryStats.UsageBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.UsageBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryUsageObservable is an instrument used to record metric values
// conforming to the "k8s.pod.memory.usage" semantic conventions. It represents
// the memory usage of the Pod.
//
// Stability: Release_Candidate
type PodMemoryUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Memory usage of the Pod."),
	metric.WithUnit("By"),
}

// NewPodMemoryUsageObservable returns a new PodMemoryUsageObservable instrument.
//
// Stability: Release_Candidate
func NewPodMemoryUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.usage",
		opt...,
	)
	if err != nil {
		return PodMemoryUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryUsageObservable) Name() string {
	return "k8s.pod.memory.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryUsageObservable) Description() string {
	return "Memory usage of the Pod."
}

// PodMemoryWorkingSet is an instrument used to record metric values conforming
// to the "k8s.pod.memory.working_set" semantic conventions. It represents the
// pod memory working set.
//
// Stability: Release_Candidate
type PodMemoryWorkingSet struct {
	metric.Int64UpDownCounter
}

var newPodMemoryWorkingSetOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod memory working set."),
	metric.WithUnit("By"),
}

// NewPodMemoryWorkingSet returns a new PodMemoryWorkingSet instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodMemoryWorkingSetOpts...)
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
// dirty memory, and kernel memory. On Linux, the working set is calculated by
// excluding inactive file-backed memory from memory usage (`inactive_file` on
// cgroup v2, `total_inactive_file` on cgroup v1), so active file-backed memory
// remains included. The kubelet uses the Pod working set when ranking Pods for
// eviction under memory pressure. `WorkingSetBytes` is less than or equal to
// `UsageBytes`.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryWorkingSet) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// The amount of working set memory. This includes recently accessed memory,
// dirty memory, and kernel memory. On Linux, the working set is calculated by
// excluding inactive file-backed memory from memory usage (`inactive_file` on
// cgroup v2, `total_inactive_file` on cgroup v1), so active file-backed memory
// remains included. The kubelet uses the Pod working set when ranking Pods for
// eviction under memory pressure. `WorkingSetBytes` is less than or equal to
// `UsageBytes`.
// This metric is derived from the [MemoryStats.WorkingSetBytes] field of the
// [PodStats.Memory] of the Kubelet's stats API.
//
// [MemoryStats.WorkingSetBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodMemoryWorkingSet) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodMemoryWorkingSetObservable is an instrument used to record metric values
// conforming to the "k8s.pod.memory.working_set" semantic conventions. It
// represents the pod memory working set.
//
// Stability: Release_Candidate
type PodMemoryWorkingSetObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodMemoryWorkingSetObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod memory working set."),
	metric.WithUnit("By"),
}

// NewPodMemoryWorkingSetObservable returns a new PodMemoryWorkingSetObservable
// instrument.
//
// Stability: Release_Candidate
func NewPodMemoryWorkingSetObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodMemoryWorkingSetObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodMemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodMemoryWorkingSetObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodMemoryWorkingSetObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.memory.working_set",
		opt...,
	)
	if err != nil {
		return PodMemoryWorkingSetObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodMemoryWorkingSetObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodMemoryWorkingSetObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodMemoryWorkingSetObservable) Name() string {
	return "k8s.pod.memory.working_set"
}

// Unit returns the semantic convention unit of the instrument
func (PodMemoryWorkingSetObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodMemoryWorkingSetObservable) Description() string {
	return "Pod memory working set."
}

// PodNetworkErrors is an instrument used to record metric values conforming to
// the "k8s.pod.network.errors" semantic conventions. It represents the pod
// network errors.
//
// Stability: Release_Candidate
type PodNetworkErrors struct {
	metric.Int64Counter
}

var newPodNetworkErrorsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Pod network errors."),
	metric.WithUnit("{error}"),
}

// NewPodNetworkErrors returns a new PodNetworkErrors instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodNetworkErrorsOpts...)
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (PodNetworkErrors) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// PodNetworkErrorsObservable is an instrument used to record metric values
// conforming to the "k8s.pod.network.errors" semantic conventions. It represents
// the pod network errors.
//
// Stability: Release_Candidate
type PodNetworkErrorsObservable struct {
	metric.Int64ObservableCounter
}

var newPodNetworkErrorsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Pod network errors."),
	metric.WithUnit("{error}"),
}

// NewPodNetworkErrorsObservable returns a new PodNetworkErrorsObservable
// instrument.
//
// Stability: Release_Candidate
func NewPodNetworkErrorsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (PodNetworkErrorsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodNetworkErrorsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodNetworkErrorsObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodNetworkErrorsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"k8s.pod.network.errors",
		opt...,
	)
	if err != nil {
		return PodNetworkErrorsObservable{noop.Int64ObservableCounter{}}, err
	}
	return PodNetworkErrorsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodNetworkErrorsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (PodNetworkErrorsObservable) Name() string {
	return "k8s.pod.network.errors"
}

// Unit returns the semantic convention unit of the instrument
func (PodNetworkErrorsObservable) Unit() string {
	return "{error}"
}

// Description returns the semantic convention description of the instrument
func (PodNetworkErrorsObservable) Description() string {
	return "Pod network errors."
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (PodNetworkErrorsObservable) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (PodNetworkErrorsObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// PodNetworkIO is an instrument used to record metric values conforming to the
// "k8s.pod.network.io" semantic conventions. It represents the network bytes for
// the Pod.
//
// Stability: Release_Candidate
type PodNetworkIO struct {
	metric.Int64Counter
}

var newPodNetworkIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Network bytes for the Pod."),
	metric.WithUnit("By"),
}

// NewPodNetworkIO returns a new PodNetworkIO instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodNetworkIOOpts...)
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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (PodNetworkIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// PodNetworkIOObservable is an instrument used to record metric values
// conforming to the "k8s.pod.network.io" semantic conventions. It represents the
// network bytes for the Pod.
//
// Stability: Release_Candidate
type PodNetworkIOObservable struct {
	metric.Int64ObservableCounter
}

var newPodNetworkIOObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Network bytes for the Pod."),
	metric.WithUnit("By"),
}

// NewPodNetworkIOObservable returns a new PodNetworkIOObservable instrument.
//
// Stability: Release_Candidate
func NewPodNetworkIOObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (PodNetworkIOObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodNetworkIOObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodNetworkIOObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodNetworkIOObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"k8s.pod.network.io",
		opt...,
	)
	if err != nil {
		return PodNetworkIOObservable{noop.Int64ObservableCounter{}}, err
	}
	return PodNetworkIOObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodNetworkIOObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (PodNetworkIOObservable) Name() string {
	return "k8s.pod.network.io"
}

// Unit returns the semantic convention unit of the instrument
func (PodNetworkIOObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodNetworkIOObservable) Description() string {
	return "Network bytes for the Pod."
}

// AttrNetworkInterfaceName returns an optional attribute for the
// "network.interface.name" semantic convention. It represents the network
// interface name.
func (PodNetworkIOObservable) AttrNetworkInterfaceName(val string) attribute.KeyValue {
	return attribute.String("network.interface.name", val)
}

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the direction of
// traffic from the perspective of the observing host's physical or virtual
// network interface. It should not be used to represent the logical direction of
// a stateful connection or network flow.
func (PodNetworkIOObservable) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// PodPagingFaults is an instrument used to record metric values conforming to
// the "k8s.pod.paging.faults" semantic conventions. It represents the pod memory
// paging faults.
//
// Stability: Release_Candidate
type PodPagingFaults struct {
	metric.Int64Counter
}

var newPodPagingFaultsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Pod memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewPodPagingFaults returns a new PodPagingFaults instrument.
//
// Stability: Release_Candidate
func NewPodPagingFaults(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (PodPagingFaults, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodPagingFaults{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodPagingFaultsOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodPagingFaultsOpts...)
	}

	i, err := m.Int64Counter(
		"k8s.pod.paging.faults",
		opt...,
	)
	if err != nil {
		return PodPagingFaults{noop.Int64Counter{}}, err
	}
	return PodPagingFaults{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodPagingFaults) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (PodPagingFaults) Name() string {
	return "k8s.pod.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (PodPagingFaults) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (PodPagingFaults) Description() string {
	return "Pod memory paging faults."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Cumulative number of major and minor page faults.
// The [MemoryStats.MajorPageFaults] field of the [PodStats.Memory] of the
// Kubelet's stats API maps directly to `system.paging.fault.type=major`.
// The [MemoryStats.PageFaults] field reports the total number of page faults
// (equivalent to `pgfault` in the cgroup `memory.stat`), with minor and major
// combined. The minor page faults MUST therefore be derived as
// `MemoryStats.PageFaults - MemoryStats.MajorPageFaults`, so that summing the
// metric over `system.paging.fault.type` yields the total number of page faults.
//
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
func (m PodPagingFaults) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// Cumulative number of major and minor page faults.
// The [MemoryStats.MajorPageFaults] field of the [PodStats.Memory] of the
// Kubelet's stats API maps directly to `system.paging.fault.type=major`.
// The [MemoryStats.PageFaults] field reports the total number of page faults
// (equivalent to `pgfault` in the cgroup `memory.stat`), with minor and major
// combined. The minor page faults MUST therefore be derived as
// `MemoryStats.PageFaults - MemoryStats.MajorPageFaults`, so that summing the
// metric over `system.paging.fault.type` yields the total number of page faults.
//
// [MemoryStats.MajorPageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
// [PodStats.Memory]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
// [MemoryStats.PageFaults]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#MemoryStats
func (m PodPagingFaults) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64Counter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (PodPagingFaults) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// PodPagingFaultsObservable is an instrument used to record metric values
// conforming to the "k8s.pod.paging.faults" semantic conventions. It represents
// the pod memory paging faults.
//
// Stability: Release_Candidate
type PodPagingFaultsObservable struct {
	metric.Int64ObservableCounter
}

var newPodPagingFaultsObservableOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Pod memory paging faults."),
	metric.WithUnit("{fault}"),
}

// NewPodPagingFaultsObservable returns a new PodPagingFaultsObservable
// instrument.
//
// Stability: Release_Candidate
func NewPodPagingFaultsObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (PodPagingFaultsObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodPagingFaultsObservable{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodPagingFaultsObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodPagingFaultsObservableOpts...)
	}

	i, err := m.Int64ObservableCounter(
		"k8s.pod.paging.faults",
		opt...,
	)
	if err != nil {
		return PodPagingFaultsObservable{noop.Int64ObservableCounter{}}, err
	}
	return PodPagingFaultsObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodPagingFaultsObservable) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}

// Name returns the semantic convention name of the instrument.
func (PodPagingFaultsObservable) Name() string {
	return "k8s.pod.paging.faults"
}

// Unit returns the semantic convention unit of the instrument
func (PodPagingFaultsObservable) Unit() string {
	return "{fault}"
}

// Description returns the semantic convention description of the instrument
func (PodPagingFaultsObservable) Description() string {
	return "Pod memory paging faults."
}

// AttrSystemPagingFaultType returns an optional attribute for the
// "system.paging.fault.type" semantic convention. It represents the paging fault
// type.
func (PodPagingFaultsObservable) AttrSystemPagingFaultType(val SystemPagingFaultTypeAttr) attribute.KeyValue {
	return attribute.String("system.paging.fault.type", string(val))
}

// PodStatusPhase is an instrument used to record metric values conforming to the
// "k8s.pod.status.phase" semantic conventions. It represents the describes
// number of K8s Pods that are currently in a given phase.
//
// Stability: Development
type PodStatusPhase struct {
	metric.Int64UpDownCounter
}

var newPodStatusPhaseOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes number of K8s Pods that are currently in a given phase."),
	metric.WithUnit("{pod}"),
}

// NewPodStatusPhase returns a new PodStatusPhase instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodStatusPhaseOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.pod.status.phase", string(podStatusPhase)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodStatusPhaseObservable is an instrument used to record metric values
// conforming to the "k8s.pod.status.phase" semantic conventions. It represents
// the describes number of K8s Pods that are currently in a given phase.
//
// Stability: Development
type PodStatusPhaseObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodStatusPhaseObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Describes number of K8s Pods that are currently in a given phase."),
	metric.WithUnit("{pod}"),
}

// NewPodStatusPhaseObservable returns a new PodStatusPhaseObservable instrument.
//
// Stability: Development
func NewPodStatusPhaseObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodStatusPhaseObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodStatusPhaseObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodStatusPhaseObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodStatusPhaseObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.status.phase",
		opt...,
	)
	if err != nil {
		return PodStatusPhaseObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodStatusPhaseObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodStatusPhaseObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodStatusPhaseObservable) Name() string {
	return "k8s.pod.status.phase"
}

// Unit returns the semantic convention unit of the instrument
func (PodStatusPhaseObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (PodStatusPhaseObservable) Description() string {
	return "Describes number of K8s Pods that are currently in a given phase."
}

// AttrPodStatusPhase returns a required attribute for the "k8s.pod.status.phase"
// semantic convention. It represents the phase for the pod. Corresponds to the
// `phase` field of the: [K8s PodStatus].
//
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#podstatus-v1-core
func (PodStatusPhaseObservable) AttrPodStatusPhase(val PodStatusPhaseAttr) attribute.KeyValue {
	return attribute.String("k8s.pod.status.phase", string(val))
}

// PodStatusReason is an instrument used to record metric values conforming to
// the "k8s.pod.status.reason" semantic conventions. It represents the describes
// the number of K8s Pods that are currently in a state for a given reason.
//
// Stability: Development
type PodStatusReason struct {
	metric.Int64UpDownCounter
}

var newPodStatusReasonOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Describes the number of K8s Pods that are currently in a state for a given reason."),
	metric.WithUnit("{pod}"),
}

// NewPodStatusReason returns a new PodStatusReason instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodStatusReasonOpts...)
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.pod.status.reason", string(podStatusReason)),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// PodStatusReasonObservable is an instrument used to record metric values
// conforming to the "k8s.pod.status.reason" semantic conventions. It represents
// the describes the number of K8s Pods that are currently in a state for a given
// reason.
//
// Stability: Development
type PodStatusReasonObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodStatusReasonObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Describes the number of K8s Pods that are currently in a state for a given reason."),
	metric.WithUnit("{pod}"),
}

// NewPodStatusReasonObservable returns a new PodStatusReasonObservable
// instrument.
//
// Stability: Development
func NewPodStatusReasonObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodStatusReasonObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodStatusReasonObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodStatusReasonObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodStatusReasonObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.status.reason",
		opt...,
	)
	if err != nil {
		return PodStatusReasonObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodStatusReasonObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodStatusReasonObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodStatusReasonObservable) Name() string {
	return "k8s.pod.status.reason"
}

// Unit returns the semantic convention unit of the instrument
func (PodStatusReasonObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (PodStatusReasonObservable) Description() string {
	return "Describes the number of K8s Pods that are currently in a state for a given reason."
}

// AttrPodStatusReason returns a required attribute for the
// "k8s.pod.status.reason" semantic convention. It represents the reason for the
// pod state. Corresponds to the `reason` field of the: [K8s PodStatus].
//
// [K8s PodStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.33/#podstatus-v1-core
func (PodStatusReasonObservable) AttrPodStatusReason(val PodStatusReasonAttr) attribute.KeyValue {
	return attribute.String("k8s.pod.status.reason", string(val))
}

// PodUptime is an instrument used to record metric values conforming to the
// "k8s.pod.uptime" semantic conventions. It represents the time the Pod has been
// running.
//
// Stability: Release_Candidate
type PodUptime struct {
	metric.Float64Gauge
}

var newPodUptimeOpts = []metric.Float64GaugeOption{
	metric.WithDescription("The time the Pod has been running."),
	metric.WithUnit("s"),
}

// NewPodUptime returns a new PodUptime instrument.
//
// Stability: Release_Candidate
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
		opt = append(opt[:len(opt):len(opt)], newPodUptimeOpts...)
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
// seconds as a floating point number with the highest precision available. The
// actual accuracy would depend on the instrumentation and operating system.
// This metric is derived from the [PodStats.StartTime] field of Kubelet's stats
// API, where `StartTime` represents the pod sandbox creation time as
// [noted in the Kubelet source]. This value resets when the node restarts, as
// the pod sandbox is recreated.
//
// [PodStats.StartTime]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L122
// [noted in the Kubelet source]: https://github.com/kubernetes/kubernetes/blob/v1.36.3/pkg/kubelet/stats/cri_stats_provider.go#L299
func (m PodUptime) Record(ctx context.Context, val float64, attrs ...attribute.KeyValue) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// Instrumentations SHOULD use a gauge with type `double` and measure uptime in
// seconds as a floating point number with the highest precision available. The
// actual accuracy would depend on the instrumentation and operating system.
// This metric is derived from the [PodStats.StartTime] field of Kubelet's stats
// API, where `StartTime` represents the pod sandbox creation time as
// [noted in the Kubelet source]. This value resets when the node restarts, as
// the pod sandbox is recreated.
//
// [PodStats.StartTime]: https://github.com/kubernetes/kubelet/blob/v0.35.2/pkg/apis/stats/v1alpha1/types.go#L122
// [noted in the Kubelet source]: https://github.com/kubernetes/kubernetes/blob/v1.36.3/pkg/kubelet/stats/cri_stats_provider.go#L299
func (m PodUptime) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Float64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Gauge.Record(ctx, val, *o...)
}

// PodUptimeObservable is an instrument used to record metric values conforming
// to the "k8s.pod.uptime" semantic conventions. It represents the time the Pod
// has been running.
//
// Stability: Release_Candidate
type PodUptimeObservable struct {
	metric.Float64ObservableGauge
}

var newPodUptimeObservableOpts = []metric.Float64ObservableGaugeOption{
	metric.WithDescription("The time the Pod has been running."),
	metric.WithUnit("s"),
}

// NewPodUptimeObservable returns a new PodUptimeObservable instrument.
//
// Stability: Release_Candidate
func NewPodUptimeObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableGaugeOption,
) (PodUptimeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodUptimeObservable{noop.Float64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodUptimeObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodUptimeObservableOpts...)
	}

	i, err := m.Float64ObservableGauge(
		"k8s.pod.uptime",
		opt...,
	)
	if err != nil {
		return PodUptimeObservable{noop.Float64ObservableGauge{}}, err
	}
	return PodUptimeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodUptimeObservable) Inst() metric.Float64ObservableGauge {
	return m.Float64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (PodUptimeObservable) Name() string {
	return "k8s.pod.uptime"
}

// Unit returns the semantic convention unit of the instrument
func (PodUptimeObservable) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (PodUptimeObservable) Description() string {
	return "The time the Pod has been running."
}

// PodVolumeAvailable is an instrument used to record metric values conforming to
// the "k8s.pod.volume.available" semantic conventions. It represents the pod
// volume storage space available.
//
// Stability: Development
type PodVolumeAvailable struct {
	metric.Int64UpDownCounter
}

var newPodVolumeAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod volume storage space available."),
	metric.WithUnit("By"),
}

// NewPodVolumeAvailable returns a new PodVolumeAvailable instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodVolumeAvailableOpts...)
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
// [VolumeStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeAvailable) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.volume.name", volumeName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [VolumeStats.AvailableBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeAvailable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeAvailableObservable is an instrument used to record metric values
// conforming to the "k8s.pod.volume.available" semantic conventions. It
// represents the pod volume storage space available.
//
// Stability: Development
type PodVolumeAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodVolumeAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod volume storage space available."),
	metric.WithUnit("By"),
}

// NewPodVolumeAvailableObservable returns a new PodVolumeAvailableObservable
// instrument.
//
// Stability: Development
func NewPodVolumeAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodVolumeAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodVolumeAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.volume.available",
		opt...,
	)
	if err != nil {
		return PodVolumeAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodVolumeAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeAvailableObservable) Name() string {
	return "k8s.pod.volume.available"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeAvailableObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeAvailableObservable) Description() string {
	return "Pod volume storage space available."
}

// AttrVolumeName returns a required attribute for the "k8s.volume.name" semantic
// convention. It represents the name of the K8s volume.
func (PodVolumeAvailableObservable) AttrVolumeName(val string) attribute.KeyValue {
	return attribute.String("k8s.volume.name", val)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeAvailableObservable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeCapacity is an instrument used to record metric values conforming to
// the "k8s.pod.volume.capacity" semantic conventions. It represents the pod
// volume total capacity.
//
// Stability: Development
type PodVolumeCapacity struct {
	metric.Int64UpDownCounter
}

var newPodVolumeCapacityOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod volume total capacity."),
	metric.WithUnit("By"),
}

// NewPodVolumeCapacity returns a new PodVolumeCapacity instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodVolumeCapacityOpts...)
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
// [VolumeStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeCapacity) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.volume.name", volumeName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [VolumeStats.CapacityBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeCapacity) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeCapacity) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeCapacityObservable is an instrument used to record metric values
// conforming to the "k8s.pod.volume.capacity" semantic conventions. It
// represents the pod volume total capacity.
//
// Stability: Development
type PodVolumeCapacityObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodVolumeCapacityObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod volume total capacity."),
	metric.WithUnit("By"),
}

// NewPodVolumeCapacityObservable returns a new PodVolumeCapacityObservable
// instrument.
//
// Stability: Development
func NewPodVolumeCapacityObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodVolumeCapacityObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeCapacityObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeCapacityObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodVolumeCapacityObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.volume.capacity",
		opt...,
	)
	if err != nil {
		return PodVolumeCapacityObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodVolumeCapacityObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeCapacityObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeCapacityObservable) Name() string {
	return "k8s.pod.volume.capacity"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeCapacityObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeCapacityObservable) Description() string {
	return "Pod volume total capacity."
}

// AttrVolumeName returns a required attribute for the "k8s.volume.name" semantic
// convention. It represents the name of the K8s volume.
func (PodVolumeCapacityObservable) AttrVolumeName(val string) attribute.KeyValue {
	return attribute.String("k8s.volume.name", val)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeCapacityObservable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeCount is an instrument used to record metric values conforming
// to the "k8s.pod.volume.inode.count" semantic conventions. It represents the
// total inodes in the filesystem of the Pod's volume.
//
// Stability: Development
type PodVolumeInodeCount struct {
	metric.Int64UpDownCounter
}

var newPodVolumeInodeCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The total inodes in the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeCount returns a new PodVolumeInodeCount instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodVolumeInodeCountOpts...)
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
// [VolumeStats.Inodes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeCount) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.volume.name", volumeName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [VolumeStats.Inodes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeCount) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeCountObservable is an instrument used to record metric values
// conforming to the "k8s.pod.volume.inode.count" semantic conventions. It
// represents the total inodes in the filesystem of the Pod's volume.
//
// Stability: Development
type PodVolumeInodeCountObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodVolumeInodeCountObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The total inodes in the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeCountObservable returns a new PodVolumeInodeCountObservable
// instrument.
//
// Stability: Development
func NewPodVolumeInodeCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodVolumeInodeCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeCountObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeInodeCountObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodVolumeInodeCountObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.volume.inode.count",
		opt...,
	)
	if err != nil {
		return PodVolumeInodeCountObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodVolumeInodeCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeInodeCountObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeInodeCountObservable) Name() string {
	return "k8s.pod.volume.inode.count"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeInodeCountObservable) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeInodeCountObservable) Description() string {
	return "The total inodes in the filesystem of the Pod's volume."
}

// AttrVolumeName returns a required attribute for the "k8s.volume.name" semantic
// convention. It represents the name of the K8s volume.
func (PodVolumeInodeCountObservable) AttrVolumeName(val string) attribute.KeyValue {
	return attribute.String("k8s.volume.name", val)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeCountObservable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeFree is an instrument used to record metric values conforming to
// the "k8s.pod.volume.inode.free" semantic conventions. It represents the free
// inodes in the filesystem of the Pod's volume.
//
// Stability: Development
type PodVolumeInodeFree struct {
	metric.Int64UpDownCounter
}

var newPodVolumeInodeFreeOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The free inodes in the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeFree returns a new PodVolumeInodeFree instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodVolumeInodeFreeOpts...)
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
// [VolumeStats.InodesFree]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeFree) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.volume.name", volumeName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [VolumeStats.InodesFree]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeFree) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeFree) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeFreeObservable is an instrument used to record metric values
// conforming to the "k8s.pod.volume.inode.free" semantic conventions. It
// represents the free inodes in the filesystem of the Pod's volume.
//
// Stability: Development
type PodVolumeInodeFreeObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodVolumeInodeFreeObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The free inodes in the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeFreeObservable returns a new PodVolumeInodeFreeObservable
// instrument.
//
// Stability: Development
func NewPodVolumeInodeFreeObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodVolumeInodeFreeObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeFreeObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeInodeFreeObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodVolumeInodeFreeObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.volume.inode.free",
		opt...,
	)
	if err != nil {
		return PodVolumeInodeFreeObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodVolumeInodeFreeObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeInodeFreeObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeInodeFreeObservable) Name() string {
	return "k8s.pod.volume.inode.free"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeInodeFreeObservable) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeInodeFreeObservable) Description() string {
	return "The free inodes in the filesystem of the Pod's volume."
}

// AttrVolumeName returns a required attribute for the "k8s.volume.name" semantic
// convention. It represents the name of the K8s volume.
func (PodVolumeInodeFreeObservable) AttrVolumeName(val string) attribute.KeyValue {
	return attribute.String("k8s.volume.name", val)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeFreeObservable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeUsed is an instrument used to record metric values conforming to
// the "k8s.pod.volume.inode.used" semantic conventions. It represents the inodes
// used by the filesystem of the Pod's volume.
//
// Stability: Development
type PodVolumeInodeUsed struct {
	metric.Int64UpDownCounter
}

var newPodVolumeInodeUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The inodes used by the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeUsed returns a new PodVolumeInodeUsed instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodVolumeInodeUsedOpts...)
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
// [VolumeStats.InodesUsed]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeUsed) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.volume.name", volumeName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [VolumeStats.InodesUsed]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeInodeUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeUsed) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeInodeUsedObservable is an instrument used to record metric values
// conforming to the "k8s.pod.volume.inode.used" semantic conventions. It
// represents the inodes used by the filesystem of the Pod's volume.
//
// Stability: Development
type PodVolumeInodeUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodVolumeInodeUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The inodes used by the filesystem of the Pod's volume."),
	metric.WithUnit("{inode}"),
}

// NewPodVolumeInodeUsedObservable returns a new PodVolumeInodeUsedObservable
// instrument.
//
// Stability: Development
func NewPodVolumeInodeUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodVolumeInodeUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeInodeUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodVolumeInodeUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.volume.inode.used",
		opt...,
	)
	if err != nil {
		return PodVolumeInodeUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodVolumeInodeUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeInodeUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeInodeUsedObservable) Name() string {
	return "k8s.pod.volume.inode.used"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeInodeUsedObservable) Unit() string {
	return "{inode}"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeInodeUsedObservable) Description() string {
	return "The inodes used by the filesystem of the Pod's volume."
}

// AttrVolumeName returns a required attribute for the "k8s.volume.name" semantic
// convention. It represents the name of the K8s volume.
func (PodVolumeInodeUsedObservable) AttrVolumeName(val string) attribute.KeyValue {
	return attribute.String("k8s.volume.name", val)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeInodeUsedObservable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeUsage is an instrument used to record metric values conforming to the
// "k8s.pod.volume.usage" semantic conventions. It represents the pod volume
// usage.
//
// Stability: Development
type PodVolumeUsage struct {
	metric.Int64UpDownCounter
}

var newPodVolumeUsageOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Pod volume usage."),
	metric.WithUnit("By"),
}

// NewPodVolumeUsage returns a new PodVolumeUsage instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newPodVolumeUsageOpts...)
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
// [VolumeStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeUsage) Add(
	ctx context.Context,
	incr int64,
	volumeName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.volume.name", volumeName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [VolumeStats.UsedBytes]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#VolumeStats
// [PodStats]: https://pkg.go.dev/k8s.io/kubelet@v0.34.0/pkg/apis/stats/v1alpha1#PodStats
func (m PodVolumeUsage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeUsage) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// PodVolumeUsageObservable is an instrument used to record metric values
// conforming to the "k8s.pod.volume.usage" semantic conventions. It represents
// the pod volume usage.
//
// Stability: Development
type PodVolumeUsageObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newPodVolumeUsageObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Pod volume usage."),
	metric.WithUnit("By"),
}

// NewPodVolumeUsageObservable returns a new PodVolumeUsageObservable instrument.
//
// Stability: Development
func NewPodVolumeUsageObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (PodVolumeUsageObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeUsageObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newPodVolumeUsageObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newPodVolumeUsageObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.pod.volume.usage",
		opt...,
	)
	if err != nil {
		return PodVolumeUsageObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return PodVolumeUsageObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m PodVolumeUsageObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (PodVolumeUsageObservable) Name() string {
	return "k8s.pod.volume.usage"
}

// Unit returns the semantic convention unit of the instrument
func (PodVolumeUsageObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (PodVolumeUsageObservable) Description() string {
	return "Pod volume usage."
}

// AttrVolumeName returns a required attribute for the "k8s.volume.name" semantic
// convention. It represents the name of the K8s volume.
func (PodVolumeUsageObservable) AttrVolumeName(val string) attribute.KeyValue {
	return attribute.String("k8s.volume.name", val)
}

// AttrVolumeType returns an optional attribute for the "k8s.volume.type"
// semantic convention. It represents the type of the K8s volume.
func (PodVolumeUsageObservable) AttrVolumeType(val VolumeTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.volume.type", string(val))
}

// ReplicaSetPodAvailable is an instrument used to record metric values
// conforming to the "k8s.replicaset.pod.available" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this replicaset.
//
// Stability: Development
type ReplicaSetPodAvailable struct {
	metric.Int64UpDownCounter
}

var newReplicaSetPodAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset."),
	metric.WithUnit("{pod}"),
}

// NewReplicaSetPodAvailable returns a new ReplicaSetPodAvailable instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newReplicaSetPodAvailableOpts...)
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
// [K8s ReplicaSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicasetstatus-v1-apps
func (m ReplicaSetPodAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicaSetStatus].
//
// [K8s ReplicaSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicasetstatus-v1-apps
func (m ReplicaSetPodAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicaSetPodAvailableObservable is an instrument used to record metric values
// conforming to the "k8s.replicaset.pod.available" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this replicaset.
//
// Stability: Development
type ReplicaSetPodAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newReplicaSetPodAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset."),
	metric.WithUnit("{pod}"),
}

// NewReplicaSetPodAvailableObservable returns a new
// ReplicaSetPodAvailableObservable instrument.
//
// Stability: Development
func NewReplicaSetPodAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ReplicaSetPodAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicaSetPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicaSetPodAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newReplicaSetPodAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.replicaset.pod.available",
		opt...,
	)
	if err != nil {
		return ReplicaSetPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ReplicaSetPodAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicaSetPodAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicaSetPodAvailableObservable) Name() string {
	return "k8s.replicaset.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicaSetPodAvailableObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicaSetPodAvailableObservable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset."
}

// ReplicaSetPodDesired is an instrument used to record metric values conforming
// to the "k8s.replicaset.pod.desired" semantic conventions. It represents the
// number of desired replica pods in this replicaset.
//
// Stability: Development
type ReplicaSetPodDesired struct {
	metric.Int64UpDownCounter
}

var newReplicaSetPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this replicaset."),
	metric.WithUnit("{pod}"),
}

// NewReplicaSetPodDesired returns a new ReplicaSetPodDesired instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newReplicaSetPodDesiredOpts...)
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
// [K8s ReplicaSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicasetspec-v1-apps
func (m ReplicaSetPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicaSetSpec].
//
// [K8s ReplicaSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicasetspec-v1-apps
func (m ReplicaSetPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicaSetPodDesiredObservable is an instrument used to record metric values
// conforming to the "k8s.replicaset.pod.desired" semantic conventions. It
// represents the number of desired replica pods in this replicaset.
//
// Stability: Development
type ReplicaSetPodDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newReplicaSetPodDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this replicaset."),
	metric.WithUnit("{pod}"),
}

// NewReplicaSetPodDesiredObservable returns a new ReplicaSetPodDesiredObservable
// instrument.
//
// Stability: Development
func NewReplicaSetPodDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ReplicaSetPodDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicaSetPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicaSetPodDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newReplicaSetPodDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.replicaset.pod.desired",
		opt...,
	)
	if err != nil {
		return ReplicaSetPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ReplicaSetPodDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicaSetPodDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicaSetPodDesiredObservable) Name() string {
	return "k8s.replicaset.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicaSetPodDesiredObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicaSetPodDesiredObservable) Description() string {
	return "Number of desired replica pods in this replicaset."
}

// ReplicationControllerPodAvailable is an instrument used to record metric
// values conforming to the "k8s.replicationcontroller.pod.available" semantic
// conventions. It represents the total number of available replica pods (ready
// for at least minReadySeconds) targeted by this replication controller.
//
// Stability: Development
type ReplicationControllerPodAvailable struct {
	metric.Int64UpDownCounter
}

var newReplicationControllerPodAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller."),
	metric.WithUnit("{pod}"),
}

// NewReplicationControllerPodAvailable returns a new
// ReplicationControllerPodAvailable instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newReplicationControllerPodAvailableOpts...)
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
// [K8s ReplicationControllerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicationcontrollerstatus-v1-core
func (m ReplicationControllerPodAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicationControllerStatus]
//
// [K8s ReplicationControllerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicationcontrollerstatus-v1-core
func (m ReplicationControllerPodAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicationControllerPodAvailableObservable is an instrument used to record
// metric values conforming to the "k8s.replicationcontroller.pod.available"
// semantic conventions. It represents the total number of available replica pods
// (ready for at least minReadySeconds) targeted by this replication controller.
//
// Stability: Development
type ReplicationControllerPodAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newReplicationControllerPodAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller."),
	metric.WithUnit("{pod}"),
}

// NewReplicationControllerPodAvailableObservable returns a new
// ReplicationControllerPodAvailableObservable instrument.
//
// Stability: Development
func NewReplicationControllerPodAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ReplicationControllerPodAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicationControllerPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicationControllerPodAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newReplicationControllerPodAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.replicationcontroller.pod.available",
		opt...,
	)
	if err != nil {
		return ReplicationControllerPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ReplicationControllerPodAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicationControllerPodAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerPodAvailableObservable) Name() string {
	return "k8s.replicationcontroller.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerPodAvailableObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerPodAvailableObservable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller."
}

// ReplicationControllerPodDesired is an instrument used to record metric values
// conforming to the "k8s.replicationcontroller.pod.desired" semantic
// conventions. It represents the number of desired replica pods in this
// replication controller.
//
// Stability: Development
type ReplicationControllerPodDesired struct {
	metric.Int64UpDownCounter
}

var newReplicationControllerPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this replication controller."),
	metric.WithUnit("{pod}"),
}

// NewReplicationControllerPodDesired returns a new
// ReplicationControllerPodDesired instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newReplicationControllerPodDesiredOpts...)
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
// [K8s ReplicationControllerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicationcontrollerspec-v1-core
func (m ReplicationControllerPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicationControllerSpec]
//
// [K8s ReplicationControllerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#replicationcontrollerspec-v1-core
func (m ReplicationControllerPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ReplicationControllerPodDesiredObservable is an instrument used to record
// metric values conforming to the "k8s.replicationcontroller.pod.desired"
// semantic conventions. It represents the number of desired replica pods in this
// replication controller.
//
// Stability: Development
type ReplicationControllerPodDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newReplicationControllerPodDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this replication controller."),
	metric.WithUnit("{pod}"),
}

// NewReplicationControllerPodDesiredObservable returns a new
// ReplicationControllerPodDesiredObservable instrument.
//
// Stability: Development
func NewReplicationControllerPodDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ReplicationControllerPodDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ReplicationControllerPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newReplicationControllerPodDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newReplicationControllerPodDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.replicationcontroller.pod.desired",
		opt...,
	)
	if err != nil {
		return ReplicationControllerPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ReplicationControllerPodDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ReplicationControllerPodDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ReplicationControllerPodDesiredObservable) Name() string {
	return "k8s.replicationcontroller.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (ReplicationControllerPodDesiredObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (ReplicationControllerPodDesiredObservable) Description() string {
	return "Number of desired replica pods in this replication controller."
}

// ResourceQuotaCPULimitHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.limit.hard" semantic conventions. It
// represents the CPU limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPULimitHard struct {
	metric.Float64UpDownCounter
}

var newResourceQuotaCPULimitHardOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPULimitHard returns a new ResourceQuotaCPULimitHard
// instrument.
//
// Stability: Development
func NewResourceQuotaCPULimitHard(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ResourceQuotaCPULimitHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPULimitHard{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPULimitHardOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPULimitHardOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.resourcequota.cpu.limit.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPULimitHard{noop.Float64UpDownCounter{}}, err
	}
	return ResourceQuotaCPULimitHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPULimitHard) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitHard) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitHard) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPULimitHardObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.cpu.limit.hard" semantic
// conventions. It represents the CPU limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPULimitHardObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newResourceQuotaCPULimitHardObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPULimitHardObservable returns a new
// ResourceQuotaCPULimitHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaCPULimitHardObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ResourceQuotaCPULimitHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPULimitHardObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPULimitHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPULimitHardObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.resourcequota.cpu.limit.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPULimitHardObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaCPULimitHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPULimitHardObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPULimitHardObservable) Name() string {
	return "k8s.resourcequota.cpu.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPULimitHardObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPULimitHardObservable) Description() string {
	return "The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// ResourceQuotaCPULimitUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.limit.used" semantic conventions. It
// represents the CPU limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPULimitUsed struct {
	metric.Float64UpDownCounter
}

var newResourceQuotaCPULimitUsedOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPULimitUsed returns a new ResourceQuotaCPULimitUsed
// instrument.
//
// Stability: Development
func NewResourceQuotaCPULimitUsed(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ResourceQuotaCPULimitUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPULimitUsed{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPULimitUsedOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPULimitUsedOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.resourcequota.cpu.limit.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPULimitUsed{noop.Float64UpDownCounter{}}, err
	}
	return ResourceQuotaCPULimitUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPULimitUsed) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitUsed) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPULimitUsed) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPULimitUsedObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.cpu.limit.used" semantic
// conventions. It represents the CPU limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPULimitUsedObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newResourceQuotaCPULimitUsedObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPULimitUsedObservable returns a new
// ResourceQuotaCPULimitUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaCPULimitUsedObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ResourceQuotaCPULimitUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPULimitUsedObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPULimitUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPULimitUsedObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.resourcequota.cpu.limit.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPULimitUsedObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaCPULimitUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPULimitUsedObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPULimitUsedObservable) Name() string {
	return "k8s.resourcequota.cpu.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPULimitUsedObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPULimitUsedObservable) Description() string {
	return "The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// ResourceQuotaCPURequestHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.request.hard" semantic conventions.
// It represents the CPU requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPURequestHard struct {
	metric.Float64UpDownCounter
}

var newResourceQuotaCPURequestHardOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPURequestHard returns a new ResourceQuotaCPURequestHard
// instrument.
//
// Stability: Development
func NewResourceQuotaCPURequestHard(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ResourceQuotaCPURequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPURequestHard{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPURequestHardOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPURequestHardOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.resourcequota.cpu.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPURequestHard{noop.Float64UpDownCounter{}}, err
	}
	return ResourceQuotaCPURequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPURequestHard) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestHard) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestHard) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPURequestHardObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.cpu.request.hard" semantic
// conventions. It represents the CPU requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPURequestHardObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newResourceQuotaCPURequestHardObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPURequestHardObservable returns a new
// ResourceQuotaCPURequestHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaCPURequestHardObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ResourceQuotaCPURequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPURequestHardObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPURequestHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPURequestHardObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.resourcequota.cpu.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPURequestHardObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaCPURequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPURequestHardObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPURequestHardObservable) Name() string {
	return "k8s.resourcequota.cpu.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPURequestHardObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPURequestHardObservable) Description() string {
	return "The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// ResourceQuotaCPURequestUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.cpu.request.used" semantic conventions.
// It represents the CPU requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPURequestUsed struct {
	metric.Float64UpDownCounter
}

var newResourceQuotaCPURequestUsedOpts = []metric.Float64UpDownCounterOption{
	metric.WithDescription("The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPURequestUsed returns a new ResourceQuotaCPURequestUsed
// instrument.
//
// Stability: Development
func NewResourceQuotaCPURequestUsed(
	m metric.Meter,
	opt ...metric.Float64UpDownCounterOption,
) (ResourceQuotaCPURequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPURequestUsed{noop.Float64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPURequestUsedOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPURequestUsedOpts...)
	}

	i, err := m.Float64UpDownCounter(
		"k8s.resourcequota.cpu.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPURequestUsed{noop.Float64UpDownCounter{}}, err
	}
	return ResourceQuotaCPURequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPURequestUsed) Inst() metric.Float64UpDownCounter {
	return m.Float64UpDownCounter
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestUsed) Add(ctx context.Context, incr float64, attrs ...attribute.KeyValue) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaCPURequestUsed) AddSet(ctx context.Context, incr float64, set attribute.Set) {
	if !m.Float64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Float64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaCPURequestUsedObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.cpu.request.used" semantic
// conventions. It represents the CPU requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaCPURequestUsedObservable struct {
	metric.Float64ObservableUpDownCounter
}

var newResourceQuotaCPURequestUsedObservableOpts = []metric.Float64ObservableUpDownCounterOption{
	metric.WithDescription("The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{cpu}"),
}

// NewResourceQuotaCPURequestUsedObservable returns a new
// ResourceQuotaCPURequestUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaCPURequestUsedObservable(
	m metric.Meter,
	opt ...metric.Float64ObservableUpDownCounterOption,
) (ResourceQuotaCPURequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaCPURequestUsedObservable{noop.Float64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaCPURequestUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaCPURequestUsedObservableOpts...)
	}

	i, err := m.Float64ObservableUpDownCounter(
		"k8s.resourcequota.cpu.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaCPURequestUsedObservable{noop.Float64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaCPURequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaCPURequestUsedObservable) Inst() metric.Float64ObservableUpDownCounter {
	return m.Float64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaCPURequestUsedObservable) Name() string {
	return "k8s.resourcequota.cpu.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaCPURequestUsedObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaCPURequestUsedObservable) Description() string {
	return "The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// ResourceQuotaEphemeralStorageLimitHard is an instrument used to record metric
// values conforming to the "k8s.resourcequota.ephemeral_storage.limit.hard"
// semantic conventions. It represents the sum of local ephemeral storage limits
// in the namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageLimitHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageLimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageLimitHard returns a new
// ResourceQuotaEphemeralStorageLimitHard instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageLimitHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageLimitHardObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.limit.hard" semantic conventions. It
// represents the sum of local ephemeral storage limits in the namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageLimitHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaEphemeralStorageLimitHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageLimitHardObservable returns a new
// ResourceQuotaEphemeralStorageLimitHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaEphemeralStorageLimitHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaEphemeralStorageLimitHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageLimitHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageLimitHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.ephemeral_storage.limit.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageLimitHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageLimitHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageLimitHardObservable) Name() string {
	return "k8s.resourcequota.ephemeral_storage.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageLimitHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageLimitHardObservable) Description() string {
	return "The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace."
}

// ResourceQuotaEphemeralStorageLimitUsed is an instrument used to record metric
// values conforming to the "k8s.resourcequota.ephemeral_storage.limit.used"
// semantic conventions. It represents the sum of local ephemeral storage limits
// in the namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageLimitUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageLimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageLimitUsed returns a new
// ResourceQuotaEphemeralStorageLimitUsed instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageLimitUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageLimitUsedObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.limit.used" semantic conventions. It
// represents the sum of local ephemeral storage limits in the namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageLimitUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaEphemeralStorageLimitUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageLimitUsedObservable returns a new
// ResourceQuotaEphemeralStorageLimitUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaEphemeralStorageLimitUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaEphemeralStorageLimitUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageLimitUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageLimitUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.ephemeral_storage.limit.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageLimitUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageLimitUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageLimitUsedObservable) Name() string {
	return "k8s.resourcequota.ephemeral_storage.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageLimitUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageLimitUsedObservable) Description() string {
	return "The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace."
}

// ResourceQuotaEphemeralStorageRequestHard is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.request.hard" semantic conventions. It
// represents the sum of local ephemeral storage requests in the namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageRequestHard returns a new
// ResourceQuotaEphemeralStorageRequestHard instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageRequestHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageRequestHardObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.request.hard" semantic conventions. It
// represents the sum of local ephemeral storage requests in the namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaEphemeralStorageRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageRequestHardObservable returns a new
// ResourceQuotaEphemeralStorageRequestHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaEphemeralStorageRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaEphemeralStorageRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageRequestHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.ephemeral_storage.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageRequestHardObservable) Name() string {
	return "k8s.resourcequota.ephemeral_storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageRequestHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageRequestHardObservable) Description() string {
	return "The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace."
}

// ResourceQuotaEphemeralStorageRequestUsed is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.request.used" semantic conventions. It
// represents the sum of local ephemeral storage requests in the namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaEphemeralStorageRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageRequestUsed returns a new
// ResourceQuotaEphemeralStorageRequestUsed instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageRequestUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaEphemeralStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaEphemeralStorageRequestUsedObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.ephemeral_storage.request.used" semantic conventions. It
// represents the sum of local ephemeral storage requests in the namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaEphemeralStorageRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaEphemeralStorageRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaEphemeralStorageRequestUsedObservable returns a new
// ResourceQuotaEphemeralStorageRequestUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaEphemeralStorageRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaEphemeralStorageRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaEphemeralStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaEphemeralStorageRequestUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaEphemeralStorageRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.ephemeral_storage.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaEphemeralStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaEphemeralStorageRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaEphemeralStorageRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaEphemeralStorageRequestUsedObservable) Name() string {
	return "k8s.resourcequota.ephemeral_storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaEphemeralStorageRequestUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaEphemeralStorageRequestUsedObservable) Description() string {
	return "The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace."
}

// ResourceQuotaHugepageCountRequestHard is an instrument used to record metric
// values conforming to the "k8s.resourcequota.hugepage_count.request.hard"
// semantic conventions. It represents the huge page requests in a specific
// namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaHugepageCountRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaHugepageCountRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{hugepage}"),
}

// NewResourceQuotaHugepageCountRequestHard returns a new
// ResourceQuotaHugepageCountRequestHard instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaHugepageCountRequestHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestHard) Add(
	ctx context.Context,
	incr int64,
	hugepageSize string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.hugepage.size", hugepageSize),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaHugepageCountRequestHardObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.hugepage_count.request.hard" semantic conventions. It
// represents the huge page requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaHugepageCountRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaHugepageCountRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{hugepage}"),
}

// NewResourceQuotaHugepageCountRequestHardObservable returns a new
// ResourceQuotaHugepageCountRequestHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaHugepageCountRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaHugepageCountRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaHugepageCountRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaHugepageCountRequestHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaHugepageCountRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.hugepage_count.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaHugepageCountRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaHugepageCountRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaHugepageCountRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaHugepageCountRequestHardObservable) Name() string {
	return "k8s.resourcequota.hugepage_count.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaHugepageCountRequestHardObservable) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaHugepageCountRequestHardObservable) Description() string {
	return "The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// AttrHugepageSize returns a required attribute for the "k8s.hugepage.size"
// semantic convention. It represents the size (identifier) of the K8s huge page.
func (ResourceQuotaHugepageCountRequestHardObservable) AttrHugepageSize(val string) attribute.KeyValue {
	return attribute.String("k8s.hugepage.size", val)
}

// ResourceQuotaHugepageCountRequestUsed is an instrument used to record metric
// values conforming to the "k8s.resourcequota.hugepage_count.request.used"
// semantic conventions. It represents the huge page requests in a specific
// namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaHugepageCountRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaHugepageCountRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{hugepage}"),
}

// NewResourceQuotaHugepageCountRequestUsed returns a new
// ResourceQuotaHugepageCountRequestUsed instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaHugepageCountRequestUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestUsed) Add(
	ctx context.Context,
	incr int64,
	hugepageSize string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.hugepage.size", hugepageSize),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaHugepageCountRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaHugepageCountRequestUsedObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.hugepage_count.request.used" semantic conventions. It
// represents the huge page requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaHugepageCountRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaHugepageCountRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{hugepage}"),
}

// NewResourceQuotaHugepageCountRequestUsedObservable returns a new
// ResourceQuotaHugepageCountRequestUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaHugepageCountRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaHugepageCountRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaHugepageCountRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaHugepageCountRequestUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaHugepageCountRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.hugepage_count.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaHugepageCountRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaHugepageCountRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaHugepageCountRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaHugepageCountRequestUsedObservable) Name() string {
	return "k8s.resourcequota.hugepage_count.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaHugepageCountRequestUsedObservable) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaHugepageCountRequestUsedObservable) Description() string {
	return "The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// AttrHugepageSize returns a required attribute for the "k8s.hugepage.size"
// semantic convention. It represents the size (identifier) of the K8s huge page.
func (ResourceQuotaHugepageCountRequestUsedObservable) AttrHugepageSize(val string) attribute.KeyValue {
	return attribute.String("k8s.hugepage.size", val)
}

// ResourceQuotaMemoryLimitHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.limit.hard" semantic conventions.
// It represents the memory limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryLimitHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryLimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryLimitHard returns a new ResourceQuotaMemoryLimitHard
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryLimitHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryLimitHardObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.memory.limit.hard" semantic
// conventions. It represents the memory limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryLimitHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaMemoryLimitHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryLimitHardObservable returns a new
// ResourceQuotaMemoryLimitHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaMemoryLimitHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaMemoryLimitHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryLimitHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryLimitHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.memory.limit.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaMemoryLimitHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryLimitHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryLimitHardObservable) Name() string {
	return "k8s.resourcequota.memory.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryLimitHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryLimitHardObservable) Description() string {
	return "The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// ResourceQuotaMemoryLimitUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.limit.used" semantic conventions.
// It represents the memory limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryLimitUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryLimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryLimitUsed returns a new ResourceQuotaMemoryLimitUsed
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryLimitUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryLimitUsedObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.memory.limit.used" semantic
// conventions. It represents the memory limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryLimitUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaMemoryLimitUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryLimitUsedObservable returns a new
// ResourceQuotaMemoryLimitUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaMemoryLimitUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaMemoryLimitUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryLimitUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryLimitUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.memory.limit.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaMemoryLimitUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryLimitUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryLimitUsedObservable) Name() string {
	return "k8s.resourcequota.memory.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryLimitUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryLimitUsedObservable) Description() string {
	return "The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// ResourceQuotaMemoryRequestHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.request.hard" semantic
// conventions. It represents the memory requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryRequestHard returns a new ResourceQuotaMemoryRequestHard
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryRequestHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `hard` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryRequestHardObservable is an instrument used to record
// metric values conforming to the "k8s.resourcequota.memory.request.hard"
// semantic conventions. It represents the memory requests in a specific
// namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaMemoryRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryRequestHardObservable returns a new
// ResourceQuotaMemoryRequestHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaMemoryRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaMemoryRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryRequestHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.memory.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaMemoryRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryRequestHardObservable) Name() string {
	return "k8s.resourcequota.memory.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryRequestHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryRequestHardObservable) Description() string {
	return "The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// ResourceQuotaMemoryRequestUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.memory.request.used" semantic
// conventions. It represents the memory requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaMemoryRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryRequestUsed returns a new ResourceQuotaMemoryRequestUsed
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryRequestUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `used` field of the
// [K8s ResourceQuotaStatus].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaMemoryRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaMemoryRequestUsedObservable is an instrument used to record
// metric values conforming to the "k8s.resourcequota.memory.request.used"
// semantic conventions. It represents the memory requests in a specific
// namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaMemoryRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaMemoryRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaMemoryRequestUsedObservable returns a new
// ResourceQuotaMemoryRequestUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaMemoryRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaMemoryRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaMemoryRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaMemoryRequestUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaMemoryRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.memory.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaMemoryRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaMemoryRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaMemoryRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaMemoryRequestUsedObservable) Name() string {
	return "k8s.resourcequota.memory.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaMemoryRequestUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaMemoryRequestUsedObservable) Description() string {
	return "The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// ResourceQuotaObjectCountHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.object_count.hard" semantic conventions.
// It represents the object count limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaObjectCountHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaObjectCountHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{object}"),
}

// NewResourceQuotaObjectCountHard returns a new ResourceQuotaObjectCountHard
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaObjectCountHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountHard) Add(
	ctx context.Context,
	incr int64,
	resourcequotaResourceName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.resourcequota.resource_name", resourcequotaResourceName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaObjectCountHardObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.object_count.hard" semantic
// conventions. It represents the object count limits in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaObjectCountHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaObjectCountHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{object}"),
}

// NewResourceQuotaObjectCountHardObservable returns a new
// ResourceQuotaObjectCountHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaObjectCountHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaObjectCountHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaObjectCountHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaObjectCountHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaObjectCountHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.object_count.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaObjectCountHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaObjectCountHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaObjectCountHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaObjectCountHardObservable) Name() string {
	return "k8s.resourcequota.object_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaObjectCountHardObservable) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaObjectCountHardObservable) Description() string {
	return "The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// AttrResourceQuotaResourceName returns a required attribute for the
// "k8s.resourcequota.resource_name" semantic convention. It represents the name
// of the K8s resource a resource quota defines.
func (ResourceQuotaObjectCountHardObservable) AttrResourceQuotaResourceName(val string) attribute.KeyValue {
	return attribute.String("k8s.resourcequota.resource_name", val)
}

// ResourceQuotaObjectCountUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.object_count.used" semantic conventions.
// It represents the object count limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaObjectCountUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaObjectCountUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{object}"),
}

// NewResourceQuotaObjectCountUsed returns a new ResourceQuotaObjectCountUsed
// instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaObjectCountUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountUsed) Add(
	ctx context.Context,
	incr int64,
	resourcequotaResourceName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.resourcequota.resource_name", resourcequotaResourceName),
		))
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaObjectCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// ResourceQuotaObjectCountUsedObservable is an instrument used to record metric
// values conforming to the "k8s.resourcequota.object_count.used" semantic
// conventions. It represents the object count limits in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaObjectCountUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaObjectCountUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{object}"),
}

// NewResourceQuotaObjectCountUsedObservable returns a new
// ResourceQuotaObjectCountUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaObjectCountUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaObjectCountUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaObjectCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaObjectCountUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaObjectCountUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.object_count.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaObjectCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaObjectCountUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaObjectCountUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaObjectCountUsedObservable) Name() string {
	return "k8s.resourcequota.object_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaObjectCountUsedObservable) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaObjectCountUsedObservable) Description() string {
	return "The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// AttrResourceQuotaResourceName returns a required attribute for the
// "k8s.resourcequota.resource_name" semantic convention. It represents the name
// of the K8s resource a resource quota defines.
func (ResourceQuotaObjectCountUsedObservable) AttrResourceQuotaResourceName(val string) attribute.KeyValue {
	return attribute.String("k8s.resourcequota.resource_name", val)
}

// ResourceQuotaPersistentvolumeclaimCountHard is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.persistentvolumeclaim_count.hard" semantic conventions. It
// represents the total number of PersistentVolumeClaims that can exist in the
// namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaPersistentvolumeclaimCountHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaPersistentvolumeclaimCountHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewResourceQuotaPersistentvolumeclaimCountHard returns a new
// ResourceQuotaPersistentvolumeclaimCountHard instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaPersistentvolumeclaimCountHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountHard) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaPersistentvolumeclaimCountHard) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaPersistentvolumeclaimCountHardObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.persistentvolumeclaim_count.hard" semantic conventions. It
// represents the total number of PersistentVolumeClaims that can exist in the
// namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaPersistentvolumeclaimCountHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaPersistentvolumeclaimCountHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewResourceQuotaPersistentvolumeclaimCountHardObservable returns a new
// ResourceQuotaPersistentvolumeclaimCountHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaPersistentvolumeclaimCountHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaPersistentvolumeclaimCountHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaPersistentvolumeclaimCountHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaPersistentvolumeclaimCountHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaPersistentvolumeclaimCountHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.persistentvolumeclaim_count.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaPersistentvolumeclaimCountHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaPersistentvolumeclaimCountHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaPersistentvolumeclaimCountHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaPersistentvolumeclaimCountHardObservable) Name() string {
	return "k8s.resourcequota.persistentvolumeclaim_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaPersistentvolumeclaimCountHardObservable) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaPersistentvolumeclaimCountHardObservable) Description() string {
	return "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace."
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaPersistentvolumeclaimCountHardObservable) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaPersistentvolumeclaimCountUsed is an instrument used to record
// metric values conforming to the
// "k8s.resourcequota.persistentvolumeclaim_count.used" semantic conventions. It
// represents the total number of PersistentVolumeClaims that can exist in the
// namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaPersistentvolumeclaimCountUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaPersistentvolumeclaimCountUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewResourceQuotaPersistentvolumeclaimCountUsed returns a new
// ResourceQuotaPersistentvolumeclaimCountUsed instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaPersistentvolumeclaimCountUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountUsed) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaPersistentvolumeclaimCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaPersistentvolumeclaimCountUsed) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaPersistentvolumeclaimCountUsedObservable is an instrument used to
// record metric values conforming to the
// "k8s.resourcequota.persistentvolumeclaim_count.used" semantic conventions. It
// represents the total number of PersistentVolumeClaims that can exist in the
// namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaPersistentvolumeclaimCountUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaPersistentvolumeclaimCountUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewResourceQuotaPersistentvolumeclaimCountUsedObservable returns a new
// ResourceQuotaPersistentvolumeclaimCountUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaPersistentvolumeclaimCountUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaPersistentvolumeclaimCountUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaPersistentvolumeclaimCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaPersistentvolumeclaimCountUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaPersistentvolumeclaimCountUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.persistentvolumeclaim_count.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaPersistentvolumeclaimCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaPersistentvolumeclaimCountUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaPersistentvolumeclaimCountUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaPersistentvolumeclaimCountUsedObservable) Name() string {
	return "k8s.resourcequota.persistentvolumeclaim_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaPersistentvolumeclaimCountUsedObservable) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaPersistentvolumeclaimCountUsedObservable) Description() string {
	return "The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace."
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaPersistentvolumeclaimCountUsedObservable) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaStorageRequestHard is an instrument used to record metric values
// conforming to the "k8s.resourcequota.storage.request.hard" semantic
// conventions. It represents the storage requests in a specific namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaStorageRequestHard struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaStorageRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaStorageRequestHard returns a new
// ResourceQuotaStorageRequestHard instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaStorageRequestHardOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestHard) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaStorageRequestHard) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaStorageRequestHardObservable is an instrument used to record
// metric values conforming to the "k8s.resourcequota.storage.request.hard"
// semantic conventions. It represents the storage requests in a specific
// namespace.
// The value represents the configured quota limit of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaStorageRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaStorageRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaStorageRequestHardObservable returns a new
// ResourceQuotaStorageRequestHardObservable instrument.
//
// Stability: Development
func NewResourceQuotaStorageRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaStorageRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaStorageRequestHardObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaStorageRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.storage.request.hard",
		opt...,
	)
	if err != nil {
		return ResourceQuotaStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaStorageRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaStorageRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaStorageRequestHardObservable) Name() string {
	return "k8s.resourcequota.storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaStorageRequestHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaStorageRequestHardObservable) Description() string {
	return "The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaStorageRequestHardObservable) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaStorageRequestUsed is an instrument used to record metric values
// conforming to the "k8s.resourcequota.storage.request.used" semantic
// conventions. It represents the storage requests in a specific namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaStorageRequestUsed struct {
	metric.Int64UpDownCounter
}

var newResourceQuotaStorageRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaStorageRequestUsed returns a new
// ResourceQuotaStorageRequestUsed instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaStorageRequestUsedOpts...)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestUsed) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
func (m ResourceQuotaStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaStorageRequestUsed) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ResourceQuotaStorageRequestUsedObservable is an instrument used to record
// metric values conforming to the "k8s.resourcequota.storage.request.used"
// semantic conventions. It represents the storage requests in a specific
// namespace.
// The value represents the current observed total usage of the resource in the
// namespace.
//
// Stability: Development
type ResourceQuotaStorageRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newResourceQuotaStorageRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
	metric.WithUnit("By"),
}

// NewResourceQuotaStorageRequestUsedObservable returns a new
// ResourceQuotaStorageRequestUsedObservable instrument.
//
// Stability: Development
func NewResourceQuotaStorageRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ResourceQuotaStorageRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ResourceQuotaStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newResourceQuotaStorageRequestUsedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newResourceQuotaStorageRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.resourcequota.storage.request.used",
		opt...,
	)
	if err != nil {
		return ResourceQuotaStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ResourceQuotaStorageRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ResourceQuotaStorageRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ResourceQuotaStorageRequestUsedObservable) Name() string {
	return "k8s.resourcequota.storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ResourceQuotaStorageRequestUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ResourceQuotaStorageRequestUsedObservable) Description() string {
	return "The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."
}

// AttrStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ResourceQuotaStorageRequestUsedObservable) AttrStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ServiceEndpointCount is an instrument used to record metric values conforming
// to the "k8s.service.endpoint.count" semantic conventions. It represents the
// number of endpoints for a service by condition and address type.
//
// Stability: Development
type ServiceEndpointCount struct {
	metric.Int64Gauge
}

var newServiceEndpointCountOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Number of endpoints for a service by condition and address type."),
	metric.WithUnit("{endpoint}"),
}

// NewServiceEndpointCount returns a new ServiceEndpointCount instrument.
//
// Stability: Development
func NewServiceEndpointCount(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (ServiceEndpointCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServiceEndpointCount{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newServiceEndpointCountOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newServiceEndpointCountOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.service.endpoint.count",
		opt...,
	)
	if err != nil {
		return ServiceEndpointCount{noop.Int64Gauge{}}, err
	}
	return ServiceEndpointCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServiceEndpointCount) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ServiceEndpointCount) Name() string {
	return "k8s.service.endpoint.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServiceEndpointCount) Unit() string {
	return "{endpoint}"
}

// Description returns the semantic convention description of the instrument
func (ServiceEndpointCount) Description() string {
	return "Number of endpoints for a service by condition and address type."
}

// Record records val to the current distribution for attrs.
//
// The serviceEndpointAddressType is the the address type of the service
// endpoint.
//
// The serviceEndpointCondition is the the condition of the service endpoint.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is derived from the Kubernetes [EndpointSlice API].
// It reports the number of network endpoints backing a Service, broken down by
// their condition and address type.
//
// In dual-stack or multi-protocol clusters, separate counts are reported for
// each address family (`IPv4`, `IPv6`, `FQDN`).
//
// When the optional `zone` attribute is enabled, counts are further broken down
// by availability zone for zone-aware monitoring.
//
// An endpoint may be reported under multiple conditions simultaneously (e.g.,
// both `serving` and `terminating` during a graceful shutdown).
// See [K8s EndpointConditions] for more details.
//
// The conditions represent:
//
//   - `ready`: Endpoints capable of receiving new connections.
//   - `serving`: Endpoints currently handling traffic.
//   - `terminating`: Endpoints that are being phased out but may still be
//     handling existing connections.
//
// For Services with `publishNotReadyAddresses` enabled (common for headless
// StatefulSets),
// this metric will include endpoints that are published despite not being ready.
// The `k8s.service.publish_not_ready_addresses` resource attribute indicates
// this setting.
//
// [EndpointSlice API]: https://kubernetes.io/docs/reference/kubernetes-api/service-resources/endpoint-slice-v1/
// [K8s EndpointConditions]: https://kubernetes.io/docs/reference/kubernetes-api/service-resources/endpoint-slice-v1/
func (m ServiceEndpointCount) Record(
	ctx context.Context,
	val int64,
	serviceEndpointAddressType ServiceEndpointAddressTypeAttr,
	serviceEndpointCondition ServiceEndpointConditionAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val, metric.WithAttributes(
			attribute.String("k8s.service.endpoint.address_type", string(serviceEndpointAddressType)),
			attribute.String("k8s.service.endpoint.condition", string(serviceEndpointCondition)),
		))
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(
		*o,
		metric.WithAttributes(
			append(
				attrs[:len(attrs):len(attrs)],
				attribute.String("k8s.service.endpoint.address_type", string(serviceEndpointAddressType)),
				attribute.String("k8s.service.endpoint.condition", string(serviceEndpointCondition)),
			)...,
		),
	)

	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// This metric is derived from the Kubernetes [EndpointSlice API].
// It reports the number of network endpoints backing a Service, broken down by
// their condition and address type.
//
// In dual-stack or multi-protocol clusters, separate counts are reported for
// each address family (`IPv4`, `IPv6`, `FQDN`).
//
// When the optional `zone` attribute is enabled, counts are further broken down
// by availability zone for zone-aware monitoring.
//
// An endpoint may be reported under multiple conditions simultaneously (e.g.,
// both `serving` and `terminating` during a graceful shutdown).
// See [K8s EndpointConditions] for more details.
//
// The conditions represent:
//
//   - `ready`: Endpoints capable of receiving new connections.
//   - `serving`: Endpoints currently handling traffic.
//   - `terminating`: Endpoints that are being phased out but may still be
//     handling existing connections.
//
// For Services with `publishNotReadyAddresses` enabled (common for headless
// StatefulSets),
// this metric will include endpoints that are published despite not being ready.
// The `k8s.service.publish_not_ready_addresses` resource attribute indicates
// this setting.
//
// [EndpointSlice API]: https://kubernetes.io/docs/reference/kubernetes-api/service-resources/endpoint-slice-v1/
// [K8s EndpointConditions]: https://kubernetes.io/docs/reference/kubernetes-api/service-resources/endpoint-slice-v1/
func (m ServiceEndpointCount) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// AttrServiceEndpointZone returns an optional attribute for the
// "k8s.service.endpoint.zone" semantic convention. It represents the zone of the
// service endpoint.
func (ServiceEndpointCount) AttrServiceEndpointZone(val string) attribute.KeyValue {
	return attribute.String("k8s.service.endpoint.zone", val)
}

// ServiceEndpointCountObservable is an instrument used to record metric values
// conforming to the "k8s.service.endpoint.count" semantic conventions. It
// represents the number of endpoints for a service by condition and address
// type.
//
// Stability: Development
type ServiceEndpointCountObservable struct {
	metric.Int64ObservableGauge
}

var newServiceEndpointCountObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Number of endpoints for a service by condition and address type."),
	metric.WithUnit("{endpoint}"),
}

// NewServiceEndpointCountObservable returns a new ServiceEndpointCountObservable
// instrument.
//
// Stability: Development
func NewServiceEndpointCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (ServiceEndpointCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServiceEndpointCountObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newServiceEndpointCountObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newServiceEndpointCountObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"k8s.service.endpoint.count",
		opt...,
	)
	if err != nil {
		return ServiceEndpointCountObservable{noop.Int64ObservableGauge{}}, err
	}
	return ServiceEndpointCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServiceEndpointCountObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (ServiceEndpointCountObservable) Name() string {
	return "k8s.service.endpoint.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServiceEndpointCountObservable) Unit() string {
	return "{endpoint}"
}

// Description returns the semantic convention description of the instrument
func (ServiceEndpointCountObservable) Description() string {
	return "Number of endpoints for a service by condition and address type."
}

// AttrServiceEndpointAddressType returns a required attribute for the
// "k8s.service.endpoint.address_type" semantic convention. It represents the
// address type of the service endpoint.
func (ServiceEndpointCountObservable) AttrServiceEndpointAddressType(val ServiceEndpointAddressTypeAttr) attribute.KeyValue {
	return attribute.String("k8s.service.endpoint.address_type", string(val))
}

// AttrServiceEndpointCondition returns a required attribute for the
// "k8s.service.endpoint.condition" semantic convention. It represents the
// condition of the service endpoint.
func (ServiceEndpointCountObservable) AttrServiceEndpointCondition(val ServiceEndpointConditionAttr) attribute.KeyValue {
	return attribute.String("k8s.service.endpoint.condition", string(val))
}

// AttrServiceEndpointZone returns an optional attribute for the
// "k8s.service.endpoint.zone" semantic convention. It represents the zone of the
// service endpoint.
func (ServiceEndpointCountObservable) AttrServiceEndpointZone(val string) attribute.KeyValue {
	return attribute.String("k8s.service.endpoint.zone", val)
}

// ServiceLoadBalancerIngressCount is an instrument used to record metric values
// conforming to the "k8s.service.load_balancer.ingress.count" semantic
// conventions. It represents the number of load balancer ingress points
// (external IPs/hostnames) assigned to the service.
//
// Stability: Development
type ServiceLoadBalancerIngressCount struct {
	metric.Int64Gauge
}

var newServiceLoadBalancerIngressCountOpts = []metric.Int64GaugeOption{
	metric.WithDescription("Number of load balancer ingress points (external IPs/hostnames) assigned to the service."),
	metric.WithUnit("{ingress}"),
}

// NewServiceLoadBalancerIngressCount returns a new
// ServiceLoadBalancerIngressCount instrument.
//
// Stability: Development
func NewServiceLoadBalancerIngressCount(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (ServiceLoadBalancerIngressCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServiceLoadBalancerIngressCount{noop.Int64Gauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newServiceLoadBalancerIngressCountOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newServiceLoadBalancerIngressCountOpts...)
	}

	i, err := m.Int64Gauge(
		"k8s.service.load_balancer.ingress.count",
		opt...,
	)
	if err != nil {
		return ServiceLoadBalancerIngressCount{noop.Int64Gauge{}}, err
	}
	return ServiceLoadBalancerIngressCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServiceLoadBalancerIngressCount) Inst() metric.Int64Gauge {
	return m.Int64Gauge
}

// Name returns the semantic convention name of the instrument.
func (ServiceLoadBalancerIngressCount) Name() string {
	return "k8s.service.load_balancer.ingress.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServiceLoadBalancerIngressCount) Unit() string {
	return "{ingress}"
}

// Description returns the semantic convention description of the instrument
func (ServiceLoadBalancerIngressCount) Description() string {
	return "Number of load balancer ingress points (external IPs/hostnames) assigned to the service."
}

// Record records val to the current distribution for attrs.
//
// This metric reports the number of external ingress points (IP addresses or
// hostnames)
// assigned to a LoadBalancer Service.
//
// It is only emitted for Services of type `LoadBalancer` and reflects the
// assignments
// made by the underlying infrastructure's load balancer controller in the
// [.status.loadBalancer.ingress] field.
//
// A value of `0` indicates that no ingress points have been assigned yet (e.g.,
// during provisioning).
// A value greater than `1` may occur when multiple IPs or hostnames are assigned
// (e.g., dual-stack configurations).
//
// This metric signals that external endpoints have been assigned by the load
// balancer controller, but it does not
// guarantee that the load balancer is healthy.
//
// [.status.loadBalancer.ingress]: https://kubernetes.io/docs/reference/kubernetes-api/service-resources/service-v1/#ServiceStatus
func (m ServiceLoadBalancerIngressCount) Record(ctx context.Context, val int64, attrs ...attribute.KeyValue) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// This metric reports the number of external ingress points (IP addresses or
// hostnames)
// assigned to a LoadBalancer Service.
//
// It is only emitted for Services of type `LoadBalancer` and reflects the
// assignments
// made by the underlying infrastructure's load balancer controller in the
// [.status.loadBalancer.ingress] field.
//
// A value of `0` indicates that no ingress points have been assigned yet (e.g.,
// during provisioning).
// A value greater than `1` may occur when multiple IPs or hostnames are assigned
// (e.g., dual-stack configurations).
//
// This metric signals that external endpoints have been assigned by the load
// balancer controller, but it does not
// guarantee that the load balancer is healthy.
//
// [.status.loadBalancer.ingress]: https://kubernetes.io/docs/reference/kubernetes-api/service-resources/service-v1/#ServiceStatus
func (m ServiceLoadBalancerIngressCount) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if !m.Int64Gauge.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64Gauge.Record(ctx, val)
		return
	}

	o := metricpool.RecordOptions()
	defer metricpool.PutRecordOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// ServiceLoadBalancerIngressCountObservable is an instrument used to record
// metric values conforming to the "k8s.service.load_balancer.ingress.count"
// semantic conventions. It represents the number of load balancer ingress points
// (external IPs/hostnames) assigned to the service.
//
// Stability: Development
type ServiceLoadBalancerIngressCountObservable struct {
	metric.Int64ObservableGauge
}

var newServiceLoadBalancerIngressCountObservableOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Number of load balancer ingress points (external IPs/hostnames) assigned to the service."),
	metric.WithUnit("{ingress}"),
}

// NewServiceLoadBalancerIngressCountObservable returns a new
// ServiceLoadBalancerIngressCountObservable instrument.
//
// Stability: Development
func NewServiceLoadBalancerIngressCountObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (ServiceLoadBalancerIngressCountObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServiceLoadBalancerIngressCountObservable{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newServiceLoadBalancerIngressCountObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newServiceLoadBalancerIngressCountObservableOpts...)
	}

	i, err := m.Int64ObservableGauge(
		"k8s.service.load_balancer.ingress.count",
		opt...,
	)
	if err != nil {
		return ServiceLoadBalancerIngressCountObservable{noop.Int64ObservableGauge{}}, err
	}
	return ServiceLoadBalancerIngressCountObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServiceLoadBalancerIngressCountObservable) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}

// Name returns the semantic convention name of the instrument.
func (ServiceLoadBalancerIngressCountObservable) Name() string {
	return "k8s.service.load_balancer.ingress.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServiceLoadBalancerIngressCountObservable) Unit() string {
	return "{ingress}"
}

// Description returns the semantic convention description of the instrument
func (ServiceLoadBalancerIngressCountObservable) Description() string {
	return "Number of load balancer ingress points (external IPs/hostnames) assigned to the service."
}

// StatefulSetPodAvailable is an instrument used to record metric values
// conforming to the "k8s.statefulset.pod.available" semantic conventions. It
// represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this statefulset.
//
// Stability: Development
type StatefulSetPodAvailable struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodAvailableOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodAvailable returns a new StatefulSetPodAvailable instrument.
//
// Stability: Development
func NewStatefulSetPodAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (StatefulSetPodAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodAvailable{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodAvailableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodAvailableOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"k8s.statefulset.pod.available",
		opt...,
	)
	if err != nil {
		return StatefulSetPodAvailable{noop.Int64UpDownCounter{}}, err
	}
	return StatefulSetPodAvailable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodAvailable) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodAvailable) Name() string {
	return "k8s.statefulset.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodAvailable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodAvailable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodAvailable) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodAvailable) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodAvailableObservable is an instrument used to record metric
// values conforming to the "k8s.statefulset.pod.available" semantic conventions.
// It represents the total number of available replica pods (ready for at least
// minReadySeconds) targeted by this statefulset.
//
// Stability: Development
type StatefulSetPodAvailableObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newStatefulSetPodAvailableObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodAvailableObservable returns a new
// StatefulSetPodAvailableObservable instrument.
//
// Stability: Development
func NewStatefulSetPodAvailableObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (StatefulSetPodAvailableObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodAvailableObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodAvailableObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.statefulset.pod.available",
		opt...,
	)
	if err != nil {
		return StatefulSetPodAvailableObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return StatefulSetPodAvailableObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodAvailableObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodAvailableObservable) Name() string {
	return "k8s.statefulset.pod.available"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodAvailableObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodAvailableObservable) Description() string {
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this statefulset."
}

// StatefulSetPodCurrent is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.current" semantic conventions. It represents the
// number of replica pods created by the statefulset controller from the
// statefulset version indicated by currentRevision.
//
// Stability: Development
type StatefulSetPodCurrent struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodCurrentOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodCurrent returns a new StatefulSetPodCurrent instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodCurrentOpts...)
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
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodCurrent) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodCurrent) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodCurrentObservable is an instrument used to record metric values
// conforming to the "k8s.statefulset.pod.current" semantic conventions. It
// represents the number of replica pods created by the statefulset controller
// from the statefulset version indicated by currentRevision.
//
// Stability: Development
type StatefulSetPodCurrentObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newStatefulSetPodCurrentObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodCurrentObservable returns a new
// StatefulSetPodCurrentObservable instrument.
//
// Stability: Development
func NewStatefulSetPodCurrentObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (StatefulSetPodCurrentObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodCurrentObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodCurrentObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodCurrentObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.statefulset.pod.current",
		opt...,
	)
	if err != nil {
		return StatefulSetPodCurrentObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return StatefulSetPodCurrentObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodCurrentObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodCurrentObservable) Name() string {
	return "k8s.statefulset.pod.current"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodCurrentObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodCurrentObservable) Description() string {
	return "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision."
}

// StatefulSetPodDesired is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.desired" semantic conventions. It represents the
// number of desired replica pods in this statefulset.
//
// Stability: Development
type StatefulSetPodDesired struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodDesiredOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this statefulset."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodDesired returns a new StatefulSetPodDesired instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodDesiredOpts...)
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
// [K8s StatefulSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetspec-v1-apps
func (m StatefulSetPodDesired) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s StatefulSetSpec].
//
// [K8s StatefulSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetspec-v1-apps
func (m StatefulSetPodDesired) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodDesiredObservable is an instrument used to record metric values
// conforming to the "k8s.statefulset.pod.desired" semantic conventions. It
// represents the number of desired replica pods in this statefulset.
//
// Stability: Development
type StatefulSetPodDesiredObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newStatefulSetPodDesiredObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of desired replica pods in this statefulset."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodDesiredObservable returns a new
// StatefulSetPodDesiredObservable instrument.
//
// Stability: Development
func NewStatefulSetPodDesiredObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (StatefulSetPodDesiredObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodDesiredObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodDesiredObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.statefulset.pod.desired",
		opt...,
	)
	if err != nil {
		return StatefulSetPodDesiredObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return StatefulSetPodDesiredObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodDesiredObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodDesiredObservable) Name() string {
	return "k8s.statefulset.pod.desired"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodDesiredObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodDesiredObservable) Description() string {
	return "Number of desired replica pods in this statefulset."
}

// StatefulSetPodReady is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.ready" semantic conventions. It represents the
// number of replica pods created for this statefulset with a Ready Condition.
//
// Stability: Development
type StatefulSetPodReady struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodReadyOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The number of replica pods created for this statefulset with a Ready Condition."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodReady returns a new StatefulSetPodReady instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodReadyOpts...)
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
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodReady) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `readyReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodReady) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodReadyObservable is an instrument used to record metric values
// conforming to the "k8s.statefulset.pod.ready" semantic conventions. It
// represents the number of replica pods created for this statefulset with a
// Ready Condition.
//
// Stability: Development
type StatefulSetPodReadyObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newStatefulSetPodReadyObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The number of replica pods created for this statefulset with a Ready Condition."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodReadyObservable returns a new StatefulSetPodReadyObservable
// instrument.
//
// Stability: Development
func NewStatefulSetPodReadyObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (StatefulSetPodReadyObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodReadyObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodReadyObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodReadyObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.statefulset.pod.ready",
		opt...,
	)
	if err != nil {
		return StatefulSetPodReadyObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return StatefulSetPodReadyObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodReadyObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodReadyObservable) Name() string {
	return "k8s.statefulset.pod.ready"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodReadyObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodReadyObservable) Description() string {
	return "The number of replica pods created for this statefulset with a Ready Condition."
}

// StatefulSetPodUpdated is an instrument used to record metric values conforming
// to the "k8s.statefulset.pod.updated" semantic conventions. It represents the
// number of replica pods created by the statefulset controller from the
// statefulset version indicated by updateRevision.
//
// Stability: Development
type StatefulSetPodUpdated struct {
	metric.Int64UpDownCounter
}

var newStatefulSetPodUpdatedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodUpdated returns a new StatefulSetPodUpdated instrument.
//
// Stability: Development
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
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodUpdatedOpts...)
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
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodUpdated) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `updatedReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#statefulsetstatus-v1-apps
func (m StatefulSetPodUpdated) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if set.Len() == 0 {
		m.Int64UpDownCounter.Add(ctx, incr)
		return
	}

	o := metricpool.AddOptions()
	defer metricpool.PutAddOptions(o)

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// StatefulSetPodUpdatedObservable is an instrument used to record metric values
// conforming to the "k8s.statefulset.pod.updated" semantic conventions. It
// represents the number of replica pods created by the statefulset controller
// from the statefulset version indicated by updateRevision.
//
// Stability: Development
type StatefulSetPodUpdatedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newStatefulSetPodUpdatedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision."),
	metric.WithUnit("{pod}"),
}

// NewStatefulSetPodUpdatedObservable returns a new
// StatefulSetPodUpdatedObservable instrument.
//
// Stability: Development
func NewStatefulSetPodUpdatedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (StatefulSetPodUpdatedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return StatefulSetPodUpdatedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStatefulSetPodUpdatedObservableOpts
	} else {
		opt = append(opt[:len(opt):len(opt)], newStatefulSetPodUpdatedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"k8s.statefulset.pod.updated",
		opt...,
	)
	if err != nil {
		return StatefulSetPodUpdatedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return StatefulSetPodUpdatedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StatefulSetPodUpdatedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (StatefulSetPodUpdatedObservable) Name() string {
	return "k8s.statefulset.pod.updated"
}

// Unit returns the semantic convention unit of the instrument
func (StatefulSetPodUpdatedObservable) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (StatefulSetPodUpdatedObservable) Description() string {
	return "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision."
}
