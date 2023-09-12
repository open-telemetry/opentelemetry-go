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

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

var urls = []string{
	"https://opentelemetry.io/schemas/1.21.0",
	"https://opentelemetry.io/schemas/1.20.0",
	"https://opentelemetry.io/schemas/1.19.0",
	"https://opentelemetry.io/schemas/1.18.0",
	"https://opentelemetry.io/schemas/1.17.0",
	"https://opentelemetry.io/schemas/1.16.0",
	"https://opentelemetry.io/schemas/1.15.0",
	"https://opentelemetry.io/schemas/1.14.0",
	"https://opentelemetry.io/schemas/1.13.0",
	"https://opentelemetry.io/schemas/1.12.0",
	"https://opentelemetry.io/schemas/1.11.0",
	"https://opentelemetry.io/schemas/1.10.0",
	"https://opentelemetry.io/schemas/1.9.0",
	"https://opentelemetry.io/schemas/1.8.0",
	"https://opentelemetry.io/schemas/1.7.0",
	"https://opentelemetry.io/schemas/1.6.1",
	// Does not exist: "https://opentelemetry.io/schemas/1.6.0"
	"https://opentelemetry.io/schemas/1.5.0",
	"https://opentelemetry.io/schemas/1.4.0",
	// Does not exist: "https://opentelemetry.io/schemas/1.3.0"
	// Does not exist: "https://opentelemetry.io/schemas/1.2.0"
	// Does not exist: "https://opentelemetry.io/schemas/1.1.0"
	// Does not exist: "https://opentelemetry.io/schemas/1.0.1"
	// Does not exist: "https://opentelemetry.io/schemas/1.0.0"
	// Does not exist: "https://opentelemetry.io/schemas/0.7.0"
	// Does not exist: "https://opentelemetry.io/schemas/0.6.0"
	// Does not exist: "https://opentelemetry.io/schemas/0.5.0"
	// Does not exist: "https://opentelemetry.io/schemas/0.4.0"
	// Does not exist: "https://opentelemetry.io/schemas/0.3"
	// Does not exist: "https://opentelemetry.io/schemas/0.2"
	// Does not exist: "https://opentelemetry.io/schemas/0.1"
}

func download(u string) error {
	resp, err := http.Get(u)
	if err != nil {
		return fmt.Errorf("failed to get %q: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request error: %d", resp.StatusCode)
	}

	fName := fmt.Sprintf("%s.yaml", path.Base(u))
	f, err := os.Create(fName)
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", fName, err)
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write %q: %w", fName, err)
	}
	return nil
}

func main() {
	for _, u := range urls {
		err := download(u)
		if err != nil {
			log.Printf("failed to download %q: %s", u, err)
			continue
		}
		log.Printf("downloaded %q", u)
	}
}
