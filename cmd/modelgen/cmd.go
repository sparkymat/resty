package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/sparkymat/resty/cmd/modelgen/golang"

	"bitbucket.org/pkg/inflect"
)

var supportedTypes = map[string]reflect.Type{
	"int64":   reflect.TypeOf(int64(0)),
	"int32":   reflect.TypeOf(int32(0)),
	"string":  reflect.TypeOf(string("")),
	"boolean": reflect.TypeOf(bool(true)),
	"float64": reflect.TypeOf(float64(0.0)),
	"float32": reflect.TypeOf(float32(0.0)),
}

func main() {
	modelName, fields := processArgs()
	tempPath := fmt.Sprintf("model/%v.go.temp", inflect.Underscore(modelName))
	outPath := fmt.Sprintf("model/%v.go", inflect.Underscore(modelName))

	fieldLines := []string{}

	for _, field := range fields {
		fieldLines = append(fieldLines, fmt.Sprintf("%v %v `db:\"%v\"`", field.FieldName, field.FieldType.Name(), field.ColumnName))
	}

	fieldLine := strings.Join(fieldLines, "\n")

	os.MkdirAll("model", 0755)
	fp, err := os.Create(tempPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to create model file. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(fp, `package model

type %v struct {
	%v
}

func (v %v) Find%vById(id int64) (*%v, error) {
}
	`, modelName, fieldLine, modelName, modelName, modelName)
	fp.Close()

	err = golang.Fmt(tempPath, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to gofmt the model file. Reason: %v\n", err)
		os.Exit(1)
	}
	os.Remove(tempPath)
}
