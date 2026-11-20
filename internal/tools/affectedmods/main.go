// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Affectedmods filters a list of benchmark shards to those whose code
// changed in the current diff.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var (
	baseRef  = flag.String("base", "main", "git ref to diff against")
	forceAll = flag.Bool("all", false, "emit every input shard without diffing")
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("affectedmods: ")
	flag.Parse()

	shards, err := readShards(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if *forceAll {
		emit(shards)
		return
	}

	changed, err := changedFiles(*baseRef)
	if err != nil {
		log.Fatal(err)
	}
	if len(changed) == 0 {
		emit(nil)
		return
	}

	// Sort longest-first so owningModule returns the first prefix match.
	sort.Slice(shards, func(i, j int) bool { return len(shards[i]) > len(shards[j]) })

	affected := map[string]bool{}
	for _, f := range changed {
		if !strings.HasSuffix(f, ".go") {
			continue
		}
		if strings.HasPrefix(f, "internal/tools/") {
			continue // tooling has its own go.mod, never a benchmark target
		}
		if shard := owningModule(f, shards); shard != "" {
			affected[shard] = true
		}
	}

	out := make([]string, 0, len(affected))
	for s := range affected {
		out = append(out, s)
	}
	sort.Strings(out)
	emit(out)
}

// readShards decodes a JSON array of shard labels from stdin.
func readShards(r *os.File) ([]string, error) {
	var shards []string
	if err := json.NewDecoder(r).Decode(&shards); err != nil {
		return nil, fmt.Errorf("decoding shards: %w", err)
	}
	if len(shards) == 0 {
		return nil, fmt.Errorf("no shards on stdin")
	}
	return shards, nil
}

func changedFiles(base string) ([]string, error) {
	if err := exec.Command("git", "rev-parse", "--verify", base).Run(); err != nil {
		return nil, fmt.Errorf("base ref %q not found", base)
	}
	out, err := exec.Command("git", "diff", "--name-only", base+"...HEAD").Output()
	if err != nil {
		return nil, fmt.Errorf("git diff: %w", err)
	}
	var files []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// owningModule returns the longest shard-prefix containing p. Inputs
// are like "./sdk/log" or "./sdk/metric/." and must be sorted
// longest-first. A trailing "/." on a shard label is stripped for
// matching -- it's a marker for the catch-all sub-shard of a split
// module. If no nested shard matches, the root module "." is returned
// when present in modules.
func owningModule(p string, modules []string) string {
	root := ""
	for _, m := range modules {
		if m == "." {
			root = "."
			continue
		}
		bare := strings.TrimPrefix(m, "./")
		bare = strings.TrimSuffix(bare, "/.")
		if p == bare || strings.HasPrefix(p, bare+"/") {
			return m
		}
	}
	return root
}

func emit(shards []string) {
	if shards == nil {
		shards = []string{}
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(shards); err != nil {
		log.Fatal(err)
	}
}
