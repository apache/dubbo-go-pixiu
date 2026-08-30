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
	for _, issue := range e {
		messages = append(messages, fmt.Sprintf("%s: %s", issue.Path, issue.Message))
	}
	return strings.Join(messages, "; ")
}

func (r *Registry) Validate(object ConfigObject) error {
	_, err := r.Normalize(object)
	return err
}

// Normalize applies registered defaults to a defensive copy and validates the
// result. The input object is never mutated.
func (r *Registry) Normalize(object ConfigObject) (ConfigObject, error) {
	normalized := object.Clone()
	if normalized.Metadata.Lifecycle == "" {
		normalized.Metadata.Lifecycle = LifecycleDraft
	}
	if normalized.Spec == nil {
		normalized.Spec = make(map[string]any)
	}

	issues := validateEnvelope(normalized)
	objectSchema, validators, exists := r.lookupEntry(normalized.APIVersion, normalized.Kind)
	if !exists {
		issues = append(issues, ValidationIssue{
			Path:    "kind",
			Code:    "schema_not_found",
			Message: fmt.Sprintf("schema %s/%s is not registered", normalized.APIVersion, normalized.Kind),
		})
		return normalized, ValidationErrors(issues)
	}

	applyDefaultsToFields(normalized.Spec, objectSchema.Fields)
	issues = append(issues, validateFields("spec", normalized.Spec, objectSchema.Fields, objectSchema.AllowUnknownFields)...)
	for _, validator := range validators {
		issues = append(issues, validator(normalized)...)
	}
	if len(issues) != 0 {
		return normalized, ValidationErrors(issues)
	}
	return normalized, nil
}

// NormalizeConfigSet validates a complete snapshot, guarantees object IDs are
// unique within each kind, and runs collection validators for relationships.
func (r *Registry) NormalizeConfigSet(configSet ConfigSet) (ConfigSet, error) {
	normalized := configSet.Clone()
	issues := make([]ValidationIssue, 0)
	if normalized.Metadata.Lifecycle == "" {
		normalized.Metadata.Lifecycle = LifecycleDraft
	}
	if strings.TrimSpace(normalized.APIVersion) == "" {
		issues = append(issues, issue("apiVersion", "required", "apiVersion is required"))
	}
	if normalized.Kind != KindConfigSet {
		issues = append(issues, issue("kind", "invalid", fmt.Sprintf("kind must be %s", KindConfigSet)))
	}
	if strings.TrimSpace(normalized.Metadata.ID) == "" {
		issues = append(issues, issue("metadata.id", "required", "metadata.id is required"))
	}
	if strings.TrimSpace(normalized.Metadata.Name) == "" {
		issues = append(issues, issue("metadata.name", "required", "metadata.name is required"))
	}
	issues = append(issues, validateLifecycle("metadata.lifecycle", normalized.Metadata.Lifecycle)...)
	if normalized.Metadata.Revision < 0 {
		issues = append(issues, issue("metadata.revision", "minimum", "revision cannot be negative"))
	}

	seen := make(map[string]int, len(normalized.Objects))
	for i := range normalized.Objects {
		if normalized.Objects[i].APIVersion == "" {
			normalized.Objects[i].APIVersion = normalized.APIVersion
		}
		object, err := r.Normalize(normalized.Objects[i])
		normalized.Objects[i] = object
		prefix := fmt.Sprintf("objects[%d].", i)
		if err != nil {
			if validationErrors, ok := err.(ValidationErrors); ok {
				for _, validationIssue := range validationErrors {
					validationIssue.Path = prefix + validationIssue.Path
					issues = append(issues, validationIssue)
				}
			} else {
				issues = append(issues, issue(prefix+"spec", "invalid", err.Error()))
			}
		}

		identity := object.Kind + "/" + object.Metadata.ID
		if previous, exists := seen[identity]; exists {
			issues = append(issues, issue(
				prefix+"metadata.id",
				"duplicate",
				fmt.Sprintf("duplicates objects[%d] identity %s", previous, identity),
			))
		} else if object.Metadata.ID != "" {
			seen[identity] = i
		}
	}
	for _, validator := range r.configSetValidatorSnapshot() {
		issues = append(issues, validator(normalized)...)
	}

	if len(issues) != 0 {
		return normalized, ValidationErrors(issues)
	}
	return normalized, nil
}

func validateEnvelope(object ConfigObject) []ValidationIssue {
	issues := make([]ValidationIssue, 0, 6)
	if strings.TrimSpace(object.APIVersion) == "" {
		issues = append(issues, issue("apiVersion", "required", "apiVersion is required"))
	}
	if strings.TrimSpace(object.Kind) == "" {
		issues = append(issues, issue("kind", "required", "kind is required"))
	}
	if strings.TrimSpace(object.Metadata.ID) == "" {
		issues = append(issues, issue("metadata.id", "required", "metadata.id is required"))
	}
	if strings.TrimSpace(object.Metadata.Name) == "" {
		issues = append(issues, issue("metadata.name", "required", "metadata.name is required"))
	}
	issues = append(issues, validateLifecycle("metadata.lifecycle", object.Metadata.Lifecycle)...)
	if object.Metadata.Revision < 0 {
		issues = append(issues, issue("metadata.revision", "minimum", "revision cannot be negative"))
	}
	return issues
}

func validateLifecycle(path string, lifecycle Lifecycle) []ValidationIssue {
	switch lifecycle {
	case LifecycleDraft, LifecyclePublished, LifecycleArchived:
		return nil
	default:
		return []ValidationIssue{issue(path, "enum", fmt.Sprintf(
			"must be one of %q, %q, or %q",
			LifecycleDraft,
			LifecyclePublished,
			LifecycleArchived,
		))}
	}
}

