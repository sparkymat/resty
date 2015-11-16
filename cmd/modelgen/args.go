package main

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/pkg/inflect"
	"github.com/sparkymat/resty/cmd/modelgen/field"
)

func processArgs() (string, string, []field.Type) {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: modelgen model-name(:table-name) primary-key:type(:column_name) (optional: field1:type field2:type:column-name ... )\n")
		os.Exit(1)
	}

	args := os.Args[2:]
	var fields []field.Type

	for _, arg := range args {
		pieces := strings.Split(arg, ":")

		if len(pieces) < 2 {
			fmt.Fprintf(os.Stderr, "Error: Field '%v' is not formatted correctly\n", pieces[0])
			os.Exit(1)
		}

		fieldType, ok := supportedTypes[pieces[1]]

		if !ok {
			fmt.Fprintf(os.Stderr, "Error: Field '%v' doesn't have a supported type\n", pieces[0])
			os.Exit(1)
		}

		fieldName := inflect.Camelize(pieces[0])
		columnName := ""

		if len(pieces) > 2 {
			columnName = pieces[2]
		} else {
			columnName = inflect.Underscore(pieces[0])
		}

		var field field.Type
		field.FieldName = fieldName
		field.FieldType = fieldType
		field.ColumnName = columnName

		fields = append(fields, field)
	}

	modelName := inflect.Camelize(os.Args[1])
	tableName := inflect.Pluralize(inflect.Underscore(modelName))

	modelParts := strings.Split(os.Args[1], ":")
	if len(modelParts) > 1 {
		modelName = inflect.Camelize(modelParts[0])
		tableName = modelParts[1]
	}

	return modelName, tableName, fields
}
