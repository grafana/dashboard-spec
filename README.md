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

#### Style Guide

All properties of an object should be defined alphabetically.

All properties should have a description. This is used for API doc generation.

File names should be in camel case. All files referenced as a schema component
in a `spec.yml` should begin with a capital letter (PascalCase).

Most objects are either a "panel", "datasource", or "template". Each object's
definition should live in its respective directory of the spec version it's
modeling.

### [bundle/](./bundle)

Single file JSON specs generated from the YAML in [specs/](./specs). This is
what should be consumed by OpenAPI tooling like code generators.
