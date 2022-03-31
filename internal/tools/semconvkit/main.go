package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"text/template"
)

var (
	out     = flag.String("output", "./", "output directory")
	version = flag.String("version", "", "semantic convention version")

	templates      = template.Must(template.ParseGlob("templates/*"))
	templateToFile = map[string]string{
		"doc.go.tmpl": "doc.go",
	}
)

// SemanticConventions are information about the semantic conventions being
// generated.
type SemanticConventions struct {
	Version string
}

func main() {
	flag.Parse()

	if *version == "" {
		log.Fatalf("invalid version: %q", *version)
	}

	sc := &SemanticConventions{
		Version: *version,
	}

	templates, err := template.ParseGlob("./templates/*")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(*out); os.IsNotExist(err) {
		err = os.Mkdir(*out, os.ModeDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	for tmpl, fName := range templateToFile {
		path := filepath.Join(*out, fName)
		fWriter, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		templates.ExecuteTemplate(fWriter, tmpl, sc)
	}
}
