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

// Clock is responsible for providing time for span's start/end timestamp.
// By default the standard "time" package will be used. User can replace
// it with customized clock implementation (e.g. has additional clock
// synchronization logic) by using the `WithClock` option.
type Clock interface {
	Now() time.Time
}

type clockWrapper struct {
	nowFunc   func() time.Time
	sinceFunc func(time.Time) time.Duration
}

func (c clockWrapper) Now() time.Time {
	return c.nowFunc()
}

func (c clockWrapper) Since(t time.Time) time.Duration {
	return c.sinceFunc(t)
}

func defaultClock() *clockWrapper {
	return &clockWrapper{
		nowFunc:   time.Now,
		sinceFunc: time.Since,
	}
}

// monotonicEndTime returns the end time at present
// but offset from start, monotonically.
//
// The monotonic clock is used in subtractions hence
// the duration since start added back to start gives
// end as a monotonic time.
// See https://golang.org/pkg/time/#hdr-Monotonic_Clocks
func monotonicEndTime(start time.Time, clock *clockWrapper) time.Time {
	return start.Add(clock.Since(start))
}
