package models

import (
	"errors"
	"path/filepath"
	"strings"
)

type SchemaType string

const (
	JSON SchemaType = "json"
	YAML SchemaType = "yaml"
	TXT  SchemaType = "txt"
)

func (SchemaType) GetSchemaFromFilename(filename string) (SchemaType, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return JSON, nil
	case ".yaml", ".yml":
		return YAML, nil
	case "txt":
		return TXT, nil
	default:
		return "", errors.New("unsupported file extension")
	}
}
