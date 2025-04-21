// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/http"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

// ConnectionStateAttr is an attribute conforming to the http.connection.state
// semantic conventions. It represents the state of the HTTP connection in the
// HTTP connection pool.
type ConnectionStateAttr string

var (
	// ConnectionStateActive is the active state.
	ConnectionStateActive ConnectionStateAttr = "active"
	// ConnectionStateIdle is the idle state.
	ConnectionStateIdle ConnectionStateAttr = "idle"
)

// RequestMethodAttr is an attribute conforming to the http.request.method
// semantic conventions. It represents the HTTP request method.
type RequestMethodAttr string

var (
	// RequestMethodConnect is the CONNECT method.
	RequestMethodConnect RequestMethodAttr = "CONNECT"
	// RequestMethodDelete is the DELETE method.
	RequestMethodDelete RequestMethodAttr = "DELETE"
	// RequestMethodGet is the GET method.
	RequestMethodGet RequestMethodAttr = "GET"
	// RequestMethodHead is the HEAD method.
	RequestMethodHead RequestMethodAttr = "HEAD"
	// RequestMethodOptions is the OPTIONS method.
	RequestMethodOptions RequestMethodAttr = "OPTIONS"
	// RequestMethodPatch is the PATCH method.
	RequestMethodPatch RequestMethodAttr = "PATCH"
	// RequestMethodPost is the POST method.
	RequestMethodPost RequestMethodAttr = "POST"
	// RequestMethodPut is the PUT method.
	RequestMethodPut RequestMethodAttr = "PUT"
	// RequestMethodTrace is the TRACE method.
	RequestMethodTrace RequestMethodAttr = "TRACE"
	// RequestMethodOther is the any HTTP method that the instrumentation has no
	// prior knowledge of.
	RequestMethodOther RequestMethodAttr = "_OTHER"
)

// UserAgentSyntheticTypeAttr is an attribute conforming to the
// user_agent.synthetic.type semantic conventions. It represents the specifies
// the category of synthetic traffic, such as tests or bots.
type UserAgentSyntheticTypeAttr string

var (
	// UserAgentSyntheticTypeBot is the bot source.
	UserAgentSyntheticTypeBot UserAgentSyntheticTypeAttr = "bot"
	// UserAgentSyntheticTypeTest is the synthetic test source.
	UserAgentSyntheticTypeTest UserAgentSyntheticTypeAttr = "test"
)

// HTTPClientActiveRequests is an instrument used to record metric values
// conforming to the "http.client.active_requests" semantic conventions. It
// represents the number of active HTTP requests.
type ClientActiveRequests struct {
	inst metric.Int64UpDownCounter
}

// NewClientActiveRequests returns a new ClientActiveRequests instrument.
func NewClientActiveRequests(m metric.Meter) (ClientActiveRequests, error) {
	i, err := m.Int64UpDownCounter(
	    "http.client.active_requests",
	    metric.WithDescription("Number of active HTTP requests."),
	    metric.WithUnit("{request}"),
	)
	if err != nil {
	    return ClientActiveRequests{}, err
	}
	return ClientActiveRequests{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientActiveRequests) Name() string {
	return "http.client.active_requests"
}

// Unit returns the semantic convention unit of the instrument
func (ClientActiveRequests) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ClientActiveRequests) Description() string {
	return "Number of active HTTP requests."
}

// Add adds incr to the existing count.
//
// The serverAddress is the server domain name if available without reverse DNS
// lookup; otherwise, IP address or Unix domain socket name.
//
// The serverPort is the port identifier of the ["URI origin"] HTTP request is
// sent to.
//
// All additional attrs passed are included in the recorded value.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientActiveRequests) Add(
    ctx context.Context,
    incr int64,
	serverAddress string,
	serverPort int,
	attrs ...ClientActiveRequestsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("server.address", serverAddress),
				attribute.Int("server.port", serverPort),
			)...,
		),
	)
}

func (m ClientActiveRequests) conv(in []ClientActiveRequestsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientActiveRequestsAttr()
	}
	return out
}

// ClientActiveRequestsAttr is an optional attribute for the ClientActiveRequests
// instrument.
type ClientActiveRequestsAttr interface {
    clientActiveRequestsAttr() attribute.KeyValue
}

type clientActiveRequestsAttr struct {
	kv attribute.KeyValue
}

func (a clientActiveRequestsAttr) clientActiveRequestsAttr() attribute.KeyValue {
    return a.kv
}

