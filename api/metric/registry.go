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

// Copyright 2018, OpenCensus Authors
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
	"errors"
	"sync"
)

type (
	// Registry is a mechanism for avoiding duplicate registration
	// of different-type pre-aggregated metrics (in one process).
	Registry interface {
		RegisterMetric(Metric) (Metric, error)
		ForeachMetric(func(string, Metric))
	}

	registry struct {
		nameType sync.Map // map[string]Metric
	}
)

var (
	registryLock   sync.Mutex
	registryGlobal Registry = &registry{}

	errDuplicateMetricTypeConflict = errors.New("Duplicate metric registration with conflicting type")
)

// SetRegistry may be used to reset the global metric registry, which should not be
// needed unless for testing purposes.
func SetRegistry(r Registry) {
	registryLock.Lock()
	defer registryLock.Unlock()
	registryGlobal = r
}

// GetRegistry may be used to access a global list of metric definitions.
func GetRegistry() Registry {
	registryLock.Lock()
	defer registryLock.Unlock()
	return registryGlobal
}

func (r *registry) RegisterMetric(newMet Metric) (Metric, error) {
	name := newMet.Measure().Name()
	has, ok := r.nameType.Load(name)

	if ok {
		m := has.(Metric)
		if m.Type() != newMet.Type() {
			return nil, errDuplicateMetricTypeConflict
		}
		return m, nil
	}

	r.nameType.Store(name, newMet)
	return newMet, nil
}

func (r *registry) ForeachMetric(f func(string, Metric)) {
	r.nameType.Range(func(key, value interface{}) bool {
		f(key.(string), value.(Metric))
		return true
	})
}
