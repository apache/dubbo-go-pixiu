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

package model

import (
	"sync"
)

import (
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
)

type (
	// LLMMeta LLM metadata for llm call
	LLMMeta struct {
		Provider   string      `yaml:"provider" json:"provider"`                                              // Provider the cluster unique name
		APIKeys    []LLMAPIKey `yaml:"api_keys" json:"api_keys" mapstructure:"api_keys"`                      // APIKey the cluster unique name
		RetryTimes uint        `yaml:"retry_times" json:"retry_times" mapstructure:"retry_times" default:"0"` // Retry times for failed call
		Fallback   bool        `yaml:"fallback" json:"fallback" mapstructure:"fallback"`                      // Fallback to other provider if failed
	}

	LLMAPIKey struct {
		Name string `yaml:"name" json:"name"` // Name of the api key
		Key  string `yaml:"key" json:"key"`   // Real Key
	}

	LLMProviderDomains struct {
		Providers map[string]LLMProvider `yaml:"providers" mapstructure:"providers"`
	}

	LLMProvider struct {
		Name        string            `yaml:"name" json:"name"` // provider' name
		Description string            `yaml:"description" json:"description"`
		BaseUrl     string            `yaml:"base_url" json:"base_url"`                            // Target domain
		Endpoints   map[string]string `yaml:"endpoints" json:"endpoints" mapstructure:"endpoints"` // Endpoints for the provider
	}
)

var (
	loadLLMProviderDomains sync.Once
	domains                *LLMProviderDomains
	err                    error
)

// GetLLMProviderDomains get llm provider domains
func GetLLMProviderDomains(id string) (*LLMProvider, error) {
	loadLLMProviderDomains.Do(func() {
		domains = &LLMProviderDomains{}
		err = yaml.UnmarshalYMLConfig("pkg/model/llmprovider.yaml", domains)
	})
	if err != nil {
		return nil, perrors.Wrap(err, "failed to load llm provider domains")
	}

	if p, ok := domains.Providers[id]; ok {
		return &p, nil
	}
	return nil, perrors.Errorf("provider %s not found", id)
}
