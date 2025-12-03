// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
// Package nfsconv provides types and functionality for OpenTelemetry semantic
// conventions in the "nfs" namespace.
package nfsconv

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

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the describes a class of error the operation ended
// with.
type ErrorTypeAttr string

var (
	// ErrorTypeOther is a fallback error value to be used when the instrumentation
	// doesn't define a custom value.
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
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

// NetworkTransportAttr is an attribute conforming to the network.transport
// semantic conventions. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
type NetworkTransportAttr string

var (
	// NetworkTransportTCP is the TCP.
	NetworkTransportTCP NetworkTransportAttr = "tcp"
	// NetworkTransportUDP is the UDP.
	NetworkTransportUDP NetworkTransportAttr = "udp"
	// NetworkTransportPipe is the named or anonymous pipe.
	NetworkTransportPipe NetworkTransportAttr = "pipe"
	// NetworkTransportUnix is the unix domain socket.
	NetworkTransportUnix NetworkTransportAttr = "unix"
	// NetworkTransportQUIC is the QUIC.
	NetworkTransportQUIC NetworkTransportAttr = "quic"
)

// ClientNetCount is an instrument used to record metric values conforming to the
// "nfs.client.net.count" semantic conventions. It represents the reports the
// count of kernel NFS client TCP segments and UDP datagrams handled.
type ClientNetCount struct {
	metric.Int64Counter
}

var newClientNetCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS client TCP segments and UDP datagrams handled."),
	metric.WithUnit("{record}"),
}

// NewClientNetCount returns a new ClientNetCount instrument.
func NewClientNetCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientNetCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientNetCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientNetCountOpts
	} else {
		opt = append(opt, newClientNetCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.client.net.count",
		opt...,
	)
	if err != nil {
	    return ClientNetCount{noop.Int64Counter{}}, err
	}
	return ClientNetCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientNetCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientNetCount) Name() string {
	return "nfs.client.net.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientNetCount) Unit() string {
	return "{record}"
}

// Description returns the semantic convention description of the instrument
func (ClientNetCount) Description() string {
	return "Reports the count of kernel NFS client TCP segments and UDP datagrams handled."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.netudpcnt and
// svc_stat.nettcpcnt
func (m ClientNetCount) Add(
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
// Linux: this metric is taken from the Linux kernel's svc_stat.netudpcnt and
// svc_stat.nettcpcnt
func (m ClientNetCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ClientNetCount) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// ClientNetTCPConnectionAccepted is an instrument used to record metric values
// conforming to the "nfs.client.net.tcp.connection.accepted" semantic
// conventions. It represents the reports the count of kernel NFS client TCP
// connections accepted.
type ClientNetTCPConnectionAccepted struct {
	metric.Int64Counter
}

var newClientNetTCPConnectionAcceptedOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS client TCP connections accepted."),
	metric.WithUnit("{connection}"),
}

// NewClientNetTCPConnectionAccepted returns a new ClientNetTCPConnectionAccepted
// instrument.
func NewClientNetTCPConnectionAccepted(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientNetTCPConnectionAccepted, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientNetTCPConnectionAccepted{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientNetTCPConnectionAcceptedOpts
	} else {
		opt = append(opt, newClientNetTCPConnectionAcceptedOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.client.net.tcp.connection.accepted",
		opt...,
	)
	if err != nil {
	    return ClientNetTCPConnectionAccepted{noop.Int64Counter{}}, err
	}
	return ClientNetTCPConnectionAccepted{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientNetTCPConnectionAccepted) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientNetTCPConnectionAccepted) Name() string {
	return "nfs.client.net.tcp.connection.accepted"
}

// Unit returns the semantic convention unit of the instrument
func (ClientNetTCPConnectionAccepted) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientNetTCPConnectionAccepted) Description() string {
	return "Reports the count of kernel NFS client TCP connections accepted."
}

// Add adds incr to the existing count for attrs.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.nettcpconn
func (m ClientNetTCPConnectionAccepted) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.nettcpconn
func (m ClientNetTCPConnectionAccepted) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClientOperationCount is an instrument used to record metric values conforming
// to the "nfs.client.operation.count" semantic conventions. It represents the
// reports the count of kernel NFSv4+ client operations.
type ClientOperationCount struct {
	metric.Int64Counter
}

var newClientOperationCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFSv4+ client operations."),
	metric.WithUnit("{operation}"),
}

