// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

//go:generate gotmpl --body=../../internal/shared/log/record.go.tmpl "--data={\"package\": \"log\", \"type\": \"Record\"}"  --out=record.go
//go:generate gotmpl --body=../../internal/shared/log/record_test.go.tmpl "--data={\"package\": \"log\", \"type\": \"Record\"}" --out=record_test.go
