// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/rpc"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// RPCClientDuration is an instrument used to record metric values conforming to
// the "rpc.client.duration" semantic conventions. It represents the measures the
// duration of outbound RPC.
type ClientDuration struct {
	inst metric.Int64Histogram
}

// NewClientDuration returns a new ClientDuration instrument.
func NewClientDuration(m metric.Meter) (ClientDuration, error) {
	i, err := m.Int64Histogram(
	    "rpc.client.duration",
	    metric.WithDescription("Measures the duration of outbound RPC."),
	    metric.WithUnit("ms"),
	)
	if err != nil {
	    return ClientDuration{}, err
	}
	return ClientDuration{i}, nil
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

func (m ClientDuration) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCClientRequestSize is an instrument used to record metric values conforming
// to the "rpc.client.request.size" semantic conventions. It represents the
// measures the size of RPC request messages (uncompressed).
type ClientRequestSize struct {
	inst metric.Int64Histogram
}

// NewClientRequestSize returns a new ClientRequestSize instrument.
func NewClientRequestSize(m metric.Meter) (ClientRequestSize, error) {
	i, err := m.Int64Histogram(
	    "rpc.client.request.size",
	    metric.WithDescription("Measures the size of RPC request messages (uncompressed)."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ClientRequestSize{}, err
	}
	return ClientRequestSize{i}, nil
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

func (m ClientRequestSize) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCClientRequestsPerRPC is an instrument used to record metric values
// conforming to the "rpc.client.requests_per_rpc" semantic conventions. It
// represents the measures the number of messages received per RPC.
type ClientRequestsPerRPC struct {
	inst metric.Int64Histogram
}

// NewClientRequestsPerRPC returns a new ClientRequestsPerRPC instrument.
func NewClientRequestsPerRPC(m metric.Meter) (ClientRequestsPerRPC, error) {
	i, err := m.Int64Histogram(
	    "rpc.client.requests_per_rpc",
	    metric.WithDescription("Measures the number of messages received per RPC."),
	    metric.WithUnit("{count}"),
	)
	if err != nil {
	    return ClientRequestsPerRPC{}, err
	}
	return ClientRequestsPerRPC{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientRequestsPerRPC) Name() string {
	return "rpc.client.requests_per_rpc"
}

// Unit returns the semantic convention unit of the instrument
func (ClientRequestsPerRPC) Unit() string {
	return "{count}"
}

// Description returns the semantic convention description of the instrument
func (ClientRequestsPerRPC) Description() string {
	return "Measures the number of messages received per RPC."
}

func (m ClientRequestsPerRPC) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCClientResponseSize is an instrument used to record metric values conforming
// to the "rpc.client.response.size" semantic conventions. It represents the
// measures the size of RPC response messages (uncompressed).
type ClientResponseSize struct {
	inst metric.Int64Histogram
}

// NewClientResponseSize returns a new ClientResponseSize instrument.
func NewClientResponseSize(m metric.Meter) (ClientResponseSize, error) {
	i, err := m.Int64Histogram(
	    "rpc.client.response.size",
	    metric.WithDescription("Measures the size of RPC response messages (uncompressed)."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ClientResponseSize{}, err
	}
	return ClientResponseSize{i}, nil
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

func (m ClientResponseSize) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCClientResponsesPerRPC is an instrument used to record metric values
// conforming to the "rpc.client.responses_per_rpc" semantic conventions. It
// represents the measures the number of messages sent per RPC.
type ClientResponsesPerRPC struct {
	inst metric.Int64Histogram
}

// NewClientResponsesPerRPC returns a new ClientResponsesPerRPC instrument.
func NewClientResponsesPerRPC(m metric.Meter) (ClientResponsesPerRPC, error) {
	i, err := m.Int64Histogram(
	    "rpc.client.responses_per_rpc",
	    metric.WithDescription("Measures the number of messages sent per RPC."),
	    metric.WithUnit("{count}"),
	)
	if err != nil {
	    return ClientResponsesPerRPC{}, err
	}
	return ClientResponsesPerRPC{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientResponsesPerRPC) Name() string {
	return "rpc.client.responses_per_rpc"
}

// Unit returns the semantic convention unit of the instrument
func (ClientResponsesPerRPC) Unit() string {
	return "{count}"
}

// Description returns the semantic convention description of the instrument
func (ClientResponsesPerRPC) Description() string {
	return "Measures the number of messages sent per RPC."
}

func (m ClientResponsesPerRPC) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCServerDuration is an instrument used to record metric values conforming to
// the "rpc.server.duration" semantic conventions. It represents the measures the
// duration of inbound RPC.
type ServerDuration struct {
	inst metric.Int64Histogram
}

// NewServerDuration returns a new ServerDuration instrument.
func NewServerDuration(m metric.Meter) (ServerDuration, error) {
	i, err := m.Int64Histogram(
	    "rpc.server.duration",
	    metric.WithDescription("Measures the duration of inbound RPC."),
	    metric.WithUnit("ms"),
	)
	if err != nil {
	    return ServerDuration{}, err
	}
	return ServerDuration{i}, nil
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

func (m ServerDuration) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCServerRequestSize is an instrument used to record metric values conforming
// to the "rpc.server.request.size" semantic conventions. It represents the
// measures the size of RPC request messages (uncompressed).
type ServerRequestSize struct {
	inst metric.Int64Histogram
}

// NewServerRequestSize returns a new ServerRequestSize instrument.
func NewServerRequestSize(m metric.Meter) (ServerRequestSize, error) {
	i, err := m.Int64Histogram(
	    "rpc.server.request.size",
	    metric.WithDescription("Measures the size of RPC request messages (uncompressed)."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ServerRequestSize{}, err
	}
	return ServerRequestSize{i}, nil
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

func (m ServerRequestSize) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCServerRequestsPerRPC is an instrument used to record metric values
// conforming to the "rpc.server.requests_per_rpc" semantic conventions. It
// represents the measures the number of messages received per RPC.
type ServerRequestsPerRPC struct {
	inst metric.Int64Histogram
}

// NewServerRequestsPerRPC returns a new ServerRequestsPerRPC instrument.
func NewServerRequestsPerRPC(m metric.Meter) (ServerRequestsPerRPC, error) {
	i, err := m.Int64Histogram(
	    "rpc.server.requests_per_rpc",
	    metric.WithDescription("Measures the number of messages received per RPC."),
	    metric.WithUnit("{count}"),
	)
	if err != nil {
	    return ServerRequestsPerRPC{}, err
	}
	return ServerRequestsPerRPC{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerRequestsPerRPC) Name() string {
	return "rpc.server.requests_per_rpc"
}

// Unit returns the semantic convention unit of the instrument
func (ServerRequestsPerRPC) Unit() string {
	return "{count}"
}

// Description returns the semantic convention description of the instrument
func (ServerRequestsPerRPC) Description() string {
	return "Measures the number of messages received per RPC."
}

func (m ServerRequestsPerRPC) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCServerResponseSize is an instrument used to record metric values conforming
// to the "rpc.server.response.size" semantic conventions. It represents the
// measures the size of RPC response messages (uncompressed).
type ServerResponseSize struct {
	inst metric.Int64Histogram
}

// NewServerResponseSize returns a new ServerResponseSize instrument.
func NewServerResponseSize(m metric.Meter) (ServerResponseSize, error) {
	i, err := m.Int64Histogram(
	    "rpc.server.response.size",
	    metric.WithDescription("Measures the size of RPC response messages (uncompressed)."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ServerResponseSize{}, err
	}
	return ServerResponseSize{i}, nil
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

func (m ServerResponseSize) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}

// RPCServerResponsesPerRPC is an instrument used to record metric values
// conforming to the "rpc.server.responses_per_rpc" semantic conventions. It
// represents the measures the number of messages sent per RPC.
type ServerResponsesPerRPC struct {
	inst metric.Int64Histogram
}

// NewServerResponsesPerRPC returns a new ServerResponsesPerRPC instrument.
func NewServerResponsesPerRPC(m metric.Meter) (ServerResponsesPerRPC, error) {
	i, err := m.Int64Histogram(
	    "rpc.server.responses_per_rpc",
	    metric.WithDescription("Measures the number of messages sent per RPC."),
	    metric.WithUnit("{count}"),
	)
	if err != nil {
	    return ServerResponsesPerRPC{}, err
	}
	return ServerResponsesPerRPC{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerResponsesPerRPC) Name() string {
	return "rpc.server.responses_per_rpc"
}

// Unit returns the semantic convention unit of the instrument
func (ServerResponsesPerRPC) Unit() string {
	return "{count}"
}

// Description returns the semantic convention description of the instrument
func (ServerResponsesPerRPC) Description() string {
	return "Measures the number of messages sent per RPC."
}

func (m ServerResponsesPerRPC) Record(ctx context.Context, val int64) {
    m.inst.Record(ctx, val)
}