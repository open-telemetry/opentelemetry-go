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

package otlp

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

// RetrySettings defines configuration for retrying batches in case of export failure.
// The current supported strategy is exponential backoff.
type RetrySettings struct {
	// Enabled indicates whether to not retry sending batches in case of export failure.
	Enabled bool
	// InitialInterval the time to wait after the first failure before retrying.
	InitialInterval time.Duration
	// MaxInterval is the upper bound on backoff interval. Once this value is reached the delay between
	// consecutive retries will always be `MaxInterval`.
	MaxInterval time.Duration
	// MaxElapsedTime is the maximum amount of time (including retries) spent trying to send a request/batch.
	// Once this value is reached, the data is discarded.
	MaxElapsedTime time.Duration
}

// DefaultRetrySettings returns the default settings for RetrySettings.
func DefaultRetrySettings() RetrySettings {
	return RetrySettings{
		Enabled:         true,
		InitialInterval: 5 * time.Second,
		MaxInterval:     30 * time.Second,
		MaxElapsedTime:  time.Minute,
	}
}

func NewExponentialConfig(rs RetrySettings) *backoff.ExponentialBackOff {
	// Do not use NewExponentialBackOff since it calls Reset and the code here must
	// call Reset after changing the InitialInterval (this saves an unnecessary call to Now).
	expBackoff := &backoff.ExponentialBackOff{
		InitialInterval:     rs.InitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         rs.MaxInterval,
		MaxElapsedTime:      rs.MaxElapsedTime,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	expBackoff.Reset()

	return expBackoff
}
