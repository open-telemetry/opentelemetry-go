// Code generated from semantic convention specification. DO NOT EDIT.

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package mcpconv provides types and functionality for OpenTelemetry semantic
// conventions in the "mcp" namespace.
package mcpconv

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

// GenAIOperationNameAttr is an attribute conforming to the gen_ai.operation.name
// semantic conventions. It represents the name of the GenAI operation being
// performed.
type GenAIOperationNameAttr string

var (
	// GenAIOperationNameChat is the chat completion operation such as
	// [OpenAI Chat API].
	//
	// [OpenAI Chat API]: https://platform.openai.com/docs/api-reference/chat
	GenAIOperationNameChat GenAIOperationNameAttr = "chat"
	// GenAIOperationNameGenerateContent is the multimodal content generation
	// operation such as [Gemini Generate Content].
	//
	// [Gemini Generate Content]: https://ai.google.dev/api/generate-content
	GenAIOperationNameGenerateContent GenAIOperationNameAttr = "generate_content"
	// GenAIOperationNameTextCompletion is the text completions operation such as
	// [OpenAI Completions API (Legacy)].
	//
	// [OpenAI Completions API (Legacy)]: https://platform.openai.com/docs/api-reference/completions
	GenAIOperationNameTextCompletion GenAIOperationNameAttr = "text_completion"
	// GenAIOperationNameEmbeddings is the embeddings operation such as
	// [OpenAI Create embeddings API].
	//
	// [OpenAI Create embeddings API]: https://platform.openai.com/docs/api-reference/embeddings/create
	GenAIOperationNameEmbeddings GenAIOperationNameAttr = "embeddings"
	// GenAIOperationNameRetrieval is the retrieval operation such as
	// [OpenAI Search Vector Store API].
	//
	// [OpenAI Search Vector Store API]: https://platform.openai.com/docs/api-reference/vector-stores/search
	GenAIOperationNameRetrieval GenAIOperationNameAttr = "retrieval"
	// GenAIOperationNameCreateAgent is the create GenAI agent.
	GenAIOperationNameCreateAgent GenAIOperationNameAttr = "create_agent"
	// GenAIOperationNameInvokeAgent is the invoke GenAI agent.
	GenAIOperationNameInvokeAgent GenAIOperationNameAttr = "invoke_agent"
	// GenAIOperationNameExecuteTool is the execute a tool.
	GenAIOperationNameExecuteTool GenAIOperationNameAttr = "execute_tool"
)

// MethodNameAttr is an attribute conforming to the mcp.method.name semantic
// conventions. It represents the name of the request or notification method.
type MethodNameAttr string

