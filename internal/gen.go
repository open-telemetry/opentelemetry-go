// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/internal"

//go:generate gotmpl --body=./shared/matchers/expectation.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=matchers/expectation.go
//go:generate gotmpl --body=./shared/matchers/expecter.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=matchers/expecter.go
//go:generate gotmpl --body=./shared/matchers/temporal_matcher.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=matchers/temporal_matcher.go

//go:generate gotmpl --body=./shared/internaltest/alignment.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=internaltest/alignment.go
//go:generate gotmpl --body=./shared/internaltest/env.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=internaltest/env.go
//go:generate gotmpl --body=./shared/internaltest/env_test.go.tmpl "--data={}" --out=internaltest/env_test.go
//go:generate gotmpl --body=./shared/internaltest/errors.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=internaltest/errors.go
//go:generate gotmpl --body=./shared/internaltest/harness.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\", \"matchersImportPath\": \"go.opentelemetry.io/otel/internal/matchers\"}" --out=internaltest/harness.go
//go:generate gotmpl --body=./shared/internaltest/text_map_carrier.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=internaltest/text_map_carrier.go
//go:generate gotmpl --body=./shared/internaltest/text_map_carrier_test.go.tmpl "--data={}" --out=internaltest/text_map_carrier_test.go
//go:generate gotmpl --body=./shared/internaltest/text_map_propagator.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/internal\"}" --out=internaltest/text_map_propagator.go
//go:generate gotmpl --body=./shared/internaltest/text_map_propagator_test.go.tmpl "--data={}" --out=internaltest/text_map_propagator_test.go
