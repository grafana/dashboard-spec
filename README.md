# Grafana Dashboard Spec

Specification for [Grafana Dashboard
JSON](https://grafana.com/docs/grafana/latest/reference/dashboard/) and core
panels using the [OpenAPI
Specification](https://github.com/OAI/OpenAPI-Specification).

This can used for generating models in a variety of programming languages. The
models facilitate writing dashboards as code.

## Repo Layout

### [specs/](./specs)

Human-managed specification YAML files.

#### Style Guide for Spec Files

All properties of an object should be defined alphabetically.

All properties should have a description. This is used for API doc generation.

If a property's name could be more descriptive in code or it collides with
another name nested in the same object, use the `title` field to indicate what
that object should be called in code. For example, most panels have a top-level
array called `links` and also a nested array called `links`. The top-level array
is referring to [panel
links](https://grafana.com/docs/grafana/latest/linking/panel-links/) while the
nested array is referring to [data
links](https://grafana.com/docs/grafana/latest/linking/data-links/), therefore,
the properties have `title` set to "Panel Link" and "Data Link". The code
generator should use this field instead for deciding method names. Depending on
what the language has set for its object inflection property, this will result
in methods like, `addPanelLinks()` and `addDataLinks()`.

File names should be in camel case. All files referenced as a schema component
in a `spec.yml` should begin with a capital letter (PascalCase).

Most objects are either a "panel", "datasource", or "template". Each object's
definition should live in its respective directory of the spec version it's
modeling.

### [bundle/](./bundle)

Single file JSON specs generated from the YAML in [specs/](./specs). This is
what should be consumed by OpenAPI tooling like code generators.

### [templates/](./templates)

Templates for the code generator. Child directories are named after the language
they contain templates for.

Each language must implment the following templates:

* `main.tmpl`: this is the main library file. It's intended that this file be
  imported when implementing the generated code.
* `dashboard.tmpl`: for generating the dashboard object and file.
* `panel.tmpl`: for generating panel objects and files.
* `target.tmpl`: for generating target objects and files.
* `template.tmpl`: for generating template objects and files.

#### Style Guide for Templates

Arrays of objects should use mutator functions to append to them. For example
`addLink()`.

Arrays of single values should be set as top level arguemnts.

First level nested objects should also use mutator functions with all non-array
fields as arguemnts. For example, `feildConfig(min=0, max=100)`.

If fields need special processing, set them as readOnly and implement static
functions. For example if you need to add an incrementing `id` field like we do
for panels.

If a property is `readOnly` and also has default, set the default as a static
value on the object.
