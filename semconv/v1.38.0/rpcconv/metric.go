// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package rpcconv provides types and functionality for OpenTelemetry semantic
// conventions in the "rpc" namespace.
package rpcconv

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

// SystemAttr is an attribute conforming to the rpc.system semantic conventions.
// It represents a string identifying the remoting system. See below for a list
// of well-known identifiers.
type SystemAttr string

var (
	// SystemGRPC is the gRPC.
	SystemGRPC SystemAttr = "grpc"
	// SystemJavaRmi is the java RMI.
	SystemJavaRmi SystemAttr = "java_rmi"
	// SystemDotnetWcf is the NET WCF.
	SystemDotnetWcf SystemAttr = "dotnet_wcf"
	// SystemApacheDubbo is the apache Dubbo.
	SystemApacheDubbo SystemAttr = "apache_dubbo"
	// SystemConnectRPC is the connect RPC.
	SystemConnectRPC SystemAttr = "connect_rpc"
	// SystemOncRPC is the [ONC RPC (Sun RPC)].
	//
	// [ONC RPC (Sun RPC)]: https://datatracker.ietf.org/doc/html/rfc5531
	SystemOncRPC SystemAttr = "onc_rpc"
	// SystemJSONRPC is the JSON-RPC.
	SystemJSONRPC SystemAttr = "jsonrpc"
)

// ClientDuration is an instrument used to record metric values conforming to the
// "rpc.client.duration" semantic conventions. It represents the measures the
// duration of outbound RPC.
type ClientDuration struct {
	metric.Float64Histogram
}

var newClientDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("Measures the duration of outbound RPC."),
	metric.WithUnit("ms"),
}

// NewClientDuration returns a new ClientDuration instrument.
func NewClientDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientDurationOpts
	} else {
		opt = append(opt, newClientDurationOpts...)
	}

	i, err := m.Float64Histogram(
		"rpc.client.duration",
		opt...,
	)
	if err != nil {
		return ClientDuration{noop.Float64Histogram{}}, err
	}
	return ClientDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientDuration) Name() string {
	return "rpc.client.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientDuration) Unit() string {
	return "ms"
}

// Description returns the semantic convention description of the instrument
func (ClientDuration) Description() string {
	return "Measures the duration of outbound RPC."
}

// Record records val to the current distribution for attrs.
//
// The system is the a string identifying the remoting system. See below for a
// list of well-known identifiers.
//
// The serverAddress is the server domain name if available without reverse DNS
// lookup; otherwise, IP address or Unix domain socket name.
//
// All additional attrs passed are included in the recorded value.
//
// While streaming RPCs may record this metric as start-of-batch
// to end-of-batch, it's hard to interpret in practice.
//
// **Streaming**: N/A.
func (m ClientDuration) Record(
	ctx context.Context,
	val float64,
	system SystemAttr,
	serverAddress string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val)
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
			append(
				attrs,
				attribute.String("rpc.system", string(system)),
				attribute.String("server.address", serverAddress),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// While streaming RPCs may record this metric as start-of-batch
// to end-of-batch, it's hard to interpret in practice.
//
// **Streaming**: N/A.
func (m ClientDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientDuration) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientDuration) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ClientDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrMethod returns an optional attribute for the "rpc.method" semantic
// convention. It represents the this is the logical name of the method from the
// RPC interface perspective.
func (ClientDuration) AttrMethod(val string) attribute.KeyValue {
	return attribute.String("rpc.method", val)
}

// AttrService returns an optional attribute for the "rpc.service" semantic
// convention. It represents the full (logical) name of the service being called,
// including its package name, if applicable.
func (ClientDuration) AttrService(val string) attribute.KeyValue {
	return attribute.String("rpc.service", val)
}

// ClientRequestSize is an instrument used to record metric values conforming to
// the "rpc.client.request.size" semantic conventions. It represents the measures
// the size of RPC request messages (uncompressed).
type ClientRequestSize struct {
	metric.Int64Histogram
}

var newClientRequestSizeOpts = []metric.Int64HistogramOption{
	metric.WithDescription("Measures the size of RPC request messages (uncompressed)."),
	metric.WithUnit("By"),
}

