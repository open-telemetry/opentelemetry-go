// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func basicTracerProvider(t *testing.T) *TracerProvider {
	tp := NewTracerProvider(WithSampler(AlwaysSample()))
	t.Cleanup(func() {
		assert.NoError(t, tp.Shutdown(context.Background()))
	})
	return tp
}
