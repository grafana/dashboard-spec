package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/trotttrotttrott/dashboard-spec/pkg/jsonnet"
)

func main() {
	const lvar = "GDS_GEN_LANG"
	l, exists := os.LookupEnv(lvar)
	if !exists {
		log.Fatalf("Set `%s` environment variable to indicate which language you'd like to generate models for.", lvar)
	}
	spec := loadSpec("bundle/7.0/spec.json")
	s := spec["components"].(map[string]interface{})["schemas"]
	switch l {
	case "jsonnet":
		jsonnet.Generate(s)
	default:
		log.Fatalf("Unsupported language: %s=%s.", lvar, l)
	}
}

func loadSpec(file string) (s map[string]interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	if err := json.Unmarshal(data, &s); err != nil {
		log.Fatalln(err)
	}
	return
}