// NewClientRequestSize returns a new ClientRequestSize instrument.
func NewClientRequestSize(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (ClientRequestSize, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientRequestSize{noop.Int64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientRequestSizeOpts
	} else {
		opt = append(opt, newClientRequestSizeOpts...)
	}

	i, err := m.Int64Histogram(
		"rpc.client.request.size",
		opt...,
	)
	if err != nil {
		return ClientRequestSize{noop.Int64Histogram{}}, err
	}
	return ClientRequestSize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientRequestSize) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientRequestSize) Name() string {
	return "rpc.client.request.size"
}

// Unit returns the semantic convention unit of the instrument
func (ClientRequestSize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClientRequestSize) Description() string {
	return "Measures the size of RPC request messages (uncompressed)."
}

// Record records val to the current distribution for attrs.
//
// The system is the a string identifying the remoting system. See below for a
// list of well-known identifiers.
//
// The serverAddress is the server domain name if available without reverse DNS
// lookup; otherwise, IP address or Unix domain socket name.
//
// All additional attrs passed are included in the recorded value.
//
// **Streaming**: Recorded per message in a streaming batch
func (m ClientRequestSize) Record(
	ctx context.Context,
	val int64,
	system SystemAttr,
	serverAddress string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Histogram.Record(ctx, val)
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
			append(
				attrs,
				attribute.String("rpc.system", string(system)),
				attribute.String("server.address", serverAddress),
			)...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// **Streaming**: Recorded per message in a streaming batch
func (m ClientRequestSize) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientRequestSize) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientRequestSize) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientRequestSize) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientRequestSize) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ClientRequestSize) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrMethod returns an optional attribute for the "rpc.method" semantic
// convention. It represents the this is the logical name of the method from the
// RPC interface perspective.
func (ClientRequestSize) AttrMethod(val string) attribute.KeyValue {
	return attribute.String("rpc.method", val)
}

// AttrService returns an optional attribute for the "rpc.service" semantic
// convention. It represents the full (logical) name of the service being called,
// including its package name, if applicable.
func (ClientRequestSize) AttrService(val string) attribute.KeyValue {
	return attribute.String("rpc.service", val)
}

// ClientResponseSize is an instrument used to record metric values conforming to
// the "rpc.client.response.size" semantic conventions. It represents the
// measures the size of RPC response messages (uncompressed).
type ClientResponseSize struct {
	metric.Int64Histogram
}

var newClientResponseSizeOpts = []metric.Int64HistogramOption{
	metric.WithDescription("Measures the size of RPC response messages (uncompressed)."),
	metric.WithUnit("By"),
}

// NewClientResponseSize returns a new ClientResponseSize instrument.
func NewClientResponseSize(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (ClientResponseSize, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientResponseSize{noop.Int64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientResponseSizeOpts
	} else {
		opt = append(opt, newClientResponseSizeOpts...)
	}

	i, err := m.Int64Histogram(
		"rpc.client.response.size",
		opt...,
	)
	if err != nil {
		return ClientResponseSize{noop.Int64Histogram{}}, err
	}
	return ClientResponseSize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientResponseSize) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientResponseSize) Name() string {
	return "rpc.client.response.size"
}

// Unit returns the semantic convention unit of the instrument
func (ClientResponseSize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClientResponseSize) Description() string {
	return "Measures the size of RPC response messages (uncompressed)."
}

// Record records val to the current distribution for attrs.
//
// The system is the a string identifying the remoting system. See below for a
// list of well-known identifiers.
//
// The serverAddress is the server domain name if available without reverse DNS
// lookup; otherwise, IP address or Unix domain socket name.
//
// All additional attrs passed are included in the recorded value.
//
// **Streaming**: Recorded per response in a streaming batch
func (m ClientResponseSize) Record(
	ctx context.Context,
	val int64,
	system SystemAttr,
	serverAddress string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Histogram.Record(ctx, val)
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
			append(
				attrs,
				attribute.String("rpc.system", string(system)),
				attribute.String("server.address", serverAddress),
			)...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// **Streaming**: Recorded per response in a streaming batch
func (m ClientResponseSize) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientResponseSize) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientResponseSize) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientResponseSize) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientResponseSize) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ClientResponseSize) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrMethod returns an optional attribute for the "rpc.method" semantic
// convention. It represents the this is the logical name of the method from the
// RPC interface perspective.
func (ClientResponseSize) AttrMethod(val string) attribute.KeyValue {
	return attribute.String("rpc.method", val)
}