// URLTemplate returns an optional attribute for the "url.template" semantic
// convention. It represents the low-cardinality template of an
// [absolute path reference].
//
// [absolute path reference]: https://www.rfc-editor.org/rfc/rfc3986#section-4.2
func (ClientActiveRequests) URLTemplate(val string) ClientActiveRequestsAttr {
	return clientActiveRequestsAttr{kv: attribute.String("url.template", val)}
}

// RequestMethod returns an optional attribute for the "http.request.method"
// semantic convention. It represents the HTTP request method.
func (ClientActiveRequests) RequestMethod(val RequestMethodAttr) ClientActiveRequestsAttr {
	return clientActiveRequestsAttr{kv: attribute.String("http.request.method", string(val))}
}

// URLScheme returns an optional attribute for the "url.scheme" semantic
// convention. It represents the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (ClientActiveRequests) URLScheme(val string) ClientActiveRequestsAttr {
	return clientActiveRequestsAttr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientConnectionDuration is an instrument used to record metric values
// conforming to the "http.client.connection.duration" semantic conventions. It
// represents the duration of the successfully established outbound HTTP
// connections.
type ClientConnectionDuration struct {
	inst metric.Float64Histogram
}

// NewClientConnectionDuration returns a new ClientConnectionDuration instrument.
func NewClientConnectionDuration(m metric.Meter) (ClientConnectionDuration, error) {
	i, err := m.Float64Histogram(
	    "http.client.connection.duration",
	    metric.WithDescription("The duration of the successfully established outbound HTTP connections."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientConnectionDuration{}, err
	}
	return ClientConnectionDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientConnectionDuration) Name() string {
	return "http.client.connection.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientConnectionDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientConnectionDuration) Description() string {
	return "The duration of the successfully established outbound HTTP connections."
}

// Record records incr to the existing count.
//
// The serverAddress is the server domain name if available without reverse DNS
// lookup; otherwise, IP address or Unix domain socket name.
//
// The serverPort is the port identifier of the ["URI origin"] HTTP request is
// sent to.
//
// All additional attrs passed are included in the recorded value.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientConnectionDuration) Record(
    ctx context.Context,
    val float64,
	serverAddress string,
	serverPort int,
	attrs ...ClientConnectionDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("server.address", serverAddress),
				attribute.Int("server.port", serverPort),
			)...,
		),
	)
}

func (m ClientConnectionDuration) conv(in []ClientConnectionDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientConnectionDurationAttr()
	}
	return out
}

// ClientConnectionDurationAttr is an optional attribute for the
// ClientConnectionDuration instrument.
type ClientConnectionDurationAttr interface {
    clientConnectionDurationAttr() attribute.KeyValue
}

type clientConnectionDurationAttr struct {
	kv attribute.KeyValue
}

func (a clientConnectionDurationAttr) clientConnectionDurationAttr() attribute.KeyValue {
    return a.kv
}

// NetworkPeerAddress returns an optional attribute for the
// "network.peer.address" semantic convention. It represents the peer address of
// the network connection - IP address or Unix domain socket name.
func (ClientConnectionDuration) NetworkPeerAddress(val string) ClientConnectionDurationAttr {
	return clientConnectionDurationAttr{kv: attribute.String("network.peer.address", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientConnectionDuration) NetworkProtocolVersion(val string) ClientConnectionDurationAttr {
	return clientConnectionDurationAttr{kv: attribute.String("network.protocol.version", val)}
}

// URLScheme returns an optional attribute for the "url.scheme" semantic
// convention. It represents the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (ClientConnectionDuration) URLScheme(val string) ClientConnectionDurationAttr {
	return clientConnectionDurationAttr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientOpenConnections is an instrument used to record metric values
// conforming to the "http.client.open_connections" semantic conventions. It
// represents the number of outbound HTTP connections that are currently active
// or idle on the client.
type ClientOpenConnections struct {
	inst metric.Int64UpDownCounter
}

// NewClientOpenConnections returns a new ClientOpenConnections instrument.
func NewClientOpenConnections(m metric.Meter) (ClientOpenConnections, error) {
	i, err := m.Int64UpDownCounter(
	    "http.client.open_connections",
	    metric.WithDescription("Number of outbound HTTP connections that are currently active or idle on the client."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return ClientOpenConnections{}, err
	}
	return ClientOpenConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientOpenConnections) Name() string {
	return "http.client.open_connections"
}

// Unit returns the semantic convention unit of the instrument
func (ClientOpenConnections) Unit() string {
	return "{connection}"
}

// Description returns the semantic convention description of the instrument
func (ClientOpenConnections) Description() string {
	return "Number of outbound HTTP connections that are currently active or idle on the client."
}

// Add adds incr to the existing count.
//
// The httpConnectionState is the state of the HTTP connection in the HTTP
// connection pool.
//
// The serverAddress is the server domain name if available without reverse DNS
// lookup; otherwise, IP address or Unix domain socket name.
//
// The serverPort is the port identifier of the ["URI origin"] HTTP request is
// sent to.
//
// All additional attrs passed are included in the recorded value.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientOpenConnections) Add(
    ctx context.Context,
    incr int64,
	connectionState ConnectionStateAttr,
	serverAddress string,
	serverPort int,
	attrs ...ClientOpenConnectionsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.connection.state", string(connectionState)),
				attribute.String("server.address", serverAddress),
				attribute.Int("server.port", serverPort),
			)...,
		),
	)
}

