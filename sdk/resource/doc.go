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
// # Environment Variables
//
// The environment variables described below are read by [Default],
// [Environment], and the detector added by [WithFromEnv].
// [Default] reads them only once and caches the resulting Resource.
//
// OTEL_RESOURCE_ATTRIBUTES (default: none) -
// a comma-separated list of key=value pairs added as string resource attributes.
// Example value: "key1=value1,key2=value2".
// Keys and values are trimmed of surrounding whitespace, and values are percent-decoded.
// Pairs without "=" are dropped and reported as an error wrapping [ErrPartialResource].
//
// OTEL_SERVICE_NAME (default: none) -
// the value of the "service.name" resource attribute.
// OTEL_SERVICE_NAME takes precedence over a "service.name" set in OTEL_RESOURCE_ATTRIBUTES.
// If neither variable sets "service.name", [Default] uses "unknown_service:"
// followed by the name of the executable.
//
// While this package provides a stable API,
// the attributes added by resource detectors may change.
package resource
