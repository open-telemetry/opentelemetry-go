// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0/http"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type ErrorTypeAttr string

var (
	ErrorTypeOther ErrorTypeAttr = "_OTHER"
)

type ConnectionStateAttr string

var (
	ConnectionStateActive ConnectionStateAttr = "active"
	ConnectionStateIdle ConnectionStateAttr = "idle"
)

type RequestMethodAttr string

var (
	RequestMethodConnect RequestMethodAttr = "CONNECT"
	RequestMethodDelete RequestMethodAttr = "DELETE"
	RequestMethodGet RequestMethodAttr = "GET"
	RequestMethodHead RequestMethodAttr = "HEAD"
	RequestMethodOptions RequestMethodAttr = "OPTIONS"
	RequestMethodPatch RequestMethodAttr = "PATCH"
	RequestMethodPost RequestMethodAttr = "POST"
	RequestMethodPut RequestMethodAttr = "PUT"
	RequestMethodTrace RequestMethodAttr = "TRACE"
	RequestMethodOther RequestMethodAttr = "_OTHER"
)

type UserAgentSyntheticTypeAttr string

var (
	UserAgentSyntheticTypeBot UserAgentSyntheticTypeAttr = "bot"
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

// Name returns the semantic convention unit of the instrument
func (ClientActiveRequests) Unit() string {
	return "{request}"
}

// Name returns the semantic convention description of the instrument
func (ClientActiveRequests) Description() string {
	return "Number of active HTTP requests."
}

// Add adds incr to the existing count.
//
// The serverAddress parameter is the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
//
// The serverPort parameter is the port identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientActiveRequests) Add(
    ctx context.Context,
    incr int64,
    serverAddress string,
    serverPort int,
    attrs ...ClientActiveRequestsAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientActiveRequestsAttrToAttrs(attrs),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientActiveRequestsAttrToAttrs(in []ClientActiveRequestsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientActiveRequestsAttr()
	}
	return out
}

type ClientActiveRequestsAttr interface {
    httpClientActiveRequestsAttr() attribute.KeyValue
}

func (a attr) httpClientActiveRequestsAttr() attribute.KeyValue {
    return a.kv
}

func (ClientActiveRequests) URLTemplate(val string) ClientActiveRequestsAttr {
	return attr{kv: attribute.String("url.template", val)}
}

func (ClientActiveRequests) RequestMethod(val RequestMethodAttr) ClientActiveRequestsAttr {
	return attr{kv: attribute.String("http.request.method", string(val))}
}

func (ClientActiveRequests) URLScheme(val string) ClientActiveRequestsAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientConnectionDuration is an instrument used to record metric values
// conforming to the "http.client.connection.duration" semantic conventions. It
// represents the duration of the successfully established outbound HTTP
// connections.
type ClientConnectionDuration struct {
	inst metric.Int64Histogram
}

// NewClientConnectionDuration returns a new ClientConnectionDuration instrument.
func NewClientConnectionDuration(m metric.Meter) (ClientConnectionDuration, error) {
	i, err := m.Int64Histogram(
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

// Name returns the semantic convention unit of the instrument
func (ClientConnectionDuration) Unit() string {
	return "s"
}

// Name returns the semantic convention description of the instrument
func (ClientConnectionDuration) Description() string {
	return "The duration of the successfully established outbound HTTP connections."
}

// Add adds incr to the existing count.
//
// The serverAddress parameter is the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
//
// The serverPort parameter is the port identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientConnectionDuration) Add(
    ctx context.Context,
    incr int64,
    serverAddress string,
    serverPort int,
    attrs ...ClientConnectionDurationAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientConnectionDurationAttrToAttrs(attrs),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientConnectionDurationAttrToAttrs(in []ClientConnectionDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientConnectionDurationAttr()
	}
	return out
}

type ClientConnectionDurationAttr interface {
    httpClientConnectionDurationAttr() attribute.KeyValue
}

func (a attr) httpClientConnectionDurationAttr() attribute.KeyValue {
    return a.kv
}

func (ClientConnectionDuration) NetworkPeerAddress(val string) ClientConnectionDurationAttr {
	return attr{kv: attribute.String("network.peer.address", val)}
}

func (ClientConnectionDuration) NetworkProtocolVersion(val string) ClientConnectionDurationAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ClientConnectionDuration) URLScheme(val string) ClientConnectionDurationAttr {
	return attr{kv: attribute.String("url.scheme", val)}
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

// Name returns the semantic convention unit of the instrument
func (ClientOpenConnections) Unit() string {
	return "{connection}"
}

// Name returns the semantic convention description of the instrument
func (ClientOpenConnections) Description() string {
	return "Number of outbound HTTP connections that are currently active or idle on the client."
}

// Add adds incr to the existing count.
//
// The httpConnectionState parameter is the state of the HTTP connection in the
// HTTP connection pool.
//
// The serverAddress parameter is the server domain name if available without
// reverse DNS lookup; otherwise, IP address or Unix domain socket name.
//
// The serverPort parameter is the port identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientOpenConnections) Add(
    ctx context.Context,
    incr int64,
    httpConnectionState ConnectionStateAttr,
    serverAddress string,
    serverPort int,
    attrs ...ClientOpenConnectionsAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientOpenConnectionsAttrToAttrs(attrs),
			attribute.String("http.connection.state", string(httpConnectionState)),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientOpenConnectionsAttrToAttrs(in []ClientOpenConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientOpenConnectionsAttr()
	}
	return out
}

