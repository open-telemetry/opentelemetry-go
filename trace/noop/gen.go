// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package noop // import "go.opentelemetry.io/otel/trace/noop"

//go:generate gotmpl --body=../../internal/shared/noop_helper_test.go.tmpl "--data={\"packageName\": \"noop\"}" --out=noop_helper_test.go
