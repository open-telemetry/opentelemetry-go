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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

// Copied from https://github.com/MrAlias/flow for demo purposes.

import (
	"context"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
)

const (
	startedState = "started"
	endedState   = "ended"

	// DefaultListenPort is the port the HTTP server listens on if not
	// configured with the WithListenAddress option.
	DefaultListenPort = 41820
	// DefaultListenAddress is the listen address of the HTTP server if not
	// configured with the WithListenAddress option.
	DefaultListenAddress = ":41820"
)

type spanProcessor struct {
	wrapped SpanExporter

	idleConnsClosed  chan struct{}
	server           *http.Server
	spanCounter      *prometheus.CounterVec
	exportErrCounter *prometheus.CounterVec
}

// Wrap returns a wrapped version of the downstream SpanExporter with
// telemetry flow reporting. All calls to the returned SpanProcessor will
// introspected for telemetry data and then forwarded to downstream.
func Wrap(downstream SpanExporter, options ...Option) SpanProcessor {
	mux := http.NewServeMux()
	registry := prometheus.NewRegistry()
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	))

	c := newConfig(options)
	sp := &spanProcessor{
		wrapped:         downstream,
		idleConnsClosed: make(chan struct{}),
		server:          &http.Server{Addr: c.address, Handler: mux},
		spanCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "spans_total",
			Help: "The total number of processed spans",
		}, []string{"state"}),
		exportErrCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "failed_export_spans_total",
			Help: "The total number of spans that failed to export",
		}, []string{}),
	}
	registry.MustRegister(sp.spanCounter)
	registry.MustRegister(sp.exportErrCounter)

	go func() {
		switch err := sp.server.ListenAndServe(); err {
		case nil, http.ErrServerClosed:
		default:
			otel.Handle(err)
		}
		close(sp.idleConnsClosed)
	}()

	return sp
}

// OnStart is called when a span is started.
func (sp *spanProcessor) OnStart(parent context.Context, s ReadWriteSpan) {
	sp.spanCounter.WithLabelValues(startedState).Inc()
}

// OnEnd is called when span is finished.
func (sp *spanProcessor) OnEnd(s ReadOnlySpan) {
	sp.spanCounter.WithLabelValues(endedState).Inc()
	spans := []ReadOnlySpan{s}
	err := sp.wrapped.ExportSpans(context.TODO(), spans)
	var errPart *PartialExportError
	if errors.As(err, &errPart) {
		sp.exportErrCounter.WithLabelValues().Add(float64(errPart.RejectedN))
	} else {
		global.Error(err, "failed export", "span-count", len(spans))
	}
}

// Shutdown is called when the SDK shuts down. The telemetry reporting process
// will be halted when this is called.
func (sp *spanProcessor) Shutdown(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- sp.wrapped.Shutdown(ctx)
	}()

	err := sp.server.Shutdown(ctx)
	select {
	case <-ctx.Done():
		// Abandon idle conns if context has expired.
		if err == nil {
			return ctx.Err()
		}
		return err
	case <-sp.idleConnsClosed:
	}

	// Downstream honors ctx timeout, no need to include in select above.
	if e := <-errCh; e != nil {
		// Prioritize downstream error over server shutdown error.
		err = e
	}
	return err
}

// ForceFlush dones nothing.
func (sp *spanProcessor) ForceFlush(ctx context.Context) error { return nil }

type config struct {
	// address is the listen address for the HTTP server.
	address string
}

func newConfig(options []Option) config {
	c := config{
		address: DefaultListenAddress,
	}

	for _, opt := range options {
		c = opt.apply(c)
	}

	return c
}

// Option configures the flow SpanProcessor.
type Option interface {
	apply(config) config
}

type addressOpt string

func (o addressOpt) apply(c config) config {
	c.address = string(o)
	return c
}

// WithListenAddress sets the listen address of the HTTP server.
func WithListenAddress(addr string) Option {
	return addressOpt(addr)
}
