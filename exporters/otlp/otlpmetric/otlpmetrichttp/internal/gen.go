// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal"

//go:generate gotmpl --body=../../../../../internal/shared/otlp/partialsuccess.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp\"}" --out=partialsuccess.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/partialsuccess_test.go.tmpl "--data={}" --out=partialsuccess_test.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/retry/retry.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=retry/retry.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/retry/retry_test.go.tmpl "--data={}" --out=retry/retry_test.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/envconfig/envconfig.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=envconfig/envconfig.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/envconfig/envconfig_test.go.tmpl "--data={}" --out=envconfig/envconfig_test.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/oconf/envconfig.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\", \"envconfigImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/envconfig\"}" --out=oconf/envconfig.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/oconf/envconfig_test.go.tmpl "--data={}" --out=oconf/envconfig_test.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/oconf/options.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\", \"retryImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/retry\"}" --out=oconf/options.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/oconf/options_test.go.tmpl "--data={\"envconfigImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/envconfig\"}" --out=oconf/options_test.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/oconf/optiontypes.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=oconf/optiontypes.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/oconf/tls.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=oconf/tls.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/otest/client.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=otest/client.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/otest/client_test.go.tmpl "--data={\"internalImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=otest/client_test.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/otest/collector.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\", \"oconfImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal/oconf\"}" --out=otest/collector.go

//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/transform/attribute.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=transform/attribute.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/transform/attribute_test.go.tmpl "--data={}" --out=transform/attribute_test.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/transform/error.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=transform/error.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/transform/error_test.go.tmpl "--data={}" --out=transform/error_test.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/transform/metricdata.go.tmpl "--data={\"parentPkg\": \"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp/internal\"}" --out=transform/metricdata.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlpmetric/transform/metricdata_test.go.tmpl "--data={}" --out=transform/metricdata_test.go
