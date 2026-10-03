// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package zipkin contains an OpenTelemetry tracing exporter for Zipkin.
//
// The environment variables described below can be used for configuration.
//
// OTEL_EXPORTER_ZIPKIN_ENDPOINT (default: "http://localhost:9411/api/v2/spans") -
// target URL to which the exporter sends spans.
// The value must contain a scheme and host.
// The environment variable is only used when [New] is called with an empty collectorURL.
//
// Deprecated: The zipkin exporter is deprecated and will be removed in early 2027.
// See the blog post "[Deprecating Zipkin Exporter]".
//
// [Deprecating Zipkin Exporter]: https://opentelemetry.io/blog/2025/deprecating-zipkin-exporters/
package zipkin
