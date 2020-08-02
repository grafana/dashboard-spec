package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/go-openapi/inflect"
)

type Inflection func(string) string

// Metadata for a supported language.
type Language struct {
	Directory          string
	FileExtension      string
	FileNameInflection Inflection
	OjectInflection    Inflection
}

func main() {
	const lvar = "GDS_GEN_LANG"
	lang, exists := os.LookupEnv(lvar)
	if !exists {
		log.Fatalf("Set `%s` environment variable to indicate which language you'd like to generate models for.", lvar)
	}
	l := map[string]Language{
		"jsonnet": {
			Directory:          "jsonnet",
			FileExtension:      "libsonnet",
			FileNameInflection: inflect.CamelizeDownFirst,
			OjectInflection:    inflect.CamelizeDownFirst,
		},
	}[lang]
	s := loadSpec("bundle/7.0/spec.json")
	err := generate(l, s)
	if err != nil {
		log.Fatalln(err)
	}
}

func loadSpec(file string) (s Spec) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(data, &s); err != nil {
		log.Fatalln(err)
	}
	return
}

func generate(l Language, s Spec) error {

	// Create directories.
	dir := path.Join("_gen", s.Info.Version, l.Directory)
	os.MkdirAll(dir, os.ModePerm)
	for _, d := range []string{"panel", "target", "template"} {
		os.MkdirAll(path.Join(dir, d), os.ModePerm)
	}

	// Function that renders templates and writes them to files.
	g := func(name string, tmplType string, data interface{}) error {
		tmplFile := fmt.Sprintf("%s.tmpl", tmplType)
		tmpl, err := template.New(tmplFile).Funcs(
			template.FuncMap{
				"objectInflection": l.OjectInflection,
			},
		).ParseFiles(path.Join("templates", l.Directory, tmplFile))
		if err != nil {
			return err
		}
		fileName := fmt.Sprintf("%s.%s", l.FileNameInflection(name), l.FileExtension)
		dest := dir
		if tmplType != "dashboard" && tmplType != "main" {
			dest = path.Join(dir, tmplType)
		}
		f, err := os.Create(path.Join(dest, fileName))
		if err != nil {
			return err
		}
		return tmpl.Execute(f, data)
	}

	// Generate dashboard file.
	err := g("dashboard", "dashboard", s.Components.Schemas["Dashboard"])
	if err != nil {
		return err
	}

	// Generate panel files.
	for n, sc := range s.Components.Schemas["Panel"].Properties {
		err = g(n, "panel", *sc)
		if err != nil {
			return err
		}
	}

	// Generate target files.
	for n, sc := range s.Components.Schemas["Target"].Properties {
		err = g(n, "target", *sc)
		if err != nil {
			return err
		}
	}

	// Generate template files.
	for n, sc := range s.Components.Schemas["Template"].Properties {
		err = g(n, "template", *sc)
		if err != nil {
			return err
		}
	}

	return g("grafana", "main", s.Components.Schemas)
}
