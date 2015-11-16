package field

import "reflect"

type Type struct {
	FieldName  string
	FieldType  reflect.Type
	ColumnName string
}
