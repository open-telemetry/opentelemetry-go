// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package stdoutlog provides an exporter for OpenTelemetry log
// telemetry.
//
// The exporter is intended to be used for testing and debugging, it is not
// meant for production use. Additionally, it does not provide an interchange
// format for OpenTelemetry that is supported with any stability or
// compatibility guarantees. If these are needed features, please use the OTLP
// exporter instead.
//
// # Self-Observability
//
// The exporter provides a self-observability feature that allows you to monitor
// the exporter itself. To enable this feature, set the environment variable
// OTEL_GO_X_SELF_OBSERVABILITY to the case-insensitive string value of "true"
// or use the WithSelfObservability option when creating the exporter.
//
// When enabled, the exporter will create the following metrics using the global
// MeterProvider:
//
//   - otel.sdk.exporter.log.inflight
//   - otel.sdk.exporter.log.exported
//   - otel.sdk.exporter.operation.duration
//
// Please see the [Semantic conventions for OpenTelemetry SDK metrics] documentation
// for more details on these metrics.
//
// [Semantic conventions for OpenTelemetry SDK metrics]: https://github.com/open-telemetry/semantic-conventions/blob/v1.36.0/docs/otel/sdk-metrics.md
package stdoutlog // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
