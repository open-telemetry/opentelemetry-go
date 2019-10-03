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

type Observer struct {
	Descriptor
}

type Float64Observer struct {
	Observer
}

type Int64Observer struct {
	Observer
}

func NewObserver(name string, valueKind MetricValueKind, mos ...Option) (o Observer) {
	registerDescriptor(name, ObserverKind, valueKind, mos, &o.Descriptor)
	return
}

func NewFloat64Observer(name string, mos ...Option) (o Float64Observer) {
	o.Observer = NewObserver(name, Float64ValueKind, mos...)
	return
}

func NewInt64Observer(name string, mos ...Option) (o Int64Observer) {
	o.Observer = NewObserver(name, Int64ValueKind, mos...)
	return
}
