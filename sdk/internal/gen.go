// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the sdk package.
package internal // import "go.opentelemetry.io/otel/sdk/internal"

//go:generate gotmpl --body=../../internal/shared/x/x.go.tmpl "--data={ \"pkg\": \"go.opentelemetry.io/otel/sdk\" }" --out=x/x.go
//go:generate gotmpl --body=../../internal/shared/x/x_test.go.tmpl "--data={}" --out=x/x_test.go