// NewClientOperationCount returns a new ClientOperationCount instrument.
func NewClientOperationCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientOperationCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientOperationCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientOperationCountOpts
	} else {
		opt = append(opt, newClientOperationCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.client.operation.count",
		opt...,
	)
	if err != nil {
	    return ClientOperationCount{noop.Int64Counter{}}, err
	}
	return ClientOperationCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientOperationCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientOperationCount) Name() string {
	return "nfs.client.operation.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOperationCount) Unit() string {
	return "{operation}"
}

// Description returns the semantic convention description of the instrument
func (ClientOperationCount) Description() string {
	return "Reports the count of kernel NFSv4+ client operations."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m ClientOperationCount) Add(
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
func (m ClientOperationCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrOperationName returns an optional attribute for the "nfs.operation.name"
// semantic convention. It represents the NFSv4+ operation name.
func (ClientOperationCount) AttrOperationName(val string) attribute.KeyValue {
	return attribute.String("nfs.operation.name", val)
}

// AttrOncRPCVersion returns an optional attribute for the "onc_rpc.version"
// semantic convention. It represents the ONC/Sun RPC program version.
func (ClientOperationCount) AttrOncRPCVersion(val int) attribute.KeyValue {
	return attribute.Int("onc_rpc.version", val)
}

// ClientProcedureCount is an instrument used to record metric values conforming
// to the "nfs.client.procedure.count" semantic conventions. It represents the
// reports the count of kernel NFS client procedures.
type ClientProcedureCount struct {
	metric.Int64Counter
}

var newClientProcedureCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS client procedures."),
	metric.WithUnit("{procedure}"),
}

// NewClientProcedureCount returns a new ClientProcedureCount instrument.
func NewClientProcedureCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientProcedureCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientProcedureCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientProcedureCountOpts
	} else {
		opt = append(opt, newClientProcedureCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.client.procedure.count",
		opt...,
	)
	if err != nil {
	    return ClientProcedureCount{noop.Int64Counter{}}, err
	}
	return ClientProcedureCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientProcedureCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientProcedureCount) Name() string {
	return "nfs.client.procedure.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientProcedureCount) Unit() string {
	return "{procedure}"
}

// Description returns the semantic convention description of the instrument
func (ClientProcedureCount) Description() string {
	return "Reports the count of kernel NFS client procedures."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m ClientProcedureCount) Add(
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
func (m ClientProcedureCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrOncRPCProcedureName returns an optional attribute for the
// "onc_rpc.procedure.name" semantic convention. It represents the ONC/Sun RPC
// procedure name.
func (ClientProcedureCount) AttrOncRPCProcedureName(val string) attribute.KeyValue {
	return attribute.String("onc_rpc.procedure.name", val)
}

// AttrOncRPCVersion returns an optional attribute for the "onc_rpc.version"
// semantic convention. It represents the ONC/Sun RPC program version.
func (ClientProcedureCount) AttrOncRPCVersion(val int) attribute.KeyValue {
	return attribute.Int("onc_rpc.version", val)
}

// ClientRPCAuthrefreshCount is an instrument used to record metric values
// conforming to the "nfs.client.rpc.authrefresh.count" semantic conventions. It
// represents the reports the count of kernel NFS client RPC authentication
// refreshes.
type ClientRPCAuthrefreshCount struct {
	metric.Int64Counter
}

var newClientRPCAuthrefreshCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS client RPC authentication refreshes."),
	metric.WithUnit("{authrefresh}"),
}

// NewClientRPCAuthrefreshCount returns a new ClientRPCAuthrefreshCount
// instrument.
func NewClientRPCAuthrefreshCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientRPCAuthrefreshCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientRPCAuthrefreshCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientRPCAuthrefreshCountOpts
	} else {
		opt = append(opt, newClientRPCAuthrefreshCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.client.rpc.authrefresh.count",
		opt...,
	)
	if err != nil {
	    return ClientRPCAuthrefreshCount{noop.Int64Counter{}}, err
	}
	return ClientRPCAuthrefreshCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientRPCAuthrefreshCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientRPCAuthrefreshCount) Name() string {
	return "nfs.client.rpc.authrefresh.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientRPCAuthrefreshCount) Unit() string {
	return "{authrefresh}"
}

// Description returns the semantic convention description of the instrument
func (ClientRPCAuthrefreshCount) Description() string {
	return "Reports the count of kernel NFS client RPC authentication refreshes."
}

