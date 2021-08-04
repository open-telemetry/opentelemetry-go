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

package internaltest

import (
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

type VersionTestSuite struct {
	suite.Suite
}

func (s *VersionTestSuite) TestVersionSemver() {
	testCases := []struct {
		name     string
		version  string
		expected string
	}{
		{
			name:     "sdk trace version",
			version:  trace.Version(),
			expected: "1.0.0-RC2",
		},
		{
			name:     "sdk metric version",
			version:  metric.Version(),
			expected: "v0.22.0",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.Assert().Equal(tc.version, tc.expected)
		})
	}
}