// AttrService returns an optional attribute for the "rpc.service" semantic
// convention. It represents the full (logical) name of the service being called,
// including its package name, if applicable.
func (ClientResponseSize) AttrService(val string) attribute.KeyValue {
	return attribute.String("rpc.service", val)
}

// ServerDuration is an instrument used to record metric values conforming to the
// "rpc.server.duration" semantic conventions. It represents the measures the
// duration of inbound RPC.
type ServerDuration struct {
	metric.Float64Histogram
}

var newServerDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("Measures the duration of inbound RPC."),
	metric.WithUnit("ms"),
}

// NewServerDuration returns a new ServerDuration instrument.
func NewServerDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ServerDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerDurationOpts
	} else {
		opt = append(opt, newServerDurationOpts...)
	}

	i, err := m.Float64Histogram(
		"rpc.server.duration",
		opt...,
	)
	if err != nil {
		return ServerDuration{noop.Float64Histogram{}}, err
	}
	return ServerDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ServerDuration) Name() string {
	return "rpc.server.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ServerDuration) Unit() string {
	return "ms"
}

// Description returns the semantic convention description of the instrument
func (ServerDuration) Description() string {
	return "Measures the duration of inbound RPC."
}

// Record records val to the current distribution for attrs.
//
// The system is the a string identifying the remoting system. See below for a
// list of well-known identifiers.
//
// All additional attrs passed are included in the recorded value.
//
// While streaming RPCs may record this metric as start-of-batch
// to end-of-batch, it's hard to interpret in practice.
//
// **Streaming**: N/A.
func (m ServerDuration) Record(
	ctx context.Context,
	val float64,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val)
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
			append(
				attrs,
				attribute.String("rpc.system", string(system)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// While streaming RPCs may record this metric as start-of-batch
// to end-of-batch, it's hard to interpret in practice.
//
// **Streaming**: N/A.
func (m ServerDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerDuration) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerDuration) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ServerDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrMethod returns an optional attribute for the "rpc.method" semantic
// convention. It represents the this is the logical name of the method from the
// RPC interface perspective.
func (ServerDuration) AttrMethod(val string) attribute.KeyValue {
	return attribute.String("rpc.method", val)
}

// AttrService returns an optional attribute for the "rpc.service" semantic
// convention. It represents the full (logical) name of the service being called,
// including its package name, if applicable.
func (ServerDuration) AttrService(val string) attribute.KeyValue {
	return attribute.String("rpc.service", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ServerDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ServerDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ServerRequestSize is an instrument used to record metric values conforming to
// the "rpc.server.request.size" semantic conventions. It represents the measures
// the size of RPC request messages (uncompressed).
type ServerRequestSize struct {
	metric.Int64Histogram
}

var newServerRequestSizeOpts = []metric.Int64HistogramOption{
	metric.WithDescription("Measures the size of RPC request messages (uncompressed)."),
	metric.WithUnit("By"),
}

// NewServerRequestSize returns a new ServerRequestSize instrument.
func NewServerRequestSize(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (ServerRequestSize, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerRequestSize{noop.Int64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerRequestSizeOpts
	} else {
		opt = append(opt, newServerRequestSizeOpts...)
	}

	i, err := m.Int64Histogram(
		"rpc.server.request.size",
		opt...,
	)
	if err != nil {
		return ServerRequestSize{noop.Int64Histogram{}}, err
	}
	return ServerRequestSize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerRequestSize) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ServerRequestSize) Name() string {
	return "rpc.server.request.size"
}

// Unit returns the semantic convention unit of the instrument
func (ServerRequestSize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ServerRequestSize) Description() string {
	return "Measures the size of RPC request messages (uncompressed)."
}

// Record records val to the current distribution for attrs.
//
// The system is the a string identifying the remoting system. See below for a
// list of well-known identifiers.
//
// All additional attrs passed are included in the recorded value.
//
// **Streaming**: Recorded per message in a streaming batch
func (m ServerRequestSize) Record(
	ctx context.Context,
	val int64,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Histogram.Record(ctx, val)
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
			append(
				attrs,
				attribute.String("rpc.system", string(system)),
			)...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// **Streaming**: Recorded per message in a streaming batch
func (m ServerRequestSize) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerRequestSize) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerRequestSize) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerRequestSize) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ServerRequestSize) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrMethod returns an optional attribute for the "rpc.method" semantic
// convention. It represents the this is the logical name of the method from the
// RPC interface perspective.
func (ServerRequestSize) AttrMethod(val string) attribute.KeyValue {
	return attribute.String("rpc.method", val)
}

