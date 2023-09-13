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

package cmd

import (
	"fmt"
	"os"
)

func Run(dest, local string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to open desination %q: %w", dest, err)
	}
	defer f.Close()

	data, err := load(local)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}

	err = render(f, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}
