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
}

type RealClock interface {
	Clock
	Since(t lib.Time) lib.Duration
}

type FallbackRealClock struct {
	NowFunc   func() lib.Time
	SinceFunc func(lib.Time) lib.Duration
}

func (c FallbackRealClock) Now() lib.Time {
	return c.NowFunc()
}

func (c FallbackRealClock) Since(t lib.Time) lib.Duration {
	return c.SinceFunc(t)
}

func Default() RealClock {
	return FallbackRealClock{
		NowFunc:   lib.Now,
		SinceFunc: lib.Since,
	}
}

func ConvertClockToRealClock(clk Clock) RealClock {
	if rclk, ok := clk.(RealClock); ok {
		return rclk
	}
	return FallbackRealClock{
		NowFunc: clk.Now,
		SinceFunc: func(t lib.Time) lib.Duration {
			return clk.Now().Sub(t)
		},
	}
}
