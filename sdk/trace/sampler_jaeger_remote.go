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

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"

	jaeger_api_v2 "github.com/jaegertracing/jaeger/proto-gen/api_v2"
	"go.opentelemetry.io/otel"
)

const (
	jaegerRemoteDefaultSamplingRefreshInterval = time.Minute
	jaegerRemoteDefaultSamplerFraction         = 0.001
)

type jaegerRemoteSampler struct {
	samplingRefreshInterval time.Duration

	fetcher jaegerSamplingStrategyFetcher

	sync.RWMutex
	lastStrategyResponse jaeger_api_v2.SamplingStrategyResponse
	sampler              Sampler
}

func (s *jaegerRemoteSampler) ShouldSample(p SamplingParameters) SamplingResult {
	s.RLock()
	defer s.RUnlock()

	return s.sampler.ShouldSample(p)
}

func (s *jaegerRemoteSampler) Description() string {
	return "JaegerRemoteSampler"
}

func (s *jaegerRemoteSampler) pollSamplingStrategies() {
	ticker := time.NewTicker(s.samplingRefreshInterval)
	for {
		<-ticker.C
		err := s.updateSamplingStrategies()
		if err != nil {
			otel.Handle(fmt.Errorf("updating jaeger remote sampling strategies failed: %w", err))
		}
	}
}

// updateSamplingStrategies fetches the sampling strategy from backend server.
// This function is called automatically on a timer, but can also be safely
// manually, e.g. from tests.
func (s *jaegerRemoteSampler) updateSamplingStrategies() error {
	strategies, err := s.fetcher.Fetch()
	if err != nil {
		return fmt.Errorf("fetching failed: %w", err)
	}

	if !s.hasChanges(strategies) {
		return nil
	}

	err = s.loadSamplingStrategies(strategies)
	if err != nil {
		return fmt.Errorf("loading failed: %w", err)
	}

	return nil
}

func (s *jaegerRemoteSampler) hasChanges(other jaeger_api_v2.SamplingStrategyResponse) bool {
	s.RLock()
	defer s.RUnlock()

	return s.lastStrategyResponse.StrategyType != other.StrategyType ||
		s.lastStrategyResponse.ProbabilisticSampling != other.ProbabilisticSampling ||
		s.lastStrategyResponse.RateLimitingSampling != other.RateLimitingSampling ||
		s.lastStrategyResponse.OperationSampling != other.OperationSampling
}

func (s *jaegerRemoteSampler) loadSamplingStrategies(strategies jaeger_api_v2.SamplingStrategyResponse) error {
	// TODO add support for rate limiting
	if strategies.StrategyType != jaeger_api_v2.SamplingStrategyType_PROBABILISTIC {
		return fmt.Errorf("only strategy type PROBABILISTC is supported, got %s", strategies.StrategyType)
	}
	// TODO add support for per operation sampling
	if strategies.OperationSampling != nil {
		return fmt.Errorf("per operation sampling is not supported")
	}

	// TODO should we implement this validation ourselves?
	if strategies.ProbabilisticSampling == nil {
		return fmt.Errorf("strategy is probabilistic, but struct is empty")
	}

	s.Lock()
	defer s.Unlock()

	s.lastStrategyResponse = strategies

	s.sampler = TraceIDRatioBased(strategies.ProbabilisticSampling.SamplingRate)

	return nil
}

// JaegerRemoteSampler returns a Sampler that consults a Jaeger remote agent
// for the sampling strategies for this service.
// TODO add option for samplingRefreshInterval
// TODO should we make serviceName an option as well?
func JaegerRemoteSampler(serviceName string, samplingServerURL string) Sampler {
	sampler := &jaegerRemoteSampler{
		fetcher: jaegerSamplingStrategyFetcherImpl{
			serviceName:       serviceName,
			samplingServerUrl: samplingServerURL,
			httpClient: &http.Client{
				Timeout: 10 * time.Second,
			},
		},
		samplingRefreshInterval: jaegerRemoteDefaultSamplingRefreshInterval,
		sampler:                 TraceIDRatioBased(jaegerRemoteDefaultSamplerFraction),
	}
	// TODO we spawn a go routine, should we clean it up before quitting?
	go sampler.pollSamplingStrategies()
	return sampler
}

type jaegerSamplingStrategyFetcher interface {
	Fetch() (jaeger_api_v2.SamplingStrategyResponse, error)
}

type jaegerSamplingStrategyFetcherImpl struct {
	serviceName       string
	samplingServerUrl string
	httpClient        *http.Client
}

func (f jaegerSamplingStrategyFetcherImpl) Fetch() (s jaeger_api_v2.SamplingStrategyResponse, err error) {
	uri := f.samplingServerUrl + "?service=" + url.QueryEscape(f.serviceName)

	resp, err := f.httpClient.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		return s, fmt.Errorf("request failed (%d): %s", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &s)
	return
}
