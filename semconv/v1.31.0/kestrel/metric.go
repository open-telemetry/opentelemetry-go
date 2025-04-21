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

// KestrelActiveConnections is an instrument used to record metric values
// conforming to the "kestrel.active_connections" semantic conventions. It
// represents the number of connections that are currently active on the server.
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
	attrs ...ActiveConnectionsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ActiveConnections) conv(in []ActiveConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.activeConnectionsAttr()
	}
	return out
}

// ActiveConnectionsAttr is an optional attribute for the ActiveConnections
// instrument.
type ActiveConnectionsAttr interface {
    activeConnectionsAttr() attribute.KeyValue
}

type activeConnectionsAttr struct {
	kv attribute.KeyValue
}

func (a activeConnectionsAttr) activeConnectionsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ActiveConnections) NetworkTransport(val NetworkTransportAttr) ActiveConnectionsAttr {
	return activeConnectionsAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (ActiveConnections) NetworkType(val NetworkTypeAttr) ActiveConnectionsAttr {
	return activeConnectionsAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (ActiveConnections) ServerAddress(val string) ActiveConnectionsAttr {
	return activeConnectionsAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ActiveConnections) ServerPort(val int) ActiveConnectionsAttr {
	return activeConnectionsAttr{kv: attribute.Int("server.port", val)}
}

// KestrelActiveTLSHandshakes is an instrument used to record metric values
// conforming to the "kestrel.active_tls_handshakes" semantic conventions. It
// represents the number of TLS handshakes that are currently in progress on the
// server.
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
	attrs ...ActiveTLSHandshakesAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ActiveTLSHandshakes) conv(in []ActiveTLSHandshakesAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.activeTlsHandshakesAttr()
	}
	return out
}

// ActiveTLSHandshakesAttr is an optional attribute for the ActiveTLSHandshakes
// instrument.
type ActiveTLSHandshakesAttr interface {
    activeTlsHandshakesAttr() attribute.KeyValue
}

type activeTlsHandshakesAttr struct {
	kv attribute.KeyValue
}

func (a activeTlsHandshakesAttr) activeTlsHandshakesAttr() attribute.KeyValue {
    return a.kv
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ActiveTLSHandshakes) NetworkTransport(val NetworkTransportAttr) ActiveTLSHandshakesAttr {
	return activeTlsHandshakesAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (ActiveTLSHandshakes) NetworkType(val NetworkTypeAttr) ActiveTLSHandshakesAttr {
	return activeTlsHandshakesAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (ActiveTLSHandshakes) ServerAddress(val string) ActiveTLSHandshakesAttr {
	return activeTlsHandshakesAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ActiveTLSHandshakes) ServerPort(val int) ActiveTLSHandshakesAttr {
	return activeTlsHandshakesAttr{kv: attribute.Int("server.port", val)}
}

// KestrelConnectionDuration is an instrument used to record metric values
// conforming to the "kestrel.connection.duration" semantic conventions. It
// represents the duration of connections on the server.
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

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ConnectionDuration) Record(
    ctx context.Context,
    val float64,
	attrs ...ConnectionDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ConnectionDuration) conv(in []ConnectionDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.connectionDurationAttr()
	}
	return out
}

// ConnectionDurationAttr is an optional attribute for the ConnectionDuration
// instrument.
type ConnectionDurationAttr interface {
    connectionDurationAttr() attribute.KeyValue
}

type connectionDurationAttr struct {
	kv attribute.KeyValue
}

func (a connectionDurationAttr) connectionDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the full name of exception type.
func (ConnectionDuration) ErrorType(val ErrorTypeAttr) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.String("error.type", string(val))}
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ConnectionDuration) NetworkProtocolName(val string) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.String("network.protocol.name", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ConnectionDuration) NetworkProtocolVersion(val string) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.String("network.protocol.version", val)}
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (ConnectionDuration) NetworkTransport(val NetworkTransportAttr) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (ConnectionDuration) NetworkType(val NetworkTypeAttr) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (ConnectionDuration) ServerAddress(val string) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ConnectionDuration) ServerPort(val int) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.Int("server.port", val)}
}

