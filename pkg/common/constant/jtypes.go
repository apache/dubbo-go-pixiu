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
type Object interface{}

// JTypeMapper maps the java basic types to golang types
var JTypeMapper = map[string]reflect.Type{
	"string":           reflect.TypeOf(""),
	"java.lang.String": reflect.TypeOf(""),
	"char":             reflect.TypeOf(""),
	"short":            reflect.TypeOf(int16(0)),
	"int":              reflect.TypeOf(int(0)),
	"long":             reflect.TypeOf(int64(0)),
	"float":            reflect.TypeOf(float32(0)),
	"double":           reflect.TypeOf(float64(0)),
	"boolean":          reflect.TypeOf(true),
	"java.util.Date":   reflect.TypeOf(time.Time{}),
	"date":             reflect.TypeOf(time.Time{}),
	"object":           reflect.TypeOf([]Object{}).Elem(),
	"java.lang.Object": reflect.TypeOf([]Object{}).Elem(),
}
