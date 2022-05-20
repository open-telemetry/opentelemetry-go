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

package handler // import "go.opentelemetry.io/otel/example/passthrough/handler"

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Handler is a minimal implementation of the handler and client from
// go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp for demonstration purposes.
// It handles an incoming http request, and makes an outgoing http request.
type Handler struct {
	propagators propagation.TextMapPropagator
	tracer      trace.Tracer
	next        func(r *http.Request)
}

func New(next func(r *http.Request)) *Handler {
	// Like most instrumentation packages, this handler defaults to using the
	// global progatators and tracer providers.
	return &Handler{
		propagators: otel.GetTextMapPropagator(),
		tracer:      otel.Tracer("examples/passthrough/handler"),
		next:        next,
	}
}

// HandleHTTPReq mimics what an instrumented http server does.
func (h *Handler) HandleHTTPReq(r *http.Request) {
	ctx := h.propagators.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	var span trace.Span
	log.Println("The \"handle passthrough request\" span should NOT be recorded, because it is recorded by a TracerProvider not backed by the SDK.")
	ctx, span = h.tracer.Start(ctx, "handle passthrough request")
	defer span.End()

	// Pretend to do work
	time.Sleep(time.Second)

	h.makeOutgoingRequest(ctx)
}

// makeOutgoingRequest mimics what an instrumented http client does.
func (h *Handler) makeOutgoingRequest(ctx context.Context) {
	// make a new http request
	r, err := http.NewRequest("", "", nil)
	if err != nil {
		panic(err)
	}

	log.Println("The \"make outgoing request from passthrough\" span should NOT be recorded, because it is recorded by a TracerProvider not backed by the SDK.")
	ctx, span := h.tracer.Start(ctx, "make outgoing request from passthrough")
	defer span.End()
	r = r.WithContext(ctx)
	h.propagators.Inject(ctx, propagation.HeaderCarrier(r.Header))
	h.next(r)
}
