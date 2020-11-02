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

package config

import (
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/yaml"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"sync"
)

var (
	apiConfig *APIConfig
	once      sync.Once
)

// HTTPVerb defines the restful api http verb
type HTTPVerb string

const (
	// MethodAny any method
	MethodAny HTTPVerb = "ANY"
	// MethodGet get
	MethodGet HTTPVerb = "GET"
	// MethodHead head
	MethodHead HTTPVerb = "HEAD"
	// MethodPost post
	MethodPost HTTPVerb = "POST"
	// MethodPut put
	MethodPut HTTPVerb = "PUT"
	// MethodPatch patch
	MethodPatch HTTPVerb = "PATCH" // RFC 5789
	// MethodDelete delete
	MethodDelete HTTPVerb = "DELETE"
	// MethodOptions options
	MethodOptions HTTPVerb = "OPTIONS"
)

// RequestType describes the type of the request. could be DUBBO/HTTP and others that we might implement in the future
type RequestType string

const (
	// DubboRequest represents the dubbo request
	DubboRequest RequestType = "dubbo"
	// HTTPRequest represents the http request
	HTTPRequest RequestType = "http"
)

// APIConfig defines the data structure of the api gateway configuration
type APIConfig struct {
	Name        string       `json:"name" yaml:"name"`
	Description string       `json:"description" yaml:"description"`
	Resources   []Resource   `json:"resources" yaml:"resources"`
	Definitions []Definition `json:"definitions" yaml:"definitions"`
}

// Resource defines the API path
type Resource struct {
	Type        string     `json:"type" yaml:"type"` // Restful, Dubbo
	Path        string     `json:"path" yaml:"path"`
	Description string     `json:"description" yaml:"description"`
	Filters     []string   `json:"filters" yaml:"filters"`
	Methods     []Method   `json:"methods" yaml:"methods"`
	Resources   []Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// Method defines the method of the api
type Method struct {
	OnAir              bool     `json:"onAir" yaml:"onAir"` // true means the method is up and false means method is down
	Filters            []string `json:"filters" yaml:"filters"`
	HTTPVerb           `json:"httpVerb" yaml:"httpVerb"`
	InboundRequest     `json:"inboundRequest" yaml:"inboundRequest"`
	IntegrationRequest `json:"integrationRequest" yaml:"integrationRequest"`
}

// InboundRequest defines the details of the inbound
type InboundRequest struct {
	RequestType  `json:"requestType" yaml:"requestType"` //http, TO-DO: dubbo
	Headers      []Params                                `json:"headers" yaml:"headers"`
	QueryStrings []Params                                `json:"queryStrings" yaml:"queryStrings"`
	RequestBody  []BodyDefinition                        `json:"requestBody" yaml:"requestBody"`
}

// Params defines the simple parameter definition
type Params struct {
	Name     string `json:"name" yaml:"name"`
	Required bool   `json:"required" yaml:"required"`
}

// BodyDefinition connects the request body to the definitions
type BodyDefinition struct {
	DefinitionName string `json:"definitionName" yaml:"definitionName"`
}

// IntegrationRequest defines the backend request format and target
type IntegrationRequest struct {
	RequestType        `json:"requestType" yaml:"requestType"` // dubbo, TO-DO: http
	DubboBackendConfig `json:"dubboBackendConfig,inline,omitempty" yaml:"dubboBackendConfig,inline,omitempty"`
	HTTPBackendConfig  `json:"httpBackendConfig" yaml:"httpBackendConfig"`
	MappingParams      []MappingParam `json:"mappingParams,omitempty" yaml:"mappingParams,omitempty"`
}

// MappingParam defines the mapping rules of headers and queryStrings
type MappingParam struct {
	Name  string `json:"name" yaml:"name"`
	MapTo string `json:"mapTo" yaml:"mapTo"`
}

// DubboBackendConfig defines the basic dubbo backend config
type DubboBackendConfig struct {
	ClusterName     string   `yaml:"clusterName" json:"clusterName"`
	ApplicationName string   `yaml:"applicationName" json:"applicationName"`
	Protocol        string   `yaml:"protocol" json:"protocol,omitempty" default:"dubbo"`
	Group           string   `yaml:"group" json:"group"`
	Version         string   `yaml:"version" json:"version"`
	Interface       string   `yaml:"interface" json:"interface"`
	Method          string   `yaml:"method" json:"method"`
	ParamTypes      []string `yaml:"paramTypes" json:"paramTypes"`
	Retries         string   `yaml:"retries" json:"retries,omitempty"`
}

// HTTPBackendConfig defines the basic dubbo backend config
type HTTPBackendConfig struct {
	//TODO  resolve from yaml config needed.   head mapping needed. paramMappingneeded
	TargetUrl       string `yaml:"targetUrl" json:"targetUrl"`
	ApplicationName string `yaml:"applicationName" json:"applicationName"`
	Protocol        string `yaml:"protocol" json:"protocol,schema" default:"https"`
	Version         string `yaml:"version" json:"version"`
	Retries         string `yaml:"retries" json:"retries,omitempty"`
}

// Definition defines the complex json request body
type Definition struct {
	Name   string `json:"name" yaml:"name"`
	Schema string `json:"schema" yaml:"schema"` // use json schema
}

// LoadAPIConfigFromFile load the api config from file
func LoadAPIConfigFromFile(path string) (*APIConfig, error) {
	if len(path) == 0 {
		return nil, perrors.Errorf("Config file not specified")
	}
	logger.Infof("Load API configuration file form %s", path)
	apiConf := &APIConfig{}
	err := yaml.UnmarshalYMLConfig(path, apiConf)
	if err != nil {
		return nil, perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	once.Do(func() {
		apiConfig = apiConf
	})
	return apiConf, nil
}

// GetAPIConf returns the initted api config
func GetAPIConf() APIConfig {
	return *apiConfig
}
