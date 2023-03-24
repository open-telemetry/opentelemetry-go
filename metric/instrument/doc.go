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
Package instrument provides the OpenTelemetry API instruments used to make
measurements.

Each instrument is designed to make measurements of a particular type. Broadly,
all instruments fall into two overlapping logical categories: asynchronous or
synchronous, and int64 or float64.

All synchronous instruments ([Int64Counter], [Int64UpDownCounter],
[Int64Histogram], [Float64Counter], [Float64UpDownCounter], [Float64Histogram])
are used to measure the operation and performance of source code during the
source code execution. These instruments only make measurements when the source
code they instrument is run.

All asynchronous instruments ([Int64ObservableCounter],
[Int64ObservableUpDownCounter], [Int64ObservableGauge],
[Float64ObservableCounter], [Float64ObservableUpDownCounter],
[Float64ObservableGauge]) are used to measure metrics outside of the execution
of source code. They are said to make "observations" via a callback function
called once every measurement collection cycle.

Each instrument is also grouped by the value type it measures. Either int64 or
float64. The value being measured will dictate which instrument in these
categories to use.

Outside of these two broad categories, instruments are described by the
function they are designed to serve. All Counters ([Int64Counter],
[Float64Counter], [Int64ObservableCounter], [Float64ObservableCounter]) are
designed to measure values that never decrease in value, but instead only
incrementally increase in value. UpDownCounters ([Int64UpDownCounter],
[Float64UpDownCounter], [Int64ObservableUpDownCounter],
[Float64ObservableUpDownCounter]) on the other hand, are designed to measure
values that can increase and decrease. When more information
needs to be conveyed about all the synchronous measurements made during a
collection cycle, a Histogram ([Int64Histogram], [Float64Histogram]) should be
used. Finally, when just the most recent measurement needs to be conveyed about an
asynchronous measurement, a Gauge ([Int64ObservableGauge],
[Float64ObservableGauge]) should be used.

See the [OpenTelemetry documentation] for more information about instruments
and their intended use.

[OpenTelemetry documentation]: https://opentelemetry.io/docs/concepts/signals/metrics/
*/
package instrument // import "go.opentelemetry.io/otel/metric/instrument"
