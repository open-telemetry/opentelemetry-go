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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

func main() {
	cfg := config{}
	flag.StringVarP(&cfg.inputPath, "input", "i", "", "Path to semantic convention definition YAML")
	flag.StringVarP(&cfg.outputPath, "output", "o", "semconv", "Path to output target. Must be either an absolute path or relative to the repository root.")
	flag.StringVarP(&cfg.containerImage, "container", "c", "otel/semconvgen", "Container image ID")
	flag.StringVarP(&cfg.outputFilename, "filename", "f", "", "Filename for templated output. If not specified 'basename(inputPath).go' will be used.")
	flag.StringVarP(&cfg.templateFilename, "template", "t", "template.j2", "Template filename")
	flag.Parse()

	cfg, err := validateConfig(cfg)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(-1)
	}

	err = render(cfg)
	if err != nil {
		panic(err)
	}

	err = fixIdentifiers(cfg.outputFilename)
	if err != nil {
		panic(err)
	}

	err = format(cfg.outputFilename)
	if err != nil {
		panic(err)
	}
}

type config struct {
	inputPath        string
	outputPath       string
	outputFilename   string
	templateFilename string
	containerImage   string
}

func validateConfig(cfg config) (config, error) {
	if cfg.inputPath == "" {
		return config{}, errors.New("input path must be provided")
	}

	if cfg.outputFilename == "" {
		cfg.outputFilename = fmt.Sprintf("%s.go", path.Base(cfg.inputPath))
	}

	if !path.IsAbs(cfg.outputPath) {
		root, err := findRepoRoot()
		if err != nil {
			return config{}, err
		}
		cfg.outputPath = path.Join(root, cfg.outputPath)
	}

	cfg.outputFilename = path.Join(cfg.outputPath, cfg.outputFilename)

	if !path.IsAbs(cfg.templateFilename) {
		pwd, err := os.Getwd()
		if err != nil {
			return config{}, err
		}
		cfg.templateFilename = path.Join(pwd, cfg.templateFilename)
	}

	return cfg, nil
}

func render(cfg config) error {
	tmpDir, err := os.MkdirTemp("", "otel_semconvgen")
	if err != nil {
		return fmt.Errorf("unable to create temporary directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	inputPath := path.Join(tmpDir, "input")
	err = os.Mkdir(inputPath, 0700)
	if err != nil {
		return fmt.Errorf("unable to create input directory: %w", err)
	}

	outputPath := path.Join(tmpDir, "output")
	err = os.Mkdir(outputPath, 0700)
	if err != nil {
		return fmt.Errorf("unable to create output directory: %w", err)
	}

	err = exec.Command("cp", "-a", cfg.inputPath, inputPath).Run()
	if err != nil {
		return fmt.Errorf("unable to copy input to temp directory: %w", err)
	}

	err = exec.Command("cp", cfg.templateFilename, tmpDir).Run()
	if err != nil {
		return fmt.Errorf("unable to copy template to temp directory: %w", err)
	}

	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data", tmpDir),
		cfg.containerImage,
		"--yaml-root", path.Join("/data/input", path.Base(cfg.inputPath)),
		"code",
		"--template", path.Join("/data", path.Base(cfg.templateFilename)),
		"--output", path.Join("/data/output", path.Base(cfg.outputFilename)),
	)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to render template: %w", err)
	}

	err = exec.Command("cp", path.Join(tmpDir, "output", path.Base(cfg.outputFilename)), cfg.outputPath).Run()
	if err != nil {
		return fmt.Errorf("unable to copy result to target: %w", err)
	}

	return nil
}

func findRepoRoot() (string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := start
	for {
		_, err := os.Stat(filepath.Join(dir, ".git"))
		if errors.Is(err, os.ErrNotExist) {
			dir = filepath.Dir(dir)
			// From https://golang.org/pkg/path/filepath/#Dir:
			// The returned path does not end in a separator unless it is the root directory.
			if strings.HasSuffix(dir, string(filepath.Separator)) {
				return "", fmt.Errorf("unable to find git repository enclosing working dir %s", start)
			}
			continue
		}

		if err != nil {
			return "", err
		}

		return dir, nil
	}
}

var capitalizations = []string{
	"ACL",
	"AIX",
	"AKS",
	"AMD64",
	"API",
	"ARM32",
	"ARM64",
	"ARN",
	"ARNs",
	"ASCII",
	"AWS",
	"CPU",
	"CSS",
	"DB",
	"DC",
	"DNS",
	"EC2",
	"ECS",
	"EDB",
	"EKS",
	"EOF",
	"GCP",
	"GRPC",
	"GUID",
	"HPUX",
	"HSQLDB",
	"HTML",
	"HTTP",
	"HTTPS",
	"IA64",
	"ID",
	"IP",
	"JDBC",
	"JSON",
	"K8S",
	"LHS",
	"MSSQL",
	"OS",
	"PHP",
	"PID",
	"PPC32",
	"PPC64",
	"QPS",
	"QUIC",
	"RAM",
	"RHS",
	"RPC",
	"SDK",
	"SLA",
	"SMTP",
	"SPDY",
	"SQL",
	"SSH",
	"TCP",
	"TLS",
	"TTL",
	"UDP",
	"UID",
	"UI",
	"UUID",
	"URI",
	"URL",
	"UTF8",
	"VM",
	"XML",
	"XMPP",
	"XSRF",
	"XSS",
	"ZOS",
	"CronJob",
	"WebEngine",
	"MySQL",
	"PostgreSQL",
	"MariaDB",
	"MaxDB",
	"FirstSQL",
	"InstantDB",
	"HBase",
	"MongoDB",
	"CouchDB",
	"CosmosDB",
	"DynamoDB",
	"HanaDB",
	"FreeBSD",
	"NetBSD",
	"OpenBSD",
	"DragonflyBSD",
	"InProc",
	"FaaS",
}

// These are not simple capitalization fixes, but require string replacement.
// All occurrences of the key will be replaced with the corresponding value.
var replacements = map[string]string{
	"RedisDatabase": "RedisDB",
	"IPTCP":         "TCP",
	"IPUDP":         "UDP",
	"Lineno":        "LineNumber",
}

func fixIdentifiers(fn string) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	for _, init := range capitalizations {
		// Match the title-cased capitalization target, asserting that its followed by
		// either a capital letter, whitespace, a digit, or the end of text.
		// This is to avoid, e.g., turning "Identifier" into "IDentifier".
		re := regexp.MustCompile(strings.Title(strings.ToLower(init)) + `([A-Z\s\d]|$)`)
		// RE2 does not support zero-width lookahead assertions, so we have to replace
		// the last character that may have matched the first capture group in the
		// expression constructed above.
		data = re.ReplaceAll(data, []byte(init+`$1`))
	}

	for cur, repl := range replacements {
		data = bytes.ReplaceAll(data, []byte(cur), []byte(repl))
	}

	err = ioutil.WriteFile(fn, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write updated file: %w", err)
	}

	return nil
}

func format(fn string) error {
	cmd := exec.Command("gofmt", "-w", "-s", fn)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to format updated file: %w", err)
	}

	return nil
}