var (
	// MethodNameNotificationsCancelled is the notification cancelling a
	// previously-issued request.
	MethodNameNotificationsCancelled MethodNameAttr = "notifications/cancelled"
	// MethodNameInitialize is the request to initialize the MCP client.
	MethodNameInitialize MethodNameAttr = "initialize"
	// MethodNameNotificationsInitialized is the notification indicating that the
	// MCP client has been initialized.
	MethodNameNotificationsInitialized MethodNameAttr = "notifications/initialized"
	// MethodNameNotificationsProgress is the notification indicating the progress
	// for a long-running operation.
	MethodNameNotificationsProgress MethodNameAttr = "notifications/progress"
	// MethodNamePing is the request to check that the other party is still alive.
	MethodNamePing MethodNameAttr = "ping"
	// MethodNameResourcesList is the request to list resources available on server.
	MethodNameResourcesList MethodNameAttr = "resources/list"
	// MethodNameResourcesTemplatesList is the request to list resource templates
	// available on server.
	MethodNameResourcesTemplatesList MethodNameAttr = "resources/templates/list"
	// MethodNameResourcesRead is the request to read a resource.
	MethodNameResourcesRead MethodNameAttr = "resources/read"
	// MethodNameNotificationsResourcesListChanged is the notification indicating
	// that the list of resources has changed.
	MethodNameNotificationsResourcesListChanged MethodNameAttr = "notifications/resources/list_changed"
	// MethodNameResourcesSubscribe is the request to subscribe to a resource.
	MethodNameResourcesSubscribe MethodNameAttr = "resources/subscribe"
	// MethodNameResourcesUnsubscribe is the request to unsubscribe from resource
	// updates.
	MethodNameResourcesUnsubscribe MethodNameAttr = "resources/unsubscribe"
	// MethodNameNotificationsResourcesUpdated is the notification indicating that a
	// resource has been updated.
	MethodNameNotificationsResourcesUpdated MethodNameAttr = "notifications/resources/updated"
	// MethodNamePromptsList is the request to list prompts available on server.
	MethodNamePromptsList MethodNameAttr = "prompts/list"
	// MethodNamePromptsGet is the request to get a prompt.
	MethodNamePromptsGet MethodNameAttr = "prompts/get"
	// MethodNameNotificationsPromptsListChanged is the notification indicating that
	// the list of prompts has changed.
	MethodNameNotificationsPromptsListChanged MethodNameAttr = "notifications/prompts/list_changed"
	// MethodNameToolsList is the request to list tools available on server.
	MethodNameToolsList MethodNameAttr = "tools/list"
	// MethodNameToolsCall is the request to call a tool.
	MethodNameToolsCall MethodNameAttr = "tools/call"
	// MethodNameNotificationsToolsListChanged is the notification indicating that
	// the list of tools has changed.
	MethodNameNotificationsToolsListChanged MethodNameAttr = "notifications/tools/list_changed"
	// MethodNameLoggingSetLevel is the request to set the logging level.
	MethodNameLoggingSetLevel MethodNameAttr = "logging/setLevel"
	// MethodNameNotificationsMessage is the notification indicating that a message
	// has been received.
	MethodNameNotificationsMessage MethodNameAttr = "notifications/message"
	// MethodNameSamplingCreateMessage is the request to create a sampling message.
	MethodNameSamplingCreateMessage MethodNameAttr = "sampling/createMessage"
	// MethodNameCompletionComplete is the request to complete a prompt.
	MethodNameCompletionComplete MethodNameAttr = "completion/complete"
	// MethodNameRootsList is the request to list roots available on server.
	MethodNameRootsList MethodNameAttr = "roots/list"
	// MethodNameNotificationsRootsListChanged is the notification indicating that
	// the list of roots has changed.
	MethodNameNotificationsRootsListChanged MethodNameAttr = "notifications/roots/list_changed"
	// MethodNameElicitationCreate is the request from the server to elicit
	// additional information from the user via the client.
	MethodNameElicitationCreate MethodNameAttr = "elicitation/create"
)

// NetworkTransportAttr is an attribute conforming to the network.transport
// semantic conventions. It represents the transport protocol used for the MCP
// session.
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

// ClientOperationDuration is an instrument used to record metric values
// conforming to the "mcp.client.operation.duration" semantic conventions. It
// represents the duration of the MCP request or notification as observed on the
// sender from the time it was sent until the response or ack is received.
type ClientOperationDuration struct {
	metric.Float64Histogram
}

var newClientOperationDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("The duration of the MCP request or notification as observed on the sender from the time it was sent until the response or ack is received."),
	metric.WithUnit("s"),
}

// NewClientOperationDuration returns a new ClientOperationDuration instrument.
func NewClientOperationDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientOperationDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientOperationDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientOperationDurationOpts
	} else {
		opt = append(opt, newClientOperationDurationOpts...)
	}

	i, err := m.Float64Histogram(
		"mcp.client.operation.duration",
		opt...,
	)
	if err != nil {
		return ClientOperationDuration{noop.Float64Histogram{}}, err
	}
	return ClientOperationDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientOperationDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientOperationDuration) Name() string {
	return "mcp.client.operation.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOperationDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientOperationDuration) Description() string {
	return "The duration of the MCP request or notification as observed on the sender from the time it was sent until the response or ack is received."
}

