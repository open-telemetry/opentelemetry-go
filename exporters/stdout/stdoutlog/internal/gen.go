// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the stdoutlog
// package.
package internal // import "go.opentelemetry.io/otel/exporters/stdout/stdoutlog/internal"

//go:generate gotmpl --body=../../../../internal/shared/counter/counter.go.tmpl "--data={}" --out=counter/counter.go
//go:generate gotmpl --body=../../../../internal/shared/counter/counter_test.go.tmpl "--data={}" --out=counter/counter_test.go
