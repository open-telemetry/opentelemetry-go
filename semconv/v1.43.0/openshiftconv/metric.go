// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package openshiftconv provides types and functionality for OpenTelemetry semantic
// conventions in the "openshift" namespace.
package openshiftconv

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/semconv/internal/metricpool"
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaCPULimitHardObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.cpu.limit.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaCPULimitHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaCPULimitHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPULimitHardObservable returns a new
// ClusterquotaCPULimitHardObservable instrument.
func NewClusterquotaCPULimitHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaCPULimitHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPULimitHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPULimitHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaCPULimitHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.cpu.limit.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaCPULimitHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaCPULimitHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPULimitHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPULimitHardObservable) Name() string {
	return "openshift.clusterquota.cpu.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPULimitHardObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPULimitHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPULimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaCPULimitUsedObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.cpu.limit.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaCPULimitUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaCPULimitUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPULimitUsedObservable returns a new
// ClusterquotaCPULimitUsedObservable instrument.
func NewClusterquotaCPULimitUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaCPULimitUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPULimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPULimitUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaCPULimitUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.cpu.limit.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaCPULimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaCPULimitUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPULimitUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPULimitUsedObservable) Name() string {
	return "openshift.clusterquota.cpu.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPULimitUsedObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPULimitUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaCPURequestHardObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.cpu.request.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaCPURequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaCPURequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPURequestHardObservable returns a new
// ClusterquotaCPURequestHardObservable instrument.
func NewClusterquotaCPURequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaCPURequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPURequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPURequestHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaCPURequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.cpu.request.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaCPURequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaCPURequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPURequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPURequestHardObservable) Name() string {
	return "openshift.clusterquota.cpu.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPURequestHardObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPURequestHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaCPURequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaCPURequestUsedObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.cpu.request.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaCPURequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaCPURequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{cpu}"),
}

// NewClusterquotaCPURequestUsedObservable returns a new
// ClusterquotaCPURequestUsedObservable instrument.
func NewClusterquotaCPURequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaCPURequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaCPURequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaCPURequestUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaCPURequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.cpu.request.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaCPURequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaCPURequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaCPURequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaCPURequestUsedObservable) Name() string {
	return "openshift.clusterquota.cpu.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaCPURequestUsedObservable) Unit() string {
	return "{cpu}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaCPURequestUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageLimitHardObservable is an instrument used to
// record metric values conforming to the
// "openshift.clusterquota.ephemeral_storage.limit.hard" semantic conventions. It
// represents the enforced hard limit of the resource across all projects.
type ClusterquotaEphemeralStorageLimitHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaEphemeralStorageLimitHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageLimitHardObservable returns a new
// ClusterquotaEphemeralStorageLimitHardObservable instrument.
func NewClusterquotaEphemeralStorageLimitHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaEphemeralStorageLimitHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageLimitHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageLimitHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.ephemeral_storage.limit.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaEphemeralStorageLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageLimitHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageLimitHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageLimitHardObservable) Name() string {
	return "openshift.clusterquota.ephemeral_storage.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageLimitHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageLimitHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageLimitUsedObservable is an instrument used to
// record metric values conforming to the
// "openshift.clusterquota.ephemeral_storage.limit.used" semantic conventions. It
// represents the current observed total usage of the resource across all
// projects.
type ClusterquotaEphemeralStorageLimitUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaEphemeralStorageLimitUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageLimitUsedObservable returns a new
// ClusterquotaEphemeralStorageLimitUsedObservable instrument.
func NewClusterquotaEphemeralStorageLimitUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaEphemeralStorageLimitUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageLimitUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageLimitUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.ephemeral_storage.limit.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaEphemeralStorageLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageLimitUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageLimitUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageLimitUsedObservable) Name() string {
	return "openshift.clusterquota.ephemeral_storage.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageLimitUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageLimitUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageRequestHardObservable is an instrument used to
// record metric values conforming to the
// "openshift.clusterquota.ephemeral_storage.request.hard" semantic conventions.
// It represents the enforced hard limit of the resource across all projects.
type ClusterquotaEphemeralStorageRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaEphemeralStorageRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageRequestHardObservable returns a new
// ClusterquotaEphemeralStorageRequestHardObservable instrument.
func NewClusterquotaEphemeralStorageRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaEphemeralStorageRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageRequestHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.ephemeral_storage.request.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaEphemeralStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageRequestHardObservable) Name() string {
	return "openshift.clusterquota.ephemeral_storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageRequestHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageRequestHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaEphemeralStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaEphemeralStorageRequestUsedObservable is an instrument used to
// record metric values conforming to the
// "openshift.clusterquota.ephemeral_storage.request.used" semantic conventions.
// It represents the current observed total usage of the resource across all
// projects.
type ClusterquotaEphemeralStorageRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaEphemeralStorageRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaEphemeralStorageRequestUsedObservable returns a new
// ClusterquotaEphemeralStorageRequestUsedObservable instrument.
func NewClusterquotaEphemeralStorageRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaEphemeralStorageRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaEphemeralStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaEphemeralStorageRequestUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaEphemeralStorageRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.ephemeral_storage.request.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaEphemeralStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaEphemeralStorageRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaEphemeralStorageRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaEphemeralStorageRequestUsedObservable) Name() string {
	return "openshift.clusterquota.ephemeral_storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaEphemeralStorageRequestUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaEphemeralStorageRequestUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestHard) Add(
	ctx context.Context,
	incr int64,
	k8sHugepageSize string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.hugepage.size", k8sHugepageSize),
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaHugepageCountRequestHardObservable is an instrument used to record
// metric values conforming to the
// "openshift.clusterquota.hugepage_count.request.hard" semantic conventions. It
// represents the enforced hard limit of the resource across all projects.
type ClusterquotaHugepageCountRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaHugepageCountRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{hugepage}"),
}

