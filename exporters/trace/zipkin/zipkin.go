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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel/api/global"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Exporter exports SpanData to the zipkin collector. It implements
// the SpanBatcher interface, so it needs to be used together with the
// WithBatcher option when setting up the exporter pipeline.
type Exporter struct {
	url         string
	serviceName string
	client      *http.Client
	logger      *log.Logger
	o           options
}

var (
	_ export.SpanBatcher = &Exporter{}
)

// Options contains configuration for the exporter.
type options struct {
	client *http.Client
	logger *log.Logger
	config *sdktrace.Config
}

// Option defines a function that configures the exporter.
type Option func(*options)

// WithLogger configures the exporter to use the passed logger.
func WithLogger(logger *log.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

// WithClient configures the exporter to use the passed HTTP client.
func WithClient(client *http.Client) Option {
	return func(opts *options) {
		opts.client = client
	}
}

// WithSDK sets the SDK config for the exporter pipeline.
func WithSDK(config *sdktrace.Config) Option {
	return func(o *options) {
		o.config = config
	}
}

// NewRawExporter creates a new Zipkin exporter.
func NewRawExporter(collectorURL, serviceName string, opts ...Option) (*Exporter, error) {
	if collectorURL == "" {
		return nil, errors.New("collector URL cannot be empty")
	}
	u, err := url.Parse(collectorURL)
	if err != nil {
		return nil, fmt.Errorf("invalid collector URL: %v", err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid collector URL")
	}

	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	if o.client == nil {
		o.client = http.DefaultClient
	}
	return &Exporter{
		url:         collectorURL,
		client:      o.client,
		logger:      o.logger,
		serviceName: serviceName,
		o:           o,
	}, nil
}

// NewExportPipeline sets up a complete export pipeline
// with the recommended setup for trace provider
func NewExportPipeline(collectorURL, serviceName string, opts ...Option) (*sdktrace.Provider, error) {
	exp, err := NewRawExporter(collectorURL, serviceName, opts...)
	if err != nil {
		return nil, err
	}

	batcher := sdktrace.WithBatcher(exp)
	tp, err := sdktrace.NewProvider(batcher)
	if err != nil {
		return nil, err
	}
	if exp.o.config != nil {
		tp.ApplyConfig(*exp.o.config)
	}

	return tp, err
}

// InstallNewPipeline instantiates a NewExportPipeline with the
// recommended configuration and registers it globally.
func InstallNewPipeline(collectorURL, serviceName string, opts ...Option) error {
	tp, err := NewExportPipeline(collectorURL, serviceName, opts...)
	if err != nil {
		return err
	}

	global.SetTraceProvider(tp)
	return nil
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
	req.Header.Set("Content-Type", "application/json")
	resp, err := e.client.Do(req)
	if err != nil {
		e.logf("request to %s failed: %v", e.url, err)
		return
	}
	e.logf("zipkin responded with status %d", resp.StatusCode)

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		e.logf("failed to read response body: %v", err)
	}

	err = resp.Body.Close()
	if err != nil {
		e.logf("failed to close response body: %v", err)
	}
}

func (e *Exporter) logf(format string, args ...interface{}) {
	if e.logger != nil {
		e.logger.Printf(format, args...)
	}
}
