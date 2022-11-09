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

package metric // import "go.opentelemetry.io/otel/metric"

// Meter creates instrument instances for an instrumentation library.
//
// Warning: methods may be added to this interface in minor releases.
type Meter interface {
	Float64Counter(name string, opts ...InstrumentOption) (Float64Counter, error)
	Float64UpDownCounter(name string, opts ...InstrumentOption) (Float64UpDownCounter, error)
	Float64Histogram(name string, opts ...InstrumentOption) (Float64Histogram, error)
	Float64ObservableCounter(name string, opts ...ObservableOption) (Float64ObservableCounter, error)
	Float64ObservableUpDownCounter(name string, opts ...ObservableOption) (Float64ObservableUpDownCounter, error)
	Float64ObservableGauge(name string, opts ...ObservableOption) (Float64ObservableGauge, error)

	Int64Counter(name string, opts ...InstrumentOption) (Int64Counter, error)
	Int64UpDownCounter(name string, opts ...InstrumentOption) (Int64UpDownCounter, error)
	Int64Histogram(name string, opts ...InstrumentOption) (Int64Histogram, error)
	Int64ObservableCounter(name string, opts ...ObservableOption) (Int64ObservableCounter, error)
	Int64ObservableUpDownCounter(name string, opts ...ObservableOption) (Int64ObservableUpDownCounter, error)
	Int64ObservableGauge(name string, opts ...ObservableOption) (Int64ObservableGauge, error)

	// RegisterCallback registers f to be called when instrument is collected.
	// The callback can also be registered for additional Observable
	// instruments that it will updated when executed.
	//
	// When the Unregister method of the returned Unregisterer is called f is
	// unregistered from the Meter for instrument(s). Calling Unregister after
	// the first time will have no effect.
	RegisterCallback(f Callback, instrument Observable, additional ...Observable) (Unregisterer, error)
}

// Unregisterer undoes the process of registering.
//
// The behavior of Unregister after the first call is undefined. Specific
// implementations may document their own behavior.
type Unregisterer interface {
	Unregister() error
}