func (m ClientOpenConnections) conv(in []ClientOpenConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientOpenConnectionsAttr()
	}
	return out
}

// ClientOpenConnectionsAttr is an optional attribute for the
// ClientOpenConnections instrument.
type ClientOpenConnectionsAttr interface {
    clientOpenConnectionsAttr() attribute.KeyValue
}

type clientOpenConnectionsAttr struct {
	kv attribute.KeyValue
}

func (a clientOpenConnectionsAttr) clientOpenConnectionsAttr() attribute.KeyValue {
    return a.kv
}

// NetworkPeerAddress returns an optional attribute for the
// "network.peer.address" semantic convention. It represents the peer address of
// the network connection - IP address or Unix domain socket name.
func (ClientOpenConnections) NetworkPeerAddress(val string) ClientOpenConnectionsAttr {
	return clientOpenConnectionsAttr{kv: attribute.String("network.peer.address", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientOpenConnections) NetworkProtocolVersion(val string) ClientOpenConnectionsAttr {
	return clientOpenConnectionsAttr{kv: attribute.String("network.protocol.version", val)}
}

// URLScheme returns an optional attribute for the "url.scheme" semantic
// convention. It represents the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (ClientOpenConnections) URLScheme(val string) ClientOpenConnectionsAttr {
	return clientOpenConnectionsAttr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientRequestBodySize is an instrument used to record metric values
// conforming to the "http.client.request.body.size" semantic conventions. It
// represents the size of HTTP client request bodies.
type ClientRequestBodySize struct {
	inst metric.Int64Histogram
}

// NewClientRequestBodySize returns a new ClientRequestBodySize instrument.
func NewClientRequestBodySize(m metric.Meter) (ClientRequestBodySize, error) {
	i, err := m.Int64Histogram(
	    "http.client.request.body.size",
	    metric.WithDescription("Size of HTTP client request bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ClientRequestBodySize{}, err
	}
	return ClientRequestBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientRequestBodySize) Name() string {
	return "http.client.request.body.size"
}

// Unit returns the semantic convention unit of the instrument
func (ClientRequestBodySize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClientRequestBodySize) Description() string {
	return "Size of HTTP client request bodies."
}

// Record records incr to the existing count.
//
// The httpRequestMethod is the HTTP request method.
//
// The serverAddress is the host identifier of the ["URI origin"] HTTP request is
// sent to.
//
// The serverPort is the port identifier of the ["URI origin"] HTTP request is
// sent to.
//
// All additional attrs passed are included in the recorded value.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientRequestBodySize) Record(
    ctx context.Context,
    val int64,
	requestMethod RequestMethodAttr,
	serverAddress string,
	serverPort int,
	attrs ...ClientRequestBodySizeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.request.method", string(requestMethod)),
				attribute.String("server.address", serverAddress),
				attribute.Int("server.port", serverPort),
			)...,
		),
	)
}

func (m ClientRequestBodySize) conv(in []ClientRequestBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientRequestBodySizeAttr()
	}
	return out
}

// ClientRequestBodySizeAttr is an optional attribute for the
// ClientRequestBodySize instrument.
type ClientRequestBodySizeAttr interface {
    clientRequestBodySizeAttr() attribute.KeyValue
}

type clientRequestBodySizeAttr struct {
	kv attribute.KeyValue
}

