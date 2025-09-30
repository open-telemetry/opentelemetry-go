package x

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type testInstrument struct{}

func (t *testInstrument) Enabled(ctx context.Context) bool {
	return true
}

func TestEnabledInstrument(t *testing.T) {
	var ei EnabledInstrument = &testInstrument{}

	ctx := context.Background()
	enabled := ei.Enabled(ctx)

	require.True(t, enabled, "Enabled() should return true")
}
