package httptrace

import (
	"net/http"

	"go.opentelemetry.io/api/tag"
	"go.opentelemetry.io/api/trace"
)

type Handler struct {
	handler http.Handler
	tracer  trace.Tracer
}

func NewHandler(handler http.Handler, opts ...HandlerOption) *Handler {
	c := newHandlerConfig(opts...)

	return &Handler{
		handler: handler,
		tracer:  c.tracer,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// TODO: what if sctx == EmptySpanContext?
	attrs, tags, sctx := Extract(ctx, r)

	// TODO: what's going on here?
	r = r.WithContext(tag.WithMap(ctx, tag.NewMap(tag.MapUpdate{
		MultiKV: tags,
	})))

	// TODO: flesh this out
	operationName := r.URL.String()

	spanOpts := []trace.SpanOption{
		trace.WithAttributes(attrs...),
		trace.ChildOf(sctx),
	}
	ctx, span := h.tracer.Start(ctx, operationName, spanOpts...)
	defer span.Finish()

	r = r.WithContext(ctx)

	h.handler.ServeHTTP(w, r)
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
