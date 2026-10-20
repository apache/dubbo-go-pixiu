/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package schema

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var supportedDubboParamTypes = []string{
	"string", "char", "short", "int", "long", "float", "double", "boolean", "byte", "date", "object",
	"java.lang.String", "java.lang.Character", "java.lang.Short", "java.lang.Integer", "java.lang.Long",
	"java.lang.Float", "java.lang.Double", "java.lang.Boolean", "java.lang.Byte", "java.lang.Object", "java.util.Date",
}

// RegisterBuiltinSchemas installs only the Admin-facing route binding. Legacy
// Resource and Method are compiler outputs rather than form objects.
func RegisterBuiltinSchemas(registry *Registry) error {
	if registry == nil {
		return fmt.Errorf("admin object schema registry is nil")
	}
	if err := registry.Register(adminRouteBindingSchema()); err != nil {
		return err
	}
	return registry.RegisterValidator(KindAdminRouteBinding, validateAdminRouteBinding)
}

func adminRouteBindingSchema() ObjectSchema {
	return ObjectSchema{
		Kind:        KindAdminRouteBinding,
		Description: "High-level HTTP entry to Dubbo invocation binding",
		Fields: map[string]*FieldSchema{
			"entry": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"protocol": {
						Type:    FieldTypeString,
						Default: "http",
						Enum:    []any{"http"},
						UI:      UIHints{Component: "select", Order: 10},
					},
					"path": {
						Type:        FieldTypeString,
						Required:    true,
						Default:     "/api/v1/example",
						Pattern:     `^/`,
						Description: "HTTP path pattern exposed by Pixiu",
						UI:          UIHints{Component: "text", Order: 20, Placeholder: "/api/v1/users/:id"},
					},
					"method": {
						Type:     FieldTypeString,
						Required: true,
						Default:  "GET",
						Enum:     []any{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
						UI:       UIHints{Component: "select", Order: 30},
					},
				},
				UI: UIHints{Group: "entry", Order: 10},
			},
			"backendRef": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"type": {
						Type:        FieldTypeString,
						Default:     "registry",
						Enum:        []any{"registry"},
						Description: "Dubbo services are resolved from the registries configured on the listener.",
						UI:          UIHints{Component: "text", Order: 10},
					},
					"clusterName": {
						Type:        FieldTypeString,
						Description: "Preserved only when migrating legacy target.cluster to integrationRequest.clusterName.",
					},
				},
				UI: UIHints{Group: "backend", Order: 20},
			},
			"serviceSelection": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"application": {Type: FieldTypeString, Required: true, UI: UIHints{Component: "text", Order: 10}},
					"interface":   {Type: FieldTypeString, Required: true, UI: UIHints{Component: "text", Order: 20}},
					"group":       {Type: FieldTypeString, Default: "", UI: UIHints{Component: "text", Order: 30}},
					"version":     {Type: FieldTypeString, Default: "", UI: UIHints{Component: "text", Order: 40}},
				},
				UI: UIHints{Group: "service", Order: 30},
			},
			"invocation": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"method":  {Type: FieldTypeString, Required: true, UI: UIHints{Component: "text", Order: 10}},
					"timeout": {Type: FieldTypeString, Default: "1s", Description: "Positive Go duration string", UI: UIHints{Component: "text", Order: 20}},
					"parameters": {
						Type:        FieldTypeArray,
						Default:     []any{},
						Description: "Ordered HTTP source to Dubbo argument mappings",
						Items: &FieldSchema{
							Type: FieldTypeObject,
							Properties: map[string]*FieldSchema{
								"from": {
									Type:        FieldTypeString,
									Required:    true,
									Pattern:     `^(uri|queryStrings|headers|requestBody)\..+`,
									Description: "legacy mappingParams.name",
								},
								"to": {
									Type:        FieldTypeInteger,
									Required:    true,
									Minimum:     floatPointer(0),
									Description: "zero-based Dubbo argument index",
								},
								"type": {
									Type:        FieldTypeString,
									Required:    true,
									Enum:        supportedDubboParamTypeEnum(),
									Description: "legacy mapType and parameterTypes entry",
								},
							},
						},
						UI: UIHints{Component: "parameter-binding-table", Group: "invocation", Order: 30},
					},
				},
				UI: UIHints{Group: "invocation", Order: 40},
			},
			"enabled": {
				Type:        FieldTypeBoolean,
				Default:     true,
				Description: "Whether Pixiu accepts requests for the published route",
				UI:          UIHints{Component: "switch", Group: "lifecycle", Order: 10},
			},
		},
	}
}

