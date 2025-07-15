// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package trace contains support for OpenTelemetry distributed tracing.

The following assumes a basic familiarity with OpenTelemetry concepts.
See https://opentelemetry.io.

# Self-Observability (Experimental)

The SDK provides a self-observability feature that allows you to monitor the SDK itself.

This feature is experimental and may change in future releases.

To opt-in, set the environment variable OTEL_GO_X_SELF_OBSERVABILITY to "true".

When enabled, the SDK will create following metrics using the global MeterProvider:
  - otel.sdk.span.live
  - otel.sdk.span.started

Please see the [Semantic conventions for OpenTelemetry SDK metrics] documentation for more details on these metrics.

[Semantic conventions for OpenTelemetry SDK metrics]: https://github.com/open-telemetry/semantic-conventions/blob/v1.36.0/docs/otel/sdk-metrics.md
*/
package trace // import "go.opentelemetry.io/otel/sdk/trace"
