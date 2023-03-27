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
Package metric provides the OpenTelemetry API used to measure metrics about
source code operation.

This API is separate from its implementation so the instrumentation built from
it is reusable. See [go.opentelemetry.io/otel/sdk/metric] for the official
OpenTelemetry implementation of this API.

All measurements made with this package are made via instruments. These
instruments are created by a [Meter] which itself is created by a
[MeterProvider]. Applications need to accept a [MeterProvider] implementation
as a starting point when instrumenting. This can be done directly, or by using
the OpenTelemetry global MeterProvider via [GetMeterProvider]. Using an
appropriately named [Meter] from the accepted [MeterProvider], instrumentation
can then be built from the [Meter]'s instruments. See
[go.opentelemetry.io/otel/metric/instrument] for documentation on each
instrument and its intended use.

# API Implementations

This package does not conform to the standard Go versioning policy, all of its
interfaces may have methods added to them without a package major version bump.
So an implementation author is not surprised by this fact, the API is designed
to require a default behavior be defined for unimplemented interface methods.
There are three behavior choices:

  - Compilation failure
  - Panic
  - Default to another implementation

All interfaces in this API embed a corresponding interface from
[go.opentelemetry.io/otel/metric/embedded]. If an author want the default
behavior of their implementations to be a compilation failure, signaling to
their users they need to update to the latest version of that implementation,
they need to embed the corresponding interface from
[go.opentelemetry.io/otel/metric/embedded] in their implementation. For
example,

	import "go.opentelemetry.io/otel/metric/embedded"

	type MeterProvider struct {
		embedded.MeterProvider
		// ...
	}

If an author wants the default behavior of their implementations to a panic,
they need to embed the API interface directly.

	import "go.opentelemetry.io/otel/metric"

	type MeterProvider struct {
		metric.MeterProvider
		// ...
	}

This is not a recommended behavior as it could lead to publishing packages that
contain runtime panics when users update other package that use newer versions
of [go.opentelemetry.io/otel/metric].

Finally, an author can embed another implementation in theirs. The embedded
implementation will be used for methods not defined by the author. For example,
an author who want to default to silently dropping the call can use
[go.opentelemetry.io/otel/metric/noop]:

	import "go.opentelemetry.io/otel/metric/noop"

	type MeterProvider struct {
		noop.MeterProvider
		// ...
	}

[GetMeterProvider]: https://pkg.go.dev/go.opentelemetry.io/otel#GetMeterProvider
*/
package metric // import "go.opentelemetry.io/otel/metric"
