// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExporterForceFlush(t *testing.T) {
	ctx := context.Background()
	e, err := New(ctx)
	require.NoError(t, err, "New")

	assert.NoError(t, e.ForceFlush(ctx), "ForceFlush")
}