// Record records val to the current distribution for attrs.
//
// The methodName is the the name of the request or notification method.
//
// All additional attrs passed are included in the recorded value.
func (m ClientOperationDuration) Record(
	ctx context.Context,
	val float64,
	methodName MethodNameAttr,
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
				attribute.String("mcp.method.name", string(methodName)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientOperationDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
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
func (ClientOperationDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrGenAIPromptName returns an optional attribute for the "gen_ai.prompt.name"
// semantic convention. It represents the name of the prompt or prompt template
// provided in the request or response.
func (ClientOperationDuration) AttrGenAIPromptName(val string) attribute.KeyValue {
	return attribute.String("gen_ai.prompt.name", val)
}

// AttrGenAIToolName returns an optional attribute for the "gen_ai.tool.name"
// semantic convention. It represents the name of the tool utilized by the agent.
func (ClientOperationDuration) AttrGenAIToolName(val string) attribute.KeyValue {
	return attribute.String("gen_ai.tool.name", val)
}

// AttrRPCResponseStatusCode returns an optional attribute for the
// "rpc.response.status_code" semantic convention. It represents the error code
// from the JSON-RPC response.
func (ClientOperationDuration) AttrRPCResponseStatusCode(val string) attribute.KeyValue {
	return attribute.String("rpc.response.status_code", val)
}

// AttrGenAIOperationName returns an optional attribute for the
// "gen_ai.operation.name" semantic convention. It represents the name of the
// GenAI operation being performed.
func (ClientOperationDuration) AttrGenAIOperationName(val GenAIOperationNameAttr) attribute.KeyValue {
	return attribute.String("gen_ai.operation.name", string(val))
}

// AttrJSONRPCProtocolVersion returns an optional attribute for the
// "jsonrpc.protocol.version" semantic convention. It represents the protocol
// version, as specified in the `jsonrpc` property of the request and its
// corresponding response.
func (ClientOperationDuration) AttrJSONRPCProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("jsonrpc.protocol.version", val)
}

// AttrProtocolVersion returns an optional attribute for the
// "mcp.protocol.version" semantic convention. It represents the [version] of the
// Model Context Protocol used.
//
// [version]: https://modelcontextprotocol.io/specification/versioning
func (ClientOperationDuration) AttrProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("mcp.protocol.version", val)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientOperationDuration) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientOperationDuration) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the transport protocol used for the MCP
// session.
func (ClientOperationDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientOperationDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientOperationDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrResourceURI returns an optional attribute for the "mcp.resource.uri"
// semantic convention. It represents the value of the resource uri.
func (ClientOperationDuration) AttrResourceURI(val string) attribute.KeyValue {
	return attribute.String("mcp.resource.uri", val)
}

// ClientSessionDuration is an instrument used to record metric values conforming
// to the "mcp.client.session.duration" semantic conventions. It represents the
// duration of the MCP session as observed on the MCP client.
type ClientSessionDuration struct {
	metric.Float64Histogram
}

var newClientSessionDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("The duration of the MCP session as observed on the MCP client."),
	metric.WithUnit("s"),
}

// NewClientSessionDuration returns a new ClientSessionDuration instrument.
func NewClientSessionDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ClientSessionDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ClientSessionDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientSessionDurationOpts
	} else {
		opt = append(opt, newClientSessionDurationOpts...)
	}

	i, err := m.Float64Histogram(
		"mcp.client.session.duration",
		opt...,
	)
	if err != nil {
		return ClientSessionDuration{noop.Float64Histogram{}}, err
	}
	return ClientSessionDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientSessionDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ClientSessionDuration) Name() string {
	return "mcp.client.session.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientSessionDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientSessionDuration) Description() string {
	return "The duration of the MCP session as observed on the MCP client."
}

