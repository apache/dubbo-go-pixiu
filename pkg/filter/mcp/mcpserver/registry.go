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

package mcpserver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// ToolRegistry tool registry, thread-safe (optimized with single indexing)
type ToolRegistry struct {
	mu                sync.RWMutex
	toolSnapshot      atomic.Value                            // *ToolCatalogSnapshot
	resources         map[string]model.ResourceConfig         // indexed by URI
	resourceTemplates map[string]model.ResourceTemplateConfig // indexed by name
	prompts           map[string]model.PromptConfig
}

// ToolCatalogSnapshot is an immutable, versioned view of the tool catalog.
// The registry publishes a fully built snapshot atomically after validation.
type ToolCatalogSnapshot struct {
	Version     string
	Generation  uint64
	Fingerprint string
	ordered     []model.ToolConfig
	byName      map[string]model.ToolConfig
}

// EmptyFingerprint is the stable fingerprint for an empty tool catalog.
const EmptyFingerprint = "00000000"

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	r := &ToolRegistry{
		resources:         make(map[string]model.ResourceConfig),
		resourceTemplates: make(map[string]model.ResourceTemplateConfig),
		prompts:           make(map[string]model.PromptConfig),
	}
	r.toolSnapshot.Store(&ToolCatalogSnapshot{
		Version:     "0:" + EmptyFingerprint,
		Fingerprint: EmptyFingerprint,
		byName:      map[string]model.ToolConfig{},
	})
	return r
}

// RegisterTool registers a tool
func (r *ToolRegistry) RegisterTool(tool model.ToolConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	snap := r.snapshotUnsafe()
	if _, exists := snap.byName[tool.Name]; exists {
		return fmt.Errorf("tool %s already exists", tool.Name)
	}
	next := append(snap.orderedToolsUnsafe(), *tool.DeepCopy())
	newSnap, err := buildToolCatalogSnapshot(next, snap.Generation+1, "")
	if err != nil {
		return err
	}
	r.toolSnapshot.Store(newSnap)
	return nil
}

// ReplaceAllTools replaces the entire tools set with the provided slice (full sync).
// It validates the full replacement first and leaves the current registry
// unchanged when duplicate names are present.
func (r *ToolRegistry) ReplaceAllTools(tools []model.ToolConfig) error {
	return r.replaceAllToolsWithFingerprint(tools, "")
}

func (r *ToolRegistry) replaceAllToolsWithFingerprint(tools []model.ToolConfig, fingerprint string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	old := r.snapshotUnsafe()
	newSnap, err := buildToolCatalogSnapshot(tools, old.Generation+1, fingerprint)
	if err != nil {
		return err
	}
	r.toolSnapshot.Store(newSnap)
	return nil
}

// RegisterResource registers a resource (indexed by URI as per MCP specification)
func (r *ToolRegistry) RegisterResource(resource model.ResourceConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[resource.URI]; exists {
		return fmt.Errorf("resource with URI %s already exists", resource.URI)
	}

	// Register resource by URI as per MCP specification
	r.resources[resource.URI] = resource
	return nil
}

// GetTool gets tool configuration
func (r *ToolRegistry) GetTool(name string) (model.ToolConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.snapshotUnsafe().byName[name]
	if !exists {
		return model.ToolConfig{}, false
	}
	return *tool.DeepCopy(), true
}

// GetResourceByURI gets resource configuration (by URI, O(1) lookup)
func (r *ToolRegistry) GetResourceByURI(uri string) (model.ResourceConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Direct O(1) lookup by URI
	resource, exists := r.resources[uri]
	return resource, exists
}

// RegisterResourceTemplate registers a resource template
func (r *ToolRegistry) RegisterResourceTemplate(template model.ResourceTemplateConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resourceTemplates[template.Name]; exists {
		return fmt.Errorf("resource template %s already exists", template.Name)
	}

	r.resourceTemplates[template.Name] = template
	return nil
}

// GetResourceTemplate gets resource template configuration
func (r *ToolRegistry) GetResourceTemplate(name string) (model.ResourceTemplateConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.resourceTemplates[name]
	return template, exists
}

// ListResourceTemplates lists all resource templates
func (r *ToolRegistry) ListResourceTemplates() []model.ResourceTemplateConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	templates := make([]model.ResourceTemplateConfig, 0, len(r.resourceTemplates))
	for _, template := range r.resourceTemplates {
		templates = append(templates, template)
	}
	return templates
}

// ListTools lists all tools
func (r *ToolRegistry) ListTools() []model.ToolConfig {
	return r.ToolCatalogSnapshot().OrderedTools()
}

