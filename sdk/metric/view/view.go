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

package view

// Config contains configuration options for a view.
type Config struct {
	// TODO (#2837): implement.
}

// TODO (#2837): add getter functions for all the internal fields of a Config.

// Option applies a configuration option value to a view Config.
type Option interface {
	apply(Config) Config
}

// TODO (#2837): implement view match options.
// TODO (#2837): implement view annotation options.

// New returns a new and configured view Config.
func New(opts ...Option) Config {
	return Config{}
}
