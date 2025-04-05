// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package verifyreadmes is used to verify that all go modules in the repository
// have a README.md file.
package main

import (
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

// verifyReadme is a [os.WalkFunc] that checks if a README.md exists in the same directory as the go.mod file.
func verifyReadme(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() || info.Name() != "go.mod" {
		return nil

	}

	for _, dir := range excludedDirs {
		if strings.Contains(path, dir) {
			return nil
		}
	}

	// Check that a README.md exists in the same directory as the go.mod file.
	readme := filepath.Join(filepath.Dir(path), "README.md")
	_, err = os.Stat(readme)
	if os.IsNotExist(err) {
		err = fmt.Errorf("couldn't find README.md for %q", filepath.Dir(path))
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

	fmt.Println("Verifying READMEs in", root)

	var errs []string
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		err = verifyReadme(path, info, err)
		if err != nil {
			errs = append(errs, err.Error())
		}
		return nil // continue walking
	})

	if len(errs) == 0 {
		return
	}

	fmt.Println("Error: some READMEs couldn't be found.")
	fmt.Println(strings.Join(errs, "\n"))
	os.Exit(1)
}
