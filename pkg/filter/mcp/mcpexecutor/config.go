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
 * distributed under the License is is "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mcpexecutor

type (
	// Config is the config for mcp_executor filter.
	Config struct {
		Tools []Tool `yaml:"tools" json:"tools"`
	}

	// Tool defines a tool that can be called by MCP.
	Tool struct {
		Name        string      `yaml:"name" json:"name"`
		Description string      `yaml:"description" json:"description"`
		Cluster     string      `yaml:"cluster" json:"cluster"`
		Request     Request     `yaml:"request" json:"request"`
		Parameters  []Parameter `yaml:"parameters" json:"parameters"`
	}

	// Request defines the HTTP request to be made.
	Request struct {
		Method       string `yaml:"method" json:"method"`
		PathTemplate string `yaml:"path_template" json:"path_template"`
	}

	// Parameter defines a parameter for a tool.
	// This structure is designed to map directly to the features of mcp.Tool builder.
	Parameter struct {
		Name        string   `yaml:"name" json:"name"`
		In          string   `yaml:"in" json:"in"`
		Description string   `yaml:"description,omitempty" json:"description,omitempty"`
		Type        string   `yaml:"type" json:"type"` // e.g., "string", "number", "integer", "boolean"
		Required    bool     `yaml:"required,omitempty" json:"required,omitempty"`
		Enum        []string `yaml:"enum,omitempty" json:"enum,omitempty"`
	}
)