// Record records val to the current distribution for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m ClientSessionDuration) Record(
	ctx context.Context,
	val float64,
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
			attrs...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ClientSessionDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
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
func (ClientSessionDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrJSONRPCProtocolVersion returns an optional attribute for the
// "jsonrpc.protocol.version" semantic convention. It represents the protocol
// version, as specified in the `jsonrpc` property of the request and its
// corresponding response.
func (ClientSessionDuration) AttrJSONRPCProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("jsonrpc.protocol.version", val)
}

// AttrProtocolVersion returns an optional attribute for the
// "mcp.protocol.version" semantic convention. It represents the [version] of the
// Model Context Protocol used.
//
// [version]: https://modelcontextprotocol.io/specification/versioning
func (ClientSessionDuration) AttrProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("mcp.protocol.version", val)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientSessionDuration) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientSessionDuration) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the transport protocol used for the MCP
// session.
func (ClientSessionDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention. It represents the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
func (ClientSessionDuration) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the server port number.
func (ClientSessionDuration) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// ServerOperationDuration is an instrument used to record metric values
// conforming to the "mcp.server.operation.duration" semantic conventions. It
// represents the MCP request or notification duration as observed on the
// receiver from the time it was received until the result or ack is sent.
type ServerOperationDuration struct {
	metric.Float64Histogram
}

var newServerOperationDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("MCP request or notification duration as observed on the receiver from the time it was received until the result or ack is sent."),
	metric.WithUnit("s"),
}

// NewServerOperationDuration returns a new ServerOperationDuration instrument.
func NewServerOperationDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ServerOperationDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerOperationDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerOperationDurationOpts
	} else {
		opt = append(opt, newServerOperationDurationOpts...)
	}

	i, err := m.Float64Histogram(
		"mcp.server.operation.duration",
		opt...,
	)
	if err != nil {
		return ServerOperationDuration{noop.Float64Histogram{}}, err
	}
	return ServerOperationDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerOperationDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ServerOperationDuration) Name() string {
	return "mcp.server.operation.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ServerOperationDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ServerOperationDuration) Description() string {
	return "MCP request or notification duration as observed on the receiver from the time it was received until the result or ack is sent."
}

// Record records val to the current distribution for attrs.
//
// The methodName is the the name of the request or notification method.
//
// All additional attrs passed are included in the recorded value.
func (m ServerOperationDuration) Record(
	ctx context.Context,
	val float64,
	methodName MethodNameAttr,
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
				attribute.String("mcp.method.name", string(methodName)),
			)...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ServerOperationDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
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
func (ServerOperationDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrGenAIPromptName returns an optional attribute for the "gen_ai.prompt.name"
// semantic convention. It represents the name of the prompt or prompt template
// provided in the request or response.
func (ServerOperationDuration) AttrGenAIPromptName(val string) attribute.KeyValue {
	return attribute.String("gen_ai.prompt.name", val)
}

// AttrGenAIToolName returns an optional attribute for the "gen_ai.tool.name"
// semantic convention. It represents the name of the tool utilized by the agent.
func (ServerOperationDuration) AttrGenAIToolName(val string) attribute.KeyValue {
	return attribute.String("gen_ai.tool.name", val)
}

// AttrRPCResponseStatusCode returns an optional attribute for the
// "rpc.response.status_code" semantic convention. It represents the error code
// from the JSON-RPC response.
func (ServerOperationDuration) AttrRPCResponseStatusCode(val string) attribute.KeyValue {
	return attribute.String("rpc.response.status_code", val)
}

// AttrGenAIOperationName returns an optional attribute for the
// "gen_ai.operation.name" semantic convention. It represents the name of the
// GenAI operation being performed.
func (ServerOperationDuration) AttrGenAIOperationName(val GenAIOperationNameAttr) attribute.KeyValue {
	return attribute.String("gen_ai.operation.name", string(val))
}

// AttrJSONRPCProtocolVersion returns an optional attribute for the
// "jsonrpc.protocol.version" semantic convention. It represents the protocol
// version, as specified in the `jsonrpc` property of the request and its
// corresponding response.
func (ServerOperationDuration) AttrJSONRPCProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("jsonrpc.protocol.version", val)
}

