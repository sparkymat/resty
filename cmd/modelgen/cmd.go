package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"bitbucket.org/pkg/inflect"
)

type modelField struct {
	fieldName  string
	fieldType  reflect.Type
	columnName string
}

var supportedTypes = map[string]reflect.Type{
	"int64":   reflect.TypeOf(int64(0)),
	"int32":   reflect.TypeOf(int32(0)),
	"string":  reflect.TypeOf(string("")),
	"boolean": reflect.TypeOf(bool(true)),
	"float64": reflect.TypeOf(float64(0.0)),
	"float32": reflect.TypeOf(float32(0.0)),
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: modelgen model-name (optional: field1:type field2:type:column-name ... )\n")
		os.Exit(1)
	}

	modelName := os.Args[1]
	args := os.Args[2:]

	var fields []modelField

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

		var field modelField
		field.fieldName = fieldName
		field.fieldType = fieldType
		field.columnName = columnName

		fields = append(fields, field)
	}

	os.MkdirAll("model", 0755)
	filePath := fmt.Sprintf("model/%v.go.temp", inflect.Underscore(modelName))
	fp, err := os.Create(filePath)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to create model file. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	fieldLines := []string{}

	for _, field := range fields {
		fieldLines = append(fieldLines, fmt.Sprintf("%v %v `db:\"%v\"`", field.fieldName, field.fieldType.Name(), field.columnName))
	}

	fieldLine := strings.Join(fieldLines, "\n")

	fmt.Fprintf(fp, `package model

type %v struct {
	%v
}

func (v %v) Find%vById(id int64) (*%v, error) {
}
	`, modelName, fieldLine, modelName, modelName, modelName)
	fp.Close()

	cmdArgs := []string{
		fmt.Sprintf("model/%v.go.temp", inflect.Underscore(modelName)),
	}
	cmd := exec.Command("gofmt", cmdArgs...)
	outPipe, err := cmd.StdoutPipe()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to run gofmt. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	outFile, err := os.Create(fmt.Sprintf("model/%v.go", inflect.Underscore(modelName)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to create model file. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Unable to create model file. Reason: %v\n", err.Error())
		os.Exit(1)
	}

	io.Copy(outFile, outPipe)
	outFile.Close()

	os.Remove(fmt.Sprintf("model/%v.go.temp", inflect.Underscore(modelName)))
}
