// Code generated from semantic convention specification. DO NOT EDIT.

// Package httpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "signalr" namespace.
package signalrconv

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

// ServerActiveConnections is an instrument used to record metric values
// conforming to the "signalr.server.active_connections" semantic conventions. It
// represents the number of connections that are currently active on the server.
type ServerActiveConnections struct {
	metric.Int64UpDownCounter
}

// NewServerActiveConnections returns a new ServerActiveConnections instrument.
func NewServerActiveConnections(
	m metric.Meter,
	opt ...metric.Int64UpDownCounterOption,
) (ServerActiveConnections, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerActiveConnections{noop.Int64UpDownCounter{}}, nil
	}

	i, err := m.Int64UpDownCounter(
		"signalr.server.active_connections",
		append([]metric.Int64UpDownCounterOption{
			metric.WithDescription("Number of connections that are currently active on the server."),
			metric.WithUnit("{connection}"),
		}, opt...)...,
	)
	if err != nil {
	    return ServerActiveConnections{noop.Int64UpDownCounter{}}, err
	}
	return ServerActiveConnections{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerActiveConnections) Inst() metric.Int64UpDownCounter {
	return m.Int64UpDownCounter
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
//
// Meter name: `Microsoft.AspNetCore.Http.Connections`; Added in: ASP.NET Core
// 8.0
func (m ServerActiveConnections) Add(
	ctx context.Context,
	incr int64,
	attrs ...attribute.KeyValue,
) {
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

// AttrConnectionStatus returns an optional attribute for the
// "signalr.connection.status" semantic convention. It represents the signalR
// HTTP connection closure status.
func (ServerActiveConnections) AttrConnectionStatus(val ConnectionStatusAttr) attribute.KeyValue {
	return attribute.String("signalr.connection.status", string(val))
}

// AttrTransport returns an optional attribute for the "signalr.transport"
// semantic convention. It represents the [SignalR transport type].
//
// [SignalR transport type]: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md
func (ServerActiveConnections) AttrTransport(val TransportAttr) attribute.KeyValue {
	return attribute.String("signalr.transport", string(val))
}

// ServerConnectionDuration is an instrument used to record metric values
// conforming to the "signalr.server.connection.duration" semantic conventions.
// It represents the duration of connections on the server.
type ServerConnectionDuration struct {
	metric.Float64Histogram
}

// NewServerConnectionDuration returns a new ServerConnectionDuration instrument.
func NewServerConnectionDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ServerConnectionDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerConnectionDuration{noop.Float64Histogram{}}, nil
	}

	i, err := m.Float64Histogram(
		"signalr.server.connection.duration",
		append([]metric.Float64HistogramOption{
			metric.WithDescription("The duration of connections on the server."),
			metric.WithUnit("s"),
		}, opt...)...,
	)
	if err != nil {
	    return ServerConnectionDuration{noop.Float64Histogram{}}, err
	}
	return ServerConnectionDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerConnectionDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
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

// Record records val to the current distribution.
//
// All additional attrs passed are included in the recorded value.
//
// Meter name: `Microsoft.AspNetCore.Http.Connections`; Added in: ASP.NET Core
// 8.0
func (m ServerConnectionDuration) Record(
	ctx context.Context,
	val float64,
	attrs ...attribute.KeyValue,
) {
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

	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrConnectionStatus returns an optional attribute for the
// "signalr.connection.status" semantic convention. It represents the signalR
// HTTP connection closure status.
func (ServerConnectionDuration) AttrConnectionStatus(val ConnectionStatusAttr) attribute.KeyValue {
	return attribute.String("signalr.connection.status", string(val))
}

// AttrTransport returns an optional attribute for the "signalr.transport"
// semantic convention. It represents the [SignalR transport type].
//
// [SignalR transport type]: https://github.com/dotnet/aspnetcore/blob/main/src/SignalR/docs/specs/TransportProtocols.md
func (ServerConnectionDuration) AttrTransport(val TransportAttr) attribute.KeyValue {
	return attribute.String("signalr.transport", string(val))
}