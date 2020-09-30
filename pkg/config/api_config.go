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
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/yaml"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
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

// APIConfig defines the data structure of the api gateway configuration
type APIConfig struct {
	Name        string       `yaml:"name"`
	Description string       `yaml:"description"`
	Resources   []Resource   `yaml:"resources"`
	Definitions []Definition `yaml:"definitions"`
}

// Resource defines the API path
type Resource struct {
	Type        string     `yaml:"type"` // Restful, Dubbo
	Path        string     `yaml:"path"`
	Description string     `yaml:"description"`
	Filters     []string   `yaml:"filters"`
	Methods     []Method   `yaml:"methods"`
	Resources   []Resource `yaml:"resources,omitempty"`
}

// Method defines the method of the api
type Method struct {
	OnAir              bool     `yaml:"onAir"` // true means the method is up and false means method is down
	Filters            []string `yaml:"filters"`
	HTTPVerb           `yaml:"httpVerb"`
	InboundRequest     `yaml:"inboundRequest"`
	IntegrationRequest `yaml:"integrationRequest"`
}

// InboundRequest defines the details of the inbound
type InboundRequest struct {
	RequestType  string           `yaml:"requestType"` //http, TO-DO: dubbo
	Headers      []Params         `yaml:"headers"`
	QueryStrings []Params         `yaml:"queryStrings"`
	RequestBody  []BodyDefinition `yaml:"requestBody"`
}

// Params defines the simple parameter definition
type Params struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required"`
}

// BodyDefinition connects the request body to the definitions
type BodyDefinition struct {
	DefinitionName string `yaml:"definitionName"`
}

// IntegrationRequest defines the backend request format and target
type IntegrationRequest struct {
	RequestType         string         `yaml:"requestType"` // dubbo, TO-DO: http
	MappingParams       []MappingParam `yaml:"mappingParams,omitempty"`
	dubbo.DubboMetadata `yaml:"dubboMetaData,inline,omitempty"`
}

// MappingParam defines the mapping rules of headers and queryStrings
type MappingParam struct {
	Name  string `yaml:"name"`
	MapTo string `yaml:"mapTo"`
}

// Definition defines the complex json request body
type Definition struct {
	Name   string `yaml:"name"`
	Schema string `yaml:"schema"` // use json schema
}

// LoadAPIConfigFromFile load the api config from file
func LoadAPIConfigFromFile(path string) (*APIConfig, error) {
	if len(path) == 0 {
		return nil, perrors.Errorf("Config file not specified")
	}
	logger.Info("Load API configuration file form %s", path)
	apiConf := &APIConfig{}
	err := yaml.UnmarshalYMLConfig(path, apiConf)
	if err != nil {
		return nil, perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	return apiConf, nil
}
