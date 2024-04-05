// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package otlploghttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	testcases := []struct {
		name    string
		options []Option
		envars  map[string]string
		want    config
	}{
		{
			name: "Defaults",
			want: config{
				endpoint:    newSetting(defaultEndpoint),
				path:        newSetting(defaultPath),
				insecure:    newSetting(defaultInsecure),
				tlsCfg:      newSetting(defaultTlsCfg),
				headers:     newSetting(defaultHeaders),
				compression: newSetting(defaultCompression),
				timeout:     newSetting(defaultTimeout),
				proxy:       newSetting(defaultProxy),
				retryCfg:    newSetting(defaultRetryCfg),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, newConfig(tc.options))
		})
	}
}
