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

package cmd // import "go.opentelemetry.io/otel/sdk/resource/internal/schema/generate/cmd"

import (
	"fmt"
	"net/http"
	"os"
	"sort"

	"github.com/Masterminds/semver/v3"

	ast10 "go.opentelemetry.io/otel/schema/v1.0/ast"
	sUtil "go.opentelemetry.io/otel/schema/v1.1"
	"go.opentelemetry.io/otel/schema/v1.1/ast"
	"go.opentelemetry.io/otel/schema/v1.1/types"
)

var schemaURL = "https://opentelemetry.io/schemas/1.21.0"

type entry struct {
	Version   *semver.Version
	All       ast10.Attributes
	Resources ast10.Attributes
}

func newEntry(ver types.TelemetryVersion, def ast.VersionDef) (entry, error) {
	v, err := semver.NewVersion(string(ver))
	return entry{
		Version:   v,
		All:       def.All,
		Resources: def.Resources,
	}, err
}

func entries(s *ast.Schema) ([]entry, error) {
	if s == nil {
		return nil, nil
	}

	es := make([]entry, 0, len(s.Versions))
	for v, def := range s.Versions {
		e, err := newEntry(v, def)
		if err != nil {
			return nil, err
		}

		es = append(es, e)
	}

	return es, nil
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
	f, err := os.Open(local)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s, err := sUtil.Parse(f)
	if err != nil {
		return nil, err
	}
	return entries(s)
}

func loadRemote() (data []entry, err error) {
	resp, err := http.Get(schemaURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request error: %d", resp.StatusCode)
	}

	s, err := sUtil.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return entries(s)
}
