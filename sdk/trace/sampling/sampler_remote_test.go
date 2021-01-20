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

// // This test is possibility too large for the race detector.
// +build !race

package sampling_test // import "go.opentelemetry.io/otel/sdk/trace/sampling_test"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/label"
	trace2 "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/sampling"
	"go.opentelemetry.io/otel/trace"
)

var firstStrategiesJSON = `
		{
		  "strategyType": 0,
		  "probabilisticSampling": {
			"samplingRate": 0.8
		  },
		  "operationSampling": {
			"defaultSamplingProbability": 0.8,
			"defaultLowerBoundTracesPerSecond": 0,
			"perOperationStrategies": [
			  {
				"operation": "rack(GET /health)",
				"probabilisticSampling": {
				  "samplingRate": 0.0001
				}
			  },
			  {
				"operation": "rack(GET /v6/orders)",
				"probabilisticSampling": {
				  "samplingRate": 1
				}
			  },
			  {
				"operation": "rack(GET /v6/rivals)",
				"probabilisticSampling": {
				  "samplingRate": 0.001
				}
			  },
			  {
				"operation": "rack(GET /v6/comments)",
				"probabilisticSampling": {
				  "samplingRate": 0.001
				}
			  }
			],
			"defaultUpperBoundTracesPerSecond": 0
		  }
		}
	`

var secondStrategiesJSON = `
		{
		  "strategyType": 1,
		  "rateLimitingSampling": {
			"maxTracesPerSecond": 0.4
		  }
		}
	`

const appName = "test_service"

var _ trace2.Sampler = (*sampling.RemotelyControlledSampler)(nil)

func Test_RemotelyControlledSampler(t *testing.T) {
	t.Run("test remote sampler apply Strategy per operation", func(t *testing.T) {
		operations := []string{"rack(GET /health)", "rack(GET /v6/orders)", "rack(GET /v6/rivals)", "rack(GET /v6/comments)"}

		var wg sync.WaitGroup
		mu := &sync.Mutex{}
		defer wg.Wait()

		samplingStrategyServer := mockStrategyServer(t, appName, &wg)
		defer samplingStrategyServer.Close()

		go switchStrategiesAfterSleep(time.Second*3, t, mu, &wg)

		buf := &bytes.Buffer{}
		logger := log.New(buf, "", 0)
		samplerConfig := sampling.NewConfig(appName)
		samplerConfig.RefreshInterval = time.Second
		samplerConfig.SamplerStrategyEndpoint = samplingStrategyServer.URL
		sampler, err := sampling.NewRemotelyControlledSampler(samplerConfig, logger)
		assert.NoError(t, err)
		defer sampler.Close()

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		ctx, cancelFunc := context.WithCancel(context.Background())

		wg.Add(1)
		go func() {
			id := randStringBytes(32)
			traceID, err := trace.TraceIDFromHex(id)
			assert.NoError(t, err)
			for {
				select {
				case <-ticker.C:
					rand.Seed(time.Now().UnixNano())
					randInt := rand.Intn(len(operations))
					operationName := operations[randInt]
					params := trace2.SamplingParameters{TraceID: traceID, Name: operationName}
					samplingResult := sampler.ShouldSample(params)
					t.Logf("operationName: %s, isSampled: %v\n", operationName, samplingResult.Decision)
					for _, attr := range samplingResult.Attributes {
						t.Logf("\tattr: %v", attr)
					}
				case <-ctx.Done():
					wg.Done()
					return
				}
			}
		}()
		time.Sleep(5 * time.Second)
		cancelFunc()

		// verify the default sampler strategy has been changed
		id := randStringBytes(32)
		traceID, err := trace.TraceIDFromHex(id)
		assert.NoError(t, err)
		operationName := "new operation"
		params := trace2.SamplingParameters{TraceID: traceID, Name: operationName}
		samplingResult := sampler.ShouldSample(params)
		t.Logf("operationName: %s, isSampled: %v\n", operationName, samplingResult.Decision)
		for _, attr := range samplingResult.Attributes {
			t.Logf("\tattr: %v", attr)
		}

		want := label.String(sampling.SamplerTypeTagKey, sampling.SamplerTypeRateLimiting)
		assert.Contains(t, samplingResult.Attributes, want)
		t.Logf("logger: %s", buf.String())
	})
}

// Returns sampling strategies for the service given
func mockStrategyServer(t *testing.T, appName string, wg *sync.WaitGroup) *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()
		q := r.URL.Query()

		if q.Get("service") != appName {
			_, err := w.Write([]byte(fmt.Sprintf("wrong service name in query params: %s", q.Encode())))
			assert.NoError(t, err)
			return
		}

		serviceStrategy := parseStrategies(&firstStrategiesJSON, t)
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		mu := &sync.Mutex{}
		mu.Lock()
		strategyBytes, err := json.Marshal(serviceStrategy)
		mu.Unlock()
		assert.NoError(t, err)
		_, err = w.Write(strategyBytes)
		assert.NoError(t, err)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}

// Parses given strategies json string to map[string]*serviceStrategy
func parseStrategies(strategiesJSON *string, t *testing.T) *sampling.ServiceStrategy {
	strategies := &sampling.ServiceStrategy{}

	err := json.Unmarshal([]byte(*strategiesJSON), &strategies)
	assert.NoError(t, err)

	return strategies
}

func switchStrategiesAfterSleep(sleepTime time.Duration, t *testing.T, mu *sync.Mutex, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	time.Sleep(sleepTime)
	mu.Lock()
	firstStrategiesJSON = secondStrategiesJSON
	mu.Unlock()
	t.Log("the strategy has been switched")
}
