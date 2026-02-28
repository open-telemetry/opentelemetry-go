// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/sdk/metric/internal/x"

import "context"

// EnabledInstrument informs whether the instrument is enabled.
type EnabledInstrument interface {
	// Enabled reports whether the instrument will process measurements for the given context.
	//
	// This function can be used in places where measuring an instrument
	// would result in computationally expensive operations.
	Enabled(context.Context) bool
}
