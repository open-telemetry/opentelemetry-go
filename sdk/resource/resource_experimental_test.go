// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package resource

import (
	"regexp"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

func TestDefaultExperimental(t *testing.T) {
	// Experimental attributes aren't present by default
	res := Default()

	require.False(t, res.Set().HasValue(semconv.ServiceInstanceIDKey))

	// Reset cache and enable experimental resources
	defaultResourceOnce = sync.Once{}
	t.Setenv("OTEL_GO_X_RESOURCE", "true")

	res = Default()

	require.True(t, res.Set().HasValue(semconv.ServiceInstanceIDKey))

	serviceInstanceID, ok := res.Set().Value(semconv.ServiceInstanceIDKey)
	require.True(t, ok)
	matched, err := regexp.MatchString(
		"^[A-Fa-f0-9]{8}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{4}-[A-Fa-f0-9]{12}$",
		serviceInstanceID.AsString(),
	)
	require.NoError(t, err)
	require.True(t, matched)
}