// Add adds incr to the existing count for attrs.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.rpcauthrefresh
func (m ClientRPCAuthrefreshCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.rpcauthrefresh
func (m ClientRPCAuthrefreshCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClientRPCCount is an instrument used to record metric values conforming to the
// "nfs.client.rpc.count" semantic conventions. It represents the reports the
// count of kernel NFS client RPCs sent, regardless of whether they're
// accepted/rejected by the server.
type ClientRPCCount struct {
	metric.Int64Counter
}

var newClientRPCCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS client RPCs sent, regardless of whether they're accepted/rejected by the server."),
	metric.WithUnit("{request}"),
}

// NewClientRPCCount returns a new ClientRPCCount instrument.
func NewClientRPCCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientRPCCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientRPCCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientRPCCountOpts
	} else {
		opt = append(opt, newClientRPCCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.client.rpc.count",
		opt...,
	)
	if err != nil {
	    return ClientRPCCount{noop.Int64Counter{}}, err
	}
	return ClientRPCCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientRPCCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientRPCCount) Name() string {
	return "nfs.client.rpc.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientRPCCount) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ClientRPCCount) Description() string {
	return "Reports the count of kernel NFS client RPCs sent, regardless of whether they're accepted/rejected by the server."
}

// Add adds incr to the existing count for attrs.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.rpccnt
func (m ClientRPCCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.rpccnt
func (m ClientRPCCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ClientRPCRetransmitCount is an instrument used to record metric values
// conforming to the "nfs.client.rpc.retransmit.count" semantic conventions. It
// represents the reports the count of kernel NFS client RPC retransmits.
type ClientRPCRetransmitCount struct {
	metric.Int64Counter
}

var newClientRPCRetransmitCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS client RPC retransmits."),
	metric.WithUnit("{retransmit}"),
}

// NewClientRPCRetransmitCount returns a new ClientRPCRetransmitCount instrument.
func NewClientRPCRetransmitCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientRPCRetransmitCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientRPCRetransmitCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientRPCRetransmitCountOpts
	} else {
		opt = append(opt, newClientRPCRetransmitCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.client.rpc.retransmit.count",
		opt...,
	)
	if err != nil {
	    return ClientRPCRetransmitCount{noop.Int64Counter{}}, err
	}
	return ClientRPCRetransmitCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientRPCRetransmitCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ClientRPCRetransmitCount) Name() string {
	return "nfs.client.rpc.retransmit.count"
}

// Unit returns the semantic convention unit of the instrument
func (ClientRPCRetransmitCount) Unit() string {
	return "{retransmit}"
}

// Description returns the semantic convention description of the instrument
func (ClientRPCRetransmitCount) Description() string {
	return "Reports the count of kernel NFS client RPC retransmits."
}

// Add adds incr to the existing count for attrs.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.rpcretrans
func (m ClientRPCRetransmitCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.rpcretrans
func (m ClientRPCRetransmitCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ServerFhStaleCount is an instrument used to record metric values conforming to
// the "nfs.server.fh.stale.count" semantic conventions. It represents the
// reports the count of kernel NFS server stale file handles.
type ServerFhStaleCount struct {
	metric.Int64Counter
}

var newServerFhStaleCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS server stale file handles."),
	metric.WithUnit("{fh}"),
}

// NewServerFhStaleCount returns a new ServerFhStaleCount instrument.
func NewServerFhStaleCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerFhStaleCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerFhStaleCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerFhStaleCountOpts
	} else {
		opt = append(opt, newServerFhStaleCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.fh.stale.count",
		opt...,
	)
	if err != nil {
	    return ServerFhStaleCount{noop.Int64Counter{}}, err
	}
	return ServerFhStaleCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerFhStaleCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerFhStaleCount) Name() string {
	return "nfs.server.fh.stale.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServerFhStaleCount) Unit() string {
	return "{fh}"
}

// Description returns the semantic convention description of the instrument
func (ServerFhStaleCount) Description() string {
	return "Reports the count of kernel NFS server stale file handles."
}

