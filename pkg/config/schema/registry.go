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
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

var (
	ErrSchemaNotFound          = errors.New("config schema not found")
	ErrSchemaAlreadyRegistered = errors.New("config schema already registered")
	ErrFieldAlreadyRegistered  = errors.New("config field already registered")
)

type ObjectValidator func(ConfigObject) []ValidationIssue

// ConfigSetValidator validates relationships that cannot be decided from one
// object alone, such as Method -> Resource and route -> Cluster references.
type ConfigSetValidator func(ConfigSet) []ValidationIssue

// Registry is an in-memory, concurrency-safe registry. Core and plugin
// schemas use the same registration path, while duplicate fields are rejected
// so extensions cannot silently override core semantics.
type Registry struct {
	mu                  sync.RWMutex
	schemas             map[string]*ObjectSchema
	validators          map[string][]ObjectValidator
	configSetValidators []ConfigSetValidator
}

func NewRegistry() *Registry {
	return &Registry{
		schemas:    make(map[string]*ObjectSchema),
		validators: make(map[string][]ObjectValidator),
	}
}

func NewBuiltinRegistry() (*Registry, error) {
	registry := NewRegistry()
	if err := RegisterBuiltinSchemas(registry); err != nil {
		return nil, err
	}
	return registry, nil
}

func (r *Registry) Register(objectSchema ObjectSchema) error {
	if err := validateSchemaDefinition(objectSchema); err != nil {
		return err
	}

	key := makeSchemaKey(objectSchema.APIVersion, objectSchema.Kind)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.schemas[key]; exists {
		return fmt.Errorf("%w: %s", ErrSchemaAlreadyRegistered, key)
	}
	cloned := objectSchema.clone()
	r.schemas[key] = &cloned
	return nil
}

// RegisterField inserts a field at a dotted path rooted at spec. Intermediate
// object fields must already exist. This makes extension ownership explicit
// and prevents a typo from creating an unintended schema branch.
func (r *Registry) RegisterField(apiVersion, kind, path string, field FieldSchema) error {
	parts, err := splitFieldPath(path)
	if err != nil {
		return err
	}
	if err := validateFieldDefinition(path, &field); err != nil {
		return err
	}

	key := makeSchemaKey(apiVersion, kind)
	r.mu.Lock()
	defer r.mu.Unlock()

	objectSchema, exists := r.schemas[key]
	if !exists {
		return fmt.Errorf("%w: %s", ErrSchemaNotFound, key)
	}

	fields := objectSchema.Fields
	for _, part := range parts[:len(parts)-1] {
		parent, exists := fields[part]
		if !exists {
			return fmt.Errorf("schema field %q does not exist", part)
		}
		if parent.Type != FieldTypeObject {
			return fmt.Errorf("schema field %q is %s, not object", part, parent.Type)
		}
		if parent.Properties == nil {
			parent.Properties = make(map[string]*FieldSchema)
		}
		fields = parent.Properties
	}

	name := parts[len(parts)-1]
	if _, exists := fields[name]; exists {
		return fmt.Errorf("%w: %s", ErrFieldAlreadyRegistered, path)
	}
	fields[name] = field.clone()
	return nil
}

func (r *Registry) RegisterValidator(apiVersion, kind string, validator ObjectValidator) error {
	if validator == nil {
		return errors.New("config object validator is nil")
	}
	key := makeSchemaKey(apiVersion, kind)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.schemas[key]; !exists {
		return fmt.Errorf("%w: %s", ErrSchemaNotFound, key)
	}
	r.validators[key] = append(r.validators[key], validator)
	return nil
}

func (r *Registry) RegisterConfigSetValidator(validator ConfigSetValidator) error {
	if validator == nil {
		return errors.New("config set validator is nil")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.configSetValidators = append(r.configSetValidators, validator)
	return nil
}

func (r *Registry) Lookup(apiVersion, kind string) (ObjectSchema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	objectSchema, exists := r.schemas[makeSchemaKey(apiVersion, kind)]
	if !exists {
		return ObjectSchema{}, false
	}
	return objectSchema.clone(), true
}

func (r *Registry) List() []ObjectSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keys := make([]string, 0, len(r.schemas))
	for key := range r.schemas {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]ObjectSchema, 0, len(keys))
	for _, key := range keys {
		result = append(result, r.schemas[key].clone())
	}
	return result
}

func (r *Registry) lookupEntry(apiVersion, kind string) (ObjectSchema, []ObjectValidator, bool) {
	key := makeSchemaKey(apiVersion, kind)
	r.mu.RLock()
	defer r.mu.RUnlock()
	objectSchema, exists := r.schemas[key]
	if !exists {
		return ObjectSchema{}, nil, false
	}
	validators := append([]ObjectValidator(nil), r.validators[key]...)
	return objectSchema.clone(), validators, true
}

func (r *Registry) configSetValidatorSnapshot() []ConfigSetValidator {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]ConfigSetValidator(nil), r.configSetValidators...)
}

func makeSchemaKey(apiVersion, kind string) string {
	return strings.TrimSpace(apiVersion) + "/" + strings.TrimSpace(kind)
}

func splitFieldPath(path string) ([]string, error) {
	path = strings.TrimSpace(strings.TrimPrefix(path, "spec."))
	if path == "" {
		return nil, errors.New("schema field path is empty")
	}
	parts := strings.Split(path, ".")
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			return nil, fmt.Errorf("schema field path %q contains an empty segment", path)
		}
	}
	return parts, nil
}
