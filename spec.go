package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
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
	AllOf       []*Schema
	Default     interface{}
	Description string
	Items       *Schema
	Properties  map[string]*Schema
	ReadOnly    bool
	Required    []string
	Title       string
	Type        string
}

// Used for the purpose of flattening the properties of a schema. The location
// field makes it possible to reconstruct later. This facilitates generating
// setter/appender methods for deeply nested properties.
type FlatSchema struct {
	Name     string
	Location []string
	Schema   *Schema
}

// Return a schema's default value as JSON.
func (s Schema) DefaultJSON() string {
	b, err := json.Marshal(s.Default)
	if err != nil {
		log.Fatalln(err)
	}
	return string(b)
}

// IsDefault checks if the value matches the schema's default value
func (s Schema) IsDefault(value interface{}, name string) (bool, error) {
	switch value.(type) {
	case string:
		b := fmt.Sprintf("%v", s.Default)
		fmt.Fprintf(os.Stderr, "%s: %s<==>%s, %s\n", name, b, value, reflect.TypeOf(s.Default))
		fmt.Fprintf(os.Stderr, "  %t\n", b == value)
		return b == value, nil
	case float64:
		var bb float64
		switch b := s.Default.(type) {
		case float64:
			bb = b
		case int:
			bb = float64(b)
		case nil:
			return true, nil
		default:
			return false, fmt.Errorf("%s: Unknown type for IsDefault: %s", name, reflect.TypeOf(s.Default))
		}
		fmt.Fprintf(os.Stderr, "%s: %f<==>%f, %s\n", name, bb, value, reflect.TypeOf(s.Default))
		fmt.Fprintf(os.Stderr, "  %t\n", bb == value.(float64))
		return bb == value.(float64), nil
	case bool:
		var bb bool
		switch b := s.Default.(type) {
		case bool:
			bb = b
		case int:
			bb = int(b) != 0
		case nil:
			return true, nil
		default:
			return false, fmt.Errorf("%s: Unknown type for IsDefault: %s", name, reflect.TypeOf(s.Default))
		}
		fmt.Fprintf(os.Stderr, "%s: %t<==>%t, %s\n", name, bb, value, reflect.TypeOf(s.Default))
		fmt.Fprintf(os.Stderr, "  %t\n", bb == value.(bool))
		return bb == value.(bool), nil
	case nil:
		fmt.Fprintf(os.Stderr, "%s: %v<==>%v, %s\n", name, nil, value, reflect.TypeOf(s.Default))
		fmt.Fprintf(os.Stderr, "  %t\n", value == nil)
		return value == nil, nil
	case []interface{}:
		b, err := json.Marshal(s.Default)
		if err != nil {
			return false, err
		}
		c, err := json.Marshal(value)
		if err != nil {
			return false, err
		}
		fmt.Fprintf(os.Stderr, "%s: %s<==>%s, %s\n", name, b, c, reflect.TypeOf(s.Default))
		fmt.Fprintf(os.Stderr, "  %t\n", string(b) == string(c))
		return string(b) == string(c), nil
	case interface{}:
		b, err := json.Marshal(s.Default)
		if err != nil {
			return false, err
		}
		c, err := json.Marshal(value)
		if err != nil {
			return false, err
		}
		fmt.Fprintf(os.Stderr, "%s: %s<==>%s, %s\n", name, b, c, reflect.TypeOf(s.Default))
		fmt.Fprintf(os.Stderr, "  %t\n", string(b) == string(c))
		return string(b) == string(c), nil

	default:
		return false, fmt.Errorf("%s: Unknown value type for IsDefault: %s", name, reflect.TypeOf(value))
	}
}

// If title is set, it's assumed it carries more meaning than the property name
// itself. And therefore more suitable for humans. This is useful for naming
// arguments and functions.
func (s Schema) HumanName(name string) string {
	if s.Title != "" {
		return s.Title
	} else {
		return name
	}
}

// Combined schemas from Properties and AllOf.
func (s Schema) AllProperties() map[string]*Schema {
	ap := s.Properties
	if ap == nil {
		ap = map[string]*Schema{}
	}
	for _, aos := range s.AllOf {
		for n, subsc := range aos.Properties {
			ap[n] = subsc
		}
	}
	return ap
}

// TopLevelSimpleProperties returns all top-level properties except objects and arrays of
// objects. These are intended to be used as arguments for the schema object's constructor.
func (s Schema) TopLevelSimpleProperties() map[string]*Schema {
	p := map[string]*Schema{}
	for n, s := range s.AllProperties() {
		if !s.ReadOnly && s.Type != "object" &&
			(s.Type != "array" || s.Type == "array" && s.Items != nil && s.Items.Type != "object") {
			p[n] = s
		}
	}
	return p
}

