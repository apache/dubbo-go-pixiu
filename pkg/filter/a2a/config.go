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

package a2a

// Config represents the configuration for the A2A Server Filter
type Config struct {
	// Endpoint specifies the A2A service endpoint path
	Endpoint string `yaml:"endpoint" json:"endpoint" default:"/a2a"`

	// AgentInfo contains information about the current agent
	AgentInfo *AgentConfig `yaml:"agent_info" json:"agent_info"`

	// KnownAgents contains a list of statically configured known agents
	KnownAgents []AgentConfig `yaml:"known_agents" json:"known_agents,omitempty"`

	// TaskConfig contains task management related configuration
	TaskConfig *TaskConfig `yaml:"task_config" json:"task_config,omitempty"`
}

// AgentConfig represents configuration for an agent
type AgentConfig struct {
	// AgentID is the unique identifier for the agent
	AgentID string `yaml:"agent_id" json:"agent_id"`

	// Name is the human-readable name of the agent
	Name string `yaml:"name" json:"name"`

	// Version is the version of the agent
	Version string `yaml:"version" json:"version" default:"1.0.0"`

	// Description is an optional description of the agent
	Description string `yaml:"description" json:"description,omitempty"`

	// Endpoint is the HTTP endpoint where the agent can be reached
	Endpoint string `yaml:"endpoint" json:"endpoint"`

	// Status is the current status of the agent
	Status AgentStatus `yaml:"status" json:"status" default:"online"`

	// Capabilities lists the capabilities that the agent provides
	Capabilities []CapabilityConfig `yaml:"capabilities" json:"capabilities,omitempty"`

	// Metadata contains additional metadata about the agent
	Metadata map[string]interface{} `yaml:"metadata" json:"metadata,omitempty"`
}

// CapabilityConfig represents configuration for an agent capability
type CapabilityConfig struct {
	// Name is the unique name of the capability
	Name string `yaml:"name" json:"name"`

	// Description is an optional description of the capability
	Description string `yaml:"description" json:"description,omitempty"`

	// InputTypes specifies the types of input this capability accepts
	InputTypes []string `yaml:"input_types" json:"input_types,omitempty"`

	// OutputTypes specifies the types of output this capability produces
	OutputTypes []string `yaml:"output_types" json:"output_types,omitempty"`

	// Tags are labels that can be used to categorize and discover capabilities
	Tags []string `yaml:"tags" json:"tags,omitempty"`

	// Parameters defines the parameters that this capability accepts
	Parameters []ParameterConfig `yaml:"parameters" json:"parameters,omitempty"`
}

// ParameterConfig represents configuration for a capability parameter
type ParameterConfig struct {
	// Name is the parameter name
	Name string `yaml:"name" json:"name"`

	// Type is the parameter type (string, number, boolean, object, array)
	Type string `yaml:"type" json:"type" default:"string"`

	// Description is an optional description of the parameter
	Description string `yaml:"description" json:"description,omitempty"`

	// Required indicates whether this parameter is required
	Required bool `yaml:"required" json:"required" default:"false"`

	// Default is the default value for the parameter
	Default interface{} `yaml:"default" json:"default,omitempty"`

	// Enum lists the allowed values for the parameter (if applicable)
	Enum []interface{} `yaml:"enum" json:"enum,omitempty"`
}

// TaskConfig represents configuration for task management
type TaskConfig struct {
	// DefaultTimeout is the default timeout for tasks in milliseconds
	DefaultTimeout int64 `yaml:"default_timeout" json:"default_timeout" default:"30000"`

	// MaxConcurrentTasks is the maximum number of concurrent tasks
	MaxConcurrentTasks int `yaml:"max_concurrent_tasks" json:"max_concurrent_tasks" default:"100"`

	// CleanupInterval is the interval for cleaning up completed/expired tasks in milliseconds
	CleanupInterval int64 `yaml:"cleanup_interval" json:"cleanup_interval" default:"300000"`

	// TaskRetention is how long to keep completed tasks in milliseconds
	TaskRetention int64 `yaml:"task_retention" json:"task_retention" default:"3600000"`
}

// GetDefaultAgentInfo returns a default agent configuration
func (c *Config) GetDefaultAgentInfo() *AgentInfo {
	if c.AgentInfo == nil {
		return &AgentInfo{
			AgentID:      "default-a2a-agent",
			Name:         "Default A2A Agent",
			Version:      DefaultAgentVersion,
			Description:  "A default A2A agent",
			Endpoint:     c.Endpoint,
			Status:       StatusOnline,
			Capabilities: []Capability{},
			Metadata:     make(map[string]interface{}),
		}
	}

	return c.AgentInfo.ToAgentInfo(c.Endpoint)
}

// ToAgentInfo converts AgentConfig to AgentInfo
func (ac *AgentConfig) ToAgentInfo(defaultEndpoint string) *AgentInfo {
	endpoint := ac.Endpoint
	if endpoint == "" {
		endpoint = defaultEndpoint
	}

	capabilities := make([]Capability, len(ac.Capabilities))
	for i, cap := range ac.Capabilities {
		capabilities[i] = cap.ToCapability()
	}

	return &AgentInfo{
		AgentID:      ac.AgentID,
		Name:         ac.Name,
		Version:      ac.Version,
		Description:  ac.Description,
		Endpoint:     endpoint,
		Status:       ac.Status,
		Capabilities: capabilities,
		Metadata:     ac.Metadata,
	}
}

// ToCapability converts CapabilityConfig to Capability
func (cc *CapabilityConfig) ToCapability() Capability {
	parameters := make([]Parameter, len(cc.Parameters))
	for i, param := range cc.Parameters {
		parameters[i] = Parameter{
			Name:        param.Name,
			Type:        param.Type,
			Description: param.Description,
			Required:    param.Required,
			Default:     param.Default,
		}
	}

	return Capability{
		Name:        cc.Name,
		Description: cc.Description,
		InputTypes:  cc.InputTypes,
		OutputTypes: cc.OutputTypes,
		Tags:        cc.Tags,
		Parameters:  parameters,
	}
}

// GetTaskConfig returns the task configuration or default values
func (c *Config) GetTaskConfig() *TaskConfig {
	if c.TaskConfig == nil {
		return &TaskConfig{
			DefaultTimeout:     DefaultTaskTimeout,
			MaxConcurrentTasks: DefaultMaxConcurrentTasks,
			CleanupInterval:    DefaultTaskCleanupInterval,
			TaskRetention:      DefaultTaskCleanupInterval * 12, // 1 hour default
		}
	}
	return c.TaskConfig
}
