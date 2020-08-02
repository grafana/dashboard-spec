package main

import (
	"encoding/json"
	"log"
)

// OpenAPI 3.0 spec document.
type Spec struct {
	Version string `json:"openapi"`
	Info    struct {
		Title   string
		Version string
	}
	Components struct {
		Schemas map[string]Schema
	}
}

// OpenAPI 3.0 schema.
type Schema struct {
	Default     interface{}
	Description string
	Items       *Schema
	Properties  map[string]*Schema
	ReadOnly    bool
	Required    []string
	Type        string
}

func (s Schema) MutableProperties() map[string]*Schema {
	p := map[string]*Schema{}
	for n, s := range s.Properties {
		if !s.ReadOnly {
			p[n] = s
		}
	}
	return p
}

func (s Schema) DefaultJSON() string {
	b, err := json.Marshal(s.Default)
	if err != nil {
		log.Fatalln(err)
	}
	return string(b)
}

// Recursively flattens objects in a schema's properties. This is used for
// simplifying the interfaces of objects with many levels of nesting.
func (s Schema) FlattenedNonArrayProperties() map[string]map[string]interface{} {
	p := map[string]map[string]interface{}{}
	var flatten func(*Schema, []string)
	flatten = func(s *Schema, locationPrefix []string) {
		for n, s := range s.Properties {
			if s.Type == "array" {
				continue
			}
			if s.Type == "object" {
				flatten(s, append(locationPrefix, n))
			} else {
				p[n] = map[string]interface{}{
					"location": append(locationPrefix, n),
					"schema":   s,
				}
			}
		}
	}
	flatten(&s, []string{})
	return p
}
