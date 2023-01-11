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

package zipkin // import "go.opentelemetry.io/otel/exporters/zipkin"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	defaultCollectorURL = "http://localhost:9411/api/v2/spans"
)

// Exporter exports spans to the zipkin collector.
type Exporter struct {
	url    string
	client *http.Client
	logger logr.Logger

	stoppedMu sync.RWMutex
	stopped   bool
}

var _ sdktrace.SpanExporter = &Exporter{}

var emptyLogger = logr.Logger{}

// Options contains configuration for the exporter.
type config struct {
	client *http.Client
	logger logr.Logger
}

// Option defines a function that configures the exporter.
type Option interface {
	apply(config) config
}

type optionFunc func(config) config

func (fn optionFunc) apply(cfg config) config {
	return fn(cfg)
}

// WithLogger configures the exporter to use the passed logger.
// WithLogger and WithLogr will overwrite each other.
func WithLogger(logger *log.Logger) Option {
	return WithLogr(stdr.New(logger))
}

// WithLogr configures the exporter to use the passed logr.Logger.
// WithLogr and WithLogger will overwrite each other.
func WithLogr(logger logr.Logger) Option {
	return optionFunc(func(cfg config) config {
		cfg.logger = logger
		return cfg
	})
}

// WithClient configures the exporter to use the passed HTTP client.
func WithClient(client *http.Client) Option {
	return optionFunc(func(cfg config) config {
		cfg.client = client
		return cfg
	})
}

// New creates a new Zipkin exporter.
func New(collectorURL string, opts ...Option) (*Exporter, error) {
	if collectorURL == "" {
		// Use endpoint from env var or default collector URL.
		collectorURL = envOr(envEndpoint, defaultCollectorURL)
	}
	u, err := url.Parse(collectorURL)
	if err != nil {
		return nil, fmt.Errorf("invalid collector URL %q: %v", collectorURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid collector URL %q: no scheme or host", collectorURL)
	}

	cfg := config{}
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}

	if cfg.client == nil {
		cfg.client = http.DefaultClient
	}
	return &Exporter{
		url:    collectorURL,
		client: cfg.client,
		logger: cfg.logger,
	}, nil
}

// ExportSpans exports spans to a Zipkin receiver.
func (e *Exporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	e.stoppedMu.RLock()
	stopped := e.stopped
	e.stoppedMu.RUnlock()
	if stopped {
		e.logf("exporter stopped, not exporting span batch")
		return nil
	}

	if len(spans) == 0 {
		e.logf("no spans to export")
		return nil
	}
	models := SpanModels(spans)
	body, err := json.Marshal(models)
	if err != nil {
		return e.errf("failed to serialize zipkin models to JSON: %v", err)
	}
	e.logf("about to send a POST request to %s with body %s", e.url, body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url, bytes.NewBuffer(body))
	if err != nil {
		return e.errf("failed to create request to %s: %v", e.url, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		return e.errf("request to %s failed: %v", e.url, err)
	}
	defer resp.Body.Close()

	// Zipkin API returns a 202 on success and the content of the body isn't interesting
	// but it is still being read because according to https://golang.org/pkg/net/http/#Response
	// > The default HTTP client's Transport may not reuse HTTP/1.x "keep-alive" TCP connections
	// > if the Body is not read to completion and closed.
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return e.errf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		return e.errf("failed to send spans to zipkin server with status %d", resp.StatusCode)
	}

	return nil
}

// Shutdown stops the exporter flushing any pending exports.
func (e *Exporter) Shutdown(ctx context.Context) error {
	e.stoppedMu.Lock()
	e.stopped = true
	e.stoppedMu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	return nil
}

func (e *Exporter) logf(format string, args ...interface{}) {
	if e.logger != emptyLogger {
		e.logger.Info(format, args)
	}
}

func (e *Exporter) errf(format string, args ...interface{}) error {
	e.logf(format, args...)
	return fmt.Errorf(format, args...)
}

// MarshalLog is the marshaling function used by the logging system to represent this exporter.
func (e *Exporter) MarshalLog() interface{} {
	return struct {
		Type string
		URL  string
	}{
		Type: "zipkin",
		URL:  e.url,
	}
}
