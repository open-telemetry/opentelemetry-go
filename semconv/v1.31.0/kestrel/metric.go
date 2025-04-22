// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/kestrel"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// ErrorTypeAttr is an attribute conforming to the error.type semantic
// conventions. It represents the full name of exception type.
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

// NetworkTypeAttr is an attribute conforming to the network.type semantic
// conventions. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
type NetworkTypeAttr string

var (
	// NetworkTypeIpv4 is the IPv4.
	NetworkTypeIpv4 NetworkTypeAttr = "ipv4"
	// NetworkTypeIpv6 is the IPv6.
	NetworkTypeIpv6 NetworkTypeAttr = "ipv6"
)

// ActiveConnections is an instrument used to record metric values conforming to
// the "kestrel.active_connections" semantic conventions. It represents the
// number of connections that are currently active on the server.
type ActiveConnections struct {
	inst metric.Int64UpDownCounter
}

// NewActiveConnections returns a new ActiveConnections instrument.
func NewActiveConnections(m metric.Meter) (ActiveConnections, error) {
	i, err := m.Int64UpDownCounter(
	    "kestrel.active_connections",
	    metric.WithDescription("Number of connections that are currently active on the server."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ActiveConnections{}, err
	}
	return ActiveConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ActiveConnections) Name() string {
	return "kestrel.active_connections"
}

// Unit returns the semantic convention unit of the instrument
func (ActiveConnections) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ActiveConnections) Description() string {
	return "Number of connections that are currently active on the server."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ActiveConnections) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ActiveConnections) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (ActiveConnections) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ActiveConnections) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ActiveConnections) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ActiveTLSHandshakes is an instrument used to record metric values conforming
// to the "kestrel.active_tls_handshakes" semantic conventions. It represents the
// number of TLS handshakes that are currently in progress on the server.
type ActiveTLSHandshakes struct {
	inst metric.Int64UpDownCounter
}

// NewActiveTLSHandshakes returns a new ActiveTLSHandshakes instrument.
func NewActiveTLSHandshakes(m metric.Meter) (ActiveTLSHandshakes, error) {
	i, err := m.Int64UpDownCounter(
	    "kestrel.active_tls_handshakes",
	    metric.WithDescription("Number of TLS handshakes that are currently in progress on the server."),
	    metric.WithUnit("{handshake}"),
	)
	if err != nil {
	    return ActiveTLSHandshakes{}, err
	}
	return ActiveTLSHandshakes{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ActiveTLSHandshakes) Name() string {
	return "kestrel.active_tls_handshakes"
}

// Unit returns the semantic convention unit of the instrument
func (ActiveTLSHandshakes) Unit() string {
	return "{handshake}"
}

// Description returns the semantic convention description of the instrument
func (ActiveTLSHandshakes) Description() string {
	return "Number of TLS handshakes that are currently in progress on the server."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ActiveTLSHandshakes) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ActiveTLSHandshakes) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (ActiveTLSHandshakes) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ActiveTLSHandshakes) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ActiveTLSHandshakes) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ConnectionDuration is an instrument used to record metric values conforming to
// the "kestrel.connection.duration" semantic conventions. It represents the
// duration of connections on the server.
type ConnectionDuration struct {
	inst metric.Float64Histogram
}

// NewConnectionDuration returns a new ConnectionDuration instrument.
func NewConnectionDuration(m metric.Meter) (ConnectionDuration, error) {
	i, err := m.Float64Histogram(
	    "kestrel.connection.duration",
	    metric.WithDescription("The duration of connections on the server."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ConnectionDuration{}, err
	}
	return ConnectionDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ConnectionDuration) Name() string {
	return "kestrel.connection.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ConnectionDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ConnectionDuration) Description() string {
	return "The duration of connections on the server."
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m ConnectionDuration) Record(
	ctx context.Context,
	val float64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the full name of exception type.
func (ConnectionDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ConnectionDuration) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ConnectionDuration) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ConnectionDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (ConnectionDuration) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ConnectionDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ConnectionDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrTLSProtocolVersion returns an optional attribute for the
// "tls.protocol.version" semantic convention. It represents the numeric part of
// the version parsed from the original string of the negotiated
// [SSL/TLS protocol version].
//
// [SSL/TLS protocol version]: https://docs.openssl.org/1.1.1/man3/SSL_get_version/#return-values
func (ConnectionDuration) AttrTLSProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("tls.protocol.version", val)
}

// QueuedConnections is an instrument used to record metric values conforming to
// the "kestrel.queued_connections" semantic conventions. It represents the
// number of connections that are currently queued and are waiting to start.
type QueuedConnections struct {
	inst metric.Int64UpDownCounter
}