func (a clientRequestBodySizeAttr) clientRequestBodySizeAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientRequestBodySize) ErrorType(val ErrorTypeAttr) ClientRequestBodySizeAttr {
	return clientRequestBodySizeAttr{kv: attribute.String("error.type", string(val))}
}

// ResponseStatusCode returns an optional attribute for the
// "http.response.status_code" semantic convention. It represents the
// [HTTP response status code].
//
// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
func (ClientRequestBodySize) ResponseStatusCode(val int) ClientRequestBodySizeAttr {
	return clientRequestBodySizeAttr{kv: attribute.Int("http.response.status_code", val)}
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientRequestBodySize) NetworkProtocolName(val string) ClientRequestBodySizeAttr {
	return clientRequestBodySizeAttr{kv: attribute.String("network.protocol.name", val)}
}

// URLTemplate returns an optional attribute for the "url.template" semantic
// convention. It represents the low-cardinality template of an
// [absolute path reference].
//
// [absolute path reference]: https://www.rfc-editor.org/rfc/rfc3986#section-4.2
func (ClientRequestBodySize) URLTemplate(val string) ClientRequestBodySizeAttr {
	return clientRequestBodySizeAttr{kv: attribute.String("url.template", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientRequestBodySize) NetworkProtocolVersion(val string) ClientRequestBodySizeAttr {
	return clientRequestBodySizeAttr{kv: attribute.String("network.protocol.version", val)}
}

// URLScheme returns an optional attribute for the "url.scheme" semantic
// convention. It represents the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (ClientRequestBodySize) URLScheme(val string) ClientRequestBodySizeAttr {
	return clientRequestBodySizeAttr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientRequestDuration is an instrument used to record metric values
// conforming to the "http.client.request.duration" semantic conventions. It
// represents the duration of HTTP client requests.
type ClientRequestDuration struct {
	inst metric.Float64Histogram
}

// NewClientRequestDuration returns a new ClientRequestDuration instrument.
func NewClientRequestDuration(m metric.Meter) (ClientRequestDuration, error) {
	i, err := m.Float64Histogram(
	    "http.client.request.duration",
	    metric.WithDescription("Duration of HTTP client requests."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ClientRequestDuration{}, err
	}
	return ClientRequestDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientRequestDuration) Name() string {
	return "http.client.request.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ClientRequestDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ClientRequestDuration) Description() string {
	return "Duration of HTTP client requests."
}

// Record records incr to the existing count.
//
// The httpRequestMethod is the HTTP request method.
//
// The serverAddress is the host identifier of the ["URI origin"] HTTP request is
// sent to.
//
// The serverPort is the port identifier of the ["URI origin"] HTTP request is
// sent to.
//
// All additional attrs passed are included in the recorded value.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientRequestDuration) Record(
    ctx context.Context,
    val float64,
	requestMethod RequestMethodAttr,
	serverAddress string,
	serverPort int,
	attrs ...ClientRequestDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.request.method", string(requestMethod)),
				attribute.String("server.address", serverAddress),
				attribute.Int("server.port", serverPort),
			)...,
		),
	)
}

func (m ClientRequestDuration) conv(in []ClientRequestDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientRequestDurationAttr()
	}
	return out
}

// ClientRequestDurationAttr is an optional attribute for the
// ClientRequestDuration instrument.
type ClientRequestDurationAttr interface {
    clientRequestDurationAttr() attribute.KeyValue
}

type clientRequestDurationAttr struct {
	kv attribute.KeyValue
}

func (a clientRequestDurationAttr) clientRequestDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientRequestDuration) ErrorType(val ErrorTypeAttr) ClientRequestDurationAttr {
	return clientRequestDurationAttr{kv: attribute.String("error.type", string(val))}
}

