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

package sampling // import "go.opentelemetry.io/otel/sdk/trace/sampling"

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.opentelemetry.io/otel/label"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	strategyTypeProbabilistic = 0
	strategyTypeRateLimiting  = 1
)

// Sampler is an interface for ProbabilisticSampler & RateLimitingSampler.
type Sampler interface {
	IsSampled(id trace.TraceID) (bool, []label.KeyValue)
}

// RemotelyControlledSampler is an otelTrace.Sampler implementation.
type RemotelyControlledSampler struct {
	once sync.Once // use to close doneChan

	sync.RWMutex // used to serialize access to perOperationSamplers

	waitGroup *sync.WaitGroup

	serviceName             string
	doneChan                chan struct{}
	defaultSampler          Sampler
	perOperationSamplers    map[string]Sampler
	remoteStrategyEndpoint  string
	samplingRefreshInterval time.Duration
	maxPerOperationSamplers int
	logger                  *log.Logger
}

// NewRemotelyControlledSampler creates new remotely controlled sampler with given Strategy server endpoint & appName.
func NewRemotelyControlledSampler(
	config Config,
	logger *log.Logger,
) (
	*RemotelyControlledSampler,
	error,
) {
	sampler := &RemotelyControlledSampler{
		waitGroup:               &sync.WaitGroup{},
		serviceName:             config.AppName,
		doneChan:                make(chan struct{}),
		defaultSampler:          NewProbabilisticSampler(config.Rate),
		perOperationSamplers:    make(map[string]Sampler),
		remoteStrategyEndpoint:  config.SamplerStrategyEndpoint,
		samplingRefreshInterval: config.RefreshInterval,
		maxPerOperationSamplers: config.MaxOperations,
		logger:                  logger,
	}

	sampler.waitGroup.Add(1)

	go sampler.pollController()

	return sampler, nil
}

// pollController updates sampler every sampling refresh interval until closed.
func (rSampler *RemotelyControlledSampler) pollController() {
	defer rSampler.waitGroup.Done()

	ticker := time.NewTicker(rSampler.samplingRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rSampler.Update()
		case <-rSampler.doneChan:
			return
		}
	}
}

// Update forces the sampler to update sampling strategy from the sampling strategy serving endpoint.
// This function is called automatically on a timer, but can also be safely called manually, e.g. from tests.
func (rSampler *RemotelyControlledSampler) Update() {
	strategy, err := rSampler.fetchStrategy()
	if err != nil {
		rSampler.logger.Printf("failed to fetch sampling strategy %s", err)

		return
	}

	rSampler.updateDefaultSampler(strategy)

	if strategy.OperationSampling != nil {
		rSampler.updatePerOpSamplers(strategy.OperationSampling)
	}

	rSampler.logger.Printf("remote sampler has been updated")
}

// fetch strategy json from the sampling strategy serving endpoint.
func (rSampler *RemotelyControlledSampler) fetchStrategy() (*ServiceStrategy, error) {
	v := url.Values{}
	v.Set("service", rSampler.serviceName)
	uri := rSampler.remoteStrategyEndpoint + "?" + v.Encode()

	req, err := http.NewRequestWithContext(context.Background(), "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed prepare request: %w", err)
	}

	httpClient := http.DefaultClient
	httpClient.Timeout = rSampler.samplingRefreshInterval

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed send request: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		err := &remotelyControlledSamplerErr{message: fmt.Sprintf("status code: %d", resp.StatusCode)}

		return nil, fmt.Errorf("received wrong status code: %d %w", resp.StatusCode, err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			rSampler.logger.Printf("failed to close HTTP response body: %s", err.Error())
		}
	}()

	strategy := &ServiceStrategy{}

	if err := json.NewDecoder(resp.Body).Decode(strategy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal strategyJSON: %w", err)
	}

	return strategy, nil
}

// updateDefaultSampler updates default sampler according to the strategy type given.
func (rSampler *RemotelyControlledSampler) updateDefaultSampler(strategy *ServiceStrategy) {
	rSampler.Lock()
	defer rSampler.Unlock()

	if strategy.StrategyType == strategyTypeRateLimiting && strategy.RateLimitingSampling != nil {
		rSampler.defaultSampler = NewRateLimitingSampler(strategy.RateLimitingSampling.MaxTracesPerSecond)
	}

	if strategy.StrategyType == strategyTypeProbabilistic && strategy.ProbabilisticSampling != nil {
		rSampler.defaultSampler = NewProbabilisticSampler(strategy.ProbabilisticSampling.SamplingRate)
	}
}

// updatePerOpSamplers updates perOperationSamplers if maxPerOperationSamplers has not been reached.
func (rSampler *RemotelyControlledSampler) updatePerOpSamplers(strategies *PerOperationSamplingStrategies) {
	rSampler.Lock()
	defer rSampler.Unlock()

	for _, perOpStrategy := range strategies.PerOperationStrategies {
		// maxPerOperationSamplers has not been reached
		if len(rSampler.perOperationSamplers) <= rSampler.maxPerOperationSamplers {
			if perOpStrategy.ProbabilisticSampling != nil {
				rSampler.perOperationSamplers[perOpStrategy.Operation] =
					NewProbabilisticSampler(perOpStrategy.ProbabilisticSampling.SamplingRate)
			} else {
				rSampler.perOperationSamplers[perOpStrategy.Operation] = rSampler.defaultSampler
			}
		} else {
			// maxPerOperationSamplers has been reached — updating only existing operation samplers
			if _, ok := rSampler.perOperationSamplers[perOpStrategy.Operation]; ok {
				rSampler.perOperationSamplers[perOpStrategy.Operation] =
					NewProbabilisticSampler(perOpStrategy.ProbabilisticSampling.SamplingRate)
			}
		}
	}
}

// Close stops sampling strategy updates.
func (rSampler *RemotelyControlledSampler) Close() {
	rSampler.once.Do(
		func() {
			close(rSampler.doneChan)
		},
	)
	rSampler.logger.Print("remote sampler has been closed")
	rSampler.waitGroup.Wait()
}

// ShouldSample implements the Sampler's ShouldSample().
//nolint:gocritic
func (rSampler *RemotelyControlledSampler) ShouldSample(param trace2.SamplingParameters) trace2.SamplingResult {
	sampler := rSampler.getSamplerForOperation(param.Name)
	isSampled, tags := sampler.IsSampled(param.TraceID)
	tags = append(tags, label.Bool("IsRecording", true), label.Bool("Sampled", isSampled))

	return trace2.SamplingResult{
		Decision:   trace2.RecordAndSample,
		Attributes: tags,
	}
}

// getSamplerForOperation returns sampler for given operation if exists, otherwise returns default sampler.
func (rSampler *RemotelyControlledSampler) getSamplerForOperation(operation string) Sampler {
	sampler, ok := rSampler.perOperationSamplers[operation]
	if !ok {
		return rSampler.defaultSampler
	}

	return sampler
}

// Description implements the Sampler's Description().
func (rSampler *RemotelyControlledSampler) Description() string {
	return fmt.Sprintf("%s, remotely controlled sampler", rSampler.serviceName)
}

type remotelyControlledSamplerErr struct {
	message string
}

func (se *remotelyControlledSamplerErr) Error() string {
	return se.message
}
