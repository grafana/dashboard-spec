package main

import (
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-openapi/inflect"
)

func generateLibrary(s Spec, l Language) error {

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
				"singularize":      inflect.Singularize,
				"add": func(x int, y int) int {
					return x + y
				},
				"subtract": func(x int, y int) int {
					return x - y
				},
				"repeat": func(s string, n int) string {
					return strings.Repeat(s, n)
				},
			},
		).ParseFiles(
			path.Join("templates", l.Directory, "library", tmplFile),
			path.Join("templates", l.Directory, "library", "_shared.tmpl"),
		)
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

	// Generate main.
	err = g("grafana", "main", s.Components.Schemas)
	if err != nil {
		return err
	}

	// Generate docs.
	tmpl, err := template.New("docs.tmpl").Funcs(
		template.FuncMap{
			"objectInflection": l.OjectInflection,
			"singularize":      inflect.Singularize,
			"inflectJoin": func(elems ...string) (s string) {
				for i, e := range elems {
					elems[i] = l.OjectInflection(e)
				}
				return strings.Join(elems, ".")
			},
			"indent": func(spaces int, s string) string {
				pad := strings.Repeat(" ", spaces)
				return pad + strings.Replace(s, "\n", "\n"+pad, -1)
			},
		},
	).ParseFiles(path.Join("templates", "docs.tmpl"))
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(dir, "DOCS.md"))
	if err != nil {
		return err
	}
	return tmpl.Execute(f, s.Components.Schemas)
}