// ToolSnapshot returns a stable ordered copy of the current tool catalog with
// the catalog version computed on the last registry mutation.
func (r *ToolRegistry) ToolSnapshot() ([]model.ToolConfig, string) {
	snap := r.ToolCatalogSnapshot()
	return snap.OrderedTools(), snap.Version
}

// ToolCatalogSnapshot returns the immutable current snapshot. Package-internal
// hot paths may read its unexported fields without copying.
func (r *ToolRegistry) ToolCatalogSnapshot() *ToolCatalogSnapshot {
	return r.snapshotUnsafe().clone()
}

func (r *ToolRegistry) toolCatalogSnapshotUnsafe() *ToolCatalogSnapshot {
	return r.snapshotUnsafe()
}

func (r *ToolRegistry) snapshotUnsafe() *ToolCatalogSnapshot {
	snap, _ := r.toolSnapshot.Load().(*ToolCatalogSnapshot)
	if snap == nil {
		return &ToolCatalogSnapshot{Version: "0:" + EmptyFingerprint, Fingerprint: EmptyFingerprint, byName: map[string]model.ToolConfig{}}
	}
	return snap
}

func buildToolCatalogSnapshot(tools []model.ToolConfig, generation uint64, fingerprint string) (*ToolCatalogSnapshot, error) {
	if fingerprint == "" {
		fingerprint = toolCatalogFingerprint(tools)
	}
	ordered := make([]model.ToolConfig, len(tools))
	byName := make(map[string]model.ToolConfig, len(tools))
	for i := range tools {
		t := *tools[i].DeepCopy()
		if t.Name == "" {
			return nil, fmt.Errorf("tool name is required at index %d", i)
		}
		if _, exists := byName[t.Name]; exists {
			return nil, fmt.Errorf("duplicate tool name %q at index %d", t.Name, i)
		}
		ordered[i] = t
		byName[t.Name] = t
	}
	return &ToolCatalogSnapshot{
		Version:     fmt.Sprintf("%d:%s", generation, firstNonEmptyFingerprint(fingerprint)),
		Generation:  generation,
		Fingerprint: firstNonEmptyFingerprint(fingerprint),
		ordered:     ordered,
		byName:      byName,
	}, nil
}

func firstNonEmptyFingerprint(fingerprint string) string {
	if fingerprint == "" {
		return EmptyFingerprint
	}
	return fingerprint
}

// OrderedTools returns a deep copy of the ordered catalog for callers outside
// hot authorization paths.
func (s *ToolCatalogSnapshot) OrderedTools() []model.ToolConfig {
	if s == nil || len(s.ordered) == 0 {
		return nil
	}
	out := make([]model.ToolConfig, len(s.ordered))
	for i := range s.ordered {
		out[i] = *s.ordered[i].DeepCopy()
	}
	return out
}

func (s *ToolCatalogSnapshot) clone() *ToolCatalogSnapshot {
	if s == nil {
		return &ToolCatalogSnapshot{Version: "0:" + EmptyFingerprint, Fingerprint: EmptyFingerprint, byName: map[string]model.ToolConfig{}}
	}
	cp := &ToolCatalogSnapshot{
		Version:     s.Version,
		Generation:  s.Generation,
		Fingerprint: s.Fingerprint,
		ordered:     s.OrderedTools(),
		byName:      make(map[string]model.ToolConfig, len(s.byName)),
	}
	for name, tool := range s.byName {
		cp.byName[name] = *tool.DeepCopy()
	}
	return cp
}

func (s *ToolCatalogSnapshot) orderedToolsUnsafe() []model.ToolConfig {
	if s == nil {
		return nil
	}
	return s.ordered
}

func (s *ToolCatalogSnapshot) lookup(name string) (model.ToolConfig, bool) {
	if s == nil {
		return model.ToolConfig{}, false
	}
	t, ok := s.byName[name]
	return t, ok
}

func toolCatalogFingerprint(tools []model.ToolConfig) string {
	if len(tools) == 0 {
		return EmptyFingerprint
	}

	items := make([]string, len(tools))
	for i, tool := range tools {
		data, err := json.Marshal(tool)
		if err != nil {
			data = []byte(fmt.Sprintf("%#v", tool))
		}
		items[i] = string(data)
	}
	sort.Strings(items)

	hash := sha256.New()
	for _, item := range items {
		_, _ = hash.Write([]byte(item))
		_, _ = hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))[:8]
}

// ListResources lists all resources
func (r *ToolRegistry) ListResources() []model.ResourceConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resources := make([]model.ResourceConfig, 0, len(r.resources))
	for _, resource := range r.resources {
		resources = append(resources, resource)
	}
	return resources
}

