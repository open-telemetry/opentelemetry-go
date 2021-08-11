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

package ratelimit // import "go.opentelemetry.io/otel/sdk/trace/ratelimit"

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
)

type config struct {
	interval time.Duration
	nowfunc  func() time.Time
}

type Option interface {
	apply(*config)
}

type window struct {
	// start is the beginning of this window
	start time.Time

	// lowS is the lower power-of-two probability
	lowS int
	// lowS is the higher power-of-two probability
	highS int

	// lowProb is the probability of sampling at the lowProb
	lowProb float64

	count   int64
	compute sync.Once
}

type Sampler struct {
	targetCount int
	interval    time.Duration

	current atomic.Value

	priorCount    int64
	priorDuration time.Duration
}

const (
	DefaultInterval = 10 * time.Second
	MinInterval     = 10 * time.Millisecond
	MinRate         = 0.00001
)

type intervalOption time.Duration
type nowfuncOption func() time.Time

func WithInterval(d time.Duration) Option {
	return intervalOption(d)
}

func WithNowFunc(f func() time.Time) Option {
	return nowfuncOption(f)
}

func (i intervalOption) apply(cfg *config) {
	cfg.interval = time.Duration(i)
}

func (n nowfuncOption) apply(cfg *config) {
	cfg.nowfunc = n
}

var _ trace.Sampler = &Sampler{}

func NewSampler(maxRate float64, opts ...Option) *Sampler {
	cfg := config{
		interval: DefaultInterval,
		nowfunc:  time.Now,
	}
	for _, opt := range opts {
		opt.apply(&cfg)
	}

	if cfg.interval < MinInterval {
		cfg.interval = MinInterval
	}

	if maxRate < MinRate {
		maxRate = MinRate
	}

	target := int(cfg.interval.Seconds() / maxRate)

	if target < 1 {
		target = 1
	}

	sampler := &Sampler{
		interval:    cfg.interval,
		targetCount: target,
	}
	sampler.current.Store(&window{
		start:   time.Now(),
		lowS:    0,
		highS:   -1,
		lowProb: 1,
	})
	return sampler
}

const (
	offsetExponentMask = 0x7ff0000000000000
	offsetExponentBias = 1023
	significandBits    = 52
)

func expFromFloat64(p float64) int {
	return int((math.Float64bits(p)&offsetExponentMask)>>significandBits) - offsetExponentBias
}

func expToFloat64(e int) float64 {
	return math.Float64frombits(uint64(offsetExponentBias+e) << significandBits)
}

func splitProb(p float64) (int, int, float64) {
	// Return the two values of log-adjusted-count nearest to p
	// Example:
	//   splitProb(0.75) returns (1, 0, 0.5)
	// meaning to sample with probability (2^-1) 50% of the time
	// and 2^0 50% of the time.
	exp := expFromFloat64(p)

	low := -exp
	high := low - 1

	lowP := expToFloat64(-low)
	highP := expToFloat64(-high)

	return low, high, (highP - p) / (highP - lowP)
}

func (s *Sampler) ShouldSample(params trace.SamplingParameters) trace.SamplingResult {
	state := s.current.Load().(*window)
	now := time.Now()

	_ = atomic.AddInt64(&state.count, 1)

	if now.Sub(state.start) >= s.interval {
		state.compute.Do(func() {
			s.updateWindow(state, now)
		})
	}

	var S int
	if rand.Float64() < state.lowProb {
		S = state.lowS
	} else {
		S = state.highS
	}

	// Note: this prototype demonstrates the intended result, not
	// an efficient implementation. Just do whatever
	// TraceIDRatioBased() would do at this point.  Make a
	// consistent sampling decision using S according to OTEP 168.
	return trace.TraceIDRatioBased(1.0 / float64(int64(1)<<S)).ShouldSample(params)
}

func (s *Sampler) Description() string {
	return fmt.Sprintf("RateLimited{%g}", float64(s.targetCount)/s.interval.Seconds())
}

func (s *Sampler) updateWindow(expired *window, now time.Time) {
	count := atomic.LoadInt64(&expired.count)
	duration := now.Sub(expired.start)

	totalCount := count + s.priorCount
	totalDuration := duration + s.priorDuration

	countFactor := float64(totalCount) / float64(s.targetCount)
	durationFactor := float64(totalDuration) / float64(s.interval)
	probability := durationFactor / countFactor

	if probability > 1 {
		probability = 1
	}

	lowS, highS, lowProb := splitProb(probability)

	next := &window{
		start:   now,
		lowS:    lowS,
		highS:   highS,
		lowProb: lowProb,
	}

	s.current.Store(next)
}