// NewClusterquotaHugepageCountRequestHardObservable returns a new
// ClusterquotaHugepageCountRequestHardObservable instrument.
func NewClusterquotaHugepageCountRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaHugepageCountRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaHugepageCountRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaHugepageCountRequestHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaHugepageCountRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.hugepage_count.request.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaHugepageCountRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaHugepageCountRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaHugepageCountRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaHugepageCountRequestHardObservable) Name() string {
	return "openshift.clusterquota.hugepage_count.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaHugepageCountRequestHardObservable) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaHugepageCountRequestHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// AttrK8SHugepageSize returns a required attribute for the "k8s.hugepage.size"
// semantic convention. It represents the size (identifier) of the K8s huge page.
func (ClusterquotaHugepageCountRequestHardObservable) AttrK8SHugepageSize(val string) attribute.KeyValue {
	return attribute.String("k8s.hugepage.size", val)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestUsed) Add(
	ctx context.Context,
	incr int64,
	k8sHugepageSize string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.hugepage.size", k8sHugepageSize),
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaHugepageCountRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaHugepageCountRequestUsedObservable is an instrument used to record
// metric values conforming to the
// "openshift.clusterquota.hugepage_count.request.used" semantic conventions. It
// represents the current observed total usage of the resource across all
// projects.
type ClusterquotaHugepageCountRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaHugepageCountRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{hugepage}"),
}

// NewClusterquotaHugepageCountRequestUsedObservable returns a new
// ClusterquotaHugepageCountRequestUsedObservable instrument.
func NewClusterquotaHugepageCountRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaHugepageCountRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaHugepageCountRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaHugepageCountRequestUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaHugepageCountRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.hugepage_count.request.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaHugepageCountRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaHugepageCountRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaHugepageCountRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaHugepageCountRequestUsedObservable) Name() string {
	return "openshift.clusterquota.hugepage_count.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaHugepageCountRequestUsedObservable) Unit() string {
	return "{hugepage}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaHugepageCountRequestUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// AttrK8SHugepageSize returns a required attribute for the "k8s.hugepage.size"
// semantic convention. It represents the size (identifier) of the K8s huge page.
func (ClusterquotaHugepageCountRequestUsedObservable) AttrK8SHugepageSize(val string) attribute.KeyValue {
	return attribute.String("k8s.hugepage.size", val)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryLimitHardObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.memory.limit.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaMemoryLimitHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaMemoryLimitHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryLimitHardObservable returns a new
// ClusterquotaMemoryLimitHardObservable instrument.
func NewClusterquotaMemoryLimitHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaMemoryLimitHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryLimitHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaMemoryLimitHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.memory.limit.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaMemoryLimitHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaMemoryLimitHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryLimitHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryLimitHardObservable) Name() string {
	return "openshift.clusterquota.memory.limit.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryLimitHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryLimitHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryLimitUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryLimitUsedObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.memory.limit.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaMemoryLimitUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaMemoryLimitUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryLimitUsedObservable returns a new
// ClusterquotaMemoryLimitUsedObservable instrument.
func NewClusterquotaMemoryLimitUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaMemoryLimitUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryLimitUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaMemoryLimitUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.memory.limit.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaMemoryLimitUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaMemoryLimitUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryLimitUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryLimitUsedObservable) Name() string {
	return "openshift.clusterquota.memory.limit.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryLimitUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryLimitUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestHard) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryRequestHardObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.memory.request.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaMemoryRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaMemoryRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryRequestHardObservable returns a new
// ClusterquotaMemoryRequestHardObservable instrument.
func NewClusterquotaMemoryRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaMemoryRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryRequestHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaMemoryRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.memory.request.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaMemoryRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaMemoryRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryRequestHardObservable) Name() string {
	return "openshift.clusterquota.memory.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryRequestHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryRequestHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestUsed) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaMemoryRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaMemoryRequestUsedObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.memory.request.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaMemoryRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaMemoryRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaMemoryRequestUsedObservable returns a new
// ClusterquotaMemoryRequestUsedObservable instrument.
func NewClusterquotaMemoryRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaMemoryRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaMemoryRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaMemoryRequestUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaMemoryRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.memory.request.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaMemoryRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaMemoryRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaMemoryRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaMemoryRequestUsedObservable) Name() string {
	return "openshift.clusterquota.memory.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaMemoryRequestUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaMemoryRequestUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountHard) Add(
	ctx context.Context,
	incr int64,
	k8sResourcequotaResourceName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.resourcequota.resource_name", k8sResourcequotaResourceName),
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaObjectCountHardObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.object_count.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaObjectCountHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaObjectCountHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{object}"),
}

// NewClusterquotaObjectCountHardObservable returns a new
// ClusterquotaObjectCountHardObservable instrument.
func NewClusterquotaObjectCountHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaObjectCountHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaObjectCountHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaObjectCountHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaObjectCountHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.object_count.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaObjectCountHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaObjectCountHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaObjectCountHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaObjectCountHardObservable) Name() string {
	return "openshift.clusterquota.object_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaObjectCountHardObservable) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaObjectCountHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// AttrK8SResourceQuotaResourceName returns a required attribute for the
// "k8s.resourcequota.resource_name" semantic convention. It represents the name
// of the K8s resource a resource quota defines.
func (ClusterquotaObjectCountHardObservable) AttrK8SResourceQuotaResourceName(val string) attribute.KeyValue {
	return attribute.String("k8s.resourcequota.resource_name", val)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountUsed) Add(
	ctx context.Context,
	incr int64,
	k8sResourcequotaResourceName string,
	attrs ...attribute.KeyValue,
) {
	if !m.Int64UpDownCounter.Enabled(ctx) {
		return
	}
	if len(attrs) == 0 {
		m.Int64UpDownCounter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("k8s.resourcequota.resource_name", k8sResourcequotaResourceName),
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaObjectCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClusterquotaObjectCountUsedObservable is an instrument used to record metric
// values conforming to the "openshift.clusterquota.object_count.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaObjectCountUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaObjectCountUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{object}"),
}

// NewClusterquotaObjectCountUsedObservable returns a new
// ClusterquotaObjectCountUsedObservable instrument.
func NewClusterquotaObjectCountUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaObjectCountUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaObjectCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaObjectCountUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaObjectCountUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.object_count.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaObjectCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaObjectCountUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaObjectCountUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaObjectCountUsedObservable) Name() string {
	return "openshift.clusterquota.object_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaObjectCountUsedObservable) Unit() string {
	return "{object}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaObjectCountUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// AttrK8SResourceQuotaResourceName returns a required attribute for the
// "k8s.resourcequota.resource_name" semantic convention. It represents the name
// of the K8s resource a resource quota defines.
func (ClusterquotaObjectCountUsedObservable) AttrK8SResourceQuotaResourceName(val string) attribute.KeyValue {
	return attribute.String("k8s.resourcequota.resource_name", val)
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountHard) Add(
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaPersistentvolumeclaimCountHard) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ClusterquotaPersistentvolumeclaimCountHardObservable is an instrument used to
// record metric values conforming to the
// "openshift.clusterquota.persistentvolumeclaim_count.hard" semantic
// conventions. It represents the enforced hard limit of the resource across all
// projects.
type ClusterquotaPersistentvolumeclaimCountHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaPersistentvolumeclaimCountHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewClusterquotaPersistentvolumeclaimCountHardObservable returns a new
// ClusterquotaPersistentvolumeclaimCountHardObservable instrument.
func NewClusterquotaPersistentvolumeclaimCountHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaPersistentvolumeclaimCountHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaPersistentvolumeclaimCountHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaPersistentvolumeclaimCountHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaPersistentvolumeclaimCountHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.persistentvolumeclaim_count.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaPersistentvolumeclaimCountHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaPersistentvolumeclaimCountHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaPersistentvolumeclaimCountHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaPersistentvolumeclaimCountHardObservable) Name() string {
	return "openshift.clusterquota.persistentvolumeclaim_count.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaPersistentvolumeclaimCountHardObservable) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaPersistentvolumeclaimCountHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaPersistentvolumeclaimCountHardObservable) AttrK8SStorageclassName(val string) attribute.KeyValue {
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountUsed) Add(
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaPersistentvolumeclaimCountUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaPersistentvolumeclaimCountUsed) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ClusterquotaPersistentvolumeclaimCountUsedObservable is an instrument used to
// record metric values conforming to the
// "openshift.clusterquota.persistentvolumeclaim_count.used" semantic
// conventions. It represents the current observed total usage of the resource
// across all projects.
type ClusterquotaPersistentvolumeclaimCountUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaPersistentvolumeclaimCountUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("{persistentvolumeclaim}"),
}

