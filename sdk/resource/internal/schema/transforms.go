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

package schema

import (
	"github.com/Masterminds/semver/v3"
	"go.opentelemetry.io/otel/schema/v1.0/ast"
)

// TODO: generate this.
var transforms = []transform{
	{Version: semver.New(1, 4, 0, "", "")},
	{Version: semver.New(1, 5, 0, "", "")},
	{Version: semver.New(1, 6, 1, "", "")},
	{Version: semver.New(1, 7, 0, "", "")},
	{Version: semver.New(1, 8, 0, "", "")},
	{Version: semver.New(1, 9, 0, "", "")},
	{Version: semver.New(1, 10, 0, "", "")},
	{Version: semver.New(1, 11, 0, "", "")},
	{Version: semver.New(1, 12, 0, "", "")},
	{Version: semver.New(1, 13, 0, "", "")},
	{Version: semver.New(1, 14, 0, "", "")},
	{Version: semver.New(1, 15, 0, "", "")},
	{Version: semver.New(1, 16, 0, "", "")},
	{Version: semver.New(1, 17, 0, "", "")},
	{Version: semver.New(1, 18, 0, "", "")},
	{
		Version: semver.New(1, 19, 0, "", ""),
		Resources: ast.Attributes{
			Changes: []ast.AttributeChange{
				{
					RenameAttributes: &ast.RenameAttributes{
						AttributeMap: map[string]string{
							"browser.user_agent": "user_agent.original",
						},
					},
				},
			},
		},
	},
	{Version: semver.New(1, 20, 0, "", "")},
	{Version: semver.New(1, 21, 0, "", "")},
}
