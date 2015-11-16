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
