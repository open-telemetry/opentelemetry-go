// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireErrorString(t *testing.T, expect string, err error) {
	t.Helper()
	require.Error(t, err)
	require.ErrorIs(t, err, PartialSuccess{})

	const pfx = "OTLP partial success: "

	msg := err.Error()
	require.True(t, strings.HasPrefix(msg, pfx))
	require.Equal(t, expect, msg[len(pfx):])
}

func TestPartialSuccessFormat(t *testing.T) {
	requireErrorString(t, "empty message (0 logs rejected)", LogPartialSuccessError(0, ""))
	requireErrorString(t, "help help (0 logs rejected)", LogPartialSuccessError(0, "help help"))
	requireErrorString(
		t,
		"what happened (10 logs rejected)",
		LogPartialSuccessError(10, "what happened"),
	)
	requireErrorString(t, "what happened (15 logs rejected)", LogPartialSuccessError(15, "what happened"))
}
