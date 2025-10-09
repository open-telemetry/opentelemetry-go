// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package observ

import (
	"testing"
)

func TestParseTarget(t *testing.T) {
	testcases := []struct {
		name   string
		target string
		host   string
		port   int
	}{
		{"Simple host and port", "example.com:4318", "example.com", 4318},
		{"Host without port", "example.com", "example.com", -1},
		{"HTTP prefix", "http://example.com:4318", "example.com", 4318},
		{"HTTPS prefix", "https://example.com:4318", "example.com", 4318},
		{"HTTP prefix without port", "http://example.com", "example.com", -1},
		{"HTTPS prefix without port", "https://example.com", "example.com", -1},
		{"IPv4 with port", "192.168.1.1:4318", "192.168.1.1", 4318},
		{"IPv4 without port", "192.168.1.1", "192.168.1.1", -1},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			host, port, err := ParseTarget(tc.target)
			if err != nil {
				t.Errorf("unexpected error for target %q: %v", tc.target, err)
				return
			}

			if host != tc.host {
				t.Errorf("expected host %q, got %q for target %q", tc.host, host, tc.target)
			}

			if port != tc.port {
				t.Errorf("expected port %d, got %d for target %q", tc.port, port, tc.target)
			}
		})
	}
}

func TestParseTargetErrors(t *testing.T) {
	targets := []string{
		"example.com:",              // Host with port separator.
		"",                          // Empty target.
		"example.com:invalid",       // Non-numeric port in URL.
		"example.com:-1",            // Port out of range
		"http://localhost:invalid",  // Non-numeric port in URL for http
		"https://localhost:invalid", // Non-numeric port in URL for http
	}
	for _, target := range targets {
		host, port, err := ParseTarget(target)
		if err == nil {
			t.Errorf("parseTarget(%q) expected error, got nil", target)
		}

		if host != "" {
			t.Errorf("parseTarget(%q) host = %q, want empty", target, host)
		}

		if port != -1 {
			t.Errorf("parseTarget(%q) port = %d, want -1", target, port)
		}
	}
}

func TestGetEndpoint(t *testing.T) {
	cases := []struct {
		name     string
		target   string
		expected string
	}{
		{
			name:     "HTTP prefix",
			target:   "http://example.com:4318",
			expected: "example.com:4318",
		},
		{
			name:     "HTTPS prefix",
			target:   "https://example.com:4318",
			expected: "example.com:4318",
		},
		{
			name:     "No prefix",
			target:   "example.com:4318",
			expected: "example.com:4318",
		},
		{
			name:     "Empty string",
			target:   "",
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := getEndpoint(tc.target)
			if result != tc.expected {
				t.Errorf("expected %q, got %q for target %q", tc.expected, result, tc.target)
			}
		})
	}
}
