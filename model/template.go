package model

import (
	"fmt"
	"strings"

	"bitbucket.org/pkg/inflect"
	"github.com/sparkymat/resty/field"
)

type modelTemplateValues struct {
	ModelName              string
	PrimaryKey             field.Type
	Fields                 []field.Type
	ResourceCollectionName string
	TableName              string
}

func (v modelTemplateValues) ModelNamePlural() string {
	return inflect.Camelize(inflect.Pluralize(inflect.Underscore(v.ModelName)))
}

func (v modelTemplateValues) ColumnNamesCSV() string {
	columnNameStrings := []string{}
	for _, field := range v.Fields {
		columnNameStrings = append(columnNameStrings, field.ColumnName)
	}
	return strings.Join(columnNameStrings, ",")
}

func (v modelTemplateValues) ValuePlacehodersCSV() string {
	valueStrings := []string{}

	for _, _ = range v.Fields {
		valueStrings = append(valueStrings, "?")
	}

	return strings.Join(valueStrings, ",")
}

func (v modelTemplateValues) ColumnNamesAndValuePlaceholdersCSV() string {
	displayStrings := []string{}

	for _, field := range v.Fields {
		displayStrings = append(displayStrings, fmt.Sprintf("%v=?", field.ColumnName))
	}

	return strings.Join(displayStrings, ",")
}

var modelTemplate = `package model

type {{.ModelName}} struct {
{{range $field := .Fields}} {{$field.FieldName}} {{$field.FieldType.Name}} {{$field.DbColumnNameAnnotation}}
{{end}}
}

// A {{.ModelName}}List aliases an slice of {{.ModelName}}, so that methods can be defined on the collection.
type {{.ModelName}}List []{{.ModelName}}

// FindAll{{.ModelNamePlural}} fetches all {{.ModelName}} objects from the database.
// TODO: You might want to add a default filter, as well as custom Find methods, for performance reasons.
func FindAll{{.ModelNamePlural}}() ({{.ModelName}}List, error) {
	collection := []{{.ModelName}}{}

	sql := "SELECT * FROM {{.TableName}}"
	rows, err := db.DB.Queryx(sql)
	if err != nil {
		return collection, err
	}
	
	defer rows.Close()

	for rows.Next() {
		var v {{.ModelName}}
		err := rows.StructScan(&v)
		if err == nil {
			collection = append(collection, v)
		}
	}

	return collection, nil
}

// Find{{.ModelName}}By{{.PrimaryKey.FieldName}} will fetch the {{.ModelName}} object with the specified key.
func Find{{.ModelName}}By{{.PrimaryKey.FieldName}}(key {{.PrimaryKey.FieldType.String}}) (*{{.ModelName}}, error) {
	sql := "SELECT * FROM {{.TableName}} WHERE {{.PrimaryKey.ColumnName}} = ?"
	rows, err := db.DB.Queryx(sql, identifier)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var t {{.ModelName}}
	rows.Next()
	err = rows.StructScan(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// New{{.ModelName}} will create a new instance of {{.ModelName}} and populate it with the values in params.
func New{{.ModelName}}(params map[string][]string) ({{.ModelName}}, error) {
	var v {{.ModelName}}

	{{range $field := .Fields}} {{$field.ParamParsingStatements}} 
	{{end}}

	return v, nil
}

// Save will INSERT a new entry into the table, if the primary key field is empty. If it's not empty, it will update the existing row.
func (t {{.ModelName}}) Save() error {
	var sql string

	if t.{{.PrimaryKey.FieldName}} == {{.PrimaryKey.ZeroValue}} {
		sql = "INSERT INTO people ({{.ColumnNamesCSV}}) VALUES ({{.ValuePlacehodersCSV}})"
	} else {
		sql = "UPDATE people SET ({{.ColumnNamesAndValuePlaceholdersCSV}}) WHERE {{.PrimaryKey.ColumnName}} = ?"
	}
}

// AsMap will return a map representation of the {{.ModelName}} object, which can be used to generate a JSON representation
func (t {{.ModelName}}) AsMap() map[string]interface{} {
	return map[string]interface{} {
	{{range $field := .Fields}} "{{$field.FieldNameUnderscored}}": t.{{$field.FieldName}},
	{{end}}
	}
}

// AsMap will return a map representation of the {{.ModelName}}List object (which is a slice of {{.ModelName}} objects), which can be used to generate a JSON representation
func (l {{.ModelName}}List) AsMap() map[string]interface{} {
	nodes := make([]map[string]interface{}, 0)

	for _, node := range l {
		nodes = append(nodes, node.AsMap())
	}

	return map[string]interface{}{
		"{{.ResourceCollectionName}}": nodes,
	}
}
`