type ClientOpenConnectionsAttr interface {
    httpClientOpenConnectionsAttr() attribute.KeyValue
}

func (a attr) httpClientOpenConnectionsAttr() attribute.KeyValue {
    return a.kv
}

func (ClientOpenConnections) NetworkPeerAddress(val string) ClientOpenConnectionsAttr {
	return attr{kv: attribute.String("network.peer.address", val)}
}

func (ClientOpenConnections) NetworkProtocolVersion(val string) ClientOpenConnectionsAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ClientOpenConnections) URLScheme(val string) ClientOpenConnectionsAttr {
	return attr{kv: attribute.String("url.scheme", val)}
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

// Name returns the semantic convention unit of the instrument
func (ClientRequestBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (ClientRequestBodySize) Description() string {
	return "Size of HTTP client request bodies."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod parameter is the hTTP request method.
//
// The serverAddress parameter is the host identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
//
// The serverPort parameter is the port identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientRequestBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod RequestMethodAttr,
    serverAddress string,
    serverPort int,
    attrs ...ClientRequestBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientRequestBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", string(httpRequestMethod)),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientRequestBodySizeAttrToAttrs(in []ClientRequestBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientRequestBodySizeAttr()
	}
	return out
}

type ClientRequestBodySizeAttr interface {
    httpClientRequestBodySizeAttr() attribute.KeyValue
}

func (a attr) httpClientRequestBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (ClientRequestBodySize) ErrorType(val ErrorTypeAttr) ClientRequestBodySizeAttr {
	return attr{kv: attribute.String("error.type", string(val))}
}

func (ClientRequestBodySize) ResponseStatusCode(val int) ClientRequestBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (ClientRequestBodySize) NetworkProtocolName(val string) ClientRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (ClientRequestBodySize) URLTemplate(val string) ClientRequestBodySizeAttr {
	return attr{kv: attribute.String("url.template", val)}
}

func (ClientRequestBodySize) NetworkProtocolVersion(val string) ClientRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ClientRequestBodySize) URLScheme(val string) ClientRequestBodySizeAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientRequestDuration is an instrument used to record metric values
// conforming to the "http.client.request.duration" semantic conventions. It
// represents the duration of HTTP client requests.
type ClientRequestDuration struct {
	inst metric.Int64Histogram
}

