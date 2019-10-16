// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package othttp

import (
	"io"
	"net/http"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/propagation"
	"go.opentelemetry.io/api/trace"
	prop "go.opentelemetry.io/propagation"
)

var _ http.Handler = &HTTPHandler{}

type httpEvent int

// Possible message events that can be enabled via WithMessageEvents
const (
	EventRead  httpEvent = iota // An event that records the number of bytes read is created for every Read
	EventWrite                  // an event that records the number of bytes written is created for every Write
)

// Attribute keys that HTTPHandler could write out.
const (
	HostKeyName       = "http.host"        // the http host (http.Request.Host)
	MethodKeyName     = "http.method"      // the http method (http.Request.Method)
	PathKeyName       = "http.path"        // the http path (http.Request.URL.Path)
	URLKeyName        = "http.url"         // the http url (http.Request.URL.String())
	UserAgentKeyName  = "http.user_agent"  // the http user agent (http.Request.UserAgent())
	RouteKeyName      = "http.route"       // the http route (ex: /users/:id)
	StatusCodeKeyName = "http.status_code" // if set, the http status
	ReadBytesKeyName  = "http.read_bytes"  // if anything was read from the request body, the total number of bytes read
	ReadErrorKeyName  = "http.read_error"  // If an error occurred while reading a request, the string of the error (io.EOF is not recorded)
	WroteBytesKeyName = "http.wrote_bytes" // if anything was written to the response writer, the total number of bytes written
	WriteErrorKeyName = "http.write_error" // if an error occurred while writing a reply, the string of the error (io.EOF is not recorded)
)

// HTTPHandler provides http middleware that corresponds to the http.Handler interface
type HTTPHandler struct {
	operation string
	handler   http.Handler

	tracer      trace.Tracer
	prop        propagation.TextFormatPropagator
	spanOptions []trace.SpanOption
	public      bool
	readEvent   bool
	writeEvent  bool
}

type HandlerOption func(*HTTPHandler)

// WithTracer configures the HTTPHandler with a specific tracer. If this option
// isn't specified then global tracer is used.
func WithTracer(tracer trace.Tracer) HandlerOption {
	return func(h *HTTPHandler) {
		h.tracer = tracer
	}
}

// WithPublicEndpoint configures the HTTPHandler to link the span with an
// incoming span context. If this option is not provided (the default), then the
// association is a child association (instead of a link).
func WithPublicEndpoint() HandlerOption {
	return func(h *HTTPHandler) {
		h.public = true
	}
}

// WithPropagator configures the HTTPHandler with a specific propagator. If this
// option isn't specificed then a w3c trace context propagator.
func WithPropagator(p propagation.TextFormatPropagator) HandlerOption {
	return func(h *HTTPHandler) {
		h.prop = p
	}
}

// WithSpanOptions configures the HTTPHandler with an additional set of
// trace.SpanOptions, which are applied to each new span.
func WithSpanOptions(opts ...trace.SpanOption) HandlerOption {
	return func(h *HTTPHandler) {
		h.spanOptions = opts
	}
}

// WithMessageEvents configures the HTTPHandler with a set of message events. By
// default only the summary attributes are added at the end of the request.
func WithMessageEvents(events ...httpEvent) HandlerOption {
	return func(h *HTTPHandler) {
		for _, e := range events {
			switch e {
			case EventRead:
				h.readEvent = true
			case EventWrite:
				h.writeEvent = true
			}
		}
	}
}

