// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

/*
Package x documents experimental features for [go.opentelemetry.io/otel/sdk/trace].

# Compatibility and Stability

Experimental features do not fall within the scope of the [OpenTelemetry Go versioning and stability policy].
These features may change in backwards incompatible ways as feedback is applied.
These features may be removed or modified in successive version releases, including patch versions.

When an experimental feature is promoted to a stable feature, a migration path will be included in the changelog entry of the release.
There is no guarantee that any environment variable feature flags that enabled the experimental feature will be supported by the stable version.
If they are supported, they may be accompanied with a deprecation notice stating a timeline for the removal of that support.

# Self-Observability

The SDK provides a self-observability feature that allows you to monitor the SDK itself.

To opt-in, set the environment variable OTEL_GO_X_SELF_OBSERVABILITY to "true".

When enabled, the SDK will create following metrics using the global MeterProvider:
  - otel.sdk.span.live
  - otel.sdk.span.started

Please see the [Semantic conventions for OpenTelemetry SDK metrics] documentation for more details on these metrics.

[OpenTelemetry Go versioning and stability policy]: https://github.com/open-telemetry/opentelemetry-go/blob/main/VERSIONING.md
[Semantic conventions for OpenTelemetry SDK metrics]: https://github.com/open-telemetry/semantic-conventions/blob/v1.36.0/docs/otel/sdk-metrics.md
*/
package x // import "go.opentelemetry.io/otel/sdk/trace/internal/x"
