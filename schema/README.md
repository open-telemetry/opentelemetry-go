# Telemetry Schema Files

The `schema` module contains packages that help to parse and validate
[schema files](https://github.com/open-telemetry/oteps/blob/main/text/0152-telemetry-schemas.md).

Each `major.minor` schema file format version is implemented as a separate package, with
the name of the package in the `vmajor.minor` form.

To parse a schema file, first decide what file format version you want to parse,
then import the corresponding package and use the `Parse` or `ParseFile` functions
like this:

```go
import schema "go.opentelemetry.io/otel/schema/v1.1"

// Load the schema from a file in v1.1.x file format.
func loadSchemaFromFile() error {
	telSchema, err := schema.ParseFile("schema-file.yaml")
	if err != nil {
		return err
	}
	// Use telSchema struct here.
}

// Alternatively use schema.Parse to read the schema file from io.Reader.
func loadSchemaFromReader(r io.Reader) error {
	telSchema, err := schema.Parse(r)
	if err != nil {
		return err
	}
	// Use telSchema struct here.
}
```
