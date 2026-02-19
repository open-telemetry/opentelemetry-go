// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package errorhandler // import "go.opentelemetry.io/otel/internal/errorhandler"

import (
	"sync"
	"testing"
)

// ResetForTest configures the test to restore the initial global error handler
// state during its Cleanup step.
func ResetForTest(t testing.TB) {
	t.Cleanup(func() {
		globalErrorHandler = defaultErrorHandler()
		delegateErrorHandlerOnce = sync.Once{}
	})
}