// NewHandler wraps the passed handler, functioning like middleware, in a span
// named after the operation and with any provided HandlerOptions.
func NewHandler(handler http.Handler, operation string, opts ...HandlerOption) http.Handler {
	h := HTTPHandler{handler: handler}
	defaultOpts := []HandlerOption{
		WithTracer(trace.GlobalTracer()),
		WithPropagator(prop.HTTPTraceContextPropagator{}),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(&h)
	}
	return &h
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	opts := append([]trace.SpanOption{}, h.spanOptions...) // start with the configured options

	sc := h.prop.Extract(r.Context(), r.Header)
	if sc.IsValid() { // not a valid span context, so no link / parent relationship to establish
		var opt trace.SpanOption
		if h.public {
			// TODO: If the endpoint is a public endpoint, it should start a new trace
			// and incoming remote sctx should be added as a link
			// (WithLinks(links...), this option doesn't exist yet). Replace ChildOf
			// below with something like: opt = trace.WithLinks(sc)
			opt = trace.ChildOf(sc)
		} else { // not a private endpoint, so assume child relationship
			opt = trace.ChildOf(sc)
		}
		opts = append(opts, opt)
	}

	ctx, span := h.tracer.Start(r.Context(), h.operation, opts...)
	defer span.End()

	readRecordFunc := func(int) {}
	if h.readEvent {
		readRecordFunc = func(n int) {
			span.AddEvent(ctx, "read", core.KeyValue{
				Key: core.Key{Name: ReadBytesKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: int64(n),
				}})
		}
	}
	bw := bodyWrapper{ReadCloser: r.Body, record: readRecordFunc}
	r.Body = &bw

	writeRecordFunc := func(int) {}
	if h.writeEvent {
		writeRecordFunc = func(n int) {
			span.AddEvent(ctx, "write", core.KeyValue{
				Key: core.Key{Name: WroteBytesKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: int64(n),
				},
			})
		}
	}
	rww := &respWriterWrapper{ResponseWriter: w, record: writeRecordFunc}

	setBeforeServeAttributes(span, r.Host, r.Method, r.URL.Path, r.URL.String(), r.UserAgent())
	// inject the response header before calling ServeHTTP because a Write in
	// ServeHTTP will cause all headers to be written out.
	h.prop.Inject(ctx, rww.Header())

	h.handler.ServeHTTP(rww, r.WithContext(ctx))
	setAfterServeAttributes(span, bw.read, rww.written, int64(rww.statusCode), bw.err, rww.err)
}

func setBeforeServeAttributes(span trace.Span, host, method, path, url, uagent string) {
	// Setup basic span attributes before calling handler.ServeHTTP so that they
	// are available to be mutated by the handler if needed.
	span.SetAttributes(
		core.KeyValue{
			Key: core.Key{Name: HostKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: host,
			}},
		core.KeyValue{
			Key: core.Key{Name: MethodKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: method,
			}},
		core.KeyValue{
			Key: core.Key{Name: PathKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: path,
			}},
		core.KeyValue{
			Key: core.Key{Name: URLKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: url,
			}},
		core.KeyValue{
			Key: core.Key{Name: UserAgentKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: uagent,
			}},
	)
}

func setAfterServeAttributes(span trace.Span, read, wrote, statusCode int64, rerr, werr error) {
	kv := make([]core.KeyValue, 0, 5)
	// TODO: Consider adding an event after each read and write, possibly as an
	// option (defaulting to off), so at to not create needlesly verbose spans.
	if read > 0 {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: ReadBytesKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: read,
				}})
	}

	if rerr != nil && rerr != io.EOF {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: ReadErrorKeyName},
				Value: core.Value{
					Type:   core.STRING,
					String: rerr.Error(),
				}})
	}

	if wrote > 0 {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: WroteBytesKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: wrote,
				}})
	}

	if statusCode > 0 {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: StatusCodeKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: statusCode,
				}})
	}

	if werr != nil && werr != io.EOF {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: WriteErrorKeyName},
				Value: core.Value{
					Type:   core.STRING,
					String: werr.Error(),
				}})
	}

	span.SetAttributes(kv...)
}

func WithRouteTag(route string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.CurrentSpan(r.Context())
		//TODO: Why doesn't tag.Upsert work?
		span.SetAttribute(
			core.KeyValue{
				Key: core.Key{Name: RouteKeyName},
				Value: core.Value{
					Type:   core.STRING,
					String: route,
				},
			},
		)

		h.ServeHTTP(w, r.WithContext(trace.SetCurrentSpan(r.Context(), span)))
	})
}
