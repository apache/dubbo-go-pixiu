/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package constant

import (
	"reflect"
	"time"
)

// Object represents the java.lang.Object type
type Object any

const (
	JavaLangStringClassName  = "java.lang.String"
	JavaLangBooleanClassName = "java.lang.Boolean"
	JavaLangByteClassName    = "java.lang.Byte"
	JavaLangCharClassName    = "java.lang.Character"
	JavaLangShortClassName   = "java.lang.Short"
	JavaLangIntegerClassName = "java.lang.Integer"
	JavaLangLongClassName    = "java.lang.Long"
	JavaLangFloatClassName   = "java.lang.Float"
	JavaLangDoubleClassName  = "java.lang.Double"
	JavaLangObjectClassName  = "java.lang.Object"
	JavaUtilDateClassName    = "java.util.Date"
)

const (
	JavaPrimitiveString  = "string"
	JavaPrimitiveChar    = "char"
	JavaPrimitiveShort   = "short"
	JavaPrimitiveInt     = "int"
	JavaPrimitiveLong    = "long"
	JavaPrimitiveFloat   = "float"
	JavaPrimitiveDouble  = "double"
	JavaPrimitiveBoolean = "boolean"
	JavaPrimitiveByte    = "byte"
	JavaPrimitiveObject  = "object"
	JavaPrimitiveDate    = "date"
)

// JTypeMapper maps the java basic types to golang types
var JTypeMapper = map[string]reflect.Type{
	JavaPrimitiveString:     reflect.TypeOf(""),
	JavaLangStringClassName: reflect.TypeOf(""),
	JavaPrimitiveChar:       reflect.TypeOf(""),
	JavaPrimitiveShort:      reflect.TypeOf(int16(0)),
	JavaPrimitiveInt:        reflect.TypeOf(int(0)),
	JavaPrimitiveLong:       reflect.TypeOf(int64(0)),
	JavaPrimitiveFloat:      reflect.TypeOf(float32(0)),
	JavaPrimitiveDouble:     reflect.TypeOf(float64(0)),
	JavaPrimitiveBoolean:    reflect.TypeOf(true),
	JavaUtilDateClassName:   reflect.TypeOf(time.Time{}),
	JavaPrimitiveDate:       reflect.TypeOf(time.Time{}),
	JavaPrimitiveObject:     reflect.TypeOf([]Object{}).Elem(),
	JavaLangObjectClassName: reflect.TypeOf([]Object{}).Elem(),
}
