// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated from semantic convention specification. DO NOT EDIT.

package semconv // import "go.opentelemetry.io/otel/semconv/v1.31.0"

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type ErrorTypeAttr string

var (
	ErrorTypeOther = ErrorTypeAttr("_OTHER")
)

type HTTPConnectionStateAttr string

var (
	HTTPConnectionStateActive = HTTPConnectionStateAttr("active")
	HTTPConnectionStateIdle = HTTPConnectionStateAttr("idle")
)

type HTTPRequestMethodAttr string

var (
	HTTPRequestMethodConnect = HTTPRequestMethodAttr("CONNECT")
	HTTPRequestMethodDelete = HTTPRequestMethodAttr("DELETE")
	HTTPRequestMethodGet = HTTPRequestMethodAttr("GET")
	HTTPRequestMethodHead = HTTPRequestMethodAttr("HEAD")
	HTTPRequestMethodOptions = HTTPRequestMethodAttr("OPTIONS")
	HTTPRequestMethodPatch = HTTPRequestMethodAttr("PATCH")
	HTTPRequestMethodPost = HTTPRequestMethodAttr("POST")
	HTTPRequestMethodPut = HTTPRequestMethodAttr("PUT")
	HTTPRequestMethodTrace = HTTPRequestMethodAttr("TRACE")
	HTTPRequestMethodOther = HTTPRequestMethodAttr("_OTHER")
)

type UserAgentSyntheticTypeAttr string

var (
	UserAgentSyntheticTypeBot = UserAgentSyntheticTypeAttr("bot")
	UserAgentSyntheticTypeTest = UserAgentSyntheticTypeAttr("test")
)

// HTTPClientActiveRequests is an instrument used to record metric values
// conforming to the "http.client.active_requests" semantic conventions. It
// represents the number of active HTTP requests.
type HTTPClientActiveRequests struct {
	inst metric.Int64UpDownCounter
}

