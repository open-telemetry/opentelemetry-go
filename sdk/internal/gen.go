// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the sdk package.
package internal // import "go.opentelemetry.io/otel/sdk/internal"

//go:generate gotmpl --body=../../internal/shared/matchers/expectation.go.tmpl "--data={}" --out=matchers/expectation.go
//go:generate gotmpl --body=../../internal/shared/matchers/expecter.go.tmpl "--data={}" --out=matchers/expecter.go
//go:generate gotmpl --body=../../internal/shared/matchers/temporal_matcher.go.tmpl "--data={}" --out=matchers/temporal_matcher.go