// ResponseStatusCode returns an optional attribute for the
// "http.response.status_code" semantic convention. It represents the
// [HTTP response status code].
//
// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
func (ClientRequestDuration) ResponseStatusCode(val int) ClientRequestDurationAttr {
	return clientRequestDurationAttr{kv: attribute.Int("http.response.status_code", val)}
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientRequestDuration) NetworkProtocolName(val string) ClientRequestDurationAttr {
	return clientRequestDurationAttr{kv: attribute.String("network.protocol.name", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientRequestDuration) NetworkProtocolVersion(val string) ClientRequestDurationAttr {
	return clientRequestDurationAttr{kv: attribute.String("network.protocol.version", val)}
}

// URLScheme returns an optional attribute for the "url.scheme" semantic
// convention. It represents the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (ClientRequestDuration) URLScheme(val string) ClientRequestDurationAttr {
	return clientRequestDurationAttr{kv: attribute.String("url.scheme", val)}
}

// URLTemplate returns an optional attribute for the "url.template" semantic
// convention. It represents the low-cardinality template of an
// [absolute path reference].
//
// [absolute path reference]: https://www.rfc-editor.org/rfc/rfc3986#section-4.2
func (ClientRequestDuration) URLTemplate(val string) ClientRequestDurationAttr {
	return clientRequestDurationAttr{kv: attribute.String("url.template", val)}
}

// HTTPClientResponseBodySize is an instrument used to record metric values
// conforming to the "http.client.response.body.size" semantic conventions. It
// represents the size of HTTP client response bodies.
type ClientResponseBodySize struct {
	inst metric.Int64Histogram
}

// NewClientResponseBodySize returns a new ClientResponseBodySize instrument.
func NewClientResponseBodySize(m metric.Meter) (ClientResponseBodySize, error) {
	i, err := m.Int64Histogram(
	    "http.client.response.body.size",
	    metric.WithDescription("Size of HTTP client response bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ClientResponseBodySize{}, err
	}
	return ClientResponseBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ClientResponseBodySize) Name() string {
	return "http.client.response.body.size"
}

// Unit returns the semantic convention unit of the instrument
func (ClientResponseBodySize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ClientResponseBodySize) Description() string {
	return "Size of HTTP client response bodies."
}

// Record records incr to the existing count.
//
// The httpRequestMethod is the HTTP request method.
//
// The serverAddress is the host identifier of the ["URI origin"] HTTP request is
// sent to.
//
// The serverPort is the port identifier of the ["URI origin"] HTTP request is
// sent to.
//
// All additional attrs passed are included in the recorded value.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientResponseBodySize) Record(
    ctx context.Context,
    val int64,
	requestMethod RequestMethodAttr,
	serverAddress string,
	serverPort int,
	attrs ...ClientResponseBodySizeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.request.method", string(requestMethod)),
				attribute.String("server.address", serverAddress),
				attribute.Int("server.port", serverPort),
			)...,
		),
	)
}

func (m ClientResponseBodySize) conv(in []ClientResponseBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.clientResponseBodySizeAttr()
	}
	return out
}

// ClientResponseBodySizeAttr is an optional attribute for the
// ClientResponseBodySize instrument.
type ClientResponseBodySizeAttr interface {
    clientResponseBodySizeAttr() attribute.KeyValue
}

type clientResponseBodySizeAttr struct {
	kv attribute.KeyValue
}

func (a clientResponseBodySizeAttr) clientResponseBodySizeAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ClientResponseBodySize) ErrorType(val ErrorTypeAttr) ClientResponseBodySizeAttr {
	return clientResponseBodySizeAttr{kv: attribute.String("error.type", string(val))}
}

// ResponseStatusCode returns an optional attribute for the
// "http.response.status_code" semantic convention. It represents the
// [HTTP response status code].
//
// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
func (ClientResponseBodySize) ResponseStatusCode(val int) ClientResponseBodySizeAttr {
	return clientResponseBodySizeAttr{kv: attribute.Int("http.response.status_code", val)}
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ClientResponseBodySize) NetworkProtocolName(val string) ClientResponseBodySizeAttr {
	return clientResponseBodySizeAttr{kv: attribute.String("network.protocol.name", val)}
}

