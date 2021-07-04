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

package time

import (
	lib "time"
)

type Clock interface {
	Now() lib.Time
	Since(t lib.Time) lib.Duration
}

type DefaultClock struct {
}

func (c DefaultClock) Now() lib.Time {
	return lib.Now()
}
func (c DefaultClock) Since(t lib.Time) lib.Duration {
	return lib.Since(t)
}

var _ Clock = DefaultClock{}
