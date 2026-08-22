// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the propagation package.
package internal

//go:generate gotmpl --body=../../internal/shared/hextable/hextable.go.tmpl "--data={ \"pkg\": \"go.opentelemetry.io/otel/propagation\" }" --out=hextable/hextable.go
