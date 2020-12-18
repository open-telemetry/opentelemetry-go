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

package otlpenv

import (
	"os"
)

type config struct {
	env []string
}

func newConfig(opts ...Option) config {
	cfg := config{
		env: os.Environ(),
	}
	for _, opt := range opts {
		opt.Apply(&cfg)
	}
	return cfg
}

// Option applies an option to the env driver.
type Option interface {
	Apply(*config)
}

type environmentOption []string

func (o environmentOption) Apply(cfg *config) {
	cfg.env = ([]string)(o)
}

// WithEnvironment tells the driver to use the passed list of strings
// as a source of environment variables. Each string in the list
// should be in format `key=value`.
func WithEnvironment(env []string) Option {
	return (environmentOption)(env)
}