// NewHTTPClientActiveRequests returns a new HTTPClientActiveRequests instrument.
func NewHTTPClientActiveRequests(m metric.Meter) (HTTPClientActiveRequests, error) {
	i, err := m.Int64UpDownCounter(
	    "http.client.active_requests",
	    metric.WithDescription("Number of active HTTP requests."),
	    metric.WithUnit("{request}"),
	)
	if err != nil {
	    return HTTPClientActiveRequests{}, err
	}
	return HTTPClientActiveRequests{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPClientActiveRequests) Name() string {
	return "http.client.active_requests"
}

// Name returns the semantic convention unit of the instrument
func (HTTPClientActiveRequests) Unit() string {
	return "{request}"
}

// Name returns the semantic convention description of the instrument
func (HTTPClientActiveRequests) Description() string {
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
func (m HTTPClientActiveRequests) Add(
    ctx context.Context,
    incr int64,
    serverAddress string,
    serverPort int,
    attrs ...HTTPClientActiveRequestsAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientActiveRequestsAttrToAttrs(attrs),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientActiveRequestsAttrToAttrs(in []HTTPClientActiveRequestsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientActiveRequestsAttr()
	}
	return out
}

type HTTPClientActiveRequestsAttr interface {
    httpClientActiveRequestsAttr() attribute.KeyValue
}

func (a attr) httpClientActiveRequestsAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPClientActiveRequests) URLTemplate(val string) HTTPClientActiveRequestsAttr {
	return attr{kv: attribute.String("url.template", val)}
}

func (HTTPClientActiveRequests) HTTPRequestMethod(val HTTPRequestMethodAttr) HTTPClientActiveRequestsAttr {
	return attr{kv: attribute.String("http.request.method", val)}
}

func (HTTPClientActiveRequests) URLScheme(val string) HTTPClientActiveRequestsAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientConnectionDuration is an instrument used to record metric values
// conforming to the "http.client.connection.duration" semantic conventions. It
// represents the duration of the successfully established outbound HTTP
// connections.
type HTTPClientConnectionDuration struct {
	inst metric.histogram
}

// NewHTTPClientConnectionDuration returns a new HTTPClientConnectionDuration
// instrument.
func NewHTTPClientConnectionDuration(m metric.Meter) (HTTPClientConnectionDuration, error) {
	i, err := m.histogram(
	    "http.client.connection.duration",
	    metric.WithDescription("The duration of the successfully established outbound HTTP connections."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return HTTPClientConnectionDuration{}, err
	}
	return HTTPClientConnectionDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPClientConnectionDuration) Name() string {
	return "http.client.connection.duration"
}

// Name returns the semantic convention unit of the instrument
func (HTTPClientConnectionDuration) Unit() string {
	return "s"
}

// Name returns the semantic convention description of the instrument
func (HTTPClientConnectionDuration) Description() string {
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
func (m HTTPClientConnectionDuration) Add(
    ctx context.Context,
    incr int64,
    serverAddress string,
    serverPort int,
    attrs ...HTTPClientConnectionDurationAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientConnectionDurationAttrToAttrs(attrs),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientConnectionDurationAttrToAttrs(in []HTTPClientConnectionDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientConnectionDurationAttr()
	}
	return out
}

type HTTPClientConnectionDurationAttr interface {
    httpClientConnectionDurationAttr() attribute.KeyValue
}

func (a attr) httpClientConnectionDurationAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPClientConnectionDuration) NetworkPeerAddress(val string) HTTPClientConnectionDurationAttr {
	return attr{kv: attribute.String("network.peer.address", val)}
}

func (HTTPClientConnectionDuration) NetworkProtocolVersion(val string) HTTPClientConnectionDurationAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPClientConnectionDuration) URLScheme(val string) HTTPClientConnectionDurationAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientOpenConnections is an instrument used to record metric values
// conforming to the "http.client.open_connections" semantic conventions. It
// represents the number of outbound HTTP connections that are currently active
// or idle on the client.
type HTTPClientOpenConnections struct {
	inst metric.Int64UpDownCounter
}

// NewHTTPClientOpenConnections returns a new HTTPClientOpenConnections
// instrument.
func NewHTTPClientOpenConnections(m metric.Meter) (HTTPClientOpenConnections, error) {
	i, err := m.Int64UpDownCounter(
	    "http.client.open_connections",
	    metric.WithDescription("Number of outbound HTTP connections that are currently active or idle on the client."),
	    metric.WithUnit("{connection}"),
	)
	if err != nil {
	    return HTTPClientOpenConnections{}, err
	}
	return HTTPClientOpenConnections{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPClientOpenConnections) Name() string {
	return "http.client.open_connections"
}

// Name returns the semantic convention unit of the instrument
func (HTTPClientOpenConnections) Unit() string {
	return "{connection}"
}

// Name returns the semantic convention description of the instrument
func (HTTPClientOpenConnections) Description() string {
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
func (m HTTPClientOpenConnections) Add(
    ctx context.Context,
    incr int64,
    httpConnectionState HTTPConnectionStateAttr,
    serverAddress string,
    serverPort int,
    attrs ...HTTPClientOpenConnectionsAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientOpenConnectionsAttrToAttrs(attrs),
			attribute.String("http.connection.state", httpConnectionState),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientOpenConnectionsAttrToAttrs(in []HTTPClientOpenConnectionsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientOpenConnectionsAttr()
	}
	return out
}

type HTTPClientOpenConnectionsAttr interface {
    httpClientOpenConnectionsAttr() attribute.KeyValue
}

func (a attr) httpClientOpenConnectionsAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPClientOpenConnections) NetworkPeerAddress(val string) HTTPClientOpenConnectionsAttr {
	return attr{kv: attribute.String("network.peer.address", val)}
}

func (HTTPClientOpenConnections) NetworkProtocolVersion(val string) HTTPClientOpenConnectionsAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPClientOpenConnections) URLScheme(val string) HTTPClientOpenConnectionsAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientRequestBodySize is an instrument used to record metric values
// conforming to the "http.client.request.body.size" semantic conventions. It
// represents the size of HTTP client request bodies.
type HTTPClientRequestBodySize struct {
	inst metric.histogram
}

// NewHTTPClientRequestBodySize returns a new HTTPClientRequestBodySize
// instrument.
func NewHTTPClientRequestBodySize(m metric.Meter) (HTTPClientRequestBodySize, error) {
	i, err := m.histogram(
	    "http.client.request.body.size",
	    metric.WithDescription("Size of HTTP client request bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return HTTPClientRequestBodySize{}, err
	}
	return HTTPClientRequestBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPClientRequestBodySize) Name() string {
	return "http.client.request.body.size"
}

// Name returns the semantic convention unit of the instrument
func (HTTPClientRequestBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (HTTPClientRequestBodySize) Description() string {
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
func (m HTTPClientRequestBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod HTTPRequestMethodAttr,
    serverAddress string,
    serverPort int,
    attrs ...HTTPClientRequestBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientRequestBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", httpRequestMethod),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientRequestBodySizeAttrToAttrs(in []HTTPClientRequestBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientRequestBodySizeAttr()
	}
	return out
}

type HTTPClientRequestBodySizeAttr interface {
    httpClientRequestBodySizeAttr() attribute.KeyValue
}

func (a attr) httpClientRequestBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPClientRequestBodySize) ErrorType(val ErrorTypeAttr) HTTPClientRequestBodySizeAttr {
	return attr{kv: attribute.String("error.type", val)}
}

func (HTTPClientRequestBodySize) HTTPResponseStatusCode(val int) HTTPClientRequestBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (HTTPClientRequestBodySize) NetworkProtocolName(val string) HTTPClientRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (HTTPClientRequestBodySize) URLTemplate(val string) HTTPClientRequestBodySizeAttr {
	return attr{kv: attribute.String("url.template", val)}
}

func (HTTPClientRequestBodySize) NetworkProtocolVersion(val string) HTTPClientRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPClientRequestBodySize) URLScheme(val string) HTTPClientRequestBodySizeAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

// HTTPClientRequestDuration is an instrument used to record metric values
// conforming to the "http.client.request.duration" semantic conventions. It
// represents the duration of HTTP client requests.
type HTTPClientRequestDuration struct {
	inst metric.histogram
}

// NewHTTPClientRequestDuration returns a new HTTPClientRequestDuration
// instrument.
func NewHTTPClientRequestDuration(m metric.Meter) (HTTPClientRequestDuration, error) {
	i, err := m.histogram(
	    "http.client.request.duration",
	    metric.WithDescription("Duration of HTTP client requests."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return HTTPClientRequestDuration{}, err
	}
	return HTTPClientRequestDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPClientRequestDuration) Name() string {
	return "http.client.request.duration"
}

// Name returns the semantic convention unit of the instrument
func (HTTPClientRequestDuration) Unit() string {
	return "s"
}

// Name returns the semantic convention description of the instrument
func (HTTPClientRequestDuration) Description() string {
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
func (m HTTPClientRequestDuration) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod HTTPRequestMethodAttr,
    serverAddress string,
    serverPort int,
    attrs ...HTTPClientRequestDurationAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientRequestDurationAttrToAttrs(attrs),
			attribute.String("http.request.method", httpRequestMethod),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientRequestDurationAttrToAttrs(in []HTTPClientRequestDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientRequestDurationAttr()
	}
	return out
}

type HTTPClientRequestDurationAttr interface {
    httpClientRequestDurationAttr() attribute.KeyValue
}

func (a attr) httpClientRequestDurationAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPClientRequestDuration) ErrorType(val ErrorTypeAttr) HTTPClientRequestDurationAttr {
	return attr{kv: attribute.String("error.type", val)}
}

func (HTTPClientRequestDuration) HTTPResponseStatusCode(val int) HTTPClientRequestDurationAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (HTTPClientRequestDuration) NetworkProtocolName(val string) HTTPClientRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (HTTPClientRequestDuration) NetworkProtocolVersion(val string) HTTPClientRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPClientRequestDuration) URLScheme(val string) HTTPClientRequestDurationAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