// migrateLegacyRouteBindingSpec converts the target/params shape emitted by
// PR #1046 to the separated backend/service/invocation contract. It only
// migrates the known legacy fields and rejects mixed or unknown input rather
// than silently discarding configuration.
func migrateLegacyRouteBindingSpec(object *AdminObject) []ValidationIssue {
	if extensions, exists := object.Spec["extensions"]; exists {
		values, ok := extensions.(map[string]any)
		if !ok || len(values) != 0 {
			return []ValidationIssue{issue("spec.extensions", "unsupported", "route extensions are not supported; remove this field")}
		}
		delete(object.Spec, "extensions")
	}

	_, hasTarget := object.Spec["target"]
	_, hasParams := object.Spec["params"]
	if !hasTarget && !hasParams {
		return nil
	}
	for _, field := range []string{"backendRef", "serviceSelection", "invocation"} {
		if _, exists := object.Spec[field]; exists {
			return []ValidationIssue{issue("spec."+field, "mixed_model", "cannot combine the legacy target/params fields with the separated route model")}
		}
	}
	issues := make([]ValidationIssue, 0)
	target, ok := object.Spec["target"].(map[string]any)
	if !hasTarget || !ok {
		return []ValidationIssue{issue("spec.target", "migration", "legacy target must be an object")}
	}
	allowedTargetFields := map[string]bool{"protocol": true, "application": true, "interface": true, "method": true, "version": true, "group": true, "cluster": true}
	for key := range target {
		if !allowedTargetFields[key] {
			issues = append(issues, issue("spec.target."+key, "migration", "legacy field has no AdminRouteBinding migration"))
		}
	}
	if protocol, _ := target["protocol"].(string); protocol != "" && protocol != "dubbo" {
		issues = append(issues, issue("spec.target.protocol", "migration", "only the legacy dubbo protocol can be migrated"))
	}
	if len(issues) != 0 {
		return issues
	}
	backendRef := map[string]any{"type": "registry"}
	if clusterName, exists := target["cluster"]; exists {
		backendRef["clusterName"] = clusterName
	}
	object.Spec["backendRef"] = backendRef
	object.Spec["serviceSelection"] = map[string]any{
		"application": target["application"],
		"interface":   target["interface"],
		"group":       target["group"],
		"version":     target["version"],
	}
	parameters, exists := object.Spec["params"]
	if !exists {
		parameters = []any{}
	}
	object.Spec["invocation"] = map[string]any{
		"method":     target["method"],
		"timeout":    "1s",
		"parameters": parameters,
	}
	delete(object.Spec, "target")
	delete(object.Spec, "params")
	return nil
}

func validateAdminRouteBinding(object AdminObject) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	entry, _ := object.Spec["entry"].(map[string]any)
	path, _ := entry["path"].(string)
	service, _ := object.Spec["serviceSelection"].(map[string]any)
	invocation, _ := object.Spec["invocation"].(map[string]any)

	issues = append(issues, validateRouteBindingService(service, invocation)...)
	if params, ok := invocation["parameters"].([]any); ok {
		issues = append(issues, validateRouteBindingParameters(params, routePathParameters(path))...)
	}
	return issues
}

func validateRouteBindingService(service, invocation map[string]any) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	for _, field := range []string{"application", "interface"} {
		value, _ := service[field].(string)
		if strings.TrimSpace(value) == "" {
			issues = append(issues, issue("spec.serviceSelection."+field, "required", "must not be empty"))
		}
	}
	method, _ := invocation["method"].(string)
	if strings.TrimSpace(method) == "" {
		issues = append(issues, issue("spec.invocation.method", "required", "must not be empty"))
	}
	timeout, _ := invocation["timeout"].(string)
	duration, err := time.ParseDuration(timeout)
	if err != nil || duration <= 0 {
		issues = append(issues, issue("spec.invocation.timeout", "duration", "must be a positive duration such as 500ms or 2s"))
	}
	return issues
}

