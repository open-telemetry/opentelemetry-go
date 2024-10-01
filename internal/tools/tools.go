// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build tools
// +build tools

package tools // import "go.opentelemetry.io/otel/internal/tools"

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/gogo/protobuf/protoc-gen-gogofast"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/jcchavezs/porto/cmd/porto"
	_ "github.com/wadey/gocovmerge"
	_ "go.opentelemetry.io/build-tools/crosslink"
	_ "go.opentelemetry.io/build-tools/gotmpl"
	_ "go.opentelemetry.io/build-tools/multimod"
	_ "go.opentelemetry.io/build-tools/semconvgen"
	_ "golang.org/x/exp/cmd/gorelease"
	_ "golang.org/x/tools/cmd/stringer"
	_ "golang.org/x/vuln/cmd/govulncheck"
)
