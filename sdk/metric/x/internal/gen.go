// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides generation hooks for the experimental metric SDK.
package internal // import "go.opentelemetry.io/otel/sdk/metric/x/internal"

//go:generate gotmpl --body=../../../../internal/shared/metricx/boundsum.go.tmpl "--data={}" --out=boundsum/boundsum.go