func (HTTPClientRequestDuration) URLTemplate(val string) HTTPClientRequestDurationAttr {
	return attr{kv: attribute.String("url.template", val)}
}

// HTTPClientResponseBodySize is an instrument used to record metric values
// conforming to the "http.client.response.body.size" semantic conventions. It
// represents the size of HTTP client response bodies.
type HTTPClientResponseBodySize struct {
	inst metric.histogram
}

// NewHTTPClientResponseBodySize returns a new HTTPClientResponseBodySize
// instrument.
func NewHTTPClientResponseBodySize(m metric.Meter) (HTTPClientResponseBodySize, error) {
	i, err := m.histogram(
	    "http.client.response.body.size",
	    metric.WithDescription("Size of HTTP client response bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return HTTPClientResponseBodySize{}, err
	}
	return HTTPClientResponseBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPClientResponseBodySize) Name() string {
	return "http.client.response.body.size"
}

// Name returns the semantic convention unit of the instrument
func (HTTPClientResponseBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (HTTPClientResponseBodySize) Description() string {
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
func (m HTTPClientResponseBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod HTTPRequestMethodAttr,
    serverAddress string,
    serverPort int,
    attrs ...HTTPClientResponseBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpClientResponseBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", httpRequestMethod),
			attribute.String("server.address", serverAddress),
			attribute.Int("server.port", serverPort),
		)...,
    ))
}

