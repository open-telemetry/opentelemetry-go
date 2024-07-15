// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

//go:generate gotmpl --body=../../internal/shared/noop_helper_test.go.tmpl "--data={\"packageName\": \"log\"}" --out=noop_helper_test.go
