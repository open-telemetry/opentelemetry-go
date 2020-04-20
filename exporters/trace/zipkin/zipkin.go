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

package zipkin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	export "go.opentelemetry.io/otel/sdk/export/trace"
)

// Exporter exports SpanData to the zipkin collector. It implements
// the SpanBatcher interface, so it needs to be used together with the
// WithBatcher option when setting up the exporter pipeline.
type Exporter struct {
	url         string
	serviceName string
	client      *http.Client
	logger      *log.Logger
}

var (
	_ export.SpanBatcher = &Exporter{}
)

// Options contains configuration for the exporter.
type Options struct {
	client *http.Client
	logger *log.Logger
}

// Option defines a function that configures the exporter.
type Option func(*Options)

// WithLogger configures the exporter to use the passed logger.
func WithLogger(logger *log.Logger) Option {
	return func(opts *Options) {
		opts.logger = logger
	}
}

// WithClient configures the exporter to use the passed HTTP client.
func WithClient(client *http.Client) Option {
	return func(opts *Options) {
		opts.client = client
	}
}

// NewExporter creates a new zipkin exporter.
func NewExporter(collectorURL string, serviceName string, os ...Option) (*Exporter, error) {
	if _, err := url.Parse(collectorURL); err != nil {
		return nil, fmt.Errorf("invalid collector URL: %v", err)
	}
	if serviceName == "" {
		return nil, fmt.Errorf("service name must be non-empty string")
	}
	opts := Options{}
	for _, o := range os {
		o(&opts)
	}
	if opts.client == nil {
		opts.client = http.DefaultClient
	}
	return &Exporter{
		url:         collectorURL,
		client:      opts.client,
		logger:      opts.logger,
		serviceName: serviceName,
	}, nil
}

// ExportSpans is a part of an implementation of the SpanBatcher
// interface.
func (e *Exporter) ExportSpans(ctx context.Context, batch []*export.SpanData) {
	if len(batch) == 0 {
		e.logf("no spans to export")
		return
	}
	models := toZipkinSpanModels(batch, e.serviceName)
	body, err := json.Marshal(models)
	if err != nil {
		e.logf("failed to serialize zipkin models to JSON: %v", err)
		return
	}
	e.logf("about to send a POST request to %s with body %s", e.url, body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.url, bytes.NewBuffer(body))
	if err != nil {
		e.logf("failed to create request to %s: %v", e.url, err)
		return
	}
	resp, err := e.client.Do(req)
	if err != nil {
		e.logf("request to %s failed: %v", e.url, err)
		return
	}
	e.logf("zipkin responded with status %d", resp.StatusCode)
}

func (e *Exporter) logf(format string, args ...interface{}) {
	if e.logger != nil {
		e.logger.Printf(format, args...)
	}
}
