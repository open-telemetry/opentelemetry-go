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

/*
Package log provides an implementation of the logging part of the
OpenTelemetry API.

To participate in distributed traces a LogRecord needs to be created for the
operation being performed as part of a traced workflow. In its simplest form:

	var logger log.Logger

	func init() {
		logger = otel.Logger("instrumentation/package/name")
	}

	func operation(ctx context.Context) {
		logger.Emit(ctx, log.WithAttributes(...))
		// ...
	}

A Logger is unique to the instrumentation and is used to create LogRecords.
Instrumentation should be designed to accept a LoggerProvider from which it
can create its own unique Logger. Alternatively, the registered global
LoggerProvider from the go.opentelemetry.io/otel package can be used as
a default.

	const (
		name    = "instrumentation/package/name"
		version = "0.1.0"
	)

	type Instrumentation struct {
		logger log.Logger
	}

	func NewInstrumentation(tp log.LoggerProvider) *Instrumentation {
		if tp == nil {
			tp = otel.LoggerProvider()
		}
		return &Instrumentation{
			logger: tp.Logger(name, log.WithInstrumentationVersion(version)),
		}
	}

	func operation(ctx context.Context, inst *Instrumentation) {
		inst.tracer.Emit(ctx, log.WithAttributes(...))
		// ...
	}
*/
package log // import "go.opentelemetry.io/otel/log"
