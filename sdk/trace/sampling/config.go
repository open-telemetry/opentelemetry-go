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
	"time"

	mowCli "github.com/jawher/mow.cli"
	"github.com/urfave/cli/v2"
)

// Config for constructor.
type Config struct {
	AppName                 string
	SamplerStrategyEndpoint string
	RefreshInterval         time.Duration
	MaxOperations           int
	Rate                    float64
}

const (
	defaultRefreshInterval = 60 * time.Second
	defaultMaxOperations   = 42
	defaultRate            = 0.5
)

// NewConfig default constructor.
func NewConfig(appName string) Config {
	return Config{
		AppName:         appName,
		RefreshInterval: defaultRefreshInterval,
		MaxOperations:   defaultMaxOperations,
		Rate:            defaultRate,
	}
}

// BuildFlags of config to cli.Flag.
func (config *Config) BuildFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "remote-sampling-addr",
			Usage:       "Remote sampling HOST:PORT",
			EnvVars:     []string{"JAEGER_REMOTE_SAMPLER_HOST"},
			Destination: &config.SamplerStrategyEndpoint,
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "sampler-app-name",
			Usage:       "Sampler app name prefix",
			EnvVars:     []string{"JAEGER_SERVICE_NAME"},
			Destination: &config.AppName,
		},
	}
}

// BuildMowCliFlags of config to mow.cli.Flag.
func (config *Config) BuildMowCliFlags(cmd *mowCli.Cmd, addSpecFlag bool) {
	if addSpecFlag {
		cmd.Spec = cmd.Spec + " --remote-sampling-addr" + " --sampler-app-name"
	}

	cmd.StringPtr(&config.SamplerStrategyEndpoint, mowCli.StringOpt{
		Name:   "remote-sampling-addr",
		Desc:   "Remote sampling HOST:PORT",
		EnvVar: "JAEGER_REMOTE_SAMPLER_HOST",
		Value:  config.SamplerStrategyEndpoint,
	})
	cmd.StringPtr(&config.AppName, mowCli.StringOpt{
		Name:   "sampler-app-name",
		Desc:   "Sampler app name prefix",
		EnvVar: "JAEGER_SERVICE_NAME",
		Value:  config.AppName,
	})
}
