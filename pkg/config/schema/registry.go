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
	ErrSchemaNotFound          = errors.New("admin object schema not found")
	ErrSchemaAlreadyRegistered = errors.New("admin object schema already registered")
)

type ObjectValidator func(AdminObject) []ValidationIssue

// Registry stores Admin-facing schemas and semantic validators. Definitions
// are cloned at the boundary so callers cannot mutate active validation rules.
type Registry struct {
	mu         sync.RWMutex
	schemas    map[string]*ObjectSchema
	validators map[string][]ObjectValidator
}

var defaultRegistryState struct {
	sync.Once
	registry *Registry
	err      error
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

// DefaultRegistry returns the process-wide registry used by Admin route
// handlers. Extensions registered during process startup are therefore
// visible to both schema discovery and route validation.
func DefaultRegistry() (*Registry, error) {
	defaultRegistryState.Do(func() {
		defaultRegistryState.registry, defaultRegistryState.err = NewBuiltinRegistry()
	})
	return defaultRegistryState.registry, defaultRegistryState.err
}

func (r *Registry) Register(objectSchema ObjectSchema) error {
	if err := validateSchemaDefinition(objectSchema); err != nil {
		return err
	}

	kind := strings.TrimSpace(objectSchema.Kind)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.schemas[kind]; exists {
		return fmt.Errorf("%w: %s", ErrSchemaAlreadyRegistered, kind)
	}
	cloned := objectSchema.clone()
	r.schemas[kind] = &cloned
	return nil
}

func (r *Registry) RegisterValidator(kind string, validator ObjectValidator) error {
	if validator == nil {
		return errors.New("admin object validator is nil")
	}
	kind = strings.TrimSpace(kind)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.schemas[kind]; !exists {
		return fmt.Errorf("%w: %s", ErrSchemaNotFound, kind)
	}
	r.validators[kind] = append(r.validators[kind], validator)
	return nil
}

func (r *Registry) Lookup(kind string) (ObjectSchema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	objectSchema, exists := r.schemas[strings.TrimSpace(kind)]
	if !exists {
		return ObjectSchema{}, false
	}
	return objectSchema.clone(), true
}

func (r *Registry) List() []ObjectSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()
	kinds := make([]string, 0, len(r.schemas))
	for kind := range r.schemas {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	result := make([]ObjectSchema, 0, len(kinds))
	for _, kind := range kinds {
		result = append(result, r.schemas[kind].clone())
	}
	return result
}

func (r *Registry) lookupEntry(kind string) (ObjectSchema, []ObjectValidator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	kind = strings.TrimSpace(kind)
	objectSchema, exists := r.schemas[kind]
	if !exists {
		return ObjectSchema{}, nil, false
	}
	validators := append([]ObjectValidator(nil), r.validators[kind]...)
	return objectSchema.clone(), validators, true
}
