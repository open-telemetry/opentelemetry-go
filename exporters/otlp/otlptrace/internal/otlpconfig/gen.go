// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlpconfig // import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/otlpconfig"

//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlptrace/otlpconfig/envconfig.go.tmpl "--data={\"envconfigImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/envconfig\"}" --out=envconfig.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlptrace/otlpconfig/options.go.tmpl "--data={\"retryImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/retry\"}" --out=options.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlptrace/otlpconfig/options_test.go.tmpl "--data={\"envconfigImportPath\": \"go.opentelemetry.io/otel/exporters/otlp/otlptrace/internal/envconfig\"}" --out=options_test.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlptrace/otlpconfig/optiontypes.go.tmpl "--data={}" --out=optiontypes.go
//go:generate gotmpl --body=../../../../../internal/shared/otlp/otlptrace/otlpconfig/tls.go.tmpl "--data={}" --out=tls.go
