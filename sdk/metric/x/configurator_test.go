// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type meterConfiguratorOptionExtractor interface {
	MeterConfigurator() func(instrumentation.Scope) any
}

func TestWithMeterConfiguratorContract(t *testing.T) {
	opt := WithMeterConfigurator(func(s instrumentation.Scope) MeterConfig {
		return MeterConfig{}
	})

	type experimental interface{ Experimental() }
	_, ok := opt.(experimental)
	require.True(t, ok, "must implement Experimental()")

	_, ok = opt.(meterConfiguratorOptionExtractor)
	require.True(t, ok, "must implement MeterConfigurator() func(scope) any")
}

func TestWithMeterConfiguratorBehavior(t *testing.T) {
	opt := WithMeterConfigurator(func(s instrumentation.Scope) MeterConfig {
		return NewMeterConfig(WithMeterEnabled(s.Name != "disabled"))
	})

	ex := opt.(meterConfiguratorOptionExtractor)

	for _, tc := range []struct {
		name    string
		scope   instrumentation.Scope
		enabled bool
	}{
		{"scope/enabled", instrumentation.Scope{Name: "test"}, true},
		{"scope/disabled", instrumentation.Scope{Name: "disabled"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := ex.MeterConfigurator()(tc.scope)
			cfg, ok := result.(interface{ Enabled() bool })
			require.True(t, ok, "result must implement Enabled() bool")
			assert.Equal(t, tc.enabled, cfg.Enabled())
		})
	}
}
