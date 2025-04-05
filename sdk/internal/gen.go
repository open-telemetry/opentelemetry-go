// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal provides internal functionality for the sdk package.
package internal // import "go.opentelemetry.io/otel/sdk/internal"

//go:generate gotmpl --body=../../internal/shared/matchers/expectation.go.tmpl "--data={}" --out=matchers/expectation.go
//go:generate gotmpl --body=../../internal/shared/matchers/expecter.go.tmpl "--data={}" --out=matchers/expecter.go
//go:generate gotmpl --body=../../internal/shared/matchers/temporal_matcher.go.tmpl "--data={}" --out=matchers/temporal_matcher.go

//go:generate gotmpl --body=../../internal/shared/internaltest/doc.go.tmpl "--data={}" --out=internaltest/doc.go
//go:generate gotmpl --body=../../internal/shared/internaltest/env.go.tmpl "--data={}" --out=internaltest/env.go
//go:generate gotmpl --body=../../internal/shared/internaltest/env_test.go.tmpl "--data={}" --out=internaltest/env_test.go
//go:generate gotmpl --body=../../internal/shared/internaltest/errors.go.tmpl "--data={}" --out=internaltest/errors.go
//go:generate gotmpl --body=../../internal/shared/internaltest/harness.go.tmpl "--data={\"matchersImportPath\": \"go.opentelemetry.io/otel/sdk/internal/matchers\"}" --out=internaltest/harness.go
//go:generate gotmpl --body=../../internal/shared/internaltest/text_map_carrier.go.tmpl "--data={}" --out=internaltest/text_map_carrier.go
//go:generate gotmpl --body=../../internal/shared/internaltest/text_map_carrier_test.go.tmpl "--data={}" --out=internaltest/text_map_carrier_test.go
//go:generate gotmpl --body=../../internal/shared/internaltest/text_map_propagator.go.tmpl "--data={}" --out=internaltest/text_map_propagator.go
//go:generate gotmpl --body=../../internal/shared/internaltest/text_map_propagator_test.go.tmpl "--data={}" --out=internaltest/text_map_propagator_test.go