// NewClusterquotaPersistentvolumeclaimCountUsedObservable returns a new
// ClusterquotaPersistentvolumeclaimCountUsedObservable instrument.
func NewClusterquotaPersistentvolumeclaimCountUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaPersistentvolumeclaimCountUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaPersistentvolumeclaimCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaPersistentvolumeclaimCountUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaPersistentvolumeclaimCountUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.persistentvolumeclaim_count.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaPersistentvolumeclaimCountUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaPersistentvolumeclaimCountUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaPersistentvolumeclaimCountUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaPersistentvolumeclaimCountUsedObservable) Name() string {
	return "openshift.clusterquota.persistentvolumeclaim_count.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaPersistentvolumeclaimCountUsedObservable) Unit() string {
	return "{persistentvolumeclaim}"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaPersistentvolumeclaimCountUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaPersistentvolumeclaimCountUsedObservable) AttrK8SStorageclassName(val string) attribute.KeyValue {
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestHard) Add(
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
// This metric is retrieved from the `Status.Total.Hard` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestHard) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaStorageRequestHard) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ClusterquotaStorageRequestHardObservable is an instrument used to record
// metric values conforming to the "openshift.clusterquota.storage.request.hard"
// semantic conventions. It represents the enforced hard limit of the resource
// across all projects.
type ClusterquotaStorageRequestHardObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaStorageRequestHardObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The enforced hard limit of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaStorageRequestHardObservable returns a new
// ClusterquotaStorageRequestHardObservable instrument.
func NewClusterquotaStorageRequestHardObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaStorageRequestHardObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaStorageRequestHardObservableOpts
	} else {
		opt = append(opt, newClusterquotaStorageRequestHardObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.storage.request.hard",
		opt...,
	)
	if err != nil {
		return ClusterquotaStorageRequestHardObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaStorageRequestHardObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaStorageRequestHardObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaStorageRequestHardObservable) Name() string {
	return "openshift.clusterquota.storage.request.hard"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaStorageRequestHardObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaStorageRequestHardObservable) Description() string {
	return "The enforced hard limit of the resource across all projects."
}

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaStorageRequestHardObservable) AttrK8SStorageclassName(val string) attribute.KeyValue {
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
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestUsed) Add(
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
// This metric is retrieved from the `Status.Total.Used` field of the
// [K8s ResourceQuotaStatus]
// of the
// [ClusterResourceQuota].
//
// The `k8s.storageclass.name` should be required when a resource quota is
// defined for a specific
// storage class.
//
// [K8s ResourceQuotaStatus]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#resourcequotastatus-v1-core
// [ClusterResourceQuota]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html/schedule_and_quota_apis/clusterresourcequota-quota-openshift-io-v1#status-total
func (m ClusterquotaStorageRequestUsed) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaStorageRequestUsed) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}