func httpClientResponseBodySizeAttrToAttrs(in []HTTPClientResponseBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpClientResponseBodySizeAttr()
	}
	return out
}

type HTTPClientResponseBodySizeAttr interface {
    httpClientResponseBodySizeAttr() attribute.KeyValue
}

func (a attr) httpClientResponseBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPClientResponseBodySize) ErrorType(val ErrorTypeAttr) HTTPClientResponseBodySizeAttr {
	return attr{kv: attribute.String("error.type", val)}
}

func (HTTPClientResponseBodySize) HTTPResponseStatusCode(val int) HTTPClientResponseBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (HTTPClientResponseBodySize) NetworkProtocolName(val string) HTTPClientResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (HTTPClientResponseBodySize) URLTemplate(val string) HTTPClientResponseBodySizeAttr {
	return attr{kv: attribute.String("url.template", val)}
}

func (HTTPClientResponseBodySize) NetworkProtocolVersion(val string) HTTPClientResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPClientResponseBodySize) URLScheme(val string) HTTPClientResponseBodySizeAttr {
	return attr{kv: attribute.String("url.scheme", val)}
}

// HTTPServerActiveRequests is an instrument used to record metric values
// conforming to the "http.server.active_requests" semantic conventions. It
// represents the number of active HTTP server requests.
type HTTPServerActiveRequests struct {
	inst metric.Int64UpDownCounter
}

// NewHTTPServerActiveRequests returns a new HTTPServerActiveRequests instrument.
func NewHTTPServerActiveRequests(m metric.Meter) (HTTPServerActiveRequests, error) {
	i, err := m.Int64UpDownCounter(
	    "http.server.active_requests",
	    metric.WithDescription("Number of active HTTP server requests."),
	    metric.WithUnit("{request}"),
	)
	if err != nil {
	    return HTTPServerActiveRequests{}, err
	}
	return HTTPServerActiveRequests{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPServerActiveRequests) Name() string {
	return "http.server.active_requests"
}

// Name returns the semantic convention unit of the instrument
func (HTTPServerActiveRequests) Unit() string {
	return "{request}"
}

// Name returns the semantic convention description of the instrument
func (HTTPServerActiveRequests) Description() string {
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
func (m HTTPServerActiveRequests) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod HTTPRequestMethodAttr,
    urlScheme string,
    attrs ...HTTPServerActiveRequestsAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerActiveRequestsAttrToAttrs(attrs),
			attribute.String("http.request.method", httpRequestMethod),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerActiveRequestsAttrToAttrs(in []HTTPServerActiveRequestsAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerActiveRequestsAttr()
	}
	return out
}

type HTTPServerActiveRequestsAttr interface {
    httpServerActiveRequestsAttr() attribute.KeyValue
}

func (a attr) httpServerActiveRequestsAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPServerActiveRequests) ServerAddress(val string) HTTPServerActiveRequestsAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (HTTPServerActiveRequests) ServerPort(val int) HTTPServerActiveRequestsAttr {
	return attr{kv: attribute.Int("server.port", val)}
}

