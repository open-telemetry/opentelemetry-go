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

// The crosslink tool generates and maintains replace directives in all
// the go.mod files within this repository. Some directives are superfluous
// (e.g. because the replaced module doesn't occur in the dependency tree),
// but we generate them anyway for the sake of consistency (#1529 tracks
// pruning this to a mininal set).
//
// In particular, we generate a replace directive from each module to itself
// (i.e., the target path "./"). This is actually necessary in the presence of
// cyclic dependencies between modules.

package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.opentelemetry.io/otel/internal/tools"
	"golang.org/x/mod/modfile"
)

func crossLink(m []*modfile.File) error {
	for _, from := range m {
		basepath := filepath.Dir(from.Syntax.Name)
		for _, to := range m {
			newPath, err := filepath.Rel(basepath, filepath.Dir(to.Syntax.Name))
			if err != nil {
				return err
			}
			switch {
			case newPath == ".", newPath == "..":
				newPath += "/"
			case !strings.HasPrefix(newPath, ".."):
				newPath = "./" + newPath
			}
			from.AddReplace(to.Module.Mod.Path, "", newPath, "")
		}

		from.Cleanup()

		f, err := os.OpenFile(from.Syntax.Name, os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		if _, err = f.Write(modfile.Format(from.Syntax)); err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	root, err := tools.FindRepoRoot()
	if err != nil {
		log.Fatalf("unable to find repo root: %v", err)
	}

	mods, err := root.FindModules()
	if err != nil {
		log.Fatalf("unable to list modules: %v", err)
	}

	if err := tools.PrintModFiles(os.Stdout, mods); err != nil {
		log.Fatalf("unable to print modules: %v", err)
	}

	if err := crossLink(mods); err != nil {
		log.Fatalf("unable to crosslink: %v", err)
	}
}