// Add adds incr to the existing count for attrs.
//
// Linux: this metric is taken from the Linux kernel NFSD_STATS_FH_STALE counter
// in the nfsd_net struct
func (m ServerFhStaleCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Linux: this metric is taken from the Linux kernel NFSD_STATS_FH_STALE counter
// in the nfsd_net struct
func (m ServerFhStaleCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ServerIO is an instrument used to record metric values conforming to the
// "nfs.server.io" semantic conventions. It represents the reports the count of
// kernel NFS server bytes returned to receive and transmit (read and write)
// requests.
type ServerIO struct {
	metric.Int64Counter
}

var newServerIOOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS server bytes returned to receive and transmit (read and write) requests."),
	metric.WithUnit("By"),
}

// NewServerIO returns a new ServerIO instrument.
func NewServerIO(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerIO, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerIO{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerIOOpts
	} else {
		opt = append(opt, newServerIOOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.io",
		opt...,
	)
	if err != nil {
	    return ServerIO{noop.Int64Counter{}}, err
	}
	return ServerIO{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerIO) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerIO) Name() string {
	return "nfs.server.io"
}

// Unit returns the semantic convention unit of the instrument
func (ServerIO) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ServerIO) Description() string {
	return "Reports the count of kernel NFS server bytes returned to receive and transmit (read and write) requests."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Linux: this metric is taken from the Linux kernel NFSD_STATS_IO_READ and
// NFSD_STATS_IO_WRITE counters in the nfsd_net struct
func (m ServerIO) Add(
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
// Linux: this metric is taken from the Linux kernel NFSD_STATS_IO_READ and
// NFSD_STATS_IO_WRITE counters in the nfsd_net struct
func (m ServerIO) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrNetworkIODirection returns an optional attribute for the
// "network.io.direction" semantic convention. It represents the network IO
// operation direction.
func (ServerIO) AttrNetworkIODirection(val NetworkIODirectionAttr) attribute.KeyValue {
	return attribute.String("network.io.direction", string(val))
}

// ServerNetCount is an instrument used to record metric values conforming to the
// "nfs.server.net.count" semantic conventions. It represents the reports the
// count of kernel NFS server TCP segments and UDP datagrams handled.
type ServerNetCount struct {
	metric.Int64Counter
}

var newServerNetCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS server TCP segments and UDP datagrams handled."),
	metric.WithUnit("{record}"),
}

// NewServerNetCount returns a new ServerNetCount instrument.
func NewServerNetCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerNetCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerNetCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerNetCountOpts
	} else {
		opt = append(opt, newServerNetCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.net.count",
		opt...,
	)
	if err != nil {
	    return ServerNetCount{noop.Int64Counter{}}, err
	}
	return ServerNetCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerNetCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerNetCount) Name() string {
	return "nfs.server.net.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServerNetCount) Unit() string {
	return "{record}"
}

// Description returns the semantic convention description of the instrument
func (ServerNetCount) Description() string {
	return "Reports the count of kernel NFS server TCP segments and UDP datagrams handled."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.nettcpcnt and
// svc_stat.netudpcnt
func (m ServerNetCount) Add(
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
// Linux: this metric is taken from the Linux kernel's svc_stat.nettcpcnt and
// svc_stat.netudpcnt
func (m ServerNetCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ServerNetCount) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// ServerNetTCPConnectionAccepted is an instrument used to record metric values
// conforming to the "nfs.server.net.tcp.connection.accepted" semantic
// conventions. It represents the reports the count of kernel NFS server TCP
// connections accepted.
type ServerNetTCPConnectionAccepted struct {
	metric.Int64Counter
}

var newServerNetTCPConnectionAcceptedOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS server TCP connections accepted."),
	metric.WithUnit("{connection}"),
}

// NewServerNetTCPConnectionAccepted returns a new ServerNetTCPConnectionAccepted
// instrument.
func NewServerNetTCPConnectionAccepted(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerNetTCPConnectionAccepted, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerNetTCPConnectionAccepted{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerNetTCPConnectionAcceptedOpts
	} else {
		opt = append(opt, newServerNetTCPConnectionAcceptedOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.net.tcp.connection.accepted",
		opt...,
	)
	if err != nil {
	    return ServerNetTCPConnectionAccepted{noop.Int64Counter{}}, err
	}
	return ServerNetTCPConnectionAccepted{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerNetTCPConnectionAccepted) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerNetTCPConnectionAccepted) Name() string {
	return "nfs.server.net.tcp.connection.accepted"
}

// Unit returns the semantic convention unit of the instrument
func (ServerNetTCPConnectionAccepted) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ServerNetTCPConnectionAccepted) Description() string {
	return "Reports the count of kernel NFS server TCP connections accepted."
}

// Add adds incr to the existing count for attrs.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.nettcpconn
func (m ServerNetTCPConnectionAccepted) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption)
	defer func() {
		*o = (*o)[:0]
		addOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributes(attrs...))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.nettcpconn
func (m ServerNetTCPConnectionAccepted) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// ServerOperationCount is an instrument used to record metric values conforming
// to the "nfs.server.operation.count" semantic conventions. It represents the
// reports the count of kernel NFSv4+ server operations.
type ServerOperationCount struct {
	metric.Int64Counter
}

var newServerOperationCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFSv4+ server operations."),
	metric.WithUnit("{operation}"),
}

