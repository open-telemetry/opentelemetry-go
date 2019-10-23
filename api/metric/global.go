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
	"sync/atomic"
)

var global atomic.Value

// GlobalMeter returns a meter registered as a global meter. If no
// meter is registered then an instance of noop Meter is returned.
func GlobalMeter() Meter {
	if t := global.Load(); t != nil {
		return t.(Meter)
	}
	return noopMeter{}
}

// SetGlobalMeter sets provided meter as a global meter.
func SetGlobalMeter(t Meter) {
	global.Store(t)
}
