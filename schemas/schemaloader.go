package schemas

import (
	"os"
	"fmt"
	"path/filepath"
	"github.com/santhosh-tekuri/jsonschema/v5"
)



func NewSchemaLoader(prefixPath string) (*SchemaLoader, error) {
	compiler := jsonschema.NewCompiler()

	error := filepath.Walk(prefixPath, func(path string, fileInfo os.FileInfo, error error) error {
		if error != nil {
			return error
		}

		if filepath.Ext(path) == ".json" {
			file, openError := os.Open(path)

			if openError != nil {
				return fmt.Errorf("Error opening schema %s: %w", path, openError)
			}

			relativePath, relativePathError := filepath.Rel(prefixPath, path)

			if relativePathError != nil {
				file.Close()
				return fmt.Errorf("Error resolving schema path %s: %w", path, relativePathError)
			}

			schemaPath := filepath.ToSlash(filepath.Join(filepath.Base(prefixPath), relativePath))
			addResourceError := compiler.AddResource(schemaPath, file)
			file.Close()

			if addResourceError != nil {
				return fmt.Errorf("Error adding schema %s: %w", path, addResourceError)
			}
		}

		return nil
	})

	if error != nil {
		return nil, fmt.Errorf("Error walking schemas: %w", error)
	}

	return &SchemaLoader{compiler: compiler}, nil
}

func (schemaLoader *SchemaLoader) Compile(path string) (*jsonschema.Schema, error) {
	schema, error  := schemaLoader.compiler.Compile(path)

	if error != nil {
		return nil, fmt.Errorf("Error compiling schema %s: %w", path, error)
	}

	return schema, nil
}

type SchemaLoader struct {
	compiler *jsonschema.Compiler
}