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
	Title       string
	Default     interface{}
	Description string
	Items       *Schema
	Properties  map[string]*Schema
	ReadOnly    bool
	Required    []string
	Type        string
}

// Return a schema's default property as JSON.
func (s Schema) DefaultJSON() string {
	b, err := json.Marshal(s.Default)
	if err != nil {
		log.Fatalln(err)
	}
	return string(b)
}

// Returns all top-level properties that are not an array or object. These are
// intended to be used as function arguments for the object's constructor.
func (s Schema) TopLevelSingleValProperties() map[string]*Schema {
	p := map[string]*Schema{}
	for n, s := range s.Properties {
		if s.Type != "array" && s.Type != "object" && !s.ReadOnly {
			p[n] = s
		}
	}
	return p
}

// Returns all top-level object properties. It's intended that these are
// implmented as methods.
func (s Schema) TopLevelObjectProperties() map[string]*Schema {
	p := map[string]*Schema{}
	for n, s := range s.Properties {
		if s.Type == "object" && !s.ReadOnly {
			p[n] = s
		}
	}
	return p
}

// Returns all properties that are readOnly and have a default property. It's
// intended that these are set, but not explicitly configurable.
func (s Schema) ReadOnlyWithDefaultProperties() map[string]map[string]interface{} {
	return flatten(&s, func(s *Schema) bool {
		return s.ReadOnly && s.Default != nil
	})
}

// Returns nested objects and arrays that should be part of a constructor
// method. This includes all objects and flat arrays. This is used to simplify
// object interfaces on those with many levels of nesting.
func (s Schema) ConstructableProperties() map[string]map[string]interface{} {
	return flatten(&s, func(s *Schema) bool {
		return !s.ReadOnly || s.Type != "array" && s.Items.Type != "object"
	})
}

// Returns nested arrays of objects. This is used to create methods for
// constructing those objects and appending them to an array.
func (s Schema) AppendableProperties() map[string]map[string]interface{} {
	return flatten(&s, func(s *Schema) bool {
		return !s.ReadOnly && s.Type == "array" && s.Items.Type == "object"
	})
}

// Recursively flattens nested properties.
// Returns a nested map to model the structure. Location is set so the schema
// can get pieced back together.
func flatten(s *Schema, filter func(*Schema) bool) map[string]map[string]interface{} {
	p := map[string]map[string]interface{}{}
	var flatten func(*Schema, []string)
	flatten = func(s *Schema, locationPrefix []string) {
		for n, s := range s.Properties {
			if filter(s) {
				p[n] = map[string]interface{}{
					"location": append(locationPrefix, n),
					"schema":   s,
				}
			} else if s.Type == "object" {
				flatten(s, append(locationPrefix, n))
			}
		}
	}
	flatten(s, []string{})
	return p
}
