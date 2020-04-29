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
	"context"
	"io"
	"net/http"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/propagation"
	"go.opentelemetry.io/otel/api/trace"

	"google.golang.org/grpc/codes"
)

// Transport implements the http.RoundTripper interface and wraps
// outbound HTTP(S) requests with a span.
type Transport struct {
	rt http.RoundTripper

	tracer            trace.Tracer
	propagators       propagation.Propagators
	spanStartOptions  []trace.StartOption
	filters           []Filter
	spanNameFormatter func(string, *http.Request) string
}

var _ http.RoundTripper = &Transport{}

// NewTransport wraps the provided http.RoundTripper with one that
// starts a span and injects the span context into the outbound request headers.
func NewTransport(base http.RoundTripper, opts ...Option) *Transport {
	t := Transport{
		rt: base,
	}

	defaultOpts := []Option{
		WithTracer(global.Tracer("go.opentelemetry.io/plugin/othttp")),
		WithPropagators(global.Propagators()),
		WithSpanOptions(trace.WithSpanKind(trace.SpanKindClient)),
		WithSpanNameFormatter(defaultTransportFormatter),
	}

	c := NewConfig(append(defaultOpts, opts...)...)
	t.configure(c)

	return &t
}

func (t *Transport) configure(c *Config) {
	t.tracer = c.Tracer
	t.propagators = c.Propagators
	t.spanStartOptions = c.SpanStartOptions
	t.filters = c.Filters
	t.spanNameFormatter = c.SpanNameFormatter
}

func defaultTransportFormatter(_ string, r *http.Request) string {
	return r.Method
}

// RoundTrip creates a Span and propagates its context via the provided request's headers
// before handing the request to the configured base RoundTripper. The created span will
// end when the response body is closed or when a read from the body returns io.EOF.
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	for _, f := range t.filters {
		if !f(r) {
			// Simply pass through to the base RoundTripper if a filter rejects the request
			return t.rt.RoundTrip(r)
		}
	}

	opts := append([]trace.StartOption{}, t.spanStartOptions...) // start with the configured options

	ctx, span := t.tracer.Start(r.Context(), t.spanNameFormatter("", r), opts...)

	r = r.WithContext(ctx)
	setBasicAttributes(span, r)
	propagation.InjectHTTP(ctx, t.propagators, r.Header)

	res, err := t.rt.RoundTrip(r)
	if err != nil {
		span.RecordError(ctx, err, trace.WithErrorStatus(codes.Internal))
		span.End()
		return res, err
	}

	span.SetAttributes(StatusCodeKey.Int(res.StatusCode))
	res.Body = &wrappedBody{ctx: ctx, span: span, body: res.Body}

	return res, err
}

type wrappedBody struct {
	ctx  context.Context
	span trace.Span
	body io.ReadCloser
}

var _ io.ReadCloser = &wrappedBody{}

func (wb *wrappedBody) Read(b []byte) (int, error) {
	n, err := wb.body.Read(b)

	switch err {
	case nil:
		// nothing to do here but fall through to the return
	case io.EOF:
		wb.span.End()
	default:
		wb.span.RecordError(wb.ctx, err, trace.WithErrorStatus(codes.Internal))
	}
	return n, err
}

func (wb *wrappedBody) Close() error {
	wb.span.End()
	return wb.body.Close()
}
