// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type testInstrument struct{}

func (*testInstrument) Enabled(_ context.Context) bool {
	return true
}

func TestEnabledInstrument(t *testing.T) {
	var ei EnabledInstrument = &testInstrument{}

	ctx := t.Context()
	enabled := ei.Enabled(ctx)

	require.True(t, enabled, "Enabled() should return true")
}
