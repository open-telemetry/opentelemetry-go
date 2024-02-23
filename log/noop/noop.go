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

// Package noop provides an implementation of the [OpenTelemetry Logs Bridge
// API] that produces no telemetry and minimizes used computation resources.
//
// Using this package to implement the [OpenTelemetry Logs Bridge API] will
// effectively disable OpenTelemetry.
//
// This implementation can be embedded in other implementations of the
// [OpenTelemetry Logs Bridge API]. Doing so will mean the implementation
// defaults to no operation for methods it does not implement.
//
// [OpenTelemetry Logs Bridge API]: https://pkg.go.dev/go.opentelemetry.io/otel/log
package noop // import "go.opentelemetry.io/otel/log/noop"

import (
	"context"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/embedded"
)

var (
	// Compile-time check this implements the OpenTelemetry API.
	_ log.LoggerProvider = LoggerProvider{}
	_ log.Logger         = Logger{}
)

// LoggerProvider is an OpenTelemetry No-Op LoggerProvider.
type LoggerProvider struct{ embedded.LoggerProvider }

// NewLoggerProvider returns a LoggerProvider that does not record any telemetry.
func NewLoggerProvider() LoggerProvider {
	return LoggerProvider{}
}

// Logger returns an OpenTelemetry Logger that does not record any telemetry.
func (LoggerProvider) Logger(string, ...log.LoggerOption) log.Logger {
	return Logger{}
}

// Logger is an OpenTelemetry No-Op Logger.
type Logger struct{ embedded.Logger }

// Emit does nothing.
func (Logger) Emit(context.Context, log.Record) {}
