// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package profile

import (
	"context"

	"go.opentelemetry.io/otel/sdk/profile/profiledata"
)

// Producer produces profile data from a profiling source.
type Producer interface {
	// Produce collects and returns profile data for the current collection
	// cycle.
	//
	// Implementations must be safe to call concurrently and should honor
	// cancellation on ctx.
	Produce(ctx context.Context) ([]profiledata.ScopeProfiles, error)
}

type producerShutdowner interface {
	Shutdown(ctx context.Context) error
}

func shutdownProducers(ctx context.Context, producers []Producer) error {
	var err error
	for _, p := range producers {
		if s, ok := p.(producerShutdowner); ok {
			if e := s.Shutdown(ctx); e != nil {
				err = joinErrors(err, e)
			}
		}
	}
	return err
}
