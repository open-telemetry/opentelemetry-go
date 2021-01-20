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

package sampling_test // import "go.opentelemetry.io/otel/sdk/trace/sampling_test"

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/sdk/trace/sampling"
	"go.opentelemetry.io/otel/trace"
)

var _ sampling.Sampler = (*sampling.RateLimitingSampler)(nil)

const letters = "abcdef0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}
	return string(b)
}

func TestNewNewRateLimitingSampler(t *testing.T) {
	testSamplers := struct {
		samplerMax0                        *sampling.RateLimitingSampler
		samplerMax0ExpectedIsSampledCount  int
		samplerMax05                       *sampling.RateLimitingSampler
		samplerMax05ExpectedIsSampledCount int
		samplerMax1                        *sampling.RateLimitingSampler
		samplerMax1ExpectedIsSampledCount  int
	}{
		samplerMax0:                        sampling.NewRateLimitingSampler(0),
		samplerMax0ExpectedIsSampledCount:  0,
		samplerMax05:                       sampling.NewRateLimitingSampler(0.5),
		samplerMax05ExpectedIsSampledCount: 3,
		samplerMax1:                        sampling.NewRateLimitingSampler(1),
		samplerMax1ExpectedIsSampledCount:  5,
	}

	samplerMax0ActualIsSampledCount := 0
	samplerMax05ActualIsSampledCount := 0
	samplerMax1ActualIsSampledCount := 0

	t.Run("test new rate limiting samplers", func(t *testing.T) {
		tick := time.NewTicker(time.Millisecond * 500)
		ctx, cancelFunc := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for {
				select {
				case <-tick.C:
					id := randStringBytes(32)
					traceID, err := trace.TraceIDFromHex(id)
					assert.NoError(t, err)

					if isDecisionTrue, _ := testSamplers.samplerMax0.IsSampled(traceID); isDecisionTrue {
						samplerMax0ActualIsSampledCount++
					}
					if isDecisionTrue, _ := testSamplers.samplerMax05.IsSampled(traceID); isDecisionTrue {
						samplerMax05ActualIsSampledCount++
					}
					if isDecisionTrue, _ := testSamplers.samplerMax1.IsSampled(traceID); isDecisionTrue {
						samplerMax1ActualIsSampledCount++
					}

				case <-ctx.Done():
					tick.Stop()
					wg.Done()
					return
				}
			}
		}()
		time.Sleep(5 * time.Second)
		cancelFunc()
		wg.Wait()
		assert.Equal(t, testSamplers.samplerMax0ExpectedIsSampledCount, samplerMax0ActualIsSampledCount)
		assert.True(t, testSamplers.samplerMax05ExpectedIsSampledCount >= samplerMax05ActualIsSampledCount)
		assert.True(t, samplerMax05ActualIsSampledCount >= 0)
		assert.True(t, testSamplers.samplerMax1ExpectedIsSampledCount >= samplerMax1ActualIsSampledCount)
		assert.True(t, samplerMax1ActualIsSampledCount >= 0)
	})
}
