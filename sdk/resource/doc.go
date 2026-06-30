// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package resource provides detecting and representing resources.
//
// The fundamental struct is a Resource which holds identifying information
// about the entities for which telemetry is exported.
//
// To automatically construct Resources from an environment a Detector
// interface is defined. Implementations of this interface can be passed to
// the Detect function to generate a Resource from the merged information.
//
// To load a user defined Resource from environment variables, use
// [Environment], [EnvironmentWithContext], or [WithFromEnv].
//
// # Environment Variables
//
// OTEL_RESOURCE_ATTRIBUTES -
// comma delimited key/value pairs used to describe the entity producing
// telemetry (e.g. `<key1>=<value1>,<key2>=<value2>,...`).
//
// OTEL_SERVICE_NAME -
// sets the `service.name` resource attribute. If `service.name` is also set in
// OTEL_RESOURCE_ATTRIBUTES, OTEL_SERVICE_NAME takes precedence.
//
// See [go.opentelemetry.io/otel/sdk/internal/x] for information about the
// experimental OTEL_GO_X_RESOURCE environment variable.
//
// While this package provides a stable API,
// the attributes added by resource detectors may change.
package resource // import "go.opentelemetry.io/otel/sdk/resource"
