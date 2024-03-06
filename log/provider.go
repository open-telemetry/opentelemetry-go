// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/log"

import "go.opentelemetry.io/otel/log/embedded"

// LoggerProvider provides access to [Logger].
//
// Warning: Methods may be added to this interface in minor releases. See
// package documentation on API implementation for information on how to set
// default behavior for unimplemented methods.
type LoggerProvider interface {
	// Users of the interface can ignore this. This embedded type is only used
	// by implementations of this interface. See the "API Implementations"
	// section of the package documentation for more information.
	embedded.LoggerProvider

	// Logger returns a new [Logger] with the provided name and configuration.
	//
	// If name is empty, implementations need to provide a default name.
	//
	// Implementations of this method need to be safe for a user to call
	// concurrently.
	Logger(name string, options ...LoggerOption) Logger
}
