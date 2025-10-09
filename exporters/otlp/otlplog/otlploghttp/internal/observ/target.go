// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ // import "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp/internal/observ"

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ParseTarget parses a target string and returns the extracted host
// (domain address or IP), the target port, or an error.
//
// If no port is specified, -1 is returned.
//
// If no host is specified, an empty string is returned.
//
// The target string is expected to always have the form
//   - "example.com:4318"
//   - "http://example.com:4318"
//   - "https://example.com:4318"
func ParseTarget(target string) (string, int, error) {
	endpoint := getEndpoint(target)
	if endpoint == "" {
		return "", -1, fmt.Errorf("invalid target %q: requires a non-empty endpoint", target)
	}

	if !strings.Contains(endpoint, ":") {
		return endpoint, -1, nil
	}

	host, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return "", -1, fmt.Errorf("invalid host:port %q: %w", endpoint, err)
	}

	const base, bitSize = 10, 16
	port16, err := strconv.ParseUint(portStr, base, bitSize)
	if err != nil {
		return "", -1, fmt.Errorf("invalid port %q: %w", portStr, err)
	}
	port := int(port16)

	return host, port, nil
}

func getEndpoint(target string) string {
	switch {
	case strings.HasPrefix(target, "http://"):
		return strings.TrimPrefix(target, "http://")
	case strings.HasPrefix(target, "https://"):
		return strings.TrimPrefix(target, "https://")
	default:
		return target
	}
}