// ClusterquotaStorageRequestUsedObservable is an instrument used to record
// metric values conforming to the "openshift.clusterquota.storage.request.used"
// semantic conventions. It represents the current observed total usage of the
// resource across all projects.
type ClusterquotaStorageRequestUsedObservable struct {
	metric.Int64ObservableUpDownCounter
}

var newClusterquotaStorageRequestUsedObservableOpts = []metric.Int64ObservableUpDownCounterOption{
	metric.WithDescription("The current observed total usage of the resource across all projects."),
	metric.WithUnit("By"),
}

// NewClusterquotaStorageRequestUsedObservable returns a new
// ClusterquotaStorageRequestUsedObservable instrument.
func NewClusterquotaStorageRequestUsedObservable(
	m metric.Meter,
	opt ...metric.Int64ObservableUpDownCounterOption,
) (ClusterquotaStorageRequestUsedObservable, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClusterquotaStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClusterquotaStorageRequestUsedObservableOpts
	} else {
		opt = append(opt, newClusterquotaStorageRequestUsedObservableOpts...)
	}

	i, err := m.Int64ObservableUpDownCounter(
		"openshift.clusterquota.storage.request.used",
		opt...,
	)
	if err != nil {
		return ClusterquotaStorageRequestUsedObservable{noop.Int64ObservableUpDownCounter{}}, err
	}
	return ClusterquotaStorageRequestUsedObservable{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClusterquotaStorageRequestUsedObservable) Inst() metric.Int64ObservableUpDownCounter {
	return m.Int64ObservableUpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ClusterquotaStorageRequestUsedObservable) Name() string {
	return "openshift.clusterquota.storage.request.used"
}

// Unit returns the semantic convention unit of the instrument
func (ClusterquotaStorageRequestUsedObservable) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClusterquotaStorageRequestUsedObservable) Description() string {
	return "The current observed total usage of the resource across all projects."
}

// AttrK8SStorageclassName returns an optional attribute for the
// "k8s.storageclass.name" semantic convention. It represents the name of K8s
// [StorageClass] object.
//
// [StorageClass]: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.34/#storageclass-v1-storage-k8s-io
func (ClusterquotaStorageRequestUsedObservable) AttrK8SStorageclassName(val string) attribute.KeyValue {
	return attribute.String("k8s.storageclass.name", val)
}
