// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build aix || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix dragonfly freebsd linux netbsd openbsd solaris zos

package resource // import "go.opentelemetry.io/otel/sdk/resource"

var (
	ParseOSReleaseFile = parseOSReleaseFile
	Skip               = skip
	Parse              = parse
	Unquote            = unquote
	Unescape           = unescape
	BuildOSRelease     = buildOSRelease
)