// NewClientRequestDuration returns a new ClientRequestDuration instrument.
func NewClientRequestDuration(m metric.Meter) (ClientRequestDuration, error) {
	i, err := m.Int64Histogram(
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

// Name returns the semantic convention unit of the instrument
func (ClientRequestDuration) Unit() string {
	return "s"
}

// Name returns the semantic convention description of the instrument
func (ClientRequestDuration) Description() string {
	return "Duration of HTTP client requests."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod parameter is the hTTP request method.
//
// The serverAddress parameter is the host identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
//
// The serverPort parameter is the port identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientRequestDuration) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod RequestMethodAttr,
    serverAddress string,
    serverPort int,
    attrs ...ClientRequestDurationAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientRequestDurationAttrToAttrs(attrs),
			attribute.String("http.request.method", string(httpRequestMethod)),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientRequestDurationAttrToAttrs(in []ClientRequestDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientRequestDurationAttr()
	}
	return out
}

type ClientRequestDurationAttr interface {
    httpClientRequestDurationAttr() attribute.KeyValue
}

func (a attr) httpClientRequestDurationAttr() attribute.KeyValue {
    return a.kv
}

func (ClientRequestDuration) ErrorType(val ErrorTypeAttr) ClientRequestDurationAttr {
	return attr{kv: attribute.String("error.type", string(val))}
}

func (ClientRequestDuration) ResponseStatusCode(val int) ClientRequestDurationAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (ClientRequestDuration) NetworkProtocolName(val string) ClientRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (ClientRequestDuration) NetworkProtocolVersion(val string) ClientRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ClientRequestDuration) URLScheme(val string) ClientRequestDurationAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

func (ClientRequestDuration) URLTemplate(val string) ClientRequestDurationAttr {
	return attr{kv: attribute.String("url.template", val)}
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

// Name returns the semantic convention unit of the instrument
func (ClientResponseBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (ClientResponseBodySize) Description() string {
	return "Size of HTTP client response bodies."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod parameter is the hTTP request method.
//
// The serverAddress parameter is the host identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
//
// The serverPort parameter is the port identifier of the ["URI origin"] HTTP
// request is sent to.
//
// ["URI origin"]: https://www.rfc-editor.org/rfc/rfc9110.html#name-uri-origin
func (m ClientResponseBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod RequestMethodAttr,
    serverAddress string,
    serverPort int,
    attrs ...ClientResponseBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientResponseBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", string(httpRequestMethod)),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientResponseBodySizeAttrToAttrs(in []ClientResponseBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientResponseBodySizeAttr()
	}
	return out
}

type ClientResponseBodySizeAttr interface {
    httpClientResponseBodySizeAttr() attribute.KeyValue
}

func (a attr) httpClientResponseBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (ClientResponseBodySize) ErrorType(val ErrorTypeAttr) ClientResponseBodySizeAttr {
	return attr{kv: attribute.String("error.type", string(val))}
}

func (ClientResponseBodySize) ResponseStatusCode(val int) ClientResponseBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (ClientResponseBodySize) NetworkProtocolName(val string) ClientResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (ClientResponseBodySize) URLTemplate(val string) ClientResponseBodySizeAttr {
	return attr{kv: attribute.String("url.template", val)}
}

func (ClientResponseBodySize) NetworkProtocolVersion(val string) ClientResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ClientResponseBodySize) URLScheme(val string) ClientResponseBodySizeAttr {
	return attr{kv: attribute.String("url.scheme", val)}
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

// Name returns the semantic convention unit of the instrument
func (ServerActiveRequests) Unit() string {
	return "{request}"
}

// Name returns the semantic convention description of the instrument
func (ServerActiveRequests) Description() string {
	return "Number of active HTTP server requests."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod parameter is the hTTP request method.
//
// The urlScheme parameter is the the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerActiveRequests) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod RequestMethodAttr,
    urlScheme string,
    attrs ...ServerActiveRequestsAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerActiveRequestsAttrToAttrs(attrs),
			attribute.String("http.request.method", string(httpRequestMethod)),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerActiveRequestsAttrToAttrs(in []ServerActiveRequestsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerActiveRequestsAttr()
	}
	return out
}

type ServerActiveRequestsAttr interface {
    httpServerActiveRequestsAttr() attribute.KeyValue
}

func (a attr) httpServerActiveRequestsAttr() attribute.KeyValue {
    return a.kv
}

func (ServerActiveRequests) ServerAddress(val string) ServerActiveRequestsAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (ServerActiveRequests) ServerPort(val int) ServerActiveRequestsAttr {
	return attr{kv: attribute.Int("server.port", val)}
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

// Name returns the semantic convention unit of the instrument
func (ServerRequestBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (ServerRequestBodySize) Description() string {
	return "Size of HTTP server request bodies."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod parameter is the hTTP request method.
//
// The urlScheme parameter is the the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerRequestBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod RequestMethodAttr,
    urlScheme string,
    attrs ...ServerRequestBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerRequestBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", string(httpRequestMethod)),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerRequestBodySizeAttrToAttrs(in []ServerRequestBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerRequestBodySizeAttr()
	}
	return out
}

type ServerRequestBodySizeAttr interface {
    httpServerRequestBodySizeAttr() attribute.KeyValue
}

func (a attr) httpServerRequestBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (ServerRequestBodySize) ErrorType(val ErrorTypeAttr) ServerRequestBodySizeAttr {
	return attr{kv: attribute.String("error.type", string(val))}
}

func (ServerRequestBodySize) ResponseStatusCode(val int) ServerRequestBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (ServerRequestBodySize) Route(val string) ServerRequestBodySizeAttr {
	return attr{kv: attribute.String("http.route", val)}
}

func (ServerRequestBodySize) NetworkProtocolName(val string) ServerRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (ServerRequestBodySize) NetworkProtocolVersion(val string) ServerRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ServerRequestBodySize) ServerAddress(val string) ServerRequestBodySizeAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (ServerRequestBodySize) ServerPort(val int) ServerRequestBodySizeAttr {
	return attr{kv: attribute.Int("server.port", val)}
}

func (ServerRequestBodySize) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) ServerRequestBodySizeAttr {
	return attr{kv: attribute.String("user_agent.synthetic.type", string(val))}
}

// HTTPServerRequestDuration is an instrument used to record metric values
// conforming to the "http.server.request.duration" semantic conventions. It
// represents the duration of HTTP server requests.
type ServerRequestDuration struct {
	inst metric.Int64Histogram
}

// NewServerRequestDuration returns a new ServerRequestDuration instrument.
func NewServerRequestDuration(m metric.Meter) (ServerRequestDuration, error) {
	i, err := m.Int64Histogram(
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

// Name returns the semantic convention unit of the instrument
func (ServerRequestDuration) Unit() string {
	return "s"
}

// Name returns the semantic convention description of the instrument
func (ServerRequestDuration) Description() string {
	return "Duration of HTTP server requests."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod parameter is the hTTP request method.
//
// The urlScheme parameter is the the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerRequestDuration) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod RequestMethodAttr,
    urlScheme string,
    attrs ...ServerRequestDurationAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerRequestDurationAttrToAttrs(attrs),
			attribute.String("http.request.method", string(httpRequestMethod)),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerRequestDurationAttrToAttrs(in []ServerRequestDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerRequestDurationAttr()
	}
	return out
}

type ServerRequestDurationAttr interface {
    httpServerRequestDurationAttr() attribute.KeyValue
}

func (a attr) httpServerRequestDurationAttr() attribute.KeyValue {
    return a.kv
}

func (ServerRequestDuration) ErrorType(val ErrorTypeAttr) ServerRequestDurationAttr {
	return attr{kv: attribute.String("error.type", string(val))}
}

func (ServerRequestDuration) ResponseStatusCode(val int) ServerRequestDurationAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (ServerRequestDuration) Route(val string) ServerRequestDurationAttr {
	return attr{kv: attribute.String("http.route", val)}
}

func (ServerRequestDuration) NetworkProtocolName(val string) ServerRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (ServerRequestDuration) NetworkProtocolVersion(val string) ServerRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ServerRequestDuration) ServerAddress(val string) ServerRequestDurationAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (ServerRequestDuration) ServerPort(val int) ServerRequestDurationAttr {
	return attr{kv: attribute.Int("server.port", val)}
}

func (ServerRequestDuration) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) ServerRequestDurationAttr {
	return attr{kv: attribute.String("user_agent.synthetic.type", string(val))}
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

// Name returns the semantic convention unit of the instrument
func (ServerResponseBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (ServerResponseBodySize) Description() string {
	return "Size of HTTP server response bodies."
}

// Add adds incr to the existing count.
//
// The httpRequestMethod parameter is the hTTP request method.
//
// The urlScheme parameter is the the [URI scheme] component identifying the used
// protocol.
//
// [URI scheme]: https://www.rfc-editor.org/rfc/rfc3986#section-3.1
func (m ServerResponseBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod RequestMethodAttr,
    urlScheme string,
    attrs ...ServerResponseBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerResponseBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", string(httpRequestMethod)),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerResponseBodySizeAttrToAttrs(in []ServerResponseBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerResponseBodySizeAttr()
	}
	return out
}

type ServerResponseBodySizeAttr interface {
    httpServerResponseBodySizeAttr() attribute.KeyValue
}

func (a attr) httpServerResponseBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (ServerResponseBodySize) ErrorType(val ErrorTypeAttr) ServerResponseBodySizeAttr {
	return attr{kv: attribute.String("error.type", string(val))}
}

func (ServerResponseBodySize) ResponseStatusCode(val int) ServerResponseBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (ServerResponseBodySize) Route(val string) ServerResponseBodySizeAttr {
	return attr{kv: attribute.String("http.route", val)}
}

func (ServerResponseBodySize) NetworkProtocolName(val string) ServerResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (ServerResponseBodySize) NetworkProtocolVersion(val string) ServerResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (ServerResponseBodySize) ServerAddress(val string) ServerResponseBodySizeAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (ServerResponseBodySize) ServerPort(val int) ServerResponseBodySizeAttr {
	return attr{kv: attribute.Int("server.port", val)}
}

func (ServerResponseBodySize) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) ServerResponseBodySizeAttr {
	return attr{kv: attribute.String("user_agent.synthetic.type", string(val))}
}

// This is used as a helper for all optional attributes.
type attr struct {
	kv attribute.KeyValue
}