// AttrService returns an optional attribute for the "rpc.service" semantic
// convention. It represents the full (logical) name of the service being called,
// including its package name, if applicable.
func (ServerRequestSize) AttrService(val string) attribute.KeyValue {
	return attribute.String("rpc.service", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ServerRequestSize) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ServerRequestSize) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ServerResponseSize is an instrument used to record metric values conforming to
// the "rpc.server.response.size" semantic conventions. It represents the
// measures the size of RPC response messages (uncompressed).
type ServerResponseSize struct {
	metric.Int64Histogram
}

var newServerResponseSizeOpts = []metric.Int64HistogramOption{
	metric.WithDescription("Measures the size of RPC response messages (uncompressed)."),
	metric.WithUnit("By"),
}

// NewServerResponseSize returns a new ServerResponseSize instrument.
func NewServerResponseSize(
	m metric.Meter,
	opt ...metric.Int64HistogramOption,
) (ServerResponseSize, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerResponseSize{noop.Int64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerResponseSizeOpts
	} else {
		opt = append(opt, newServerResponseSizeOpts...)
	}

	i, err := m.Int64Histogram(
		"rpc.server.response.size",
		opt...,
	)
	if err != nil {
		return ServerResponseSize{noop.Int64Histogram{}}, err
	}
	return ServerResponseSize{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerResponseSize) Inst() metric.Int64Histogram {
	return m.Int64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ServerResponseSize) Name() string {
	return "rpc.server.response.size"
}

// Unit returns the semantic convention unit of the instrument
func (ServerResponseSize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ServerResponseSize) Description() string {
	return "Measures the size of RPC response messages (uncompressed)."
}

// Record records val to the current distribution for attrs.
//
// The system is the a string identifying the remoting system. See below for a
// list of well-known identifiers.
//
// All additional attrs passed are included in the recorded value.
//
// **Streaming**: Recorded per response in a streaming batch
func (m ServerResponseSize) Record(
	ctx context.Context,
	val int64,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		m.Int64Histogram.Record(ctx, val)
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
			append(
				attrs,
				attribute.String("rpc.system", string(system)),
			)...,
		),
	)

	m.Int64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
//
// **Streaming**: Recorded per response in a streaming batch
func (m ServerResponseSize) RecordSet(ctx context.Context, val int64, set attribute.Set) {
	if set.Len() == 0 {
		m.Int64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption)
	defer func() {
		*o = (*o)[:0]
		recOptPool.Put(o)
	}()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerResponseSize) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerResponseSize) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerResponseSize) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ServerResponseSize) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrMethod returns an optional attribute for the "rpc.method" semantic
// convention. It represents the this is the logical name of the method from the
// RPC interface perspective.
func (ServerResponseSize) AttrMethod(val string) attribute.KeyValue {
	return attribute.String("rpc.method", val)
}

// AttrService returns an optional attribute for the "rpc.service" semantic
// convention. It represents the full (logical) name of the service being called,
// including its package name, if applicable.
func (ServerResponseSize) AttrService(val string) attribute.KeyValue {
	return attribute.String("rpc.service", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ServerResponseSize) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ServerResponseSize) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}