// HTTPServerRequestBodySize is an instrument used to record metric values
// conforming to the "http.server.request.body.size" semantic conventions. It
// represents the size of HTTP server request bodies.
type HTTPServerRequestBodySize struct {
	inst metric.histogram
}

// NewHTTPServerRequestBodySize returns a new HTTPServerRequestBodySize
// instrument.
func NewHTTPServerRequestBodySize(m metric.Meter) (HTTPServerRequestBodySize, error) {
	i, err := m.histogram(
	    "http.server.request.body.size",
	    metric.WithDescription("Size of HTTP server request bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return HTTPServerRequestBodySize{}, err
	}
	return HTTPServerRequestBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPServerRequestBodySize) Name() string {
	return "http.server.request.body.size"
}

// Name returns the semantic convention unit of the instrument
func (HTTPServerRequestBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (HTTPServerRequestBodySize) Description() string {
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
func (m HTTPServerRequestBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod HTTPRequestMethodAttr,
    urlScheme string,
    attrs ...HTTPServerRequestBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerRequestBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", httpRequestMethod),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerRequestBodySizeAttrToAttrs(in []HTTPServerRequestBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerRequestBodySizeAttr()
	}
	return out
}

type HTTPServerRequestBodySizeAttr interface {
    httpServerRequestBodySizeAttr() attribute.KeyValue
}

func (a attr) httpServerRequestBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPServerRequestBodySize) ErrorType(val ErrorTypeAttr) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.String("error.type", val)}
}

func (HTTPServerRequestBodySize) HTTPResponseStatusCode(val int) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (HTTPServerRequestBodySize) HTTPRoute(val string) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.String("http.route", val)}
}

func (HTTPServerRequestBodySize) NetworkProtocolName(val string) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (HTTPServerRequestBodySize) NetworkProtocolVersion(val string) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPServerRequestBodySize) ServerAddress(val string) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (HTTPServerRequestBodySize) ServerPort(val int) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.Int("server.port", val)}
}

func (HTTPServerRequestBodySize) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) HTTPServerRequestBodySizeAttr {
	return attr{kv: attribute.String("user_agent.synthetic.type", val)}
}

// HTTPServerRequestDuration is an instrument used to record metric values
// conforming to the "http.server.request.duration" semantic conventions. It
// represents the duration of HTTP server requests.
type HTTPServerRequestDuration struct {
	inst metric.histogram
}

// NewHTTPServerRequestDuration returns a new HTTPServerRequestDuration
// instrument.
func NewHTTPServerRequestDuration(m metric.Meter) (HTTPServerRequestDuration, error) {
	i, err := m.histogram(
	    "http.server.request.duration",
	    metric.WithDescription("Duration of HTTP server requests."),
	    metric.WithUnit("s"),
	)
	if err != nil {
	    return HTTPServerRequestDuration{}, err
	}
	return HTTPServerRequestDuration{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPServerRequestDuration) Name() string {
	return "http.server.request.duration"
}

// Name returns the semantic convention unit of the instrument
func (HTTPServerRequestDuration) Unit() string {
	return "s"
}

// Name returns the semantic convention description of the instrument
func (HTTPServerRequestDuration) Description() string {
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
func (m HTTPServerRequestDuration) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod HTTPRequestMethodAttr,
    urlScheme string,
    attrs ...HTTPServerRequestDurationAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerRequestDurationAttrToAttrs(attrs),
			attribute.String("http.request.method", httpRequestMethod),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerRequestDurationAttrToAttrs(in []HTTPServerRequestDurationAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerRequestDurationAttr()
	}
	return out
}

type HTTPServerRequestDurationAttr interface {
    httpServerRequestDurationAttr() attribute.KeyValue
}

func (a attr) httpServerRequestDurationAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPServerRequestDuration) ErrorType(val ErrorTypeAttr) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.String("error.type", val)}
}

