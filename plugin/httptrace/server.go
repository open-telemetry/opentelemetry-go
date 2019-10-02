package httptrace

import (
	"io"
	"net/http"

	"go.opentelemetry.io/api/core"
	"go.opentelemetry.io/api/tag"
	"go.opentelemetry.io/api/trace"
)

var _ http.Handler = &httpHandler{}

// HTTPHandler provides http middleware that corresponds to the http.Handler interface
type httpHandler struct {
	operation string
	handler   http.Handler
	tracer    trace.Tracer
}

// NewHandler wraps the passed handler in an span named after the operation and
// with provided options, and functions as http middleware. The span is tagged
// with:
//   * "read_bytes" - if anything was read from the request body
//   * "wrote_bytes" - if anything was written to the response writer
//   * "http_status" - the http status, if set
// TODO: Add these aatributes: https://github.com/census-instrumentation/opencensus-go/blob/fa651b05963cfb6060755dc887e7d156ba66e792/plugin/ochttp/trace.go#L33-L40
func NewHandler(handler http.Handler, operation string, opts ...HandlerOption) http.Handler {
	c := newHandlerConfig(opts...)

	return &httpHandler{
		operation: operation,
		handler:   handler,
		tracer:    c.tracer,
	}
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO: what if sctx == EmptySpanContext?
	attrs, tags, sctx := Extract(ctx, r)

	// TODO: What is this doing exactly?
	ctx = tag.WithMap(ctx, tag.NewMap(tag.MapUpdate{MultiKV: tags}))

	spanOpts := []trace.SpanOption{
		trace.WithAttributes(attrs...),
		// TODO: If the endpoint is a public endpoint, it should start a new trace
		// and incoming remote sctx should be added as a link (WithLinks(links...) -
		// this option doesn't exist). OR it should somehow indicate to SDK that it
		// is a public endpoint and SDK takes care of using sctx as either a parent
		// or a Link.
		trace.ChildOf(sctx),
	}

	var span trace.Span
	ctx, span = h.tracer.Start(ctx, h.operation, spanOpts...)
	defer span.End()

	r = r.WithContext(ctx)
	bw := &bodyWrapper{rc: r.Body}
	r.Body = wrapBody(bw, r.Body)
	rw := &respWriterWrapper{w: w}

	h.handler.ServeHTTP(rw, r)

	span.SetAttributes(attributes(bw, rw)...)
}

func attributes(bw *bodyWrapper, rw *respWriterWrapper) []core.KeyValue {
	kv := make([]core.KeyValue, 0, 5)
	// TODO: Consider adding an event after each read and write, possibly as an
	// option (defaulting to off), so at to not create needlesly verbose spans.
	if bw.read > 0 {
		kv = append(kv,
			core.KeyValue{
				Key: core.Key{Name: "read_bytes"},
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
				Key: core.Key{Name: "read_error"},
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
				Key: core.Key{Name: "wrote_bytes"},
				Value: core.Value{
					Type:  core.INT64,
					Int64: rw.written,
				},
			},
			core.KeyValue{
				Key: core.Key{Name: "http_status"},
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
				Key: core.Key{Name: "write_error"},
				Value: core.Value{
					Type:   core.STRING,
					String: rw.err.Error(),
				},
			},
		)
	}

	return kv
}

type HandlerOption func(*handlerConfig)

func WithTracer(tracer trace.Tracer) HandlerOption {
	return func(c *handlerConfig) {
		c.tracer = tracer
	}
}

type handlerConfig struct {
	tracer trace.Tracer
}

func newHandlerConfig(opts ...HandlerOption) handlerConfig {
	var c handlerConfig
	defaultOpts := []HandlerOption{
		WithTracer(trace.GlobalTracer()),
	}

	for _, opt := range append(defaultOpts, opts...) {
		opt(&c)
	}

	return c
}
