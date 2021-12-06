// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otlpmetrichttp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnreasonableBackoff(t *testing.T) {
	cIface := NewClient(
		WithEndpoint("http://localhost"),
		WithInsecure(),
		WithBackoff(-time.Microsecond),
	)
	require.IsType(t, &client{}, cIface)
	c := cIface.(*client)
	assert.True(t, c.generalCfg.RetryConfig.Enabled)
	assert.Equal(t, 5*time.Second, c.generalCfg.RetryConfig.InitialInterval)
	assert.Equal(t, 300*time.Millisecond, c.generalCfg.RetryConfig.MaxInterval)
	assert.Equal(t, time.Minute, c.generalCfg.RetryConfig.MaxElapsedTime)
}

func TestUnreasonableMaxAttempts(t *testing.T) {
	type testcase struct {
		name        string
		maxAttempts int
	}
	for _, tc := range []testcase{
		{
			name:        "negative max attempts",
			maxAttempts: -3,
		},
		{
			name:        "too large max attempts",
			maxAttempts: 10,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cIface := NewClient(
				WithEndpoint("http://localhost"),
				WithInsecure(),
				WithMaxAttempts(tc.maxAttempts),
			)
			require.IsType(t, &client{}, cIface)
			c := cIface.(*client)
			assert.True(t, c.generalCfg.RetryConfig.Enabled)
			assert.Equal(t, 5*time.Second, c.generalCfg.RetryConfig.InitialInterval)
			assert.Equal(t, 30*time.Second, c.generalCfg.RetryConfig.MaxInterval)
			assert.Equal(t, 145*time.Second, c.generalCfg.RetryConfig.MaxElapsedTime)
		})
	}
}
