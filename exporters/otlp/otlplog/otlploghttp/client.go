// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"

import (
	"context"

	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
)

type client interface {
	UploadMetrics(context.Context, *logpb.ResourceLogs) error
	Shutdown(context.Context) error
}

type shutdownClient struct{}

// Compile-time check shutdownClient implements client.
var _ client = shutdownClient{}

func (c shutdownClient) UploadMetrics(context.Context, *logpb.ResourceLogs) error {
	return nil
}

func (c shutdownClient) Shutdown(context.Context) error {
	return nil
}

type httpClient struct {
	// TODO: implement.
}

// newHTTPClient creates a new HTTP log client.
func newHTTPClient(cfg config) (*httpClient, error) {
	// TODO: implement.
	return &httpClient{}, nil
}

func (c *httpClient) UploadMetrics(context.Context, *logpb.ResourceLogs) error {
	// TODO: implement.
	return nil
}

func (c *httpClient) Shutdown(context.Context) error {
	// TODO: implement.
	return nil
}
