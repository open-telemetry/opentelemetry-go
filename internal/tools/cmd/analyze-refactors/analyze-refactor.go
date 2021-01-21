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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"text/template"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type (
	mod     string
	modInfo struct{ Stable bool }
	pkg     string
	pkgInfo struct {
		Mod  mod
		Deps []pkg
	}
	goListPkg struct {
		ImportPath string
		Module     struct{ Path string }
		Imports    []string
	}
	info struct {
		Mods map[mod]modInfo
		Pkgs map[pkg]pkgInfo
	}
)

func (i *info) incorporate(desc goListPkg, prefix, dir string) error {
	if !strings.HasPrefix(desc.Module.Path, prefix) {
		return nil
	}

	m := mod("." + strings.TrimPrefix(desc.Module.Path, prefix))
	_, ok := i.Mods[m]
	if !ok {
		i.Mods[m] = modInfo{}
	}

	if !strings.HasPrefix(desc.ImportPath, prefix) {
		return fmt.Errorf("inconsistent prefix: pkg=%s, mod=%s", desc.ImportPath, desc.Module.Path)
	}

	p := pkg("." + strings.TrimPrefix(desc.ImportPath, prefix))
	_, ok = i.Pkgs[p]
	if ok {
		return fmt.Errorf("duplicate go list record for package %s", p)
	}

	var deps []pkg
	for _, imp := range desc.Imports {
		if !strings.HasPrefix(imp, prefix) {
			continue
		}
		q := pkg("." + strings.TrimPrefix(imp, prefix))
		deps = append(deps, q)
	}
	i.Pkgs[p] = pkgInfo{Mod: m, Deps: deps}
	return nil
}

func (i *info) incorporateDir(prefix, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.Name() != "go.mod" {
			return nil
		}
		return i.incorporateMod(prefix, filepath.Dir(path))
	})
}

func (i *info) incorporateMod(prefix, dir string) error {
	cmd := exec.Command("go", "list", "-json", "./...")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var desc goListPkg
		err := dec.Decode(&desc)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("decoding error: %v", err)
		}

		err = i.incorporate(desc, prefix, dir)
		if err != nil {
			return fmt.Errorf("processing error: %v", err)
		}
	}
}

var dotTemplate = template.Must(template.New("").Parse(`
digraph {
	node[shape=box];
	rankdir = "LR";
	fontfamily = "Helvetica";

	{{range $mod, $info := .Mods}}
	subgraph "cluster_{{$mod}}" {
		fontsize = 20;
		label = "{{$mod}}";
		style=filled;
		{{if $info.Stable}}
			color=lightgreen;
		{{else}}
			color=lightgrey;
		{{end}}
	}
	{{end}}

	{{range $pkg, $info := .Pkgs}}
		subgraph "cluster_{{$info.Mod}}" {
			"{{$pkg}}";
		}
		{{range $info.Deps}}
			"{{$pkg}}" -> "{{.}}";
		{{end}}
	{{end}}
}
`))

func (i info) dot() {
	if err := dotTemplate.Execute(os.Stdout, i); err != nil {
		log.Fatal(err)
	}
}

var dumpTemplate = template.Must(template.New("").Parse(`
MODULE	STABLE
{{range $mod, $info := .Mods -}}
{{$mod}}	{{$info.Stable}}
{{end -}}
{{"\f"}}
PACKAGE	MODULE	DEPS
{{range $pkg, $info := .Pkgs -}}
{{$pkg}}	{{$info.Mod}}	{{range $info.Deps}}{{.}} {{end}}
{{end -}}
{{"\f"}}
`))

func (i info) dump() {
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if err := dumpTemplate.Execute(tw, i); err != nil {
		log.Fatal(err)
	}
}

