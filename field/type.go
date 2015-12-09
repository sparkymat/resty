package field

import (
	"fmt"
	"reflect"

	"bitbucket.org/pkg/inflect"
)

type Type struct {
	FieldName  string
	FieldType  reflect.Type
	ColumnName string
}

func (t Type) FieldNameUnderscored() string {
	return inflect.Underscore(t.FieldName)
}

func (t Type) DbColumnNameAnnotation() string {
	return fmt.Sprintf("`db:\"%v\"`", t.ColumnName)
}

func (t Type) ZeroValue() string {
	if t.FieldType == reflect.TypeOf("") {
		return "\"\""
	}

	return fmt.Sprintf("%v", reflect.Zero(t.FieldType))
}

func (t Type) ParamParsingStatements() string {
	lowerFieldName := inflect.CamelizeDownFirst(t.FieldName)

	switch t.FieldType.String() {
	case "int64":
		return fmt.Sprintf(`%vStrings, ok := params["%v"]
if ok && len(%vStrings) >= 1 {
       %v, err := strconv.ParseInt(%vStrings[0], 10, 64) 
       return v, err
       v.%v = %v
}
`, lowerFieldName, t.FieldNameUnderscored(), lowerFieldName, lowerFieldName, lowerFieldName, t.FieldName, lowerFieldName)
	case "int32":
		return fmt.Sprintf(`%vStrings, ok := params["%v"]
if ok && len(%vStrings) >= 1 {
       %v, err := strconv.ParseInt(%vStrings, 10, 32) 
       if err != nil {
               return v, err
       }
       v.%v = %v
}
`, lowerFieldName, t.FieldNameUnderscored(), lowerFieldName, lowerFieldName, lowerFieldName, t.FieldName, lowerFieldName)
	case "string":
		return fmt.Sprintf(`if _, ok := params["%v"]; ok {
  if len(params["%v"]) > 0 {
         v.%v = params["%v"][0]
  }
}
`, t.FieldNameUnderscored(), t.FieldNameUnderscored(), t.FieldName, t.FieldNameUnderscored())
	case "bool":
		return fmt.Sprintf(`if _, ok := params["%v"]; ok {
  if len(params["%v"]) > 0 {
		if parmas["%v"][0] == "true" {
			v.%v = true
		} else if params["%v"][0] == "false" {
      v.%v = false
		}
  }
}
`, t.FieldNameUnderscored(), t.FieldNameUnderscored(), t.FieldNameUnderscored(), t.FieldName, t.FieldNameUnderscored(), t.FieldName)
	case "float64":
		return fmt.Sprintf(`%vStrings, ok := params["%v"]
if ok && len(%vStrings) >= 1 {
       %v, err := strconv.ParseFloat(%vStrings[0], 64) 
       return v, err
       v.%v = %v
}
`, lowerFieldName, t.FieldNameUnderscored(), lowerFieldName, lowerFieldName, lowerFieldName, t.FieldName, lowerFieldName)
	case "float32":
		return fmt.Sprintf(`%vStrings, ok := params["%v"]
if ok && len(%vStrings) >= 1 {
       %v, err := strconv.ParseFloat(%vStrings[0], 32) 
       return v, err
       v.%v = %v
}
`, lowerFieldName, t.FieldNameUnderscored(), lowerFieldName, lowerFieldName, lowerFieldName, t.FieldName, lowerFieldName)
	}

	return ""
}
