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

package trace

import (
	"time"
)

// Clock is the entrypoint for providing time to span's start/end timestamp.
// By default the standard "time" package will be used. User can replace
// it with customized clock implementation (e.g. has additional clock
// synchronization logic) by using the `WithClock` option.
type Clock interface {
	// Stopwatch returns a started Stopwatch measuring a time interval.
	Stopwatch() Stopwatch
}

type Stopwatch interface {
	// Started returns the time the Stopwatch was started
	Started() time.Time
	// Elapsed returns the duration from when this Stopwatch was started and now using a consistent monotonic time.
	Elapsed() time.Duration
}

func defaultClock() Clock {
	return standardClock{}
}

func defaultStopwatch(t time.Time) Stopwatch {
	return standardStopwatch{
		startTime: t,
	}
}

func (standardClock) Stopwatch() Stopwatch {
	return standardStopwatch{
		startTime: time.Now(),
	}
}

type standardClock struct{}
type standardStopwatch struct {
	startTime time.Time
}

func (w standardStopwatch) Started() time.Time {
	return w.startTime
}

func (w standardStopwatch) Elapsed() time.Duration {
	return time.Since(w.startTime)
}