// URLTemplate returns an optional attribute for the "url.template" semantic
// convention. It represents the low-cardinality template of an
// [absolute path reference].
//
// [absolute path reference]: https://www.rfc-editor.org/rfc/rfc3986#section-4.2
func (ClientResponseBodySize) URLTemplate(val string) ClientResponseBodySizeAttr {
	return clientResponseBodySizeAttr{kv: attribute.String("url.template", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ClientResponseBodySize) NetworkProtocolVersion(val string) ClientResponseBodySizeAttr {
	return clientResponseBodySizeAttr{kv: attribute.String("network.protocol.version", val)}
}

// URLScheme returns an optional attribute for the "url.scheme" semantic
// convention. It represents the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (ClientResponseBodySize) URLScheme(val string) ClientResponseBodySizeAttr {
	return clientResponseBodySizeAttr{kv: attribute.String("url.scheme", val)}
}

// HTTPServerActiveRequests is an instrument used to record metric values
// conforming to the "http.server.active_requests" semantic conventions. It
// represents the number of active HTTP server requests.
type ServerActiveRequests struct {
	inst metric.Int64UpDownCounter
}

// NewServerActiveRequests returns a new ServerActiveRequests instrument.
func NewServerActiveRequests(m metric.Meter) (ServerActiveRequests, error) {
	i, err := m.Int64UpDownCounter(
	    "http.server.active_requests",
	    metric.WithDescription("Number of active HTTP server requests."),
	    metric.WithUnit("{request}"),
	)
	if err != nil {
	    return ServerActiveRequests{}, err
	}
	return ServerActiveRequests{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerActiveRequests) Name() string {
	return "http.server.active_requests"
}

// Unit returns the semantic convention unit of the instrument
func (ServerActiveRequests) Unit() string {
	return "{request}"
}

// Description returns the semantic convention description of the instrument
func (ServerActiveRequests) Description() string {
	return "Number of active HTTP server requests."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod is the HTTP request method.
//
// The urlScheme is the the [URI scheme] component identifying the used protocol.
//
// All additional attrs passed are included in the recorded value.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerActiveRequests) Add(
    ctx context.Context,
    incr int64,
	requestMethod RequestMethodAttr,
	urlScheme string,
	attrs ...ServerActiveRequestsAttr,
) {
	m.inst.Add(
		ctx,
		incr,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.request.method", string(requestMethod)),
				attribute.String("url.scheme", urlScheme),
			)...,
		),
	)
}

func (m ServerActiveRequests) conv(in []ServerActiveRequestsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverActiveRequestsAttr()
	}
	return out
}

// ServerActiveRequestsAttr is an optional attribute for the ServerActiveRequests
// instrument.
type ServerActiveRequestsAttr interface {
    serverActiveRequestsAttr() attribute.KeyValue
}

type serverActiveRequestsAttr struct {
	kv attribute.KeyValue
}

func (a serverActiveRequestsAttr) serverActiveRequestsAttr() attribute.KeyValue {
    return a.kv
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the local HTTP server that received the
// request.
func (ServerActiveRequests) ServerAddress(val string) ServerActiveRequestsAttr {
	return serverActiveRequestsAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the port of the local HTTP server that received the
// request.
func (ServerActiveRequests) ServerPort(val int) ServerActiveRequestsAttr {
	return serverActiveRequestsAttr{kv: attribute.Int("server.port", val)}
}

// HTTPServerRequestBodySize is an instrument used to record metric values
// conforming to the "http.server.request.body.size" semantic conventions. It
// represents the size of HTTP server request bodies.
type ServerRequestBodySize struct {
	inst metric.Int64Histogram
}

// NewServerRequestBodySize returns a new ServerRequestBodySize instrument.
func NewServerRequestBodySize(m metric.Meter) (ServerRequestBodySize, error) {
	i, err := m.Int64Histogram(
	    "http.server.request.body.size",
	    metric.WithDescription("Size of HTTP server request bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ServerRequestBodySize{}, err
	}
	return ServerRequestBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerRequestBodySize) Name() string {
	return "http.server.request.body.size"
}

// Unit returns the semantic convention unit of the instrument
func (ServerRequestBodySize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ServerRequestBodySize) Description() string {
	return "Size of HTTP server request bodies."
}

// Record records incr to the existing count.
//
// The httpRequestMethod is the HTTP request method.
//
// The urlScheme is the the [URI scheme] component identifying the used protocol.
//
// All additional attrs passed are included in the recorded value.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerRequestBodySize) Record(
    ctx context.Context,
    val int64,
	requestMethod RequestMethodAttr,
	urlScheme string,
	attrs ...ServerRequestBodySizeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.request.method", string(requestMethod)),
				attribute.String("url.scheme", urlScheme),
			)...,
		),
	)
}

func (m ServerRequestBodySize) conv(in []ServerRequestBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverRequestBodySizeAttr()
	}
	return out
}

// ServerRequestBodySizeAttr is an optional attribute for the
// ServerRequestBodySize instrument.
type ServerRequestBodySizeAttr interface {
    serverRequestBodySizeAttr() attribute.KeyValue
}

type serverRequestBodySizeAttr struct {
	kv attribute.KeyValue
}

func (a serverRequestBodySizeAttr) serverRequestBodySizeAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerRequestBodySize) ErrorType(val ErrorTypeAttr) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.String("error.type", string(val))}
}

