// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the metric package.
package internal // import "go.opentelemetry.io/otel/sdk/metric/internal"

//go:generate gotmpl --body=../../../internal/shared/counter/counter.go.tmpl "--data={ \"pkg\": \"go.opentelemetry.io/otel/sdk/metric/internal/counter\" }" --out=counter/counter.go
//go:generate gotmpl --body=../../../internal/shared/counter/counter_test.go.tmpl "--data={}" --out=counter/counter_test.go
