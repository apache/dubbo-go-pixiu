package constant

import (
	"reflect"
	"time"
)

// Object represents the java.lang.Object type
type Object interface{}

// JTypeMapper maps the java basic types to golang types
var JTypeMapper = map[string]reflect.Type{
	"string":           reflect.TypeOf(""),
	"java.lang.String": reflect.TypeOf(""),
	"char":             reflect.TypeOf(""),
	"short":            reflect.TypeOf(int32(0)),
	"int":              reflect.TypeOf(int32(0)),
	"long":             reflect.TypeOf(int64(0)),
	"float":            reflect.TypeOf(float64(0)),
	"double":           reflect.TypeOf(float64(0)),
	"boolean":          reflect.TypeOf(true),
	"java.util.Date":   reflect.TypeOf(time.Time{}),
	"date":             reflect.TypeOf(time.Time{}),
	"object":           reflect.TypeOf([]Object{}).Elem(),
	"java.lang.Object": reflect.TypeOf([]Object{}).Elem(),
}