// NewServerOperationCount returns a new ServerOperationCount instrument.
func NewServerOperationCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerOperationCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerOperationCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerOperationCountOpts
	} else {
		opt = append(opt, newServerOperationCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.operation.count",
		opt...,
	)
	if err != nil {
	    return ServerOperationCount{noop.Int64Counter{}}, err
	}
	return ServerOperationCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerOperationCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerOperationCount) Name() string {
	return "nfs.server.operation.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServerOperationCount) Unit() string {
	return "{operation}"
}

// Description returns the semantic convention description of the instrument
func (ServerOperationCount) Description() string {
	return "Reports the count of kernel NFSv4+ server operations."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m ServerOperationCount) Add(
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
func (m ServerOperationCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrOperationName returns an optional attribute for the "nfs.operation.name"
// semantic convention. It represents the NFSv4+ operation name.
func (ServerOperationCount) AttrOperationName(val string) attribute.KeyValue {
	return attribute.String("nfs.operation.name", val)
}

// AttrOncRPCVersion returns an optional attribute for the "onc_rpc.version"
// semantic convention. It represents the ONC/Sun RPC program version.
func (ServerOperationCount) AttrOncRPCVersion(val int) attribute.KeyValue {
	return attribute.Int("onc_rpc.version", val)
}

// ServerProcedureCount is an instrument used to record metric values conforming
// to the "nfs.server.procedure.count" semantic conventions. It represents the
// reports the count of kernel NFS server procedures.
type ServerProcedureCount struct {
	metric.Int64Counter
}

var newServerProcedureCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS server procedures."),
	metric.WithUnit("{procedure}"),
}

// NewServerProcedureCount returns a new ServerProcedureCount instrument.
func NewServerProcedureCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerProcedureCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerProcedureCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerProcedureCountOpts
	} else {
		opt = append(opt, newServerProcedureCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.procedure.count",
		opt...,
	)
	if err != nil {
	    return ServerProcedureCount{noop.Int64Counter{}}, err
	}
	return ServerProcedureCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerProcedureCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerProcedureCount) Name() string {
	return "nfs.server.procedure.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServerProcedureCount) Unit() string {
	return "{procedure}"
}