func validateRouteBindingParameters(params []any, pathParams map[string]struct{}) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	seenIndexes := make(map[int]int, len(params))
	usedPathParams := make(map[string]int)
	for index, rawParam := range params {
		param, ok := rawParam.(map[string]any)
		if !ok {
			continue
		}
		issues = append(issues, validateRouteBindingParameter(index, param, pathParams, seenIndexes, usedPathParams)...)
	}
	issues = append(issues, validatePathParameterMappings(pathParams, usedPathParams)...)
	issues = append(issues, validateContiguousParameterIndexes(params, seenIndexes)...)
	return issues
}

func validateRouteBindingParameter(index int, param map[string]any, pathParams map[string]struct{}, seenIndexes map[int]int, usedPathParams map[string]int) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	toNumber, ok := numericValue(param["to"])
	if !ok {
		return issues
	}
	to := int(toNumber)
	if previous, exists := seenIndexes[to]; exists {
		issues = append(issues, issue(
			fmt.Sprintf("spec.invocation.parameters[%d].to", index),
			"duplicate",
			fmt.Sprintf("duplicates params[%d] argument index %d", previous, to),
		))
	} else {
		seenIndexes[to] = index
	}

	if paramType, ok := param["type"].(string); ok && !isSupportedParamType(paramType) {
		issues = append(issues, issue(
			fmt.Sprintf("spec.invocation.parameters[%d].type", index),
			"unsupported",
			fmt.Sprintf("unsupported scalar Dubbo map type %q", paramType),
		))
	}
	if source, ok := param["from"].(string); ok && strings.HasPrefix(source, "uri.") {
		issues = append(issues, validatePathParameterSource(index, source, pathParams, usedPathParams)...)
	}
	return issues
}

func validatePathParameterSource(index int, source string, pathParams map[string]struct{}, usedPathParams map[string]int) []ValidationIssue {
	name := strings.TrimPrefix(source, "uri.")
	if _, exists := pathParams[name]; !exists {
		return []ValidationIssue{issue(fmt.Sprintf("spec.invocation.parameters[%d].from", index), "path_mismatch", fmt.Sprintf("URI parameter %q does not appear in spec.entry.path", name))}
	}
	if previous, exists := usedPathParams[name]; exists {
		return []ValidationIssue{issue(fmt.Sprintf("spec.invocation.parameters[%d].from", index), "duplicate_uri", fmt.Sprintf("URI parameter %q is already mapped by parameters[%d]", name, previous))}
	}
	usedPathParams[name] = index
	return nil
}

func validatePathParameterMappings(pathParams map[string]struct{}, usedPathParams map[string]int) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	for name := range pathParams {
		if _, exists := usedPathParams[name]; !exists {
			issues = append(issues, issue("spec.invocation.parameters", "unmapped_uri", fmt.Sprintf("path parameter %q must have a uri.%s mapping", name, name)))
		}
	}
	return issues
}

func validateContiguousParameterIndexes(params []any, seenIndexes map[int]int) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	for expected := 0; expected < len(params); expected++ {
		if _, exists := seenIndexes[expected]; !exists {
			issues = append(issues, issue(
				"spec.invocation.parameters",
				"non_contiguous",
				fmt.Sprintf("argument indexes must be contiguous from 0; index %d is missing", expected),
			))
		}
	}
	return issues
}

var routePathParameterPattern = regexp.MustCompile(`:([A-Za-z_][A-Za-z0-9_]*)`)

func routePathParameters(path string) map[string]struct{} {
	parameters := make(map[string]struct{})
	for _, match := range routePathParameterPattern.FindAllStringSubmatch(path, -1) {
		parameters[match[1]] = struct{}{}
	}
	return parameters
}

func isSupportedParamType(value string) bool {
	for _, supportedType := range supportedDubboParamTypes {
		if value == supportedType {
			return true
		}
	}
	return false
}

func supportedDubboParamTypeEnum() []any {
	values := make([]any, len(supportedDubboParamTypes))
	for index, value := range supportedDubboParamTypes {
		values[index] = value
	}
	return values
}

func floatPointer(value float64) *float64 {
	return &value
}
