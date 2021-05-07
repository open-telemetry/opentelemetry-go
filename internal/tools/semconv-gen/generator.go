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

var outputFn = flag.StringP("filename", "f", "", "Filename for template output. If not specified 'basename(inputPath).go' will be used.")
var inputPath = flag.StringP("input", "i", "", "Path to semantic convention definition YAML")
var outputPath = flag.StringP("output", "o", "semconv", "Path to output target. Must be either an absolute path or relative to the repository root.")
var templateFn = flag.StringP("template", "t", "template.j2", "Template filename")
var containerImage = flag.StringP("container", "c", "otel/semconvgen", "Container image ID")

func main() {
	flag.Parse()

	cfg, err := buildConfig()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(-1)
	}

	err = cfg.render()
	if err != nil {
		panic(err)
	}

	err = fixInitialisms(path.Join(cfg.output, cfg.target))
	if err != nil {
		panic(err)
	}
}

type config struct {
	input string
	output string
	target string
	template string
	image string
}

func buildConfig() (*config, error) {
	if *inputPath == "" {
		return nil, errors.New("input path must be provided")
	}

	if *outputFn == "" {
		*outputFn = fmt.Sprintf("%s.go", path.Base(*inputPath))
	}

	if !path.IsAbs(*outputPath) {
		root, err := findRepoRoot()
		if err != nil {
			return nil, err
		}
		*outputPath = path.Join(root, *outputPath)
	}

	if !path.IsAbs(*templateFn) {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		*templateFn = path.Join(pwd, *templateFn)
	}

	return &config{
		input: *inputPath,
		output: *outputPath,
		target: *outputFn,
		template: *templateFn,
		image: *containerImage,
	}, nil
}

func (cfg config) render() error {
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

	err = exec.Command("cp", "-a", cfg.input, path.Join(tmpDir, "input")).Run()
	if err != nil {
		return fmt.Errorf("unable to copy input to temp directory: %w", err)
	}

	err = exec.Command("cp", cfg.template, tmpDir).Run()
	if err != nil {
		return fmt.Errorf("unable to copy template to temp directory: %w", err)
	}

	cmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data", tmpDir),
		cfg.image,
		"--yaml-root", path.Join("/data/input", path.Base(cfg.input)),
		"code",
		"--template", path.Join("/data", path.Base(cfg.template)),
		"--output", path.Join("/data/output", path.Base(cfg.target)),
	)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to render template: %w", err)
	}

	err = exec.Command("cp", path.Join(tmpDir, "output", path.Base(cfg.target)), cfg.output).Run()
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

var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DB":    true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"OS":    true,
	"PID":   true,
	"QPS":   true,
	"QUIC":  true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SDK":   true,
	"SLA":   true,
	"SMTP":  true,
	"SPDY":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

var replacements = map[string]string{
	"Inproc": "InProc",
	"IPTCP": "TCP",
	"IPUDP": "UDP",
}

func fixInitialisms(fn string) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}

	for init := range commonInitialisms {
		re := regexp.MustCompile(strings.Title(strings.ToLower(init)))
		data = re.ReplaceAllLiteral(data, []byte(init))
	}

	for old, new := range replacements {
		data = bytes.ReplaceAll(data, []byte(old), []byte(new))
	}

	err = ioutil.WriteFile(fn, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write updated file: %w", err)
	}

	cmd := exec.Command("gofmt", "-w", "-s", fn)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to format updated file: %w", err)
	}

	return nil
}