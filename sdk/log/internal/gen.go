// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the sdk/log package.
package internal // import "go.opentelemetry.io/otel/sdk/log/internal"

//go:generate gotmpl --body=../../../internal/shared/x/x.go.tmpl "--data={ \"pkg\": \"go.opentelemetry.io/otel/sdk/log\" }" --out=x/x.go
//go:generate gotmpl --body=../../../internal/shared/x/x_test.go.tmpl "--data={}" --out=x/x_test.go

//go:generate gotmpl --body=../../../internal/shared/counter/counter.go.tmpl "--data={ \"pkg\": \"go.opentelemetry.io/otel/sdk/log\" }" --out=counter/counter.go
//go:generate gotmpl --body=../../../internal/shared/counter/counter_test.go.tmpl "--data={}" --out=counter/counter_test.go
