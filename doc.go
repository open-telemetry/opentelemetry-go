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

/*Package otel contains the implementation of the OpenTelemetry API.

The OpenTelemetry API is used to instrument your code to participate in
distributed traces and capture measurements of the state of your system in the
form of metric events.

TODO (MrAlias): docs on how to instrument for tracing.

Metric instruments to the way you can instrument your code to start reporting
measurements. To understand what instrument to use when instrumenting your
code it is important to understand the semantic meaning of these instruments.

Instruments are categorized as by being synchronous or asynchronous, and if
they are additive or non-additive. Synchronous instruments are called by the
user with a Context. Asynchronous instruments are called by the SDK during
collection. Additive instruments are semantically intended for capturing a
sum. Non-additive instruments are intended for capturing a distribution.

In addition to this categorization of instruments there is a further
refinement of additive instruments, they may be monotonic or not. If they are
monotonic they are non-decreasing and naturally define a rate.

The synchronous instrument names are:

  Counter:           additive, monotonic
  UpDownCounter:     additive
  ValueRecorder:     non-additive

and the asynchronous instruments are:

  SumObserver:       additive, monotonic
  UpDownSumObserver: additive
  ValueObserver:     non-additive

All instruments are provided with support for either float64 or int64 input
values.

The Meter interface supports allocating new instruments as well as interfaces
for recording batches of synchronous measurements or asynchronous
observations.  To obtain a Meter, use a MeterProvider.

The MeterProvider interface supports obtaining a named Meter interface.  To
obtain a MeterProvider implementation, initialize and configure any compatible
SDK.

*/
package otel // import "go.opentelemetry.io/otel"
