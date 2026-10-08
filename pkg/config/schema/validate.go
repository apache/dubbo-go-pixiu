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
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

type ValidationIssue struct {
	Path    string `json:"path" yaml:"path"`
	Code    string `json:"code" yaml:"code"`
	Message string `json:"message" yaml:"message"`
}

type ValidationErrors []ValidationIssue

func (e ValidationErrors) Error() string {
	messages := make([]string, 0, len(e))
	for _, validationIssue := range e {
		messages = append(messages, fmt.Sprintf("%s: %s", validationIssue.Path, validationIssue.Message))
	}
	return strings.Join(messages, "; ")
}

func (r *Registry) Validate(object AdminObject) error {
	_, err := r.Normalize(object)
	return err
}

// Normalize applies schema defaults to a defensive copy, then validates both
// field-level and semantic rules. The caller's object is never mutated.
func (r *Registry) Normalize(object AdminObject) (AdminObject, error) {
	normalized := object.Clone()
	if normalized.Spec == nil {
		normalized.Spec = make(map[string]any)
	}
	// The publish object used to contain UI-only lifecycle preferences. Drop it
	// while normalizing older drafts so removing those no-op fields does not
	// block an existing route from publishing.
	if normalized.Kind == KindAdminRouteBinding {
		delete(normalized.Spec, "publish")
	}

	issues := validateEnvelope(normalized)
	objectSchema, validators, exists := r.lookupEntry(normalized.Kind)
	if !exists {
		issues = append(issues, issue("kind", "schema_not_found", fmt.Sprintf(
			"schema for kind %q is not registered",
			normalized.Kind,
		)))
		return normalized, ValidationErrors(issues)
	}

	applyDefaultsToFields(normalized.Spec, objectSchema.Fields)
	issues = append(issues, validateFields("spec", normalized.Spec, objectSchema.Fields, false)...)
	for _, validator := range validators {
		issues = append(issues, validator(normalized)...)
	}
	if len(issues) != 0 {
		return normalized, ValidationErrors(issues)
	}
	return normalized, nil
}

func validateEnvelope(object AdminObject) []ValidationIssue {
	issues := make([]ValidationIssue, 0, 2)
	if strings.TrimSpace(object.Kind) == "" {
		issues = append(issues, issue("kind", "required", "kind is required"))
	}
	if strings.TrimSpace(object.Metadata.Name) == "" {
		issues = append(issues, issue("metadata.name", "required", "metadata.name is required"))
	}
	return issues
}

func validateFields(path string, values map[string]any, fields map[string]*FieldSchema, allowUnknown bool) []ValidationIssue {
	issues := validateKnownFields(path, values, fields)
	if allowUnknown {
		return issues
	}
	return append(issues, validateUnknownFields(path, values, fields)...)
}

func validateKnownFields(path string, values map[string]any, fields map[string]*FieldSchema) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	names := make([]string, 0, len(fields))
	for name := range fields {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		field := fields[name]
		value, exists := values[name]
		fieldPath := joinPath(path, name)
		if (!exists || value == nil) && field.Required && field.Default == nil {
			issues = append(issues, issue(fieldPath, "required", "field is required"))
			continue
		}
		if !exists || value == nil {
			continue
		}
		issues = append(issues, validateValue(fieldPath, value, field)...)
	}
	return issues
}

func validateUnknownFields(path string, values map[string]any, fields map[string]*FieldSchema) []ValidationIssue {
	unknown := make([]string, 0)
	for name := range values {
		if _, exists := fields[name]; !exists {
			unknown = append(unknown, name)
		}
	}
	sort.Strings(unknown)
	issues := make([]ValidationIssue, 0, len(unknown))
	for _, name := range unknown {
		issues = append(issues, issue(joinPath(path, name), "unknown", "field is not registered"))
	}
	return issues
}

func validateValue(path string, value any, field *FieldSchema) []ValidationIssue {
	switch field.Type {
	case FieldTypeString:
		return validateStringValue(path, value, field)
	case FieldTypeInteger:
		return validateIntegerValue(path, value, field)
	case FieldTypeBoolean:
		return validateBooleanValue(path, value, field)
	case FieldTypeObject:
		return validateObjectValue(path, value, field)
	case FieldTypeArray:
		return validateArrayValue(path, value, field)
	case FieldTypeMap:
		return validateMapValue(path, value, field)
	default:
		return []ValidationIssue{issue(path, "schema", fmt.Sprintf("unsupported field type %q", field.Type))}
	}
}

func validateStringValue(path string, value any, field *FieldSchema) []ValidationIssue {
	valueString, ok := value.(string)
	if !ok {
		return []ValidationIssue{typeIssue(path, field.Type, value)}
	}
	if field.Pattern != "" {
		matched, _ := regexp.MatchString(field.Pattern, valueString)
		if !matched {
			return []ValidationIssue{issue(path, "pattern", fmt.Sprintf("must match %q", field.Pattern))}
		}
	}
	return validateEnumValue(path, value, field)
}

func validateIntegerValue(path string, value any, field *FieldSchema) []ValidationIssue {
	number, ok := numericValue(value)
	if !ok || math.Trunc(number) != number {
		return []ValidationIssue{typeIssue(path, field.Type, value)}
	}
	if field.Minimum != nil && number < *field.Minimum {
		return []ValidationIssue{issue(path, "minimum", fmt.Sprintf("must be at least %v", *field.Minimum))}
	}
	return validateEnumValue(path, value, field)
}