// ResponseStatusCode returns an optional attribute for the
// "http.response.status_code" semantic convention. It represents the
// [HTTP response status code].
//
// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
func (ServerRequestBodySize) ResponseStatusCode(val int) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.Int("http.response.status_code", val)}
}

// Route returns an optional attribute for the "http.route" semantic convention.
// It represents the matched route, that is, the path template in the format used
// by the respective server framework.
func (ServerRequestBodySize) Route(val string) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.String("http.route", val)}
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerRequestBodySize) NetworkProtocolName(val string) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.String("network.protocol.name", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerRequestBodySize) NetworkProtocolVersion(val string) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.String("network.protocol.version", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the local HTTP server that received the
// request.
func (ServerRequestBodySize) ServerAddress(val string) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the port of the local HTTP server that received the
// request.
func (ServerRequestBodySize) ServerPort(val int) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.Int("server.port", val)}
}

// UserAgentSyntheticType returns an optional attribute for the
// "user_agent.synthetic.type" semantic convention. It represents the specifies
// the category of synthetic traffic, such as tests or bots.
func (ServerRequestBodySize) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) ServerRequestBodySizeAttr {
	return serverRequestBodySizeAttr{kv: attribute.String("user_agent.synthetic.type", string(val))}
}

// HTTPServerRequestDuration is an instrument used to record metric values
// conforming to the "http.server.request.duration" semantic conventions. It
// represents the duration of HTTP server requests.
type ServerRequestDuration struct {
	inst metric.Float64Histogram
}

// NewServerRequestDuration returns a new ServerRequestDuration instrument.
func NewServerRequestDuration(m metric.Meter) (ServerRequestDuration, error) {
	i, err := m.Float64Histogram(
	    "http.server.request.duration",
	    metric.WithDescription("Duration of HTTP server requests."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return ServerRequestDuration{}, err
	}
	return ServerRequestDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerRequestDuration) Name() string {
	return "http.server.request.duration"
}

// Unit returns the semantic convention unit of the instrument
func (ServerRequestDuration) Unit() string {
	return "s"
}

// Description returns the semantic convention description of the instrument
func (ServerRequestDuration) Description() string {
	return "Duration of HTTP server requests."
}

// Record records incr to the existing count.
//
// The httpRequestMethod is the HTTP request method.
//
// The urlScheme is the the [URI scheme] component identifying the used protocol.
//
// All additional attrs passed are included in the recorded value.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerRequestDuration) Record(
    ctx context.Context,
    val float64,
	requestMethod RequestMethodAttr,
	urlScheme string,
	attrs ...ServerRequestDurationAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.request.method", string(requestMethod)),
				attribute.String("url.scheme", urlScheme),
			)...,
		),
	)
}

func (m ServerRequestDuration) conv(in []ServerRequestDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverRequestDurationAttr()
	}
	return out
}

// ServerRequestDurationAttr is an optional attribute for the
// ServerRequestDuration instrument.
type ServerRequestDurationAttr interface {
    serverRequestDurationAttr() attribute.KeyValue
}

type serverRequestDurationAttr struct {
	kv attribute.KeyValue
}

func (a serverRequestDurationAttr) serverRequestDurationAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerRequestDuration) ErrorType(val ErrorTypeAttr) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("error.type", string(val))}
}

// ResponseStatusCode returns an optional attribute for the
// "http.response.status_code" semantic convention. It represents the
// [HTTP response status code].
//
// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
func (ServerRequestDuration) ResponseStatusCode(val int) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.Int("http.response.status_code", val)}
}

// Route returns an optional attribute for the "http.route" semantic convention.
// It represents the matched route, that is, the path template in the format used
// by the respective server framework.
func (ServerRequestDuration) Route(val string) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("http.route", val)}
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerRequestDuration) NetworkProtocolName(val string) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("network.protocol.name", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerRequestDuration) NetworkProtocolVersion(val string) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("network.protocol.version", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the local HTTP server that received the
// request.
func (ServerRequestDuration) ServerAddress(val string) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the port of the local HTTP server that received the
// request.
func (ServerRequestDuration) ServerPort(val int) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.Int("server.port", val)}
}

// UserAgentSyntheticType returns an optional attribute for the
// "user_agent.synthetic.type" semantic convention. It represents the specifies
// the category of synthetic traffic, such as tests or bots.
func (ServerRequestDuration) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) ServerRequestDurationAttr {
	return serverRequestDurationAttr{kv: attribute.String("user_agent.synthetic.type", string(val))}
}

