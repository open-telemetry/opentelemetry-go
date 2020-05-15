// Copyright The OpenTelemetry Authors
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

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"
)

var _ http.Handler = &Handler{}

// Handler is http middleware that corresponds to the http.Handler interface and
// is designed to wrap a http.Mux (or equivalent), while individual routes on
// the mux are wrapped with WithRouteTag. A Handler will add various attributes
// to the span using the kv.Keys defined in this package.
type Handler struct {
	operation string
	handler   http.Handler

	tracer            trace.Tracer
	propagators       propagation.Propagators
	spanStartOptions  []trace.StartOption
	readEvent         bool
	writeEvent        bool
	filters           []Filter
	spanNameFormatter func(string, *http.Request) string
}

func defaultHandlerFormatter(operation string, _ *http.Request) string {
	return operation
}

// NewHandler wraps the passed handler, functioning like middleware, in a span
// named after the operation and with any provided Options.
func NewHandler(handler http.Handler, operation string, opts ...Option) http.Handler {
	h := Handler{
		handler:   handler,
		operation: operation,
	}

	defaultOpts := []Option{
		WithTracer(global.Tracer("go.opentelemetry.io/plugin/othttp")),
		WithPropagators(global.Propagators()),
		WithSpanOptions(trace.WithSpanKind(trace.SpanKindServer)),
		WithSpanNameFormatter(defaultHandlerFormatter),
	}

	c := NewConfig(append(defaultOpts, opts...)...)
	h.configure(c)

	return &h
}

func (h *Handler) configure(c *Config) {
	h.tracer = c.Tracer
	h.propagators = c.Propagators
	h.spanStartOptions = c.SpanStartOptions
	h.readEvent = c.ReadEvent
	h.writeEvent = c.WriteEvent
	h.filters = c.Filters
	h.spanNameFormatter = c.SpanNameFormatter
}

// ServeHTTP serves HTTP requests (http.Handler)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, f := range h.filters {
		if !f(r) {
			// Simply pass through to the handler if a filter rejects the request
			h.handler.ServeHTTP(w, r)
			return
		}
	}

	opts := append([]trace.StartOption{}, h.spanStartOptions...) // start with the configured options

	ctx := propagation.ExtractHTTP(r.Context(), h.propagators, r.Header)
	ctx, span := h.tracer.Start(ctx, h.spanNameFormatter(h.operation, r), opts...)
	defer span.End()

	readRecordFunc := func(int64) {}
	if h.readEvent {
		readRecordFunc = func(n int64) {
			span.AddEvent(ctx, "read", ReadBytesKey.Int64(n))
		}
	}
	bw := bodyWrapper{ReadCloser: r.Body, record: readRecordFunc}
	r.Body = &bw

	writeRecordFunc := func(int64) {}
	if h.writeEvent {
		writeRecordFunc = func(n int64) {
			span.AddEvent(ctx, "write", WroteBytesKey.Int64(n))
		}
	}

	rww := &respWriterWrapper{ResponseWriter: w, record: writeRecordFunc, ctx: ctx, props: h.propagators}

	setBasicAttributes(span, r)
	span.SetAttributes(RemoteAddrKey.String(r.RemoteAddr))

	h.handler.ServeHTTP(rww, r.WithContext(ctx))

	setAfterServeAttributes(span, bw.read, rww.written, int64(rww.statusCode), bw.err, rww.err)
}

func setAfterServeAttributes(span trace.Span, read, wrote, statusCode int64, rerr, werr error) {
	kv := make([]kv.KeyValue, 0, 5)
	// TODO: Consider adding an event after each read and write, possibly as an
	// option (defaulting to off), so as to not create needlessly verbose spans.
	if read > 0 {
		kv = append(kv, ReadBytesKey.Int64(read))
	}
	if rerr != nil && rerr != io.EOF {
		kv = append(kv, ReadErrorKey.String(rerr.Error()))
	}
	if wrote > 0 {
		kv = append(kv, WroteBytesKey.Int64(wrote))
	}
	if statusCode > 0 {
		kv = append(kv, StatusCodeKey.Int64(statusCode))
	}
	if werr != nil && werr != io.EOF {
		kv = append(kv, WriteErrorKey.String(werr.Error()))
	}
	span.SetAttributes(kv...)
}

// WithRouteTag annotates a span with the provided route name using the
// RouteKey Tag.
func WithRouteTag(route string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		span.SetAttributes(RouteKey.String(route))
		h.ServeHTTP(w, r)
	})
}
