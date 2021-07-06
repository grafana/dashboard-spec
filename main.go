package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/go-openapi/inflect"
)

type Inflection func(string) string

// Language holds metadata for a supported language.
type Language struct {
	Directory          string
	FileExtension      string
	FileNameInflection Inflection
	OjectInflection    Inflection
	Type               int
}

func main() {
	args := os.Args[1:]
	action, specVersion, language := args[0], args[1], args[2]
	if action == "instance" {
		instanceFile := args[3]
		err := generateInstance(
			loadSpec(path.Join("_gen", specVersion, "spec.json")),
			loadLanguage(language),
			instanceFile,
		)
		if err != nil {
			log.Fatalln(err)
		}

	} else if action == "library" {
		err := generateLibrary(
			loadSpec(path.Join("_gen", specVersion, "spec.json")),
			loadLanguage(language),
		)
		if err != nil {
			log.Fatalln(err)
		}
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

func loadLanguage(l string) Language {
	languages := map[string]Language{
		"jsonnet": {
			Directory:          "jsonnet",
			FileExtension:      "libsonnet",
			FileNameInflection: inflect.CamelizeDownFirst,
			OjectInflection:    inflect.CamelizeDownFirst,
		},
	}
	lang, ok := languages[l]
	if !ok {
		log.Fatalf("%q is not a supported language.", l)
	}
	return lang
}
