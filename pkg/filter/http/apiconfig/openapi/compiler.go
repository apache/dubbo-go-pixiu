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

package openapi

import (
	"fmt"
	"regexp"
	"strings"
)

import (
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"go.yaml.in/yaml/v4"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

var openAPIPathParamPattern = regexp.MustCompile(`\{([^}/]+)\}`)

type CompiledRoute struct {
	Route      router.API
	Validation *ValidationPlan
}

func LoadDocumentFromBytes(spec []byte) (libopenapi.Document, error) {
	return libopenapi.NewDocument(spec)
}

func CompileRoutes(doc libopenapi.Document) ([]CompiledRoute, error) {
	model, err := doc.BuildV3Model()
	if err != nil {
		return nil, err
	}
	if model == nil || model.Model.Paths == nil {
		return nil, fmt.Errorf("openapi document does not contain paths")
	}

	var compiled []CompiledRoute
	for path, item := range model.Model.Paths.PathItems.FromOldest() {
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "GET"}, item, item.Get)
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "PUT"}, item, item.Put)
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "POST"}, item, item.Post)
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "DELETE"}, item, item.Delete)
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "PATCH"}, item, item.Patch)
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "HEAD"}, item, item.Head)
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "OPTIONS"}, item, item.Options)
		compiled = appendCompiledRoute(compiled, path, config.Method{Enable: true, HTTPVerb: "TRACE"}, item, item.Trace)
	}

	return compiled, nil
}

func appendCompiledRoute(compiled []CompiledRoute, path string, method config.Method, pathItem *v3.PathItem, operation *v3.Operation) []CompiledRoute {
	if operation == nil {
		return compiled
	}

	validationPlan := compileValidationPlan(path, pathItem, operation)

	return append(compiled, CompiledRoute{
		Route: router.API{
			URLPattern: normalizeOpenAPIPath(path),
			Method:     method,
			Metadata: map[string]any{
				ValidationPlanMetadataKey: validationPlan,
			},
		},
		Validation: validationPlan,
	})
}

func compileValidationPlan(path string, pathItem *v3.PathItem, operation *v3.Operation) *ValidationPlan {
	plan := &ValidationPlan{}
	if operation == nil {
		return plan
	}
	plan.RoutePattern = normalizeOpenAPIPath(path)

	for _, parameter := range mergedParameters(pathItem, operation) {
		if parameter == nil {
			continue
		}
		compiled := compileParameter(parameter)
		switch strings.ToLower(parameter.In) {
		case "path":
			plan.PathParameters = append(plan.PathParameters, compiled)
		case "query":
			plan.QueryParameters = append(plan.QueryParameters, compiled)
		case "header":
			plan.HeaderParameters = append(plan.HeaderParameters, compiled)
		}
	}

	plan.RequestBody = compileRequestBody(operation.RequestBody)
	return plan
}

func mergedParameters(pathItem *v3.PathItem, operation *v3.Operation) []*v3.Parameter {
	var parameters []*v3.Parameter
	if pathItem != nil {
		parameters = append(parameters, pathItem.Parameters...)
	}
	if operation == nil {
		return parameters
	}

	for _, parameter := range operation.Parameters {
		if parameter == nil {
			parameters = append(parameters, nil)
			continue
		}
		replaced := false
		for idx, existing := range parameters {
			if existing == nil {
				continue
			}
			if strings.EqualFold(existing.Name, parameter.Name) && strings.EqualFold(existing.In, parameter.In) {
				parameters[idx] = parameter
				replaced = true
				break
			}
		}
		if !replaced {
			parameters = append(parameters, parameter)
		}
	}

	return parameters
}

func compileParameter(parameter *v3.Parameter) ParameterValidation {
	compiled := ParameterValidation{
		Name: parameter.Name,
	}
	if parameter.Required != nil {
		compiled.Required = *parameter.Required
	}
	if parameter.Schema != nil {
		schema, err := parameter.Schema.BuildSchema()
		if err == nil && schema != nil {
			compiled.Type = firstType(schema)
			compiled.Enum = decodeEnum(schema.Enum)
			if schema.MinLength != nil {
				value := int(*schema.MinLength)
				compiled.MinLength = &value
			}
			if schema.MaxLength != nil {
				value := int(*schema.MaxLength)
				compiled.MaxLength = &value
			}
			if schema.Minimum != nil {
				value := *schema.Minimum
				compiled.Minimum = &value
			}
			if schema.Maximum != nil {
				value := *schema.Maximum
				compiled.Maximum = &value
			}
		}
	}
	return compiled
}

func compileRequestBody(body *v3.RequestBody) *BodyValidation {
	if body == nil {
		return nil
	}

	compiled := &BodyValidation{}
	if body.Required != nil {
		compiled.Required = *body.Required
	}
	if body.Content == nil {
		return compiled
	}

	mediaType := body.Content.GetOrZero("application/json")
	if mediaType == nil || mediaType.Schema == nil {
		return compiled
	}

	schema, err := mediaType.Schema.BuildSchema()
	if err != nil || schema == nil {
		return compiled
	}
	compiled.Schema = compileSchema(schema)
	return compiled
}

func compileSchema(schema *base.Schema) *SchemaValidation {
	if schema == nil {
		return nil
	}

	compiled := &SchemaValidation{
		Type:     firstType(schema),
		Required: append([]string(nil), schema.Required...),
		Enum:     decodeEnum(schema.Enum),
	}

	if schema.MinLength != nil {
		value := int(*schema.MinLength)
		compiled.MinLength = &value
	}
	if schema.MaxLength != nil {
		value := int(*schema.MaxLength)
		compiled.MaxLength = &value
	}
	if schema.Minimum != nil {
		value := *schema.Minimum
		compiled.Minimum = &value
	}
	if schema.Maximum != nil {
		value := *schema.Maximum
		compiled.Maximum = &value
	}
	if schema.Properties != nil {
		compiled.Properties = make(map[string]*SchemaValidation, schema.Properties.Len())
		for name, property := range schema.Properties.FromOldest() {
			if property == nil {
				continue
			}
			propertySchema, err := property.BuildSchema()
			if err != nil || propertySchema == nil {
				continue
			}
			compiled.Properties[name] = compileSchema(propertySchema)
		}
	}

	return compiled
}

func firstType(schema *base.Schema) string {
	if schema == nil || len(schema.Type) == 0 {
		return ""
	}
	return schema.Type[0]
}

func normalizeOpenAPIPath(path string) string {
	return openAPIPathParamPattern.ReplaceAllString(path, `:$1`)
}

func decodeEnum(values []*yaml.Node) []string {
	if len(values) == 0 {
		return nil
	}

	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == nil {
			continue
		}
		var decoded string
		if err := value.Decode(&decoded); err == nil {
			result = append(result, decoded)
			continue
		}
		result = append(result, value.Value)
	}
	return result
}
