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

package otlphttp

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func applyEnvConfigs(cfg *config, getEnv func(string) string) *config {
	opts := getOptionsFromEnv(getEnv)
	for _, opt := range opts {
		opt.Apply(cfg)
	}
	return cfg
}

func getOptionsFromEnv(env func(string) string) []Option {
	var opts []Option

	// Endpoint
	if v, ok := getEnv(env, "ENDPOINT"); ok {
		opts = append(opts, WithEndpoint(v))
	}
	if v, ok := getEnv(env, "TRACES_ENDPOINT"); ok {
		opts = append(opts, WithTracesEndpoint(v))
	}
	if v, ok := getEnv(env, "METRICS_ENDPOINT"); ok {
		opts = append(opts, WithMetricsEndpoint(v))
	}

	// Certificate File
	// TODO: add certificate file env config support

	// Headers
	if h, ok := getEnv(env, "HEADERS"); ok {
		opts = append(opts, WithHeaders(stringToHeader(h)))
	}
	if h, ok := getEnv(env, "TRACES_HEADERS"); ok {
		opts = append(opts, WithTracesHeaders(stringToHeader(h)))
	}
	if h, ok := getEnv(env, "METRICS_HEADERS"); ok {
		opts = append(opts, WithMetricsHeaders(stringToHeader(h)))
	}

	// Compression
	if c, ok := getEnv(env, "COMPRESSION"); ok {
		opts = append(opts, WithCompression(stringToCompression(c)))
	}
	if c, ok := getEnv(env, "TRACES_COMPRESSION"); ok {
		opts = append(opts, WithTracesCompression(stringToCompression(c)))
	}
	if c, ok := getEnv(env, "METRICS_COMPRESSION"); ok {
		opts = append(opts, WithMetricsCompression(stringToCompression(c)))
	}

	// Timeout
	if t, ok := getEnv(env, "TIMEOUT"); ok {
		if d, err := strconv.Atoi(t); err == nil {
			opts = append(opts, WithTimeout(time.Duration(d)*time.Millisecond))
		}
	}
	if t, ok := getEnv(env, "TRACES_TIMEOUT"); ok {
		if d, err := strconv.Atoi(t); err == nil {
			opts = append(opts, WithTracesTimeout(time.Duration(d)*time.Millisecond))
		}
	}
	if t, ok := getEnv(env, "METRICS_TIMEOUT"); ok {
		if d, err := strconv.Atoi(t); err == nil {
			opts = append(opts, WithMetricsTimeout(time.Duration(d)*time.Millisecond))
		}
	}

	return opts
}

// getEnv gets an OTLP environment variable value of the specified key using the env function.
// This function already prepends the OTLP prefix to all key lookup.
func getEnv(env func(string) string, key string) (string, bool) {
	v := strings.TrimSpace(env(fmt.Sprintf("OTEL_EXPORTER_OTLP_%s", key)))
	return v, v != ""
}

func stringToCompression(value string) Compression {
	switch value {
	case "gzip":
		return GzipCompression
	}

	return NoCompression
}

func stringToHeader(value string) map[string]string {
	headersPairs := strings.Split(value, ",")
	headers := make(map[string]string)

	for _, header := range headersPairs {
		nameValue := strings.SplitN(header, "=", 2)
		if len(nameValue) < 2 {
			continue
		}
		name, err := url.QueryUnescape(nameValue[0])
		if err != nil {
			continue
		}
		trimmedName := strings.TrimSpace(name)
		value, err := url.QueryUnescape(nameValue[1])
		if err != nil {
			continue
		}
		trimmedValue := strings.TrimSpace(value)

		headers[trimmedName] = trimmedValue
	}

	return headers
}
