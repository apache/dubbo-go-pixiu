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

import "github.com/apache/dubbo-go-pixiu/pkg/common/copyutil"

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInteger FieldType = "integer"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeObject  FieldType = "object"
	FieldTypeArray   FieldType = "array"
	FieldTypeMap     FieldType = "map"
)

// ObjectSchema is the source of truth for validation and form generation for
// one Admin-facing object kind.
type ObjectSchema struct {
	Kind        string                  `json:"kind" yaml:"kind"`
	Description string                  `json:"description,omitempty" yaml:"description,omitempty"`
	Fields      map[string]*FieldSchema `json:"fields" yaml:"fields"`
}

// FieldSchema intentionally covers the smaller contract needed by Admin
// forms instead of attempting to reimplement the full JSON Schema standard.
type FieldSchema struct {
	Type                 FieldType               `json:"type" yaml:"type"`
	Description          string                  `json:"description,omitempty" yaml:"description,omitempty"`
	Required             bool                    `json:"required,omitempty" yaml:"required,omitempty"`
	Default              any                     `json:"default,omitempty" yaml:"default,omitempty"`
	Enum                 []any                   `json:"enum,omitempty" yaml:"enum,omitempty"`
	Pattern              string                  `json:"pattern,omitempty" yaml:"pattern,omitempty"`
	Minimum              *float64                `json:"minimum,omitempty" yaml:"minimum,omitempty"`
	Properties           map[string]*FieldSchema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items                *FieldSchema            `json:"items,omitempty" yaml:"items,omitempty"`
	AdditionalProperties *FieldSchema            `json:"additionalProperties,omitempty" yaml:"additionalProperties,omitempty"`
	AllowUnknown         bool                    `json:"allowUnknown,omitempty" yaml:"allowUnknown,omitempty"`
	UI                   UIHints                 `json:"ui,omitempty" yaml:"ui,omitempty"`
}

type UIHints struct {
	Component   string         `json:"component,omitempty" yaml:"component,omitempty"`
	Group       string         `json:"group,omitempty" yaml:"group,omitempty"`
	Order       int            `json:"order,omitempty" yaml:"order,omitempty"`
	Placeholder string         `json:"placeholder,omitempty" yaml:"placeholder,omitempty"`
	Advanced    bool           `json:"advanced,omitempty" yaml:"advanced,omitempty"`
	Options     map[string]any `json:"options,omitempty" yaml:"options,omitempty"`
}

func (s ObjectSchema) clone() ObjectSchema {
	cloned := s
	cloned.Fields = cloneFieldMap(s.Fields)
	return cloned
}

func (s *FieldSchema) clone() *FieldSchema {
	if s == nil {
		return nil
	}
	cloned := *s
	cloned.Default = copyutil.CloneJSONLike(s.Default)
	cloned.Enum = make([]any, len(s.Enum))
	for i := range s.Enum {
		cloned.Enum[i] = copyutil.CloneJSONLike(s.Enum[i])
	}
	cloned.Properties = cloneFieldMap(s.Properties)
	cloned.Items = s.Items.clone()
	cloned.AdditionalProperties = s.AdditionalProperties.clone()
	cloned.UI.Options = copyutil.CloneStringAnyMap(s.UI.Options)
	if s.Minimum != nil {
		minimum := *s.Minimum
		cloned.Minimum = &minimum
	}
	return &cloned
}

func cloneFieldMap(fields map[string]*FieldSchema) map[string]*FieldSchema {
	if fields == nil {
		return nil
	}
	cloned := make(map[string]*FieldSchema, len(fields))
	for name, field := range fields {
		cloned[name] = field.clone()
	}
	return cloned
}
