// Copyright 2019, OpenTelemetry Authors
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

package metric

func WrapInt64CounterInstrument(instrument Instrument) Int64Counter {
	return Int64Counter{commonMetric: newCommonMetric(instrument)}
}

func WrapFloat64CounterInstrument(instrument Instrument) Float64Counter {
	return Float64Counter{commonMetric: newCommonMetric(instrument)}
}

func WrapInt64GaugeInstrument(instrument Instrument) Int64Gauge {
	return Int64Gauge{commonMetric: newCommonMetric(instrument)}
}

func WrapFloat64GaugeInstrument(instrument Instrument) Float64Gauge {
	return Float64Gauge{commonMetric: newCommonMetric(instrument)}
}

func WrapInt64MeasureInstrument(instrument Instrument) Int64Measure {
	return Int64Measure{commonMetric: newCommonMetric(instrument)}
}

func WrapFloat64MeasureInstrument(instrument Instrument) Float64Measure {
	return Float64Measure{commonMetric: newCommonMetric(instrument)}
}

func ApplyCounterOptions(opts *Options, cos ...CounterOptionApplier) {
	for _, o := range cos {
		o.ApplyCounterOption(opts)
	}
}

func ApplyGaugeOptions(opts *Options, gos ...GaugeOptionApplier) {
	for _, o := range gos {
		o.ApplyGaugeOption(opts)
	}
}

func ApplyMeasureOptions(opts *Options, mos ...MeasureOptionApplier) {
	for _, o := range mos {
		o.ApplyMeasureOption(opts)
	}
}