// ToMCPTools converts tool configurations to tool list
func (r *ToolRegistry) ToMCPTools() ([]map[string]any, error) {
	tools := r.ListTools()
	mcpTools := make([]map[string]any, 0, len(tools))

	for _, tool := range tools {
		// Build tool according to MCP protocol specification
		mcpTool := map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": r.convertToInputSchema(tool),
		}
		mcpTools = append(mcpTools, mcpTool)
	}

	return mcpTools, nil
}

// convertToInputSchema converts tool parameters to MCP inputSchema format
func (r *ToolRegistry) convertToInputSchema(tool model.ToolConfig) map[string]any {
	allParams, err := tool.GetAllParameters()
	if err != nil {
		logger.Errorf("failed to get parameters for tool %s: %v", tool.Name, err)
		return map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}
	}

	properties := make(map[string]any)
	required := make([]string, 0)

	for _, param := range allParams {
		propSchema := map[string]any{
			"type":        param.Type,
			"description": param.Description,
		}

		if len(param.Enum) > 0 {
			propSchema["enum"] = param.Enum
		}

		if param.Default != nil {
			propSchema["default"] = param.Default
		}

		properties[param.Name] = propSchema

		if param.Required {
			required = append(required, param.Name)
		}
	}

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// ToMCPResources converts resource configurations to MCP resource list using mcp-go structures
func (r *ToolRegistry) ToMCPResources() ([]mcp.Resource, error) {
	resources := r.ListResources()
	mcpResources := make([]mcp.Resource, 0, len(resources))

	for _, resource := range resources {
		// Use mcp-go Resource structure
		mcpResource := mcp.Resource{
			URI:         resource.URI,
			Name:        resource.Name,
			Description: resource.Description,
			MIMEType:    resource.MIMEType,
		}
		mcpResources = append(mcpResources, mcpResource)
	}

	return mcpResources, nil
}

// ToMCPResourceTemplates converts resource template configurations to MCP resource template list
func (r *ToolRegistry) ToMCPResourceTemplates() ([]map[string]any, error) {
	templates := r.ListResourceTemplates()
	mcpTemplates := make([]map[string]any, 0, len(templates))

	for _, template := range templates {
		mcpTemplate := map[string]any{
			"uriTemplate": template.URITemplate,
			"name":        template.Name,
			"description": template.Description,
			"mimeType":    template.MIMEType,
		}

		// Add optional fields
		if template.Title != "" {
			mcpTemplate["title"] = template.Title
		}

		// Add annotations
		if template.Annotations != nil {
			annotations := make(map[string]any)
			if len(template.Annotations.Audience) > 0 {
				annotations["audience"] = template.Annotations.Audience
			}
			if template.Annotations.Priority != nil {
				annotations["priority"] = *template.Annotations.Priority
			}
			if template.Annotations.LastModified != "" {
				annotations["lastModified"] = template.Annotations.LastModified
			}
			if len(annotations) > 0 {
				mcpTemplate["annotations"] = annotations
			}
		}

		mcpTemplates = append(mcpTemplates, mcpTemplate)
	}

	return mcpTemplates, nil
}

// RegisterPrompt registers a prompt
func (r *ToolRegistry) RegisterPrompt(prompt model.PromptConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prompts[prompt.Name]; exists {
		return fmt.Errorf("prompt %s already exists", prompt.Name)
	}

	r.prompts[prompt.Name] = prompt
	return nil
}

// GetPrompt gets a prompt
func (r *ToolRegistry) GetPrompt(name string) (model.PromptConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompt, exists := r.prompts[name]
	return prompt, exists
}

// ListPrompts lists all prompts
func (r *ToolRegistry) ListPrompts() []model.PromptConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompts := make([]model.PromptConfig, 0, len(r.prompts))
	for _, prompt := range r.prompts {
		prompts = append(prompts, prompt)
	}

	return prompts
}

// ToMCPPrompts converts prompt configurations to MCP prompt list
func (r *ToolRegistry) ToMCPPrompts() ([]map[string]any, error) {
	prompts := r.ListPrompts()
	mcpPrompts := make([]map[string]any, 0, len(prompts))

	for _, prompt := range prompts {
		mcpPrompt := map[string]any{
			"name":        prompt.Name,
			"description": prompt.Description,
		}

		if prompt.Title != "" {
			mcpPrompt["title"] = prompt.Title
		}

		if len(prompt.Arguments) > 0 {
			args := make([]map[string]any, 0, len(prompt.Arguments))
			for _, arg := range prompt.Arguments {
				argMap := map[string]any{
					"name":        arg.Name,
					"description": arg.Description,
					"required":    arg.Required,
				}
				args = append(args, argMap)
			}
			mcpPrompt["arguments"] = args
		}

		mcpPrompts = append(mcpPrompts, mcpPrompt)
	}

	return mcpPrompts, nil
}
