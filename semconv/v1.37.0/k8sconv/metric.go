// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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

// ContainerCPULimit is an instrument used to record metric values conforming to
// the "k8s.container.cpu.limit" semantic conventions. It represents the maximum
// CPU resource limit set for the container.
type ContainerCPULimit struct {
	metric.Int64UpDownCounter
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

	i, err := m.Int64UpDownCounter(
		"k8s.container.cpu.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Maximum CPU resource limit set for the container."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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

// ContainerCPURequest is an instrument used to record metric values conforming
// to the "k8s.container.cpu.request" semantic conventions. It represents the CPU
// resource requested for the container.
type ContainerCPURequest struct {
	metric.Int64UpDownCounter
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

	i, err := m.Int64UpDownCounter(
		"k8s.container.cpu.request",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("CPU resource requested for the container."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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

// ContainerEphemeralStorageLimit is an instrument used to record metric values
// conforming to the "k8s.container.ephemeral_storage.limit" semantic
// conventions. It represents the maximum ephemeral storage resource limit set
// for the container.
type ContainerEphemeralStorageLimit struct {
	metric.Int64UpDownCounter
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

	i, err := m.Int64UpDownCounter(
		"k8s.container.ephemeral_storage.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Maximum ephemeral storage resource limit set for the container."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.container.ephemeral_storage.request",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Ephemeral storage resource requested for the container."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewContainerMemoryLimit returns a new ContainerMemoryLimit instrument.
func NewContainerMemoryLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Maximum memory resource limit set for the container."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewContainerMemoryRequest returns a new ContainerMemoryRequest instrument.
func NewContainerMemoryRequest(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerMemoryRequest, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerMemoryRequest{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.memory.request",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Memory resource requested for the container."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewContainerReady returns a new ContainerReady instrument.
func NewContainerReady(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerReady, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerReady{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.ready",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Indicates whether the container is currently marked as ready to accept traffic, based on its readiness probe (1 = ready, 0 = not ready)."),
			metric.WithUnit("{container}"),
		}, opt...)...,
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

// NewContainerRestartCount returns a new ContainerRestartCount instrument.
func NewContainerRestartCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerRestartCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerRestartCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.restart.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Describes how many times the container has restarted (since the last counter reset)."),
			metric.WithUnit("{restart}"),
		}, opt...)...,
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

// NewContainerStatusReason returns a new ContainerStatusReason instrument.
func NewContainerStatusReason(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStatusReason, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStatusReason{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.status.reason",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Describes the number of K8s containers that are currently in a state for a given reason."),
			metric.WithUnit("{container}"),
		}, opt...)...,
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
// [K8s ContainerStateWaiting]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstatewaiting-v1-core
// [K8s ContainerStateTerminated]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstateterminated-v1-core
//
// All possible container state reasons will be reported at each time interval to
// avoid missing metrics.
// Only the value corresponding to the current state reason will be non-zero.
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

// NewContainerStatusState returns a new ContainerStatusState instrument.
func NewContainerStatusState(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStatusState, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStatusState{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.status.state",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Describes the number of K8s containers that are currently in a given state."),
			metric.WithUnit("{container}"),
		}, opt...)...,
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
// [K8s ContainerState]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#containerstate-v1-core
//
// All possible container states will be reported at each time interval to avoid
// missing metrics.
// Only the value corresponding to the current state will be non-zero.
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

// NewContainerStorageLimit returns a new ContainerStorageLimit instrument.
func NewContainerStorageLimit(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStorageLimit, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStorageLimit{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.storage.limit",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Maximum storage resource limit set for the container."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewContainerStorageRequest returns a new ContainerStorageRequest instrument.
func NewContainerStorageRequest(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ContainerStorageRequest, error) {
	// Check if the meter is nil.
	if m == nil {
		return ContainerStorageRequest{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.container.storage.request",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Storage resource requested for the container."),
			metric.WithUnit("By"),
		}, opt...)...,
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
			metric.WithDescription("The number of actively running jobs for a cronjob."),
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
	return "The number of actively running jobs for a cronjob."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `active` field of the
// [K8s CronJobStatus].
//
// [K8s CronJobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#cronjobstatus-v1-batch
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `active` field of the
// [K8s CronJobStatus].
//
// [K8s CronJobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#cronjobstatus-v1-batch
func (m CronJobActiveJobs) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod."),
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
	return "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `currentNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetCurrentScheduledNodes) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)."),
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
	return "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `desiredNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `desiredNumberScheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetDesiredScheduledNodes) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod."),
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
	return "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `numberMisscheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `numberMisscheduled` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetMisscheduledNodes) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready."),
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
	return "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `numberReady` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `numberReady` field of the
// [K8s DaemonSetStatus].
//
// [K8s DaemonSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#daemonsetstatus-v1-apps
func (m DaemonSetReadyNodes) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment."),
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
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this deployment."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s DeploymentStatus].
//
// [K8s DeploymentStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s DeploymentStatus].
//
// [K8s DeploymentStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentstatus-v1-apps
func (m DeploymentAvailablePods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of desired replica pods in this deployment."),
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
	return "Number of desired replica pods in this deployment."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s DeploymentSpec].
//
// [K8s DeploymentSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentspec-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s DeploymentSpec].
//
// [K8s DeploymentSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#deploymentspec-v1-apps
func (m DeploymentDesiredPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler."),
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
	return "Current number of replica pods managed by this horizontal pod autoscaler, as last seen by the autoscaler."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPACurrentPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler."),
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
	return "Desired number of replica pods managed by this horizontal pod autoscaler, as last calculated by the autoscaler."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `desiredReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `desiredReplicas` field of the
// [K8s HorizontalPodAutoscalerStatus]
//
// [K8s HorizontalPodAutoscalerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerstatus-v2-autoscaling
func (m HPADesiredPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("The upper limit for the number of replica pods to which the autoscaler can scale up."),
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
	return "The upper limit for the number of replica pods to which the autoscaler can scale up."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `maxReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `maxReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAMaxPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

	i, err := m.Int64Gauge(
		"k8s.hpa.metric.target.cpu.average_utilization",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Target average utilization, in percentage, for CPU resource in HPA config."),
			metric.WithUnit("1"),
		}, opt...)...,
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

	i, err := m.Int64Gauge(
		"k8s.hpa.metric.target.cpu.average_value",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Target average value for CPU resource in HPA config."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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

// NewHPAMetricTargetCPUValue returns a new HPAMetricTargetCPUValue instrument.
func NewHPAMetricTargetCPUValue(
	m metric.Meter,
	opt ...metric.Int64GaugeOption,
) (HPAMetricTargetCPUValue, error) {
	// Check if the meter is nil.
	if m == nil {
		return HPAMetricTargetCPUValue{noop.Int64Gauge{}}, nil
	}

	i, err := m.Int64Gauge(
		"k8s.hpa.metric.target.cpu.value",
		append([]metric.Int64GaugeOption{
			metric.WithDescription("Target value for CPU resource in HPA config."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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
			metric.WithDescription("The lower limit for the number of replica pods to which the autoscaler can scale down."),
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
	return "The lower limit for the number of replica pods to which the autoscaler can scale down."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `minReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `minReplicas` field of the
// [K8s HorizontalPodAutoscalerSpec]
//
// [K8s HorizontalPodAutoscalerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#horizontalpodautoscalerspec-v2-autoscaling
func (m HPAMinPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("The number of pending and actively running pods for a job."),
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
	return "The number of pending and actively running pods for a job."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `active` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `active` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobActivePods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("The desired number of successfully finished pods the job should be run with."),
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
	return "The desired number of successfully finished pods the job should be run with."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `completions` field of the
// [K8s JobSpec]..
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `completions` field of the
// [K8s JobSpec]..
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
func (m JobDesiredSuccessfulPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("The number of pods which reached phase Failed for a job."),
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
	return "The number of pods which reached phase Failed for a job."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `failed` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `failed` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobFailedPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("The max desired number of pods the job should run at any given time."),
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
	return "The max desired number of pods the job should run at any given time."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `parallelism` field of the
// [K8s JobSpec].
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `parallelism` field of the
// [K8s JobSpec].
//
// [K8s JobSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobspec-v1-batch
func (m JobMaxParallelPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("The number of pods which reached phase Succeeded for a job."),
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
	return "The number of pods which reached phase Succeeded for a job."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `succeeded` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `succeeded` field of the
// [K8s JobStatus].
//
// [K8s JobStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#jobstatus-v1-batch
func (m JobSuccessfulPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NodeAllocatableCPU is an instrument used to record metric values conforming to
// the "k8s.node.allocatable.cpu" semantic conventions. It represents the amount
// of cpu allocatable on the node.
type NodeAllocatableCPU struct {
	metric.Int64UpDownCounter
}

// NewNodeAllocatableCPU returns a new NodeAllocatableCPU instrument.
func NewNodeAllocatableCPU(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeAllocatableCPU, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeAllocatableCPU{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.allocatable.cpu",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Amount of cpu allocatable on the node."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeAllocatableCPU{noop.Int64UpDownCounter{}}, err
	}
	return NodeAllocatableCPU{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeAllocatableCPU) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeAllocatableCPU) Name() string {
	return "k8s.node.allocatable.cpu"
}

// Unit returns the semantic convention unit of the instrument
func (NodeAllocatableCPU) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (NodeAllocatableCPU) Description() string {
	return "Amount of cpu allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeAllocatableCPU) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
func (m NodeAllocatableCPU) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NodeAllocatableEphemeralStorage is an instrument used to record metric values
// conforming to the "k8s.node.allocatable.ephemeral_storage" semantic
// conventions. It represents the amount of ephemeral-storage allocatable on the
// node.
type NodeAllocatableEphemeralStorage struct {
	metric.Int64UpDownCounter
}

// NewNodeAllocatableEphemeralStorage returns a new
// NodeAllocatableEphemeralStorage instrument.
func NewNodeAllocatableEphemeralStorage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeAllocatableEphemeralStorage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeAllocatableEphemeralStorage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.allocatable.ephemeral_storage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Amount of ephemeral-storage allocatable on the node."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeAllocatableEphemeralStorage{noop.Int64UpDownCounter{}}, err
	}
	return NodeAllocatableEphemeralStorage{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeAllocatableEphemeralStorage) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeAllocatableEphemeralStorage) Name() string {
	return "k8s.node.allocatable.ephemeral_storage"
}

// Unit returns the semantic convention unit of the instrument
func (NodeAllocatableEphemeralStorage) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeAllocatableEphemeralStorage) Description() string {
	return "Amount of ephemeral-storage allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeAllocatableEphemeralStorage) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
func (m NodeAllocatableEphemeralStorage) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NodeAllocatableMemory is an instrument used to record metric values conforming
// to the "k8s.node.allocatable.memory" semantic conventions. It represents the
// amount of memory allocatable on the node.
type NodeAllocatableMemory struct {
	metric.Int64UpDownCounter
}

// NewNodeAllocatableMemory returns a new NodeAllocatableMemory instrument.
func NewNodeAllocatableMemory(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeAllocatableMemory, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeAllocatableMemory{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.allocatable.memory",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Amount of memory allocatable on the node."),
			metric.WithUnit("By"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeAllocatableMemory{noop.Int64UpDownCounter{}}, err
	}
	return NodeAllocatableMemory{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeAllocatableMemory) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeAllocatableMemory) Name() string {
	return "k8s.node.allocatable.memory"
}

// Unit returns the semantic convention unit of the instrument
func (NodeAllocatableMemory) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (NodeAllocatableMemory) Description() string {
	return "Amount of memory allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeAllocatableMemory) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
func (m NodeAllocatableMemory) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NodeAllocatablePods is an instrument used to record metric values conforming
// to the "k8s.node.allocatable.pods" semantic conventions. It represents the
// amount of pods allocatable on the node.
type NodeAllocatablePods struct {
	metric.Int64UpDownCounter
}

// NewNodeAllocatablePods returns a new NodeAllocatablePods instrument.
func NewNodeAllocatablePods(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeAllocatablePods, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeAllocatablePods{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.allocatable.pods",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Amount of pods allocatable on the node."),
			metric.WithUnit("{pod}"),
		}, opt...)...,
	)
	if err != nil {
	    return NodeAllocatablePods{noop.Int64UpDownCounter{}}, err
	}
	return NodeAllocatablePods{i}, nil
}

// Inst returns the underlying metric instrument.
func (m NodeAllocatablePods) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (NodeAllocatablePods) Name() string {
	return "k8s.node.allocatable.pods"
}

// Unit returns the semantic convention unit of the instrument
func (NodeAllocatablePods) Unit() string {
	return "{pod}"
}

// Description returns the semantic convention description of the instrument
func (NodeAllocatablePods) Description() string {
	return "Amount of pods allocatable on the node."
}

// Add adds incr to the existing count for attrs.
func (m NodeAllocatablePods) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
func (m NodeAllocatablePods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// NewNodeConditionStatus returns a new NodeConditionStatus instrument.
func NewNodeConditionStatus(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeConditionStatus, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeConditionStatus{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.condition.status",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Describes the condition of a particular Node."),
			metric.WithUnit("{node}"),
		}, opt...)...,
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
			metric.WithDescription("Total CPU time consumed."),
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
			metric.WithDescription("Node's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."),
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
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Gauge.Record(ctx, val, *o...)
}

// NodeFilesystemAvailable is an instrument used to record metric values
// conforming to the "k8s.node.filesystem.available" semantic conventions. It
// represents the node filesystem available bytes.
type NodeFilesystemAvailable struct {
	metric.Int64UpDownCounter
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

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.available",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Node filesystem available bytes."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewNodeFilesystemCapacity returns a new NodeFilesystemCapacity instrument.
func NewNodeFilesystemCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeFilesystemCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemCapacity{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.capacity",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Node filesystem capacity."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewNodeFilesystemUsage returns a new NodeFilesystemUsage instrument.
func NewNodeFilesystemUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (NodeFilesystemUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return NodeFilesystemUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.node.filesystem.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Node filesystem usage."),
			metric.WithUnit("By"),
		}, opt...)...,
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
			metric.WithDescription("Memory usage of the Node."),
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
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
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
			metric.WithDescription("Node network errors."),
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
			metric.WithDescription("Network bytes for the Node."),
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
			metric.WithDescription("The time the Node has been running."),
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
			metric.WithDescription("Total CPU time consumed."),
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
			metric.WithDescription("Pod's CPU usage, measured in cpus. Range from 0 to the number of allocatable CPUs."),
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

// NewPodFilesystemAvailable returns a new PodFilesystemAvailable instrument.
func NewPodFilesystemAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodFilesystemAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemAvailable{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.filesystem.available",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Pod filesystem available bytes."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewPodFilesystemCapacity returns a new PodFilesystemCapacity instrument.
func NewPodFilesystemCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodFilesystemCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemCapacity{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.filesystem.capacity",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Pod filesystem capacity."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewPodFilesystemUsage returns a new PodFilesystemUsage instrument.
func NewPodFilesystemUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodFilesystemUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodFilesystemUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.filesystem.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Pod filesystem usage."),
			metric.WithUnit("By"),
		}, opt...)...,
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
			metric.WithDescription("Memory usage of the Pod."),
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
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
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
			metric.WithDescription("Pod network errors."),
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
			metric.WithDescription("Network bytes for the Pod."),
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
			metric.WithDescription("The time the Pod has been running."),
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

// NewPodVolumeAvailable returns a new PodVolumeAvailable instrument.
func NewPodVolumeAvailable(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeAvailable, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeAvailable{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.available",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Pod volume storage space available."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewPodVolumeCapacity returns a new PodVolumeCapacity instrument.
func NewPodVolumeCapacity(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeCapacity, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeCapacity{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.capacity",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Pod volume total capacity."),
			metric.WithUnit("By"),
		}, opt...)...,
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

// NewPodVolumeInodeCount returns a new PodVolumeInodeCount instrument.
func NewPodVolumeInodeCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeInodeCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeCount{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.inode.count",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The total inodes in the filesystem of the Pod's volume."),
			metric.WithUnit("{inode}"),
		}, opt...)...,
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

// NewPodVolumeInodeFree returns a new PodVolumeInodeFree instrument.
func NewPodVolumeInodeFree(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeInodeFree, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeFree{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.inode.free",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The free inodes in the filesystem of the Pod's volume."),
			metric.WithUnit("{inode}"),
		}, opt...)...,
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

// NewPodVolumeInodeUsed returns a new PodVolumeInodeUsed instrument.
func NewPodVolumeInodeUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeInodeUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeInodeUsed{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.inode.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The inodes used by the filesystem of the Pod's volume."),
			metric.WithUnit("{inode}"),
		}, opt...)...,
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

// NewPodVolumeUsage returns a new PodVolumeUsage instrument.
func NewPodVolumeUsage(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (PodVolumeUsage, error) {
	// Check if the meter is nil.
	if m == nil {
		return PodVolumeUsage{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"k8s.pod.volume.usage",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Pod volume usage."),
			metric.WithUnit("By"),
		}, opt...)...,
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
			metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset."),
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
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replicaset."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicaSetStatus].
//
// [K8s ReplicaSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicaSetStatus].
//
// [K8s ReplicaSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetstatus-v1-apps
func (m ReplicaSetAvailablePods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of desired replica pods in this replicaset."),
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
	return "Number of desired replica pods in this replicaset."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicaSetSpec].
//
// [K8s ReplicaSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetspec-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicaSetSpec].
//
// [K8s ReplicaSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicasetspec-v1-apps
func (m ReplicaSetDesiredPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller."),
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
	return "Total number of available replica pods (ready for at least minReadySeconds) targeted by this replication controller."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicationControllerStatus]
//
// [K8s ReplicationControllerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerstatus-v1-core
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `availableReplicas` field of the
// [K8s ReplicationControllerStatus]
//
// [K8s ReplicationControllerStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerstatus-v1-core
func (m ReplicationControllerAvailablePods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of desired replica pods in this replication controller."),
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
	return "Number of desired replica pods in this replication controller."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicationControllerSpec]
//
// [K8s ReplicationControllerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerspec-v1-core
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s ReplicationControllerSpec]
//
// [K8s ReplicationControllerSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#replicationcontrollerspec-v1-core
func (m ReplicationControllerDesiredPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.limit.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The CPU limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.limit.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The CPU limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.request.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The CPU requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.cpu.request.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The CPU requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("{cpu}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.limit.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.limit.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The sum of local ephemeral storage limits in the namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.request.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.ephemeral_storage.request.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The sum of local ephemeral storage requests in the namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.hugepage_count.request.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The huge page requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("{hugepage}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.hugepage_count.request.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The huge page requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("{hugepage}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.limit.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The memory limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.limit.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The memory limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.request.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The memory requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.memory.request.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The memory requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.object_count.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The object count limits in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("{object}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.object_count.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The object count limits in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("{object}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.persistentvolumeclaim_count.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("{persistentvolumeclaim}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.persistentvolumeclaim_count.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The total number of PersistentVolumeClaims that can exist in the namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("{persistentvolumeclaim}"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.storage.request.hard",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The storage requests in a specific namespace. The value represents the configured quota limit of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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

	i, err := m.Int64UpDownCounter(
		"k8s.resourcequota.storage.request.used",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("The storage requests in a specific namespace. The value represents the current observed total usage of the resource in the namespace."),
			metric.WithUnit("By"),
		}, opt...)...,
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
			metric.WithDescription("The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision."),
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
	return "The number of replica pods created by the statefulset controller from the statefulset version indicated by currentRevision."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `currentReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetCurrentPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of desired replica pods in this statefulset."),
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
	return "Number of desired replica pods in this statefulset."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `replicas` field of the
// [K8s StatefulSetSpec].
//
// [K8s StatefulSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetspec-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `replicas` field of the
// [K8s StatefulSetSpec].
//
// [K8s StatefulSetSpec]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetspec-v1-apps
func (m StatefulSetDesiredPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("The number of replica pods created for this statefulset with a Ready Condition."),
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
	return "The number of replica pods created for this statefulset with a Ready Condition."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `readyReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `readyReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetReadyPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
			metric.WithDescription("Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision."),
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
	return "Number of replica pods created by the statefulset controller from the statefulset version indicated by updateRevision."
}

// Add adds incr to the existing count for attrs.
//
// This metric aligns with the `updatedReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
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

// AddSet adds incr to the existing count for set.
//
// This metric aligns with the `updatedReplicas` field of the
// [K8s StatefulSetStatus].
//
// [K8s StatefulSetStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#statefulsetstatus-v1-apps
func (m StatefulSetUpdatedPods) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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