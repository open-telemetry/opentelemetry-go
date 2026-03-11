// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package openshiftconv provides types and functionality for OpenTelemetry semantic
// conventions in the "openshift" namespace.
package openshiftconv

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

// ClusterquotaCPULimitHard is an instrument used to record metric values
// conforming to the "openshift.clusterquota.cpu.limit.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaCPULimitHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaCPULimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPULimitHard returns a new ClusterquotaCPULimitHard instrument.
func NewClusterquotaCPULimitHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaCPULimitHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPULimitHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPULimitHardOpts
	} else {
		opt = append(opt, newClusterquotaCPULimitHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.cpu.limit.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaCPULimitHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaCPULimitHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPULimitHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPULimitHard) Name() string {
	return "openshift.clusterquota.cpu.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPULimitHard) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPULimitHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaCPULimitUsed is an instrument used to record metric values
// conforming to the "openshift.clusterquota.cpu.limit.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaCPULimitUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaCPULimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPULimitUsed returns a new ClusterquotaCPULimitUsed instrument.
func NewClusterquotaCPULimitUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaCPULimitUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPULimitUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPULimitUsedOpts
	} else {
		opt = append(opt, newClusterquotaCPULimitUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.cpu.limit.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaCPULimitUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaCPULimitUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPULimitUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPULimitUsed) Name() string {
	return "openshift.clusterquota.cpu.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPULimitUsed) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPULimitUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaCPURequestHard is an instrument used to record metric values
// conforming to the "openshift.clusterquota.cpu.request.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaCPURequestHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaCPURequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPURequestHard returns a new ClusterquotaCPURequestHard
// instrument.
func NewClusterquotaCPURequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaCPURequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPURequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPURequestHardOpts
	} else {
		opt = append(opt, newClusterquotaCPURequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.cpu.request.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaCPURequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaCPURequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPURequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPURequestHard) Name() string {
	return "openshift.clusterquota.cpu.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPURequestHard) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPURequestHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaCPURequestUsed is an instrument used to record metric values
// conforming to the "openshift.clusterquota.cpu.request.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaCPURequestUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaCPURequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPURequestUsed returns a new ClusterquotaCPURequestUsed
// instrument.
func NewClusterquotaCPURequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaCPURequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPURequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPURequestUsedOpts
	} else {
		opt = append(opt, newClusterquotaCPURequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.cpu.request.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaCPURequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaCPURequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPURequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPURequestUsed) Name() string {
	return "openshift.clusterquota.cpu.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPURequestUsed) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPURequestUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageLimitHard is an instrument used to record metric
// values conforming to the "openshift.clusterquota.ephemeral_storage.limit.hard"
// semantic conventions. It represents the enforced hard limit of the resource
// across all projects.
type ClusterquotaEphemeralStorageLimitHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaEphemeralStorageLimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageLimitHard returns a new
// ClusterquotaEphemeralStorageLimitHard instrument.
func NewClusterquotaEphemeralStorageLimitHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaEphemeralStorageLimitHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageLimitHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageLimitHardOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageLimitHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.ephemeral_storage.limit.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaEphemeralStorageLimitHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageLimitHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageLimitHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageLimitHard) Name() string {
	return "openshift.clusterquota.ephemeral_storage.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageLimitHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageLimitHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageLimitUsed is an instrument used to record metric
// values conforming to the "openshift.clusterquota.ephemeral_storage.limit.used"
// semantic conventions. It represents the current observed total usage of the
// resource across all projects.
type ClusterquotaEphemeralStorageLimitUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaEphemeralStorageLimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageLimitUsed returns a new
// ClusterquotaEphemeralStorageLimitUsed instrument.
func NewClusterquotaEphemeralStorageLimitUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaEphemeralStorageLimitUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageLimitUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageLimitUsedOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageLimitUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.ephemeral_storage.limit.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaEphemeralStorageLimitUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageLimitUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageLimitUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageLimitUsed) Name() string {
	return "openshift.clusterquota.ephemeral_storage.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageLimitUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageLimitUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageRequestHard is an instrument used to record metric
// values conforming to the
// "openshift.clusterquota.ephemeral_storage.request.hard" semantic conventions.
// It represents the enforced hard limit of the resource across all projects.
type ClusterquotaEphemeralStorageRequestHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaEphemeralStorageRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageRequestHard returns a new
// ClusterquotaEphemeralStorageRequestHard instrument.
func NewClusterquotaEphemeralStorageRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaEphemeralStorageRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageRequestHardOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.ephemeral_storage.request.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaEphemeralStorageRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageRequestHard) Name() string {
	return "openshift.clusterquota.ephemeral_storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageRequestHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageRequestHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageRequestUsed is an instrument used to record metric
// values conforming to the
// "openshift.clusterquota.ephemeral_storage.request.used" semantic conventions.
// It represents the current observed total usage of the resource across all
// projects.
type ClusterquotaEphemeralStorageRequestUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaEphemeralStorageRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageRequestUsed returns a new
// ClusterquotaEphemeralStorageRequestUsed instrument.
func NewClusterquotaEphemeralStorageRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaEphemeralStorageRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageRequestUsedOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.ephemeral_storage.request.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaEphemeralStorageRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageRequestUsed) Name() string {
	return "openshift.clusterquota.ephemeral_storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageRequestUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageRequestUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaHugepageCountRequestHard is an instrument used to record metric
// values conforming to the "openshift.clusterquota.hugepage_count.request.hard"
// semantic conventions. It represents the enforced hard limit of the resource
// across all projects.
type ClusterquotaHugepageCountRequestHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaHugepageCountRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{hugepage}"),
}

// NewClusterquotaHugepageCountRequestHard returns a new
// ClusterquotaHugepageCountRequestHard instrument.
func NewClusterquotaHugepageCountRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaHugepageCountRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaHugepageCountRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaHugepageCountRequestHardOpts
	} else {
		opt = append(opt, newClusterquotaHugepageCountRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.hugepage_count.request.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaHugepageCountRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaHugepageCountRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaHugepageCountRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaHugepageCountRequestHard) Name() string {
	return "openshift.clusterquota.hugepage_count.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaHugepageCountRequestHard) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaHugepageCountRequestHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// The k8sHugepageSize is the the size (identifier) of the K8s huge page.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestHard) Add(
	ctx context.Context,
	incr int64,
	k8sHugepageSize string,
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
				attribute.String("k8s.hugepage.size", k8sHugepageSize),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaHugepageCountRequestUsed is an instrument used to record metric
// values conforming to the "openshift.clusterquota.hugepage_count.request.used"
// semantic conventions. It represents the current observed total usage of the
// resource across all projects.
type ClusterquotaHugepageCountRequestUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaHugepageCountRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{hugepage}"),
}

// NewClusterquotaHugepageCountRequestUsed returns a new
// ClusterquotaHugepageCountRequestUsed instrument.
func NewClusterquotaHugepageCountRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaHugepageCountRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaHugepageCountRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaHugepageCountRequestUsedOpts
	} else {
		opt = append(opt, newClusterquotaHugepageCountRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.hugepage_count.request.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaHugepageCountRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaHugepageCountRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaHugepageCountRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaHugepageCountRequestUsed) Name() string {
	return "openshift.clusterquota.hugepage_count.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaHugepageCountRequestUsed) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaHugepageCountRequestUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// The k8sHugepageSize is the the size (identifier) of the K8s huge page.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestUsed) Add(
	ctx context.Context,
	incr int64,
	k8sHugepageSize string,
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
				attribute.String("k8s.hugepage.size", k8sHugepageSize),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryLimitHard is an instrument used to record metric values
// conforming to the "openshift.clusterquota.memory.limit.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaMemoryLimitHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaMemoryLimitHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryLimitHard returns a new ClusterquotaMemoryLimitHard
// instrument.
func NewClusterquotaMemoryLimitHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaMemoryLimitHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryLimitHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryLimitHardOpts
	} else {
		opt = append(opt, newClusterquotaMemoryLimitHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.memory.limit.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaMemoryLimitHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaMemoryLimitHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryLimitHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryLimitHard) Name() string {
	return "openshift.clusterquota.memory.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryLimitHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryLimitHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryLimitUsed is an instrument used to record metric values
// conforming to the "openshift.clusterquota.memory.limit.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaMemoryLimitUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaMemoryLimitUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryLimitUsed returns a new ClusterquotaMemoryLimitUsed
// instrument.
func NewClusterquotaMemoryLimitUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaMemoryLimitUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryLimitUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryLimitUsedOpts
	} else {
		opt = append(opt, newClusterquotaMemoryLimitUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.memory.limit.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaMemoryLimitUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaMemoryLimitUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryLimitUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryLimitUsed) Name() string {
	return "openshift.clusterquota.memory.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryLimitUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryLimitUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryRequestHard is an instrument used to record metric values
// conforming to the "openshift.clusterquota.memory.request.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaMemoryRequestHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaMemoryRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryRequestHard returns a new ClusterquotaMemoryRequestHard
// instrument.
func NewClusterquotaMemoryRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaMemoryRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryRequestHardOpts
	} else {
		opt = append(opt, newClusterquotaMemoryRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.memory.request.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaMemoryRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaMemoryRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryRequestHard) Name() string {
	return "openshift.clusterquota.memory.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryRequestHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryRequestHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryRequestUsed is an instrument used to record metric values
// conforming to the "openshift.clusterquota.memory.request.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaMemoryRequestUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaMemoryRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryRequestUsed returns a new ClusterquotaMemoryRequestUsed
// instrument.
func NewClusterquotaMemoryRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaMemoryRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryRequestUsedOpts
	} else {
		opt = append(opt, newClusterquotaMemoryRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.memory.request.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaMemoryRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaMemoryRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryRequestUsed) Name() string {
	return "openshift.clusterquota.memory.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryRequestUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryRequestUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaObjectCountHard is an instrument used to record metric values
// conforming to the "openshift.clusterquota.object_count.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaObjectCountHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaObjectCountHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{object}"),
}

// NewClusterquotaObjectCountHard returns a new ClusterquotaObjectCountHard
// instrument.
func NewClusterquotaObjectCountHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaObjectCountHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaObjectCountHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaObjectCountHardOpts
	} else {
		opt = append(opt, newClusterquotaObjectCountHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.object_count.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaObjectCountHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaObjectCountHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaObjectCountHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaObjectCountHard) Name() string {
	return "openshift.clusterquota.object_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaObjectCountHard) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaObjectCountHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// The k8sResourcequotaResourceName is the the name of the K8s resource a
// resource quota defines.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountHard) Add(
	ctx context.Context,
	incr int64,
	k8sResourcequotaResourceName string,
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
				attribute.String("k8s.resourcequota.resource_name", k8sResourcequotaResourceName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaObjectCountUsed is an instrument used to record metric values
// conforming to the "openshift.clusterquota.object_count.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaObjectCountUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaObjectCountUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{object}"),
}

// NewClusterquotaObjectCountUsed returns a new ClusterquotaObjectCountUsed
// instrument.
func NewClusterquotaObjectCountUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaObjectCountUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaObjectCountUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaObjectCountUsedOpts
	} else {
		opt = append(opt, newClusterquotaObjectCountUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.object_count.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaObjectCountUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaObjectCountUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaObjectCountUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaObjectCountUsed) Name() string {
	return "openshift.clusterquota.object_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaObjectCountUsed) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaObjectCountUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// The k8sResourcequotaResourceName is the the name of the K8s resource a
// resource quota defines.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountUsed) Add(
	ctx context.Context,
	incr int64,
	k8sResourcequotaResourceName string,
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
				attribute.String("k8s.resourcequota.resource_name", k8sResourcequotaResourceName),
			)...,
		),
	)

	m.Int64UpDownCounter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaPersistentvolumeclaimCountHard is an instrument used to record
// metric values conforming to the
// "openshift.clusterquota.persistentvolumeclaim_count.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaPersistentvolumeclaimCountHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaPersistentvolumeclaimCountHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewClusterquotaPersistentvolumeclaimCountHard returns a new
// ClusterquotaPersistentvolumeclaimCountHard instrument.
func NewClusterquotaPersistentvolumeclaimCountHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaPersistentvolumeclaimCountHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaPersistentvolumeclaimCountHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaPersistentvolumeclaimCountHardOpts
	} else {
		opt = append(opt, newClusterquotaPersistentvolumeclaimCountHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.persistentvolumeclaim_count.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaPersistentvolumeclaimCountHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaPersistentvolumeclaimCountHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaPersistentvolumeclaimCountHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaPersistentvolumeclaimCountHard) Name() string {
	return "openshift.clusterquota.persistentvolumeclaim_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaPersistentvolumeclaimCountHard) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaPersistentvolumeclaimCountHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountHard) Add(
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ClusterquotaPersistentvolumeclaimCountHard) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ClusterquotaPersistentvolumeclaimCountUsed is an instrument used to record
// metric values conforming to the
// "openshift.clusterquota.persistentvolumeclaim_count.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaPersistentvolumeclaimCountUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaPersistentvolumeclaimCountUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewClusterquotaPersistentvolumeclaimCountUsed returns a new
// ClusterquotaPersistentvolumeclaimCountUsed instrument.
func NewClusterquotaPersistentvolumeclaimCountUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaPersistentvolumeclaimCountUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaPersistentvolumeclaimCountUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaPersistentvolumeclaimCountUsedOpts
	} else {
		opt = append(opt, newClusterquotaPersistentvolumeclaimCountUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.persistentvolumeclaim_count.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaPersistentvolumeclaimCountUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaPersistentvolumeclaimCountUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaPersistentvolumeclaimCountUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaPersistentvolumeclaimCountUsed) Name() string {
	return "openshift.clusterquota.persistentvolumeclaim_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaPersistentvolumeclaimCountUsed) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaPersistentvolumeclaimCountUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountUsed) Add(
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ClusterquotaPersistentvolumeclaimCountUsed) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ClusterquotaStorageRequestHard is an instrument used to record metric values
// conforming to the "openshift.clusterquota.storage.request.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaStorageRequestHard struct {
	metric.Int64UpDownCounter
}

var newClusterquotaStorageRequestHardOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaStorageRequestHard returns a new ClusterquotaStorageRequestHard
// instrument.
func NewClusterquotaStorageRequestHard(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaStorageRequestHard, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaStorageRequestHard{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaStorageRequestHardOpts
	} else {
		opt = append(opt, newClusterquotaStorageRequestHardOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.storage.request.hard",
		opt...,
	)
	if err != nil {
	    return ClusterquotaStorageRequestHard{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaStorageRequestHard{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaStorageRequestHard) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaStorageRequestHard) Name() string {
	return "openshift.clusterquota.storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaStorageRequestHard) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaStorageRequestHard) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestHard) Add(
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ClusterquotaStorageRequestHard) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ClusterquotaStorageRequestUsed is an instrument used to record metric values
// conforming to the "openshift.clusterquota.storage.request.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaStorageRequestUsed struct {
	metric.Int64UpDownCounter
}

var newClusterquotaStorageRequestUsedOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaStorageRequestUsed returns a new ClusterquotaStorageRequestUsed
// instrument.
func NewClusterquotaStorageRequestUsed(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ClusterquotaStorageRequestUsed, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaStorageRequestUsed{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaStorageRequestUsedOpts
	} else {
		opt = append(opt, newClusterquotaStorageRequestUsedOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"openshift.clusterquota.storage.request.used",
		opt...,
	)
	if err != nil {
	    return ClusterquotaStorageRequestUsed{noop.Int64UpDownCounter{}}, err
	}
	return ClusterquotaStorageRequestUsed{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaStorageRequestUsed) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaStorageRequestUsed) Name() string {
	return "openshift.clusterquota.storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaStorageRequestUsed) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaStorageRequestUsed) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestUsed) Add(
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#storageclass-v1-storage-k8s-io
func (ClusterquotaStorageRequestUsed) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}