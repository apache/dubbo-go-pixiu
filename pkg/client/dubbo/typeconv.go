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

package dubbo

import (
	"net/url"
	"reflect"
	"strings"
	"time"
)

import (
	"github.com/pkg/errors"

	"github.com/spf13/cast"
)

import (
	cst "github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

var javaWrapperTypeAliases = map[string]string{
	cst.JavaLangBooleanClassName: cst.JavaPrimitiveBoolean,
	cst.JavaLangByteClassName:    cst.JavaPrimitiveByte,
	cst.JavaLangCharClassName:    cst.JavaPrimitiveChar,
	cst.JavaLangDoubleClassName:  cst.JavaPrimitiveDouble,
	cst.JavaLangFloatClassName:   cst.JavaPrimitiveFloat,
	cst.JavaLangIntegerClassName: cst.JavaPrimitiveInt,
	cst.JavaLangLongClassName:    cst.JavaPrimitiveLong,
	cst.JavaLangShortClassName:   cst.JavaPrimitiveShort,
}

// MapTypes converts a declared java type into a Go value using the dubbo mapper.
func MapTypes(jType string, originVal any) (any, error) {
	normalized := normalizeJavaTypeName(jType)
	if normalized == "" {
		return originVal, nil
	}
	targetType, ok := cst.JTypeMapper[normalized]
	if !ok {
		return nil, errors.Errorf("Invalid parameter type: %s", normalized)
	}
	switch targetType {
	case reflect.TypeOf(""):
		return cast.ToStringE(originVal)
	case reflect.TypeOf(int(0)):
		return cast.ToIntE(originVal)
	case reflect.TypeOf(int8(0)):
		return cast.ToInt8E(originVal)
	case reflect.TypeOf(int16(16)):
		return cast.ToInt16E(originVal)
	case reflect.TypeOf(int32(0)):
		return cast.ToInt32E(originVal)
	case reflect.TypeOf(int64(0)):
		return cast.ToInt64E(originVal)
	case reflect.TypeOf(float32(0)):
		return cast.ToFloat32E(originVal)
	case reflect.TypeOf(float64(0)):
		return cast.ToFloat64E(originVal)
	case reflect.TypeOf(true):
		return cast.ToBoolE(originVal)
	case reflect.TypeOf(time.Time{}):
		return cast.ToTimeE(originVal)
	default:
		return originVal, nil
	}
}

// CoerceDirectInvokeValue coerces direct invoke values according to parameter types.
func CoerceDirectInvokeValue(parameterType string, value any) (any, error) {
	trimmed := strings.TrimSpace(parameterType)
	if trimmed == "" {
		return value, nil
	}

	if strings.HasSuffix(trimmed, "[]") {
		elementType := normalizeJavaTypeName(strings.TrimSuffix(trimmed, "[]"))
		if _, ok := cst.JTypeMapper[elementType]; !ok {
			return value, nil
		}

		items, ok := value.([]any)
		if !ok {
			return value, nil
		}

		result := make([]any, len(items))
		for i, item := range items {
			mapped, err := MapTypes(elementType, item)
			if err != nil {
				return nil, err
			}
			result[i] = mapped
		}
		return result, nil
	}

	if _, ok := cst.JTypeMapper[normalizeJavaTypeName(trimmed)]; ok {
		return MapTypes(trimmed, value)
	}

	return value, nil
}

// NormalizeReferenceProtocol canonicalizes dubbo reference protocols.
func NormalizeReferenceProtocol(protocol string) string {
	normalized := strings.ToLower(strings.TrimSpace(protocol))
	switch normalized {
	case "triple":
		return "tri"
	default:
		return normalized
	}
}

// DirectURLProtocol returns the normalized protocol from a direct URL.
func DirectURLProtocol(rawURL string) (string, error) {
	parsedURL, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", errors.Wrapf(err, "parse direct url %q", rawURL)
	}
	if strings.TrimSpace(parsedURL.Scheme) == "" {
		return "", errors.New("direct url requires scheme")
	}
	return NormalizeReferenceProtocol(parsedURL.Scheme), nil
}

// InferJavaClassNames infers java class names from Go values.
func InferJavaClassNames(values []any) []string {
	types := make([]string, len(values))
	for i, val := range values {
		switch val.(type) {
		case string:
			types[i] = cst.JavaLangStringClassName
		case bool:
			types[i] = cst.JavaLangBooleanClassName
		case float32:
			types[i] = cst.JavaLangFloatClassName
		case float64:
			types[i] = cst.JavaLangDoubleClassName
		case int16:
			types[i] = cst.JavaLangShortClassName
		case int32:
			types[i] = cst.JavaLangIntegerClassName
		case int, int8, int64:
			types[i] = cst.JavaLangLongClassName
		case uint, uint8, uint16, uint32, uint64:
			types[i] = cst.JavaLangLongClassName
		default:
			types[i] = cst.JavaLangObjectClassName
		}
	}
	return types
}

func normalizeJavaTypeName(jType string) string {
	trimmed := strings.TrimSpace(jType)
	if trimmed == "" {
		return ""
	}
	if alias, ok := javaWrapperTypeAliases[trimmed]; ok {
		return alias
	}
	return trimmed
}
