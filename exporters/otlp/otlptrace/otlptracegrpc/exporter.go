// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlptracegrpc // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
)

// New constructs a new Exporter and starts it.
func New(ctx context.Context, opts ...Option) (*otlptrace.Exporter, error) {
	if err := validateGRPCEndpontFromEnv(); err != nil {
		return nil, err
	}
	client := newClient(opts...)
	return otlptrace.New(ctx, client)
}

// NewUnstarted constructs a new Exporter and does not start it.
func NewUnstarted(opts ...Option) *otlptrace.Exporter {
	return otlptrace.NewUnstarted(NewClient(opts...))
}

func validateGRPCEndpontFromEnv() error {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	if u.Path != "" && u.Path != "/" {
		return fmt.Errorf(
			"invalid OTLP enpoint %q: gRPC endpoint must not include a path", endpoint)

	}
	return nil
}
