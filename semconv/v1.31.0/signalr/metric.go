// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/signalr"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// ConnectionStatusAttr is an attribute conforming to the
// signalr.connection.status semantic conventions. It represents the signalR HTTP
// connection closure status.
type ConnectionStatusAttr string

var (
	// ConnectionStatusNormalClosure is the connection was closed normally.
	ConnectionStatusNormalClosure ConnectionStatusAttr = "normal_closure"
	// ConnectionStatusTimeout is the connection was closed due to a timeout.
	ConnectionStatusTimeout ConnectionStatusAttr = "timeout"
	// ConnectionStatusAppShutdown is the connection was closed because the app is
	// shutting down.
	ConnectionStatusAppShutdown ConnectionStatusAttr = "app_shutdown"
)

// TransportAttr is an attribute conforming to the signalr.transport semantic
// conventions. It represents the [SignalR transport type].
//
// [SignalR transport type]: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md
type TransportAttr string

var (
	// TransportServerSentEvents is the serverSentEvents protocol.
	TransportServerSentEvents TransportAttr = "server_sent_events"
	// TransportLongPolling is the longPolling protocol.
	TransportLongPolling TransportAttr = "long_polling"
	// TransportWebSockets is the webSockets protocol.
	TransportWebSockets TransportAttr = "web_sockets"
)

// SignalrServerActiveConnections is an instrument used to record metric values
// conforming to the "signalr.server.active_connections" semantic conventions. It
// represents the number of connections that are currently active on the server.
type ServerActiveConnections struct {
	inst metric.Int64UpDownCounter
}

// NewServerActiveConnections returns a new ServerActiveConnections instrument.
func NewServerActiveConnections(m metric.Meter) (ServerActiveConnections, error) {
	i, err := m.Int64UpDownCounter(
	    "signalr.server.active_connections",
	    metric.WithDescription("Number of connections that are currently active on the server."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ServerActiveConnections{}, err
	}
	return ServerActiveConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerActiveConnections) Name() string {
	return "signalr.server.active_connections"
}

// Unit returns the semantic convention unit of the instrument
func (ServerActiveConnections) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ServerActiveConnections) Description() string {
	return "Number of connections that are currently active on the server."
}

// Add adds incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ServerActiveConnections) Add(
    ctx context.Context,
    incr int64,
	attrs ...ServerActiveConnectionsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ServerActiveConnections) conv(in []ServerActiveConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverActiveConnectionsAttr()
	}
	return out
}

// ServerActiveConnectionsAttr is an optional attribute for the
// ServerActiveConnections instrument.
type ServerActiveConnectionsAttr interface {
    serverActiveConnectionsAttr() attribute.KeyValue
}

type serverActiveConnectionsAttr struct {
	kv attribute.KeyValue
}

func (a serverActiveConnectionsAttr) serverActiveConnectionsAttr() attribute.KeyValue {
    return a.kv
}

// ConnectionStatus returns an optional attribute for the
// "signalr.connection.status" semantic convention. It represents the signalR
// HTTP connection closure status.
func (ServerActiveConnections) ConnectionStatus(val ConnectionStatusAttr) ServerActiveConnectionsAttr {
	return serverActiveConnectionsAttr{kv: attribute.String("signalr.connection.status", string(val))}
}

// Transport returns an optional attribute for the "signalr.transport" semantic
// convention. It represents the [SignalR transport type].
//
// [SignalR transport type]: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md
func (ServerActiveConnections) Transport(val TransportAttr) ServerActiveConnectionsAttr {
	return serverActiveConnectionsAttr{kv: attribute.String("signalr.transport", string(val))}
}

// SignalrServerConnectionDuration is an instrument used to record metric values
// conforming to the "signalr.server.connection.duration" semantic conventions.
// It represents the duration of connections on the server.
type ServerConnectionDuration struct {
	inst metric.Float64Histogram
}

// NewServerConnectionDuration returns a new ServerConnectionDuration instrument.
func NewServerConnectionDuration(m metric.Meter) (ServerConnectionDuration, error) {
	i, err := m.Float64Histogram(
	    "signalr.server.connection.duration",
	    metric.WithDescription("The duration of connections on the server."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ServerConnectionDuration{}, err
	}
	return ServerConnectionDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerConnectionDuration) Name() string {
	return "signalr.server.connection.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ServerConnectionDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ServerConnectionDuration) Description() string {
	return "The duration of connections on the server."
}

// Record records incr to the existing count.
//
// All additional attrs passed are included in the recorded value.
func (m ServerConnectionDuration) Record(
    ctx context.Context,
    val float64,
	attrs ...ServerConnectionDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			m.conv(attrs)...,
		),
	)
}

func (m ServerConnectionDuration) conv(in []ServerConnectionDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverConnectionDurationAttr()
	}
	return out
}

// ServerConnectionDurationAttr is an optional attribute for the
// ServerConnectionDuration instrument.
type ServerConnectionDurationAttr interface {
    serverConnectionDurationAttr() attribute.KeyValue
}

type serverConnectionDurationAttr struct {
	kv attribute.KeyValue
}

func (a serverConnectionDurationAttr) serverConnectionDurationAttr() attribute.KeyValue {
    return a.kv
}

// ConnectionStatus returns an optional attribute for the
// "signalr.connection.status" semantic convention. It represents the signalR
// HTTP connection closure status.
func (ServerConnectionDuration) ConnectionStatus(val ConnectionStatusAttr) ServerConnectionDurationAttr {
	return serverConnectionDurationAttr{kv: attribute.String("signalr.connection.status", string(val))}
}

// Transport returns an optional attribute for the "signalr.transport" semantic
// convention. It represents the [SignalR transport type].
//
// [SignalR transport type]: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md
func (ServerConnectionDuration) Transport(val TransportAttr) ServerConnectionDurationAttr {
	return serverConnectionDurationAttr{kv: attribute.String("signalr.transport", string(val))}
}