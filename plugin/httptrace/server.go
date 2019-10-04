package httptrace

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

// IsPublicEndpoint configures the HTTPHandler to link the span with an
// incoming span context. If this option is not provided (the default), then the
// association is a child association (instead of a link).
func IsPublicEndpoint() HandlerOption {
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
	var h HTTPHandler
	defaultOpts := []HandlerOption{
		WithTracer(trace.GlobalTracer()),
		WithPropagator(prop.HttpTraceContextPropagator()),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(&h)
	}

	h.handler = handler
	return &h
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	opts := append([]trace.SpanOption{}, h.spanOptions...) // start with the configured options

	sc := h.prop.Extract(ctx, r.Header)
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
	var span trace.Span
	ctx, span = h.tracer.Start(ctx, h.operation, opts...)
	defer span.End()

	r = r.WithContext(ctx)

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
	bw := &bodyWrapper{rc: r.Body, record: readRecordFunc}
	r.Body = wrapBody(bw, r.Body)

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
	rw := &respWriterWrapper{w: w, record: writeRecordFunc}

	span.SetAttributes(
		core.KeyValue{
			Key: core.Key{Name: HostKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: r.Host,
			}},
		core.KeyValue{
			Key: core.Key{Name: MethodKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: r.Method,
			}},
		core.KeyValue{
			Key: core.Key{Name: PathKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: r.URL.Path,
			}},
		core.KeyValue{
			Key: core.Key{Name: URLKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: r.URL.String(),
			}},
		core.KeyValue{
			Key: core.Key{Name: UserAgentKeyName},
			Value: core.Value{
				Type:   core.STRING,
				String: r.UserAgent(),
			}},
	)

	// inject the response header before because calling ServeHTTP because a
	// Write in ServeHTTP will cause all headers to be written out.
	h.prop.Inject(ctx, rw.Header())

	h.handler.ServeHTTP(rw, r)
	span.SetAttributes(afterServeAttributes(bw, rw)...)
}

func afterServeAttributes(bw *bodyWrapper, rw *respWriterWrapper) []core.KeyValue {
	kv := make([]core.KeyValue, 0, 5)
	// TODO: Consider adding an event after each read and write, possibly as an
	// option (defaulting to off), so at to not create needlesly verbose spans.
	if bw.read > 0 {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: ReadBytesKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: bw.read,
				},
			},
		)
	}

	if bw.err != nil && bw.err != io.EOF {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: ReadErrorKeyName},
				Value: core.Value{
					Type:   core.STRING,
					String: bw.err.Error(),
				},
			},
		)
	}

	if rw.wroteHeader {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: WroteBytesKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: rw.written,
				},
			},
			core.KeyValue{
				Key: core.Key{Name: StatusCodeKeyName},
				Value: core.Value{
					Type:  core.INT64,
					Int64: int64(rw.statusCode),
				},
			},
		)
	}

	if rw.err != nil && rw.err != io.EOF {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: WriteErrorKeyName},
				Value: core.Value{
					Type:   core.STRING,
					String: rw.err.Error(),
				},
			},
		)
	}

	return kv
}

func WithRouteTag(route string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.CurrentSpan(ctx)
		//TODO: Why doesn't tag.Upset work?
		span.SetAttribute(
			core.KeyValue{
				Key: core.Key{Name: RouteKeyName},
				Value: core.Value{
					Type:   core.STRING,
					String: route,
				},
			},
		)
		r = r.WithContext(trace.SetCurrentSpan(ctx, span))
		h.ServeHTTP(w, r)
	})
}