func debug(msg string, value interface{}) {
	b, _ := json.MarshalIndent(value, "  ", "  ")
	fmt.Fprintf(os.Stderr, "LOG: %s %v\n", msg, string(b))
}

// TopLevelSimpleNonDefaultProperties returns all top-level properties except objects and arrays of
// objects. These are intended to be used as arguments for the schema object's constructor.
func (s Schema) TopLevelSimpleNonDefaultProperties(values map[string]interface{}) map[string]*Schema {
	p := map[string]*Schema{}
	for n, s := range s.AllProperties() {
		debug("NON DEFAULT SIMPLE NAME", n)
		debug("NON DEFAULT SIMPLE", s)

		if !s.ReadOnly && s.Type != "object" &&
			(s.Type != "array" || s.Type == "array" && s.Items != nil && s.Items.Type != "object") {
			debug("IS DEFAULT COMPARING", values[n])
			debug("                 AND", s.Default)
			isDefault, err := s.IsDefault(values[n], n)
			if err == nil && !isDefault {
				p[n] = s
			}
		}
		debug("DONE", n)
	}
	debug("FINISHED", p)
	return p
}

// Returns all properties that are readOnly and have a default property. It's
// intended that these are set, but not explicitly configurable. For example, a
// panel's "type" field.
func (s Schema) ReadOnlyWithDefaultProperties() []FlatSchema {
	return flatten(&s, func(n string, s *Schema) bool {
		return s.ReadOnly && s.Default != nil
	})
}

// Returns all top-level object properties. It's anticipated that these have
// setter methods nested inside their parent schema object.
func (s Schema) TopLevelObjectProperties() map[string]*Schema {
	p := map[string]*Schema{}
	for n, s := range s.AllProperties() {
		if !s.ReadOnly && s.Type == "object" {
			p[n] = s
		}
	}
	return p
}

// TopLevelArrayProperties returns all top-level array properties.
func (s Schema) TopLevelArrayProperties() map[string]*Schema {
	p := map[string]*Schema{}
	for n, s := range s.AllProperties() {
		if !s.ReadOnly && s.Type == "array" {
			fmt.Fprintf(os.Stderr, "ARRAY TYPE: %v\n", reflect.TypeOf(s))
			p[n] = s
		}
	}
	return p
}

//NestedSimpleProperties Returns all nested properties except arrays of objects. It's anticipated
// that the parent schema object is a top-level object property and that the
// properties returned here will be arguments in the parent's setter method.
func (s Schema) NestedSimpleProperties() []FlatSchema {
	return flatten(&s, func(n string, s *Schema) bool {
		return !s.ReadOnly && s.Type != "object" &&
			(s.Type != "array" || s.Type == "array" && s.Items != nil && s.Items.Type != "object")
	})
}

// NestedSimpleNonDefaultProperties Returns all nested properties except arrays of objects. It's anticipated
// that the parent schema object is a top-level object property and that the
// properties returned here will be arguments in the parent's setter method.
func (s Schema) NestedSimpleNonDefaultProperties(values map[string]interface{}) []FlatSchema {
	return flatten(&s, func(n string, s *Schema) bool {
		isDefault, err := s.IsDefault(values[n], n)
		fmt.Fprintf(os.Stderr, "NestedSimpleNonDefaultProperties: %s=%t\n", n, isDefault)
		return !s.ReadOnly && s.Type != "object" &&
			(s.Type != "array" || s.Type == "array" && s.Items != nil && s.Items.Type != "object") &&
			err == nil && !isDefault
	})
}

// Returns nested properties that are arrays of objects. It's anticipated that
// these are used to create appender methods for constructing those objects and
// appending them.
func (s Schema) NestedComplexArrayProperties() []FlatSchema {
	return flatten(&s, func(n string, s *Schema) bool {
		return !s.ReadOnly && s.Type == "array" && s.Items != nil && s.Items.Type == "object"
	})
}

// Recursively flattens nested properties.
func flatten(s *Schema, filter func(string, *Schema) bool) (fs []FlatSchema) {
	var flatten func(*Schema, []string)
	flatten = func(s *Schema, locationPrefix []string) {
		for n, s := range s.AllProperties() {
			if filter(n, s) {
				fs = append(fs, FlatSchema{
					Name:     n,
					Location: append(locationPrefix, n),
					Schema:   s,
				})
			} else if s.Type == "object" {
				flatten(s, append(locationPrefix, n))
			}
		}
	}
	flatten(s, []string{})
	sort.SliceStable(fs, func(i, j int) bool { return fs[i].Name < fs[j].Name })
	return fs
}
