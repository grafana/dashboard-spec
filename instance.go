package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/go-openapi/inflect"
)

func generateInstance(s Spec, l Language, instanceFile string) error {

	// Create directories.
	dir := path.Join("_gen", s.Info.Version, l.Directory)
	os.MkdirAll(dir, os.ModePerm)
	for _, d := range []string{"panel", "target", "template"} {
		os.MkdirAll(path.Join(dir, d), os.ModePerm)
	}

	// Function that renders templates and writes them to files.
	//g := func(name string, tmplType string, data interface{}) error {

	dashboardString, err := ioutil.ReadFile(instanceFile)
	if err != nil {
		return err
	}
	dashboard := map[string]interface{}{}
	err = json.Unmarshal(dashboardString, &dashboard)
	if err != nil {
		return err
	}

	outFile := strings.ReplaceAll(instanceFile, ".json", "") + ".libsonnet"

	f, err := os.Create(outFile)
	if err != nil {
		return err
	}

	tmplFile := "dashboard.tmpl"
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
			"toString": func(value interface{}, valueType, name string) (string, error) {
				switch valueType {
				case "string":
					if value == nil {
						return "", nil
					}
					fmt.Fprintf(os.Stderr, "%s: %v\n", name, value)
					return fmt.Sprintf("\"%v\"", value.(string)), nil
				case "integer":
					return fmt.Sprintf("%v", value), nil
				case "boolean":
					return fmt.Sprintf("%v", value), nil
				default:
					return "", nil //fmt.Errorf("Cannot convert %v (%s) to string", v, reflect.TypeOf(v))
				}
			},
			"toKey": func(s string) string {
				return inflect.CamelizeDownFirst(s)
			},
			"toCamel": func(s string) string {
				return inflect.Camelize(s)
			},
			"debug": func(message string, value interface{}) string {
				message = fmt.Sprintf("DEBUG: %s %v\n", message, value)
				fmt.Fprintf(os.Stderr, message)
				return message
			},
			"debugx": func(message string, value interface{}) string {
				b, _ := json.MarshalIndent(value, "  ", "  ")
				message = fmt.Sprintf("DEBUG: %s %v\n", message, string(b))
				fmt.Fprintf(os.Stderr, message)
				return ""
			}, // wrap overcomes Golang template's inability to pass in multiple arguments to a template
			"wrap": func(values ...interface{}) []interface{} {
				return values
			},
		},
	).ParseFiles(
		path.Join("templates", l.Directory, "instance", tmplFile),
		path.Join("templates", l.Directory, "instance", "_shared.tmpl"),
	)
	if err != nil {
		return err
	}
	data := map[string]interface{}{}
	data["spec"] = s.Components.Schemas
	data["instance"] = dashboard

	err = tmpl.Execute(f, data)
	return err
}
