// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package internal contains utility functions and internal implementations
package internal // import "go.opentelemetry.io/otel/internal"

//go:generate gotmpl --body=./shared/internaltest/doc.go.tmpl "--data={}" --out=internaltest/doc.go
//go:generate gotmpl --body=./shared/internaltest/text_map_carrier.go.tmpl "--data={}" --out=internaltest/text_map_carrier.go
//go:generate gotmpl --body=./shared/internaltest/text_map_carrier_test.go.tmpl "--data={}" --out=internaltest/text_map_carrier_test.go
//go:generate gotmpl --body=./shared/internaltest/text_map_propagator.go.tmpl "--data={}" --out=internaltest/text_map_propagator.go
//go:generate gotmpl --body=./shared/internaltest/text_map_propagator_test.go.tmpl "--data={}" --out=internaltest/text_map_propagator_test.go
