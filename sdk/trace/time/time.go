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

type SimpleClock interface {
	Now() lib.Time
}

type Clock interface {
	SimpleClock
	Since(t lib.Time) lib.Duration
}

type ClockWrapper struct {
	NowFunc   func() lib.Time
	SinceFunc func(lib.Time) lib.Duration
}

func (c ClockWrapper) Now() lib.Time {
	return c.NowFunc()
}

func (c ClockWrapper) Since(t lib.Time) lib.Duration {
	return c.SinceFunc(t)
}

func Default() Clock {
	return ClockWrapper{
		NowFunc:   lib.Now,
		SinceFunc: lib.Since,
	}
}

func ConvertFromSimpleClock(clk SimpleClock) Clock {
	if rclk, ok := clk.(Clock); ok {
		return rclk
	}
	return ClockWrapper{
		NowFunc: clk.Now,
		SinceFunc: func(t lib.Time) lib.Duration {
			return clk.Now().Sub(t)
		},
	}
}
