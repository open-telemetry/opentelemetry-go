// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package decls provides a set of functions to parse and analyze Go source
// code and get the declarations within it.
package decls // import "go.opentelemetry.io/otel/internal/tools/semconvkit/decls"

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// GetNames parses the Go source code in the specified package path and returns
// the names extracted from the declarations using the provided parser
// function.
//
// The names are returned as a map where the keys are the names fully
// lowercased form of the name and the values are the original format of the
// name.
func GetNames(pkgPath string, f Parser) (Names, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, nil, 0)
	if err != nil {
		return nil, err
	}

	out := make(Names)
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				for _, name := range f(decl) {
					out[NewCanonicalName(name)] = Name(name)
				}
			}
		}
	}
	return out, nil
}

// Parser is a function type that takes an [ast.Decl] and returns a slice of
// parsed string identifiers.
type Parser func(ast.Decl) []string

// CanonicalName is the canonical form of a name (lowercase).
type CanonicalName string

// NewCanonicalName returns name as a [CanonicalName].
func NewCanonicalName(name string) CanonicalName {
	return CanonicalName(strings.ToLower(name))
}

// Name is the original form of a name (case-sensitive).
type Name string

// Names is a map of canonical names to their original names.
type Names map[CanonicalName]Name
