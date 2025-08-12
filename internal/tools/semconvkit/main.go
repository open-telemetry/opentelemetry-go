// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package semconvkit is used to generate opentelemetry-go specific semantic
// convention code.
package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"

	"go.opentelemetry.io/otel/internal/tools/semconvkit/decls"
)

var (
	logLevel   = flag.String("log-level", "", `Logging level ("debug", "info", "warn", "error")`)
	semconvPkg = flag.String("semconv", "./", "semconv package directory")
	tag        = flag.String("tag", "", "OpenTelemetry tagged version")
	prev       = flag.String("prev", "", "previous semconv version")

	//go:embed templates/*.tmpl
	rootFS embed.FS
)

func main() {
	flag.Parse()

	slog.SetDefault(newLogger(*logLevel))

	if *tag == "" {
		slog.Error("invalid tag", "tag", *tag)
		os.Exit(1)
	}

	sc := &SemanticConventions{TagVer: *tag}

	out := filepath.Join(*semconvPkg, *tag)

	// Render all other files before the MIGRATION file. That file needs the
	// full package declaration so it can determine compatibility accurately.
	entries, err := rootFS.ReadDir("templates")
	if err != nil {
		slog.Error("error reading templates", "err", err)
		os.Exit(1)
	}

	for _, entry := range entries {
		if entry.Name() == "MIGRATION.md.tmpl" {
			continue
		}

		src := filepath.Join("templates", entry.Name())
		err := render(src, out, sc)
		if err != nil {
			slog.Error("error rendering template", "err", err, "template", entry.Name())
			os.Exit(1)
		}
	}

	prevPkg, err := prevVer(*semconvPkg, *tag, *prev)
	if err != nil {
		slog.Error("previous version not found, skipping migration", "err", err)
		os.Exit(1)
	}

	slog.Debug("previous version found", "prev", prevPkg)
	m, err := newMigration(out, filepath.Join(*semconvPkg, prevPkg))
	if err != nil {
		slog.Error("error getting migration, skipping", "err", err)
		os.Exit(1)
	}

	if err := render("templates/MIGRATION.md.tmpl", out, m); err != nil {
		slog.Error("error rendering migration template", "err", err)
		os.Exit(1)
	}
}

func newLogger(lvlStr string) *slog.Logger {
	levelVar := new(slog.LevelVar) // Default value of info.
	opts := &slog.HandlerOptions{AddSource: true, Level: levelVar}
	h := slog.NewTextHandler(os.Stderr, opts)
	logger := slog.New(h)

	if lvlStr == "" {
		return logger
	}

	var level slog.Level
	if err := level.UnmarshalText([]byte(lvlStr)); err != nil {
		logger.Error("failed to parse log level", "error", err, "log-level", lvlStr)
	} else {
		levelVar.Set(level)
	}

	return logger
}

// render renders all templates to the dest directory using the data.
func render(src, dest string, data any) error {
	tmpls, err := template.ParseFS(rootFS, src)
	if err != nil {
		return err
	}
	for _, tmpl := range tmpls.Templates() {
		slog.Debug("rendering template", "name", tmpl.Name())
		target := filepath.Join(dest, strings.TrimSuffix(tmpl.Name(), ".tmpl"))
		wr, err := os.Create(target)
		if err != nil {
			return err
		}

		err = tmpl.Execute(wr, data)
		if err != nil {
			return err
		}
	}

	return nil
}

// prevVer returns the previous version of the semantic conventions package.
// It will first check for hint within root and return that value if found. If
// not found, it will find all directories in root with a version name and
// return the version that is less than and closest to the curr version.
func prevVer(root, cur, hint string) (string, error) {
	slog.Debug("prevVer", "root", root, "current", cur, "hint", hint)
	info, err := os.Stat(root)
	if err != nil {
		return "", fmt.Errorf("root directory %q not found: %w", root, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("root %q is not a directory", root)
	}

	if hint != "" {
		sub := filepath.Join(root, hint)
		slog.Debug("looking for hint", "path", sub)
		info, err = os.Stat(sub)
		if err == nil && info.IsDir() {
			return hint, nil
		}
	}

	v, err := semver.NewVersion(cur)
	if err != nil {
		return "", fmt.Errorf("invalid current version %q: %w", cur, err)
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return "", fmt.Errorf("error reading root %q: %w", root, err)
	}

	var prev *semver.Version
	for _, entry := range entries {
		slog.Debug("root entry", "name", entry.Name())
		if !entry.IsDir() {
			continue
		}

		ver, err := semver.NewVersion(entry.Name())
		if err != nil {
			slog.Debug("not a version dir", "name", entry.Name())
			// Ignore errors for non-semver directories.
			continue
		}
		slog.Debug("found version dir", "prev", ver)
		if ver.LessThan(v) && (prev == nil || ver.GreaterThan(prev)) {
			slog.Debug("new previous version", "version", ver)
			prev = ver
		}
	}

	if prev == nil {
		return "", errors.New("no previous version found")
	}
	return prev.Original(), nil
}

// SemanticConventions are information about the semantic conventions being
// generated.
type SemanticConventions struct {
	// TagVer is the tagged version (i.e. v.7.0 and not 1.7.0).
	TagVer string
}

func (sc SemanticConventions) SemVer() string {
	return strings.TrimPrefix(*tag, "v")
}

// Migration contains the details about the migration from the previous
// semantic conventions to the current one.
type Migration struct {
	CurVer   string
	PrevVer  string
	Removals []string
	Renames  []Rename
}

// Remove is a semantic convention declaration that has been renamed.
type Rename struct {
	Old, New string
}

func newMigration(cur, prev string) (*Migration, error) {
	cDecl, err := decls.GetNames(cur, parse)
	if err != nil {
		return nil, fmt.Errorf("error parsing current version %q: %w", cur, err)
	}

	pDecl, err := decls.GetNames(prev, parse)
	if err != nil {
		return nil, fmt.Errorf("error parsing previous version %q: %w", prev, err)
	}

	m := Migration{
		CurVer:   filepath.Base(cur),
		PrevVer:  filepath.Base(prev),
		Removals: inAnotB(pDecl, cDecl),
		Renames:  renames(pDecl, cDecl),
	}

	sort.Strings(m.Removals)
	sort.Slice(m.Renames, func(i, j int) bool {
		return m.Renames[i].Old < m.Renames[j].Old
	})

	return &m, nil
}

func parse(d ast.Decl) []string {
	var out []string
	switch decl := d.(type) {
	case *ast.FuncDecl:
		out = []string{decl.Name.Name}
	case *ast.GenDecl:
		if decl.Tok == token.CONST || decl.Tok == token.VAR {
			for _, spec := range decl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						out = append(out, name.Name)
					}
				}
			}
		}
	}
	return out
}

// inAnotB returns the canonical names in a that are not in b.
func inAnotB(a, b decls.Names) []string {
	var diff []string
	for key, name := range a {
		if _, ok := b[key]; !ok {
			diff = append(diff, string(name))
		}
	}
	return diff
}

// renames returns the renames between the old and current names.
func renames(old, current decls.Names) []Rename {
	var renames []Rename
	for key, name := range old {
		if otherName, ok := current[key]; ok && name != otherName {
			renames = append(renames, Rename{
				Old: string(name),
				New: string(otherName),
			})
		}
	}
	return renames
}
