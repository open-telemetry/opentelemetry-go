// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package zipkin // import "go.opentelemetry.io/otel/exporters/zipkin"

import "os"

// Environment variable names.
const (
	// Endpoint for Zipkin collector.
	envEndpoint = "OTEL_EXPORTER_ZIPKIN_ENDPOINT"
)

// envOr returns an env variable's value if it is exists or the default if not.
func envOr(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
