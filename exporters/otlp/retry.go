package otlp

import (
	"github.com/cenkalti/backoff/v4"
	"time"
)

// RetrySettings defines configuration for retrying batches in case of export failure.
// The current supported strategy is exponential backoff.
type RetrySettings struct {
	// Enabled indicates whether to not retry sending batches in case of export failure.
	Enabled bool `mapstructure:"enabled"`
	// InitialInterval the time to wait after the first failure before retrying.
	InitialInterval time.Duration `mapstructure:"initial_interval"`
	// MaxInterval is the upper bound on backoff interval. Once this value is reached the delay between
	// consecutive retries will always be `MaxInterval`.
	MaxInterval time.Duration `mapstructure:"max_interval"`
	// MaxElapsedTime is the maximum amount of time (including retries) spent trying to send a request/batch.
	// Once this value is reached, the data is discarded.
	MaxElapsedTime time.Duration `mapstructure:"max_elapsed_time"`
}

// DefaultRetrySettings returns the default settings for RetrySettings.
func DefaultRetrySettings() RetrySettings {
	return RetrySettings{
		Enabled:         true,
		InitialInterval: 5 * time.Second,
		MaxInterval:     30 * time.Second,
		MaxElapsedTime:  5 * time.Minute,
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