func hasPrefix(s string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func (i *info) prune(prefix ...string) info {
	mods := map[mod]modInfo{}
	pkgs := map[pkg]pkgInfo{}
	for m, mi := range i.Mods {
		if hasPrefix(string(m), prefix...) {
			continue
		}
		mods[m] = mi
	}
	for p, pi := range i.Pkgs {
		if hasPrefix(string(p), prefix...) {
			continue
		}
		var deps []pkg
		for _, d := range pi.Deps {
			if hasPrefix(string(d), prefix...) {
				continue
			}
			deps = append(deps, d)
		}
		pi.Deps = deps
		pkgs[p] = pi
	}
	return info{Mods: mods, Pkgs: pkgs}
}

func (i info) checkPackageCycles() {
	g := simple.NewDirectedGraph()
	nodes := map[pkg]graph.Node{}
	revNodes := map[int64]pkg{}
	for p := range i.Pkgs {
		n := g.NewNode()
		nodes[p] = n
		revNodes[n.ID()] = p
		g.AddNode(n)
	}
	for p, pi := range i.Pkgs {
		from := nodes[p]
		for _, d := range pi.Deps {
			g.SetEdge(g.NewEdge(from, nodes[d]))
		}
	}

	cycles := topo.DirectedCyclesIn(g)
	if len(topo.DirectedCyclesIn(g)) != 0 {
		var sb strings.Builder
		if len(cycles) != 0 {
			for _, cyc := range cycles {
				for i, el := range cyc {
					if i != 0 {
						sb.WriteString(" -> ")
					}
					sb.WriteString(string(revNodes[el.ID()]))
				}
				sb.WriteString("\n")
			}
			log.Fatalf("found package import cycle(s):\n%s", sb.String())
		}
	}
}

func (i info) checkModuleCycles() {
	g := simple.NewDirectedGraph()
	nodes := map[mod]graph.Node{}
	revNodes := map[int64]mod{}
	for m := range i.Mods {
		n := g.NewNode()
		nodes[m] = n
		revNodes[n.ID()] = m
		g.AddNode(n)
	}
	for _, pi := range i.Pkgs {
		from := nodes[pi.Mod]
		for _, d := range pi.Deps {
			dp := i.Pkgs[d]
			if dp.Mod == pi.Mod {
				continue
			}
			g.SetEdge(g.NewEdge(from, nodes[dp.Mod]))
		}
	}

	cycles := topo.DirectedCyclesIn(g)
	var sb strings.Builder
	if len(cycles) != 0 {
		for _, cyc := range cycles {
			for i, el := range cyc {
				if i != 0 {
					sb.WriteString(" -> ")
				}
				sb.WriteString(string(revNodes[el.ID()]))
			}
			sb.WriteString("\n")
		}
		log.Fatalf("found module dependency cycle(s):\n%s", sb.String())
	}
}

func (i info) checkUnstableImports() {
	for p, pi := range i.Pkgs {
		if !i.Mods[pi.Mod].Stable {
			continue
		}

		for _, d := range pi.Deps {
			dp := i.Pkgs[d]
			if !i.Mods[dp.Mod].Stable {
				log.Fatalf("package %s in stable module imports package %s in unstable module", p, d)
			}
		}
	}
}

func (i info) checkInternalImport() bool {
	for p, pi := range i.Pkgs {
		for _, d := range pi.Deps {
			ix := strings.Index(string(d), "internal")
			if ix == -1 {
				continue
			}

			if !strings.HasPrefix(string(p), string(d)[:ix-1]) {
				log.Fatalf("package %s imports inaccessible package %s", p, d)
			}
		}
	}

	return false
}

func (i info) check() {
	i.checkPackageCycles()
	i.checkInternalImport()
	i.checkModuleCycles()
	i.checkUnstableImports()
}

func (i *info) createModule(path string) {
	m := mod(path)
	i.Mods[m] = modInfo{}
	for p, pi := range i.Pkgs {
		if strings.HasPrefix(string(p), path) {
			pi.Mod = m
			i.Pkgs[p] = pi
		}
	}
	i.check()
}

func (i *info) createPackage(p string) {
	for m := range i.Mods {
		if strings.HasPrefix(p, string(m)) {
			i.Pkgs[pkg(p)] = pkgInfo{Mod: m}
			break
		}
	}
	i.check()

}

func (i *info) copyPackage(from, to string) {
	pi := i.Pkgs[pkg(from)]
	var nearest mod
	for m := range i.Mods {
		if !strings.HasPrefix(to, string(m)) {
			continue
		}
		if len(m) < len(nearest) {
			continue
		}
		nearest = m
	}
	if len(nearest) == 0 {
		log.Fatalf("copy destination %s is outside all known modules", to)
	}
	pi.Mod = nearest
	i.Pkgs[pkg(to)] = pi
	i.check()
}

func (i *info) movePackage(from, to string) {
	i.copyPackage(from, to)
	for p, pi := range i.Pkgs {
		for ix, d := range pi.Deps {
			if string(d) == from {
				deps := pi.Deps
				pi.Deps = append(pi.Deps[:ix], pkg(to))
				pi.Deps = append(pi.Deps, deps[ix+1:]...)
				i.Pkgs[p] = pi
				continue
			}
		}
	}
	delete(i.Pkgs, pkg(from))
	i.check()
}

func (i *info) duplicateDependency(from, toBefore, toAfter string) {
	i.copyPackage(toBefore, toAfter)
	i.replaceDependency(from, toBefore, toAfter)
}

func (i *info) replaceDependency(from, toBefore, toAfter string) {
	pi := i.Pkgs[pkg(from)]
	for ix, d := range pi.Deps {
		if string(d) == toBefore {
			deps := pi.Deps
			pi.Deps = append(pi.Deps[:ix], pkg(toAfter))
			pi.Deps = append(pi.Deps, deps[ix+1:]...)
			i.Pkgs[pkg(from)] = pi
			return
		}
	}
	i.check()
}

func (i *info) removeDependency(from, to string) {
	pi := i.Pkgs[pkg(from)]
	for ix, d := range pi.Deps {
		if string(d) == to {
			pi.Deps = append(pi.Deps[:ix], pi.Deps[ix+1:]...)
			i.Pkgs[pkg(from)] = pi
			return
		}
	}
}

func (i *info) stabilize(path string) {
	mi := i.Mods[mod(path)]
	mi.Stable = true
	i.Mods[mod(path)] = mi
	i.check()
}

func main() {
	i := info{
		Mods: map[mod]modInfo{},
		Pkgs: map[pkg]pkgInfo{},
	}
	err := i.incorporateDir("go.opentelemetry.io/otel", "../..")
	if err != nil {
		log.Fatalf("unable to incorporate: %v", err)
	}

	i.duplicateDependency("./label", "./internal", "./label/internal")
	i.createModule("./label")
	i.stabilize("./label")

	i.removeDependency("./semconv", "./codes") // move semconv.SpanStatusFromHTTPStatusCode into codes
	i.createModule("./semconv")
	i.stabilize("./semconv")

	i.movePackage("./codes", "./trace/codes")
	i.createModule("./trace")
	i.stabilize("./trace")

	i.duplicateDependency("./propagation", "./internal/baggage", "./propagation/internal/baggage")
	i.createModule("./propagation")
	i.stabilize("./propagation")

	i.duplicateDependency("./baggage", "./internal/baggage", "./baggage/internal")
	i.createModule("./baggage")
	i.stabilize("./baggage")

	i.duplicateDependency("./metric/number", "./internal", "./metric/internal")
	i.movePackage("./unit", "./metric/unit")
	i.createModule("./metric")

	i.createPackage("./errors") // split from .
	i.replaceDependency("./internal/metric", ".", "./errors")
	i.createModule("./errors")
	i.stabilize("./errors")

	i.removeDependency(".", "./metric")
	i.removeDependency("./internal/global", "./metric")
	i.removeDependency("./internal/global", "./metric/number")
	i.removeDependency("./internal/global", "./metric/registry")
	i.duplicateDependency("./sdk/metric", "./internal/metric", "./sdk/internal/metric")
	i.movePackage("./internal/metric", "./oteltest/internal/metric")

	i.createModule("./oteltest")
	i.stabilize(".")

	i.prune("./example", "./exporters", "./sdk", "./internal/tools", "./bridge").dot()
}
