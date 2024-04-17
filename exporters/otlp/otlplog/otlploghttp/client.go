// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"

import (
	"context"

	logpb "go.opentelemetry.io/proto/otlp/logs/v1"
)

type client struct {
	uploadLogs func(context.Context, []*logpb.ResourceLogs) error
}

func (c *client) UploadLogs(ctx context.Context, rl []*logpb.ResourceLogs) error {
	if c.uploadLogs != nil {
		return c.uploadLogs(ctx, rl)
	}
	return nil
}

func newNoopClient() *client {
	return &client{}
}

// newHTTPClient creates a new HTTP log client.
func newHTTPClient(cfg config) (*client, error) {
	// TODO: implement.
	return &client{}, nil
}
