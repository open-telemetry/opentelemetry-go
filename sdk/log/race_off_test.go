// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !race
// +build !race

package log

func race() bool {
	return false
}
