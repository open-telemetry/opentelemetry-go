// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewMigrationIgnoresUnexportedDeclarations(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	goMod := filepath.Join(root, "go.mod")
	if err := os.WriteFile(
		goMod,
		[]byte("module example.com/semconvtest\n\ngo 1.25.0\n"),
		0o600,
	); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", goMod, err)
	}

	prev := filepath.Join(root, "v1.40.0")
	cur := filepath.Join(root, "v1.41.0")

	writeGoFile(t, prev, "semconv", `package semconv

func ErrorType(err error) {}
func DeploymentEnvironmentName(val string) {}
`)
	writeGoFile(t, cur, "semconv", `package semconv

func ErrorType(err error) {}
func errorType(err error) {}
`)

	m, err := newMigration(cur, prev)
	if err != nil {
		t.Fatalf("newMigration() error = %v", err)
	}

	if len(m.Renames) != 0 {
		t.Fatalf("newMigration() renames = %#v, want none", m.Renames)
	}

	wantRemoval := []string{"DeploymentEnvironmentName"}
	if len(m.Removals) != len(wantRemoval) || m.Removals[0] != wantRemoval[0] {
		t.Fatalf("newMigration() removals = %#v, want %#v", m.Removals, wantRemoval)
	}
}

// TestNoInvalidObservableHistogramTypes asserts that the semconv code
// generator never emits types containing "HistogramObservable" or
// "ObservableObservable", which would indicate the histogram-exclusion
// rule in metric.go.j2 was accidentally removed or bypassed.
func TestNoInvalidObservableHistogramTypes(t *testing.T) {
	t.Parallel()

	semconvDir := filepath.Join("..", "..", "..", "semconv")

	forbidden := [][]byte{
		[]byte("HistogramObservable"),
		[]byte("ObservableObservable"),
	}

	err := filepath.WalkDir(semconvDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() ||
			filepath.Base(path) != "metric.go" ||
			!strings.HasSuffix(filepath.Base(filepath.Dir(path)), "conv") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for _, term := range forbidden {
			if bytes.Contains(data, term) {
				t.Errorf("%s: contains forbidden substring %q", path, term)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir(%q) error = %v", semconvDir, err)
	}
}

// TestGeneratedMetricHelpersReturnAfterEmptyAttrs checks every retained helper
// affected by the missing return in semconv v1.32.0 through v1.34.0.
func TestGeneratedMetricHelpersReturnAfterEmptyAttrs(t *testing.T) {
	t.Parallel()

	packages := map[string][]string{
		"containerconv": {"Uptime"},
		"goconv":        {"ScheduleDuration"},
		"k8sconv":       {"NodeCPUUsage", "NodeMemoryUsage", "NodeUptime", "PodCPUUsage", "PodMemoryUsage", "PodUptime"},
		"processconv":   {"Uptime"},
		"rpcconv":       {"ClientDuration", "ClientRequestSize", "ClientRequestsPerRPC", "ClientResponseSize", "ClientResponsesPerRPC", "ServerDuration", "ServerRequestSize", "ServerRequestsPerRPC", "ServerResponseSize", "ServerResponsesPerRPC"},
		"systemconv":    {"Uptime"},
	}
	versions := []string{"v1.32.0", "v1.33.0", "v1.34.0"}
	fset := token.NewFileSet()

	for _, version := range versions {
		for pkg, functions := range packages {
			path := filepath.Join("..", "..", "..", "semconv", version, pkg, "metric.go")
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				t.Fatalf("ParseFile(%q) error = %v", path, err)
			}

			for _, function := range functions {
				matches := 0
				ast.Inspect(file, func(node ast.Node) bool {
					fn, ok := node.(*ast.FuncDecl)
					if !ok || fn.Recv == nil || fn.Name.Name != "Record" {
						return true
					}
					if len(fn.Recv.List) != 1 || len(fn.Recv.List[0].Names) != 1 || fn.Recv.List[0].Names[0].Name != "m" {
						return true
					}
					receiver, ok := fn.Recv.List[0].Type.(*ast.Ident)
					if !ok || receiver.Name != function {
						return true
					}

					fastPaths := 0
					ast.Inspect(fn.Body, func(node ast.Node) bool {
						ifStmt, ok := node.(*ast.IfStmt)
						if !ok {
							return true
						}
						condition, ok := ifStmt.Cond.(*ast.BinaryExpr)
						if !ok || condition.Op != token.EQL {
							return true
						}
						call, ok := condition.X.(*ast.CallExpr)
						zero, zeroOK := condition.Y.(*ast.BasicLit)
						if !ok || !zeroOK || zero.Kind != token.INT || zero.Value != "0" || len(call.Args) != 1 {
							return true
						}
						fn, fnOK := call.Fun.(*ast.Ident)
						arg, argOK := call.Args[0].(*ast.Ident)
						if !fnOK || !argOK || fn.Name != "len" || arg.Name != "attrs" {
							return true
						}
						if len(ifStmt.Body.List) > 0 {
							if _, ok := ifStmt.Body.List[len(ifStmt.Body.List)-1].(*ast.ReturnStmt); ok {
								fastPaths++
							}
						}
						return true
					})
					if fastPaths == 1 {
						matches++
					}
					return false
				})
				if matches != 1 {
					t.Errorf("%s: %s.Record empty-attrs fast paths = %d, want 1", path, function, matches)
				}
			}
		}
	}
}

func writeGoFile(t *testing.T, dir, _, src string) {
	t.Helper()

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", dir, err)
	}

	modPath := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(modPath, []byte("module example.com/semconvtest\n\ngo 1.25.0\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", modPath, err)
	}

	path := filepath.Join(dir, "decls.go")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}
