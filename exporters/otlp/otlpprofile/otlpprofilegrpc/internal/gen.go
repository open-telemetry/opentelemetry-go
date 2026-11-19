// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the otlpprofilegrpc
// package.
package internal

//go:generate gotmpl --body=../../../../../internal/shared/otlp/grpcendpoint/grpcendpoint.go.tmpl "--data={ \"pkg\": \"grpcendpoint\" }" --out=grpcendpoint/grpcendpoint.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/grpcendpoint/grpcendpoint_test.go.tmpl "--data={ \"pkg\": \"grpcendpoint\" }" --out=grpcendpoint/grpcendpoint_test.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/retry/retry.go.tmpl "--data={}" --out=retry/retry.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/retry/retry_test.go.tmpl "--data={}" --out=retry/retry_test.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpprofile/transform/profile.go.tmpl "--data={ \"protoRoot\": \"go.opentelemetry.io/proto/otlp\" }" --out=transform/profile.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpprofile/transform/profile_test.go.tmpl "--data={ \"protoRoot\": \"go.opentelemetry.io/proto/otlp\" }" --out=transform/profile_test.go