// TLSProtocolVersion returns an optional attribute for the
// "tls.protocol.version" semantic convention. It represents the numeric part of
// the version parsed from the original string of the negotiated
// [SSL/TLS protocol version].
//
// [SSL/TLS protocol version]: https://docs.openssl.org/1.1.1/man3/SSL_get_version/#return-values
func (ConnectionDuration) TLSProtocolVersion(val string) ConnectionDurationAttr {
	return connectionDurationAttr{kv: attribute.String("tls.protocol.version", val)}
}

// KestrelQueuedConnections is an instrument used to record metric values
// conforming to the "kestrel.queued_connections" semantic conventions. It
// represents the number of connections that are currently queued and are waiting
// to start.
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
	attrs ...QueuedConnectionsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m QueuedConnections) conv(in []QueuedConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.queuedConnectionsAttr()
	}
	return out
}

// QueuedConnectionsAttr is an optional attribute for the QueuedConnections
// instrument.
type QueuedConnectionsAttr interface {
    queuedConnectionsAttr() attribute.KeyValue
}

type queuedConnectionsAttr struct {
	kv attribute.KeyValue
}

func (a queuedConnectionsAttr) queuedConnectionsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (QueuedConnections) NetworkTransport(val NetworkTransportAttr) QueuedConnectionsAttr {
	return queuedConnectionsAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (QueuedConnections) NetworkType(val NetworkTypeAttr) QueuedConnectionsAttr {
	return queuedConnectionsAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (QueuedConnections) ServerAddress(val string) QueuedConnectionsAttr {
	return queuedConnectionsAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (QueuedConnections) ServerPort(val int) QueuedConnectionsAttr {
	return queuedConnectionsAttr{kv: attribute.Int("server.port", val)}
}

// KestrelQueuedRequests is an instrument used to record metric values conforming
// to the "kestrel.queued_requests" semantic conventions. It represents the
// number of HTTP requests on multiplexed connections (HTTP/2 and HTTP/3) that
// are currently queued and are waiting to start.
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
	attrs ...QueuedRequestsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m QueuedRequests) conv(in []QueuedRequestsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.queuedRequestsAttr()
	}
	return out
}

// QueuedRequestsAttr is an optional attribute for the QueuedRequests instrument.
type QueuedRequestsAttr interface {
    queuedRequestsAttr() attribute.KeyValue
}

type queuedRequestsAttr struct {
	kv attribute.KeyValue
}

func (a queuedRequestsAttr) queuedRequestsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (QueuedRequests) NetworkProtocolName(val string) QueuedRequestsAttr {
	return queuedRequestsAttr{kv: attribute.String("network.protocol.name", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (QueuedRequests) NetworkProtocolVersion(val string) QueuedRequestsAttr {
	return queuedRequestsAttr{kv: attribute.String("network.protocol.version", val)}
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (QueuedRequests) NetworkTransport(val NetworkTransportAttr) QueuedRequestsAttr {
	return queuedRequestsAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (QueuedRequests) NetworkType(val NetworkTypeAttr) QueuedRequestsAttr {
	return queuedRequestsAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (QueuedRequests) ServerAddress(val string) QueuedRequestsAttr {
	return queuedRequestsAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (QueuedRequests) ServerPort(val int) QueuedRequestsAttr {
	return queuedRequestsAttr{kv: attribute.Int("server.port", val)}
}

// KestrelRejectedConnections is an instrument used to record metric values
// conforming to the "kestrel.rejected_connections" semantic conventions. It
// represents the number of connections rejected by the server.
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
	attrs ...RejectedConnectionsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m RejectedConnections) conv(in []RejectedConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.rejectedConnectionsAttr()
	}
	return out
}

// RejectedConnectionsAttr is an optional attribute for the RejectedConnections
// instrument.
type RejectedConnectionsAttr interface {
    rejectedConnectionsAttr() attribute.KeyValue
}

type rejectedConnectionsAttr struct {
	kv attribute.KeyValue
}

func (a rejectedConnectionsAttr) rejectedConnectionsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (RejectedConnections) NetworkTransport(val NetworkTransportAttr) RejectedConnectionsAttr {
	return rejectedConnectionsAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (RejectedConnections) NetworkType(val NetworkTypeAttr) RejectedConnectionsAttr {
	return rejectedConnectionsAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (RejectedConnections) ServerAddress(val string) RejectedConnectionsAttr {
	return rejectedConnectionsAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (RejectedConnections) ServerPort(val int) RejectedConnectionsAttr {
	return rejectedConnectionsAttr{kv: attribute.Int("server.port", val)}
}

// KestrelTLSHandshakeDuration is an instrument used to record metric values
// conforming to the "kestrel.tls_handshake.duration" semantic conventions. It
// represents the duration of TLS handshakes on the server.
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

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m TLSHandshakeDuration) Record(
    ctx context.Context,
    val float64,
	attrs ...TLSHandshakeDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m TLSHandshakeDuration) conv(in []TLSHandshakeDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.tlsHandshakeDurationAttr()
	}
	return out
}

// TLSHandshakeDurationAttr is an optional attribute for the TLSHandshakeDuration
// instrument.
type TLSHandshakeDurationAttr interface {
    tlsHandshakeDurationAttr() attribute.KeyValue
}

type tlsHandshakeDurationAttr struct {
	kv attribute.KeyValue
}

func (a tlsHandshakeDurationAttr) tlsHandshakeDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the full name of exception type.
func (TLSHandshakeDuration) ErrorType(val ErrorTypeAttr) TLSHandshakeDurationAttr {
	return tlsHandshakeDurationAttr{kv: attribute.String("error.type", string(val))}
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (TLSHandshakeDuration) NetworkTransport(val NetworkTransportAttr) TLSHandshakeDurationAttr {
	return tlsHandshakeDurationAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (TLSHandshakeDuration) NetworkType(val NetworkTypeAttr) TLSHandshakeDurationAttr {
	return tlsHandshakeDurationAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (TLSHandshakeDuration) ServerAddress(val string) TLSHandshakeDurationAttr {
	return tlsHandshakeDurationAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (TLSHandshakeDuration) ServerPort(val int) TLSHandshakeDurationAttr {
	return tlsHandshakeDurationAttr{kv: attribute.Int("server.port", val)}
}

// TLSProtocolVersion returns an optional attribute for the
// "tls.protocol.version" semantic convention. It represents the numeric part of
// the version parsed from the original string of the negotiated
// [SSL/TLS protocol version].
//
// [SSL/TLS protocol version]: https://docs.openssl.org/1.1.1/man3/SSL_get_version/#return-values
func (TLSHandshakeDuration) TLSProtocolVersion(val string) TLSHandshakeDurationAttr {
	return tlsHandshakeDurationAttr{kv: attribute.String("tls.protocol.version", val)}
}

// KestrelUpgradedConnections is an instrument used to record metric values
// conforming to the "kestrel.upgraded_connections" semantic conventions. It
// represents the number of connections that are currently upgraded (WebSockets).
// .
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
	attrs ...UpgradedConnectionsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m UpgradedConnections) conv(in []UpgradedConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.upgradedConnectionsAttr()
	}
	return out
}

// UpgradedConnectionsAttr is an optional attribute for the UpgradedConnections
// instrument.
type UpgradedConnectionsAttr interface {
    upgradedConnectionsAttr() attribute.KeyValue
}

type upgradedConnectionsAttr struct {
	kv attribute.KeyValue
}

func (a upgradedConnectionsAttr) upgradedConnectionsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the [OSI transport layer] or
// [inter-process communication method].
//
// [OSI transport layer]: https://wikipedia.org/wiki/Transport_layer
// [inter-process communication method]: https://wikipedia.org/wiki/Inter-process_communication
func (UpgradedConnections) NetworkTransport(val NetworkTransportAttr) UpgradedConnectionsAttr {
	return upgradedConnectionsAttr{kv: attribute.String("network.transport", string(val))}
}

// NetworkType returns an optional attribute for the "network.type" semantic
// convention. It represents the [OSI network layer] or non-OSI equivalent.
//
// [OSI network layer]: https://wikipedia.org/wiki/Network_layer
func (UpgradedConnections) NetworkType(val NetworkTypeAttr) UpgradedConnectionsAttr {
	return upgradedConnectionsAttr{kv: attribute.String("network.type", string(val))}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the server domain name if available without reverse
// DNS lookup; otherwise, IP address or Unix domain socket name.
func (UpgradedConnections) ServerAddress(val string) UpgradedConnectionsAttr {
	return upgradedConnectionsAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (UpgradedConnections) ServerPort(val int) UpgradedConnectionsAttr {
	return upgradedConnectionsAttr{kv: attribute.Int("server.port", val)}
}