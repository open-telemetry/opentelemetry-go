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
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"

	sUtil "go.opentelemetry.io/otel/schema/v1.1"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
)

var schemaURLs = []string{
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

type entry struct {
	Version *semver.Version
	Schema  *ast.Schema
}

func newEntry(s *ast.Schema) (entry, error) {
	if s == nil {
		return entry{}, errors.New("nil schema")
	}

	// https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md#schema-url
	u, err := url.Parse(s.SchemaURL)
	if err != nil {
		return entry{}, fmt.Errorf("invalid schema URL %q: %w", s.SchemaURL, err)
	}

	verStr := u.Path[strings.LastIndex(u.Path, "/")+1:]
	ver, err := semver.NewVersion(verStr)
	if err != nil {
		return entry{}, fmt.Errorf("invalid schema URL version %q: %w", verStr, err)
	}

	return entry{Version: ver, Schema: s}, nil
}

func load(local string) (any, error) {
	var (
		data []entry
		err  error
	)
	if local != "" {
		data, err = loadLocal(local)
	} else {
		data, err = loadRemote()
	}
	if err != nil {
		return nil, err
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i].Version.LessThan(data[j].Version)
	})
	return data, nil
}

func loadLocal(local string) (data []entry, err error) {
	const suffix = ".yaml"
	err = filepath.WalkDir(local, func(p string, _ fs.DirEntry, err error) error {
		if err != nil || path.Ext(p) != suffix {
			return err
		}

		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()

		s, err := sUtil.Parse(f)
		if err != nil {
			return err
		}

		e, err := newEntry(s)
		if err != nil {
			return err
		}

		data = append(data, e)
		return nil
	})
	return data, err
}

func loadRemote() (data []entry, err error) {
	var e entry
	for _, u := range schemaURLs {
		e, err = download(u)
		if err != nil {
			return data, fmt.Errorf("failed to download %q: %w", u, err)
		}
		data = append(data, e)
	}
	return data, nil
}

func download(u string) (entry, error) {
	resp, err := http.Get(u)
	if err != nil {
		return entry{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return entry{}, fmt.Errorf("request error: %d", resp.StatusCode)
	}

	s, err := sUtil.Parse(resp.Body)
	if err != nil {
		return entry{}, err
	}

	return newEntry(s)
}
