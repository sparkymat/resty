package model

import (
	"fmt"
	"os"
	"reflect"
	"text/template"

	"github.com/sparkymat/resty/golang"

	"bitbucket.org/pkg/inflect"
)

var supportedTypes = map[string]reflect.Type{
	"int64":   reflect.TypeOf(int64(0)),
	"int32":   reflect.TypeOf(int32(0)),
	"string":  reflect.TypeOf(string("")),
	"bool":    reflect.TypeOf(bool(true)),
	"float64": reflect.TypeOf(float64(0.0)),
	"float32": reflect.TypeOf(float32(0.0)),
}

func Generate(args []string) {
	modelName, tableName, fields := processArgs(args)
	tempPath := fmt.Sprintf("model/%v.go.temp", inflect.Underscore(modelName))
	outPath := fmt.Sprintf("model/%v.go", inflect.Underscore(modelName))

	os.MkdirAll("model", 0755)

	tpl := template.New("model")
	_, err := tpl.Parse(modelTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to load template. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	fp, err := os.Create(tempPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to create model file. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	values := modelTemplateValues{}
	values.ModelName = modelName
	values.TableName = tableName
	values.Fields = fields
	values.PrimaryKey = fields[0]
	values.ResourceCollectionName = inflect.Pluralize(inflect.Underscore(modelName))

	err = tpl.Execute(fp, values)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to generate model file. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	fp.Close()

	err = golang.Fmt(tempPath, outPath)
	os.Remove(tempPath)
	if err != nil {
		os.Remove(outPath)
		fmt.Fprintf(os.Stderr, "Error: Unable to gofmt the model file. Reason: %v\n", err.Error())
		os.Exit(1)
	}
}