func validateFields(path string, values map[string]any, fields map[string]*FieldSchema, allowUnknown bool) []ValidationIssue {
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

	if !allowUnknown {
		unknown := make([]string, 0)
		for name := range values {
			if _, exists := fields[name]; !exists {
				unknown = append(unknown, name)
			}
		}
		sort.Strings(unknown)
		for _, name := range unknown {
			issues = append(issues, issue(joinPath(path, name), "unknown", "field is not registered"))
		}
	}
	return issues
}

func validateValue(path string, value any, field *FieldSchema) []ValidationIssue {
	if field == nil || field.Type == FieldTypeAny {
		return nil
	}

	switch field.Type {
	case FieldTypeString:
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
	case FieldTypeInteger:
		number, ok := numericValue(value)
		if !ok || math.Trunc(number) != number {
			return []ValidationIssue{typeIssue(path, field.Type, value)}
		}
		if rangeIssue := validateNumberRange(path, number, field); rangeIssue != nil {
			return []ValidationIssue{*rangeIssue}
		}
	case FieldTypeNumber:
		number, ok := numericValue(value)
		if !ok {
			return []ValidationIssue{typeIssue(path, field.Type, value)}
		}
		if rangeIssue := validateNumberRange(path, number, field); rangeIssue != nil {
			return []ValidationIssue{*rangeIssue}
		}
	case FieldTypeBoolean:
		if _, ok := value.(bool); !ok {
			return []ValidationIssue{typeIssue(path, field.Type, value)}
		}
	case FieldTypeObject:
		object, ok := value.(map[string]any)
		if !ok {
			return []ValidationIssue{typeIssue(path, field.Type, value)}
		}
		return validateFields(path, object, field.Properties, field.AllowUnknown)
	case FieldTypeArray:
		items, ok := value.([]any)
		if !ok {
			return []ValidationIssue{typeIssue(path, field.Type, value)}
		}
		issues := make([]ValidationIssue, 0)
		if field.MinItems != nil && len(items) < *field.MinItems {
			issues = append(issues, issue(path, "min_items", fmt.Sprintf("must contain at least %d item(s)", *field.MinItems)))
		}
		for i, item := range items {
			issues = append(issues, validateValue(fmt.Sprintf("%s[%d]", path, i), item, field.Items)...)
		}
		return issues
	case FieldTypeMap:
		values, ok := value.(map[string]any)
		if !ok {
			return []ValidationIssue{typeIssue(path, field.Type, value)}
		}
		issues := make([]ValidationIssue, 0)
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			issues = append(issues, validateValue(joinPath(path, key), values[key], field.AdditionalProperties)...)
		}
		return issues
	case FieldTypeUnion:
		object, ok := value.(map[string]any)
		if !ok {
			return []ValidationIssue{typeIssue(path, field.Type, value)}
		}
		discriminator, ok := object[field.Discriminator].(string)
		if !ok || strings.TrimSpace(discriminator) == "" {
			return []ValidationIssue{issue(joinPath(path, field.Discriminator), "required", "union discriminator is required")}
		}
		variant, exists := field.Variants[discriminator]
		if !exists {
			return []ValidationIssue{issue(joinPath(path, field.Discriminator), "unsupported", fmt.Sprintf("unsupported variant %q", discriminator))}
		}
		return validateValue(path, value, variant)
	default:
		return []ValidationIssue{issue(path, "schema", fmt.Sprintf("unsupported field type %q", field.Type))}
	}

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
	if field == nil {
		return
	}
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
		if values, ok := value.(map[string]any); ok {
			for _, item := range values {
				applyDefaultsToValue(item, field.AdditionalProperties)
			}
		}
	case FieldTypeUnion:
		if object, ok := value.(map[string]any); ok {
			if discriminator, ok := object[field.Discriminator].(string); ok {
				applyDefaultsToValue(value, field.Variants[discriminator])
			}
		}
	}
}

func validateSchemaDefinition(objectSchema ObjectSchema) error {
	if strings.TrimSpace(objectSchema.APIVersion) == "" {
		return fmt.Errorf("object schema apiVersion is required")
	}
	if strings.TrimSpace(objectSchema.Kind) == "" {
		return fmt.Errorf("object schema kind is required")
	}
	if objectSchema.Fields == nil {
		return fmt.Errorf("object schema %s/%s fields are required", objectSchema.APIVersion, objectSchema.Kind)
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
	case FieldTypeAny, FieldTypeString, FieldTypeInteger, FieldTypeNumber, FieldTypeBoolean:
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
		return validateFieldDefinition(path+"[]", field.Items)
	case FieldTypeMap:
		if field.AdditionalProperties == nil {
			return fmt.Errorf("schema field %s map value schema is required", path)
		}
		return validateFieldDefinition(path+".*", field.AdditionalProperties)
	case FieldTypeUnion:
		if strings.TrimSpace(field.Discriminator) == "" {
			return fmt.Errorf("schema field %s union discriminator is required", path)
		}
		if len(field.Variants) == 0 {
			return fmt.Errorf("schema field %s union variants are required", path)
		}
		for name, variant := range field.Variants {
			if err := validateFieldDefinition(path+"<"+name+">", variant); err != nil {
				return err
			}
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

func validateNumberRange(path string, number float64, field *FieldSchema) *ValidationIssue {
	if field.Minimum != nil && number < *field.Minimum {
		result := issue(path, "minimum", fmt.Sprintf("must be at least %v", *field.Minimum))
		return &result
	}
	if field.Maximum != nil && number > *field.Maximum {
		result := issue(path, "maximum", fmt.Sprintf("must be at most %v", *field.Maximum))
		return &result
	}
	return nil
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