// AttrProtocolVersion returns an optional attribute for the
// "mcp.protocol.version" semantic convention. It represents the [version] of the
// Model Context Protocol used.
//
// [version]: https://modelcontextprotocol.io/specification/versioning
func (ServerOperationDuration) AttrProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("mcp.protocol.version", val)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerOperationDuration) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerOperationDuration) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the transport protocol used for the MCP
// session.
func (ServerOperationDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}

// AttrResourceURI returns an optional attribute for the "mcp.resource.uri"
// semantic convention. It represents the value of the resource uri.
func (ServerOperationDuration) AttrResourceURI(val string) attribute.KeyValue {
	return attribute.String("mcp.resource.uri", val)
}

// ServerSessionDuration is an instrument used to record metric values conforming
// to the "mcp.server.session.duration" semantic conventions. It represents the
// duration of the MCP session as observed on the MCP server.
type ServerSessionDuration struct {
	metric.Float64Histogram
}

var newServerSessionDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("The duration of the MCP session as observed on the MCP server."),
	metric.WithUnit("s"),
}

// NewServerSessionDuration returns a new ServerSessionDuration instrument.
func NewServerSessionDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (ServerSessionDuration, error) {
	// Check if the meter is nil.
	if m == nil {
		return ServerSessionDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newServerSessionDurationOpts
	} else {
		opt = append(opt, newServerSessionDurationOpts...)
	}

	i, err := m.Float64Histogram(
		"mcp.server.session.duration",
		opt...,
	)
	if err != nil {
		return ServerSessionDuration{noop.Float64Histogram{}}, err
	}
	return ServerSessionDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ServerSessionDuration) Inst() metric.Float64Histogram {
	return m.Float64Histogram
}

// Name returns the semantic convention name of the instrument.
func (ServerSessionDuration) Name() string {
	return "mcp.server.session.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ServerSessionDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ServerSessionDuration) Description() string {
	return "The duration of the MCP session as observed on the MCP server."
}

// Record records val to the current distribution for attrs.
//
// All additional attrs passed are included in the recorded value.
func (m ServerSessionDuration) Record(
	ctx context.Context,
	val float64,
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
			attrs...,
		),
	)

	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m ServerSessionDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
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
func (ServerSessionDuration) AttrErrorType(val ErrorTypeAttr) attribute.KeyValue {
	return attribute.String("error.type", string(val))
}

// AttrJSONRPCProtocolVersion returns an optional attribute for the
// "jsonrpc.protocol.version" semantic convention. It represents the protocol
// version, as specified in the `jsonrpc` property of the request and its
// corresponding response.
func (ServerSessionDuration) AttrJSONRPCProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("jsonrpc.protocol.version", val)
}

// AttrProtocolVersion returns an optional attribute for the
// "mcp.protocol.version" semantic convention. It represents the [version] of the
// Model Context Protocol used.
//
// [version]: https://modelcontextprotocol.io/specification/versioning
func (ServerSessionDuration) AttrProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("mcp.protocol.version", val)
}

// AttrNetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerSessionDuration) AttrNetworkProtocolName(val string) attribute.KeyValue {
	return attribute.String("network.protocol.name", val)
}

// AttrNetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerSessionDuration) AttrNetworkProtocolVersion(val string) attribute.KeyValue {
	return attribute.String("network.protocol.version", val)
}

// AttrNetworkTransport returns an optional attribute for the "network.transport"
// semantic convention. It represents the transport protocol used for the MCP
// session.
func (ServerSessionDuration) AttrNetworkTransport(val NetworkTransportAttr) attribute.KeyValue {
	return attribute.String("network.transport", string(val))
}
