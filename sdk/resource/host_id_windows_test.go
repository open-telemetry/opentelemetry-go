// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows
// +build windows

package resource // import "go.opentelemetry.io/otel/sdk/resource"

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	reader := &hostIDReaderWindows{}
	result, err := reader.read()

	require.NoError(t, err)
	require.NotEmpty(t, result)
}
