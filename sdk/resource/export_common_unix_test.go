// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris zos

package resource // import "go.opentelemetry.io/otel/sdk/resource"

var (
	Uname                 = uname
	GetFirstAvailableFile = getFirstAvailableFile
)

var (
	SetUnameProvider        = setUnameProvider
	SetDefaultUnameProvider = setDefaultUnameProvider
)