// Description returns the semantic convention description of the instrument
func (ServerProcedureCount) Description() string {
	return "Reports the count of kernel NFS server procedures."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m ServerProcedureCount) Add(
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
func (m ServerProcedureCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrOncRPCProcedureName returns an optional attribute for the
// "onc_rpc.procedure.name" semantic convention. It represents the ONC/Sun RPC
// procedure name.
func (ServerProcedureCount) AttrOncRPCProcedureName(val string) attribute.KeyValue {
	return attribute.String("onc_rpc.procedure.name", val)
}

// AttrOncRPCVersion returns an optional attribute for the "onc_rpc.version"
// semantic convention. It represents the ONC/Sun RPC program version.
func (ServerProcedureCount) AttrOncRPCVersion(val int) attribute.KeyValue {
	return attribute.Int("onc_rpc.version", val)
}

// ServerRepcacheRequests is an instrument used to record metric values
// conforming to the "nfs.server.repcache.requests" semantic conventions. It
// represents the reports the kernel NFS server reply cache request count by
// cache hit status.
type ServerRepcacheRequests struct {
	metric.Int64Counter
}

var newServerRepcacheRequestsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the kernel NFS server reply cache request count by cache hit status."),
	metric.WithUnit("{request}"),
}

// NewServerRepcacheRequests returns a new ServerRepcacheRequests instrument.
func NewServerRepcacheRequests(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerRepcacheRequests, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerRepcacheRequests{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerRepcacheRequestsOpts
	} else {
		opt = append(opt, newServerRepcacheRequestsOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.repcache.requests",
		opt...,
	)
	if err != nil {
	    return ServerRepcacheRequests{noop.Int64Counter{}}, err
	}
	return ServerRepcacheRequests{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerRepcacheRequests) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerRepcacheRequests) Name() string {
	return "nfs.server.repcache.requests"
}

// Unit returns the semantic convention unit of the instrument
func (ServerRepcacheRequests) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ServerRepcacheRequests) Description() string {
	return "Reports the kernel NFS server reply cache request count by cache hit status."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m ServerRepcacheRequests) Add(
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
func (m ServerRepcacheRequests) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrServerRepcacheStatus returns an optional attribute for the
// "nfs.server.repcache.status" semantic convention. It represents the linux: one
// of "hit" (NFSD_STATS_RC_HITS), "miss" (NFSD_STATS_RC_MISSES), or "nocache"
// (NFSD_STATS_RC_NOCACHE -- uncacheable).
func (ServerRepcacheRequests) AttrServerRepcacheStatus(val string) attribute.KeyValue {
	return attribute.String("nfs.server.repcache.status", val)
}

// ServerRPCCount is an instrument used to record metric values conforming to the
// "nfs.server.rpc.count" semantic conventions. It represents the reports the
// count of kernel NFS server RPCs handled.
type ServerRPCCount struct {
	metric.Int64Counter
}

var newServerRPCCountOpts = []metric.Int64CounterOption{
	metric.WithDescription("Reports the count of kernel NFS server RPCs handled."),
	metric.WithUnit("{request}"),
}

// NewServerRPCCount returns a new ServerRPCCount instrument.
func NewServerRPCCount(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ServerRPCCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerRPCCount{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerRPCCountOpts
	} else {
		opt = append(opt, newServerRPCCountOpts...)
	}

	i, err := m.Int64Counter(
		"nfs.server.rpc.count",
		opt...,
	)
	if err != nil {
	    return ServerRPCCount{noop.Int64Counter{}}, err
	}
	return ServerRPCCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerRPCCount) Inst() metric.Int64Counter {
	return m.Int64Counter
}

// Name returns the semantic convention name of the instrument.
func (ServerRPCCount) Name() string {
	return "nfs.server.rpc.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServerRPCCount) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ServerRPCCount) Description() string {
	return "Reports the count of kernel NFS server RPCs handled."
}

// Add adds incr to the existing count for attrs.
//
// All additional attrs passed are included in the recorded value.
//
// Linux: this metric is taken from the Linux kernel's svc_stat.rpccnt, the count
// of good RPCs. This metric can have
// an error.type of "format", "auth", or "client" for svc_stat.badfmt,
// svc_stat.badauth, and svc_stat.badclnt.
func (m ServerRPCCount) Add(
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
// Linux: this metric is taken from the Linux kernel's svc_stat.rpccnt, the count
// of good RPCs. This metric can have
// an error.type of "format", "auth", or "client" for svc_stat.badfmt,
// svc_stat.badauth, and svc_stat.badclnt.
func (m ServerRPCCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerRPCCount) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// ServerThreadCount is an instrument used to record metric values conforming to
// the "nfs.server.thread.count" semantic conventions. It represents the reports
// the count of kernel NFS server available threads.
type ServerThreadCount struct {
	metric.Int64UpDownCounter
}

var newServerThreadCountOpts = []metric.Int64UpDownCounterOption{
	metric.WithDescription("Reports the count of kernel NFS server available threads."),
	metric.WithUnit("{thread}"),
}

// NewServerThreadCount returns a new ServerThreadCount instrument.
func NewServerThreadCount(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ServerThreadCount, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerThreadCount{noop.Int64UpDownCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerThreadCountOpts
	} else {
		opt = append(opt, newServerThreadCountOpts...)
	}

	i, err := m.Int64UpDownCounter(
		"nfs.server.thread.count",
		opt...,
	)
	if err != nil {
	    return ServerThreadCount{noop.Int64UpDownCounter{}}, err
	}
	return ServerThreadCount{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerThreadCount) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
}

// Name returns the semantic convention name of the instrument.
func (ServerThreadCount) Name() string {
	return "nfs.server.thread.count"
}

// Unit returns the semantic convention unit of the instrument
func (ServerThreadCount) Unit() string {
	return "{thread}"
}

// Description returns the semantic convention description of the instrument
func (ServerThreadCount) Description() string {
	return "Reports the count of kernel NFS server available threads."
}

// Add adds incr to the existing count for attrs.
//
// Linux: this metric is taken from the Linux kernel nfsd_th_cnt variable
func (m ServerThreadCount) Add(ctx context.Context, incr int64, attrs ...attribute.KeyValue) {
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
// Linux: this metric is taken from the Linux kernel nfsd_th_cnt variable
func (m ServerThreadCount) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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