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
	"encoding/json"
	"os/exec"
	"sort"
	"strings"
)

type depTracker struct {
	goListResult map[string]goListResult
	buildAndTest map[string]map[string]struct{}
}

type goListResult struct {
	Deps         []string
	TestImports  []string
	XTestImports []string
}

func (dt *depTracker) goList(pkg string, m mod) (goListResult, error) {
	if dt.goListResult == nil {
		dt.goListResult = make(map[string]goListResult)
	}

	if cached, ok := dt.goListResult[pkg]; ok {
		return cached, nil
	}

	cmd := exec.Command("go", "list", "-json", pkg)
	cmd.Dir = m.filePath
	out, err := cmd.Output()
	if err != nil {
		return goListResult{}, err
	}

	var result goListResult
	err = json.Unmarshal(out, &result)
	if err != nil {
		return result, err
	}

	dt.goListResult[pkg] = result
	return result, err
}

func (dt *depTracker) depsIncludingTest(pkg string, m mod) (map[string]struct{}, error) {
	if dt.buildAndTest == nil {
		dt.buildAndTest = make(map[string]map[string]struct{})
	}

	if cached, ok := dt.buildAndTest[pkg]; ok {
		return cached, nil
	}

	res, err := dt.goList(pkg, m)
	if err != nil {
		return nil, err
	}

	set := make(map[string]struct{})
	for _, pkg := range res.Deps {
		set[pkg] = struct{}{}
	}

	// Necessary because of https://github.com/golang/go/issues/23806
	for _, pkg := range res.TestImports {
		set[pkg] = struct{}{}

		res, err := dt.goList(pkg, m)
		if err != nil {
			return nil, err
		}

		for _, pkg := range res.Deps {
			set[pkg] = struct{}{}
		}
	}
	for _, pkg := range res.XTestImports {
		set[pkg] = struct{}{}

		res, err := dt.goList(pkg, m)
		if err != nil {
			return nil, err
		}

		for _, pkg := range res.Deps {
			set[pkg] = struct{}{}
		}
	}
	return set, nil
}

func (dt *depTracker) modDeps(m mod, ms mods) ([]mod, error) {
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = m.filePath
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	modOfPkg := map[string]mod{}
	set := map[mod]struct{}{}
	for _, pkg := range strings.Split(string(out), "\n") {
		pkgs, err := dt.depsIncludingTest(pkg, m)
		if err != nil {
			return nil, err
		}

		for pkg := range pkgs {
			mod, ok := modOfPkg[pkg]
			if !ok {
				for _, m := range ms {
					if strings.HasPrefix(pkg, m.importPath) && len(m.importPath) > len(mod.importPath) {
						mod = m
					}
				}
				modOfPkg[pkg] = mod
			}
			if mod.importPath != "" {
				set[mod] = struct{}{}
			}
		}
	}

	var result []mod
	for m := range set {
		result = append(result, m)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].importPath < result[j].importPath })
	return result, nil
}
