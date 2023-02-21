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
	"encoding/xml"
	"html/template"
	"log"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// File path to UCUM specification XML.
const specification = "./specification/2.1/ucum-essence.xml"

// Only units listed here are rendered.
var unitAllowList = map[string]bool{
	"bit":            true,
	"byte":           true,
	"baud":           true,
	"hertz":          true,
	"watt":           true,
	"volt":           true,
	"degree Celsius": true,
}

// renderManifest si the mapping from source template to target Go file.
var renderManifest = []struct {
	Source *template.Template
	Dest   string
}{
	{
		Source: parseTmpl("./prefix.tmpl"),
		Dest:   "../prefix.go",
	},
	{
		Source: parseTmpl("./prefix_test.tmpl"),
		Dest:   "../prefix_test.go",
	},
	{
		Source: parseTmpl("./base.tmpl"),
		Dest:   "../base.go",
	},
	{
		Source: parseTmpl("./metric.tmpl"),
		Dest:   "../metric.go",
	},
}

var funcMap = template.FuncMap{
	"Title": func(s string) string {
		caser := cases.Title(language.AmericanEnglish)
		return strings.ReplaceAll(caser.String(s), " ", "")
	},
	"AllowedMetricUnits": func(r *Root) []Unit {
		var allowed []Unit
		for _, u := range r.Units {
			if _, ok := unitAllowList[u.Name]; ok && u.IsMetric == "yes" {
				allowed = append(allowed, u)
			}
		}
		return allowed
	},
}

func parseTmpl(src string) *template.Template {
	name := path.Base(src)
	return template.Must(
		template.New(name).Funcs(funcMap).ParseFiles(src),
	)
}

type UnparsedElementMap map[string][]string

type element struct {
	XMLName    xml.Name
	Attributes []string `xml:",any,attr"`
	Content    string   `xml:",innerxml"`
}

func (u *UnparsedElementMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t := element{}
	if err := d.DecodeElement(&t, &start); err != nil {
		return err
	}
	if *u == nil {
		*u = UnparsedElementMap{}
	}
	(*u)[t.XMLName.Local] = append((*u)[t.XMLName.Local], t.Content)
	return nil
}

type Root struct {
	XMLName      xml.Name   `xml:"root"`
	Namespace    string     `xml:"xmlns,attr"`
	Version      string     `xml:"version,attr"`
	Revision     string     `xml:"revision,attr"`
	RevisionDate string     `xml:"revision-date,attr"`
	Prefixes     []Prefix   `xml:"prefix"`
	BaseUnits    []BaseUnit `xml:"base-unit"`
	Units        []Unit     `xml:"unit"`

	UnparsedAttributes []string           `xml:",any,attr"`
	UnparsedElements   UnparsedElementMap `xml:",any"`
}

type Value struct {
	XMLName             xml.Name `xml:"value"`
	CaseSensitiveUnit   string   `xml:"Unit,attr"`
	CaseInsensitiveUnit string   `xml:"UNIT,attr"`
	Value               string   `xml:"value,attr"` // nolint:revive  // XMLName overlaps.
	Raw                 string   `xml:",chardata"`

	UnparsedAttributes []string           `xml:",any,attr"`
	UnparsedElements   UnparsedElementMap `xml:",any"`
}

type Prefix struct {
	XMLName             xml.Name `xml:"prefix"`
	Name                string   `xml:"name"`
	CaseSensitiveCode   string   `xml:"Code,attr"`
	CaseInsensitiveCode string   `xml:"CODE,attr"`
	PrintSymbol         string   `xml:"printSymbol"`
	Value               Value    `xml:"value"`

	UnparsedAttributes []string           `xml:",any,attr"`
	UnparsedElements   UnparsedElementMap `xml:",any"`
}

type BaseUnit struct {
	XMLName             xml.Name `xml:"base-unit"`
	Name                string   `xml:"name"`
	CaseSensitiveCode   string   `xml:"Code,attr"`
	CaseInsensitiveCode string   `xml:"CODE,attr"`
	Dimension           string   `xml:"dim,attr"`
	PrintSymbol         string   `xml:"printSymbol"`
	Property            string   `xml:"property"`

	UnparsedAttributes []string           `xml:",any,attr"`
	UnparsedElements   UnparsedElementMap `xml:",any"`
}

type Unit struct {
	XMLName             xml.Name `xml:"unit"`
	Name                string   `xml:"name"`
	CaseSensitiveCode   string   `xml:"Code,attr"`
	CaseInsensitiveCode string   `xml:"CODE,attr"`
	IsMetric            string   `xml:"isMetric,attr"`
	IsArbitrary         string   `xml:"isArbitrary,attr"`
	IsSpecial           string   `xml:"isSpecial,attr"`
	Class               string   `xml:"class,attr"`
	PrintSymbol         string   `xml:"printSymbol"`
	Property            string   `xml:"property"`
	Value               Value    `xml:"value"`

	UnparsedAttributes []string           `xml:",any,attr"`
	UnparsedElements   UnparsedElementMap `xml:",any"`
}

func decodedUCUMSpec(spec string) (*Root, error) {
	xmlFile, err := os.Open(spec)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	root := new(Root)
	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(root)
	return root, err
}

func main() {
	root, err := decodedUCUMSpec(specification)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range renderManifest {
		dest, err := os.Create(r.Dest)
		if err != nil {
			log.Fatal(err)
		}
		if err := r.Source.Execute(dest, root); err != nil {
			log.Fatal(err)
		}
	}
}
