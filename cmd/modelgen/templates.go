package main

import "github.com/sparkymat/resty/cmd/modelgen/field"

type modelTemplateValues struct {
	ModelName              string
	PrimaryKey             field.Type
	Fields                 []field.Type
	ResourceCollectionName string
	BackTick               string
	TableName              string
}

var modelTemplate = `package model
type {{.ModelName}} struct {
{{range $field := .Fields}} {{$field.FieldName}} {{$field.FieldType.Name}} {{$field.DbColumnNameAnnotation}}
{{end}}
}

type {{.ModelName}}List []{{.ModelName}}

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

func (t {{.ModelName}}) AsMap() map[string]interface{} {
	return map[string]interface{} {
	{{range $field := .Fields}} "{{$field.FieldNameUnderscored}}": t.{{$field.FieldName}},
	{{end}}
	}
}

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