func (HTTPServerRequestDuration) HTTPResponseStatusCode(val int) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (HTTPServerRequestDuration) HTTPRoute(val string) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.String("http.route", val)}
}

func (HTTPServerRequestDuration) NetworkProtocolName(val string) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (HTTPServerRequestDuration) NetworkProtocolVersion(val string) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPServerRequestDuration) ServerAddress(val string) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (HTTPServerRequestDuration) ServerPort(val int) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.Int("server.port", val)}
}

func (HTTPServerRequestDuration) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) HTTPServerRequestDurationAttr {
	return attr{kv: attribute.String("user_agent.synthetic.type", val)}
}

// HTTPServerResponseBodySize is an instrument used to record metric values
// conforming to the "http.server.response.body.size" semantic conventions. It
// represents the size of HTTP server response bodies.
type HTTPServerResponseBodySize struct {
	inst metric.histogram
}

// NewHTTPServerResponseBodySize returns a new HTTPServerResponseBodySize
// instrument.
func NewHTTPServerResponseBodySize(m metric.Meter) (HTTPServerResponseBodySize, error) {
	i, err := m.histogram(
	    "http.server.response.body.size",
	    metric.WithDescription("Size of HTTP server response bodies."),
	    metric.WithUnit("By"),
	)
	if err != nil {
	    return HTTPServerResponseBodySize{}, err
	}
	return HTTPServerResponseBodySize{i}, nil
}

// Name returns the semantic convention name of the instrument.
func (HTTPServerResponseBodySize) Name() string {
	return "http.server.response.body.size"
}

// Name returns the semantic convention unit of the instrument
func (HTTPServerResponseBodySize) Unit() string {
	return "By"
}

// Name returns the semantic convention description of the instrument
func (HTTPServerResponseBodySize) Description() string {
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
func (m HTTPServerResponseBodySize) Add(
    ctx context.Context,
    incr int64,
    httpRequestMethod HTTPRequestMethodAttr,
    urlScheme string,
    attrs ...HTTPServerResponseBodySizeAttr,
) {
    m.inst.Add(ctx, incr, metric.WithAttributes(
		append(
			httpServerResponseBodySizeAttrToAttrs(attrs),
			attribute.String("http.request.method", httpRequestMethod),
			attribute.String("url.scheme", urlScheme),
		)...,
    ))
}

func httpServerResponseBodySizeAttrToAttrs(in []HTTPServerResponseBodySizeAttr) []attribute.KeyValue {
	if len(in) == 0 {
		return nil
	}

	out := make([]attribute.KeyValue, len(in))
	for i, a := range in {
		out[i] = a.httpServerResponseBodySizeAttr()
	}
	return out
}

type HTTPServerResponseBodySizeAttr interface {
    httpServerResponseBodySizeAttr() attribute.KeyValue
}

func (a attr) httpServerResponseBodySizeAttr() attribute.KeyValue {
    return a.kv
}

func (HTTPServerResponseBodySize) ErrorType(val ErrorTypeAttr) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.String("error.type", val)}
}

func (HTTPServerResponseBodySize) HTTPResponseStatusCode(val int) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.Int("http.response.status_code", val)}
}

func (HTTPServerResponseBodySize) HTTPRoute(val string) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.String("http.route", val)}
}

func (HTTPServerResponseBodySize) NetworkProtocolName(val string) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.name", val)}
}

func (HTTPServerResponseBodySize) NetworkProtocolVersion(val string) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.String("network.protocol.version", val)}
}

func (HTTPServerResponseBodySize) ServerAddress(val string) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.String("server.address", val)}
}

func (HTTPServerResponseBodySize) ServerPort(val int) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.Int("server.port", val)}
}

func (HTTPServerResponseBodySize) UserAgentSyntheticType(val UserAgentSyntheticTypeAttr) HTTPServerResponseBodySizeAttr {
	return attr{kv: attribute.String("user_agent.synthetic.type", val)}
}

// This is used as a helper for all optional attributes.
type attr struct {
	kv attribute.KeyValue
}