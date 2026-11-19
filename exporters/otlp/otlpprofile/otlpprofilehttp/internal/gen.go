// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the otlpprofilehttp
// package.
package internal

//go:generate gotmpl --body=../../../../../internal/shared/otlp/response.go.tmpl "--data={}" --out=response.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/response_test.go.tmpl "--data={}" --out=response_test.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/retry/retry.go.tmpl "--data={}" --out=retry/retry.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/retry/retry_test.go.tmpl "--data={}" --out=retry/retry_test.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpprofile/transform/profile.go.tmpl "--data={ \"protoRoot\": \"go.opentelemetry.io/proto/slim/otlp\" }" --out=transform/profile.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpprofile/transform/profile_test.go.tmpl "--data={ \"protoRoot\": \"go.opentelemetry.io/proto/slim/otlp\" }" --out=transform/profile_test.go
