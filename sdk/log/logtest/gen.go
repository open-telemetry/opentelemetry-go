// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logtest // import "go.opentelemetry.io/otel/sdk/log/logtest"

//go:generate gotmpl --body=../../../internal/shared/log/record.go.tmpl "--data={\"package\": \"logtest\", \"type\": \"record\"}"  --out=record.go
//go:generate gotmpl --body=../../../internal/shared/log/record_test.go.tmpl "--data={\"package\": \"logtest\", \"type\": \"record\"}" --out=record_test.go
