package model

import (
	"fmt"
	"os"
	"strings"

	"bitbucket.org/pkg/inflect"
	"github.com/sparkymat/resty/field"
)

func processArgs(args []string) (string, string, []field.Type) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: Missing input args\n")
		os.Exit(1)
	}

	var fields []field.Type

	for _, arg := range args[1:] {
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

	modelName := inflect.Camelize(args[0])
	tableName := inflect.Pluralize(inflect.Underscore(modelName))

	modelParts := strings.Split(args[0], ":")
	if len(modelParts) > 1 {
		modelName = inflect.Camelize(modelParts[0])
		tableName = modelParts[1]
	}

	return modelName, tableName, fields
}
