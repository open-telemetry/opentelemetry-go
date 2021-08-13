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

type window struct {
	// start is the beginning of this window
	start time.Time

	// low is the lower power-of-two probability
	low trace.Sampler
	// high is the higher power-of-two probability
	high trace.Sampler

	// lowProb is the probability of sampling at the lowProb
	lowProb float64

	// count is updated with atomic.AddInt64
	count int64

	// compute is called when the Sampler's current window
	// expires. the first caller computes a new value and updates
	// the Sampler's current window.
	compute sync.Once
}

// Sampler dynamically adjusts its sampling rate based on the observed
// arrival rate to produce an (expected) rate of sample spans.  This Sampler
// is probabilistic in nature and does not ensure a hard rate limit.
type Sampler struct {
	targetRate float64
	interval   time.Duration
	nowfunc    func() time.Time

	// current is an atomic variable storing the current *window
	// used for estimating the next window's probability.  the
	// first caller to discover a *window after the interval has
	// expired will replace it with a new window.  updates are
	// synchronized via the sync.Once of the expiring window.
	current atomic.Value

	priorCount    int64
	priorDuration time.Duration
}

const (
	DefaultInterval = 10 * time.Second
	MinimumInterval = 10 * time.Millisecond
)

type config struct {
	interval time.Duration
	nowfunc  func() time.Time
}

type Option interface {
	apply(*config)
}

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

// NewSampler returns a Sampler that adjusts its sampling probability
// to achieve an expected rate.
func NewSampler(targetRate float64, opts ...Option) trace.Sampler {
	// Negatigve or zero rate means do not sample.
	if targetRate <= 0 {
		return trace.NeverSample()
	}
	cfg := config{
		interval: DefaultInterval,
		nowfunc:  time.Now,
	}
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	// MinimumInterval avoids bad configurations near the
	// resolution of the runtime scheduler.
	if cfg.interval < MinimumInterval {
		cfg.interval = MinimumInterval
	}

	sampler := &Sampler{
		interval:   cfg.interval,
		nowfunc:    cfg.nowfunc,
		targetRate: targetRate,
	}
	sampler.current.Store(&window{
		start:   sampler.nowfunc(),
		low:     tidSamplerForLogAdjustedCount(0), // starting probability is 1
		high:    nil,
		lowProb: 1,
	})
	return sampler
}

// These are IEEE 754 double-width floating point constants used with
// math.Float64bits.
const (
	offsetExponentMask = 0x7ff0000000000000
	offsetExponentBias = 1023
	significandBits    = 52
)

// expFromFloat64 returns floor(log2(x)).
func expFromFloat64(x float64) int {
	return int((math.Float64bits(x)&offsetExponentMask)>>significandBits) - offsetExponentBias
}

// expToFloat64 returns 2^x.
func expToFloat64(x int) float64 {
	return math.Float64frombits(uint64(offsetExponentBias+x) << significandBits)
}

// splitProb returns the two values of log-adjusted-count nearest to p
// Example:
//
//   splitProb(0.375) => (2, 1, 0.5)
//
// indicates to sample with probability (2^-2) 50% of the time
// and (2^-1) 50% of the time.
func splitProb(p float64) (int, int, float64) {
	// Take the exponent and drop the significand to locate the
	// smaller of two powers of two.
	exp := expFromFloat64(p)

	// Low is the smaller of two log-adjusted counts, the negative
	// of the exponent computed above.
	low := -exp
	// High is the greater of two log-adjusted counts (i.e., one
	// less than low, a smaller adjusted count means a larger
	// probability).
	high := low - 1

	// Return these to probability values and use linear
	// interpolation to compute the required probability of
	// choosing the low-probability Sampler.
	lowP := expToFloat64(-low)
	highP := expToFloat64(-high)
	lowProb := (highP - p) / (highP - lowP)

	return low, high, lowProb
}

// tidSamplerForLogAdjustedCount
func tidSamplerForLogAdjustedCount(logAdjustedCount int) trace.Sampler {
	return trace.TraceIDRatioBased(expToFloat64(-logAdjustedCount))
}

func (s *Sampler) ShouldSample(params trace.SamplingParameters) trace.SamplingResult {
	state := s.current.Load().(*window)
	now := s.nowfunc()

	if now.Sub(state.start) >= s.interval {
		// If the window has expired, update it and re-load.
		state.compute.Do(func() {
			s.updateWindow(state, now)
		})
		state = s.current.Load().(*window)
	}

	// Count the span in this window's rate estimate.
	_ = atomic.AddInt64(&state.count, 1)

	// Compare a uniform random with lowProb, choose either the
	// low or high probability TraceIDRatio Sampler.
	var tid trace.Sampler
	if rand.Float64() < state.lowProb {
		tid = state.low
	} else {
		tid = state.high
	}

	return tid.ShouldSample(params)
}

func (s *Sampler) Description() string {
	return fmt.Sprintf("RateLimited{%g}", s.targetRate)
}

func (s *Sampler) updateWindow(expired *window, now time.Time) {
	// Capture the actual count and the corresponding interval
	// that was measured since the probability was last updated.
	count := atomic.LoadInt64(&expired.count)
	duration := now.Sub(expired.start)

	// Combine the new data and the old data.  In Bayesian terms,
	// this is justified by modelling the arrival of spans as a
	// Poisson process.  The maximum-a-posteriori estimate of the
	// rate based on the observed data equals totalCount divided
	// by totalDuration.
	totalCount := count + s.priorCount
	totalDuration := duration + s.priorDuration
	predictedRate := float64(totalCount) / totalDuration.Seconds()

	// Compute the probability that will yield the target rate.
	probability := s.targetRate / predictedRate

	if probability > 1 {
		probability = 1
	}

	// update the Sampler state, save this window's count and
	// duration for the next window's update.
	lowS, highS, lowProb := splitProb(probability)

	next := &window{
		start:   now,
		low:     tidSamplerForLogAdjustedCount(lowS),
		high:    tidSamplerForLogAdjustedCount(highS),
		lowProb: lowProb,
	}

	s.priorDuration = duration
	s.priorCount = count

	s.current.Store(next)
}
