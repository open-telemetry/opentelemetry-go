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

import (
	"fmt"
	"sync/atomic"
)

var (
	descriptorID uint64
)

func registerDescriptor(name string, kind Kind, valueKind MetricValueKind, opts []Option, d *Descriptor) {
	d.Name = name
	d.Kind = kind
	d.ValueKind = valueKind
	d.ID = DescriptorID(atomic.AddUint64(&descriptorID, 1))

	for _, opt := range opts {
		opt(d)
	}
	ensureValidDescriptor(d)
}

func ensureValidDescriptor(d *Descriptor) {
	checkNonMonotonic := false
	checkMonotonic := false
	checkSigned := false
	switch d.Kind {
	case Invalid:
		panic("tried to register a metric descriptor with invalid kind")
	case CounterKind:
		checkMonotonic, checkSigned = true, true
	case GaugeKind:
		checkNonMonotonic, checkSigned = true, true
	case MeasureKind:
		checkMonotonic, checkNonMonotonic = true, true
	}
	if checkNonMonotonic && d.NonMonotonic {
		panicBadField(d.Kind, "NonMonotonic")
	}
	if checkMonotonic && d.Monotonic {
		panicBadField(d.Kind, "Monotonic")
	}
	if checkSigned && d.Signed {
		panicBadField(d.Kind, "Signed")
	}
}

func panicBadField(kind Kind, field string) {
	panic(fmt.Sprintf("invalid %s descriptor, has set %s field", kind.String(), field))
}