// NewQueuedConnections returns a new QueuedConnections instrument.
func NewQueuedConnections(m metric.Meter) (QueuedConnections, error) {
	i, err := m.Int64UpDownCounter(
	    "kestrel.queued_connections",
	    metric.WithDescription("Number of connections that are currently queued and are waiting to start."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return QueuedConnections{}, err
	}
	return QueuedConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (QueuedConnections) Name() string {
	return "kestrel.queued_connections"
}

// Unit returns the semantic convention unit of the instrument
func (QueuedConnections) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (QueuedConnections) Description() string {
	return "Number of connections that are currently queued and are waiting to start."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m QueuedConnections) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (QueuedConnections) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (QueuedConnections) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (QueuedConnections) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (QueuedConnections) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// QueuedRequests is an instrument used to record metric values conforming to the
// "kestrel.queued_requests" semantic conventions. It represents the number of
// HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are
// currently queued and are waiting to start.
type QueuedRequests struct {
	inst metric.Int64UpDownCounter
}

// NewQueuedRequests returns a new QueuedRequests instrument.
func NewQueuedRequests(m metric.Meter) (QueuedRequests, error) {
	i, err := m.Int64UpDownCounter(
	    "kestrel.queued_requests",
	    metric.WithDescription("Number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are currently queued and are waiting to start."),
	    metric.WithUnit("{request}"),
	)
	if err != nil {
	    return QueuedRequests{}, err
	}
	return QueuedRequests{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (QueuedRequests) Name() string {
	return "kestrel.queued_requests"
}

// Unit returns the semantic convention unit of the instrument
func (QueuedRequests) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (QueuedRequests) Description() string {
	return "Number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that are currently queued and are waiting to start."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m QueuedRequests) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (QueuedRequests) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (QueuedRequests) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (QueuedRequests) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (QueuedRequests) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (QueuedRequests) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (QueuedRequests) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// RejectedConnections is an instrument used to record metric values conforming
// to the "kestrel.rejected_connections" semantic conventions. It represents the
// number of connections rejected by the server.
type RejectedConnections struct {
	inst metric.Int64Counter
}

// NewRejectedConnections returns a new RejectedConnections instrument.
func NewRejectedConnections(m metric.Meter) (RejectedConnections, error) {
	i, err := m.Int64Counter(
	    "kestrel.rejected_connections",
	    metric.WithDescription("Number of connections rejected by the server."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return RejectedConnections{}, err
	}
	return RejectedConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (RejectedConnections) Name() string {
	return "kestrel.rejected_connections"
}

// Unit returns the semantic convention unit of the instrument
func (RejectedConnections) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (RejectedConnections) Description() string {
	return "Number of connections rejected by the server."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m RejectedConnections) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (RejectedConnections) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (RejectedConnections) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (RejectedConnections) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (RejectedConnections) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// TLSHandshakeDuration is an instrument used to record metric values conforming
// to the "kestrel.tls_handshake.duration" semantic conventions. It represents
// the duration of TLS handshakes on the server.
type TLSHandshakeDuration struct {
	inst metric.Float64Histogram
}

// NewTLSHandshakeDuration returns a new TLSHandshakeDuration instrument.
func NewTLSHandshakeDuration(m metric.Meter) (TLSHandshakeDuration, error) {
	i, err := m.Float64Histogram(
	    "kestrel.tls_handshake.duration",
	    metric.WithDescription("The duration of TLS handshakes on the server."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return TLSHandshakeDuration{}, err
	}
	return TLSHandshakeDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (TLSHandshakeDuration) Name() string {
	return "kestrel.tls_handshake.duration"
}

// Unit returns the semantic convention unit of the instrument
func (TLSHandshakeDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (TLSHandshakeDuration) Description() string {
	return "The duration of TLS handshakes on the server."
}

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
func (m TLSHandshakeDuration) Record(
	ctx context.Context,
	val float64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the full name of exception type.
func (TLSHandshakeDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (TLSHandshakeDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (TLSHandshakeDuration) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (TLSHandshakeDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (TLSHandshakeDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrTLSProtocolVersion returns an optional attribute for the
// "tls.protocol.version" semantic convention. It represents the numeric part of
// the version parsed from the original string of the negotiated
// [SSL/TLS protocol version].
//
// [SSL/TLS protocol version]: https://docs.openssl.org/1.1.1/man3/SSL_get_version/#return-values
func (TLSHandshakeDuration) AttrTLSProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("tls.protocol.version", val)
}

// UpgradedConnections is an instrument used to record metric values conforming
// to the "kestrel.upgraded_connections" semantic conventions. It represents the
// number of connections that are currently upgraded (WebSockets). .
type UpgradedConnections struct {
	inst metric.Int64UpDownCounter
}

// NewUpgradedConnections returns a new UpgradedConnections instrument.
func NewUpgradedConnections(m metric.Meter) (UpgradedConnections, error) {
	i, err := m.Int64UpDownCounter(
	    "kestrel.upgraded_connections",
	    metric.WithDescription("Number of connections that are currently upgraded (WebSockets). ."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return UpgradedConnections{}, err
	}
	return UpgradedConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (UpgradedConnections) Name() string {
	return "kestrel.upgraded_connections"
}

// Unit returns the semantic convention unit of the instrument
func (UpgradedConnections) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (UpgradedConnections) Description() string {
	return "Number of connections that are currently upgraded (WebSockets). ."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m UpgradedConnections) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			attrs...,
		),
	)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (UpgradedConnections) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrNetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (UpgradedConnections) AttrNetworkType(val NetworkTypeAttr) attribute.KeyValue {
	return attribute.String("network.type", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (UpgradedConnections) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (UpgradedConnections) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}