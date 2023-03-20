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

[GetMeterProvider]: https://pkg.go.dev/go.opentelemetry.io/otel#GetMeterProvider
*/
package metric // import "go.opentelemetry.io/otel/metric"
