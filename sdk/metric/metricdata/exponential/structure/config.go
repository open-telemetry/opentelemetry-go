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

package structure // import "go.opentelemetry.io/otel/sdk/metric/metricdata/exponential/structure"

import "fmt"

// DefaultMaxSize is the default maximum number of buckets per
// positive or negative number range.  The value 160 is specified by
// OpenTelemetry--yields a maximum relative error of less than 5% for
// data with contrast 10**5 (e.g., latencies in the range 1ms to 100s).
// See the derivation here:
// https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/sdk.md#exponential-bucket-histogram-aggregation
const DefaultMaxSize int32 = 160

// MinSize is the smallest reasonable configuration, which is small
// enough to contain the entire normal floating point range at
// MinScale.
const MinSize = 2

// MaximumMaxSize is an arbitrary limit meant to limit accidental use
// of giant histograms.
const MaximumMaxSize = 16384

// Config contains configuration for exponential histogram creation.
type Config struct {
	maxSize int32
}

// Option is the interface that applies a configuration option.
type Option interface {
	// apply sets the Option value of a config.
	apply(Config) Config
}

// WithMaxSize sets the maximum size of each range (positive and/or
// negative) in the histogram.
func WithMaxSize(size int32) Option {
	return maxSize(size)
}

// maxSize is an option to set the maximum histogram size.
type maxSize int32

// apply implements Option.
func (ms maxSize) apply(cfg Config) Config {
	cfg.maxSize = int32(ms)
	return cfg
}

// NewConfig returns an exponential histogram configuration with
// defaults and limits applied.
func NewConfig(opts ...Option) Config {
	var cfg Config
	for _, opt := range opts {
		cfg = opt.apply(cfg)
	}
	return cfg
}

// Validate returns true for valid configurations.
func (c Config) Valid() bool {
	_, err := c.Validate()
	return err == nil
}

// Validate returns the nearest valid Config object to the input and a
// boolean indicating whether the the input was a valid
// configurations.
func (c Config) Validate() (Config, error) {
	if c.maxSize >= MinSize && c.maxSize <= MaximumMaxSize {
		return c, nil
	}
	if c.maxSize == 0 {
		c.maxSize = DefaultMaxSize
		return c, nil
	}
	err := fmt.Errorf("invalid histogram size: %d", c.maxSize)
	if c.maxSize < 0 {
		c.maxSize = DefaultMaxSize
	} else if c.maxSize < MinSize {
		c.maxSize = MinSize
	} else if c.maxSize > MaximumMaxSize {
		c.maxSize = MaximumMaxSize
	}
	return c, err
}
