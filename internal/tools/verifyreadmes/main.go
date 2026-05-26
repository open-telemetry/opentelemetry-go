// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package verifyreadmes is used to verify that all go modules in the repository
// have a README.md file.
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// excludedDirs is a list of directories to exclude from the README check if the full path contains any of these strings.
var excludedDirs = []string{
	"internal",
	"test",
	"example",
	"/.",
}

const readmeFilename = "README.md"

// verifyReadme is a [fs.WalkDirFunc] that checks if a README.md exists in the same directory as the go.mod file.
func verifyReadme(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !d.Type().IsRegular() || d.Name() != "go.mod" {
		return nil
	}

	for _, dir := range excludedDirs {
		if strings.Contains(path, dir) {
			return nil
		}
	}

	// Check that a README.md exists in the same directory as the go.mod file.
	readme := filepath.Join(filepath.Dir(path), readmeFilename)
	_, err = os.Stat(readme)
	if os.IsNotExist(err) {
		return fmt.Errorf("couldn't find %s for %q", readmeFilename, filepath.Dir(path))
	}

	return err
}

func main() {
	root, err := os.Getwd()
	if len(os.Args) == 2 {
		root, err = filepath.Abs(os.Args[1])
	}

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// Clean the path to prevent path traversal issues
	root = filepath.Clean(root)

	fmt.Println("Verifying READMEs in", root)

	var errs []error
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err := verifyReadme(path, d, err); err != nil {
			errs = append(errs, err)
		}
		return nil // continue walking
	})

	if err := errors.Join(errs...); err != nil {
		fmt.Println("Some readme files couldn't be found.")
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
