package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"text/template"
)

var (
	out = flag.String("output", "./", "output directory")
	tag = flag.String("tag", "", "OpenTelemetry tagged version")
)

// SemanticConventions are information about the semantic conventions being
// generated.
type SemanticConventions struct {
	// SemVer is the semantic version (i.e. 1.7.0 and not v1.7.0).
	SemVer string
	// TagVer is the tagged version (i.e. v1.7.0 and not 1.7.0).
	TagVer string
}

func render(dest string, sc *SemanticConventions) error {
	const templateDir = "templates/"

	f, err := os.Open(templateDir)
	if err != nil {
		return err
	}
	files, err := f.Readdir(0)
	if err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(templateDir, file.Name())
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		target := filepath.Join(dest, strings.TrimSuffix(file.Name(), ".tmpl"))
		wr, err := os.Create(target)
		if err != nil {
			return err
		}

		tmpl.Execute(wr, sc)
	}

	return nil
}

func main() {
	flag.Parse()

	if *tag == "" {
		log.Fatalf("invalid tag: %q", *tag)
	}

	sc := &SemanticConventions{
		SemVer: strings.TrimPrefix(*tag, "v"),
		TagVer: *tag,
	}

	if err := render(*out, sc); err != nil {
		log.Fatal(err)
	}
}