func validateBooleanValue(path string, value any, field *FieldSchema) []ValidationIssue {
	if _, ok := value.(bool); !ok {
		return []ValidationIssue{typeIssue(path, field.Type, value)}
	}
	return validateEnumValue(path, value, field)
}

func validateObjectValue(path string, value any, field *FieldSchema) []ValidationIssue {
	object, ok := value.(map[string]any)
	if !ok {
		return []ValidationIssue{typeIssue(path, field.Type, value)}
	}
	return validateFields(path, object, field.Properties, field.AllowUnknown)
}

func validateArrayValue(path string, value any, field *FieldSchema) []ValidationIssue {
	items, ok := value.([]any)
	if !ok {
		return []ValidationIssue{typeIssue(path, field.Type, value)}
	}
	issues := make([]ValidationIssue, 0)
	for index, item := range items {
		itemPath := fmt.Sprintf("%s[%d]", path, index)
		issues = append(issues, validateValue(itemPath, item, field.Items)...)
	}
	return issues
}

func validateMapValue(path string, value any, field *FieldSchema) []ValidationIssue {
	entries, ok := value.(map[string]any)
	if !ok {
		return []ValidationIssue{typeIssue(path, field.Type, value)}
	}
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	issues := make([]ValidationIssue, 0)
	for _, key := range keys {
		issues = append(issues, validateValue(joinPath(path, key), entries[key], field.AdditionalProperties)...)
	}
	return issues
}

func validateEnumValue(path string, value any, field *FieldSchema) []ValidationIssue {
	if len(field.Enum) != 0 && !enumContains(field.Enum, value) {
		return []ValidationIssue{issue(path, "enum", fmt.Sprintf("must be one of %v", field.Enum))}
	}
	return nil
}

func applyDefaultsToFields(values map[string]any, fields map[string]*FieldSchema) {
	for name, field := range fields {
		value, exists := values[name]
		if (!exists || value == nil) && field.Default != nil {
			value = cloneDefault(field.Default)
			values[name] = value
			exists = true
		}
		if !exists || value == nil {
			continue
		}
		applyDefaultsToValue(value, field)
	}
}

func applyDefaultsToValue(value any, field *FieldSchema) {
	switch field.Type {
	case FieldTypeObject:
		if object, ok := value.(map[string]any); ok {
			applyDefaultsToFields(object, field.Properties)
		}
	case FieldTypeArray:
		if items, ok := value.([]any); ok {
			for _, item := range items {
				applyDefaultsToValue(item, field.Items)
			}
		}
	case FieldTypeMap:
		if entries, ok := value.(map[string]any); ok {
			for _, item := range entries {
				applyDefaultsToValue(item, field.AdditionalProperties)
			}
		}
	}
}

func validateSchemaDefinition(objectSchema ObjectSchema) error {
	if strings.TrimSpace(objectSchema.Kind) == "" {
		return fmt.Errorf("object schema kind is required")
	}
	if objectSchema.Fields == nil {
		return fmt.Errorf("object schema %s fields are required", objectSchema.Kind)
	}
	for name, field := range objectSchema.Fields {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("object schema field name is empty")
		}
		if err := validateFieldDefinition("spec."+name, field); err != nil {
			return err
		}
	}
	return nil
}

func validateFieldDefinition(path string, field *FieldSchema) error {
	if field == nil {
		return fmt.Errorf("schema field %s is nil", path)
	}
	switch field.Type {
	case FieldTypeString, FieldTypeInteger, FieldTypeBoolean:
	case FieldTypeObject:
		for name, child := range field.Properties {
			if err := validateFieldDefinition(joinPath(path, name), child); err != nil {
				return err
			}
		}
	case FieldTypeArray:
		if field.Items == nil {
			return fmt.Errorf("schema field %s array items are required", path)
		}
		if err := validateFieldDefinition(path+"[]", field.Items); err != nil {
			return err
		}
	case FieldTypeMap:
		if field.AdditionalProperties == nil {
			return fmt.Errorf("schema field %s map value schema is required", path)
		}
		if err := validateFieldDefinition(path+".*", field.AdditionalProperties); err != nil {
			return err
		}
	default:
		return fmt.Errorf("schema field %s has unsupported type %q", path, field.Type)
	}
	if field.Pattern != "" {
		if _, err := regexp.Compile(field.Pattern); err != nil {
			return fmt.Errorf("schema field %s pattern: %w", path, err)
		}
	}
	return nil
}

func numericValue(value any) (float64, bool) {
	switch number := value.(type) {
	case int:
		return float64(number), true
	case int8:
		return float64(number), true
	case int16:
		return float64(number), true
	case int32:
		return float64(number), true
	case int64:
		return float64(number), true
	case uint:
		return float64(number), true
	case uint8:
		return float64(number), true
	case uint16:
		return float64(number), true
	case uint32:
		return float64(number), true
	case uint64:
		return float64(number), true
	case float32:
		return float64(number), true
	case float64:
		return number, true
	case json.Number:
		parsed, err := number.Float64()
		return parsed, err == nil
	default:
		return 0, false
	}
}

func enumContains(values []any, target any) bool {
	for _, value := range values {
		if reflect.DeepEqual(value, target) {
			return true
		}
	}
	return false
}

func typeIssue(path string, expected FieldType, value any) ValidationIssue {
	return issue(path, "type", fmt.Sprintf("must be %s, got %T", expected, value))
}

func issue(path, code, message string) ValidationIssue {
	return ValidationIssue{Path: path, Code: code, Message: message}
}

func joinPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}

func cloneDefault(value any) any {
	field := &FieldSchema{Default: value}
	return field.clone().Default
}
