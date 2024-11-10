// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/exporters/zipkin/internal"

//go:generate gotmpl --body=../../../internal/shared/matchers/expectation.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=matchers/expectation.go
//go:generate gotmpl --body=../../../internal/shared/matchers/expecter.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=matchers/expecter.go
//go:generate gotmpl --body=../../../internal/shared/matchers/temporal_matcher.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=matchers/temporal_matcher.go

//go:generate gotmpl --body=../../../internal/shared/internaltest/alignment.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=internaltest/alignment.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/env.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=internaltest/env.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/env_test.go.tmpl "--data={}" --out=internaltest/env_test.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/errors.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=internaltest/errors.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/harness.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\", \"matchersImportPath\": \"go.opentelemetry.io/otel/exporters/zipkin/internal/matchers\"}" --out=internaltest/harness.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/text_map_carrier.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=internaltest/text_map_carrier.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/text_map_carrier_test.go.tmpl "--data={}" --out=internaltest/text_map_carrier_test.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/text_map_propagator.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/zipkin/internal\"}" --out=internaltest/text_map_propagator.go
//go:generate gotmpl --body=../../../internal/shared/internaltest/text_map_propagator_test.go.tmpl "--data={}" --out=internaltest/text_map_propagator_test.go
