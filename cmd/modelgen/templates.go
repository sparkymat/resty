package main

var modelTemplate = `package model

type {{.ModelName}} struct {
  {{.FieldLines}}
}

func (t {{.ModelName}}) Find{{.ModelName}}By{{.PrimaryKey}}(key {{.PrimaryKeyType}}) (*{{.ModelName}}, error) {
}
`