// HTTPServerResponseBodySize is an instrument used to record metric values
// conforming to the "http.server.response.body.size" semantic conventions. It
// represents the size of HTTP server response bodies.
type ServerResponseBodySize struct {
	inst metric.Int64Histogram
}

// NewServerResponseBodySize returns a new ServerResponseBodySize instrument.
func NewServerResponseBodySize(m metric.Meter) (ServerResponseBodySize, error) {
	i, err := m.Int64Histogram(
	    "http.server.response.body.size",
	    metric.WithDescription("Size of HTTP server response bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return ServerResponseBodySize{}, err
	}
	return ServerResponseBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (ServerResponseBodySize) Name() string {
	return "http.server.response.body.size"
}

// Unit returns the semantic convention unit of the instrument
func (ServerResponseBodySize) Unit() string {
	return "By"
}

// Description returns the semantic convention description of the instrument
func (ServerResponseBodySize) Description() string {
	return "Size of HTTP server response bodies."
}

// Record records incr to the existing count.
//
// The httpRequestMethod is the HTTP request method.
//
// The urlScheme is the the [URI scheme] component identifying the used protocol.
//
// All additional attrs passed are included in the recorded value.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerResponseBodySize) Record(
    ctx context.Context,
    val int64,
	requestMethod RequestMethodAttr,
	urlScheme string,
	attrs ...ServerResponseBodySizeAttr,
) {
	m.inst.Record(
		ctx,
		val,
		metric.WithAttributes(
			append(
				m.conv(attrs),
				attribute.String("http.request.method", string(requestMethod)),
				attribute.String("url.scheme", urlScheme),
			)...,
		),
	)
}

func (m ServerResponseBodySize) conv(in []ServerResponseBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.serverResponseBodySizeAttr()
	}
	return out
}

// ServerResponseBodySizeAttr is an optional attribute for the
// ServerResponseBodySize instrument.
type ServerResponseBodySizeAttr interface {
    serverResponseBodySizeAttr() attribute.KeyValue
}

type serverResponseBodySizeAttr struct {
	kv attribute.KeyValue
}

func (a serverResponseBodySizeAttr) serverResponseBodySizeAttr() attribute.KeyValue {
    return a.kv
}

// ErrorType returns an optional attribute for the "error.type" semantic
// convention. It represents the describes a class of error the operation ended
// with.
func (ServerResponseBodySize) ErrorType(val ErrorTypeAttr) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.String("error.type", string(val))}
}

// ResponseStatusCode returns an optional attribute for the
// "http.response.status_code" semantic convention. It represents the
// [HTTP response status code].
//
// [HTTP response status code]: https://tools.ietf.org/html/rfc7231#section-6
func (ServerResponseBodySize) ResponseStatusCode(val int) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.Int("http.response.status_code", val)}
}

// Route returns an optional attribute for the "http.route" semantic convention.
// It represents the matched route, that is, the path template in the format used
// by the respective server framework.
func (ServerResponseBodySize) Route(val string) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.String("http.route", val)}
}

// NetworkProtocolName returns an optional attribute for the
// "network.protocol.name" semantic convention. It represents the
// [OSI application layer] or non-OSI equivalent.
//
// [OSI application layer]: https://wikipedia.org/wiki/Application_layer
func (ServerResponseBodySize) NetworkProtocolName(val string) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.String("network.protocol.name", val)}
}

// NetworkProtocolVersion returns an optional attribute for the
// "network.protocol.version" semantic convention. It represents the actual
// version of the protocol used for network communication.
func (ServerResponseBodySize) NetworkProtocolVersion(val string) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.String("network.protocol.version", val)}
}

// ServerAddress returns an optional attribute for the "server.address" semantic
// convention. It represents the name of the local HTTP server that received the
// request.
func (ServerResponseBodySize) ServerAddress(val string) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.String("server.address", val)}
}

// ServerPort returns an optional attribute for the "server.port" semantic
// convention. It represents the port of the local HTTP server that received the
// request.
func (ServerResponseBodySize) ServerPort(val int) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.Int("server.port", val)}
}

// UserAgentSyntheticType returns an optional attribute for the
// "user_agent.synthetic.type" semantic convention. It represents the specifies
// the category of synthetic traffic, such as tests or bots.
func (ServerResponseBodySize) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) ServerResponseBodySizeAttr {
	return serverResponseBodySizeAttr{kv: attribute.String("user_agent.synthetic.type", string(val))}
}